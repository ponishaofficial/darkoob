package utils

import (
	"slices"
	"time"
)

type ScenarioStep struct {
	Order     int               `yaml:"order"`
	URL       string            `yaml:"url"`
	Verb      string            `yaml:"verb"`
	Pause     time.Duration     `yaml:"pause"`
	Variables map[string]any    `yaml:"variables"`
	Headers   map[string]string `yaml:"headers"`
	Body      map[string]any    `yaml:"body"`
}

type Scenario struct {
	Name        string                   `yaml:"name"`
	Concurrency int                      `yaml:"concurrency"`
	SilentRun   bool                     `yaml:"silent"`
	Iteration   int                      `yaml:"iteration"`
	Steps       map[string]*ScenarioStep `yaml:"steps"`
	Verbose     bool                     `yaml:"verbose"`

	sortedScenario []*SortInfo
}

type SortInfo struct {
	Order int
	Name  string
}

func (s *Scenario) Sorted() []*SortInfo {
	if s.sortedScenario == nil {
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
		s.sortedScenario = sorted
	}

	return s.sortedScenario
}
