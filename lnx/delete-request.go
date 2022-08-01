package lnx

type deleteRequest struct {
	Query query `json:"query"`
	Limit int   `json:"limit"`
}
