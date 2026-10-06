# Lesson 5 — Modules and tests

**Duration:** 90 minutes, or 2 sessions of 45.

**By the end you will be able to:**

- Explain what `go.mod` solves and why the `revisor` doesn't need `go.sum`.
- Reorganize a one-file program into packages under `internal/`, and explain what that folder gives you
  that a normal folder doesn't.
- Write tests with a table of cases, without any external library, and read what they report when they
  fail.
- Run `go test` with `-v`, `-run`, `-cover`, and `-race`, and explain what each flag measures.
- Recognize the "zero tests" trap that looks identical to "all of them passed".
- Decide, with a real coverage number in front of you, which part of that gap is worth fixing and which
  isn't.

---

## Why it matters

Up to lesson 4 you had a single file, `main.go`, with everything inside: the `Servicio` and `Estado`
structs, the `Revisor` interface, the function that builds the report. It worked, and it worked well —but
a single file has a ceiling. As soon as you want to **test** one piece without running the whole program
(without touching the network, without reading a real file), a single `main.go` won't let you: everything
is mixed up with everything else.

This is the lesson where the `revisor` stops being a one-file exercise and becomes **a real project**:
several packages, each one with one responsibility, and a test suite that shows that each piece does what
it says without needing the others. It is the same leap you took in lesson 1 from "a program that fits in
your head" to "a program that lives in a folder with `go.mod`" — now that `go.mod` is going to organize
more than one file.

🔑 **And it is not an organizational whim.** A test that needs the network, a file on disk, or a running
server in order to run is a slow, fragile test that nobody runs often. Splitting into packages is what lets
you write tests that run in milliseconds, without touching anything external — and that is what makes you
actually run them, on every change, not only when you remember.

---

## The concepts

### 5.1 `go.mod`, and why this project has no `go.sum`

You already used `go mod init` in lesson 1 for the `hola` program. The `revisor`'s `go.mod` is just as
simple:

```
module github.com/habil/revisor

go 1.27
```

Two lines: the name of the module (that is how another program would import it, if some day you publish
one of its packages) and the minimum Go version it needs.

**If you search for module tutorials on the internet, almost all of them are going to mention `go.sum`
right away** — the file with the cryptographic hashes of each external dependency, so that nobody can slip
you a different version of a package you use. The `revisor` **has no `go.sum`**, and it is neither an
error nor an oversight:

```bash
$ ls go.sum
ls: go.sum: No such file or directory
```

**It doesn't exist because the `revisor` doesn't import a single external package.** Check the `import`s
of any file in the project and you are going to find only packages from the standard library: `net/http`,
`encoding/json`, `context`, `sync`, `flag`, `os`, `time`, `strings`, `sort`. It is a deliberate decision,
not a limitation: lesson 0 already warned about it —*"learn the standard library before any
framework"*— and the `revisor` is the proof that the standard library is enough for a complete program,
with concurrency, HTTP, JSON, and tests, without adding a single third-party dependency. If some day you
add one (for example, a real YAML client), at that moment `go get` is going to create `go.sum` for you,
and both files —`go.mod` and `go.sum`— go into the repository.

To really see what would have changed, I ran the test in a separate project (not in the `revisor`, which
still has no dependencies): a real `go get` of a small external package, `gopkg.in/yaml.v3`.

```
$ go get gopkg.in/yaml.v3
go: downloading gopkg.in/yaml.v3 v3.0.1
go: added gopkg.in/yaml.v3 v3.0.1

$ cat go.mod
module demo
go 1.27.1
require gopkg.in/yaml.v3 v3.0.1 // indirect

$ cat go.sum
gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
```

That is what appears as soon as you add **a single** external dependency: `go.mod` gains a `require`
line, and `go.sum` is born with the cryptographic hashes (the long texts that start with `h1:`) of that
dependency and of its own dependencies (`check.v1` is an indirect dependency of `yaml.v3`, not something
you asked for). Those hashes are what make `go build` refuse to compile if someone tried to slip you a
different version of the package with the same name and version — it is an integrity guarantee, not just
a record. The `revisor` doesn't have them because it doesn't need them: zero external dependencies, zero
surface for that kind of risk.

### 5.2 Packages by responsibility, not by layer

This is how the `revisor` ended up organized in this lesson:

