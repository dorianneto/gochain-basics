package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/dorianneto/gochain/blockchain"
)

type Payload struct {
	BPM int
}

func handleGetBlockchain(response http.ResponseWriter, request *http.Request) {
	bytes, err := json.MarshalIndent(blockchain.Blockchain, "", " ")

	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}

	io.WriteString(response, string(bytes))
}

func handleWriteBlock(response http.ResponseWriter, request *http.Request) {
	var payload Payload

	decoder := json.NewDecoder(request.Body)

	if err := decoder.Decode(&payload); err != nil {
		respondWithJson(response, request, http.StatusBadRequest, request.Body)
		return
	}

	defer request.Body.Close()

	oldblock := blockchain.Blockchain[len(blockchain.Blockchain)-1]
	newBlock, err := blockchain.GenerateBlock(oldblock, payload.BPM)

	if err != nil {
		respondWithJson(response, request, http.StatusInternalServerError, payload)
		return
	}

	if blockchain.IsBlockValid(newBlock, oldblock) {
		newBlockChain := append(blockchain.Blockchain, newBlock)
		blockchain.ReplaceChain(newBlockChain)
		spew.Dump(blockchain.Blockchain)
	}

	respondWithJson(response, request, http.StatusCreated, newBlock)
}
