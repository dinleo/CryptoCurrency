package blockchain

import (
	"CryptoCurrency/db"
	"CryptoCurrency/utils"
	"errors"
	"strings"
	"time"
)

type Block struct {
	Hash         string `json:"hash"`
	PrevHash     string `json:"prevHash,omitempty"`
	Height       int    `json:"height"`
	Difficulty   int    `json:"difficulty"`
	Nonce        int    `json:"nonce"`
	Timestamp    int    `json:"timestamp"`
	Transactions []*Tx  `json:"transactions"`
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

// Mine block by finding Nonce that can make collect hash and saving Timestamp when done
func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		b.Timestamp = int(time.Now().Unix())
		hash := utils.Hash(b)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}

// Create New Block and save to DB
func createBlock(prevHash string, height int) *Block {
	// Creating
	block := Block{

		Hash:         "",
		PrevHash:     prevHash,
		Height:       height,
		Difficulty:   Blockchain().difficulty(),
		Nonce:        0,
		Transactions: []*Tx{makeCoinbaseTx("leo")},
	}

	// Mining
	block.mine()
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
