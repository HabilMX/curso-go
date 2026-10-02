// fig02_08.go
// Define funciones que reciben y devuelven valores.
package main

import "fmt"

// sumar recibe dos enteros y devuelve su suma.
func sumar(a int, b int) int {
    return a + b
}

// etiqueta arma un texto descriptivo del servicio.
func etiqueta(nombre string, puerto int) string {
    return fmt.Sprintf("%s:%d", nombre, puerto)
}

// esSeguro dice si el puerto corresponde a HTTPS.
func esSeguro(puerto int) bool {
    return puerto == 443
}

func main() {
    fmt.Println("2 + 3 =", sumar(2, 3))
    fmt.Println(etiqueta("catalogo", 443))
    fmt.Println("catalogo es seguro:", esSeguro(443))
    fmt.Println("local es seguro:", esSeguro(8080))
}
