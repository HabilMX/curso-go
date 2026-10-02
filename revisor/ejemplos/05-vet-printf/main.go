// Ejemplo de la lección 5, sección 5.11: un Printf con el verbo de formato
// equivocado. Compila y corre sin panic — solo produce una salida rota.
// go vet SÍ lo detecta; go build no.
package main

import "fmt"

func main() {
	nombre := "catalogo"
	puerto := 443
	fmt.Printf("servicio %s en el puerto %s\n", nombre, puerto) // %s para un int
}
