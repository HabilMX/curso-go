package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestEjecutar_reportaOKyFalla(t *testing.T) {
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/mal") {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer servidor.Close()

	archivoConfig, err := os.CreateTemp(t.TempDir(), "servicios-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	contenido := "sano " + servidor.URL + "/bien\nmalo " + servidor.URL + "/mal\n"
	if _, err := archivoConfig.WriteString(contenido); err != nil {
		t.Fatal(err)
	}
	archivoConfig.Close()

	var salida, errores bytes.Buffer
	codigo := correr(t, []string{"--config", archivoConfig.Name()}, &salida)

	if codigo != 1 {
		t.Errorf("código de salida = %d, quería 1 (hay un servicio con falla)", codigo)
	}
	if !strings.Contains(salida.String(), "sano") || !strings.Contains(salida.String(), "malo") {
		t.Errorf("la tabla no menciona a los dos servicios:\n%s", salida.String())
	}
	_ = errores
}

func TestEjecutar_faltaConfig(t *testing.T) {
	var salida bytes.Buffer
	codigo := correr(t, []string{}, &salida)
	if codigo != 2 {
		t.Errorf("código de salida = %d, quería 2 (falta --config)", codigo)
	}
}

// correr es un envoltorio de ejecutar() que redirige stdout a un archivo
// temporal para poder capturarlo, porque ejecutar() escribe en *os.File, no
// en io.Writer (una limitación real de esta versión del programa: se explica
// en la lección).
func correr(t *testing.T, args []string, capturaEnBuffer *bytes.Buffer) int {
	t.Helper()
	tmp, err := os.CreateTemp(t.TempDir(), "salida-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer tmp.Close()

	codigo := ejecutar(args, tmp, os.Stderr)

	tmp.Seek(0, 0)
	datos := make([]byte, 1<<16)
	n, _ := tmp.Read(datos)
	capturaEnBuffer.Write(datos[:n])
	return codigo
}
