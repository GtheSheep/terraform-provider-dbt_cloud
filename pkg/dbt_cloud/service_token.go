package dbt_cloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type ServiceTokenPermission struct {
	ID             *int                  `json:"id,omitempty"`
	AccountID      int                   `json:"account_id"`
	ServiceTokenID int                   `json:"service_token_id"`
	ProjectID      int                   `json:"project_id,omitempty"`
	AllProjects    bool                  `json:"all_projects"`
	State          int                   `json:"state,omitempty"`
	Set            string                `json:"permission_set,omitempty"`
	WritableEnvs   []EnvironmentCategory `json:"writable_environment_categories,omitempty"`
}

type ServiceToken struct {
	ID          *int                     `json:"id"`
	AccountID   int                      `json:"account_id"`
	UID         string                   `json:"uid"`
	Name        string                   `json:"name"`
	TokenString *string                  `json:"token_string,omitempty"`
	State       int                      `json:"state"`
	Permissions []ServiceTokenPermission `json:"service_token_permissions,omitempty"`
}

type ServiceTokenResponse struct {
	Data   ServiceToken   `json:"data"`
	Status ResponseStatus `json:"status"`
}

type ServiceTokenListResponse struct {
	Data   []ServiceToken `json:"data"`
	Status ResponseStatus `json:"status"`
}

type ServiceTokenPermissionListResponse struct {
	Data   []ServiceTokenPermission `json:"data"`
	Status ResponseStatus           `json:"status"`
}

func (c *Client) GetServiceTokenPermissions(serviceTokenID int) (*[]ServiceTokenPermission, error) {

	allServiceTokenPermissionsRaw, err := c.GetRawData(fmt.Sprintf("%s/v3/accounts/%s/service-tokens/%s/permissions/", c.HostURL, strconv.Itoa(c.AccountID), strconv.Itoa(serviceTokenID)))
	if err != nil {
		return nil, err
	}

	allPermissions := make([]ServiceTokenPermission, len(allServiceTokenPermissionsRaw))

	for i, permission := range allServiceTokenPermissionsRaw {
		err := json.Unmarshal(permission, &allPermissions[i])
		if err != nil {
			return nil, err
		}
	}

	return &allPermissions, nil
}

func (c *Client) GetServiceToken(serviceTokenID int) (*ServiceToken, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v3/accounts/%s/service-tokens/%s/", c.HostURL, strconv.Itoa(c.AccountID), strconv.Itoa(serviceTokenID)), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	serviceTokenResponse := ServiceTokenResponse{}
	err = json.Unmarshal(body, &serviceTokenResponse)
	if err != nil {
		return nil, err
	}

	// the endpoint returns service tokens when their state is inactive, so we need to check for the state
	if serviceTokenResponse.Data.State != STATE_ACTIVE {
		return nil, fmt.Errorf("resource-not-found: service token %d is not active", serviceTokenID)
	}

	permissions, err := c.GetServiceTokenPermissions(serviceTokenID)
	if err != nil {
		return nil, err
	}
	serviceTokenResponse.Data.Permissions = *permissions

	return &serviceTokenResponse.Data, nil
}

func (c *Client) CreateServiceToken(
	name string,
	state int,
) (*ServiceToken, error) {
	newServiceToken := ServiceToken{
		AccountID: c.AccountID,
		State:     state,
		Name:      name,
	}
	newServiceTokenData, err := json.Marshal(newServiceToken)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v3/accounts/%d/service-tokens/", c.HostURL, c.AccountID), strings.NewReader(string(newServiceTokenData)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	serviceTokenResponse := ServiceTokenResponse{}
	err = json.Unmarshal(body, &serviceTokenResponse)
	if err != nil {
		return nil, err
	}

	return &serviceTokenResponse.Data, nil
}

func (c *Client) UpdateServiceToken(serviceTokenID int, serviceToken ServiceToken) (*ServiceToken, error) {
	serviceTokenData, err := json.Marshal(serviceToken)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v3/accounts/%s/service-tokens/%d/", c.HostURL, strconv.Itoa(c.AccountID), serviceTokenID), strings.NewReader(string(serviceTokenData)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	serviceTokenResponse := ServiceTokenResponse{}
	err = json.Unmarshal(body, &serviceTokenResponse)
	if err != nil {
		return nil, err
	}

	return &serviceTokenResponse.Data, nil
}

func (c *Client) UpdateServiceTokenPermissions(serviceTokenID int, serviceTokenPermissions []ServiceTokenPermission) (*[]ServiceTokenPermission, error) {
	serviceTokenPermissionData, err := json.Marshal(serviceTokenPermissions)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v3/accounts/%s/service-tokens/%d/permissions/", c.HostURL, strconv.Itoa(c.AccountID), serviceTokenID), strings.NewReader(string(serviceTokenPermissionData)))
	if err != nil {
		return nil, err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return nil, err
	}

	return c.GetServiceTokenPermissions(serviceTokenID)
}

func (c *Client) DeleteServiceToken(serviceTokenID int) (string, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v3/accounts/%d/service-tokens/%d/", c.HostURL, c.AccountID, serviceTokenID), nil)
	if err != nil {
		return "", err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return "", err
	}

	return "", err
}
