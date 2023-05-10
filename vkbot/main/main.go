package main

import (
	"vkbot/bot"
	"vkbot/utils"
)

var config configuration

func main() {

	utils.ReadJSON("config.json", &config)

	bot.HandleAdvancedMessage("начать", helpHandler)
	bot.HandleAdvancedMessage("/help", helpHandler)

	bot.HandleAdvancedMessage("/info", infoHandler)
	bot.HandleAdvancedMessage("/botstat", infoHandler)
	bot.HandleAdvancedMessage("/userstat", infoHandler)

	bot.HandleAdvancedMessage("/admin", adminHandler)
	bot.HandleAdvancedMessage("/semen", adminHandler)
	bot.HandleAdvancedMessage("/skylorun", adminHandler)

	bot.HandleAdvancedMessage("/quotes", quoteHandler)
	bot.HandleAdvancedMessage("/funny", quoteHandler)
	bot.HandleAdvancedMessage("/serious", quoteHandler)

	bot.HandleAction("friend_add", addFriendHandler)
	bot.HandleAction("friend_delete", deleteFriendHandler)

	bot.HandleError(errorHandler)
	bot.SetAutoFriend(true) // enable auto accept/delete friends

	bot.SetDebug(true)                                 // log debug messages
	bot.Listen(config.VKToken, "", "", config.AdminID) // start bot

}
