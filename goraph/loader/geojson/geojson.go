package geojson

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	// "fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/takoyaki-3/mmGraphpt/goraph"
	"github.com/takoyaki-3/mmGraphpt/goraph/geometry"
)

type Sha256 [32]byte

func Node2Hash(node goraph.LatLon) Sha256 {
	flat := strconv.FormatFloat(node.Lat, 'f', -1, 64)
	flon := strconv.FormatFloat(node.Lon, 'f', -1, 64)
	b := bytes.Join([][]byte{[]byte(flat), []byte(flon)}, []byte{})
	return sha256.Sum256(b)
}

func random(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()*(max-min) + min
}

func Load(filename string) goraph.Graph {
	g := goraph.Graph{}

	latlon2id := map[Sha256]int{}

	// JSONファイル読み込み
	rawJSON, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	fc, _ := geojson.UnmarshalFeatureCollection(rawJSON)

	rand.Seed(time.Now().UnixNano())

	func() {
		for _, f := range fc.Features {
			line := []int{}
			for k, v := range f.Geometry.(orb.LineString) {
				node := goraph.LatLon{v[1], v[0]}
				h := Node2Hash(node)

				// fmt.Println(k)
				if k < 5 || k > len(f.Geometry.(orb.LineString))-6 {
					if id, ok := latlon2id[h]; !ok {
						id = int(g.AddLatLon(node))
						latlon2id[h] = id
						line = append(line, id)
					}
					line = append(line, latlon2id[h])
				} else {
					line = append(line, int(g.AddLatLon(node)))
				}
			}

			for k, _ := range line {
				if k == 0 {
					continue
				}
				e := goraph.Edge{}

				node1, node2 := int64(line[k-1]), int64(line[k])
				for len(g.Edges) <= int(node1) || len(g.Edges) <= int(node2) {
					g.Edges = append(g.Edges, []goraph.Edge{})
				}

				e.Cost = geometry.HubenyDistance(g.LatLons[node1], g.LatLons[node2])
				e.To = node2
				g.Edges[node1] = append(g.Edges[node1], e)
				e.To = node1
				g.Edges[node2] = append(g.Edges[node2], e)
			}
		}
	}()

	fmt.Println("ok")

	return g
}

func LoadCurve(filename string) goraph.Graph {
	g := goraph.Graph{}

	latlon2id := map[Sha256]int{}

	// JSONファイル読み込み
	rawJSON, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	fc, _ := geojson.UnmarshalFeatureCollection(rawJSON)

	func() {
		for _, f := range fc.Features {
			line := []int{}

			for _, v := range f.Geometry.(orb.LineString) {
				node := goraph.LatLon{v[1], v[0]}
				h := Node2Hash(node)
				if id, ok := latlon2id[h]; !ok {
					id = len(g.LatLons)
					latlon2id[h] = id
					g.LatLons = append(g.LatLons, node)
					line = append(line, id)
				}
				line = append(line, latlon2id[h])
			}

			if len(line) == 0 {
				continue
			}

			e1, e2 := goraph.Edge{}, goraph.Edge{}
			e1.To = int64(line[len(line)-1])
			e2.To = int64(line[0])

			for k, _ := range line {
				if k == 0 {
					continue
				}

				node1, node2 := int64(line[k-1]), int64(line[k])
				if node1 == node2 {
					continue
				}

				e1.Cost += geometry.HubenyDistance(g.LatLons[node1], g.LatLons[node2])

				if node1 != e1.To {
					e1.LatLons = append(e1.LatLons, node1)
				}
				if node2 != e2.To {
					e2.LatLons = append([]int64{node2}, e2.LatLons...)
				}
			}

			e2.Cost = e1.Cost

			// fmt.Println(e1)

			g.AddEdge(e1, e2.To)
			g.AddEdge(e2, e1.To)
		}
	}()

	return g
}
