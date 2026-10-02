# Lección 7 — El programa terminado

**Duración:** 90-120 minutos, o 2 sesiones de 60. Es la última lección del `revisor`: al terminar tienes
un binario que de verdad corre, se configura desde la línea de comandos, habla HTTP con timeouts reales,
produce dos formatos de salida, y se compila para otra plataforma sin salir de tu máquina.

**Al terminar vas a poder:**

- Escribir un cliente HTTP con timeout explícito, y explicar por qué `http.DefaultClient` es peligroso en
  producción.
- Separar la lógica de un programa (`ejecutar`) de su punto de entrada (`main`), para poder probarla sin
  tocar `os.Exit`.
- Definir banderas de línea de comandos con el paquete estándar `flag`, sin ninguna librería externa.
- Serializar una estructura a JSON con `encoding/json`, y explicar qué campos se pierden y por qué.
- Correr el `revisor` de principio a fin contra un servidor real (uno de prueba, hecho por ti) y leer su
  salida.
- Compilar el mismo código para otra arquitectura y otro sistema operativo sin salir de tu máquina, y
  explicar por qué eso es posible.

---

## Por qué importa

Ya tienes las cuatro piezas del `revisor` construidas por separado: `servicio` define el vocabulario
(lección 2-3), `config` lee la configuración (lección 5), `revisar.Todos` consulta todo a la vez (lección
6). Lo único que falta es lo que convierte esas piezas en **un programa que alguien más pueda usar sin
leer el código fuente**: que hable HTTP de verdad (hasta ahora solo lo probamos con `Falso`), que reciba
banderas desde la terminal en vez de valores fijos en el código, que produzca un formato que otro
programa pueda consumir, y que se pueda compilar y repartir como un solo archivo.

Ésta es la lección donde el `revisor` deja de ser "código que funciona si lo corro yo, en mi máquina, con
mis datos de prueba" y se vuelve un binario que puedes copiarle a alguien más, con la confianza de que va
a hacer exactamente lo que la línea de comandos le pida.

---

## Los conceptos

### 7.1 El cliente HTTP: por qué nunca `http.DefaultClient`

Así quedó el `Revisor` de verdad, el que sí toca la red (`revisor/internal/revisar/revisar.go`):

<!-- verificar:extracto:internal/revisar/revisar.go -->
```go
// HTTP es el Revisor de verdad: hace una petición GET y mira qué contesta.
type HTTP struct {
	Cliente *http.Client
}

// NuevoHTTP construye el Revisor con un cliente propio.
//
// 🔴 El cliente es propio y no http.DefaultClient a propósito: el cliente por
// omisión de Go NO tiene tiempo límite. Una petición contra un servicio que
// acepta la conexión y luego no contesta nada se queda colgada para siempre, sin
// error y sin síntoma.
//
// El Timeout del cliente es la red de seguridad de último recurso. El límite que
// de verdad manda es el del context, uno por servicio, que pone Todos.
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
			Err:      fmt.Errorf("armando la peticion: %w", err),
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
```

🔴 **`http.DefaultClient` —el que usarías si escribieras `http.Get(url)` directo— no tiene ningún
timeout.** Si el servidor del otro lado acepta la conexión y nunca contesta nada, esa llamada se queda
esperando **para siempre**, sin error, sin panic, solo un programa que un día deja de avanzar y nadie
sabe por qué. Es uno de los errores más caros de Go en producción precisamente porque no da ningún
síntoma hasta que ya está pasando.

El `revisor` se protege dos veces, no una: el `http.Client` propio tiene un timeout general (el límite de
todo el reporte, pasado a `NuevoHTTP`), y además cada consulta individual corre bajo el
`context.WithTimeout` por servicio que armó `Todos` en la lección 6 — ese segundo, más corto, es el que
manda en la práctica. El timeout del cliente es la red de seguridad de último recurso, no el mecanismo
principal.

⚠️ **`defer resp.Body.Close()` va después de comprobar el error, nunca antes.** Si `err != nil`, `resp`
es `nil`, y llamar a un método sobre un puntero nulo hace panic. Y **`io.Copy(io.Discard, resp.Body)`
antes de cerrar** no es decorativo: sin drenar el cuerpo de la respuesta, la conexión TCP subyacente no
se puede reutilizar para la siguiente petición al mismo servidor, y el programa termina abriendo una
conexión nueva cada vez en vez de reusar las que ya tiene.

**`traducir` cambia el error crudo de `net/http` por uno legible en una tabla:**

