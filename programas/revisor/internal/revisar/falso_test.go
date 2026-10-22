package revisar

import (
	"context"
	"sync"
	"time"

	"github.com/habil/revisor/internal/servicio"
)

// revisorFalso es un Revisor que nunca toca la red. Existe para poder probar
// Todos —la función concurrente— con resultados exactos y predecibles, en vez
// de depender de que un servidor de verdad esté arriba y conteste rápido.
//
// Vive en un archivo _test.go a propósito: go build no lo compila, así que no
// forma parte del programa ni de lo que el paquete exporta. Es un doble de
// prueba, y un doble de prueba solo existe para las pruebas.
//
// Se configura con respuestas: un mapa de nombre de servicio -> lo que el
// falso debe "contestar" cuando lo consulten, incluyendo cuánto tardar en
// hacerlo. Si un servicio no está en el mapa, contesta 200 de inmediato.
//
// 🔒 Lleva su propio mutex porque las pruebas lo consultan desde varias
// goroutines a la vez (eso es justo lo que Todos hace) y contador se
// incrementa en cada llamada: sin el candado, el detector de carreras (`go
// test -race`) lo señala de inmediato. Ésa es la prueba de que hace falta,
// no una suposición.
type revisorFalso struct {
	respuestas map[string]respuestaFalsa

	mu       sync.Mutex
	contador int // cuántas veces se llamó a Revisar
}

// respuestaFalsa es lo que revisorFalso va a devolver para un servicio en
// particular.
type respuestaFalsa struct {
	codigo  int
	espera  time.Duration // simula latencia de red
	err     error
	colgado bool // si es true, Revisar no responde hasta que el context lo cancele
}

// Revisar cumple la interfaz Revisor sin tocar la red.
func (f *revisorFalso) Revisar(ctx context.Context, s servicio.Servicio) servicio.Estado {
	f.mu.Lock()
	f.contador++
	f.mu.Unlock()

	r, existe := f.respuestas[s.Nombre]
	if !existe {
		return servicio.Estado{Servicio: s, Codigo: 200}
	}

	if r.colgado {
		// Simula un servicio que acepta la conexión y nunca contesta: la única
		// forma de salir de aquí es que el context de arriba se cancele.
		<-ctx.Done()
		return servicio.Estado{Servicio: s, Err: ctx.Err()}
	}

	if r.espera > 0 {
		temporizador := time.NewTimer(r.espera)
		defer temporizador.Stop()
		select {
		case <-temporizador.C:
			// pasó el tiempo simulado; sigue abajo con la respuesta normal
		case <-ctx.Done():
			return servicio.Estado{Servicio: s, Err: ctx.Err()}
		}
	}

	if r.err != nil {
		return servicio.Estado{Servicio: s, Err: r.err}
	}
	return servicio.Estado{Servicio: s, Codigo: r.codigo}
}
