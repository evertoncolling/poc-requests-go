package dto

type ViewReference struct {
	Type       string `json:"type"`
	Space      string `json:"space"`
	ExternalId string `json:"externalId"`
	Version    string `json:"version"`
}

type DataModelItem struct {
	Space           string          `json:"space"`
	ExternalId      string          `json:"externalId"`
	Name            string          `json:"name,omitempty"`
	Description     string          `json:"description,omitempty"`
	Version         string          `json:"version"`
	Views           []ViewReference `json:"views"`
	CreatedTime     int64           `json:"createdTime"`
	LastUpdatedTime int64           `json:"lastUpdatedTime"`
	IsGlobal        bool            `json:"isGlobal"`
}

type DataModelList struct {
	Items      []DataModelItem `json:"items"`
	NextCursor *string         `json:"nextCursor,omitempty"`
}

type NodeDefinition struct {
	InstanceType    string                 `json:"instanceType"`
	Version         int64                  `json:"version"`
	Space           string                 `json:"space"`
	ExternalId      string                 `json:"externalId"`
	Type            *InstanceId            `json:"type,omitempty"`
	CreatedTime     int64                  `json:"createdTime"`
	LastUpdatedTime int64                  `json:"lastUpdatedTime"`
	DeletedTime     *int64                 `json:"deletedTime,omitempty"`
	Properties      map[string]interface{} `json:"properties"`
}

type NodeList struct {
	Items  []NodeDefinition       `json:"items"`
	Typing map[string]interface{} `json:"typing"`
}

type UnitReferenceDM struct {
	ExternalId     *string `json:"externalId,omitempty"`
	UnitSystemName *string `json:"unitSystemName,omitempty"`
}

type TargetUnitsDM struct {
	Property string          `json:"property"`
	Unit     UnitReferenceDM `json:"unit"`
}

type SearchSort struct {
	Property  string  `json:"property"`
	Direction *string `json:"direction,omitempty"`
}
