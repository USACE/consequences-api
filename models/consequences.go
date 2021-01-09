package models

import (
	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/crops"
	"github.com/USACE/go-consequences/hazards"
)

// RunConsequencesByBoundingBox Runs the Consequences by bounding box
func RunConsequencesByBoundingBox(cbb compute.BboxCompute) (consequences.ConsequenceDamageResult, error) {
	var r = compute.NSIStructureSimulation{}
	cbb.HazardArgs = hazards.DepthEvent{Depth: 8.73}
	args := compute.RequestArgs{Args: cbb}
	r.Compute(args) //how do i get this to go on a separate thread?
	output := r.Result
	//r.Status == "Complete?"
	return output, nil
}
func RunConsequencesByFips(fips compute.FipsCodeCompute, depth float64) (consequences.ConsequenceDamageResult, error) {
	var r = compute.NSIStructureSimulation{}
	fips.HazardArgs = hazards.DepthEvent{Depth: depth}
	args := compute.RequestArgs{Args: fips}
	r.Compute(args) //how do i get this to go on a separate thread?
	output := r.Result
	//r.Status == "Complete?"
	return output, nil
}
func RunAgConsequencesByXY(x string, y, string, event hazards.ArrivalandDurationEvent) (consequences.ConsequenceDamageResult, error) {
	//var r = compute.NSIStructureSimulation{}
	//fips.HazardArgs = hazards.DepthEvent{Depth: depth}
	//args := compute.RequestArgs{Args: fips}
	//r.Compute(args) //how do i get this to go on a separate thread?
	//output := r.Result
	//r.Status == "Complete?"
	cropFromNass := crops.GetCDLValue("2018", x, y)
	path := "C:\\Temp\\agtesting\\" + cropFromNass.GetCropName() + ".crop" //this wont work.
	c := crops.ReadFromXML(path)
	//compute
	cd := c.ComputeConsequences(event)
	return cd, nil
}
