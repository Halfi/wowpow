package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	"wowpow/pkg/api/message"
)

func TestProto(t *testing.T) {
	p := &message.Message{
		Header: message.Message_resource,
		Response: &message.Message_Payload{
			Payload: "Never close your hands in the face of evil fate; get up and rebel against it\n",
		},
	}

	bin, err := proto.Marshal(p)
	assert.Nil(t, err)

	msg := new(message.Message)
	err = proto.Unmarshal(bin, msg)
	assert.Nil(t, err)
}
