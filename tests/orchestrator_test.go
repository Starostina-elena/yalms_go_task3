package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"github.com/Starostina-elena/yalms_go_task2/internal/application"
)

func TestCalculateHandler(t *testing.T) {
    orchestrator := application.New()
    server := httptest.NewServer(http.HandlerFunc(application.ErrorHandling(orchestrator.CalculateHandler)))
    defer server.Close()

	testCases := []struct {
        statusCode int
        expression string
    }{
        {201, "3+4"},
        {201, "3*5*6-4/3+3*7"},
        {201, "5+(3-2)*5*(8-3)"},
        {422, "2jdbjfb+a"},
        {422, ")("},
    }

	for _, tc := range testCases {
		payload := map[string]string{"expression": tc.expression}
		payloadBytes, _ := json.Marshal(payload)
		resp, err := http.Post(server.URL+"/api/v1/calculate", "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != tc.statusCode {
			t.Fatalf("Expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
		}

	}
}

func TestExpressionsHandler(t *testing.T) {
    orchestrator := application.New()
    server := httptest.NewServer(http.HandlerFunc(application.ErrorHandling(orchestrator.ExpressionsHandler)))
    defer server.Close()

    resp, err := http.Get(server.URL + "/api/v1/expressions")
    if err != nil {
        t.Fatalf("Failed to send request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
    }

    var result map[string][]application.Expression
    err = json.NewDecoder(resp.Body).Decode(&result)
    if err != nil {
        t.Fatalf("Failed to decode response: %v", err)
    }

    if _, ok := result["expressions"]; !ok {
        t.Fatalf("Expected response to contain 'expressions'")
    }
}

func TestExpressionIDHandler(t *testing.T) {
    orchestrator := application.New()
    mux := http.NewServeMux()
    mux.HandleFunc("/api/v1/calculate", application.ErrorHandling(orchestrator.CalculateHandler))
    mux.HandleFunc("/api/v1/expressions/", application.ErrorHandling(orchestrator.ExpressionIDHandler))
    server := httptest.NewServer(mux)
    defer server.Close()

    payload := map[string]string{"expression": "2*5"}
    payloadBytes, _ := json.Marshal(payload)
    resp, err := http.Post(server.URL+"/api/v1/calculate", "application/json", bytes.NewBuffer(payloadBytes))
    if err != nil {
        t.Fatalf("Failed to send request: %v", err)
    }
    defer resp.Body.Close()

    var result map[string]int64
    err = json.NewDecoder(resp.Body).Decode(&result)
    if err != nil {
        t.Fatalf("Failed to decode response: %v", err)
    }

    exprID := result["id"]

    resp, err = http.Get(server.URL + "/api/v1/expressions/" + strconv.FormatInt(exprID, 10))
    if err != nil {
        t.Fatalf("Failed to send request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
    }

    var exprResult map[string]application.Expression
    err = json.NewDecoder(resp.Body).Decode(&exprResult)
    if err != nil {
        t.Fatalf("Failed to decode response: %v", err)
    }

    if _, ok := exprResult["expression"]; !ok {
        t.Fatalf("Expected response to contain 'expression'")
    }
}

func TestTaskHandler(t *testing.T) {
    orchestrator := application.New()
    mux := http.NewServeMux()
    mux.HandleFunc("/api/v1/calculate", application.ErrorHandling(orchestrator.CalculateHandler))
    mux.HandleFunc("/internal/task", application.ErrorHandling(orchestrator.TaskHandler))
    server := httptest.NewServer(mux)
    defer server.Close()

    payload := map[string]string{"expression": "3*5"}
    payloadBytes, _ := json.Marshal(payload)
    resp, err := http.Post(server.URL+"/api/v1/calculate", "application/json", bytes.NewBuffer(payloadBytes))
    if err != nil {
        t.Fatalf("Failed to send request: %v", err)
    }
    defer resp.Body.Close()

    var result map[string]int64
    err = json.NewDecoder(resp.Body).Decode(&result)
    if err != nil {
        t.Fatalf("Failed to decode response: %v", err)
    }

    resp, err = http.Get(server.URL + "/internal/task")
    if err != nil {
        t.Fatalf("Failed to send request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
    }

    var taskResult map[string]application.Task
    err = json.NewDecoder(resp.Body).Decode(&taskResult)
    if err != nil {
        t.Fatalf("Failed to decode response: %v", err)
    }

    if _, ok := taskResult["task"]; !ok {
        t.Fatalf("Expected response to contain 'task'")
    }
}
