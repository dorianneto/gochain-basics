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

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var Blockchain []Block

type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
}

type Payload struct {
	BPM int
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
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

func handleGetBlockchain(response http.ResponseWriter, request *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", " ")

	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}

	io.WriteString(response, string(bytes))
}

func respondWithJson(response http.ResponseWriter, request *http.Request, code int, payload interface{}) {
	output, err := json.MarshalIndent(payload, "", " ")

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}

	response.WriteHeader(code)
	response.Write(output)
}

func handleWriteBlock(response http.ResponseWriter, request *http.Request) {
	var payload Payload

	decoder := json.NewDecoder(request.Body)

	if err := decoder.Decode(&payload); err != nil {
		respondWithJson(response, request, http.StatusBadRequest, request.Body)
		return
	}

	defer request.Body.Close()

	oldblock := Blockchain[len(Blockchain)-1]
	newBlock, err := generateBlock(oldblock, payload.BPM)

	if err != nil {
		respondWithJson(response, request, http.StatusInternalServerError, payload)
		return
	}

	if isBlockValid(newBlock, oldblock) {
		newBlockChain := append(Blockchain, newBlock)
		replaceChain(newBlockChain)
		spew.Dump(Blockchain)
	}

	respondWithJson(response, request, http.StatusCreated, newBlock)
}

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

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		timestamp := time.Now()
		genesisBlock := Block{0, timestamp.String(), 0, "", ""}

		spew.Dump(genesisBlock)

		Blockchain = append(Blockchain, genesisBlock)
	}()

	log.Fatal(run())
}
