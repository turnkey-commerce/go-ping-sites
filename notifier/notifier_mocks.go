package notifier

func SendEmailMock(recipient string, message string, subject string) error {
	return nil
}

func SendSmsMock(smsNumber string, message string) error {
	return nil
}
