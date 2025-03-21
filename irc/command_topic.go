package irc

import (
	"strings"
)

type ChangeTopic struct{}

func (c ChangeTopic) GetName() string {
	return "topic"
}

func (c ChangeTopic) GetHelp() string {
	return "Changes the topic"
}

func (c ChangeTopic) Execute(cm *ConnectionManager, server *Connection, channel *Channel, input string) error {
	if server == nil {
		return NoServerError
	}
	split := strings.SplitN(input, " ", 2)
	if len(split) != 2 {
		return NoChannelError
	}
	return server.connection.Send("TOPIC", split[0], split[1])
}
