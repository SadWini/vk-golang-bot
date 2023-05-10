package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"vkbot/structs"
)

type H map[string]string

type VkAPI struct {
	Token           string
	URL             string
	Ver             string
	UID             int
	GroupID         int
	Lang            string
	HTTPS           bool
	AdminID         int
	MessagesCount   int
	RequestInterval int
	DEBUG           bool
}

const (
	apiUsersGet                       = "users.get"
	apiGroupsGet                      = "groups.getById"
	apiMessagesGet                    = "messages.get"
	apiMessagesGetChat                = "messages.getChat"
	apiMessagesGetConversationsById   = "messages.getConversationsById"
	apiMessagesGetChatUsers           = "messages.getChatUsers"
	apiMessagesGetConversationMembers = "messages.getConversationMembers"
	apiMessagesSend                   = "messages.send"
	apiMessagesMarkAsRead             = "messages.markAsRead"
	apiFriendsGetRequests             = "friends.getRequests"
	apiFriendsAdd                     = "friends.add"
	apiFriendsDelete                  = "friends.delete"
)

func (api *VkAPI) IsGroup() bool {
	if api.GroupID != 0 {
		return true
	} else if api.UID != 0 {
		return false
	}

	g, err := API.CurrentGroup()
	if err != nil || g.ID == 0 {
		u, err := API.Me()
		if err != nil || u == nil {
			fmt.Printf("Get current user/group error %+v\n", err)
		} else {
			api.UID = u.ID
		}
	} else {
		api.GroupID = g.ID
	}
	return api.GroupID != 0
}

func (api *VkAPI) Call(method string, params map[string]string) ([]byte, error) {
	DebugPrint("vk req: %+v params: %+v\n", api.URL+method, params)
	params["access_token"] = api.Token
	params["v"] = api.Ver
	if api.Lang != "" {
		params["lang"] = api.Lang
	}
	if api.HTTPS {
		params["https"] = "1"
	}

	parameters := url.Values{}
	for k, v := range params {
		parameters.Add(k, v)
	}

	if api.URL == "test" {
		content, err := ioutil.ReadFile("./mocks/" + method + ".json")
		return content, err
	}
	resp, err := http.PostForm(api.URL+method, parameters)
	if err != nil {
		DebugPrint("%+v\n", err.Error())
		time.Sleep(time.Duration(time.Millisecond * time.Duration(api.RequestInterval)))
		return nil, err
	}
	buf, err := ioutil.ReadAll(resp.Body)
	time.Sleep(time.Duration(time.Millisecond * time.Duration(api.RequestInterval)))
	DebugPrint("vk resp: %+v\n", string(buf))

	return buf, err
}
func (api *VkAPI) Me() (*structs.User, error) {

	r := structs.UsersResponse{}
	err := api.CallMethod(apiUsersGet, H{"fields": "screen_name"}, &r)

	if len(r.Response) > 0 {
		DebugPrint("me: %+v - %+v\n", r.Response[0].ID, r.Response[0].ScreenName)
		return r.Response[0], err
	}
	return nil, err
}

func (api *VkAPI) CallMethod(method string, params map[string]string, result interface{}) error {
	buf, err := api.Call(method, params)
	if err != nil {
		return err
	}
	r := structs.ErrorResponse{}
	err = json.Unmarshal(buf, &r)
	if err != nil {
		return &structs.ResponseError{
			Err:     errors.New("vkapi: vk response is not json"),
			Content: string(buf)}
	}
	if r.Error != nil {
		DebugPrint("%+v\n", r.Error.ErrorMsg)
		return r.Error
	}

	err = json.Unmarshal(buf, result)
	return err
}

func (api *VkAPI) CurrentGroup() (*structs.User, error) {

	r := structs.UsersResponse{}
	err := api.CallMethod(apiGroupsGet, H{"fields": "screen_name"}, &r)

	if len(r.Response) > 0 {
		DebugPrint("me: %+v - %+v\n", r.Response[0].ID, r.Response[0].ScreenName)
		return r.Response[0], err
	}
	return nil, err
}

func (api *VkAPI) GetRandomID() string {
	return strconv.FormatUint(uint64(rand.Uint32()), 10)
}

// SendChatMessage sending a message to chat
func (api *VkAPI) SendChatMessage(chatID int, msg string) (id int, err error) {
	r := structs.SimpleResponse{}
	err = api.CallMethod(apiMessagesSend, H{
		"chat_id":          strconv.Itoa(chatID),
		"message":          msg,
		"dont_parse_links": "1",
		"random_id":        api.GetRandomID(),
	}, &r)
	return r.Response, err
}

