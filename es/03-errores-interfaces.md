# Lección 3 — Structs, métodos, errores e interfaces

> **Aquí Go se vuelve Go.** Ésta es la lección más importante del curso.

**Duración:** dos sesiones de 60 minutos. No la hagas de corrido.

**Al terminar vas a poder:**

- Agrupar datos relacionados en un **struct** y explicar por qué es mejor que variables sueltas.
- Escribir **métodos** y decidir cuándo usar receptor de valor y cuándo de puntero.
- Explicar qué es un **puntero** sin asustarte.
- Manejar errores como valores, envolverlos y preguntar por ellos.
- Definir una **interfaz** y entender por qué en Go no se declara que se cumple.
- Escribir código que se puede probar **sin conectarse a nada**.

---

## Por qué importa

En la lección 2 guardaste los datos de un servicio en tres variables separadas (`nombre`, `url`, `puerto`).
Funciona con **un** servicio. Con diez son treinta variables sueltas, y nada en el código dice cuáles van
juntas — si te equivocas y combinas el nombre de uno con el puerto de otro, **el programa compila igual** y
produce un reporte incorrecto. Ese es el primer problema que resuelve esta lección: el **struct**.

El segundo problema es más de fondo. El `revisor` va a tener que consultar servicios de verdad, y las
consultas de verdad **fallan**: se cae la red, el servicio tarda, responde con un código que no esperabas.
Un programa que no sabe manejar eso con orden no sirve para nada en producción. Go resuelve esto de una
forma que probablemente no hayas visto si vienes de otro lenguaje: **los errores son valores**, no
excepciones que interrumpen el flujo.

Y el tercer problema es el que hace que las dos piezas anteriores encajen entre sí sin que el `revisor`
tenga que saber, de antemano, **cómo** se va a revisar cada servicio (HTTP hoy, quizá una base de datos
mañana). Eso lo resuelve la **interfaz**, que es también lo que te va a permitir, en la lección 5, probar
todo el programa **sin conectarte a nada real**.

Structs, errores e interfaces son, en ese orden, la columna vertebral de todo lo que sigue en el curso.

---

## Los conceptos

### 3.1 Structs: agrupar lo que va junto

Lo que necesitas es decirle al lenguaje: *«estas tres cosas son un servicio»*.

**Fig. 3.1** | Definir y usar un struct.

```go
 1  // fig03_01.go
 2  // Agrupa los datos de un servicio en un solo tipo.
 3  package main
 4
 5  import "fmt"
 6
 7  // Servicio agrupa todo lo que describe a un servicio que vamos a revisar.
 8  type Servicio struct {
 9      Nombre string
10      URL    string
11      Puerto int
12  }
13
14  func main() {
15      s := Servicio{
16          Nombre: "catalogo",
17          URL:    "https://catalogo.example.com",
18          Puerto: 443,
19      }
20
21      fmt.Println("nombre:", s.Nombre)
22      fmt.Println("url:", s.URL)
23      fmt.Println("puerto:", s.Puerto)
24
25      s.Puerto = 8443          // se puede modificar
26      fmt.Println("nuevo puerto:", s.Puerto)
27
28      fmt.Println(s)           // e imprimir completo
29  }
```

```bash
$ go run fig03_01.go
nombre: catalogo
url: https://catalogo.example.com
puerto: 443
nuevo puerto: 8443
{catalogo https://catalogo.example.com 8443}
```

**Línea 8: `type Servicio struct {`.** Se lee: *«define un tipo nuevo llamado `Servicio`, que es una
estructura»*. **Acabas de crear un tipo que no venía con el lenguaje**, y desde ahora se usa igual que
`int` o `string`.

**Líneas 15-19.** Crea un valor de ese tipo. Los nombres de campo con dos puntos son opcionales —podrías
escribir solo los valores en orden— pero **ponlos siempre**:

> [!TIP]
> ✅ **Buena práctica 3.1**
> Escribe siempre los nombres de los campos al crear un struct: `Servicio{Nombre: "x", Puerto: 443}` en
> vez de `Servicio{"x", "", 443}`. La versión corta se rompe en silencio si alguien agrega un campo o
> cambia el orden, y el compilador no puede avisarte porque los tipos siguen cuadrando.

**Línea 21: `s.Nombre`.** El punto accede a un campo. Se lee «el `Nombre` de `s`».

#### 3.1.1 La mayúscula es permiso, no estética

Fíjate que los campos empiezan con **mayúscula**: `Nombre`, `URL`, `Puerto`. En Go eso no es un estilo
elegido: **es el control de acceso del lenguaje.**

| | |
|---|---|
| `Nombre` (mayúscula) | **exportado**: visible desde otros paquetes |
| `nombre` (minúscula) | **no exportado**: solo visible dentro de su propio paquete |

No hay `public`, `private` ni `protected`. **La letra inicial es la regla completa.**

#### 3.1.2 El valor cero de un struct

**Fig. 3.2** | Un struct sin inicializar.

```go
 1  // fig03_02.go
 2  package main
 3
 4  import "fmt"
 5
 6  type Servicio struct {
 7      Nombre string
 8      URL    string
 9      Puerto int
10  }
11
12  func main() {
13      var vacio Servicio
14      fmt.Printf("%+v\n", vacio)
15      fmt.Println("¿el nombre está vacío?", vacio.Nombre == "")
16      fmt.Println("puerto:", vacio.Puerto)
17  }
```

