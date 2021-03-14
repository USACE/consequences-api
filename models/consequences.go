package models

import (
	"github.com/USACE/go-consequences/compute"
)

type Compute struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	DepthFilePath string   `json:"depthfilepath"`
}

func ComputeByStructureFromFile(c Compute) (string, error){
	s, err := compute.FromFile(c.DepthFilePath)
	//compute.StreamFromFile(c.DepthFilePath, c.Writer)
	return s, err
}