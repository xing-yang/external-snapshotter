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

package controller

import (
	"fmt"
	//"strings"
	"time"

	crdv1 "github.com/kubernetes-csi/execution-hook/pkg/apis/executionhook/v1alpha1"
	//"k8s.io/api/core/v1"
	//storagev1 "k8s.io/api/storage/v1"
	//storage "k8s.io/api/storage/v1beta1"
	//apierrs "k8s.io/apimachinery/pkg/api/errors"
	//"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/labels"
	//"k8s.io/client-go/kubernetes/scheme"
	//ref "k8s.io/client-go/tools/reference"
	"k8s.io/klog"
	//"k8s.io/kubernetes/pkg/util/goroutinemap"
	//"k8s.io/kubernetes/pkg/util/goroutinemap/exponentialbackoff"
	//"k8s.io/kubernetes/pkg/util/slice"
)

const controllerUpdateFailMsg = "hook controller failed to update"

// syncContent deals with one key off the queue.  It returns false when it's time to quit.
func (ctrl *hookController) syncHookTemplate(hookTemplate *crdv1.ExecutionHookTemplate) error {
	klog.V(5).Infof("synchronizing ExecutionHookTemplate[%s]", hookTemplate.Name)

	return nil
}

// syncSnapshot is the main controller method to decide what to do with a snapshot.
// It's invoked by appropriate cache.Controller callbacks when a snapshot is
// created, updated or periodically synced. We do not differentiate between
// these events.
// For easier readability, it is split into syncUnreadySnapshot and syncReadySnapshot
// Wait for PreActionSucceed/PostActionSucceed to be true
func (ctrl *hookController) syncHook(hook *crdv1.ExecutionHook) error {
	klog.V(5).Infof("synchonizing ExecutionHook[%s]", hook.Name) //snapshotKey(snapshot), getSnapshotStatusForLogging(snapshot))

	return ctrl.syncPreActionHook(hook)
}

