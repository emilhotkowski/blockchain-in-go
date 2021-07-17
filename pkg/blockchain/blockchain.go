package blockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"
)

//Block is a generic, minimal block
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

//Create new block
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Nonce = nonce
	block.Hash = hash[:]

	return block
}

//Blockchain is a simple(only mine), array based database
type Blockchain struct {
	Blocks []*Block
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	nextBlock := NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, nextBlock)
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

//Proof-of-work

const targetBits = 24
const maxNonce = math.MaxInt64

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	return &ProofOfWork{b, target}
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			[]byte(strconv.FormatInt(pow.block.Timestamp, 16)),
			[]byte(strconv.FormatInt(targetBits, 16)),
			[]byte(strconv.FormatInt(int64(nonce), 16)),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hash [32]byte
	var hashInt big.Int
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		fmt.Printf("\r%x", hash)

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}

	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt *big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt = big.NewInt(1)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}
