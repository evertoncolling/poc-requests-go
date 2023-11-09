package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"poc-requests-go/pkg/dto"
)

func ListUnits(
	project string,
	token string,
	baseURL string,
) (dto.UnitList, error) {
	endpoint := fmt.Sprintf("/api/v1/projects/%s/units", project)
	url := baseURL + endpoint

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return dto.UnitList{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("cdf-version", "beta")

	client := &http.Client{}
	resp, err := client.Do(req)
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
