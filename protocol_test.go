package texto

import (
	"testing"

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
