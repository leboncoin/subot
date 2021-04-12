package globals

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

// GetType returns the type of a Message
func (m Message) GetType() MessageType {
	return m.Type
}

// JSONData returns the json bytes of the message
func (m Message) JSONData() []byte {
	j, err := json.Marshal(m)
	if err != nil {
		log.Error("Error marshaling message")
	}
	return j
}

// GetType always returns MessageReaction type
func (r Reaction) GetType() MessageType {
	return MessageReaction
}

// JSONData returns the Reaction data in json bytes
func (r Reaction) JSONData() []byte {
	j, err := json.Marshal(r)
	if err != nil {
		log.Error("Error marshaling reaction")
	}
	return j
}

// GetType always returns Thread MessageType
func (r Reply) GetType() MessageType {
	return Thread
}

// JSONData returns the Reply object in json bytes
func (r Reply) JSONData() []byte {
	j, err := json.Marshal(r)
	if err != nil {
		log.Error("Error marshaling reply")
	}
	return j
}

// JSONData returns the Interaction object in json bytes
func (i Interaction) JSONData() []byte {
	j, err := json.Marshal(i)
	if err != nil {
		log.Error("Error marshaling reply")
	}
	return j
}
