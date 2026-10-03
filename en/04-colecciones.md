# Lesson 4 — Collections: slices and maps

**Duration:** 90 minutes (or 2 sessions of 45).

**By the end you will be able to:**

- Tell an **array** from a **slice** and say why you will always use the second one.
- Explain why copying a slice **doesn't copy the data**, and when that is going to bite you.
- Use a **map** and know why iterating over it comes out unordered on purpose.
- Iterate over collections with `for range` and know what the `_` means.
- Sort results so your program always produces the same output.
- Read a file and turn it into a list of data.

---

## Why it matters

**This is the lesson where almost everyone stumbles.** Not because it is hard, but because there is a trap
that no course explains until it has already bitten you: in Go, copying a slice doesn't copy its data, and
that behavior doesn't produce any error — it corrupts data silently.

The `revisor` needs collections for two very concrete things: a **list** of services you are going to
query (those are slices) and a **report** that associates each service with its state (those are maps).
Without slices and maps there is no program: you can't have "several services" or "the state of each one"
with what you saw up to Lesson 3.

This is the map of the lesson:

| | |
|---|---|
| **4.1** | Arrays: the ones you will hardly use |
| **4.2** | Slices: the ones you will always use |
| **4.3** | The shared-memory trap |
| **4.4** | Maps |
| **4.5** | The random order, and why it is deliberate |
| **4.6** | Reading a real file |

---

## The concepts

### 4.1 Arrays: the ones you will hardly use

An **array** is a list of **fixed** size. The size is part of the type.

**Fig. 4.1** | Arrays, so you can recognize them.

```go
 1  // fig04_01.go
 2  // Muestra arreglos, que en Go casi no se usan directamente.
 3  package main
 4
 5  import "fmt"
 6
 7  func main() {
 8      var puertos [3]int              // tres enteros, todos en 0
 9      fmt.Println("recien creado:", puertos)
10
11      puertos[0] = 443
12      puertos[1] = 8080
13      fmt.Println("con valores:", puertos)
14      fmt.Println("cuantos caben:", len(puertos))
15
16      otros := [3]int{80, 443, 8443}
17      fmt.Println("otro arreglo:", otros)
18
19      copia := otros                  // los arreglos SÍ se copian completos
20      copia[0] = 999
21      fmt.Println("original:", otros)
22      fmt.Println("copia:   ", copia)
23  }
```

```bash
$ go run fig04_01.go
recien creado: [0 0 0]
con valores: [443 8080 0]
cuantos caben: 3
otro arreglo: [80 443 8443]
original: [80 443 8443]
copia:    [999 443 8443]
```

Notice two things:

**Line 8.** The array is born with zero values, like everything in Go. There is no garbage.

**Lines 19-22.** When you assign an array to another variable, **it is copied in full**: modifying the copy
doesn't touch the original. Remember this, because with slices the same thing does **not** happen, and
that is where this lesson's trap lies.

> [!NOTE]
> 🔧 **Software engineering observation 4.1**
> `[3]int` and `[4]int` are **different types**. A function that receives `[3]int` doesn't accept a
> `[4]int`. That makes arrays impractical, and it is the reason almost nobody in Go uses them directly:
> slices are used instead, which do grow. Arrays are there because slices are built on top of them.

### 4.2 Slices: the ones you will always use

A **slice** is a list of **variable** size. It is what in other languages you would call a list or a
dynamic array.

**Fig. 4.2** | Creating and growing a slice.

```go
 1  // fig04_02.go
 2  // Crea un slice y lo hace crecer con append.
 3  package main
 4
 5  import "fmt"
 6
 7  type Servicio struct {
 8      Nombre string
 9      Puerto int
10  }
11
12  func main() {
13      // forma 1: vacio, y crece
14      var nombres []string
15      fmt.Printf("vacio: %v · largo %d · es nil: %t\n", nombres, len(nombres), nombres == nil)
16
17      nombres = append(nombres, "catalogo")
18      nombres = append(nombres, "pagos", "inventario")
19      fmt.Printf("con datos: %v · largo %d\n", nombres, len(nombres))
20
21      // forma 2: con valores desde el inicio
22      servicios := []Servicio{
23          {Nombre: "catalogo", Puerto: 443},
24          {Nombre: "pagos", Puerto: 443},
25          {Nombre: "local", Puerto: 8080},
26      }
27      fmt.Println("cuantos servicios:", len(servicios))
28
29      // acceder por posicion, empezando en 0
30      fmt.Println("el primero:", servicios[0].Nombre)
31      fmt.Println("el ultimo: ", servicios[len(servicios)-1].Nombre)
32  }
```

