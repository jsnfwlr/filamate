package stats

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/db"
	"github.com/jsnfwlr/filamate/internal/server/oapi"
)

type chartsQuerier interface {
	FindSpools(ctx context.Context) ([]db.Spool, error)
	GetMaterialChartData(ctx context.Context) ([]db.GetMaterialChartDataRow, error)
}

// CheckStorageChart does CORS preflight for storage chart by ID
// (OPTIONS /api/chart/storage)
func CheckStorageChart(ctx context.Context, dbq chartsQuerier, r oapi.CheckStorageChartRequestObject) (response oapi.CheckStorageChartResponseObject, fault error) {
	return oapi.CheckStorageChart204Response{}, nil
}

type tallies struct {
	Min     int
	Tallies []tally
}

type tally struct {
	Label     string
	YearMonth int
	Purchased int64
	Used      int64
	Stored    int64
}

func (t tallies) Months() []string {
	months := []string{}
	for _, tally := range t.Tallies {
		months = append(months, tally.Label)
	}
	return months
}

func (t tallies) Purchased() []int64 {
	purchased := []int64{}
	for _, tally := range t.Tallies {
		purchased = append(purchased, tally.Purchased)
	}
	return purchased
}

func (t tallies) Used() []int64 {
	used := []int64{}
	for _, tally := range t.Tallies {
		used = append(used, tally.Used)
	}
	return used
}

func (t tallies) Stored() []int64 {
	stored := []int64{}
	for _, tally := range t.Tallies {
		stored = append(stored, tally.Stored)
	}
	return stored
}

func (t *tallies) Set(yearMonth int, month string) {
	t.Tallies = append(t.Tallies, tally{
		Label:     month,
		YearMonth: yearMonth,
	})
}

func (t *tallies) Calc(yearMonth int, month string, purchased, used int64) {
	if yearMonth < t.Min {
		t.Tallies[0].Stored += purchased - used
		return
	}

	for i, x := range t.Tallies {
		if x.YearMonth == yearMonth {
			t.Tallies[i].Purchased += purchased
			t.Tallies[i].Used += used

			return
		}
	}
}

func (t *tallies) Finalize() {
	var used int64 = 0
	for i := range t.Tallies {
		if i != 0 {
			t.Tallies[i].Stored = t.Tallies[i-1].Stored + (t.Tallies[i-1].Purchased - used)
		}

		used = t.Tallies[i].Used

		t.Tallies[i].Used = t.Tallies[i].Used * -1
	}
}

// GetStorageChart retrieves storage chart
// (GET /api/chart/storage)
func GetStorageChart(ctx context.Context, dbq chartsQuerier, r oapi.GetStorageChartRequestObject) (response oapi.GetStorageChartResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	spools, err := dbq.FindSpools(ctx)
	if err != nil {
		o.Error("failed to find chart", err, go11y.SeverityHigh)

		return oapi.GetStorageChart500JSONResponse{
			Message: "Failed to find chart",
			Code:    500,
		}, err
	}

	o.Info("retrieved storage chart")

	endYear := time.Now().Year()
	endMonth := int(time.Now().Month())

	startYear := endYear - 1
	startMonth := endMonth + 1

	tSet := tallies{
		Min: startYear*100 + startMonth,
	}

	tSet.Set(startYear*100+startMonth, time.Month(startMonth).String())

	for i := 1; i < 12; i++ {
		month := startMonth + i
		year := startYear
		if month > 12 {
			month = month - 12
			year = year + 1
		}
		tSet.Set(year*100+month, time.Month(month).String())
	}

	for _, s := range spools {
		tSet.Calc(s.CreatedAt.Year()*100+int(s.CreatedAt.Month()), s.CreatedAt.Month().String(), 1, 0)

		if s.EmptiedAt != nil && !s.EmptiedAt.IsZero() {
			tSet.Calc(s.EmptiedAt.Year()*100+int(s.EmptiedAt.Month()), s.EmptiedAt.Month().String(), 0, 1)
		}
	}

	tSet.Finalize()

	results := oapi.StorageChartItem{
		Labels:    tSet.Months(),
		Used:      tSet.Used(),
		Purchased: tSet.Purchased(),
		Stored:    tSet.Stored(),
	}

	return oapi.GetStorageChart200JSONResponse(results), nil
}

