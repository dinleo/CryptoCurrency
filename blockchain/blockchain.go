package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type block struct {
	Data     string
	Hash     string
	PrevHash string
}

type blockchain struct {
	blocks []*block
}

var bc *blockchain
var once sync.Once

func (b *block) getHash() {
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash)
}

func getLastHash() string {
	totalBlocks := len(GetBlockchain().blocks)
	if totalBlocks == 0 {
		return ""
	}
	return GetBlockchain().blocks[totalBlocks-1].Hash
}

func createBlock(data string) *block {
	newBlock := block{data, "", getLastHash()}
	newBlock.getHash()
	return &newBlock
}

func (bc *blockchain) AddBlock(data string) {
	bc.blocks = append(bc.blocks, createBlock(data))
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

func (bc *blockchain) AllBlocks() []*block {
	return bc.blocks
}
