// fig03_07.go
// Muestra las tres formas de producir un error.
package main

import (
    "errors"
    "fmt"
)

// ErrPuertoInvalido es un error CENTINELA: se declara una vez y se compara.
var ErrPuertoInvalido = errors.New("el puerto debe estar entre 1 y 65535")

type Servicio struct {
    Nombre string
    Puerto int
}

// Validar revisa que el servicio tenga sentido.
func Validar(s Servicio) error {
    if s.Nombre == "" {
        // 1. un error sencillo, creado al momento
        return errors.New("el nombre no puede estar vacio")
    }
    if s.Puerto < 1 || s.Puerto > 65535 {
        // 2. un error centinela, para poder compararlo despues
        return ErrPuertoInvalido
    }
    if s.Puerto == 80 {
        // 3. un error con datos dentro
        return fmt.Errorf("el puerto %d no usa cifrado, usa 443", s.Puerto)
    }
    return nil          // nil significa: todo bien
}

func main() {
    casos := []Servicio{
        {Nombre: "catalogo", Puerto: 443},
        {Nombre: "", Puerto: 443},
        {Nombre: "pagos", Puerto: 99999},
        {Nombre: "viejo", Puerto: 80},
    }

    for _, s := range casos {
        if err := Validar(s); err != nil {
            fmt.Printf("%-10s ❌ %v\n", s.Nombre, err)
        } else {
            fmt.Printf("%-10s ✅ valido\n", s.Nombre)
        }
    }
}
