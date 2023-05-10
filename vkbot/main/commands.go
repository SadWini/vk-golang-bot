package main

const (
	availableCommands = "Available commands: /help, /info, /admin "
	infoText          = "You can click botstat to see statistics of bot or userstat to get your stats"
	adminText         = "You can dm admin for any collaborations"
	quoteText         = "You can get quotes funny or serious by your choice"
)

var Commands = []string{"/help", "/info", "/admin", "/quotes"}
var infoCommand = []string{"/botstat", "/userstat", "/help"}
var adminCommand = []string{"/skylorun", "/semen", "/help"}
var quoteCommand = []string{"/funny", "/serious", "/help"}
