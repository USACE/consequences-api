package models

import (
	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/consequences"
)

// RunConsequencesByBoundingBox Runs the Consequences by bounding box
func RunConsequencesByBoundingBox(cbb compute.Bbox) (consequences.ConsequenceDamageResult, error) {
	var r = compute.NSIStructureSimulation{}
	args := compute.ComputeArgs{Args: cbb}
	r.Compute(args) //how do i get this to go on a separate thread?
	output := r.Result
	//r.Status == "Complete?"
	return output, nil
}
func RunConsequencesByFips(fips compute.FipsCode) (consequences.ConsequenceDamageResult, error) {
	var r = compute.NSIStructureSimulation{}
	args := compute.ComputeArgs{Args: fips}
	r.Compute(args) //how do i get this to go on a separate thread?
	output := r.Result
	//r.Status == "Complete?"
	return output, nil
}
