package main

import (
	"fmt"
	"github.com/takoyaki-3/mmGraphpt/loader"
)

func main(){
	ptg := loader.Load("map.ptgoraph.pbf")
	fmt.Println(len(ptg.Map.LatLons))
}