<!-- verificar:extracto:internal/revisar/revisar.go -->
```go
func traducir(ctx context.Context, err error) error {
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("se acabo el tiempo de espera")
	}
	if ctx.Err() == context.Canceled {
		return fmt.Errorf("consulta cancelada")
	}
	return fmt.Errorf("no responde: %w", err)
}
```

El error que da `net/http` de fábrica trae la URL completa y la palabra `Get`, que en una tabla de
reporte solo ocupan espacio y no le dicen nada nuevo al lector — por eso se traduce antes de mostrarse.

### 7.2 `main` no se prueba; `ejecutar` sí

El patrón que separa el `revisor` de un programa de un solo archivo:

<!-- verificar:extracto:cmd/revisor/main.go -->
```go
func main() {
	os.Exit(ejecutar(os.Args[1:], os.Stdout, os.Stderr))
}
```

`ejecutar` es donde vive la lógica real —las secciones 7.3, 7.4 y 7.6 la van mostrando por partes—; su
firma:

<!-- verificar:fragmento -->
```go
func ejecutar(args []string, salida, errores *os.File) int {
	// toda la lógica real vive aquí
}
```

**`go test` no puede probar una función que llama a `os.Exit`**, porque `os.Exit` termina el proceso
entero de inmediato — incluido el propio proceso de pruebas, que jamás llegaría a reportar el resultado.
Separando "decidir qué código de salida corresponde" (`ejecutar`, que **devuelve** un `int`) de "salir de
verdad con ese código" (`main`, la única línea que llama a `os.Exit`), toda la lógica se puede probar
llamando a `ejecutar` directo, con argumentos de prueba, y revisando el número que devuelve — exactamente
lo que hace `cmd/revisor/main_test.go`:

<!-- verificar:extracto:cmd/revisor/main_test.go -->
```go
func TestEjecutar_faltaConfig(t *testing.T) {
	var salida bytes.Buffer
	codigo := correr(t, []string{}, &salida)
	if codigo != 2 {
		t.Errorf("código de salida = %d, quería 2 (falta --config)", codigo)
	}
}
```

```
$ go test ./cmd/revisor/... -v -run TestEjecutar_faltaConfig
=== RUN   TestEjecutar_faltaConfig
revisor: falta --config
Usage of revisor:
  -config string
    	ruta al archivo de servicios (obligatorio)
  -formato string
    	tabla o json (default "tabla")
  -limite duration
    	tiempo máximo para todo el reporte (default 10s)
  -paralelo int
    	cuántos servicios consultar a la vez (0 = por omisión)
--- PASS: TestEjecutar_faltaConfig (0.00s)
```

Esa es la salida real de `--help` (generada automáticamente por `flag`, sección 7.3) capturada en medio
de una prueba, sin abrir una terminal ni ejecutar el binario compilado.

### 7.3 Banderas de línea de comandos, sin librerías

<!-- verificar:extracto:cmd/revisor/main.go -->
```go
banderas := flag.NewFlagSet("revisor", flag.ContinueOnError)
rutaConfig := banderas.String("config", "", "ruta al archivo de servicios (obligatorio)")
formato := banderas.String("formato", "tabla", "tabla o json")
paralelo := banderas.Int("paralelo", 0, "cuántos servicios consultar a la vez (0 = por omisión)")
limiteGeneral := banderas.Duration("limite", 10*time.Second, "tiempo máximo para todo el reporte")

if err := banderas.Parse(args); err != nil {
	return 2 // flag ya imprimió el error y el uso
}
```

**`flag.NewFlagSet` en vez del paquete `flag` a nivel global** (que usarías con `flag.String(...)`
directo) es lo que permite tener una función `ejecutar(args []string, ...)` que recibe sus argumentos
como parámetro, en vez de leerlos siempre de `os.Args` — otro detalle que existe específicamente para
poder probarla, como en la sección 7.2.

`flag` viene en la biblioteca estándar y alcanza para un programa de este tamaño: cuatro banderas, tipos
básicos (`string`, `int`, `time.Duration`), ayuda automática. Si algún día el `revisor` necesitara
subcomandos (`revisor check`, `revisor list`, cada uno con sus propias banderas), ahí sí valdría la pena
mirar una librería como `cobra` — pero no antes de necesitarlo de verdad.

### 7.4 Los dos formatos de salida: tabla y JSON

Así quedó `Tabla`, la función que produce la salida legible por humanos que has visto en toda la lección
(`revisor/internal/reporte/tabla.go`):

