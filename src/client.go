package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/c-fandango/rocketchat-term/creds"
	"github.com/c-fandango/rocketchat-term/utils"
	"github.com/gorilla/websocket"
	"net/url"
	"os"
	"os/signal"
	"time"
)

var homeDir, _ = os.UserHomeDir()
var cachePath = homeDir + "/.rocketchat-term"

type userSchema struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type cacheSchema struct {
	Host      string `json:"host"`
	Token     string `json:"token"`
	ExpiresAt int    `json:"expiresAt"`
}

type timestampSchema struct {
	TS int `json:"$date"`
}

type messageSchema struct {
	ID       string          `json:"_id"`
	RoomID   string          `json:"rid"`
	Content  string          `json:"msg"`
	SentTS   timestampSchema `json:"ts"`
	UpdateTS timestampSchema `json:"_updatedAt"`
	Sender   userSchema      `json:"u"`
}

type roomSchema struct {
	ID        string   `json:"_id"`
	ReadOnly  bool     `json:"ro"`
	Name      string   `json:"name"`
	Usernames []string `json:"usernames"`
	Messages  []messageSchema
}

type errorResponse struct {
	Error   int    `json:"error"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type wssRequest struct {
	ID      string `json:"id"`
	Message string `json:"msg"`
	Method  string `json:"method"`
	Name    string `json:"name"`
}

type wssResponse struct {
	ID         string        `json:"id"`
	Message    string        `json:"msg"`
	Error      errorResponse `json:"error"`
	Collection string        `json:"collection"`
}

type authResponse struct {
	wssResponse
	Result struct {
		Token   string          `json:"token"`
		Type    string          `json:"type"`
		Expires timestampSchema `json:"tokenExpires"`
	} `json:"result"`
	host string
}

func (a *authResponse) authenticateLdap(username string, password string) string {

	type ldapParams struct {
		Ldap        bool              `json:"ldap"`
		Username    string            `json:"username"`
		LdapPass    string            `json:"ldapPass"`
		LdapOptions map[string]string `json:"ldapOptions"`
	}

	a.ID = utils.RandID(5)

	request := struct {
		wssRequest
		Params []ldapParams `json:"params"`
	}{
		wssRequest: wssRequest{
			ID:      a.ID,
			Message: "method",
			Method:  "login",
		},
		Params: []ldapParams{
			ldapParams{
				Ldap:        true,
				Username:    username,
				LdapPass:    password,
				LdapOptions: map[string]string{},
			},
		}}

	message, _ := json.Marshal(request)

	return string(message)
}

func (a *authResponse) authenticateToken(token string) string {

	a.ID = utils.RandID(5)

	request := struct {
		wssRequest
		Params []map[string]string `json:"params"`
	}{
		wssRequest: wssRequest{
			ID:      a.ID,
			Message: "method",
			Method:  "login",
		},
		Params: []map[string]string{
			map[string]string{"resume": token},
		}}

	message, _ := json.Marshal(request)

	return string(message)
}

func (a *authResponse) handleResponse(response []byte) error {
	json.Unmarshal(response, a)

	if a.Error != (errorResponse{}) {
		creds.ClearCache(cachePath)
		return errors.New("authorisation failed")
	}

	tokenCache := cacheSchema{
		Host:      a.host,
		Token:     a.Result.Token,
		ExpiresAt: a.Result.Expires.TS,
	}

	cache, _ := json.Marshal(tokenCache)

	err := creds.WriteCache(cachePath, cache)

	if err != nil {
		return err
	}

	return nil
}

type rooms struct {
	wssResponse
	Result struct {
		Rooms []roomSchema `json:"update"`
	} `json:"result"`
}

func (r *rooms) handleResponse(response []byte) error {
	json.Unmarshal(response, r)

	if r.Error != (errorResponse{}) {
		return errors.New("failed to fetch room data")
	}

	for i, room := range r.Result.Rooms {
		if room.Name == "" {
			r.Result.Rooms[i].Name = makeRoomName(room.Usernames)
		}
	}

	return nil
}

func (r *rooms) constructRequest() string {
	r.ID = utils.RandID(5)

	request := struct {
		wssRequest
		Params []map[string]int `json:"params"`
	}{
		wssRequest: wssRequest{
			ID:      r.ID,
			Message: "method",
			Method:  "rooms/get",
		},
		Params: []map[string]int{
			map[string]int{
				"$date": 0,
			},
		},
	}
	message, _ := json.Marshal(request)

	return string(message)
}

type subscription struct {
	wssResponse
	Fields struct {
		EventName string          `json:"eventName"`
		Messages  []messageSchema `json:"args"`
	} `json:"fields"`
}

func (s *subscription) handleResponse(response []byte, allRooms []roomSchema) error {
	const newMessageAllowedDelayMS = 400

	json.Unmarshal(response, s)
	for _, message := range s.Fields.Messages {
		var roomName string

		if message.UpdateTS.TS > message.SentTS.TS+newMessageAllowedDelayMS {
			continue
		}

		for _, room := range allRooms {
			if room.ID == message.RoomID {
				roomName = room.Name
				break
			}
		}

		if message.Content != "" {
			printMessage(roomName, message.Sender.Name, message.Content, message.SentTS.TS)
		}

	}

	if s.Error != (errorResponse{}) {
		return errors.New("failed to fetch room data")
	}

	return nil
}

func (s *subscription) constructRequest(roomID string) string {
	s.ID = utils.RandID(5)

	request := struct {
		wssRequest
		Params []string `json:"params"`
	}{
		wssRequest: wssRequest{
			ID:      s.ID,
			Message: "sub",
			Name:    "stream-room-messages",
		},
		Params: []string{
			roomID,
			"false",
		},
	}
	message, _ := json.Marshal(request)

	return string(message)
}

func main() {
	fmt.Println("hello world")

	var auth authResponse

	credentials, err := getCredentials(cachePath)
	if err != nil {
		fmt.Println(err)
	}
	auth.host = credentials["host"]

	u := url.URL{Scheme: "wss", Host: credentials["host"], Path: "/websocket"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		panic("invalid host")
	}

	defer c.Close()

	var allRooms rooms
	var roomSub subscription
	roomSub.Collection = "stream-room-messages"
	connectMessage := `{"msg": "connect","version": "1","support": ["1"]}`
	pongMessage := `{"msg": "pong"}`

	done := make(chan struct{})
	messageOut := make(chan string)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		defer close(done)
		messageOut <- connectMessage

		// TODO don't busy loop
		for {
			_, response, err := c.ReadMessage()

			if err != nil {
				fmt.Println("error in reading incoming message ", err)
				return
			}

			var data wssResponse
			json.Unmarshal(response, &data)

			if data.Message == "connected" {
				if credentials["token"] != "" {
					messageOut <- auth.authenticateToken(credentials["token"])
				} else {
					messageOut <- auth.authenticateLdap(credentials["username"], credentials["password"])
				}
			} else if data.ID == auth.ID && data.Message == "result" {
				err := auth.handleResponse(response)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("authenticated")

				messageOut <- allRooms.constructRequest()

			} else if data.ID == allRooms.ID && data.Message == "result" {
				err := allRooms.handleResponse(response)
				if err != nil {
					fmt.Println(err)
					return
				}

				messageOut <- roomSub.constructRequest("__my_messages__")

			} else if data.Collection == roomSub.Collection && data.Message == "changed" {
				err := roomSub.handleResponse(response, allRooms.Result.Rooms)
				if err != nil {
					fmt.Println(err)
					return
				}

			} else if data.Message == "ping" {
				messageOut <- pongMessage
			}
		}
	}()

eventLoop:
	for {
		select {
		case <-done:
			break eventLoop
		case m := <-messageOut:

			err := c.WriteMessage(websocket.TextMessage, []byte(m))

			if err != nil {
				fmt.Println("error sending websocket message ", err)
			}
		case <-interrupt:

			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

			if err != nil {
				fmt.Println("error closing websocket", err)
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			break eventLoop
		}
	}
}