// syncPreActionHook is the main controller method to decide what to do with a hook whose PreActionSucceed summary status is not set to true.
func (ctrl *hookController) syncPreActionHook(hook *crdv1.ExecutionHook) error {
	// TODO: After Hook controller knows what containers and pods are selected,
	// it should create ContainerExecutionHookStatuses; otherwise we have to
	// save the container/pod list in cache
	// If PreActionSucceed summary status in ExecutionHookStatus is not set (nil),
	// loop around and wait until it is set

	statusFalse := false
	statusTrue := true
	klog.V(5).Infof("Entering syncPreActionHook with hook [%+v]", hook)

	// TODO: Make sure hook.Spec.PodContainerNamesList is populated before the next step
	if len(hook.Spec.PodContainerNamesList) == 0 {
		return fmt.Errorf("hook.Spec.PodContainerNamesList cannot be empty [%+v]", hook.Spec.PodContainerNamesList)
	}

	// hook.Status is of type crdv1.ExecutionHookStatus
	if hook.Status.PreActionSucceed == nil {
		// Initialize ContainerExecutionHookStatuses if it is empty based on
		// hook.Spec.PodContainerNamesList
		//var hookClone *crdv1.ExecutionHook
		//var hookUpdate *crdv1.ExecutionHook
		var err error
		klog.V(5).Infof("Xing 1: syncPreActionHook with hook [%+v]", hook)
		if len(hook.Status.ContainerExecutionHookStatuses) == 0 {
			hookClone := hook.DeepCopy()
			klog.V(5).Infof("Xing 1b: syncPreActionHook with hookClone [%+v]", hookClone)
			// TODO: Move the code to build initial ContainerExecutionHookStatuses
			// outside of syncPreActionHook so we don't need to retrieve the hook
			// from hookLister again
			// podContainers is of type crdv1.PodContainerNames
			for _, podContainers := range hookClone.Spec.PodContainerNamesList {
				for _, containerName := range podContainers.ContainerNames {
					containerStatus := crdv1.ContainerExecutionHookStatus{}
					klog.V(5).Infof("Xing 2: syncPreActionHook with PodContainerNamesList [%+v]", hook.Spec.PodContainerNamesList)
					containerStatus.PodName = podContainers.PodName
					containerStatus.ContainerName = containerName
					//containerStatus.PreActionStatus = nil
					klog.V(5).Infof("syncPreActionHook: add crdv1.ContainerExecutionHookStatus [%+v]", containerStatus)
					hookClone.Status.ContainerExecutionHookStatuses = append(hookClone.Status.ContainerExecutionHookStatuses, containerStatus)
					klog.V(5).Infof("Xing 3: syncPreActionHook with containerStatuses [%+v]", hookClone.Status.ContainerExecutionHookStatuses)
				}
			}
			_, err = ctrl.updateHookStatus(hookClone)
			if err != nil {
				return fmt.Errorf("failed to update hook status [%+v]", hookClone)
			}
		}

		//if hookUpdate0 == nil {
		//        return fmt.Errorf("Here! failed to get updated hook status [%+v]", hookClone)
		//}

		klog.V(5).Infof("Xing syncPreActionHook before hookUpdate0.DeepCopy with hook [%+v]", hook)
		// syncPreAction will be entered multiple times. When it enters for the 2nd time,
		// hook.Status.ContainerExecutionHookStatuses is no longer 0
		hookUpdate0, err := ctrl.hookLister.ExecutionHooks(hook.Namespace).Get(hook.Name)
		//hookNew, err := ctrl.clientset.ExecutionhookV1alpha1().ExecutionHooks(hookUpdate.Namespace).Get(hookUpdate.Name, metav1.GetOptions{})
		if err != nil {
			klog.Errorf("failed to retrieve hook %s from the informer: %q", hook.Name, err)
			return fmt.Errorf("failed to retrieve hook %s from the informer: %q", hook.Name, err)
		}
		hookUpdate := hookUpdate0.DeepCopy()
		statuses := []crdv1.ContainerExecutionHookStatus{}
		klog.V(5).Infof("Xing syncPreActionHook after hookUpdate0.DeepCopy with hook [%+v]", hook)
		// Loop through ContainerExecutionHookStatuses
		for i, _ := range hookUpdate0.Status.ContainerExecutionHookStatuses {
			// Hook PreAction has not started running if timestamp is nil
			// so run it now
			containerHookStatus := hookUpdate.Status.ContainerExecutionHookStatuses[i]
			if containerHookStatus.PreActionStatus == nil {
				containerHookStatus.PreActionStatus = &crdv1.ExecutionHookActionStatus{}
			}

			if containerHookStatus.PreActionStatus.ActionTimestamp == nil {
				// Set timestamp now
				var timeNow = metav1.Time{
					Time: time.Now(),
				}
				containerHookStatus.PreActionStatus.ActionTimestamp = &timeNow
				klog.V(5).Infof("Xing syncPreActionHook containerHookStatus actionTimestamp[%+v]", containerHookStatus.PreActionStatus.ActionTimestamp)
				// Run execution hook action
				// ExecPodContainerCommand will wait until it is done
				// TODO: add retry and timeout logic
				//err := ExecPodContainerCommand(hook.Namespace, containerHookStatus.PodName, containerHookStatus.ContainerName, hook.Name, hook.Spec.PreAction.Action.Exec.Command, hook.Spec.PreAction.ActionTimeoutSeconds)
				// TODO: Remove this
				var err error = nil

				// wait for it to come back
				// if success, set ActionSucceed = true
				// else, set ActionSucceed = false and
				// set hook.Status.PreActionSucceed = false
				// TODO: log an event
				// and bail out
				if err != nil { // failed on one command, no need to proceed, bail out
					containerHookStatus.PreActionStatus.ActionSucceed = &statusFalse
					hookUpdate.Status.PreActionSucceed = &statusFalse
					ctrl.updateHookStatus(hookUpdate)
					return fmt.Errorf("Failed to run PreAction %s in container %s in pod %s/%s", hookUpdate.Spec.PreAction.Action.Exec.Command, containerHookStatus.ContainerName, hookUpdate.Namespace, containerHookStatus.PodName)
				}
				// success
				successTrue := true
				containerHookStatus.PreActionStatus.ActionSucceed = &successTrue
				statuses = append(statuses, containerHookStatus)
				klog.V(5).Infof("Xing syncPreActionHook containerHookStatus actionTimestamp[%+v] actionStatus [%+v]", containerHookStatus.PreActionStatus.ActionTimestamp, containerHookStatus.PreActionStatus)
				klog.V(5).Infof("Xing syncPreActionHook containerHookStatus actionTimestamp[%+v] actionSucceed [%t]", containerHookStatus.PreActionStatus.ActionTimestamp, *(containerHookStatus.PreActionStatus.ActionSucceed))

				//klog.V(5).Infof("Xing5: syncPreActionHook update containerHookStatus [%+v]", containerHookStatus)
				//ctrl.updateHookStatus(hookUpdate)

			} else if containerHookStatus.PreActionStatus.ActionTimestamp != nil && containerHookStatus.PreActionStatus.ActionSucceed == nil {
				// Hook started but not complete yet
				// Wait for it to come back
				// if success, set ActionSucceed = true
				// else, set ActionSucceed = false and
				// set hook.Status.PreActionSucceed = false
				// log an event and bail out
			}
			if containerHookStatus.PreActionStatus.ActionTimestamp != nil && containerHookStatus.PreActionStatus.ActionSucceed != nil && *(containerHookStatus.PreActionStatus.ActionSucceed) == false {
				// It failed for this container preaction, bail out
				ctrl.updateHookStatus(hookUpdate)
				return fmt.Errorf("Failed to run PreAction %s in container %s in pod %s/%s", hookUpdate.Spec.PreAction.Action.Exec.Command, containerHookStatus.ContainerName, hookUpdate.Namespace, containerHookStatus.PodName)
			}
			//klog.V(5).Infof("Xing5: syncPreActionHook update containerHookStatus [%+v]", containerHookStatus)
			//ctrl.updateHookStatus(hookUpdate)
		}

		// successful
		if len(statuses) > 0 {
			hookUpdate.Status.ContainerExecutionHookStatuses = statuses
			ctrl.updateHookStatus(hookUpdate)
		}

		// Done with all containerHookStatus, set summary status to true
		// if it is not set to false yet, it means not failure occurred for a hook action on any container
		if hookUpdate.Status.PreActionSucceed == nil {
			//hookNew, err := ctrl.hookLister.ExecutionHooks(hookUpdate.Namespace).Get(hookUpdate.Name)
			hookNew, err := ctrl.clientset.ExecutionhookV1alpha1().ExecutionHooks(hookUpdate.Namespace).Get(hookUpdate.Name, metav1.GetOptions{})
			if err != nil {
				klog.Errorf("failed to retrieve hook %s from the informer: %q", hookUpdate.Name, err)
				return fmt.Errorf("failed to retrieve hook %s from the informer: %q", hookUpdate.Name, err)
			}
			klog.V(5).Infof("Xing: syncPreActionHook update PreAction summary status with hook [%+v]", hookNew)
			hookNew2 := hookNew.DeepCopy()
			hookNew2.Status.PreActionSucceed = &statusTrue
			ctrl.updateHookStatus(hookNew2)
			klog.V(5).Infof("PreAction %s ran successfully in selected containers and pods", hookNew.Spec.PreAction.Action.Exec.Command)
			klog.V(5).Infof("Xing Finished syncPreActionHook with hook [%+v]", hookNew2)
			// Need an event?
		}
	}

	return nil
}

