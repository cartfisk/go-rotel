package ha

import (
	"encoding/json"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Switch struct {
	component
	Icon string `json:"icon,omitempty"`
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewSwitch(topic, Id, objectId string, name, icon string) (*Switch, error) {
	self := new(Switch)
	if err := self.Init(topic, "switch", Id, objectId, name, true, true); err != nil {
		return nil, err
	}
	self.Icon = icon

	// Return success
	return self, nil
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (self *Switch) JSON() ([]byte, error) {
	return json.Marshal(self)
}
