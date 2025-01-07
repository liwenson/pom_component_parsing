package utils

// String 返回一个指向给定字符串的指针。
// 参数 a: 一个字符串。
// 返回值: 字符串 a 的指针。
func String(a string) *string {
	return &a
}

// KeysOfMap 获取映射 (map) 中所有的键，并以切片的形式返回。
// 参数 m: 一个 map 类型的变量，其中 K 是键的类型，V 是值的类型。
// 返回值: 一个包含 map 中所有键的切片。
func KeysOfMap[K comparable, V any](m map[K]V) []K {
	// 创建一个切片 r 用来存储所有的键，初始容量设置为 map 的长度，以优化性能
	var r = make([]K, 0, len(m))
	// 遍历 map 中的每一个键，将其添加到切片 r 中
	for k := range m {
		r = append(r, k)
	}
	// 返回包含所有键的切片
	return r
}

// DistinctSlice 返回一个去重后的切片。
// 参数 input: 一个切片，其中的元素类型为 T。
// 返回值: 一个去重后的切片，包含了原切片中不重复的元素。
func DistinctSlice[T comparable](input []T) []T {
	// 如果输入的切片为 nil，直接返回 nil
	if input == nil {
		return nil
	}
	// 创建一个空切片 r，用来存储去重后的结果
	var r = make([]T, 0)
	// 创建一个 map 来存储元素，确保每个元素唯一
	var m = make(map[T]struct{})
	// 遍历输入切片中的每个元素
	for _, i := range input {
		// 如果元素已存在于 map 中，跳过
		if _, ok := m[i]; ok {
			continue
		}
		// 如果元素不在 map 中，加入 map 并添加到结果切片 r 中
		m[i] = struct{}{}
		r = append(r, i)
	}
	// 返回去重后的切片
	return r
}
