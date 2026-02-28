package stats

import (
	"context"
	"fmt"
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

// var baseColors = []hsl{
// 	{hue: 0, saturation: 100, lightness: 40},   // "hsl(0 100%, 40%)"
// 	{hue: 15, saturation: 100, lightness: 40},  // "hsl(15 100%, 40%)"
// 	{hue: 30, saturation: 100, lightness: 40},  // "hsl(30 100%, 40%)"
// 	{hue: 45, saturation: 100, lightness: 40},  // "hsl(45 100%, 40%)"
// 	{hue: 60, saturation: 100, lightness: 40},  // "hsl(60 100%, 40%)"
// 	{hue: 75, saturation: 100, lightness: 40},  // "hsl(75 100%, 40%)"
// 	{hue: 90, saturation: 100, lightness: 40},  // "hsl(90 100%, 40%)"
// 	{hue: 105, saturation: 100, lightness: 40}, // "hsl(105, 100%, 40%)"
// 	{hue: 120, saturation: 100, lightness: 40}, // "hsl(120, 100%, 40%)"
// 	{hue: 135, saturation: 100, lightness: 40}, // "hsl(135, 100%, 40%)"
// 	{hue: 150, saturation: 100, lightness: 40}, // "hsl(150, 100%, 40%)"
// 	{hue: 165, saturation: 100, lightness: 40}, // "hsl(165, 100%, 40%)"
// 	{hue: 180, saturation: 100, lightness: 40}, // "hsl(180, 100%, 40%)"
// 	{hue: 195, saturation: 100, lightness: 40}, // "hsl(195, 100%, 40%)"
// 	{hue: 210, saturation: 100, lightness: 40}, // "hsl(210, 100%, 40%)"
// 	{hue: 225, saturation: 100, lightness: 40}, // "hsl(225, 100%, 40%)"
// 	{hue: 240, saturation: 100, lightness: 40}, // "hsl(240, 100%, 40%)"
// 	{hue: 255, saturation: 100, lightness: 40}, // "hsl(255, 100%, 40%)"
// 	{hue: 270, saturation: 100, lightness: 40}, // "hsl(270, 100%, 40%)"
// 	{hue: 285, saturation: 100, lightness: 40}, // "hsl(285, 100%, 40%)"
// 	{hue: 300, saturation: 100, lightness: 40}, // "hsl(300, 100%, 40%)"
// 	{hue: 315, saturation: 100, lightness: 40}, // "hsl(315, 100%, 40%)"
// 	{hue: 330, saturation: 100, lightness: 40}, // "hsl(330, 100%, 40%)"

// }

func makeBaseColors(details []db.GetMaterialChartDataRow) []hsl {
	count := 0
	for _, d := range details {
		if d.Brand != "" {
			count++
		}
	}

	colors := []hsl{}

	for i := range count {
		colors = append(colors, hsl{
			hue:        (i * 360 / count) % 360,
			saturation: 100,
			lightness:  60,
		})
	}

	return colors
}

type assignedColors struct {
	first hsl
	last  hsl
}

type DataSet struct {
	BackgroundColor []string
	Data            []int64
	Label           []string
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

	// slices.Reverse(details)

	o.Info("retrieved material chart", "details", details)

	datasets := [3]DataSet{}
	labels := []string{}
	dataSetColors := map[string]assignedColors{}
	usedColors := 0

	baseColors := makeBaseColors(details)

	for _, d := range details {
		parts := []string{}
		label := ""
		{ // build parts and label
			if d.Class != "" {
				parts = append(parts, d.Class)
			}

			if d.Material != "" {
				parts = append(parts, d.Material)
			}

			if d.Brand != "" {
				parts = append(parts, d.Brand)
			}

			label = strings.Join(parts, "/")
		}
		if label == "" {
			continue
		}

		idx := len(parts) - 1
		color := hsl{hue: 0, saturation: 0, lightness: 0}

		switch len(parts) {
		case 3:
			color = baseColors[usedColors]
			usedColors++

			if x, ok := dataSetColors[strings.Join(parts[:2], "/")]; !ok {
				dataSetColors[strings.Join(parts[:2], "/")] = assignedColors{
					first: color,
					last:  color,
				}
			} else {
				x.last = color
				dataSetColors[strings.Join(parts[:2], "/")] = x
			}
		case 2:
			colors := dataSetColors[strings.Join(parts, "/")]
			color = hsl{
				hue:        (colors.first.hue + colors.last.hue) / 2,
				saturation: (colors.first.saturation + colors.last.saturation) / 2,
				lightness:  (colors.first.lightness - 15),
			}

			if x, ok := dataSetColors[strings.Join(parts[:1], "/")]; !ok {
				dataSetColors[strings.Join(parts[:1], "/")] = assignedColors{
					first: color,
					last:  color,
				}
			} else {
				x.last = color
				dataSetColors[strings.Join(parts[:1], "/")] = x
			}
		case 1:
			colors := dataSetColors[strings.Join(parts, "/")]
			color = hsl{
				hue:        (colors.first.hue + colors.last.hue) / 2,
				saturation: (colors.first.saturation + colors.last.saturation) / 2,
				lightness:  (colors.first.lightness - 15),
			}
		}

		// if len(datasets) < len(parts) {
		// 	datasets = append(datasets, oapi.MaterialChartDataset{
		// 		Data: []int64{d.Count},
		// 		BackgroundColor: []string{
		// 			"hsl(" + fmt.Sprint(rune(color.hue)) + " " + fmt.Sprint(rune(color.saturation)) + "% " + fmt.Sprint(rune(color.lightness)) + "%)",
		// 		},
		// 	})
		// } else {

		datasets[idx].Label = append(datasets[idx].Label, label)
		datasets[idx].Data = append(datasets[idx].Data, d.Count)
		datasets[idx].BackgroundColor = append(datasets[idx].BackgroundColor, "hsl("+fmt.Sprint(rune(color.hue))+" "+fmt.Sprint(rune(color.saturation))+"% "+fmt.Sprint(rune(color.lightness))+"%)")
		// }
	}

	// slices.Reverse(dataSetLabels)

	mcds := []oapi.MaterialChartDataset{}

	for i := 2; i >= 0; i-- {
		fmt.Printf("%d = %d vs %d vs %d\n", i, len(datasets[i].Label), len(datasets[i].Data), len(datasets[i].BackgroundColor))
		labels = append(labels, datasets[i].Label...)
		mcds = append(mcds, oapi.MaterialChartDataset{
			Data:            datasets[i].Data,
			BackgroundColor: datasets[i].BackgroundColor,
		})
	}

	results := oapi.MaterialChart{
		Labels:   labels,
		Datasets: mcds,
	}

	return oapi.GetMaterialChart200JSONResponse(results), nil
}
