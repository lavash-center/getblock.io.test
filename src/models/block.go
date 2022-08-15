package models

type BlockByNumberResponse struct {
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	To    string `json:"to"`
	From  string `json:"from"`
	Value string `json:"value"`
}
