package memtable

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
