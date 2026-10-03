# Lección 6 — Concurrencia

**Duración:** 2 sesiones de 60-90 minutos. Es la razón por la que Go existe, y la única lección del curso
que se da en pasadas: goroutines, luego canales y `WaitGroup`, luego `context` y el límite de paralelismo
aplicado al `revisor`. En una sola pasada no cuaja — lo confirma el mismo orden que usa el curso de Go con
más alumnos del mundo.

**Al terminar vas a poder:**

- Lanzar una goroutine y explicar, con una prueba real, por qué el programa puede terminar antes de que
  termine.
- Usar `sync.WaitGroup` para esperar a un grupo de goroutines, y reconocer el mensaje exacto que da Go
  cuando el `Add`/`Done` no cuadran.
- Explicar por qué el `revisor` reparte resultados por un canal en vez de escribir un mapa compartido, y
  provocar tú mismo el error que pasaría si no lo hiciera.
- Leer, literal, la salida de `go test -race` cuando encuentra una carrera de datos real.
- Usar `context.WithTimeout` para poner un límite de tiempo a una operación, y explicar por qué
  `defer cancelar()` no es opcional.
- Limitar cuántas goroutines corren a la vez con un semáforo de canal, y medir que el límite de verdad se
  respeta.

---

## Por qué importa

Hasta la lección 5 el `revisor` consulta los servicios **uno por uno**: si tienes diez servicios y cada
uno tarda un segundo en contestar, el reporte completo tarda diez segundos, sea cual sea el orden. Con
cien servicios, casi dos minutos. Y el tiempo de espera no lo decides tú: lo decide el servicio más lento
de la lista, multiplicado por cuántos haya.

Esto no tiene por qué ser así, porque **consultar un servicio es esperar, no trabajar**: casi todo el
tiempo de esa espera la CPU no está haciendo nada, solo aguardando una respuesta de red. Go fue diseñado
por gente que en Google pasaba los días esperando exactamente eso —respuestas de red entre miles de
servicios— y por eso la concurrencia no es una librería que agregaste después: es parte del lenguaje
desde la primera línea (`go`, una palabra reservada, no una función de una biblioteca).

**El resultado que vas a construir en esta lección:** los mismos diez servicios de un segundo cada uno,
consultados **todos a la vez**, en poco más de un segundo total en vez de diez. Ese número —la diferencia
entre "en serie" y "en paralelo"— es el premio de la lección, y lo vas a medir tú mismo, no a creerlo de
oídas.

🔑 **Y la advertencia que hace esta lección distinta de las anteriores:** en las lecciones 2 a 5, un
programa que compila y cuyas pruebas pasan casi siempre está bien. **En concurrencia no.** Un programa
con una carrera de datos puede compilar, correr, pasar sus pruebas noventa y nueve veces, y fallar la
centésima — o nunca fallar en tu máquina y fallar todos los días en producción, con más núcleos y más
carga. Los errores de esta lección son los que **no se ven a simple vista**, y por eso el punto 4 de esta
lección —los errores reales, provocados y capturados— es el más importante del curso entero.

---

## Los conceptos

### 6.1 Goroutines: arrancar es trivial, esperar no

Lanzar una goroutine es una palabra:

<!-- verificar:fragmento -->
```go
go revisar(s)
```

Y ya. Eso basta para que `revisar(s)` corra **concurrentemente** con el resto del programa, sin esperar a
que termine para seguir con la siguiente línea. Cuestan aproximadamente 2 KB de memoria cada una para
empezar (crecen si hace falta), no los megabytes de un hilo del sistema operativo — por eso un programa
en Go puede tener cientos de miles de goroutines vivas sin planear nada especial, algo que sería
impensable con hilos del sistema.

Pero mira este programa, con el bug más común del día uno de concurrencia en Go (completo en
`programas/revisor/ejemplos/06-goroutines-sin-esperar/main.go`):

<!-- verificar:ejemplo:ejemplos/06-goroutines-sin-esperar:nodeterminista -->
```go
func main() {
	servicios := []string{"catalogo", "pagos", "inventario"}
	for _, s := range servicios {
		go fmt.Println("revisando", s)
	}
	// el programa termina AQUÍ, sin haber esperado a ninguna goroutine
}
```

Lo corrí **260 veces** en esta máquina para medir qué tan seguido se VE el problema, no para adivinarlo:
40 con `go run`, 20 más forzando un solo procesador lógico (`GOMAXPROCS=1`, para quitarle al programa
cualquier chance de que una goroutine alcance a correr en otro núcleo mientras `main` termina), y 200 con
el binario ya compilado (`go build` + `./gor61`, para quitar de en medio el tiempo que toma compilar). El
resultado:

```
$ go run main.go
$ go run main.go
$ go run main.go
```

**Vacío. Las tres, y prácticamente todas las 260** (una sola de las 200 corridas del binario compilado
imprimió algo — el resto, nada).

