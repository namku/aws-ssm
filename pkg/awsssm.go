package pkg

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go/aws"
)

type AWSSSM struct {
	awsSession aws.Config
	SSM        *ssm.Client
}

func NewSSM(profile string, region string) *ssm.Client {
	// initialize aws session using config files
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile),
		config.WithRegion(region),
	)

	if err != nil {
		panic(fmt.Sprintf("failed loading config, %v", err))
	}

	ssmClient := ssm.NewFromConfig(cfg)
	return ssmClient

	//return &AWSSSM{
	//	//awsSession: cfg,
	//	SSM: ssmClient,
	//}
}
