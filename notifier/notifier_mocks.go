package notifier

import "errors"

// SendEmailMock mocks the normal response of a successful send with no error.
func SendEmailMock(recipient string, message string, subject string) error {
	return nil
}

// SendEmailErrorMock mocks an error response about no response from the server.
func SendEmailErrorMock(recipient string, message string, subject string) error {
	return errors.New("Error - no response from server.")
}

// SendSmsMock mocks the normal response of a successful send with no error.
func SendSmsMock(smsNumber string, message string) error {
	return nil
}

// SendSmsErrorMock mocks an error response about no response from the server.
func SendSmsErrorMock(smsNumber string, message string) error {
	return errors.New("Error - no response from server.")
}