```bash
$ go run fig04_02.go
vacio: [] · largo 0 · es nil: true
con datos: [catalogo pagos inventario] · largo 3
cuantos servicios: 3
el primero: catalogo
el ultimo:  local
```

**Line 15.** An uninitialized slice is **`nil`** and its length is `0`. But notice something important:

> [!TIP]
> ✅ **Good practice 4.1**
> **A `nil` slice can be used with `append`, with `len`, and with `for range` without any problem.** That
> is why there is no need to initialize it with `[]string{}`: `var nombres []string` and straight to
> `append`. It is the idiomatic way, and another sign that Go's zero value is designed to be useful.

**Line 17: `nombres = append(nombres, "catalogo")`.** Notice that you have to **reassign**. `append`
doesn't modify the slice: **it returns a new one**, and if you don't store the result, you lose it.

> [!WARNING]
> ⚠️ **Common programming error 4.1**
> Writing `append(nombres, "x")` without assigning the result. **The compiler detects it** because the
> returned value is not used (`append(...) evaluated but not used`), so here Go saves you. But if you
> assign it to another variable by mistake, it compiles and the original slice doesn't change.

**Line 31: `servicios[len(servicios)-1]`.** In Go there is no negative index as in Python: for the last
element you compute it. And if you go past the end, the program blows up:

**Fig. 4.3** | Going out of range.

```go
 1  // fig04_03.go — este programa COMPILA pero truena
 2  package main
 3
 4  import "fmt"
 5
 6  func main() {
 7      puertos := []int{443, 8080}
 8      fmt.Println(puertos[0])
 9      fmt.Println(puertos[5])      // solo hay 2 elementos
10  }
```

```bash
$ go run fig04_03.go
443
panic: runtime error: index out of range [5] with length 2

goroutine 1 [running]:
main.main()
	/tmp/fig04_03.go:9 +0x1c
exit status 2
```

> [!WARNING]
> ⚠️ **Common programming error 4.2**
> The `index out of range` is Go's second most frequent panic, after the nil pointer. **Notice how useful
> the message is:** it tells you the index you asked for (`[5]`), the real length (`length 2`), and the
> line. Before indexing something that comes from outside —a file, a request, an argument— check `len`.

### 4.3 The shared-memory trap

**This is the most important section of the lesson.** It is a behavior that surprises everyone, that
doesn't produce any error, and that can corrupt data silently.

#### 4.3.1 The problem

**Fig. 4.4** | Copying a slice does **not** copy the data.

```go
 1  // fig04_04.go
 2  // Demuestra que asignar un slice NO copia sus datos.
 3  package main
 4
 5  import "fmt"
 6
 7  func main() {
 8      original := []int{443, 8080, 3000}
 9      copia := original              // parece una copia...
10
11      copia[0] = 999                 // modifico la "copia"
12
13      fmt.Println("original:", original)
14      fmt.Println("copia:   ", copia)
15  }
```

```bash
$ go run fig04_04.go
original: [999 8080 3000]
copia:    [999 8080 3000]
```

**You modified `copia` and `original` changed.** There was no error, there was no warning, and both slices
show the same thing.

Compare it with Fig. 4.1, where the same code with an **array** did copy. The difference is fundamental.

#### 4.3.2 Why it happens

A slice **doesn't contain** the data: it is a window that points to it. Inside it stores three things:

| | |
|---|---|
| a **pointer** | to where the real data is |
| the **length** (`len`) | how many elements it has now |
| the **capacity** (`cap`) | how many fit before it has to move |

When you write `copia := original`, Go copies **those three things** — not the data. Both variables end up
pointing to the **same** place.

It is like giving you the address of a house instead of building you an identical one: if you paint
"your" house, mine changes, because it is the same one.

#### 4.3.3 The solution

**Fig. 4.5** | Copying for real.

```go
 1  // fig04_05.go
 2  // Las dos formas de copiar un slice de verdad.
 3  package main
 4
 5  import "fmt"
 6
 7  func main() {
 8      original := []int{443, 8080, 3000}
 9
10      // forma 1: make + copy
11      copia1 := make([]int, len(original))
12      copy(copia1, original)
13      copia1[0] = 111
14
15      // forma 2: append sobre un slice nil (mas corta)
16      copia2 := append([]int(nil), original...)
17      copia2[1] = 222
18
19      fmt.Println("original:", original)
20      fmt.Println("copia1:  ", copia1)
21      fmt.Println("copia2:  ", copia2)
22  }
```

