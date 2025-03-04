package tests

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
	"bytes"
    "github.com/Starostina-elena/yalms_go_task2/internal/application"
)

func TestAgentStart(t *testing.T) {
    orchestrator := application.New()
    mux := http.NewServeMux()
    mux.HandleFunc("/internal/task", application.ErrorHandling(orchestrator.TaskHandler))
    server := httptest.NewServer(mux)
    defer server.Close()

    agent := application.NewAgent()
    agent.OrchestratorURL = server.URL

    go agent.Start()

    time.Sleep(1 * time.Second)

    resp, err := http.Get(server.URL + "/internal/task")
    if err != nil {
        t.Fatalf("Failed to send request: %v", err)
    }
    defer resp.Body.Close()
}

func TestAgentRunThread(t *testing.T) {
    orchestrator := application.New()
    mux := http.NewServeMux()
    mux.HandleFunc("/api/v1/calculate", application.ErrorHandling(orchestrator.CalculateHandler))
    mux.HandleFunc("/internal/task", application.ErrorHandling(orchestrator.TaskHandler))
    server := httptest.NewServer(mux)
    defer server.Close()

    payload := map[string]string{"expression": "2*5"}
    payloadBytes, _ := json.Marshal(payload)
    resp, err := http.Post(server.URL+"/api/v1/calculate", "application/json", bytes.NewBuffer(payloadBytes))
    if err != nil {
        t.Fatalf("Failed to send request: %v", err)
    }
    defer resp.Body.Close()

    agent := application.NewAgent()
    agent.OrchestratorURL = server.URL

    go agent.RunThread()

    time.Sleep(1 * time.Second)

    resp, err = http.Get(server.URL + "/internal/task")
    if err != nil {
        t.Fatalf("Failed to send request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusNotFound {
        t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
    }
}
