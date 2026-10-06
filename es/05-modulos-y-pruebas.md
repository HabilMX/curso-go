# Lección 5 — Módulos y pruebas

**Duración:** 90 minutos, o 2 sesiones de 45.

**Al terminar vas a poder:**

- Explicar qué resuelve `go.mod` y por qué el `revisor` no necesita `go.sum`.
- Reorganizar un programa de un archivo en paquetes bajo `internal/`, y explicar qué te da esa carpeta
  que una carpeta normal no da.
- Escribir pruebas con tabla de casos, sin ninguna librería externa, y leer lo que reportan cuando fallan.
- Correr `go test` con `-v`, `-run`, `-cover` y `-race`, y explicar qué mide cada bandera.
- Reconocer la trampa de "cero pruebas" que se ve idéntica a "todas pasaron".
- Decidir, con un número real de cobertura delante, qué parte de ese hueco importa arreglar y cuál no.

---

## Por qué importa

Hasta la lección 4 tenías un único archivo, `main.go`, con todo adentro: los structs `Servicio` y
`Estado`, la interfaz `Revisor`, la función que arma el reporte. Funcionaba, y funcionaba bien —pero un
solo archivo tiene un techo. En cuanto quieres **probar** una pieza sin ejecutar el programa completo
(sin tocar la red, sin leer un archivo real), un solo `main.go` no te deja: todo está mezclado con todo.

Ésta es la lección donde el `revisor` deja de ser un ejercicio de un archivo y se convierte en **un
proyecto de verdad**: varios paquetes, cada uno con una responsabilidad, y una suite de pruebas que
demuestra que cada pieza hace lo que dice sin necesitar las demás. Es el mismo salto que diste en la
lección 1 de "un programa que cabe en la cabeza" a "un programa que vive en una carpeta con `go.mod`" —
ahora ese `go.mod` va a organizar más de un archivo.

🔑 **Y no es un capricho de organización.** Una prueba que necesita red, un archivo de disco o un
servidor levantado para correr es una prueba lenta, frágil y que nadie corre seguido. Separar en paquetes
es lo que te permite escribir pruebas que corren en milisegundos, sin tocar nada externo — y eso es lo
que hace que sí las corras, en cada cambio, no solo cuando te acuerdas.

---

## Los conceptos

### 5.1 `go.mod`, y por qué este proyecto no tiene `go.sum`

Ya usaste `go mod init` en la lección 1 para el programa `hola`. El `go.mod` del `revisor` es igual de
simple:

```
module github.com/habil/revisor

go 1.27
```

Dos líneas: el nombre del módulo (así es como otro programa lo importaría, si algún día publicas alguno
de sus paquetes) y la versión mínima de Go que necesita.

**Si buscas tutoriales de módulos en internet, casi todos van a mencionar `go.sum` enseguida** — el
archivo con las huellas criptográficas de cada dependencia externa, para que nadie te cuele una versión
distinta de un paquete que usas. El `revisor` **no tiene `go.sum`**, y no es un error ni un descuido:

```bash
$ ls go.sum
ls: go.sum: No such file or directory
```

**No existe porque el `revisor` no importa ni un solo paquete externo.** Revisa los `import` de cualquier
archivo del proyecto y vas a encontrar únicamente paquetes de la biblioteca estándar: `net/http`,
`encoding/json`, `context`, `sync`, `flag`, `os`, `time`, `strings`, `sort`. Es una decisión deliberada,
no una limitación: la lección 0 ya lo advertía —*"aprende la biblioteca estándar antes que cualquier
framework"*— y el `revisor` es la prueba de que la estándar alcanza para un programa completo, con
concurrencia, HTTP, JSON y pruebas, sin agregar una sola dependencia de terceros. Si algún día agregas
una (por ejemplo, un cliente de YAML de verdad), en ese momento `go get` va a crear `go.sum` por ti, y
los dos archivos —`go.mod` y `go.sum`— van al repositorio.

Para ver de verdad qué habría cambiado, hice la prueba en un proyecto aparte (no en el `revisor`, que
sigue sin dependencias): un `go get` real de un paquete externo pequeño, `gopkg.in/yaml.v3`.

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

