package ipc

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"time"

	npipe "gopkg.in/natefinch/npipe.v2"
)

// https://github.com/hugolgst/rich-go/blob/master/ipc/ipc_windows.go

var socket net.Conn

func init() {
	sock, err := npipe.DialTimeout(`\\.\pipe\discord-ipc-0`, time.Second*2)

	if err != nil {
		panic(err)
	}

	log.Print("connected to Discord IPC")
	socket = sock
}

// https://github.com/hugolgst/rich-go/blob/master/ipc/ipc.go

func Read() string {
	buf := make([]byte, 512)
	payloadlength, err := socket.Read(buf)

	if err != nil {
		panic(err)
	}

	buffer := new(bytes.Buffer)

	for i := 8; i < payloadlength; i++ {
		buffer.WriteByte(buf[i])
	}

	return buffer.String()
}


func Send(opcode int, payload []byte) string {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, int32(opcode))

	if err != nil {
		panic(err)
	}

	err = binary.Write(buf, binary.LittleEndian, int32(len(payload)))

	if err != nil {
		panic(err)
	}

	buf.Write(payload)

	_, err = socket.Write(buf.Bytes())

	if err != nil {
		panic(err)
	}

	return Read()
}