🔴 **Y hay que leer ese número con mucho cuidado, porque es fácil sacar la conclusión equivocada. El
defecto no está presente «1 de cada 260 veces»: está presente las 260 de 260, sin excepción.** En
**ninguna** de las 260 corridas el programa esperó a sus goroutines — ni siquiera en la única que
imprimió algo. Lo que cambia entre corridas no es si el programa espera (nunca lo hace): es si, por pura
suerte de cómo el sistema operativo reparte el tiempo entre procesos, alguna goroutine alcanza a ejecutar
`fmt.Println` en la rendija de tiempo que hay entre que se lanza y que `main` mata al programa entero. Esa
rendija casi nunca alcanza — por eso casi nunca se ve nada — pero el programa está **igual de roto** en
las 259 corridas silenciosas que en la que sí imprimió.

**Así hay que quedarse con esto: siempre roto, casi nunca visible.** No es "hay una probabilidad chica de
que falle" —leerlo así es exactamente el error que hay que evitar—: es que el programa **nunca** hace lo
correcto, y la mayoría de las veces el síntoma no llega a asomarse a tiempo para que lo notes. Así es como
este tipo de bug se cuela: pasa la revisión de código (compila, corre, no truena), pasa las pruebas
manuales (¿quién corre un programa 260 veces para confiar en él?), pasa el CI (que corre una vez, tal vez
dos) — y avisa hasta que ya está en producción, con más carga, más núcleos, y una rendija de tiempo que un
día sí alcanza a abrirse, en el peor momento posible para descubrirlo.

⚠️ **No lo tomes tampoco como "esto nunca imprime nada, en ninguna máquina".** Con más carga, otro sistema
operativo, o simple mala suerte, la rendija se puede abrir más seguido — de hecho se abrió una vez en
estas mismas 260 corridas. Lo que no cambia entre máquinas es que el programa **nunca** espera: eso es
estructural, no depende de la suerte. Corre tú mismo el experimento en tu equipo, con suficientes
repeticiones para que el número signifique algo, y compara.

**No es un bug intermitente del programa: es que `main` termina en cuanto llega al final de su cuerpo,
sin importarle si hay goroutines todavía corriendo** — y cuando `main` termina, el programa entero
termina con él, goroutines vivas incluidas, sin aviso y sin error. Lanzar una goroutine es fácil; el
trabajo de verdad es asegurarte de que el programa espera a que terminen antes de seguir.

### 6.2 `sync.WaitGroup`: la forma correcta de esperar

<!-- verificar:fragmento -->
```go
var wg sync.WaitGroup
for _, s := range servicios {
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("revisando", s)
	}()
}
wg.Wait() // aquí sí espera a que las tres hayan llamado Done
```

El patrón es siempre el mismo, y conviene memorizarlo en este orden:

1. **`wg.Add(1)` antes de lanzar la goroutine**, no adentro de ella. Si lo pones adentro, `Wait()` puede
   ejecutarse antes de que la goroutine alcance a hacer su propio `Add`, y entonces no cuenta con esa
   espera.
2. **`defer wg.Done()` como primera línea de la goroutine.** El `defer` garantiza que se ejecute pase lo
   que pase adentro —incluso si la goroutine hace panic—, y ponerlo primero evita que un `return`
   temprano en medio del código se te olvide de contarlo.
3. **`wg.Wait()` donde de verdad necesitas el resultado**, normalmente justo antes de usar lo que las
   goroutines produjeron.

**Qué pasa si te saltas el paso 1 — `Add` faltante —, provocado a propósito** (completo en
`programas/revisor/ejemplos/06-waitgroup-sin-add/main.go`):

<!-- verificar:ejemplo:ejemplos/06-waitgroup-sin-add -->
```go
var wg sync.WaitGroup
for i := 0; i < 5; i++ {
	go func(n int) { // <- sin wg.Add(1) antes de esta línea
		defer wg.Done()
		time.Sleep(10 * time.Millisecond)
	}(i)
}
wg.Wait()
fmt.Println("listo")
```

Lo corrí **seis veces seguidas, con `-race`**, para que veas lo que de verdad pasa, no solo el error
final:

```
$ go run -race sinesperar.go
listo
panic: sync: negative WaitGroup counter
...
exit status 2

$ go run -race sinesperar.go
listo
panic: sync: negative WaitGroup counter
...
exit status 2

(cuatro corridas más, idénticas: "listo" primero, el panic después — 6 de 6, CON -race)
```

🔑 **Lee eso con cuidado, porque es el error más engañoso de la lección: el programa imprime `listo`
ANTES de tronar, las seis veces.** No es una casualidad de esta corrida: `wg.Wait()` no espera nada
porque el contador nunca se incrementó (nunca hubo un `Add(1)`), así que `Wait()` regresa **de
inmediato**, `main` imprime `listo` como si todo hubiera salido bien, y **recién entonces** una de las
goroutines rezagadas —que sí estaba corriendo, solo que Go nunca esperó por ella— termina su `Sleep` y
llama a `Done()`, restándole uno a un contador que ya estaba en cero. El panic no es lo primero que ves:
es lo último, después de que el programa ya te dijo que todo estaba bien.

