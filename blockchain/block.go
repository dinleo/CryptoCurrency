package blockchain

import (
	"CryptoCurrency/db"
	"CryptoCurrency/utils"
	"crypto/sha256"
	"errors"
	"fmt"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height   int    `json:"height"`
}

var ErrNotFound = errors.New("block not found")

// Save Hash and BlockData to DB
func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToByte(b))
}

// Restore Block by decoding []byte of data
func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

// Create New Block and save to DB
func createBlock(data string, prevHash string, height int) *Block {
	// Creating
	block := Block{
		Data:     data,
		Hash:     "",
		PrevHash: prevHash,
		Height:   height,
	}

	// Hashing
	payload := block.Data + block.PrevHash + fmt.Sprint(block.Height)
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))

	// Persist
	block.persist()

	return &block
}

// FindBlock return Block which found from DB using Hash
func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}
