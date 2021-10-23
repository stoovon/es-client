package externalModels

type Query struct {
	Match map[string]string `json:"match"`
}

type ESQuery struct {
	Query Query `json:"query"`
}
