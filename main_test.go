package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

type Result struct {
	C Consequence `json:"consequence"`
}
type Consequence struct {
	Fdid    string  `json:"fd_id"`
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	DamCat  string  `json:"damage category"`
	OccType string  `json:"occupancy type"`
	SDamage float64 `json:"structure damage"`
	CDamage float64 `json:"content damage"`
}

func Test_Consequences_IsLive(t *testing.T) {
	response, err := http.Get("http://host.docker.internal:8000/consequences")
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", b)
}
func Test_Structure_NSI_Stream(t *testing.T) {
	requestBody := Compute{DepthFilePath: "/vsis3/media/clipped_sample.tif", InventorySource: "NSI"} //default of stream
	b, _ := json.Marshal(requestBody)
	response, err := http.Post(
		"http://host.docker.internal:8000/consequences/structure/compute",
		"application/json; charset=UTF-8",
		bytes.NewReader(b),
	)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	dec := json.NewDecoder(response.Body)

	for {
		var n Result
		if err := dec.Decode(&n); err == io.EOF {
			break
		} else if err != nil {
			panic("Error unmarshalling JSON record: %s.  Stopping Compute.\n")
		}
		fmt.Println(n)
	}
}
func Test_Structure_SHP_Stream(t *testing.T) {
	requestBody := Compute{DepthFilePath: "/vsis3/media/CERA_Adv29_maxwaterelev_4326_90m.tif", InventorySource: "/vsis3/media/ORNLcentroids_LBattributes.shp"}

	b, _ := json.Marshal(requestBody)
	response, err := http.Post(
		"http://host.docker.internal:8000/consequences/structure/compute",
		"application/json; charset=UTF-8",
		bytes.NewReader(b),
	)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	dec := json.NewDecoder(response.Body)

	for {
		var n Result
		if err := dec.Decode(&n); err == io.EOF {
			break
		} else if err != nil {
			panic("Error unmarshalling JSON record: %s.  Stopping Compute.\n")
		}
		fmt.Println(n)
	}
}
func Test_Structure_NSI_GeoJSON(t *testing.T) {
	requestBody := Compute{DepthFilePath: "/vsis3/media/clipped_sample.tif", InventorySource: "NSI", OutputType: "GeoJson"}
	b, _ := json.Marshal(requestBody)
	response, err := http.Post(
		"http://host.docker.internal:8000/consequences/structure/compute",
		"application/json; charset=UTF-8",
		bytes.NewReader(b),
	)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", result)
}
func Test_Structure_SHP_GeoJSON(t *testing.T) {
	requestBody := Compute{DepthFilePath: "/vsis3/media/CERA_Adv29_maxwaterelev_4326_90m.tif", InventorySource: "/vsis3/media/ORNLcentroids_LBattributes.shp", OutputType: "GeoJson"}
	b, _ := json.Marshal(requestBody)
	response, err := http.Post(
		"http://host.docker.internal:8000/consequences/structure/compute",
		"application/json; charset=UTF-8",
		bytes.NewReader(b),
	)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", result)
}
func Test_Structure_NSI_Shp(t *testing.T) {
	requestBody := Compute{DepthFilePath: "/vsis3/media/clipped_sample.tif", InventorySource: "NSI", OutputType: "ESRI SHP"} //default of stream
	b, _ := json.Marshal(requestBody)
	response, err := http.Post(
		"http://host.docker.internal:8000/consequences/structure/compute",
		"application/json; charset=UTF-8",
		bytes.NewReader(b),
	)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", result)
}
func Test_Summary_NSI(t *testing.T) {
	requestBody := Compute{DepthFilePath: "/vsis3/media/clipped_sample.tif", InventorySource: "NSI", OutputType: "Summary"}
	b, _ := json.Marshal(requestBody)
	response, err := http.Post(
		"http://host.docker.internal:8000/consequences/summary/compute",
		"application/json; charset=UTF-8",
		bytes.NewReader(b),
	)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", result)
}
func Test_Summary_SHP(t *testing.T) {
	requestBody := Compute{DepthFilePath: "/vsis3/media/CERA_Adv29_maxwaterelev_4326_90m.tif", InventorySource: "/vsis3/media/ORNLcentroids_LBattributes.shp", OutputType: "Summary"}
	b, _ := json.Marshal(requestBody)
	response, err := http.Post(
		"http://host.docker.internal:8000/consequences/summary/compute",
		"application/json; charset=UTF-8",
		bytes.NewReader(b),
	)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", result)
}
