package main

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type Transaction struct {
    From   string `json:"from"`
    To     string `json:"to"`
    Amount int    `json:"amount"`
}

type Block struct {
    Index        int
    Timestamp    string
    Transactions []Transaction
    PrevHash     string
    Hash         string
}

var chain []Block
var mempool []Transaction
var balances = map[string]int{
    "0x21a66166ce1C33CE016aFCd90C6d09b697afdD58": 1000000000,
}

func calculateHash(block Block) string {
    record := fmt.Sprintf("%d%s%v%s", block.Index, block.Timestamp, block.Transactions, block.PrevHash)
    h := sha256.Sum256([]byte(record))
    return hex.EncodeToString(h[:])
}

func createBlock() Block {
    prevBlock := chain[len(chain)-1]

    for _, tx := range mempool {
        if tx.From != "mint" {
            balances[tx.From] -= tx.Amount
        }
        balances[tx.To] += tx.Amount
    }

    block := Block{
        Index:        len(chain),
        Timestamp:    time.Now().String(),
        Transactions: mempool,
        PrevHash:     prevBlock.Hash,
    }

    block.Hash = calculateHash(block)
    mempool = []Transaction{}
    return block
}

func mine(w http.ResponseWriter, r *http.Request) {
    block := createBlock()
    chain = append(chain, block)
    json.NewEncoder(w).Encode(block)
}

func getBlocks(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(chain)
}

func faucet(w http.ResponseWriter, r *http.Request) {
    addr := r.URL.Query().Get("address")
    mempool = append(mempool, Transaction{"mint", addr, 100})
    json.NewEncoder(w).Encode("faucet queued")
}

func main() {
    genesis := Block{Index: 0, Timestamp: time.Now().String(), PrevHash: "0"}
    genesis.Hash = calculateHash(genesis)
    chain = append(chain, genesis)

    http.HandleFunc("/mine", mine)
    http.HandleFunc("/blocks", getBlocks)
    http.HandleFunc("/faucet", faucet)

    fmt.Println("NovaChain running on :8080")
    http.ListenAndServe(":8080", nil)
}
