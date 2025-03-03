package rpn

import (
 "errors"
 "strconv"
 "strings"
 "unicode"
)

func Calc(expression string) (float64, error) {
	tokens, err := tokenize(expression)
	if err != nil {
		return 0, err
 	}

	rpn, err := toRPN(tokens)
 	if err != nil {
  		return 0, err
 	}

 	return evaluateRPN(rpn)
}

func tokenize(expression string) ([]string, error) {
	var tokens []string
 	var number strings.Builder

    for _, ch := range expression {
        if unicode.IsSpace(ch) {
            continue
        } else if unicode.IsDigit(ch) || ch == '.' {
            number.WriteRune(ch)
        } else {
            if number.Len() > 0 {
            tokens = append(tokens, number.String())
            number.Reset()
        }
        if strings.ContainsRune("+-*/()", ch) {
            tokens = append(tokens, string(ch))
        } else {
            return nil, errors.New("недопустимый символ в выражении")
        }
        }
    }

    if number.Len() > 0 {
        tokens = append(tokens, number.String())
    }

    return tokens, nil
}

func toRPN(tokens []string) ([]string, error) {
    var output []string
    var stack []string

    precedence := map[string]int{
        "+": 1,
        "-": 1,
        "*": 2,
        "/": 2,
    }

    for _, token := range tokens {
        switch token {
            case "+", "-", "*", "/":
                for len(stack) > 0 && precedence[stack[len(stack)-1]] >= precedence[token] {
                    output = append(output, stack[len(stack)-1])
                    stack = stack[:len(stack)-1]
                }
                stack = append(stack, token)
            case "(":
                stack = append(stack, token)
            case ")":
                for len(stack) > 0 && stack[len(stack)-1] != "(" {
                output = append(output, stack[len(stack)-1])
                stack = stack[:len(stack)-1]
                }
                if len(stack) == 0 {
                    return nil, errors.New("несоответствующие скобки")
                }
                stack = stack[:len(stack)-1]
            default:
                output = append(output, token)
        }
    }

    for len(stack) > 0 {
        if stack[len(stack)-1] == "(" {
            return nil, errors.New("несоответствующие скобки")
        }
        output = append(output, stack[len(stack)-1])
        stack = stack[:len(stack)-1]
    }

    return output, nil
}

func evaluateRPN(tokens []string) (float64, error) {
    var stack []float64

    for _, token := range tokens {
        switch token {
            case "+", "-", "*", "/":
            if len(stack) < 2 {
                return 0, errors.New("недостаточно операндов")
            }
            b := stack[len(stack)-1]
            a := stack[len(stack)-2]
            stack = stack[:len(stack)-2]

            switch token {
                case "+":
                stack = append(stack, a+b)
            case "-":
                stack = append(stack, a-b)
            case "*":
                stack = append(stack, a*b)
            case "/":
                if b == 0 {
                    return 0, errors.New("деление на ноль")
                }
                stack = append(stack, a/b)
            }
            default:
                value, err := strconv.ParseFloat(token, 64)
                if err != nil {
                    return 0, errors.New("неверный операнд")
                }
                stack = append(stack, value)
        }
    }

    if len(stack) != 1 {
        return 0, errors.New("неверное выражение")
    }

    return stack[0], nil
}
