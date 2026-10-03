# Lesson 3 — Structs, methods, errors, and interfaces

> **This is where Go becomes Go.** This is the most important lesson of the course.

**Duration:** two sessions of 60 minutes. Don't do it in one go.

**By the end you will be able to:**

- Group related data in a **struct** and explain why it is better than loose variables.
- Write **methods** and decide when to use a value receiver and when a pointer receiver.
- Explain what a **pointer** is without getting scared.
- Handle errors as values, wrap them, and ask about them.
- Define an **interface** and understand why in Go you don't declare that it is satisfied.
- Write code that can be tested **without connecting to anything**.

---

## Why it matters

In lesson 2 you stored a service's data in three separate variables (`nombre`, `url`, `puerto`). It works
with **one** service. With ten that is thirty loose variables, and nothing in the code says which ones go
together — if you make a mistake and combine one service's name with another's port, **the program
compiles anyway** and produces an incorrect report. That is the first problem this lesson solves: the
**struct**.

The second problem goes deeper. The `revisor` will have to query real services, and real queries
**fail**: the network goes down, the service is slow, it answers with a code you didn't expect. A program
that doesn't know how to handle that in an orderly way is useless in production. Go solves this in a way
you have probably not seen if you come from another language: **errors are values**, not exceptions that
interrupt the flow.

And the third problem is the one that makes the previous two pieces fit together without the `revisor`
having to know, in advance, **how** each service is going to be checked (HTTP today, maybe a database
tomorrow). That is solved by the **interface**, which is also what will let you, in lesson 5, test the
whole program **without connecting to anything real**.

Structs, errors, and interfaces are, in that order, the backbone of everything that follows in the
course.

---

## The concepts

### 3.1 Structs: grouping what goes together

What you need is to tell the language: *"these three things are a service"*.

**Fig. 3.1** | Defining and using a struct.

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

**Line 8: `type Servicio struct {`.** It reads: *"define a new type called `Servicio`, which is a
structure"*. **You just created a type that didn't come with the language**, and from now on it is used
just like `int` or `string`.

**Lines 15-19.** Creates a value of that type. The field names with a colon are optional —you could
write just the values in order— but **always include them**:

> [!TIP]
> ✅ **Good practice 3.1**
> Always write the field names when creating a struct: `Servicio{Nombre: "x", Puerto: 443}` instead of
> `Servicio{"x", "", 443}`. The short version breaks silently if someone adds a field or changes the
> order, and the compiler can't warn you because the types still line up.

**Line 21: `s.Nombre`.** The dot accesses a field. It reads "the `Nombre` of `s`".

#### 3.1.1 The capital letter is permission, not aesthetics

Notice that the fields start with a **capital letter**: `Nombre`, `URL`, `Puerto`. In Go that is not a
chosen style: **it is the language's access control.**

| | |
|---|---|
| `Nombre` (capital) | **exported**: visible from other packages |
| `nombre` (lowercase) | **unexported**: only visible inside its own package |

There is no `public`, `private`, or `protected`. **The initial letter is the entire rule.**

#### 3.1.2 The zero value of a struct

**Fig. 3.2** | An uninitialized struct.

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

Each field receives **its own** zero value: texts end up as `""` and numbers as `0`. **There is no garbage
or "undefined"**, and that is why a freshly created struct can already be used without fear.

**And here `%+v` appears**, which is very useful for debugging: it prints the struct **with the field
names**. Compare it with `%v`, which prints only the values.

> [!TIP]
> 🧪 **Testing and debugging tip 3.1**
> When you don't understand what a struct holds, print it with `%+v`. It is the fastest way to see all its
> fields with their names, and it will save you a great deal of time.

### 3.2 Methods

A **method** is a function that belongs to a type. Instead of writing `etiqueta(s)`, you write
`s.Etiqueta()`, and the function is tied to the type it corresponds to.

**Fig. 3.3** | Methods on a struct.

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

**Line 14: `func (s Servicio) Etiqueta() string`.** That `(s Servicio)` between `func` and the name is
called the **receiver**, and it is what turns a function into a method. It reads: *"this function belongs
to the `Servicio` type, and inside it I am going to refer to the value as `s`"*.

