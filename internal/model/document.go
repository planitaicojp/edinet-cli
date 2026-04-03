package model

type Document struct {
	SeqNumber            int     `json:"seqNumber"`
	DocID                string  `json:"docID"`
	EdinetCode           string  `json:"edinetCode"`
	SecCode              string  `json:"secCode"`
	JCN                  string  `json:"JCN"`
	FilerName            string  `json:"filerName"`
	FundCode             string  `json:"fundCode"`
	OrdinanceCode        string  `json:"ordinanceCode"`
	FormCode             string  `json:"formCode"`
	DocTypeCode          string  `json:"docTypeCode"`
	PeriodStart          *string `json:"periodStart"`
	PeriodEnd            *string `json:"periodEnd"`
	SubmitDateTime       string  `json:"submitDateTime"`
	DocDescription       string  `json:"docDescription"`
	IssuerEdinetCode     *string `json:"issuerEdinetCode"`
	SubjectEdinetCode    *string `json:"subjectEdinetCode"`
	SubsidiaryEdinetCode *string `json:"subsidiaryEdinetCode"`
	CurrentReportReason  *string `json:"currentReportReason"`
	ParentDocID          *string `json:"parentDocID"`
	OpeDateTime          *string `json:"opeDateTime"`
	WithdrawalStatus     string  `json:"withdrawalStatus"`
	DocInfoEditStatus    string  `json:"docInfoEditStatus"`
	DisclosureStatus     string  `json:"disclosureStatus"`
	XbrlFlag             string  `json:"xbrlFlag"`
	PdfFlag              string  `json:"pdfFlag"`
	AttachDocFlag        string  `json:"attachDocFlag"`
	EnglishDocFlag       string  `json:"englishDocFlag"`
	CsvFlag              string  `json:"csvFlag"`
	LegalStatus          string  `json:"legalStatus"`
}
