package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/wavyllama/chat/core"
	"github.com/wavyllama/chat/db"
	"testing"
	"time"
)

const (
	fakeMessage = "Hello world"
)

func onReceiveFriendRequest(m *FriendMessage) {
}

func onAcceptFriend(displayName string) {

}

func onReceiveChatMessage(message []byte, friend *db.Friend, time time.Time) {

}

func onProtocolFinish(messageToDisplay string) {

}

func startUpServer(t *testing.T) Server {
	var server Server
	server.InitUIHandlers(onReceiveFriendRequest, onAcceptFriend, onReceiveChatMessage, onProtocolFinish)
	mac, ip, _ := core.GetAddresses()
	db.SetupEmptyTestDatabase()

	user := new(db.User)
	user.Username = "archillin"
	user.IP = ip
	user.MAC = mac
	assert.NoError(t, server.Start(user))
	// Let time pass for handshake to complete
	time.Sleep(2000 * time.Millisecond)
	return server
}

func TestServer_Start(t *testing.T) {
	server := startUpServer(t)
	assert.NotEqual(t, nil, server.Listener)
	assert.NoError(t, server.Shutdown())
}

// TODO uncomment when websocket server tests are fixed
//func TestServer_GetSessionsWithFriend(t *testing.T) {
//	server := startUpServer(t)
//	sessions := server.GetSessionsWithFriend(server.User.MAC, server.User.Username)
//	assert.Equal(t, 2, len(sessions))
//
//	msg := []byte("Hello world")
//	cyp, _ := sessions[0].Proto.Encrypt(msg)
//
//	msgBack, _ := sessions[1].Proto.Decrypt(cyp[0], onProtocolFinish)
//	assert.Equal(t, msgBack[0], msg)
//	assert.NoError(t, server.Shutdown())
//}

func sendAFakeMessage(server Server) {
	sessions := server.GetSessionsWithFriend(server.User.MAC, server.User.Username)

	user1Proto := sessions[0].Proto
	db.InsertMessage(user1Proto.GetSessionID(), []byte(fakeMessage), core.GetFormattedTime(time.Now()), db.Sent)
}

//func TestUser_GetConversationHistory(t *testing.T) {
//	server := startUpServer(t)
//
//	sessions := server.GetSessionsWithFriend(server.User.MAC, server.User.Username)
//	assert.Equal(t, 2, len(sessions))
//	sessions[0].Save()
//
//	sendAFakeMessage(server)
//
//	conversations := server.User.GetConversationHistory(db.Self)
//	assert.Equal(t, 1, len(conversations))
//	assert.Equal(t, []byte(fakeMessage), db.AESDecrypt(conversations[0].Message.Text, []byte(config.AESKey)))
//	assert.Equal(t, db.Sent, conversations[0].Message.SentOrReceived)
//	assert.NoError(t, server.Shutdown())
//}