// syncPostActionHook is the main controller method to decide what to do with a hook whose PostActionSucceed is not set to true.
func (ctrl *hookController) syncPostActionHook(hook *crdv1.ExecutionHook) error {
	statusFalse := false
	statusTrue := true
	klog.V(5).Infof("Entering syncPostActionHook with hook [%+v]", hook)

	// TODO: Make sure hook.Spec.PodContainerNamesList is populated before the next step
	if len(hook.Spec.PodContainerNamesList) == 0 {
		return fmt.Errorf("hook.Spec.PodContainerNamesList cannot be empty [%+v]", hook.Spec.PodContainerNamesList)
	}

	// hook.Status is of type crdv1.ExecutionHookStatus
	if hook.Status.PostActionSucceed == nil {
		if len(hook.Status.ContainerExecutionHookStatuses) == 0 {
			return fmt.Errorf("hook.Status.ContainerExecutionHookStatuses should not be empty at this point [%+v]. It should have been initialized in syncPostActionHook. Something is wrong.", hook.Status.ContainerExecutionHookStatuses)
		}

		hookClone := hook.DeepCopy()

		// Loop through ContainerExecutionHookStatuses
		for _, containerHookStatus := range hookClone.Status.ContainerExecutionHookStatuses {
			// Hook PostAction has not started running if timestamp is nil
			// so run it now
			if containerHookStatus.PostActionStatus == nil {
				containerHookStatus.PostActionStatus = &crdv1.ExecutionHookActionStatus{}
			}

			if containerHookStatus.PostActionStatus.ActionTimestamp == nil {
				// Set timestamp now
				var timeNow = metav1.Time{
					Time: time.Now(),
				}
				containerHookStatus.PostActionStatus.ActionTimestamp = &timeNow
				// Run execution hook action
				// ExecPodContainerCommand will wait until it is done
				// TODO: add retry and timeout logic
				err := ExecPodContainerCommand(hook.Namespace, containerHookStatus.PodName, containerHookStatus.ContainerName, hook.Name, hook.Spec.PostAction.Action.Exec.Command, hook.Spec.PostAction.ActionTimeoutSeconds)

				// wait for it to come back
				// if success, set ActionSucceed = true
				// else, set ActionSucceed = false and
				// set hook.Status.PostActionSucceed = false
				// TODO: log an event
				// and bail out
				if err != nil { // failed on one command, no need to proceed, bail out
					containerHookStatus.PostActionStatus.ActionSucceed = &statusFalse
					hookClone.Status.PostActionSucceed = &statusFalse
					ctrl.updateHookStatus(hookClone)
					return fmt.Errorf("Failed to run PostAction %s in container %s in pod %s/%s", hookClone.Spec.PostAction.Action.Exec.Command, containerHookStatus.ContainerName, hookClone.Namespace, containerHookStatus.PodName)
				}
				// success
				containerHookStatus.PostActionStatus.ActionSucceed = &statusTrue
			} else if containerHookStatus.PostActionStatus.ActionTimestamp != nil && containerHookStatus.PostActionStatus.ActionSucceed == nil {
				// Hook started but not complete yet
				// Wait for it to come back
				// if success, set ActionSucceed = true
				// else, set ActionSucceed = false and
				// set hook.Status.PostActionSucceed = false
				// log an event and bail out
			}
			if containerHookStatus.PostActionStatus.ActionTimestamp != nil && containerHookStatus.PostActionStatus.ActionSucceed != nil && *(containerHookStatus.PostActionStatus.ActionSucceed) == false {
				// It failed for this container postaction, bail out
				ctrl.updateHookStatus(hookClone)
				return fmt.Errorf("Failed to run PostAction %s in container %s in pod %s/%s", hookClone.Spec.PostAction.Action.Exec.Command, containerHookStatus.ContainerName, hookClone.Namespace, containerHookStatus.PodName)
			}
			ctrl.updateHookStatus(hookClone)
		}

		// Done with all containerHookStatus, set summary status to true
		// if it is not set to false yet, it means not failure occurred for a hook action on any container
		if hookClone.Status.PostActionSucceed == nil {
			hookClone.Status.PostActionSucceed = &statusTrue
			ctrl.updateHookStatus(hookClone)
			klog.V(5).Infof("PostAction %s ran successfully in selected containers and pods", hookClone.Spec.PostAction.Action.Exec.Command)
			// Need an event?
		}
	}

	return nil
}

