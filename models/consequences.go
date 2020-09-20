package models

import (
	"encoding/json"

	"github.com/USACE/go-consequences/consequences"
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

// RunConsequences Runs the Consequences
func RunConsequences(d ConsequencesInputCollection) ([]ConsequencesInputAndResult, error) {
	ss := make([]ConsequencesInputAndResult, len(d.Items))
	for idx, item := range d.Items {
		result := consequences.BaseStructure().ComputeConsequences(item.Depth)
		ss[idx] = ConsequencesInputAndResult{
			ConsequencesInput: item,
			ConsequencesResult: ConsequencesResult{
				Name:   "Test",
				Result: result.String(),
			},
		}
	}
	return ss, nil
}

// RunConsequences Runs the Consequences
func GetInventory(d ConsequencesBoundingBox) (ConsequencesStructureInventory, error) {
	structures := nsi.GetByBbox(d.BoundingBox) //i'd like to save this in memory while I wait on IFIM to request damages...
	result := make([]ConsequencesStructure, len(structures))
	for idx, structure := range structures {
		result[idx] = ConsequencesStructure{
			X:     structure.X,
			Y:     structure.Y,
			FD_ID: structure.Name,
		}
	}
	inventory := ConsequencesStructureInventory{Inventory: result}
	return inventory, nil
}