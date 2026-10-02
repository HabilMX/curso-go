// fig04_06.go
// append comparte o no memoria segun la capacidad disponible.
package main

import "fmt"

func main() {
    // CASO A: hay capacidad de sobra (cap 5, largo 3)
    a := make([]int, 3, 5)
    a[0], a[1], a[2] = 1, 2, 3
    fmt.Printf("A: len=%d cap=%d %v\n", len(a), cap(a), a)

    b := append(a, 99)          // cabe: NO se muda, comparte memoria
    b[0] = 777
    fmt.Println("   tras modificar b, a vale:", a, "← cambió")

    // CASO B: no hay capacidad (cap 3, largo 3)
    c := make([]int, 3, 3)
    c[0], c[1], c[2] = 1, 2, 3
    fmt.Printf("B: len=%d cap=%d %v\n", len(c), cap(c), c)

    d := append(c, 99)          // NO cabe: se muda a otra memoria
    d[0] = 777
    fmt.Println("   tras modificar d, c vale:", c, "← NO cambió")
}