⚠️ **Y aquí viene la parte que medí mal la primera vez, y vale la pena dejar escrita la corrección: sin
`-race`, este mismo programa no truena nunca.** Las mismas seis corridas, sin `-race`:

```
$ go run sinesperar.go
listo

$ go run sinesperar.go
listo

(cuatro corridas más, idénticas: "listo" y nada más — 6 de 6, SIN -race, saliendo con código 0)
```

**El `time.Sleep(10 * time.Millisecond)` no es lo que HACE aparecer el panic — es lo que lo IMPIDE**, sin
`-race`. `main` llega a `wg.Wait()`, que regresa de inmediato, imprime `listo`, y el programa entero
termina **mucho antes** de que se cumplan los 10 milisegundos — así que ninguna de las cinco goroutines
rezagadas alcanza siquiera a despertar de su `Sleep` y llamar a `Done()`. El programa se ve, literalmente,
"terminado con éxito": código de salida 0, sin ningún rastro de que algo estaba mal armado. Lo que sí hace
aparecer el panic, de forma consistente, es la instrumentación de `-race`: vigilar cada acceso a memoria
para detectar carreras hace que el programa corra notablemente más lento, y esa lentitud extra es
exactamente lo que le da tiempo a alguna goroutine rezagada de completar su `Sleep` y llamar a `Done()`
antes de que el proceso termine.

**Esto es exactamente el tipo de error que no falla siempre, y por eso es el más valioso de esta
lección — y ahora con el dato correcto: acabas de medir que ni siquiera hace falta "un trabajo más corto"
para que pase inadvertido. Corriéndolo tal cual, SIN `-race`, ya pasa inadvertido las 6 de 6 veces.** Un
alumno que corra este ejemplo sin `-race` —lo obvio, si nadie se lo advierte— va a ver `listo` y nada
más, y va a concluir, con toda razón, que el programa funciona. **No funciona: el `Add(1)` sigue faltando
exactamente igual.** Lo único que cambió es si algo alcanzó a asomarse a tiempo para delatarlo — el mismo
patrón "siempre roto, casi nunca visible" de la sección 6.1, aquí con un actor distinto: no es la carga de
la máquina, es si corriste con el detector de carreras encendido o no.

**`Done()` resta uno al contador interno del `WaitGroup`.** Si nunca hiciste `Add(1)`, el contador
empieza en cero, y restarle uno lo vuelve negativo — y Go, en vez de dejarlo pasar en silencio, hace
panic con un mensaje que dice exactamente qué está mal: `negative WaitGroup counter`. Ese mensaje literal
es tu pista: si lo ves, casi siempre falta un `Add(1)` en algún lugar, o hay más `Done()` de los que
debería haber. Pero la lección real no es el mensaje del panic — es que **el programa ya había dicho
`listo` antes de que aparezca.**

### 6.3 Por qué el `revisor` no comparte un mapa: canales

El lema del lenguaje, y vale la pena memorizarlo tal cual: **«no comuniques compartiendo memoria; comparte
memoria comunicando».** En vez de que cada goroutine escriba su resultado en un mapa o slice compartido
—lo que exigiría proteger cada acceso con un candado—, cada una manda su resultado por un **canal**, y
una sola goroutine (o el código principal) los recoge del otro lado.

<!-- verificar:fragmento -->
```go
ch := make(chan servicio.Estado)         // sin buffer: cada envío espera a que alguien reciba
ch := make(chan servicio.Estado, 10)     // con buffer: hasta 10 caben sin esperar a que nadie reciba

ch <- estado        // enviar
e := <-ch           // recibir
close(ch)            // cerrar: después de cerrado, recibir sigue funcionando hasta vaciarlo
for e := range ch { ... }   // recibe hasta que el canal se cierre Y se vacíe
```

**Comprobé qué pasa si en vez de un canal uso un slice compartido sin protección**, exactamente el error
que los canales evitan. Este programa lanza mil goroutines que incrementan la misma variable (completo en
`programas/revisor/ejemplos/06-carrera-de-datos/main.go`):

<!-- verificar:ejemplo:ejemplos/06-carrera-de-datos:nodeterminista -->
```go
contador := 0
var wg sync.WaitGroup
for i := 0; i < 1000; i++ {
	wg.Add(1)
	go func() {
		defer wg.Done()
		contador++ // dos goroutines pueden leer el mismo valor antes de que la otra escriba
	}()
}
wg.Wait()
fmt.Println("contador:", contador)
```

Sin el detector de carreras, **el programa no truena — solo da un resultado incorrecto**, y distinto cada
vez:

```
$ go run carrera.go
contador: 956
```

