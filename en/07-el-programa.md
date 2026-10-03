# Lesson 7 — The finished program

**Duration:** 90-120 minutes, or 2 sessions of 60. It is the last lesson of the `revisor`: by the end you
have a binary that really runs, is configured from the command line, speaks HTTP with real timeouts,
produces two output formats, and is compiled for another platform without leaving your machine.

**By the end you will be able to:**

- Write an HTTP client with an explicit timeout, and explain why `http.DefaultClient` is dangerous in
  production.
- Separate a program's logic (`ejecutar`) from its entry point (`main`), so you can test it without
  touching `os.Exit`.
- Define command-line flags with the standard `flag` package, without any external library.
- Serialize a structure to JSON with `encoding/json`, and explain which fields are lost and why.
- Run the `revisor` from start to finish against a real server (a test one, made by you) and read its
  output.
- Compile the same code for another architecture and another operating system without leaving your
  machine, and explain why that is possible.

---

## Why it matters

You already have the four pieces of the `revisor` built separately: `servicio` defines the vocabulary
(lesson 2-3), `config` reads the configuration (lesson 5), `revisar.Todos` queries everything at once
(lesson 6). The only thing missing is what turns those pieces into **a program someone else can use
without reading the source code**: that it speak HTTP for real (until now we only tested it with
`Falso`), that it receive flags from the terminal instead of values fixed in the code, that it produce a
format another program can consume, and that it can be compiled and distributed as a single file.

This is the lesson where the `revisor` stops being "code that works if I run it, on my machine, with my
test data" and becomes a binary you can copy to someone else, with the confidence that it is going to do
exactly what the command line asks of it.

---

## The concepts

### 7.1 The HTTP client: why never `http.DefaultClient`

This is how the real `Revisor` ended up, the one that does touch the network
(`programas/revisor/internal/revisar/revisar.go`):

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

🔴 **`http.DefaultClient` —the one you would use if you wrote `http.Get(url)` directly— has no timeout
at all.** If the server on the other side accepts the connection and never answers anything, that call
stays waiting **forever**, without an error, without a panic, just a program that one day stops making
progress and nobody knows why. It is one of the most expensive errors in Go in production precisely
because it gives no symptom until it is already happening.

The `revisor` protects itself twice, not once: its own `http.Client` has a general timeout (the limit for
the whole report, passed to `NuevoHTTP`), and in addition each individual query runs under the
per-service `context.WithTimeout` that `Todos` set up in lesson 6 — that second one, shorter, is the one
that rules in practice. The client's timeout is the last-resort safety net, not the main mechanism.

⚠️ **`defer resp.Body.Close()` goes after checking the error, never before.** If `err != nil`, `resp` is
`nil`, and calling a method on a nil pointer panics. And **`io.Copy(io.Discard, resp.Body)` before
closing** is not decorative: without draining the response body, the underlying TCP connection can't be
reused for the next request to the same server, and the program ends up opening a new connection every
time instead of reusing the ones it already has.

**`traducir` changes the raw error from `net/http` into one that is readable in a table:**

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

The error `net/http` gives out of the box carries the full URL and the word `Get`, which in a report table
only take up space and tell the reader nothing new — that is why it is translated before being shown.

### 7.2 `main` isn't tested; `ejecutar` is

The pattern that separates the `revisor` from a one-file program:

<!-- verificar:extracto:cmd/revisor/main.go -->
```go
func main() {
	os.Exit(ejecutar(os.Args[1:], os.Stdout, os.Stderr))
}
```

`ejecutar` is where the real logic lives —sections 7.3, 7.4, and 7.6 show it piece by piece—; its
signature:

<!-- verificar:fragmento -->
```go
func ejecutar(args []string, salida, errores *os.File) int {
	// toda la lógica real vive aquí
}
```

**`go test` can't test a function that calls `os.Exit`**, because `os.Exit` ends the entire process
immediately — including the test process itself, which would never get to report the result. By
separating "deciding which exit code applies" (`ejecutar`, which **returns** an `int`) from "really exiting
with that code" (`main`, the only line that calls `os.Exit`), all the logic can be tested by calling
`ejecutar` directly, with test arguments, and checking the number it returns — exactly what
`cmd/revisor/main_test.go` does:

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