// CheckMaterialChart does CORS preflight for material chart by ID
// (OPTIONS /api/chart/material)
func CheckMaterialChart(ctx context.Context, dbq chartsQuerier, r oapi.CheckMaterialChartRequestObject) (response oapi.CheckMaterialChartResponseObject, fault error) {
	return oapi.CheckMaterialChart204Response{}, nil
}

type hsl struct {
	hue        int
	saturation int
	lightness  int
}

var baseColors = []hsl{
	{hue: 1, saturation: 100, lightness: 15},   // "hsl(1 100%, 15%)"
	{hue: 1, saturation: 100, lightness: 30},   // "hsl(1 100%, 30%)"
	{hue: 1, saturation: 100, lightness: 50},   // "hsl(1 100%, 50%)"
	{hue: 1, saturation: 100, lightness: 70},   // "hsl(1 100%, 70%)"
	{hue: 1, saturation: 100, lightness: 85},   // "hsl(1 100%, 85%)"
	{hue: 20, saturation: 100, lightness: 15},  // "hsl(20 100%, 15%)"
	{hue: 20, saturation: 100, lightness: 30},  // "hsl(20 100%, 30%)"
	{hue: 20, saturation: 100, lightness: 50},  // "hsl(20 100%, 50%)"
	{hue: 20, saturation: 100, lightness: 70},  // "hsl(20 100%, 70%)"
	{hue: 20, saturation: 100, lightness: 85},  // "hsl(20 100%, 85%)"
	{hue: 40, saturation: 100, lightness: 15},  // "hsl(40 100%, 15%)"
	{hue: 40, saturation: 100, lightness: 30},  // "hsl(40 100%, 30%)"
	{hue: 40, saturation: 100, lightness: 50},  // "hsl(40 100%, 50%)"
	{hue: 40, saturation: 100, lightness: 70},  // "hsl(40 100%, 70%)"
	{hue: 40, saturation: 100, lightness: 85},  // "hsl(40 100%, 85%)"
	{hue: 60, saturation: 100, lightness: 15},  // "hsl(60 100%, 15%)"
	{hue: 60, saturation: 100, lightness: 30},  // "hsl(60 100%, 30%)"
	{hue: 60, saturation: 100, lightness: 50},  // "hsl(60 100%, 50%)"
	{hue: 60, saturation: 100, lightness: 70},  // "hsl(60 100%, 70%)"
	{hue: 60, saturation: 100, lightness: 85},  // "hsl(60 100%, 85%)"
	{hue: 80, saturation: 100, lightness: 15},  // "hsl(80 100%, 15%)"
	{hue: 80, saturation: 100, lightness: 30},  // "hsl(80 100%, 30%)"
	{hue: 80, saturation: 100, lightness: 50},  // "hsl(80 100%, 50%)"
	{hue: 80, saturation: 100, lightness: 70},  // "hsl(80 100%, 70%)"
	{hue: 80, saturation: 100, lightness: 85},  // "hsl(80 100%, 85%)"
	{hue: 100, saturation: 100, lightness: 15}, // "hsl(100, 100%, 15%)"
	{hue: 100, saturation: 100, lightness: 30}, // "hsl(100, 100%, 30%)"
	{hue: 100, saturation: 100, lightness: 50}, // "hsl(100, 100%, 50%)"
	{hue: 100, saturation: 100, lightness: 70}, // "hsl(100, 100%, 70%)"
	{hue: 100, saturation: 100, lightness: 85}, // "hsl(100, 100%, 85%)"
	{hue: 120, saturation: 100, lightness: 15}, // "hsl(120, 100%, 15%)"
	{hue: 120, saturation: 100, lightness: 30}, // "hsl(120, 100%, 30%)"
	{hue: 120, saturation: 100, lightness: 50}, // "hsl(120, 100%, 50%)"
	{hue: 120, saturation: 100, lightness: 70}, // "hsl(120, 100%, 70%)"
	{hue: 120, saturation: 100, lightness: 85}, // "hsl(120, 100%, 85%)"
	{hue: 140, saturation: 100, lightness: 15}, // "hsl(140, 100%, 15%)"
	{hue: 140, saturation: 100, lightness: 30}, // "hsl(140, 100%, 30%)"
	{hue: 140, saturation: 100, lightness: 50}, // "hsl(140, 100%, 50%)"
	{hue: 140, saturation: 100, lightness: 70}, // "hsl(140, 100%, 70%)"
	{hue: 140, saturation: 100, lightness: 85}, // "hsl(140, 100%, 85%)"
	{hue: 160, saturation: 100, lightness: 15}, // "hsl(160, 100%, 15%)"
	{hue: 160, saturation: 100, lightness: 30}, // "hsl(160, 100%, 30%)"
	{hue: 160, saturation: 100, lightness: 50}, // "hsl(160, 100%, 50%)"
	{hue: 160, saturation: 100, lightness: 70}, // "hsl(160, 100%, 70%)"
	{hue: 160, saturation: 100, lightness: 85}, // "hsl(160, 100%, 85%)"
	{hue: 180, saturation: 100, lightness: 15}, // "hsl(180, 100%, 15%)"
	{hue: 180, saturation: 100, lightness: 30}, // "hsl(180, 100%, 30%)"
	{hue: 180, saturation: 100, lightness: 50}, // "hsl(180, 100%, 50%)"
	{hue: 180, saturation: 100, lightness: 70}, // "hsl(180, 100%, 70%)"
	{hue: 180, saturation: 100, lightness: 85}, // "hsl(180, 100%, 85%)"
	{hue: 200, saturation: 100, lightness: 15}, // "hsl(200, 100%, 15%)"
	{hue: 200, saturation: 100, lightness: 30}, // "hsl(200, 100%, 30%)"
	{hue: 200, saturation: 100, lightness: 50}, // "hsl(200, 100%, 50%)"
	{hue: 200, saturation: 100, lightness: 70}, // "hsl(200, 100%, 70%)"
	{hue: 200, saturation: 100, lightness: 85}, // "hsl(200, 100%, 85%)"
	{hue: 220, saturation: 100, lightness: 15}, // "hsl(220, 100%, 15%)"
	{hue: 220, saturation: 100, lightness: 30}, // "hsl(220, 100%, 30%)"
	{hue: 220, saturation: 100, lightness: 50}, // "hsl(220, 100%, 50%)"
	{hue: 220, saturation: 100, lightness: 70}, // "hsl(220, 100%, 70%)"
	{hue: 220, saturation: 100, lightness: 85}, // "hsl(220, 100%, 85%)"
	{hue: 240, saturation: 100, lightness: 15}, // "hsl(240, 100%, 15%)"
	{hue: 240, saturation: 100, lightness: 30}, // "hsl(240, 100%, 30%)"
	{hue: 240, saturation: 100, lightness: 50}, // "hsl(240, 100%, 50%)"
	{hue: 240, saturation: 100, lightness: 70}, // "hsl(240, 100%, 70%)"
	{hue: 240, saturation: 100, lightness: 85}, // "hsl(240, 100%, 85%)"
	{hue: 260, saturation: 100, lightness: 15}, // "hsl(260, 100%, 15%)"
	{hue: 260, saturation: 100, lightness: 30}, // "hsl(260, 100%, 30%)"
	{hue: 260, saturation: 100, lightness: 50}, // "hsl(260, 100%, 50%)"
	{hue: 260, saturation: 100, lightness: 70}, // "hsl(260, 100%, 70%)"
	{hue: 260, saturation: 100, lightness: 85}, // "hsl(260, 100%, 85%)"
	{hue: 280, saturation: 100, lightness: 15}, // "hsl(280, 100%, 15%)"
	{hue: 280, saturation: 100, lightness: 30}, // "hsl(280, 100%, 30%)"
	{hue: 280, saturation: 100, lightness: 50}, // "hsl(280, 100%, 50%)"
	{hue: 280, saturation: 100, lightness: 70}, // "hsl(280, 100%, 70%)"
	{hue: 280, saturation: 100, lightness: 85}, // "hsl(280, 100%, 85%)"
	{hue: 300, saturation: 100, lightness: 15}, // "hsl(300, 100%, 15%)"
	{hue: 300, saturation: 100, lightness: 30}, // "hsl(300, 100%, 30%)"
	{hue: 300, saturation: 100, lightness: 50}, // "hsl(300, 100%, 50%)"
	{hue: 300, saturation: 100, lightness: 70}, // "hsl(300, 100%, 70%)"
	{hue: 300, saturation: 100, lightness: 85}, // "hsl(300, 100%, 85%)"
	{hue: 320, saturation: 100, lightness: 15}, // "hsl(320, 100%, 15%)"
	{hue: 320, saturation: 100, lightness: 30}, // "hsl(320, 100%, 30%)"
	{hue: 320, saturation: 100, lightness: 50}, // "hsl(320, 100%, 50%)"
	{hue: 320, saturation: 100, lightness: 70}, // "hsl(320, 100%, 70%)"
	{hue: 320, saturation: 100, lightness: 85}, // "hsl(320, 100%, 85%)"
	{hue: 340, saturation: 100, lightness: 15}, // "hsl(340, 100%, 15%)"
	{hue: 340, saturation: 100, lightness: 30}, // "hsl(340, 100%, 30%)"
	{hue: 340, saturation: 100, lightness: 50}, // "hsl(340, 100%, 50%)"
	{hue: 340, saturation: 100, lightness: 70}, // "hsl(340, 100%, 70%)"
	{hue: 340, saturation: 100, lightness: 85}, // "hsl(340, 100%, 85%)"
}

