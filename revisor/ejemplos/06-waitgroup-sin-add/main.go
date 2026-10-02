// Ejemplo de la lección 6, sección 6.2: qué pasa si te saltas wg.Add(1)
// antes de lanzar la goroutine. El programa imprime "listo" ANTES del
// panic, las veces que se corra — ver la lección para la explicación.
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		go func(n int) { // <- sin wg.Add(1) antes de esta línea
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
		}(i)
	}
	wg.Wait()
	fmt.Println("listo")
}