```bash
$ go run fig04_05.go
original: [443 8080 3000]
copia1:   [111 8080 3000]
copia2:   [443 222 3000]
```

**Now they really are independent.**

**Line 11: `make([]int, len(original))`.** `make` creates a slice with reserved space. It is the way to
say "I want a slice of this length, with its zero values".

**Line 16: `original...`.** Those three dots mean "spread the elements one by one". Without them you would
be trying to append the whole slice as a single element, and it doesn't compile.

#### 4.3.4 And the truly treacherous part: `append`

**Fig. 4.6** | The same code behaves differently depending on the capacity.

```go
 1  // fig04_06.go
 2  // append comparte o no memoria segun la capacidad disponible.
 3  package main
 4
 5  import "fmt"
 6
 7  func main() {
 8      // CASO A: hay capacidad de sobra (cap 5, largo 3)
 9      a := make([]int, 3, 5)
10      a[0], a[1], a[2] = 1, 2, 3
11      fmt.Printf("A: len=%d cap=%d %v\n", len(a), cap(a), a)
12
13      b := append(a, 99)          // cabe: NO se muda, comparte memoria
14      b[0] = 777
15      fmt.Println("   tras modificar b, a vale:", a, "← cambió")
16
17      // CASO B: no hay capacidad (cap 3, largo 3)
18      c := make([]int, 3, 3)
19      c[0], c[1], c[2] = 1, 2, 3
20      fmt.Printf("B: len=%d cap=%d %v\n", len(c), cap(c), c)
21
22      d := append(c, 99)          // NO cabe: se muda a otra memoria
23      d[0] = 777
24      fmt.Println("   tras modificar d, c vale:", c, "← NO cambió")
25  }
```

```bash
$ go run fig04_06.go
A: len=3 cap=5 [1 2 3]
   tras modificar b, a vale: [777 2 3] ← cambió
B: len=3 cap=3 [1 2 3]
   tras modificar d, c vale: [1 2 3] ← NO cambió
```

**Read that output twice.** It is **the same code** —`append` and then modify— and the result is different
depending on how much spare capacity there was.

> [!WARNING]
> 🔴 **Common programming error 4.3 — the worst one in this lesson**
> Storing a slice you received as a parameter, or keeping the result of `append` assuming it is
> independent. **The behavior depends on the capacity**, which you almost never control and which changes
> depending on how much the slice has grown before. In other words: **your program can work in the tests
> and corrupt data in production**, with the same code.
>
> **The rule that saves you:** if you are going to **store** a slice you received from someone else, copy
> it first. If you are only going to read it, there is no need.

> [!NOTE]
> 🚀 **Performance tip 4.1**
> When you know how many elements you are going to add, give it the capacity from the start:
> `make([]Estado, 0, len(servicios))`. That way `append` doesn't have to move or copy anything while it
> grows. With small lists it is irrelevant; with thousands of elements, it shows. **And measure it before
> believing it.**

### 4.4 Maps

A **map** stores key-value pairs. It is what in other languages is called a dictionary or an associative
table.

**Fig. 4.7** | Creating, writing, and reading a map.

```go
 1  // fig04_07.go
 2  // Operaciones basicas con un map.
 3  package main
 4
 5  import "fmt"
 6
 7  type Estado struct {
 8      Codigo int
 9      Ms     int
10  }
11
12  func main() {
13      // se crea con make, o vacio con {}
14      estados := make(map[string]Estado)
15
16      estados["catalogo"] = Estado{Codigo: 200, Ms: 142}
17      estados["pagos"] = Estado{Codigo: 200, Ms: 87}
18      estados["inventario"] = Estado{Codigo: 503, Ms: 2310}
19
20      fmt.Println("cuantos:", len(estados))
21
22      // leer una llave que SÍ existe
23      e := estados["catalogo"]
24      fmt.Println("catalogo:", e.Codigo, e.Ms)
25
26      // 🔑 leer una que NO existe: devuelve el valor cero, SIN error
27      fantasma := estados["no-existe"]
28      fmt.Printf("no-existe: %+v  ← el valor cero\n", fantasma)
29
30      // la forma correcta de preguntar: el segundo valor
31      if e, hay := estados["pagos"]; hay {
32          fmt.Println("pagos si esta:", e.Codigo)
33      }
34      if _, hay := estados["no-existe"]; !hay {
35          fmt.Println("no-existe NO esta")
36      }
37
38      // borrar
39      delete(estados, "inventario")
40      fmt.Println("tras borrar:", len(estados))
41  }
```