> [!NOTE]
> 🔧 **Software engineering observation 3.1**
> If you come from Java, C#, or Python, a Go method is similar to a class method — with one important
> difference: **in Go the method is written outside the type.** You can have the `type` in one file and its
> methods in another. And since there are no classes, there is no inheritance either: in Go code is reused
> in another way, and you will see it in section 3.6.

### 3.3 Pointers, in ten minutes

Pointers have a reputation for being hard. In Go they are much simpler than in C, and you need to
understand **just one idea** for what comes next.

#### 3.3.1 The problem

**Fig. 3.4** | A method that **doesn't work** the way you would expect.

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

**Nothing changed.** The method ran, there was no error, and the value is still the same.

**Why?** Because when the receiver is `(s Servicio)`, Go hands the method **a copy** of the struct. The
method modifies the copy, the copy is discarded when it finishes, and the original never found out. It is
such a frequent antipattern that it has its own entry in "What goes wrong", below.

#### 3.3.2 The solution: the pointer receiver

A **pointer** is a variable that stores **the address** of another, instead of a copy of its contents.

Think of the difference between handing you **a photocopy** of a document and handing you **the address
of the filing cabinet where the original is**. With the photocopy you can write all you want: the original
doesn't change. With the address, you can go and modify the original.

- `Servicio` is the photocopy.
- `*Servicio` is the address of the original.

**Fig. 3.5** | The same method, now with a pointer receiver.

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

**The only difference between Fig. 3.4 and 3.5 is an asterisk on line 13.** That is all.

And notice line 21: you wrote `x.CambiarPuerto(8443)` just like before, **without `&` or anything odd**.
Go realizes that the method needs a pointer and takes it on its own. That is the reason Go's pointers are
much less work than C's.

> [!TIP]
> ✅ **Good practice 3.2**
> The rule is simple: **if the method modifies the struct, pointer receiver (`*T`); if it only reads,
> value receiver (`T`)**. And be consistent within a single type: if most of your methods need a pointer,
> use a pointer on all of them, even if some don't require it. Mixing them confuses whoever reads the code.

> [!NOTE]
> 🚀 **Performance tip 3.1**
> There is a second reason to use pointers: **avoiding copying**. If a struct has twenty fields, every
> call with a value receiver copies all twenty. For small structs it is irrelevant; for large ones or in
> loops of millions of iterations, it matters. **Don't optimize this without measuring:** clarity is worth
> more than a 40-byte copy.

#### 3.3.3 The only dangerous pointer: `nil`

There is still a third situation with pointers, and this one does produce a real error message —so real
that it has its own dedicated section: **"The error you will see"**, further down in this lesson. We
anticipate the idea: the zero value of a pointer is **`nil`** ("I point to nothing"), and reading a field
through a nil pointer makes the program **blow up** instead of silently misbehaving as in Fig. 3.4.

### 3.4 Errors are values

You already saw in lesson 2 that a Go function can return `(result, error)`. Now we get to the bottom of
it, because **this is what most distinguishes Go from what you have probably seen.**

#### 3.4.1 `error` is an interface, not something magic

In Go, `error` is simply a type with a method:

```go
type error interface {
    Error() string
}
```

That means: *"anything that has an `Error()` method that returns text is an error"*. There is no hierarchy
of exception classes, no `throw`, no call stack that unwinds. **An error is an ordinary value that travels
like any other.**

#### 3.4.2 Creating errors

**Fig. 3.7** | The three ways to create an error.

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

**Line 11: the sentinel error.** It is declared **once**, at package level, with the `Err` prefix. It
lets whoever calls your function ask *"was it this particular error?"*, as you will see in 3.5.

**Line 30: `fmt.Errorf`.** Like `Printf`, but it produces an error instead of printing. Use it when the
message needs data.

**Line 44: `if err := Validar(s); err != nil`.** This pattern declares `err` **inside** the `if`, so it
only exists there. It is very idiomatic in Go and keeps the code clean.