That is the real output of `--help` (generated automatically by `flag`, section 7.3) captured in the middle
of a test, without opening a terminal or running the compiled binary.

### 7.3 Command-line flags, without libraries

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

**`flag.NewFlagSet` instead of the `flag` package at the global level** (which you would use with
`flag.String(...)` directly) is what makes it possible to have an `ejecutar(args []string, ...)` function
that receives its arguments as a parameter, instead of always reading them from `os.Args` — another detail
that exists specifically so it can be tested, as in section 7.2.

`flag` comes in the standard library and is enough for a program of this size: four flags, basic types
(`string`, `int`, `time.Duration`), automatic help. If some day the `revisor` needed subcommands
(`revisor check`, `revisor list`, each one with its own flags), then it would be worth looking at a library
like `cobra` — but not before you really need it.

### 7.4 The two output formats: table and JSON

This is how `Tabla` ended up, the function that produces the human-readable output you have seen
throughout the lesson (`programas/revisor/internal/reporte/tabla.go`):

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

**The width of the first column is not fixed in the code — it is calculated.** The loop above goes over
all the states once before printing anything, and keeps the longest name (starting from the width of the
word "SERVICIO" itself, in case some service name were shorter than that). That number goes in as the `*`
in `%-*s`: a formatting verb with a **variable** width, where the width itself is one more argument of
`Fprintf`, not a number written by hand. It is the reason the real table in this lesson has straight
columns regardless of whether the service names are short (`sano`) or long (`inventario`): with a fixed
width, one of the two cases would always look crooked.

`ordenarPorNombre` (sections 6.3 and 6.7 already explained why the order of arrival can't be trusted:
several goroutines deliver results through a channel, and whoever answers first wins) makes a copy of the
slice and sorts it alphabetically before printing — without this, the same query would produce a table in
a different order every time you ran the program, even if the data were identical.

`etiquetaEstado`, `formatoTiempo`, and `detalle` are the three small functions that decide which word goes
in each column (`OK`/`LENTO`/`FALLA`, the time rounded to the millisecond or a dash if there was never a
response, and the HTTP code or the reason for the failure) — they are the ones you already saw at work in
the real table of section 7.5, now with the code that produces them.

🔑 **Why `text/tabwriter` wasn't used, which is the tool the standard library offers precisely for
aligning columns:** for a four-column table where only one has a variable width (the service name; the
other three are short and predictable), calculating the width by hand is simpler to read than introducing
a `tabwriter.Writer` with its own `Flush()` and tab separators. With a table with more variable columns,
`tabwriter` would be the right tool — it is worth knowing it (it is in the "Further reading" section of
this lesson) even though the `revisor` doesn't need it.

With the table understood, the other output format — JSON — is the other half of this section:

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

⚠️ **`servicio.Estado` is not serialized directly — on purpose.** `Estado` carries an `Err error` field, and
`error` is an interface: `encoding/json` doesn't know how to convert it to text on its own (trying it
produces an empty object `{}`, not a compilation error, which is worse: it fails silently). The `revisor`
solves this with an intermediate type, `lineaJSON`, which only lives inside the `reporte` package: it
converts the error to its text (`e.Motivo()`) before encoding, and along the way decouples the public
format (what other programs are going to read) from the internal model (`Estado`) — if tomorrow `Estado`
gains a new field, the JSON already in circulation doesn't change just because something internal changed.

`omitempty` on `Codigo` and `Detalle` removes those fields from the JSON when they are zero or an empty
string — that way a healthy service doesn't drag along a meaningless `"detalle": ""`. Actually encoded:

<!-- verificar:extracto:internal/reporte/json.go -->
```go
codificador := json.NewEncoder(w)
codificador.SetIndent("", "  ")
return codificador.Encode(lineas)
```

### 7.5 End to end: the `revisor` running against a real server

To test the whole program together —not each piece separately— I built `cmd/servidor-demo`: a minimal
HTTP server with four routes that behave like the four cases the `revisor` has to be able to report.

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

With that server up on `:8091` and a configuration file pointing to its four routes, I ran the real
binary:

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

