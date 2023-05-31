package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"golang.org/x/term"
	"log"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type userSchema struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type messageSchema struct {
	ID        string `json:"_id"`
	RoomID    string `json:"rid"`
	Content   string `json:"msg"`
	Timestamp struct {
		Timestamp int `json:"$date"`
	} `json:"ts"`
	Sender userSchema `json:"u"`
}

type roomSubSchema struct {
	Message    string `json:"msg"`
	ID         string `json:"id"`
	Collection string `json:"collection"`
	Fields     struct {
		EventName string          `json:"eventName"`
		Messages  []messageSchema `json:"args"`
	} `json:"fields"`
}

type roomSchema struct {
	ID        string   `json:"_id"`
	ReadOnly  bool     `json:"ro"`
	Name      string   `json:"name"`
	Usernames []string `json:"usernames"`
        Messages []messageSchema
}

type roomsSchema struct {
	Message string `json:"msg"`
	ID      string `json:"id"`
	Result  struct {
		Rooms []roomSchema `json:"update"`
	} `json:"result"`
}

type wssResponse struct {
	Message string `json:"msg"`
	Session string `json:"session"`
	Version string `json:"version"`
	Support []int  `json:"support"`
}

func getCredentials() (string, string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Host Domain: ")
	host, err := reader.ReadString('\n')
	if err != nil {
		return "", "", "", err
	}

	fmt.Print("Enter Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", "", err
	}

	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", "", err
	}

	return strings.TrimSpace(host), strings.TrimSpace(username), string(bytePassword), nil
}

func printMessage(room string, user string, content string, timestamp int) {

	colourTextCodes := []string{
		"31", "32", "33", "34", "35", "36",
	}
	colourBackgroundCodes := []string{"41", "42", "43", "44", "45", "46"}

	resetColour := "\033[0m"
	blackText := "\033[30m"

	userColour := "\033[" + colourTextCodes[len(user)%(len(colourTextCodes)-1)] + "m"
	roomColour := "\033[" + colourBackgroundCodes[len(room)%(len(colourBackgroundCodes)-1)] + "m"

	user = userColour + user + resetColour
	room = roomColour + blackText + " " + room + " " + resetColour

        ts:=time.UnixMilli(int64(timestamp))
	timePretty := ts.Format(time.Kitchen)

	//fmt.Println("        -------")
	fmt.Println("       ", timePretty, "     ", room, "      ", user, "     ", content)
}

func main() {
	fmt.Println("hello world")
	f, err := os.OpenFile("./rocketchat.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	messageOut := make(chan string)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	host, username, password, _ := getCredentials()

	u := url.URL{Scheme: "wss", Host: host, Path: "/websocket"}
	log.Printf("connecting to %s", u.String())

	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		log.Printf("handshake failed with status %d", resp.StatusCode)
		log.Fatal("dial:", err)
	}

	defer c.Close()

	connectMessage := "{\"msg\": \"connect\",\"version\": \"1\",\"support\": [\"1\"]}"
	pongMessage := "{\"msg\": \"pong\"}"
	var sessionId string

	state := "resting"
	rooms := []roomSchema{}
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, response, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", response)
			fmt.Println("recv: %s", string(response))


			var data wssResponse

			json.Unmarshal(response, &data)

			if data.Message == "" && state == "resting" {
				state = "pre-connect"
				messageOut <- connectMessage
			} else if data.Message == "ping" {
				log.Printf("ping message recieved: %s", data.Message)
				messageOut <- pongMessage
			} else if data.Message == "connected" && state == "pre-connect" {
				state = "connected"
				sessionId = data.Session
				log.Printf("connected!, session id: %s", sessionId)
				authMessage := "{\"msg\": \"method\",\"method\": \"login\",\"id\":\"42\",\"params\":[{\"ldap\": true,   \"username\": \"" + username + "\",\"ldapPass\": \"" + password + "\",\"ldapOptions\": {} }]}"
				messageOut <- authMessage
			} else if data.Message == "result" && state == "connected" {
				state = "authenticated"
				fmt.Println("\nauthenticated!")
				roomRequest := "{\"msg\": \"method\",\"method\": \"rooms/get\",\"id\": \"42\",\"params\": [ { \"$date\": 0} ]}"
				messageOut <- roomRequest
			} else if data.Message == "result" && state == "authenticated" {
				state = "rooms fetched"
				var roomData roomsSchema
				json.Unmarshal(response, &roomData)
				for _, room := range roomData.Result.Rooms {
					if room.Name == "" {
						room.Name = strings.Join(room.Usernames, ", ")
					}
					rooms = append(rooms, room)
				}
				updateRequest := "{\"msg\":\"sub\",\"id\":\"" + strconv.Itoa(rand.Int()) + "\",\"name\":\"stream-room-messages\",\"params\":[\"" + "__my_messages__" + "\", false]}"
				messageOut <- updateRequest
			} else if data.Message == "ready" && state == "rooms fetched" {
				state = "subscribed"
			} else if data.Message == "changed" && state == "subscribed" {
				var roomSubData roomSubSchema
				json.Unmarshal(response, &roomSubData)

				for _, message := range roomSubData.Fields.Messages {

					if message.Content != "" {
						var roomName string
						for _, room := range rooms {
							if room.ID == message.RoomID {
								roomName = room.Name
								break
							}
						}
						printMessage(roomName, message.Sender.Name, message.Content, message.Timestamp.Timestamp)
					}
				}
			} else {
				log.Printf("next thing")
			}
		}
	}()

        eventLoop:
	for {
		select {
		case <-done:
		case m := <-messageOut:

			log.Printf("Send Message %s", m)
			err := c.WriteMessage(websocket.TextMessage, []byte(m))

			if err != nil {
				log.Println("write:", err)
			}
		case <-interrupt:

			log.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

			if err != nil {
				log.Println("write close:", err)
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
                        break eventLoop
		}
	}
}
