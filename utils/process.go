package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
)

func RunScenario(ctx context.Context, wg *sync.WaitGroup, statCh chan<- *Stats, round int, scenario *Scenario) {
	defer wg.Done()
	for i := 1; i <= scenario.Concurrency; i++ {
		wg.Add(1)
		go runSteps(ctx, wg, statCh, round, i, scenario)
	}
}

func runSteps(ctx context.Context, wg *sync.WaitGroup, statCh chan<- *Stats, round, clientID int, scenario *Scenario) {
	defer wg.Done()
	data := make(map[string]map[string]any)
	sorted := scenario.Sorted()
	isNotSilent := !scenario.SilentRun
	count := 0
	defer func(scenario *Scenario, round int, count *int) {
		if isNotSilent {
			log.Printf("🛈 On %s, round %d, %d request sent.\n", scenario.Name, round, *count)
		}
	}(scenario, round, &count)

	for s := range sorted {
		count++
		name := sorted[s].Name
		step := processStep(data, scenario.Steps[name])
		req, err := makeRequest(ctx, step)
		if err != nil && isNotSilent {
			log.Printf("[%s][%s][R%dC%d] error on making request, %s\n", scenario.Name, name, round, clientID, err)
			sendStat(wg, statCh, scenario.Name, name, "RequestError")
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil && isNotSilent {
			log.Printf("[%s][%s][R%dC%d] error on doing request, %s\n", scenario.Name, name, round, clientID, err)
			sendStat(wg, statCh, scenario.Name, name, "ClientError")
			return
		}

		sendStat(wg, statCh, scenario.Name, name, strconv.Itoa(resp.StatusCode))
		respText := make(map[string]any)
		err = json.NewDecoder(resp.Body).Decode(&respText)
		if err != nil && isNotSilent {
			fmt.Printf("[%s][%s][R%dC%d] error on decoding response body, %s\n", scenario.Name, name, round, clientID, err)
			return
		}
		if resp.StatusCode >= 400 && isNotSilent {
			respErr := "could be shown using verbose flag."
			if len(respText) > 0 && scenario.Verbose {
				r, _ := json.Marshal(respText)
				respErr = fmt.Sprintf("is %s", r)
			}
			log.Printf("[%s][%s][R%dC%d] error on sending request to %s, status is %d and response %v\n", scenario.Name, name, round, clientID, step.URL, resp.StatusCode, respErr)
			if step.Pause >= 1*time.Millisecond && isNotSilent {
				log.Printf("[%s][%s][R%dC%d] Sleep for %v.\n", scenario.Name, name, round, clientID, step.Pause)
				time.Sleep(step.Pause)
			}
			return
		}
		data[name] = respText
		resp.Body.Close()
		if isNotSilent {
			log.Printf("[%s][%s][R%dC%d] Done.\n", scenario.Name, name, round, clientID)
		}
		if step.Pause >= 1*time.Millisecond && isNotSilent {
			log.Printf("[%s][%s][R%dC%d] Sleep for %v.\n", scenario.Name, name, round, clientID, step.Pause)
			time.Sleep(step.Pause)
		}
	}
}

func makeRequest(ctx context.Context, scenario *ScenarioStep) (*http.Request, error) {
	data, err := json.Marshal(scenario.Body)
	if err != nil && len(scenario.Body) > 0 {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(scenario.Verb), scenario.URL, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	for key, value := range scenario.Headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

func processStep(data map[string]map[string]any, step *ScenarioStep) *ScenarioStep {
	return &ScenarioStep{
		URL:     fmt.Sprintf("%v", processLine(data, step.URL)),
		Verb:    fmt.Sprintf("%v", processLine(data, step.Verb)),
		Headers: processMap(data, step.Headers),
		Body:    processMap(data, step.Body),
		Pause:   step.Pause,
	}
}

func processLine[T ~string | any](data map[string]map[string]any, line T) T {
	switch v := any(line).(type) {
	case string:
		tmpl, err := template.New("line").Funcs(FuncMaps).Parse(v)
		if err != nil {
			log.Println("error on creating template", err)
			return any("").(T)
		}

		var outputBuffer bytes.Buffer
		err = tmpl.Execute(&outputBuffer, data)
		if err != nil {
			log.Println("error on interpolating", err)
			return any("").(T)
		}

		return any(outputBuffer.String()).(T)

	default:
		return v.(T)
	}
}

func processMap[T ~string | any](data map[string]map[string]any, m map[string]T) map[string]T {
	result := make(map[string]T)
	for k := range m {
		switch v := any(m[k]).(type) {
		case string:
			result[k] = any(processLine(data, v)).(T)
		default:
			result[k] = v.(T)
		}
	}

	return result
}

func sendStat(wg *sync.WaitGroup, ch chan<- *Stats, scenario, step, status string) {
	wg.Add(1)
	go func(wg *sync.WaitGroup, ch chan<- *Stats, scenario, step, status string) {
		defer wg.Done()
		ch <- &Stats{
			Scenario: scenario,
			Step:     step,
			Status:   status,
		}
	}(wg, ch, scenario, step, status)
}
