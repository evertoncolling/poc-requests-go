package dto

type UnitConversion struct {
	Multiplier float64 `json:"multiplier"`
	Offset     float64 `json:"offset"`
}

type Unit struct {
	ExternalId      string         `json:"externalId"`
	Name            string         `json:"name"`
	LongName        string         `json:"longName"`
	Symbol          string         `json:"symbol"`
	AliasNames      []string       `json:"aliasNames"`
	Quantity        string         `json:"quantity"`
	Conversion      UnitConversion `json:"conversion"`
	Source          string         `json:"source,omitempty"`
	SourceReference string         `json:"sourceReference,omitempty"`
}

type UnitList struct {
	Items []Unit `json:"items"`
}
