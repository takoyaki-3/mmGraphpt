package mmgraphpt

import (
	"github.com/takoyaki-3/goraph"
)

type Stop struct {
	StopId string
	ZoneId string
	Name   string
}

type PTGraph struct {
	Map          goraph.Graph   // 任意のタイミングで移動可能な道
	Stops        map[int64]Stop // 停留所に指定されているポイント
	StopId2Place map[string]int64
	// SameNode		 map[int64][]MultiNode // マルチグラフ化した場合の共通ノードリスト
	GraphId int // マルチグラフ化した場合のグラフID
}

func NewPTGraph() *PTGraph {
	ptg := PTGraph{}
	ptg.Stops = map[int64]Stop{}
	return &ptg
}

