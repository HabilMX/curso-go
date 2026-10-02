// fig04_03.go — este programa COMPILA pero truena
package main

import "fmt"

func main() {
    puertos := []int{443, 8080}
    fmt.Println(puertos[0])
    fmt.Println(puertos[5])      // solo hay 2 elementos
}
