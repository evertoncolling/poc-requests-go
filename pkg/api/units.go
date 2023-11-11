package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/evertoncolling/poc-requests-go/pkg/dto"
)

func (u *Units) List() (dto.UnitList, error) {
	endpoint := fmt.Sprintf("/api/v1/projects/%s/units", u.Client.ClientConfig.Project)
	url := u.Client.BaseURL + endpoint

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return dto.UnitList{}, err
	}

	for key, value := range u.Client.Headers {
		req.Header.Set(key, value)
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return dto.UnitList{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dto.UnitList{}, fmt.Errorf("failed to fetch units: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.UnitList{}, err
	}

	var unitList dto.UnitList
	if err := json.Unmarshal(body, &unitList); err != nil {
		return dto.UnitList{}, err
	}

	return unitList, nil
}
