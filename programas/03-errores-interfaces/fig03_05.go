// fig03_05.go
// Con receptor de PUNTERO, el metodo si modifica el original.
package main

import "fmt"

type Servicio struct {
    Nombre string
    Puerto int
}

// CambiarPuerto modifica el servicio. El * es la diferencia.
func (s *Servicio) CambiarPuerto(nuevo int) {
    s.Puerto = nuevo
}

func main() {
    x := Servicio{Nombre: "catalogo", Puerto: 443}
    fmt.Println("antes:", x.Puerto)

    x.CambiarPuerto(8443)
    fmt.Println("despues:", x.Puerto)
}
