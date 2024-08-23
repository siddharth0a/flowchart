package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"

	"github.com/awalterschulze/gographviz"
	// 추가로 필요한 패키지 임포트
)

func parseSource(filename string) ([]Node, []Edge) {
	fset := token.NewFileSet()

	// 파일을 파싱하여 AST 생성
	node, err := parser.ParseFile(fset, filename, nil, parser.AllErrors)
	if err != nil {
		log.Fatalf("Failed to parse file: %v", err)
	}

	var nodes []Node
	var edges []Edge
	nodeMap := make(map[string]Node)

	// AST를 순회하면서 함수 선언과 함수 호출을 찾아 노드와 엣지 생성
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			// 함수 선언 노드 생성
			funcNode := Node{
				ID:    x.Name.Name,
				Label: "Func: " + x.Name.Name,
			}
			nodes = append(nodes, funcNode)
			nodeMap[x.Name.Name] = funcNode

		case *ast.CallExpr:
			// 함수 호출 노드 생성
			if ident, ok := x.Fun.(*ast.Ident); ok {
				caller := ident.Name

				// 호출된 함수가 이미 노드로 존재하는지 확인하고 없다면 추가
				if _, exists := nodeMap[caller]; !exists {
					calledNode := Node{
						ID:    caller,
						Label: "Func: " + caller,
					}
					nodes = append(nodes, calledNode)
					nodeMap[caller] = calledNode
				}

				edge := Edge{
					From: nodeMap[fset.Position(x.Pos()).String()], // 수정된 부분
					To:   nodeMap[caller],
				}
				edges = append(edges, edge)

			}
		}

		return true
	})

	return nodes, edges
}

// 플로우차트를 시각화하는 함수
func generateGraph(nodes []Node, edges []Edge) string {
	graphAst, _ := gographviz.ParseString(`digraph G {}`)
	graph := gographviz.NewGraph()
	gographviz.Analyse(graphAst, graph)

	for _, node := range nodes {
		graph.AddNode("G", node.ID, map[string]string{
			"label": node.Label,
		})
	}

	for _, edge := range edges {
		graph.AddEdge(edge.From.ID, edge.To.ID, true, nil)
	}

	return graph.String() // Graphviz dot 언어로 반환
}
