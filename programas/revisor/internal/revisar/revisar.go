// Package revisar consulta muchos servicios a la vez. Es el corazón del
// programa.
//
// Lo que este paquete USA lo declara él mismo: la interfaz Revisor vive aquí,
// junto a Todos, que es quien la necesita. La consulta de verdad por HTTP vive
// en el paquete consulta y cumple esta interfaz sin saberlo, y las pruebas usan
// un revisor falso que vive en un archivo _test.go, fuera del programa que se
// compila.
package revisar

import (
	"context"

	"github.com/habil/revisor/internal/servicio"
)

// Revisor es cualquier cosa capaz de consultar un servicio y devolver su estado.
//
// La interfaz es deliberadamente diminuta: un método. Eso es lo que permite que
// Todos —la función que hace el trabajo concurrente— se pueda probar entera con
// un doble, sin levantar un solo servidor.
type Revisor interface {
	Revisar(ctx context.Context, s servicio.Servicio) servicio.Estado
}
