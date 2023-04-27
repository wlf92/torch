package packet

import (
	"bytes"
	"encoding/binary"
)

type Message struct {
	Route  uint32
	Buffer []byte
}

func Pack(message *Message) []byte {
	var buf bytes.Buffer

	length := len(message.Buffer) + 4 + 2
	buf.Grow(length)

	binary.Write(&buf, binary.BigEndian, int16(length))
	binary.Write(&buf, binary.BigEndian, message.Route)
	binary.Write(&buf, binary.BigEndian, message.Buffer)

	return buf.Bytes()
}

func Unpack(data []byte) (*Message, error) {
	reader := bytes.NewReader(data)

	var length int16
	binary.Read(reader, binary.BigEndian, &length)

	message := &Message{Buffer: make([]byte, length-4-2)}
	binary.Read(reader, binary.BigEndian, &message.Route)
	binary.Read(reader, binary.BigEndian, &message.Buffer)

	return message, nil
}
