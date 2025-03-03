package application

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	Addr string
	TimePlus int
	TimeMinus int
	TimeMultiply int
	TimeDivide int
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	config.TimePlus, _ = strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
	if config.TimePlus == 0 {
		config.TimePlus = 5000
	}
	config.TimeMinus, _ = strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
	if config.TimeMinus == 0 {
		config.TimeMinus = 5000
	}
	config.TimeMultiply, _ = strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
	if config.TimeMultiply == 0 {
		config.TimeMultiply = 5000
	}
	config.TimeDivide, _ = strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))
	if config.TimeDivide == 0 {
		config.TimeDivide = 5000
	}
	return config
}

type Orchestrator struct {
	config *Config
}

func New() *Orchestrator {
	return &Orchestrator{
		config: ConfigFromEnv(),
	}
}

func (o *Orchestrator) CalculateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CalculateHandler")
}

func (o *Orchestrator) ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ExpressionsHandler")
}

func (o *Orchestrator) ExpressionIDHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ExpressionIDHandler")
}

func (o *Orchestrator) TaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Println("TaskHandler GET")
	} else if r.Method == http.MethodPost {
		fmt.Println("TaskHandler POST")
	}
}

func (o *Orchestrator) RunServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", o.CalculateHandler)
	mux.HandleFunc("/api/v1/expressions", o.ExpressionsHandler)
	mux.HandleFunc("/api/v1/expressions/", o.ExpressionIDHandler)
	mux.HandleFunc("/internal/task", o.TaskHandler)
	return http.ListenAndServe(":"+o.config.Addr, mux)
}
