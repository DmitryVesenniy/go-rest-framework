package notifications

type NotificationsInterface interface {
	SendEmail(string) error
	SendSMS(string) error
}
