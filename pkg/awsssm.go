package pkg

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type AWSSSM struct {
	//awsSession aws.Config
	SSM *ssm.Client
}

func NewSSM(profile string, region string) *AWSSSM {
	// initialize aws session using config files
	awsSession, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile),
		config.WithRegion(region),
	)

	if err != nil {
		panic(fmt.Sprintf("failed loading config, %v", err))
	}

	ssmClient := ssm.NewFromConfig(awsSession)

	return &AWSSSM{
		//awsSession: awsSession,
		SSM: ssmClient,
	}
}