```bash
$ go run fig03_02.go
{Nombre: URL: Puerto:0}
¿el nombre está vacío? true
puerto: 0
```

Cada campo recibe **su** valor cero: los textos quedan en `""` y los números en `0`. **No hay basura ni
«indefinido»**, y por eso un struct recién creado ya se puede usar sin miedo.

**Y aquí aparece `%+v`**, que es muy útil para depurar: imprime el struct **con los nombres de los
campos**. Compáralo con `%v`, que imprime solo los valores.

> [!TIP]
> 🧪 **Tip de prueba y depuración 3.1**
> Cuando no entiendas qué tiene un struct, imprímelo con `%+v`. Es la forma más rápida de ver todos sus
> campos con su nombre, y te va a ahorrar muchísimo tiempo.

### 3.2 Métodos

Un **método** es una función que pertenece a un tipo. En vez de escribir `etiqueta(s)`, escribes
`s.Etiqueta()`, y la función queda ligada al tipo al que corresponde.

**Fig. 3.3** | Métodos sobre un struct.

```go
 1  // fig03_03.go
 2  // Define metodos que pertenecen al tipo Servicio.
 3  package main
 4
 5  import "fmt"
 6
 7  type Servicio struct {
 8      Nombre string
 9      URL    string
10      Puerto int
11  }
12
13  // Etiqueta devuelve una descripcion legible del servicio.
14  func (s Servicio) Etiqueta() string {
15      return fmt.Sprintf("%s (%s:%d)", s.Nombre, s.URL, s.Puerto)
16  }
17
18  // EsSeguro indica si el servicio usa el puerto de HTTPS.
19  func (s Servicio) EsSeguro() bool {
20      return s.Puerto == 443
21  }
22
23  func main() {
24      a := Servicio{Nombre: "catalogo", URL: "https://catalogo.example.com", Puerto: 443}
25      b := Servicio{Nombre: "local", URL: "http://localhost", Puerto: 8080}
26
27      fmt.Println(a.Etiqueta(), "— seguro:", a.EsSeguro())
28      fmt.Println(b.Etiqueta(), "— seguro:", b.EsSeguro())
29  }
```

```bash
$ go run fig03_03.go
catalogo (https://catalogo.example.com:443) — seguro: true
local (http://localhost:8080) — seguro: false
```

**Línea 14: `func (s Servicio) Etiqueta() string`.** Ese `(s Servicio)` entre `func` y el nombre se llama
**receptor**, y es lo que convierte una función en un método. Se lee: *«esta función pertenece al tipo
`Servicio`, y dentro me voy a referir al valor como `s`»*.

> [!NOTE]
> 🔧 **Observación de ingeniería de software 3.1**
> Si vienes de Java, C# o Python, un método de Go es parecido a un método de clase — con una diferencia
> importante: **en Go el método se escribe fuera del tipo.** Puedes tener el `type` en un archivo y sus
> métodos en otro. Y como no hay clases, tampoco hay herencia: en Go se reutiliza código de otra forma, y
> la vas a ver en la sección 3.6.

### 3.3 Punteros, en diez minutos

Los punteros tienen fama de difíciles. En Go son mucho más simples que en C, y necesitas entender **una
sola idea** para lo que viene.

#### 3.3.1 El problema

**Fig. 3.4** | Un método que **no funciona** como esperarías.

```go
 1  // fig03_04.go
 2  // Muestra por que un receptor de VALOR no puede modificar el original.
 3  package main
 4
 5  import "fmt"
 6
 7  type Servicio struct {
 8      Nombre string
 9      Puerto int
10  }
11
12  // CambiarPuerto intenta modificar el servicio... y no lo logra.
13  func (s Servicio) CambiarPuerto(nuevo int) {
14      s.Puerto = nuevo
15  }
16
17  func main() {
18      x := Servicio{Nombre: "catalogo", Puerto: 443}
19      fmt.Println("antes:", x.Puerto)
20
21      x.CambiarPuerto(8443)
22      fmt.Println("despues:", x.Puerto)      // ¿cambió?
23  }
24
```

```bash
$ go run fig03_04.go
antes: 443
despues: 443
```

**No cambió nada.** El método se ejecutó, no hubo error, y el valor sigue igual.

**¿Por qué?** Porque cuando el receptor es `(s Servicio)`, Go le entrega al método **una copia** del
struct. El método modifica la copia, la copia se descarta al terminar, y el original nunca se enteró. Es
un antipatrón tan frecuente que tiene su propia entrada en «Lo que se hace mal», más abajo.

#### 3.3.2 La solución: el receptor de puntero

Un **puntero** es una variable que guarda **la dirección** de otra, en vez de una copia de su contenido.

Piensa en la diferencia entre darte **una fotocopia** de un documento y darte **la dirección del archivero
donde está el original**. Con la fotocopia puedes escribir todo lo que quieras: el original no cambia. Con
la dirección, puedes ir y modificar el original.

- `Servicio` es la fotocopia.
- `*Servicio` es la dirección del original.

**Fig. 3.5** | El mismo método, ahora con receptor de puntero.

```go
 1  // fig03_05.go
 2  // Con receptor de PUNTERO, el metodo si modifica el original.
 3  package main
 4
 5  import "fmt"
 6
 7  type Servicio struct {
 8      Nombre string
 9      Puerto int
10  }
11
12  // CambiarPuerto modifica el servicio. El * es la diferencia.
13  func (s *Servicio) CambiarPuerto(nuevo int) {
14      s.Puerto = nuevo
15  }
16
17  func main() {
18      x := Servicio{Nombre: "catalogo", Puerto: 443}
19      fmt.Println("antes:", x.Puerto)
20
21      x.CambiarPuerto(8443)
22      fmt.Println("despues:", x.Puerto)
23  }
```

