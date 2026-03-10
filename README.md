# json

[![Go Reference](https://pkg.go.dev/badge/github.com/zhgqiang/json.svg)](https://pkg.go.dev/github.com/zhgqiang/json)
[![Go Report Card](https://goreportcard.com/badge/github.com/zhgqiang/json)](https://goreportcard.com/report/github.com/zhgqiang/json)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

English | [中文](#中文文档)

A Go library providing JSON operations and deep copy utilities.

## Installation

```bash
go get github.com/zhgqiang/json
```

## Features

### JSON Operations

Standard library compatible JSON functions:

```go
import "github.com/zhgqiang/json"

// Encoding
data, err := json.Marshal(obj)
data, err := json.MarshalIndent(obj, "", "  ")
str := json.MarshalToString(obj)  // Returns string directly

// Decoding
err := json.Unmarshal(data, &obj)

// Streaming
decoder := json.NewDecoder(reader)
encoder := json.NewEncoder(writer)
```

### Object Copy

#### CopyByJson - Deep copy via JSON serialization

```go
src := map[string]int{"a": 1, "b": 2}
var dst map[string]int
err := json.CopyByJson(&dst, src)
```

#### Copy - Struct field copy

Uses `github.com/jinzhu/copier` for struct field copying:

```go
type User struct {
    Name string
    Age  int
}
src := User{Name: "Alice", Age: 30}
var dst User
err := json.Copy(&dst, src)
```

### Map Operations

#### MapToStruct - Map to Struct conversion

```go
type Config struct {
    Host string `mapstructure:"host"`
    Port int    `mapstructure:"port"`
}
src := map[string]interface{}{"host": "localhost", "port": 8080}
var dst Config
err := json.MapToStruct(&dst, src)
```

#### MapCopy - Shallow copy

```go
src := map[string]interface{}{"a": 1}
dst := json.MapCopy(src)  // Shallow copy, nested maps share references
```

#### DeepCopyMap - Deep copy

Recursively deep copies map, supporting nested maps, slices, pointers, and structs:

```go
src := map[string]interface{}{
    "a": 1,
    "b": map[string]interface{}{"c": 2},
    "d": []int{1, 2, 3},
}
dst := json.DeepCopyMap(src)  // Completely independent copy
dst["b"].(map[string]interface{})["c"] = 999  // Does not affect src
```

### Deep Copy Functions

```go
// Deep copy slice
newSlice := json.DeepCopySlice(slice)

// Deep copy pointer
newPtr := json.DeepCopyPointer(ptr)

// Deep copy struct
newStruct := json.DeepCopyStruct(structure)

// Deep copy via gob encoding
err := json.DeepCopyGob(src, &dst)
```

## API Reference

| Function | Description |
|----------|-------------|
| `Marshal` | JSON encode to bytes |
| `Unmarshal` | JSON decode from bytes |
| `MarshalIndent` | JSON encode with indentation |
| `MarshalToString` | JSON encode to string |
| `NewDecoder` | Create JSON decoder |
| `NewEncoder` | Create JSON encoder |
| `CopyByJson` | Deep copy via JSON serialization |
| `Copy` | Struct field copy using copier |
| `MapToStruct` | Convert map to struct |
| `MapCopy` | Shallow copy a map |
| `DeepCopyMap` | Deep copy a map recursively |
| `DeepCopySlice` | Deep copy a slice |
| `DeepCopyPointer` | Deep copy a pointer |
| `DeepCopyStruct` | Deep copy a struct |
| `DeepCopyGob` | Deep copy via gob encoding |

## Dependencies

- `encoding/json` - Standard JSON encoding/decoding
- `github.com/json-iterator/go` - High-performance JSON library
- `github.com/jinzhu/copier` - Struct copying
- `github.com/mitchellh/mapstructure` - Map to struct conversion

## License

[MIT](LICENSE)

---

# 中文文档

Go JSON 工具库，提供 JSON 操作和对象深拷贝功能。

## 安装

```bash
go get github.com/zhgqiang/json
```

## 功能

### JSON 操作

标准库兼容的 JSON 操作函数：

```go
import "github.com/zhgqiang/json"

// 编码
data, err := json.Marshal(obj)
data, err := json.MarshalIndent(obj, "", "  ")
str := json.MarshalToString(obj)  // 直接返回字符串

// 解码
err := json.Unmarshal(data, &obj)

// 流式操作
decoder := json.NewDecoder(reader)
encoder := json.NewEncoder(writer)
```

### 对象复制

#### CopyByJson - 通过 JSON 序列化进行深拷贝

```go
src := map[string]int{"a": 1, "b": 2}
var dst map[string]int
err := json.CopyByJson(&dst, src)
```

#### Copy - 结构体复制

使用 `github.com/jinzhu/copier` 进行结构体字段复制：

```go
type User struct {
    Name string
    Age  int
}
src := User{Name: "Alice", Age: 30}
var dst User
err := json.Copy(&dst, src)
```

### Map 操作

#### MapToStruct - Map 转 Struct

```go
type Config struct {
    Host string `mapstructure:"host"`
    Port int    `mapstructure:"port"`
}
src := map[string]interface{}{"host": "localhost", "port": 8080}
var dst Config
err := json.MapToStruct(&dst, src)
```

#### MapCopy - 浅拷贝

```go
src := map[string]interface{}{"a": 1}
dst := json.MapCopy(src)  // 浅拷贝，嵌套 map 会共享引用
```

#### DeepCopyMap - 深拷贝

递归深拷贝 map，支持嵌套 map、slice、指针和结构体：

```go
src := map[string]interface{}{
    "a": 1,
    "b": map[string]interface{}{"c": 2},
    "d": []int{1, 2, 3},
}
dst := json.DeepCopyMap(src)  // 完全独立副本
dst["b"].(map[string]interface{})["c"] = 999  // 不影响 src
```

### 深拷贝函数

```go
// 深拷贝 slice
newSlice := json.DeepCopySlice(slice)

// 深拷贝指针
newPtr := json.DeepCopyPointer(ptr)

// 深拷贝结构体
newStruct := json.DeepCopyStruct(structure)

// 通过 gob 编码进行深拷贝
err := json.DeepCopyGob(src, &dst)
```

## API 参考

| 函数 | 说明 |
|------|------|
| `Marshal` | JSON 编码为字节 |
| `Unmarshal` | JSON 解码 |
| `MarshalIndent` | JSON 编码（带缩进） |
| `MarshalToString` | JSON 编码为字符串 |
| `NewDecoder` | 创建 JSON 解码器 |
| `NewEncoder` | 创建 JSON 编码器 |
| `CopyByJson` | 通过 JSON 序列化深拷贝 |
| `Copy` | 结构体字段复制 |
| `MapToStruct` | Map 转 Struct |
| `MapCopy` | Map 浅拷贝 |
| `DeepCopyMap` | Map 深拷贝 |
| `DeepCopySlice` | Slice 深拷贝 |
| `DeepCopyPointer` | 指针深拷贝 |
| `DeepCopyStruct` | 结构体深拷贝 |
| `DeepCopyGob` | 通过 gob 编码深拷贝 |

## 依赖

- `encoding/json` - 标准 JSON 编解码
- `github.com/json-iterator/go` - 高性能 JSON 库
- `github.com/jinzhu/copier` - 结构体复制
- `github.com/mitchellh/mapstructure` - Map 转 Struct

## 许可证

[MIT](LICENSE)