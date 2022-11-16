package buckis

import (
	murmur "github.com/aviddiviner/go-murmur"
	"log"
	"os"
	"sync"
	"time"
)

const (
	Strings = iota
	Hashes
	Set
	SortedSet
	Search
	List
	Bloom
	Json

	// Graph
)

type DictEntry struct {
	Key    string
	Values any // possible types are string, int, hash(map[string]any)
	Next   *DictEntry
}

type Config struct {
	IntervalToSave time.Duration
	PathToAOF      string
	PathToDump     string
	AppendOnly     bool
}

type dict struct {
	Ht               [9][1000]*DictEntry
	hexastore        []string
	commandLoadQueue chan command
	commandChan      chan command
	waiter           *sync.WaitGroup
	aof              *os.File
	config           *Config
}

func newDict(config *Config) *dict {
	// open or create append-only file
	aof, err := os.OpenFile("./buckis.aof", os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatal(err)
	}

	return &dict{
		Ht:               [9][1000]*DictEntry{},
		commandLoadQueue: make(chan command, 10),
		commandChan:      make(chan command),
		waiter:           &sync.WaitGroup{},
		aof:              aof,
		config:           config,
	}
}

func (d *dict) hash(key string) uint32 {
	b := []byte(key)

	h := murmur.MurmurHash2(b, 0)

	return h % uint32(len(d.Ht))
}
