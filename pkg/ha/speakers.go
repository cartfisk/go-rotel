package ha

import (
	"encoding/json"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Speakers struct {
	component
	Icon    string   `json:"icon,omitempty"`
	Options []string `json:"options,omitempty"`
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewSpeakers(topic, Id, objectId string, options []string) (*Speakers, error) {
	self := new(Speakers)
	if err := self.Init(topic, "select", Id, objectId, "Speakers", true, true); err != nil {
		return nil, err
	}
	self.Icon = "mdi:speaker-multiple"
	self.Options = options

	// Return success
	return self, nil
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (self *Speakers) JSON() ([]byte, error) {
	return json.Marshal(self)
}
