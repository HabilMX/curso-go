# Lección 2 — Variables, funciones y tipos

**Duración:** 90 minutos, o dos sesiones de 45.

**Al terminar vas a poder:**

- Explicar qué hace el compilador de Go con tu archivo de texto.
- Declarar variables y saber cuándo usar `:=` y cuándo `=`.
- Nombrar los cuatro tipos básicos y decir por qué Go no te deja mezclarlos.
- Escribir funciones que reciben y devuelven valores.
- Entender por qué una función de Go puede devolver **dos** cosas a la vez.
- Leer los mensajes de error del compilador y saber qué te está pidiendo.
- Escribir, compilar y ejecutar un programa completo por tu cuenta.

---

## Por qué importa

El `revisor` —el programa que vas a construir durante todo el curso— necesita, desde su primera línea,
tres cosas: **datos** (el nombre de un servicio, su puerto, cuánto tardó en responder), **lógica** que
decida algo con esos datos (¿está rápido o lento?) y una forma de **avisar cuando algo salió mal** (¿el
servicio no contestó?). Eso es, en ese orden, exactamente lo que esta lección enseña: variables para
guardar datos, funciones para la lógica, y el patrón de dos valores de retorno para los errores.

No es casualidad que Go sea **compilado** y de **tipos estáticos**: las dos decisiones existen para que la
computadora atrape tus errores **antes** de que el programa corra, en vez de que los descubra un usuario en
producción. Un lenguaje interpretado te deja escribir `puerto = "muchos"` y solo truena cuando esa línea se
ejecuta —que puede ser diez minutos después, o nunca, si esa ruta del código casi no se usa—. Go se niega a
producir el ejecutable. Esa diferencia es la razón de fondo de casi todo lo que vas a ver en esta lección:
por qué el compilador es estricto, por qué existe el «valor cero», y por qué una función que puede fallar
**tiene que decirlo en su firma**, no como una sorpresa.

Al final de la lección vas a escribir el primer paso real del `revisor` (ejercicio de proyecto): una función
que recibe el nombre de un servicio, su puerto y su tiempo de respuesta, y devuelve una línea de reporte
clasificada. Todavía sin structs, sin concurrencia y sin HTTP —eso viene después— pero ya es código que se
parece al programa final.

---

## Los conceptos

### 2.1 Qué hace el compilador

Escribes un archivo de texto. Ese archivo, por sí solo, no hace nada: es texto. Para que la computadora
lo ejecute, algo tiene que traducirlo a las instrucciones que el procesador entiende.

Hay dos maneras de hacer esa traducción:

| | Cómo funciona | Lenguajes |
|---|---|---|
| **Interpretado** | Un programa lee tu texto línea por línea y hace lo que dice, **cada vez** que lo ejecutas | Python, JavaScript, PHP |
| **Compilado** | Un programa traduce **todo** tu texto **una sola vez** y guarda el resultado en un archivo ejecutable | **Go**, C, C++, Rust |

**Go es compilado**, y eso tiene tres consecuencias que vas a notar desde hoy:

1. **Hay un paso previo a ejecutar: compilar.** Si escribiste mal el nombre de una variable, te enteras
   ahí — antes de que el programa corra. En un lenguaje interpretado te enterarías cuando la ejecución
   llegara a esa línea, que puede ser diez minutos después o nunca.
2. **El resultado es un archivo que funciona solo.** No necesita que la otra computadora tenga Go
   instalado.
3. **Es rápido al ejecutar**, porque la traducción ya está hecha.

> [!NOTE]
> 🔧 **Observación de ingeniería de software 2.1**
> El compilador es tu primera línea de defensa, no un obstáculo. Cada error que atrapa es un error que no
> vas a estar buscando de noche con el programa ya entregado. Cuando Go te detenga, lee el mensaje con
> calma: te está ahorrando trabajo.

### 2.2 Tu primer programa, línea por línea

Vamos a escribir, compilar y ejecutar un programa completo. **Todos los programas de este curso son
completos y ejecutables**: nada de fragmentos que no se pueden correr.

Prepara la carpeta:

```bash
mkdir -p ~/w/curso-go/cap01 && cd ~/w/curso-go/cap01
go mod init cap01
```

