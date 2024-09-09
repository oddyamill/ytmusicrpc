package discord

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"

	"github.com/oddyamill/ytmusicrpc/ipc"
)

func SendHandshake(clientId string) {
	ipc.Send(0, []byte(`{"v":1,"client_id":"`+clientId+`","nonce":"1"}`))
}

type Packet struct {
	Cmd   string `json:"cmd"`
	Args  Args   `json:"args"`
	Nonce string `json:"nonce"`
}

type Args struct {
	Pid      int       `json:"pid"`
	Activity *Activity `json:"activity"`
}

type Activity struct {
	Type       int         `json:"type"`
	Details    string      `json:"details"`
	State      string      `json:"state"`
	Assets     Assets      `json:"assets"`
	Timestamps *Timestamps `json:"timestamps,omitempty"`
	Buttons    []Button    `json:"buttons"`
}

type Assets struct {
	LargeImage string `json:"large_image"`
	LargeText  string `json:"large_text,omitempty"`
}

type Timestamps struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

type Button struct {
	Label string `json:"label"`
	Url   string `json:"url"`
}

func UpdatePresence(activity Activity) {
	payload, err := json.Marshal(Packet{
		"SET_ACTIVITY",
		Args{
			os.Getpid(),
			&activity,
		},
		getNonce(),
	})

	if err != nil {
		panic(err)
	}

	ipc.Send(1, payload)
}

func DeletePresence() {
	payload, err := json.Marshal(Packet{
		"SET_ACTIVITY",
		Args{
			os.Getpid(),
			nil,
		},
		getNonce(),
	})

	if err != nil {
		panic(err)
	}

	ipc.Send(1, payload)
}

// https://github.com/hugolgst/rich-go/blob/74618cc1ace23ea759a4be6a77ebc928b4d8c996/client/client.go#L67-L77

func getNonce() string {
	buf := make([]byte, 16)

	_, err := rand.Read(buf)

	if err != nil {
		panic(err)
	}

	buf[6] = (buf[6] & 0x0f) | 0x40

	return fmt.Sprintf("%x-%x-%x-%x-%x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:])
}
