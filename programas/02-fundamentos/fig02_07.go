// fig02_07.go
// Muestra el valor cero de cada tipo basico.
package main

import "fmt"

func main() {
    var texto string
    var entero int
    var decimal float64
    var logico bool

    fmt.Printf("string:  %q\n", texto)
    fmt.Printf("int:     %d\n", entero)
    fmt.Printf("float64: %v\n", decimal)
    fmt.Printf("bool:    %t\n", logico)
}
