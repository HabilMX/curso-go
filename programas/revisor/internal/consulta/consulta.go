// Package consulta sabe consultar UN servicio por HTTP de verdad.
//
// No sabe nada de concurrencia ni de reportes, y no conoce la interfaz que
// cumple: revisar.Todos la declara del lado de quien la usa, y main le pasa un
// consulta.HTTP. Así este paquete no depende de revisar, y revisar no depende
// de la red.
package consulta

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/habil/revisor/internal/servicio"
)

// HTTP consulta un servicio de verdad: hace una petición GET y mira qué
// contesta.
//
// Cumple la interfaz revisar.Revisor sin declararlo y sin importar revisar: en
// Go basta con tener el método Revisar (lección 3). La interfaz la declara quien
// la usa, revisar.Todos; este paquete solo devuelve un tipo concreto.
type HTTP struct {
	Cliente *http.Client
}

// NuevoHTTP construye la consulta HTTP con un cliente propio.
//
// 🔴 El cliente es propio y no http.DefaultClient a propósito: el cliente por
// omisión de Go NO tiene tiempo límite. Una petición contra un servicio que
// acepta la conexión y luego no contesta nada se queda colgada para siempre, sin
// error y sin síntoma.
//
// El Timeout del cliente es la red de seguridad de último recurso. El límite que
// de verdad manda es el del context, uno por servicio, que pone revisar.Todos.
func NuevoHTTP(limiteGeneral time.Duration) HTTP {
	return HTTP{Cliente: &http.Client{Timeout: limiteGeneral}}
}

// Revisar consulta un servicio. Nunca devuelve error: el fracaso es parte del
// resultado, porque «este servicio no responde» es justo lo que el programa
// quiere reportar, no una excepción al trabajo.
func (r HTTP) Revisar(ctx context.Context, s servicio.Servicio) servicio.Estado {
	cliente := r.Cliente
	if cliente == nil {
		cliente = &http.Client{Timeout: s.TimeoutEfectivo()}
	}

	inicio := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.URL, nil)
	if err != nil {
		return servicio.Estado{
			Servicio: s,
			Duracion: time.Since(inicio),
			Err:      fmt.Errorf("armando la peticion: %w", sinURL(err)),
		}
	}

	resp, err := cliente.Do(req)
	if err != nil {
		return servicio.Estado{
			Servicio: s,
			Duracion: time.Since(inicio),
			Err:      traducir(ctx, err),
		}
	}
	// Después de comprobar el error, nunca antes: si err != nil, resp es nil y
	// cerrarlo es un pánico.
	defer resp.Body.Close()

	// Se drena el cuerpo aunque no lo queramos leer. Sin esto la conexión no se
	// puede reutilizar y el programa abre una nueva por cada consulta.
	io.Copy(io.Discard, resp.Body)

	return servicio.Estado{
		Servicio: s,
		Codigo:   resp.StatusCode,
		Duracion: time.Since(inicio),
	}
}

// traducir cambia el error del paquete http por uno que se entienda en un
// reporte, y sin la URL (ver sinURL).
func traducir(ctx context.Context, err error) error {
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("se acabo el tiempo de espera")
	}
	if ctx.Err() == context.Canceled {
		return fmt.Errorf("consulta cancelada")
	}
	return fmt.Errorf("no responde: %w", sinURL(err))
}

// sinURL quita la URL de un error de net/http o de net/url.
//
// 🔒 Esos errores se imprimen como «Get "https://host/ruta?token=...": causa»:
// llevan la URL entera, con el usuario y el query string, y el motivo sale
// tal cual en el reporte. Nos quedamos solo con la causa.
func sinURL(err error) error {
	var ue *url.Error
	if errors.As(err, &ue) {
		return ue.Err
	}
	return err
}
