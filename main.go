package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Operacion struct {
	Operando1 float64 `json:"operando1"`
	Operando2 float64 `json:"operando2"`
	Operador  string  `json:"operador"`
	Resultado float64 `json:"resultado,omitempty"`
	Error     string  `json:"error,omitempty"`
}

func (op *Operacion) Calcular() {
	switch op.Operador {
	case "+":
		op.Resultado = op.Operando1 + op.Operando2
	case "-":
		op.Resultado = op.Operando1 - op.Operando2
	case "*":
		op.Resultado = op.Operando1 * op.Operando2
	case "/":
		if op.Operando2 == 0 {
			op.Error = "Error: División por cero"
		} else {
			op.Resultado = op.Operando1 / op.Operando2
		}
	default:
		op.Error = "Operador no válido. Usa +, -, * o /"
	}
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
