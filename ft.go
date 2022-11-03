package buckis

import (
	"errors"
	"fmt"
	"github.com/google/btree"
	"github.com/samber/lo"
	"regexp"
	"strings"
)

var ErrIndexNotFound = errors.New("index not found")
var ErrWordNotFound = errors.New("word not found")

type invertedIndexEntry struct {
	term          string
	frequency     int
	documentsList []string // keys
}

func (iie *invertedIndexEntry) Less(than btree.Item) bool {

	return !iie.Less(than) && !than.Less(iie)
}

func lessFunc(a, b *invertedIndexEntry) bool {
	less := strings.Compare(a.term, b.term)

	switch less {
	case -1:
		return true
	case 1:
		return false
	default:
		return false
	}
}

func (d *dict) FTCreate(indexKey, collectionsKey string, schema ...string) error {
	i := d.hash(indexKey)

	invertedIndex := d.ht[Search][i]

	de := &dictEntry{
		key:    indexKey,
		values: btree.NewG(2, lessFunc),
	}

	if invertedIndex == nil {
		de.next = nil
	} else {
		de.next = invertedIndex
	}

	// TODO: overload the DE.VALUES with indexes

	for _, v := range d.ht[Hashes] {
		if v == nil {
			continue
		}

		currentNode := v
		for {

			hash := currentNode.values.(map[string]any)

			lo.ForEach(schema, func(item string, index int) {
				if sentence, ok := hash[item]; ok {

					//process the strings
					wordlist := tokenize(sentence.(string))

					de.values = d.populateIndexTable(de.values.(*btree.BTreeG[*invertedIndexEntry]), wordlist, v.key)
				}
			})

			if currentNode.next == nil {
				break
			}

			currentNode = currentNode.next

		}
	}

	d.ht[Search][i] = de

	return nil
}

func (d *dict) FTSearch(indexKey, query string) (result []string, err error) {
	// index dict entry
	ide, err := d.indexKeyLookup(indexKey)
	if err != nil {
		return
	}

	tree := ide.values.(*btree.BTreeG[*invertedIndexEntry])

	qEntry := &invertedIndexEntry{term: query}

	tree.AscendGreaterOrEqual(qEntry, func(entry *invertedIndexEntry) bool {
		if strings.HasPrefix(entry.term, query) {
			result = append(result, entry.term)
		}

		return true
	})

	return
}

func (d *dict) FTFind(indexKey, word string) (result []string, err error) {
	// index dict entry
	ide, err := d.indexKeyLookup(indexKey)
	if err != nil {
		return
	}

	tree := ide.values.(*btree.BTreeG[*invertedIndexEntry])

	if item, ok := tree.Get(&invertedIndexEntry{term: word}); !ok {
		err = fmt.Errorf("word not found")
	} else {
		result = item.documentsList
	}

	return
}

func tokenize(doc string) (wordList []string) {

	// The following regexp finds individual
	// words in a sentence
	r := regexp.MustCompile("[^\\s]+")
	wordList = r.FindAllString(doc, -1)

	for i := 0; i < len(wordList); i++ {
		wordList[i] = strings.ToLower(wordList[i])
	}

	return
}

func (d *dict) indexKeyLookup(key string) (*dictEntry, error) {
	i := d.hash(key)
	currentEntry := d.ht[Search][i]

	if currentEntry == nil {
		return &dictEntry{}, ErrIndexNotFound
	}

	for {
		if currentEntry.key == key {
			return currentEntry, nil
		}

		if currentEntry.next == nil {
			return &dictEntry{}, ErrIndexNotFound
		}

		currentEntry = currentEntry.next
	}
}

func (d *dict) populateIndexTable(indexTree *btree.BTreeG[*invertedIndexEntry], wordList []string, hashKey string) *btree.BTreeG[*invertedIndexEntry] {
	for _, word := range wordList {

		entry := &invertedIndexEntry{
			term:          word,
			frequency:     1,
			documentsList: []string{hashKey},
		}

		// try and check if it is already in the list
		if item, ok := indexTree.Get(&invertedIndexEntry{term: word}); !ok {
			indexTree.ReplaceOrInsert(entry)
		} else {
			if _, found := lo.Find(item.documentsList, func(document string) bool {
				return document == hashKey
			}); !found {
				item.documentsList = append(item.documentsList, hashKey)
			}

			item.frequency = item.frequency + 1

			indexTree.ReplaceOrInsert(item)
		}

	}

	return indexTree
}