**Fig. 2.1** | Un programa que imprime un mensaje.

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

Compílalo y ejecútalo:

```bash
$ go run fig02_01.go
Bienvenido a Go!
```

Ahora **cada línea**, porque todas tienen una razón:

**Líneas 1-2.** Comentarios. Empiezan con `//` y el compilador los ignora por completo. Están para quien
lea el código — incluido tú en tres semanas.

**Línea 3: `package main`.** Todo archivo de Go pertenece a un *paquete*, que es una agrupación de código.
**El paquete `main` es especial: es el único que produce un programa ejecutable.** Si escribieras
`package utilidades`, tendrías una biblioteca que otros programas pueden usar, pero que no se puede
ejecutar por sí sola.

**Línea 5: `import "fmt"`.** Le dice al compilador que vas a usar código del paquete `fmt` (de *format*),
que viene con Go y contiene las funciones para imprimir y dar formato a texto. **Go no carga nada por
omisión**: todo lo que uses, lo pides.

**Línea 7: `func main() {`.** Declara la función `main`, que es **donde comienza la ejecución del
programa**. Cuando ejecutas un programa de Go, el sistema busca exactamente esta función. Si la llamaras
`principal` o `inicio`, el programa no arrancaría.

**Línea 8: `fmt.Println("Bienvenido a Go!")`.** Llama a la función `Println` que está dentro del paquete
`fmt`. El punto se lee como «de»: *«la función `Println` **de** `fmt`»*. `Println` imprime lo que le des y
salta a la línea siguiente (*print line*).

**Línea 9: `}`.** Cierra el cuerpo de la función. En Go las llaves delimitan bloques, igual que en C,
C++, Java o JavaScript.

> [!TIP]
> ✅ **Buena práctica 2.1**
> Pon siempre un comentario al inicio del archivo con su nombre y qué hace. Toma cinco segundos y ahorra
> minutos a quien lo abra después.

#### 2.2.1 La diferencia entre `go run` y `go build`

**Fig. 2.2** | Las dos formas de ejecutar tu programa.

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

Fíjate en los tamaños: tu texto son **104 bytes**; el programa compilado, **1.8 megabytes**. La diferencia
es que el ejecutable trae adentro todo lo que necesita para funcionar.

> [!NOTE]
> 🚀 **Tip de portabilidad 2.1**
> Ese archivo `fig02_01` se puede copiar a **cualquier** computadora con Linux de la misma arquitectura y
> funciona, aunque esa máquina no tenga Go instalado. Es la razón por la que tantas herramientas de
> servidores están escritas en Go.

> [!TIP]
> 🧪 **Tip de prueba y depuración 2.1**
> Usa `go run` mientras programas: es más rápido y no te llena la carpeta de ejecutables. Usa `go build`
> cuando el programa ya sirva y quieras entregarlo.

### 2.3 Variables

Una **variable** es un espacio de memoria con un nombre, donde guardas un dato que tu programa va a usar.

**Fig. 2.3** | Declarar variables e imprimirlas.

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

**El operador `:=`** se lee *«declara una variable nueva y guárdale esto»*. Go mira el valor de la derecha
y **deduce el tipo solo**: `"catalogo"` está entre comillas, así que es texto; `443` no las tiene y no tiene
punto, así que es entero.

#### 2.3.1 `:=` contra `=`

Una vez que la variable existe, para cambiarle el valor usas `=` **sin** los dos puntos:

**Fig. 2.4** | Crear y modificar una variable.

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
> ✅ **Buena práctica 2.2**
> Si de verdad necesitas declarar algo que no vas a usar todavía, usa el identificador vacío `_`. En Go,
> `_` significa *«esto lo descarto a propósito»*, y el compilador lo acepta.

#### 2.3.2 Go se niega a compilar si no usas una variable

**Fig. 2.5** | Un programa que **no compila**, a propósito.

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

**Go no lo permite.** No es una advertencia que puedas ignorar: es un error que detiene la compilación.

**¿Por qué tan estricto?** Porque una variable que no se usa casi siempre significa una de dos cosas: te
equivocaste al escribir el nombre en otra línea, o quedó basura de código que borraste a medias. Las dos
son problemas. Go prefiere molestarte hoy que dejar código muerto acumulándose para siempre.

