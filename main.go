package main

import (
	"fmt"
	"multimessenger_bot/internal/bot"
	"multimessenger_bot/internal/config"
	ma "multimessenger_bot/internal/messenger_adapter"
	"multimessenger_bot/internal/telegram"
	"multimessenger_bot/internal/whatsapp"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func main() {

	//gCalendarClient, err := google_calendar.NewGoogleCalendarClient("service_credentials.json")
	//if err != nil {
	//	  fmt.Println(err)
	//	  return
	//}

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	cfg, err := config.Load("config.toml")
	if err != nil {
		fmt.Print(err)
		return
	}

	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New("sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}

	//clientLog := waLog.Stdout("Client", "DEBUG", true)

	recvMsgChan := make(chan *ma.Message)
	tgClient, _ := telegram.NewTelegramClient(cfg, recvMsgChan)
	waClient, _ := whatsapp.NewWhatsAppClient(nil, cfg, container, recvMsgChan)

	bot, _ := bot.NewBot([]ma.ClientInterface{waClient, tgClient}, recvMsgChan)
	bot.Run()

	signalHandler := setupSignalHandler()
	<-signalHandler

	bot.Shutdown()
}

func setupSignalHandler() chan os.Signal {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	return ch
}