Eso es lo que aparece en cuanto agregas **una sola** dependencia externa: `go.mod` gana una línea
`require`, y `go.sum` nace con las huellas criptográficas (los textos largos que empiezan con `h1:`) de
esa dependencia y de las suyas propias (`check.v1` es una dependencia indirecta de `yaml.v3`, no algo que
tú pediste). Esas huellas son lo que hace que, si alguien intentara colarte una versión distinta del
paquete con el mismo nombre y versión, `go build` se niegue a compilar — es una garantía de integridad,
no solo un registro. El `revisor` no las tiene porque no las necesita: cero dependencias externas, cero
superficie de esa clase de riesgo.

### 5.2 Paquetes por responsabilidad, no por capa

Así quedó organizado el `revisor` en esta lección:

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

**Por responsabilidad, no por capa.** No hay un paquete `modelos` con todos los structs del programa ni
un paquete `utilidades` con funciones sueltas: cada paquete tiene una sola pregunta que sabe responder.
`servicio` sabe qué es un servicio y un estado. `config` sabe leer la configuración. `revisar` sabe
consultar. `reporte` sabe imprimir. Si mañana cambias cómo se ve la tabla, tocas un archivo, no cinco.

🔑 **`internal/` es una regla del compilador, no una convención de buena educación.** Cualquier paquete
que viva bajo una carpeta llamada `internal/` solo puede ser importado por código que esté **dentro del
mismo módulo**, en cualquier nivel arriba de ese `internal/`. Compruébalo: si otro módulo de Go —cualquiera,
no solo uno tuyo— intenta `import "github.com/habil/revisor/internal/servicio"`, el compilador se niega
a compilar, con un mensaje explícito de que ese paquete es interno. No es una recomendación que puedas
ignorar bajo presión: es una restricción real, la misma clase de garantía que la mayúscula te da dentro
de un struct (lección 3), pero a nivel de paquete completo.

⚠️ **Lo que NO hicimos, a propósito: un paquete `utils`, `helpers` o `common`.** Es el antipatrón más
repetido en proyectos reales: alguien crea una carpeta para "cosas que no encajan en otro lado", y esa
carpeta crece sin límite hasta que nadie sabe qué hay adentro ni por qué. Cada vez que sientas la tentación
de poner algo en un paquete así, pregúntate de qué **responsabilidad** es esa función, y ponla en el
paquete dueño de esa responsabilidad — o si de verdad no encaja en ninguno, es una señal de que falta
nombrar un concepto nuevo, no de que falta un cajón de sastre.

### 5.3 El archivo de configuración, y sus errores reales

El `config` de esta lección lee un formato de texto simple, una línea por servicio:

```
nombre  url  [tiempo-limite]
```

```
# las lineas que empiezan con # se ignoran
catalogo   https://catalogo.interno.mx
pagos      https://pagos.interno.mx      500ms
```

Esto lo hace la función `Interpretar`, que **nunca lee un archivo**: recibe los bytes ya leídos, y por
eso se puede probar con cien variantes de contenido sin crear un solo archivo temporal (`Cargar`, que sí
toca disco, es una capa delgada encima que solo lee el archivo y le pasa el contenido a `Interpretar`).
Cada línea mal escrita produce un error real, con el número de línea y qué se esperaba — probado, no
supuesto:

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

Fíjate en el último: **envuelve el error de `time.ParseDuration` con `%w`** (lección 3) en vez de
inventar su propio texto — así, si alguna vez necesitas distinguir programáticamente "duración inválida"
de otro tipo de error con `errors.As`, la información original sigue ahí.

### 5.4 Pruebas: sin librerías, con tabla de casos

