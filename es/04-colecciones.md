# Lección 4 — Colecciones: slices y maps

**Duración:** 90 minutos (o 2 sesiones de 45).

**Al terminar vas a poder:**

- Distinguir un **arreglo** de un **slice** y decir por qué siempre usarás el segundo.
- Explicar por qué copiar un slice **no copia los datos**, y cuándo eso te va a morder.
- Usar un **map** y saber por qué su recorrido sale desordenado a propósito.
- Recorrer colecciones con `for range` y saber qué significa el `_`.
- Ordenar resultados para que tu programa produzca la misma salida siempre.
- Leer un archivo y convertirlo en una lista de datos.

---

## Por qué importa

**Esta es la lección donde casi todos se tropiezan.** Y no por difícil, sino porque hay una trampa que
ningún curso explica hasta que ya te mordió: en Go, copiar un slice no copia sus datos, y ese
comportamiento no produce ningún error — corrompe datos en silencio.

El `revisor` necesita colecciones para dos cosas muy concretas: una **lista** de servicios que vas a
consultar (eso son slices) y un **reporte** que asocia cada servicio con su estado (eso son maps). Sin
slices ni maps no hay programa: no puedes tener "varios servicios" ni "el estado de cada uno" con lo que
viste hasta la Lección 3.

Este es el mapa de la lección:

| | |
|---|---|
| **4.1** | Arreglos: los que casi no usarás |
| **4.2** | Slices: los que usarás siempre |
| **4.3** | La trampa de la memoria compartida |
| **4.4** | Maps |
| **4.5** | El orden aleatorio, y por qué es deliberado |
| **4.6** | Leer un archivo de verdad |

---

## Los conceptos

### 4.1 Arreglos: los que casi no usarás

Un **arreglo** es una lista de tamaño **fijo**. El tamaño es parte del tipo.

**Fig. 4.1** | Arreglos, para que los reconozcas.

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

Fíjate en dos cosas:

**Línea 8.** El arreglo nace con los valores cero, como todo en Go. No hay basura.

**Líneas 19-22.** Al asignar un arreglo a otra variable, **se copia completo**: modificar la copia no toca
al original. Recuerda esto, porque con los slices **no** pasa lo mismo, y ahí está la trampa de esta
lección.

> [!NOTE]
> 🔧 **Observación de ingeniería de software 4.1**
> `[3]int` y `[4]int` son **tipos distintos**. Una función que recibe `[3]int` no acepta un `[4]int`. Eso
> hace a los arreglos poco prácticos, y es la razón por la que en Go casi nadie los usa directamente: se
> usan slices, que sí crecen. Los arreglos están ahí porque los slices se construyen sobre ellos.

### 4.2 Slices: los que usarás siempre

Un **slice** es una lista de tamaño **variable**. Es lo que en otros lenguajes llamarías lista o arreglo
dinámico.

**Fig. 4.2** | Crear y hacer crecer un slice.

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

**Línea 15.** Un slice sin inicializar vale **`nil`** y su largo es `0`. Pero fíjate en algo importante:

> [!TIP]
> ✅ **Buena práctica 4.1**
> **Un slice `nil` se puede usar con `append`, con `len` y con `for range` sin ningún problema.** Por eso
> no hace falta inicializarlo con `[]string{}`: `var nombres []string` y directo a `append`. Es la forma
> idiomática, y otra muestra de que el valor cero de Go está pensado para ser útil.

**Línea 17: `nombres = append(nombres, "catalogo")`.** Fíjate que hay que **reasignar**. `append` no
modifica el slice: **devuelve uno nuevo**, y si no guardas el resultado, lo pierdes.

> [!WARNING]
> ⚠️ **Error común de programación 4.1**
> Escribir `append(nombres, "x")` sin asignar el resultado. **El compilador lo detecta** porque el valor
> devuelto no se usa (`append(...) evaluated but not used`), así que aquí Go te salva. Pero si lo asignas a
> otra variable por error, se compila y el slice original no cambia.

