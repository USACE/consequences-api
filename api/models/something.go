package models

import (
	"consequences-api/api/consequences"
	"encoding/json"
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
	for idx, dd := range d.Items {
		ret := consequences.BaseStructure().ComputeConsequences(dd.Depth)
		s := ConsequencesInputAndResult{}
		s.ConsequencesInput = dd
		s.ConsequencesResult = ConsequencesResult{Name: "Test", Result: ret.String()}
		ss[idx] = s
	}
	return ss, nil
}
