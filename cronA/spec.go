package cronA

import "time"

type bounds struct {
	min, max uint64
	names    map[string]uint64
}

// https://web.archive.org/web/20111025080042/http://www.quartz-scheduler.org/documentation/quartz-1.x/tutorials/crontrigger#:~:text=FRI%202002%2D2010-,Special%20characters,-*%20(%22all

var (
	seconds = bounds{0, 59, nil}
	minutes = bounds{0, 59, nil}
	hours   = bounds{0, 23, nil}
	dom     = bounds{1, 31, nil} // day of month
	months  = bounds{1, 12, map[string]uint64{
		"jan": 1, "feb": 2, "mar": 3, "apr": 4,
		"may": 5, "jun": 6, "jul": 7, "aug": 8,
		"sep": 9, "oct": 10, "nov": 11, "dec": 12,
	}}
	dow = bounds{0, 6, map[string]uint64{ // day of week
		"sun": 0, "mon": 1, "tue": 2, "wed": 3,
		"thu": 4, "fri": 5, "sat": 6,
	}}
)

func (r bounds) all() uint64 { return getBits(r.min, r.max, 1) }

func getBits(min, max, step uint64) uint64 {
	var bits uint64
	for i := min; i <= max; i += step {
		bits |= 1 << i
	}
	return bits
}

type Schedule struct {
	Second, Minute, Hour, Dom, Month, Dow uint64
	Location                              *time.Location
}

func (s *Schedule) Next() {}

func (s *Schedule) Prev() {}