**956, no 1000.** Ningún error, ningún panic, ninguna señal de que algo salió mal — solo un número más
chico del que debería ser, porque algunas de las mil sumas se perdieron cuando dos goroutines leyeron
`contador` al mismo tiempo, ambas sumaron uno al mismo valor viejo, y una de las dos escrituras sobrescribió
a la otra sin que nadie se enterara. **Este es el bug más peligroso de la concurrencia: no falla, da un
resultado plausible y equivocado.** Un programa así puede pasar revisión, pasar pruebas superficiales, y
fallar en producción de forma esporádica durante meses antes de que alguien lo note.

### 6.4 `go test -race`: ver el detector encontrar el bug

El mismo programa, corrido con `-race`, sí lo delata — con evidencia exacta de dónde:

```
$ go run -race carrera.go
==================
WARNING: DATA RACE
Read at 0x00c00013c018 by goroutine 11:
  main.main.func1()
      carrera.go:15 +0x68

Previous write at 0x00c00013c018 by goroutine 8:
  main.main.func1()
      carrera.go:15 +0x78

Goroutine 11 (running) created at:
  main.main()
      carrera.go:13 +0x6c

Goroutine 8 (finished) created at:
  main.main()
      carrera.go:13 +0x6c
==================
contador: 847
Found 2 data race(s)
exit status 66
```

Léelo de arriba abajo, porque cada bloque contesta una pregunta distinta:

- **`Read at ... by goroutine 11`** y **`Previous write at ... by goroutine 8`**: dos goroutines
  distintas tocaron la **misma dirección de memoria** (`0x00c00013c018`, la variable `contador`), una
  leyendo y otra escribiendo, sin ningún mecanismo que garantice que una espere a la otra.
- **`carrera.go:15`** en ambas: la línea exacta del código donde pasó (`contador++`), la misma línea para
  las dos, porque es la única línea que toca esa variable.
- **`Goroutine 11 (running) created at ... carrera.go:13`**: de dónde salió esa goroutine —la línea del
  `go func() {...}()` dentro del `for`—, para que puedas rastrear cuál lanzamiento fue.
- **`Found 2 data race(s)`** y **`exit status 66`**: el detector no se detiene en la primera carrera que
  encuentra; sigue corriendo y las reporta todas, y el programa termina con un código de salida distinto
  de 0 y de 1 (66 es el código que Go reserva para esto), para que un pipeline de integración continua lo
  pueda distinguir de una falla normal.

**Y el número final, `contador: 847`, cambió respecto a la corrida sin `-race` (956).** No por
casualidad: `-race` hace que el programa corra más lento y con más instrumentación, lo que cambia el
orden exacto en que las goroutines se entrelazan — otra razón por la que este bug es tan traicionero: el
número equivocado ni siquiera es el mismo equivocado cada vez.

🔴 **La regla que se sigue de esto, sin excepción:** un programa concurrente que pasa sus pruebas **sin**
`-race` no está probado. Corre `-race` desde la primera prueba de concurrencia que escribas, no cuando
"algo se vea raro" — porque, como acabas de ver, nada se ve raro hasta que ya es tarde.

### 6.5 Los otros dos modos de fallar: panic y deadlock

**Panic, si cierras mal un canal:**

- Enviar a un canal ya cerrado hace panic: `panic: send on closed channel`.
- Cerrar un canal dos veces hace panic: `panic: close of closed channel`.

La regla que evita los dos: **cierra el canal quien envía, nunca quien recibe**, y ciérralo una sola vez,
normalmente desde una goroutine dedicada que sabe cuándo ya no va a haber más envíos (lo vas a ver en la
sección 6.7, con `wg.Wait()` seguido de `close`).

**Deadlock, si nadie del otro lado está escuchando.** Lo provoqué con el programa más corto posible
(completo en `programas/revisor/ejemplos/06-deadlock/main.go`):

<!-- verificar:ejemplo:ejemplos/06-deadlock:nodeterminista -->
```go
func main() {
	canal := make(chan int)
	canal <- 1 // nadie del otro lado está leyendo: se bloquea
	fmt.Println(<-canal)
}
```

```
$ go run candado.go
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send]:
main.main()
	candado.go:7 +0x38
exit status 2
```

**Esto corrige algo que muchas explicaciones dan por hecho: un deadlock en Go no siempre "se cuelga para
siempre".** El propio runtime de Go tiene un detector de deadlocks: si en algún momento **todas** las
goroutines del programa están dormidas esperando algo que nunca va a pasar, Go lo nota y mata el programa
con `fatal error: all goroutines are asleep - deadlock!`, con la línea exacta donde se atoró
(`[chan send]`, en este caso, porque estaba enviando). Si tu programa parece "colgado para siempre" en
vez de terminar con este error, casi seguro **no** es un deadlock verdadero: es más probable que tengas
al menos una goroutine viva haciendo otra cosa (por ejemplo, un temporizador o un servidor HTTP
escuchando), y ésa evita que el runtime declare que "todas" están dormidas.

### 6.6 `context`: cómo se cancela y se pone un límite de tiempo