```bash
$ go run fig04_07.go
cuantos: 3
catalogo: 200 142
no-existe: {Codigo:0 Ms:0}  ← el valor cero
pagos si esta: 200
no-existe NO esta
tras borrar: 2
```

**Lines 27-28: here is what you need to understand.** Reading a key that doesn't exist **doesn't give an
error**: it returns the zero value of the type. For an `int` that is `0`, which could be a legitimate
value.

> [!WARNING]
> ⚠️ **Common programming error 4.4**
> Using `m[llave]` without checking whether it exists, when the zero value is indistinguishable from real
> data. If you store `map[string]int` with response times and read a nonexistent key, you get `0` — which
> looks like "it responded instantly" instead of "we didn't measure it". **Always use the two-value form
> (`v, hay := m[k]`) when absence means something other than zero.**

> [!TIP]
> ✅ **Good practice 4.2**
> Call the second variable `hay`, `existe`, or `ok`. In Go code you will see `ok` a lot, and it is the
> convention: `if v, ok := m[k]; ok { … }`.

#### 4.4.1 The nil map: the only one that bites

**Fig. 4.8** | A `nil` map can be read but not written.

```go
 1  // fig04_08.go — este programa truena al escribir
 2  package main
 3
 4  import "fmt"
 5
 6  func main() {
 7      var m map[string]int          // nil: NO inicializado
 8
 9      fmt.Println("leer de un map nil:", m["x"])    // esto funciona
10      fmt.Println("su largo:", len(m))              // esto tambien
11
12      m["x"] = 1                                     // esto truena
13  }
```

```bash
$ go run fig04_08.go
leer de un map nil: 0
su largo: 0
panic: assignment to entry in nil map

goroutine 1 [running]:
main.main()
	/tmp/fig04_08.go:12 +0x3c
exit status 2
```

> [!WARNING]
> 🔴 **Common programming error 4.5**
> **This is Go's great asymmetry and you have to memorize it:** a `nil` **slice** can be used with
> `append` without problem; a `nil` **map** blows up when you write to it. Maps **have to be created** with
> `make(map[K]V)` or `map[K]V{}`. And notice the tricky part: reading from the nil map **does work**, so
> the program can move along for a while before blowing up.

### 4.5 The random order, and why it is deliberate

**Fig. 4.9** | The same program, two runs, two orders.

```go
 1  // fig04_09.go
 2  // El recorrido de un map sale en orden distinto cada vez.
 3  package main
 4
 5  import "fmt"
 6
 7  func main() {
 8      puertos := map[string]int{
 9          "catalogo":   443,
10          "pagos":      443,
11          "inventario": 8080,
12          "reportes":   3001,
13      }
14
15      for nombre, puerto := range puertos {
16          fmt.Printf("%s=%d ", nombre, puerto)
17      }
18      fmt.Println()
19  }
```

```bash
$ go run fig04_09.go
reportes=3001 catalogo=443 pagos=443 inventario=8080

$ go run fig04_09.go
pagos=443 inventario=8080 reportes=3001 catalogo=443

$ go run fig04_09.go
inventario=8080 reportes=3001 catalogo=443 pagos=443
```

**Three runs, three different orders.** And it is not a defect: **Go does it on purpose**, randomizing the
starting point on each iteration.

> [!NOTE]
> 🔧 **Software engineering observation 4.2**
> Why bother the programmer with this? Because **the order of a map was never guaranteed**, neither in Go
> nor in most languages. If Go left it "almost always the same", there would be programs that work for
> years and one day, as the data grows, change order and break — and nobody would understand why. **By
> randomizing it, it forces you to find out today.** It is the same philosophy as refusing to compile with
> an unused variable: an early, annoying error is worth more than a late, incomprehensible one.

#### 4.5.1 Sorting so the output is stable

If you need order, you ask for it explicitly.

**Fig. 4.10** | Iterating over a map in alphabetical order.

