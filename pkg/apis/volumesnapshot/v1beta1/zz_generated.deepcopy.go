// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VolumeSnapshot) DeepCopyInto(out *VolumeSnapshot) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VolumeSnapshot.
func (in *VolumeSnapshot) DeepCopy() *VolumeSnapshot {
	if in == nil {
		return nil
	}
	out := new(VolumeSnapshot)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VolumeSnapshot) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VolumeSnapshotClass) DeepCopyInto(out *VolumeSnapshotClass) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Parameters != nil {
		in, out := &in.Parameters, &out.Parameters
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.DeletionPolicy != nil {
		in, out := &in.DeletionPolicy, &out.DeletionPolicy
		*out = new(DeletionPolicy)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VolumeSnapshotClass.
func (in *VolumeSnapshotClass) DeepCopy() *VolumeSnapshotClass {
	if in == nil {
		return nil
	}
	out := new(VolumeSnapshotClass)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VolumeSnapshotClass) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VolumeSnapshotClassList) DeepCopyInto(out *VolumeSnapshotClassList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]VolumeSnapshotClass, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VolumeSnapshotClassList.
func (in *VolumeSnapshotClassList) DeepCopy() *VolumeSnapshotClassList {
	if in == nil {
		return nil
	}
	out := new(VolumeSnapshotClassList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VolumeSnapshotClassList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VolumeSnapshotContent) DeepCopyInto(out *VolumeSnapshotContent) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VolumeSnapshotContent.
func (in *VolumeSnapshotContent) DeepCopy() *VolumeSnapshotContent {
	if in == nil {
		return nil
	}
	out := new(VolumeSnapshotContent)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VolumeSnapshotContent) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VolumeSnapshotContentList) DeepCopyInto(out *VolumeSnapshotContentList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]VolumeSnapshotContent, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VolumeSnapshotContentList.
func (in *VolumeSnapshotContentList) DeepCopy() *VolumeSnapshotContentList {
	if in == nil {
		return nil
	}
	out := new(VolumeSnapshotContentList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VolumeSnapshotContentList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VolumeSnapshotContentSpec) DeepCopyInto(out *VolumeSnapshotContentSpec) {
	*out = *in
	if in.VolumeSnapshotRef != nil {
		in, out := &in.VolumeSnapshotRef, &out.VolumeSnapshotRef
		*out = new(v1.ObjectReference)
		**out = **in
	}
	if in.DeletionPolicy != nil {
		in, out := &in.DeletionPolicy, &out.DeletionPolicy
		*out = new(DeletionPolicy)
		**out = **in
	}
	if in.SnapshotHandle != nil {
		in, out := &in.SnapshotHandle, &out.SnapshotHandle
		*out = new(string)
		**out = **in
	}
	if in.CreationTime != nil {
		in, out := &in.CreationTime, &out.CreationTime
		*out = new(int64)
		**out = **in
	}
	if in.RestoreSize != nil {
		in, out := &in.RestoreSize, &out.RestoreSize
		*out = new(int64)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VolumeSnapshotContentSpec.
func (in *VolumeSnapshotContentSpec) DeepCopy() *VolumeSnapshotContentSpec {
	if in == nil {
		return nil
	}
	out := new(VolumeSnapshotContentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VolumeSnapshotContentStatus) DeepCopyInto(out *VolumeSnapshotContentStatus) {
	*out = *in
	if in.Error != nil {
		in, out := &in.Error, &out.Error
		*out = new(VolumeSnapshotError)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VolumeSnapshotContentStatus.
func (in *VolumeSnapshotContentStatus) DeepCopy() *VolumeSnapshotContentStatus {
	if in == nil {
		return nil
	}
	out := new(VolumeSnapshotContentStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VolumeSnapshotError) DeepCopyInto(out *VolumeSnapshotError) {
	*out = *in
	if in.Time != nil {
		in, out := &in.Time, &out.Time
		*out = (*in).DeepCopy()
	}
	if in.Message != nil {
		in, out := &in.Message, &out.Message
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VolumeSnapshotError.
func (in *VolumeSnapshotError) DeepCopy() *VolumeSnapshotError {
	if in == nil {
		return nil
	}
	out := new(VolumeSnapshotError)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VolumeSnapshotList) DeepCopyInto(out *VolumeSnapshotList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]VolumeSnapshot, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VolumeSnapshotList.
func (in *VolumeSnapshotList) DeepCopy() *VolumeSnapshotList {
	if in == nil {
		return nil
	}
	out := new(VolumeSnapshotList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VolumeSnapshotList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VolumeSnapshotSource) DeepCopyInto(out *VolumeSnapshotSource) {
	*out = *in
	if in.PersistentVolumeClaimName != nil {
		in, out := &in.PersistentVolumeClaimName, &out.PersistentVolumeClaimName
		*out = new(string)
		**out = **in
	}
	if in.VolumeSnapshotContentName != nil {
		in, out := &in.VolumeSnapshotContentName, &out.VolumeSnapshotContentName
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VolumeSnapshotSource.
func (in *VolumeSnapshotSource) DeepCopy() *VolumeSnapshotSource {
	if in == nil {
		return nil
	}
	out := new(VolumeSnapshotSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VolumeSnapshotSpec) DeepCopyInto(out *VolumeSnapshotSpec) {
	*out = *in
	if in.Source != nil {
		in, out := &in.Source, &out.Source
		*out = new(VolumeSnapshotSource)
		(*in).DeepCopyInto(*out)
	}
	if in.VolumeSnapshotClassName != nil {
		in, out := &in.VolumeSnapshotClassName, &out.VolumeSnapshotClassName
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VolumeSnapshotSpec.
func (in *VolumeSnapshotSpec) DeepCopy() *VolumeSnapshotSpec {
	if in == nil {
		return nil
	}
	out := new(VolumeSnapshotSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VolumeSnapshotStatus) DeepCopyInto(out *VolumeSnapshotStatus) {
	*out = *in
	if in.BoundVolumeSnapshotContentName != nil {
		in, out := &in.BoundVolumeSnapshotContentName, &out.BoundVolumeSnapshotContentName
		*out = new(string)
		**out = **in
	}
	if in.CreationTime != nil {
		in, out := &in.CreationTime, &out.CreationTime
		*out = (*in).DeepCopy()
	}
	if in.RestoreSize != nil {
		in, out := &in.RestoreSize, &out.RestoreSize
		x := (*in).DeepCopy()
		*out = &x
	}
	if in.Error != nil {
		in, out := &in.Error, &out.Error
		*out = new(VolumeSnapshotError)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VolumeSnapshotStatus.
func (in *VolumeSnapshotStatus) DeepCopy() *VolumeSnapshotStatus {
	if in == nil {
		return nil
	}
	out := new(VolumeSnapshotStatus)
	in.DeepCopyInto(out)
	return out
}