type assignedColors struct {
	first hsl
	last  hsl
}

func GetMaterialChart(ctx context.Context, dbq chartsQuerier, r oapi.GetMaterialChartRequestObject) (response oapi.GetMaterialChartResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	details, err := dbq.GetMaterialChartData(ctx)
	if err != nil {
		o.Error("failed to find chart", err, go11y.SeverityHigh)

		return oapi.GetMaterialChart500JSONResponse{
			Message: "Failed to find chart",
			Code:    500,
		}, err
	}

	o.Info("retrieved material chart", "details", details)

	labels := []string{}
	datasets := []oapi.MaterialChartDataset{}

	dataSetColors := []assignedColors{}
	usedColors := 0

	dataSetLabels := [][]string{}

	for _, d := range details {
		parts := []string{}
		if d.Class != "" {
			parts = append(parts, d.Class)
		}

		if d.Material != "" {
			parts = append(parts, d.Material)
		}

		if d.Brand != "" {
			parts = append(parts, d.Brand)
		}

		label := strings.Join(parts, "/")

		if label == "" {
			continue
		}

		if len(dataSetLabels) < len(parts) {
			dataSetLabels = append(dataSetLabels, []string{label})
		} else {
			dataSetLabels[len(parts)-1] = append(dataSetLabels[len(parts)-1], label)
		}

		if len(datasets) < len(parts) {
			color := baseColors[usedColors]
			usedColors++
			dataSetColors = append(dataSetColors, assignedColors{
				first: color,
				last:  color,
			})
			datasets = append(datasets, oapi.MaterialChartDataset{
				Data: []int64{d.Count},
				BackgroundColor: []string{
					"hsl(" + fmt.Sprint(rune(color.hue)) + " " + fmt.Sprint(rune(color.saturation)) + "% " + fmt.Sprint(rune(color.lightness)) + "%)",
				},
			})
		} else {
			color := baseColors[usedColors]
			dataSetColors[len(parts)-1].last = color
			datasets[len(parts)-1].Data = append(datasets[len(parts)-1].Data, d.Count)
			datasets[len(parts)-1].BackgroundColor = append(datasets[len(parts)-1].BackgroundColor, "hsl("+fmt.Sprint(rune(color.hue))+" "+fmt.Sprint(rune(color.saturation))+"% "+fmt.Sprint(rune(color.lightness))+"%)")
			usedColors++
		}
	}

	for d := range dataSetLabels {
		if dataSetLabels[d] == nil {
			continue
		}

		for _, label := range dataSetLabels[d] {
			if !slices.Contains(labels, label) {
				labels = append(labels, label)
			}
		}
	}

	results := oapi.MaterialChart{
		Labels:   labels,
		Datasets: datasets,
	}

	return oapi.GetMaterialChart200JSONResponse(results), nil
}
