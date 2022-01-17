package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height   int    `json:"height"`
}

type blockchain struct {
	blocks []*Block
}

var bc *blockchain
var once sync.Once
var ErrNotFound = errors.New("block not found")

func getLastHash() string {
	totalBlocks := len(GetBlockchain().blocks)
	if totalBlocks == 0 {
		return ""
	}
	return GetBlockchain().blocks[totalBlocks-1].Hash
}

func createBlock(data string) *Block {
	newBlock := Block{data, "", getLastHash(), len(GetBlockchain().blocks) + 1}
	newBlock.getHash()
	return &newBlock
}

func GetBlockchain() *blockchain {
	if bc == nil {
		once.Do(func() {
			bc = &blockchain{}
			bc.AddBlock("Genesis Block")
		})
	}
	return bc
}

func (b *Block) getHash() {
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash)
}

func (bc *blockchain) AddBlock(data string) {
	bc.blocks = append(bc.blocks, createBlock(data))
}

func (bc *blockchain) AllBlocks() []*Block {
	return bc.blocks
}

func (bc *blockchain) GetBlock(height int) (*Block, error) {
	if height > len(bc.blocks) {
		return nil, ErrNotFound
	}
	return bc.blocks[height-1], nil
}