```bash
$ go run fig03_05.go
antes: 443
despues: 8443
```

**La única diferencia entre la Fig. 3.4 y la 3.5 es un asterisco en la línea 13.** Eso es todo.

Y fíjate en la línea 21: escribiste `x.CambiarPuerto(8443)` igual que antes, **sin `&` ni nada raro**. Go
se da cuenta de que el método necesita un puntero y lo toma solo. Ésa es la razón por la que los punteros
de Go dan mucho menos trabajo que los de C.

> [!TIP]
> ✅ **Buena práctica 3.2**
> La regla es sencilla: **si el método modifica el struct, receptor de puntero (`*T`); si solo lee, de
> valor (`T`)**. Y sé consistente dentro de un mismo tipo: si la mayoría de tus métodos necesitan puntero,
> usa puntero en todos, aunque algunos no lo requieran. Mezclarlos confunde a quien lea el código.

> [!NOTE]
> 🚀 **Tip de rendimiento 3.1**
> Hay una segunda razón para usar punteros: **evitar copiar**. Si un struct tiene veinte campos, cada
> llamada con receptor de valor copia los veinte. Para structs pequeños es irrelevante; para grandes o en
> ciclos de millones de vueltas, importa. **No optimices esto sin medir:** la claridad vale más que una
> copia de 40 bytes.

#### 3.3.3 El único puntero peligroso: `nil`

Todavía hay una tercera situación con punteros, y ésta sí produce un mensaje de error real —tan real que
tiene su propia sección dedicada: **«El error que vas a ver»**, más abajo en esta lección. Adelantamos la
idea: el valor cero de un puntero es **`nil`** («no apunto a nada»), y leer un campo a través de un
puntero nulo hace que el programa **truene** en vez de comportarse mal en silencio como en la Fig. 3.4.

### 3.4 Los errores son valores

Ya viste en la lección 2 que una función de Go puede devolver `(resultado, error)`. Ahora vamos al fondo,
porque **esto es lo que más distingue a Go de lo que probablemente hayas visto.**

#### 3.4.1 `error` es una interfaz, no algo mágico

En Go, `error` es simplemente un tipo con un método:

```go
type error interface {
    Error() string
}
```

Eso significa: *«cualquier cosa que tenga un método `Error()` que devuelva texto, es un error»*. No hay
jerarquía de clases de excepción, no hay `throw`, no hay pila de llamadas que se desenrolla. **Un error es
un valor común que viaja como cualquier otro.**

#### 3.4.2 Crear errores

**Fig. 3.7** | Las tres formas de crear un error.

```go
 1  // fig03_07.go
 2  // Muestra las tres formas de producir un error.
 3  package main
 4
 5  import (
 6      "errors"
 7      "fmt"
 8  )
 9
10  // ErrPuertoInvalido es un error CENTINELA: se declara una vez y se compara.
11  var ErrPuertoInvalido = errors.New("el puerto debe estar entre 1 y 65535")
12
13  type Servicio struct {
14      Nombre string
15      Puerto int
16  }
17
18  // Validar revisa que el servicio tenga sentido.
19  func Validar(s Servicio) error {
20      if s.Nombre == "" {
21          // 1. un error sencillo, creado al momento
22          return errors.New("el nombre no puede estar vacio")
23      }
24      if s.Puerto < 1 || s.Puerto > 65535 {
25          // 2. un error centinela, para poder compararlo despues
26          return ErrPuertoInvalido
27      }
28      if s.Puerto == 80 {
29          // 3. un error con datos dentro
30          return fmt.Errorf("el puerto %d no usa cifrado, usa 443", s.Puerto)
31      }
32      return nil          // nil significa: todo bien
33  }
34
35  func main() {
36      casos := []Servicio{
37          {Nombre: "catalogo", Puerto: 443},
38          {Nombre: "", Puerto: 443},
39          {Nombre: "pagos", Puerto: 99999},
40          {Nombre: "viejo", Puerto: 80},
41      }
42
43      for _, s := range casos {
44          if err := Validar(s); err != nil {
45              fmt.Printf("%-10s ❌ %v\n", s.Nombre, err)
46          } else {
47              fmt.Printf("%-10s ✅ valido\n", s.Nombre)
48          }
49      }
50  }
```

```bash
$ go run fig03_07.go
catalogo   ✅ valido
           ❌ el nombre no puede estar vacio
pagos      ❌ el puerto debe estar entre 1 y 65535
viejo      ❌ el puerto 80 no usa cifrado, usa 443
```

**Línea 11: el error centinela.** Se declara **una vez**, a nivel de paquete, con el prefijo `Err`. Sirve
para que quien llame a tu función pueda preguntar *«¿fue este error en particular?»*, como verás en 3.5.

**Línea 30: `fmt.Errorf`.** Como `Printf`, pero produce un error en vez de imprimir. Úsalo cuando el
mensaje necesita datos.

**Línea 44: `if err := Validar(s); err != nil`.** Este patrón declara `err` **dentro** del `if`, así que
solo existe ahí. Es muy idiomático en Go y mantiene el código limpio.

