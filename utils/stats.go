package utils

import (
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
)

type Stats struct {
	Name   string
	Status int
}

var _stats = make([]*Stats, 0)

func DoStats(ch <-chan *Stats) {
	for v := range ch {
		_stats = append(_stats, v)
	}
}

func GetStats() map[string]map[int]int {
	res := make(map[string]map[int]int)
	for i := range _stats {
		if _, ok := res[_stats[i].Name]; ok {
			if n, ok := res[_stats[i].Name][_stats[i].Status]; ok {
				res[_stats[i].Name][_stats[i].Status] = n + 1
			} else {
				res[_stats[i].Name][_stats[i].Status] = 1
			}
		} else {
			res[_stats[i].Name] = map[int]int{_stats[i].Status: 1}
		}
	}

	return res
}

func ShowStats() {
	data := make([][]string, 0)
	res := GetStats()
	total := 0
	for i := range res {
		for j := range res[i] {
			total += res[i][j]
			data = append(data, []string{i, strconv.Itoa(j), strconv.Itoa(res[i][j])})
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Step", "HTTP Status", "Count"})
	table.SetFooter([]string{"", "Total", strconv.Itoa(total)})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.AppendBulk(data)
	table.Render()
}
