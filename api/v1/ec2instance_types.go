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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EC2InstanceSpec defines the desired state of EC2Instance
type EC2InstanceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// The following markers will use OpenAPI v3 schema to validate the value
	// More info: https://book.kubebuilder.io/reference/markers/crd-validation.html

	// foo is an example field of EC2Instance. Edit ec2instance_types.go to remove/update
	// +optional
	AmiId              string            `json:"amiId"`
	InstanceType       string            `json:"instanceType"`
	KeyPair            string            `json:"keyPair,omitempty"`
	SecurityGroupIds   []string          `json:"securityGroupIds,omitempty"`
	Storage            []Storage         `json:"storage"`
	Name               string            `json:"name,omitempty"`
	Tags               map[string]string `json:"tags,omitempty"`
	SubnetId           string            `json:"subnetId,omitempty"`
	AdditionalStorage  []Storage         `json:"additionalStorage,omitempty"`
	Region             string            `json:"region"`
	AssociatedPublicIP bool              `json:"associatedPublicIP,omitempty"`
	UserData           string            `json:"userData,omitempty"`
}

// Storage defines the storage configuration for the EC2Instance.
type Storage struct {
	VolumeType string `json:"volumeType"`
	VolumeSize int    `json:"volumeSize"`
	DeviceName string `json:"deviceName"`
	Encrypted  bool   `json:"encrypted"`
}

// EC2InstanceStatus defines the observed state of EC2Instance.
type EC2InstanceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	State      string `json:"state,omitempty"`
	PublicIP   string `json:"publicIP,omitempty"`
	InstanceID string `json:"instanceID,omitempty"`
	PrivateIP  string `json:"privateIP,omitempty"`
	PublicDNS  string `json:"publicDNS,omitempty"`
	PrivateDNS string `json:"privateDNS,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// EC2Instance is the Schema for the ec2instances API
type EC2Instance struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of EC2Instance
	// +required
	Spec EC2InstanceSpec `json:"spec"`

	// status defines the observed state of EC2Instance
	// +optional
	Status EC2InstanceStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// EC2InstanceList contains a list of EC2Instance
type EC2InstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EC2Instance `json:"items"`
}

type CreatedInstanceInfo struct {
	InstanceID string `json:"instanceID"`
	PublicIP   string `json:"publicIP"`
	PrivateIP  string `json:"privateIP"`
	PublicDNS  string `json:"publicDNS"`
	PrivateDNS string `json:"privateDNS"`
	State      string `json:"state"`
}

func init() {
	SchemeBuilder.Register(&EC2Instance{}, &EC2InstanceList{})
}
