package dto

type TimeSeries struct {
	Id                 int64      `json:"id"`
	ExternalId         string     `json:"externalId,omitempty"`
	InstanceId         InstanceId `json:"instanceId,omitempty"`
	Name               string     `json:"name,omitempty"`
	IsString           bool       `json:"isString"`
	Metadata           Metadata   `json:"metadata,omitempty"`
	Unit               string     `json:"unit,omitempty"`
	UnitExternalId     string     `json:"unitExternalId,omitempty"`
	AssetId            int64      `json:"assetId,omitempty"`
	IsStep             bool       `json:"isStep"`
	Description        string     `json:"description,omitempty"`
	SecurityCategories []int64    `json:"securityCategories,omitempty"`
	DataSetID          int64      `json:"dataSetId,omitempty"`
	CreatedTime        int64      `json:"createdTime"`
	LastUpdatedTime    int64      `json:"lastUpdatedTime"`
}

type TimeSeriesSortItem struct {
	Property string `json:"property"`
	Order    string `json:"order,omitempty"`
	Nulls    string `json:"nulls,omitempty"`
}

type TimeSeriesList struct {
	Items      []TimeSeries `json:"items"`
	NextCursor *string      `json:"nextCursor,omitempty"`
}

type TimeSeriesFilter struct {
	Name             string          `json:"name,omitempty"`
	Unit             string          `json:"unit,omitempty"`
	UnitExternalId   string          `json:"unitExternalId,omitempty"`
	UnitQuantity     string          `json:"unitQuantity,omitempty"`
	IsString         bool            `json:"isString,omitempty"`
	IsStep           bool            `json:"isStep,omitempty"`
	Metadata         Metadata        `json:"metadata,omitempty"`
	AssetIDs         []int64         `json:"assetIds,omitempty"`
	AssetExternalIDs []string        `json:"assetExternalIds,omitempty"`
	RootAssetIDs     []int64         `json:"rootAssetIds,omitempty"`
	AssetSubtreeIds  []Identity      `json:"assetSubtreeIds,omitempty"`
	DataSetIds       []Identity      `json:"dataSetIds,omitempty"`
	ExternalIdPrefix string          `json:"externalIdPrefix,omitempty"`
	CreatedTime      *TimestampRange `json:"createdTime,omitempty"`
	LastUpdatedTime  *TimestampRange `json:"lastUpdatedTime,omitempty"`
}

type DataPointsQueryItem struct {
	Id                   int64       `json:"id,omitempty"`
	ExternalId           string      `json:"externalId,omitempty"`
	InstanceId           *InstanceId `json:"instanceId,omitempty"`
	Start                string      `json:"start,omitempty"`
	End                  string      `json:"end,omitempty"`
	Limit                int64       `json:"limit,omitempty"`
	Aggregates           []string    `json:"aggregates,omitempty"`
	Granularity          string      `json:"granularity,omitempty"`
	TargetUnit           string      `json:"targetUnit,omitempty"`
	TargetUnitSystem     string      `json:"targetUnitSystem,omitempty"`
	IncludeOutsidePoints bool        `json:"includeOutsidePoints,omitempty"`
	IncludeStatus        bool        `json:"includeStatus,omitempty"`
	IgnoreBadDataPoints  bool        `json:"ignoreBadDataPoints,omitempty"`
	TreatUncertainAsBad  bool        `json:"treatUncertainAsBad,omitempty"`
	TimeZone             string      `json:"timeZone,omitempty"`
	Cursor               string      `json:"cursor,omitempty"`
}

type LatestDataPointsQueryItem struct {
	Id                  int64       `json:"id,omitempty"`
	ExternalId          string      `json:"externalId,omitempty"`
	InstanceId          *InstanceId `json:"instanceId,omitempty"`
	Before              string      `json:"before,omitempty"`
	TargetUnit          string      `json:"targetUnit,omitempty"`
	TargetUnitSystem    string      `json:"targetUnitSystem,omitempty"`
	IncludeStatus       bool        `json:"includeStatus,omitempty"`
	IgnoreBadDataPoints bool        `json:"ignoreBadDataPoints,omitempty"`
	TreatUncertainAsBad bool        `json:"treatUncertainAsBad,omitempty"`
}
