// fig04_07.go
// Operaciones basicas con un map.
package main

import "fmt"

type Estado struct {
    Codigo int
    Ms     int
}

func main() {
    // se crea con make, o vacio con {}
    estados := make(map[string]Estado)

    estados["catalogo"] = Estado{Codigo: 200, Ms: 142}
    estados["pagos"] = Estado{Codigo: 200, Ms: 87}
    estados["inventario"] = Estado{Codigo: 503, Ms: 2310}

    fmt.Println("cuantos:", len(estados))

    // leer una llave que SÍ existe
    e := estados["catalogo"]
    fmt.Println("catalogo:", e.Codigo, e.Ms)

    // 🔑 leer una que NO existe: devuelve el valor cero, SIN error
    fantasma := estados["no-existe"]
    fmt.Printf("no-existe: %+v  ← el valor cero\n", fantasma)

    // la forma correcta de preguntar: el segundo valor
    if e, hay := estados["pagos"]; hay {
        fmt.Println("pagos si esta:", e.Codigo)
    }
    if _, hay := estados["no-existe"]; !hay {
        fmt.Println("no-existe NO esta")
    }

    // borrar
    delete(estados, "inventario")
    fmt.Println("tras borrar:", len(estados))
}
