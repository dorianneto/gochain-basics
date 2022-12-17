package main

import (
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/dorianneto/gochain/blockchain"
	"github.com/dorianneto/gochain/http"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		timestamp := time.Now()
		genesisBlock := blockchain.Block{0, timestamp.String(), 0, "", ""}

		spew.Dump(genesisBlock)

		blockchain.Blockchain = append(blockchain.Blockchain, genesisBlock)
	}()

	log.Fatal(http.Run())
}
