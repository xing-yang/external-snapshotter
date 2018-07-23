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

package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/golang/glog"
	crdv1 "github.com/kubernetes-csi/external-snapshotter/pkg/apis/volumesnapshot/v1alpha1"
	clientset "github.com/kubernetes-csi/external-snapshotter/pkg/client/clientset/versioned"
	"github.com/kubernetes-csi/external-snapshotter/pkg/connection"
	"k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
)

// Handler is responsible for handling VolumeSnapshot events from informer.
type Handler interface {
	CreateSnapshotOperation(snapshot *crdv1.VolumeSnapshot) (*crdv1.VolumeSnapshot, error)
	DeleteSnapshotContentOperation(vsd *crdv1.VolumeSnapshotContent) error
	ListSnapshots(vsd *crdv1.VolumeSnapshotContent) (*csi.SnapshotStatus, int64, error)
	BindandUpdateVolumeSnapshot(snapshotContent *crdv1.VolumeSnapshotContent, snapshot *crdv1.VolumeSnapshot, status *crdv1.VolumeSnapshotStatus) (*crdv1.VolumeSnapshot, error)
	GetClassFromVolumeSnapshot(snapshot *crdv1.VolumeSnapshot) (*crdv1.VolumeSnapshotClass, error)
	UpdateVolumeSnapshotStatus(snapshot *crdv1.VolumeSnapshot, status *crdv1.VolumeSnapshotStatus) (*crdv1.VolumeSnapshot, error)
	GetSimplifiedSnapshotStatus(status *crdv1.VolumeSnapshotStatus) string
}

// csiHandler is a handler that calls CSI to create/delete volume snapshot.
type csiHandler struct {
	clientset                       clientset.Interface
	client                          kubernetes.Interface
	snapshotterName                 string
	eventRecorder                   record.EventRecorder
	csiConnection                   connection.CSIConnection
	timeout                         time.Duration
	createSnapshotContentRetryCount int
	createSnapshotContentInterval   time.Duration
}

func NewCSIHandler(
	clientset clientset.Interface,
	client kubernetes.Interface,
	snapshotterName string,
	eventRecorder record.EventRecorder,
	csiConnection connection.CSIConnection,
	timeout time.Duration,
	createSnapshotContentRetryCount int,
	createSnapshotContentInterval time.Duration,
) Handler {
	return &csiHandler{
		clientset:       clientset,
		client:          client,
		snapshotterName: snapshotterName,
		eventRecorder:   eventRecorder,
		csiConnection:   csiConnection,
		timeout:         timeout,
		createSnapshotContentRetryCount: createSnapshotContentRetryCount,
		createSnapshotContentInterval:   createSnapshotContentInterval,
	}
}

func (handler *csiHandler) takeSnapshot(snapshot *crdv1.VolumeSnapshot,
	volume *v1.PersistentVolume, parameters map[string]string) (*crdv1.VolumeSnapshotContent, *crdv1.VolumeSnapshotStatus, error) {
	glog.V(5).Infof("takeSnapshot: [%s]", snapshot.Name)
	ctx, cancel := context.WithTimeout(context.Background(), handler.timeout)
	defer cancel()

	driverName, snapshotId, timestamp, csiSnapshotStatus, err := handler.csiConnection.CreateSnapshot(ctx, snapshot, volume, parameters)
	if err != nil {
		// TODO: Handle Uploading gRPC error code 10
		return nil, nil, fmt.Errorf("failed to take snapshot of the volume %s: %q", volume.Name, err)
	}

	snapDataName := GetSnapshotContentNameForSnapshot(snapshot)

	// Create VolumeSnapshotContent in the database
	snapshotContent := &crdv1.VolumeSnapshotContent{
		ObjectMeta: metav1.ObjectMeta{
			Name: snapDataName,
		},
		Spec: crdv1.VolumeSnapshotContentSpec{
			VolumeSnapshotRef: &v1.ObjectReference{
				Kind:       "VolumeSnapshot",
				Namespace:  snapshot.Namespace,
				Name:       snapshot.Name,
				UID:        snapshot.UID,
				APIVersion: "v1alpha1",
			},
			PersistentVolumeRef: &v1.ObjectReference{
				Kind: "PersistentVolume",
				Name: volume.Name,
			},
			VolumeSnapshotSource: crdv1.VolumeSnapshotSource{
				CSI: &crdv1.CSIVolumeSnapshotSource{
					Driver:         driverName,
					SnapshotHandle: snapshotId,
					CreatedAt:      timestamp,
				},
			},
		},
	}

	status := ConvertSnapshotStatus(csiSnapshotStatus, timestamp)

	glog.V(5).Infof("takeSnapshot: Created snapshot [%s]. Snapshot object [%#v] Status [%#v]", snapshot.Name, snapshotContent, status)
	return snapshotContent, &status, nil
}

