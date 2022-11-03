package buckis

import (
	murmur "github.com/aviddiviner/go-murmur"
	"log"
	"os"
	"sync"
)

const (
	Strings = iota
	Hashes
	Set
	SortedSet
	Search
	List

	// Graph
)

type dictEntry struct {
	key    string
	values any // possible types are string, int, hash(map[string]any)
	next   *dictEntry
}

type dict struct {
	ht               [6][100]*dictEntry
	commandLoadQueue chan command
	commandQueueQuit chan chan error
	commandChan      chan command
	waiter           *sync.WaitGroup
	locker           *sync.Mutex
	aof              *os.File
}

func newDict() *dict {
	// open or create append-only file
	aof, err := os.OpenFile("./buckis.aof", os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatal(err)
	}

	return &dict{
		ht:               [6][100]*dictEntry{},
		commandLoadQueue: make(chan command, 10),
		commandQueueQuit: make(chan chan error),
		commandChan:      make(chan command),
		waiter:           &sync.WaitGroup{},
		locker:           &sync.Mutex{},
		aof:              aof,
	}
}

func (d *dict) hash(key string) uint32 {
	b := []byte(key)

	h := murmur.MurmurHash2(b, 0)

	return h % uint32(len(d.ht))
}
