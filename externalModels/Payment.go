package externalModels

type Amount struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type Party struct {
	Name          string `json:"Name"`
	AccountNumber string `json:"AccountNumber"`
	BankId        string `json:"BankId"`
}

type Payment struct {
	Id          string `json:"Id"`
	Amount      Amount `json:"Amount"`
	Beneficiary Party  `json:"Beneficiary"`
	Debtor      Party  `json:"Debtor"`
}
