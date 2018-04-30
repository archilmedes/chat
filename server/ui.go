package server

import (
	"fmt"
	"github.com/marcusolsson/tui-go"
	"github.com/wavyllama/chat/db"
	"log"
	"strings"
)

const (
	exit        = ":exit"
	friend      = ":friend"
	displayName = "displayName"
	unfriend    = ":delete"
	reject      = ":reject"
	chat        = ":chat"
)

var activeFriend = ""

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
		if len(words) == 2 && ui.Program.User.IsFriendsWith(words[1]) {
			activeFriend = words[1]
			ui.Program.StartOTRSession(activeFriend)
		} else {
			fmt.Printf("Format for commands: '%s %s'\n", chat, displayName)
		}
	default:
		if len(words) == 1 && ui.Program.User.IsFriendsWith(words[0][1:]) {
			activeFriend = words[0][1:]
			ui.Program.StartOTRSession(activeFriend)
		} else {
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
		// messages := ui.Program.User.GetConversationHistory(list.SelectedItem())
		// TODO: WIP
		/*
			for _, message := range messages {
				chatMessage := new(ReceiveChat)
				chatMessage.New(string(message), )
				DisplayChatMessage(ui, )
			}
		*/
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
	root := tui.NewHBox(sidebar, ui.Chat)
	internalUI, err := tui.New(root)
	if err != nil {
		return nil, err
	}
	ui.UI = internalUI
	internalUI.SetKeybinding("Esc", func() { ui.UI.Quit() })
	if err := ui.UI.Run(); err != nil {
		return nil, err
	}
	return ui, nil
}

func addFriend(ui *UI, newFriend string) {
	ui.List.AddItems(newFriend)
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
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", m.Sender))),
			tui.NewLabel(m.Body()),
			tui.NewSpacer(),
		))
	}
}

func DisplayFriendRequest(ui *UI, m *FriendRequest) {
	// TODO
}
