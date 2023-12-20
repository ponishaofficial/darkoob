package utils

import (
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
	"sync"
)

type Stats struct {
	Scenario string
	Step     string
	Status   string
}

var _stats = make(map[string]map[string]map[string]int)
var mu = new(sync.Mutex)

func DoStats(ch <-chan *Stats) {
	for v := range ch {
		mu.Lock()
		if _, ok := _stats[v.Scenario]; ok {
			if _, ok := _stats[v.Scenario][v.Step]; ok {
				if n, ok := _stats[v.Scenario][v.Step][v.Status]; ok {
					_stats[v.Scenario][v.Step][v.Status] = n + 1
				} else {
					_stats[v.Scenario][v.Step][v.Status] = 1
				}
			} else {
				_stats[v.Scenario][v.Step] = map[string]int{v.Status: 1}
			}
		} else {
			_stats[v.Scenario] = map[string]map[string]int{
				v.Step: {
					v.Status: 1,
				},
			}
		}
		mu.Unlock()
	}
}

func ShowStats(scenarios map[string]*Scenario) {
	data := make([][]string, 0)
	total := 0
	mu.Lock()
	for scenarioName := range scenarios {
		scenario := scenarios[scenarioName].Name
		if _, ok := _stats[scenario]; ok {
			sorted := scenarios[scenarioName].Sorted()
			for i := range sorted {
				if _, ok := _stats[scenario][sorted[i].Name]; ok {
					for status := range _stats[scenario][sorted[i].Name] {
						total += _stats[scenario][sorted[i].Name][status]
						data = append(data, []string{scenario, sorted[i].Name, status, strconv.Itoa(_stats[scenario][sorted[i].Name][status])})
					}
				}
			}
		}
	}
	mu.Unlock()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Scenario", "Step", "HTTP Status", "Count"})
	table.SetFooter([]string{"", "", "Total", strconv.Itoa(total)})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.AppendBulk(data)
	table.Render()
}
