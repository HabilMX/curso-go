// fig02_06.go — este programa NO compila
package main

import "fmt"

func main() {
    puerto := 443
    puerto = "muchos"       // intenta guardar texto en un entero
    fmt.Println(puerto)
}
