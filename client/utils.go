package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/ahmetalpbalkan/go-cursor"
	"github.com/fatih/color"
	"golang.org/x/term"
)

type ip struct {
	Ip string `json:"ip"`
}

func getTermSize() (x int, y int) {
   x,y,err := term.GetSize(int(os.Stdout.Fd()))
   if err != nil {
	  log.Fatal(err)
	}
   return x,y
}

func printIp(ip string) {
	x, y := getTermSize()
		_ = y
		fmt.Print(cursor.MoveTo(1, (x/2 - (len(ip) / 2))))
		color.Green("Your IP is: %s", ip)
}

func getIp(write chan string) string {
	resp, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var ip ip
	json.Unmarshal(data, &ip)
	write <- ip.Ip
	return ip.Ip
}