package main

import (
	"crypto/sha256"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/syndtr/goleveldb/leveldb"
	"sync"
)

var (
	DB *leveldb.DB

	HashLock *sync.RWMutex
	Hashes   map[string]*FileHash
)

type FileHash struct {
	PathOnDisk string
	Hash       []byte
}

// Update cross-references the hash at PathOnDisk with the one in the db, updating if necessary. Returns true if updated.
func (fh *FileHash) Update() bool {
	if DB == nil || fh == nil {
		return false
	}
	dbhash, err := DB.Get([]byte(fh.PathOnDisk), nil)
	if err != nil || string(dbhash) != string(fh.Hash) {
		DB.Put([]byte(fh.PathOnDisk), fh.Hash, nil)
		return true
	}
	return false
}

// Delete removes the PathOnFisk:Hash from the db.
func (fh *FileHash) Delete() {
	if DB == nil || fh == nil {
		return
	}
	DB.Delete([]byte(fh.PathOnDisk), nil)
}

// Recalculate simply recalculates the Hash, updating Hash and PathOnDisk, and returning a copy of the pointer.
func (fh *FileHash) Recalculate(PathOnDisk string) *FileHash {
	fh.PathOnDisk = PathOnDisk
	f, err := os.Open(fh.PathOnDisk)
	if err != nil {
		return nil
	}
	hash := sha256.New224()
	if _, err := io.Copy(hash, f); err != nil {
		f.Close()
		return nil
	}
	f.Close()

	fh.Hash = hash.Sum(nil)
	return fh
}

// HashDir recursively searches through a directory, hashing every file, and returning them as a list []*FileHash.
func HashDir(path string) (map[string]*FileHash, error) {
	files, err := filePathWalkDir(path)
	if err != nil {
		return nil, err
	}
	hashes := make(map[string]*FileHash, len(files))
	for _, file := range files {
		splitName := strings.Split(file, ".")
		if findInStringSlice(Ignore, splitName[len(splitName)-1]) > -1 {
			continue
		}
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		hash := sha256.New224()
		if _, err := io.Copy(hash, f); err != nil {
			f.Close()
			return nil, err
		}
		f.Close()
		hashes[file] = &FileHash{PathOnDisk: file, Hash: hash.Sum(nil)}
	}
	return hashes, nil
}

// InitDB initializes a database at `path`.
func InitDB(path string) {
	Hashes = make(map[string]*FileHash)
	HashLock = new(sync.RWMutex)
	tdb, err := leveldb.OpenFile(path, nil)
	if err != nil {
		log.Fatalln(err)
	}
	DB = tdb
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)
	go func() {
		<-c
		DB.Close()
		os.Exit(1)
	}()
}