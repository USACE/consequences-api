package models

import (
	"errors"
	"io"

	"github.com/USACE/go-consequences/compute"
	"github.com/google/uuid"
)

type Compute struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	DepthFilePath string    `json:"depthfilepath"`
}

func ComputeByStructureFromFile(c Compute, w io.Writer) (string, error) {
	if c.valid() {
		compute.StreamFromFile(c.DepthFilePath, w)
		//compute.FromFile(c.DepthFilePath)
		return "", nil
	}
	return "", errors.New("aint no way you are getting that to work...")
}
func (c Compute) valid() bool {
	return true //@TODO implement me!
}