```
revisor/
  go.mod
  cmd/
    revisor/            → package main: arranca, parsea banderas, imprime
    servidor-demo/       → package main: un servidor de prueba, no parte del programa
  internal/
    servicio/           → package servicio: los tipos (Servicio, Estado) y sus métodos
    config/             → package config: lee y valida el archivo de servicios
    revisar/            → package revisar: la interfaz Revisor y quien la implementa
    reporte/            → package reporte: convierte estados en tabla o en JSON
```

**By responsibility, not by layer.** There is no `modelos` package with all the program's structs, nor a
`utilidades` package with loose functions: each package has a single question it knows how to answer.
`servicio` knows what a service and a state are. `config` knows how to read the configuration. `revisar`
knows how to query. `reporte` knows how to print. If tomorrow you change how the table looks, you touch one
file, not five.

🔑 **`internal/` is a compiler rule, not a convention of good manners.** Any package that lives under a
folder called `internal/` can only be imported by code that is **inside the same module**, at any level
above that `internal/`. Check it: if another Go module —any of them, not only one of yours— tries
`import "github.com/habil/revisor/internal/servicio"`, the compiler refuses to compile, with an explicit
message that this package is internal. It is not a recommendation you can ignore under pressure: it is a
real restriction, the same kind of guarantee that the capital letter gives you inside a struct (lesson 3),
but at the level of a whole package.

⚠️ **What we did NOT do, on purpose: a `utils`, `helpers`, or `common` package.** It is the most repeated
antipattern in real projects: someone creates a folder for "things that don't fit anywhere else", and that
folder grows without limit until nobody knows what is inside or why. Every time you feel tempted to put
something in a package like that, ask yourself which **responsibility** that function belongs to, and put
it in the package that owns that responsibility — or if it really doesn't fit in any of them, that is a
sign that a new concept is missing a name, not that a catch-all drawer is missing.

### 5.3 The configuration file, and its real errors

This lesson's `config` reads a simple text format, one line per service:

```
nombre  url  [tiempo-limite]
```

```
# las lineas que empiezan con # se ignoran
catalogo   https://catalogo.interno.mx
pagos      https://pagos.interno.mx      500ms
```

This is done by the `Interpretar` function, which **never reads a file**: it receives the bytes already
read, and that is why it can be tested with a hundred content variants without creating a single temporary
file (`Cargar`, which does touch the disk, is a thin layer on top that only reads the file and passes the
content to `Interpretar`). Each badly written line produces a real error, with the line number and what
was expected — tested, not assumed:

```
$ (línea: "catalogo")
prueba.txt:1: esperaba «nombre url [tiempo]», hay 1 campo(s)

$ (línea: "catalogo catalogo.interno.mx")
prueba.txt:1: la URL "catalogo.interno.mx" debe empezar con http:// o https://

$ (línea: "catalogo https://a.mx nombas")
prueba.txt:1: tiempo limite "nombas" invalido (usa 500ms, 2s, 1m): time: invalid duration "nombas"

$ (dos líneas con el mismo nombre "catalogo")
prueba.txt:2: el nombre "catalogo" ya estaba en la linea 1
```

Notice the last one: **it wraps the error from `time.ParseDuration` with `%w`** (lesson 3) instead of
inventing its own text — that way, if you ever need to programmatically distinguish "invalid duration"
from another kind of error with `errors.As`, the original information is still there.

### 5.4 Tests: no libraries, with a table of cases

This is what a real test of the `servicio` package looks like (the complete file lives in
`programas/revisor/internal/servicio/servicio_test.go`):

<!-- verificar:extracto:internal/servicio/servicio_test.go -->
```go
func TestTimeoutEfectivo(t *testing.T) {
	casos := []struct {
		nombre  string
		timeout time.Duration
		quiere  time.Duration
	}{
		{"declarado", 5 * time.Second, 5 * time.Second},
		{"cero (valor por omision)", 0, TimeoutPorOmision},
		{"negativo (dato corrupto, no debe pasar)", -1 * time.Second, TimeoutPorOmision},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			s := Servicio{Timeout: c.timeout}
			if got := s.TimeoutEfectivo(); got != c.quiere {
				t.Errorf("TimeoutEfectivo() = %v, quería %v", got, c.quiere)
			}
		})
	}
}
```

