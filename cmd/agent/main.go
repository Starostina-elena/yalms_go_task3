package main

import (
	"github.com/Starostina-elena/yalms_go_task2/internal/application"
)

func main() {
	agent := application.NewAgent()
	agent.Start()
}