### 2.4 Los tipos básicos

Un **tipo** es la respuesta a la pregunta *«¿qué clase de dato es esto?»*. Los cuatro que vas a usar todo
el tiempo:

| Tipo | Qué guarda | Ejemplos | Valor cero |
|---|---|---|---|
| `string` | texto | `"hola"`, `"https://catalogo.example.com"` | `""` (vacío) |
| `int` | números enteros | `42`, `-7`, `0` | `0` |
| `float64` | números con decimales | `3.14`, `0.142` | `0` |
| `bool` | verdadero o falso | `true`, `false` | `false` |

#### 2.4.1 Go no mezcla tipos

**Fig. 2.6** | Otro programa que **no compila**, a propósito.

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

En Python esto funcionaría sin queja. **En Go no compila**, y eso es una ventaja enorme:

> [!NOTE]
> 🔧 **Observación de ingeniería de software 2.2**
> Buena parte de los errores que llegan a producción son de esta clase: alguien creyó que una variable
> tenía un número y tenía texto. En un lenguaje interpretado eso truena **cuando un cliente está usando el
> programa**. En Go no llega a compilar: el error aparece en tu máquina, no en la del cliente.

#### 2.4.2 El valor cero: en Go nunca hay basura

**Fig. 2.7** | Variables declaradas sin valor.

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

La forma `var nombre tipo` declara una variable **sin darle valor**. En otros lenguajes eso dejaría
«basura» —lo que hubiera en esa memoria— o un estado especial de «indefinido». **En Go siempre recibe un
valor válido**, llamado **valor cero**.

**Y aquí aparece `Printf`**, que es distinto de `Println`: recibe una plantilla con huecos y los va
llenando. Cada hueco empieza con `%`:

| | Para qué |
|---|---|
| `%d` | un número entero (*decimal*) |
| `%s` | texto (*string*) |
| `%q` | texto **con comillas** (*quoted*) — útil para ver si algo está vacío |
| `%t` | un `bool` (*true/false*) |
| `%v` | cualquier cosa, en su formato natural (*value*) |
| `\n` | salto de línea (`Printf` **no** salta solo, `Println` sí) |

> [!TIP]
> ✅ **Buena práctica 2.3**
> Usa `%q` cuando imprimas texto que estés depurando. `%s` con un texto vacío no imprime nada y parece
> que la línea falló; `%q` muestra `""` y se ve claramente que el texto está vacío.

### 2.5 Funciones

Una **función** es un bloque de código con nombre, al que puedes darle datos y que puede devolverte un
resultado. Ya usaste dos: `fmt.Println` y `fmt.Printf`.

**Fig. 2.8** | Definir y llamar funciones.

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

**Anatomía de la línea 8**, que es donde está todo:

```
func  sumar  (a int, b int)  int  {
 │      │          │          │
 │      │          │          └── lo que DEVUELVE
 │      │          └── lo que RECIBE (los parámetros), con su tipo
 │      └── el nombre
 └── palabra clave: «voy a definir una función»
```

**Línea 14: `fmt.Sprintf`.** Es como `Printf` pero en vez de imprimir, **devuelve** el texto armado. La
`S` es de *string*. Se usa muchísimo.

**Línea 19: `puerto == 443`.** El doble igual **compara** y da `true` o `false`. Un solo igual `=`
**asigna**. Confundirlos es clásico (ver «Lo que se hace mal», más abajo).

> [!TIP]
> ✅ **Buena práctica 2.4**
> Cuando dos parámetros son del mismo tipo puedes abreviar: `func sumar(a, b int) int`. Es lo idiomático y
> más corto. Aquí lo escribimos completo para que se vea la estructura.

### 2.6 Devolver dos valores: la firma de Go

**Esto es lo más distintivo de Go**, y lo vas a escribir miles de veces en tu carrera. Presta atención.

Piensa en una función que divide dos números. ¿Qué hace si el divisor es cero? No puede devolver un
número, porque el resultado no existe.