> [!TIP]
> ✅ **Buena práctica 3.3**
> Los mensajes de error se escriben **en minúscula y sin punto final**: `"no se pudo abrir el archivo"`,
> no `"No se pudo abrir el archivo."`. La razón es práctica: los errores se **envuelven** unos en otros
> (sección 3.5), y al concatenarse quedan como `"revisando catalogo: no se pudo abrir el archivo"`. Con
> mayúsculas y puntos, el resultado se vería roto.

### 3.5 Envolver errores: `%w`, `errors.Is` y `errors.As`

Un error sin contexto es poco útil. Si tu programa dice `connection refused`, no sabes **qué** servicio
falló.

**Envolver** un error es agregarle contexto **sin perder el original**.

**Fig. 3.8** | Envolver errores y preguntar por ellos.

```go
 1  // fig03_08.go
 2  // Envuelve un error para agregar contexto sin perder el original.
 3  package main
 4
 5  import (
 6      "errors"
 7      "fmt"
 8  )
 9
10  var ErrNoResponde = errors.New("no responde")
11
12  type Servicio struct{ Nombre string }
13
14  // consultar simula la consulta de bajo nivel.
15  func consultar(s Servicio) error {
16      return ErrNoResponde
17  }
18
19  // Revisar agrega contexto al error de consultar.
20  func Revisar(s Servicio) error {
21      if err := consultar(s); err != nil {
22          // el %w ENVUELVE el error original
23          return fmt.Errorf("revisando %s: %w", s.Nombre, err)
24      }
25      return nil
26  }
27
28  func main() {
29      err := Revisar(Servicio{Nombre: "catalogo"})
30
31      fmt.Println("1. el mensaje completo:")
32      fmt.Println("  ", err)
33
34      fmt.Println("2. ¿es un ErrNoResponde, aunque este envuelto?")
35      fmt.Println("  ", errors.Is(err, ErrNoResponde))
36
37      fmt.Println("3. ¿y si pregunto por otro error?")
38      fmt.Println("  ", errors.Is(err, errors.New("otra cosa")))
39
40      fmt.Println("4. el error original, desenvuelto:")
41      fmt.Println("  ", errors.Unwrap(err))
42  }
```

```bash
$ go run fig03_08.go
1. el mensaje completo:
   revisando catalogo: no responde
2. ¿es un ErrNoResponde, aunque este envuelto?
   true
3. ¿y si pregunto por otro error?
   false
4. el error original, desenvuelto:
   no responde
```

**Línea 23: `%w`.** Éste es el verbo clave. Se ve igual que `%v` al imprimir, **pero conserva el error
original adentro** para que se pueda preguntar por él después.

**Línea 35: `errors.Is`.** Pregunta *«¿en algún punto de esta cadena está ese error?»*. Funciona aunque
haya cinco capas de envoltura.

> [!WARNING]
> 🔴 **Error común de programación 3.3 — el más silencioso de esta lección**
> Escribir `%v` en vez de `%w` al envolver:
>
> ```go
> return fmt.Errorf("revisando %s: %v", s.Nombre, err)   // ❌ con %v
> ```
>
> **El mensaje impreso es idéntico.** No hay error, no hay aviso, todo parece funcionar. Pero
> `errors.Is` deja de encontrar el error original y devuelve `false`, así que el código que decidía qué
> hacer según el tipo de fallo empieza a tomar el camino equivocado. **Es el tipo de error que no falla:
> devuelve un dato peor.** Cuando envuelvas, usa `%w`. Este es el segundo antipatrón de la sección «Lo
> que se hace mal».

#### 3.5.1 `errors.As`, cuando necesitas los datos del error

`errors.Is` responde «¿es este error?». `errors.As` responde «¿es de este **tipo**? dámelo para leer sus
campos».

**Fig. 3.9** | Un error propio con datos, recuperado con `errors.As`.

```go
 1  // fig03_09.go
 2  // Define un tipo de error propio y recupera sus datos.
 3  package main
 4
 5  import (
 6      "errors"
 7      "fmt"
 8  )
 9
10  // ErrorHTTP es un error que lleva datos adentro.
11  type ErrorHTTP struct {
12      Codigo int
13      URL    string
14  }
15
16  // Error hace que ErrorHTTP cumpla la interfaz error.
17  func (e *ErrorHTTP) Error() string {
18      return fmt.Sprintf("el servidor respondio %d", e.Codigo)
19  }
20
21  func consultar(url string) error {
22      return &ErrorHTTP{Codigo: 503, URL: url}
23  }
24
25  func Revisar(nombre, url string) error {
26      if err := consultar(url); err != nil {
27          return fmt.Errorf("revisando %s: %w", nombre, err)
28      }
29      return nil
30  }
31
32  func main() {
33      err := Revisar("catalogo", "https://catalogo.example.com")
34      fmt.Println("mensaje:", err)
35
36      var errHTTP *ErrorHTTP
37      if errors.As(err, &errHTTP) {
38          fmt.Println("es un ErrorHTTP")
39          fmt.Println("  codigo:", errHTTP.Codigo)
40          fmt.Println("  url:   ", errHTTP.URL)
41
42          if errHTTP.Codigo >= 500 {
43              fmt.Println("  → es culpa del servidor, conviene reintentar")
44          }
45      }
46  }
```

```bash
$ go run fig03_09.go
mensaje: revisando catalogo: el servidor respondio 503
es un ErrorHTTP
  codigo: 503
  url:    https://catalogo.example.com
  → es culpa del servidor, conviene reintentar
```

