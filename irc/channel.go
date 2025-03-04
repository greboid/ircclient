package irc

type Channel struct {
	id       string
	name     string
	messages []*Message
	topic    *Topic
}

func (c *Channel) GetID() string {
	return c.id
}

func (c *Channel) GetName() string {
	return c.name
}

func (c *Channel) GetMessages() []*Message {
	var messages []*Message
	for _, message := range c.messages {
		messages = append(messages, message)
	}
	return messages
}

func (c *Channel) SetTopic(topic *Topic) {
	c.topic = topic
}

func (c *Channel) GetTopic() *Topic {
	if c.topic == nil {
		return NewTopic("")
	}
	return c.topic
}