```go
 1  // fig04_10.go
 2  // Ordena las llaves para producir una salida estable.
 3  package main
 4
 5  import (
 6      "fmt"
 7      "sort"
 8  )
 9
10  type Estado struct {
11      Codigo int
12      Ms     int
13  }
14
15  func main() {
16      estados := map[string]Estado{
17          "reportes":   {Codigo: 0, Ms: 0},
18          "catalogo":   {Codigo: 200, Ms: 142},
19          "inventario": {Codigo: 503, Ms: 2310},
20          "pagos":      {Codigo: 200, Ms: 87},
21      }
22
23      // 1. saca las llaves a un slice
24      nombres := make([]string, 0, len(estados))
25      for nombre := range estados {
26          nombres = append(nombres, nombre)
27      }
28
29      // 2. ordenalas
30      sort.Strings(nombres)
31
32      // 3. recorre el slice ordenado, no el map
33      fmt.Println("SERVICIO     CODIGO  TIEMPO")
34      fmt.Println("-----------------------------")
35      for _, nombre := range nombres {
36          e := estados[nombre]
37          fmt.Printf("%-12s %6d  %5dms\n", nombre, e.Codigo, e.Ms)
38      }
39  }
```

```bash
$ go run fig04_10.go
SERVICIO     CODIGO  TIEMPO
-----------------------------
catalogo        200    142ms
inventario      503   2310ms
pagos           200     87ms
reportes          0      0ms
```

**And now, yes: the same output, always.** Run the program ten times and it doesn't change.

**Line 24: `make([]string, 0, len(estados))`.** Length `0`, capacity `len(estados)`. You know how many keys
you are going to add, so the space is reserved all at once.

**Line 25: `for nombre := range estados`.** Over a map, `range` gives **the key** in the first variable.
Over a slice it gives **the index**. It is a difference worth keeping in mind:

| Collection | first variable | second variable |
|---|---|---|
| slice | the **index** (0, 1, 2…) | the element |
| map | the **key** | the value |

**Line 35: `for _, nombre := range nombres`.** Here `nombres` is a slice, so the first variable is the
index, and I don't need it: that is why `_`.

> [!TIP]
> ✅ **Good practice 4.3**
> **If your program prints results, sort them.** Output that changes order between runs is impossible to
> compare, impossible to test automatically, and confuses whoever reads it. This three-step pattern
> —extract keys, sort, iterate— is idiomatic and you are going to use it often.

### 4.6 Reading a real file

Until now the data has been written in the program. That is no good: every change requires recompiling.

**Fig. 4.11** | Reading services from a text file.

First the data file, `servicios.txt`:

```
catalogo https://catalogo.example.com 443
pagos https://pagos.example.com 443
inventario https://inventario.example.com 8080
```

And the program:

```go
 1  // fig04_11.go
 2  // Lee los servicios desde un archivo de texto.
 3  package main
 4
 5  import (
 6      "fmt"
 7      "os"
 8      "strconv"
 9      "strings"
10  )
11
12  type Servicio struct {
13      Nombre string
14      URL    string
15      Puerto int
16  }
17
18  // Cargar lee el archivo y devuelve los servicios, o un error.
19  func Cargar(ruta string) ([]Servicio, error) {
20      datos, err := os.ReadFile(ruta)
21      if err != nil {
22          return nil, fmt.Errorf("leyendo %s: %w", ruta, err)
23      }
24
25      var servicios []Servicio
26      lineas := strings.Split(strings.TrimSpace(string(datos)), "\n")
27
28      for i, linea := range lineas {
29          linea = strings.TrimSpace(linea)
30          if linea == "" || strings.HasPrefix(linea, "#") {
31              continue                       // salta vacias y comentarios
32          }
33
34          campos := strings.Fields(linea)     // separa por espacios
35          if len(campos) != 3 {
36              return nil, fmt.Errorf("%s linea %d: esperaba 3 campos, hay %d",
37                  ruta, i+1, len(campos))
38          }
39
40          puerto, err := strconv.Atoi(campos[2])
41          if err != nil {
42              return nil, fmt.Errorf("%s linea %d: puerto invalido %q: %w",
43                  ruta, i+1, campos[2], err)
44          }
45
46          servicios = append(servicios, Servicio{
47              Nombre: campos[0],
48              URL:    campos[1],
49              Puerto: puerto,
50          })
51      }
52      return servicios, nil
53  }
54
55  func main() {
56      servicios, err := Cargar("servicios.txt")
57      if err != nil {
58          fmt.Fprintln(os.Stderr, "error:", err)
59          os.Exit(1)
60      }
61
62      fmt.Printf("cargados %d servicios:\n", len(servicios))
63      for _, s := range servicios {
64          fmt.Printf("  %-12s %-38s :%d\n", s.Nombre, s.URL, s.Puerto)
65      }
66  }
```

```bash
$ go run fig04_11.go
cargados 3 servicios:
  catalogo     https://catalogo.example.com           :443
  pagos        https://pagos.example.com              :443
  inventario   https://inventario.example.com         :8080
```

And let's try what happens when something goes wrong:

