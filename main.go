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

type Block struct {
	Index int 
	Timestamp string 
	Hash string
	PrevHash string 
}
var BlockChain []Block


func calculateHash(block Block) string  {
record := string(block.Index) + block.TimeStamp + string(block.BPM) + block.PrevHash
h :=  sha256.New()
h.Write([]byte(record))
hashed := h.Sum(nil)
return hex.EncodeToString(hashed)
}

func generateBlock(oldBlock block, BPM int)(Block, error) {
var newBlock block

t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.TimeStamp =t.String()
	newBlock.BPM=BPM
	newBlock.PrevHash = oldBlock.Hash 
	newBlock.Hash=calculateHash(newBlock)

	return newBlock,nil
}


	func isBlockValid(newBlock, oldBlock Block) bool {
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
func replaceChain (newBlocks []Block) {
	if len(newBlocks) >len(BlockChain) {
	BlockChain =newBlocks
	}
}
func run() error {
	mux := makeMuxRouter()
	httpAddr := os.Getenv("ADDR")
	log.Println("listening on ",os.Getenv("ADDR"))
	s := &http.Server{
		Addr: ":" + httpAddr,
		Handler:   mux,
		RealTimeout:  10 *time.Second,
		WriteTimeout: 10 *time.Second,
		MaxHeaderBytes: 1 << 20,
}
		if err:= s.ListenAndServer(); err!=nil {
	return err
	}

	return nil
}

func makeMuxRouter() http.Handler {
muxRouter := mux.NewRouter()
muxRouter.HandleFunc("/", handleGetBlokChain).Methods("GET")
muxRouter.HandleFunc("/",handleWriteBlock).Methods("POST")
return muxRouter
}

func handleGetBlockChain(w http.ResponseWriter, r*http.Request){
bytes, err := json.MarshalIndent(BlockChain, "", " ") 
if err != nil {
http.Error(w, err.Error(), httpStatusInternalServerError)
return 
}
io.WriteString(w, string(bytes))
}

type Message struct {
BPM int 
}
func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], m.BPM)
	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}
	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		newBlockchain := append(Blockchain, newBlock)
		replaceChain(newBlockchain)
		spew.Dump(Blockchain)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

func respondWithJSON( w http.ResponseWriter, r *http.Request, code int, payload interface{} {
response, err := json.MarshalIndent(payload, "", " "}
		if err != nil {
		w.WriteReader(http.StatusInternalServerError) 
		w.Writ([] byte(" HTTP 500: Internal Server Error"))
		return 
	} 
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		t := time.Now()
		genesisBlock := Block{0, t.String(), 0, "", ""}
		spew.Dump(genesisBlock)
		Blockchain = append(Blockchain, genesisBlock)
	}()
	log.Fatal(run())

}