**Línea 31: `servicios[len(servicios)-1]`.** En Go no hay índice negativo como en Python: para el último
elemento se calcula. Y si te pasas del final, el programa truena:

**Fig. 4.3** | Salirse del rango.

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
> ⚠️ **Error común de programación 4.2**
> El `index out of range` es el segundo panic más frecuente de Go, después del puntero nulo. **Fíjate en lo
> útil que es el mensaje:** te dice el índice que pediste (`[5]`), el largo real (`length 2`) y la línea.
> Antes de indexar algo que venga de fuera —un archivo, una petición, un argumento— comprueba `len`.

### 4.3 La trampa de la memoria compartida

**Ésta es la sección más importante de la lección.** Es un comportamiento que sorprende a todo el mundo,
que no produce ningún error, y que puede corromper datos en silencio.

#### 4.3.1 El problema

**Fig. 4.4** | Copiar un slice **no** copia los datos.

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

**Modificaste `copia` y cambió `original`.** No hubo error, no hubo aviso, y los dos slices muestran lo
mismo.

Compáralo con la Fig. 4.1, donde el mismo código con un **arreglo** sí copiaba. La diferencia es de fondo.

#### 4.3.2 Por qué pasa

Un slice **no contiene** los datos: es una ventana que apunta a ellos. Por dentro guarda tres cosas:

| | |
|---|---|
| un **puntero** | a dónde están los datos de verdad |
| el **largo** (`len`) | cuántos elementos tiene ahora |
| la **capacidad** (`cap`) | cuántos caben antes de tener que mudarse |

Cuando escribes `copia := original`, Go copia **esas tres cosas** — no los datos. Las dos variables quedan
apuntando al **mismo** lugar.

Es como darte la dirección de una casa en vez de construirte una igual: si pintas «tu» casa, la mía
cambia, porque es la misma.

#### 4.3.3 La solución

**Fig. 4.5** | Copiar de verdad.

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

**Ahora sí son independientes.**

**Línea 11: `make([]int, len(original))`.** `make` crea un slice con espacio reservado. Es la forma de
decir «quiero un slice de este largo, con sus valores cero».

**Línea 16: `original...`.** Esos tres puntos significan «esparce los elementos uno por uno». Sin ellos
estarías intentando agregar el slice completo como un solo elemento, y no compila.

#### 4.3.4 Y la parte de verdad traicionera: `append`

**Fig. 4.6** | El mismo código se comporta distinto según la capacidad.

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

**Lee esa salida dos veces.** Es **el mismo código** —`append` y luego modificar— y el resultado es
distinto según cuánta capacidad hubiera de sobra.

> [!WARNING]
> 🔴 **Error común de programación 4.3 — el peor de esta lección**
> Guardar un slice que recibiste como parámetro, o quedarte con el resultado de `append` asumiendo que es
> independiente. **El comportamiento depende de la capacidad**, que casi nunca controlas y que cambia según
> cuánto haya crecido el slice antes. O sea: **tu programa puede funcionar en las pruebas y corromper datos
> en producción**, con el mismo código.
>
> **La regla que te salva:** si vas a **guardar** un slice que recibiste de alguien más, cópialo primero.
> Si solo lo vas a leer, no hace falta.

> [!NOTE]
> 🚀 **Tip de rendimiento 4.1**
> Cuando sabes cuántos elementos vas a agregar, dale la capacidad desde el inicio:
> `make([]Estado, 0, len(servicios))`. Así `append` no tiene que mudarse ni copiar nada mientras crece.
> Con listas chicas es irrelevante; con miles de elementos, se nota. **Y mídelo antes de creerlo.**

### 4.4 Maps

Un **map** guarda pares de llave y valor. Es lo que en otros lenguajes se llama diccionario o tabla
asociativa.

**Fig. 4.7** | Crear, escribir y leer un map.

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

