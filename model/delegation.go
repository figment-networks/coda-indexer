package model

import (
	"github.com/figment-networks/mina-indexer/model/types"
)

type Delegation struct {
	PublicKey string           `json:"public_key"`
	Delegate  string           `json:"delegate"`
	Balance   types.Amount     `json:"balance"`
	Weight    types.Percentage `json:"weight"`
}
