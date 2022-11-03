package buckis

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/aviddiviner/go-murmur"
	"log"
	"reflect"
	"time"
)

var errNotFoundAndFirstEntry = fmt.Errorf("first entry on linked list")
var errChainNotFound = fmt.Errorf("chain not found")

const MineRate = 1000
const sizeOfChainHT = 100

type blockchainDict[T any] struct {
	ht [sizeOfChainHT]*blockchain[T]
}

func NewBlockchainDict[T any]() *blockchainDict[T] {
	return &blockchainDict[T]{
		[sizeOfChainHT]*blockchain[T]{},
	}
}

type p2pNet[T any] struct {
	nodes []*blockchain[T]
}

type block[T any] struct {
	timestamp int64
	data      []T
	hash      string
	lastHash  string
}

type blockchain[T any] struct {
	key string
	//password string
	chain []block[T]

	//net *p2pNet[T]

	next *blockchain[T]
}

func genesis[T any]() *block[T] {
	h := sha256.New()
	h.Write([]byte("genesis block"))
	hash := fmt.Sprintf("%x", h.Sum(nil))

	return &block[T]{
		timestamp: 1,
		data:      []T{}, // had to put nil for import cycles reasons
		hash:      hash,
		lastHash:  hash,
	}
}

func makeNewChain[T any](key string) *blockchain[T] {
	genesisBlock := genesis[T]()

	return &blockchain[T]{
		key:   key,
		chain: []block[T]{*genesisBlock},
	}
}

func blockchainHash(key string) uint32 {
	b := []byte(key)

	h := murmur.MurmurHash2(b, 0)

	return h % sizeOfChainHT
}

func (b *blockchainDict[T]) BCAdd(key string, data []T) error {
	// first find the chain from hashtable
	bl, err := b.chainLookup(key)
	i := blockchainHash(key)

	if err != nil {
		// there was no chain so we are to create one
		newChain := makeNewChain[T](key)
		newChain.addBlock(data)

		// add the new chain to hashtable
		if b.ht[i] == nil {
			newChain.next = nil
			b.ht[i] = newChain
		} else {
			newChain.next = b.ht[i]
			b.ht[i] = newChain
		}

		return nil
	}

	bl.addBlock(data)

	return nil
}

func (b *blockchainDict[T]) BCView(key string) (blocks [][]T, err error) {
	bl, err := b.chainLookup(key)
	if err != nil {
		return
	}

	for _, v := range bl.chain {
		blocks = append(blocks, v.data)
	}

	return
}

func (b *blockchainDict[T]) BCIsValid(key string) bool {
	bl, err := b.chainLookup(key)

	if err != nil {
		return false
	}

	firstBlock := genesis[T]()

	if reflect.DeepEqual(bl.chain[0], firstBlock) {
		fmt.Println("not valid because the first block is not the standard genesis")
		return false
	}

	for i := 1; i < len(bl.chain); i++ {
		chain := bl.chain[i]
		actualLastHash := bl.chain[i-1].hash
		if actualLastHash != chain.lastHash {
			return false
		}

		h := hashable[T]{chain.timestamp, chain.data, chain.lastHash}
		validatedHash := cryptoHash(h)

		if chain.hash != validatedHash {
			return false
		}
	}

	return true
}

func (b *blockchainDict[T]) BCLast(key string) (data []T, err error) {
	bl, err := b.chainLookup(key)

	if err != nil {
		return
	}

	data = bl.chain[len(bl.chain)-1].data

	return
}

func (b *blockchainDict[T]) chainLookup(key string) (*blockchain[T], error) {
	i := blockchainHash(key)

	currentEntry := b.ht[i]

	if currentEntry == nil {
		return &blockchain[T]{}, errNotFoundAndFirstEntry
	}

	for {
		if currentEntry.key == key {
			return currentEntry, nil
		}

		if currentEntry.next == nil {
			return &blockchain[T]{}, errChainNotFound
		}

		currentEntry = currentEntry.next
	}

}

func (bl *blockchain[T]) addBlock(data []T) {
	fmt.Println(len(bl.chain))
	dataBlock := generateBlock[T](bl.chain[len(bl.chain)-1], data)
	bl.chain = append(bl.chain, *dataBlock)
}

type hashable[T any] struct {
	timestamp int64
	data      []T
	lastHash  string
}

func generateBlock[T any](lastBlock block[T], data []T) *block[T] {
	t := time.Now().UnixMilli()
	lh := lastBlock.hash

	var hash string

	h := hashable[T]{t, data, lh}
	hash = cryptoHash(h)

	return &block[T]{
		timestamp: t,
		data:      data,
		lastHash:  lh,
		hash:      hash,
	}
}

func cryptoHash(o any) string {
	h := sha256.New()
	// h.Write([]byte(fmt.Sprintf("%v", o)))
	marshal, err := json.Marshal(o)
	if err != nil {
		log.Fatal(err)
	}
	h.Write(marshal)

	hash := hex.EncodeToString(h.Sum(nil))

	//fmt.Println(len(h.Sum(nil)))   // 32 bytes
	//fmt.Println(len(hex.DecodeString(h.Sum(nil))))   // 32 bytes
	//fmt.Println(len([]byte(hash))) // 64 bytes

	return hash
}
