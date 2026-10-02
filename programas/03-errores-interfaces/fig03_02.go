// fig03_02.go
package main

import "fmt"

type Servicio struct {
    Nombre string
    URL    string
    Puerto int
}

func main() {
    var vacio Servicio
    fmt.Printf("%+v\n", vacio)
    fmt.Println("¿el nombre está vacío?", vacio.Nombre == "")
    fmt.Println("puerto:", vacio.Puerto)
}
