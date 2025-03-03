package application

import (
	"fmt"
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

func (o *Orchestrator) RunServer() error {
	// http.HandleFunc("/api/v1/calculate", CheckMethodIsPost(Answer500(RPNHandler)))
	// return http.ListenAndServe(":"+a.config.Addr, nil)
	fmt.Println("Lets go")
	return nil
}