// SendAdvancedPeerMessage sending a message to chat
func (api *VkAPI) SendAdvancedPeerMessage(peerID int64, message structs.Reply) (id int, err error) {
	r := structs.SimpleResponse{}
	params := H{
		"peer_id":          strconv.FormatInt(peerID, 10),
		"message":          message.Msg,
		"dont_parse_links": "1",
		"random_id":        api.GetRandomID(),
	}
	if message.Keyboard != nil {
		keyboard, err := json.Marshal(message.Keyboard)
		if err != nil {
			fmt.Printf("ERROR encode keyboard %+v\n", message.Keyboard)
		} else {
			params["keyboard"] = string(keyboard)
		}
	}
	err = api.CallMethod(apiMessagesSend, params, &r)
	return r.Response, err
}

// SendPeerMessage sending a message to chat
func (api *VkAPI) SendPeerMessage(peerID int64, msg string) (id int, err error) {
	r := structs.SimpleResponse{}
	err = api.CallMethod(apiMessagesSend, H{
		"peer_id":          strconv.FormatInt(peerID, 10),
		"message":          msg,
		"dont_parse_links": "1",
		"random_id":        api.GetRandomID(),
	}, &r)
	return r.Response, err
}

func (api *VkAPI) SendMessage(userID int, msg string) (id int, err error) {
	r := structs.SimpleResponse{}
	if msg != "" {
		err = api.CallMethod(apiMessagesSend, H{
			"user_id":          strconv.Itoa(userID),
			"message":          msg,
			"dont_parse_links": "1",
			"random_id":        api.GetRandomID(),
		}, &r)
	}
	return r.Response, err
}
func (api *VkAPI) GetMessages(count int, offset int) (*structs.Messages, error) {

	m := structs.MessagesResponse{}
	err := api.CallMethod(apiMessagesGet, H{
		"count":  strconv.Itoa(count),
		"offset": strconv.Itoa(offset),
	}, &m)

	return &m.Response, err
}

func NewButton(label string, payload interface{}) structs.Button {
	button := structs.Button{}
	button.Action.Type = "text"
	button.Action.Label = label
	button.Action.Payload = "{}"
	if payload != nil {
		jPayoad, err := json.Marshal(payload)
		if err == nil {
			button.Action.Payload = string(jPayoad)
		}
	}
	button.Color = "default"
	return button
}

// GetFriendRequests - get friend requests
func (api *VkAPI) GetFriendRequests(out bool) (friends []int, err error) {
	p := H{}
	if out {
		p["out"] = "1"
	}

	r := structs.FriendRequestsResponse{}
	err = api.CallMethod(apiFriendsGetRequests, p, &r)

	return r.Response.Items, err
}

// User - get simple user info
func (api *VkAPI) User(uid int) (*structs.User, error) {

	r := structs.UsersResponse{}
	err := api.CallMethod(apiUsersGet, H{
		"user_ids": strconv.Itoa(uid),
		"fields":   "sex,screen_name, city, country, bdate",
	}, &r)

	if err != nil {
		return nil, err
	}
	if len(r.Response) > 0 {
		return r.Response[0], err
	}
	return nil, errors.New("no users returned")
}

// CheckFriends checking friend invites and matсhes and deletes mutual
func (bot *VKBot) CheckFriends() {
	uids, _ := API.GetFriendRequests(false)
	if len(uids) > 0 {
		for _, uid := range uids {
			API.AddFriend(uid)
			for k, v := range bot.actionRoutes {
				if k == "friend_add" {
					m := structs.Message{Action: "friend_add", UserID: uid}
					v(&m)
				}
			}
		}
	}
	uids, _ = API.GetFriendRequests(true)
	if len(uids) > 0 {
		for _, uid := range uids {
			API.DeleteFriend(uid)
			for k, v := range bot.actionRoutes {
				if k == "friend_delete" {
					m := structs.Message{Action: "friend_delete", UserID: uid}
					v(&m)
				}
			}
		}
	}
}

// AddFriend - add friend
func (api *VkAPI) AddFriend(uid int) bool {

	r := structs.SimpleResponse{}
	err := api.CallMethod(apiFriendsAdd, H{"user_id": strconv.Itoa(uid)}, &r)
	if err != nil {
		return false
	}

	return r.Response == 1
}

// DeleteFriend - delete friend
func (api *VkAPI) DeleteFriend(uid int) bool {

	u := structs.FriendDeleteResponse{}
	err := api.CallMethod(apiFriendsDelete, H{"user_id": strconv.Itoa(uid)}, &u)

	if err != nil {
		return false
	}

	return u.Response["success"] == 1
}

func (api *VkAPI) NotifyAdmin(msg string) (err error) {
	if api.AdminID != 0 {
		_, err = api.SendMessage(api.AdminID, msg)
	}
	return err
}
