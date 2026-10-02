// fig03_06.go — este programa COMPILA pero truena al correr
package main

import "fmt"

type Servicio struct{ Puerto int }

func main() {
    var p *Servicio          // un puntero sin apuntar a nada: vale nil
    fmt.Println(p)           // esto sí funciona
    fmt.Println(p.Puerto)    // esto truena
}