<!-- verificar:extracto:internal/reporte/tabla.go -->
```go
func Tabla(w io.Writer, estados []servicio.Estado) {
	ordenados := ordenarPorNombre(estados)

	anchoNombre := len("SERVICIO")
	for _, e := range ordenados {
		if n := len(e.Servicio.Etiqueta()); n > anchoNombre {
			anchoNombre = n
		}
	}

	fmt.Fprintf(w, "%-*s  %-6s  %8s  %s\n", anchoNombre, "SERVICIO", "ESTADO", "TIEMPO", "DETALLE")
	for _, e := range ordenados {
		fmt.Fprintf(w, "%-*s  %-6s  %8s  %s\n",
			anchoNombre,
			e.Servicio.Etiqueta(),
			etiquetaEstado(e),
			formatoTiempo(e),
			detalle(e),
		)
	}
}
```

**El ancho de la primera columna no está fijo en el código — se calcula.** El bucle de arriba recorre
todos los estados una vez antes de imprimir nada, y se queda con el nombre más largo (empezando desde el
ancho de la propia palabra "SERVICIO", por si algún nombre de servicio fuera más corto que eso). Ese
número entra como el `*` en `%-*s`: un verbo de formato de ancho **variable**, donde el ancho mismo es un
argumento más de `Fprintf`, no un número escrito a mano. Es la razón por la que la tabla real de esta
lección tiene columnas rectas sin importar si los nombres de los servicios son cortos (`sano`) o largos
(`inventario`): con un ancho fijo, uno de los dos casos siempre se vería torcido.

`ordenarPorNombre` (sección 6.3 y 6.7 ya explicaron por qué el orden de llegada no es de fiar: varias
goroutines entregan resultados por un canal, y gana quien conteste primero) hace una copia del slice y la
ordena alfabéticamente antes de imprimir — sin esto, la misma consulta produciría una tabla en un orden
distinto cada vez que corrieras el programa, aunque los datos fueran idénticos.

`etiquetaEstado`, `formatoTiempo` y `detalle` son las tres funciones chicas que deciden qué palabra va en
cada columna (`OK`/`LENTO`/`FALLA`, el tiempo redondeado al milisegundo o un guion si nunca hubo
respuesta, y el código HTTP o el motivo del fallo) — son las que ya viste actuando en la tabla real de la
sección 7.5, ahora con el código que las produce.

🔑 **Por qué no se usó `text/tabwriter`, que es la herramienta que la biblioteca estándar ofrece
justamente para alinear columnas:** para una tabla de cuatro columnas donde solo una tiene ancho variable
(el nombre del servicio; las otras tres son cortas y predecibles), calcular el ancho a mano es más simple
de leer que introducir un `tabwriter.Writer` con sus propios `Flush()` y separadores por tabulador. Con
una tabla de más columnas variables, `tabwriter` sí sería la herramienta correcta — vale la pena conocerlo
(está en la sección "Para leer más" de esta lección) aunque el `revisor` no lo necesite.

Con la tabla ya entendida, el otro formato de salida — JSON — es la otra mitad de esta sección:

<!-- verificar:extracto:internal/reporte/json.go -->
```go
type lineaJSON struct {
	Servicio string `json:"servicio"`
	OK       bool   `json:"ok"`
	Codigo   int    `json:"codigo,omitempty"`
	TiempoMs int64  `json:"tiempo_ms"`
	Detalle  string `json:"detalle,omitempty"`
}
```

⚠️ **`servicio.Estado` no se serializa directo — a propósito.** `Estado` trae un campo `Err error`, y
`error` es una interfaz: `encoding/json` no sabe cómo convertirla a texto por sí sola (intentarlo produce
un objeto vacío `{}`, no un error de compilación, que es peor: falla en silencio). El `revisor` resuelve
esto con un tipo intermedio, `lineaJSON`, que solo vive dentro del paquete `reporte`: convierte el error
a su texto (`e.Motivo()`) antes de codificar, y de paso desacopla el formato público (lo que otros
programas van a leer) del modelo interno (`Estado`) — si mañana `Estado` gana un campo nuevo, el JSON que
ya circula no cambia solo porque cambió algo interno.

`omitempty` en `Codigo` y `Detalle` quita esos campos del JSON cuando valen cero o cadena vacía — así un
servicio sano no arrastra un `"detalle": ""` sin sentido. Codificado de verdad:

<!-- verificar:extracto:internal/reporte/json.go -->
```go
codificador := json.NewEncoder(w)
codificador.SetIndent("", "  ")
return codificador.Encode(lineas)
```

### 7.5 De extremo a extremo: el `revisor` corriendo contra un servidor real

Para probar todo el programa junto —no cada pieza por separado— construí `cmd/servidor-demo`: un
servidor HTTP mínimo con cuatro rutas que se comportan como los cuatro casos que el `revisor` tiene que
saber reportar.