// updateHookStatus updates the status of the hook
// Note that hook passed in should be a clone of the original hook to make sure
// object update is successful
func (ctrl *hookController) updateHookStatus(hook *crdv1.ExecutionHook) (*crdv1.ExecutionHook, error) {
	klog.V(5).Infof("updateHookStatus[%s]", hookKey(hook))

	newHook, err := ctrl.clientset.ExecutionhookV1alpha1().ExecutionHooks(hook.Namespace).Update(hook)
	if err != nil {
		klog.V(4).Infof("updating ExecutionHook[%s] error status failed %v", hookKey(hook), err)
		return nil, err
	}

	_, err = ctrl.storeHookUpdate(newHook)
	if err != nil {
		klog.V(4).Infof("updating ExecutionHook[%s] error status: cannot update internal cache %v", hookKey(hook), err)
		return nil, err
	}
	//ctrl.eventRecorder.Event(newHook, eventtype, reason, message)
	//ctrl.eventRecorder.Event(newHook, eventtype, reason, message)

	return newHook, nil
}

/*
// scheduleOperation starts given asynchronous operation on given volume. It
// makes sure the operation is already not running.
func (ctrl *hookController) scheduleOperation(operationName string, operation func() error) {
	klog.V(5).Infof("scheduleOperation[%s]", operationName)

	err := ctrl.runningOperations.Run(operationName, operation)
	if err != nil {
		switch {
		case goroutinemap.IsAlreadyExists(err):
			klog.V(4).Infof("operation %q is already running, skipping", operationName)
		case exponentialbackoff.IsExponentialBackoff(err):
			klog.V(4).Infof("operation %q postponed due to exponential backoff", operationName)
		default:
			klog.Errorf("error scheduling operation %q: %v", operationName, err)
		}
	}
}
*/

func (ctrl *hookController) storeHookUpdate(hook interface{}) (bool, error) {
	return storeObjectUpdate(ctrl.hookStore, hook, "hook")
}

func (ctrl *hookController) storeHookTemplateUpdate(hookTemplate interface{}) (bool, error) {
	return storeObjectUpdate(ctrl.hookTemplateStore, hookTemplate, "hookTemplate")
}

