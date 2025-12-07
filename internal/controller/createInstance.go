package controller

import (
	"context"
	"fmt"
	"time"

	computev1 "github.com/BadmusAnu/kube-operator"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func CreateInstance(instanceSpec *computev1.Ec2Instance) (InstanceOutput *computev1.CreatedInstanceInfo, err error) {

	l := log.Log.WithName("createInstance")
	l.Info("=== STARTING EC2 INSTANCE CREATION ===")

	RunInstancesInput := &ec2.RunInstancesInput{
		ImageId:      aws.String(instanceSpec.Spec.AMIId),
		InstanceType: ec2types.InstanceType(instanceSpec.Spec.InstanceType),
		KeyName:      aws.String(instanceSpec.Spec.KeyPair),
		SubnetId:     aws.String(instanceSpec.Spec.Subnet),
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
		UserData:     aws.String(instanceSpec.Spec.UserData),
	}

	l.Info("=== CALLING AWS API to create EC2 instance===")

	result, err := ec2Client.RunInstances(context.TODO(), RunInstancesInput)
	if err != nil {
		l.Error(err, "Failed to create EC2 instance")
		return nil, fmt.Errorf("failed to create EC2 instance: %w", err)
	}

	if len(result.Instances) == 0 {
		l.Error(nil, "No instances returned in RunInstancesOutput")

		return nil, fmt.Errorf("no instances returned in RunInstancesOutput")

	}

	inst := result.Instances[0]
	l.Info("=== EC2 INSTANCE CREATED SUCCESSFULLY ===", "instanceID", *inst.InstanceId)

	l.Info("=== WAITING FOR INSTANCE TO BE RUNNING ===")

	runWaiter := ec2.NewInstanceRunningWaiter(ec2Client)
	maxWaitTime := 3 * time.Minute // Increased from 10 seconds - instances typically take 30-60 seconds

	err = runWaiter.Wait(context.TODO(), &ec2.DescribeInstancesInput{
		InstanceIds: []string{*inst.InstanceId},
	}, maxWaitTime)
	if err != nil {
		l.Error(err, "Failed to wait for instance to be running")
		return nil, fmt.Errorf("failed to wait for instance to be running: %w", err)
	}

	l.Info("=== CALLING AWS DescribeInstances API TO GET INSTANCE DETAILS ===")
	describeInput := &ec2.DescribeInstancesInput{
		InstanceIds: []string{*inst.InstanceId},
	}

	describeResult, err := ec2Client.DescribeInstances(context.TODO(), describeInput)
	if err != nil {
		l.Error(err, "Failed to describe EC2 instance")
		return nil, fmt.Errorf("failed to describe EC2 instance: %w", err)
	}

	fmt.Println("Describe result", "public ip", *describeResult.Reservations[0].Instances[0].PublicDnsName, "state", describeResult.Reservations[0].Instances[0].State.Name)

	l.Info("=== INSTANCE IS RUNNING ===")

	instance := describeResult.Reservations[0].Instances[0]
	InstanceOutput = &computev1.CreatedInstanceInfo{
		InstanceID: *inst.InstanceId,
		State:      string(instance.State.Name),
		PublicIP:   derefString(instance.PublicIpAddress),
		PrivateIP:  derefString(instance.PrivateIpAddress),
		PublicDNS:  derefString(instance.PublicDnsName),
		PrivateDNS: derefString(instance.PrivateDnsName),
	}

	return InstanceOutput, nil

}

func derefString(s *string) string {
	if s != nil {
		return *s
	}
	return "<nil>"
}
