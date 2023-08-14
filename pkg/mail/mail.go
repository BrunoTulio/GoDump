package mail

type (
	Message struct {
		Recipient   []string
		Subject     string
		Body        string
		Attachments map[string][]byte
	}

	Mail interface {
		Send(message Message) error
	}
)