Otros lenguajes usan **excepciones**: la función «lanza» un error y alguien más arriba lo «atrapa» con
`try / catch`. El problema es que **no es obvio quién lo atrapa ni dónde**, y es facilísimo que nadie lo
haga y el programa muera.

**Go hace otra cosa: devuelve dos valores.** El resultado, y un error.

**Fig. 2.9** | Una función que devuelve un resultado y un posible error.

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

**Línea 11: `(int, error)`.** Los paréntesis con dos tipos significan *«esta función devuelve dos
cosas»*: un entero y un error.

**Línea 13: `return 0, errors.New(...)`.** Cuando hay problema, devuelve un valor cualquiera (aquí `0`,
que no se va a usar) **y** un error describiendo qué pasó.

**Línea 15: `return a / b, nil`.** Cuando todo va bien, devuelve el resultado **y `nil`**.

🔑 **`nil` significa «nada», «vacío», «no hay».** Entonces la pregunta `if err != nil` se lee:
*«¿el error es distinto de nada?»*, o en español normal: **«¿hubo error?»**

**Líneas 20-21: el patrón que define a Go.**

```go
resultado, err := dividir(10, 2)
if err != nil {
    // algo salió mal: aquí se atiende
}
// si llegaste aquí, todo bien
```

Este patrón `if err != nil` es **la línea más escrita en toda la historia de Go**. Aparece por todas
partes, y hay gente que lo critica por repetitivo.

> [!NOTE]
> 🔧 **Observación de ingeniería de software 2.3**
> La crítica es cierta: es verboso. La ventaja es que **es imposible ignorar un error por descuido**,
> porque está ahí, en tu cara, en la línea siguiente a la llamada. Con excepciones puedes olvidar un
> `catch` y no enterarte hasta que el programa muera en producción. Go cambió brevedad por seguridad, a
> propósito.

> [!TIP]
> 🧪 **Tip de prueba y depuración 2.2**
> Fíjate en la línea 28: dice `resultado, err = dividir(10, 0)` con `=`, no `:=`. Es porque las dos
> variables **ya existen** de la línea 20. Si pusieras `:=` ahí, el compilador te lo diría.

### 2.7 Poniéndolo todo junto

**Fig. 2.10** | Programa completo que usa todo lo de la lección.

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

Este programa ya es el antepasado directo del `revisor`: `clasificar` y `reporte` son, en miniatura, lo
que en la lección 3 se convertirá en un método sobre un struct `Servicio`, y en la lección 6 correrá para
varios servicios **a la vez**.

**Línea 21: `%-10s` y `%-5d`.** El número dice **cuántos espacios de ancho** reservar, y el signo menos
significa **alineado a la izquierda**. Así las columnas quedan derechas. Sin eso, la tabla saldría torcida.

> [!TIP]
> ✅ **Buena práctica 2.5**
> Cuando una llamada es muy larga, córtala en varias líneas como en las líneas 21-22. Go lo permite
> siempre que la coma quede al final de la línea anterior.

> [!NOTE]
> 🚀 **Tip de rendimiento 2.1**
> `Sprintf` es cómodo pero no es gratis: arma texto nuevo en memoria cada vez. Para cinco líneas no
> importa en absoluto. Si algún día armas miles, existe `strings.Builder`, que es mucho más rápido. **No
> optimices antes de medir**: esto es información para después, no algo que hacer hoy.

---

## El error que vas a ver

Casi todos los mensajes de error de Go tienen la misma forma:

```
<archivo>:<línea>:<columna>: <qué esperaba el compilador y qué encontró>
```

Aprender a leer esa línea es la mitad del trabajo. Dos ejemplos que ya viste en esta lección:

```
./fig02_05.go:8:5: declared and not used: edad
```

Se lee: *«en el archivo `fig02_05.go`, línea 8, columna 5, declaraste `edad` y nunca la usaste»*. El
arreglo es sencillo: o usas la variable, o la borras, o —si de verdad la necesitas declarada pero todavía
no la usas— la cambias por el identificador vacío `_`.

```
./fig02_06.go:8:14: cannot use "muchos" (untyped string constant) as int value in assignment
```

