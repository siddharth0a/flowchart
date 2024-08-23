package main

// 노드(Node) 및 엣지(Edge) 구조체 정의
type Node struct {
	ID    string
	Label string
}

type Edge struct {
	From Node
	To   Node
}

// 그래프(플로우차트)를 나타내는 데이터 구조
type Graph struct {
	Nodes []Node
	Edges []Edge
}
