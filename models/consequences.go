package models

import (
	"io"

	"github.com/USACE/go-consequences/compute"
)

type Compute struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	DepthFilePath string   `json:"depthfilepath"`
}

func ComputeByStructureFromFile(c Compute, w io.Writer) (string, error){
	s, err := compute.StreamFromFile(c.DepthFilePath,w)
	//compute.FromFile(c.DepthFilePath)
	return s, err
}