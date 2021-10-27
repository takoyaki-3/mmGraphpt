package goraph

type Edge struct{
	To int64
	Cost float64
	LatLons []int64
}

type LatLon struct{
	Lat float64
	Lon float64
}

type Graph struct{
	Edges [][]Edge
	LatLons []LatLon
}

func (s *Graph)AddEdge(e Edge,from int64){
	for int64(len(s.Edges)) <= from {
		s.Edges = append(s.Edges,[]Edge{})
	}
	s.Edges[from] = append(s.Edges[from],e)
}
func (s *Graph)SetLatLon(n LatLon,id int64){
	for int64(len(s.LatLons)) <= id {
		s.LatLons = append(s.LatLons,LatLon{})
	}
	s.LatLons[id] = n
}
func (s *Graph)AddLatLon(n LatLon)int64{
	id := int64(len(s.LatLons))
	s.LatLons = append(s.LatLons,n)
	for len(s.Edges) < len(s.LatLons){
		s.Edges = append(s.Edges, []Edge{})
	}
	return id
}