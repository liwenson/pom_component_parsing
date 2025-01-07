package model

import "encoding/json"

// IsOnline 结构体用于表示依赖项是否在线，包括有效性标志和值
type IsOnline struct {
	Valid bool `json:"valid"` // 标识该字段是否有效
	Value bool `json:"value"` // 表示依赖项是否在线，true 表示在线，false 表示离线
}

// UnmarshalJSON 实现自定义的 JSON 反序列化逻辑
func (i *IsOnline) UnmarshalJSON(bytes []byte) error {
	i.Valid = true
	return json.Unmarshal(bytes, &i.Value)
}

// MarshalJSON 实现自定义的 JSON 序列化逻辑
func (i *IsOnline) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("true"), nil // 默认返回 true
	}
	return json.Marshal(i.Value)
}

// SetOnline 设置 IsOnline 的值并标记为有效
func (i *IsOnline) SetOnline(b bool) {
	i.Valid = true
	i.Value = b
}

// 确保 IsOnline 实现了 json.Marshaler 接口
var _ json.Marshaler = (*IsOnline)(nil)

// 确保 IsOnline 实现了 json.Unmarshaler 接口
var _ json.Unmarshaler = (*IsOnline)(nil)

// IsOnlineTrue 返回一个标记为有效且值为 true 的 IsOnline 实例
func IsOnlineTrue() IsOnline {
	return IsOnline{
		Value: true,
		Valid: true,
	}
}

// IsOnlineFalse 返回一个标记为有效且值为 false 的 IsOnline 实例
func IsOnlineFalse() IsOnline {
	return IsOnline{
		Value: false,
		Valid: true,
	}
}
