// El binario revisor: lee una configuración, consulta todos los servicios que
// declara a la vez, y produce un reporte en tabla o en JSON.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/habil/revisor/internal/config"
	"github.com/habil/revisor/internal/reporte"
	"github.com/habil/revisor/internal/revisar"
)

func main() {
	os.Exit(ejecutar(os.Args[1:], os.Stdout, os.Stderr))
}

// ejecutar hace todo el trabajo real y devuelve el código de salida. Está
// separado de main a propósito: main no se puede probar con go test (llama a
// os.Exit), ejecutar sí.
func ejecutar(args []string, salida, errores *os.File) int {
	banderas := flag.NewFlagSet("revisor", flag.ContinueOnError)
	rutaConfig := banderas.String("config", "", "ruta al archivo de servicios (obligatorio)")
	formato := banderas.String("formato", "tabla", "tabla o json")
	paralelo := banderas.Int("paralelo", 0, "cuántos servicios consultar a la vez (0 = por omisión)")
	limiteGeneral := banderas.Duration("limite", 10*time.Second, "tiempo máximo para todo el reporte")

	if err := banderas.Parse(args); err != nil {
		return 2 // flag ya imprimió el error y el uso
	}
	if *rutaConfig == "" {
		fmt.Fprintln(errores, "revisor: falta --config")
		banderas.Usage()
		return 2
	}
	if *formato != "tabla" && *formato != "json" {
		fmt.Fprintf(errores, "revisor: --formato debe ser \"tabla\" o \"json\", no %q\n", *formato)
		return 2
	}

	servicios, err := config.Cargar(*rutaConfig)
	if err != nil {
		fmt.Fprintln(errores, "revisor:", err)
		return 1
	}

	ctx, cancelar := context.WithTimeout(context.Background(), *limiteGeneral)
	defer cancelar()

	r := revisar.NuevoHTTP(*limiteGeneral)
	estados := revisar.Todos(ctx, r, servicios, *paralelo)

	switch *formato {
	case "json":
		if err := reporte.JSON(salida, estados); err != nil {
			fmt.Fprintln(errores, "revisor: escribiendo el reporte:", err)
			return 1
		}
	default:
		reporte.Tabla(salida, estados)
	}

	for _, e := range estados {
		if !e.OK() {
			return 1 // al menos un servicio falló: el código de salida lo refleja
		}
	}
	return 0
}