Así se ve una prueba real del paquete `servicio` (el archivo completo vive en
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

**No hay `assert`, ni `expect`, ni ninguna librería de aserciones — y es a propósito.** Go compara
valores con un `if` normal y reporta con `t.Errorf`, escribiendo tú mismo qué esperabas y qué obtuviste.
Al principio se siente más verboso que un `assert.Equal(t, esperado, obtenido)` de otros lenguajes; la
ganancia es que el mensaje de falla lo controlas tú, en vez de heredar el formato genérico de una
librería, y que no hay que instalar ni aprender nada extra para escribir la prueba más simple.

**`t.Run` le da nombre a cada caso de la tabla**, y eso importa cuando algo falla: en vez de un genérico
"`TestTimeoutEfectivo` falló", el reporte dice exactamente cuál de los tres casos fue:

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

Esa salida es real: la corrí sobre el código de este mismo proyecto antes de escribir esta línea.
**Agregar un caso nuevo a la tabla es agregar una línea al slice `casos`** — no una función nueva, no
repetir el cuerpo de la prueba. Es la forma idiomática de probar una función con muchas entradas en Go, y
la vas a usar en cada paquete de aquí en adelante.

### 5.5 Probar los errores, no solo los aciertos

Una tabla de casos también sirve para probar que algo **falla como debe**, no solo que funciona (versión
abreviada aquí, con 3 de los 6 casos reales y struct literal posicional en vez de con nombre de campo,
para que quepa; la completa vive en `internal/config/config_test.go`):

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

**`t.Fatalf` en vez de `t.Errorf`** en el primer `if`: si `Interpretar` no devolvió error cuando debía,
seguir revisando `err.Error()` en la siguiente línea haría panic (`err` sería `nil`). `Fatalf` detiene
esa prueba en particular ahí mismo; `Errorf` deja que la prueba siga corriendo y acumule más fallas antes
de reportar. La regla práctica: usa `Fatalf` cuando seguir no tiene sentido sin lo que acabas de
comprobar, `Errorf` cuando sí lo tiene.

### 5.6 Los comandos que vas a usar siempre

```bash
go test ./...                          # todo el proyecto
go test ./internal/config/... -v       # verboso, un paquete
go test ./... -run TestInterpretar     # solo las pruebas cuyo nombre haga match
go test ./... -race                    # detector de carreras (se explica a fondo en la lección 6)
go test ./... -cover                   # porcentaje de líneas ejercitadas por las pruebas
```

Corridos de verdad sobre el `revisor`, el 30-sep-2026:

```
$ go test ./internal/servicio/... ./internal/config/... -cover
ok  	github.com/habil/revisor/internal/servicio	0.195s	coverage: 92.3% of statements
ok  	github.com/habil/revisor/internal/config	0.192s	coverage: 97.4% of statements
```

### 5.7 Qué hacer con un número de cobertura

**92.3% y 97.4% no son metas, son puntos de partida para una pregunta: ¿qué es ese 8% y ese 3% que no se
ejecutó, y me importa?** Con `go test -coverprofile` puedes ver exactamente qué líneas quedaron sin
tocar:

```bash
go test ./internal/servicio/... -coverprofile=/tmp/cobertura.out
go tool cover -func=/tmp/cobertura.out
```

En el `revisor`, el hueco de `servicio` es la rama de `Motivo()` que arma el mensaje cuando el código no
es ni 0 ni un éxito con un texto específico (`fmt.Sprintf("codigo %d", ...)` para un 404 sin más
contexto) — una rama que las pruebas existentes no ejercitan con ese código exacto. **La decisión correcta
no es perseguir el 100%** llenando cada rama con una prueba forzada que no enseña nada nuevo: es mirar el
hueco, decidir si importa (aquí, un poco: agregarías un caso con 404 en la tabla) y anotarlo, en vez de
fingir que no existe.

🔴 **Y una trampa real, medida en este mismo proyecto: la cobertura por paquete puede subestimar una
función central sin avisarte.** Corrido solo, el paquete `revisar` del `revisor` (que vas a conocer a
fondo en la lección 6) reporta:

```
$ go test ./internal/revisar/... -cover
ok  	github.com/habil/revisor/internal/revisar	1.33s	coverage: 61.9% of statements
```

**61.9% suena a que más de un tercio de ese paquete nunca se ejecuta en ninguna prueba — y es falso.** Su
función más importante, `Revisar` (la que de verdad habla HTTP), sí está bien probada: solo que la prueba
que la ejercita de verdad —`TestEjecutar_reportaOKyFalla`, con un servidor `httptest` real, en la lección
7— vive en el paquete `cmd/revisor`, no en `internal/revisar`. `go test ./internal/revisar/...` **solo
cuenta lo que las pruebas DE ESE PAQUETE ejercitan**; no ve lo que una prueba de otro paquete recorre por
el camino, aunque pase por el mismo código. Con `-coverpkg`, que le pide a Go medir la cobertura de un
paquete contando **todas** las pruebas del proyecto, no solo las suyas:

```
$ go test ./... -coverpkg=./... -coverprofile=/tmp/cov.out
$ go tool cover -func=/tmp/cov.out | grep 'revisar.go.*Revisar'
github.com/habil/revisor/internal/revisar/revisar.go:48:  Revisar    86.4%

$ go tool cover -func=/tmp/cov.out | tail -1
total:                                                    (statements)    81.5%
```

**86.4% para `Revisar`, 81.5% para el proyecto completo — no 61.9%.** La lección no es "ignora el número
por paquete": es que un número de cobertura siempre responde a una pregunta implícita —¿cobertura de qué,
medida contra las pruebas de dónde?— y `go test ./paquete/... -cover` calla esa segunda mitad de la
pregunta. Antes de decidir que algo "no está probado" por un número bajo, corre `-coverpkg=./...` sobre
todo el proyecto y compara.

### 5.8 Provocar una falla, para saber que la prueba sirve

Una prueba que nunca has visto fallar es una prueba de la que no sabes si funciona — puede estar
comparando dos cosas que siempre son iguales por accidente. Rómpela a propósito: cambia
`TimeoutEfectivo()` para que devuelva `s.Timeout` siempre, sin el `if`, y corre la prueba:

```
$ go test ./internal/servicio/... -run TestTimeoutEfectivo -v
=== RUN   TestTimeoutEfectivo/cero_(valor_por_omision)
    servicio_test.go:64: TimeoutEfectivo() = 0s, quería 2s
--- FAIL: TestTimeoutEfectivo (0.00s)
    --- FAIL: TestTimeoutEfectivo/cero_(valor_por_omision) (0.00s)
FAIL
```

**Ese `FAIL`, con el caso exacto que se rompió y el valor que obtuvo contra el que esperaba, es la prueba
de que la prueba funciona.** Devuelve el `if` y confirma que vuelve a estar en verde antes de seguir.

### 5.9 Diseñar para poder probar: separar la lógica de la E/S

Fíjate en algo que ya mencionamos de pasada en la sección 5.3 y que merece su propio espacio: `config`
tiene **dos** funciones, no una.

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

`Interpretar` (la firma nada más; el cuerpo completo, con toda la lógica de parseo, está en la sección
5.3, con sus errores reales):

<!-- verificar:fragmento -->
```go
func Interpretar(datos []byte, origen string) ([]servicio.Servicio, error) {
	// ... toda la lógica de parseo, sin tocar el disco ...
}
```

**Si hubiera una sola función `Cargar(ruta string)` que leyera el archivo y parseara todo junto**, cada
prueba de "línea mal escrita", "nombre repetido" o "tiempo límite inválido" tendría que empezar creando
un archivo temporal en disco con `os.CreateTemp`, escribirle el contenido del caso, pasarle la ruta, y
borrarlo al final. Funciona, pero es lento (toca el sistema de archivos real) y ensucia la prueba con
código que no tiene nada que ver con lo que en realidad se está probando: **si tu configuración se
interpreta bien o mal**, no si sabes crear archivos temporales.

Separando **la parte que decide** (`Interpretar`, una función pura: mismos bytes de entrada, mismo
resultado siempre, sin acceso a nada externo) de **la parte que obtiene los bytes** (`Cargar`, la única
que toca disco), las nueve pruebas de la sección 5.5 corren en microsegundos y sin crear un solo archivo.
`Cargar` en sí casi no necesita pruebas propias: solo comprobar que devuelve error si el archivo no
existe (ejercicio 6), porque toda la lógica interesante ya está en `Interpretar` y ya está probada.

🔑 **La regla general, útil mucho más allá de este proyecto:** cuando una función es difícil de probar,
casi siempre es porque mezcla "decidir algo" con "tocar el mundo exterior" (un archivo, la red, el
reloj). Separarlas no es una regla de estilo: es lo que determina si vas a poder escribir la prueba en
tres líneas o en veinte.

### 5.10 Benchmarks: medir, no adivinar

Además de `Test...`, Go reconoce funciones `Benchmark...` que miden cuánto tarda tu código, no si es
correcto:

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

`b.N` no lo eliges tú: Go corre el ciclo con valores de `N` cada vez más grandes hasta que la medición es
estable, y reporta el tiempo por operación. Corrido de verdad sobre el `revisor` (Apple M5, 200,000
repeticiones):

```
$ go test ./internal/config/... -bench=. -run '^$'
goos: darwin
goarch: arm64
pkg: github.com/habil/revisor/internal/config
cpu: Apple M5
BenchmarkInterpretar-10    	  200000	       432.5 ns/op
PASS
```

**`-run '^$'`** le dice a `go test` que no corra ninguna prueba normal (una expresión regular que no
hace match con ningún nombre), para que el reporte del benchmark no se mezcle con el de las pruebas. El
número —432.5 nanosegundos por llamada, en esta máquina, en este momento— no es para memorizarlo: es para
compararlo **contra sí mismo** después de un cambio. Si mañana reescribes `Interpretar` y el benchmark
sube a 4,000 ns/op, tienes una señal objetiva de que algo se puso más lento, sin necesitar opinar al
respecto.

### 5.11 `go vet`: el que encuentra lo que compila pero está mal

`go test` te dice si tu lógica hace lo que esperabas. **`go vet` te dice si tu código tiene un error que
el compilador no atrapa porque, técnicamente, es válido** — pero casi seguro no es lo que quisiste
escribir. El caso más común es un verbo de formato (lección 2) que no corresponde al tipo del argumento:

<!-- verificar:ejemplo:ejemplos/05-vet-printf -->
```go
puerto := 443
fmt.Printf("servicio %s en el puerto %s\n", nombre, puerto) // %s para un int
```

Esto **compila** y **corre**, sin panic ni error — y produce una salida rota (programa completo en
`programas/revisor/ejemplos/05-vet-printf/main.go`):

```
$ go run ./ejemplos/05-vet-printf/
servicio catalogo en el puerto %!s(int=443)
```

Y `go vet`, sobre ese mismo archivo, sí lo detecta:

```
$ go vet ./ejemplos/05-vet-printf/
ejemplos/05-vet-printf/main.go:11:39: fmt.Printf format %s has arg puerto of wrong type int
```

`go vet` sí lo detecta, porque analiza la cadena de formato contra los tipos reales de los argumentos, un
paso que el compilador de Go no da por sí solo. Corre `go vet ./...` junto con `go test ./...` como
rutina: la mayoría de los editores con la extensión de Go (lección 1) ya lo hacen por ti mientras
escribes, subrayando el problema antes de que llegues a ejecutar nada.

### 5.12 Pruebas de frontera: donde de verdad viven los bugs

`Estado.OK()` decide que un código está bien si cae **entre 200 y 299**. Es tentador probarlo con un
caso "obvio" (200) y un caso "obvio" de falla (500) y darlo por bueno. La tabla real del `revisor`
prueba también los dos valores que están **justo en el borde del rango**:

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

**¿Por qué importa, si 199 y 300 "obviamente" no están OK?** Porque un error de un solo carácter en la
condición —`>=` en vez de `>`, o `<=` en vez de `<`— es exactamente el tipo de bug que un caso "obvio"
nunca detecta, y que un caso de frontera detecta siempre. Si alguien cambiara `OK()` a
`e.Codigo >= 200 && e.Codigo <= 300` (incluyendo el 300 por error), los casos 200/299/404/500 seguirían
pasando igual de bien — **solo el caso de 300 lo delataría.** Ésa es la razón de fondo detrás de "prueba
los bordes, no solo el centro": los bordes son donde los errores de comparación se esconden, y son
invisibles para cualquier prueba que solo use valores muy adentro o muy afuera del rango.

---

## El error que vas a ver

| El síntoma | Mensaje literal | Qué pasa y qué hacer |
|---|---|---|
| Import de un paquete `internal` ajeno | `use of internal package github.com/habil/revisor/internal/servicio not allowed` | El compilador impide importar algo bajo `internal/` desde fuera del módulo. No es un permiso que puedas dar: hay que exponer el tipo desde un paquete público si de verdad hace falta |
| Servicio duplicado en la configuración | `prueba.txt:2: el nombre "catalogo" ya estaba en la linea 1` | Dos servicios con el mismo nombre; corrige el archivo de configuración |
| URL sin esquema | `prueba.txt:1: la URL "catalogo.interno.mx" debe empezar con http:// o https://` | Falta `http://` o `https://` al inicio de la URL |
| Tiempo límite mal escrito | `prueba.txt:1: tiempo limite "nombas" invalido (usa 500ms, 2s, 1m): time: invalid duration "nombas"` | El formato de duración de Go no es libre: usa un número seguido de una unidad (`ms`, `s`, `m`, `h`) |
| Una prueba que rompiste a propósito | `servicio_test.go:64: TimeoutEfectivo() = 0s, quería 2s` | Ver sección 5.8: así se ve un `t.Errorf` señalando exactamente qué caso falló y por qué |
| Cero pruebas en un paquete | `?   	github.com/habil/revisor/cmd/servidor-demo	[no test files]` | No es una falla: `go test` avisa explícitamente cuando un paquete no tiene ningún archivo `_test.go`, en vez de fingir que corrió algo |

**Y el más importante de aprender a leer**, porque no es un error sino la ausencia de uno:

```
$ go test ./...
ok  	github.com/habil/revisor/internal/vacio	0.001s
```

Si `internal/vacio` no tuviera **ninguna** función `Test...`, esa línea se ve exactamente igual:
`ok`, en verde, sin ninguna marca de que no se ejecutó nada. La única forma de distinguir "todo pasó, en
serio" de "no había nada que correr" es mirar el conteo con `-v` (que sí imprime cada `RUN`) o, mejor,
nunca confiar en un paquete que no tiene ningún archivo `_test.go` — eso `go test` sí lo dice, como en la
fila de la tabla de arriba.

---

## Lo que se hace mal

- **Crear un paquete `utils`, `helpers` o `common`.** Ya lo viste en la sección 5.2: es el antipatrón más
  repetido y el que más rápido pierde su propósito. Nombra la responsabilidad, no el hecho de que "no
  encaja en otro lado".
- **Confundir "compiló" con "las pruebas pasaron".** `go build` verifica que el código es válido; no
  ejecuta ni una sola prueba. Son dos comandos distintos con dos preguntas distintas.
- **Leer el código de salida de `go test` en vez del conteo.** Un paquete sin archivos de prueba sale con
  código 0 (éxito), igual que uno con 50 pruebas que sí pasaron. El código de salida responde "¿algo
  falló?", no "¿algo se probó?" — para la segunda pregunta hace falta leer la salida, no solo el `$?`.
- **Perseguir el 100% de cobertura como si fuera el objetivo.** Una cobertura alta con aserciones débiles
  (comprobar que algo no truena, sin comprobar qué devuelve) da una cifra bonita y una prueba que no
  detecta casi nada. El número es una guía de dónde mirar, no una meta en sí misma — la sección 5.7 lo
  muestra con un hueco real del propio `revisor`.
- **Escribir una prueba y nunca verla fallar.** Si nunca rompiste el código a propósito para confirmar
  que la prueba se pone roja, no sabes si esa prueba prueba algo. La sección 5.8 lo hizo con datos reales
  del proyecto: hazlo tú también con al menos una prueba propia antes de darla por buena.

---

## Ejercicios

1. Clona la estructura de paquetes de la sección 5.2 (`internal/servicio`, `internal/config`) para tu
   propia copia del `revisor`, moviendo el código que ya tenías de las lecciones 2 a 4.
2. Escribe la tabla de casos de `TestEtiqueta` para el método `Etiqueta()` de `Servicio`, con al menos un
   caso de nombre normal y uno de struct vacío.
3. Agrega un caso a `TestInterpretar_casosDeError` para una línea con **cuatro** campos (más de los tres
   que el formato permite). Verifica el mensaje exacto contra el código de `config.go`.
4. Corre `go test ./... -cover` sobre tu copia y anota el porcentaje de cada paquete. Elige **un** hueco
   de cobertura y decide, por escrito en tu bitácora, si te importa cerrarlo y por qué.
5. (Como en la sección 5.8) Rompe a propósito una función que ya probaste, corre la prueba, lee el
   `FAIL` completo, y repáralo. Pega los dos resultados —el rojo y el verde— en tu bitácora.
6. (Un poco más difícil) Escribe `TestCargar_archivoInexistente`, que confirme que `Cargar` (no
   `Interpretar`) devuelve error cuando la ruta no existe. Pista: no necesitas crear ningún archivo para
   esta prueba, solo pasar una ruta que sabes que no existe.

### Soluciones

1 y 2 no tienen solución de referencia única: depende de cómo tenías organizado tu propio programa de la
lección 4. Compara tu resultado contra el código real de `programas/revisor/internal/servicio/servicio.go` del
proyecto de este curso.

3. Con `"catalogo https://a.mx 500ms extra\n"`, el mensaje esperado es
   `prueba.txt:1: esperaba «nombre url [tiempo]», hay 4 campo(s)` — el mismo camino de código que ya
   maneja "faltan campos", porque `len(campos) > 3` cubre ambos casos con una sola comprobación.

4. No hay una respuesta única: lo que importa es que la decisión quede escrita con su razón, no el
   porcentaje en sí.

5. Ver la sección 5.8 completa: el patrón es siempre "cambia el código, corre la prueba, lee el `FAIL`
   con el caso exacto, repara, corre otra vez".

6. Así:

   <!-- verificar:extracto:internal/config/config_test.go -->
   ```go
   func TestCargar_archivoInexistente(t *testing.T) {
   	_, err := Cargar("/ruta/que/no/existe.txt")
   	if err == nil {
   		t.Fatal("Cargar() no devolvió error con una ruta inexistente")
   	}
   }
   ```

   Esta prueba sí existe en el proyecto real (`config_test.go`) y pasa porque `os.ReadFile`, dentro de
   `Cargar`, devuelve un error de sistema operativo que `Cargar` envuelve con `%w` antes de propagarlo.

---

## Cómo sé que lo logré

- [ ] Mi `revisor` está organizado en paquetes bajo `internal/`, cada uno con una sola responsabilidad.
- [ ] `go test ./...` corre y **el conteo, no solo el color**, me dice cuántas pruebas se ejecutaron.
- [ ] Escribí al menos una tabla de casos con `t.Run`, y sé leer cuál caso falló cuando algo se rompe.
- [ ] Corrí `go test -cover` y puedo decir qué porcentaje salió y qué hay en el hueco.
- [ ] Rompí una función a propósito, vi el `FAIL` exacto, y la reparé — lo tengo en mi bitácora.
- [ ] Sé explicar por qué el `revisor` no tiene `go.sum` y qué lo generaría si algún día lo necesitara.
- [ ] Sé por qué no debería crear un paquete `utils`.

---

## Para leer más

1. [Writing tests](https://go.dev/doc/tutorial/add-a-test) — el tutorial oficial de pruebas, con tabla de
   casos incluida.
2. [Go Wiki: TableDrivenTests](https://go.dev/wiki/TableDrivenTests) — la referencia canónica del patrón
   usado en toda esta lección.
3. [Package internal](https://go.dev/doc/go1.4#internalpackages) — el anuncio original de la regla de
   `internal/`, directo de las notas de versión de Go.
4. [Go Blog: The Cover Story](https://go.dev/blog/cover) — cómo funciona `go test -cover` por dentro, y
   por qué un número alto no siempre significa pruebas buenas.
5. [Liberar sin miedo en su propia infraestructura](https://www.habil.mx/es/blog/cicd-devsecops-infraestructura-propia/) — artículo sobre un camino a producción donde las pruebas, entre otras compuertas, pueden detener una liberación.

### Términos de esta lección

| Término | Qué significa |
|---|---|
| `go.sum` | archivo con las huellas criptográficas de las dependencias externas; el `revisor` no lo tiene porque no usa ninguna |
| `internal/` | carpeta especial que el compilador de Go impide importar desde fuera del módulo |
| tabla de casos | patrón de prueba donde una lista de entradas y salidas esperadas se recorre con un solo cuerpo de prueba |
| `t.Run` | ejecuta un subcaso con su propio nombre, para que el reporte diga exactamente cuál falló |
| cobertura | porcentaje de líneas del código que las pruebas ejercitaron al correr |

---

**Anterior:** [Lección 4 — Colecciones](04-colecciones.md) ·
**Siguiente:** [Lección 6 — Concurrencia](06-concurrencia.md)
