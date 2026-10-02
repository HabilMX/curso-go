// Ejemplo de la lección 6, sección 6.5: un deadlock real. El runtime de Go
// lo detecta y mata el programa con "fatal error: all goroutines are
// asleep - deadlock!", en vez de colgarse para siempre.
package main

import "fmt"

func main() {
	canal := make(chan int)
	canal <- 1 // nadie del otro lado está leyendo: se bloquea
	fmt.Println(<-canal)
}