```bash
$ go run fig04_11.go            # con el archivo renombrado
error: leyendo servicios.txt: open servicios.txt: no such file or directory
$ echo $?
1
```

**Line 20: `os.ReadFile`.** Reads the whole file and returns its bytes. For configuration files it is
perfect; for a file of several gigabytes there are other ways, because this loads everything into memory.

**Line 26: `string(datos)`.** Converts the bytes to text. `strings.TrimSpace` removes spaces and line
breaks from the beginning and the end — without it, the last empty line would produce a phantom service.

**Line 34: `strings.Fields`.** Splits on spaces, **collapsing repeated ones**. It is more robust than
`strings.Split(linea, " ")`, which with two consecutive spaces would give you an empty field.

**Line 40: `strconv.Atoi`.** Converts text to an integer (*ASCII to integer*). It returns an error if it is
not a number, and **we are checking it**.

**Lines 58-59: the difference between a script and a program.**

> [!TIP]
> ✅ **Good practice 4.4**
> Errors go to **`os.Stderr`**, not to the normal output, and the program ends with a **nonzero exit
> code**. That allows using it in a pipe: `mi-programa | otro-programa` keeps working because the error
> doesn't contaminate the output, and `mi-programa && echo ok` doesn't print "ok" if it failed. **That is
> what separates a program from a script.**

> [!TIP]
> 🧪 **Testing and debugging tip 4.1**
> Notice the error messages on lines 36 and 42: they say **the file, the line number, and what they
> expected**. Compare that with a bare `invalid syntax`. When someone uses your program with a fifty-line
> file, the difference between the two messages is twenty minutes of their life.

---

## The error you will see

This lesson leaves two `panic` messages etched in your memory, because they are the two most frequent ones
of the whole course after the nil pointer:

**`panic: runtime error: index out of range [N] with length M`** (Fig. 4.3). You asked for a position `N`
that doesn't exist in a collection of length `M`. **What it means:** the highest valid index is `M-1`, not
`M`. **How to fix it:** check `len(coleccion)` before indexing, above all if the index comes from outside
(a file, an argument, a request) and you didn't write it yourself in the code.

**`panic: assignment to entry in nil map`** (Fig. 4.8). You tried to write to a map that was never created
— `var m map[string]Estado` leaves `m` as `nil`, and a `nil` map **can be read but not written**. **What
it means:** the `make(map[K]V)` or the `map[K]V{}` that reserves the table inside is missing. **How to fix
it:** create the map before using it as the target of an assignment; if you are only going to read it,
`nil` causes no problem at all.

And there is a third case that does **not** blow up, and that is why it is the most dangerous of the
three: modifying a shared slice (Fig. 4.4 and Fig. 4.6) doesn't produce any error message. The program
keeps running with incorrect data. That is exactly the reason section 4.3 exists.

## What goes wrong

- **Assuming that `copia := original` copies a slice's data.** It copies the pointer, the length, and the
  capacity — not the data. Both variables end up sharing memory (Fig. 4.4).
- **Storing a slice you received as a parameter without copying it first.** The behavior of `append`
  depends on the capacity the slice brings, which whoever receives it almost never controls (Fig. 4.6).
  If you are going to **store** someone else's slice, copy it; if you are only going to read it, there is
  no need.
- **Trusting that a `for range` over a map will always come out in the same order.** It was never
  guaranteed, and Go randomizes it on purpose (Fig. 4.9) so the error shows up today and not the day the
  database grows in production.
- **Reading `m[llave]` without checking `ok` when the zero value is ambiguous.** A `map[string]int` that
  stores response times can't distinguish "responded in 0 ms" from "was never queried" if you don't use
  the two-value form.
- **Indexing something that comes from outside —a file, `os.Args`, a network response— without checking
  `len` first.** It is the most common cause of this lesson's `index out of range`, and it is avoided with
  one line.
- **Writing to a map before creating it with `make` or with `{}`.** The compiler doesn't detect it because
  `var m map[K]V` is valid code; the `panic` doesn't appear until the program runs and tries to write.

---

## Exercises

### Questions

**4.1** What is the difference between `[3]int` and `[]int`?

**4.2** Why are `[3]int` and `[4]int` different types, and what practical consequence does that have?

**4.3** What does this print and why?

```go
a := []int{1, 2, 3}
b := a
b[0] = 99
fmt.Println(a[0])
```

**4.4** Name the two ways to really copy a slice.

**4.5** Why does the same `append` sometimes share memory with the original and sometimes not?

