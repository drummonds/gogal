package main

import (
	"fmt"
	"time"

	"git.bytestone.uk/hum3/gogal"
)

// transaction is one posting to the account; balance is the running
// balance after it.
type transaction struct {
	day     int
	desc    string
	balance float64
}

var ledger = []transaction{
	{1, "Opening balance", 1200},
	{1, "Rent", 250},
	{3, "Groceries", 185},
	{7, "Salary", 2385},
	{8, "Car insurance", 2075},
	{12, "Groceries", 1990},
	{15, "Electricity", 1880},
	{19, "Dinner out", 1815},
	{22, "Groceries", 1735},
	{25, "Refund", 1795},
	{28, "Credit card", 1145},
	{30, "Closing balance", 1145},
}

func balanceChart() *gogal.Chart {
	var points []gogal.DataPoint
	for _, tx := range ledger {
		points = append(points, gogal.DataPoint{
			Time:  time.Date(2026, 9, tx.day, 12, 0, 0, 0, time.UTC),
			Y:     tx.balance,
			Label: fmt.Sprintf("%s: £%.0f", tx.desc, tx.balance),
		})
	}
	chart := gogal.NewStepChart(
		gogal.WithTitle("Current account — September 2026"),
		gogal.WithSize(700, 320),
		gogal.WithYTitle("£"),
		gogal.WithYFormat("%.0f"),
		gogal.WithTimeFormat("2 Jan"),
		gogal.WithGrid(true),
		gogal.WithTooltips(true),
	)
	chart.Add("Balance", points)
	return chart
}
