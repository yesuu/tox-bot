package main

import (
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/codedust/go-tox"
)

const (
	savepath  = "./bot.tox"
	address   = "192.254.75.102"
	port      = 33445
	publicKey = "951C88B7E75C867418ACDB5D273821372BB5BD652740BCDF623A4FA293E75D2F"
)

func main() {
	var bot *gotox.Tox

	savedata, err := ioutil.ReadFile(savepath)
	if err != nil {
		log.Println("Warn:", err)
		bot, err = newBot()
		if err != nil {
			log.Panicln("Error:", err)
		}
	} else {
		bot, err = gotox.New(&gotox.Options{
			IPv6Enabled:  true,
			UDPEnabled:   true,
			ProxyType:    gotox.TOX_PROXY_TYPE_NONE,
			StartPort:    0,
			EndPort:      0,
			TcpPort:      0,
			SaveDataType: gotox.TOX_SAVEDATA_TYPE_TOX_SAVE,
			SaveData:     savedata,
		})
		if err != nil {
			log.Panicln("Error:", err)
		}
	}

	var botId []byte
	botId, err = bot.SelfGetAddress()
	if err != nil {
		log.Panicln("Error:", err)
	}
	log.Println("Info: Tox Id", hex.EncodeToString(botId))

	bot.CallbackFriendRequest(onFriendRequest)
	bot.CallbackFriendMessage(onFriendMessage)

	publicKeyHex, _ := hex.DecodeString(publicKey)
	err = bot.Bootstrap(address, port, publicKeyHex)
	if err != nil {
		log.Panicln("Error:", err)
	}

	log.Println("Info: Run")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ticker := time.NewTicker(30 * time.Millisecond)

	for {
		select {
		case <-c:
			log.Println("Info: SaveData")
			var data []byte
			data, err = bot.GetSavedata()
			if err != nil {
				log.Panicln("Error:", err)
			}
			err = ioutil.WriteFile(savepath, data, 0600)
			if err != nil {
				log.Panicln("Error:", err)
			}

			log.Println("Info: Killing")
			bot.Kill()
			return
		case <-ticker.C:
			bot.Iterate()
		}
	}
}

func newBot() (*gotox.Tox, error) {
	bot, err := gotox.New(nil)
	if err != nil {
		return nil, err
	}

	bot.SelfSetName("emailBot")
	bot.SelfSetStatusMessage("😁")
	bot.SelfSetStatus(gotox.TOX_USERSTATUS_NONE)

	log.Println("Info: SaveData")
	var data []byte
	data, err = bot.GetSavedata()
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(savepath, data, 0600)
	if err != nil {
		return nil, err
	}
	return bot, nil

}

func onFriendRequest(t *gotox.Tox, publicKey []byte, message string) {
	log.Println("Info: 好友请求", hex.EncodeToString(publicKey))
	log.Println("Info: 好友请求消息", message)

	t.FriendAddNorequest(publicKey)
}

func onFriendMessage(t *gotox.Tox, friendNumber uint32, messageType gotox.ToxMessageType, message string) {
	if messageType == gotox.TOX_MESSAGE_TYPE_ACTION {
		log.Println("Info: 新消息（action）", friendNumber, message)
		return
	}

	log.Println("Info: 新消息", friendNumber, message)

	switch message {
	case "ping":
		t.FriendSendMessage(friendNumber, gotox.TOX_MESSAGE_TYPE_NORMAL, "pong")
	}
}
