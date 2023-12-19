package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/meysampg/testapi/utils"
	"log"
	"net/http"
	"strings"
	"sync"
	"text/template"
	"time"
)

var scenarioFolder string

func init() {
	flag.StringVar(&scenarioFolder, "scenarios", "./scenarios", "location of scenarios folder")
	flag.Parse()
}

func main() {
	yamls, err := utils.ReadYAMLFiles(scenarioFolder)
	if err != nil {
		log.Fatal(err)
	}

	wg := new(sync.WaitGroup)
	statCh := make(chan *utils.Stats)

	go utils.DoStats(statCh)

	for file := range yamls {
		if yamls[file].Name == "" {
			yamls[file].Name = file
		}
		if yamls[file].Concurrency == 0 {
			yamls[file].Concurrency = 1
		}

		for i := 1; i <= yamls[file].Iteration; i++ {
			wg.Add(1)
			go runScenario(wg, statCh, i, yamls[file])
		}
	}

	wg.Wait()
	close(statCh)
	utils.ShowStats()
}

func runScenario(wg *sync.WaitGroup, statCh chan<- *utils.Stats, round int, scenario *utils.Scenario) {
	defer wg.Done()
	for i := 1; i <= scenario.Concurrency; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, statCh chan<- *utils.Stats, round, i int, scenario *utils.Scenario) {
			defer wg.Done()
			runSteps(statCh, round, i, scenario)
		}(wg, statCh, round, i, scenario)
	}
}

func runSteps(statCh chan<- *utils.Stats, round, n int, scenario *utils.Scenario) {
	data := make(map[string]map[string]any)
	sorted := scenario.Sorted()

	for s := range sorted {
		name := sorted[s].Name
		step := processStep(data, scenario.Steps[name])
		req, err := makeRequest(step)
		if err != nil {
			log.Printf("[%s][%d-%d] error on making request, %s\n", name, round, n, err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("[%s][%d-%d] error on doing request, %s\n", name, round, n, err)
			return
		}

		respText := make(map[string]any)
		err = json.NewDecoder(resp.Body).Decode(&respText)
		if err != nil {
			fmt.Printf("[%s][%d-%d] error on decoding response body, %s\n", name, round, n, err)
			return
		}
		statCh <- &utils.Stats{
			Name:   name,
			Status: resp.StatusCode,
		}
		if resp.StatusCode >= 400 {
			respErr := "could be shown using verbose flag."
			if len(respText) > 0 && scenario.Verbose {
				r, _ := json.Marshal(respText)
				respErr = fmt.Sprintf("is %s", r)
			}
			log.Printf("[%s][%d-%d] error on sending request to %s, status is %d and response %v\n", name, round, n, step.URL, resp.StatusCode, respErr)
			if step.Pause >= 1*time.Millisecond {
				log.Printf("[%s][%d-%d] Sleep for %v.\n", name, round, n, step.Pause)
				time.Sleep(step.Pause)
			}
			return
		}
		data[name] = respText
		resp.Body.Close()
		log.Printf("[%s][%d-%d] Done.\n", name, round, n)
		if step.Pause >= 1*time.Millisecond {
			log.Printf("[%s][%d-%d] Sleep for %v.\n", name, round, n, step.Pause)
			time.Sleep(step.Pause)
		}
	}
}

func makeRequest(scenario *utils.ScenarioStep) (*http.Request, error) {
	data, err := json.Marshal(scenario.Body)
	if err != nil && len(scenario.Body) > 0 {
		return nil, err
	}
	req, err := http.NewRequest(strings.ToUpper(scenario.Verb), scenario.URL, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	for key, value := range scenario.Headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

func processStep(data map[string]map[string]any, step *utils.ScenarioStep) *utils.ScenarioStep {
	return &utils.ScenarioStep{
		URL:     fmt.Sprintf("%v", processLine(data, step.URL)),
		Verb:    fmt.Sprintf("%v", processLine(data, step.Verb)),
		Headers: step.Headers,
		Body:    processMap(data, step.Body),
		Pause:   step.Pause,
	}
}

func processLine(data map[string]map[string]any, line any) any {
	switch v := line.(type) {
	case string:
		tmpl, err := template.New("line").Funcs(utils.FuncMaps).Parse(v)
		if err != nil {
			log.Println("error on creating template", err)
			return ""
		}

		var outputBuffer bytes.Buffer
		err = tmpl.Execute(&outputBuffer, data)
		if err != nil {
			log.Println("error on interpolating", err)
			return ""
		}

		return outputBuffer.String()

	default:
		return v
	}
}

func processMap(data map[string]map[string]any, m map[string]any) map[string]any {
	result := make(map[string]any)
	for k := range m {
		switch v := m[k].(type) {
		case string:
			result[k] = processLine(data, v)
		default:
			result[k] = v
		}
	}

	return result
}
