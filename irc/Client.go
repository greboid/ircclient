package irc

import (
	"context"
	"fmt"
	"github.com/ergochat/irc-go/ircevent"
	"github.com/ergochat/irc-go/ircmsg"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"strings"
)

const (
	v3TimestampFormat = "2006-01-02T15:04:05.000Z"
)

type ConnectableServer struct {
	Server       string             `yaml:"server" json:"server"`
	TLS          bool               `yaml:"tls" json:"tls"`
	SaslMech     string             `yaml:"saslMech,omitempty" json:"saslMech,omitempty"`
	Saslusername string             `yaml:"saslUsername,omitempty" json:"saslUsername,omitempty"`
	Saslpassword string             `yaml:"saslPassword,omitempty" json:"saslPassword,omitempty"`
	Profile      ConnectableProfile `yaml:"profile" json:"profile"`
}

type ConnectableProfile struct {
	Nick     string `yaml:"nick" json:"nick"`
	User     string `yaml:"user,omitempty" json:"user,omitempty"`
	Realname string `yaml:"realname,omitempty" json:"realname,omitempty"`
}

type Channel struct {
	Name string `json:"name" yaml:"name"`
}

type Message struct {
	Source  string `json:"source" yaml:"source"`
	Target  string `json:"target" yaml:"target"`
	Message string `json:"message" yaml:"message"`
}

type ChannelMessage struct {
	Message
}

type DirectMessage struct {
	Message
}

type ServerMessage struct {
	Message
}

type Client struct {
	Ctx                 context.Context `json:"-" yaml:"-"`
	ircevent.Connection `json:"-" yaml:"-"`
	ConnectableServer   ConnectableServer `json:"-" yaml:"-"`
}

func (c *Client) Connect(server ConnectableServer) error {
	c.ConnectableServer = server
	c.Connection.Server = fmt.Sprintf("%s", server.Server)
	c.Connection.UseTLS = server.TLS
	c.Connection.SASLLogin = server.Saslusername
	c.Connection.SASLPassword = server.Saslpassword
	c.Connection.Nick = server.Profile.Nick
	c.Connection.Debug = true
	c.Connection.RequestCaps = []string{
		"message-tags",
		"echo-message",
		"server-time",
	}
	c.AddListeners()
	return c.Connection.Connect()
}

func (c *Client) AddListeners() {
	c.AddCallback("PRIVMSG", c.handlePrivMsg)
	c.AddCallback("JOIN", c.handleJoin)
}

func (c *Client) handlePrivMsg(message ircmsg.Message) {
	if c.isChannel(message.Params[0]) {
		go runtime.EventsEmit(c.Ctx, "channelMessage", ChannelMessage{Message{
			Source:  message.Source,
			Target:  message.Params[0],
			Message: message.Params[1],
		}})
	} else {
		go runtime.EventsEmit(c.Ctx, "directMessage", DirectMessage{Message{
			Source:  message.Source,
			Target:  message.Params[0],
			Message: message.Params[1],
		}})
	}
}

func (c *Client) handleJoin(message ircmsg.Message) {
	go runtime.EventsEmit(c.Ctx, "channelAdded", &Channel{Name: message.Params[0]})
}

func (c *Client) isChannel(source string) bool {
	chanTypes := c.ISupport()["CHANTYPES"]
	if chanTypes == "" {
		chanTypes = "#"
	}
	for i := 0; i < len(chanTypes); i++ {
		if strings.HasPrefix(source, chanTypes[i:i+1]) {
			return true
		}
	}
	return false
}

func (c *Client) Quit() {
	disconnected := make(chan bool)
	c.AddDisconnectCallback(func(message ircmsg.Message) {
		disconnected <- true
	})
	go c.Connection.Quit()
	<-disconnected
}
