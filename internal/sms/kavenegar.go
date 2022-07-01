package absms

import (
	"comonitoring/pkg/logger"
	"github.com/kavenegar/kavenegar-go"
)

// Kavenegar is an implementation for SMS. This uses
// the Kavenegar APIs to send SMS to the users.
type Kavenegar struct {
	api *kavenegar.Kavenegar
	key string
}

var log, _ = logger.NewLogger(logger.Config{})

// NewKavenegar creates an instance of SMS connected to the
// Kavenegar APIs.
func NewKavenegar(apiKey string) SMS {
	return &Kavenegar{
		api: kavenegar.New(apiKey),
		key: apiKey,
	}
}

// SendMessage uses the Kavenegar APIs to send a new message
// to the user.
func (s *Kavenegar) SendMessage(sender, phoneNumber, content string) error {
	msg, err := s.api.Message.Send(
		sender,
		[]string{phoneNumber},
		content,
		nil,
	)
	log.Infof("%+v", msg)

	return err
}
