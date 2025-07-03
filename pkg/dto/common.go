package dto

type Metadata map[string]string

type Identity struct {
	Id         int64  `json:"id,omitempty"`
	ExternalId string `json:"externalId,omitempty"`
}

type TimestampRange struct {
	Min int64 `json:"min,omitempty"`
	Max int64 `json:"max,omitempty"`
}
