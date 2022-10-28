package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ahmetalpbalkan/go-cursor"
	"github.com/c-bata/go-prompt"
	"github.com/fatih/color"
	"github.com/gen2brain/beeep"
	"github.com/gorilla/websocket"
)

type baseType struct {
	Type string `json:"type"`
}

type wordResp struct {
	Type string `json:"type"`
	Word string `json:"word"`
}

type leaked struct {
	Type string `json:"type"`
	Ip string `json:"ip"`
}

type wordGuess struct {
	Type string
	Guess string
}

var IpRes chan string = make(chan string)

func connectToServer(urlToConnect string) {
	u := url.URL{Scheme: "ws", Host: urlToConnect}
	fmt.Print(cursor.Hide())
	fmt.Print(cursor.ClearEntireScreen())
	fmt.Print(cursor.MoveTo(0, 0))
	color.Red("Connecting you to our servers...")

	// VALUE INITIALIZERS
	go getIp(IpRes)
	ip := <-IpRes
	c, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{"ip": {ip}})
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	ticker := time.NewTicker(time.Second)

	defer c.Close()

	if err != nil {
		color.Red("Error connecting to servers... Trying again...")
	} else {
		fmt.Print(cursor.ClearEntireScreen())
		fmt.Print(cursor.MoveTo(0, 0))
		color.Blue("You are connected to our servers!")
		time.Sleep(2 * time.Second)
		fmt.Print(cursor.ClearEntireScreen())
		color.Green("Your IP is: %s", ip)
		printIp(ip)

		defer ticker.Stop()
		defer func() {
			fmt.Print(cursor.ClearEntireScreen())
			fmt.Print(cursor.MoveTo(0, 0))
			color.Red("You have been disconnected from our servers. The game has been ended and your IP has not been leaked. (Dont worry, we care)")
		}()

		for {
			select {
			case <-interrupt:
				fmt.Println("Why do you want to leave...")

				err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					fmt.Println("write close:", err)
					return
				}
				time.Sleep(1)
				return
			case <-ticker.C:
				if err != nil {
					return
				}
				_, data, err := c.ReadMessage()
				if err != nil {
					return
				}
				var unknown baseType
				json.Unmarshal(data, &unknown)
				if unknown.Type == "word" {
					var word wordResp
					json.Unmarshal(data, &word)
					printDefaults(ip, word.Word)
					go func() {
					promptForWord(c)
					}()
					printDefaults(ip, word.Word)
					} else if unknown.Type == "heartbeat" {

				} else if unknown.Type == "incorrect" {
					go func() {
                        fmt.Print(cursor.ClearEntireScreen())
					fmt.Print(cursor.MoveTo(0, 0))
					color.Blue("Your word was incorrect...")
					time.Sleep(4 * time.Second)
					fmt.Print(cursor.ClearEntireScreen())
					fmt.Print(cursor.MoveTo(0, 0))
					color.Green("Your IP is: %s", ip)
					time.Sleep(4 * time.Second)
					fmt.Print(cursor.ClearEntireScreen())
					fmt.Print(cursor.MoveTo(0, 0))
					color.Blue("You have been disconnected from our servers... Dont worry, your ip has been leaked.")
					time.Sleep(4 * time.Second)
					os.Exit(1)
					}()
				} else if unknown.Type == "correct" {
					go func() {
                        fmt.Print(cursor.ClearEntireScreen())
					fmt.Print(cursor.MoveTo(0, 0))
					color.Blue("You got it correct...")
					time.Sleep(4 * time.Second)
					fmt.Print(cursor.ClearEntireScreen())
					fmt.Print(cursor.MoveTo(0, 0))
					color.Green("Your IP is: %s", ip)
					time.Sleep(4 * time.Second)
					fmt.Print(cursor.ClearEntireScreen())
					fmt.Print(cursor.MoveTo(0, 0))
					color.Blue("You have been disconnected from our servers... Dont worry, only you know your ip...")
					time.Sleep(4 * time.Second)
					os.Exit(1)
					}()
				} else if unknown.Type == "leaked" {
					var leaked leaked
					json.Unmarshal(data, &leaked)
					go func() {
						beeep.Notify("An ip has been leaked!", leaked.Ip, "")
					}()
				} else {
					fmt.Print(cursor.ClearEntireScreen())
					return
				}
			}
		}
	}
}

func printDefaults(ip string, word string) {
	fmt.Print(cursor.ClearEntireScreen())
	fmt.Print(cursor.MoveTo(0, 0))
	color.Green("Your word is: %s", word)
	fmt.Print(cursor.Show())
	printIp(ip)
	fmt.Print(cursor.MoveTo(5, 1))
}

func promptForWord(socket *websocket.Conn) string {
	fmt.Print(cursor.Show())
	text := prompt.Input("What word do you think you got >", completer)
	fmt.Println(len(text))
	if len(text) > 5 {
		fmt.Print(cursor.ClearEntireScreen())
		fmt.Print(cursor.Hide())
		color.Red("Your word is too long! Please try again.")
		time.Sleep(2 * time.Second)
		print(cursor.MoveUp(1))
		print(cursor.ClearEntireLine())

		return promptForWord(socket)
	} else if len(text) < 5 {
		fmt.Print(cursor.ClearEntireScreen())
		fmt.Print(cursor.Hide())
		color.Red("Your word is too short! Please try again.")
		time.Sleep(2 * time.Second)
		print(cursor.MoveUp(1))
		print(cursor.ClearEntireLine())

		return promptForWord(socket)
	}
	message, err := json.Marshal(wordGuess{Guess: text, Type: "guess"})
	if err != nil {
		color.Red("Error sending word to server... Im sorry.")
		os.Exit(1)
	}
	err2 := socket.WriteMessage(websocket.TextMessage, message)
	if err2 != nil {
		color.Red("Error sending word to server... Im sorry.")
		os.Exit(1)
	}
	return text
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}
