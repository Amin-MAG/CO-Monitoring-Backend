package absms

// SMS is a wrapper for the SMS provider
type SMS interface {
	SendMessage(sender, phoneNumber, content string) error
}

// ShortMessage is a simple message that an SMS
// provider sends to the users.
type ShortMessage struct {
	// Sender
	Sender string
	// Receptor a string array of phone numbers that we want to send a message to them.
	Receptor []string
	// Content the text message we like to send to our users.
	Content string
}