**Línea 17.** Al escribir un método `Error() string`, tu tipo **ya es un `error`**. No declaraste nada: el
tipo cumple la interfaz porque tiene el método. Eso es lo que explica la sección siguiente.

**Línea 42.** Y aquí está el valor real: el programa puede **decidir** según el tipo de fallo. Un 503 se
reintenta; un 404, no. Con errores de solo texto eso sería imposible sin comparar cadenas, que es frágil.

### 3.6 Interfaces: el corazón de Go

Una **interfaz** es una lista de métodos. Cualquier tipo que tenga esos métodos **la cumple**, y no hay
que declararlo en ninguna parte.

**Fig. 3.10** | Una interfaz y dos implementaciones.

```go
 1  // fig03_10.go
 2  // Una interfaz con dos implementaciones: la real y una de prueba.
 3  package main
 4
 5  import (
 6      "errors"
 7      "fmt"
 8  )
 9
10  type Servicio struct {
11      Nombre string
12      URL    string
13  }
14
15  type Estado struct {
16      Servicio Servicio
17      Codigo   int
18      Err      error
19  }
20
21  // Revisor es una INTERFAZ: cualquier cosa con este metodo la cumple.
22  type Revisor interface {
23      Revisar(s Servicio) Estado
24  }
25
26  // ---- primera implementacion: la de verdad (simplificada) ----
27  type RevisorHTTP struct{}
28
29  func (r RevisorHTTP) Revisar(s Servicio) Estado {
30      // En la leccion 7 esto hara una peticion HTTP real.
31      return Estado{Servicio: s, Codigo: 200}
32  }
33
34  // ---- segunda implementacion: para probar, sin red ----
35  type RevisorFalso struct {
36      Respuesta Estado
37  }
38
39  func (r RevisorFalso) Revisar(s Servicio) Estado {
40      r.Respuesta.Servicio = s
41      return r.Respuesta
42  }
43
44  // RevisarTodos acepta CUALQUIER Revisor. No sabe ni le importa cual.
45  func RevisarTodos(r Revisor, servicios []Servicio) []Estado {
46      var estados []Estado
47      for _, s := range servicios {
48          estados = append(estados, r.Revisar(s))
49      }
50      return estados
51  }
52
53  func main() {
54      servicios := []Servicio{
55          {Nombre: "catalogo", URL: "https://catalogo.example.com"},
56          {Nombre: "pagos", URL: "https://pagos.example.com"},
57      }
58
59      fmt.Println("--- con el revisor real ---")
60      for _, e := range RevisarTodos(RevisorHTTP{}, servicios) {
61          fmt.Printf("  %-10s codigo %d\n", e.Servicio.Nombre, e.Codigo)
62      }
63
64      fmt.Println("--- con el falso, simulando una falla ---")
65      falso := RevisorFalso{
66          Respuesta: Estado{Err: errors.New("no responde")},
67      }
68      for _, e := range RevisarTodos(falso, servicios) {
69          fmt.Printf("  %-10s error: %v\n", e.Servicio.Nombre, e.Err)
70      }
71  }
```

```bash
$ go run fig03_10.go
--- con el revisor real ---
  catalogo   codigo 200
  pagos      codigo 200
--- con el falso, simulando una falla ---
  catalogo   error: no responde
  pagos      error: no responde
```

**Léelo despacio, porque aquí está la idea que sostiene todo Go.**

**Línea 22-24: la interfaz.** Dice: *«un `Revisor` es cualquier cosa que tenga un método `Revisar` que
reciba un `Servicio` y devuelva un `Estado`»*.

**Líneas 29 y 39.** `RevisorHTTP` y `RevisorFalso` **no declaran en ninguna parte** que cumplen `Revisor`.
No hay `implements`, no hay `: Revisor`, nada. **Cumplen la interfaz porque tienen el método**, y eso lo
verifica el compilador solo.

**Línea 45: `func RevisarTodos(r Revisor, ...)`.** Esta función **no sabe** con qué está trabajando. Solo
sabe que puede llamar `.Revisar()`. Y por eso las líneas 60 y 68 le pasan dos cosas completamente
distintas **sin cambiar una sola línea de `RevisarTodos`**.

> [!NOTE]
> 🔧 **Observación de ingeniería de software 3.2 — la más importante de la lección**
> Esa satisfacción implícita tiene una consecuencia que no existe en Java ni en C#: **puedes definir una
> interfaz para código que no escribiste tú.** Si una biblioteca ajena tiene un tipo con un método
> `Revisar`, ese tipo cumple *tu* interfaz sin que el autor lo supiera ni tuviera que cooperar. En otros
> lenguajes, si el autor no declaró la interfaz, no hay nada que hacer.

> [!TIP]
> ✅ **Buena práctica 3.4 — las dos reglas de oro de las interfaces en Go**
> **1. Hazlas pequeñas.** Uno o dos métodos. Una interfaz de diez métodos casi siempre viene de otro
> lenguaje — y es, de hecho, el tercer antipatrón de la sección «Lo que se hace mal». Las de la biblioteca
> estándar más usadas —`io.Reader`, `io.Writer`— tienen **un** método.
> **2. Defínelas donde se USAN, no donde se implementan.** La interfaz `Revisor` pertenece al código que
> necesita revisar cosas, no al que sabe cómo revisarlas. Eso invierte la dependencia: quien consume
> declara lo que necesita.
>
> Y el resumen que vas a oír mucho: **«acepta interfaces, devuelve structs»**.

