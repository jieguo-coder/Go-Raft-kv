package memtable

import (
	"math/rand"
)

// Node 相当于跳表里的“地铁站”（存实际的 Key 和 value）
type Node struct {
	key   string
	value string
	next  []*Node // 跳表的核心
}

// SkipList 代表整个 memtable
type SkipList struct {
	head     *Node
	maxLevel int
}

// randLevel 抛硬币决定“地铁”建几层（随机算法）
func (s *SkipList) randomLevel() int {
	level := 1
	// 如果抛硬币是正面 (概率0.5)，且没超过交通局规定的最高层数，就再加一层
	for rand.Float32() < 0.5 && level < s.maxLevel {
		level++
	}
	return level
}