**Each line of that table corresponds exactly to the behavior programmed into the test server:** `sano`
answers fast and with 200; `lento` really takes 1.5 seconds and that is why it comes out `LENTO`; `malo`
answers instantly but with 500; `colgado` never answers, and its individual timeout of 200ms —declared in
the configuration file, column three— cut off the wait at 202ms, almost exactly. **The exit code, 1, is
real**: at least one service was not `OK`, so `ejecutar` returns 1 instead of 0 (section 7.6) — the same
binary, used from a script, can tell whoever invokes it whether there were problems without anyone having
to read the table.

And in JSON, the same report:

```
$ /tmp/revisor --config /tmp/servicios-demo.txt --formato json
[
  {"servicio": "colgado", "ok": false, "tiempo_ms": 205, "detalle": "se acabo el tiempo de espera"},
  {"servicio": "lento", "ok": true, "codigo": 200, "tiempo_ms": 1503},
  {"servicio": "malo", "ok": false, "codigo": 500, "tiempo_ms": 1, "detalle": "codigo 500"},
  {"servicio": "sano", "ok": true, "codigo": 200, "tiempo_ms": 1}
]
```

(Here without the `SetIndent` indentation so that it fits on one line per service; the real program
produces it with line breaks, as you saw in section 7.4.)

### 7.6 The exit code, with meaning

<!-- verificar:extracto:cmd/revisor/main.go -->
```go
for _, e := range estados {
	if !e.OK() {
		return 1 // al menos un servicio falló: el código de salida lo refleja
	}
}
return 0
```

Three possible codes, and each one answers a different question for whoever invokes the program from a
script or a pipeline: **0** — everything went well; **1** — the program ran completely, but at least one
service was not healthy; **2** — the program couldn't even start (`--config` is missing, the file doesn't
exist, the configuration format is badly written). It is the difference between "I did the work and found
problems" and "I couldn't even start working" — and a continuous integration pipeline can react
differently to each one.

### 7.7 Cross-compilation: the same code, another platform

```bash
GOOS=linux GOARCH=amd64 go build -o revisor-linux-amd64 ./cmd/revisor
```

Actually confirmed, compiling from this Mac (arm64) to Linux x86-64:

```
$ GOOS=linux GOARCH=amd64 go build -o revisor-linux-amd64 ./cmd/revisor
$ file revisor-linux-amd64
revisor-linux-amd64: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, ...
```

**No Docker, no virtual machine, nothing additional installed: two environment variables are enough.**
This is possible because the Go compiler brings, out of the box, the code needed to generate binaries for
any combination of operating system and architecture it supports — it doesn't depend on running on the
target system to compile for it, unlike how many other compiled languages work.

**The size of the binary, measured, with and without debug symbols:**

```
$ go build -o revisor-normal ./cmd/revisor
$ wc -c revisor-normal
9464130 revisor-normal

$ go build -ldflags="-s -w" -o revisor-chico ./cmd/revisor
$ wc -c revisor-chico
6386194 revisor-chico
```

**From 9.46 MB to 6.39 MB, a real reduction of 32%**, by removing the symbol table (`-s`) and the DWARF
debugging information (`-w`) that the normal binary includes so that a debugger can inspect it. For a
binary you are going to distribute to production and are not going to debug right there, that data is of
no use and only takes up space — for one you are actively developing, keep it.

### 7.8 Testing code that does HTTP, without a real network: `httptest`

Up to lesson 6, the concurrency tests used `Falso` (a `Revisor` that never touches the network). To test
the **complete** program —flags, reading the configuration, the real HTTP client, the report— without
depending on an external service or on something listening on a fixed port, `net/http/httptest` starts a
real server, on a port the operating system assigns by itself, inside the test process itself:

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

Real run:

```
$ go test ./cmd/revisor/... -run TestEjecutar_reportaOKyFalla -v
=== RUN   TestEjecutar_reportaOKyFalla
--- PASS: TestEjecutar_reportaOKyFalla (0.00s)
```