<!-- verificar:extracto:cmd/servidor-demo/main.go -->
```go
mux.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
})
mux.HandleFunc("/lento", func(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1500 * time.Millisecond)
	w.WriteHeader(http.StatusOK)
})
mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
})
mux.HandleFunc("/nunca-contesta", func(w http.ResponseWriter, r *http.Request) {
	<-r.Context().Done() // no responde nunca por su cuenta; solo cede cuando el cliente se cansa
})
```

Con ese servidor levantado en `:8091` y un archivo de configuración apuntando a sus cuatro rutas, corrí
el binario real:

```
$ /tmp/servidor-demo -puerto 8091 &
servidor-demo escuchando en :8091 (/ok, /lento, /error, /nunca-contesta)

$ cat /tmp/servicios-demo.txt
sano      http://127.0.0.1:8091/ok
lento     http://127.0.0.1:8091/lento
malo      http://127.0.0.1:8091/error
colgado   http://127.0.0.1:8091/nunca-contesta   200ms

$ /tmp/revisor --config /tmp/servicios-demo.txt --formato tabla
SERVICIO  ESTADO    TIEMPO  DETALLE
colgado   FALLA      202ms  se acabo el tiempo de espera
lento     LENTO     1.503s  200
malo      FALLA        1ms  codigo 500
sano      OK           1ms  200
$ echo $?
1
```

**Cada línea de esa tabla corresponde exactamente al comportamiento que se programó en el servidor de
prueba:** `sano` contesta rápido y con 200; `lento` tarda 1.5 segundos de verdad y por eso sale `LENTO`;
`malo` contesta al instante pero con 500; `colgado` nunca contesta, y su timeout individual de 200ms
—declarado en el archivo de configuración, columna tres— cortó la espera en 202ms, casi exacto. **El
código de salida, 1, es real**: al menos un servicio no estaba `OK`, así que `ejecutar` devuelve 1 en vez
de 0 (sección 7.6) — el mismo binario, usado desde un script, le puede decir a quien lo invoque si hubo
problemas sin que nadie tenga que leer la tabla.

Y en JSON, el mismo reporte:

```
$ /tmp/revisor --config /tmp/servicios-demo.txt --formato json
[
  {"servicio": "colgado", "ok": false, "tiempo_ms": 205, "detalle": "se acabo el tiempo de espera"},
  {"servicio": "lento", "ok": true, "codigo": 200, "tiempo_ms": 1503},
  {"servicio": "malo", "ok": false, "codigo": 500, "tiempo_ms": 1, "detalle": "codigo 500"},
  {"servicio": "sano", "ok": true, "codigo": 200, "tiempo_ms": 1}
]
```

(Aquí sin la indentación de `SetIndent` para que quepa en una línea por servicio; el programa real la
produce con saltos de línea, como viste en la sección 7.4.)

### 7.6 El código de salida, con significado

<!-- verificar:extracto:cmd/revisor/main.go -->
```go
for _, e := range estados {
	if !e.OK() {
		return 1 // al menos un servicio falló: el código de salida lo refleja
	}
}
return 0
```

Tres códigos posibles, y cada uno responde una pregunta distinta a quien invoque el programa desde un
script o un pipeline: **0** — todo salió bien; **1** — el programa corrió completo, pero al menos un
servicio no estaba sano; **2** — el programa ni siquiera pudo arrancar (falta `--config`, el archivo no
existe, el formato de configuración está mal escrito). Es la diferencia entre "hice el trabajo y encontré
problemas" y "no pude ni empezar a trabajar" — y un pipeline de integración continua puede reaccionar
distinto a cada uno.

### 7.7 Compilación cruzada: el mismo código, otra plataforma

```bash
GOOS=linux GOARCH=amd64 go build -o revisor-linux-amd64 ./cmd/revisor
```

Confirmado de verdad, compilando desde esta Mac (arm64) hacia Linux x86-64:

```
$ GOOS=linux GOARCH=amd64 go build -o revisor-linux-amd64 ./cmd/revisor
$ file revisor-linux-amd64
revisor-linux-amd64: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, ...
```

**Sin Docker, sin máquina virtual, sin instalar nada adicional: dos variables de entorno bastan.** Esto es
posible porque el compilador de Go trae, de fábrica, el código necesario para generar binarios de
cualquier combinación de sistema operativo y arquitectura que soporte — no depende de estar corriendo en
el sistema destino para compilar hacia él, al revés de cómo funcionan muchos otros lenguajes compilados.

