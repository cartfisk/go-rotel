package ha

import (
	"fmt"
	"strings"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type HAStatus uint

type HA struct {
	HAStatus

	topic      string               // Base topic
	components map[string]Component // home assistant components
	commands   map[string]Component // home assistant command topics
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	HAStatusUnknown HAStatus = iota
	HAStatusOnline
	HAStatusOffline
)

const (
	HAStatusOnlineStr  = "online"
	HAStatusOfflineStr = "offline"
)

const (
	topicSeparator = "/"
	topicStatus    = "status"
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func New(topic string) (*HA, error) {
	self := new(HA)

	// Set topic
	if topic == "" {
		return nil, ErrBadParameter.Withf("invalid topic")
	} else {
		self.topic = topic
	}

	// Create a map of configured components
	self.components = make(map[string]Component)
	self.commands = make(map[string]Component)

	// Return success
	return self, nil
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (self *HA) String() string {
	str := "<ha"
	if self.topic != "" {
		str += fmt.Sprintf(" topic=%q", self.topic)
	}
	switch self.HAStatus {
	case HAStatusOnline:
		str += fmt.Sprintf(" status=%q", HAStatusOnlineStr)
	case HAStatusOffline:
		str += fmt.Sprintf(" status=%q", HAStatusOfflineStr)
	default:
		str += fmt.Sprintf(" status=%q", "unknown")
	}
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (self *HA) TopicStatus() string {
	return topicName(self.topic, topicStatus)
}

func (self *HA) SetStatus(v string) error {
	switch v {
	case HAStatusOnlineStr:
		self.HAStatus = HAStatusOnline
	case HAStatusOfflineStr:
		self.HAStatus = HAStatusOffline
	default:
		return ErrBadParameter.Withf("invalid home assistant status %q", v)
	}

	// Return success
	return nil
}

func (self *HA) AddPowerButton(id, object_id string) (Component, error) {
	component, err := NewPowerButton(self.topic, id, object_id)
	if err != nil {
		return nil, err
	}
	if err := self.AddComponent(component); err != nil {
		return nil, err
	}
	return component, nil
}

func (self *HA) AddVolume(id, object_id string) (Component, error) {
	component, err := NewVolume(self.topic, id, object_id)
	if err != nil {
		return nil, err
	}
	if err := self.AddComponent(component); err != nil {
		return nil, err
	}
	return component, nil
}

func (self *HA) AddInput(id, object_id string, options []string) (Component, error) {
	component, err := NewInput(self.topic, id, object_id, options)
	if err != nil {
		return nil, err
	}
	if err := self.AddComponent(component); err != nil {
		return nil, err
	}
	return component, nil
}

func (self *HA) AddComponent(component Component) error {
	key := component.Id()
	if _, exists := self.components[key]; exists {
		return ErrDuplicateEntry.Withf("key %q", key)
	}

	// Save component by unique_id
	self.components[key] = component

	// Save component by command topic
	if key := component.CommandTopic(); key != "" {
		self.commands[key] = component
	}

	// Return success
	return nil
}

func (self *HA) Command(topic string, data []byte) error {
	component, exists := self.commands[topic]
	if !exists {
		return ErrNotFound.Withf("command topic %q", topic)
	}
	if component.SetState(string(data)) {
		fmt.Println("State has changed:", component, string(data))
	}
	return nil
}

func (self *HA) Components() []Component {
	components := make([]Component, 0, len(self.components))
	for _, component := range self.components {
		components = append(components, component)
	}
	return components
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func topicName(topic string, parts ...string) string {
	return strings.Join(append([]string{topic}, parts...), topicSeparator)
}
