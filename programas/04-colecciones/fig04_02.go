// fig04_02.go
// Crea un slice y lo hace crecer con append.
package main

import "fmt"

type Servicio struct {
    Nombre string
    Puerto int
}

func main() {
    // forma 1: vacio, y crece
    var nombres []string
    fmt.Printf("vacio: %v · largo %d · es nil: %t\n", nombres, len(nombres), nombres == nil)

    nombres = append(nombres, "catalogo")
    nombres = append(nombres, "pagos", "inventario")
    fmt.Printf("con datos: %v · largo %d\n", nombres, len(nombres))

    // forma 2: con valores desde el inicio
    servicios := []Servicio{
        {Nombre: "catalogo", Puerto: 443},
        {Nombre: "pagos", Puerto: 443},
        {Nombre: "local", Puerto: 8080},
    }
    fmt.Println("cuantos servicios:", len(servicios))

    // acceder por posicion, empezando en 0
    fmt.Println("el primero:", servicios[0].Nombre)
    fmt.Println("el ultimo: ", servicios[len(servicios)-1].Nombre)
}
