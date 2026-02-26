package stats

import (
	"context"
	"time"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/db"
	"github.com/jsnfwlr/filamate/internal/server/oapi"
)

type chartsQuerier interface {
	FindSpools(ctx context.Context) ([]db.Spool, error)
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
