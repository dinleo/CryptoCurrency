package blockchain

import (
	"CryptoCurrency/db"
	"CryptoCurrency/utils"
	"sync"
)

type blockchain struct {
	LastHash          string `json:"lastHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
}

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5
	blockInterval      int = 2
	allowedRange       int = 2
)

var bc *blockchain
var once sync.Once

// Method

// Restore blockchain by decoding []byte of data
func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

// AddBlock create new Block from the data and change LastHash and Height to Block's
func (b *blockchain) AddBlock(data string) {
	NewBlock := createBlock(b.LastHash, b.Height+1, data, getDifficulty(b))
	b.LastHash = NewBlock.Hash
	b.Height = NewBlock.Height
	b.CurrentDifficulty = NewBlock.Difficulty
	persist(b)
}

// Function

// Save LastHash and Height to DB checkpoint
func persist(b *blockchain) {
	db.SaveBlockchain(utils.ToByte(b))
}

// Recalculate getDifficulty based on block creating speed
func recalculateDifficulty(b *blockchain) int {
	allBlocks := Blocks(b)
	newestBlock := allBlocks[0]
	lastRecalculatedBlock := allBlocks[difficultyInterval-1]
	actualTime := (newestBlock.Timestamp - lastRecalculatedBlock.Timestamp) / 60
	expectedTime := difficultyInterval * blockInterval
	if actualTime < (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTime > (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1
	} else {
		return b.CurrentDifficulty
	}
}

// Adjust getDifficulty for every blockInterval block
func getDifficulty(b *blockchain) int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		return recalculateDifficulty(b)
	} else {
		return b.CurrentDifficulty
	}
}

// UTxOutsByAddress return []*UTxOut of all unspent Tx for an address
func UTxOutsByAddress(address string, b *blockchain) []*UTxOut {
	var uTxOuts []*UTxOut

	usedTxId := make(map[string]bool)
	for _, block := range Blocks(b) {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Owner == address {
					// Find TxId already used as TxIns until just before Tx
					usedTxId[input.TxId] = true
				}
			}
			for index, output := range tx.TxOuts {
				if output.Owner == address {
					_, used := usedTxId[tx.Id]
					if !used {
						// Append only not used TxOut yet
						uTxOut := &UTxOut{tx.Id, index, output.Amount}
						if !isOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}
	return uTxOuts
}

// BalanceByAddress return Balance for an address
func BalanceByAddress(address string, b *blockchain) int {
	txOuts := UTxOutsByAddress(address, b)
	var amount int
	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
}

// Wallets return all wallets exists
func Wallets(b *blockchain) []string {
	walletsMap := make(map[string]bool)
	var wallets []string
	for _, block := range Blocks(b) {
		for _, tx := range block.Transactions {
			for _, output := range tx.TxOuts {
				walletsMap[output.Owner] = true
			}
		}
	}
	for w, _ := range walletsMap {
		wallets = append(wallets, w)
	}
	return wallets
}

// Blocks return []*Block of all Block exist
func Blocks(b *blockchain) []*Block {
	var blocks []*Block
	hashCursor := b.LastHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}

// Blockchain return blockchain object or create if not exist
func Blockchain() *blockchain {
	// Create blockchain if not exist else return bc
	once.Do(func() {
		bc = &blockchain{Height: 0}
		checkpoint := db.Checkpoint()

		// Create block & DB if not exist else get LastHash & Height from checkpoint
		if checkpoint == nil {
			bc.AddBlock("Genesis Block")
		} else {
			bc.restore(checkpoint)
		}
	})
	return bc
}
