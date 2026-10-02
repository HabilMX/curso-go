// fig04_10.go
// Ordena las llaves para producir una salida estable.
package main

import (
    "fmt"
    "sort"
)

type Estado struct {
    Codigo int
    Ms     int
}

func main() {
    estados := map[string]Estado{
        "reportes":   {Codigo: 0, Ms: 0},
        "catalogo":   {Codigo: 200, Ms: 142},
        "inventario": {Codigo: 503, Ms: 2310},
        "pagos":      {Codigo: 200, Ms: 87},
    }

    // 1. saca las llaves a un slice
    nombres := make([]string, 0, len(estados))
    for nombre := range estados {
        nombres = append(nombres, nombre)
    }

    // 2. ordenalas
    sort.Strings(nombres)

    // 3. recorre el slice ordenado, no el map
    fmt.Println("SERVICIO     CODIGO  TIEMPO")
    fmt.Println("-----------------------------")
    for _, nombre := range nombres {
        e := estados[nombre]
        fmt.Printf("%-12s %6d  %5dms\n", nombre, e.Codigo, e.Ms)
    }
}
