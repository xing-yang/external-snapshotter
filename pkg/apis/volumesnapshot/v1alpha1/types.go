/*
Copyright 2018 The Kubernetes Authors.

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
	"encoding/json"

	core_v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// VolumeSnapshotContentResourcePlural is "volumesnapshotcontents"
	VolumeSnapshotContentResourcePlural = "volumesnapshotcontents"
	// VolumeSnapshotResourcePlural is "volumesnapshots"
	VolumeSnapshotResourcePlural = "volumesnapshots"
	// SnapshotClassResourcePlural is "snapshotclasss"
	SnapshotClassResourcePlural = "snapshotclasses"
)

type VolumeSnapshotCreated struct {
	// The time the snapshot was successfully created
	// +optional
	CreatedAt metav1.Time `json:"createdAt" protobuf:"bytes,1,opt,name=createdAt"`
}

type VolumeSnapshotAvailable struct {
	// The time the snapshot was successfully created and available for use
	// +optional
	AvailableAt metav1.Time `json:"availableAt" protobuf:"bytes,1,opt,name=availableAt"`
}

type VolumeSnapshotError struct {
	// A brief CamelCase string indicating details about why the snapshot is in error state.
	// +optional
	Reason string
	// A human-readable message indicating details about why the snapshot is in error state.
	// +optional
	Message string
	// The time the error occurred during the snapshot creation (or uploading) process
	// +optional
	FailedAt metav1.Time `json:"failedAt" protobuf:"bytes,1,opt,name=failedAt"`
}

// VolumeSnapshotStatus is the status of the VolumeSnapshot
type VolumeSnapshotStatus struct {
	// VolumeSnapshotCreated indicates whether the snapshot was successfully created.
	// If the timestamp CreateAt is set, it means the snapshot was created;
	// Otherwise the snapshot was not created.
	// +optional
	Created VolumeSnapshotCreated

	// VolumeSnapshotAvailable indicates whether the snapshot was available for use.
	// A snapshot MUST have already been created before it can be available.
	// If a snapshot was available, it indicates the snapshot was created.
	// When the snapshot was created but not available yet, the application can be
	// resumed if it was previously frozen before taking the snapshot. In this case,
	// it is possible that the snapshot is being uploaded to the cloud. For example,
	// both GCE and AWS support uploading of the snapshot after it is cut as part of
	// the Create Snapshot process.
	// If the timestamp AvailableAt is set, it means the snapshot was available;
	// Otherwise the snapshot was not available.
	// +optional
	Available VolumeSnapshotAvailable

	// VolumeSnapshotError indicates an error occurred during the snapshot creation
	// (or uploading) process.
	// +optional
	Error VolumeSnapshotError
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VolumeSnapshot is the volume snapshot object accessible to the user. Upon successful creation of the actual
// snapshot by the volume provider it is bound to the corresponding VolumeSnapshotContent through
// the VolumeSnapshotSpec
type VolumeSnapshot struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec represents the desired state of the snapshot
	Spec VolumeSnapshotSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`

	// Status represents the latest observer state of the snapshot
	// +optional
	Status VolumeSnapshotStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VolumeSnapshotList is a list of VolumeSnapshot objects
type VolumeSnapshotList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is the list of VolumeSnapshots
	Items []VolumeSnapshot `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// VolumeSnapshotSpec is the desired state of the volume snapshot
type VolumeSnapshotSpec struct {
	// PersistentVolumeClaimName is the name of the PVC being snapshotted
	// +optional
	PersistentVolumeClaimName string `json:"persistentVolumeClaimName" protobuf:"bytes,1,opt,name=persistentVolumeClaimName"`

	// SnapshotDataName binds the VolumeSnapshot object with the VolumeSnapshotContent
	// +optional
	SnapshotDataName string `json:"snapshotDataName" protobuf:"bytes,2,opt,name=snapshotDataName"`

	// Name of the SnapshotClass required by the volume snapshot.
	// +optional
	SnapshotClassName string `json:"snapshotClassName" protobuf:"bytes,3,opt,name=snapshotClassName"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SnapshotClass describes the parameters used by storage system when
// provisioning VolumeSnapshots from PVCs.
//
// SnapshotClasses are non-namespaced; the name of the snapshot class
// according to etcd is in ObjectMeta.Name.
type SnapshotClass struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Snapshotter is the driver expected to handle this SnapshotClass.
	Snapshotter string `json:"snapshotter" protobuf:"bytes,2,opt,name=snapshotter"`

	// Parameters holds parameters for the snapshotter.
	// These values are opaque to the system and are passed directly
	// to the snapshotter.
	// +optional
	Parameters map[string]string `json:"parameters,omitempty" protobuf:"bytes,3,rep,name=parameters"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SnapshotClassList is a collection of snapshot classes.
type SnapshotClassList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is the list of SnapshotClasses
	Items []SnapshotClass `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VolumeSnapshotContent represents the actual "on-disk" snapshot object
type VolumeSnapshotContent struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec represents the desired state of the snapshot data
	Spec VolumeSnapshotContentSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VolumeSnapshotContentList is a list of VolumeSnapshotContent objects
type VolumeSnapshotContentList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is the list of VolumeSnapshotContents
	Items []VolumeSnapshotContent `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// VolumeSnapshotContentSpec is the spec of the volume snapshot data
type VolumeSnapshotContentSpec struct {
	// Source represents the location and type of the volume snapshot
	VolumeSnapshotSource `json:",inline" protobuf:"bytes,1,opt,name=volumeSnapshotSource"`

	// VolumeSnapshotRef is part of bi-directional binding between VolumeSnapshot
	// and VolumeSnapshotContent. It becomes non-nil when bound.
	// +optional
	VolumeSnapshotRef *core_v1.ObjectReference `json:"volumeSnapshotRef" protobuf:"bytes,2,opt,name=volumeSnapshotRef"`

	// PersistentVolumeRef represents the PersistentVolume that the snapshot has been
	// taken from. It becomes non-nil when VolumeSnapshot and VolumeSnapshotContent are bound.
	// +optional
	PersistentVolumeRef *core_v1.ObjectReference `json:"persistentVolumeRef" protobuf:"bytes,3,opt,name=persistentVolumeRef"`
}

// VolumeSnapshotSource represents the actual location and type of the snapshot. Only one of its members may be specified.
type VolumeSnapshotSource struct {
	// CSI (Container Storage Interface) represents storage that handled by an external CSI Volume Driver (Alpha feature).
	// +optional
	CSI *CSIVolumeSnapshotSource `json:"csiVolumeSnapshotSource,omitempty"`
}

// Represents the source from CSI volume snapshot
type CSIVolumeSnapshotSource struct {
	// Driver is the name of the driver to use for this snapshot.
	// Required.
	Driver string `json:"driver"`

	// SnapshotHandle is the unique snapshot id returned by the CSI volume
	// plugin’s CreateSnapshot to refer to the snapshot on all subsequent calls.
	// Required.
	SnapshotHandle string `json:"snapshotHandle"`

	// Timestamp when the point-in-time snapshot is taken on the storage
	// system. This timestamp will be generated by the CSI volume driver after
	// the snapshot is cut. The format of this field should be a Unix nanoseconds
	// time encoded as an int64. On Unix, the command `date +%s%N` returns
	// the  current time in nanoseconds since 1970-01-01 00:00:00 UTC.
	// This field is REQUIRED.
	CreatedAt int64 `json:"createdAt,omitempty" protobuf:"varint,3,opt,name=createdAt"`
}

// GetObjectKind is required to satisfy Object interface
func (v *VolumeSnapshotContent) GetObjectKind() schema.ObjectKind {
	return &v.TypeMeta
}

// GetObjectMeta is required to satisfy ObjectMetaAccessor interface
func (v *VolumeSnapshotContent) GetObjectMeta() metav1.Object {
	return &v.ObjectMeta
}

// GetObjectKind is required to satisfy Object interface
func (vd *VolumeSnapshotContentList) GetObjectKind() schema.ObjectKind {
	return &vd.TypeMeta
}

// GetObjectKind is required to satisfy Object interface
func (v *VolumeSnapshot) GetObjectKind() schema.ObjectKind {
	return &v.TypeMeta
}

// GetObjectMeta is required to satisfy ObjectMetaAccessor interface
func (v *VolumeSnapshot) GetObjectMeta() metav1.Object {
	return &v.ObjectMeta
}

// GetObjectKind is required to satisfy Object interface
func (vd *VolumeSnapshotList) GetObjectKind() schema.ObjectKind {
	return &vd.TypeMeta
}

// VolumeSnapshotContentListCopy is a VolumeSnapshotContentList type
type VolumeSnapshotContentListCopy VolumeSnapshotContentList

// VolumeSnapshotContentCopy is a VolumeSnapshotContent type
type VolumeSnapshotContentCopy VolumeSnapshotContent

// VolumeSnapshotListCopy is a VolumeSnapshotList type
type VolumeSnapshotListCopy VolumeSnapshotList

// VolumeSnapshotCopy is a VolumeSnapshot type
type VolumeSnapshotCopy VolumeSnapshot

// UnmarshalJSON unmarshalls json data
func (v *VolumeSnapshot) UnmarshalJSON(data []byte) error {
	tmp := VolumeSnapshotCopy{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tmp2 := VolumeSnapshot(tmp)
	*v = tmp2
	return nil
}

// UnmarshalJSON unmarshals json data
func (vd *VolumeSnapshotList) UnmarshalJSON(data []byte) error {
	tmp := VolumeSnapshotListCopy{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tmp2 := VolumeSnapshotList(tmp)
	*vd = tmp2
	return nil
}

// UnmarshalJSON unmarshals json data
func (v *VolumeSnapshotContent) UnmarshalJSON(data []byte) error {
	tmp := VolumeSnapshotContentCopy{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tmp2 := VolumeSnapshotContent(tmp)
	*v = tmp2
	return nil
}

// UnmarshalJSON unmarshals json data
func (vd *VolumeSnapshotContentList) UnmarshalJSON(data []byte) error {
	tmp := VolumeSnapshotContentListCopy{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tmp2 := VolumeSnapshotContentList(tmp)
	*vd = tmp2
	return nil
}