**Líneas 27-28: aquí está lo que hay que entender.** Leer una llave que no existe **no da error**:
devuelve el valor cero del tipo. Para un `int` eso es `0`, que podría ser un valor legítimo.

> [!WARNING]
> ⚠️ **Error común de programación 4.4**
> Usar `m[llave]` sin comprobar si existe, cuando el valor cero es indistinguible de un dato real. Si
> guardas `map[string]int` con tiempos de respuesta y lees una llave inexistente, obtienes `0` — que
> parece «respondió instantáneo» en vez de «no lo medimos». **Usa siempre la forma de dos valores
> (`v, hay := m[k]`) cuando la ausencia signifique algo distinto de cero.**

> [!TIP]
> ✅ **Buena práctica 4.2**
> A la segunda variable llámala `hay`, `existe` u `ok`. En código de Go verás mucho `ok`, y es la
> convención: `if v, ok := m[k]; ok { … }`.

#### 4.4.1 El map nil: el único que muerde

**Fig. 4.8** | Un map `nil` se puede leer pero no escribir.

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
> 🔴 **Error común de programación 4.5**
> **Ésta es la gran asimetría de Go y hay que memorizarla:** un **slice** `nil` se puede usar con `append`
> sin problema; un **map** `nil` truena al escribirle. Los maps **hay que crearlos** con `make(map[K]V)` o
> `map[K]V{}`. Y fíjate en lo tramposo: leer del map nil **sí funciona**, así que el programa puede
> avanzar un rato antes de reventar.

### 4.5 El orden aleatorio, y por qué es deliberado

**Fig. 4.9** | El mismo programa, dos ejecuciones, dos órdenes.

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

**Tres ejecuciones, tres órdenes distintos.** Y no es un defecto: **Go lo hace a propósito**, aleatorizando
el punto de inicio en cada recorrido.

> [!NOTE]
> 🔧 **Observación de ingeniería de software 4.2**
> ¿Por qué molestar al programador con esto? Porque **el orden de un map nunca estuvo garantizado**, ni en
> Go ni en la mayoría de los lenguajes. Si Go lo dejara «casi siempre igual», habría programas que
> funcionan durante años y un día, al crecer los datos, cambian de orden y se rompen — y nadie entendería
> por qué. **Aleatorizándolo, te obliga a enterarte hoy.** Es la misma filosofía que negarse a compilar con
> una variable sin usar: un error temprano y molesto vale más que uno tardío e incomprensible.

#### 4.5.1 Ordenar para que la salida sea estable

Si necesitas orden, se pide explícitamente.

**Fig. 4.10** | Recorrer un map en orden alfabético.

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

**Y ahora sí: la misma salida, siempre.** Corre el programa diez veces y no cambia.

**Línea 24: `make([]string, 0, len(estados))`.** Largo `0`, capacidad `len(estados)`. Sabes cuántas llaves
vas a agregar, así que se reserva el espacio de una vez.

**Línea 25: `for nombre := range estados`.** Sobre un map, `range` da **la llave** en la primera variable.
Sobre un slice da **el índice**. Es una diferencia que conviene tener presente:

| Colección | primera variable | segunda variable |
|---|---|---|
| slice | el **índice** (0, 1, 2…) | el elemento |
| map | la **llave** | el valor |

**Línea 35: `for _, nombre := range nombres`.** Aquí `nombres` es un slice, así que la primera variable es
el índice, y no lo necesito: por eso `_`.

> [!TIP]
> ✅ **Buena práctica 4.3**
> **Si tu programa imprime resultados, ordénalos.** Una salida que cambia de orden entre ejecuciones es
> imposible de comparar, imposible de probar automáticamente y confunde a quien la lee. Este patrón de
> tres pasos —sacar llaves, ordenar, recorrer— es idiomático y lo vas a usar seguido.

### 4.6 Leer un archivo de verdad

Hasta ahora los datos han estado escritos en el programa. Eso no sirve: cada cambio exige recompilar.

**Fig. 4.11** | Leer servicios desde un archivo de texto.