Se lee: *«en la línea 8, columna 12, intentaste usar el texto `"muchos"` donde se esperaba un `int`, y no
te dejo»*. El arreglo tampoco es sencillo de adivinar si nunca lo viste: **no existe una conversión
automática** entre `string` y `int` en una asignación así. Tendrías que convertir explícitamente (con
`strconv`, que se ve más adelante) o, más probable, darte cuenta de que mezclaste el tipo por error.

Y un tercero, el más común de la primera semana, que aparece si reusas `:=` sobre una variable que ya
existe:

```
./fig02_04.go:10:9: no new variables on left side of :=
```

Se lee: *«del lado izquierdo del `:=` no hay ninguna variable nueva»* — todas ya existían, así que Go no
sabe qué declarar. El arreglo: usa `=` en vez de `:=` cuando ya declaraste la variable antes.

**La regla general:** el compilador de Go casi nunca dice «algo salió mal» en abstracto. Dice el archivo,
la línea, la columna y una frase que —aunque suene rara la primera vez— describe exactamente el problema.
Leerla completa, sin asustarse, resuelve la mayoría de los casos sin tener que buscar nada en internet.

## Lo que se hace mal

- **Escribir `println` en minúscula en vez de `Println`.** Go **distingue mayúsculas de minúsculas**:
  `Println` y `println` son nombres distintos (`println` es una función interna del compilador, pensada
  para depurar el propio Go, no tu programa). El error es `undefined: fmt.println`.
- **Usar `import "fmt"` y no usarlo, o usar `fmt.Println` sin el `import`.** El primero da
  `"fmt" imported and not used`; el segundo, `undefined: fmt`. Los dos son el mismo tipo de error: le
  dijiste al compilador algo que no coincide con lo que hiciste.
- **Confundir `=` (asignar) con `==` (comparar)** dentro de un `if`. En C esto compilaría y haría algo
  inesperado (asignar el valor y seguir de largo). **En Go es un error de compilación** —
  `cannot use puerto = 443 (...) as value` —, lo cual es una buena noticia: te detiene antes de que el
  programa haga algo que no pediste.
- **Olvidar el `\n` en `Printf`.** El programa **sí compila y sí corre**, pero todo sale pegado en una
  sola línea. `Println` salta de línea automáticamente; `Printf` no. Es un error silencioso, no uno que el
  compilador te señale.
- **Ignorar el error con `_`**, así: `resultado, _ := dividir(10, 0)`. **Esto también compila**, y el
  programa sigue adelante con `resultado` valiendo `0` como si todo estuviera bien. Es la forma más rápida
  de crear un error imposible de encontrar después, porque no hay ningún mensaje ni caída: el programa
  simplemente sigue con un dato incorrecto. **Si alguna vez descartas un error con `_`, que sea una
  decisión consciente y comentada, nunca un reflejo para que el compilador deje de quejarse.**

## Ejercicios

### Preguntas de repaso

Contesta sin ver las respuestas. Están al final de esta sección.

**2.1** ¿Cuál es la diferencia entre un lenguaje compilado y uno interpretado, y cuál de los dos es Go?

**2.2** ¿Por qué la función donde empieza el programa tiene que llamarse `main` exactamente?

**2.3** ¿Qué diferencia hay entre `:=` y `=`?

**2.4** ¿Qué error da Go si declaras una variable y no la usas? ¿Por qué crees que lo hace?

**2.5** ¿Qué es el *valor cero* y cuál es el de `string`, `int` y `bool`?

**2.6** ¿Qué imprime este programa?

```go
package main

import "fmt"

func main() {
    var n int
    var s string
    fmt.Printf("[%d] [%q]\n", n, s)
}
```

**2.7** ¿Qué diferencia hay entre `Println`, `Printf` y `Sprintf`?

**2.8** ¿Qué significa `nil` y qué se pregunta `if err != nil`?

**2.9** Este programa tiene **tres** errores que impiden que compile. Encuéntralos sin correrlo.

```go
package main

func main() {
    nombre := "catalogo"
    puerto := 443
    puerto := 8080
    fmt.Println(nombre)
}
```

#### Respuestas

**2.1** Un lenguaje **interpretado** se traduce línea por línea en cada ejecución; uno **compilado** se
traduce completo una sola vez y produce un ejecutable. **Go es compilado.**