**4.6** What does `m["no-existe"]` return in a `map[string]int`, and why is it dangerous?

**4.7** What is the asymmetry between a `nil` slice and a `nil` map?

**4.8** Why does iterating over a map come out unordered, and how do you get a fixed order?

**4.9** In `for a, b := range x`, what is `a` if `x` is a slice? And if it is a map?

**4.10** What does `strings.Fields` do that `strings.Split(s, " ")` doesn't?

**4.11** Why do errors go to `os.Stderr` and not to the normal output?

**4.12** This program has **two** problems. Find them.

```go
func main() {
    var m map[string]int
    m["catalogo"] = 443
    lista := []int{1, 2, 3}
    fmt.Println(lista[3])
}
```

### Code exercises

**4.13** Write `func Contar(estados map[string]Estado) (ok, fallas int)` that counts how many states have a
200-299 code and how many don't. Go returns several values: use that.

**4.14** Write `func Nombres(servicios []Servicio) []string` that returns only the names, **sorted**.

**4.15 (The trap)** Write a function `func Guardar(s []int)` that stores the slice in a global variable and
then prints it. Call it like this:

```go
datos := make([]int, 3, 10)
datos[0], datos[1], datos[2] = 1, 2, 3
Guardar(datos)
datos = append(datos, 4)
datos[0] = 999
// ahora imprime lo que guardó Guardar
```

**Predict what it is going to print before running it.** Then run it. If you got it right, you understood
section 4.3; if not, read it again — it is the one that costs the most.

**4.16** Extend `Cargar` from Fig. 4.11 so it accepts an optional fourth field with the timeout in
milliseconds. If it is not there, use 5000. **Make sure a file with three fields keeps working.**

**4.17** Write `func Agrupar(servicios []Servicio) map[int][]Servicio` that groups by port. Test that a
port with three services has all of them.

**4.18 (Find the error)** Say what is wrong in each one:

```go
// (a)
var m map[string]Estado
m["catalogo"] = Estado{}

// (b)
append(servicios, nuevo)

// (c)
for i := range servicios {
    fmt.Println(servicios[i+1].Nombre)
}

// (d)
ms := estados["catalogo"].Ms
if ms == 0 {
    fmt.Println("respondio instantaneo")
}
```

**4.19 (Course project)** Take your program to the next step:
1. `Cargar(ruta string) ([]Servicio, error)` that reads the file, like Fig. 4.11.
2. `RevisarTodos` that returns `map[string]Estado`.
3. `Reporte(estados map[string]Estado) string` that produces the table **sorted alphabetically**.
4. A final count: "3 of 5 responded".
5. Make the program accept the file path as an argument: `os.Args[1]`. 🔴 **Check `len(os.Args)`
   before reading it**, or you will get an `index out of range` when someone runs it without arguments.
6. Make it exit with code 0 if everything responded and 1 if something failed.

### Solutions

**4.1** `[3]int` is an **array**: fixed size, part of the type, and it is copied when assigned. `[]int` is
a **slice**: variable size, and when assigned it **shares the data**.

**4.2** Because the size is part of the type. The consequence: a function that receives `[3]int` doesn't
accept a `[4]int`, which makes arrays impractical. That is why slices are used.

**4.3** It prints **99**. `b := a` doesn't copy the data: both variables point to the same memory.

**4.4** `make` + `copy`, or `append([]T(nil), original...)`.

**4.5** Because it depends on the **capacity**. If the slice has spare room, `append` writes right there
and keeps sharing memory; if it doesn't fit, it moves to new memory and stops sharing. Since you almost
never control the capacity, **the same code can behave differently**.

**4.6** It returns `0`, the zero value, **without an error**. It is dangerous because `0` can be legitimate
data, so you can't distinguish "it is zero" from "it isn't there". It is solved with `v, ok := m[k]`.

**4.7** A `nil` slice accepts `append` without problem. A `nil` map **blows up when you write to it**
(`assignment to entry in nil map`), although it can be read. Maps have to be created with `make`.

**4.8** Because its order **was never guaranteed**, and Go randomizes it on purpose so you don't write
programs that depend on an accidental order. For a fixed order: extract the keys into a slice, sort them
with `sort.Strings`, and iterate over the slice.

**4.9** If `x` is a slice, `a` is the **index**. If it is a map, `a` is the **key**.

**4.10** `Fields` splits on spaces **collapsing repeated ones**, so two consecutive spaces don't produce an
empty field. It also treats tabs as a separator.