**El tamaño del binario, medido, con y sin símbolos de depuración:**

```
$ go build -o revisor-normal ./cmd/revisor
$ wc -c revisor-normal
9464130 revisor-normal

$ go build -ldflags="-s -w" -o revisor-chico ./cmd/revisor
$ wc -c revisor-chico
6386194 revisor-chico
```

**De 9.46 MB a 6.39 MB, una reducción real del 32%**, quitando la tabla de símbolos (`-s`) y la
información de depuración de DWARF (`-w`) que el binario normal incluye para que un depurador pueda
inspeccionarlo. Para un binario que vas a repartir a producción y no vas a depurar ahí mismo, esos datos
no sirven de nada y solo ocupan espacio — para uno que estás desarrollando activamente, consérvalos.

### 7.8 Probar código que hace HTTP, sin red de verdad: `httptest`

Hasta la lección 6, las pruebas de concurrencia usaban `Falso` (un `Revisor` que nunca toca la red). Para
probar el programa **completo** —banderas, lectura de configuración, el cliente HTTP real, el reporte—
sin depender de un servicio externo ni de que algo esté escuchando en un puerto fijo, `net/http/httptest`
levanta un servidor real, en un puerto que el sistema operativo asigna solo, dentro del propio proceso de
pruebas:

<!-- verificar:extracto:cmd/revisor/main_test.go -->
```go
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
```

Corrida real:

```
$ go test ./cmd/revisor/... -run TestEjecutar_reportaOKyFalla -v
=== RUN   TestEjecutar_reportaOKyFalla
--- PASS: TestEjecutar_reportaOKyFalla (0.00s)
```

**`servidor.URL` es una URL real** (`http://127.0.0.1:PUERTO`, con un puerto libre que el sistema eligió),
así que el `Revisor` HTTP de la sección 7.1 la consulta exactamente igual que consultaría cualquier
servicio de producción — la única diferencia es que el "servicio" es un `http.HandlerFunc` de quince
líneas que vive en el mismo proceso que la prueba. Esto es lo que permite probar de extremo a extremo
—banderas, configuración, cliente HTTP, reporte, código de salida— en milisegundos, sin abrir ningún
puerto fijo que pudiera chocar con otra cosa corriendo en la máquina, y sin dejar nada corriendo después
de que la prueba termina (`defer servidor.Close()` se encarga).

**`t.TempDir()`** crea una carpeta temporal que Go borra solo al terminar la prueba (a diferencia de
`os.CreateTemp("/tmp", ...)` a secas, que dejaría el archivo ahí para siempre si nadie lo borra a mano) —
la combinación de `httptest` para la red y `t.TempDir()` para el disco es lo que permite que
`TestEjecutar_reportaOKyFalla` no deje ningún rastro en el sistema después de correr.

### 7.9 Por qué dos límites de tiempo, no uno

El `revisor` acepta `--limite` (el tiempo máximo para **todo** el reporte) y cada línea de configuración
puede declarar su propio tiempo límite **por servicio**. No es redundante: resuelven dos preguntas
distintas. `--limite` responde *«¿cuánto tiempo, como máximo, estoy dispuesto a esperar el reporte
completo?»* — útil si el `revisor` corre dentro de otro proceso con su propio plazo (un chequeo de salud
que un orquestador espera cada cierto tiempo, por ejemplo). El tiempo por servicio responde *«¿cuánto es
razonable esperar a ESTE servicio en particular?»* — un servicio interno de baja latencia y otro que cruza
a un proveedor externo no deberían compartir el mismo límite.

La demostración de la sección 7.5 lo confirma con datos reales: `colgado` tenía un límite propio de
200ms declarado en el archivo de configuración, y se cortó en 202ms — mucho antes de que el `--limite`
general (10 segundos por omisión) tuviera oportunidad de intervenir. El límite más corto de los dos es
siempre el que manda, y eso es exactamente lo que quieres: que un solo servicio lento no consuma todo el
presupuesto de tiempo del reporte completo.

### 7.10 El binario "no necesita nada" — salvo una cosa: certificados

La lección 1 demostró que un binario de Go corre en una máquina Linux sin Go instalado, gracias al
enlazado estático. Es tentador extender esa idea a "no necesita nada del sistema, punto" — y lo comprobé
llevándola al extremo: metí el binario del `revisor` en una imagen de Docker construida `FROM scratch`,
la base más vacía que existe (ni siquiera tiene una cáscara de sistema operativo, solo el binario).

```dockerfile
FROM scratch
COPY revisor /revisor
COPY servicios.txt /servicios.txt
ENTRYPOINT ["/revisor", "--config", "/servicios.txt"]
```

