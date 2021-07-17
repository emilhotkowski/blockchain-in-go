package blockchain

import (
	"log"

	"github.com/boltdb/bolt"
)

//Blockchain is a simple(only mine), array based database
var bucketName = []byte("blockchainBucket")

type Blockchain struct {
	tip []byte
	Db  *bolt.DB
}

func (bc *Blockchain) AddBlock(data string) {
	nextBlock := NewBlock(data, bc.tip)

	bc.Db.Update(func(t *bolt.Tx) error {
		b := t.Bucket(bucketName)
		b.Put([]byte("l"), nextBlock.Hash)
		b.Put(nextBlock.Hash, nextBlock.Serialize())

		return nil
	})
}

func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open("blockchain.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	db.Update(func(t *bolt.Tx) error {
		b := t.Bucket(bucketName)

		if b == nil {
			b, err = t.CreateBucket(bucketName)
			if err != nil {
				log.Fatal(err)
				return nil
			}

			genesis := NewGenesisBlock()
			err = b.Put(genesis.Hash, genesis.Serialize())
			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				log.Fatal(err)
				return nil
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	return &Blockchain{tip, db}
}

//Block chain iterator
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.tip, bc.Db}
}

func (bci *BlockchainIterator) Next() *Block {
	var lastBlock *Block
	bci.db.View(func(t *bolt.Tx) error {
		b := t.Bucket(bucketName)
		lastBlockData := b.Get(bci.currentHash)
		lastBlock = DeserializeBlock(lastBlockData)

		return nil
	})
	bci.currentHash = lastBlock.PrevBlockHash
	return lastBlock
}
