package message

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-clog/clog"
)

// Message struct that has aws sqsmessage and body is parsed "json format"
type Message struct {
	*sqs.Message
	body map[string]interface{}
}

// NewMessage return the message with pre-conf defined
func NewMessage(sqs *sqs.Message) *Message {
	msg := &Message{
		Message: sqs,
		body:    make(map[string]interface{}),
	}

	if err := json.Unmarshal(msg.Bytes(), &msg.body); err != nil {
		clog.Warn("error: %v", err.Error())
	}

	return msg
}

// String returns the body message as string
func (m *Message) String() string {
	return *m.Body
}

// Bytes returns the body message as bytes
func (m *Message) Bytes() []byte {
	return []byte(m.String())
}

// Get if body is json this return the field with name.
func (m *Message) Get(name string) (interface{}, bool) {
	rs, has := m.body[name]
	return rs, has
}

// Topic try to return topic if it is defined
func (m *Message) Topic() string {
	msgTopic, has := m.Get("TopicArn")
	if !has {
		return ""
	}

	msg, ok := msgTopic.(string)
	if !ok {
		return ""
	}

	slices := strings.Split(msg, ":")
	return slices[len(slices)-1]
}
