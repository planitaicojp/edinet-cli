package model

type Metadata struct {
	Title           string         `json:"title"`
	Parameter       MetadataParam  `json:"parameter"`
	ResultSet       MetadataResult `json:"resultset"`
	ProcessDateTime string         `json:"processDateTime"`
	Status          string         `json:"status"`
	Message         string         `json:"message"`
}

type MetadataParam struct {
	Date string `json:"date"`
	Type string `json:"type"`
}

type MetadataResult struct {
	Count int `json:"count"`
}

type DocumentListResponse struct {
	Metadata Metadata   `json:"metadata"`
	Results  []Document `json:"results"`
}
