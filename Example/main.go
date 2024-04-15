package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Distortions81/EbitUI"
)

func main() {

	//Detect sig int/term
	signalHandle = make(chan os.Signal, 1)
	signal.Notify(signalHandle, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	//Start EbitUI
	EbitUI.Start(defaultWindowWidth, defaultWindowHeight)

	//Create window data
	windowName := "Test window!"
	nw := EbitUI.DefaultWinSettings
	nw.Title = windowName
	nw.Size = EbitUI.V2i{X: 400, Y: 400}
	nw.Position = EbitUI.V2i{X: 25, Y: 25}

	//Add the window
	err := EbitUI.AddWindow(windowName, nw)
	if err != nil {
		log.Fatal(err)
	}

	//Open it
	err = EbitUI.OpenWindow(windowName)
	if err != nil {
		log.Fatal(err)
	}

	//Start Ebiten on a goroutine so we can detect control-c or window close
	go startEbiten()

	//Wait here for signal
	<-signalHandle
}
