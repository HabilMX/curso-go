// fig03_03.go
// Define metodos que pertenecen al tipo Servicio.
package main

import "fmt"

type Servicio struct {
    Nombre string
    URL    string
    Puerto int
}

// Etiqueta devuelve una descripcion legible del servicio.
func (s Servicio) Etiqueta() string {
    return fmt.Sprintf("%s (%s:%d)", s.Nombre, s.URL, s.Puerto)
}

// EsSeguro indica si el servicio usa el puerto de HTTPS.
func (s Servicio) EsSeguro() bool {
    return s.Puerto == 443
}

func main() {
    a := Servicio{Nombre: "catalogo", URL: "https://catalogo.example.com", Puerto: 443}
    b := Servicio{Nombre: "local", URL: "http://localhost", Puerto: 8080}

    fmt.Println(a.Etiqueta(), "— seguro:", a.EsSeguro())
    fmt.Println(b.Etiqueta(), "— seguro:", b.EsSeguro())
}
