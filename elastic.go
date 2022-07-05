package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
)

type ElasticData struct {
	TeamID   string `json:"TeamID"`
	PlayerID string `json:"PlayerID"`
	ProfilID string `json:"ProfilID"`
	Name     string `json:"Name"`

	MMR string `json:"MMR"`

	HadBounty  bool `json:"HadBounty"`
	KilledByMe bool `json:"KilledByMe"`
	KilledMe   bool `json:"KilledMe"`

	MatchHash int    `json:"MatchHash"`
	UUID      string `json:"uuid"`
}

func parsePlayerAndSendToElastic(players []Player, hashMatch int, uuid string) {
	for _, p := range players {
		var elasticData ElasticData

		elasticData.TeamID = p.TeamID
		elasticData.PlayerID = p.PlayerID
		elasticData.ProfilID = p.ProfilID
		elasticData.Name = p.Name
		elasticData.MMR = p.MMR
		if p.HadBounty == "X" {
			elasticData.HadBounty = true
		}
		if p.KilledByMe == "X" {
			elasticData.KilledByMe = true
		}
		if p.KilledMe == "X" {
			elasticData.KilledMe = true
		}
		elasticData.MatchHash = hashMatch
		elasticData.UUID = uuid

		SendToElastic(elasticData)
	}
}

func SendToElastic(dataPlayer ElasticData) {

	var cfgFile Config = loadConfig()

	cfg := elasticsearch.Config{
		Addresses: []string{
			cfgFile.ElasticUrl,
		},
		Username: cfgFile.ElacticUser,
		Password: cfgFile.ElasticPassword,
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
		Index:    "hunt_mmr",
		Body:     bytes.NewReader(data),
		Refresh:  "true",
		Pipeline: "timestamp",
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
		}
	}
}
