package main

import (
	"fmt"
	"io"
	"me/memtable"
	"me/wal"
)

func main() {
	w, err := wal.OpenWAL("my_log.txt", ".")

	if err != nil {
		fmt.Println("打开文件失败：", err)
		return
	}

	// 1. 写两条数据
	w.Write([]byte("我是第一条机密数据"))
	w.Write([]byte("我是第二条机密数据"))

	// 2. 将文件指针移回文件开头，准备从头开始读
	// 在实际项目中，通常会专门写个 ResetReader 方法，为了简单先在在 main 里直接操作底层 file
	// (需要把 wal 结构体里的 file 改成首字母大写的 File 才能在外部访问

	sl := memtable.NewSkipList()
	sl.Put("apple", "苹果")

	// 3. 循环读取，直到读完
	for {
		data, err := w.ReadNext()
		if err == io.EOF {
			fmt.Println(" 全部读完了")
			break
		}
		if err != nil {
			fmt.Println(" 读取遇到严重错误，停止恢复:", err)
			break
		}

		fmt.Printf("成功读出数据: %s\n", string(data))
	}
}