func (handler *csiHandler) deleteSnapshot(vsd *crdv1.VolumeSnapshotContent) error {
	if vsd.Spec.CSI == nil {
		return fmt.Errorf("CSISnapshot not defined in spec")
	}
	ctx, cancel := context.WithTimeout(context.Background(), handler.timeout)
	defer cancel()

	err := handler.csiConnection.DeleteSnapshot(ctx, vsd.Spec.CSI.SnapshotHandle)
	if err != nil {
		return fmt.Errorf("failed to delete snapshot data %s: %q", vsd.Name, err)
	}

	return nil
}

func (handler *csiHandler) ListSnapshots(vsd *crdv1.VolumeSnapshotContent) (*csi.SnapshotStatus, int64, error) {
	if vsd.Spec.CSI == nil {
		return nil, 0, fmt.Errorf("CSISnapshot not defined in spec")
	}
	ctx, cancel := context.WithTimeout(context.Background(), handler.timeout)
	defer cancel()

	csiSnapshotStatus, timestamp, err := handler.csiConnection.ListSnapshots(ctx, vsd.Spec.CSI.SnapshotHandle)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list snapshot data %s: %q", vsd.Name, err)
	}

	return csiSnapshotStatus, timestamp, nil
}

// The function goes through the whole snapshot creation process.
// 1. Update VolumeSnapshot metadata to include the snapshotted PV name, timestamp and snapshot uid, also generate tag for cloud provider
// 2. Trigger the snapshot through cloud provider and attach the tag to the snapshot.
// 3. Create the VolumeSnapshotContent object with the snapshot id information returned from step 2.
// 4. Bind the VolumeSnapshot and VolumeSnapshotContent object
// 5. Query the snapshot status through cloud provider and update the status until snapshot is ready or fails.
func (handler *csiHandler) CreateSnapshotOperation(snapshot *crdv1.VolumeSnapshot) (*crdv1.VolumeSnapshot, error) {
	glog.Infof("createSnapshot: Creating snapshot %s through the plugin ...", vsToVsKey(snapshot))
	var result *crdv1.VolumeSnapshot
	var err error
	class, err := handler.GetClassFromVolumeSnapshot(snapshot)
	if err != nil {
		glog.Errorf("creatSnapshotOperation failed to getClassFromVolumeSnapshot %s", err)
		return nil, err
	}

	//  A previous createSnapshot may just have finished while we were waiting for
	//  the locks. Check that snapshot data (with deterministic name) hasn't been created
	//  yet.
	snapDataName := GetSnapshotContentNameForSnapshot(snapshot)

	vsd, err := handler.clientset.VolumesnapshotV1alpha1().VolumeSnapshotContents().Get(snapDataName, metav1.GetOptions{})
	if err == nil && vsd != nil {
		// Volume snapshot data has been already created, nothing to do.
		glog.V(4).Infof("createSnapshot [%s]: volume snapshot data already exists, skipping", vsToVsKey(snapshot))
		return nil, nil
	}
	glog.V(5).Infof("createSnapshotOperation [%s]: VolumeSnapshotContent does not exist  yet", vsToVsKey(snapshot))

	volume, err := handler.getVolumeFromVolumeSnapshot(snapshot)
	if err != nil {
		glog.Errorf("createSnapshotOperation failed [%s]: Error: [%#v]", snapshot.Name, err)
		return nil, err
	}
	snapshotContent, status, err := handler.takeSnapshot(snapshot, volume, class.Parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to take snapshot of the volume %s: %q", volume.Name, err)
	}

	// Try to create the VSD object several times
	for i := 0; i < handler.createSnapshotContentRetryCount; i++ {
		glog.V(4).Infof("createSnapshot [%s]: trying to save volume snapshot data %s", vsToVsKey(snapshot), snapshotContent.Name)
		if _, err = handler.clientset.VolumesnapshotV1alpha1().VolumeSnapshotContents().Create(snapshotContent); err == nil || apierrs.IsAlreadyExists(err) {
			// Save succeeded.
			if err != nil {
				glog.V(3).Infof("volume snapshot data %q for snapshot %q already exists, reusing", snapshotContent.Name, vsToVsKey(snapshot))
				err = nil
			} else {
				glog.V(3).Infof("volume snapshot data %q for snapshot %q saved", snapshotContent.Name, vsToVsKey(snapshot))
			}
			break
		}
		// Save failed, try again after a while.
		glog.V(3).Infof("failed to save volume snapshot data %q for snapshot %q: %v", snapshotContent.Name, vsToVsKey(snapshot), err)
		time.Sleep(handler.createSnapshotContentInterval)
	}

	if err != nil {
		// Save failed. Now we have a storage asset outside of Kubernetes,
		// but we don't have appropriate volumesnapshotdata object for it.
		// Emit some event here and try to delete the storage asset several
		// times.
		strerr := fmt.Sprintf("Error creating volume snapshot data object for snapshot %s: %v. Deleting the snapshot data.", vsToVsKey(snapshot), err)
		glog.Error(strerr)
		handler.eventRecorder.Event(snapshot, v1.EventTypeWarning, "CreateSnapshotContentFailed", strerr)

		for i := 0; i < handler.createSnapshotContentRetryCount; i++ {
			if err = handler.deleteSnapshot(snapshotContent); err == nil {
				// Delete succeeded
				glog.V(4).Infof("createSnapshot [%s]: cleaning snapshot data %s succeeded", vsToVsKey(snapshot), snapshotContent.Name)
				break
			}
			// Delete failed, try again after a while.
			glog.Infof("failed to delete snapshot data %q: %v", snapshotContent.Name, err)
			time.Sleep(handler.createSnapshotContentInterval)
		}

		if err != nil {
			// Delete failed several times. There is an orphaned volume snapshot data and there
			// is nothing we can do about it.
			strerr := fmt.Sprintf("Error cleaning volume snapshot data for snapshot %s: %v. Please delete manually.", vsToVsKey(snapshot), err)
			glog.Error(strerr)
			handler.eventRecorder.Event(snapshot, v1.EventTypeWarning, "SnapshotContentCleanupFailed", strerr)
		}
	} else {
		// save succeeded, bind and update status for snapshot.
		result, err = handler.BindandUpdateVolumeSnapshot(snapshotContent, snapshot, status)
		if err != nil {
			return nil, err
		}
	}

	return result, err
}