> [!TIP]
> ✅ **Good practice 3.3**
> Error messages are written **in lowercase and without a final period**: `"no se pudo abrir el archivo"`,
> not `"No se pudo abrir el archivo."`. The reason is practical: errors get **wrapped** in one another
> (section 3.5), and when concatenated they end up as `"revisando catalogo: no se pudo abrir el archivo"`.
> With capitals and periods, the result would look broken.

### 3.5 Wrapping errors: `%w`, `errors.Is`, and `errors.As`

An error without context is of little use. If your program says `connection refused`, you don't know
**which** service failed.

**Wrapping** an error is adding context to it **without losing the original**.

**Fig. 3.8** | Wrapping errors and asking about them.

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

**Line 23: `%w`.** This is the key verb. It looks the same as `%v` when printed, **but it keeps the original
error inside** so it can be asked about later.

**Line 35: `errors.Is`.** It asks *"is that error anywhere in this chain?"*. It works even with five layers
of wrapping.

> [!WARNING]
> 🔴 **Common programming error 3.3 — the most silent one in this lesson**
> Writing `%v` instead of `%w` when wrapping:
>
> ```go
> return fmt.Errorf("revisando %s: %v", s.Nombre, err)   // ❌ con %v
> ```
>
> **The printed message is identical.** There is no error, no warning, everything seems to work. But
> `errors.Is` stops finding the original error and returns `false`, so the code that decided what to do
> according to the kind of failure starts taking the wrong path. **It is the kind of error that doesn't
> fail: it returns worse data.** When you wrap, use `%w`. This is the second antipattern in the "What
> goes wrong" section.

#### 3.5.1 `errors.As`, when you need the error's data

`errors.Is` answers "is it this error?". `errors.As` answers "is it of this **type**? give it to me so I
can read its fields".

**Fig. 3.9** | A custom error with data, recovered with `errors.As`.

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

**Line 17.** By writing an `Error() string` method, your type **is already an `error`**. You didn't
declare anything: the type satisfies the interface because it has the method. That is what the next
section explains.

**Line 42.** And here is the real value: the program can **decide** according to the kind of failure. A
503 is retried; a 404 is not. With text-only errors that would be impossible without comparing strings,
which is fragile.

### 3.6 Interfaces: the heart of Go

An **interface** is a list of methods. Any type that has those methods **satisfies it**, and you don't
have to declare it anywhere.

**Fig. 3.10** | An interface and two implementations.

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

**Read it slowly, because this is the idea that holds up all of Go.**

**Lines 22-24: the interface.** It says: *"a `Revisor` is anything that has a `Revisar` method that
receives a `Servicio` and returns an `Estado`"*.

**Lines 29 and 39.** `RevisorHTTP` and `RevisorFalso` **don't declare anywhere** that they satisfy
`Revisor`. There is no `implements`, no `: Revisor`, nothing. **They satisfy the interface because they
have the method**, and the compiler verifies that on its own.

**Line 45: `func RevisarTodos(r Revisor, ...)`.** This function **doesn't know** what it is working with. It
only knows that it can call `.Revisar()`. And that is why lines 60 and 68 pass it two completely
different things **without changing a single line of `RevisarTodos`**.

> [!NOTE]
> 🔧 **Software engineering observation 3.2 — the most important one in the lesson**
> That implicit satisfaction has a consequence that doesn't exist in Java or C#: **you can define an
> interface for code you didn't write.** If a third-party library has a type with a `Revisar` method, that
> type satisfies *your* interface without the author knowing or having to cooperate. In other languages,
> if the author didn't declare the interface, there is nothing to be done.

> [!TIP]
> ✅ **Good practice 3.4 — the two golden rules of interfaces in Go**
> **1. Keep them small.** One or two methods. A ten-method interface almost always comes from another
> language — and is, in fact, the third antipattern in the "What goes wrong" section. The most-used ones in
> the standard library —`io.Reader`, `io.Writer`— have **one** method.
> **2. Define them where they are USED, not where they are implemented.** The `Revisor` interface belongs
> to the code that needs to check things, not to the code that knows how to check them. That inverts the
> dependency: whoever consumes declares what it needs.
>
> And the summary you will hear a lot: **"accept interfaces, return structs"**.

### 3.7 Why this makes your code testable

