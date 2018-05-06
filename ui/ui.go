package ui

import (
	"fmt"
	"github.com/marcusolsson/tui-go"
	"github.com/wavyllama/chat/core"
	"github.com/wavyllama/chat/db"
	"github.com/wavyllama/chat/server"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type UI struct {
	Program       *server.Server
	UI            tui.UI
	Chat, History *tui.Box
	Input         *tui.Entry
	List          *tui.List
}

const (
	exit        = ":exit"
	friend      = ":friend"
	displayName = "displayName"
	unfriend    = ":delete"
	deleteSelf  = ":deleteSelf"
	chat        = ":chat"
	errorPrefix = "Error"
	logFile     = "chat-debug.log"
)

var activeFriend = db.Self

// Updates the history in the same thread. Use with caution! Must be wrapped inside a UI.Update, or used in an event handler
func (ui *UI) displayMessage(message string) {
	ui.History.Append(tui.NewHBox(
		tui.NewLabel(message),
		tui.NewSpacer(),
	))
}

// Handles special input for commands
// Since this function is only called when a string that is a command is typed, it does not have to be executed in a
// separate thread to update UI
func (ui *UI) handleSpecialString(words []string) {
	switch words[0] {
	case exit:
		ui.UI.Quit()
	case friend:
		if len(words) == 2 {
			friendInfo := strings.Split(words[1], "@")
			if len(friendInfo) == 2 {
				ui.displayMessage(fmt.Sprintf("Sent friend request to %s", friendInfo))
				username, ip := friendInfo[0], friendInfo[1]
				if !core.IsIPOnline(ip) {
					ui.displayMessage(fmt.Sprintf("Friend is not currently online! Please check that you have typed the address correctly"))
				}
				// Send a friend request in another thread
				go func() {
					err := ui.Program.SendFriendRequest(ip, username)
					if err != nil {
						ui.displayMessage(fmt.Sprintf("Error sending friend request: %s", err))
					}
				}()
				return
			}
		}
		ui.displayMessage(fmt.Sprintf("Format to add friend: '%s username@ipaddr'", friend))
	case unfriend:
		if len(words) != 2 {
			ui.displayMessage(fmt.Sprintf("Format to delete a friend: '%s displayName", unfriend))
			return
		}
		if ui.Program.User.DeleteFriend(words[1]) {
			// TODO fix this logic
			curr := ui.List.SelectedItem()
			if curr == words[1] {
				curr = db.Self
			}
			ui.List.RemoveItem(1)
		} else {
			ui.displayMessage("Error deleting friend")
		}
	case chat:
		if len(words) == 2 {
			nextFriend, err := strconv.Atoi(words[1])
			if err != nil || nextFriend < 0 || nextFriend >= ui.List.Length() {
				ui.displayMessage("You entered a number outside the list of your friends.")
				return
			}
			ui.List.Select(nextFriend)
			activeFriend = ui.List.SelectedItem()
			err = ui.Program.StartOTRSession(activeFriend)
			if err != nil {
				ui.displayMessage(fmt.Sprintf("%s %s\n", errorPrefix, err.Error()))
			}
			return
		}
		ui.displayMessage(fmt.Sprintf("Format for commands: '%s %s'", chat, displayName))
	case deleteSelf:
		if !ui.Program.User.Delete() {
			ui.displayMessage(fmt.Sprintf("%s: Failed to delete your account", errorPrefix))
		}
		ui.displayMessage("Successfully deleted all of your data")
		ui.UI.Quit()
	default:
		log.Printf("Display name is %s", words[0][1:])
		if len(words) == 1 {
			ui.Program.AcceptedFriend(words[0][1:])
			return
		}
		ui.displayMessage(fmt.Sprintf("Format to accept friend request ':%s", displayName))
	}
}

// Handle user input
func (ui *UI) setInputReader() {
	ui.Input.OnSubmit(func(e *tui.Entry) {
		message := e.Text()
		words := strings.Fields(message)
		if len(message) == 0 {
			return
		}
		if message[0] == ':' {
			ui.handleSpecialString(words)
		} else {
			if activeFriend != "" {
				sendMessage := new(ReceiveChat)
				sendMessage.Message = message
				sendMessage.Time = time.Now()
				sendMessage.Sender = db.Self
				if activeFriend != db.Self {
					ui.displayChatMessage(sendMessage)
				}
				err := ui.Program.SendChatMessage(activeFriend, message)
				if err != nil {
					ui.displayMessage(fmt.Sprintf("%s %s\n", errorPrefix, err.Error()))
				}
			} else {
				ui.displayMessage(fmt.Sprintf("Please set active friend: '%s %s'\n", chat, displayName))
			}
		}
		ui.Input.SetText("")
	})
}

func (ui *UI) displayConversationsWithFriend(displayName string) {
	conversations := ui.Program.User.GetConversationHistory(displayName)
	for _, conv := range conversations {
		var sender string // Get who sent the message
		if conv.Message.SentOrReceived == db.Sent {
			sender = db.Self
		} else {
			if conv.Session.FriendDisplayName == db.Self {
				continue
			}
			sender = conv.Session.FriendDisplayName
		}
		ui.displayChatMessage(NewReceiveChatMessage(string(conv.Message.Text), sender, conv.Message.Timestamp))
	}
}

// Show previous messages of friend upon choosing friend to chat with
func (ui *UI) setPersonChange() {
	ui.List.OnSelectionChanged(func(list *tui.List) {
		// clear history
		ui.UI.Update(func() {
			for i := 0; i < ui.History.Length(); i++ {
				ui.History.Remove(i)
			}
		})
		ui.displayConversationsWithFriend(ui.List.SelectedItem())
	})
}

// Help from: https://github.com/marcusolsson/tui-go/tree/master/example/chat
func NewUI(program *server.Server) (*UI, error) {
	var ui = new(UI)

	ui.Program = program
	friends := program.User.GetFriends()
	ui.List = tui.NewList()
	for i, f := range friends {
		friendName := f.DisplayName
		if online, _ := program.User.IsFriendOnline(f.DisplayName); online {
			friendName = fmt.Sprintf("%s (active)\n", friendName)
		}
		ui.List.AddItems(friendName)
		if strings.ToLower(f.DisplayName) == db.Self {
			ui.List.SetSelected(i)
		}
	}
	sidebar := tui.NewVBox(
		tui.NewLabel("FRIENDS"),
		ui.List,
	)
	sidebar.SetBorder(true)
	ui.History = tui.NewVBox()
	historyScroll := tui.NewScrollArea(ui.History)
	historyScroll.SetAutoscrollToBottom(true)
	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)
	ui.Input = tui.NewEntry()
	ui.Input.SetFocused(true)
	ui.Input.SetSizePolicy(tui.Expanding, tui.Maximum)
	inputBox := tui.NewHBox(ui.Input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)
	ui.Chat = tui.NewVBox(historyBox, inputBox)
	ui.Chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	// Put all friends in sidebar
	for i, f := range friends {
		// Select self as current friend
		if strings.ToLower(f.DisplayName) == db.Self {
			ui.List.Select(i)
		}
	}
	root := tui.NewHBox(sidebar, ui.Chat)
	internalUI, err := tui.New(root)
	if err != nil {
		return nil, err
	}
	ui.UI = internalUI
	internalUI.SetKeybinding("Esc", func() { ui.UI.Quit() })

	// Set event handlers
	ui.setInputReader()
	ui.setPersonChange()
	program.InitUIHandlers(ui.onReceiveFriendRequest, ui.onAcceptFriend, ui.onReceiveChatMessage, ui.onInfoMessage)

	f, _ := os.Create(logFile)
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	tui.SetLogger(logger)

	return ui, nil
}

