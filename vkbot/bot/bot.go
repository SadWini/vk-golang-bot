package bot

import (
	"vkbot/api"
	"vkbot/structs"
)

type VKBot struct {
	msgRoutes        map[string]msgRoute
	actionRoutes     map[string]func(*structs.Message) string
	cmdHandlers      map[string]func(*structs.Message) string
	msgHandlers      map[string]func(*structs.Message) string
	errorHandler     func(*structs.Message, error)
	LastMsg          int
	lastUserMessages map[int]int
	lastChatMessages map[int]int
	autoFriend       bool
	IgnoreBots       bool
	LongPoll         LongPollServer
	API              *api.VkAPI
}

type msgRoute struct {
	SimpleHandler func(*structs.Message) string
	Handler       func(*structs.Message) Reply
}

func (api *api.VkAPI) NewBot() *VKBot {
	if api.IsGroup() {
		return &VKBot{
			msgRoutes:        make(map[string]msgRoute),
			actionRoutes:     make(map[string]func(*api.Message) string),
			lastUserMessages: make(map[int]int),
			lastChatMessages: make(map[int]int),
			LongPoll:         NewGroupLongPollServer(API.RequestInterval),
			API:              api,
		}
	}
	return &VKBot{
		msgRoutes:        make(map[string]msgRoute),
		actionRoutes:     make(map[string]func(*api.Message) string),
		lastUserMessages: make(map[int]int),
		lastChatMessages: make(map[int]int),
		LongPoll:         NewUserLongPollServer(false, longPollVersion, API.RequestInterval),
		API:              api,
	}
}