Con un servicio real apuntando a `https://example.com`:

```
$ docker build -t revisor-demo:scratch .
$ docker run --rm revisor-demo:scratch --formato tabla
SERVICIO  ESTADO    TIEMPO  DETALLE
ejemplo   FALLA      176ms  no responde: Get "https://example.com": tls: failed to verify certificate: x509: certificate signed by unknown authority
```

**Falló — y no por el `revisor`, sino por algo que ninguna imagen `scratch` trae: la lista de autoridades
certificadoras.** Para validar un certificado TLS (el candado del `https://`), el sistema operativo
normalmente provee una lista de quién tiene permiso de firmar certificados válidos — típicamente en
`/etc/ssl/certs/`. Una imagen `scratch` no tiene ni esa carpeta, así que Go no puede verificar ningún
certificado y se niega a continuar, con toda razón: seguir sin verificar sería aceptar cualquier
certificado, válido o falso.

**La solución, verificada, es copiar solo ese archivo** desde una imagen que sí lo tenga, sin arrastrar
el resto del sistema operativo:

```dockerfile
FROM alpine:3.20 AS certs
RUN apk add --no-cache ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY revisor /revisor
COPY servicios.txt /servicios.txt
ENTRYPOINT ["/revisor", "--config", "/servicios.txt"]
```

```
$ docker build -f Dockerfile.certs -t revisor-demo:scratch-certs .
$ docker run --rm revisor-demo:scratch-certs --formato tabla
SERVICIO  ESTADO    TIEMPO  DETALLE
ejemplo   OK         238ms  200

$ docker images revisor-demo:scratch-certs --format "{{.Size}}"
14.7MB
```

**14.7 MB en total, contra 14.4 MB sin los certificados — menos de 300 KB por resolver el problema
correctamente.** El binario sí es autosuficiente para todo lo que es *código Go*; lo que nunca trae por sí
solo es la confianza de qué autoridades son legítimas, porque eso es información del mundo (qué
certificadoras existen y siguen vigentes), no algo que el compilador pueda decidir por ti.

---

## El error que vas a ver

| El síntoma | Mensaje / evidencia literal | Qué pasa y qué hacer |
|---|---|---|
| Falta la bandera obligatoria | `revisor: falta --config` seguido del uso automático de `flag` | `ejecutar` valida explícitamente que `--config` no venga vacío antes de hacer cualquier otra cosa, y sale con código 2 |
| Formato de salida inválido | `revisor: --formato debe ser "tabla" o "json", no "xml"` | Solo dos formatos existen; cualquier otro valor se rechaza antes de intentar nada |
| Servicio que nunca contesta | En la tabla: `FALLA` con detalle `se acabo el tiempo de espera`, tiempo cercano al timeout configurado | El `context.WithTimeout` por servicio (lección 6) cortó la espera; no es un bug, es el mecanismo funcionando |
| Cerrar el cuerpo de una respuesta nula | (si se hiciera mal) `panic: runtime error: invalid memory address or nil pointer dereference` | Pasaría si `defer resp.Body.Close()` se pusiera antes de comprobar `err`; el `revisor` lo evita comprobando el error primero (sección 7.1) |
| `error` serializado directo a JSON sin traducir | `{}` (objeto vacío, sin ningún dato útil, sin ningún error de compilación que avise) | `encoding/json` no sabe convertir una interfaz `error`; por eso `reporte` usa `lineaJSON` con el motivo ya convertido a texto (sección 7.4) |
| Código de salida sin revisar en un script | El script sigue como si nada, aunque un servicio haya fallado | `ejecutar` sí distingue 0/1/2 (sección 7.6); el que integra el `revisor` en un pipeline debe leer `$?`, no solo la salida impresa |
| Binario en una imagen `FROM scratch`, contra HTTPS | `no responde: Get "https://...": tls: failed to verify certificate: x509: certificate signed by unknown authority` | Falta la lista de autoridades certificadoras (`ca-certificates.crt`); cópiala desde una imagen que la tenga (sección 7.10) |

---

## Lo que se hace mal

- **Usar `http.DefaultClient` o `http.Get` directo en producción.** Sección 7.1: sin timeout propio, una
  conexión que nunca responde cuelga el programa para siempre, sin ningún síntoma previo.
- **Poner toda la lógica dentro de `main` y llamar a `os.Exit` en medio del código.** Vuelve la función
  imposible de probar con `go test`, porque `os.Exit` termina el proceso de pruebas junto con el
  programa. Separa "decidir el código de salida" de "salir de verdad" (sección 7.2).