<!-- verificar:fragmento -->
```go
ctx, cancelar := context.WithTimeout(context.Background(), 5*time.Second)
defer cancelar() // SIEMPRE, incluso si terminas antes de que se cumplan los 5 segundos

req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
resp, err := http.DefaultClient.Do(req) // se aborta solo si pasan los 5 s
```

`context` es el estándar de la casa en Go para dos cosas: **poner un límite de tiempo** a una operación
que podría tardar demasiado, y **propagar una cancelación** hacia abajo (si el de arriba se cancela,
todo lo que dependa de él se cancela también, sin que cada función tenga que reinventar su propio
mecanismo).

⚠️ **`defer cancelar()` no es opcional, incluso cuando la operación ya terminó por su cuenta.**
`context.WithTimeout` arranca un temporizador interno; si nunca llamas a la función de cancelar, ese
temporizador sigue vivo hasta que se cumple el plazo original, reteniendo memoria durante todo ese tiempo
— y si tu programa crea contextos así todo el tiempo (uno por cada servicio consultado, como el
`revisor`), sin `defer cancelar()` acumulas fugas de memoria proporcionales a cuántas consultas hiciste.
Es la fuga más común de programas en Go que usan `context`, y por eso conviene escribir el `defer` en la
misma línea donde creas el contexto, antes de escribir cualquier otra cosa.

### 6.7 `Todos`: la función que junta las tres piezas

Así quedó, de verdad, la función central del `revisor` (`programas/revisor/internal/revisar/todos.go`), después de
juntar goroutines, canales, semáforo y `context`:

<!-- verificar:extracto:internal/revisar/todos.go -->
```go
func Todos(ctx context.Context, r Revisor, servicios []servicio.Servicio, paralelo int) []servicio.Estado {
	if len(servicios) == 0 {
		return nil
	}
	if paralelo <= 0 {
		paralelo = ParaleloPorOmision
	}

	resultados := make(chan servicio.Estado, len(servicios))
	semaforo := make(chan struct{}, paralelo)

	var wg sync.WaitGroup
	for _, s := range servicios {
		wg.Add(1)
		go func() {
			defer wg.Done()

			semaforo <- struct{}{}        // pide turno; se bloquea si ya hay «paralelo» corriendo
			defer func() { <-semaforo }() // devuelve el turno pase lo que pase

			// Cada servicio tiene su propio tiempo límite, hijo del general. Si
			// el de arriba se cancela, este muere con él.
			propio, cancelar := context.WithTimeout(ctx, s.TimeoutEfectivo())
			defer cancelar() // sin esto, el temporizador no se libera: es la fuga más común de Go

			resultados <- r.Revisar(propio, s)
		}()
	}

	// Quien envía cierra, nunca quien recibe. Y se cierra desde otra goroutine
	// porque Wait tiene que poder esperar mientras el bucle de abajo ya está
	// recibiendo: si cerráramos aquí mismo, con el canal lleno nos trabaríamos.
	go func() {
		wg.Wait()
		close(resultados)
	}()

	estados := make([]servicio.Estado, 0, len(servicios))
	for e := range resultados {
		estados = append(estados, e)
	}
	return estados
}
```

Léela con las secciones anteriores frescas, porque cada pieza responde a un problema que ya viste:

- **Una goroutine por servicio** (6.1), contada con **`wg.Add(1)`/`defer wg.Done()`** (6.2) para que el
  programa sepa cuándo terminaron todas.
- **Un canal `resultados` con buffer** (6.3) recoge los estados: nadie escribe en un slice o mapa
  compartido, así que no hace falta ningún candado para esa parte.
- **`semaforo := make(chan struct{}, paralelo)`** es un canal usado como cupo: tiene espacio para
  `paralelo` valores, así que la `paralelo + 1`-ésima goroutine que intenta escribir en él (`semaforo <-
  struct{}{}`) se bloquea hasta que otra libere su lugar (`<-semaforo`, en el `defer`). Es el mismo
  canal con buffer de la sección 6.3, usado no para llevar datos sino para llevar la cuenta de cuántos
  "turnos" quedan.
- **Un `context.WithTimeout` propio por servicio** (6.6), hijo del `ctx` general: si el de arriba se
  cancela (por ejemplo, si se agota el tiempo total del reporte), todos los hijos se cancelan con él.
- **`close(resultados)` desde OTRA goroutine, después de `wg.Wait()`** — y esto merece explicarse,
  porque es la parte que menos se ve a simple vista: si cerráramos el canal en la misma goroutine que
  hace el `for e := range resultados` de abajo, nos trabaríamos, porque `Wait()` necesita que todas las
  goroutines lean del canal para liberar espacio y poder devolver su turno, pero el bucle de lectura
  nunca arrancaría porque estaríamos esperando a `Wait()` primero. Lanzándolo aparte, el cierre y la
  lectura ocurren **al mismo tiempo**, no uno después del otro.

