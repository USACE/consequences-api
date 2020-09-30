package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/nsi"
)

//ConsequencesBoundingBox is a list of x,y representing a square
type ConsequencesBoundingBox struct {
	BoundingBox string `json:"bbox"`
}

//ConsequencesStructure is a structure FDID and x,y location
type ConsequencesStructure struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	FD_ID string  `json:"fd_id"`
}

//ConsequencesStructureInventory is a list of ConsequencesStructure
type ConsequencesStructureInventory struct {
	Inventory []ConsequencesStructure `json:"structures"`
}

// ConsequencesInput is input
type ConsequencesInput struct {
	Structure ConsequencesStructure `json:"structure"`
	Depth     float64               `json:"depth"`
}

// ConsequencesInputCollection is many consequences inputs
type ConsequencesInputCollection struct {
	Items []ConsequencesInput `json:"items"`
}

// UnmarshalJSON implements UnmarshalJSON interface
func (c *ConsequencesInputCollection) UnmarshalJSON(b []byte) error {

	switch JSONType(b) {
	case "ARRAY":
		if err := json.Unmarshal(b, &c.Items); err != nil {
			return err
		}
	case "OBJECT":
		var n ConsequencesInput
		if err := json.Unmarshal(b, &n); err != nil {
			return err
		}
		c.Items = []ConsequencesInput{n}
	default:
		c.Items = make([]ConsequencesInput, 0)
	}
	return nil
}

// ConsequencesResult is a Something
type ConsequencesResult struct {
	Name   string `json:"name"`
	Result string `json:"result"`
}

// ConsequencesInputAndResult includes all fields from both
type ConsequencesInputAndResult struct {
	ConsequencesInput
	ConsequencesResult
}

// RunConsequencesByBoundingBox Runs the Consequences by bounding box
func RunConsequencesByBoundingBox(cbb ConsequencesBoundingBox) ([]ConsequencesInputAndResult, error) {
	structures := nsi.GetByBbox(cbb.BoundingBox) //i'd like to save this in memory while I wait on IFIM to request damages...
	ifimRequest := make([]ConsequencesStructure, len(structures))
	for idx, structure := range structures {
		ifimRequest[idx] = ConsequencesStructure{
			X:     structure.X,
			Y:     structure.Y,
			FD_ID: structure.Name,
		}
	}

	//query IFIM for depths
	ifimResponse := make([]ConsequencesInput, len(ifimRequest))
	for idx, location := range ifimRequest {
		ifimResponse[idx] = ConsequencesInput{Structure: location, Depth: 5.0}
	}

	output := make([]ConsequencesInputAndResult, len(ifimResponse))
	for idx, item := range ifimResponse {
		d := hazards.DepthEvent{Depth: item.Depth}
		result := consequences.BaseStructure().ComputeConsequences(d)
		output[idx] = ConsequencesInputAndResult{
			ConsequencesInput: item,
			ConsequencesResult: ConsequencesResult{
				Name:   "Test",
				Result: result.String(),
			},
		}
	}
	return output, nil
}
func RunConsequencesByFips(fips string) (string, error) {
	startnsi := time.Now()
	structures := nsi.GetByFips(fips) //i'd like to save this in memory while I wait on IFIM to request damages...
	ifimRequest := make([]ConsequencesStructure, len(structures))
	for idx, structure := range structures {
		ifimRequest[idx] = ConsequencesStructure{
			X:     structure.X,
			Y:     structure.Y,
			FD_ID: structure.Name,
		}
	}
	elapsedNsi := time.Since(startnsi)
	//query IFIM for depths
	ifimResponse := make([]ConsequencesInput, len(ifimRequest))
	for idx, location := range ifimRequest {
		ifimResponse[idx] = ConsequencesInput{Structure: location, Depth: 5.46}
	}

	startcompute := time.Now()
	var count = 0
	for idx, item := range ifimResponse {
		d := hazards.DepthEvent{Depth: item.Depth}
		consequences.BaseStructure().ComputeConsequences(d)
		count = idx
	}
	count += 1
	elapsed := time.Since(startcompute)
	output := fmt.Sprintf("NSI Fetching took %s Compute took %s for %d structures", elapsedNsi, elapsed, count)
	return output, nil
}