// Delete a snapshot
// 1. Find the SnapshotContent corresponding to Snapshot
//   1a: Not found => finish (it's been deleted already)
// 2. Ask the backend to remove the snapshot device
// 3. Delete the SnapshotContent object
// 4. Remove the Snapshot from vsStore
// 5. Finish
func (handler *csiHandler) DeleteSnapshotContentOperation(vsd *crdv1.VolumeSnapshotContent) error {
	glog.V(4).Infof("deleteSnapshotOperation [%s] started", vsd.Name)

	err := handler.deleteSnapshot(vsd)
	if err != nil {
		return fmt.Errorf("failed to delete snapshot %#v, err: %v", vsd.Name, err)
	}

	err = handler.clientset.VolumesnapshotV1alpha1().VolumeSnapshotContents().Delete(vsd.Name, &metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete VolumeSnapshotContent %s from API server: %q", vsd.Name, err)
	}

	return nil
}

func (handler *csiHandler) BindandUpdateVolumeSnapshot(snapshotContent *crdv1.VolumeSnapshotContent, snapshot *crdv1.VolumeSnapshot, status *crdv1.VolumeSnapshotStatus) (*crdv1.VolumeSnapshot, error) {
	glog.V(4).Infof("bindandUpdateVolumeSnapshot for snapshot [%s]: snapshotContent [%s] status [%#v]", snapshot.Name, snapshotContent.Name, status)
	snapshotObj, err := handler.clientset.VolumesnapshotV1alpha1().VolumeSnapshots(snapshot.Namespace).Get(snapshot.Name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("error get snapshot %s from api server: %v", vsToVsKey(snapshot), err)
	}

	// Copy the snapshot object before updating it
	snapshotCopy := snapshotObj.DeepCopy()
	var updateSnapshot *crdv1.VolumeSnapshot
	if snapshotObj.Spec.SnapshotContentName == snapshotContent.Name {
		glog.Infof("bindVolumeSnapshotContentToVolumeSnapshot: VolumeSnapshot %s already bind to volumeSnapshotContent [%s]", snapshot.Name, snapshotContent.Name)
	} else {
		glog.Infof("bindVolumeSnapshotContentToVolumeSnapshot: before bind VolumeSnapshot %s to volumeSnapshotContent [%s]", snapshot.Name, snapshotContent.Name)
		snapshotCopy.Spec.SnapshotContentName = snapshotContent.Name
		updateSnapshot, err = handler.clientset.VolumesnapshotV1alpha1().VolumeSnapshots(snapshot.Namespace).Update(snapshotCopy)
		if err != nil {
			glog.Infof("bindVolumeSnapshotContentToVolumeSnapshot: Error binding VolumeSnapshot %s to volumeSnapshotContent [%s]. Error [%#v]", snapshot.Name, snapshotContent.Name, err)
			return nil, fmt.Errorf("error updating snapshot object %s on the API server: %v", vsToVsKey(updateSnapshot), err)
		}
		snapshotCopy = updateSnapshot
	}

	if status != nil && (status.CreatedAt != nil || status.AvailableAt != nil || status.Error != nil) {
		snapshotCopy.Status = *(status.DeepCopy())
		updateSnapshot2, err := handler.clientset.VolumesnapshotV1alpha1().VolumeSnapshots(snapshot.Namespace).Update(snapshotCopy)
		if err != nil {
			return nil, fmt.Errorf("error updating snapshot object %s on the API server: %v", vsToVsKey(snapshotCopy), err)
		}
		snapshotCopy = updateSnapshot2
	}

	glog.V(5).Infof("bindandUpdateVolumeSnapshot for snapshot completed [%#v]", snapshotCopy)
	return snapshotCopy, nil
}

