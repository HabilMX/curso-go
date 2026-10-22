// Package config traduce un archivo de texto a la lista de servicios que el
// programa va a revisar.
//
// Está separado de la lectura del disco a propósito: Interpretar recibe bytes,
// así que se puede probar con cien entradas distintas sin crear un solo archivo
// temporal.
package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/habil/revisor/internal/servicio"
)

// Cargar lee el archivo de la ruta dada y devuelve los servicios que declara.
func Cargar(ruta string) ([]servicio.Servicio, error) {
	datos, err := os.ReadFile(ruta)
	if err != nil {
		return nil, fmt.Errorf("leyendo la configuracion: %w", err)
	}
	return Interpretar(datos, ruta)
}

// Interpretar convierte el contenido de la configuracion en servicios.
//
// El parametro origen solo se usa para los mensajes de error: es el nombre que
// el usuario va a leer cuando algo esté mal escrito.
//
// Formato, una linea por servicio:
//
//	nombre  url  [tiempo-limite]
//
// Las lineas vacias y las que empiezan con # se ignoran.
func Interpretar(datos []byte, origen string) ([]servicio.Servicio, error) {
	var servicios []servicio.Servicio
	vistos := make(map[string]int) // nombre -> numero de linea donde aparecio

	for i, linea := range strings.Split(string(datos), "\n") {
		numero := i + 1
		linea = strings.TrimSpace(linea)
		if linea == "" || strings.HasPrefix(linea, "#") {
			continue
		}

		campos := strings.Fields(linea)
		if len(campos) < 2 || len(campos) > 3 {
			return nil, fmt.Errorf("%s:%d: esperaba «nombre url [tiempo]», hay %d campo(s)",
				origen, numero, len(campos))
		}

		nombre, direccion := campos[0], campos[1]

		if antes, repetido := vistos[nombre]; repetido {
			return nil, fmt.Errorf("%s:%d: el nombre %q ya estaba en la linea %d",
				origen, numero, nombre, antes)
		}
		if err := validarURL(direccion); err != nil {
			return nil, fmt.Errorf("%s:%d: %w", origen, numero, err)
		}

		timeout := time.Duration(0)
		if len(campos) == 3 {
			d, err := time.ParseDuration(campos[2])
			if err != nil {
				return nil, fmt.Errorf("%s:%d: tiempo limite %q invalido (usa 500ms, 2s, 1m): %w",
					origen, numero, campos[2], err)
			}
			if d <= 0 {
				return nil, fmt.Errorf("%s:%d: el tiempo limite %q tiene que ser mayor que cero",
					origen, numero, campos[2])
			}
			timeout = d
		}

		vistos[nombre] = numero
		servicios = append(servicios, servicio.Servicio{Nombre: nombre, URL: direccion, Timeout: timeout})
	}

	if len(servicios) == 0 {
		return nil, fmt.Errorf("%s: no declara ningun servicio", origen)
	}
	return servicios, nil
}

// validarURL exige http o https y un servidor, y rechaza usuario y clave dentro
// de la URL.
//
// 🔒 Una URL como https://usuario:clave@host lleva un secreto escrito en el
// archivo, y una URL termina en logs, reportes y mensajes de error. Si un
// servicio pide credenciales, van en otro lado (una variable de entorno, un
// gestor de secretos), nunca en la URL. Por eso ningun error de esta funcion
// repite la URL: solo dice que esta mal, para no imprimir el secreto que se
// esta rechazando.
func validarURL(crudo string) error {
	u, err := url.Parse(crudo)
	if err != nil {
		return errors.New("la URL no se puede leer")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("la URL debe empezar con http:// o https://")
	}
	if u.Host == "" {
		return errors.New("a la URL le falta el servidor")
	}
	if u.User != nil {
		return errors.New("la URL trae usuario o clave (usuario:clave@...); las credenciales no van en la URL")
	}
	return nil
}
