package longpollServer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"vkbot/bot"
	"vkbot/debug"
	"vkbot/router"
	"vkbot/structs"
)

const LongPollVersion = 3

type GroupLongPollServer struct {
	Key             string
	Server          string
	Ts              string
	Wait            int
	Mode            int
	Version         int
	RequestInterval int
	NeedPts         bool
	API             *bot.VkAPI
	LpVersion       int
	ReadMessages    map[int]time.Time
}

type GroupLongPollServerResponse struct {
	Response GroupLongPollServer
}

type GroupLongPollMessage struct {
	Type   string
	Object struct {
		Date                  int                    `json:"date"`
		FromID                int                    `json:"from_id"`
		ID                    int                    `json:"id"`
		Out                   int                    `json:"out"`
		PeerID                int                    `json:"peer_id"`
		Text                  string                 `json:"text"`
		ConversationMessageID int                    `json:"conversation_message_id"`
		FwdMessages           []GroupLongPollMessage `json:"fwd_messages"`
		Important             bool                   `json:"important"`
		RandomID              int                    `json:"random_id"`
		// structs.Attachment[]  `json:"attachments"`
		IsHidden bool `json:"is_hidden"`
		Action   struct {
			Type     string
			MemberID int `json:"member_id"`
		}
	}
	GroupID int
}

type GroupLongPollEvent struct {
	Type    string
	GroupID int
}

type GroupLongPollResponse struct {
	Ts       string
	Messages []*structs.Message
}

func NewGroupLongPollServer(requestInterval int) (resp *GroupLongPollServer) {
	server := GroupLongPollServer{}
	server.Wait = DefaultWait
	server.Mode = DefaultMode
	server.Version = DefaultVersion
	server.RequestInterval = requestInterval
	server.ReadMessages = make(map[int]time.Time)
	return &server
}

func (server *GroupLongPollServer) Init() (err error) {
	r := GroupLongPollServerResponse{}
	err = router.API.CallMethod("groups.getLongPollServer", bot.H{
		"group_id": strconv.Itoa(router.API.GroupID),
	}, &r)
	server.Wait = DefaultWait
	server.Mode = DefaultMode
	server.Version = DefaultVersion
	server.RequestInterval = router.API.RequestInterval
	server.Server = r.Response.Server
	server.Ts = r.Response.Ts
	server.Key = r.Response.Key
	return err
}

func (server *GroupLongPollServer) Request() ([]byte, error) {
	var err error

	if server.Server == "" {
		err = server.Init()
		if err != nil {
			log.Fatal(err)
		}
	}

	parameters := url.Values{}
	parameters.Add("act", "a_check")
	parameters.Add("ts", server.Ts)
	parameters.Add("wait", strconv.Itoa(server.Wait))
	parameters.Add("key", server.Key)
	query := server.Server + "?" + parameters.Encode()
	if server.Server == "test" {
		content, err := ioutil.ReadFile("./mocks/longpoll.json")
		return content, err
	}
	resp, err := http.Get(query)
	if err != nil {
		debug.DebugPrint("%+v\n", err.Error())
		time.Sleep(time.Duration(time.Millisecond * time.Duration(server.RequestInterval)))
		return nil, err
	}
	buf, err := ioutil.ReadAll(resp.Body)
	time.Sleep(time.Duration(time.Millisecond * time.Duration(server.RequestInterval)))
	//debugPrint("longpoll vk resp: %+v\n", string(buf))

	failResp := structs.GroupFailResponse{}
	err = json.Unmarshal(buf, &failResp)
	if err != nil {
		return nil, err
	}
	switch failResp.Failed {
	case 1:
		server.Ts = failResp.Ts
		return server.Request()
	case 2:
		err = server.Init()
		if err != nil {
			log.Fatal(err)
		}
		return server.Request()
	case 3:
		err = server.Init()
		if err != nil {
			log.Fatal(err)
		}
		return server.Request()
	case 4:
		return nil, errors.New("vkapi: wrong longpoll version")
	default:
		server.Ts = failResp.Ts
		return buf, nil
	}
}

func (server *GroupLongPollServer) GetLongPollMessages() ([]*structs.Message, error) {
	resp, err := server.Request()
	if err != nil {
		return nil, err
	}
	messages, err := server.ParseLongPollMessages(string(resp))
	return messages.Messages, nil
}

func (server *GroupLongPollServer) ParseMessage(obj map[string]interface{}) structs.Message {
	msg := structs.Message{}
	msg.ID = structs.getJSONInt(obj["id"])
	msg.Body = obj["text"].(string)
	userID := structs.getJSONInt(obj["from_id"])
	if userID != 0 {
		msg.UserID = userID
	}
	msg.PeerID = structs.getJSONInt64(obj["peer_id"])
	if int64(msg.UserID) == msg.PeerID {
		msg.ChatID = 0
	} else {
		msg.ChatID = int(msg.PeerID)
	}
	msg.Date = structs.getJSONInt(obj["date"])

	fmt.Printf("%+v\n", msg)
	return msg
}

func (server *GroupLongPollServer) ParseLongPollMessages(j string) (*GroupLongPollResponse, error) {
	//fmt.Printf("\n>>>>>>>>>>>>>updates0: %+v\n\n", j)
	d := json.NewDecoder(strings.NewReader(j))
	d.UseNumber()
	var lp interface{}
	if err := d.Decode(&lp); err != nil {
		return nil, err
	}
	lpMap := lp.(map[string]interface{})
	result := GroupLongPollResponse{Messages: []*structs.Message{}}
	result.Ts = lpMap["ts"].(string)
	updates := lpMap["updates"].([]interface{})
	for _, event := range updates {
		eventType := event.(map[string]interface{})["type"].(string)
		if eventType == "message_new" {
			obj := event.(map[string]interface{})["object"].(map[string]interface{})
			out := structs.getJSONInt(obj["out"])
			if out == 0 {
				msg := server.ParseMessage(obj)
				result.Messages = append(result.Messages, &msg)
				fmt.Printf("\n>>>>>>>>>>>>>msg: %+v\n\n", msg)
			}
		}
	}
	if len(result.Messages) == 0 {
		fmt.Println(j)
	}
	fmt.Printf("\n>>>>>>>>>>>>>messages: %+v\n\n", result)
	return &result, nil
}

func (server *GroupLongPollServer) FilterReadMesages(messages []*structs.Message) (result []*structs.Message) {
	for _, m := range messages {
		t, ok := server.ReadMessages[m.ID]
		if ok {
			if time.Since(t).Minutes() > 1 {
				delete(server.ReadMessages, m.ID)
			}
		} else {
			result = append(result, m)
			server.ReadMessages[m.ID] = time.Now()
		}
	}
	return result
}