**`servidor.URL` is a real URL** (`http://127.0.0.1:PUERTO`, with a free port the system chose), so the
HTTP `Revisor` from section 7.1 queries it exactly the same way it would query any production service —
the only difference is that the "service" is a fifteen-line `http.HandlerFunc` that lives in the same
process as the test. This is what makes it possible to test end to end —flags, configuration, HTTP client,
report, exit code— in milliseconds, without opening any fixed port that could clash with something else
running on the machine, and without leaving anything running after the test ends
(`defer servidor.Close()` takes care of it).

**`t.TempDir()`** creates a temporary folder that Go deletes by itself when the test ends (unlike a bare
`os.CreateTemp("/tmp", ...)`, which would leave the file there forever if nobody deleted it by hand) —
the combination of `httptest` for the network and `t.TempDir()` for the disk is what lets
`TestEjecutar_reportaOKyFalla` leave no trace on the system after running.

### 7.9 Why two time limits, not one

The `revisor` accepts `--limite` (the maximum time for the **whole** report) and each configuration line
can declare its own time limit **per service**. It is not redundant: they solve two different questions.
`--limite` answers *"how much time, at most, am I willing to wait for the complete report?"* — useful if
the `revisor` runs inside another process with its own deadline (a health check that an orchestrator
expects every so often, for example). The per-service time answers *"how long is it reasonable to wait for
THIS particular service?"* — a low-latency internal service and another one that crosses over to an
external provider shouldn't share the same limit.

The demonstration in section 7.5 confirms it with real data: `colgado` had its own limit of 200ms
declared in the configuration file, and it was cut off at 202ms — long before the general `--limite`
(10 seconds by default) had a chance to step in. The shorter of the two limits is always the one that
rules, and that is exactly what you want: that a single slow service doesn't consume the whole time budget
of the complete report.

### 7.10 The binary "needs nothing" — except one thing: certificates

