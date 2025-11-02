# Calculadora distribuida en Go

Este proyecto es una calculadora aritmética que recibe expresiones completas como `(3 + 4) * 2` mediante una solicitud HTTP POST. Evalúa la expresión y responde con el resultado en formato JSON.

## Cómo ejecutar

1. Asegúrate de tener Go instalado.
2. Ejecuta el servidor con:

```bash
go run main.go
```

El servidor escuchará en:
http://localhost:8080/operacion

### Ejemplo de solicitud
```bash
{
"expresion": "(3 + 4) * 2"
}
```
### Ejemplo de respuesta
```bash
{
  "expresion": "(3 + 4) * 2",
  "resultado": 14
}
```