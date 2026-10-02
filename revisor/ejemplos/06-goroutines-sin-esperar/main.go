// Ejemplo de la lección 6, sección 6.1: el bug más común del día uno de
// concurrencia en Go. Lanza tres goroutines y no espera a ninguna.
package main

import "fmt"

func main() {
	servicios := []string{"catalogo", "pagos", "inventario"}
	for _, s := range servicios {
		go fmt.Println("revisando", s)
	}
	// el programa termina AQUÍ, sin haber esperado a ninguna goroutine
}
