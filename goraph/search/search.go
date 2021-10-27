package search

import (
	"math"

	"github.com/takoyaki-3/mmGraphpt/goraph"
	"github.com/takoyaki-3/mmGraphpt/goraph/pkg/priority_queue"
)

type Query struct {
	From int64
	To   int64
}

type Output struct {
	Nodes []int64
	Cost  float64
}

//
func Search(g goraph.Graph, query Query) Output {

	l := len(g.Edges)
	if l < len(g.LatLons) {
		l = len(g.LatLons)
	}

	q := priority_queue.NewMinSet()
	cost := make([]float64, l)
	flag := make([]bool, l)
	before := make([]int64, l)

	for k, _ := range cost {
		cost[k] = math.MaxFloat64
	}

	if len(g.Edges) <= int(query.From) || len(g.Edges) <= int(query.To) {
		return Output{}
	}

	cost[query.From] = 0.0
	before[query.From] = -2

	q.AddVal(query.From, 0.0)

	var pos int64
	for q.Len() > 0 {
		pos = q.GetMin()
		if flag[pos] {
			continue
		}
		flag[pos] = true

		if pos == query.To {
			break
		}

		for _, e := range g.Edges[pos] {
			eto := e.To
			if flag[eto] {
				continue
			}
			if cost[eto] <= cost[pos]+e.Cost {
				continue
			}
			cost[eto] = cost[pos] + e.Cost
			if len(e.LatLons) == 0 {
				before[eto] = pos
			} else {
				if eto != e.LatLons[len(e.LatLons)-1] {
					before[eto] = e.LatLons[len(e.LatLons)-1]
				}
				if e.LatLons[0] != pos {
					before[e.LatLons[0]] = pos
				}
				for k, v := range e.LatLons {
					if k == 0 {
						continue
					}
					if v != e.LatLons[k-1] {
						before[v] = e.LatLons[k-1]
					}
				}
			}
			q.AddVal(eto, cost[eto])
		}
	}

	// 出力
	out := Output{}
	out.Cost = cost[pos]
	out.Nodes = append(out.Nodes, pos)

	bef := before[pos]
	for bef != -2 {
		out.Nodes = append([]int64{bef}, out.Nodes...)
		bef = before[bef]
	}
	return out
}

func Voronoi(g goraph.Graph, bases []int64) map[int64]int64 {
	// initialization
	q := priority_queue.NewMinSet()
	cost := make([]float64, len(g.Edges))
	flag := make([]bool, len(g.Edges))
	start_group := map[int64]int64{}

	counter := int64(0)
	for k, _ := range cost {
		cost[k] = math.MaxFloat64
	}

	for _, v := range bases {
		cost[v] = 0.0
		q.AddVal(v, 0.0)
		start_group[int64(v)] = counter % 20
		counter++
	}

	for q.Len() > 0 {
		pos := q.GetMin()
		if flag[pos] {
			continue
		}
		flag[pos] = true

		// グラフ拡張処理
		for _, e := range g.Edges[pos] {
			eto := e.To
			if flag[eto] {
				continue
			}
			if cost[eto] <= cost[pos]+e.Cost {
				continue
			}
			cost[eto] = cost[pos] + e.Cost
			start_group[eto] = start_group[pos]
			q.AddVal(eto, cost[pos]+e.Cost)
		}
	}

	return start_group
}

func AllDistance(g goraph.Graph, base []int64) []float64 {
	q := priority_queue.NewMinSet()
	cost := make([]float64, len(g.Edges))
	flag := make([]bool, len(g.Edges))

	for k, _ := range cost {
		cost[k] = math.MaxFloat64
	}

	for _, v := range base {
		cost[v] = 0.0
		q.AddVal(v, 0.0)
	}

	for q.Len() > 0 {
		pos := q.GetMin()
		if flag[pos] {
			continue
		}
		flag[pos] = true
		for _, e := range g.Edges[pos] {
			eto := e.To
			if flag[eto] {
				continue
			}
			if cost[eto] <= cost[pos]+e.Cost {
				continue
			}
			cost[eto] = cost[pos] + e.Cost
			q.AddVal(eto, cost[eto])
		}
	}
	return cost
}