**2.2** Porque es una convención del lenguaje: al ejecutar un programa, Go busca la función `main` del
paquete `main` para empezar. Con otro nombre no encuentra por dónde arrancar.

**2.3** `:=` **declara** una variable nueva y le asigna valor; `=` solo **asigna** a una que ya existe.

**2.4** `declared and not used`. Lo hace porque una variable sin usar casi siempre indica un error —un
nombre mal escrito o restos de código borrado— y Go prefiere detenerte a dejar código muerto.

**2.5** Es el valor que recibe automáticamente una variable declarada sin valor. `string` → `""`,
`int` → `0`, `bool` → `false`. Garantiza que **nunca** hay basura ni «indefinido».

**2.6** `[0] [""]` — el valor cero de `int` y de `string`, y `%q` muestra las comillas.

**2.7** `Println` imprime y salta de línea. `Printf` imprime con una plantilla de `%` y **no** salta
solo. `Sprintf` usa la misma plantilla pero **devuelve** el texto en vez de imprimirlo.

**2.8** `nil` significa «nada» / «vacío». `if err != nil` pregunta *«¿el error no está vacío?»*, es decir
**«¿hubo un error?»**.

**2.9** Los tres:
1. Falta `import "fmt"`, y se usa `fmt.Println`.
2. `puerto := 8080` usa `:=` sobre una variable que ya existe → debe ser `=`.
3. `puerto` se declara y **nunca se usa** (solo se imprime `nombre`) → `declared and not used`.

### Ejercicios de código

> 🔴 **Estos siete ejercicios todavía no traen solución de referencia** — se agrega en una entrega
> posterior. No se inventó código de solución para no publicar un ejemplo sin ejecutar primero.

**2.10** Escribe un programa que declare tu nombre, tu edad y si estudias o trabajas, y los imprima en
tres líneas con `Printf`, usando el verbo correcto para cada tipo.

**2.11** Escribe `func celsiusAFahrenheit(c float64) float64` y pruébala con 0, 37 y 100. La fórmula es
`f = c*9/5 + 32`. **Cuidado:** si escribes `c*9/5` con enteros el resultado se trunca. ¿Por qué aquí no
pasa?

**2.12** Escribe `func esPar(n int) bool`. Pista: el operador `%` da el resto de una división, así que
`n % 2 == 0` es verdadero para los pares.

**2.13** Escribe `func raizCuadrada(n float64) (float64, error)` que devuelva error si `n` es negativo.
Usa `math.Sqrt` (necesitas `import "math"`). Pruébala con 16 y con -4, atendiendo el error en los dos
casos.

**2.14** Toma la **Fig. 2.10** y agrégale una columna que diga `SI` o `NO` según si el puerto es 443.
Necesitas una función nueva y ajustar la plantilla de `Sprintf`.

**2.15 (Encuentra el error)** Cada uno de estos fragmentos tiene un problema. Dilo sin compilar, y después
compila para confirmar:

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

**2.16 (Proyecto del curso)** Éste es el primer paso del programa que vas a construir durante todo el
curso. Escribe un programa que:

1. Tenga una función `revisar(nombre string, puerto int, ms int) string` que devuelva una línea de reporte.
2. Tenga una función `clasificar(ms int) string`.
3. Imprima un encabezado y cinco servicios.
4. Al final, imprima cuántos de los cinco fueron `"rapido"`. Necesitarás una variable contador y un `if`.

Guárdalo: en la lección 3 lo vas a reorganizar con structs.

## Cómo sé que lo logré

- Compilaste y corriste, tú mismo y sin copiar/pegar, las figuras 2.1 a 2.10 de esta lección, y obtuviste
  exactamente la salida mostrada.
- Puedes explicar, sin ver el texto, por qué Go es un lenguaje compilado y qué consecuencia práctica tiene
  eso para atrapar errores antes de producción.
- Sabes de memoria cuándo usar `:=` y cuándo `=`, y qué error da el compilador si te equivocas.
- Puedes nombrar los cuatro tipos básicos (`string`, `int`, `float64`, `bool`) y el valor cero de cada uno,
  sin consultar la tabla.