Look again at line 65 of Fig. 3.10. You just tested `RevisarTodos` **simulating a service that is down**,
without shutting anything down, without a network, and in a millisecond.

That is what people mean when they talk about "testable code", and in Go it is achieved **without mock
libraries, without annotations, and without frameworks**. Just with a small interface.

> [!NOTE]
> 🔧 **Software engineering observation 3.3**
> The question worth asking when designing: *"can I test this without the outside world existing?"* If
> the answer is no, an interface is usually missing. In lesson 5 you will write real tests, and you will be
> grateful that you did this now.

---

## The error you will see

The most frequent error of this lesson is the **`nil pointer dereference`**, and unlike the antipattern of
Fig. 3.4 (which fails silently), this one does warn you — with a message that at first scares more than it
should.

**Fig. 3.6** | A program that **compiles perfectly** and **blows up when run**.

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

**The zero value of a pointer is `nil`**: "I point to nothing". Trying to read a field through a nil
pointer causes a **panic**, which is the way Go aborts a program when something unrecoverable happens.

**How to read this message, line by line:**

- `panic: runtime error: invalid memory address or nil pointer dereference` — the name of the problem.
  When you see it, the cause is almost always a pointer, a map, or an interface that is `nil` that you tried
  to use as if it had something inside.
- `goroutine 1 [running]:` — which Go thread of execution was running when it blew up. For now, in a
  program without concurrency, it will always be 1 (the `main` function); in lesson 6, with several
  goroutines, this number starts to really matter.
- `main.main() … fig03_06.go:11` — **the exact file and line** where it happened. Always start there, not
  with the message above: the message tells you *what kind* of error it was, this line tells you *where*.
- `exit status 2` — the program ended with an error (an `exit status 0` is success).

> [!WARNING]
> ⚠️ **Common programming error 3.2**
> The `nil pointer dereference` is the most frequent panic in Go. Notice that it **compiles perfectly**:
> the compiler can't know whether a pointer will be `nil` at run time. When you see it, go straight to the
> file and line that the trace gives you — in the example, `fig03_06.go:11`.

---

## What goes wrong

**1. Writing a method that should modify the struct, with a value receiver.** It is the error of
Fig. 3.4 (section 3.3.1): it compiles, it runs, and **it does absolutely nothing** — without any message.
It is one of the most disconcerting stumbles for beginners, precisely because there is no symptom at all.
The rule that avoids it is in Good practice 3.2: if the method modifies, pointer receiver.

**2. Wrapping an error with `%v` instead of `%w`.** It is the error of section 3.5: the printed message
stays identical, so nothing warns you, but `errors.Is` stops recognizing the original error and the code
that decided according to the kind of failure starts taking the wrong path. Of the two antipatterns in
this lesson, this one is the most expensive, because nobody notices the symptom until the program has
already made a bad decision with it.

**3. Writing large interfaces, with many methods, copying the habit of another language.** In Go
interfaces are made small —one or two methods, as you saw in Good practice 3.4— and are defined on the
side of whoever consumes them, not whoever implements them. A ten-method interface is almost never
necessary: normally the method or methods that the function receiving it really calls are enough.

---

## Exercises

### Review exercises

**3.1** What advantage does a struct have over three loose variables?

**3.2** What does it mean for a field to start with a capital letter?

**3.3** What does `fmt.Printf("%+v\n", Servicio{})` print if the struct has `Nombre string` and
`Puerto int`?

**3.4** What is the difference between `func (s Servicio) X()` and `func (s *Servicio) X()`, and when is
each one used?

**3.5** This method compiles, runs, and does nothing. Why?

```go
func (s Servicio) Renombrar(nuevo string) {
    s.Nombre = nuevo
}
```

**3.6** What is `nil` for a pointer and what happens if you read a field through one?

**3.7** What does a type need in order to be an `error`?

**3.8** What is the difference between `%w` and `%v` when wrapping an error, and why is it dangerous?

**3.9** When do you use `errors.Is` and when `errors.As`?

**3.10** What do you have to write for a type to satisfy an interface in Go?

**3.11** What does this program print?

```go
type Estado struct{ Codigo int }

func (e Estado) OK() bool { return e.Codigo == 200 }

func main() {
    var e Estado
    fmt.Println(e.OK())
}
```

