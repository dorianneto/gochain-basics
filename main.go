package main

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

var BlockChain []Block

type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
}

func generateBlock(oldBlock Block, BPM int) (Block, error) {
	var newBlock Block

	timestamp := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = timestamp.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash

	hash := sha256.New()
	hash.Write([]byte(record))

	return hex.EncodeToString(hash.Sum(nil))
}

func main() {}
