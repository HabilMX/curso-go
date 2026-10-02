// fig04_09.go
// El recorrido de un map sale en orden distinto cada vez.
package main

import "fmt"

func main() {
    puertos := map[string]int{
        "catalogo":   443,
        "pagos":      443,
        "inventario": 8080,
        "reportes":   3001,
    }

    for nombre, puerto := range puertos {
        fmt.Printf("%s=%d ", nombre, puerto)
    }
    fmt.Println()
}
