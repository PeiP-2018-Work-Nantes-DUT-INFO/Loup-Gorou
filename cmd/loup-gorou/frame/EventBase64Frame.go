package frame

import (
	b64 "encoding/base64"
	"loupgorou/cmd/loup-gorou/gonest"

	"google.golang.org/protobuf/proto"
)

func EncodeEventB64(message *gonest.Event) (encoded string, err error) {
	out, err := proto.Marshal(message)
	if err != nil {
		return
	}
	encoded = b64.StdEncoding.EncodeToString(out) + "\n"
	return
}

func DecodeEventB64(encoded string) (event *gonest.Event, err error) {
	buf, err := b64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return
	}
	event = &gonest.Event{}
	err = proto.Unmarshal(buf, event)
	return
}
