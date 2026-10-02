// fig03_10.go
// Una interfaz con dos implementaciones: la real y una de prueba.
package main

import (
    "errors"
    "fmt"
)

type Servicio struct {
    Nombre string
    URL    string
}

type Estado struct {
    Servicio Servicio
    Codigo   int
    Err      error
}

// Revisor es una INTERFAZ: cualquier cosa con este metodo la cumple.
type Revisor interface {
    Revisar(s Servicio) Estado
}

// ---- primera implementacion: la de verdad (simplificada) ----
type RevisorHTTP struct{}

func (r RevisorHTTP) Revisar(s Servicio) Estado {
    // En la leccion 7 esto hara una peticion HTTP real.
    return Estado{Servicio: s, Codigo: 200}
}

// ---- segunda implementacion: para probar, sin red ----
type RevisorFalso struct {
    Respuesta Estado
}

func (r RevisorFalso) Revisar(s Servicio) Estado {
    r.Respuesta.Servicio = s
    return r.Respuesta
}

// RevisarTodos acepta CUALQUIER Revisor. No sabe ni le importa cual.
func RevisarTodos(r Revisor, servicios []Servicio) []Estado {
    var estados []Estado
    for _, s := range servicios {
        estados = append(estados, r.Revisar(s))
    }
    return estados
}

func main() {
    servicios := []Servicio{
        {Nombre: "catalogo", URL: "https://catalogo.example.com"},
        {Nombre: "pagos", URL: "https://pagos.example.com"},
    }

    fmt.Println("--- con el revisor real ---")
    for _, e := range RevisarTodos(RevisorHTTP{}, servicios) {
        fmt.Printf("  %-10s codigo %d\n", e.Servicio.Nombre, e.Codigo)
    }

    fmt.Println("--- con el falso, simulando una falla ---")
    falso := RevisorFalso{
        Respuesta: Estado{Err: errors.New("no responde")},
    }
    for _, e := range RevisarTodos(falso, servicios) {
        fmt.Printf("  %-10s error: %v\n", e.Servicio.Nombre, e.Err)
    }
}
