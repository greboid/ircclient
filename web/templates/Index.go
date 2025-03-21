package templates

import "github.com/greboid/ircclient/irc"

type Index struct {
	Connections   []*irc.Connection
	ActiveServer  *irc.Connection
	ActiveChannel *irc.Channel
	WindowInfo    string
	Messages      []*irc.Message
	Users         []*irc.User
}
