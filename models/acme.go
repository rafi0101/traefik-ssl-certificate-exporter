package models

type ProviderMdl map[string]Provider

type Provider struct {
	Account struct {
		Email        string `json:"Email"`
		Registration struct {
			Body struct {
				Status  string   `json:"status"`
				Contact []string `json:"contact"`
			} `json:"body"`
			URI string `json:"uri"`
		} `json:"Registration"`
		PrivateKey string `json:"PrivateKey"`
		KeyType    string `json:"KeyType"`
	} `json:"Account"`
	Certificates []struct {
		Domain struct {
			Main string   `json:"main"`
			Sans []string `json:"sans"`
		} `json:"domain"`
		Certificate string `json:"certificate"`
		Key         string `json:"key"`
		Store       string `json:"Store"`
	} `json:"Certificates"`
}
