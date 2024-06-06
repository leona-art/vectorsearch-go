package main

import (
	"fmt"
	"vectorsearch-go/infra"
)

func main() {
	ai, err := infra.NewGenAi()
	if err != nil {
		panic(err)
	}
	e1, err := ai.Embedding("Vertex AI は、主要な基盤モデルのための API と、迅速なプロトタイプ作成、独自のデータによるモデルの調整、アプリケーションにシームレスにデプロイするためのツールを提供します。")
	if err != nil {
		panic(err)
	}
	e2, err := ai.Embedding(``)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 768; i++ {
		if e1[i] != e2[i] {
			fmt.Println("Different")
			return
		}
	}
}
