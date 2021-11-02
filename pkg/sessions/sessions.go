package sessions

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

// CreateAWSSession returns an AWS SDK session for use by
// any service client.
func CreateAWSSession() (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}

	return sess, nil
}
