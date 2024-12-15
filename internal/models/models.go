package models

type RequestShortenAPI struct {
	URL string `json:"url"`
}

type ResponseShortenAPI struct {
	Result string `json:"result"`
}