Lesson 1 showed that a Go binary runs on a Linux machine without Go installed, thanks to static linking.
It is tempting to extend that idea to "it needs nothing from the system, period" — and I tested it by
taking it to the extreme: I put the `revisor` binary into a Docker image built `FROM scratch`, the emptiest
base there is (it doesn't even have an operating system shell, just the binary).

```dockerfile
FROM scratch
COPY revisor /revisor
COPY servicios.txt /servicios.txt
ENTRYPOINT ["/revisor", "--config", "/servicios.txt"]
```

With a real service pointing to `https://example.com`:

```
$ docker build -t revisor-demo:scratch .
$ docker run --rm revisor-demo:scratch --formato tabla
SERVICIO  ESTADO    TIEMPO  DETALLE
ejemplo   FALLA      176ms  no responde: Get "https://example.com": tls: failed to verify certificate: x509: certificate signed by unknown authority
```

**It failed — and not because of the `revisor`, but because of something no `scratch` image brings: the
list of certificate authorities.** To validate a TLS certificate (the padlock of `https://`), the
operating system normally provides a list of who has permission to sign valid certificates — typically in
`/etc/ssl/certs/`. A `scratch` image doesn't even have that folder, so Go can't verify any certificate and
refuses to continue, with good reason: continuing without verifying would mean accepting any certificate,
valid or fake.

**The solution, verified, is to copy only that file** from an image that does have it, without dragging
along the rest of the operating system:

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

**14.7 MB in total, against 14.4 MB without the certificates — less than 300 KB to solve the problem
correctly.** The binary is self-sufficient for everything that is *Go code*; what it never brings on its
own is the trust about which authorities are legitimate, because that is information about the world
(which certificate authorities exist and are still valid), not something the compiler can decide for you.

---

## The error you will see

| The symptom | Literal message / evidence | What happens and what to do |
|---|---|---|
| The mandatory flag is missing | `revisor: falta --config` followed by the automatic usage from `flag` | `ejecutar` explicitly validates that `--config` isn't empty before doing anything else, and exits with code 2 |
| Invalid output format | `revisor: --formato debe ser "tabla" o "json", no "xml"` | Only two formats exist; any other value is rejected before attempting anything |
| Service that never answers | In the table: `FALLA` with detail `se acabo el tiempo de espera`, time close to the configured timeout | The per-service `context.WithTimeout` (lesson 6) cut off the wait; it is not a bug, it is the mechanism working |
| Closing the body of a nil response | (if it were done wrong) `panic: runtime error: invalid memory address or nil pointer dereference` | It would happen if `defer resp.Body.Close()` were placed before checking `err`; the `revisor` avoids it by checking the error first (section 7.1) |
| `error` serialized directly to JSON without translating | `{}` (empty object, without any useful data, without any compilation error to warn you) | `encoding/json` doesn't know how to convert an `error` interface; that is why `reporte` uses `lineaJSON` with the reason already converted to text (section 7.4) |
| Exit code not checked in a script | The script carries on as if nothing happened, even though a service failed | `ejecutar` does distinguish 0/1/2 (section 7.6); whoever integrates the `revisor` into a pipeline must read `$?`, not just the printed output |
| Binary in a `FROM scratch` image, against HTTPS | `no responde: Get "https://...": tls: failed to verify certificate: x509: certificate signed by unknown authority` | The list of certificate authorities (`ca-certificates.crt`) is missing; copy it from an image that has it (section 7.10) |

---

## What goes wrong

- **Using `http.DefaultClient` or `http.Get` directly in production.** Section 7.1: without its own
  timeout, a connection that never responds hangs the program forever, without any prior symptom.
- **Putting all the logic inside `main` and calling `os.Exit` in the middle of the code.** It makes the
  function impossible to test with `go test`, because `os.Exit` ends the test process along with the
  program. Separate "deciding the exit code" from "really exiting" (section 7.2).
- **Serializing `error` directly to JSON, expecting "something to come out".** An empty object comes out,
  without any warning that something couldn't be converted — the hardest kind of silent failure to detect,
  because the program neither blows up nor complains.
- **Closing the body of an HTTP response before checking the error.** If the request failed, `resp` is
  `nil`, and calling a method on it panics. Check the error first, always.
- **Not draining the response body before closing it.** The connection can't be reused, and a program
  that makes many requests to the same server ends up opening a new connection every time, slower than
  necessary without any error to point it out.
- **Ignoring the binary's exit code from a script or a pipeline.** The `revisor` distinguishes "I ran
  fine, everything healthy" (0), "I ran fine, something failed" (1), and "I couldn't even start" (2) —
  wasting it by reading only the printed output is throwing away information the program is already
  giving you for free.

---

## Exercises

1. Implement (or review, if you already have it) the real HTTP `Revisor` from section 7.1, with its own
   client and timeout. Confirm that it compiles and that `go vet ./...` doesn't complain.
2. Start `cmd/servidor-demo` on a free port and run your `revisor` against its four routes, as in
   section 7.5. Paste the real table you got, not a made-up one.
3. Add to `Tabla` (section 7.4) a fifth column, `PROTOCOLO`, that says `https` or `http` according to
   `Servicio.EsSeguro()` (the method defined in lesson 3). Set the fixed width of that column by hand,
   without needing the dynamic calculation `anchoNombre` already has.
4. Run the same report with `--formato json` and validate that the result is legitimate JSON (for
   example, pipe it through `jq .` or paste it into a JSON validator).
5. Cause the invalid `--formato` error (section 7.6, table) and confirm the exit code with `echo $?`.
6. Compile your `revisor` for `GOOS=linux GOARCH=amd64` from your current machine and confirm with `file`
   that the resulting binary is the one for the right platform.
7. (A bit harder) Measure the size of your binary with and without `-ldflags="-s -w"`, as in section
   7.7, and calculate the percentage of reduction.
8. (Course wrap-up) Run the project's complete suite — `go vet ./...`, `gofmt -l .`, and
   `go test ./... -race -cover` — and paste the complete output in your logbook. If something is not
   green, fix it before considering the course closed.

### Solutions

1-7. There is no single reference output because it depends on your own implementation and your own
machine; compare the shape of your result against the corresponding sections (7.1, 7.4, 7.5, 7.6, 7.7).

8. The real run, on this course's `revisor`, on 30-sep-2026:

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

   `cmd/servidor-demo` comes out at 0.0% without `ok` or `FAIL`, not because something is broken: with
   `-cover`, a package without any `_test.go` file is still instrumented and reported, but there is no
   test that exercises that code. It is a support tool for testing the complete program, not part of the
   `revisor` that gets distributed, and it has no logic of its own that decides anything — it only answers
   what it was programmed to answer — that is why it doesn't need tests of its own.

   ⚠️ **The 61.9% of `internal/revisar` in this table is the artifact explained in section 5.7: per
   package, it doesn't count what `TestEjecutar_reportaOKyFalla` (in `cmd/revisor`, a little higher up in
   this same run) exercises of `Revisar` when it speaks real HTTP against the `httptest` server.** Measured
   with `-coverpkg=./...` over the whole project: `Revisar` goes up to 86.4% and the project total is
   81.5%, not 61.9%. If you are going to cite a coverage number to decide something, run
   `-coverpkg=./...`, don't trust the per-package number alone.

