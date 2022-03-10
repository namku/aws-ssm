package pkg

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type AWSSSM struct {
	session aws.Config
	SSM     *ssm.Client
}

func NewSSM() *AWSSSM {
	// initialize aws session using config files
	awsSession, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		panic(fmt.Sprintf("failed loading config, %v", err))
	}

	ssmClient := ssm.NewFromConfig(awsSession)

	//	session.Handlers.Send.PushFront(func(r *request.Request) {
	//		// Log every request made and its payload
	//		fmt.Printf("Request: %s/%s, Params: %s",
	//			r.ClientInfo.ServiceName, r.Operation, r.Params)
	//	})

	return &AWSSSM{
		session: awsSession,
		SSM:     ssmClient,
	}
}
