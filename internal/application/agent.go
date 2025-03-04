package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Agent struct {
	ComputingPower  int
	OrchestratorURL string
}

func NewAgent() *Agent {
	number_of_threads, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil || number_of_threads < 1 {
		number_of_threads = 2
	}
	orchestratorURL := os.Getenv("ORCHESTRATOR_URL")
	if orchestratorURL == "" {
		orchestratorURL = "http://localhost:8080"
	}
	return &Agent{
		ComputingPower:  number_of_threads,
		OrchestratorURL: orchestratorURL,
	}
}

func (a *Agent) Start() {
	for i := 0; i < a.ComputingPower; i++ {
		go a.runThread()
	}
	fmt.Println("Agent started")
	select{}
}

func (a *Agent) runThread() {
	for {
		resp, err := http.Get(a.OrchestratorURL + "/internal/task")
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			time.Sleep(time.Second)
			resp.Body.Close()
			continue
		}

		var response struct {
            Task struct {
                ID int64 `json:"id"`
                Arg1 float64 `json:"arg1"`
                Arg2 float64 `json:"arg2"`
                Operation string `json:"operation"`
                OperationTime int `json:"operation_time"`
            } `json:"task"`
        }

        err = json.NewDecoder(resp.Body).Decode(&response)
        if err != nil {
            time.Sleep(time.Second)
            continue
        }

		resp.Body.Close()

        task := response.Task

		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)

		var res float64
		if task.Operation == "+" {
			res = task.Arg1 + task.Arg2
		} else if task.Operation == "-" {
			res = task.Arg1 - task.Arg2
		} else if task.Operation == "*" {	
			res = task.Arg1 * task.Arg2
		} else if task.Operation == "/" {
			res = task.Arg1 / task.Arg2
		}

		result_to_send := map[string]interface{}{
			"id": task.ID,
			"result": res,
		}
		json_result_to_send, _ := json.Marshal(result_to_send)

		respPost, _ := http.Post(a.OrchestratorURL+"/internal/task", "application/json", bytes.NewReader(json_result_to_send))
		respPost.Body.Close()
	}
}
