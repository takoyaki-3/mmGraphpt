package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/takoyaki-3/mmGraphpt/goraph"
	"github.com/takoyaki-3/mmGraphpt/goraph/geometry"
	"github.com/takoyaki-3/mmGraphpt/goraph/geometry/h3"
	"github.com/takoyaki-3/mmGraphpt/goraph/loader"

	// "github.com/takoyaki-3/mmGraphpt/goraph/loader/osm"
	// "github.com/takoyaki-3/mmGraphpt/goraph/loader/geojson"
	"github.com/takoyaki-3/mmGraphpt/goraph/search"
)

func main() {

	// Graph load
	fmt.Println("Graph load")
	g := loader.Load("./kanto.goraph.pbf")
	// g := osm.Load("./kanto-latest.osm.pbf")
	// g := geojson.Load("./kanto-lines.geojson")

	// Make index
	fmt.Println("Make index")
	h3indexes := h3.MakeH3Index(g, 6)

	// Start server
	fmt.Println("start Server")
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		v := r.URL.Query()
		if v == nil {
			return
		}
		parm_v := map[string]float64{}
		must_parms := []string{"from_lat", "from_lon", "to_lat", "to_lon"}
		for _, k := range must_parms {
			if _, ok := v[k]; !ok {
				fmt.Fprintf(w, "{\"ErrorMessage\":\"Required parameters do not exist.\"}")
				return
			}
			f, err := strconv.ParseFloat(v[k][0], 64)
			if err != nil {
				log.Fatal(err)
				fmt.Fprintf(w, "{\"ErrorMessage\":\"Required parameters cannot be converted.\"}")
				return
			}
			parm_v[k] = f
		}

		q := search.Query{}

		// Find node
		q.To = h3.Find(g, h3indexes, goraph.LatLon{parm_v["from_lat"], parm_v["from_lon"]}, 6)
		q.From = h3.Find(g, h3indexes, goraph.LatLon{parm_v["to_lat"], parm_v["to_lon"]}, 6)

		// Search
		o := search.Search(g, q)

		// Make GeoJSON
		rawJSON := geometry.MakeLineString(g, o.Nodes)

		// Response
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintln(w, rawJSON)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bytes, err := ioutil.ReadFile("./index.html")
		if err != nil {
			panic(err)
		}
		fmt.Fprintln(w, string(bytes))
	})
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
