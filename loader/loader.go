package loader

import (
	"log"
	"io/ioutil"
	"github.com/takoyaki-3/goraph"
	"github.com/takoyaki-3/mmGraphpt/pb"
	. "github.com/takoyaki-3/mmGraphpt"
	"github.com/golang/protobuf/proto"
)

// Load Protocol Buffer
func Load(filename string) PTGraph {
	// Read the existing graph.
	in, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}
	ptgraph := &pb.PTGraph{}
	if err := proto.Unmarshal(in, ptgraph); err != nil {
		log.Fatalln("Failed to parse graph:", err)
	}

	edges := [][]goraph.Edge{}
	latlons := []goraph.LatLon{}

	for _, v := range ptgraph.Map.Edge {
		for int64(len(edges)) <= v.From {
			edges = append(edges, []goraph.Edge{})
		}
		edges[v.From] = append(edges[v.From], goraph.Edge{
			To:   v.To,
			Cost: v.Cost})
	}
	for _, v := range ptgraph.Map.Latlon {
		for int64(len(latlons)) <= v.LatlonId {
			latlons = append(latlons, goraph.LatLon{})
		}
		latlons[v.LatlonId] = goraph.LatLon{
			Lat: v.Lat,
			Lon: v.Lon}
	}

	g := PTGraph{
		Map: goraph.Graph{
			Edges:   edges,
			LatLons: latlons},
		Stops:        map[int64]Stop{},
		StopId2Place: map[string]int64{}}

	for _, v := range ptgraph.Stop {
		g.Stops[v.Id] = Stop{
			Name:   v.Name,
			StopId: v.StopId}
		g.StopId2Place[v.StopId] = v.Id
	}

	// // 共通点の読み込み
	// g.SameNode = map[int64][]MultiNode{}
	// for _,v := range ptgraph.SamePlaces{
	// 	place := v.Place
	// 	for _,v := range v.SamePlaces{
	// 		g.SameNode[place] = append(g.SameNode[place], MultiNode{
	// 			GraphId: int(v.Graphid),
	// 			Id:			 v.Place,
	// 		})
	// 	}
	// }

	return g
}

// Write to Protocol Buffer
func Write(filename string, ptg PTGraph) {
	// ...
	id := int64(0)
	edge := []*pb.Edge{}
	for k, v := range ptg.Map.Edges {
		for _, v := range v {
			edge = append(edge, &pb.Edge{
				EdgeId: id,
				From:   int64(k),
				To:     v.To,
				Cost:   v.Cost,
			})
			id++
		}
	}
	latlon := []*pb.LatLon{}
	for k, v := range ptg.Map.LatLons {
		latlon = append(latlon, &pb.LatLon{
			LatlonId: int64(k),
			Lat:      v.Lat,
			Lon:      v.Lon,
		})
	}

	// PTG
	ptgraph := &pb.PTGraph{
		Map: &pb.Graph{
			Edge:   edge,
			Latlon: latlon,
		},
	}

	for id, s := range ptg.Stops {
		ptgraph.Stop = append(ptgraph.Stop, &pb.Stop{
			Id:     id,
			Name:   s.Name,
			StopId: s.StopId,
		})
	}

	// // 共通点の書き出し
	// for place,v := range ptg.SameNode{
	// 	sps := []*pb.MultiPlace{}
	// 	for _,v:=range v{
	// 		sps = append(sps, &pb.MultiPlace{
	// 			Graphid: int32(v.GraphId),
	// 			Place: 	 v.Id,
	// 		})
	// 	}
	// 	n := pb.SamePlace{
	// 		Place: place,
	// 		SamePlaces: sps,
	// 	}
	// 	ptgraph.SamePlaces = append(ptgraph.SamePlaces, &n)
	// }

	// Write the new address book back to disk.
	out, err := proto.Marshal(ptgraph)
	if err != nil {
		log.Fatalln("Failed to encode address book:", err)
	}
	if err := ioutil.WriteFile(filename, out, 0644); err != nil {
		log.Fatalln("Failed to write address book:", err)
	}
}