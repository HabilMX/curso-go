// fig02_10.go
// Clasifica servicios por su tiempo de respuesta.
package main

import "fmt"

// clasificar traduce milisegundos a una etiqueta legible.
func clasificar(ms int) string {
    if ms < 0 {
        return "sin medir"
    } else if ms < 500 {
        return "rapido"
    } else if ms < 3000 {
        return "lento"
    }
    return "muy lento"
}

// reporte arma una linea del informe.
func reporte(nombre string, puerto int, ms int) string {
    return fmt.Sprintf("%-10s :%-5d %6dms  %s",
        nombre, puerto, ms, clasificar(ms))
}

func main() {
    fmt.Println("SERVICIO   PUERTO  TIEMPO   ESTADO")
    fmt.Println("-------------------------------------")
    fmt.Println(reporte("catalogo", 443, 142))
    fmt.Println(reporte("pagos", 443, 87))
    fmt.Println(reporte("inventario", 443, 2310))
    fmt.Println(reporte("reportes", 3001, 4500))
    fmt.Println(reporte("viejo", 8080, -1))
}
