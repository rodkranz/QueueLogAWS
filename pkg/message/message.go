package message

import (
	"encoding/json"
	"strings"
	
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-clog/clog"
)

type Message struct {
	*sqs.Message
	body map[string]interface{}
}

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

func (m *Message) String() string {
	return *m.Body
}

func (m *Message) Bytes() []byte {
	return []byte(m.String())
}

func (m *Message) Get(name string) (interface{}, bool) {
	rs, has := m.body[name]
	return rs, has
}

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
