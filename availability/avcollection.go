package availability

type AvailabilityCollection interface {
	FindAvailabilityById(id string) *Availability
	SaveAvailability(id string, av *Availability)
}

type MemAvailabilityCollection struct {
	avMap map[string]*Availability
}

func NewAvailabilityCollection() *MemAvailabilityCollection {
	return &MemAvailabilityCollection{
		avMap: make(map[string]*Availability),
	}
}

func (avc *MemAvailabilityCollection) FindAvailabilityById(id string) *Availability {
	if av, ok := avc.avMap[id]; ok {
		return av
	}
	return nil
}

func (avc *MemAvailabilityCollection) SaveAvailability(id string, av *Availability) {
	avc.avMap[id] = av
}
