package db

import (
	"me/memtable"
	"me/wal"
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
