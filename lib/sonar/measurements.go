package sonar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//MeasurementResponse is the payload layout from SonarCloud
type MeasurementResponse struct {
	Component struct {
		ID          string `json:"id"`
		Key         string `json:"key"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Qualifier   string `json:"qualifier"`
		Measures    []struct {
			Metric    string `json:"metric"`
			Value     string `json:"value"`
			BestValue bool   `json:"bestValue,omitempty"`
		} `json:"measures"`
	} `json:"component"`
}

//ProjectMeasurements retrieves the measurements we pass through about the projects.
func (c *Client) ProjectMeasurements(key string) (*MeasurementResponse, error) {

	url := fmt.Sprintf("%s/measures/component?metricKeys=ncloc,coverage,vulnerabilities,bugs,violations&component=%s", c.sonarConnectionString, key)

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
	var measurements MeasurementResponse

	err = json.Unmarshal(body, &measurements)
	if err != nil {
		return nil, err
	}

	return &measurements, nil
}
