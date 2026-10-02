// fig02_09.go
// Devuelve dos valores: el resultado y un posible error.
package main

import (
    "errors"
    "fmt"
)

// dividir devuelve el cociente y un error si el divisor es cero.
func dividir(a int, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division entre cero")
    }
    return a / b, nil
}

func main() {
    // caso que funciona
    resultado, err := dividir(10, 2)
    if err != nil {
        fmt.Println("error:", err)
    } else {
        fmt.Println("10 / 2 =", resultado)
    }

    // caso que falla
    resultado, err = dividir(10, 0)
    if err != nil {
        fmt.Println("error:", err)
    } else {
        fmt.Println("10 / 0 =", resultado)
    }
}
