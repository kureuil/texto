package texto

import (
	"github.com/satori/go.uuid"
	"encoding/json"
	"fmt"
)

const (
	// ErrorMessageKind is returned whenever an error prevents the node from processing a request.
	ErrorMessageKind = "error"
	// RegistrationKind is sent by a Client on its first connection to a server.
	RegistrationKind = "registration"
	// ConnectionMessageKind is sent a Client in response to a RegisterRequestKind. The response includes its ClientID.
	ConnectionMessageKind = "connection"
	// SendMessageKind is sent by a Client whenever it wants to transmit a message to another Client.
	SendMessageKind = "send"
	// ReceiveMessageKind is sent by a Server to a Client, when another Client wants to transmit a message to them.
	ReceiveMessageKind = "receive"
	// AcknowledgeMessageKind is sent after a SendMessageKind or a ReceiveMessageKind was properly processed by the node.
	AcknowledgeMessageKind = "ack"
)

// A ChatMessage conforms to the schema that the chat clients and servers use to communicate.
type ChatMessage struct {
	// Each client is assigned a UUID when it connects to a server.
	ClientID uuid.UUID `json:"client_id"`
	// Each message is identified by its own UUID
	ID uuid.UUID `json:"id"`
	// The kind of the message.
	Kind string `json:"kind"`
	// The actual content of the message, if any.
	Data interface{} `json:"data"`
}

// _ChatMessage is a shadow type which sole purpose is to avoid recursion in ChatMessage_UnmarshalJSON.
type _ChatMessage ChatMessage

// MessageUnmarshalJSON unmarshals a JSON description of a ChatMessage into the given instance.
func (m *ChatMessage) UnmarshalJSON(input []byte) error {
	var data json.RawMessage
	tmp := _ChatMessage{
		Data: &data,
	}
	if err := json.Unmarshal(input, &tmp); err != nil {
		return err
	}
	switch tmp.Kind {
	case ErrorMessageKind:
		var payload ErrorMessagePayload
		if err := json.Unmarshal(data, &payload); err != nil {
			return err
		}
		tmp.Data = payload
	case ConnectionMessageKind:
		var payload ConnectionMessagePayload
		if err := json.Unmarshal(data, &payload); err != nil {
			return err
		}
		tmp.Data = payload
	case SendMessageKind:
		var payload SendMessagePayload
		if err := json.Unmarshal(data, &payload); err != nil {
			return err
		}
		tmp.Data = payload
	case ReceiveMessageKind:
		var payload ReceiveMessagePayload
		if err := json.Unmarshal(data, &payload); err != nil {
			return err
		}
		tmp.Data = payload
	case RegistrationKind:
	case AcknowledgeMessageKind:
	default:
		return fmt.Errorf("Unknown message kind: %s", tmp.Kind)
	}
	*m = ChatMessage(tmp)
	return nil
}

// An ErrorMessagePayload contains the code and the human-readable description of an error.
type ErrorMessagePayload struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// A ConnectionMessagePayload contains the id of a newly registered client.
type ConnectionMessagePayload struct {
	ClientID uuid.UUID `json:"client_id"`
}

// A SendMessagePayload contains the receiver's ID and the content of the message.
type SendMessagePayload struct {
	ReceiverID uuid.UUID `json:"receiver_id"`
	Text       string `json:"text"`
}

// A ReceiveMessagePayload contains the sender's ID and the content of the message.
type ReceiveMessagePayload struct {
	SenderID uuid.UUID `json:"sender_id"`
	Text     string `json:"text"`
}

// NewErrorMessage creates a new ChatMessage of kind "error", with an ErrorMessagePayload.
func NewErrorMessage(messageID *uuid.UUID, clientID uuid.UUID, payload ErrorMessagePayload) *ChatMessage {
	var mID uuid.UUID
	if messageID == nil {
		mID = uuid.NewV4()
	} else {
		mID = *messageID
	}
	return &ChatMessage{
		ID:       mID,
		ClientID: clientID,
		Kind:     ErrorMessageKind,
		Data:     payload,
	}
}

// NewRegistrationMessage creates a new ChatMessage of kind "registration".
func NewRegistrationMessage(messageID *uuid.UUID, clientID uuid.UUID) *ChatMessage {
	var mID uuid.UUID
	if messageID == nil {
		mID = uuid.NewV4()
	} else {
		mID = *messageID
	}
	return &ChatMessage{
		ID:       mID,
		ClientID: clientID,
		Kind:     RegistrationKind,
		Data:     nil,
	}
}


// NewConnectionMessage creates a new ChatMessage of kind "connection", with a ConnectionMessagePayload.
func NewConnectionMessage(messageID *uuid.UUID, clientID uuid.UUID, payload ConnectionMessagePayload) *ChatMessage {
	var mID uuid.UUID
	if messageID == nil {
		mID = uuid.NewV4()
	} else {
		mID = *messageID
	}
	return &ChatMessage{
		ID:       mID,
		ClientID: clientID,
		Kind:     ConnectionMessageKind,
		Data:     payload,
	}
}

// NewAckMessage creates a new ChatMessage of kind "acknowledge".
func NewAckMessage(messageID *uuid.UUID, clientID uuid.UUID) *ChatMessage {
	var mID uuid.UUID
	if messageID == nil {
		mID = uuid.NewV4()
	} else {
		mID = *messageID
	}
	return &ChatMessage{
		ID:       mID,
		ClientID: clientID,
		Kind:     AcknowledgeMessageKind,
		Data:     nil,
	}
}


// NewSendMessage creates a new ChatMessage of kind "send", with a SendMessagePayload.
func NewSendMessage(messageID *uuid.UUID, clientID uuid.UUID, payload SendMessagePayload) *ChatMessage {
	var mID uuid.UUID
	if messageID == nil {
		mID = uuid.NewV4()
	} else {
		mID = *messageID
	}
	return &ChatMessage{
		ID:       mID,
		ClientID: clientID,
		Kind:     SendMessageKind,
		Data:     payload,
	}
}
