package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"time"
)

type ElasticData struct {
	TeamID   string `json:"TeamID"`
	PlayerID string `json:"PlayerID"`
	ProfilID string `json:"ProfilID"`

	MMR string `json:"MMR"`

	HadBounty  string `json:"HadBounty"`
	KilledByMe string `json:"KilledByMe"`
	KilledMe   string `json:"KilledMe"`

	MatchHash int `json:"MatchHash"`

	Timestamp int `json:"timestamp"`
}

func parsePlayerAndSendToElastic(players []Player, hashMatch int) {
	var elasticData ElasticData

	for _, p := range players {

		elasticData.TeamID = p.TeamID
		elasticData.PlayerID = p.PlayerID
		elasticData.ProfilID = p.ProfilID
		elasticData.MMR = p.MMR
		elasticData.HadBounty = p.HadBounty
		elasticData.KilledByMe = p.KilledByMe
		elasticData.KilledMe = p.KilledMe
		elasticData.MatchHash = hashMatch
		elasticData.Timestamp = int(time.Now().Unix())
	}

	SendToElastic(elasticData)
}

func SendToElastic(dataPlayer ElasticData) {

	cfg := elasticsearch.Config{
		Addresses: []string{
			"https://kibana.sevenn.fr/e",
		},
		Username: "hunt",
		Password: "huntmmr",
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	data, err := json.Marshal(dataPlayer)
	if err != nil {
		panic(err)
	}

	req := esapi.IndexRequest{
		Index:   "hunt_mmr",
		Body:    bytes.NewReader(data),
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error indexing document", res.Status())
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			log.Print("Stats successfully sent to tracking server.")
		}
	}
}
