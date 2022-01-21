package blockchain

import (
	"CryptoCurrency/utils"
	"time"
)

const (
	minerReward int = 50
)

type Tx struct {
	Id        string    `json:"id"`
	Timestamp int       `json:"timestamp"`
	TxIns     []*TxIn   `json:"txIns"`
	TxOuts    []*TxOuts `json:"txOuts"`
}

type TxIn struct {
	Owner  string
	Amount int
}
type TxOuts struct {
	Owner  string
	Amount int
}

func (t *Tx) setID() {
	t.Id = utils.Hash(t)
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"COINBASE", minerReward},
	}
	txOuts := []*TxOuts{
		{address, minerReward},
	}
	tx := Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.setID()
	return &tx
}
