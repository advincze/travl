package availability

import (
	"loveoneanother.at/tiedot/db"
)

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

type TiedotAvailabilityCollection struct {
	collection *db.Col
}

func NewTiedotAvailabilityCollection() *TiedotAvailabilityCollection {
	dir := "MyDatabase"
	database, err := db.OpenDB(dir)
	if err != nil {
		panic(err)
	}
	database.Create("av")

	return &TiedotAvailabilityCollection{
		collection: database.Use("av"),
	}
}

func (avc *TiedotAvailabilityCollection) FindAvailabilityById(id string) *Availability {
	var av *Availability
	avc.collection.Read(uint64(0), &av)
	return nil
}

func (avc *TiedotAvailabilityCollection) SaveAvailability(id string, av *Availability) {
	avc.collection.Insert(nil)
}
