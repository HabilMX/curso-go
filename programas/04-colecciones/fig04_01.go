// fig04_01.go
// Muestra arreglos, que en Go casi no se usan directamente.
package main

import "fmt"

func main() {
    var puertos [3]int              // tres enteros, todos en 0
    fmt.Println("recien creado:", puertos)

    puertos[0] = 443
    puertos[1] = 8080
    fmt.Println("con valores:", puertos)
    fmt.Println("cuantos caben:", len(puertos))

    otros := [3]int{80, 443, 8443}
    fmt.Println("otro arreglo:", otros)

    copia := otros                  // los arreglos SÍ se copian completos
    copia[0] = 999
    fmt.Println("original:", otros)
    fmt.Println("copia:   ", copia)
}
