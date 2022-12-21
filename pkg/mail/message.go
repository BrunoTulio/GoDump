package mail

type Message struct {
	Recipient   []string
	Subject     string
	Body        string
	From        string
	Attachments map[string][]byte
}

func (m Message) IsAttachments() bool {
	return len(m.Attachments) > 0
}
