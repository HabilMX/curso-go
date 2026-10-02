// fig02_03.go
// Declara variables de distintos tipos y las imprime.
package main

import "fmt"

func main() {
    nombre := "catalogo"                 // texto
    puerto := 443                      // número entero
    tiempo := 0.142                    // número con decimales
    seguro := true                     // verdadero o falso

    fmt.Println("servicio:", nombre)
    fmt.Println("puerto:", puerto)
    fmt.Println("tiempo de respuesta:", tiempo)
    fmt.Println("usa HTTPS:", seguro)
}
