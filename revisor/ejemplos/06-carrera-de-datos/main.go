// Ejemplo de la lección 6, sección 6.3-6.4: una carrera de datos real,
// provocada a propósito. Corre sin -race y no truena, solo da un número
// incorrecto. Corre con -race y el detector la encuentra.
package main

import (
	"fmt"
	"sync"
)

func main() {
	contador := 0
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			contador++ // dos goroutines pueden leer el mismo valor antes de que la otra escriba
		}()
	}
	wg.Wait()
	fmt.Println("contador:", contador)
}
