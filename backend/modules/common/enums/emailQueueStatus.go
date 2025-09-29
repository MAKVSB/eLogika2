package enums

type EmailQueueStatusEnum string

const (
	EmailQueueStatusPending EmailQueueStatusEnum = "PENDING"
	EmailQueueStatusSending EmailQueueStatusEnum = "SENDING"
	EmailQueueStatusSent    EmailQueueStatusEnum = "SENT"
	EmailQueueStatusFailed  EmailQueueStatusEnum = "FAILED"
)

var EmailQueueStatusEnumAll = []EmailQueueStatusEnum{
	EmailQueueStatusPending,
	EmailQueueStatusSending,
	EmailQueueStatusSent,
	EmailQueueStatusFailed,
}

func (w EmailQueueStatusEnum) TSName() string {
	switch w {
	case EmailQueueStatusPending:
		return "PENDING"
	case EmailQueueStatusSending:
		return "SENDING"
	case EmailQueueStatusSent:
		return "SENT"
	case EmailQueueStatusFailed:
		return "FAILED"
	default:
		return "???"
	}
}
