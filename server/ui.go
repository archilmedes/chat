package server

import (
	"fmt"
	"github.com/marcusolsson/tui-go"
	"github.com/wavyllama/chat/db"
	"log"
	"strconv"
	"strings"
)

const (
	exit        = ":exit"
	friend      = ":friend"
	displayName = "displayName"
	unfriend    = ":delete"
	chat        = ":chat"
)

var activeFriend = "me"

type UI struct {
	Program       *Server
	UI            tui.UI
	Chat, History *tui.Box
	Input         *tui.Entry
	List          *tui.List
}

func handleSpecialString(ui *UI, words []string) {
	switch words[0] {
	case exit:
		ui.UI.Quit()
	case friend:
		if len(words) == 2 {
			friendInfo := strings.Split(words[1], "@")
			if len(friendInfo) == 2 {
				err := ui.Program.SendFriendRequest(friendInfo[1], friendInfo[0])
				if err != nil {
					log.Printf("Error sending friend request: %s\n", err)
				}
				return
			}
		}
		fmt.Printf("Format to add friend: '%s username@ipaddr'\n", friend)
	case unfriend:
		if len(words) == 2 {
			if ui.Program.User.DeleteFriend(words[1]) {
				curr := ui.List.SelectedItem()
				if curr == words[1] {
					curr = db.Self
				}
				ui.List.RemoveItems()
				friends := ui.Program.User.GetFriends()
				for i, f := range friends {
					ui.List.AddItems(f.DisplayName)
					if f.DisplayName == curr {
						ui.List.Select(i)
					}
				}
			} else {
				fmt.Printf("Error deleting friend: %s\n", words[1])
			}
		} else {
			fmt.Printf("Format to delete a friend: '%s displayName\n", unfriend)
		}
	case chat:
		if len(words) == 2 {
			nextFriend, err := strconv.Atoi(words[1])
			if err != nil || nextFriend >= ui.List.Length() {
				goto Failure
			}
			ui.List.Select(nextFriend)
			activeFriend = ui.List.SelectedItem()
			ui.Program.StartOTRSession(activeFriend)
			return
		}
	Failure:
		fmt.Printf("Format for commands: '%s %s'\n", chat, displayName)
	default:
		if len(words) == 1 {
			ui.Program.AcceptedFriend(words[0][1:])
		} else {
			fmt.Printf("Format to accept friend request ':%s\n", displayName)
		}
	}
}

func setInputReader(ui *UI) {
	ui.Input.OnSubmit(func(e *tui.Entry) {
		message := e.Text()
		words := strings.Fields(message)
		if len(message) == 0 {
			return
		}
		if message[0] == ':' {
			handleSpecialString(ui, words)
		} else {
			if activeFriend != "" {
				err := ui.Program.SendChatMessage(activeFriend, message)
				if err != nil {
					log.Printf("Error sending message: %s\n", err)
				}
			} else {
				fmt.Printf("Please set active friend: '%s %s'\n", chat, displayName)
			}
		}
		ui.Input.SetText("")
	})
}

func setPersonChange(ui *UI) {
	ui.List.OnSelectionChanged(func(list *tui.List) {
		for i := 0; i < ui.History.Length(); i++ {
			ui.History.Remove(i)
		}
		conversations := ui.Program.User.GetConversationHistory(list.SelectedItem())
		for _, conv := range conversations {
			chatMessage := new(ReceiveChat)
			var sender string
			if conv.Message.SentOrReceived == db.Sent {
				sender = db.Self
			} else {
				if list.SelectedItem() == db.Self {
					continue
				}
				sender = list.SelectedItem()
			}
			chatMessage.New(string(conv.Message.Text), sender, conv.Message.Timestamp)
			DisplayChatMessage(ui, chatMessage)
		}
	})
}

func NewUI(program *Server) (*UI, error) {
	var ui = new(UI)
	program.UI = ui
	ui.Program = program
	friends := program.User.GetFriends()
	ui.List = tui.NewList()
	for i, f := range friends {
		ui.List.AddItems(f.DisplayName)
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
	setInputReader(ui)
	setPersonChange(ui)
	for i, f := range friends {
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
	showInfo := new(InfoMessage)
	showInfo.New(fmt.Sprintf("Listening on: '%s:%d'", program.User.IP, Port))
	DisplayInfoMessage(ui, showInfo)
	if err := ui.UI.Run(); err != nil {
		return nil, err
	}
	return ui, nil
}

func DisplayInfoMessage(ui *UI, m *InfoMessage) {
	ui.History.Append(tui.NewHBox(
		tui.NewLabel(m.Body()),
		tui.NewSpacer(),
	))
}

func DisplayChatMessage(ui *UI, m *ReceiveChat) {
	if strings.ToLower(m.Sender) == strings.ToLower(activeFriend) {
		ui.History.Append(tui.NewHBox(
			tui.NewLabel(m.Time),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", m.Sender))),
			tui.NewLabel(m.Body()),
			tui.NewSpacer(),
		))
	}
}

func DisplayFriendRequest(ui *UI, m *FriendRequest) {
	ui.History.Append(tui.NewHBox(
		tui.NewLabel(fmt.Sprintf("Friend request from %s@%s (':%s' or just ignore)",
			m.Username, m.IP, displayName)),
		tui.NewSpacer(),
	))
}
