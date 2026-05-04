package main

import (
	"fmt"
	"me/wal"
)

func main() {
	w, err := wal.OpenWAL("my_log.txt", ".")

	if err != nil {
		fmt.Println("工作台准备失败：", err)
		return
	}

	err = w.Write([]byte("Hello database"))

	if err != nil {
		fmt.Println("写日志失败：", err)
		return
	}

	fmt.Println("数据写入成功！")
}
