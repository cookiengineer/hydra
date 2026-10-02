package lib

import "bufio"
import "encoding/json"
import "strings"

type Result struct {
	Role    string `json:"role"`
	Name    string `json:"name"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type Summary struct {
	Passed  int
	Failed  int
	Results []Result
	Lines   []string
}

func Parse(output string) Summary {

	summary := Summary{}

	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "E2E_RESULT ") {

			payload := strings.TrimPrefix(line, "E2E_RESULT ")

			var result Result

			if err := json.Unmarshal([]byte(payload), &result); err == nil {
				summary.Results = append(summary.Results, result)
			}

			continue

		}

		if strings.HasPrefix(line, "--- PASS:") {
			summary.Passed++
			continue
		}

		if strings.HasPrefix(line, "--- FAIL:") {
			summary.Failed++
			continue
		}

		if strings.HasPrefix(line, "ok ") || line == "PASS" || line == "FAIL" || strings.HasPrefix(line, "FAIL\t") {
			summary.Lines = append(summary.Lines, line)
		}

	}

	return summary

}

func (summary Summary) OK() bool {

	if summary.Failed > 0 {
		return false
	}

	for _, result := range summary.Results {
		if !result.Success {
			return false
		}
	}

	return true

}