/*
// createSnapshot starts new asynchronous operation to create snapshot
func (ctrl *hookController) createSnapshot(snapshot *crdv1.VolumeSnapshot) error {
	klog.V(5).Infof("createSnapshot[%s]: started", snapshotKey(snapshot))
	opName := fmt.Sprintf("create-%s[%s]", snapshotKey(snapshot), string(snapshot.UID))
	ctrl.scheduleOperation(opName, func() error {
		snapshotObj, err := ctrl.createSnapshotOperation(snapshot)
		if err != nil {
			ctrl.updateSnapshotErrorStatusWithEvent(snapshot, v1.EventTypeWarning, "SnapshotCreationFailed", fmt.Sprintf("Failed to create snapshot: %v", err))
			klog.Errorf("createSnapshot [%s]: error occurred in createSnapshotOperation: %v", opName, err)
			return err
		}
		_, updateErr := ctrl.storeSnapshotUpdate(snapshotObj)
		if updateErr != nil {
			// We will get an "snapshot update" event soon, this is not a big error
			klog.V(4).Infof("createSnapshot [%s]: cannot update internal cache: %v", snapshotKey(snapshotObj), updateErr)
		}

		return nil
	})
	return nil
}

// updateSnapshotStatusWithEvent saves new snapshot.Status to API server and emits
// given event on the snapshot. It saves the status and emits the event only when
// the status has actually changed from the version saved in API server.
// Parameters:
//   snapshot - snapshot to update
//   eventtype, reason, message - event to send, see EventRecorder.Event()
func (ctrl *hookController) updateSnapshotErrorStatusWithEvent(snapshot *crdv1.VolumeSnapshot, eventtype, reason, message string) error {
	klog.V(5).Infof("updateSnapshotStatusWithEvent[%s]", snapshotKey(snapshot))

	if snapshot.Status.Error != nil && snapshot.Status.Error.Message == message {
		klog.V(4).Infof("updateSnapshotStatusWithEvent[%s]: the same error %v is already set", snapshot.Name, snapshot.Status.Error)
		return nil
	}
	snapshotClone := snapshot.DeepCopy()
	statusError := &storage.VolumeError{
		Time: metav1.Time{
			Time: time.Now(),
		},
		Message: message,
	}
	snapshotClone.Status.Error = statusError

	snapshotClone.Status.ReadyToUse = false
	newSnapshot, err := ctrl.clientset.VolumesnapshotV1alpha1().VolumeSnapshots(snapshotClone.Namespace).Update(snapshotClone)
	if err != nil {
		klog.V(4).Infof("updating VolumeSnapshot[%s] error status failed %v", snapshotKey(snapshot), err)
		return err
	}

	_, err = ctrl.storeSnapshotUpdate(newSnapshot)
	if err != nil {
		klog.V(4).Infof("updating VolumeSnapshot[%s] error status: cannot update internal cache %v", snapshotKey(snapshot), err)
		return err
	}
	// Emit the event only when the status change happens
	ctrl.eventRecorder.Event(newSnapshot, eventtype, reason, message)

	return nil
}

// Stateless functions
func getSnapshotStatusForLogging(snapshot *crdv1.VolumeSnapshot) string {
	return fmt.Sprintf("bound to: %q, Completed: %v", snapshot.Spec.SnapshotContentName, snapshot.Status.ReadyToUse)
}

// The function goes through the whole snapshot creation process.
// 1. Trigger the snapshot through csi storage provider.
// 2. Update VolumeSnapshot status with creationtimestamp information
// 3. Create the VolumeSnapshotContent object with the snapshot id information.
// 4. Bind the VolumeSnapshot and VolumeSnapshotContent object
func (ctrl *hookController) createSnapshotOperation(snapshot *crdv1.VolumeSnapshot) (*crdv1.VolumeSnapshot, error) {
	klog.Infof("createSnapshot: Creating snapshot %s through the plugin ...", snapshotKey(snapshot))

	if snapshot.Status.Error != nil && !isControllerUpdateFailError(snapshot.Status.Error) {
		klog.V(4).Infof("error is already set in snapshot, do not retry to create: %s", snapshot.Status.Error.Message)
		return snapshot, nil
	}

	class, volume, contentName, snapshotterCredentials, err := ctrl.getCreateSnapshotInput(snapshot)
	if err != nil {
		return nil, fmt.Errorf("failed to get input parameters to create snapshot %s: %q", snapshot.Name, err)
	}

	driverName, snapshotID, timestamp, size, readyToUse, err := ctrl.handler.CreateSnapshot(snapshot, volume, class.Parameters, snapshotterCredentials)
	if err != nil {
		return nil, fmt.Errorf("failed to take snapshot of the volume, %s: %q", volume.Name, err)
	}
	klog.V(5).Infof("Created snapshot: driver %s, snapshotId %s, timestamp %d, size %d, readyToUse %t", driverName, snapshotID, timestamp, size, readyToUse)

	var newSnapshot *crdv1.VolumeSnapshot
	// Update snapshot status with timestamp
	for i := 0; i < ctrl.createSnapshotContentRetryCount; i++ {
		klog.V(5).Infof("createSnapshot [%s]: trying to update snapshot creation timestamp", snapshotKey(snapshot))
		newSnapshot, err = ctrl.updateSnapshotStatus(snapshot, readyToUse, timestamp, size, false)
		if err == nil {
			break
		}
		klog.V(4).Infof("failed to update snapshot %s creation timestamp: %v", snapshotKey(snapshot), err)
	}

	if err != nil {
		return nil, err
	}
	// Create VolumeSnapshotContent in the database
	volumeRef, err := ref.GetReference(scheme.Scheme, volume)
	if err != nil {
		return nil, err
	}
	snapshotRef, err := ref.GetReference(scheme.Scheme, snapshot)
	if err != nil {
		return nil, err
	}

	if class.DeletionPolicy == nil {
		class.DeletionPolicy = new(crdv1.DeletionPolicy)
		*class.DeletionPolicy = crdv1.VolumeSnapshotContentDelete
	}
	snapshotContent := &crdv1.VolumeSnapshotContent{
		ObjectMeta: metav1.ObjectMeta{
			Name: contentName,
		},
		Spec: crdv1.VolumeSnapshotContentSpec{
			VolumeSnapshotRef:   snapshotRef,
			PersistentVolumeRef: volumeRef,
			VolumeSnapshotSource: crdv1.VolumeSnapshotSource{
				CSI: &crdv1.CSIVolumeSnapshotSource{
					Driver:         driverName,
					SnapshotHandle: snapshotID,
					CreationTime:   &timestamp,
					RestoreSize:    &size,
				},
			},
			VolumeSnapshotClassName: &(class.Name),
			DeletionPolicy:          class.DeletionPolicy,
		},
	}
	klog.V(3).Infof("volume snapshot content %v", snapshotContent)
	// Try to create the VolumeSnapshotContent object several times
	for i := 0; i < ctrl.createSnapshotContentRetryCount; i++ {
		klog.V(5).Infof("createSnapshot [%s]: trying to save volume snapshot content %s", snapshotKey(snapshot), snapshotContent.Name)
		if _, err = ctrl.clientset.VolumesnapshotV1alpha1().VolumeSnapshotContents().Create(snapshotContent); err == nil || apierrs.IsAlreadyExists(err) {
			// Save succeeded.
			if err != nil {
				klog.V(3).Infof("volume snapshot content %q for snapshot %q already exists, reusing", snapshotContent.Name, snapshotKey(snapshot))
				err = nil
			} else {
				klog.V(3).Infof("volume snapshot content %q for snapshot %q saved, %v", snapshotContent.Name, snapshotKey(snapshot), snapshotContent)
			}
			break
		}
		// Save failed, try again after a while.
		klog.V(3).Infof("failed to save volume snapshot content %q for snapshot %q: %v", snapshotContent.Name, snapshotKey(snapshot), err)
		time.Sleep(ctrl.createSnapshotContentInterval)
	}

	if err != nil {
		// Save failed. Now we have a snapshot asset outside of Kubernetes,
		// but we don't have appropriate volumesnapshot content object for it.
		// Emit some event here and controller should try to create the content in next sync period.
		strerr := fmt.Sprintf("Error creating volume snapshot content object for snapshot %s: %v.", snapshotKey(snapshot), err)
		klog.Error(strerr)
		ctrl.eventRecorder.Event(newSnapshot, v1.EventTypeWarning, "CreateSnapshotContentFailed", strerr)
		return nil, newControllerUpdateError(snapshotKey(snapshot), err.Error())
	}

	// save succeeded, bind and update status for snapshot.
	result, err := ctrl.bindandUpdateVolumeSnapshot(snapshotContent, newSnapshot)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Delete a snapshot
// 1. Find the SnapshotContent corresponding to Snapshot
//   1a: Not found => finish (it's been deleted already)
// 2. Ask the backend to remove the snapshot device
// 3. Delete the SnapshotContent object
// 4. Remove the Snapshot from store
// 5. Finish
func (ctrl *hookController) deleteSnapshotContentOperation(content *crdv1.VolumeSnapshotContent) error {
	klog.V(5).Infof("deleteSnapshotOperation [%s] started", content.Name)

	// get secrets if VolumeSnapshotClass specifies it
	var snapshotterCredentials map[string]string
	snapshotClassName := content.Spec.VolumeSnapshotClassName
	if snapshotClassName != nil {
		if snapshotClass, err := ctrl.classLister.Get(*snapshotClassName); err == nil {
			// Resolve snapshotting secret credentials.
			// No VolumeSnapshot is provided when resolving delete secret names, since the VolumeSnapshot may or may not exist at delete time.
			snapshotterSecretRef, err := getSecretReference(snapshotClass.Parameters, content.Name, nil)
			if err != nil {
				return err
			}
			snapshotterCredentials, err = getCredentials(ctrl.client, snapshotterSecretRef)
			if err != nil {
				return err
			}
		}
	}

	err := ctrl.handler.DeleteSnapshot(content, snapshotterCredentials)
	if err != nil {
		ctrl.eventRecorder.Event(content, v1.EventTypeWarning, "SnapshotDeleteError", "Failed to delete snapshot")
		return fmt.Errorf("failed to delete snapshot %#v, err: %v", content.Name, err)
	}

	err = ctrl.clientset.VolumesnapshotV1alpha1().VolumeSnapshotContents().Delete(content.Name, &metav1.DeleteOptions{})
	if err != nil {
		ctrl.eventRecorder.Event(content, v1.EventTypeWarning, "SnapshotContentObjectDeleteError", "Failed to delete snapshot content API object")
		return fmt.Errorf("failed to delete VolumeSnapshotContent %s from API server: %q", content.Name, err)
	}

	return nil
}

var _ error = controllerUpdateError{}

type controllerUpdateError struct {
	message string
}

func newControllerUpdateError(name, message string) error {
	return controllerUpdateError{
		message: fmt.Sprintf("%s %s on API server: %s", controllerUpdateFailMsg, name, message),
	}
}

func (e controllerUpdateError) Error() string {
	return e.message
}

func isControllerUpdateFailError(err *storage.VolumeError) bool {
	if err != nil {
		if strings.Contains(err.Message, controllerUpdateFailMsg) {
			return true
		}
	}
	return false
}
*/

