package analytics

import (
	"fmt"
	"os"
	"sort"

	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"

	"github.com/idlephysicist/cave-logger/internal/db"
	//"github.com/idlephysicist/cave-logger/internal/model"
)

const TempDir = `/tmp/cl`

var (
	colorWhite          = drawing.Color{R: 241, G: 241, B: 241, A: 255}
	colorMariner        = drawing.Color{R: 60, G: 100, B: 148, A: 255}
)

type Analysis struct {
	db  *db.Database
}


func New(d *db.Database) *Analysis {
	return &Analysis{db: d}
}

func (a *Analysis) Run() error {
	defer cleanup()

	if err := a.peopleAndClubFreq(); err != nil {
		return err
	}

	return nil
}

//
// INTERNAL
//

func (a *Analysis) peopleAndClubFreq() error {
	cavers, err := a.db.GetAllPeople()
	if err != nil {
		return err
	}

	people := make(map[string]float64, 0)
	club := make(map[string]float64, 0)
	for _, caver := range cavers {
		people[caver.Name] += 1.0
		club[caver.Club] += 1.0
	}

	// Club
	barChart := chart.BarChart{
		Title: "Caving Club Frequency",
		TitleStyle: chart.StyleShow(),
		Background: chart.Style{
			Padding: chart.Box{
				Top: 100,
			},
		},
		Width:  1000,
		Height: 700,
		XAxis:  chart.Style{
			Show: true,
			TextRotationDegrees: 25,
			TextHorizontalAlign: chart.TextHorizontalAlignLeft,
			TextWrap: chart.TextWrapNone,
		},
		YAxis: chart.YAxis{
			Name: "Frequency",
			NameStyle: chart.StyleShow(),
			Style: chart.StyleShow(),/*chart.Style{
				Show: true,
				TextHorizontalAlign: chart.TextHorizontalAlignLeft,
			},*/
			ValueFormatter: func(v interface{}) string {
					return fmt.Sprintf("%d", int(v.(float64)))
			},
		},
		BarSpacing: 20,
		Bars: corralBars(club),
	}

	if err = saveBarChart(barChart, `cave.png`); err != nil {
		return err
	}
	return nil
}


func corralBars(data map[string]float64) []chart.Value {
	var keys []string
	for k := range data {
		if k == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	points := make([]chart.Value, len(data))
	i := 0
	for _, k := range keys {
		points[i] = chart.Value{
			Label: k,
			Value: data[k],
			Style: chart.Style{
				StrokeWidth: .01,
				FillColor: colorMariner,
				FontColor: colorWhite,
			},
		}
		i++
	}
	return points
}

func saveBarChart(ch chart.BarChart, fn string) error {
	pngFile, err := os.Create(fn)
	if err != nil {
		return nil
	}

	if err := ch.Render(chart.PNG, pngFile); err != nil {
		return nil
	}

	if err := pngFile.Close(); err != nil {
		return nil
	}

	return nil
}

func cleanup() error {
	return nil
}

