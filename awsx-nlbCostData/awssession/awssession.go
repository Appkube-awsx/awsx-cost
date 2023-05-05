package awssession

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

var sess *session.Session

func GetSessionByRegion(region string) *session.Session {
	sess = session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	return sess
}

func GetSessionByCreds(region string, accessKey string, secretKey string,token string) (*session.Session,error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, token),
		Region:      aws.String(region),
	})
	return sess,err
}
