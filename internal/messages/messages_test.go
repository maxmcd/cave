package messages

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestClientMessage_UnmarshalJSON(t *testing.T) {
	type fields struct {
		ComponentID    string
		SubcomponentID *int
		Type           ClientType
		Name           string
		Data           []byte
	}
	type args struct {
		data []byte
	}

	zero := 0

	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		errContains string
	}{
		{
			name: "minimal",
			args: args{data: []byte(`["init", {"foo":1}]`)},
			fields: fields{
				Type: ClientTypeInit,
				Data: []byte(`{"foo":1}`),
			},
		},
		{
			name: "componentID",
			args: args{data: []byte(`["init", [1,2], "main-1"]`)},
			fields: fields{
				Type:        ClientTypeInit,
				Data:        []byte(`[1,2]`),
				ComponentID: "main-1",
			},
		},
		{
			name:        "componentID int",
			args:        args{data: []byte(`["init", [1,2], 1]`)},
			wantErr:     true,
			errContains: "unmarshal number",
		},
		{
			name: "componentID and name",
			args: args{data: []byte(`["init", [1,2], "main-1", "main"]`)},
			fields: fields{
				Type:        ClientTypeInit,
				Data:        []byte(`[1,2]`),
				ComponentID: "main-1",
				Name:        "main",
			},
		},
		{
			name: "componentID and name and subcomponentID",
			args: args{data: []byte(`["init", [1,2], "main-1", "main", "0"]`)},
			fields: fields{
				Type:           ClientTypeInit,
				Data:           []byte(`[1,2]`),
				ComponentID:    "main-1",
				Name:           "main",
				SubcomponentID: &zero,
			},
		},
		{
			name:        "componentID and name and subcomponentID as int",
			args:        args{data: []byte(`["init", [1,2], "main-1", "main", 2]`)},
			wantErr:     true,
			errContains: "unmarshal number",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := ClientMessage{}
			if err := json.Unmarshal(tt.args.data, &cm); (err != nil) != tt.wantErr {
				t.Errorf("ClientMessage.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr {
				if tt.errContains == "" || !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ClientMessage.UnmarshalJSON() error = %q, err should contain %q", err, tt.errContains)
				}
				return
			}
			reference := ClientMessage{
				ComponentID:    tt.fields.ComponentID,
				SubcomponentID: tt.fields.SubcomponentID,
				Type:           tt.fields.Type,
				Name:           tt.fields.Name,
				Data:           tt.fields.Data,
			}

			if !reflect.DeepEqual(cm, reference) {
				t.Errorf("Diff() = %s, want %s",
					spew.Sdump(cm),
					spew.Sdump(reference),
				)
			}
		})
	}
}
