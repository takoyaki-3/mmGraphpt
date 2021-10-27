package main

import (
	csvloader "github.com/takoyaki-3/mmGraphpt/goraph/loader/csv"
	"github.com/takoyaki-3/mmGraphpt/goraph/loader/geojson"
)

func main() {
	g := geojson.Load("N02-19_RailroadSection.geojson")
	csvloader.WriteEdge("edge.csv", g)
	csvloader.WriteLatLon("latlon.csv", g)
}
