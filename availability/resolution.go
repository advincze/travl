package availability

type TimeResolution int32

const (
	undefined TimeResolution = 0
	sec       TimeResolution = 1
	Minute    TimeResolution = sec * 60
	Minute5   TimeResolution = Minute * 5
	Minute15  TimeResolution = Minute * 15
	Hour      TimeResolution = Minute * 60
	Day       TimeResolution = Hour * 24
)

var resolutionStrings = map[string]TimeResolution{
	"s":     sec,
	"sec":   sec,
	"m":     Minute,
	"min":   Minute,
	"5m":    Minute5,
	"5min":  Minute5,
	"15m":   Minute15,
	"15min": Minute15,
	"h":     Hour,
	"hour":  Hour,
	"d":     Day,
	"day":   Day,
}

func (tr TimeResolution) String() string {
	switch tr {
	case sec:
		return "sec"
	case Minute:
		return "min"
	case Minute5:
		return "5 min"
	case Minute15:
		return "15 min"
	case Hour:
		return "hour"
	case Day:
		return "day"
	}
	return "undefined"
}

func ParseTimeResolution(s string) TimeResolution {
	if resolution, ok := resolutionStrings[s]; ok {
		return resolution
	}
	return undefined
}
