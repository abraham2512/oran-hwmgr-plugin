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

package hwmgrclient

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
	hwmgrapi "github.com/openshift-kni/oran-hwmgr-plugin/adaptors/dell-hwmgr/generated"
	pluginv1alpha1 "github.com/openshift-kni/oran-hwmgr-plugin/api/hwmgr-plugin/v1alpha1"
	"github.com/openshift-kni/oran-hwmgr-plugin/internal/controller/utils"
	hwmgmtv1alpha1 "github.com/openshift-kni/oran-o2ims/api/hardwaremanagement/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// TODO: Define how these labels are specified. New field in NodePool?
// For now, we'll use an enum and label and hardcode a correlation to the nodegroup list order
const (
	ResourceSelectorControllerIdx = iota
	ResourceSelectorWorkerIdx

	ResourceSelectorControllerLabel = "controller" // TODO:
	ResourceSelectorWorkerLabel     = "worker"
)

const (
	RoleKey = "role"
	Tenant  = "default_tenant" // TODO: Hardcoded, for now
)

// HardwareManagerClient provides functions for calling the hardware manager APIs
type HardwareManagerClient struct {
	rtclient    client.Client
	HwmgrClient *hwmgrapi.ClientWithResponses
	Logger      *slog.Logger
	Namespace   string
	hwmgr       *pluginv1alpha1.HardwareManager
}

// GetToken sends a request to the hardware manager to request an authentication token
func (c *HardwareManagerClient) GetToken(ctx context.Context) (string, error) {
	clientSecrets, err := utils.GetSecret(ctx, c.rtclient, c.hwmgr.Spec.DellData.AuthSecret, c.Namespace)
	if err != nil {
		return "", fmt.Errorf("failed to get client secret: %w", err)
	}

	clientId, err := utils.GetSecretField(clientSecrets, "client-id")
	if err != nil {
		return "", fmt.Errorf("failed to get client-id from secret: %s, %w", c.hwmgr.Spec.DellData.AuthSecret, err)
	}

	username, err := utils.GetSecretField(clientSecrets, corev1.BasicAuthUsernameKey)
	if err != nil {
		return "", fmt.Errorf("failed to get %s from secret: %s, %w", corev1.BasicAuthUsernameKey, c.hwmgr.Spec.DellData.AuthSecret, err)
	}

	password, err := utils.GetSecretField(clientSecrets, corev1.BasicAuthPasswordKey)
	if err != nil {
		return "", fmt.Errorf("failed to get %s from secret: %s, %w", corev1.BasicAuthPasswordKey, c.hwmgr.Spec.DellData.AuthSecret, err)
	}

	grant_type := string(pluginv1alpha1.OAuthGrantTypes.Password)

	req := hwmgrapi.GetTokenJSONRequestBody{
		ClientId:  &clientId,
		Username:  &username,
		Password:  &password,
		GrantType: &grant_type,
	}

	tokenrsp, err := c.HwmgrClient.GetTokenWithResponse(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to get token: response: %v, err: %w", tokenrsp, err)
	}

	if tokenrsp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("token request failed with status %s (%d), message=%s",
			tokenrsp.Status(), tokenrsp.StatusCode(), string(tokenrsp.Body))
	}

	var tokenData hwmgrapi.RhprotoGetTokenResponseBody
	if err := json.Unmarshal(tokenrsp.Body, &tokenData); err != nil {
		return "", fmt.Errorf("failed to parse token: response: %v, err: %w", tokenrsp, err)
	}

	if tokenData.AccessToken == nil {
		return "", fmt.Errorf("failed to get token: access_token field empty: %v", tokenrsp)
	}
	return *tokenData.AccessToken, nil
}

