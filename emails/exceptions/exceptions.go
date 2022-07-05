package exceptions 

import (
	"errors"
)

func EmailSendFailedError() error {
	return errors.New("Failed to send Email Notification.")
}

func EmailInvalidAddress() error {
	return errors.New("Email Invalid Address.")
}

func EmailSMTPFailure() error {
	return errors.New("SMTP Server Error.")
}


