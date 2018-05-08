package db

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/wavyllama/chat/config"
	"testing"
)

func TestMessages(t *testing.T) {
	SetupDatabaseForTests(t)
	MessagesSetup(t)
	InsertMessageTest(t)
	DeleteMessageTest(t)
}

func updateMessages() {
	query := fmt.Sprintf("UPDATE messages SET message=UNHEX(\"%s\") WHERE message_timestamp=\"2017-02-01 08:20:19.123456\"", hex.EncodeToString(AESEncrypt([]byte("Hello World"), []byte(db.AESKey))))
	updateMessage(query)
	query = fmt.Sprintf("UPDATE messages SET message=UNHEX(\"%s\") WHERE message_timestamp=\"2018-02-14 11:11:11.111111\"", hex.EncodeToString(AESEncrypt([]byte("Hey Sameet, its Alice <3"), []byte(db.AESKey))))
	updateMessage(query)
	query = fmt.Sprintf("UPDATE messages SET message=UNHEX(\"%s\") WHERE message_timestamp=\"2018-04-10 12:30:08.222222\"", hex.EncodeToString(AESEncrypt([]byte("Hey Andrew, I need help with 511, when are you free?"), []byte(db.AESKey))))
	updateMessage(query)
	query = fmt.Sprintf("UPDATE messages SET message=UNHEX(\"%s\") WHERE message_timestamp=\"2018-03-28 18:04:10.333333\"", hex.EncodeToString(AESEncrypt([]byte("lul"), []byte(db.AESKey))))
	updateMessage(query)
	query = fmt.Sprintf("UPDATE messages SET message=UNHEX(\"%s\") WHERE message_timestamp=\"2018-04-08 17:01:40.444444\"", hex.EncodeToString(AESEncrypt([]byte("I almost made my Mac a brick"), []byte(db.AESKey))))
	updateMessage(query)
	query = fmt.Sprintf("UPDATE messages SET message=UNHEX(\"%s\") WHERE message_timestamp=\"2018-04-12 07:56:00.555555\"", hex.EncodeToString(AESEncrypt([]byte("Why did the chicken cross the road?"), []byte(db.AESKey))))
	updateMessage(query)
	query = fmt.Sprintf("UPDATE messages SET message=UNHEX(\"%s\") WHERE message_timestamp=\"2018-04-12 07:59:13.666666\"", hex.EncodeToString(AESEncrypt([]byte("To get to the other side?"), []byte(db.AESKey))))
	updateMessage(query)
	query = fmt.Sprintf("UPDATE messages SET message=UNHEX(\"%s\") WHERE message_timestamp=\"2018-04-08 17:59:02.777777\"", hex.EncodeToString(AESEncrypt([]byte("When are we playing Fortnite?"), []byte(db.AESKey))))
	updateMessage(query)
}

func updateMessage(query string) {
	_, err := DB.Exec(query)
	if err != nil {
		fmt.Printf("Error with updating messages in tests: %s\n", err)
	}
}

func MessagesSetup(t *testing.T) {
	updateMessages()
	messages := QueryMessages()
	assert.Equal(t, 8, len(messages))
	assert.Equal(t, uint64(52), messages[3].SSID)
	assert.Equal(t, 0, messages[3].SentOrReceived)
	assert.Equal(t, []byte("lul"), messages[3].Text)
}

func InsertMessageTest(t *testing.T) {
	assert.True(t, InsertMessage(52, []byte("wassup"), "2018-04-12 05:01:10.888888", Received))
	messages := QueryMessages()
	assert.Equal(t, 9, len(messages))
	assert.Equal(t, []byte("wassup"), messages[8].Text)
}

func DeleteMessageTest(t *testing.T) {
	assert.True(t, DeleteMessage(52, []byte("wassup"), "2018-04-12 05:01:10.888888", Received))
	messages := QueryMessages()
	assert.Equal(t, 8, len(messages))
}
