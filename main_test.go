package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func Test_Consequences(t *testing.T) {

	requestBody := Compute{Name: "myname", DepthFilePath: "/workspaces/consequences-api/data/3782_COG.tif"}
	b, _ := json.Marshal(requestBody)
	response, err := http.Post(
		"http://127.0.0.1:3030/consequences/structure/compute",
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
			fmt.Printf("Error unmarshalling JSON record: %s.  Stopping Compute.\n", err)
		}
		fmt.Println(n)
	}
}
