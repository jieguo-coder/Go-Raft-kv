# LSM-Tree 存储引擎 — 练习进度

## 当前状态

### 已完成模块
- **WAL** (`wal/log.go`): 追加写（CRC32 + 长度前缀 + 金丝雀 + fsync）、顺序读、SeekToStart
- **MemTable** (`memtable/skiplist.go`): 跳表实现，支持 Put / Get
- **DB** (`db/db.go`): 组合 WAL + MemTable，Get / Put / NewDB

### 本次会话完成
- **序列化格式改进**: Put 方法从 `key:value` 分隔符方案改为 **长度前缀方案**
  - 格式: `[key长度 4字节 LittleEndian][key字节][value字节]`
  - 解决了 key/value 含 `:` 时的解析歧义
  - Put 签名改为返回 error
  - 用 `append` 拼接，代码简洁可读

### 关键技术讨论
- `append` vs 手动预分配内存：WAL 路径是 I/O 密集型，fsync 占主导，内存分配开销可忽略
- 长度前缀 vs 分隔符：固定宽度头部消除解析歧义
- 小端序：跟随 x86 CPU 原生字节序，效率最高

---

## 下一步计划

### 主线任务：NewDB + WAL 恢复

**1. 反序列化函数**（WAL 恢复的前置条件）
- 写一个函数，给定 WAL 读出的 `[]byte`，按长度前缀格式拆出 key 和 value
- 格式: 读前 4 字节 → key 长度 N → 读 N 字节 → key → 剩余 → value

**2. NewDB 加入恢复逻辑**
- 打开 WAL 后，调用 SeekToStart 从头读
- 循环 ReadNext，对每条记录调用反序列化函数，回放到 memtable
- 遇到 io.EOF 时正常结束恢复
- 注意 O_APPEND 模式：读取完后的文件指针位置不影响后续写入（O_APPEND 会自动定位到文件末尾写）

### 待讨论的问题
- 用户上次的 3 个思考题（Q1/Q2/Q3）中，Q2 和 Q3 还没有完整回答
- Q2: O_APPEND 对读写流程的影响
- Q3: WAL 文件为空时恢复逻辑如何处理
