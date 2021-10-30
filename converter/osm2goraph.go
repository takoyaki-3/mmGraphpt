package converter

import (
	"fmt"

	"github.com/takoyaki-3/goraph/loader"
	"github.com/takoyaki-3/goraph/loader/osm"
)

func Osm2goraph(inputFileName string, outputFileName string) {
	// OSMの読み込み
	g := osm.Load(inputFileName)
	fmt.Println(len(g.LatLons))
	// プロトコルバッファの書き出し
	loader.Write(outputFileName, g)
}
