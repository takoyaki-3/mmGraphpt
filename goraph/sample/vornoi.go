package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/takoyaki-3/gomc"
	"github.com/takoyaki-3/mmGraphpt/goraph"
	"github.com/takoyaki-3/mmGraphpt/goraph/geometry/h3"
	"github.com/takoyaki-3/mmGraphpt/goraph/loader/osm"
	"github.com/takoyaki-3/mmGraphpt/goraph/search"
)

func main() {

	// Graph load
	fmt.Println("Graph load")
	g := osm.Load("japan-latest.osm.pbf")

	// Make index
	fmt.Println("Make index")
	h3indexes := h3.MakeH3Index(g, 9)

	// Load base point
	fmt.Println("Load base point")
	bases := []int64{}
	titles, records := gomc.ReadCSV("./base.csv")
	for _, v := range records {
		lat, _ := strconv.ParseFloat(v[titles["lat"]], 64)
		lon, _ := strconv.ParseFloat(v[titles["lon"]], 64)
		bases = append(bases, h3.Find(g, h3indexes, goraph.LatLon{lat, lon}, 9))
	}

	// Create voronoi diagram
	fmt.Println("Create voronoi diagram")
	nodes := search.Voronoi(g, bases)

	// Output
	fmt.Println("Output")
	wf, err := os.Create("./output.csv")
	if err != nil {
		log.Println(err)
	}
	defer wf.Close()

	w := csv.NewWriter(wf) // utf8
	w.Write([]string{"base", "lat", "lon"})
	for k, v := range nodes {
		if k%100 != 0 {
			continue
		}
		lat := strconv.FormatFloat(g.LatLons[k].Lat, 'f', -1, 64)
		lon := strconv.FormatFloat(g.LatLons[k].Lon, 'f', -1, 64)
		w.Write([]string{strconv.Itoa(int(v)), lat, lon})
	}
	w.Flush()
}
