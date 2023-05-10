package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"vkbot/bot"
	"vkbot/structs"
)

// admins' id
const (
	semenID    = 301126247
	skylorunID = 477249238
)

type configuration struct {
	VKToken      string
	ServiceURL   string
	ServiceToken string
	AdminID      int
}

func getMeMessage(uid int) (reply string) {
	me, _ := bot.API.Me()
	return fmt.Sprintf("You: %+v %+v", me.FirstName, me.LastName)
}

func anyHandler(m *structs.Message) (reply string) {
	notifyAdmin(fmt.Sprintf("Command %+v by user vk.com/id%+v in chat %+v", m.Body, m.UserID, m.PeerID))
	return reply
}

func meHandler(m *structs.Message) (reply string) {
	return getMeMessage(m.UserID)
}

func infoHandler(m *structs.Message) (reply structs.Reply) {
	keyboard := keyboardByCommand(infoCommand)
	reply.Keyboard = &keyboard
	if strings.Contains(m.Body, "/botstat") {
		reply.Msg = "Bot already sent " + bot.Bot.LongPoll.GetTs() + " messages"
		return reply
	} else if strings.Contains(m.Body, "/userstat") {
		reply.Msg = "You already sent " + strconv.Itoa(m.ConvMsgID) + "to this bot"
		return reply
	} else {
		reply.Msg = infoText
		return reply
	}
}
func adminHandler(m *structs.Message) (reply structs.Reply) {
	if strings.Contains(m.Body, "/semen") {
		return adminContacts(semenID)
	} else if strings.Contains(m.Body, "/skylorun") {
		return adminContacts(skylorunID)
	} else {
		keyboard := keyboardByCommand(adminCommand)
		return structs.Reply{Msg: adminText, Keyboard: &keyboard}
	}
}

func helpHandler(m *structs.Message) (reply structs.Reply) {
	reply.Msg = "Information about buttons:\n" + "info button: " + infoText +
		"\n" + "admin button: " + adminText + "\n" + "quote button: " + quoteText
	keyboard := keyboardByCommand(Commands)
	reply.Keyboard = &keyboard
	return reply
}

func quoteHandler(m *structs.Message) (reply structs.Reply) {
	keyboard := keyboardByCommand(quoteCommand)
	reply.Keyboard = &keyboard
	if strings.Contains(m.Body, "/funny") {
		reply.Msg = funnyQuotes[rand.Int()%(len(funnyQuotes)-1)]
		return reply
	} else if strings.Contains(m.Body, "/serious") {
		reply.Msg = seriousQuotes[rand.Int()%(len(seriousQuotes)-1)]
		return reply
	} else {
		reply.Msg = quoteText
		return reply
	}
}

func errorHandler(msg *structs.Message, err error) {
	if _, ok := err.(*structs.VKError); !ok {
		notifyAdmin("VK ERROR: " + err.Error())
	}
	notifyAdmin("ERROR: " + err.Error())
}

func greetUser(uid int) (reply string) {
	u, err := bot.API.User(uid)
	if err == nil {
		reply = fmt.Sprintf("Hello, %+v", u.FullName())
	}
	return reply
}
func replyGreet() (reply string) {
	reply = "Hi all. I'am bot\n" + availableCommands
	return reply
}

func addFriendHandler(m *structs.Message) (reply string) {
	log.Printf("friend %+v added\n", m.UserID)
	notifyAdmin(fmt.Sprintf("user vk.com/id%+v add me to friends", m.UserID))
	return reply
}
func deleteFriendHandler(m *structs.Message) (reply string) {
	log.Printf("friend %+v deleted\n", m.UserID)
	notifyAdmin(fmt.Sprintf("user vk.com/id%+v delete me from friends", m.UserID))
	return reply
}

func adminContacts(adminID int) (reply structs.Reply) {
	user, err := bot.API.User(adminID)
	if err != nil {
		fmt.Sprintf("Cannot get info about admin")
	}
	reply.Msg = user.FullName() + " and my vkID is " + strconv.Itoa(semenID)
	keyboard := keyboardByCommand(Commands)
	reply.Keyboard = &keyboard
	return reply
}

func keyboardByCommand(data []string) structs.Keyboard {
	keyboard := structs.Keyboard{Buttons: make([][]structs.Button, 0)}
	row := make([]structs.Button, 0)
	for _, com := range data {
		button := bot.NewButton(com, nil)
		row = append(row, button)
	}
	keyboard.Buttons = append(keyboard.Buttons, row)
	return keyboard
}
func notifyAdmin(msg string) {
	err := NotifyAdmin(msg)
	if err != nil {
		log.Printf("VK Admin Notify ERROR: %+v\n", msg)
	}
}

func NotifyAdmin(msg string) error {
	return bot.API.NotifyAdmin(msg)
}
