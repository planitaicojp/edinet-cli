package model

// Company represents an entry in the EDINET code list CSV.
type Company struct {
	EdinetCode    string `json:"edinetCode"`
	FilerType     string `json:"filerType"`
	ListingStatus string `json:"listingStatus"`
	Consolidated  string `json:"consolidated"`
	Capital       string `json:"capital"`
	FiscalYearEnd string `json:"fiscalYearEnd"`
	FilerName     string `json:"filerName"`
	FilerNameEN   string `json:"filerNameEN"`
	FilerNameKana string `json:"filerNameKana"`
	Address       string `json:"address"`
	Industry      string `json:"industry"`
	SecCode       string `json:"secCode"`
	JCN           string `json:"JCN"`
}
