package dbt_cloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type Response struct {
	Data  []any `json:"data"`
	Extra Extra `json:"extra"`
}

type RawResponse struct {
	Data  []json.RawMessage `json:"data"`
	Extra Extra             `json:"extra"`
}

type Extra struct {
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Count      int `json:"count"`
	TotalCount int `json:"total_count"`
}

var log = logrus.New()

func (c *Client) GetEndpoint(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating a new request: %v", err)
	}

	resp, err := c.doRequest(req)
	if err != nil {
		log.Fatalf("Error fetching URL %v: %v", url, err)
	}

	return resp, err
}

func (c *Client) GetRawData(url string) ([]json.RawMessage, error) {

	// get the first page
	jsonPayload, err := c.GetEndpoint(url)
	if err != nil {
		return nil, err
	}

	var response RawResponse

	err = json.Unmarshal(jsonPayload, &response)
	if err != nil {
		return nil, err
	}

	allResponses := response.Data

	count := response.Extra.Pagination.Count
	for count < response.Extra.Pagination.TotalCount {
		// get the next page

		var newURL string
		lastPartURL, _ := lo.Last(strings.Split(url, "/"))
		if strings.Contains(lastPartURL, "?") {
			newURL = fmt.Sprintf("%s&offset=%d", url, count)
		} else {
			newURL = fmt.Sprintf("%s?offset=%d", url, count)
		}

		jsonPayload, err := c.GetEndpoint(newURL)
		if err != nil {
			return nil, err
		}
		var response RawResponse

		err = json.Unmarshal(jsonPayload, &response)
		if err != nil {
			return nil, err
		}

		if response.Extra.Pagination.Count == 0 {
			// Unlucky! one object might have been deleted since the first call
			// if we don't stop here we will loop forever!
			break
		} else {
			count += response.Extra.Pagination.Count
		}
		allResponses = append(allResponses, response.Data...)
	}

	return allResponses, nil
}

func (c *Client) GetData(url string) []any {
	rawData, err := c.GetRawData(url)
	if err != nil {
		log.Fatal(err)
	}

	allData := make([]any, len(rawData))
	for i, data := range rawData {
		err := json.Unmarshal(data, &allData[i])
		if err != nil {
			log.Fatal(err)
		}
	}
	return allData
}

func (c *Client) GetAllGroupIDsByName(groupName string) []int {
	url := fmt.Sprintf("%s/v3/accounts/%d/groups/", c.HostURL, c.AccountID)

	allGroupsRaw := c.GetData(url)

	return lo.FilterMap(allGroupsRaw, func(group any, _ int) (int, bool) {
		if group.(map[string]any)["name"].(string) == groupName {
			return int(group.(map[string]any)["id"].(float64)), true
		}
		return 0, false
	})
}

func (c *Client) GetAllEnvironments(projectID int) ([]Environment, error) {
	url := fmt.Sprintf("%s/v3/accounts/%d/environments/", c.HostURL, c.AccountID)

	if projectID != 0 {
		url = fmt.Sprintf("%s?project_id=%d", url, projectID)
	}

	allEnvironmentsRaw := c.GetData(url)

	allEnvs := []Environment{}
	for _, env := range allEnvironmentsRaw {

		data, _ := json.Marshal(env)
		currentEnv := Environment{}
		err := json.Unmarshal(data, &currentEnv)
		if err != nil {
			return nil, err
		}
		allEnvs = append(allEnvs, currentEnv)
	}
	return allEnvs, nil
}

func (c *Client) GetAllNotifications() ([]Notification, error) {
	url := fmt.Sprintf("%s/v2/accounts/%d/notifications/", c.HostURL, c.AccountID)

	allNotificationsRaw := c.GetData(url)

	allNotifications := []Notification{}
	for _, notification := range allNotificationsRaw {

		data, _ := json.Marshal(notification)
		currentNotification := Notification{}
		err := json.Unmarshal(data, &currentNotification)
		if err != nil {
			return nil, err
		}
		allNotifications = append(allNotifications, currentNotification)
	}
	return allNotifications, nil
}

func (c *Client) GetAllServiceTokens() ([]ServiceToken, error) {
	url := fmt.Sprintf("%s/v3/accounts/%d/service-tokens/?state=1", c.HostURL, c.AccountID)

	allServiceTokensRaw := c.GetData(url)

	allServiceTokens := []ServiceToken{}
	for _, notification := range allServiceTokensRaw {

		data, _ := json.Marshal(notification)
		currentServiceToken := ServiceToken{}
		err := json.Unmarshal(data, &currentServiceToken)
		if err != nil {
			return nil, err
		}
		allServiceTokens = append(allServiceTokens, currentServiceToken)
	}
	return allServiceTokens, nil
}

func (c *Client) GetAllLicenseMaps() ([]LicenseMap, error) {
	url := fmt.Sprintf("%s/v3/accounts/%d/license-maps/", c.HostURL, c.AccountID)

	allLicenseMapsRaw := c.GetData(url)

	allLicenseMaps := []LicenseMap{}
	for _, notification := range allLicenseMapsRaw {

		data, _ := json.Marshal(notification)
		currentLicenseMap := LicenseMap{}
		err := json.Unmarshal(data, &currentLicenseMap)
		if err != nil {
			return nil, err
		}
		allLicenseMaps = append(allLicenseMaps, currentLicenseMap)
	}
	return allLicenseMaps, nil
}

func (c *Client) GetAllJobs(projectID int, environmentID int) ([]JobWithEnvironment, error) {
	var url string

	if projectID != 0 && environmentID != 0 {
		return nil, fmt.Errorf("you can't filter by both project and environment")
	}

	if projectID == 0 && environmentID == 0 {
		return nil, fmt.Errorf("you must filter by either project or environment")
	}

	if projectID != 0 {
		url = fmt.Sprintf(
			"%s/v2/accounts/%d/jobs?project_id=%d&include_related=[environment]",
			c.HostURL,
			c.AccountID,
			projectID,
		)
	}

	if environmentID != 0 {
		url = fmt.Sprintf(
			"%s/v2/accounts/%d/jobs?environment_id=%d&include_related=[environment]",
			c.HostURL,
			c.AccountID,
			environmentID,
		)
	}

	allJobsRaw := c.GetData(url)

	allJobs := []JobWithEnvironment{}
	for _, job := range allJobsRaw {

		data, _ := json.Marshal(job)
		currentJob := JobWithEnvironment{}
		err := json.Unmarshal(data, &currentJob)
		if err != nil {
			return nil, err
		}
		allJobs = append(allJobs, currentJob)
	}
	return allJobs, nil
}
