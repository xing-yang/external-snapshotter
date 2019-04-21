/*
Copyright 2019 The Kubernetes Authors.

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

package v1alpha1

import (
	core_v1 "k8s.io/api/core/v1"
	storage "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
        // ExecutionHookTemplateResourcePlural is "executionhooktemplates"
        ExecutionHookTemplateResourcePlural = "executionhooktemplates"
        // ExecutionHookResourcePlural is "executionhooks"
        ExecutionHookResourcePlural = "executionhooks"
)

// ExecutionHookSpec defines the desired state of ExecutionHook
type ExecutionHookSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Without PodSelector/ContainerSelector, Snapshot controller will automatically detect
	// pods/containers needed to run the hook from the volume source of the snapshot.
	// If PodSelector/ContainerSelector are specified in ExecutionHookTemplate, snapshot
	// controller will copy them to ExecutionHookSpec. ExecutionHook controller will find all
	// pods and containers based on PodSelector and containerSelector.
	// +optional
	PodContainerNamesList []PodContainerNames `json:"podContainerNamesList,omitempty" protobuf:"bytes,1,rep,name=podContainerNamesList"`

	// PodSelector is a selector which must be true for the hook to run on the pod
	// Hook controller will find all pods based on PodSelector, if specified
	// This can either be copied from the ExecutionHookTemplate by the controller or
	// new PodSelector can be created in the spec to override the one in the template
	// +optional
	PodSelector *metav1.LabelSelector `json:"podSelector,omitempty" protobuf:"bytes,2,opt,name=podSelector"`

	// ContainerSelector is a selector which must be true for the hook to run on the container
	// Hook controller will find all containers in the selected pods based on ContainerSelector,
	// if specified
	// This can either be copied from the ExecutionHookTemplate by the controller or
	// new ContainerSelector can be created in the spec to override the one in the template
	// If ContainerSeletor is not specified, the first container on the selected pod will be
	// chosen by default.
	// +optional
	ContainerSelector map[string]string `json:"containerSelector,omitempty" protobuf:"bytes,3,rep,name=containerSelector"`

	// Policy is not needed in the spec. If ExecuteOnce, only one pod name will be added here.
	// If ExecuteAll, all pod names will be added here.
	// This field is required.
	PreAction *ExecutionHookAction `json:"preAction" protobuf:"bytes,4,opt,name=preAction"`

	// +optional
	PostAction *ExecutionHookAction `json:"postAction,omitempty" protobuf:"bytes,5,opt,name=postAction"`

	// Time between PreAction and PostAction
	// ExecutionHook controller will check this timeout. If API server dies, hook controller will still trigger
	// PostAction command if this field is set, timeout happens, and PreAction is called (there is a
	// PreActionStartTimestamp)
	// +optional
	ExpirationSeconds *int64 `json:"expirationSeconds,omitempty" protobuf:"varint,6,opt,name=expirationSeconds"`
}

type PodContainerNames struct {
	// This field is required
	PodName string `json:"podName" protobuf:"bytes,1,opt,name=podName"`

	// +optional
	ContainerNames []string `json:"containerNames,omitempty" protobuf:"bytes,2,rep,name=containerNames"`
}

type ExecutionHookAction struct {
	// This field is required.
	Action *core_v1.Handler `json:"action" protobuf:"bytes,1,opt,name=action"`

	// +optional
	ActionTimeoutSeconds *int64 `json:"actionTimeoutSeconds,omitempty" protobuf:"bytes,2,opt,name=actionTimeoutSeconds"`
}

// ExecutionHookStatus defines the observed state of ExecutionHook
type ExecutionHookStatus struct {
	// This is a list of ContainerExecutionHookStatus, with each status representing
	// information about how hook is executed in a pod, including pod name, container name,
	// preaction timestamp, preaction succeed, etc.)
	// +optional
	ContainerExecutionHookStatuses []ContainerExecutionHookStatus `json:"containerExecutionHookStatuses,omitempty" protobuf:"bytes,1,rep,name=containerExecutionHookStatuses"`

	// PreAction Summary status
	// +optional
	// Default is nil
	PreActionSucceed *bool `json:"preActionSucceed,omitempty" protobuf:"varint,2,opt,name=preActionSucceed"`

	// PostAction Summary status
	// +optional
	// Default is nil
	PostActionSucceed *bool `json:"postActionSucceed,omitempty" protobuf:"varint,3,opt,name=postActionSucceed"`
}

// ContainerExecutionHookStatus represents the current state of a hook for a specific
// container in a pod
type ContainerExecutionHookStatus struct {
	// This field is required
	PodName string `json:"podName" protobuf:"bytes,1,opt,name=podName"`

	// This field is required
	ContainerName string `json:"containerName" protobuf:"bytes,2,opt,name=containerName"`

	// This represents the status of a PreAction on a container in a pod
	// +optional
	PreActionStatus *ExecutionHookActionStatus `json:"preActionStatus,omitempty" protobuf:"bytes,3,opt,name=preActionStatus"`

	// This represents the status of a PostAction on a container in a pod
	// +optional
	PostActionStatus *ExecutionHookActionStatus `json:"postActionStatus,omitempty" protobuf:"bytes,4,opt,name=postActionStatus"`
}

// The current state of an action on one container in a pod
type ExecutionHookActionStatus struct {
	// If not set, it is nil, indicating Action has not started
	// If set, it means Action has started at the specified time
	// +optional
	ActionTimestamp *metav1.Time `json:"actionTimestamp,omitempty" protobuf:"bytes,1,opt,name=actionTimestamp"`

	// If not set, it is nil, indicating Action is not complete on the specified pod
	// +optional
	ActionSucceed *bool `json:"actionSucceed,omitempty" protobuf:"varint,2,opt,name=actionSucceed"`

	// The last error encountered when running the action
	// +optional
	Error *storage.VolumeError `json:"error,omitempty" protobuf:"bytes,3,opt,name=error"`
}

type HookPolicy string

// These are valid policies of ExecutionHook
const (
	// Execute commands only once on any one selected pod
	HookExecuteOnce HookPolicy = "ExecuteOnce"
	// Execute commands on all selected pods
	HookExecuteAll HookPolicy = "ExecuteAll"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExecutionHook is the Schema for the executionhooks API
// +k8s:openapi-gen=true
type ExecutionHook struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExecutionHookSpec   `json:"spec,omitempty"`
	Status ExecutionHookStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExecutionHookList contains a list of ExecutionHook
type ExecutionHookList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ExecutionHook `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExecutionHookTemplate describes pre and post action commands to run on selected pods based
// on specified policies. If PodSelector/ContainerSelector are not specified, snapshot controller will
// figure out pods/containers based on the volume source of the snapshot. If specified, the selectors
// will be copied to the ExecutionHookSpec and the hook controller will find out what pods/containers
// should run the commands. Note that selectors in ExecutionHookSpec can override those in template.
// It is namespaced.
// ExecutionHookTemplate is the Schema for the executionhooktemplates API.
// +k8s:openapi-gen=true
type ExecutionHookTemplate struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Policy specifies whether to execute once on one selected pod or execute on all selected pods
	// +optional
	Policy *HookPolicy `json:"hookPolicy,omitempty" protobuf:"bytes,2,opt,name=hookPolicy"`

	// This field is required
	PreAction *ExecutionHookAction `json:"preAction" protobuf:"bytes,3,opt,name=preAction"`

	// PostAction is optional
	// +optional
	PostAction *ExecutionHookAction `json:"postAction" protobuf:"bytes,4,opt,name=postAction"`

	// Time between PreAction and PostAction
	// ExecutionHook controller will check this timeout. If API server dies, hook controller will still trigger
	// PostAction command if this field is set, timeout happens, and PreAction is called (there is a
	// PreActionStartTimestamp)
	// +optional
	ExpirationSeconds *int64 `json:"expirationSeconds,omitempty" protobuf:"varint,5,opt,name=expirationSeconds"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExecutionHookTemplateList contains a list of ExecutionHookTemplate
type ExecutionHookTemplateList struct {
        metav1.TypeMeta `json:",inline"`
        metav1.ListMeta `json:"metadata,omitempty"`
        Items           []ExecutionHookTemplate `json:"items"`
}

type ExecutionHookInfo struct {
        // PodSelector is a selector which must be true for the hook to run on the pod
        // PodSelector in ExecutionHook can override the selector in template
        // Hook controller will find all pods based on PodSelector
        // +optional
        PodSelector *metav1.LabelSelector `json:"podSelector,omitempty" protobuf:"bytes,1,opt,name=podSelector"`

        // ContainerSelector is a selector which must be true for the hook to run on the container
        // ContainerSelector in ExecutionHook can override the selector in template
        // Hook controller will find all containers in the selected pods based on ContainerSelector
        // +optional
        ContainerSelector map[string]string `json:"containerSelector,omitempty" protobuf:"bytes,2,rep,name=containerSelector"`

        // Without PodSelector/ContainerSelector, Snapshot controller will automatically detect
        // pods/containers needed to run the hook from the volume source of the snapshot.
        // If PodSelector/ContainerSelector are specified in ExecutionHookTemplate, snapshot
        // controller will copy them to ExecutionHookSpec. ExecutionHook controller will find all
        // pods and containers based on PodSelector and containerSelector.
        // +optional
        PodContainerNamesList []PodContainerNames `json:"podContainerNamesList,omitempty" protobuf:"bytes,3,rep,name=podContainerNamesList"`

        // Name of the ExecutionHookTemplate
        // This is required and must have a valid hook template name; an empty string is not valid
	ExecutionHookTemplateName string `json:"executionHookTemplateName",omitempty" protobuf:"bytes,4,rep,name=executionHookTemplateName"`
}