### 3.7 Por qué esto hace tu código probable

Mira otra vez la línea 65 de la Fig. 3.10. Acabas de probar `RevisarTodos` **simulando un servicio
caído**, sin apagar nada, sin red, y en un milisegundo.

Eso es lo que la gente quiere decir cuando habla de «código probable», y en Go se consigue **sin
bibliotecas de mocks, sin anotaciones y sin frameworks**. Solo con una interfaz pequeña.

> [!NOTE]
> 🔧 **Observación de ingeniería de software 3.3**
> La pregunta que conviene hacerse al diseñar: *«¿puedo probar esto sin que exista el mundo exterior?»* Si
> la respuesta es no, normalmente falta una interfaz. En la lección 5 vas a escribir pruebas de verdad, y
> vas a agradecer haber hecho esto ahora.

---

## El error que vas a ver

El error más frecuente de esta lección es el **`nil pointer dereference`**, y a diferencia del antipatrón
de la Fig. 3.4 (que falla en silencio), éste sí te avisa — con un mensaje que al principio asusta más de
lo que debería.

**Fig. 3.6** | Un programa que **compila perfecto** y **truena al ejecutarse**.

```go
 1  // fig03_06.go — este programa COMPILA pero truena al correr
 2  package main
 3
 4  import "fmt"
 5
 6  type Servicio struct{ Puerto int }
 7
 8  func main() {
 9      var p *Servicio          // un puntero sin apuntar a nada: vale nil
10      fmt.Println(p)           // esto sí funciona
11      fmt.Println(p.Puerto)    // esto truena
12  }
```

```bash
$ go run fig03_06.go
<nil>
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x...]

goroutine 1 [running]:
main.main()
	/home/usuario/fig03_06.go:11 +0x18
exit status 2
```

**El valor cero de un puntero es `nil`**: «no apunto a nada». Intentar leer un campo a través de un
puntero nulo provoca un **panic**, que es la forma en que Go aborta un programa cuando pasa algo
irrecuperable.

**Cómo leer este mensaje, línea por línea:**

- `panic: runtime error: invalid memory address or nil pointer dereference` — el nombre del problema.
  Cuando lo veas, la causa casi siempre es un puntero, un mapa o una interfaz en `nil` que intentaste usar
  como si tuviera algo adentro.
- `goroutine 1 [running]:` — qué hilo de ejecución de Go estaba corriendo cuando reventó. Por ahora, en un
  programa sin concurrencia, siempre va a ser el 1 (la función `main`); en la lección 6, con varias
  goroutines, este número empieza a importar de verdad.
- `main.main() … fig03_06.go:11` — **el archivo y la línea exactos** donde pasó. Empieza siempre por ahí,
  no por el mensaje de arriba: el mensaje te dice *qué tipo* de error fue, esta línea te dice *dónde*.
- `exit status 2` — el programa terminó con error (un `exit status 0` es éxito).

> [!WARNING]
> ⚠️ **Error común de programación 3.2**
> El `nil pointer dereference` es el panic más frecuente en Go. Fíjate en que **compila perfectamente**:
> el compilador no puede saber si un puntero va a valer `nil` en ejecución. Cuando lo veas, ve directo al
> archivo y la línea que te da la traza — en el ejemplo, `fig03_06.go:11`.

---

## Lo que se hace mal

**1. Escribir un método que debería modificar el struct, con receptor de valor.** Es el error de la
Fig. 3.4 (sección 3.3.1): compila, se ejecuta, y **no hace absolutamente nada** — sin ningún mensaje. Es
uno de los tropiezos más desconcertantes para quien empieza, precisamente porque no hay síntoma alguno. La
regla que lo evita está en la Buena práctica 3.2: si el método modifica, receptor de puntero.

**2. Envolver un error con `%v` en vez de `%w`.** Es el error de la sección 3.5: el mensaje impreso queda
idéntico, así que nada avisa, pero `errors.Is` deja de reconocer el error original y el código que decidía
según el tipo de fallo empieza a tomar el camino equivocado. De los dos antipatrones de esta lección, éste
es el más caro, porque nadie nota el síntoma hasta que el programa ya tomó una mala decisión con él.

**3. Escribir interfaces grandes, con muchos métodos, copiando el hábito de otro lenguaje.** En Go las
interfaces se hacen pequeñas —uno o dos métodos, como viste en la Buena práctica 3.4— y se definen del
lado de quien las consume, no de quien las implementa. Una interfaz de diez métodos casi nunca es
necesaria: normalmente basta con el o los métodos que la función que la recibe realmente llama.

---

## Ejercicios

### Ejercicios de repaso

**3.1** ¿Qué ventaja tiene un struct sobre tres variables sueltas?

**3.2** ¿Qué significa que un campo empiece con mayúscula?

**3.3** ¿Qué imprime `fmt.Printf("%+v\n", Servicio{})` si el struct tiene `Nombre string` y `Puerto int`?

**3.4** ¿Cuál es la diferencia entre `func (s Servicio) X()` y `func (s *Servicio) X()`, y cuándo se usa
cada uno?

**3.5** Este método compila, se ejecuta y no hace nada. ¿Por qué?

```go
func (s Servicio) Renombrar(nuevo string) {
    s.Nombre = nuevo
}
```

