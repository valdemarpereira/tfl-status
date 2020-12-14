package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

//https://stackoverflow.com/questions/4842424/list-of-ansi-color-escape-sequences
var (
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	White   = Color("\033[1;37m%s\033[0m")

	BoldGreen = Color("\033[1;38;5;46m%s\033[0m")
	BoldRed = Color("\033[1;38;5;196m%s\033[0m")
)

var (
	JubileeColor      = Color("\033[38;5;231;48;5;102m%s\033[0m")
	BakerlooColor     = Color("\033[38;5;231;48;5;94m%s\033[0m")
	CentralColor      = Color("\033[38;5;231;48;5;160m%s\033[0m")
	CircleColor       = Color("\033[38;5;16;48;5;220m%s\033[0m")
	DistrictColor     = Color("\033[38;5;231;48;5;22m%s\033[0m")
	HammersmithColor  = Color("\033[38;5;16;48;5;175m%s\033[0m")
	MetropolitanColor = Color("\033[38;5;231;48;5;89m%s\033[0m")
	NorthenColor      = Color("\033[38;5;231;48;5;16m%s\033[0m")
	PiccadillyColor   = Color("\033[38;5;231;48;5;19m%s\033[0m")
	VictoriaColor     = Color("\033[38;5;16;48;5;38m%s\033[0m")
	WaterlooColor     = Color("\033[38;5;16;48;5;115m%s\033[0m")


)

var stationsNames = []string{
	"Jubilee",
	"Bakerloo",
	"Central",
	"Circle",
	"District",
	"Hammersmith",
	"Metropolitan",
	"Northen",
	"Piccadilly",
	"Victoria",
	"Waterloo",
}

type Tube struct {
	LineName  string
	LineId    string
	LineColor func(...interface{}) string
}

var LondonTube = []Tube{
	{
		LineName:  "Jubilee",
		LineId:    "jubilee",
		LineColor: JubileeColor,
	},
	{
		LineName:  "Bakerloo",
		LineId:    "bakerloo",
		LineColor: BakerlooColor,
	},
	{
		LineName:  "Central",
		LineId:    "central",
		LineColor: CentralColor,
	},
	{
		LineName:  "District",
		LineId:    "district",
		LineColor: DistrictColor,
	},
	{
		LineName:  "Hammersmith",
		LineId:    "hammersmith-city",
		LineColor: HammersmithColor,
	},
	{
		LineName:  "Circle",
		LineId:    "circle",
		LineColor: CircleColor,
	},
	{
		LineName:  "Metropolitan",
		LineId:    "metropolitan",
		LineColor: MetropolitanColor,
	},
	{
		LineName:  "Northen",
		LineId:    "northern",
		LineColor: NorthenColor,
	},
	{
		LineName:  "Piccadilly",
		LineId:    "piccadilly",
		LineColor: PiccadillyColor,
	},
	{
		LineName:  "Victoria",
		LineId:    "victoria",
		LineColor: VictoriaColor,
	},
	{
		LineName:  "Waterloo",
		LineId:    "waterloo-city",
		LineColor: WaterlooColor,
	},
}

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

func main() {
	json, err := getTFLStatusJson()
	if err != nil {
		fmt.Println("ERROR")
	}

	for _, tube := range LondonTube {
		tubeLineNameJustify := justifyText(tube.LineName, stationsNames, 2)
		statusSeverity, statusReason := getLineStatus(tube, json)

		fmt.Print(tube.LineColor(tubeLineNameJustify))
		fmt.Println(statusSeverity)

		if len(statusReason) > 0 {
			text := fmt.Sprintf("%*v", len(tubeLineNameJustify)+2+len(statusReason), statusReason)
			fmt.Println(text)
		}
	}
}

func getLineStatus(tube Tube, tflStatus json.RawMessage) (string, string) {

	expressionStatusSeverityId := fmt.Sprintf("#(id==%s).lineStatuses.0.statusSeverity", tube.LineId)

	expressionStatusSeverity := fmt.Sprintf("#(id==%s).lineStatuses.0.statusSeverityDescription", tube.LineId)

	expressionStatusReason := fmt.Sprintf("#(id==%s).lineStatuses.0.reason", tube.LineId)

	statusId := gjson.GetBytes(tflStatus, expressionStatusSeverityId)
	status := gjson.GetBytes(tflStatus, expressionStatusSeverity)
	statusReason := gjson.GetBytes(tflStatus, expressionStatusReason)

	if statusId.Int() == 10 { //Good Service
		return fmt.Sprintf(BoldGreen("  " + status.String())), ""
	} else {
		return fmt.Sprintf(BoldRed("  " + status.String())), statusReason.String()

	}
}

func justifyText(text string, texts []string, addedPad int) string {

	maxLength := 0

	for _, str := range texts {
		if len(str) > maxLength {
			maxLength = len(str)
		}
	}

	pad := float64(maxLength-len(text)) / float64(2)

	padLeft := math.Floor(pad)

	leftPadded := fmt.Sprintf("%*v", len(text)+int(padLeft)+addedPad, text)
	leftPadded = fmt.Sprintf("%-*v", maxLength+(addedPad*2), leftPadded)

	return leftPadded

}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getTFLStatusJson() (json.RawMessage, error) {
	url := "https://api.tfl.gov.uk/line/mode/tube/status?detail=true"

	r, err := myClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		return nil, err
	}

	return body, nil
}