// UpdateVolumeSnapshotStatus update VolumeSnapshot status.
func (handler *csiHandler) UpdateVolumeSnapshotStatus(snapshot *crdv1.VolumeSnapshot, status *crdv1.VolumeSnapshotStatus) (*crdv1.VolumeSnapshot, error) {
	snapshotObj, err := handler.clientset.VolumesnapshotV1alpha1().VolumeSnapshots(snapshot.Namespace).Get(snapshot.Name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("error get volume snapshot %s from  api server: %s", vsToVsKey(snapshot), err)
	}

	snapshotObj.Status = *status
	newSnapshotObj, err := handler.clientset.VolumesnapshotV1alpha1().VolumeSnapshots(snapshot.Namespace).Update(snapshotObj)
	if err != nil {
		return nil, fmt.Errorf("error update status for volume snapshot %s: %s", vsToVsKey(snapshot), err)
	}

	glog.Infof("UpdateVolumeSnapshotStatus finishes %+v", newSnapshotObj)
	return newSnapshotObj, nil
}

// getSimplifiedSnapshotStatus get status for snapshot.
func (handler *csiHandler) GetSimplifiedSnapshotStatus(status *crdv1.VolumeSnapshotStatus) string {
	if status == nil {
		glog.Errorf("No status for this snapshot yet.")
		return statusNew
	}

	if status.Error != nil {
		return statusError
	}
	if status.AvailableAt != nil {
		return statusReady
	}
	if status.CreatedAt != nil {
		return statusUploading
	}

	return statusNew
}

// getVolumeFromVolumeSnapshot is a helper function to get PV from VolumeSnapshot.
func (handler *csiHandler) getVolumeFromVolumeSnapshot(snapshot *crdv1.VolumeSnapshot) (*v1.PersistentVolume, error) {
	pvc, err := handler.getClaimFromVolumeSnapshot(snapshot)
	if err != nil {
		return nil, err
	}

	pvName := pvc.Spec.VolumeName
	pv, err := handler.client.CoreV1().PersistentVolumes().Get(pvName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve PV %s from the API server: %q", pvName, err)
	}

	glog.V(5).Infof("getVolumeFromVolumeSnapshot: snapshot [%s] PV name [%s]", snapshot.Name, pvName)

	return pv, nil
}

// getClassFromVolumeSnapshot is a helper function to get storage class from VolumeSnapshot.
func (handler *csiHandler) GetClassFromVolumeSnapshot(snapshot *crdv1.VolumeSnapshot) (*crdv1.VolumeSnapshotClass, error) {
	className := snapshot.Spec.VolumeSnapshotClassName
	glog.V(5).Infof("getClassFromVolumeSnapshot [%s]: VolumeSnapshotClassName [%s]", snapshot.Name, className)
	class, err := handler.clientset.VolumesnapshotV1alpha1().VolumeSnapshotClasses().Get(className, metav1.GetOptions{})
	if err != nil {
		glog.Errorf("failed to retrieve storage class %s from the API server: %q", className, err)
		//return nil, fmt.Errorf("failed to retrieve storage class %s from the API server: %q", className, err)
	}
	return class, nil
}

// getClaimFromVolumeSnapshot is a helper function to get PV from VolumeSnapshot.
func (handler *csiHandler) getClaimFromVolumeSnapshot(snapshot *crdv1.VolumeSnapshot) (*v1.PersistentVolumeClaim, error) {
	pvcName := snapshot.Spec.PersistentVolumeClaimName
	if pvcName == "" {
		return nil, fmt.Errorf("the PVC name is not specified in snapshot %s", vsToVsKey(snapshot))
	}

	pvc, err := handler.client.CoreV1().PersistentVolumeClaims(snapshot.Namespace).Get(pvcName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve PVC %s from the API server: %q", pvcName, err)
	}
	if pvc.Status.Phase != v1.ClaimBound {
		return nil, fmt.Errorf("the PVC %s not yet bound to a PV, will not attempt to take a snapshot yet", pvcName)
	}

	return pvc, nil
}
