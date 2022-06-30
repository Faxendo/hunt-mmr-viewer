package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/manifoldco/promptui"
)

type Attributes struct {
	XMLName    xml.Name    `xml:"Attributes"`
	Attributes []Attribute `xml:"Attr"`
}

type Attribute struct {
	XMLName xml.Name `xml:"Attr"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value,attr"`
}

type Player struct {
	ID       string
	TeamID   string
	PlayerID string

	Name string
	MMR  string

	HadBounty  string
	KilledByMe string
	KilledMe   string
}

var (
	defaultFileDir   = "user/profiles/default/attributes.xml"
	playerAttrRegexp = "MissionBagPlayer_([0-9]{1,2})_([0-2]{1})_(.+)"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var cfg Config = loadConfig()

	prompt := promptui.Prompt{
		Label:   "Dossier contenant Hunt : Showdown",
		Default: cfg.DefaultFolder,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if result != cfg.DefaultFolder {
		saveConfig(result)
	}

	parse(result)
}

func parse(baseFolder string) {
	xmlFile, err := os.Open(filepath.Join(baseFolder, defaultFileDir))
	check(err)

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var attributes Attributes
	xml.Unmarshal(byteValue, &attributes)

	r, _ := regexp.Compile(playerAttrRegexp)

	var playerMap map[string]Player = make(map[string]Player)
	var players []Player

	var teamCount int = 0

	for _, attr := range attributes.Attributes {
		if attr.Name == "MissionBagNumTeams" {
			teamCount, _ = strconv.Atoi(attr.Value)
		}
	}

	if teamCount == 0 {
		panic("Error while parsing")
	}

	for _, attr := range attributes.Attributes {
		if r.MatchString(attr.Name) {

			submatch := r.FindStringSubmatch(attr.Name)
			teamId := submatch[1]
			playerId := submatch[2]
			attrName := submatch[3]

			teamIdInt, _ := strconv.Atoi(teamId)

			if (teamIdInt + 1) > teamCount {
				continue
			}

			completeId := teamId + "_" + playerId

			var player Player

			if _, ok := playerMap[completeId]; !ok {
				player.ID = completeId
				player.TeamID = strconv.Itoa(teamIdInt + 1)
				player.PlayerID = playerId

				playerMap[completeId] = player
			}

			player = playerMap[completeId]

			switch attrName {
			case "blood_line_name":
				player.Name = attr.Value
			case "mmr":
				player.MMR = attr.Value
			case "hadbounty":
				if attr.Value == "true" {
					player.HadBounty = "X"
				} else {
					player.HadBounty = ""
				}
			case "killedbyme":
				if attr.Value == "1" {
					player.KilledByMe = "X"
				} else {
					player.KilledByMe = ""
				}
			case "killedme":
				if attr.Value == "1" {
					player.KilledMe = "X"
				} else {
					player.KilledMe = ""
				}
			}

			playerMap[completeId] = player
		}
	}

	for _, p := range playerMap {
		// fmt.Printf("%s | [MMR : %s]\n", p.Name, p.MMR)
		players = append(players, p)
	}

	defer xmlFile.Close()

	display(players)
}

func display(players []Player) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(table.Row{
		"Team",
		"Name",
		"MMR",
		"Had Bounty",
		"Killed Me",
		"Killed by Me",
	})
	for _, p := range players {
		t.AppendRow(table.Row{
			p.TeamID,
			p.Name,
			p.MMR,
			p.HadBounty,
			p.KilledByMe,
			p.KilledMe,
		})
	}

	t.SetStyle(table.StyleColoredBlackOnRedWhite)
	t.SortBy([]table.SortBy{
		{Name: "Team", Mode: table.Asc},
	})

	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:  "Had Bounty",
			Align: text.AlignCenter,
		},
		{
			Name:  "Killed Me",
			Align: text.AlignCenter,
		},
		{
			Name:  "Killed by Me",
			Align: text.AlignCenter,
		},
	})

	t.Render()

	fmt.Println("Press any key to exit")
	fmt.Scanln()
}
