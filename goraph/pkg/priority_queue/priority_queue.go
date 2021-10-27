package priority_queue

type costque struct {
	key int64
	num float64
}

type MinSet struct {
	que      [][]costque
	last_num int
	n        int64
	C        int
}

func NewMinSet() *MinSet {
	s := new(MinSet)
	s.last_num = 0
	s.n = 0
	s.C = 2048
	s.que = make([][]costque, s.C)
	return s
}

func (s *MinSet) AddVal(key int64, val float64) {
	var t costque
	t.key = key
	t.num = val
	s.n++
	s.que[int(val)%s.C] = append(s.que[int(val)%s.C], t)
}

func (s *MinSet) GetMin() int64 {
	for s.n > 0 {
		i := s.last_num % s.C
		if len(s.que[i]) == 0 {
			s.last_num++
			continue
		}
		minI := 0
		for j := 1; j < len(s.que[i]); j++ {
			if s.que[i][j].num < s.que[i][minI].num {
				minI = j
			}
		}
		if int(s.que[i][minI].num) != int(s.last_num) {
			s.last_num++
			continue
		}

		ans := s.que[i][minI].key
		s.que[i][minI] = s.que[i][0]
		s.que[i] = s.que[i][1:]

		s.n--
		return ans
	}
	return -1
}

func (s *MinSet) Len() int64 {
	return s.n
}