**There is no `assert`, no `expect`, no assertion library of any kind — and it is on purpose.** Go compares
values with a normal `if` and reports with `t.Errorf`, with you writing what you expected and what you
got. At first it feels more verbose than an `assert.Equal(t, esperado, obtenido)` from other languages; the
gain is that you control the failure message, instead of inheriting a library's generic format, and that
there is nothing extra to install or learn in order to write the simplest test.

**`t.Run` gives a name to each case in the table**, and that matters when something fails: instead of a
generic "`TestTimeoutEfectivo` failed", the report says exactly which of the three cases it was:

```
$ go test ./internal/servicio/... -run TestTimeoutEfectivo -v
=== RUN   TestTimeoutEfectivo
=== RUN   TestTimeoutEfectivo/declarado
=== RUN   TestTimeoutEfectivo/cero_(valor_por_omision)
=== RUN   TestTimeoutEfectivo/negativo_(dato_corrupto,_no_debe_pasar)
--- PASS: TestTimeoutEfectivo (0.00s)
    --- PASS: TestTimeoutEfectivo/declarado (0.00s)
    --- PASS: TestTimeoutEfectivo/cero_(valor_por_omision) (0.00s)
    --- PASS: TestTimeoutEfectivo/negativo_(dato_corrupto,_no_debe_pasar) (0.00s)
PASS
ok  	github.com/habil/revisor/internal/servicio	0.382s
```

That output is real: I ran it on the code of this very project before writing this line. **Adding a new
case to the table is adding a line to the `casos` slice** — not a new function, not repeating the body of
the test. It is the idiomatic way of testing a function with many inputs in Go, and you are going to use it
in every package from here on.

### 5.5 Testing the errors, not only the successes

A table of cases also serves to test that something **fails as it should**, not only that it works
(abbreviated version here, with 3 of the 6 real cases and a positional struct literal instead of one with
field names, so that it fits; the complete one lives in `internal/config/config_test.go`):

<!-- verificar:fragmento -->
```go
func TestInterpretar_casosDeError(t *testing.T) {
	casos := []struct {
		nombre   string
		entrada  string
		contexto string // fragmento que el mensaje de error debe contener
	}{
		{"nombre repetido", "catalogo https://a.mx\ncatalogo https://b.mx\n", "ya estaba en la linea"},
		{"url sin esquema", "catalogo catalogo.interno.mx\n", "debe empezar con http"},
		{"tiempo limite invalido", "catalogo https://a.mx nombas\n", "invalido"},
		// ...
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			_, err := Interpretar([]byte(c.entrada), "prueba.txt")
			if err == nil {
				t.Fatalf("Interpretar() no devolvió error, y debía contener %q", c.contexto)
			}
			if !strings.Contains(err.Error(), c.contexto) {
				t.Errorf("Interpretar() error = %q, quería que contuviera %q", err.Error(), c.contexto)
			}
		})
	}
}
```

**`t.Fatalf` instead of `t.Errorf`** in the first `if`: if `Interpretar` didn't return an error when it
should have, continuing to check `err.Error()` on the next line would panic (`err` would be `nil`).
`Fatalf` stops that particular test right there; `Errorf` lets the test keep running and accumulate more
failures before reporting. The practical rule: use `Fatalf` when continuing makes no sense without what
you just checked, `Errorf` when it does.

### 5.6 The commands you are always going to use

```bash
go test ./...                          # todo el proyecto
go test ./internal/config/... -v       # verboso, un paquete
go test ./... -run TestInterpretar     # solo las pruebas cuyo nombre haga match
go test ./... -race                    # detector de carreras (se explica a fondo en la lección 6)
go test ./... -cover                   # porcentaje de líneas ejercitadas por las pruebas
```

Actually run on the `revisor`, on 30-sep-2026:

```
$ go test ./internal/servicio/... ./internal/config/... -cover
ok  	github.com/habil/revisor/internal/servicio	0.195s	coverage: 92.3% of statements
ok  	github.com/habil/revisor/internal/config	0.192s	coverage: 97.4% of statements
```

### 5.7 What to do with a coverage number

**92.3% and 97.4% are not goals, they are starting points for a question: what is that 8% and that 3%
that wasn't executed, and do I care?** With `go test -coverprofile` you can see exactly which lines were
left untouched:

