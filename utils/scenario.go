package utils

import (
	"slices"
	"time"
)

type ScenarioStep struct {
	URL     string            `yaml:"url"`
	Verb    string            `yaml:"verb"`
	Headers map[string]string `yaml:"headers"`
	Body    map[string]any    `yaml:"body"`
	Return  map[string]string `yaml:"return"`
	Pause   time.Duration     `yaml:"pause"`
	Order   int               `yaml:"order"`
}

type Scenario struct {
	Name        string                   `yaml:"name"`
	Concurrency int                      `yaml:"concurrency"`
	Iteration   int                      `yaml:"iteration"`
	Steps       map[string]*ScenarioStep `yaml:"steps"`
	Verbose     bool                     `yaml:"verbose"`
}

type SortInfo struct {
	Order int
	Name  string
}

func (s *Scenario) Sorted() []*SortInfo {
	sorted := make([]*SortInfo, 0, len(s.Steps))
	for k := range s.Steps {
		sorted = append(sorted, &SortInfo{
			Order: s.Steps[k].Order,
			Name:  k,
		})
	}
	slices.SortFunc(sorted, func(a, b *SortInfo) int {
		return a.Order - b.Order
	})

	return sorted
}
