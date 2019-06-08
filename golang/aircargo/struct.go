package main

import . "github.com/davidkhala/goutils"

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

type CargoStatus struct {
	Status int
	CurrentOwner string

}
type AirCargo struct {
	Package
	Status     CargoStatus
	ShipperID  string
	CreateTime TimeLong
	ImageUrls  map[string][]string
}

type createAirCargo struct {
	AirCargo
	EntityID   string
	ImageUrls  []string
	AirCargoID string
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