// NewClientWithResponses creates an authenticated client connected to the hardware manager
func NewClientWithResponses(
	ctx context.Context,
	logger *slog.Logger,
	rtclient client.Client,
	hwmgr *pluginv1alpha1.HardwareManager) (*HardwareManagerClient, error) {

	hwmgrClient := HardwareManagerClient{
		rtclient:  rtclient,
		Logger:    logger,
		Namespace: hwmgr.Namespace,
		hwmgr:     hwmgr,
	}

	// If the HardwareManager CR includes certificates, get the bundle to add to the client
	var caBundle string
	if hwmgr.Spec.DellData.CaBundleName != nil {
		cm, err := utils.GetConfigmap(ctx, rtclient, *hwmgr.Spec.DellData.CaBundleName, hwmgr.Namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to get configmap: %s", err.Error())
		}

		caBundle, err = utils.GetConfigMapField(cm, "ca-bundle.pem")
		if err != nil {
			return nil, fmt.Errorf("failed to get certificate bundle from configmap: %s", err.Error())
		}
	}

	config := utils.OAuthClientConfig{
		CaBundle: []byte(caBundle),
	}

	tr, err := utils.GetTransportWithCaBundle(config)
	if err != nil {
		return nil, fmt.Errorf("failed to get http transport: %w", err)
	}

	httpClient := &http.Client{Transport: tr}

	// Create the hwmgrapi client, along with a bearer token
	hwmgrClient.HwmgrClient, err = hwmgrapi.NewClientWithResponses(
		hwmgr.Spec.DellData.ApiUrl,
		hwmgrapi.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("failed to setup client to %s: %w", hwmgr.Spec.DellData.ApiUrl, err)
	}

	token, err := hwmgrClient.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token for %s: %w", hwmgr.Name, err)
	}

	bearerAuth, err := securityprovider.NewSecurityProviderBearerToken(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create security provider for %s: %w", hwmgr.Name, err)
	}

	// Create a new client with an intercept to add the bearer token
	hwmgrClient.HwmgrClient, err = hwmgrapi.NewClientWithResponses(
		hwmgr.Spec.DellData.ApiUrl,
		hwmgrapi.WithHTTPClient(httpClient),
		hwmgrapi.WithRequestEditorFn(bearerAuth.Intercept))
	if err != nil {
		return nil, fmt.Errorf("failed to setup auth client for %s: %w", hwmgr.Name, err)
	}

	return &hwmgrClient, nil
}

// GetResourceGroup queries the hardware manager to get the resource group data
func (c *HardwareManagerClient) GetResourceGroup(ctx context.Context, nodepool *hwmgmtv1alpha1.NodePool) (*hwmgrapi.RhprotoResourceGroupObjectGetResponseBody, error) {
	rg := ResourceGroupFromNodePool(nodepool)
	rgId := *rg.ResourceGroup.Id

	response, err := c.HwmgrClient.GetResourceGroupWithResponse(ctx, rgId)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource group %s: response: %v, err: %w", rgId, response, err)
	}

	if response.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("resource group get failed with status %s (%d), message=%s",
			response.Status(), response.StatusCode(), string(response.Body))
	}

	return response.JSON200, nil
}

// ResourceGroupIdFromNodePool returns the resource group identifier corresponding to the specified nodepool
func ResourceGroupIdFromNodePool(nodepool *hwmgmtv1alpha1.NodePool) string {
	return fmt.Sprintf("rhplugin-rg-%s", nodepool.Spec.CloudID)
}

