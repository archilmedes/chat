package db

import (
	"fmt"
	"log"
	"database/sql"
	"github.com/wavyllama/chat/config"
	"crypto/aes"
	"crypto/cipher"
	"bytes"
)

const (
	Sent     = 0
	Received = 1
)

// Stores a DB Message
type DBMessage struct {
	SSID           uint64
	SentOrReceived int
	Text           []byte
	Timestamp      string
}

// Inserts a message into the messages table
func InsertMessage(SSID uint64, message []byte, timestamp string, sentOrReceived int) bool {
	if sentOrReceived != Sent && sentOrReceived != Received {
		log.Fatalf("Invalid entry for sent/received - must be 0 or 1. Instead, received a %d", sentOrReceived)
	}
	encryptedMessages := AESEncrypt(string(message), []byte(db.AESKey))
	insertCommand, err := DB.Prepare("INSERT INTO messages VALUES (?, ?, ?, ?)")
	if err != nil {
		fmt.Printf("Error creating messages prepared statement for InsertMessage: %s", err)
	}
	_, err = insertCommand.Exec(SSID, encryptedMessages, timestamp, sentOrReceived)
	if err != nil {
		log.Panicf("Failed to insert into messages: %s", err)
	}
	return true
}

// Removes a message from the messages table
func DeleteMessage(SSID uint64, message string, timestamp string, sentOrReceived int) bool {
	deleteCommand, err := DB.Prepare("DELETE FROM messages WHERE SSID=? AND message=? AND message_timestamp=? AND sent_or_received=?")
	if err != nil {
		fmt.Printf("Error creating users prepared statement for UpdateUserIP: %s", err)
	}
	encrypted_message := AESEncrypt(string(message), []byte(db.AESKey))
	_, err = deleteCommand.Exec(SSID, encrypted_message, timestamp, sentOrReceived)
	if err != nil {
		log.Panicf("Failed to delete message: %s", err)
	}
	return true
}

// Returns all data in the messages table
func QueryMessages() []DBMessage {
	query, err := DB.Prepare("SELECT * FROM messages")
	if err != nil {
		fmt.Printf("Error creating messages prepared statement for QueryMessages: %s", err)
	}
	results, err :=query.Query()
	if err != nil {
		fmt.Printf("Error executing QueryMessages query: %s", err)
	}
	return ExecuteMessagesQuery(results)
}

// Returns all messages for a given session identified by SSID
func getSessionMessages(SSID uint64) []DBMessage {
	query, err := DB.Prepare("SELECT * FROM messages WHERE SSID=?")
	if err != nil {
		fmt.Printf("Error creating messages prepared statement for getSessionMessages: %s", err)
	}
	results, err :=query.Query(SSID)
	if err != nil {
		fmt.Printf("Error executing getSessionMessages query: %s", err)
	}
	return ExecuteMessagesQuery(results)
}

// Executes the specified database command
func ExecuteMessagesQuery(results *sql.Rows) []DBMessage {
	var messages []DBMessage
	msg := DBMessage{}
	for results.Next() {
		err := results.Scan(&msg.SSID, &msg.Text, &msg.Timestamp, &msg.SentOrReceived)
		if err != nil {
			log.Panicf("Failed to parse results from messages: %s", err)
			panic(err)
		}
		msg.Text = []byte(AESDecrypt(msg.Text, []byte(db.AESKey)))
		messages = append(messages, msg)
	}
	err := results.Err()
	if err != nil {
		log.Panicf("Failed to get results from conversations query: %s", err)
	}
	return messages
}

// Parts of AES encryption code from https://gist.github.com/saoin/b306746041b48a8366d0f63507a4e7f3
func AESEncrypt(src string, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("key error1", err)
	}
	if src == "" {
		fmt.Println("plain content empty")
	}
	ecb := cipher.NewCBCEncrypter(block, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	content := []byte(src)
	content = pkcs5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	return crypted
}

func AESDecrypt(crypt []byte, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("key error1", err)
	}
	if len(crypt) == 0 {
		fmt.Println("plain content empty")
	}
	ecb := cipher.NewCBCDecrypter(block, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	decrypted := make([]byte, len(crypt))
	ecb.CryptBlocks(decrypted, crypt)
	return pkcs5Trimming(decrypted)
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}