```bash
go test ./internal/servicio/... -coverprofile=/tmp/cobertura.out
go tool cover -func=/tmp/cobertura.out
```

In the `revisor`, the gap in `servicio` is the branch of `Motivo()` that builds the message when the code
is neither 0 nor a success with a specific text (`fmt.Sprintf("codigo %d", ...)` for a 404 with no further
context) — a branch the existing tests don't exercise with that exact code. **The right decision is not
to chase 100%** by filling each branch with a forced test that teaches nothing new: it is to look at the
gap, decide whether it matters (here, a little: you would add a case with 404 to the table), and write it
down, instead of pretending it doesn't exist.

🔴 **And a real trap, measured in this very project: per-package coverage can underestimate a central
function without warning you.** Run on its own, the `revisor`'s `revisar` package (which you are going to
get to know in depth in lesson 6) reports:

```
$ go test ./internal/revisar/... -cover
ok  	github.com/habil/revisor/internal/revisar	1.33s	coverage: 61.9% of statements
```

**61.9% sounds like more than a third of that package is never executed in any test — and it is false.**
Its most important function, `Revisar` (the one that really speaks HTTP), is well tested: it is just that
the test that really exercises it —`TestEjecutar_reportaOKyFalla`, with a real `httptest` server, in
lesson 7— lives in the `cmd/revisor` package, not in `internal/revisar`. `go test ./internal/revisar/...`
**only counts what the tests OF THAT PACKAGE exercise**; it doesn't see what a test from another package
goes through along the way, even if it passes through the same code. With `-coverpkg`, which asks Go to
measure a package's coverage counting **all** the project's tests, not only its own:

```
$ go test ./... -coverpkg=./... -coverprofile=/tmp/cov.out
$ go tool cover -func=/tmp/cov.out | grep 'revisar.go.*Revisar'
github.com/habil/revisor/internal/revisar/revisar.go:48:  Revisar    86.4%

$ go tool cover -func=/tmp/cov.out | tail -1
total:                                                    (statements)    81.5%
```

**86.4% for `Revisar`, 81.5% for the whole project — not 61.9%.** The lesson is not "ignore the
per-package number": it is that a coverage number always answers an implicit question —coverage of what,
measured against the tests from where?— and `go test ./paquete/... -cover` keeps quiet about that second
half of the question. Before deciding that something "isn't tested" because of a low number, run
`-coverpkg=./...` over the whole project and compare.

### 5.8 Causing a failure, to know the test works

A test you have never seen fail is a test you don't know works — it could be comparing two things that are
always equal by accident. Break it on purpose: change `TimeoutEfectivo()` so it always returns
`s.Timeout`, without the `if`, and run the test:

```
$ go test ./internal/servicio/... -run TestTimeoutEfectivo -v
=== RUN   TestTimeoutEfectivo/cero_(valor_por_omision)
    servicio_test.go:64: TimeoutEfectivo() = 0s, quería 2s
--- FAIL: TestTimeoutEfectivo (0.00s)
    --- FAIL: TestTimeoutEfectivo/cero_(valor_por_omision) (0.00s)
FAIL
```

**That `FAIL`, with the exact case that broke and the value it got against the one it expected, is the
proof that the test works.** Put the `if` back and confirm it is green again before moving on.

### 5.9 Designing to be able to test: separating logic from I/O

Notice something we already mentioned in passing in section 5.3 and that deserves its own space: `config`
has **two** functions, not one.

<!-- verificar:extracto:internal/config/config.go -->
```go
func Cargar(ruta string) ([]servicio.Servicio, error) {
	datos, err := os.ReadFile(ruta)
	if err != nil {
		return nil, fmt.Errorf("leyendo la configuracion: %w", err)
	}
	return Interpretar(datos, ruta)
}
```

`Interpretar` (just the signature; the complete body, with all the parsing logic, is in section 5.3, with
its real errors):

<!-- verificar:fragmento -->
```go
func Interpretar(datos []byte, origen string) ([]servicio.Servicio, error) {
	// ... toda la lógica de parseo, sin tocar el disco ...
}
```