/*
func (ctrl *hookController) executeHook(hook *crdv1.ExecutionHook) error {
        // Freeze before taking snapshot
        klog.Infof("createSnapshot: Freeze before Creating snapshot %s through the plugin ...", snapshotKey(snapshot))
        //e := &defaultPodCommandExecutor{}
        // command := []string{"echo Hi"}
        //command := []string{"/usr/sbin/nginx -s quit"}
        //command := []string{"mysql -u root -ppassword -D mysql -e \"FLUSH TABLES WITH READ LOCK;\""}
        //command := []string{"//usr/bin/mysql -D mysql -e \"FLUSH TABLES WITH READ LOCK;\""}
        //command := []string{"python /root/quiesce.py"}
        //command := []string{"ls"}
        command := []string{"run_quiesce.sh"}
        klog.Infof("createSnapshot: command [%#v] Exe command [%#v] ...", command, container.QuiesceUnquiesceHook.Quiesce.Handler.Exec.Command)
        //myhook := ExecHook{
        //        Container: container.Name,                                    //"mysql", //nginx",
        //        Command:   command, //container.QuiesceUnquiesceHook.Quiesce.Handler.Exec.Command, //command,
        //        OnError:   "",
        //        Timeout:   1 * time.Minute,
        //}
        //hook := &myhook
        // Snapshot can only be in the same namespace as the container and pod
        namespace := snapshot.Namespace                          //"default"
        name := container.Name                                    //"mysql" //"nginx"
        er := ExecutePodCommand(namespace, name, "quiesce", hook) // works!!!
        if er != nil {
                klog.Infof("createSnapshot: failed to freeze before Creating snapshot %s through the plugin ...", snapshotKey(snapshot))
                return nil, fmt.Errorf("failed to freeze before create snapshot %s: %q", snapshot.Name, er)
        }
        klog.Infof("createSnapshot: After freeze before Creating snapshot %s through the plugin ...", snapshotKey(snapshot))
        // After Freeze call
}*/

