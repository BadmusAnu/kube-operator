/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	cloudv1 "github.com/BadmusAnu/kube-operator/api/v1"
)

// EC2InstanceReconciler reconciles a EC2Instance object
type EC2InstanceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cloud.badmus.test,resources=ec2instances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cloud.badmus.test,resources=ec2instances/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cloud.badmus.test,resources=ec2instances/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the EC2Instance object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *EC2InstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := logf.FromContext(ctx)

	l.Info("=== Reconciling Loop for EC2Instance Started ===", "Namespace", req.Namespace, "Name", req.Name)

	ec2Instance := &cloudv1.EC2Instance{}

	if err := r.Get(ctx, req.NamespacedName, ec2Instance); err != nil {
		if errors.IsNotFound(err) {
			l.Info("Instance deleted/does not exist. No Reconciliation Required")
			// kuberenetes will not retry - done, wait for next event
			return ctrl.Result{}, nil
		}
		// kuberenetes will retry with backoff
		return ctrl.Result{}, err
	}

	// check if deletionTimestamp is not zero
	if !ec2Instance.DeletionTimestamp.IsZero() {
		l.Info("Has deletionTimestamp, Instance is being deleted")
		_, err := deleteEc2Instance(ctx, ec2Instance)
		if err != nil {
			l.Error(err, "Failed to delete EC2 instance")
			// Kubernetes will retry with backoff
			return ctrl.Result{Requeue: true}, err
		}

		// Remove the finalizer
		controllerutil.RemoveFinalizer(ec2Instance, "ec2instance.compute.cloud.com")
		if err := r.Update(ctx, ec2Instance); err != nil {
			l.Error(err, "Failed to remove finalizer")
			// Kubernetes will retry with backoff
			return ctrl.Result{Requeue: true}, err
		}
		// at this point, the instance state is terminated and the finalizer is removed
		return ctrl.Result{}, nil
	}

	// check if we already have the instance ID in status
	if ec2Instance.Status.InstanceID != "" {
		l.Info("Instance already exists", "Name", ec2Instance.Spec.Name)
		return ctrl.Result{}, nil
	}

	// l.Info("=== DRIFT DETECTION STARTED ===")

	// // if instance exists, but state is not running, leave it as is, but update the status
	// if ec2Instance.Status.InstanceID != "" {
	// 	l.Info("Instance already exists. No new instance will be created", "Name", ec2Instance.Spec.Name)
	// 	return ctrl.Result{}, nil

	// 	// update the instance status

	// }

	l.Info("Instance does not exist. Creating...", "Name", ec2Instance.Spec.Name)

	l.Info("=== ADDING FINALIZER ===")
	ec2Instance.Finalizers = append(ec2Instance.Finalizers, "ec2instance.cloud.badmus.test")
	if err := r.Update(ctx, ec2Instance); err != nil {
		l.Error(err, "Failed to add finalizer")
		return ctrl.Result{
			Requeue: true,
		}, err
	}
	l.Info("=== FINALIZER ADDED ===")

	l.Info("Reconciling EC2Instance", "Name", ec2Instance.Spec.Name)

	// create the instance
	instanceOutput, err := CreateInstance(ec2Instance)
	if err != nil {
		l.Error(err, "Failed to create EC2 instance")
		return ctrl.Result{}, err // no need for requeue, err automatically means retry with backoff
	}

	l.Info("=== INSTANCE STATUS UPDATING ===")
	// update the instance status

	ec2Instance.Status.InstanceID = instanceOutput.InstanceID
	ec2Instance.Status.State = instanceOutput.State
	ec2Instance.Status.PrivateIP = instanceOutput.PrivateIP
	ec2Instance.Status.PublicIP = instanceOutput.PublicIP
	ec2Instance.Status.PublicDNS = instanceOutput.PublicDNS
	ec2Instance.Status.PrivateDNS = instanceOutput.PrivateDNS

	err = r.Status().Update(ctx, ec2Instance)
	if err != nil {
		l.Error(err, "Failed to update instance status")
		return ctrl.Result{}, err
	}
	l.Info("=== INSTANCE STATUS UPDATED ===")
	l.Info("=== EC2 INSTANCE CREATED SUCCESSFULLY ===", "instanceID", instanceOutput.InstanceID)

	// Kubernetes will not retry - done, wait for next event
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager. or {
func (r *EC2InstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cloudv1.EC2Instance{}).
		Named("ec2instance").
		Complete(r)
}
