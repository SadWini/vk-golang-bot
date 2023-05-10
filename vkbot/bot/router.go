package bot

import (
	"log"
	"vkbot/structs"
)

const (
	vkAPIURL        = "https://api.vk.com/method/"
	vkAPIVer        = "5.131"
	messagesCount   = 200
	requestInterval = 400 // 3 requests per second VK limit
)

var API = newAPI()

var Bot = API.NewBot()

func newAPI() *VkAPI {
	return &VkAPI{
		Token:           "",
		URL:             vkAPIURL,
		Ver:             vkAPIVer,
		MessagesCount:   messagesCount,
		RequestInterval: requestInterval,
		DEBUG:           false,
		HTTPS:           true,
	}
}
func SetToken(token string) {
	API.Token = token
}

func SetAPI(token string, url string, ver string) {
	SetToken(token)
	if url != "" {
		API.URL = url
	}
	if ver != "" {
		API.Ver = ver
	}
}

func SetDebug(debug bool) {
	API.DEBUG = debug
}

func SetAutoFriend(af bool) {
	Bot.SetAutoFriend(af)
}

// SetLang - sets VK response language. Default auto. Available: en, ru, ua, be, es, fi, de, it
func SetLang(lang string) {
	API.Lang = lang
}

// HandleMessage - add substr message handler.
// Function must return string to reply or "" (if no reply)
func HandleMessage(command string, handler func(string2 *structs.Message) string) {
	Bot.HandleMessage(command, handler)
}

// HandleAdvancedMessage - add substr message handler.
// Function must return string to reply or "" (if no reply)
func HandleAdvancedMessage(command string, handler func(*structs.Message) structs.Reply) {
	Bot.HandleAdvancedMessage(command, handler)
}

// HandleAction - add action handler.
// Function must return string to reply or "" (if no reply)
func HandleAction(command string, handler func(*structs.Message) string) {
	Bot.HandleAction(command, handler)
}

// HandleError - add error handler
func HandleError(handler func(*structs.Message, error)) {
	Bot.HandleError(handler)
}

func SendError(msg *structs.Message, err error) {
	if Bot.ErrorHandler != nil {
		Bot.ErrorHandler(msg, err)
	} else {
		log.Fatalf("VKBot error: %+v\n", err.Error())
	}

}

// Listen - start server
func Listen(token string, url string, ver string, adminID int) error {
	if API.Token == "" {
		SetAPI(token, url, ver)
	}
	API.AdminID = adminID
	if Bot.API.IsGroup() {
		return Bot.ListenGroup(API)
	}
	return Bot.ListenUser(API)
}

// NotifyAdmin - notify AdminID by VK
func NotifyAdmin(msg string) error {
	return API.NotifyAdmin(msg)
}
