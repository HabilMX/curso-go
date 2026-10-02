// fig03_08.go
// Envuelve un error para agregar contexto sin perder el original.
package main

import (
    "errors"
    "fmt"
)

var ErrNoResponde = errors.New("no responde")

type Servicio struct{ Nombre string }

// consultar simula la consulta de bajo nivel.
func consultar(s Servicio) error {
    return ErrNoResponde
}

// Revisar agrega contexto al error de consultar.
func Revisar(s Servicio) error {
    if err := consultar(s); err != nil {
        // el %w ENVUELVE el error original
        return fmt.Errorf("revisando %s: %w", s.Nombre, err)
    }
    return nil
}

func main() {
    err := Revisar(Servicio{Nombre: "catalogo"})

    fmt.Println("1. el mensaje completo:")
    fmt.Println("  ", err)

    fmt.Println("2. ¿es un ErrNoResponde, aunque este envuelto?")
    fmt.Println("  ", errors.Is(err, ErrNoResponde))

    fmt.Println("3. ¿y si pregunto por otro error?")
    fmt.Println("  ", errors.Is(err, errors.New("otra cosa")))

    fmt.Println("4. el error original, desenvuelto:")
    fmt.Println("  ", errors.Unwrap(err))
}
