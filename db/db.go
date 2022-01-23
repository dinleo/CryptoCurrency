package db

import (
	"CryptoCurrency/utils"
	"github.com/boltdb/bolt"
)

const (
	dbName       = "blockchain.db"
	dataBucket   = "data"
	blocksBucket = "blocks"
	checkpoint   = "checkpoint"
)

var db *bolt.DB

// DB return bolt.DB object or create if not exist
func DB() *bolt.DB {
	if db == nil {
		// Create new DB
		dbPointer, err := bolt.Open(dbName, 0600, nil)
		utils.HandleErr(err)
		db = dbPointer
		db.Update(func(tx *bolt.Tx) error {
			// Create new Bucket
			_, err := tx.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleErr(err)
			_, err = tx.CreateBucketIfNotExists([]byte(blocksBucket))
			utils.HandleErr(err)
			return nil
		})
	}
	return db
}

// Close DB
func Close() {
	DB().Close()
}

// SaveBlock save Hash and Name to DB in {Hash:Name} pair format
func SaveBlock(hash string, data []byte) {
	DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), data)
		utils.HandleErr(err)
		return nil
	})
}

// SaveBlockchain save Name to DB in {"checkpoint":Name} pair format
func SaveBlockchain(data []byte) {
	DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkpoint), data)
		utils.HandleErr(err)
		return err
	})
}

// Checkpoint return []byte of LastHash and Height from DB checkpoint
func Checkpoint() []byte {
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})
	return data
}

// Block return []byte of BlockData from DB using Hash
func Block(hash string) []byte {
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})
	return data
}
