package model

type Balance struct {
	Current   Amount `json:"current"`
	Withdrawn Amount `json:"withdrawn"`
}
