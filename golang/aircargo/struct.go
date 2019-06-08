package main

const (
	Accepted = iota
	Transferred
	Created
	Rejected
)

type Entity struct {
	Operators []string
	Role      string
}

type Package struct {
	NoOfPieces  int
	GrossWeight string
	Content     string
}

type AirCargo struct {
	Package
	Status    int
	ShipperID string
}

type createAirCargo struct {
	AirCargo
	EntityID  string
	ImageUrls []string
}

type transferAirCargo struct {
	AirCargoID string
	EntityID   string
	ImageUrls  []string
}
type handleAirCargo struct {
	AirCargoID string
	EntityID   string
	NewStatus  int
	Comment    string
	ImageUrls  []string
}
