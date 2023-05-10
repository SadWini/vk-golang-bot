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

const DefaultWait = 25 // DefaultWait Connection will be automatically stopped after 30 seconds
const (
	LongPollModeGetAttachments     = 2
	LongPollModeGetExtendedEvents1 = 1 << 3
	LongPollModeGetPts             = 2 << (2 + iota)
	LongPollModeGetExtraData
	LongPollModeGetRandomID
)

const DefaultMode = LongPollModeGetAttachments
const DefaultVersion = 2
const ChatPrefix = 2000000000

type LongPollServer interface {
	Init() (err error)
	Request() ([]byte, error)
	GetLongPollMessages() ([]*structs.Message, error)
	FilterReadMesages(messages []*structs.Message) (result []*structs.Message)
}

type UserLongPollServer struct {
	Key             string
	Server          string
	Ts              int
	Wait            int
	Mode            int
	Version         int
	RequestInterval int
	NeedPts         bool
	API             *bot.VkAPI
	LpVersion       int
	ReadMessages    map[int]time.Time
}

type UserLongPollServerResponse struct {
	Response UserLongPollServer
}

type LongPollUpdate []interface{}
type LongPollUpdateNum []int64

type LongPollMessage struct {
	MessageType int
	MessageID   int
	Flags       int
	PeerID      int64
	Timestamp   int64
	Text        string
	Attachments []structs.Attachment
	RandomID    int
}

type LongPollResponse struct {
	Ts       uint
	Messages []*structs.Message
}

func NewUserLongPollServer(needPts bool, lpVersion int, requestInterval int) (resp *UserLongPollServer) {
	server := UserLongPollServer{}
	server.NeedPts = needPts
	server.Wait = DefaultWait
	server.Mode = DefaultMode
	server.Version = DefaultVersion
	server.RequestInterval = requestInterval
	server.LpVersion = lpVersion
	server.ReadMessages = make(map[int]time.Time)
	return &server
}
func (server *UserLongPollServer) Init() (err error) {
	r := UserLongPollServerResponse{}
	pts := 0
	if server.NeedPts {
		pts = 1
	}
	err = router.API.CallMethod("messages.getLongPollServer", bot.H{
		"need_pts": strconv.Itoa(pts),
		"message":  strconv.Itoa(server.LpVersion),
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
func (server *UserLongPollServer) Request() ([]byte, error) {
	var err error

	if server.Server == "" {
		err = server.Init()
		if err != nil {
			log.Fatal(err)
		}
	}

	parameters := url.Values{}
	parameters.Add("act", "a_check")
	parameters.Add("ts", strconv.Itoa(server.Ts))
	parameters.Add("wait", strconv.Itoa(server.Wait))
	parameters.Add("key", server.Key)
	parameters.Add("mode", strconv.Itoa(DefaultMode))
	parameters.Add("version", strconv.Itoa(server.Version))
	query := "https://" + server.Server + "?" + parameters.Encode()
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

	failResp := structs.FailResponse{}
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
		return buf, nil
	}
}
func GetLongPollMessage(resp []interface{}) *structs.Message {
	message := structs.Message{}
	mid, _ := resp[1].(json.Number).Int64()
	message.ID = int(mid)
	flags, _ := resp[2].(json.Number).Int64()
	message.Flags = int(flags)
	message.PeerID, _ = resp[3].(json.Number).Int64()
	message.Timestamp, _ = resp[4].(json.Number).Int64()
	message.Body = resp[5].(string)
	return &message
}

func (server *UserLongPollServer) GetLongPollMessages() ([]*structs.Message, error) {
	resp, err := server.Request()
	if err != nil {
		return nil, err
	}
	messages, err := server.ParseLongPollMessages(string(resp))
	return messages.Messages, nil
}

func (server *UserLongPollServer) ParseLongPollMessages(j string) (*LongPollResponse, error) {
	d := json.NewDecoder(strings.NewReader(j))
	d.UseNumber()
	var lp interface{}
	if err := d.Decode(&lp); err != nil {
		return nil, err
	}
	lpMap := lp.(map[string]interface{})
	result := LongPollResponse{Messages: []*structs.Message{}}
	ts, _ := lpMap["ts"].(json.Number).Int64()
	result.Ts = uint(ts)
	updates := lpMap["updates"].([]interface{})
	for _, event := range updates {
		el := event.([]interface{})
		eventType := getJSONInt(el[0])
		if eventType == 4 {
			out := getJSONInt(el[2]) & 2
			if out == 0 {
				msg := structs.Message{}
				msg.ID = getJSONInt(el[1])
				msg.Body = el[5].(string)
				userID := el[6].(map[string]interface{})["from"]
				if userID != nil {
					msg.UserID, _ = strconv.Atoi(userID.(string))
				}
				msg.PeerID = getJSONInt64(el[3])
				if msg.UserID == 0 {
					msg.UserID = int(msg.PeerID)
				} else {
					msg.ChatID = int(msg.PeerID - ChatPrefix)
				}
				msg.Date = getJSONInt(el[4])
				fmt.Println(msg.Body)
				result.Messages = append(result.Messages, &msg)
			}
		}
	}
	if len(result.Messages) == 0 {
		fmt.Println(j)
	}
	result.Messages = server.FilterReadMesages(result.Messages)
	return &result, nil
}

func (server *UserLongPollServer) FilterReadMesages(messages []*structs.Message) (result []*structs.Message) {
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

func getJSONInt64(el interface{}) int64 {
	if el == nil {
		return 0
	}
	v, _ := el.(json.Number).Int64()
	return v
}

func getJSONInt(el interface{}) int {
	return int(getJSONInt64(el))
}