### Code exercises

**3.12** Add to the `Servicio` struct a `TimeoutMs int` field and a `Timeout() time.Duration` method that
converts it. Hint: `time.Duration(s.TimeoutMs) * time.Millisecond`.

**3.13** Write a method `func (s *Servicio) Normalizar()` that: removes spaces from the name with
`strings.TrimSpace`, lowercases it with `strings.ToLower`, and if the port is 0 sets it to 443.
**Explain to yourself why this method needs a pointer receiver.**

**3.14** Define an `Estado` type with the fields `Servicio`, `Codigo`, `Duracion`, and `Err`. Write it an
`OK() bool` method that returns true only if there is no error **and** the code is between 200 and 299.
Test it with five cases, including the empty `Estado{}`.

**3.15** Create a sentinel error `ErrTimeout` and a function that returns it wrapped with context. Check
with `errors.Is` that it is detected. **Then change the `%w` to `%v` and check that `errors.Is` returns
`false`.** Write down in the logbook that the printed message didn't change.

**3.16** Define a custom error `ErrorValidacion` with the fields `Campo string` and `Motivo string`, make
it satisfy the `error` interface, and recover it with `errors.As` to print which field failed.

**3.17** Define the `Notificador` interface with a method `Notificar(mensaje string) error`. Implement
`NotificadorConsola` (which prints) and `NotificadorFalso` (which stores the messages in a slice so they
can be inspected). Write a function that receives a `Notificador` and use it with both.

**3.18 (Find the error)** Say what is wrong in each one **without compiling**, and then compile to
confirm:

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

**3.19 (Course project)** Reorganize your lesson 2 program using what is in this lesson:
1. The `Servicio` and `Estado` structs.
2. The `Etiqueta()`, `EsSeguro()`, and `OK()` methods.
3. The `Revisor` interface with `RevisorHTTP` (which for now returns made-up data) and `RevisorFalso`.
4. The `RevisarTodos(r Revisor, servicios []Servicio) []Estado` function.
5. A `main` that prints the report using the **fake** one, with a service that fails.

🔑 **When you finish, notice something: your program can already be tested in full without connecting to
anything, and you haven't written a single test yet.** That is what you just gained in this lesson.

### Solutions

**Review exercises (3.1 to 3.11):**

**3.1** It groups the data that goes together in a single variable, so the language —and whoever reads the
code— knows they belong to the same thing. And data from two different services can't be combined by
mistake.

**3.2** That it is **exported**: it is visible from other packages. With a lowercase letter it is only seen
inside its package. It is all the access control Go has.

**3.3** `{Nombre: Puerto:0}` — the zero values, with the field names because it is `%+v`.

**3.4** The first is a **value receiver**: it receives a copy and can't modify the original. The second is
a **pointer receiver**: it receives the address and can. A pointer is used when the method modifies, or
when the struct is large and copying it is costly.

**3.5** Because the receiver is a **value**: the method modifies a copy that is discarded when it
finishes. It must be changed to `func (s *Servicio) Renombrar(...)`.

**3.6** `nil` is the zero value of a pointer and means "I point to nothing". Reading a field through a
`nil` pointer causes a **panic** (`nil pointer dereference`) and the program aborts.

**3.7** An `Error() string` method. Nothing more: `error` is an interface with that single method.

**3.8** `%w` **wraps** the original error and keeps it inside; `%v` only converts it to text. It is
dangerous because **the printed message is identical**, so there is no symptom — but `errors.Is` stops
finding the original error and the code that decided according to the kind of failure starts getting it
wrong.

**3.9** `errors.Is` to ask *"is this particular error?"*. `errors.As` to ask *"is it of this type? give it
to me"*, when you need to read the data the error carries inside.

**3.10** **Nothing.** It is enough to have the methods the interface asks for; the compiler verifies it on
its own. There is no `implements`.

**3.11** `false` — the zero value of `Estado` has `Codigo: 0`, which is not 200. Note that the method works
perfectly on an empty struct: that is the zero value being useful.