**3.6** ¿Qué es `nil` para un puntero y qué pasa si lees un campo a través de uno?

**3.7** ¿Qué tiene que tener un tipo para ser un `error`?

**3.8** ¿Cuál es la diferencia entre `%w` y `%v` al envolver un error, y por qué es peligrosa?

**3.9** ¿Cuándo usas `errors.Is` y cuándo `errors.As`?

**3.10** ¿Qué hay que escribir para que un tipo cumpla una interfaz en Go?

**3.11** ¿Qué imprime este programa?

```go
type Estado struct{ Codigo int }

func (e Estado) OK() bool { return e.Codigo == 200 }

func main() {
    var e Estado
    fmt.Println(e.OK())
}
```

### Ejercicios de código

**3.12** Agrega al struct `Servicio` un campo `TimeoutMs int` y un método `Timeout() time.Duration` que lo
convierta. Pista: `time.Duration(s.TimeoutMs) * time.Millisecond`.

**3.13** Escribe un método `func (s *Servicio) Normalizar()` que: quite espacios del nombre con
`strings.TrimSpace`, lo pase a minúsculas con `strings.ToLower`, y si el puerto es 0 lo ponga en 443.
**Explícate por qué este método necesita receptor de puntero.**

**3.14** Define un tipo `Estado` con los campos `Servicio`, `Codigo`, `Duracion` y `Err`. Escríbele un
método `OK() bool` que devuelva verdadero solo si no hay error **y** el código está entre 200 y 299.
Pruébalo con cinco casos, incluido el `Estado{}` vacío.

**3.15** Crea un error centinela `ErrTimeout` y una función que lo devuelva envuelto con contexto.
Comprueba con `errors.Is` que se detecta. **Después cambia el `%w` por `%v` y comprueba que `errors.Is`
devuelve `false`.** Anota en la bitácora que el mensaje impreso no cambió.

**3.16** Define un error propio `ErrorValidacion` con los campos `Campo string` y `Motivo string`, hazlo
cumplir la interfaz `error`, y recupéralo con `errors.As` para imprimir qué campo falló.

**3.17** Define la interfaz `Notificador` con un método `Notificar(mensaje string) error`. Implementa
`NotificadorConsola` (que imprime) y `NotificadorFalso` (que guarda los mensajes en un slice para poder
revisarlos). Escribe una función que reciba un `Notificador` y úsala con las dos.

**3.18 (Encuentra el error)** Di qué está mal en cada uno **sin compilar**, y luego compila para
confirmar:

```go
// (a)
func (s Servicio) Renombrar(n string) { s.Nombre = n }

// (b)
return fmt.Errorf("fallo al revisar %s: %v", nombre, err)

// (c)
var p *Servicio
fmt.Println(p.Nombre)

// (d)
type Revisor interface {
    Revisar(s Servicio) Estado
}
type MiRevisor struct{}
func (m MiRevisor) revisar(s Servicio) Estado { return Estado{} }
// y luego:  var r Revisor = MiRevisor{}
```

**3.19 (Proyecto del curso)** Reorganiza tu programa de la lección 2 usando lo de esta lección:
1. Los structs `Servicio` y `Estado`.
2. Los métodos `Etiqueta()`, `EsSeguro()` y `OK()`.
3. La interfaz `Revisor` con `RevisorHTTP` (que por ahora devuelve datos inventados) y `RevisorFalso`.
4. La función `RevisarTodos(r Revisor, servicios []Servicio) []Estado`.
5. Un `main` que imprima el reporte usando el **falso**, con un servicio que falle.

🔑 **Cuando termines, fíjate en algo: tu programa ya se puede probar completo sin conectarse a nada, y
todavía no has escrito una sola prueba.** Eso es lo que acabas de ganar en esta lección.

### Soluciones

**Ejercicios de repaso (3.1 a 3.11):**

**3.1** Agrupa los datos que van juntos en una sola variable, así que el lenguaje —y quien lea el código—
sabe que pertenecen a la misma cosa. Y no se pueden combinar por error datos de dos servicios distintos.

**3.2** Que está **exportado**: es visible desde otros paquetes. Con minúscula solo se ve dentro de su
paquete. Es todo el control de acceso que tiene Go.

**3.3** `{Nombre: Puerto:0}` — los valores cero, con los nombres de los campos porque es `%+v`.

**3.4** El primero es **receptor de valor**: recibe una copia y no puede modificar el original. El segundo
es **receptor de puntero**: recibe la dirección y sí puede. Se usa puntero cuando el método modifica, o
cuando el struct es grande y copiarlo cuesta.

**3.5** Porque el receptor es de **valor**: el método modifica una copia que se descarta al terminar. Hay
que cambiarlo a `func (s *Servicio) Renombrar(...)`.

**3.6** `nil` es el valor cero de un puntero y significa «no apunto a nada». Leer un campo a través de un
puntero `nil` provoca un **panic** (`nil pointer dereference`) y el programa aborta.

**3.7** Un método `Error() string`. Nada más: `error` es una interfaz con ese único método.

**3.8** `%w` **envuelve** el error original y lo conserva dentro; `%v` solo lo convierte a texto. Es
peligrosa porque **el mensaje impreso es idéntico**, así que no hay síntoma — pero `errors.Is` deja de
encontrar el error original y el código que decidía según el tipo de fallo empieza a equivocarse.