Primero el archivo de datos, `servicios.txt`:

```
catalogo https://catalogo.example.com 443
pagos https://pagos.example.com 443
inventario https://inventario.example.com 8080
```

Y el programa:

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

Y probemos qué pasa cuando algo sale mal:

```bash
$ go run fig04_11.go            # con el archivo renombrado
error: leyendo servicios.txt: open servicios.txt: no such file or directory
$ echo $?
1
```

**Línea 20: `os.ReadFile`.** Lee el archivo completo y devuelve sus bytes. Para archivos de configuración
es perfecto; para un archivo de varios gigabytes existen otras formas, porque esto lo carga todo en
memoria.

**Línea 26: `string(datos)`.** Convierte los bytes a texto. `strings.TrimSpace` quita espacios y saltos de
línea del inicio y el final — sin eso, la última línea vacía produciría un servicio fantasma.

**Línea 34: `strings.Fields`.** Separa por espacios, **colapsando los repetidos**. Es más robusto que
`strings.Split(linea, " ")`, que con dos espacios seguidos te daría un campo vacío.

**Línea 40: `strconv.Atoi`.** Convierte texto a entero (*ASCII to integer*). Devuelve error si no es un
número, y **lo estamos revisando**.

**Líneas 58-59: la diferencia entre un script y un programa.**

> [!TIP]
> ✅ **Buena práctica 4.4**
> Los errores van a **`os.Stderr`**, no a la salida normal, y el programa termina con **código distinto de
> cero**. Eso permite usarlo en una tubería: `mi-programa | otro-programa` sigue funcionando porque el
> error no contamina la salida, y `mi-programa && echo ok` no imprime «ok» si falló. **Es lo que separa un
> programa de un script.**

> [!TIP]
> 🧪 **Tip de prueba y depuración 4.1**
> Fíjate en los mensajes de error de las líneas 36 y 42: dicen **el archivo, el número de línea y qué
> esperaban**. Compáralo con un `invalid syntax` a secas. Cuando alguien use tu programa con un archivo de
> cincuenta líneas, la diferencia entre los dos mensajes son veinte minutos de su vida.

---

## El error que vas a ver

Esta lección deja dos mensajes de `panic` grabados en la memoria, porque son los dos más frecuentes de
todo el curso después del puntero nulo:

**`panic: runtime error: index out of range [N] with length M`** (Fig. 4.3). Pediste una posición `N` que
no existe en una colección de largo `M`. **Qué significa:** el índice válido más alto es `M-1`, no `M`.
**Cómo se arregla:** comprueba `len(coleccion)` antes de indexar, sobre todo si el índice viene de fuera
(un archivo, un argumento, una petición) y no lo escribiste tú mismo en el código.

**`panic: assignment to entry in nil map`** (Fig. 4.8). Intentaste escribir en un map que nunca se creó —
`var m map[string]Estado` deja `m` en `nil`, y un map `nil` **se puede leer pero no escribir**. **Qué
significa:** falta el `make(map[K]V)` o el `map[K]V{}` que reserva la tabla por dentro. **Cómo se
arregla:** crea el map antes de usarlo como destino de una asignación; si solo lo vas a leer, `nil` no da
ningún problema.

Y hay un tercer caso que **no** truena, y por eso es el más peligroso de los tres: modificar un slice
compartido (Fig. 4.4 y Fig. 4.6) no produce ningún mensaje de error. El programa sigue corriendo con datos
incorrectos. Ese es exactamente el motivo de que la sección 4.3 exista.

## Lo que se hace mal

- **Asumir que `copia := original` copia los datos de un slice.** Copia el puntero, el largo y la
  capacidad — no los datos. Las dos variables terminan compartiendo memoria (Fig. 4.4).
- **Guardar un slice que recibiste como parámetro sin copiarlo antes.** El comportamiento de `append`
  depende de la capacidad que traiga el slice, que casi nunca controla quien lo recibe (Fig. 4.6). Si vas
  a **guardar** un slice ajeno, cópialo; si solo lo vas a leer, no hace falta.