- **Serializar `error` directo a JSON, esperando que "algo salga".** Sale un objeto vacío, sin ningún
  aviso de que algo no se pudo convertir — el tipo de falla silenciosa más difícil de detectar, porque
  el programa no truena ni se queja.
- **Cerrar el cuerpo de una respuesta HTTP antes de comprobar el error.** Si la petición falló, `resp` es
  `nil`, y llamar a un método sobre él hace panic. Comprueba el error primero, siempre.
- **No drenar el cuerpo de la respuesta antes de cerrarlo.** La conexión no se puede reutilizar, y un
  programa que hace muchas peticiones al mismo servidor termina abriendo una conexión nueva cada vez, más
  lento de lo necesario sin ningún error que lo señale.
- **Ignorar el código de salida del binario desde un script o un pipeline.** El `revisor` distingue
  "corrí bien, todo sano" (0), "corrí bien, algo falló" (1) y "no pude ni arrancar" (2) — desperdiciarlo
  leyendo solo la salida impresa es tirar información que el programa ya te está dando gratis.

---

## Ejercicios

1. Implementa (o revisa, si ya lo tienes) el `Revisor` HTTP real de la sección 7.1, con su propio cliente
   y timeout. Confirma que compila y que `go vet ./...` no se queja.
2. Levanta `cmd/servidor-demo` en un puerto libre y corre tu `revisor` contra sus cuatro rutas, como en la
   sección 7.5. Pega la tabla real que obtuviste, no una inventada.
3. Agrégale a `Tabla` (sección 7.4) una quinta columna, `PROTOCOLO`, que diga `https` o `http` según
   `Servicio.EsSeguro()` (el método que se definió en la lección 3). Ajusta el ancho fijo de esa columna
   a mano, sin necesitar el cálculo dinámico que ya tiene `anchoNombre`.
4. Corre el mismo reporte con `--formato json` y valida que el resultado es JSON legítimo (por ejemplo,
   pásalo por `jq .` o pégalo en un validador de JSON).
5. Provoca el error de `--formato` inválido (sección 7.6, tabla) y confirma el código de salida con
   `echo $?`.
6. Compila tu `revisor` para `GOOS=linux GOARCH=amd64` desde tu máquina actual y confirma con `file` que
   el binario resultante es el de la plataforma correcta.
7. (Un poco más difícil) Mide el tamaño de tu binario con y sin `-ldflags="-s -w"`, como en la sección
   7.7, y calcula el porcentaje de reducción.
8. (Cierre del curso) Corre la suite completa del proyecto — `go vet ./...`, `gofmt -l .` y
   `go test ./... -race -cover` — y pega la salida completa en tu bitácora. Si algo no está en verde,
   arréglalo antes de dar por cerrado el curso.

### Soluciones

1-7. No hay una única salida de referencia porque depende de tu propia implementación y tu propia
máquina; compara la forma de tu resultado contra las secciones correspondientes (7.1, 7.4, 7.5, 7.6, 7.7).

8. La corrida real, sobre el `revisor` de este curso, el 30-sep-2026:

   ```
   $ gofmt -l .
   (sin salida: nada por formatear)

   $ go vet ./...
   (sin salida: limpio)

   $ go test ./... -race -cover
   ok  	github.com/habil/revisor/cmd/revisor	1.21s	coverage: 79.2% of statements
   	github.com/habil/revisor/cmd/servidor-demo		coverage: 0.0% of statements
   ok  	github.com/habil/revisor/internal/config	1.17s	coverage: 97.4% of statements
   ok  	github.com/habil/revisor/internal/reporte	1.18s	coverage: 100.0% of statements
   ok  	github.com/habil/revisor/internal/revisar	1.33s	coverage: 61.9% of statements
   ok  	github.com/habil/revisor/internal/servicio	1.17s	coverage: 92.3% of statements
   ```

   `cmd/servidor-demo` sale en 0.0% sin `ok` ni `FAIL`, no porque algo esté roto: con `-cover`, un
   paquete sin ningún archivo `_test.go` igual se instrumenta y se reporta, pero no hay ninguna prueba
   que ejercite ese código. Es una herramienta de apoyo para probar el programa completo, no parte del
   `revisor` que se reparte, y no tiene lógica propia que decida nada — solo contesta lo que se le
   programó que conteste — por eso no le hacen falta pruebas propias.

   ⚠️ **El 61.9% de `internal/revisar` en esta tabla es el artefacto que explica la sección 5.7: por
   paquete, no cuenta lo que `TestEjecutar_reportaOKyFalla` (en `cmd/revisor`, un poco más arriba en esta
   misma corrida) ejercita de `Revisar` al hablar HTTP de verdad contra el servidor `httptest`.** Medido
   con `-coverpkg=./...` sobre todo el proyecto: `Revisar` sube a 86.4% y el total del proyecto es 81.5%,
   no 61.9%. Si vas a citar un número de cobertura para decidir algo, corre `-coverpkg=./...`, no confíes
   en el número por paquete solo.

