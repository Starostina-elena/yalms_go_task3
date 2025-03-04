package application

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/Starostina-elena/yalms_go_task2/pkg/rpn"
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

type Expression struct {
	ID int64 `json:"id"`
	Expr string `json:"expression"`
	Status string `json:"status"`
	Result *float64 `json:"result,omitempty"`
	RPN []string `json:"-"`
}

type Task struct {
	ID int64 `json:"id"`
	ExprID int64 `json:"-"`
	Arg1 float64 `json:"arg1"`
	Arg2 float64 `json:"arg2"`
	Operation string `json:"operation"`
	OperationTime int `json:"operation_time"`
	Assigned bool `json:"-"`
}

type Orchestrator struct {
	config *Config
	mu sync.Mutex
	exprStore map[int64]*Expression
	exprCurId int64;
	taskStore map[int64]*Task
	taskCurId int64;
}

func New() *Orchestrator {
	return &Orchestrator{
		config: ConfigFromEnv(),
		exprStore: make(map[int64]*Expression),
		taskStore: make(map[int64]*Task),
	}
}

func (o *Orchestrator) CalculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Expression string `json:"expression"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Expression == "" {
		http.Error(w, "Invalid data", http.StatusUnprocessableEntity)
		return
	}
	rpn, err := rpn.Transform(req.Expression)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	o.mu.Lock()
	o.exprCurId++
	expr := &Expression{
		ID:     o.exprCurId,
		Expr:   req.Expression,
		Status: "pending",
		RPN:    rpn,
	}
	o.exprStore[o.exprCurId] = expr
	o.mu.Unlock()
	o.manageTasks(expr)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": o.exprCurId})
	
}

func isOperator(token string) bool {
	return token == "+" || token == "-" || token == "*" || token == "/"
}

func (o *Orchestrator) getOperationTime(token string) int {
	switch token {
	case "+":
		return o.config.TimePlus
	case "-":
		return o.config.TimeMinus
	case "*":
		return o.config.TimeMultiply
	case "/":
		return o.config.TimeDivide
	}
	return 100
}

func (o *Orchestrator) createTask(expr *Expression, arg1str string, arg2str string, operation string) {
	o.mu.Lock()
	o.taskCurId++
	arg1, _ := strconv.ParseFloat(arg1str, 64)
	arg2, _ := strconv.ParseFloat(arg2str, 64)
	if operation == "/" && arg2 == 0 {
		expr.Status = "error: zero division"
		for _, task := range o.taskStore {
			if task.ExprID == expr.ID {
				delete(o.taskStore, task.ID)
			}
		}
	} else {
		task := &Task{
		ID: o.taskCurId,
		ExprID: expr.ID,
		Arg1: arg1,
		Arg2: arg2,
		Operation: operation,
		OperationTime: o.getOperationTime(operation),
	}
	o.taskStore[o.taskCurId] = task
	}
	o.mu.Unlock()
}

func (o *Orchestrator) manageTasks(expr *Expression) {
	for i, token := range expr.RPN {
		if i >= 2 && isOperator(token) && !isOperator(expr.RPN[i-1]) && !isOperator(expr.RPN[i-2]) {
			o.createTask(expr, expr.RPN[i-2], expr.RPN[i-1], token)
		}
	}
}

func (o *Orchestrator) ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	o.mu.Lock()
    defer o.mu.Unlock()

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

	all_expr := make([]*Expression, 0, len(o.exprStore))
    for _, expr := range o.exprStore {
        all_expr = append(all_expr, expr)
    }

    json.NewEncoder(w).Encode(map[string]interface{}{"expressions": all_expr})
}

func (o *Orchestrator) ExpressionIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.ParseInt(r.URL.Path[len("/api/v1/expressions/"):], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	o.mu.Lock()
	expr, ok := o.exprStore[id]
	o.mu.Unlock()
	if !ok {
		http.Error(w, "Id not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"expression": expr})
}

func (o *Orchestrator) TaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		o.mu.Lock()
        defer o.mu.Unlock()

        for _, task := range o.taskStore {
            if !task.Assigned {
                task.Assigned = true
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{"task": task})
                return
            }
        }

        http.Error(w, "No tasks available", http.StatusNotFound)
	} else if r.Method == http.MethodPost {
		var req struct {
			ID int64 `json:"id"`
			Result float64 `json:"result"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid data", http.StatusUnprocessableEntity)
			return
		}

		o.mu.Lock()
		task := o.taskStore[req.ID]
		delete(o.taskStore, req.ID)
		if task == nil {
			http.Error(w, "task not found", http.StatusNotFound)
			o.mu.Unlock()
			return
		}
		expr := o.exprStore[task.ExprID]
		o.mu.Unlock()
		arg1 := strconv.FormatFloat(task.Arg1, 'f', -1, 64)
		arg2 := strconv.FormatFloat(task.Arg2, 'f', -1, 64)
		for i, token := range expr.RPN {
			if i >= 2 && token == task.Operation && expr.RPN[i-1] == arg2 && expr.RPN[i-2] == arg1 {
				expr.RPN = append(expr.RPN[:i-2], append([]string{strconv.FormatFloat(req.Result, 'f', -1, 64)}, expr.RPN[i+1:]...)...)
				if len(expr.RPN) == 1 {
					expr.Status = "done"
					expr.Result = &req.Result
				}
				if i >= 2 && i < len(expr.RPN) && isOperator(expr.RPN[i]) && !isOperator(expr.RPN[i-1]) && !isOperator(expr.RPN[i-2]) {
					o.createTask(expr, expr.RPN[i-2], expr.RPN[i-1], expr.RPN[i])
				}
				i--
				if i >= 2 && i < len(expr.RPN) && isOperator(expr.RPN[i]) && !isOperator(expr.RPN[i-1]) && !isOperator(expr.RPN[i-2]) {
					o.createTask(expr, expr.RPN[i-2], expr.RPN[i-1], expr.RPN[i])
				}
				break
			}
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func errorHandling(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		handler(w, r)
	}
}

func (o *Orchestrator) RunServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", errorHandling(o.CalculateHandler))
	mux.HandleFunc("/api/v1/expressions", errorHandling(o.ExpressionsHandler))
	mux.HandleFunc("/api/v1/expressions/", errorHandling(o.ExpressionIDHandler))
	mux.HandleFunc("/internal/task", errorHandling(o.TaskHandler))
	fmt.Println("Orchestrator started")
	return http.ListenAndServe(":"+o.config.Addr, mux)
}
