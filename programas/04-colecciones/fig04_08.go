// fig04_08.go — este programa truena al escribir
package main

import "fmt"

func main() {
    var m map[string]int          // nil: NO inicializado

    fmt.Println("leer de un map nil:", m["x"])    // esto funciona
    fmt.Println("su largo:", len(m))              // esto tambien

    m["x"] = 1                                     // esto truena
}