// Runs the UI
func (ui *UI) Run() error {
	go func() {
		ui.Program.Start()
		ui.displayConversationsWithFriend(db.Self)
	}()
	return ui.UI.Run()
}

// Write info message to UI
func (ui *UI) displayInfoMessage(m *InfoMessage) {
	ui.UI.Update(func() {
		ui.displayMessage(m.Body())
	})
}

// Write chat message to UI
func (ui *UI) displayChatMessage(m *ReceiveChat) {
	ui.UI.Update(func() {
		ui.History.Append(tui.NewHBox(
			tui.NewLabel(m.Time.Format("3:04PM")),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", m.Sender))),
			tui.NewLabel(m.Body()),
			tui.NewSpacer(),
		))
	})
}

// Write friend request to UI
func (ui *UI) displayFriendRequest(m *FriendRequest) {
	ui.UI.Update(func() {
		ui.displayMessage(fmt.Sprintf("Friend request from %s@%s (':%s' or just ignore)",
			m.Username, m.IP, displayName))
	})
}

func (ui *UI) onReceiveFriendRequest(m *server.FriendMessage) {
	ui.displayFriendRequest(NewFriendRequestMessage("", ui.Program.LastFriend.Username, ui.Program.LastFriend.IP))
}

func (ui *UI) onAcceptFriend(displayName string) {
	ui.UI.Update(func() {
		ui.List.AddItems(displayName)
	})
}

func (ui *UI) onReceiveChatMessage(message []byte, friend *db.Friend, time time.Time) {
	ui.displayChatMessage(NewReceiveChatMessage(string(message), friend.DisplayName, time))
}

func (ui *UI) onInfoMessage(messageToDisplay string) {
	ui.displayInfoMessage(NewInfoMessage(messageToDisplay))
}

func (ui *UI) onErrorReceived(errorDescription string) {
	ui.onInfoMessage(fmt.Sprintf("%s: %s\n", errorPrefix, errorDescription))
}
