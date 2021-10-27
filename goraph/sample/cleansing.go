package main

import (
	"fmt"

	"github.com/takoyaki-3/mmGraphpt/goraph"
	"github.com/takoyaki-3/mmGraphpt/goraph/loader"
	"github.com/takoyaki-3/mmGraphpt/pkg"
)

func main() {

	// グラフデータの読み込み処理
	paths := pkg.FindFiles("../../sample/graphdata", ".goraph.pbf")
	for _, fileName := range paths {
		// goraphの読み込み
		g := loader.Load(fileName)

		flag := make([]int, len(g.LatLons))
		nums := []int{}

		for pos, _ := range g.LatLons {
			if flag[int(pos)] != 0 {
				continue
			}
			nums = append(nums, 0)
			fmt.Println(len(nums))
			stack := []int64{int64(pos)}
			for {
				if len(stack) == 0 {
					break
				}
				pos := stack[0]
				stack = stack[1:]
				if flag[pos] != 0 {
					continue
				}
				flag[pos] = len(nums)

				if len(g.Edges) <= int(pos) {
					fmt.Println("over", len(g.Edges), pos)
					continue
				}
				for _, e := range g.Edges[pos] {
					if flag[e.To] == 0 {
						stack = append(stack, e.To)
					}
				}
			}
		}
		nums = append(nums, 0)
		for k, _ := range g.LatLons {
			nums[flag[k]]++
		}

		fmt.Println(nums)
		maxFlag := 0
		for k, v := range nums {
			if nums[maxFlag] < v {
				maxFlag = k
			}
		}

		newG := goraph.Graph{}
		old2new := map[int]int{}
		for k, v := range g.LatLons {
			if flag[k] == maxFlag {
				old2new[k] = len(newG.LatLons)
				newG.LatLons = append(newG.LatLons, v)
			}
		}
		for k, _ := range g.LatLons {
			if flag[k] == maxFlag {
				for j, edge := range g.Edges[k] {
					v := old2new[int(edge.To)]
					g.Edges[k][j].To = int64(v)
				}
				newG.Edges = append(newG.Edges, g.Edges[k])
			}
		}
		loader.Write(fileName, newG)
	}
}
