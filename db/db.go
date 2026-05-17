package db

import (
	"me/memtable"
	"me/wal"
	"path/filepath"
)

type DB struct {
	wal *wal.WAL
	mem *memtable.SkipList
}

func (db *DB) Get(key string) (string, bool) {
	return db.mem.Get(key)
}

func (db *DB) Put(key string, value string) {
	data := key + ":" + value
	kv := []byte(data)

	if err := db.wal.Write(kv); err != nil {
		return
	}

	db.mem.Put(key, value)

}

func NewDB(WalPath string) (*DB, error) {
	w, err := wal.OpenWAL(WalPath, filepath.Dir(WalPath))

	if err != nil {
		return nil, err
	}

	sl := memtable.NewSkipList()

	return &DB{
		wal: w,
		mem: sl,
	}, nil

}