---

## Cómo sé que lo logré

- [ ] Mi `Revisor` HTTP tiene su propio cliente, con timeout, y nunca usa `http.DefaultClient`.
- [ ] Separé `main` (que solo llama a `os.Exit`) de una función que hace el trabajo y devuelve un `int`.
- [ ] Mi binario acepta `--config`, `--formato`, `--paralelo` y `--limite` desde la línea de comandos.
- [ ] Corrí el `revisor` real contra un servidor real (el mío o el `servidor-demo` del curso) y la tabla
      que obtuve refleja exactamente lo que ese servidor hace.
- [ ] `--formato json` produce JSON válido, verificado con una herramienta que no soy yo leyéndolo.
- [ ] `echo $?` después de correr el `revisor` da 0, 1 o 2, y sé explicar cada uno.
- [ ] Compilé para otra plataforma (`GOOS`/`GOARCH`) y confirmé con `file` que el binario es el correcto.
- [ ] `go vet ./...`, `gofmt -l .` y `go test ./... -race -cover` están limpios en mi propia copia del
      proyecto — no lo supongo, lo corrí.

---

## Para leer más

1. [net/http package](https://pkg.go.dev/net/http) — la documentación oficial completa del cliente y
   servidor HTTP de la biblioteca estándar.
2. [Command flag](https://pkg.go.dev/flag) — la documentación oficial del paquete de banderas usado en
   esta lección.
3. [encoding/json: JSON and Go](https://go.dev/blog/json) — el artículo oficial sobre cómo Go traduce
   entre structs y JSON, incluidas las reglas de las etiquetas (`json:"..."`, `omitempty`, `-`).
4. [Build constraints y compilación cruzada](https://pkg.go.dev/cmd/go#hdr-Environment_variables) — la
   referencia de `GOOS`/`GOARCH` y las demás variables de entorno que controlan `go build`.
5. [text/tabwriter](https://pkg.go.dev/text/tabwriter) — la herramienta estándar para alinear columnas
   cuando el cálculo manual de la sección 7.4 se queda corto (varias columnas de ancho variable).

### Términos de esta lección

| Término | Qué significa |
|---|---|
| `http.Client` | el tipo que hace peticiones HTTP en Go; sin configurar, no tiene límite de tiempo |
| `flag.FlagSet` | conjunto de banderas de línea de comandos que se puede parsear a partir de un slice propio, no solo de `os.Args` |
| `omitempty` | opción de etiqueta JSON que quita un campo de la salida cuando vale cero o está vacío |
| código de salida | el número entero que un programa devuelve al terminar; 0 significa éxito por convención universal |
| `GOOS` / `GOARCH` | variables de entorno que le dicen a `go build` para qué sistema operativo y arquitectura compilar |
| `-ldflags="-s -w"` | opciones de compilación que quitan símbolos y datos de depuración para reducir el tamaño del binario |

---

## Y ahora, lo que sigue

Tienes un programa completo: consulta servicios reales, todos a la vez, con límites de tiempo, y produce
un reporte que un humano o un programa pueden leer. Tres cosas que lo vuelven más sólido, en el orden en
que más conviene atacarlas:

1. **`golangci-lint`** — junta decenas de analizadores estáticos en una sola corrida. Pásalo por tu
   propio `revisor` y lee lo que dice: casi siempre enseña algo que ni `go vet` ni las pruebas atrapan.
2. **[Effective Go](https://go.dev/doc/effective_go) otra vez.** Ahora que escribiste un programa
   completo, vas a leerlo distinto que en la lección 0.
3. **Lee código ajeno bueno.** El propio paquete `net/http` de la biblioteca estándar es un buen punto de
   partida: ya tienes con qué orientarte para entender por qué está escrito como está.

**Y lo prometido desde el README:** vas a escribir este mismo programa otra vez, en Rust, en el curso
hermano. No para comparar sintaxis, sino para ver el mismo problema resuelto con otras herramientas —eso
es lo que de verdad enseña en qué se diferencian los dos lenguajes.

---

**Anterior:** [Lección 6 — Concurrencia](06-concurrencia.md)