**If there were a single function `Cargar(ruta string)` that read the file and parsed everything
together**, each test for "badly written line", "repeated name", or "invalid timeout" would have to start
by creating a temporary file on disk with `os.CreateTemp`, writing the case's content to it, passing it the
path, and deleting it at the end. It works, but it is slow (it touches the real file system) and it clutters
the test with code that has nothing to do with what is really being tested: **whether your configuration
is interpreted correctly or not**, not whether you know how to create temporary files.

By separating **the part that decides** (`Interpretar`, a pure function: same input bytes, same result
always, without access to anything external) from **the part that obtains the bytes** (`Cargar`, the only
one that touches the disk), the nine tests in section 5.5 run in microseconds and without creating a single
file. `Cargar` itself hardly needs tests of its own: only checking that it returns an error if the file
doesn't exist (exercise 6), because all the interesting logic is already in `Interpretar` and is already
tested.

🔑 **The general rule, useful far beyond this project:** when a function is hard to test, it is almost
always because it mixes "deciding something" with "touching the outside world" (a file, the network, the
clock). Separating them is not a style rule: it is what determines whether you are going to be able to
write the test in three lines or in twenty.

### 5.10 Benchmarks: measuring, not guessing

Besides `Test...`, Go recognizes `Benchmark...` functions that measure how long your code takes, not
whether it is correct:

<!-- verificar:extracto:internal/config/config_bench_test.go -->
```go
func BenchmarkInterpretar(b *testing.B) {
	entrada := []byte(`
catalogo   https://catalogo.interno.mx
pagos      https://pagos.interno.mx      500ms
inventario https://inventario.interno.mx 1s
reportes   https://reportes.interno.mx
`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Interpretar(entrada, "bench.txt"); err != nil {
			b.Fatal(err)
		}
	}
}
```

You don't choose `b.N`: Go runs the loop with larger and larger values of `N` until the measurement is
stable, and reports the time per operation. Actually run on the `revisor` (Apple M5, 200,000 repetitions):

```
$ go test ./internal/config/... -bench=. -run '^$'
goos: darwin
goarch: arm64
pkg: github.com/habil/revisor/internal/config
cpu: Apple M5
BenchmarkInterpretar-10    	  200000	       432.5 ns/op
PASS
```

**`-run '^$'`** tells `go test` not to run any normal test (a regular expression that doesn't match any
name), so that the benchmark report doesn't get mixed with the tests'. The number —432.5 nanoseconds per
call, on this machine, at this moment— is not for memorizing: it is for comparing **against itself** after
a change. If tomorrow you rewrite `Interpretar` and the benchmark goes up to 4,000 ns/op, you have an
objective signal that something got slower, without needing to have an opinion about it.

### 5.11 `go vet`: the one that finds what compiles but is wrong

`go test` tells you whether your logic does what you expected. **`go vet` tells you whether your code has
an error that the compiler doesn't catch because, technically, it is valid** — but it is almost certainly
not what you meant to write. The most common case is a formatting verb (lesson 2) that doesn't match the
type of the argument:

<!-- verificar:ejemplo:ejemplos/05-vet-printf -->
```go
puerto := 443
fmt.Printf("servicio %s en el puerto %s\n", nombre, puerto) // %s para un int
```

This **compiles** and **runs**, without a panic or an error — and produces broken output (complete
program in `programas/revisor/ejemplos/05-vet-printf/main.go`):

```
$ go run ./ejemplos/05-vet-printf/
servicio catalogo en el puerto %!s(int=443)
```

And `go vet`, on that same file, does detect it:

```
$ go vet ./ejemplos/05-vet-printf/
ejemplos/05-vet-printf/main.go:11:39: fmt.Printf format %s has arg puerto of wrong type int
```

`go vet` does detect it, because it analyzes the format string against the real types of the arguments, a
step the Go compiler doesn't take on its own. Run `go vet ./...` together with `go test ./...` as a
routine: most editors with the Go extension (lesson 1) already do it for you while you type, underlining
the problem before you even get to run anything.

### 5.12 Boundary tests: where bugs really live

`Estado.OK()` decides that a code is fine if it falls **between 200 and 299**. It is tempting to test it
with an "obvious" case (200) and an "obvious" failure case (500) and call it done. The `revisor`'s real
table also tests the two values that sit **right on the edge of the range**:

<!-- verificar:fragmento -->
```go
{"199, justo debajo del rango", Estado{Codigo: 199}, false},
{"300, justo arriba del rango", Estado{Codigo: 300}, false},
```

