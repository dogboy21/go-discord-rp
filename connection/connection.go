package connection

import (
	"os"
	"net"
	"encoding/binary"
	"bytes"
	"log"
	"time"
	"fmt"
	"crypto/rand"
)

var socket net.Conn

func GetSocketPath() string {
	path, exists := os.LookupEnv("XDG_RUNTIME_DIR")

	if exists {
		return path
	}

	path, exists = os.LookupEnv("TMPDIR")

	if exists {
		return path
	}

	path, exists = os.LookupEnv("TMP")

	if exists {
		return path
	}

	path, exists = os.LookupEnv("TEMP")

	if exists {
		return path
	}

	return "/tmp"
}

func OpenSocket(appId string) {
	c, err := net.Dial("unix", GetSocketPath() + "/discord-ipc-0")

	if err != nil {
		panic(err)
	}

	log.Println("Socket opened")

	socket = c

	go read()

	Handshake(appId)
}

func read() {
	for {
		buf := make([]byte, 512)
		nr, err := socket.Read(buf)
		if err != nil {
			return
		}

		buffer := new(bytes.Buffer)
		for i := 8; i < nr; i++ {
			buffer.WriteByte(buf[i])
		}

		log.Println(string(buffer.Bytes()))
	}
}

func getNonce() string {
	buf := make([]byte, 16)
	rand.Read(buf)
	buf[6] = (buf[6] & 0x0f) | 0x40

	return fmt.Sprintf("%x-%x-%x-%x-%x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:])
}

func SendFramed(opcode int, msg string) {
	log.Printf("> %s\n", msg)
	msgBytes := []byte(msg)

	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, int32(opcode))
	if err != nil {
		log.Println(err)
	}

	err = binary.Write(buf, binary.LittleEndian, int32(len(msgBytes)))
	if err != nil {
		log.Println(err)
	}

	buf.Write(msgBytes)

	socket.Write(buf.Bytes())
}

func Handshake(appId string) {
	SendFramed(0, "{\"v\":1, \"client_id\": \"" + appId + "\"}")
	time.Sleep(3 * time.Second)
}

func SetActivity(state, details, smallImg, smallText, largeImg, largeText string) {
	pid := os.Getpid()
	activity := "{\"cmd\": \"SET_ACTIVITY\", \"args\":{\"pid\": %d, \"activity\": {\"state\": \"%s\", \"details\": \"%s\", \"instance\": true, \"assets\": {\"small_image\": \"%s\", \"small_text\": \"%s\", \"large_image\": \"%s\", \"large_text\": \"%s\"}}}, \"nonce\": \"%s\"}"
	SendFramed(1, fmt.Sprintf(activity, pid, state, details, smallImg, smallText, largeImg, largeText, getNonce()))
}

func SetActivityText(state, details string) {
	pid := os.Getpid()
	activity := "{\"cmd\": \"SET_ACTIVITY\", \"args\":{\"pid\": %d, \"activity\": {\"state\": \"%s\", \"details\": \"%s\", \"instance\": true}}, \"nonce\": \"%s\"}"
	SendFramed(1, fmt.Sprintf(activity, pid, state, details, getNonce()))
}
