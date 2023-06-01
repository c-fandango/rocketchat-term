package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/c-fandango/rocketchat-term/utils"
	"fmt"
	"github.com/gorilla/websocket"
	"golang.org/x/term"
	"log"
	"net/url"
	"os"
	"os/signal"
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

type subscription struct {
	wssResponse
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
	Messages  []messageSchema
}

type errorResponse struct {
	Error   int    `json:"error"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type wssResponse struct {
	ID      string        `json:"id"`
	Message string        `json:"msg"`
	Error   errorResponse `json:"error"`
	//do i need this?
	Session string `json:"session"`
}

type wssRequest struct {
	ID      string `json:"id"`
	Message string `json:"msg"`
	Method  string `json:"method"`
	Name    string `json:"name"`
}

// add match room id method
type rooms struct {
	wssResponse
	Result struct {
		Rooms []roomSchema `json:"update"`
	} `json:"result"`
}

func (s *subscription) handleResponse(response []byte, allRooms []roomSchema) error {
	json.Unmarshal(response, s)
	for _, message := range s.Fields.Messages {
		if message.Content == "" {
			continue
		}

		var roomName string
		for _, room := range allRooms {
			if room.ID == message.RoomID {
				roomName = room.Name
				break
			}
		}
		printMessage(roomName, message.Sender.Name, message.Content, message.Timestamp.Timestamp)
	}

	if s.Error != (errorResponse{}) {
		return errors.New("failed to fetch room data")
	}

	return nil
}

func (s *subscription) constructRequest(roomID string) string {
	//fix this bug? rocket always responds with "id: id" regargless
	//s.ID = utils.RandID(5)
	s.ID = "id"

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

func (r *rooms) handleResponse(response []byte) error {
	json.Unmarshal(response, r)

	if r.Error != (errorResponse{}) {
		return errors.New("failed to fetch room data")
	}

	for i, room := range r.Result.Rooms {
		if room.Name == "" {
			r.Result.Rooms[i].Name = strings.Join(room.Usernames, ", ")
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

func (w *wssResponse) authenticate(username string, password string) string {

	type ldapParams struct {
		Ldap        bool              `json:"ldap"`
		Username    string            `json:"username"`
		LdapPass    string            `json:"ldapPass"`
		LdapOptions map[string]string `json:"ldapOptions"`
	}

	w.ID = utils.RandID(5)

	request := struct {
		wssRequest
		Params []ldapParams `json:"params"`
	}{
		wssRequest: wssRequest{
			ID:      w.ID,
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

// TODO cache the token
func (w *wssResponse) handleResponse(response []byte) error {
	json.Unmarshal(response, w)

	if w.Error != (errorResponse{}) {
		return errors.New("authorisation failed")
	}

	fmt.Println("authenticated")

	return nil
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

	ts := time.UnixMilli(int64(timestamp))
	timePretty := ts.Format(time.Kitchen)

	timePretty = utils.PadLeft(timePretty, "0", 7)

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

	var auth wssResponse
	var roomSub subscription
	var allRooms rooms

	done := make(chan struct{})

	go func() {
		defer close(done)
		connectMessage := `{"msg": "connect","version": "1","support": ["1"]}`
		messageOut <- connectMessage

		for {
			_, response, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", response)

			var data wssResponse

			json.Unmarshal(response, &data)

			if data.Message == "connected" {
				messageOut <- auth.authenticate(username, password)

			} else if data.ID == auth.ID && data.Message == "result" {
				err := auth.handleResponse(response)
				if err != nil {
					fmt.Println(err)
					break
				}

				messageOut <- allRooms.constructRequest()

			} else if data.ID == allRooms.ID && data.Message == "result" {
				err := allRooms.handleResponse(response)
				if err != nil {
					fmt.Println(err)
					break
				}

				messageOut <- roomSub.constructRequest("__my_messages__")

			} else if data.ID == roomSub.ID && data.Message == "changed" {
				err := roomSub.handleResponse(response, allRooms.Result.Rooms)
				if err != nil {
					fmt.Println(err)
					break
				}

			} else if data.Message == "ping" {
				messageOut <- `{"msg": "pong"}`
			}
		}
	}()

eventLoop:
	for {
		select {
		case <-done:
			break eventLoop
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
