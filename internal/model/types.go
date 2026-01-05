package model

// 对应JSON结构的数据类型
type ItemNode struct {
	Title string `json:"title"`
	CID   uint64 `json:"cid"`
	P     uint32 `json:"p"`
}

type TabNode struct {
	Name  string     `json:"name"`
	Items []ItemNode `json:"items"`
}

type TitleNode struct {
	Name string    `json:"name"`
	P    *uint32   `json:"p"`
	Tabs []TabNode `json:"tabs"`
	Open bool      `json:"-"` // 运行时状态，不序列化
}

type GroupNode struct {
	Name   string      `json:"name"`
	Titles []TitleNode `json:"titles"`
	Open   bool        `json:"-"` // 运行时状态，不序列化
}

// 根结构
type TreeData struct {
	Groups []GroupNode `json:"groups"`
}