---

## How I know I got it

- [ ] My HTTP `Revisor` has its own client, with a timeout, and never uses `http.DefaultClient`.
- [ ] I separated `main` (which only calls `os.Exit`) from a function that does the work and returns an
      `int`.
- [ ] My binary accepts `--config`, `--formato`, `--paralelo`, and `--limite` from the command line.
- [ ] I ran the real `revisor` against a real server (mine or the course's `servidor-demo`) and the table
      I got reflects exactly what that server does.
- [ ] `--formato json` produces valid JSON, verified with a tool that isn't me reading it.
- [ ] `echo $?` after running the `revisor` gives 0, 1, or 2, and I can explain each one.
- [ ] I compiled for another platform (`GOOS`/`GOARCH`) and confirmed with `file` that the binary is the
      right one.
- [ ] `go vet ./...`, `gofmt -l .`, and `go test ./... -race -cover` are clean in my own copy of the
      project — I don't assume it, I ran it.

---

## Further reading

1. [net/http package](https://pkg.go.dev/net/http) — the complete official documentation of the standard
   library's HTTP client and server.
2. [Command flag](https://pkg.go.dev/flag) — the official documentation of the flags package used in this
   lesson.
3. [encoding/json: JSON and Go](https://go.dev/blog/json) — the official article on how Go translates
   between structs and JSON, including the rules for tags (`json:"..."`, `omitempty`, `-`).
4. [Build constraints and cross-compilation](https://pkg.go.dev/cmd/go#hdr-Environment_variables) — the
   reference for `GOOS`/`GOARCH` and the other environment variables that control `go build`.
5. [text/tabwriter](https://pkg.go.dev/text/tabwriter) — the standard tool for aligning columns when the
   manual calculation from section 7.4 falls short (several columns of variable width).

### Terms from this lesson

| Term | What it means |
|---|---|
| `http.Client` | the type that makes HTTP requests in Go; unconfigured, it has no time limit |
| `flag.FlagSet` | set of command-line flags that can be parsed from a slice of your own, not only from `os.Args` |
| `omitempty` | JSON tag option that removes a field from the output when it is zero or empty |
| exit code | the integer a program returns when it ends; 0 means success by universal convention |
| `GOOS` / `GOARCH` | environment variables that tell `go build` which operating system and architecture to compile for |
| `-ldflags="-s -w"` | build options that remove symbols and debugging data to reduce the size of the binary |

---

## And now, what comes next

You have a complete program: it queries real services, all at once, with time limits, and produces a
report that a human or a program can read. Three things that make it more solid, in the order in which it
is best to tackle them:

1. **`golangci-lint`** — brings together dozens of static analyzers in a single run. Run it over your own
   `revisor` and read what it says: it almost always teaches something that neither `go vet` nor the
   tests catch.
2. **[Effective Go](https://go.dev/doc/effective_go) again.** Now that you have written a complete
   program, you are going to read it differently than in lesson 0.
3. **Read good code written by others.** The standard library's own `net/http` package is a good starting
   point: you now have what you need to find your way and understand why it is written the way it is.

**And what was promised since the README:** you are going to write this same program again, in Rust, in
the sibling course. Not to compare syntax, but to see the same problem solved with other tools —that is
what really teaches how the two languages differ.

---

**Previous:** [Lesson 6 — Concurrency](06-concurrencia.md)
