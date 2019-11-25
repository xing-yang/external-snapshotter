#!/bin/bash

########################
# include the magic
########################
. demo-magic.sh

DEMO_PROMPT="${GREEN}➜ ${CYAN}\W "

# hide the evidence
clear

# Put your stuff here
export KUBECONFIG=/var/run/kubernetes/admin.kubeconfig

echo "==> Install snapshot beta CRDs"
pe "cd ~/go/src/github.com/kubernetes-csi/external-snapshotter/config/crd"
pe "ls"
pe "vi snapshot.storage.k8s.io_volumesnapshots.yaml"
pe "kubectl create -f ~/go/src/github.com/kubernetes-csi/external-snapshotter/config/crd"

echo "==> Install snapshot controller"
pe "cd ~/go/src/github.com/kubernetes-csi/external-snapshotter/deploy/kubernetes/snapshot-controller/"
pe "ls"
pe "vi setup-snapshot-controller.yaml"
pe "kubectl create -f ~/go/src/github.com/kubernetes-csi/external-snapshotter/deploy/kubernetes/snapshot-controller/"

echo "==> Install csi-provisioner and csi-snapshotter sidecar, and CSI driver"
pe "cd ~/go/src/github.com/kubernetes-csi/external-snapshotter/deploy/kubernetes/csi-snapshotter/"
pe "ls"
pe "vi setup-csi-snapshotter.yaml"
pe "kubectl create -f ~/go/src/github.com/kubernetes-csi/external-snapshotter/deploy/kubernetes/csi-snapshotter/"

pe "kubectl get pod"

echo "==> Test dynamically creating a snapshot"
pe "cd ~/go/src/github.com/kubernetes-csi/external-snapshotter/examples/kubernetes"

pe "kubectl create -f storageclass.yaml"
pe "kubectl get sc csi-hostpath-sc"

pe "kubectl create -f pvc.yaml"
pe "kubectl get pvc"

pe "vi snapshotclass.yaml"
pe "kubectl create -f snapshotclass.yaml"
pe "kubectl describe volumesnapshotclass"

pe "vi snapshot.yaml"
pe "kubectl create -f snapshot.yaml"

pe "kubectl describe volumesnapshot"

pe "kubectl describe volumesnapshotcontent"

pe "vi restore.yaml"
pe "kubectl create -f restore.yaml"
pe "kubectl get pvc"

echo "==> Test pre-provisioned snapshots"
pe "vi content.yaml"
pe "kubectl create -f content.yaml"
pe "vi pre-provisioned-snapshot.yaml"
pe "kubectl create -f pre-provisioned-snapshot.yaml"
