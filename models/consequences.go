package models

import (
	"encoding/json"

	"github.com/USACE/go-consequences/consequences"
)

// ConsequencesInput is input
type ConsequencesInput struct {
	Depth float64 `json:"depth"`
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
