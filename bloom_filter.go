package buckis

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"strconv"
)

// my implementation on how to reduce collision ? use two hash functions (256 but append number of hash)
const numberOfHashFuncs = 2

func (d *dict) BFAdd(key, value string) error {
	return d.bfAdd(SAVE, key, value)
}

func (d *dict) BFExists(key, value string) (bool, error) {
	bfde, err := d.bloomFilterLookup(key)

	if err != nil {
		return false, ErrSetNotFound
	}

	bitArray := bfde.values.([100]bool)

	for k := 1; k <= numberOfHashFuncs; k++ {
		if bitArray[bloomHashFunc(value, k)] {
			return bitArray[bloomHashFunc(value, k)], nil
		} else {
			continue
		}
	}

	// means value is not in bloom filter
	return false, nil
}

func (d *dict) bloomFilterLookup(key string) (*dictEntry, error) {
	i := d.hash(key)

	currentEntry := d.ht[Bloom][i]

	if currentEntry == nil {
		return &dictEntry{}, ErrSetNotFound
	}

	for {
		if currentEntry.key == key {
			return currentEntry, nil
		}

		if currentEntry.next == nil {
			return &dictEntry{}, ErrSetNotFound
		}

		currentEntry = currentEntry.next
	}
}

func bloomHashFunc(value string, k int) uint64 {
	salt := strconv.Itoa(k)

	value += salt

	h := sha256.New()
	h.Write([]byte(value))

	bits := h.Sum(nil)

	buf := bytes.NewBuffer(bits)
	result, _ := binary.ReadUvarint(buf)

	//fmt.Println(result)
	return result % 100
}
