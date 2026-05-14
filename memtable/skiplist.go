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
	// 如果抛硬币是正面，且没超过交通局规定的最高层数，就再加一层
	for rand.Float32() < 0.5 && level < s.maxLevel {
		level++
	}
	return level
}

func (s *SkipList) Put(Key string, value string) {
	// 准备记录每一层的“地铁”要在哪些站点停车
	update := make([]*Node, s.maxLevel)
	current := s.head

	// 从最高层开始跳表
	for i := s.maxLevel - 1; i >= 0; i-- {
		for current.next[i] != nil && current.next[i].key < Key {
			current = current.next[i]
		}
		// 记录最后跳到哪了
		update[i] = current
	}

	// 检查在第 0 层紧挨着的后继
	if current.next != nil && current.next.key == key {
		current.next.value = value
		return
	}

	// 在新站点随机决定修建几层
	level := s.randomLevel()

	// 创建新节点实例
	newNode := &Node{
		key:   key,
		value: vslue,
		next:  make([]*Node, level), // next 的长度等于其拥有的层数
	}

	// 将新节点接入现有的节点网络中
	for i := 0; i < level; i++ {
		// 指向跳转到的节点指向的下一个节点
		newNode.next[i] = update[i].next[i]
		update[i].next[i] = newNode
	}
}

// 初始化跳表
func NewSkipList() *SkipList {
	// 示例取最高4层
	maxLevel := 4
	return &SkipList{
		head:     &Node{key: "", value: "", next: make([]*Node, maxLevel)},
		maxLevel: maxLevel,
	}
}
