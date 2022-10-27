package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/ahmetalpbalkan/go-cursor"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
)

type baseType struct {
	Type string `json:"type"`
}

type wordResp struct {
	Type string `json:"type"`
	Word string `json:"word"`
}

func connectToServer(urlToConnect string) {
	u := url.URL{Scheme: "ws", Host: urlToConnect}
	fmt.Print(cursor.ClearEntireScreen())
	fmt.Print(cursor.MoveTo(0, 0))
	color.Red("Connecting you to our servers...")

	// VALUE INITIALIZERS
	IpRes := make(chan string)
	go getIp(IpRes)
	ip := <-IpRes
	c, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{"ip": {ip}})
	interrupt := make(chan os.Signal, 1)
	ticker := time.NewTicker(time.Second)

	if err != nil {
		color.Red("Error connecting to servers... Trying again...")
	} else {
		fmt.Print(cursor.ClearEntireScreen())
		fmt.Print(cursor.MoveTo(0, 0))
		color.Blue("You are connected to our servers!")
		time.Sleep(2 * time.Second)
		fmt.Print(cursor.ClearEntireScreen())
		x, y := getTermSize()
		_ = y
		fmt.Print(cursor.MoveTo(1, (x/2 - (len(ip) / 2))))
		color.Green("Your IP is: %s", ip)
	}
	signal.Notify(interrupt, os.Interrupt)

	defer ticker.Stop()
	defer func() {
		fmt.Print(cursor.ClearEntireScreen())
		fmt.Print(cursor.MoveTo(0, 0))
		color.Red("You have been disconnected from our servers. The game has been ended and your IP has not been leaked. (Dont worry, we care)")
	}()

	for {
		select {
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
			switch unknown.Type {
			case "word":
				var word wordResp
				json.Unmarshal(data, &word)
				fmt.Print(cursor.ClearEntireScreen())
				fmt.Print(cursor.MoveTo(0, 0))
				color.Green("Your word is: %s", word.Word)
				printIp(ip)
			default:
				fmt.Print(cursor.ClearEntireScreen())
			}
		case <-interrupt:
			fmt.Println("Why do you want to leave...")

			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Println("write close:", err)
				return
			}
			time.Sleep(1)
			return
		}
	}
}
