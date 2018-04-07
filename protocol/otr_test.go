package protocol

import (
	"testing"
	"fmt"
	"bytes"
	"golang.org/x/crypto/otr"
	"encoding/hex"
)

var QueryMessage = "?OTRv2?"

var alicePrivateKeyHex = "000000000080c81c2cb2eb729b7e6fd48e975a932c638b3a9055478583afa46755683e30102447f6da2d8bec9f386bbb5da6403b0040fee8650b6ab2d7f32c55ab017ae9b6aec8c324ab5844784e9a80e194830d548fb7f09a0410df2c4d5c8bc2b3e9ad484e65412be689cf0834694e0839fb2954021521ffdffb8f5c32c14dbf2020b3ce7500000014da4591d58def96de61aea7b04a8405fe1609308d000000808ddd5cb0b9d66956e3dea5a915d9aba9d8a6e7053b74dadb2fc52f9fe4e5bcc487d2305485ed95fed026ad93f06ebb8c9e8baf693b7887132c7ffdd3b0f72f4002ff4ed56583ca7c54458f8c068ca3e8a4dfa309d1dd5d34e2a4b68e6f4338835e5e0fb4317c9e4c7e4806dafda3ef459cd563775a586dd91b1319f72621bf3f00000080b8147e74d8c45e6318c37731b8b33b984a795b3653c2cd1d65cc99efe097cb7eb2fa49569bab5aab6e8a1c261a27d0f7840a5e80b317e6683042b59b6dceca2879c6ffc877a465be690c15e4a42f9a7588e79b10faac11b1ce3741fcef7aba8ce05327a2c16d279ee1b3d77eb783fb10e3356caa25635331e26dd42b8396c4d00000001420bec691fea37ecea58a5c717142f0b804452f57"
var bobPrivateKeyHex = "000000000080a5138eb3d3eb9c1d85716faecadb718f87d31aaed1157671d7fee7e488f95e8e0ba60ad449ec732710a7dec5190f7182af2e2f98312d98497221dff160fd68033dd4f3a33b7c078d0d9f66e26847e76ca7447d4bab35486045090572863d9e4454777f24d6706f63e02548dfec2d0a620af37bbc1d24f884708a212c343b480d00000014e9c58f0ea21a5e4dfd9f44b6a9f7f6a9961a8fa9000000803c4d111aebd62d3c50c2889d420a32cdf1e98b70affcc1fcf44d59cca2eb019f6b774ef88153fb9b9615441a5fe25ea2d11b74ce922ca0232bd81b3c0fcac2a95b20cb6e6c0c5c1ace2e26f65dc43c751af0edbb10d669890e8ab6beea91410b8b2187af1a8347627a06ecea7e0f772c28aae9461301e83884860c9b656c722f0000008065af8625a555ea0e008cd04743671a3cda21162e83af045725db2eb2bb52712708dc0cc1a84c08b3649b88a966974bde27d8612c2861792ec9f08786a246fcadd6d8d3a81a32287745f309238f47618c2bd7612cb8b02d940571e0f30b96420bcd462ff542901b46109b1e5ad6423744448d20a57818a8cbb1647d0fea3b664e0000001440f9f2eb554cb00d45a5826b54bfa419b6980e48"

func setupMockUsers() (*User, *User) {
	alicePrivateKey, _ := hex.DecodeString(alicePrivateKeyHex)
	bobPrivateKey, _ := hex.DecodeString(bobPrivateKeyHex)

	alice, bob := NewSecureUser(), NewSecureUser()

	// Read private keys from above
	alice.proto.(OTRProtocol).converse.PrivateKey = new(otr.PrivateKey)
	bob.proto.(OTRProtocol).converse.PrivateKey = new(otr.PrivateKey)
	alice.proto.(OTRProtocol).converse.PrivateKey.Parse(alicePrivateKey)
	bob.proto.(OTRProtocol).converse.PrivateKey.Parse(bobPrivateKey)

	return alice, bob
}

func performHandshake(t *testing.T, alice *User, bob *User) {
	var alicesMessage, bobsMessage [][]byte
	var out []byte
	var err error
	var aliceChange, bobChange otr.SecurityChange

	alicesMessage = append(alicesMessage, []byte(QueryMessage))

	// Inspired by official OTR tests in Golang here: https://github.com/keybase/go-crypto

	for round := 0; len(alicesMessage) > 0 || len(bobsMessage) > 0; round++ {
		bobsMessage = nil
		for i, msg := range alicesMessage {
			out, _, bobChange, bobsMessage, err = bob.proto.(OTRProtocol).converse.Receive(msg)
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
			out, _, aliceChange, alicesMessage, err = alice.proto.(OTRProtocol).converse.Receive(msg)
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

	if aliceChange != otr.NewKeys {
		t.Errorf("Alice terminated without signaling new keys")
	}
	if bobChange != otr.NewKeys {
		t.Errorf("Bob terminated without signaling new keys")
	}

	if !bytes.Equal(alice.proto.(OTRProtocol).converse.SSID[:], bob.proto.(OTRProtocol).converse.SSID[:]) {
		t.Errorf("Session identifiers don't match. Alice has %x, Bob has %x",
				 alice.proto.(OTRProtocol).converse.SSID[:], bob.proto.(OTRProtocol).converse.SSID[:])
	}

	if !alice.proto.(OTRProtocol).converse.IsEncrypted() {
		t.Error("Alice doesn't believe that the conversation is secure")
	}
	if !bob.proto.(OTRProtocol).converse.IsEncrypted() {
		t.Error("Bob doesn't believe that the conversation is secure")
	}

	cypher, _ := alice.proto.(OTRProtocol).converse.Send([]byte("Test message"))
	fmt.Println(string(cypher[0]))
	message, _ , _, _, _ := bob.proto.(OTRProtocol).converse.Receive(cypher[0])
	fmt.Println(string(message))
}

func TestPerformHandshake(t *testing.T) {
	alice, bob := setupMockUsers()
	performHandshake(t, alice, bob)
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