// ResourceGroupFromNodePool transforms data from a nodepool object to a CreateResourceGroupJSONRequestBody instance
func ResourceGroupFromNodePool(nodepool *hwmgmtv1alpha1.NodePool) *hwmgrapi.CreateResourceGroupJSONRequestBody {
	rgId := ResourceGroupIdFromNodePool(nodepool)
	tenant := Tenant
	// TODO: Get these values appropriately
	resourceTypeId := "ResourceGroup~2.1.1"
	description := "Resource Group managed by O-Cloud Hardware Manager Plugin"
	excludes := make(map[string]interface{})
	roleKey := RoleKey
	// TODO: What should the roles be?
	controllerRole := ResourceSelectorControllerLabel
	workerRole := ResourceSelectorWorkerLabel

	rg := hwmgrapi.CreateResourceGroupJSONRequestBody{
		Tenant: &tenant,
		ResourceGroup: &hwmgrapi.RhprotoResourceGroupObjectRequest{
			Description:    &description,
			Id:             &rgId,
			Name:           &rgId,
			ResourceTypeId: &resourceTypeId,
			ResourceSelectors: &map[string]hwmgrapi.RhprotoResourceSelectorRequest{
				ResourceSelectorControllerLabel: {
					RpId:              &nodepool.Spec.NodeGroup[ResourceSelectorControllerIdx].Name,
					ResourceProfileId: &nodepool.Spec.NodeGroup[ResourceSelectorControllerIdx].HwProfile,
					NumResources:      &nodepool.Spec.NodeGroup[ResourceSelectorControllerIdx].Size,
					Filters: &hwmgrapi.RhprotoResourceSelectorFilter{
						Include: &hwmgrapi.RhprotoResourceSelectorFilterInclude{
							Labels: &[]hwmgrapi.RhprotoResourceSelectorFilterIncludeLabel{
								{
									Key:   &roleKey,
									Value: &controllerRole,
								},
							},
						},
						Exclude: &excludes,
					},
				},
				ResourceSelectorWorkerLabel: {
					RpId:              &nodepool.Spec.NodeGroup[ResourceSelectorWorkerIdx].Name,
					ResourceProfileId: &nodepool.Spec.NodeGroup[ResourceSelectorWorkerIdx].HwProfile,
					NumResources:      &nodepool.Spec.NodeGroup[ResourceSelectorWorkerIdx].Size,
					Filters: &hwmgrapi.RhprotoResourceSelectorFilter{
						Include: &hwmgrapi.RhprotoResourceSelectorFilterInclude{
							Labels: &[]hwmgrapi.RhprotoResourceSelectorFilterIncludeLabel{
								{
									Key:   &roleKey,
									Value: &workerRole,
								},
							},
						},
						Exclude: &excludes,
					},
				},
			},
		},
	}

	return &rg
}

// CreateResourceGroup sends a request to the hardware manager, returns a jobId
// TODO: Improve error handling for different status codes
func (c *HardwareManagerClient) CreateResourceGroup(ctx context.Context, nodepool *hwmgmtv1alpha1.NodePool) (string, error) {
	rg := ResourceGroupFromNodePool(nodepool)
	rgId := *rg.ResourceGroup.Id

	// First check whether the resource group already exists
	response, err := c.HwmgrClient.GetResourceGroupWithResponse(ctx, rgId)
	if err != nil {
		return "", fmt.Errorf("failed to query for resource group %s: response: %v, err: %w", rgId, response, err)
	}

	if response.StatusCode() == http.StatusOK {
		return "", fmt.Errorf("resource group %s already exists", rgId)
	}

	// Send a request to the hardware manager to create the resource group
	rgResponse, err := c.HwmgrClient.CreateResourceGroupWithResponse(ctx, *rg)
	if err != nil {
		return "", fmt.Errorf("failed to create resource group %s, api failure: response: %v, err: %w", rgId, response, err)
	}

	if rgResponse.StatusCode() != http.StatusOK {
		// TODO: Remove this log
		c.Logger.InfoContext(ctx, "Failure from CreateResourceGroupWithResponse", slog.String("message", *rgResponse.JSONDefault.Message), slog.Any("response", rgResponse.JSONDefault))
		return "", fmt.Errorf("failed to create resource group %s, bad status: %s, code: %d, response: %v", rgId, rgResponse.Status(), rgResponse.StatusCode(), rgResponse)
	}

	// Return the job ID for the request
	return *rgResponse.JSON200.Jobid, nil
}

// CheckResourceGroupRequest queries the hardware manager for the status of a job
func (c *HardwareManagerClient) CheckResourceGroupRequest(ctx context.Context, jobId string) (*hwmgrapi.RhprotoJobStatus, error) {
	response, err := c.HwmgrClient.VerifyRequestStatusWithResponse(ctx, jobId)
	if err != nil {
		return nil, fmt.Errorf("failed to query for resource group job status: id: %s, response: %v, err: %w", jobId, response, err)
	}

	if response.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("job query failed for %s: %s", jobId, *response.JSONDefault.Message)
	}

	return response.JSON200, nil
}

// DeleteResourceGroup asks the hardware manager to delete the resource group associated with the specified nodepool
func (c *HardwareManagerClient) DeleteResourceGroup(ctx context.Context, nodepool *hwmgmtv1alpha1.NodePool) (string, error) {
	rgId := ResourceGroupIdFromNodePool(nodepool)

	response, err := c.HwmgrClient.DeleteResourceGroupWithResponse(ctx, rgId)
	if err != nil || response.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("failed to delete resource group %s: response: %v, err: %w", rgId, response, err)
	}

	return *response.JSON200.Jobid, nil
}
