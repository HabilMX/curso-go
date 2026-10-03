# Lesson 2 — Variables, functions, and types

**Duration:** 90 minutes, or two sessions of 45.

**By the end you will be able to:**

- Explain what the Go compiler does with your text file.
- Declare variables and know when to use `:=` and when `=`.
- Name the four basic types and say why Go doesn't let you mix them.
- Write functions that receive and return values.
- Understand why a Go function can return **two** things at once.
- Read the compiler's error messages and know what they are asking of you.
- Write, compile, and run a complete program on your own.

---

## Why it matters

The `revisor` —the program you will build throughout the course— needs, from its first line, three
things: **data** (a service's name, its port, how long it took to respond), **logic** that decides
something with that data (is it fast or slow?), and a way to **signal when something went wrong** (did the
service not answer?). That is, in that order, exactly what this lesson teaches: variables to store data,
functions for the logic, and the two-return-value pattern for errors.

It is no coincidence that Go is **compiled** and **statically typed**: both decisions exist so that the
computer catches your mistakes **before** the program runs, instead of a user discovering them in
production. An interpreted language lets you write `puerto = "muchos"` and only blows up when that line
runs —which can be ten minutes later, or never, if that path of the code is hardly ever used—. Go refuses
to produce the executable. That difference is the underlying reason for almost everything you will see in
this lesson: why the compiler is strict, why the "zero value" exists, and why a function that can fail
**has to say so in its signature**, not as a surprise.

At the end of the lesson you will write the first real step of the `revisor` (project exercise): a
function that receives a service's name, its port, and its response time, and returns a classified report
line. Still without structs, without concurrency, and without HTTP —that comes later— but it is already
code that resembles the final program.

---

## The concepts

### 2.1 What the compiler does

You write a text file. That file, by itself, does nothing: it is text. For the computer to run it,
something has to translate it into the instructions the processor understands.

There are two ways to do that translation:

| | How it works | Languages |
|---|---|---|
| **Interpreted** | A program reads your text line by line and does what it says, **every time** you run it | Python, JavaScript, PHP |
| **Compiled** | A program translates **all** your text **just once** and saves the result in an executable file | **Go**, C, C++, Rust |

**Go is compiled**, and that has three consequences you will notice starting today:

1. **There is a step before running: compiling.** If you misspelled a variable name, you find out right
   there — before the program runs. In an interpreted language you would find out when execution reached
   that line, which could be ten minutes later or never.
2. **The result is a file that works on its own.** It doesn't need the other computer to have Go
   installed.
3. **It is fast to run**, because the translation is already done.

> [!NOTE]
> 🔧 **Software engineering observation 2.1**
> The compiler is your first line of defense, not an obstacle. Every error it catches is an error you
> won't be hunting for at night with the program already delivered. When Go stops you, read the message
> calmly: it is saving you work.

### 2.2 Your first program, line by line

We are going to write, compile, and run a complete program. **All the programs in this course are
complete and executable**: no fragments that can't be run.

Prepare the folder:

```bash
mkdir -p ~/w/curso-go/cap01 && cd ~/w/curso-go/cap01
go mod init cap01
```

**Fig. 2.1** | A program that prints a message.

```go
 1  // fig02_01.go
 2  // Imprime un mensaje en la pantalla.
 3  package main
 4
 5  import "fmt"
 6
 7  func main() {
 8      fmt.Println("Bienvenido a Go!")
 9  }
```

Compile it and run it:

```bash
$ go run fig02_01.go
Bienvenido a Go!
```

Now **each line**, because all of them have a reason:

**Lines 1-2.** Comments. They start with `//` and the compiler ignores them completely. They are there
for whoever reads the code — including you in three weeks.

**Line 3: `package main`.** Every Go file belongs to a *package*, which is a grouping of code. **The
`main` package is special: it is the only one that produces an executable program.** If you wrote
`package utilidades`, you would have a library that other programs can use, but that can't be run on its
own.

**Line 5: `import "fmt"`.** It tells the compiler that you are going to use code from the `fmt` package
(from *format*), which comes with Go and contains the functions for printing and formatting text. **Go
loads nothing by default**: whatever you use, you ask for.

**Line 7: `func main() {`.** Declares the `main` function, which is **where the program's execution
begins**. When you run a Go program, the system looks for exactly this function. If you called it
`principal` or `inicio`, the program wouldn't start.

**Line 8: `fmt.Println("Bienvenido a Go!")`.** Calls the `Println` function that is inside the `fmt`
package. The dot reads as "of": *"the `Println` function **of** `fmt`"*. `Println` prints whatever you
give it and moves to the next line (*print line*).

**Line 9: `}`.** Closes the function body. In Go braces delimit blocks, just like in C, C++, Java, or
JavaScript.

> [!TIP]
> ✅ **Good practice 2.1**
> Always put a comment at the top of the file with its name and what it does. It takes five seconds and
> saves minutes for whoever opens it later.

#### 2.2.1 The difference between `go run` and `go build`

**Fig. 2.2** | The two ways to run your program.

```bash
$ go run fig02_01.go        # compila en memoria, ejecuta, y no deja archivo
Bienvenido a Go!

$ go build fig02_01.go      # compila y GUARDA el ejecutable
$ ls -la
-rwxr-xr-x  1 usuario usuario 1841624  fig02_01      ← el programa
-rw-r--r--  1 usuario usuario     104  fig02_01.go   ← tu texto

$ ./fig02_01                # y se ejecuta directo
Bienvenido a Go!
```

Look at the sizes: your text is **104 bytes**; the compiled program, **1.8 megabytes**. The difference is
that the executable carries inside everything it needs to work.

> [!NOTE]
> 🚀 **Portability tip 2.1**
> That `fig02_01` file can be copied to **any** Linux computer of the same architecture and it works, even
> if that machine doesn't have Go installed. It is the reason so many server tools are written in Go.

> [!TIP]
> 🧪 **Testing and debugging tip 2.1**
> Use `go run` while you program: it is faster and doesn't fill your folder with executables. Use
> `go build` when the program already works and you want to deliver it.

### 2.3 Variables

A **variable** is a named space in memory, where you store a piece of data your program is going to use.

**Fig. 2.3** | Declaring variables and printing them.

```go
 1  // fig02_03.go
 2  // Declara variables de distintos tipos y las imprime.
 3  package main
 4
 5  import "fmt"
 6
 7  func main() {
 8      nombre := "catalogo"                 // texto
 9      puerto := 443                      // número entero
10      tiempo := 0.142                    // número con decimales
11      seguro := true                     // verdadero o falso
12
13      fmt.Println("servicio:", nombre)
14      fmt.Println("puerto:", puerto)
15      fmt.Println("tiempo de respuesta:", tiempo)
16      fmt.Println("usa HTTPS:", seguro)
17  }
```

```bash
$ go run fig02_03.go
servicio: catalogo
puerto: 443
tiempo de respuesta: 0.142
usa HTTPS: true
```

**The `:=` operator** reads *"declare a new variable and store this in it"*. Go looks at the value on the
right and **deduces the type on its own**: `"catalogo"` is between quotes, so it is text; `443` has none
and has no decimal point, so it is an integer.

#### 2.3.1 `:=` versus `=`

Once the variable exists, to change its value you use `=` **without** the colon:

**Fig. 2.4** | Creating and modifying a variable.

```go
 1  // fig02_04.go
 2  package main
 3
 4  import "fmt"
 5
 6  func main() {
 7      puerto := 443          // la CREA con valor 443
 8      fmt.Println(puerto)
 9
10      puerto = 8080          // le CAMBIA el valor
11      fmt.Println(puerto)
12  }
```

```bash
$ go run fig02_04.go
443
8080
```

> [!TIP]
> ✅ **Good practice 2.2**
> If you really need to declare something you are not going to use yet, use the blank identifier `_`. In
> Go, `_` means *"I am discarding this on purpose"*, and the compiler accepts it.

#### 2.3.2 Go refuses to compile if you don't use a variable

**Fig. 2.5** | A program that **does not compile**, on purpose.

```go
 1  // fig02_05.go — este programa NO compila
 2  package main
 3
 4  import "fmt"
 5
 6  func main() {
 7      nombre := "catalogo"
 8      edad := 30              // se declara y nunca se usa
 9      fmt.Println(nombre)
10  }
```

```bash
$ go run fig02_05.go
# command-line-arguments
./fig02_05.go:8:5: declared and not used: edad
```

**Go does not allow it.** It is not a warning you can ignore: it is an error that stops the compilation.

**Why so strict?** Because an unused variable almost always means one of two things: you misspelled the
name on another line, or leftover junk remained from code you half-deleted. Both are problems. Go
prefers to bother you today rather than let dead code pile up forever.

### 2.4 The basic types

A **type** is the answer to the question *"what kind of data is this?"*. The four you will use all the
time:

| Type | What it stores | Examples | Zero value |
|---|---|---|---|
| `string` | text | `"hola"`, `"https://catalogo.example.com"` | `""` (empty) |
| `int` | whole numbers | `42`, `-7`, `0` | `0` |
| `float64` | numbers with decimals | `3.14`, `0.142` | `0` |
| `bool` | true or false | `true`, `false` | `false` |

#### 2.4.1 Go doesn't mix types

**Fig. 2.6** | Another program that **does not compile**, on purpose.

```go
 1  // fig02_06.go — este programa NO compila
 2  package main
 3
 4  import "fmt"
 5
 6  func main() {
 7      puerto := 443
 8      puerto = "muchos"       // intenta guardar texto en un entero
 9      fmt.Println(puerto)
10  }
```

```bash
$ go run fig02_06.go
# command-line-arguments
./fig02_06.go:8:14: cannot use "muchos" (untyped string constant) as int value in assignment
```

In Python this would work without complaint. **In Go it doesn't compile**, and that is an enormous
advantage:

> [!NOTE]
> 🔧 **Software engineering observation 2.2**
> A good share of the errors that reach production are of this kind: someone believed a variable held a
> number and it held text. In an interpreted language that blows up **when a customer is using the
> program**. In Go it never gets to compile: the error shows up on your machine, not the customer's.

#### 2.4.2 The zero value: in Go there is never garbage

**Fig. 2.7** | Variables declared without a value.

```go
 1  // fig02_07.go
 2  // Muestra el valor cero de cada tipo basico.
 3  package main
 4
 5  import "fmt"
 6
 7  func main() {
 8      var texto string
 9      var entero int
10      var decimal float64
11      var logico bool
12
13      fmt.Printf("string:  %q\n", texto)
14      fmt.Printf("int:     %d\n", entero)
15      fmt.Printf("float64: %v\n", decimal)
16      fmt.Printf("bool:    %t\n", logico)
17  }
```

```bash
$ go run fig02_07.go
string:  ""
int:     0
float64: 0
bool:    false
```

The `var name type` form declares a variable **without giving it a value**. In other languages that would
leave "garbage" —whatever was in that memory— or a special "undefined" state. **In Go it always receives a
valid value**, called the **zero value**.

**And here `Printf` appears**, which is different from `Println`: it receives a template with blanks and
fills them in. Each blank starts with `%`:

| | What for |
|---|---|
| `%d` | a whole number (*decimal*) |
| `%s` | text (*string*) |
| `%q` | text **with quotes** (*quoted*) — useful for seeing whether something is empty |
| `%t` | a `bool` (*true/false*) |
| `%v` | anything, in its natural format (*value*) |
| `\n` | line break (`Printf` does **not** add one by itself, `Println` does) |

> [!TIP]
> ✅ **Good practice 2.3**
> Use `%q` when you print text you are debugging. `%s` with an empty text prints nothing and it looks
> like the line failed; `%q` shows `""` and it is clearly visible that the text is empty.

### 2.5 Functions

A **function** is a named block of code, to which you can give data and which can return a result to you.
You already used two: `fmt.Println` and `fmt.Printf`.

**Fig. 2.8** | Defining and calling functions.

```go
 1  // fig02_08.go
 2  // Define funciones que reciben y devuelven valores.
 3  package main
 4
 5  import "fmt"
 6
 7  // sumar recibe dos enteros y devuelve su suma.
 8  func sumar(a int, b int) int {
 9      return a + b
10  }
11
12  // etiqueta arma un texto descriptivo del servicio.
13  func etiqueta(nombre string, puerto int) string {
14      return fmt.Sprintf("%s:%d", nombre, puerto)
15  }
16
17  // esSeguro dice si el puerto corresponde a HTTPS.
18  func esSeguro(puerto int) bool {
19      return puerto == 443
20  }
21
22  func main() {
23      fmt.Println("2 + 3 =", sumar(2, 3))
24      fmt.Println(etiqueta("catalogo", 443))
25      fmt.Println("catalogo es seguro:", esSeguro(443))
26      fmt.Println("local es seguro:", esSeguro(8080))
27  }
```

```bash
$ go run fig02_08.go
2 + 3 = 5
catalogo:443
catalogo es seguro: true
local es seguro: false
```

**Anatomy of line 8**, which is where everything is:

```
func  sumar  (a int, b int)  int  {
 │      │          │          │
 │      │          │          └── lo que DEVUELVE
 │      │          └── lo que RECIBE (los parámetros), con su tipo
 │      └── el nombre
 └── palabra clave: «voy a definir una función»
```

**Line 14: `fmt.Sprintf`.** It is like `Printf` but instead of printing, it **returns** the assembled text.
The `S` is for *string*. It is used a great deal.

**Line 19: `puerto == 443`.** The double equals sign **compares** and gives `true` or `false`. A single equals
sign `=` **assigns**. Confusing them is a classic (see "What goes wrong", below).

> [!TIP]
> ✅ **Good practice 2.4**
> When two parameters have the same type you can abbreviate: `func sumar(a, b int) int`. It is idiomatic
> and shorter. Here we write it out in full so the structure is visible.

### 2.6 Returning two values: Go's signature

**This is the most distinctive thing about Go**, and you will write it thousands of times in your career.
Pay attention.

Think of a function that divides two numbers. What does it do if the divisor is zero? It can't return a
number, because the result doesn't exist.

Other languages use **exceptions**: the function "throws" an error and someone further up "catches" it
with `try / catch`. The problem is that **it isn't obvious who catches it or where**, and it is very easy
for nobody to do so and the program dies.

**Go does something else: it returns two values.** The result, and an error.

**Fig. 2.9** | A function that returns a result and a possible error.

```go
 1  // fig02_09.go
 2  // Devuelve dos valores: el resultado y un posible error.
 3  package main
 4
 5  import (
 6      "errors"
 7      "fmt"
 8  )
 9
10  // dividir devuelve el cociente y un error si el divisor es cero.
11  func dividir(a int, b int) (int, error) {
12      if b == 0 {
13          return 0, errors.New("division entre cero")
14      }
15      return a / b, nil
16  }
17
18  func main() {
19      // caso que funciona
20      resultado, err := dividir(10, 2)
21      if err != nil {
22          fmt.Println("error:", err)
23      } else {
24          fmt.Println("10 / 2 =", resultado)
25      }
26
27      // caso que falla
28      resultado, err = dividir(10, 0)
29      if err != nil {
30          fmt.Println("error:", err)
31      } else {
32          fmt.Println("10 / 0 =", resultado)
33      }
34  }
```

```bash
$ go run fig02_09.go
10 / 2 = 5
error: division entre cero
```

**Line 11: `(int, error)`.** The parentheses with two types mean *"this function returns two things"*: an
integer and an error.

**Line 13: `return 0, errors.New(...)`.** When there is a problem, it returns any value (here `0`, which
won't be used) **and** an error describing what happened.

**Line 15: `return a / b, nil`.** When everything goes well, it returns the result **and `nil`**.

🔑 **`nil` means "nothing", "empty", "there is none".** So the question `if err != nil` reads: *"is the
error different from nothing?"*, or in plain English: **"was there an error?"**

**Lines 20-21: the pattern that defines Go.**

```go
resultado, err := dividir(10, 2)
if err != nil {
    // algo salió mal: aquí se atiende
}
// si llegaste aquí, todo bien
```

This `if err != nil` pattern is **the most-written line in the entire history of Go**. It appears
everywhere, and some people criticize it for being repetitive.

> [!NOTE]
> 🔧 **Software engineering observation 2.3**
> The criticism is true: it is verbose. The advantage is that **it is impossible to ignore an error by
> oversight**, because it is right there, in your face, on the line after the call. With exceptions you
> can forget a `catch` and not find out until the program dies in production. Go traded brevity for
> safety, on purpose.

> [!TIP]
> 🧪 **Testing and debugging tip 2.2**
> Notice line 28: it says `resultado, err = dividir(10, 0)` with `=`, not `:=`. That is because the two
> variables **already exist** from line 20. If you put `:=` there, the compiler would tell you.

### 2.7 Putting it all together

**Fig. 2.10** | Complete program that uses everything in the lesson.

```go
 1  // fig02_10.go
 2  // Clasifica servicios por su tiempo de respuesta.
 3  package main
 4
 5  import "fmt"
 6
 7  // clasificar traduce milisegundos a una etiqueta legible.
 8  func clasificar(ms int) string {
 9      if ms < 0 {
10          return "sin medir"
11      } else if ms < 500 {
12          return "rapido"
13      } else if ms < 3000 {
14          return "lento"
15      }
16      return "muy lento"
17  }
18
19  // reporte arma una linea del informe.
20  func reporte(nombre string, puerto int, ms int) string {
21      return fmt.Sprintf("%-10s :%-5d %6dms  %s",
22          nombre, puerto, ms, clasificar(ms))
23  }
24
25  func main() {
26      fmt.Println("SERVICIO   PUERTO  TIEMPO   ESTADO")
27      fmt.Println("-------------------------------------")
28      fmt.Println(reporte("catalogo", 443, 142))
29      fmt.Println(reporte("pagos", 443, 87))
30      fmt.Println(reporte("inventario", 443, 2310))
31      fmt.Println(reporte("reportes", 3001, 4500))
32      fmt.Println(reporte("viejo", 8080, -1))
33  }
```

```bash
$ go run fig02_10.go
SERVICIO   PUERTO  TIEMPO   ESTADO
-------------------------------------
catalogo   :443      142ms  rapido
pagos      :443       87ms  rapido
inventario :443     2310ms  lento
reportes   :3001    4500ms  muy lento
viejo      :8080      -1ms  sin medir
```

This program is already the direct ancestor of the `revisor`: `clasificar` and `reporte` are, in
miniature, what in lesson 3 will become a method on a `Servicio` struct, and in lesson 6 will run for
several services **at the same time**.

**Line 21: `%-10s` and `%-5d`.** The number says **how many spaces wide** to reserve, and the minus sign
means **left-aligned**. That way the columns come out straight. Without it, the table would come out
crooked.

> [!TIP]
> ✅ **Good practice 2.5**
> When a call is very long, break it across several lines as in lines 21-22. Go allows it as long as the
> comma stays at the end of the previous line.

> [!NOTE]
> 🚀 **Performance tip 2.1**
> `Sprintf` is convenient but it is not free: it builds new text in memory each time. For five lines it
> doesn't matter at all. If someday you build thousands, there is `strings.Builder`, which is much faster.
> **Don't optimize before measuring**: this is information for later, not something to do today.

---

## The error you will see

Almost all Go error messages have the same form:

```
<archivo>:<línea>:<columna>: <qué esperaba el compilador y qué encontró>
```

Learning to read that line is half the work. Two examples you already saw in this lesson:

```
./fig02_05.go:8:5: declared and not used: edad
```

It reads: *"in file `fig02_05.go`, line 8, column 5, you declared `edad` and never used it"*. The fix is
simple: either you use the variable, or you delete it, or —if you really need it declared but are not using
it yet— you replace it with the blank identifier `_`.

```
./fig02_06.go:8:14: cannot use "muchos" (untyped string constant) as int value in assignment
```

It reads: *"on line 8, column 12, you tried to use the text `"muchos"` where an `int` was expected, and I
won't let you"*. The fix is not easy to guess either if you have never seen it: **there is no automatic
conversion** between `string` and `int` in an assignment like this. You would have to convert explicitly
(with `strconv`, which we will see later) or, more likely, realize that you mixed up the type by mistake.

And a third one, the most common in the first week, which appears if you reuse `:=` on a variable that
already exists:

```
./fig02_04.go:10:9: no new variables on left side of :=
```

It reads: *"on the left side of the `:=` there is no new variable"* — they all already existed, so Go
doesn't know what to declare. The fix: use `=` instead of `:=` when you already declared the variable
before.

**The general rule:** the Go compiler almost never says "something went wrong" in the abstract. It says
the file, the line, the column, and a phrase that —even if it sounds odd the first time— describes exactly
the problem. Reading it in full, without panicking, solves most cases without having to search for
anything online.

## What goes wrong

- **Writing `println` in lowercase instead of `Println`.** Go is **case-sensitive**: `Println` and
  `println` are different names (`println` is a built-in compiler function, meant for debugging Go
  itself, not your program). The error is `undefined: fmt.println`.
- **Using `import "fmt"` and not using it, or using `fmt.Println` without the `import`.** The first gives
  `"fmt" imported and not used`; the second, `undefined: fmt`. Both are the same kind of error: you told
  the compiler something that doesn't match what you did.
- **Confusing `=` (assign) with `==` (compare)** inside an `if`. In C this would compile and do something
  unexpected (assign the value and carry on). **In Go it is a compilation error** —
  `cannot use puerto = 443 (...) as value` —, which is good news: it stops you before the program does
  something you didn't ask for.
- **Forgetting the `\n` in `Printf`.** The program **does compile and does run**, but everything comes out
  stuck together on a single line. `Println` moves to a new line automatically; `Printf` doesn't. It is a
  silent error, not one the compiler points out to you.
- **Ignoring the error with `_`**, like this: `resultado, _ := dividir(10, 0)`. **This also compiles**, and
  the program carries on with `resultado` equal to `0` as if everything were fine. It is the fastest way to
  create an error that is impossible to find later, because there is no message or crash: the program
  simply continues with incorrect data. **If you ever discard an error with `_`, let it be a conscious,
  commented decision, never a reflex to make the compiler stop complaining.**

## Exercises

### Review questions

Answer without looking at the answers. They are at the end of this section.

**2.1** What is the difference between a compiled language and an interpreted one, and which of the two
is Go?

**2.2** Why does the function where the program starts have to be called exactly `main`?

**2.3** What is the difference between `:=` and `=`?

**2.4** What error does Go give if you declare a variable and don't use it? Why do you think it does so?

**2.5** What is the *zero value* and what is it for `string`, `int`, and `bool`?

**2.6** What does this program print?

```go
package main

import "fmt"

func main() {
    var n int
    var s string
    fmt.Printf("[%d] [%q]\n", n, s)
}
```

**2.7** What is the difference between `Println`, `Printf`, and `Sprintf`?

**2.8** What does `nil` mean and what does `if err != nil` ask?

**2.9** This program has **three** errors that prevent it from compiling. Find them without running it.

```go
package main

func main() {
    nombre := "catalogo"
    puerto := 443
    puerto := 8080
    fmt.Println(nombre)
}
```

#### Answers

**2.1** An **interpreted** language is translated line by line on each run; a **compiled** one is
translated entirely just once and produces an executable. **Go is compiled.**

**2.2** Because it is a convention of the language: when running a program, Go looks for the `main`
function of the `main` package in order to start. With another name it finds nowhere to begin.

**2.3** `:=` **declares** a new variable and assigns it a value; `=` only **assigns** to one that already
exists.

**2.4** `declared and not used`. It does so because an unused variable almost always indicates an error
—a misspelled name or leftovers from deleted code— and Go prefers to stop you rather than leave dead code.

**2.5** It is the value that a variable declared without a value automatically receives. `string` → `""`,
`int` → `0`, `bool` → `false`. It guarantees that there is **never** garbage or "undefined".

**2.6** `[0] [""]` — the zero value of `int` and of `string`, and `%q` shows the quotes.

**2.7** `Println` prints and moves to a new line. `Printf` prints with a `%` template and does **not** move
on by itself. `Sprintf` uses the same template but **returns** the text instead of printing it.

**2.8** `nil` means "nothing" / "empty". `if err != nil` asks *"is the error not empty?"*, that is,
**"was there an error?"**.

**2.9** The three:
1. `import "fmt"` is missing, and `fmt.Println` is used.
2. `puerto := 8080` uses `:=` on a variable that already exists → it must be `=`.
3. `puerto` is declared and **never used** (only `nombre` is printed) → `declared and not used`.

### Code exercises

> 🔴 **These seven exercises don't have a reference solution yet** — it will be added in a later
> delivery. No solution code was invented, so as not to publish an example without running it first.

**2.10** Write a program that declares your name, your age, and whether you study or work, and prints
them on three lines with `Printf`, using the right verb for each type.

**2.11** Write `func celsiusAFahrenheit(c float64) float64` and test it with 0, 37, and 100. The formula is
`f = c*9/5 + 32`. **Careful:** if you write `c*9/5` with integers the result is truncated. Why doesn't
that happen here?

**2.12** Write `func esPar(n int) bool`. Hint: the `%` operator gives the remainder of a division, so
`n % 2 == 0` is true for even numbers.

**2.13** Write `func raizCuadrada(n float64) (float64, error)` that returns an error if `n` is negative.
Use `math.Sqrt` (you need `import "math"`). Test it with 16 and with -4, handling the error in both
cases.

**2.14** Take **Fig. 2.10** and add a column that says `SI` or `NO` depending on whether the port is 443.
You need a new function and to adjust the `Sprintf` template.

**2.15 (Find the error)** Each of these fragments has a problem. Say what it is without compiling, and
then compile to confirm:

```go
// (a)
resultado := dividir(10, 2)

// (b)
func sumar(a, b) int { return a + b }

// (c)
var x int = "5"

// (d)
fmt.Printf("el total es %d")
```

**2.16 (Course project)** This is the first step of the program you will build throughout the course.
Write a program that:

1. Has a function `revisar(nombre string, puerto int, ms int) string` that returns a report line.
2. Has a function `clasificar(ms int) string`.
3. Prints a header and five services.
4. At the end, prints how many of the five were `"rapido"`. You will need a counter variable and an `if`.

Save it: in lesson 3 you will reorganize it with structs.

## How I know I got it

- You compiled and ran, yourself and without copy/pasting, figures 2.1 to 2.10 of this lesson, and got
  exactly the output shown.
- You can explain, without looking at the text, why Go is a compiled language and what practical
  consequence that has for catching errors before production.
- You know by heart when to use `:=` and when `=`, and what error the compiler gives if you get it wrong.
- You can name the four basic types (`string`, `int`, `float64`, `bool`) and the zero value of each one,
  without consulting the table.
- You wrote and ran your own version of exercise **2.16** (the first step of the `revisor`): it compiles,
  prints the header, the five rows, and the count of fast services.
- You can explain in two sentences why `if err != nil` appears everywhere in Go code, and what problem it
  avoids compared to the exceptions of other languages.

Open [`bitacora.md`](bitacora.md) and write down **two things**: the compiler error that was hardest for
you to understand, and something that surprised you. In three weeks it will seem obvious to you and you
won't remember why it was hard — and that is exactly what is worth having written down.

---

## Summary

- Go is a **compiled** language: it translates all your code once and produces an independent executable.
- Every file belongs to a **package**; the **`main`** package is the only one that produces an executable.
- Execution starts in the **`main`** function.
- **`import`** brings in packages; Go loads nothing by default.
- **`:=`** declares and assigns; **`=`** only assigns.
- Go **does not compile** if you declare a variable and don't use it, nor if you mix types.
- The four basic types are **`string`**, **`int`**, **`float64`**, and **`bool`**.
- The **zero value** guarantees that every variable is born with a valid value: `""`, `0`, `0`, `false`.
- **`Println`** prints with a line break; **`Printf`** uses a `%` template; **`Sprintf`** returns the text
  instead of printing it.
- A Go function can return **several values**, and the standard pattern is to return
  **`(result, error)`**.
- **`nil`** means "nothing". **`if err != nil`** is the idiomatic way to check whether there was an error,
  and the most-written line in Go.
- `go run` compiles and runs without leaving a file; `go build` leaves the executable.

---

## Further reading

1. **[A Tour of Go](https://go.dev/tour/)** — official, interactive. Do the "Basics" part as far as you
   have gotten in this lesson.
2. **[Go by Example — Variables](https://gobyexample.com/variables)** and
   **[Go by Example — Multiple Return Values](https://gobyexample.com/multiple-return-values)** — the
   same `(result, error)` pattern with minimal programs that run.
3. **[Effective Go](https://go.dev/doc/effective_go)** — you don't need to read it all yet, but consult
   it if something in this lesson felt like an arbitrary rule: that is where the why is.

If something wasn't clear: read the complete error message (Go usually says the line, the column, and
what it expected), look up the concept in Go by Example, and write the question down in the logbook
**even if you don't resolve it**. A written question can be resolved later; a forgotten one can't.

### Terms from this lesson

| | |
|---|---|
| **compiler** | program that translates source code into machine instructions |
| **compilation error** | error that prevents producing the executable; detected before running |
| **function** | named block of code that receives parameters and can return values |
| **blank identifier (`_`)** | symbol that discards a value on purpose |
| **`int`, `float64`, `string`, `bool`** | the four basic types |
| **compiled / interpreted language** | translates everything once / line by line on each run |
| **`main` (function)** | the program's entry point |
| **`main` (package)** | the only package that produces an executable |
| **`nil`** | absence of a value |
| **package** | grouping of related code |
| **parameter** | piece of data that a function receives |
| **`Printf` / `Println` / `Sprintf`** | print with a template / print with a line break / return text |
| **type** | kind of data a variable can store |
| **zero value** | valid value that every variable declared without a value receives |
| **variable** | named space in memory that stores a piece of data |
| **formatting verb (`%d`, `%s`, `%q`, `%t`, `%v`)** | blank in a `Printf` template |

---

**Previous:** [Lesson 1 — Installing Go](01-instalacion.md) ·
**Next:** [Lesson 3 — Structs, methods, errors, and interfaces](03-errores-interfaces.md)
