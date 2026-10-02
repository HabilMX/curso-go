// fig03_04.go
// Muestra por que un receptor de VALOR no puede modificar el original.
package main

import "fmt"

type Servicio struct {
    Nombre string
    Puerto int
}

// CambiarPuerto intenta modificar el servicio... y no lo logra.
func (s Servicio) CambiarPuerto(nuevo int) {
    s.Puerto = nuevo
}

func main() {
    x := Servicio{Nombre: "catalogo", Puerto: 443}
    fmt.Println("antes:", x.Puerto)

    x.CambiarPuerto(8443)
    fmt.Println("despues:", x.Puerto)      // ¿cambió?
}

