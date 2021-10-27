package osm

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"github.com/cheggaaa/pb"
	"github.com/dustin/go-humanize"
	"github.com/qedus/osmpbf"
	"github.com/takoyaki-3/mmGraphpt/goraph"
	"github.com/takoyaki-3/mmGraphpt/goraph/geometry"
)

func Load(filename string) goraph.Graph {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	stat, _ := f.Stat()
	filesiz := int(stat.Size() / 1024)

	d := osmpbf.NewDecoder(f)
	err = d.Start(runtime.GOMAXPROCS(-1))
	if err != nil {
		log.Fatal(err)
	}

	exclusionList := map[string]bool{}
	exclusionList["motorway"] = true
	exclusionList["bus_guideway"] = true
	exclusionList["raceway"] = true
	exclusionList["busway"] = true
	exclusionList["cycleway"] = true
	exclusionList["proposed"] = true
	exclusionList["construction"] = true
	exclusionList["motorway_junction"] = true
	exclusionList["platform"] = true

	// 一時記憶用変数
	latlons := map[int64]goraph.LatLon{}
	ways := map[int64][]int64{}
	usednode := map[int64]bool{}

	nc, wc, rc := int64(0), int64(0), int64(0)
	pb.New(filesiz).SetUnits(pb.U_NO)
	for i := 0; ; i++ {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				latlons[v.ID] = goraph.LatLon{v.Lat, v.Lon}
				nc++
			case *osmpbf.Way:
				if _, ok := v.Tags["highway"]; !ok {
					continue
				}
				if v, ok := exclusionList[v.Tags["highway"]]; ok {
					if v {
						continue
					}
				}

				nodes := []int64{}
				for _, v := range v.NodeIDs {
					nodes = append(nodes, v)
					usednode[v] = true
				}
				ways[v.ID] = nodes
				wc++
			case *osmpbf.Relation:
				rc++
			default:
				log.Fatalf("unknown type %T\n", v)
			}
		}
	}
	fmt.Printf("Nodes: %s, Ways: %s, Relations: %s\n", humanize.Comma(nc), humanize.Comma(wc), humanize.Comma(rc))

	g := goraph.Graph{}
	nodeid := NewReplace()

	for _, v := range ways {
		for i := 1; i < len(v); i++ {
			e := goraph.Edge{}
			e.Cost = geometry.HubenyDistance(latlons[v[i-1]], latlons[v[i]])
			node1 := nodeid.AddReplace(v[i-1])
			node2 := nodeid.AddReplace(v[i])
			for len(g.Edges) <= int(node1) || len(g.Edges) <= int(node2) {
				g.Edges = append(g.Edges, []goraph.Edge{})
			}
			e.To = node2
			g.Edges[node1] = append(g.Edges[node1], e)
			e.To = node1
			g.Edges[node2] = append(g.Edges[node2], e)
		}
	}
	for k, v := range latlons {
		if _, ok := usednode[k]; !ok {
			continue
		}
		id := nodeid.AddReplace(int64(k))
		for len(g.LatLons) <= int(id) {
			g.LatLons = append(g.LatLons, goraph.LatLon{})
		}
		g.LatLons[id] = v
		for len(g.Edges) <= int(id) {
			g.Edges = append(g.Edges, []goraph.Edge{})
		}
	}

	return g
}

// 置き換え関数
type Replace struct {
	Id2Str map[int64]int64
	Str2Id map[int64]int64
}

func (s *Replace) AddReplace(str int64) int64 {
	if val, ok := s.Str2Id[str]; ok {
		return val
	}
	id := int64(len(s.Str2Id) + 1)
	s.Str2Id[str] = id
	s.Id2Str[id] = str
	return id
}

func (s *Replace) AddReplaceIndex(str int64, index int64) int64 {
	if val, ok := s.Str2Id[str]; ok {
		return val
	}
	s.Str2Id[str] = index
	s.Id2Str[index] = str
	return index
}

func NewReplace() *Replace {
	s := new(Replace)
	s.Str2Id = map[int64]int64{}
	s.Id2Str = map[int64]int64{}
	return s
}
