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

// Save LastHash and Height to DB checkpoint
func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToByte(b))
}

// Restore blockchain by decoding []byte of data
func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

// AddBlock create new Block from the data and change LastHash and Height to Block's
func (b *blockchain) AddBlock(data string) {
	NewBlock := createBlock(data, b.LastHash, b.Height+1)
	b.LastHash = NewBlock.Hash
	b.Height = NewBlock.Height
	b.CurrentDifficulty = NewBlock.Difficulty
	b.persist()
}

// Blocks return []*Block of All Block exist
func (b *blockchain) Blocks() []*Block {
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

// Recalculate difficulty so that only blockInterval blocks are created per minute
// Increase difficulty if actualTime is less than expectedTime else decrease
func (b *blockchain) recalculateDifficulty() int {
	allBlocks := b.Blocks()
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

// Adjust difficulty for every blockInterval block
func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		return b.recalculateDifficulty()
	} else {
		return b.CurrentDifficulty
	}
}

// Blockchain return blockchain object or create if not exist
func Blockchain() *blockchain {
	// Create blockchain if not exist else return bc
	if bc == nil {
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
	}
	return bc
}
