package main

import (
	"fmt"
	"github.com/takoyaki-3/mmGraphpt/loader"
	"github.com/takoyaki-3/mmGraphpt/converter"
)

func main(){
	converter.Osm2goraph("map.osm.pbf","map.goraph.pbf")

	ptg := loader.Load("map.ptgoraph.pbf")
	fmt.Println(len(ptg.Map.LatLons))
}