- **Confiar en que un `for range` sobre un map va a salir siempre en el mismo orden.** Nunca estuvo
  garantizado, y Go lo aleatoriza a propósito (Fig. 4.9) para que el error salga hoy y no el día que crezca
  la base de datos en producción.
- **Leer `m[llave]` sin comprobar `ok` cuando el valor cero es ambiguo.** Un `map[string]int` que guarda
  tiempos de respuesta no puede distinguir «respondió en 0 ms» de «nunca se consultó» si no usas la forma
  de dos valores.
- **Indexar algo que viene de fuera —un archivo, `os.Args`, una respuesta de red— sin comprobar `len`
  primero.** Es la causa más común del `index out of range` de esta lección, y se evita con una línea.
- **Escribir en un map antes de crearlo con `make` o con `{}`.** El compilador no lo detecta porque `var m
  map[K]V` es código válido; el `panic` aparece hasta que el programa corre e intenta escribir.

---

## Ejercicios

### Preguntas

**4.1** ¿Cuál es la diferencia entre `[3]int` y `[]int`?

**4.2** ¿Por qué `[3]int` y `[4]int` son tipos distintos, y qué consecuencia práctica tiene?

**4.3** ¿Qué imprime esto y por qué?

```go
a := []int{1, 2, 3}
b := a
b[0] = 99
fmt.Println(a[0])
```

**4.4** Nombra las dos formas de copiar un slice de verdad.

**4.5** ¿Por qué el mismo `append` a veces comparte memoria con el original y a veces no?

**4.6** ¿Qué devuelve `m["no-existe"]` en un `map[string]int`, y por qué es peligroso?

**4.7** ¿Cuál es la asimetría entre un slice `nil` y un map `nil`?

**4.8** ¿Por qué el recorrido de un map sale desordenado, y cómo obtienes un orden fijo?

**4.9** En `for a, b := range x`, ¿qué es `a` si `x` es un slice? ¿Y si es un map?

**4.10** ¿Qué hace `strings.Fields` que `strings.Split(s, " ")` no hace?

**4.11** ¿Por qué los errores van a `os.Stderr` y no a la salida normal?

**4.12** Este programa tiene **dos** problemas. Encuéntralos.

```go
func main() {
    var m map[string]int
    m["catalogo"] = 443
    lista := []int{1, 2, 3}
    fmt.Println(lista[3])
}
```

### Ejercicios de código

**4.13** Escribe `func Contar(estados map[string]Estado) (ok, fallas int)` que cuente cuántos estados
tienen código 200-299 y cuántos no. Go devuelve varios valores: úsalo.

**4.14** Escribe `func Nombres(servicios []Servicio) []string` que devuelva solo los nombres, **ordenados**.

**4.15 (La trampa)** Escribe una función `func Guardar(s []int)` que guarde el slice en una variable global
y luego lo imprima. Llámala así:

```go
datos := make([]int, 3, 10)
datos[0], datos[1], datos[2] = 1, 2, 3
Guardar(datos)
datos = append(datos, 4)
datos[0] = 999
// ahora imprime lo que guardó Guardar
```

**Predice qué va a imprimir antes de correrlo.** Después córrelo. Si acertaste, entendiste la sección 4.3;
si no, vuelve a leerla — es la que más caro sale.

**4.16** Extiende `Cargar` de la Fig. 4.11 para que acepte un cuarto campo opcional con el timeout en
milisegundos. Si no viene, usa 5000. **Cuida que un archivo con tres campos siga funcionando.**

**4.17** Escribe `func Agrupar(servicios []Servicio) map[int][]Servicio` que agrupe por puerto. Prueba que
un puerto con tres servicios los tenga todos.

**4.18 (Encuentra el error)** Di qué está mal en cada uno:

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

