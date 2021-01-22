package messages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

type ServerType string

const (
	ServerTypePatch ServerType = "p"
	ServerTypeError ServerType = "error"
)

type ClientType string

const (
	ClientTypeSubmit ClientType = "submit"
	ClientTypeClick  ClientType = "click"
	ClientTypeInit   ClientType = "init"
)

func marshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

type ServerMessage struct {
	Type        ServerType
	Data        interface{}
	ComponentID string
}

// format is:
// [<event type>, [<data>], <componentID>]
func (sm ServerMessage) Serialize() ([]byte, error) {
	onWire := make([]json.RawMessage, 1, 3)
	var err error
	onWire[0], err = json.Marshal(sm.Type)
	if err != nil {
		return nil, err
	}

	if sm.Data == nil {
		return nil, nil
	}
	onWire = append(onWire, nil)
	// we must use our internal marshal so that html is not escaped
	onWire[1], err = marshal(sm.Data)
	if err != nil {
		return nil, err
	}

	if sm.ComponentID == "" {
		return nil, nil
	}
	onWire = append(onWire, nil)
	onWire[2], err = json.Marshal(sm.ComponentID)
	if err != nil {
		return nil, err
	}
	// we must use our internal marshal so that html is not escaped
	return marshal(onWire)
}

type ClientMessage struct {
	Type           ClientType
	Data           []byte
	ComponentID    string
	Name           string
	SubcomponentID *int
}

func (cm *ClientMessage) UnmarshalJSON(data []byte) error {
	// // format is:
	// // [<event type>, [<data>], <componentID>, <subcomponentID>,  <name>, ]
	raw := [5]json.RawMessage{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if err := json.Unmarshal(raw[0], &cm.Type); err != nil {
		return err
	}
	cm.Data = raw[1]
	if raw[2] == nil {
		return nil
	}
	if err := json.Unmarshal(raw[2], &cm.ComponentID); err != nil {
		return err
	}
	if raw[3] == nil {
		return nil
	}
	if err := json.Unmarshal(raw[3], &cm.Name); err != nil {
		return err
	}
	if raw[4] == nil {
		return nil
	}
	var subcomponentID string
	if err := json.Unmarshal(raw[4], &subcomponentID); err != nil {
		return err
	}
	sID, err := strconv.Atoi(subcomponentID)
	if err != nil {
		return fmt.Errorf("subcomponentID of %q is not in the correct format", subcomponentID)
	}
	cm.SubcomponentID = &sID
	return nil
}