/*
Snapshot controller watches VolumeAttachment object and search for a matching PersistentVolumeName
With Snapshot’s source PersistentVolume.
Attacher should be equal to Snapshotter.
Also check VolumeAttachmentStatus Attached = true

If it is attached, it means we need to do a freeze.
Also from NodeName in VolumeAttachmentSpec, we can find out the matching pod that is using the PVC.

Snapshot controller also watches PodSpec, and find out which pod this volume is being used.

If containers[] in the pod contains execution hook, run it before snapshotting. (Option 1: Use PodCommandExec to run it remotely)
Option 2: external-snapshotter pass the call to CSI driver to do freeze.  It needs to pass in the VolumeAttachment related info so that
The CSI driver knows which filesystem it is and how to do fsfreeze.  Using Option 2, we don’t need Exec Hook in the Pod.
*/
/*func (ctrl *hookController) executeHook(snapshot *crdv1.VolumeSnapshot) (*v1.Container, error) {
  klog.Infof("executeHook before creating snapshot %s through the plugin ...", snapshotKey(snapshot))

  //pvcList, err := ctrl.pvcLister.PersistentVolumeClaims(snapshot.Namespace).List(labels.Everything())
  //pvc, err := ctrl.pvcLister.PersistentVolumeClaims(snapshot.Namespace).Get(pvcName)

  // VolumeAttachmentLister helps list VolumeAttachments.
  //type VolumeAttachmentLister interface {
  // List lists all VolumeAttachments in the indexer.
  //List(selector labels.Selector) (ret []*v1beta1.VolumeAttachment, err error)
  // Get retrieves the VolumeAttachment from the index for a given name.
  //Get(name string) (*v1beta1.VolumeAttachment, error)
  //VolumeAttachmentListerExpansion
  //}
  //va, err := ctrl.vaLister.Get(vaName)
  pvc, err := ctrl.getClaimFromVolumeSnapshot(snapshot)
  if err != nil {
          klog.Errorf("executeHook: failed to getClaimFromVolumeSnapshot %s ...", snapshotKey(snapshot))
          return nil, err
  }
  volume, err := ctrl.getVolumeFromVolumeSnapshot(snapshot)
  if err != nil {
          klog.Errorf("executeHook: failed to getVolumeFromVolumeSnapshot %s ...", snapshotKey(snapshot))
          return nil, err
  }
  klog.Infof("executeHook: found volume %v ...", volume)
  volumeattachments, err := ctrl.vaLister.List(labels.Everything())
  if err != nil {
          klog.Errorf("executeHook: failed to get volumeattachments %s ...", snapshotKey(snapshot))
          return nil, err
  }
  var foundVA *storage.VolumeAttachment = nil
  for _, va := range volumeattachments {
          klog.Infof("executeHook: looping thru volumeattachments %v ...", va)
          if va.Spec.Source.PersistentVolumeName != nil {
                  // compare with PV name of snapshot source
                  if *(va.Spec.Source.PersistentVolumeName) == volume.Name {
                          foundVA = va.DeepCopy()
                          klog.Infof("executeHook: foundVA %v ...", foundVA)
                          break
                  }
          }
  }
  if foundVA == nil {
          klog.Infof("no attachment found when taking snapshot %s. No need to freeze", snapshot.Name)
          return nil, nil
  }
  if foundVA.Spec.Attacher != ctrl.snapshotterName {
          return nil, fmt.Errorf("Driver for attacher [%s] and snapshotter [%s] does not match for snapshot %s", foundVA.Spec.Attacher, ctrl.snapshotterName, snapshotKey(snapshot))
  }
  // TODO: Check status of va
  var needToFreeze = false
  if foundVA.Status.Attached == true {
          needToFreeze = true
  }
  klog.V(5).Infof("syncSnapshot: Need to freeze is [%t] for VolumeSnapshot[%s]", needToFreeze, snapshotKey(snapshot))

  // Go through the pods
  // PodLister helps list Pods.
  ///type PodLister interface {
  //  // List lists all Pods in the indexer.
  //  List(selector labels.Selector) (ret []*v1.Pod, err error)
  //  // Pods returns an object that can list and get Pods.
  //  Pods(namespace string) PodNamespaceLister
  //  PodListerExpansion
  //  }
  pods, err := ctrl.podLister.List(labels.Everything())
  if err != nil {
          klog.Errorf("executeHook: failed to get pods %s ...", snapshotKey(snapshot))
          return nil, err
  }
  var foundPod *v1.Pod
  for _, pd := range pods {
          for _, vol := range pd.Spec.Volumes {
                  klog.V(5).Infof("executeHook: pod volumes %+v", vol)
                  if vol.VolumeSource.PersistentVolumeClaim != nil && vol.VolumeSource.PersistentVolumeClaim.ClaimName == pvc.Name {
                          foundPod = pd.DeepCopy()
                          klog.V(5).Infof("execHook Pod: [%+v]", foundPod)
                          break
                  }
          }
  }
  for _, container := range foundPod.Spec.Containers {
          klog.V(5).Infof("Container: [%+v]", container)
          if container.QuiesceUnquiesceHook != nil {
                  klog.V(5).Infof("container.QuiesceUnquiesceHook: [%+v]", container.QuiesceUnquiesceHook)
          }
          //Run container.QuiesceUnquiesceHook.Quiesce.Exec.Command
          // Create Pod subresource
          // container.QuiesceUnquiesceHook.Quiesce.Exec.Command
          // []Hooks under PodSpec or a parameter inside Container?
          // If it is directly under PodSpec, we also need to add container name to it
          // To run unquiesce, we have to make sure all snapshots are taken
          //er := ExecutePodCommand(container.namespace, container.name, "quiesce", hook)
          //if er != nil {
          //      klog.Infof("execHook: failed to freeze before Creating snapshot %s through the plugin ...", snapshotKey(snapshot))
          //      return nil, fmt.Errorf("failed to freeze before create snapshot %s: %q", snapshot.Name, er)
          //}
          return &container, nil
  }
  //ExecutionHook:
  //  quiesceHook:
  //    exec:
  //      command: ["/usr/sbin/nginx","-s","quit";]
  //  unquiesceHook:
  //    exec:
  //      command: ["/usr/sbin/nginx","-s","start"]
*/

/*
   type ExecutionHook Struct {
       // A Trigger is in what condition the execution handler should be executed
       Trigger ExecutionHookTrigger

      // Command to execute for a particular trigger
      Handler

      // How long the controller should try/retry to execute the hook before giving up
      RetryTimeOutSeconds int
   }

   type ExecutionHookTrigger string

   const (
       PreSnapshot ExecutionHookTrigger = “PreSnapshot”
       PostSnapshot ExecutionHookTrigger = “PostSnapshot”
   )
   type Handler struct {
       // One and only one of the following should be specified.
       // Exec specifies the action to take.
       // +optional
       Exec *ExecAction
       // HTTPGet specifies the http request to perform.
       // +optional
       HTTPGet *HTTPGetAction
       // TCPSocket specifies an action involving a TCP port.
       // TODO: implement a realistic TCP lifecycle hook
       // +optional
       TCPSocket *TCPSocketAction
   }
   // ExecAction describes a "run in container" action.
   type ExecAction struct {
           // Command is the command line to execute inside the container
           Command []string
   }
*/ /*
        return nil, nil
}*/
