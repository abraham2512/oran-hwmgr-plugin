/*
Copyright 2024.

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

package loopback

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/openshift-kni/oran-hwmgr-plugin/internal/controller/utils"
	hwmgmtv1alpha1 "github.com/openshift-kni/oran-o2ims/api/hardwaremanagement/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

// AllocateNode processes a NodePool CR, allocating a free node for each specified nodegroup as needed
func (a *LoopbackAdaptor) AllocateNode(ctx context.Context, nodepool *hwmgmtv1alpha1.NodePool) error {
	cloudID := nodepool.Spec.CloudID

	// Inject a delay before allocating node
	time.Sleep(10 * time.Second)

	cm, resources, allocations, err := a.GetCurrentResources(ctx)
	if err != nil {
		return fmt.Errorf("unable to get current resources: %w", err)
	}

	var cloud *cmAllocatedCloud
	for i, iter := range allocations.Clouds {
		if iter.CloudID == cloudID {
			cloud = &allocations.Clouds[i]
			break
		}
	}
	if cloud == nil {
		// The cloud wasn't found in the list, so create a new entry
		allocations.Clouds = append(allocations.Clouds, cmAllocatedCloud{CloudID: cloudID, Nodegroups: make(map[string][]string)})
		cloud = &allocations.Clouds[len(allocations.Clouds)-1]
	}

	// Check available resources
	for _, nodegroup := range nodepool.Spec.NodeGroup {
		used := cloud.Nodegroups[nodegroup.Name]
		remaining := nodegroup.Size - len(used)
		if remaining <= 0 {
			// This group is allocated
			a.logger.InfoContext(ctx, "nodegroup is fully allocated", "nodegroup", nodegroup.Name)
			continue
		}

		freenodes := getFreeNodesInProfile(resources, allocations, nodegroup.HwProfile)
		if remaining > len(freenodes) {
			return fmt.Errorf("not enough free resources remaining in group %s", nodegroup.HwProfile)
		}

		// Grab the first node
		nodename := freenodes[0]

		nodeinfo, exists := resources.Nodes[nodename]
		if !exists {
			return fmt.Errorf("unable to find nodeinfo for %s", nodename)
		}

		if err := a.CreateBMCSecret(ctx, nodename, nodeinfo.BMC.UsernameBase64, nodeinfo.BMC.PasswordBase64); err != nil {
			return fmt.Errorf("failed to create bmc-secret when allocating node %s: %w", nodename, err)
		}

		cloud.Nodegroups[nodegroup.Name] = append(cloud.Nodegroups[nodegroup.Name], nodename)

		// Update the configmap
		yamlString, err := yaml.Marshal(&allocations)
		if err != nil {
			return fmt.Errorf("unable to marshal allocated data: %w", err)
		}
		cm.Data[allocationsKey] = string(yamlString)
		if err := a.Client.Update(ctx, cm); err != nil {
			return fmt.Errorf("failed to update configmap: %w", err)
		}

		if err := a.CreateNode(ctx, cloudID, nodename, nodegroup.Name, nodegroup.HwProfile); err != nil {
			return fmt.Errorf("failed to create allocated node (%s): %w", nodename, err)
		}

		if err := a.UpdateNodeStatus(ctx, nodename, nodeinfo); err != nil {
			return fmt.Errorf("failed to update node status (%s): %w", nodename, err)
		}
	}

	return nil
}

func bmcSecretName(nodename string) string {
	return fmt.Sprintf("%s-bmc-secret", nodename)
}

// CreateBMCSecret creates the bmc-secret for a node
func (a *LoopbackAdaptor) CreateBMCSecret(ctx context.Context, nodename, usernameBase64, passwordBase64 string) error {
	a.logger.InfoContext(ctx, "Creating bmc-secret:", "nodename", nodename)

	secretName := bmcSecretName(nodename)

	username, err := base64.StdEncoding.DecodeString(usernameBase64)
	if err != nil {
		return fmt.Errorf("failed to decode usernameBase64 string (%s) for node %s: %w", usernameBase64, nodename, err)
	}

	password, err := base64.StdEncoding.DecodeString(passwordBase64)
	if err != nil {
		return fmt.Errorf("failed to decode usernameBase64 string (%s) for node %s: %w", passwordBase64, nodename, err)
	}

	bmcSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: a.namespace,
		},
		Data: map[string][]byte{
			"username": username,
			"password": password,
		},
	}

	if err = utils.CreateK8sCR(ctx, a.Client, bmcSecret, nil, utils.UPDATE); err != nil {
		return fmt.Errorf("failed to create bmc-secret for node %s: %w", nodename, err)
	}

	return nil
}

// DeleteBMCSecret deletes the bmc-secret for a node
func (a *LoopbackAdaptor) DeleteBMCSecret(ctx context.Context, nodename string) error {
	a.logger.InfoContext(ctx, "Deleting bmc-secret:", "nodename", nodename)

	secretName := bmcSecretName(nodename)

	bmcSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: a.namespace,
		},
	}

	if err := a.Client.Delete(ctx, bmcSecret); client.IgnoreNotFound(err) != nil {
		return fmt.Errorf("failed to delete bmc-secret for node %s: %w", nodename, err)
	}

	return nil
}

// CreateNode creates a Node CR with specified attributes
func (a *LoopbackAdaptor) CreateNode(ctx context.Context, cloudID, nodename, groupname, hwprofile string) error {

	a.logger.InfoContext(ctx, "Creating node:",
		"cloudID", cloudID,
		"nodegroup name", groupname,
		"nodename", nodename,
	)

	node := &hwmgmtv1alpha1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nodename,
			Namespace: a.namespace,
		},
		Spec: hwmgmtv1alpha1.NodeSpec{
			NodePool:  cloudID,
			GroupName: groupname,
			HwProfile: hwprofile,
		},
	}

	if err := a.Client.Create(ctx, node); err != nil {
		return fmt.Errorf("failed to create Node: %w", err)
	}

	return nil
}

// UpdateNodeStatus updates a Node CR status field with additional node information from the nodelist configmap
func (a *LoopbackAdaptor) UpdateNodeStatus(ctx context.Context, nodename string, info cmNodeInfo) error {

	a.logger.InfoContext(ctx, "Updating node:",
		"nodename", nodename,
	)

	node := &hwmgmtv1alpha1.Node{}

	if err := a.Client.Get(ctx, types.NamespacedName{Name: nodename, Namespace: a.namespace}, node); err != nil {
		return fmt.Errorf("failed to create Node: %w", err)
	}

	a.logger.InfoContext(ctx, "Adding info to node", "nodename", nodename, "info", info)
	node.Status.BMC = &hwmgmtv1alpha1.BMC{
		Address:         info.BMC.Address,
		CredentialsName: bmcSecretName(nodename),
	}
	node.Status.Hostname = info.Hostname
	node.Status.Interfaces = info.Interfaces

	utils.SetStatusCondition(&node.Status.Conditions,
		hwmgmtv1alpha1.Provisioned,
		hwmgmtv1alpha1.Completed,
		metav1.ConditionTrue,
		"Provisioned")

	if err := utils.UpdateK8sCRStatus(ctx, a.Client, node); err != nil {
		return fmt.Errorf("failed to update status for node %s: %w", nodename, err)
	}

	return nil
}

// DeleteNode deletes a Node CR
func (a *LoopbackAdaptor) DeleteNode(ctx context.Context, nodename string) error {

	a.logger.InfoContext(ctx, "Deleting node:",
		"nodename", nodename,
	)

	node := &hwmgmtv1alpha1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nodename,
			Namespace: a.namespace,
		},
	}

	if err := a.Client.Delete(ctx, node); client.IgnoreNotFound(err) != nil {
		return fmt.Errorf("failed to delete Node: %w", err)
	}

	return nil
}