```
$ go test ./internal/servicio/... -run TestEstadoOK -v
=== RUN   TestEstadoOK/199,_justo_debajo_del_rango
=== RUN   TestEstadoOK/300,_justo_arriba_del_rango
--- PASS: TestEstadoOK (0.00s)
    --- PASS: TestEstadoOK/199,_justo_debajo_del_rango (0.00s)
    --- PASS: TestEstadoOK/300,_justo_arriba_del_rango (0.00s)
```

**Why does it matter, if 199 and 300 "obviously" are not OK?** Because a one-character error in the
condition —`>=` instead of `>`, or `<=` instead of `<`— is exactly the kind of bug an "obvious" case never
detects, and a boundary case always detects. If someone changed `OK()` to
`e.Codigo >= 200 && e.Codigo <= 300` (including 300 by mistake), the 200/299/404/500 cases would keep
passing just as well — **only the 300 case would give it away.** That is the underlying reason behind
"test the edges, not only the center": the edges are where comparison errors hide, and they are invisible
to any test that only uses values far inside or far outside the range.

---

## The error you will see

| The symptom | Literal message | What happens and what to do |
|---|---|---|
| Import of someone else's `internal` package | `use of internal package github.com/habil/revisor/internal/servicio not allowed` | The compiler prevents importing something under `internal/` from outside the module. It is not a permission you can grant: you have to expose the type from a public package if it is really needed |
| Duplicate service in the configuration | `prueba.txt:2: el nombre "catalogo" ya estaba en la linea 1` | Two services with the same name; fix the configuration file |
| URL without a scheme | `prueba.txt:1: la URL "catalogo.interno.mx" debe empezar con http:// o https://` | `http://` or `https://` is missing at the start of the URL |
| Badly written timeout | `prueba.txt:1: tiempo limite "nombas" invalido (usa 500ms, 2s, 1m): time: invalid duration "nombas"` | Go's duration format is not free-form: use a number followed by a unit (`ms`, `s`, `m`, `h`) |
| A test you broke on purpose | `servicio_test.go:64: TimeoutEfectivo() = 0s, quería 2s` | See section 5.8: this is what a `t.Errorf` looks like, pointing out exactly which case failed and why |
| Zero tests in a package | `?   	github.com/habil/revisor/cmd/servidor-demo	[no test files]` | It is not a failure: `go test` warns explicitly when a package has no `_test.go` file at all, instead of pretending it ran something |

**And the most important one to learn to read**, because it is not an error but the absence of one:

```
$ go test ./...
ok  	github.com/habil/revisor/internal/vacio	0.001s
```

If `internal/vacio` had **no** `Test...` function at all, that line would look exactly the same: `ok`, in
green, with no mark that nothing was executed. The only way to distinguish "everything passed, for real"
from "there was nothing to run" is to look at the count with `-v` (which does print each `RUN`) or, better,
never trust a package that has no `_test.go` file at all — that `go test` does say, as in the row of the
table above.

---

## What goes wrong

- **Creating a `utils`, `helpers`, or `common` package.** You already saw it in section 5.2: it is the most
  repeated antipattern and the one that loses its purpose fastest. Name the responsibility, not the fact
  that it "doesn't fit anywhere else".
- **Confusing "it compiled" with "the tests passed".** `go build` verifies that the code is valid; it
  doesn't run a single test. They are two different commands with two different questions.
- **Reading the exit code of `go test` instead of the count.** A package without test files exits with
  code 0 (success), just like one with 50 tests that did pass. The exit code answers "did something
  fail?", not "was something tested?" — for the second question you need to read the output, not just
  the `$?`.
