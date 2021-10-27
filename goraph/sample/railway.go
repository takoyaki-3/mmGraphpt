package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/takoyaki-3/mmGraphpt/goraph"
	"github.com/takoyaki-3/mmGraphpt/goraph/geometry"
	"github.com/takoyaki-3/mmGraphpt/goraph/loader/geojson"
	"github.com/takoyaki-3/mmGraphpt/goraph/search"
)

func main() {
	fmt.Println("start")

	g := geojson.Load("N02-19_RailroadSection.geojson")

	stationcode2latlon := map[string]goraph.LatLon{}
	func() {
		// stationcodeから緯度経度へ変換
		file, err := os.Open("stationcode2latlon.csv")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		reader := csv.NewReader(file)
		titles := map[string]int{}
		counter := -1

		for {
			counter++
			line, err := reader.Read()
			if err != nil {
				break
			}
			if counter == 0 {
				for k, v := range line {
					titles[v] = k
				}
				continue
			}
			lat, _ := strconv.ParseFloat(line[titles["lat"]], 64)
			lon, _ := strconv.ParseFloat(line[titles["lon"]], 64)
			station_code := strings.Replace(line[titles["MLIT.stationcode"]], ".0", "", 2)
			stationcode2latlon[station_code] = goraph.LatLon{lat, lon}
		}
	}()

	// station point
	station_points := map[string][]int{}
	func() {
		file, err := os.Open("station_point.csv")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		reader := csv.NewReader(file)
		var line []string
		titles := map[string]int{}
		counter := -1

		for {
			counter++
			line, err = reader.Read()
			if err != nil {
				break
			}
			if counter == 0 {
				for k, v := range line {
					titles[v] = k
				}
				continue
			}
			if _, ok := station_points[line[titles["title"]]]; !ok {
				station_points[line[titles["title"]]] = []int{}
			}
			node_id, _ := strconv.Atoi(line[titles["node_id"]])
			station_points[line[titles["title"]]] = append(station_points[line[titles["title"]]], node_id)
		}
	}()

	// 出力用ファイル
	outfile, _ := os.OpenFile("railway_geometry.csv", os.O_WRONLY|os.O_CREATE, 0600)
	defer outfile.Close()

	err := outfile.Truncate(0)
	if err != nil {
		panic(err)
	}
	writer := csv.NewWriter(outfile)
	writer.Write([]string{"from", "to", "geometry"})

	func() {
		file, err := os.Open("table.csv")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		reader := csv.NewReader(file)
		titles := map[string]int{}
		counter := -1

		for {
			counter++
			line, err := reader.Read()
			if err != nil {
				fmt.Println("finish!")
				break
			}
			if counter == 0 {
				for k, v := range line {
					titles[v] = k
				}
				continue
			}

			// search
			min_cost := math.MaxFloat64
			min_json := ""

			fmt.Println(counter, line[titles["title"]])

			from_code := strings.Replace(line[titles["from_station.MLIT.stationcode"]], ".0", "", -1)
			to_code := strings.Replace(line[titles["to_station.MLIT.stationcode"]], ".0", "", -1)

			start := stationcode2latlon[from_code]
			end := stationcode2latlon[to_code]

			if from_code == "" || to_code == "" {
				writer.Write([]string{from_code, to_code, min_json})
				continue
			}

			for _, from_node := range station_points[from_code] {
				for _, to_node := range station_points[to_code] {
					ans := search.Search(g, search.Query{int64(from_node), int64(to_node)})
					ans.Nodes = append([]int64{int64(len(g.LatLons))}, ans.Nodes...)
					g.LatLons = append(g.LatLons, start)
					ans.Nodes = append(ans.Nodes, int64(len(g.LatLons)))
					g.LatLons = append(g.LatLons, end)
					if min_cost > ans.Cost {
						min_cost = ans.Cost
						min_json = geometry.MakeLineString(g, ans.Nodes)
					}
				}
			}
			writer.Write([]string{from_code, to_code, min_json})
		}
	}()
	writer.Flush()
}
