package cmp

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func service(newObj, oldObj *v1.Service) bool {
	return false
}

func node(newObj, oldObj *v1.Node) bool {

	klog.Info("comparing node %s", newObj.Name)

	newNodeAddress := newObj.Status.Addresses
	oldNodeAddress := oldObj.Status.Addresses
	if len(newNodeAddress) != len(oldNodeAddress) {
		return false
	}

	mp := make(map[v1.NodeAddressType]string)
	for _, addr := range oldNodeAddress {
		mp[addr.Type] = addr.Address
	}
	for _, addr := range newNodeAddress {
		if mp[addr.Type] != addr.Address {
			return false
		}
	}
	return true
}

func Equal(newObj, oldObj client.Object) bool {

	klog.Infof("*****************************************")
	klog.Infof("checking equality for %s", newObj.GetName())

	if newObj.GetObjectKind().GroupVersionKind() != oldObj.GetObjectKind().GroupVersionKind() {
		return false
	}

	switch newObj.GetObjectKind().GroupVersionKind().Kind {
	case "Service":
		service(newObj.(*v1.Service).DeepCopy(), oldObj.(*v1.Service).DeepCopy())
	case "Node":
		node(newObj.(*v1.Node).DeepCopy(), oldObj.(*v1.Node).DeepCopy())
	default:
		klog.Infof("unknown object kind")
	}

	return false
}