**4.19 (Proyecto del curso)** Lleva tu programa al siguiente paso:
1. `Cargar(ruta string) ([]Servicio, error)` que lea el archivo, como la Fig. 4.11.
2. `RevisarTodos` que devuelva `map[string]Estado`.
3. `Reporte(estados map[string]Estado) string` que produzca la tabla **ordenada alfabéticamente**.
4. Un conteo final: «3 de 5 respondieron».
5. Que el programa acepte la ruta del archivo como argumento: `os.Args[1]`. 🔴 **Comprueba `len(os.Args)`
   antes de leerlo**, o tendrás un `index out of range` cuando alguien lo corra sin argumentos.
6. Que salga con código 0 si todo respondió y 1 si algo falló.

### Soluciones

**4.1** `[3]int` es un **arreglo**: tamaño fijo, parte del tipo, y se copia al asignarse. `[]int` es un
**slice**: tamaño variable, y al asignarse **comparte los datos**.

**4.2** Porque el tamaño es parte del tipo. La consecuencia: una función que recibe `[3]int` no acepta un
`[4]int`, lo que vuelve a los arreglos poco prácticos. Por eso se usan slices.

**4.3** Imprime **99**. `b := a` no copia los datos: las dos variables apuntan a la misma memoria.

**4.4** `make` + `copy`, o `append([]T(nil), original...)`.

**4.5** Porque depende de la **capacidad**. Si el slice tiene espacio de sobra, `append` escribe ahí mismo y
sigue compartiendo memoria; si no cabe, se muda a memoria nueva y deja de compartir. Como casi nunca
controlas la capacidad, **el mismo código puede comportarse distinto**.

**4.6** Devuelve `0`, el valor cero, **sin error**. Es peligroso porque `0` puede ser un dato legítimo, así
que no puedes distinguir «vale cero» de «no está». Se resuelve con `v, ok := m[k]`.

**4.7** Un slice `nil` acepta `append` sin problema. Un map `nil` **truena al escribirle**
(`assignment to entry in nil map`), aunque sí se puede leer. Los maps hay que crearlos con `make`.

**4.8** Porque su orden **nunca estuvo garantizado**, y Go lo aleatoriza a propósito para que no escribas
programas que dependan de un orden accidental. Para orden fijo: sacar las llaves a un slice, ordenarlas con
`sort.Strings`, y recorrer el slice.

**4.9** Si `x` es un slice, `a` es el **índice**. Si es un map, `a` es la **llave**.

**4.10** `Fields` separa por espacios **colapsando los repetidos**, así que dos espacios seguidos no
producen un campo vacío. También trata los tabuladores como separador.

**4.11** Para que la salida normal quede limpia y se pueda usar en una tubería, y porque quien llama al
programa espera encontrar los errores ahí. Junto con el código de salida distinto de cero, es lo que hace
que un programa se pueda automatizar.

**4.12** Los dos: (1) `m` es un map `nil` y escribirle provoca `panic: assignment to entry in nil map` —
falta `m := make(map[string]int)`; (2) `lista[3]` se sale del rango, porque los índices válidos son 0, 1 y
2 — `panic: index out of range [3] with length 3`.

**Ejercicios de código (4.13 a 4.19):** sus soluciones verificadas —compiladas y ejecutadas— se agregan en
la entrega en la que se cierre el capítulo de ejercicios con solución del curso completo; no se publican
sin haber corrido `go build` sobre cada una.

---

## Cómo sé que lo logré

- Puedes explicar, sin ver el texto, por qué `copia := original` no copia los datos de un slice — y
  dibujarlo (puntero, largo, capacidad).
- Prediciste correctamente qué imprime el ejercicio 4.15 **antes** de correrlo. Si no acertaste, repetiste
  la sección 4.3 hasta que sí.
- Tu `Cargar` (ejercicio 4.16) sigue funcionando con archivos de tres campos y acepta el cuarto opcional.
- Tu `Reporte` produce **la misma salida, en el mismo orden**, en diez corridas seguidas.
- Tu programa no truena con `index out of range` si lo corres sin argumentos: comprobaste `len(os.Args)`
  antes de leer `os.Args[1]`.
