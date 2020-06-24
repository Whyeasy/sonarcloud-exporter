package sonar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

//ProjectResponse is the payload layout from SonarCloud
type ProjectResponse struct {
	Paging struct {
		PageIndex int `json:"pageIndex"`
		PageSize  int `json:"pageSize"`
		Total     int `json:"total"`
	} `json:"paging"`
	Components []struct {
		Organization     string `json:"organization"`
		Key              string `json:"key"`
		Name             string `json:"name"`
		Qualifier        string `json:"qualifier"`
		Visibility       string `json:"visibility"`
		LastAnalysisDate string `json:"lastAnalysisDate"`
		Revision         string `json:"revision"`
	} `json:"components"`
}

//ListProjects lists the current projects in the organization.
func (c *Client) ListProjects(opt *ListOptions) (*ProjectResponse, error) {

	url := fmt.Sprintf("%s/projects/search?organization=%s&p=%s", c.sonarConnectionString, c.organization, strconv.Itoa(opt.Page))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var projects ProjectResponse

	err = json.Unmarshal(body, &projects)
	if err != nil {
		return nil, err
	}

	return &projects, nil
}
