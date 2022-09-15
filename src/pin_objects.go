package main

import (
	"time"
)

// Response used for listing pin objects matching request
type PinResults struct {
	// The total number of pin objects that exist for passed query filters
	Count int32 `json:"count"`
	// An array of PinStatus results
	Results []PinStatus `json:"results"`
}

// Pin object with status
type PinStatus struct {
	// Globally unique identifier of the pin request; can be used to check the status of ongoing pinning, or pin removal
	Requestid string `json:"requestid"`

	Status *Status `json:"status"`
	// Immutable timestamp indicating when a pin request entered a pinning service; can be used for filtering results and pagination
	Created time.Time `json:"created"`

	Pin *Pin `json:"pin"`

	Delegates *[]string `json:"delegates"`

	Info *map[string]string `json:"info,omitempty"`
}

// Pin object
type Pin struct {
	// Content Identifier (CID) to be pinned recursively
	Cid string `json:"cid"`

	Path string `json:"path"`
	// Optional name for pinned data; can be used for lookups later
	Name string `json:"name,omitempty"`

	Origins *[]string `json:"origins,omitempty"`

	Meta *map[string]string `json:"meta,omitempty"`
}

// Status : Status a pin object can have at a pinning service
type Status string

// List of Status
const (
	QUEUED  Status = "queued"
	PINNING Status = "pinning"
	PINNED  Status = "pinned"
	FAILED  Status = "failed"
)