- **Chasing 100% coverage as if it were the goal.** High coverage with weak assertions (checking that
  something doesn't blow up, without checking what it returns) gives a pretty figure and a test that
  detects almost nothing. The number is a guide to where to look, not a goal in itself — section 5.7 shows
  it with a real gap in the `revisor` itself.
- **Writing a test and never seeing it fail.** If you never broke the code on purpose to confirm that the
  test turns red, you don't know whether that test tests anything. Section 5.8 did it with real data from
  the project: do it yourself too with at least one test of your own before calling it good.

---

## Exercises

1. Clone the package structure from section 5.2 (`internal/servicio`, `internal/config`) for your own
   copy of the `revisor`, moving the code you already had from lessons 2 to 4.
2. Write the table of cases for `TestEtiqueta` for the `Etiqueta()` method of `Servicio`, with at least
   one case with a normal name and one with an empty struct.
3. Add a case to `TestInterpretar_casosDeError` for a line with **four** fields (more than the three the
   format allows). Check the exact message against the code in `config.go`.
4. Run `go test ./... -cover` on your copy and write down the percentage of each package. Choose **one**
   coverage gap and decide, in writing in your logbook, whether you care about closing it and why.
5. (As in section 5.8) Break on purpose a function you already tested, run the test, read the full
   `FAIL`, and fix it. Paste both results —the red one and the green one— in your logbook.
6. (A bit harder) Write `TestCargar_archivoInexistente`, which confirms that `Cargar` (not
   `Interpretar`) returns an error when the path doesn't exist. Hint: you don't need to create any file
   for this test, just pass a path you know doesn't exist.

### Solutions

1 and 2 don't have a single reference solution: it depends on how you had your own lesson 4 program
organized. Compare your result against the real code of `programas/revisor/internal/servicio/servicio.go`
in this course's project.

3. With `"catalogo https://a.mx 500ms extra\n"`, the expected message is
   `prueba.txt:1: esperaba «nombre url [tiempo]», hay 4 campo(s)` — the same code path that already
   handles "missing fields", because `len(campos) > 3` covers both cases with a single check.

4. There is no single answer: what matters is that the decision is written down with its reason, not the
   percentage itself.

5. See the whole of section 5.8: the pattern is always "change the code, run the test, read the `FAIL`
   with the exact case, fix it, run it again".

6. Like this:

   <!-- verificar:extracto:internal/config/config_test.go -->
   ```go
   func TestCargar_archivoInexistente(t *testing.T) {
   	_, err := Cargar("/ruta/que/no/existe.txt")
   	if err == nil {
   		t.Fatal("Cargar() no devolvió error con una ruta inexistente")
   	}
   }
   ```

   This test does exist in the real project (`config_test.go`) and it passes because `os.ReadFile`, inside
   `Cargar`, returns an operating system error that `Cargar` wraps with `%w` before propagating it.

---

## How I know I got it

- [ ] My `revisor` is organized in packages under `internal/`, each one with a single responsibility.
- [ ] `go test ./...` runs and **the count, not just the color**, tells me how many tests were executed.
- [ ] I wrote at least one table of cases with `t.Run`, and I know how to read which case failed when
  something breaks.
- [ ] I ran `go test -cover` and I can say what percentage came out and what is in the gap.
- [ ] I broke a function on purpose, saw the exact `FAIL`, and fixed it — I have it in my logbook.
- [ ] I can explain why the `revisor` has no `go.sum` and what would generate it if some day it needed one.
- [ ] I know why I shouldn't create a `utils` package.

---

## Further reading

1. [Writing tests](https://go.dev/doc/tutorial/add-a-test) — the official testing tutorial, with a table
   of cases included.
2. [Go Wiki: TableDrivenTests](https://go.dev/wiki/TableDrivenTests) — the canonical reference for the
   pattern used throughout this lesson.
3. [Package internal](https://go.dev/doc/go1.4#internalpackages) — the original announcement of the
   `internal/` rule, straight from the Go release notes.
4. [Go Blog: The Cover Story](https://go.dev/blog/cover) — how `go test -cover` works inside, and why a
   high number doesn't always mean good tests.
5. [Release without fear on your own infrastructure](https://www.habil.mx/en/blog/cicd-devsecops-own-infrastructure/) — article on a path to production where tests, among other gates, can stop a release.

### Terms from this lesson

| Term | What it means |
|---|---|
| `go.sum` | file with the cryptographic hashes of the external dependencies; the `revisor` doesn't have one because it uses none |
| `internal/` | special folder that the Go compiler prevents from being imported from outside the module |
| table of cases | test pattern where a list of inputs and expected outputs is iterated over with a single test body |
| `t.Run` | runs a subcase with its own name, so the report says exactly which one failed |
| coverage | percentage of the code's lines that the tests exercised when running |

---

**Previous:** [Lesson 4 — Collections](04-colecciones.md) ·
**Next:** [Lesson 6 — Concurrency](06-concurrencia.md)