**Code exercises (3.12 to 3.19):** their verified solutions —compiled and run— will be added in the
delivery in which the chapter of exercises with solutions for the whole course is closed; they are not
published without having run `go build` on each one.

---

## How I know I got it

- You compiled and ran figures 3.1 to 3.10 and your output matches the one shown.
- You can explain, without looking at the text, what the difference is between a value receiver and a
  pointer receiver, and give an example of when to use each.
- You did the experiment of exercise 3.15 (changing `%w` to `%v`) and saw with your own eyes that
  `errors.Is` stops finding the error **without the printed message changing**.
- You can explain to someone else, without jargon, why in Go "satisfying an interface" is declared
  nowhere.
- You finished exercise 3.19: your `revisor` program already uses structs, methods, the `Revisor`
  interface, and can be run with the `RevisorFalso` without touching the network.

**Before closing:** in [`bitacora.md`](bitacora.md) write down what was hardest for you among pointers,
errors, and interfaces, and the result of exercise 3.15 —the `%w` versus `%v` one—. That experiment is
the one most easily forgotten and the one that costs the most in a real program.

---

## Summary

- A **struct** groups related data in a type of your own. It is defined with `type Name struct { … }`.
- Fields with a **capital letter** are exported; with a lowercase letter, private to the package.
- The **zero value** of a struct fills each field with its own zero value: there is never garbage.
- **`%+v`** prints a struct with its field names: the best tool for debugging.
- A **method** is a function with a **receiver**: `func (s Servicio) X()`.
- A **value receiver** receives a copy and **can't modify** the original; a **pointer receiver**
  (`*T`) can.
- A **pointer** stores the address of another variable. Go inserts the `&` and the `*` for you when calling
  methods.
- The zero value of a pointer is **`nil`**; reading through it causes a **panic**.
- **`error` is an interface** with a single method `Error() string`. An error is an ordinary value.
- They are created with **`errors.New`** (simple ones), as **sentinels** (`var ErrX = errors.New(...)`), or
  with **`fmt.Errorf`** (with data).
- **`%w`** wraps an error keeping the original; **`%v` flattens it and breaks `errors.Is` without warning**.
- **`errors.Is`** asks whether an error is in the chain; **`errors.As`** recovers the error of a concrete
  type to read its data.
- An **interface** is a list of methods. A type satisfies it **by having the methods**, without declaring
  it.
- Go interfaces are made **small** and are defined **where they are used**.
- A small interface lets you **test without the outside world**: the real implementation is replaced by a
  fake one.

---

## Further reading

1. **[A Tour of Go — Methods](https://go.dev/tour/methods/1)** — official and interactive: receivers,
   pointers, and interfaces with exercises in the browser.
2. **[Effective Go — Interfaces and other types](https://go.dev/doc/effective_go#interfaces_and_types)**
   — the official explanation of why interface satisfaction is implicit.
3. **[Error handling and Go](https://go.dev/blog/error-handling-and-go)** — the official Go blog post on
   why errors are values and how they are handled idiomatically.
4. **[Package `errors` — official documentation](https://pkg.go.dev/errors)** — the exact reference for
   `errors.Is`, `errors.As`, and `errors.Unwrap`.

### Terms from this lesson

| | |
|---|---|
| **field** | each piece of data a struct contains |
| **wrap (an error)** | add context to it while keeping the original, with `%w` |
| **sentinel error** | error declared once at package level, to compare with `errors.Is` |
| **exported / unexported** | visible outside the package (capital letter) or not (lowercase) |
| **interface** | list of methods; a type satisfies it by having those methods |
| **method** | function associated with a type through a receiver |
| **`nil`** | absence of a value; the zero value of pointers, interfaces, and errors |
| **`nil pointer dereference`** | panic from reading through a nil pointer |
| **panic** | program abort due to an unrecoverable error |
| **pointer** | variable that stores the address of another; its type is written `*T` |
| **receiver** | the `(s Servicio)` that ties a function to a type |
| **implicit satisfaction** | satisfying an interface without declaring it |
| **struct** | type that groups several fields |

---

**Previous:** [Lesson 2 — Variables, functions, and types](02-fundamentos.md) ·
**Next:** [Lesson 4 — Collections](04-colecciones.md)
