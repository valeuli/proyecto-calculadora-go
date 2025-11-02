package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	_ "strconv"
	"strings"
	_ "strings"
	"unicode"
	_ "unicode"
)

// Estructura de datos
type Operacion struct {
	Expresion string  `json:"expresion"`
	Resultado float64 `json:"resultado,omitempty"`
	Error     string  `json:"error,omitempty"`
}

func (op *Operacion) Calcular() {
	postfija, err := procesarExpresion(op.Expresion)
	if err != nil {
		op.Error = fmt.Sprintf("Error en la expresión: %v", err)
		return
	}

	resultado, err := evaluarPostfix(postfija)
	if err != nil {
		op.Error = fmt.Sprintf("Error al calcular: %v", err)
		return
	}

	op.Resultado = resultado
}

func procesarExpresion(expr string) ([]string, error) {
	var output []string
	var stack []rune
	var number strings.Builder

	precedencia := map[rune]int{
		'+': 1,
		'-': 1,
		'*': 2,
		'/': 2,
	}

	expr = strings.ReplaceAll(expr, " ", "")
	for i, ch := range expr {
		if unicode.IsDigit(ch) || ch == '.' {
			number.WriteRune(ch)
			if i == len(expr)-1 || (!unicode.IsDigit(rune(expr[i+1])) && expr[i+1] != '.') {
				output = append(output, number.String())
				number.Reset()
			}
		} else if ch == '(' {
			stack = append(stack, ch)
		} else if ch == ')' {
			for len(stack) > 0 && stack[len(stack)-1] != '(' {
				output = append(output, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, fmt.Errorf("paréntesis no balanceados")
			}
			stack = stack[:len(stack)-1] // quitar '('
		} else if prec, esOperador := precedencia[ch]; esOperador {
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if top == '(' || precedencia[top] < prec {
					break
				}
				output = append(output, string(top))
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, ch)
		} else {
			return nil, fmt.Errorf("carácter inválido: %c", ch)
		}
	}

	// Validar la cantidad de paréntesis
	for len(stack) > 0 {
		if stack[len(stack)-1] == '(' {
			return nil, fmt.Errorf("paréntesis no balanceados")
		}
		output = append(output, string(stack[len(stack)-1]))
		stack = stack[:len(stack)-1]
	}

	return output, nil
}

func evaluarPostfix(postfix []string) (float64, error) {
	var stack []float64

	for _, token := range postfix {
		switch token {
		case "+", "-", "*", "/":
			if len(stack) < 2 {
				return 0, fmt.Errorf("expresión mal formada")
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
					return 0, fmt.Errorf("división por cero")
				}
				stack = append(stack, a/b)
			}
		default:
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, fmt.Errorf("token inválido: %s", token)
			}
			stack = append(stack, num)
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("expresión mal formada")
	}
	return stack[0], nil
}

func manejarOperacion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var op Operacion
	err := json.NewDecoder(r.Body).Decode(&op)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo", http.StatusBadRequest)
		return
	}

	op.Calcular()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(op)
}

func main() {
	http.HandleFunc("/operacion", manejarOperacion)
	fmt.Println("Servidor en ejecución en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
