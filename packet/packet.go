package packet

import (
	"bytes"
	"encoding/binary"
)

type Message struct {
	Route  uint32
	Buffer []byte
}

func Pack(message *Message) ([]byte, error) {
	var buf bytes.Buffer

	length := len(message.Buffer) + 4 + 2
	buf.Grow(length)

	binary.Write(&buf, binary.BigEndian, length)
	binary.Write(&buf, binary.BigEndian, message.Route)
	binary.Write(&buf, binary.BigEndian, message.Buffer)

	return buf.Bytes(), nil
}

func Unpack(data []byte) (*Message, error) {
	reader := bytes.NewReader(data)

	var length int16
	binary.Read(reader, binary.BigEndian, &length)

	message := &Message{Buffer: make([]byte, length)}
	binary.Read(reader, binary.BigEndian, &message.Route)
	binary.Read(reader, binary.BigEndian, &message.Buffer)

	return message, nil
}
