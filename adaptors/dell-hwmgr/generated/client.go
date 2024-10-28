// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/oapi-codegen/runtime"
)

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetTokenWithBody request with any body
	GetTokenWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	GetToken(ctx context.Context, body GetTokenJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// VerifyRequestStatusWithBody request with any body
	VerifyRequestStatusWithBody(ctx context.Context, jobid string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	VerifyRequestStatus(ctx context.Context, jobid string, body VerifyRequestStatusJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateResourceGroupWithBody request with any body
	CreateResourceGroupWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateResourceGroup(ctx context.Context, body CreateResourceGroupJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteResourceGroupWithBody request with any body
	DeleteResourceGroupWithBody(ctx context.Context, resourceGroupId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	DeleteResourceGroup(ctx context.Context, resourceGroupId string, body DeleteResourceGroupJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetResourceGroupWithBody request with any body
	GetResourceGroupWithBody(ctx context.Context, resourceGroupId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	GetResourceGroup(ctx context.Context, resourceGroupId string, body GetResourceGroupJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateResourcePoolWithBody request with any body
	CreateResourcePoolWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateResourcePool(ctx context.Context, body CreateResourcePoolJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// OnboardServerWithBody request with any body
	OnboardServerWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	OnboardServer(ctx context.Context, body OnboardServerJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreatesecretsWithBody request with any body
	CreatesecretsWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	Createsecrets(ctx context.Context, body CreatesecretsJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetTokenWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetTokenRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetToken(ctx context.Context, body GetTokenJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetTokenRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) VerifyRequestStatusWithBody(ctx context.Context, jobid string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewVerifyRequestStatusRequestWithBody(c.Server, jobid, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) VerifyRequestStatus(ctx context.Context, jobid string, body VerifyRequestStatusJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewVerifyRequestStatusRequest(c.Server, jobid, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateResourceGroupWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateResourceGroupRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateResourceGroup(ctx context.Context, body CreateResourceGroupJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateResourceGroupRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteResourceGroupWithBody(ctx context.Context, resourceGroupId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteResourceGroupRequestWithBody(c.Server, resourceGroupId, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteResourceGroup(ctx context.Context, resourceGroupId string, body DeleteResourceGroupJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteResourceGroupRequest(c.Server, resourceGroupId, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetResourceGroupWithBody(ctx context.Context, resourceGroupId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetResourceGroupRequestWithBody(c.Server, resourceGroupId, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetResourceGroup(ctx context.Context, resourceGroupId string, body GetResourceGroupJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetResourceGroupRequest(c.Server, resourceGroupId, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateResourcePoolWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateResourcePoolRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateResourcePool(ctx context.Context, body CreateResourcePoolJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateResourcePoolRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) OnboardServerWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewOnboardServerRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) OnboardServer(ctx context.Context, body OnboardServerJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewOnboardServerRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreatesecretsWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreatesecretsRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) Createsecrets(ctx context.Context, body CreatesecretsJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreatesecretsRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetTokenRequest calls the generic GetToken builder with application/json body
func NewGetTokenRequest(server string, body GetTokenJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewGetTokenRequestWithBody(server, "application/json", bodyReader)
}

// NewGetTokenRequestWithBody generates requests for GetToken with any type of body
func NewGetTokenRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/identity/v1/tenant/Fulcrum/token/create")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewVerifyRequestStatusRequest calls the generic VerifyRequestStatus builder with application/json body
func NewVerifyRequestStatusRequest(server string, jobid string, body VerifyRequestStatusJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewVerifyRequestStatusRequestWithBody(server, jobid, "application/json", bodyReader)
}

// NewVerifyRequestStatusRequestWithBody generates requests for VerifyRequestStatus with any type of body
func NewVerifyRequestStatusRequestWithBody(server string, jobid string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "jobid", runtime.ParamLocationPath, jobid)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/tenants/default_tenant/jobs/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewCreateResourceGroupRequest calls the generic CreateResourceGroup builder with application/json body
func NewCreateResourceGroupRequest(server string, body CreateResourceGroupJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateResourceGroupRequestWithBody(server, "application/json", bodyReader)
}

// NewCreateResourceGroupRequestWithBody generates requests for CreateResourceGroup with any type of body
func NewCreateResourceGroupRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/tenants/default_tenant/resourcegroups")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewDeleteResourceGroupRequest calls the generic DeleteResourceGroup builder with application/json body
func NewDeleteResourceGroupRequest(server string, resourceGroupId string, body DeleteResourceGroupJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewDeleteResourceGroupRequestWithBody(server, resourceGroupId, "application/json", bodyReader)
}

// NewDeleteResourceGroupRequestWithBody generates requests for DeleteResourceGroup with any type of body
func NewDeleteResourceGroupRequestWithBody(server string, resourceGroupId string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "resource-group-id", runtime.ParamLocationPath, resourceGroupId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/tenants/default_tenant/resourcegroups/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewGetResourceGroupRequest calls the generic GetResourceGroup builder with application/json body
func NewGetResourceGroupRequest(server string, resourceGroupId string, body GetResourceGroupJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewGetResourceGroupRequestWithBody(server, resourceGroupId, "application/json", bodyReader)
}

// NewGetResourceGroupRequestWithBody generates requests for GetResourceGroup with any type of body
func NewGetResourceGroupRequestWithBody(server string, resourceGroupId string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "resource-group-id", runtime.ParamLocationPath, resourceGroupId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/tenants/default_tenant/resourcegroups/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewCreateResourcePoolRequest calls the generic CreateResourcePool builder with application/json body
func NewCreateResourcePoolRequest(server string, body CreateResourcePoolJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateResourcePoolRequestWithBody(server, "application/json", bodyReader)
}

// NewCreateResourcePoolRequestWithBody generates requests for CreateResourcePool with any type of body
func NewCreateResourcePoolRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/tenants/default_tenant/resourcepools")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewOnboardServerRequest calls the generic OnboardServer builder with application/json body
func NewOnboardServerRequest(server string, body OnboardServerJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewOnboardServerRequestWithBody(server, "application/json", bodyReader)
}

// NewOnboardServerRequestWithBody generates requests for OnboardServer with any type of body
func NewOnboardServerRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/tenants/default_tenant/resources")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewCreatesecretsRequest calls the generic Createsecrets builder with application/json body
func NewCreatesecretsRequest(server string, body CreatesecretsJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreatesecretsRequestWithBody(server, "application/json", bodyReader)
}

// NewCreatesecretsRequestWithBody generates requests for Createsecrets with any type of body
func NewCreatesecretsRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/tenants/default_tenant/secrets")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetTokenWithBodyWithResponse request with any body
	GetTokenWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*GetTokenResponse, error)

	GetTokenWithResponse(ctx context.Context, body GetTokenJSONRequestBody, reqEditors ...RequestEditorFn) (*GetTokenResponse, error)

	// VerifyRequestStatusWithBodyWithResponse request with any body
	VerifyRequestStatusWithBodyWithResponse(ctx context.Context, jobid string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*VerifyRequestStatusResponse, error)

	VerifyRequestStatusWithResponse(ctx context.Context, jobid string, body VerifyRequestStatusJSONRequestBody, reqEditors ...RequestEditorFn) (*VerifyRequestStatusResponse, error)

	// CreateResourceGroupWithBodyWithResponse request with any body
	CreateResourceGroupWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateResourceGroupResponse, error)

	CreateResourceGroupWithResponse(ctx context.Context, body CreateResourceGroupJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateResourceGroupResponse, error)

	// DeleteResourceGroupWithBodyWithResponse request with any body
	DeleteResourceGroupWithBodyWithResponse(ctx context.Context, resourceGroupId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*DeleteResourceGroupResponse, error)

	DeleteResourceGroupWithResponse(ctx context.Context, resourceGroupId string, body DeleteResourceGroupJSONRequestBody, reqEditors ...RequestEditorFn) (*DeleteResourceGroupResponse, error)

	// GetResourceGroupWithBodyWithResponse request with any body
	GetResourceGroupWithBodyWithResponse(ctx context.Context, resourceGroupId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*GetResourceGroupResponse, error)

	GetResourceGroupWithResponse(ctx context.Context, resourceGroupId string, body GetResourceGroupJSONRequestBody, reqEditors ...RequestEditorFn) (*GetResourceGroupResponse, error)

	// CreateResourcePoolWithBodyWithResponse request with any body
	CreateResourcePoolWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateResourcePoolResponse, error)

	CreateResourcePoolWithResponse(ctx context.Context, body CreateResourcePoolJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateResourcePoolResponse, error)

	// OnboardServerWithBodyWithResponse request with any body
	OnboardServerWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*OnboardServerResponse, error)

	OnboardServerWithResponse(ctx context.Context, body OnboardServerJSONRequestBody, reqEditors ...RequestEditorFn) (*OnboardServerResponse, error)

	// CreatesecretsWithBodyWithResponse request with any body
	CreatesecretsWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreatesecretsResponse, error)

	CreatesecretsWithResponse(ctx context.Context, body CreatesecretsJSONRequestBody, reqEditors ...RequestEditorFn) (*CreatesecretsResponse, error)
}

type GetTokenResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r GetTokenResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetTokenResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type VerifyRequestStatusResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r VerifyRequestStatusResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r VerifyRequestStatusResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateResourceGroupResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r CreateResourceGroupResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateResourceGroupResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteResourceGroupResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r DeleteResourceGroupResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteResourceGroupResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetResourceGroupResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r GetResourceGroupResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetResourceGroupResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateResourcePoolResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r CreateResourcePoolResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateResourcePoolResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type OnboardServerResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r OnboardServerResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r OnboardServerResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreatesecretsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r CreatesecretsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreatesecretsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetTokenWithBodyWithResponse request with arbitrary body returning *GetTokenResponse
func (c *ClientWithResponses) GetTokenWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*GetTokenResponse, error) {
	rsp, err := c.GetTokenWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetTokenResponse(rsp)
}

func (c *ClientWithResponses) GetTokenWithResponse(ctx context.Context, body GetTokenJSONRequestBody, reqEditors ...RequestEditorFn) (*GetTokenResponse, error) {
	rsp, err := c.GetToken(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetTokenResponse(rsp)
}

// VerifyRequestStatusWithBodyWithResponse request with arbitrary body returning *VerifyRequestStatusResponse
func (c *ClientWithResponses) VerifyRequestStatusWithBodyWithResponse(ctx context.Context, jobid string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*VerifyRequestStatusResponse, error) {
	rsp, err := c.VerifyRequestStatusWithBody(ctx, jobid, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseVerifyRequestStatusResponse(rsp)
}

func (c *ClientWithResponses) VerifyRequestStatusWithResponse(ctx context.Context, jobid string, body VerifyRequestStatusJSONRequestBody, reqEditors ...RequestEditorFn) (*VerifyRequestStatusResponse, error) {
	rsp, err := c.VerifyRequestStatus(ctx, jobid, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseVerifyRequestStatusResponse(rsp)
}

// CreateResourceGroupWithBodyWithResponse request with arbitrary body returning *CreateResourceGroupResponse
func (c *ClientWithResponses) CreateResourceGroupWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateResourceGroupResponse, error) {
	rsp, err := c.CreateResourceGroupWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateResourceGroupResponse(rsp)
}

func (c *ClientWithResponses) CreateResourceGroupWithResponse(ctx context.Context, body CreateResourceGroupJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateResourceGroupResponse, error) {
	rsp, err := c.CreateResourceGroup(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateResourceGroupResponse(rsp)
}

// DeleteResourceGroupWithBodyWithResponse request with arbitrary body returning *DeleteResourceGroupResponse
func (c *ClientWithResponses) DeleteResourceGroupWithBodyWithResponse(ctx context.Context, resourceGroupId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*DeleteResourceGroupResponse, error) {
	rsp, err := c.DeleteResourceGroupWithBody(ctx, resourceGroupId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteResourceGroupResponse(rsp)
}

func (c *ClientWithResponses) DeleteResourceGroupWithResponse(ctx context.Context, resourceGroupId string, body DeleteResourceGroupJSONRequestBody, reqEditors ...RequestEditorFn) (*DeleteResourceGroupResponse, error) {
	rsp, err := c.DeleteResourceGroup(ctx, resourceGroupId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteResourceGroupResponse(rsp)
}

// GetResourceGroupWithBodyWithResponse request with arbitrary body returning *GetResourceGroupResponse
func (c *ClientWithResponses) GetResourceGroupWithBodyWithResponse(ctx context.Context, resourceGroupId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*GetResourceGroupResponse, error) {
	rsp, err := c.GetResourceGroupWithBody(ctx, resourceGroupId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetResourceGroupResponse(rsp)
}

func (c *ClientWithResponses) GetResourceGroupWithResponse(ctx context.Context, resourceGroupId string, body GetResourceGroupJSONRequestBody, reqEditors ...RequestEditorFn) (*GetResourceGroupResponse, error) {
	rsp, err := c.GetResourceGroup(ctx, resourceGroupId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetResourceGroupResponse(rsp)
}

// CreateResourcePoolWithBodyWithResponse request with arbitrary body returning *CreateResourcePoolResponse
func (c *ClientWithResponses) CreateResourcePoolWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateResourcePoolResponse, error) {
	rsp, err := c.CreateResourcePoolWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateResourcePoolResponse(rsp)
}

func (c *ClientWithResponses) CreateResourcePoolWithResponse(ctx context.Context, body CreateResourcePoolJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateResourcePoolResponse, error) {
	rsp, err := c.CreateResourcePool(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateResourcePoolResponse(rsp)
}

// OnboardServerWithBodyWithResponse request with arbitrary body returning *OnboardServerResponse
func (c *ClientWithResponses) OnboardServerWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*OnboardServerResponse, error) {
	rsp, err := c.OnboardServerWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseOnboardServerResponse(rsp)
}

func (c *ClientWithResponses) OnboardServerWithResponse(ctx context.Context, body OnboardServerJSONRequestBody, reqEditors ...RequestEditorFn) (*OnboardServerResponse, error) {
	rsp, err := c.OnboardServer(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseOnboardServerResponse(rsp)
}

// CreatesecretsWithBodyWithResponse request with arbitrary body returning *CreatesecretsResponse
func (c *ClientWithResponses) CreatesecretsWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreatesecretsResponse, error) {
	rsp, err := c.CreatesecretsWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreatesecretsResponse(rsp)
}

func (c *ClientWithResponses) CreatesecretsWithResponse(ctx context.Context, body CreatesecretsJSONRequestBody, reqEditors ...RequestEditorFn) (*CreatesecretsResponse, error) {
	rsp, err := c.Createsecrets(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreatesecretsResponse(rsp)
}

// ParseGetTokenResponse parses an HTTP response from a GetTokenWithResponse call
func ParseGetTokenResponse(rsp *http.Response) (*GetTokenResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetTokenResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseVerifyRequestStatusResponse parses an HTTP response from a VerifyRequestStatusWithResponse call
func ParseVerifyRequestStatusResponse(rsp *http.Response) (*VerifyRequestStatusResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &VerifyRequestStatusResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseCreateResourceGroupResponse parses an HTTP response from a CreateResourceGroupWithResponse call
func ParseCreateResourceGroupResponse(rsp *http.Response) (*CreateResourceGroupResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateResourceGroupResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseDeleteResourceGroupResponse parses an HTTP response from a DeleteResourceGroupWithResponse call
func ParseDeleteResourceGroupResponse(rsp *http.Response) (*DeleteResourceGroupResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteResourceGroupResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseGetResourceGroupResponse parses an HTTP response from a GetResourceGroupWithResponse call
func ParseGetResourceGroupResponse(rsp *http.Response) (*GetResourceGroupResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetResourceGroupResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseCreateResourcePoolResponse parses an HTTP response from a CreateResourcePoolWithResponse call
func ParseCreateResourcePoolResponse(rsp *http.Response) (*CreateResourcePoolResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateResourcePoolResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseOnboardServerResponse parses an HTTP response from a OnboardServerWithResponse call
func ParseOnboardServerResponse(rsp *http.Response) (*OnboardServerResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &OnboardServerResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseCreatesecretsResponse parses an HTTP response from a CreatesecretsWithResponse call
func ParseCreatesecretsResponse(rsp *http.Response) (*CreatesecretsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreatesecretsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}