### 6.8 Probarlo sin red: `Falso` y el límite de paralelismo medido

Probar `Todos` contra servicios de verdad sería lento y no determinista. El `revisor` usa un `Revisor`
falso (`programas/revisor/internal/revisar/falso.go`) que simula respuestas, demoras e incluso servicios que nunca
contestan, todo en memoria:

<!-- verificar:fragmento -->
```go
type Falso struct {
	Respuestas map[string]RespuestaFalsa
	mu         sync.Mutex
	Contador   int
}

func (f *Falso) Revisar(ctx context.Context, s servicio.Servicio) servicio.Estado {
	f.mu.Lock()
	f.Contador++
	f.mu.Unlock()
	// ...
}
```

🔒 **Aquí sí hace falta un mutex, y es la excepción a la regla de "prefiere canales" de la sección 6.3.**
`Contador` es un entero compartido que **muchas goroutines incrementan al mismo tiempo** durante las
pruebas (`Falso.Revisar` se llama concurrentemente, una vez por cada servicio que `Todos` está
consultando). Un canal serviría para *mandar* resultados, pero para un contador simple compartido, un
`sync.Mutex` alrededor de la única línea que lo toca es más simple y más claro. La regla no es "nunca uses
un mutex": es "antes de usar uno, pregúntate si un canal expresa mejor lo que estás haciendo" — y para
mandar resultados, casi siempre sí; para un contador compartido, casi siempre el mutex es la herramienta
correcta.

Con `Falso`, esta prueba mide algo que de otra forma sería casi imposible de comprobar con confianza: que
el semáforo de la sección 6.7 de verdad limita cuántas consultas corren a la vez, no solo que "funciona
en general" (versión completa, sin abreviar, en `programas/revisor/internal/revisar/todos_test.go`):

<!-- verificar:fragmento -->
```go
func TestTodos_elParaleloLimitaCuantasCorrenALaVez(t *testing.T) {
	const totalServicios = 20
	const limite = 3
	var enVuelo, maximoObservado int32
	// medidorDeConcurrencia cuenta, con un contador atómico, cuántas llamadas
	// a Revisar están abiertas AL MISMO TIEMPO, y se queda con el máximo.
	medidor := &medidorDeConcurrencia{enVuelo: &enVuelo, maximoObservado: &maximoObservado, espera: 15 * time.Millisecond}
	servicios := make([]servicio.Servicio, totalServicios)
	// ... llena servicios ...
	Todos(context.Background(), medidor, servicios, limite)
	if max := atomic.LoadInt32(&maximoObservado); max > int32(limite) {
		t.Errorf("se observaron %d consultas simultáneas; el límite era %d", max, limite)
	}
}
```

Corrida real, con el detector de carreras activo (para confirmar que ni siquiera la propia medición
introduce una carrera):

```
$ go test ./internal/revisar/... -race -v -run TestTodos_elParaleloLimitaCuantasCorrenALaVez
=== RUN   TestTodos_elParaleloLimitaCuantasCorrenALaVez
--- PASS: TestTodos_elParaleloLimitaCuantasCorrenALaVez (0.11s)
PASS
ok  	github.com/habil/revisor/internal/revisar	1.435s
```

**20 servicios, límite de 3, y la prueba confirma que nunca hubo más de 3 llamadas a `Revisar` corriendo
a la vez** — no porque lo asumamos del código, sino porque un contador atómico lo midió mientras corría.

### 6.9 El tiempo, medido: en serie contra en paralelo

Con `Falso` configurado para simular demoras reales, la prueba `TestTodos_respetaElTimeoutPorServicio`
confirma el otro lado de la moneda: un servicio que nunca contesta (`Colgado: true`) con un
`Timeout: 30 * time.Millisecond` no puede hacer que `Todos` tarde más de eso:

```
$ go test ./internal/revisar/... -run TestTodos_respetaElTimeoutPorServicio -v
=== RUN   TestTodos_respetaElTimeoutPorServicio
--- PASS: TestTodos_respetaElTimeoutPorServicio (0.03s)
```

**0.03 segundos, no más.** El `context.WithTimeout` de la sección 6.7, hijo del general, cortó esa
consulta exactamente cuando debía, sin que el resto del programa tuviera que enterarse de que pasó.

---

## El error que vas a ver

Todos estos son literales, provocados a propósito para esta lección:

