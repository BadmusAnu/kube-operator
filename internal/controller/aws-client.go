package controller

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

var ec2Client *ec2.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("configuration error, " + err.Error())
		os.Exit(1)
	}

	ec2Client = ec2.NewFromConfig(cfg)
}
