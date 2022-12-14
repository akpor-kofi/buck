package buckis

import (
	"errors"
	"github.com/google/btree"
	"github.com/samber/lo"
	"regexp"
	"strings"
)

var ErrIndexNotFound = errors.New("index not found")

type Filter struct {
	NumField string
	Min      int
	Max      int
}

type IndexOptions struct {
	Prefix    string
	StopWords []string
	NoFreqs   bool
	Schema    []string
}

type SearchOptions struct {
	NoContent bool
	Filter    Filter // on numeric
	Limit     []int
	On        int
}

type invertedIndexEntry struct {
	term          string
	frequency     int
	documentsList []doc // keys
}

type doc struct {
	id    string
	attrs string
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

func (d *dict) FTCreate(indexKey string, opts IndexOptions) error {
	i := d.hash(indexKey)

	idxTree := btree.NewG(2, lessFunc)

	for _, v := range d.Ht[Hashes] {
		if v == nil {
			continue
		}

		// if v.key

		currentNode := v

		for {
			if !strings.HasPrefix(currentNode.Key, opts.Prefix) {
				currentNode = currentNode.Next
			}

			hash := currentNode.Values.(map[string]any)

			lo.ForEach(opts.Schema, func(attr string, index int) {
				if sentence, ok := hash[attr]; ok {

					//process the strings
					wordlist := tokenize(sentence.(string))

					idxTree = d.populateIndexTable(idxTree, wordlist, v.Key, attr)
				}
			})

			if currentNode.Next == nil {
				break
			}

			currentNode = currentNode.Next

		}
	}

	de := &DictEntry{
		Key:    indexKey,
		Values: idxTree,
		Next:   d.Ht[Search][i],
	}

	d.Ht[Search][i] = de

	return nil
}

func (d *dict) FTSearch(indexKey, query string) (result []string, err error) {
	// index dict entry
	ide, err := d.indexKeyLookup(indexKey)
	if err != nil {
		return
	}

	tree := ide.Values.(*btree.BTreeG[*invertedIndexEntry])

	qEntry := &invertedIndexEntry{term: query}

	tree.AscendGreaterOrEqual(qEntry, func(entry *invertedIndexEntry) bool {
		if strings.HasPrefix(entry.term, query) {
			result = append(result, entry.term)
		}

		return true
	})

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

func (d *dict) indexKeyLookup(key string) (*DictEntry, error) {
	i := d.hash(key)
	currentEntry := d.Ht[Search][i]

	if currentEntry == nil {
		return &DictEntry{}, ErrIndexNotFound
	}

	for {
		if currentEntry.Key == key {
			return currentEntry, nil
		}

		if currentEntry.Next == nil {
			return &DictEntry{}, ErrIndexNotFound
		}

		currentEntry = currentEntry.Next
	}
}

func (d *dict) populateIndexTable(indexTree *btree.BTreeG[*invertedIndexEntry], wordList []string, hashKey string, attr string) *btree.BTreeG[*invertedIndexEntry] {
	for _, word := range wordList {

		document := doc{hashKey, attr}

		entry := &invertedIndexEntry{
			term:          word,
			frequency:     1,
			documentsList: []doc{document},
		}

		// try and check if it is already in the list
		if item, ok := indexTree.Get(&invertedIndexEntry{term: word}); !ok {
			indexTree.ReplaceOrInsert(entry)
		} else {
			item.documentsList = append(item.documentsList, doc{hashKey, attr})

			item.frequency = item.frequency + 1

			indexTree.ReplaceOrInsert(item)
		}

	}

	return indexTree
}