| El síntoma | Mensaje / evidencia literal | Qué pasa y qué hacer |
|---|---|---|
| `main` termina antes de que las goroutines corran | (ninguna salida, o una salida parcial e inconsistente entre corridas) | Falta un `sync.WaitGroup` (o un canal) que haga que `main` espere. Sección 6.1 |
| `Add`/`Done` no cuadran | El programa imprime que **todo salió bien primero**, y `panic: sync: negative WaitGroup counter` llega DESPUÉS | Falta un `wg.Add(1)` antes de lanzar alguna goroutine, o sobra un `Done()`. `Wait()` regresó sin esperar nada. Sección 6.2 |
| Datos compartidos sin protección | `WARNING: DATA RACE` + `Found 2 data race(s)` + `exit status 66` (con `-race`); un número incorrecto sin explicación, sin `-race` | Dos goroutines tocan la misma memoria sin sincronización. Usa un canal para mandar el resultado, o un mutex si de verdad necesitas un contador compartido. Secciones 6.3 y 6.4 |
| Enviar a un canal cerrado | `panic: send on closed channel` | Alguien más ya cerró el canal, o lo cerró quien no debía. Cierra siempre desde quien envía, una sola vez |
| Cerrar dos veces el mismo canal | `panic: close of closed channel` | Dos goroutines (o dos caminos de código) intentan cerrar el mismo canal. Centraliza el cierre en un solo lugar |
| Nadie del otro lado de un canal sin buffer | `fatal error: all goroutines are asleep - deadlock!` | El runtime de Go detecta que todo el programa quedó dormido esperando algo que nunca va a pasar, y lo mata con esta línea exacta. Sección 6.5 |
| Olvidaste `defer cancelar()` | (sin error inmediato; fuga de memoria acumulada con el tiempo) | Cada `context.WithTimeout` sin cancelar mantiene vivo su temporizador hasta que expira por sí solo. Sección 6.6 |

**Y la instrucción que resume esta lección entera:** si vas a escribir concurrencia, corre `-race` desde
tu primera prueba, no cuando algo huela mal — porque, como viste en la sección 6.3, nada huele mal hasta
que ya es tarde.

---

## Lo que se hace mal

- **Lanzar una goroutine por elemento, sin ningún límite, contra un recurso externo.** Con 5 servicios no
  se nota. Con 5,000, el programa abre 5,000 conexiones de golpe, y el cuello de botella deja de ser el
  servicio remoto para ser tu propia máquina (o la red, o el propio servidor al que le llueven 5,000
  peticiones simultáneas). El semáforo de la sección 6.7 existe exactamente para esto — nunca dejes que
  "cuántas goroutines lanzo" dependa solo de "cuántos elementos tengo".
- **Ignorar el `context` que te dan, o no propagarlo.** Si una función recibe un `ctx` y llama a otra
  operación que puede tardar sin pasárselo, esa operación no se va a cancelar cuando el `ctx` original lo
  haga — tienes dos relojes que no se hablan. Todo lo que pueda tardar recibe el `ctx` de quien lo llamó.
- **No cerrar lo que se abre.** Un `context.WithTimeout` sin su `cancelar()`, una conexión HTTP sin
  `resp.Body.Close()` (lección 7), un archivo sin cerrar: cada uno es una fuga distinta, pero la forma de
  evitarlas es la misma — un `defer` justo después de abrir, antes de escribir cualquier otra línea.
- **Compartir memoria en vez de comunicarla "porque es más rápido de escribir".** Un mapa compartido con
  un mutex alrededor de todo el bloque puede parecer más corto que armar un canal, pero es más fácil
  olvidar un `Lock()` en un solo punto de acceso que olvidar mandar por un canal — y el primer olvido no
  se nota hasta que el detector de carreras (o, peor, producción) lo encuentra.
- **Confiar en que "nunca ha fallado" significa "está bien".** Una carrera de datos puede pasar noventa y
  nueve corridas y fallar la centésima, o solo fallar con más núcleos que los de tu laptop. La única
  prueba de que un programa concurrente no tiene carreras es correrlo con `-race`, no verlo pasar varias
  veces sin él.

---

## Ejercicios

1. Escribe el programa de la sección 6.1 (goroutines sin esperar) y córrelo cinco veces seguidas. Anota
   cuántas veces imprimió algo y cuántas no.
2. Agrégale un `sync.WaitGroup` correcto (sección 6.2) y confirma que ahora las tres líneas se imprimen
   siempre, en las cinco corridas.
3. Reproduce el bug de `Add` faltante de la sección 6.2 y pega el `panic` completo en tu bitácora.
4. Reproduce la carrera de datos de la sección 6.3 (un contador compartido sin protección) y corre el
   mismo programa con y sin `-race`. Compara los dos números finales y pega la salida completa de
   `-race` en tu bitácora.
5. Reproduce el deadlock de la sección 6.5 con un canal sin buffer y sin receptor. Confirma que ves el
   `fatal error: all goroutines are asleep - deadlock!`, no un colgado silencioso.
6. Toma tu propio `revisor` (o el de este curso) y corre `TestTodos_elParaleloLimitaCuantasCorrenALaVez`
   cambiando el `limite` a 1 y luego a 10. Explica, con tus palabras, por qué el resultado de la prueba
   no cambia (siempre pasa) pero el **tiempo** que tarda sí.
7. (Un poco más difícil) Quita el semáforo de `Todos` (deja que se lancen todas las goroutines sin
   límite) y vuelve a correr `TestTodos_elParaleloLimitaCuantasCorrenALaVez`. Confirma que ahora falla, y
   pega el mensaje de error exacto que da `t.Errorf` con el máximo que sí se observó.

