package csv

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	"github.com/takoyaki-3/go_replace"
	"github.com/takoyaki-3/mmGraphpt/goraph"
)

// Load CSV
func LoadEdge(filename string) goraph.Graph {

	g := goraph.Graph{}

	replace_nodeid := *go_replace.NewReplace()

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	reader := csv.NewReader(file)
	defer file.Close()

	counter := -1
	titles := map[string]int{}

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
		from := replace_nodeid.AddReplace(line[titles["from"]])
		to := replace_nodeid.AddReplace(line[titles["to"]])
		cost := 1.0
		if v, ok := titles["cost"]; ok {
			cost, _ = strconv.ParseFloat(line[v], 64)
		}

		e := goraph.Edge{}
		e.To = to
		e.Cost = cost

		g.AddEdge(e, from)
	}

	return g
}
func LoadLatLon(filename string) goraph.Graph {

	g := goraph.Graph{}

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	reader := csv.NewReader(file)
	defer file.Close()

	counter := -1
	titles := map[string]int{}

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
		id, _ := strconv.ParseInt(line[titles["latlon_id"]], 10, 64)
		lat, _ := strconv.ParseFloat(line[titles["lat"]], 64)
		lon, _ := strconv.ParseFloat(line[titles["lon"]], 64)

		g.SetLatLon(goraph.LatLon{lat, lon}, id)
	}

	return g
}

func WriteEdge(filename string, g goraph.Graph) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	wr := csv.NewWriter(f)
	wr.Write([]string{"from", "to", "cost"})

	for k, v := range g.Edges {
		for _, v := range v {
			cost := strconv.FormatFloat(v.Cost, 'f', -1, 64)
			line := []string{strconv.FormatInt(int64(k), 10), strconv.FormatInt(v.To, 10), cost}
			wr.Write(line)
		}
	}

	wr.Flush()
	if err := wr.Error(); err != nil {
		log.Fatal(err)
	}
}

func WriteLatLon(filename string, g goraph.Graph) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	wr := csv.NewWriter(f)
	wr.Write([]string{"latlon_id", "lat", "lon"})

	for k, v := range g.LatLons {
		lat := strconv.FormatFloat(v.Lat, 'f', -1, 64)
		lon := strconv.FormatFloat(v.Lon, 'f', -1, 64)
		line := []string{strconv.FormatInt(int64(k), 10), lat, lon}
		wr.Write(line)
	}

	wr.Flush()
	if err := wr.Error(); err != nil {
		log.Fatal(err)
	}
}
