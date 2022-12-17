package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
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

func isBlockValid(newBlock Block, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// It needs in order to pick the right blockchain as the source of truth
func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(BlockChain) {
		BlockChain = newBlocks
	}
}

func handleGetBlockchain(response http.ResponseWriter, request *http.Request) {
	bytes, err := json.MarshalIndent(BlockChain, "", " ")

	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}

	io.WriteString(response, string(bytes))
}

func handleWriteBlock(response http.ResponseWriter, request *http.Request) {}

func router() http.Handler {
	muxRouter := mux.NewRouter()

	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")

	return muxRouter
}

func run() error {
	router := router()

	port := os.Getenv("PORT")
	log.Println("Listening on ", port)

	server := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func main() {}
