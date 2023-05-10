package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"vkbot/structs"
	"vkbot/utils"
)

type VKBot struct {
	msgRoutes        map[string]msgRoute
	actionRoutes     map[string]func(*structs.Message) string
	cmdHandlers      map[string]func(*structs.Message) string
	msgHandlers      map[string]func(*structs.Message) string
	ErrorHandler     func(*structs.Message, error)
	LastMsg          int
	lastUserMessages map[int]int
	lastChatMessages map[int]int
	autoFriend       bool
	IgnoreBots       bool
	LongPoll         LongPollServer
	API              *VkAPI
}

type msgRoute struct {
	SimpleHandler func(*structs.Message) string
	Handler       func(*structs.Message) structs.Reply
}

func (api *VkAPI) NewBot() *VKBot {
	if api.IsGroup() {
		return &VKBot{
			msgRoutes:        make(map[string]msgRoute),
			actionRoutes:     make(map[string]func(*structs.Message) string),
			lastUserMessages: make(map[int]int),
			lastChatMessages: make(map[int]int),
			LongPoll:         NewGroupLongPollServer(API.RequestInterval),
			API:              api,
		}
	}
	return &VKBot{
		msgRoutes:        make(map[string]msgRoute),
		actionRoutes:     make(map[string]func(*structs.Message) string),
		lastUserMessages: make(map[int]int),
		lastChatMessages: make(map[int]int),
		LongPoll:         NewUserLongPollServer(false, LongPollVersion, API.RequestInterval),
		API:              api,
	}
}

func (bot *VKBot) SetAutoFriend(af bool) {
	bot.autoFriend = af
}

// HandleMessage - add substr message handler.
// Function must return string to reply or "" (if no reply)
func (bot *VKBot) HandleMessage(command string, handler func(*structs.Message) string) {
	bot.msgRoutes[command] = msgRoute{SimpleHandler: handler}
}

func (bot *VKBot) HandleAdvancedMessage(command string, handler func(*structs.Message) structs.Reply) {
	bot.msgRoutes[command] = msgRoute{Handler: handler}
}

// HandleAction - add action handler.
// Function must return string to reply or "" (if no reply)
func (bot *VKBot) HandleAction(command string, handler func(*structs.Message) string) {
	bot.actionRoutes[command] = handler
}

// HandleError - add error handler
func (bot *VKBot) HandleError(handler func(*structs.Message, error)) {
	bot.ErrorHandler = handler
}
func (bot *VKBot) GetMessages() ([]*structs.Message, error) {
	var allMessages []*structs.Message
	lastMsg := bot.LastMsg
	offset := 0
	var err error
	var messages *structs.Messages
	for {
		messages, err = bot.API.GetMessages(bot.API.MessagesCount, offset)
		if len(messages.Items) > 0 {
			if messages.Items[0].ID > lastMsg {
				lastMsg = messages.Items[0].ID
			}
		}
		allMessages = append(allMessages, messages.Items...)
		if bot.LastMsg > 0 {
			if len(messages.Items) > 0 {
				if messages.Items[len(messages.Items)-1].ID <= bot.LastMsg {
					bot.LastMsg = lastMsg
					break
				}
			} else {
				break
			}
			offset += bot.API.MessagesCount
		} else {
			bot.LastMsg = lastMsg
			break
		}
	}
	if offset > 0 {
		bot.API.NotifyAdmin("many messages in interval. offset: " + strconv.Itoa(offset))
	}
	return allMessages, err
}

// RouteAction routes an action
func (bot *VKBot) RouteAction(m *structs.Message) (replies []string, err error) {
	if m.Action != "" {
		DebugPrint("route action: %+v\n", m.Action)
		for k, v := range bot.actionRoutes {
			if m.Action == k {
				msg := v(m)
				if msg != "" {
					replies = append(replies, msg)
				}
			}
		}
	}
	return replies, nil
}

// ListenGroup - listen group VK API
func (bot *VKBot) ListenGroup(api *VkAPI) error {
	bot.LongPoll = NewGroupLongPollServer(bot.API.RequestInterval)
	c := time.Tick(3 * time.Second)
	for range c {
		bot.MainRoute()
	}
	return nil
}

// ListenUser - listen User VK API (deprecated)
func (bot *VKBot) ListenUser(api *VkAPI) error {
	bot.LongPoll = NewUserLongPollServer(false, LongPollVersion, API.RequestInterval)
	go bot.friendReceiver()

	c := time.Tick(3 * time.Second)
	for range c {
		bot.MainRoute()
	}
	return nil
}

func (bot *VKBot) friendReceiver() {
	if bot.API.UID > 0 {
		bot.CheckFriends()
		c := time.Tick(30 * time.Second)
		for range c {
			bot.CheckFriends()
		}
	}
}

// MainRoute - main router func. Working cycle Listen.
func (bot *VKBot) MainRoute() {
	messages, err := bot.LongPoll.GetLongPollMessages()
	if err != nil {
		SendError(nil, err)
	}
	replies := bot.RouteMessages(messages)
	for m, msgs := range replies {
		for _, reply := range msgs {
			fmt.Println("outbox: ", reply.Msg)
			if reply.Msg != "" || reply.Keyboard != nil {
				_, err = bot.Reply(m, reply)
				if err != nil {
					log.Printf("Error sending message: '%+v'\n", reply)
					SendError(m, err)
					_, err = bot.Reply(m, structs.Reply{Msg: "Cant send message, maybe wrong/china letters?"})
					if err != nil {
						SendError(m, err)
					}
				}
			}
		}
	}
}

// RouteMessages routes inbound messages
func (bot *VKBot) RouteMessages(messages []*structs.Message) (result map[*structs.Message][]structs.Reply) {
	result = make(map[*structs.Message][]structs.Reply)
	for _, m := range messages {
		if m.ReadState == 0 {
			if bot.IgnoreBots && m.UserID < 0 {
				continue
			}
			replies, err := bot.RouteMessage(m)
			if err != nil {
				SendError(m, err)
			}
			if len(replies) > 0 {
				result[m] = replies
			}
		}
	}
	return result
}

// RouteMessage routes single message
func (bot *VKBot) RouteMessage(m *structs.Message) (replies []structs.Reply, err error) {
	message := strings.TrimSpace(strings.ToLower(m.Body))
	if utils.HasPrefix(message, "/ ") {
		message = "/" + utils.TrimPrefix(message, "/ ")
	}
	if m.Action != "" {
		actionReplies, err := bot.RouteAction(m)
		for _, r := range actionReplies {
			replies = append(replies, structs.Reply{Msg: r})
		}
		return replies, err
	}
	for k, v := range bot.msgRoutes {
		if utils.HasPrefix(message, k) {
			if v.Handler != nil {
				reply := v.Handler(m)
				if reply.Msg != "" || reply.Keyboard != nil {
					replies = append(replies, reply)
				}
			} else {
				msg := v.SimpleHandler(m)
				if msg != "" {
					replies = append(replies, structs.Reply{Msg: msg})
				}

			}
		}
	}
	return replies, nil
}

// Reply - reply message
func (bot *VKBot) Reply(m *structs.Message, reply structs.Reply) (id int, err error) {
	if m.PeerID != 0 {
		return bot.API.SendAdvancedPeerMessage(m.PeerID, reply)
	}
	if m.ChatID != 0 {
		return bot.API.SendChatMessage(m.ChatID, reply.Msg)
	}
	return bot.API.SendMessage(m.UserID, reply.Msg)
}