**3.9** `errors.Is` para preguntar *«¿es este error en particular?»*. `errors.As` para preguntar *«¿es de
este tipo? dámelo»*, cuando necesitas leer los datos que el error lleva adentro.

**3.10** **Nada.** Basta con tener los métodos que la interfaz pide; el compilador lo verifica solo. No
existe `implements`.

**3.11** `false` — el valor cero de `Estado` tiene `Codigo: 0`, que no es 200. Nota que el método funciona
perfectamente sobre un struct vacío: eso es el valor cero siendo útil.

**Ejercicios de código (3.12 a 3.19):** sus soluciones verificadas —compiladas y ejecutadas— se agregan en
la entrega en la que se cierre el capítulo de ejercicios con solución del curso completo; no se publican
sin haber corrido `go build` sobre cada una.

---

## Cómo sé que lo logré

- Compilaste y corriste las figuras 3.1 a 3.10 y tu salida coincide con la mostrada.
- Puedes explicar, sin ver el texto, qué diferencia hay entre receptor de valor y receptor de puntero, y
  dar un ejemplo de cuándo usar cada uno.
- Hiciste el experimento del ejercicio 3.15 (cambiar `%w` por `%v`) y viste con tus propios ojos que
  `errors.Is` deja de encontrar el error **sin que el mensaje impreso cambie**.
- Puedes explicarle a alguien más, sin tecnicismos, por qué en Go «cumplir una interfaz» no se declara en
  ningún lado.
- Terminaste el ejercicio 3.19: tu programa del `revisor` ya usa structs, métodos, la interfaz `Revisor` y
  se puede correr con el `RevisorFalso` sin tocar la red.

**Antes de cerrar:** en [`bitacora.md`](bitacora.md) anota qué te costó más entre punteros, errores e
interfaces, y el resultado del ejercicio 3.15 —el del `%w` contra `%v`—. Ese experimento es el que más se
olvida y el que más caro sale en un programa real.

---

## Resumen

- Un **struct** agrupa datos relacionados en un tipo propio. Se define con `type Nombre struct { … }`.
- Los campos con **mayúscula** son exportados; con minúscula, privados al paquete.
- El **valor cero** de un struct llena cada campo con su propio valor cero: nunca hay basura.
- **`%+v`** imprime un struct con los nombres de sus campos: la mejor herramienta para depurar.
- Un **método** es una función con **receptor**: `func (s Servicio) X()`.
- **Receptor de valor** recibe una copia y **no puede modificar** el original; **receptor de puntero**
  (`*T`) sí puede.
- Un **puntero** guarda la dirección de otra variable. Go inserta el `&` y el `*` por ti al llamar
  métodos.
- El valor cero de un puntero es **`nil`**; leer a través de él provoca un **panic**.
- **`error` es una interfaz** con un solo método `Error() string`. Un error es un valor común.
- Se crean con **`errors.New`** (simples), como **centinelas** (`var ErrX = errors.New(...)`) o con
  **`fmt.Errorf`** (con datos).
- **`%w`** envuelve un error conservando el original; **`%v` lo aplasta y rompe `errors.Is` sin avisar**.
- **`errors.Is`** pregunta si un error está en la cadena; **`errors.As`** recupera el error de un tipo
  concreto para leer sus datos.
- Una **interfaz** es una lista de métodos. Un tipo la cumple **por tener los métodos**, sin declararlo.
- Las interfaces de Go se hacen **pequeñas** y se definen **donde se usan**.
- Una interfaz pequeña permite **probar sin el mundo exterior**: se sustituye la implementación real por
  una falsa.

---

## Para leer más

1. **[A Tour of Go — Methods](https://go.dev/tour/methods/1)** — oficial e interactivo: receptores,
   punteros e interfaces con ejercicios en el navegador.
2. **[Effective Go — Interfaces and other types](https://go.dev/doc/effective_go#interfaces_and_types)**
   — la explicación oficial de por qué la satisfacción de interfaces es implícita.
3. **[Error handling and Go](https://go.dev/blog/error-handling-and-go)** — la entrada del blog oficial de
   Go sobre por qué los errores son valores y cómo se manejan idiomáticamente.
4. **[Paquete `errors` — documentación oficial](https://pkg.go.dev/errors)** — la referencia exacta de
   `errors.Is`, `errors.As` y `errors.Unwrap`.

### Términos de esta lección

| | |
|---|---|
| **campo** | cada dato que contiene un struct |
| **envolver (un error)** | agregarle contexto conservando el original, con `%w` |
| **error centinela** | error declarado una vez a nivel de paquete, para comparar con `errors.Is` |
| **exportado / no exportado** | visible fuera del paquete (mayúscula) o no (minúscula) |
| **interfaz** | lista de métodos; un tipo la cumple al tener esos métodos |
| **método** | función asociada a un tipo mediante un receptor |
| **`nil`** | ausencia de valor; el valor cero de punteros, interfaces y errores |
| **`nil pointer dereference`** | panic por leer a través de un puntero nulo |
| **panic** | aborto del programa por un error irrecuperable |
| **puntero** | variable que guarda la dirección de otra; su tipo se escribe `*T` |
| **receptor** | el `(s Servicio)` que liga una función a un tipo |
| **satisfacción implícita** | cumplir una interfaz sin declararlo |
| **struct** | tipo que agrupa varios campos |

---

**Anterior:** [Lección 2 — Variables, funciones y tipos](02-fundamentos.md) ·
**Siguiente:** [Lección 4 — Colecciones](04-colecciones.md)
