package mail

type Mail interface {
	Send(message Message) error
}
