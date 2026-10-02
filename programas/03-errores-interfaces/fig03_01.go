// fig03_01.go
// Agrupa los datos de un servicio en un solo tipo.
package main

import "fmt"

// Servicio agrupa todo lo que describe a un servicio que vamos a revisar.
type Servicio struct {
    Nombre string
    URL    string
    Puerto int
}

func main() {
    s := Servicio{
        Nombre: "catalogo",
        URL:    "https://catalogo.example.com",
        Puerto: 443,
    }

    fmt.Println("nombre:", s.Nombre)
    fmt.Println("url:", s.URL)
    fmt.Println("puerto:", s.Puerto)

    s.Puerto = 8443          // se puede modificar
    fmt.Println("nuevo puerto:", s.Puerto)

    fmt.Println(s)           // e imprimir completo
}
