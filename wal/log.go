package wal

import (
	"encoding/binary"
	"hash/crc32"
	"os"
	"sync"
)

type WAL struct {
	file *os.File
	dir  *os.File
	mu   sync.Mutex
}

func OpenWAL(filePath string, dirPath string) (*WAL, error) {
	//打开日志
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	//打开目录
	d, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}

	//返回需要用到的工具
	return &WAL{
		file: f,
		dir:  d,
	}, nil
}

func (w *WAL) Write(data []byte) error {
	// 先上锁
	w.mu.Lock()
	defer w.mu.Unlock()

	// 使用 Castagnoli 多项式创建一个 CRC32 计算表
	crcTable := crc32.MakeTable(crc32.Castagnoli)

	// 计算我们要写入数据的特征码，得到一个 uint32 类型的数字
	checksum := crc32.Checksum(data, crcTable)

	// 准备一个长度为 4 的字节切片，把 uint32 数字装进去
	crcBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(crcBytes, checksum)

	// 把 4 字节的“印章”（校验码） 写进文件
	_, err := w.file.Write(crcBytes)
	if err != nil {
		return err
	}

	// 再把真正的数据写进文件
	_, err = w.file.Write(data)
	if err != nil {
		return err
	}
	// 强制刷盘
	err = w.file.Sync()
	if err != nil {
		return err
	}

	return nil
}

// func (w *WAL) Write(data []byte) error {
// 	// 先上锁
// 	w.mu.Lock()
// 	defer w.mu.Unlock()

// 	// 再写数据
// 	_, err := w.file.Write(data)

// 	if err != nil {
// 		return err
// 	}

// 	// 强制刷盘
// 	err = w.file.Sync()

// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
