package wal

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"os"
	"sync"
)

// 定义用来检测不完整写入的“尾部金丝雀”魔法数字
const CanaryMagicNumber uint64 = 0xDEADBEEFFEEDFACE

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

	// 计算 CRC32
	crcTable := crc32.MakeTable(crc32.Castagnoli)
	checksum := crc32.Checksum(data, crcTable)
	crcBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(crcBytes, checksum)

	// 计算数据的字节长度
	length := uint32(len(data))
	lenBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(lenBytes, length)

	// 准备金丝雀，将魔法数字转换为8个字节
	canaryBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(canaryBytes, CanaryMagicNumber)

	// 按照顺序写入硬盘
	// 先写 CRC32 （4字节）
	if _, err := w.file.Write(crcBytes); err != nil {
		return nil
	}

	// 再写数据长度（4字节）
	if _, err := w.file.Write(lenBytes); err != nil {
		return nil
	}

	//写入真实数据（N字节）
	if _, err := w.file.Write(data); err != nil {
		return nil
	}

	// 写入尾部金丝雀（8字节）
	if _, err := w.file.Write(canaryBytes); err != nil {
		return nil
	}

	// 强制刷盘
	if err := w.file.Sync(); err != nil {
		return nil
	}

	return nil
}

// ReadNext 从日志中读取下一条完整记录
func (w *WAL) ReadNext() ([]byte, error) {
	// 先读取8字节的头部（4字节CRC + 4字节长度）
	header := make([]byte, 8)
	_, err := io.ReadFull(w.file, header)
	if err != nil {
		// 如果读到文件末尾，会返回 io.EOF,说明读完了
		return nil, err
	}

	// 从头部解析出期望的 CRC 和数据长度
	expectedCRC := binary.LittleEndian.Uint32(header[0:4])
	length := binary.LittleEndian.Uint32(header[4:8])

	// 防止内存撑爆（Double-Read 校验）
	//规定一条数据最大不超过 10 MB
	if length > 10*1024*1024 {
		return nil, errors.New("数据长度异常，拒绝读取")
	}

	// 根据安全的长度读取真正数据
	data := make([]byte, length)
	_, err = io.ReadFull(w.file, data)
	if err != nil {
		return nil, err
	}

	// 读取8字节的尾部金丝雀
	canaryBytes := make([]byte, 8)
	_, err = io.ReadFull(w.file, canaryBytes)
	if err != nil {
		return nil, err
	}

	// 解析金丝雀并验证
	canary := binary.LittleEndian.Uint64(canaryBytes)
	if canary != CanaryMagicNumber {
		// 如果金丝雀不对，说明当时写到一半停电了
		return nil, errors.New("校验失败，数据记录不全")
	}

	// 数据完好则重新计算 CRC32 校验和，与头部记录的是否一致
	crcTable := crc32.MakeTable(crc32.Castagnoli)
	actualCRC := crc32.Checksum(data, crcTable)
	if actualCRC != expectedCRC {
		return nil, errors.New("CRC32 校验失败")
	}

	return data, nil
}

// SeekToStart 将日记本的翻页指针重新拨回第一页的第一行
func (w *WAL) SeekToStart() error {
	// 0 表示偏移量为 0，io.SeekStart 表示从文件最开头算起
	_, err := w.file.Seek(0, io.SeekStart)
	return err
}
