package db

import (
	"encoding/binary"
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

func (db *DB) Put(key string, value string) error {
	length := uint32(len(key))
	lenBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(lenBytes, length)

	data := append(lenBytes, []byte(key)...)
	data = append(data, []byte(value)...)
	if err := db.wal.Write(data); err != nil {
		return err
	}

	db.mem.Put(key, value)

	return nil
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