### Soluciones

1-2. No hay número único: depende de tu máquina. Lo que importa es la comparación — sin `WaitGroup`,
inconsistente; con él, siempre las tres líneas.

3. El programa imprime `listo` primero — `Wait()` regresó de inmediato porque el contador nunca se
   incrementó — y el `panic: sync: negative WaitGroup counter` llega después, con una traza que incluye
   `sync.(*WaitGroup).Add(...)`. Corrida varias veces seguidas, el orden es siempre el mismo: `listo`,
   después el panic.

4. Sin `-race`, un número menor a lo esperado (por ejemplo, 956 de 1000), sin ningún error. Con `-race`,
   el bloque `WARNING: DATA RACE` con la línea exacta del código, más `Found N data race(s)` y
   `exit status 66`.

5. `fatal error: all goroutines are asleep - deadlock!`, con `goroutine 1 [chan send]:` (o `[chan
   receive]`, según de qué lado se atoró) y la línea exacta del canal.

6. La prueba pasa en los dos casos porque **mide** el máximo real y lo compara contra el límite que ella
   misma configuró (1 o 10) — nunca contra un número fijo esperado. El tiempo total sí cambia: con
   `limite=1` las 20 consultas son estrictamente en serie (una espera a la otra), con `limite=10` corren
   en dos tandas de 10 en vez de veinte de una.

7. Sin semáforo, todas las 20 goroutines corren a la vez, así que `maximoObservado` va a ser cercano a 20
   (no exactamente el límite de la prueba, que sigue pidiendo 3). El `t.Errorf` dice algo como
   `se observaron 20 consultas simultáneas; el límite era 3` — la prueba SÍ detecta la regresión, que es
   justamente para lo que existe.

---

## Cómo sé que lo logré

- [ ] Vi, con mis propios ojos, un programa terminar sin haber esperado a sus goroutines (sección 6.1).
- [ ] Provoqué el `panic: sync: negative WaitGroup counter` y vi que el programa imprimió `listo` ANTES
      del panic — y puedo explicar por qué.
- [ ] Provoqué una carrera de datos real y vi la diferencia entre correrla con y sin `-race`.
- [ ] Puedo leer un bloque de `WARNING: DATA RACE` y decir qué línea, qué goroutines y qué variable están
      en conflicto.
- [ ] Provoqué el `fatal error: all goroutines are asleep - deadlock!` y sé por qué Go lo detecta en vez
      de quedarse colgado para siempre.
- [ ] Puedo explicar, con el código real de `Todos`, para qué sirve cada una de sus cuatro piezas:
      goroutines, canal de resultados, semáforo y `context` por servicio.
- [ ] Medí, con una prueba real (no de memoria), que el límite de paralelismo del `revisor` sí se
      respeta.
- [ ] Sé explicar por qué `Falso.Contador` usa un mutex en vez de un canal, y por qué eso no contradice
      la regla general de la sección 6.3.

---

## Para leer más

1. [A Tour of Go: Concurrency](https://go.dev/tour/concurrency/1) — el recorrido oficial de goroutines,
   canales y `sync.WaitGroup`, interactivo.
2. [Go Blog: Share Memory By Communicating](https://go.dev/blog/codelab-share) — el artículo original
   donde se explica el lema citado en la sección 6.3.
3. [Package context](https://pkg.go.dev/context) — la documentación oficial, con los cuatro casos de uso
   canónicos (`WithCancel`, `WithTimeout`, `WithDeadline`, `WithValue`).
4. [Data Race Detector](https://go.dev/doc/articles/race_detector) — la documentación oficial de `-race`:
   qué detecta, qué NO detecta (por ejemplo, no encuentra deadlocks, solo carreras de datos), y su costo
   en tiempo de ejecución.

### Términos de esta lección

| Término | Qué significa |
|---|---|
| goroutine | una función que corre concurrentemente con el resto del programa, mucho más barata que un hilo del sistema operativo |
| `sync.WaitGroup` | mecanismo para esperar a que un grupo de goroutines termine, contando `Add`/`Done` |
| carrera de datos (*data race*) | dos goroutines acceden a la misma memoria al mismo tiempo, al menos una escribiendo, sin sincronización |
| canal (`chan`) | el mecanismo de Go para que una goroutine mande datos a otra sin memoria compartida |
| semáforo de canal | un canal con buffer usado como cupo de turnos disponibles, no para llevar datos |
| `context` | mecanismo estándar para propagar límites de tiempo y cancelación entre funciones |
| deadlock | estado en el que todas las goroutines de un programa están dormidas esperando algo que nunca va a pasar; el runtime de Go lo detecta y termina el programa |

---

**Anterior:** [Lección 5 — Módulos y pruebas](05-modulos-y-pruebas.md) ·
**Siguiente:** [Lección 7 — El programa terminado](07-el-programa.md)
