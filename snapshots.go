package buckis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

var locker sync.RWMutex

var marshal = func(ht [9][1000]*DictEntry) (io.Reader, error) {
	b, err := json.MarshalIndent(ht, "", "\t")
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(b), nil
}

var unmarshal = func(r io.Reader, ht *[9][1000]*DictEntry) error {
	return json.NewDecoder(r).Decode(ht)
}

// SaveDataset  only to persist the dict object to disk in binary file
// format for dump: dump-(timestamp)
func SaveDataset(d *Db) error {
	// 1) set filename
	locker.Lock()
	defer locker.Unlock()
	filename := "dump:" + strconv.Itoa(int(time.Now().UnixMilli())) + ".bdb"

	// 2) open file p.s only creation mode no appending
	file, err := os.OpenFile(filepath.Join(d.config.PathToDump, filename), os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// 3) write dict to disk with marshalling now TODO: would circle back to gob encoding don't think json is efficient

	r, err := marshal(d.Ht)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(file, r)
	if err != nil {
		panic(err)
	}

	// 4) flush append log buckis.aof
	err = os.Remove("buckis.aof")
	if err != nil {
		return err
	}

	return err
}

func GetSnapshot(d *Db) error {
	dirs, err := os.ReadDir(d.config.PathToDump)
	if err != nil {
		return err
	}

	// assuming recent snapshot is last in the directory
	if len(dirs) == 0 {
		return fmt.Errorf("no current dump")
	}
	lastRecentDump := dirs[len(dirs)-1]
	name := lastRecentDump.Name()
	fmt.Println(name)

	file, err := os.OpenFile(filepath.Join(d.config.PathToDump, name), os.O_RDONLY, 0600)
	if err != nil {
		return nil
	}

	//binary.Read(file, binary.LittleEndian, d)
	// TODO: retry this gob thing, wasn't working due to null pointer
	//err = gob.NewDecoder(file).Decode(&d.Ht)
	//if err != nil {
	//	return err
	//}

	err = unmarshal(file, &d.Ht)
	if err != nil {
		panic(err)
	}

	return nil
}
