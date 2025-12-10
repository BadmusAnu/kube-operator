package controller

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/log"

	computev1 "github.com/BadmusAnu/kube-operator/api/v1"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func InstanceCheck(ctx context.Context, ec2Instance *computev1.EC2Instance) (string, error) {

	l := log.FromContext(ctx)

	checkInstanceInput := &ec2.DescribeInstancesInput{
		InstanceIds: []string{ec2Instance.Status.InstanceID},
	}

	checkInstance, err := ec2Client.DescribeInstances(ctx, checkInstanceInput)

	if err != nil {
		l.Error(err, "Failed to describe EC2 instance", err)
		return "", fmt.Errorf("failed to describe EC2 instance", err)
	}

	if len(checkInstance.Reservations) == 0 {
		l.Info("EC2 instance not found", "instanceID", ec2Instance.Status.InstanceID)
		return "", fmt.Errorf("EC2 instance not found")
	}

	instance := checkInstance.Reservations[0].Instances[0]

	currentState := string(instance.State.Name)

	return currentState, nil
}

func restartEc2Instance(ctx context.Context, ec2Instance *computev1.EC2Instance) error {
	l := log.FromContext(ctx)

	restartInstanceInput := &ec2.RebootInstancesInput{
		InstanceIds: []string{ec2Instance.Status.InstanceID},
	}

	_, err := ec2Client.RebootInstances(ctx, restartInstanceInput)
	if err != nil {
		l.Error(err, "Failed to restart EC2 instance", err)
		return fmt.Errorf("failed to restart EC2 instance", err)
	}

	return nil
}
