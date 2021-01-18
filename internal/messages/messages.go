package messages

import (
	"encoding/json"
	"fmt"
)

type ServerType string

const (
	ServerTypePatch ServerType = "p"
	ServerTypeError ServerType = "error"
)

type ClientType string

const (
	ClientTypeSubmit ClientType = "submit"
	ClientTypeInit   ClientType = "init"
)

// format is:
// [<componentID>, <event type>, [<data>]]
type ServerMessage [3]json.RawMessage

func NewServerMessage(componentID string, event string, data interface{}) (ServerMessage, error) {
	var msg ServerMessage
	msg[0], _ = json.Marshal(componentID)
	msg[1], _ = json.Marshal(event)
	var err error
	msg[2], err = json.Marshal(data)
	return msg, err
}

func (m ServerMessage) String() string {
	return fmt.Sprint(string(m[0]), string(m[1]), string(m[2]))
}

// format is:
// [<componentID>, <event type>, <name>, [<data>]]
type ClientMessage [4]json.RawMessage

func (m ClientMessage) String() string {
	return fmt.Sprint(string(m[0]), string(m[1]), string(m[2]), string(m[3]))
}

func (m ClientMessage) ComponentID() string {
	var out string
	_ = json.Unmarshal(m[0], &out)
	return out
}

func (m ClientMessage) EventType() ClientType {
	var out ClientType
	_ = json.Unmarshal(m[1], &out)
	return out
}

func (m ClientMessage) Name() string {
	var out string
	_ = json.Unmarshal(m[2], &out)
	return out
}
func (m ClientMessage) UnmarshalBody(v interface{}) error {
	return json.Unmarshal(m[3], v)
}
