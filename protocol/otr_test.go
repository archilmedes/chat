package protocol

import (
	"testing"
	"fmt"
	"bytes"
)

var QueryMessage = "?OTRv2?"

func TestPerformHandshake(t *testing.T) {
	var alicesMessage, bobsMessage [][]byte
	var out []byte
	var err error
	alice, bob := NewSecureUser(), NewSecureUser()

	alice.NewSession()
	bob.NewSession()

	aliceProto, bobProto := NewOTRProtocol(), NewOTRProtocol()
	alice.proto = aliceProto
	bob.proto = bobProto

	alicesMessage = append(alicesMessage, []byte(QueryMessage))

	for round := 0; len(alicesMessage) > 0 || len(bobsMessage) > 0; round++ {
		bobsMessage = nil
		for i, msg := range alicesMessage {
			out, _, _, bobsMessage, err = bobProto.converse.Receive(msg)
			if len(out) > 0 {
				t.Errorf("Bob generated output during key exchange, round %d, message %d", round, i)
			}
			if err != nil {
				t.Fatalf("Bob returned an error, round %d, message %d (%x): %s", round, i, msg, err)
			}
			if len(bobsMessage) > 0 && i != len(alicesMessage)-1 {
				t.Errorf("Bob produced output while processing a fragment, round %d, message %d", round, i)
			}
		}

		alicesMessage = nil
		for i, msg := range bobsMessage {
			out, _, _, alicesMessage, err = aliceProto.converse.Receive(msg)
			if len(out) > 0 {
				t.Errorf("Alice generated output during key exchange, round %d, message %d", round, i)
			}
			if err != nil {
				t.Fatalf("Alice returned an error, round %d, message %d (%x): %s", round, i, msg, err)
			}
			if len(alicesMessage) > 0 && i != len(bobsMessage)-1 {
				t.Errorf("Alice produced output while processing a fragment, round %d, message %d", round, i)
			}
		}
	}

	//if aliceChange != NewKeys {
	//	t.Errorf("Alice terminated without signaling new keys")
	//}
	//if bobChange != NewKeys {
	//	t.Errorf("Bob terminated without signaling new keys")
	//}

	if !bytes.Equal(aliceProto.converse.SSID[:], bobProto.converse.SSID[:]) {
		t.Errorf("Session identifiers don't match. Alice has %x, Bob has %x", aliceProto.converse.SSID[:], bobProto.converse.SSID[:])
	}

	if !aliceProto.converse.IsEncrypted() {
		t.Error("Alice doesn't believe that the conversation is secure")
	}
	if !bobProto.converse.IsEncrypted() {
		t.Error("Bob doesn't believe that the conversation is secure")
	}
}

//func roundTrip(t *testing.T) {
//	alice, bob := NewSecureUser(), NewSecureUser()
//
//	alice.NewSession()
//	bob.NewSession()
//
//	alicesMessage, err := alice.Send(message)
//	if err != nil {
//		t.Errorf("Error from Alice sending message: %s", err)
//	}
//
//	if len(alice.oldMACs) != 0 {
//		t.Errorf("Alice has not revealed all MAC keys")
//	}
//
//	for i, msg := range alicesMessage {
//		out, encrypted, _, _, err := bob.Receive(msg)
//
//		if err != nil {
//			t.Errorf("Error generated while processing test message: %s", err.Error())
//		}
//		if len(out) > 0 {
//			if i != len(alicesMessage)-1 {
//				t.Fatal("Bob produced a message while processing a fragment of Alice's")
//			}
//			if !encrypted {
//				t.Errorf("Message was not marked as encrypted")
//			}
//			if !bytes.Equal(out, message) {
//				t.Errorf("Message corrupted: got %x, want %x", out, message)
//			}
//		}
//	}
//
//	switch macKeyCheck {
//	case firstRoundTrip:
//		if len(bob.oldMACs) != 0 {
//			t.Errorf("Bob should not have MAC keys to reveal")
//		}
//	case subsequentRoundTrip:
//		if len(bob.oldMACs) != 40 {
//			t.Errorf("Bob has %d bytes of MAC keys to reveal, but should have 40", len(bob.oldMACs))
//		}
//	}
//
//	bobsMessage, err := bob.Send(message)
//	if err != nil {
//		t.Errorf("Error from Bob sending message: %s", err)
//	}
//
//	if len(bob.oldMACs) != 0 {
//		t.Errorf("Bob has not revealed all MAC keys")
//	}
//
//	for i, msg := range bobsMessage {
//		out, encrypted, _, _, err := alice.Receive(msg)
//
//		if err != nil {
//			t.Errorf("Error generated while processing test message: %s", err.Error())
//		}
//		if len(out) > 0 {
//			if i != len(bobsMessage)-1 {
//				t.Fatal("Alice produced a message while processing a fragment of Bob's")
//			}
//			if !encrypted {
//				t.Errorf("Message was not marked as encrypted")
//			}
//			if !bytes.Equal(out, message) {
//				t.Errorf("Message corrupted: got %x, want %x", out, message)
//			}
//		}
//	}
//
//	switch macKeyCheck {
//	case firstRoundTrip:
//		if len(alice.oldMACs) != 20 {
//			t.Errorf("Alice has %d bytes of MAC keys to reveal, but should have 20", len(alice.oldMACs))
//		}
//	case subsequentRoundTrip:
//		if len(alice.oldMACs) != 40 {
//			t.Errorf("Alice has %d bytes of MAC keys to reveal, but should have 40", len(alice.oldMACs))
//		}
//	}
//}

func TestOTRProtocol_Encrypt(t *testing.T) {
	proto := NewOTRProtocol()
	// TODO mock or something
	cipherText, _ := proto.Encrypt([]byte("message"))
	fmt.Println(cipherText)
}