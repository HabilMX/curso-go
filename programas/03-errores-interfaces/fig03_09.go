// fig03_09.go
// Define un tipo de error propio y recupera sus datos.
package main

import (
    "errors"
    "fmt"
)

// ErrorHTTP es un error que lleva datos adentro.
type ErrorHTTP struct {
    Codigo int
    URL    string
}

// Error hace que ErrorHTTP cumpla la interfaz error.
func (e *ErrorHTTP) Error() string {
    return fmt.Sprintf("el servidor respondio %d", e.Codigo)
}

func consultar(url string) error {
    return &ErrorHTTP{Codigo: 503, URL: url}
}

func Revisar(nombre, url string) error {
    if err := consultar(url); err != nil {
        return fmt.Errorf("revisando %s: %w", nombre, err)
    }
    return nil
}

func main() {
    err := Revisar("catalogo", "https://catalogo.example.com")
    fmt.Println("mensaje:", err)

    var errHTTP *ErrorHTTP
    if errors.As(err, &errHTTP) {
        fmt.Println("es un ErrorHTTP")
        fmt.Println("  codigo:", errHTTP.Codigo)
        fmt.Println("  url:   ", errHTTP.URL)

        if errHTTP.Codigo >= 500 {
            fmt.Println("  → es culpa del servidor, conviene reintentar")
        }
    }
}
