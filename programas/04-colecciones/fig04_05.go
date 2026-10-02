// fig04_05.go
// Las dos formas de copiar un slice de verdad.
package main

import "fmt"

func main() {
    original := []int{443, 8080, 3000}

    // forma 1: make + copy
    copia1 := make([]int, len(original))
    copy(copia1, original)
    copia1[0] = 111

    // forma 2: append sobre un slice nil (mas corta)
    copia2 := append([]int(nil), original...)
    copia2[1] = 222

    fmt.Println("original:", original)
    fmt.Println("copia1:  ", copia1)
    fmt.Println("copia2:  ", copia2)
}
