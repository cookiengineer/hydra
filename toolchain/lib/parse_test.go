package lib

import "testing"

func TestParseResult(t *testing.T) {

	output := "=== RUN   TestE2E\n" +
		"E2E_RESULT {\"role\":\"server\",\"name\":\"register\",\"success\":true,\"message\":\"ok\"}\n" +
		"--- PASS: TestE2E (0.10s)\n" +
		"PASS\n"

	summary := Parse(output)

	if summary.Passed != 1 {
		t.Errorf("Expected 1 pass, got %d", summary.Passed)
	}

	if summary.Failed != 0 {
		t.Errorf("Expected 0 fail, got %d", summary.Failed)
	}

	if len(summary.Results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(summary.Results))
	}

	if summary.Results[0].Role != "server" {
		t.Errorf("Expected role server, got %s", summary.Results[0].Role)
	}

	if !summary.OK() {
		t.Error("Expected summary to be OK")
	}

}

func TestParseFailure(t *testing.T) {

	output := "E2E_RESULT {\"role\":\"client\",\"name\":\"connect\",\"success\":false,\"message\":\"nope\"}\n" +
		"--- FAIL: TestE2E (0.10s)\n" +
		"FAIL\n"

	summary := Parse(output)

	if summary.Failed != 1 {
		t.Errorf("Expected 1 fail, got %d", summary.Failed)
	}

	if summary.OK() {
		t.Error("Expected summary to be not OK")
	}

}

func TestParseEmpty(t *testing.T) {

	summary := Parse("")

	if !summary.OK() {
		t.Error("Expected empty summary to be OK")
	}

}
