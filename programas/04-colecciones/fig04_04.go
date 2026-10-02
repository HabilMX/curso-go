// fig04_04.go
// Demuestra que asignar un slice NO copia sus datos.
package main

import "fmt"

func main() {
    original := []int{443, 8080, 3000}
    copia := original              // parece una copia...

    copia[0] = 999                 // modifico la "copia"

    fmt.Println("original:", original)
    fmt.Println("copia:   ", copia)
}