- Escribiste y corriste tu propia versión del ejercicio **2.16** (el primer paso del `revisor`): compila,
  imprime el encabezado, las cinco filas y el conteo de servicios rápidos.
- Puedes explicar en dos frases por qué `if err != nil` aparece en todas partes en código Go, y qué
  problema evita frente a las excepciones de otros lenguajes.

Abre [`bitacora.md`](bitacora.md) y anota **dos cosas**: el error del compilador que más te costó
entender, y algo que te haya sorprendido. En tres semanas te va a parecer obvio y no vas a recordar por
qué te costó — y eso es justo lo que conviene tener escrito.

---

## Resumen

- Go es un lenguaje **compilado**: traduce todo tu código una vez y produce un ejecutable independiente.
- Todo archivo pertenece a un **paquete**; el paquete **`main`** es el único que produce un ejecutable.
- La ejecución empieza en la función **`main`**.
- **`import`** trae paquetes; Go no carga nada por omisión.
- **`:=`** declara y asigna; **`=`** solo asigna.
- Go **no compila** si declaras una variable y no la usas, ni si mezclas tipos.
- Los cuatro tipos básicos son **`string`**, **`int`**, **`float64`** y **`bool`**.
- El **valor cero** garantiza que toda variable nace con un valor válido: `""`, `0`, `0`, `false`.
- **`Println`** imprime con salto de línea; **`Printf`** usa plantilla con `%`; **`Sprintf`** devuelve el
  texto en lugar de imprimirlo.
- Una función de Go puede devolver **varios valores**, y el patrón estándar es devolver
  **`(resultado, error)`**.
- **`nil`** significa «nada». **`if err != nil`** es la forma idiomática de comprobar si hubo error, y la
  línea más escrita en Go.
- `go run` compila y ejecuta sin dejar archivo; `go build` deja el ejecutable.

---

## Para leer más

1. **[A Tour of Go](https://go.dev/tour/)** — oficial, interactivo. Haz la parte de «Basics» hasta donde
   llevas en esta lección.
2. **[Go by Example — Variables](https://gobyexample.com/variables)** y
   **[Go by Example — Multiple Return Values](https://gobyexample.com/multiple-return-values)** — el
   mismo patrón `(resultado, error)` con programas mínimos que corren.
3. **[Effective Go](https://go.dev/doc/effective_go)** — todavía no hace falta leerlo completo, pero
   consúltalo si algo de esta lección se sintió como una regla arbitraria: ahí está el porqué.

Si algo no te quedó claro: lee el mensaje de error completo (Go suele decir la línea, la columna y qué
esperaba), busca el concepto en Go by Example, y anota la duda en la bitácora **aunque no la resuelvas**.
Una duda escrita se puede resolver después; una olvidada, no.

### Términos de esta lección

| | |
|---|---|
| **compilador** | programa que traduce código fuente a instrucciones de máquina |
| **error de compilación** | error que impide producir el ejecutable; se detecta antes de ejecutar |
| **función** | bloque de código con nombre que recibe parámetros y puede devolver valores |
| **identificador vacío (`_`)** | símbolo que descarta un valor a propósito |
| **`int`, `float64`, `string`, `bool`** | los cuatro tipos básicos |
| **lenguaje compilado / interpretado** | traduce todo una vez / línea por línea en cada ejecución |
| **`main` (función)** | punto de entrada del programa |
| **`main` (paquete)** | el único paquete que produce un ejecutable |
| **`nil`** | ausencia de valor |
| **paquete** | agrupación de código relacionado |
| **parámetro** | dato que una función recibe |
| **`Printf` / `Println` / `Sprintf`** | imprimir con plantilla / imprimir con salto / devolver texto |
| **tipo** | clase de dato que una variable puede guardar |
| **valor cero** | valor válido que recibe toda variable declarada sin valor |
| **variable** | espacio de memoria con nombre que guarda un dato |
| **verbo de formato (`%d`, `%s`, `%q`, `%t`, `%v`)** | hueco en una plantilla de `Printf` |

---

**Anterior:** [Lección 1 — Instalar Go](01-instalacion.md) ·
**Siguiente:** [Lección 3 — Structs, métodos, errores e interfaces](03-errores-interfaces.md)
