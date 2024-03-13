package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	Index     int
	Timestamp string
	Data      string
	PrevHash  string
	Hash      string
}

type Blockchain struct {
	Chain []Block
}

func calculateHash(block Block) string {
	data := strconv.Itoa(block.Index) + block.Timestamp + block.Data + block.PrevHash
	hash := sha256.New()
	hash.Write([]byte(data))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(prevBlock Block, data string) Block {
	var newBlock Block
	newBlock.Index = prevBlock.Index + 1
	newBlock.Timestamp = time.Now().String()
	newBlock.Data = data
	newBlock.PrevHash = prevBlock.Hash
	newBlock.Hash = calculateHash(newBlock)
	return newBlock
}

func createGenesisBlock() Block {
	var genesisBlock Block
	genesisBlock.Index = 0
	genesisBlock.Timestamp = time.Now().String()
	genesisBlock.Data = "Genesis Block"
	genesisBlock.PrevHash = ""
	genesisBlock.Hash = calculateHash(genesisBlock)
	return genesisBlock
}

func createBlockchain() Blockchain {
	var blockchain Blockchain
	genesisBlock := createGenesisBlock()
	blockchain.Chain = append(blockchain.Chain, genesisBlock)
	return blockchain
}

func addBlock(blockchain *Blockchain, data string) {
	prevBlock := blockchain.Chain[len(blockchain.Chain)-1]
	newBlock := generateBlock(prevBlock, data)
	blockchain.Chain = append(blockchain.Chain, newBlock)
}

func listBlocks(blockchain Blockchain) {
	for _, block := range blockchain.Chain {
		fmt.Printf("Index: %d\n", block.Index)
		fmt.Printf("Timestamp: %s\n", block.Timestamp)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Previous Hash: %s\n", block.PrevHash)
		fmt.Printf("Hash: %s\n\n", block.Hash)
	}
}

func saveBlockchain(blockchain Blockchain, filename string) error {
	data, err := json.MarshalIndent(blockchain, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Blockchain saved to", filename)
	return nil
}

func main() {
	blockchain := createBlockchain()

	fmt.Println("Blockchain created")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter data for the next block (or 'q' to quit): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "q" {
			break
		}

		addBlock(&blockchain, input)
		fmt.Println("Block added")

		listBlocks(blockchain)
	}

	filename := "blockchain.go"
	err := saveBlockchain(blockchain, filename)
	if err != nil {
		fmt.Println("Error saving blockchain:", err)
		os.Exit(1)
	}
}