**4.11** So that the normal output stays clean and can be used in a pipe, and because whoever calls the
program expects to find the errors there. Together with the nonzero exit code, that is what makes a
program automatable.

**4.12** Both: (1) `m` is a `nil` map and writing to it causes `panic: assignment to entry in nil map` —
`m := make(map[string]int)` is missing; (2) `lista[3]` goes out of range, because the valid indexes are 0,
1, and 2 — `panic: index out of range [3] with length 3`.

**Code exercises (4.13 to 4.19):** their verified solutions —compiled and run— will be added in the
delivery in which the chapter of exercises with solutions for the whole course is closed; they are not
published without having run `go build` on each one.

---

## How I know I got it

- You can explain, without looking at the text, why `copia := original` doesn't copy a slice's data — and
  draw it (pointer, length, capacity).
- You correctly predicted what exercise 4.15 prints **before** running it. If you didn't get it right, you
  went over section 4.3 again until you did.
- Your `Cargar` (exercise 4.16) keeps working with three-field files and accepts the optional fourth one.
- Your `Reporte` produces **the same output, in the same order**, in ten runs in a row.
- Your program doesn't blow up with `index out of range` if you run it without arguments: you checked
  `len(os.Args)` before reading `os.Args[1]`.
- You wrote down in [`bitacora.md`](bitacora.md) the result of exercise 4.15: what you predicted and what
  really happened. That is what separates understanding slices from believing you understand them.

---

## Summary

- An **array** (`[3]int`) has a fixed size, the size is part of the type, and **it is copied** when
  assigned. It is hardly ever used directly.
- A **slice** (`[]int`) has a variable size and is what is always used.
- A slice stores **pointer, length, and capacity**: it doesn't contain the data, it points to it.
- 🔴 **Assigning a slice does NOT copy the data**: both variables share memory. To really copy,
  `make` + `copy` or `append([]T(nil), s...)`.
- 🔴 **`append` shares memory or not depending on the available capacity**, so the same code can behave
  differently. If you are going to **store** someone else's slice, copy it.
- A `nil` slice works with `append`, `len`, and `range`. **A `nil` map blows up when you write to it.**
- Reading a nonexistent key from a map **returns the zero value without an error**. Use `v, ok := m[k]`
  when absence matters.
- **Iterating over a map is random on purpose.** For a fixed order: keys into a slice, `sort`, and iterate
  over the slice.
- In `range`, the first variable is the **index** in slices and the **key** in maps. `_` discards.
- `os.ReadFile` reads a whole file; `strings.Fields` splits on spaces collapsing repeated ones;
  `strconv.Atoi` converts text to an integer and **returns an error**.
- Errors go to **`os.Stderr`** and the program exits with a **nonzero exit code**.

---

## Further reading

1. **[Go Slices: usage and internals](https://go.dev/blog/slices-intro)** — the official blog of the Go
   team, explains pointer/length/capacity with drawings. Start here.
2. **[Go maps in action](https://go.dev/blog/maps)** — official blog, covers the random order and the
   `nil` map in more detail than fits in this lesson.
3. **[Go by Example: Slices](https://gobyexample.com/slices) and [Maps](https://gobyexample.com/maps)** —
   minimal code, good as a quick reference.
4. **[Effective Go — Slices](https://go.dev/doc/effective_go#slices)** — the specific section on the
   two-step pattern (`make` + `copy`) and why `append` behaves the way it does.

### Terms from this lesson

| | |
|---|---|
| **`append`** | adds elements to a slice and **returns** the result |
| **array** | fixed-size list; the size is part of the type |
| **capacity (`cap`)** | how many elements fit in a slice before it moves |
| **`copy`** | copies elements between slices |
| **`delete`** | removes a key from a map |
| **`index out of range`** | panic from asking for a position that doesn't exist |
| **length (`len`)** | how many elements it has now |
| **map** | collection of key→value pairs |
| **`make`** | creates slices, maps, and channels with reserved space |
| **`nil map`** | map not created; it can be read but **not** written |
| **`os.Stderr`** | error output, separate from the normal one |
| **slice** | variable-size list; a window over the data |
| **`sort.Strings`** | sorts a slice of text |
| **`strings.Fields`** | splits a text on spaces, collapsing repeated ones |
| **`strconv.Atoi`** | converts text to an integer |
| **two-result value (`v, ok`)** | way of reading a map that distinguishes absence |

---

**Previous:** [Lesson 3 — Structs, errors, and interfaces](03-errores-interfaces.md) ·
**Next:** [Lesson 5 — Modules and tests](05-modulos-y-pruebas.md)
