package texto

import (
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

const invalidChatMessage = `
{
	"client_id": "b50bff94-4f43-4e24-9c71-dea94d3db825",
	"id": "b857e508-3993-46b9-b227-ca7528f2861d",
	"kind": "bleepblopimabot",
	"data": {}
}
`

const errorChatMessage = `
{
	"client_id": "b50bff94-4f43-4e24-9c71-dea94d3db825",
	"id": "b857e508-3993-46b9-b227-ca7528f2861d",
	"kind": "error",
	"data": {
		"code": "ENOMEM",
		"description": "Out-of-memory"
	}
}
`

const registrationChatMessage = `
{
	"client_id": "b50bff94-4f43-4e24-9c71-dea94d3db825",
	"id": "b857e508-3993-46b9-b227-ca7528f2861d",
	"kind": "registration",
	"data": {}
}
`

const connectionChatMessage = `
{
	"client_id": "b50bff94-4f43-4e24-9c71-dea94d3db825",
	"id": "b857e508-3993-46b9-b227-ca7528f2861d",
	"kind": "connection",
	"data": {
		"client_id": "8f718542-0e5a-4d9a-9ce9-eb8ad1912359"
	}
}
`

const sendChatMessage = `
{
	"client_id": "b50bff94-4f43-4e24-9c71-dea94d3db825",
	"id": "b857e508-3993-46b9-b227-ca7528f2861d",
	"kind": "send",
	"data": {
		"receiver_id": "8f718542-0e5a-4d9a-9ce9-eb8ad1912359",
		"text": "Lorem ipsum dolor si amet."
	}
}
`

const receiveChatMessage = `
{
	"client_id": "b50bff94-4f43-4e24-9c71-dea94d3db825",
	"id": "b857e508-3993-46b9-b227-ca7528f2861d",
	"kind": "receive",
	"data": {
		"sender_id": "8f718542-0e5a-4d9a-9ce9-eb8ad1912359",
		"text": "Lorem ipsum dolor si amet."
	}
}
`

const ackChatMessage = `
{
	"client_id": "b50bff94-4f43-4e24-9c71-dea94d3db825",
	"id": "b857e508-3993-46b9-b227-ca7528f2861d",
	"kind": "ack",
	"data": {}
}
`

func TestMessageUnmarshalJSON(t *testing.T) {
	var invalidMsg ChatMessage
	assert.NotNil(t, invalidMsg.UnmarshalJSON([]byte(invalidChatMessage)))

	var errorMsg ChatMessage
	if assert.NoError(t, errorMsg.UnmarshalJSON([]byte(errorChatMessage))) {
		assert.Equal(t, "b857e508-3993-46b9-b227-ca7528f2861d", errorMsg.ID.String())
		assert.Equal(t, "ENOMEM", errorMsg.Data.(ErrorMessagePayload).Code)
		assert.Equal(t, "Out-of-memory", errorMsg.Data.(ErrorMessagePayload).Description)
	}

	var registrationMsg ChatMessage
	if assert.NoError(t, registrationMsg.UnmarshalJSON([]byte(registrationChatMessage))) {
		assert.Equal(t, "b857e508-3993-46b9-b227-ca7528f2861d", registrationMsg.ID.String())
	}

	var connectionMsg ChatMessage
	if assert.NoError(t, connectionMsg.UnmarshalJSON([]byte(connectionChatMessage))) {
		assert.Equal(t, "b857e508-3993-46b9-b227-ca7528f2861d", connectionMsg.ID.String())
		assert.Equal(t, "8f718542-0e5a-4d9a-9ce9-eb8ad1912359", connectionMsg.Data.(ConnectionMessagePayload).ClientID.String())
	}

	var sendMsg ChatMessage
	if assert.NoError(t, sendMsg.UnmarshalJSON([]byte(sendChatMessage))) {
		assert.Equal(t, "b857e508-3993-46b9-b227-ca7528f2861d", sendMsg.ID.String())
		assert.Equal(t, "8f718542-0e5a-4d9a-9ce9-eb8ad1912359", sendMsg.Data.(SendMessagePayload).ReceiverID.String())
		assert.Equal(t, "Lorem ipsum dolor si amet.", sendMsg.Data.(SendMessagePayload).Text)
	}

	var receiveMsg ChatMessage
	if assert.NoError(t, receiveMsg.UnmarshalJSON([]byte(receiveChatMessage))) {
		assert.Equal(t, "b857e508-3993-46b9-b227-ca7528f2861d", receiveMsg.ID.String())
		assert.Equal(t, "8f718542-0e5a-4d9a-9ce9-eb8ad1912359", receiveMsg.Data.(ReceiveMessagePayload).SenderID.String())
		assert.Equal(t, "Lorem ipsum dolor si amet.", receiveMsg.Data.(ReceiveMessagePayload).Text)
	}

	var ackMsg ChatMessage
	if assert.NoError(t, ackMsg.UnmarshalJSON([]byte(ackChatMessage))) {
		assert.Equal(t, "b857e508-3993-46b9-b227-ca7528f2861d", ackMsg.ID.String())
	}
}

func TestNewErrorMessage(t *testing.T) {
	var defaultID uuid.UUID
	clientID := uuid.NewV4()
	errorPayload := ErrorMessagePayload{"ENOMEM", "Out-of-memory"}
	msg := NewErrorMessage(nil, clientID, errorPayload)
	assert.NotEqual(t, defaultID.String(), msg.ID.String())
	assert.Equal(t, clientID.String(), msg.ClientID.String())
	assert.Equal(t, ErrorMessageKind, msg.Kind)
	assert.Equal(t, errorPayload.Code, msg.Data.(ErrorMessagePayload).Code)
	assert.Equal(t, errorPayload.Description, msg.Data.(ErrorMessagePayload).Description)
	answer := NewErrorMessage(&msg.ID, clientID, errorPayload)
	assert.Equal(t, msg.ID, answer.ID)
}

func TestNewAckMessage(t *testing.T) {
	var defaultID uuid.UUID
	clientID := uuid.NewV4()
	msg := NewAckMessage(nil, clientID)
	assert.NotEqual(t, defaultID.String(), msg.ID.String())
	assert.Equal(t, clientID.String(), msg.ClientID.String())
	assert.Equal(t, AcknowledgeMessageKind, msg.Kind)
	answer := NewAckMessage(&msg.ID, clientID)
	assert.Equal(t, msg.ID, answer.ID)
}

func TestNewConnectionMessage(t *testing.T) {
	var defaultID uuid.UUID
	clientID := uuid.NewV4()
	connectionPayload := ConnectionMessagePayload{clientID}
	msg := NewConnectionMessage(nil, clientID, connectionPayload)
	assert.NotEqual(t, defaultID.String(), msg.ID.String())
	assert.Equal(t, clientID.String(), msg.ClientID.String())
	assert.Equal(t, ConnectionMessageKind, msg.Kind)
	assert.Equal(t, connectionPayload.ClientID, msg.Data.(ConnectionMessagePayload).ClientID)
	answer := NewConnectionMessage(&msg.ID, clientID, connectionPayload)
	assert.Equal(t, msg.ID, answer.ID)
}

func TestNewRegistrationMessage(t *testing.T) {
	var defaultID uuid.UUID
	clientID := uuid.NewV4()
	msg := NewRegistrationMessage(nil, clientID)
	assert.NotEqual(t, defaultID.String(), msg.ID.String())
	assert.Equal(t, clientID.String(), msg.ClientID.String())
	assert.Equal(t, RegistrationKind, msg.Kind)
	answer := NewRegistrationMessage(&msg.ID, clientID)
	assert.Equal(t, msg.ID, answer.ID)
}

func TestNewSendMessage(t *testing.T) {
	var defaultID uuid.UUID
	clientID := uuid.NewV4()
	sendPayload := SendMessagePayload{
		ReceiverID: clientID,
		Text: "Hello World!",
	}
	msg := NewSendMessage(nil, clientID, sendPayload)
	assert.NotEqual(t, defaultID.String(), msg.ID.String())
	assert.Equal(t, clientID.String(), msg.ClientID.String())
	assert.Equal(t, SendMessageKind, msg.Kind)
	assert.Equal(t, sendPayload.ReceiverID, msg.Data.(SendMessagePayload).ReceiverID)
	assert.Equal(t, sendPayload.Text, msg.Data.(SendMessagePayload).Text)
	answer := NewRegistrationMessage(&msg.ID, clientID)
	assert.Equal(t, msg.ID, answer.ID)
}