- Anotaste en [`bitacora.md`](bitacora.md) el resultado del ejercicio 4.15: qué predijiste y qué pasó
  realmente. Es lo que separa entender los slices de creer que los entiendes.

---

## Resumen

- Un **arreglo** (`[3]int`) tiene tamaño fijo, el tamaño es parte del tipo, y **se copia** al asignarse.
  Casi no se usa directamente.
- Un **slice** (`[]int`) tiene tamaño variable y es lo que se usa siempre.
- Un slice guarda **puntero, largo y capacidad**: no contiene los datos, apunta a ellos.
- 🔴 **Asignar un slice NO copia los datos**: las dos variables comparten memoria. Para copiar de verdad,
  `make` + `copy` o `append([]T(nil), s...)`.
- 🔴 **`append` comparte o no memoria según la capacidad disponible**, así que el mismo código puede
  comportarse distinto. Si vas a **guardar** un slice ajeno, cópialo.
- Un slice `nil` funciona con `append`, `len` y `range`. **Un map `nil` truena al escribirle.**
- Leer una llave inexistente de un map **devuelve el valor cero sin error**. Usa `v, ok := m[k]` cuando la
  ausencia importe.
- **El recorrido de un map es aleatorio a propósito.** Para orden fijo: llaves a un slice, `sort`, y
  recorrer el slice.
- En `range`, la primera variable es el **índice** en slices y la **llave** en maps. `_` descarta.
- `os.ReadFile` lee un archivo completo; `strings.Fields` separa por espacios colapsando repetidos;
  `strconv.Atoi` convierte texto a entero y **devuelve error**.
- Los errores van a **`os.Stderr`** y el programa sale con **código distinto de cero**.

---

## Para leer más

1. **[Go Slices: usage and internals](https://go.dev/blog/slices-intro)** — el blog oficial del equipo de
   Go, explica el puntero/largo/capacidad con dibujos. Empieza aquí.
2. **[Go maps in action](https://go.dev/blog/maps)** — blog oficial, cubre el orden aleatorio y el map
   `nil` con más detalle del que cabe en esta lección.
3. **[Go by Example: Slices](https://gobyexample.com/slices) y [Maps](https://gobyexample.com/maps)** —
   código mínimo, bueno como referencia rápida.
4. **[Effective Go — Slices](https://go.dev/doc/effective_go#slices)** — la sección específica sobre el
   patrón de dos pasos (`make` + `copy`) y por qué `append` se comporta como se comporta.

### Términos de esta lección

| | |
|---|---|
| **`append`** | agrega elementos a un slice y **devuelve** el resultado |
| **arreglo** | lista de tamaño fijo; el tamaño es parte del tipo |
| **capacidad (`cap`)** | cuántos elementos caben en un slice antes de mudarse |
| **`copy`** | copia elementos entre slices |
| **`delete`** | quita una llave de un map |
| **`index out of range`** | panic por pedir una posición que no existe |
| **largo (`len`)** | cuántos elementos tiene ahora |
| **map** | colección de pares llave→valor |
| **`make`** | crea slices, maps y canales con espacio reservado |
| **`nil map`** | map sin crear; se puede leer pero **no** escribir |
| **`os.Stderr`** | salida de errores, separada de la normal |
| **slice** | lista de tamaño variable; una ventana sobre los datos |
| **`sort.Strings`** | ordena un slice de texto |
| **`strings.Fields`** | separa un texto por espacios, colapsando repetidos |
| **`strconv.Atoi`** | convierte texto a entero |
| **valor de dos resultados (`v, ok`)** | forma de leer un map distinguiendo la ausencia |

---

**Anterior:** [Lección 3 — Structs, errores e interfaces](03-errores-interfaces.md) ·
**Siguiente:** [Lección 5 — Módulos y pruebas](05-modulos-y-pruebas.md)
