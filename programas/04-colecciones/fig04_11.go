// fig04_11.go
// Lee los servicios desde un archivo de texto.
package main

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Servicio struct {
    Nombre string
    URL    string
    Puerto int
}

// Cargar lee el archivo y devuelve los servicios, o un error.
func Cargar(ruta string) ([]Servicio, error) {
    datos, err := os.ReadFile(ruta)
    if err != nil {
        return nil, fmt.Errorf("leyendo %s: %w", ruta, err)
    }

    var servicios []Servicio
    lineas := strings.Split(strings.TrimSpace(string(datos)), "\n")

    for i, linea := range lineas {
        linea = strings.TrimSpace(linea)
        if linea == "" || strings.HasPrefix(linea, "#") {
            continue                       // salta vacias y comentarios
        }

        campos := strings.Fields(linea)     // separa por espacios
        if len(campos) != 3 {
            return nil, fmt.Errorf("%s linea %d: esperaba 3 campos, hay %d",
                ruta, i+1, len(campos))
        }

        puerto, err := strconv.Atoi(campos[2])
        if err != nil {
            return nil, fmt.Errorf("%s linea %d: puerto invalido %q: %w",
                ruta, i+1, campos[2], err)
        }

        servicios = append(servicios, Servicio{
            Nombre: campos[0],
            URL:    campos[1],
            Puerto: puerto,
        })
    }
    return servicios, nil
}

func main() {
    servicios, err := Cargar("servicios.txt")
    if err != nil {
        fmt.Fprintln(os.Stderr, "error:", err)
        os.Exit(1)
    }

    fmt.Printf("cargados %d servicios:\n", len(servicios))
    for _, s := range servicios {
        fmt.Printf("  %-12s %-38s :%d\n", s.Nombre, s.URL, s.Puerto)
    }
}
