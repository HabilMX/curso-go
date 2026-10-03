# Lesson 6 — Concurrency

**Duration:** 2 sessions of 60-90 minutes. It is the reason Go exists, and the only lesson of the course
that is taught in stages: goroutines, then channels and `WaitGroup`, then `context` and the parallelism
limit applied to the `revisor`. In a single pass it doesn't sink in — that is confirmed by the very same
order used by the Go course with the most learners in the world.

**By the end you will be able to:**

- Launch a goroutine and explain, with a real test, why the program can end before the goroutine finishes.
- Use `sync.WaitGroup` to wait for a group of goroutines, and recognize the exact message Go gives when
  the `Add`/`Done` don't add up.
- Explain why the `revisor` sends results through a channel instead of writing to a shared map, and
  trigger the error that would happen if it didn't.
- Read, literally, the output of `go test -race` when it finds a real data race.
- Use `context.WithTimeout` to put a time limit on an operation, and explain why `defer cancelar()` is not
  optional.
- Limit how many goroutines run at the same time with a channel semaphore, and measure that the limit is
  really respected.

---

## Why it matters

Up to lesson 5 the `revisor` queries the services **one by one**: if you have ten services and each one
takes one second to answer, the complete report takes ten seconds, whatever the order. With a hundred
services, almost two minutes. And you don't decide the waiting time: it is decided by the slowest service
on the list, multiplied by how many there are.

It doesn't have to be this way, because **querying a service is waiting, not working**: for almost all of
that waiting time the CPU isn't doing anything, just waiting for a network response. Go was designed by
people who at Google spent their days waiting for exactly that —network responses between thousands of
services— and that is why concurrency is not a library added afterwards: it is part of the language from
the first line (`go`, a reserved word, not a function from a library).

**The result you are going to build in this lesson:** the same ten services of one second each, queried
**all at once**, in a little over one second in total instead of ten. That number —the difference between
"in series" and "in parallel"— is the prize of the lesson, and you are going to measure it yourself, not
take it on hearsay.

🔑 **And the warning that makes this lesson different from the previous ones:** in lessons 2 to 5, a
program that compiles and whose tests pass is almost always fine. **In concurrency it isn't.** A program
with a data race can compile, run, pass its tests ninety-nine times, and fail the hundredth — or never fail
on your machine and fail every day in production, with more cores and more load. The errors in this lesson
are the ones that **can't be seen with the naked eye**, and that is why point 4 of this lesson —the real
errors, caused and captured— is the most important in the entire course.

---

## The concepts

### 6.1 Goroutines: starting is trivial, waiting isn't

Launching a goroutine takes one word:

<!-- verificar:fragmento -->
```go
go revisar(s)
```

And that's it. That is enough for `revisar(s)` to run **concurrently** with the rest of the program,
without waiting for it to finish before moving on to the next line. They cost approximately 2 KB of memory
each to start (they grow if needed), not the megabytes of an operating system thread — that is why a Go
program can have hundreds of thousands of live goroutines without planning anything special, something
that would be unthinkable with system threads.

But look at this program, with the most common day-one concurrency bug in Go (the full example is in
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

I ran it **260 times** on this machine to measure how often the problem is SEEN, not to guess it: 40 with
`go run`, 20 more forcing a single logical processor (`GOMAXPROCS=1`, to take away from the program any
chance that a goroutine manages to run on another core while `main` finishes), and 200 with the binary
already compiled (`go build` + `./gor61`, to remove the time compiling takes from the picture). The
result:

```
$ go run main.go
$ go run main.go
$ go run main.go
```

**Empty. All three, and practically all 260** (a single one of the 200 runs of the compiled binary printed
something — the rest, nothing).

🔴 **And you have to read that number very carefully, because it is easy to draw the wrong conclusion. The
defect is not present "1 out of every 260 times": it is present 260 out of 260, without exception.** In
**none** of the 260 runs did the program wait for its goroutines — not even in the only one that printed
something. What changes between runs is not whether the program waits (it never does): it is whether, by
pure luck in how the operating system divides time between processes, some goroutine manages to execute
`fmt.Println` in the sliver of time between being launched and `main` killing the whole program. That
sliver is almost never enough — that is why almost nothing is ever seen — but the program is **just as
broken** in the 259 silent runs as in the one that did print.

**So this is what you have to keep: always broken, almost never visible.** It is not "there is a small
probability that it fails" —reading it that way is exactly the error to avoid—: it is that the program
**never** does the right thing, and most of the time the symptom doesn't manage to show up in time for you
to notice it. That is how this kind of bug sneaks in: it passes code review (it compiles, it runs, it
doesn't blow up), it passes manual testing (who runs a program 260 times before trusting it?), it passes
CI (which runs once, maybe twice) — and it doesn't speak up until it is already in production, with more
load, more cores, and a sliver of time that one day does manage to open, at the worst possible moment to
discover it.

⚠️ **Don't take it as "this never prints anything, on any machine" either.** With more load, another
operating system, or simple bad luck, the sliver can open more often — in fact it opened once in these very
260 runs. What doesn't change between machines is that the program **never** waits: that is structural, it
doesn't depend on luck. Run the experiment yourself on your computer, with enough repetitions for the
number to mean something, and compare.

**It is not an intermittent bug in the program: it is that `main` ends as soon as it reaches the end of its
body, without caring whether there are goroutines still running** — and when `main` ends, the whole program
ends with it, live goroutines included, without warning and without error. Launching a goroutine is easy;
the real work is making sure the program waits for them to finish before moving on.

### 6.2 `sync.WaitGroup`: the right way to wait

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

The pattern is always the same, and it is worth memorizing in this order:

1. **`wg.Add(1)` before launching the goroutine**, not inside it. If you put it inside, `Wait()` can run
   before the goroutine manages to do its own `Add`, and then that goroutine is not counted by that wait.
2. **`defer wg.Done()` as the first line of the goroutine.** The `defer` guarantees it runs no matter what
   happens inside —even if the goroutine panics—, and putting it first keeps an early `return` in the
   middle of the code from making you forget to count it.
3. **`wg.Wait()` where you really need the result**, normally right before using what the goroutines
   produced.

**What happens if you skip step 1 — missing `Add` —, caused on purpose** (the full example is in
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

I ran it **six times in a row, with `-race`**, so you can see what really happens, not just the final
error:

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

🔑 **Read that carefully, because it is the most misleading error in the lesson: the program prints
`listo` BEFORE blowing up, all six times.** It is not a coincidence of this run: `wg.Wait()` doesn't wait
for anything because the counter was never incremented (there was never an `Add(1)`), so `Wait()` returns
**immediately**, `main` prints `listo` as if everything had gone well, and **only then** one of the
straggling goroutines —which was running, it is just that Go never waited for it— finishes its `Sleep` and
calls `Done()`, subtracting one from a counter that was already at zero. The panic is not the first thing
you see: it is the last, after the program has already told you everything was fine.

⚠️ **And here comes the part I measured wrong the first time, and it is worth leaving the correction in
writing: without `-race`, this same program never blows up.** The same six runs, without `-race`:

```
$ go run sinesperar.go
listo

$ go run sinesperar.go
listo

(cuatro corridas más, idénticas: "listo" y nada más — 6 de 6, SIN -race, saliendo con código 0)
```

**The `time.Sleep(10 * time.Millisecond)` is not what MAKES the panic appear — it is what PREVENTS it**,
without `-race`. `main` reaches `wg.Wait()`, which returns immediately, prints `listo`, and the whole
program ends **long before** the 10 milliseconds are up — so none of the five straggling goroutines even
manages to wake up from its `Sleep` and call `Done()`. The program looks, literally, "finished
successfully": exit code 0, without any trace that something was wrong. What does make the
panic appear, consistently, is the instrumentation of `-race`: watching every memory access to detect races
makes the program run noticeably slower, and that extra slowness is exactly what gives some straggling
goroutine time to complete its `Sleep` and call `Done()` before the process ends.

**This is exactly the kind of error that doesn't always fail, and that is why it is the most valuable in
this lesson — and now with the correct data: you just measured that you don't even need "a shorter task"
for it to go unnoticed. Running it as is, WITHOUT `-race`, it already goes unnoticed 6 out of 6 times.** A
learner who runs this example without `-race` —the obvious thing, if nobody warns them— is going to see
`listo` and nothing else, and is going to conclude, with every reason, that the program works. **It
doesn't work: the `Add(1)` is still missing exactly the same.** The only thing that changed is whether
something managed to show up in time to give it away — the same "always broken, almost never visible"
pattern from section 6.1, here with a different actor: it is not the machine's load, it is whether you ran
with the race detector on or not.

**`Done()` subtracts one from the `WaitGroup`'s internal counter.** If you never did `Add(1)`, the counter
starts at zero, and subtracting one makes it negative — and Go, instead of letting it pass silently,
panics with a message that says exactly what is wrong: `negative WaitGroup counter`. That literal message
is your clue: if you see it, there is almost always an `Add(1)` missing somewhere, or there are more
`Done()` than there should be. But the real lesson is not the panic message — it is that **the program had
already said `listo` before it appears.**

### 6.3 Why the `revisor` doesn't share a map: channels

The language's motto, and it is worth memorizing as is: **"don't communicate by sharing memory; share
memory by communicating."** Instead of each goroutine writing its result to a shared map or slice —which
would require protecting each access with a lock—, each one sends its result through a **channel**, and a
single goroutine (or the main function) collects them on the other side.

<!-- verificar:fragmento -->
```go
ch := make(chan servicio.Estado)         // sin buffer: cada envío espera a que alguien reciba
ch := make(chan servicio.Estado, 10)     // con buffer: hasta 10 caben sin esperar a que nadie reciba

ch <- estado        // enviar
e := <-ch           // recibir
close(ch)            // cerrar: después de cerrado, recibir sigue funcionando hasta vaciarlo
for e := range ch { ... }   // recibe hasta que el canal se cierre Y se vacíe
```

**I checked what happens if instead of a channel I use a shared slice without protection**, exactly the
error channels avoid. This program launches a thousand goroutines that increment the same variable
(the full example is in `programas/revisor/ejemplos/06-carrera-de-datos/main.go`):

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

Without the race detector, **the program doesn't blow up — it just gives an incorrect result**, and a
different one every time:

```
$ go run carrera.go
contador: 956
```

**956, not 1000.** No error, no panic, no sign that something went wrong — just a smaller number than it
should be, because some of the thousand additions were lost when two goroutines read `contador` at the same
time, both added one to the same old value, and one of the two writes overwrote the other without anyone
finding out. **This is the most dangerous bug in concurrency: it doesn't fail, it gives a plausible and
wrong result.** A program like that can pass review, pass superficial tests, and fail in production
sporadically for months before someone notices.

### 6.4 `go test -race`: watching the detector find the bug

The same program, run with `-race`, does give it away — with exact evidence of where:

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

Read it from top to bottom, because each block answers a different question:

- **`Read at ... by goroutine 11`** and **`Previous write at ... by goroutine 8`**: two different
  goroutines touched the **same memory address** (`0x00c00013c018`, the `contador` variable), one reading
  and the other writing, without any mechanism that guarantees one waits for the other.
- **`carrera.go:15`** in both: the exact line of code where it happened (`contador++`), the same line for
  both, because it is the only line that touches that variable.
- **`Goroutine 11 (running) created at ... carrera.go:13`**: where that goroutine came from —the line of
  the `go func() {...}()` inside the `for`—, so you can trace which launch it was.
- **`Found 2 data race(s)`** and **`exit status 66`**: the detector doesn't stop at the first race it
  finds; it keeps running and reports all of them, and the program ends with an exit code different from
  0 and from 1 (66 is the code Go reserves for this), so that a continuous integration pipeline can tell
  it apart from a normal failure.

**And the final number, `contador: 847`, changed with respect to the run without `-race` (956).** Not by
chance: `-race` makes the program run slower and with more instrumentation, which changes the exact order
in which the goroutines interleave — another reason this bug is so treacherous: the wrong number isn't even
the same wrong number every time.

🔴 **The rule that follows from this, without exception:** a concurrent program that passes its tests
**without** `-race` is not tested. Run `-race` from the first concurrency test you write, not when
"something looks odd" — because, as you just saw, nothing looks odd until it is too late.

### 6.5 The other two ways to fail: panic and deadlock

**Panic, if you close a channel wrongly:**

- Sending to an already closed channel panics: `panic: send on closed channel`.
- Closing a channel twice panics: `panic: close of closed channel`.

The rule that avoids both: **the sender closes the channel, never the receiver**, and closes it only once,
normally from a dedicated goroutine that knows when there are going to be no more sends (you are going to
see it in section 6.7, with `wg.Wait()` followed by `close`).

**Deadlock, if nobody on the other side is listening.** I caused it with the shortest possible program
(the full example is in `programas/revisor/ejemplos/06-deadlock/main.go`):

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

**This corrects something many explanations take for granted: a deadlock in Go doesn't always "hang
forever".** Go's own runtime has a deadlock detector: if at some point **all** the program's goroutines are
asleep waiting for something that is never going to happen, Go notices and kills the program with
`fatal error: all goroutines are asleep - deadlock!`, with the exact line where it got stuck
(`[chan send]`, in this case, because it was sending). If your program seems "hung forever" instead of
ending with this error, it is almost certainly **not** a true deadlock: it is more likely that you have at
least one live goroutine doing something else (for example, a timer or an HTTP server listening), and that
one keeps the runtime from declaring that "all" of them are asleep.

### 6.6 `context`: how to cancel and set a time limit

<!-- verificar:fragmento -->
```go
ctx, cancelar := context.WithTimeout(context.Background(), 5*time.Second)
defer cancelar() // SIEMPRE, incluso si terminas antes de que se cumplan los 5 segundos

req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
resp, err := http.DefaultClient.Do(req) // se aborta solo si pasan los 5 s
```

`context` is Go's standard mechanism for two things: **putting a time limit** on an operation that could
take too long, and **propagating a cancellation** downwards (if the one above is canceled, everything that
depends on it is canceled too, without each function having to reinvent its own mechanism).

⚠️ **`defer cancelar()` is not optional, even when the operation has already finished on its own.**
`context.WithTimeout` starts an internal timer; if you never call the cancel function, that timer stays
alive until the original deadline is reached, holding on to memory the whole time — and if your program
creates contexts like this all the time (one for each service queried, like the `revisor`), without
`defer cancelar()` you accumulate memory leaks proportional to how many queries you made. It is the most
common leak in Go programs that use `context`, and that is why it is worth writing the `defer` on the same
line where you create the context, before writing anything else.

### 6.7 `Todos`: the function that brings the three pieces together

This is how the `revisor`'s central function really ended up (`programas/revisor/internal/revisar/todos.go`),
after bringing together goroutines, channels, semaphore, and `context`:

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

Read it with the previous sections fresh, because each piece answers a problem you already saw:

- **One goroutine per service** (6.1), counted with **`wg.Add(1)`/`defer wg.Done()`** (6.2) so the program
  knows when all of them finished.
- **A buffered `resultados` channel** (6.3) collects the states: nobody writes to a shared slice or map,
  so no lock is needed for that part.
- **`semaforo := make(chan struct{}, paralelo)`** is a channel used as a quota: it has room for
  `paralelo` values, so the `paralelo + 1`-th goroutine that tries to write to it (`semaforo <-
  struct{}{}`) blocks until another one frees its place (`<-semaforo`, in the `defer`). It is the same
  buffered channel from section 6.3, used not to carry data but to keep count of how many "turns" are
  left.
- **A `context.WithTimeout` of its own per service** (6.6), a child of the general `ctx`: if the one above
  is canceled (for example, if the report's total time limit expires), all the children are canceled with
  it.
- **`close(resultados)` from ANOTHER goroutine, after `wg.Wait()`** — and this deserves an explanation,
  because it is the part that is least visible to the naked eye: if we closed the channel in the same
  goroutine that does the `for e := range resultados` below, we would get stuck, because `Wait()` needs all
  the goroutines to read from the channel to free up space and be able to give back their turn, but the reading loop would never start because we would be waiting for `Wait()` first. By
  launching it separately, the closing and the reading happen **at the same time**, not one after the
  other.

### 6.8 Testing it without a network: `Falso` and the measured parallelism limit

Testing `Todos` against real services would be slow and nondeterministic. The `revisor` uses a fake
`Revisor` (`programas/revisor/internal/revisar/falso.go`) that simulates responses, delays, and even
services that never answer, all in memory:

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

🔒 **Here a mutex is needed, and it is the exception to the "prefer channels" rule from section 6.3.**
`Contador` is a shared integer that **many goroutines increment at the same time** during the tests
(`Falso.Revisar` is called concurrently, once for each service `Todos` is querying). A channel would serve
to *send* results, but for a simple shared counter, a `sync.Mutex` around the only line that touches it is
simpler and clearer. The rule is not "never use a mutex": it is "before using one, ask yourself whether a
channel better expresses what you are doing" — and for sending results, almost always yes; for a shared
counter, the mutex is almost always the right tool.

With `Falso`, this test measures something that would otherwise be almost impossible to check with
confidence: that the semaphore from section 6.7 really limits how many queries run at the same time, not
just that it "works in general" (complete, unabridged version in
`programas/revisor/internal/revisar/todos_test.go`):

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

Real run, with the race detector active (to confirm that not even the measurement itself introduces a
race):

```
$ go test ./internal/revisar/... -race -v -run TestTodos_elParaleloLimitaCuantasCorrenALaVez
=== RUN   TestTodos_elParaleloLimitaCuantasCorrenALaVez
--- PASS: TestTodos_elParaleloLimitaCuantasCorrenALaVez (0.11s)
PASS
ok  	github.com/habil/revisor/internal/revisar	1.435s
```

**20 services, a limit of 3, and the test confirms that there were never more than 3 calls to `Revisar`
running at the same time** — not because we assume it from the code, but because an atomic counter
measured it while it ran.

### 6.9 Time, measured: in series versus in parallel

With `Falso` configured to simulate real delays, the `TestTodos_respetaElTimeoutPorServicio` test confirms
the other side of the coin: a service that never answers (`Colgado: true`) with a
`Timeout: 30 * time.Millisecond` can't make `Todos` take longer than that:

```
$ go test ./internal/revisar/... -run TestTodos_respetaElTimeoutPorServicio -v
=== RUN   TestTodos_respetaElTimeoutPorServicio
--- PASS: TestTodos_respetaElTimeoutPorServicio (0.03s)
```

**0.03 seconds, no more.** The `context.WithTimeout` from section 6.7, a child of the general one, cut off
that query exactly when it should have, without the rest of the program having to find out that it
happened.

---

## The error you will see

All of these are literal, caused on purpose for this lesson:

| The symptom | Literal message / evidence | What happens and what to do |
|---|---|---|
| `main` ends before the goroutines run | (no output, or partial output that is inconsistent between runs) | A `sync.WaitGroup` (or a channel) that makes `main` wait is missing. Section 6.1 |
| `Add`/`Done` don't add up | The program prints that **everything went well first**, and `panic: sync: negative WaitGroup counter` arrives AFTERWARDS | A `wg.Add(1)` is missing before launching some goroutine, or there is an extra `Done()`. `Wait()` returned without waiting for anything. Section 6.2 |
| Shared data without protection | `WARNING: DATA RACE` + `Found 2 data race(s)` + `exit status 66` (with `-race`); an incorrect number with no explanation, without `-race` | Two goroutines touch the same memory without synchronization. Use a channel to send the result, or a mutex if you really need a shared counter. Sections 6.3 and 6.4 |
| Sending to a closed channel | `panic: send on closed channel` | Someone else already closed the channel, or it was closed by whoever shouldn't have. Always close from the sender, only once |
| Closing the same channel twice | `panic: close of closed channel` | Two goroutines (or two code paths) try to close the same channel. Centralize the closing in a single place |
| Nobody on the other side of an unbuffered channel | `fatal error: all goroutines are asleep - deadlock!` | Go's runtime detects that the whole program is asleep waiting for something that is never going to happen, and kills it with this exact line. Section 6.5 |
| You forgot `defer cancelar()` | (no immediate error; memory leak accumulated over time) | Each `context.WithTimeout` without canceling keeps its timer alive until it expires on its own. Section 6.6 |

**And the instruction that sums up this entire lesson:** if you are going to write concurrency, run
`-race` from your first test, not when something smells wrong — because, as you saw in section 6.3,
nothing smells wrong until it is too late.

---

## What goes wrong

- **Launching one goroutine per element, without any limit, against an external resource.** With 5
  services you don't notice it. With 5,000, the program opens 5,000 connections at once, and the
  bottleneck stops being the remote service and becomes your own machine (or the network, or the server
  itself, which gets showered with 5,000 simultaneous requests). The semaphore from section 6.7 exists
  exactly for this — never let "how many goroutines I launch" depend only on "how many elements I have".
- **Ignoring the `context` you are given, or not propagating it.** If a function receives a `ctx` and
  calls another operation that can take a while without passing it along, that operation is not going to
  be canceled when the original `ctx` is — you have two clocks that don't talk to each other. Everything
  that can take a while receives the `ctx` of whoever called it.
- **Not closing what you open.** A `context.WithTimeout` without its `cancelar()`, an HTTP connection
  without `resp.Body.Close()` (lesson 7), an unclosed file: each one is a different leak, but the way to
  avoid them is the same — a `defer` right after opening, before writing any other line.
- **Sharing memory instead of communicating it "because it is quicker to write".** A shared map with a
  mutex around the whole block may seem shorter than setting up a channel, but it is easier to forget a
  `Lock()` at a single access point than to forget to send through a channel — and the first slip doesn't
  show until the race detector (or, worse, production) finds it.
- **Trusting that "it has never failed" means "it is fine".** A data race can pass ninety-nine runs and
  fail the hundredth, or only fail with more cores than your laptop has. The only proof that a concurrent
  program has no races is running it with `-race`, not watching it pass several times without it.

---

## Exercises

1. Write the program from section 6.1 (goroutines without waiting) and run it five times in a row. Write
   down how many times it printed something and how many it didn't.
2. Add a correct `sync.WaitGroup` to it (section 6.2) and confirm that now the three lines are always
   printed, in all five runs.
3. Reproduce the missing-`Add` bug from section 6.2 and paste the complete `panic` in your logbook.
4. Reproduce the data race from section 6.3 (a shared counter without protection) and run the same
   program with and without `-race`. Compare the two final numbers and paste the complete `-race` output
   in your logbook.
5. Reproduce the deadlock from section 6.5 with an unbuffered channel and no receiver. Confirm that you
   see the `fatal error: all goroutines are asleep - deadlock!`, not a silent hang.
6. Take your own `revisor` (or this course's) and run `TestTodos_elParaleloLimitaCuantasCorrenALaVez`
   changing the `limite` to 1 and then to 10. Explain, in your own words, why the result of the test
   doesn't change (it always passes) but the **time** it takes does.
7. (A bit harder) Remove the semaphore from `Todos` (let all the goroutines be launched without a limit)
   and run `TestTodos_elParaleloLimitaCuantasCorrenALaVez` again. Confirm that it now fails, and paste the
   exact error message that `t.Errorf` gives with the maximum that was actually observed.

### Solutions

1-2. There is no single number: it depends on your machine. What matters is the comparison — without
`WaitGroup`, inconsistent; with it, always all three lines.

3. The program prints `listo` first — `Wait()` returned immediately because the counter was never
   incremented — and the `panic: sync: negative WaitGroup counter` arrives afterwards, with a trace that
   includes `sync.(*WaitGroup).Add(...)`. Run several times in a row, the order is always the same:
   `listo`, then the panic.

4. Without `-race`, a number lower than expected (for example, 956 out of 1000), without any error. With
   `-race`, the `WARNING: DATA RACE` block with the exact line of code, plus `Found N data race(s)` and
   `exit status 66`.

5. `fatal error: all goroutines are asleep - deadlock!`, with `goroutine 1 [chan send]:` (or `[chan
   receive]`, depending on which side got stuck) and the exact line of the channel.

6. The test passes in both cases because it **measures** the real maximum and compares it against the
   limit it configured itself (1 or 10) — never against a fixed expected number. The total time does
   change: with `limite=1` the 20 queries are strictly in series (one waits for the other), with
   `limite=10` they run in two batches of 10 instead of twenty of one.

7. Without a semaphore, all 20 goroutines run at the same time, so `maximoObservado` is going to be close
   to 20 (not exactly the test's limit, which still asks for 3). The `t.Errorf` says something like
   `se observaron 20 consultas simultáneas; el límite era 3` — the test DOES detect the regression, which
   is precisely what it exists for.

---

## How I know I got it

- [ ] I saw, with my own eyes, a program end without having waited for its goroutines (section 6.1).
- [ ] I caused the `panic: sync: negative WaitGroup counter` and saw that the program printed `listo`
      BEFORE the panic — and I can explain why.
- [ ] I caused a real data race and saw the difference between running it with and without `-race`.
- [ ] I can read a `WARNING: DATA RACE` block and say which line, which goroutines, and which variable
      are in conflict.
- [ ] I caused the `fatal error: all goroutines are asleep - deadlock!` and I know why Go detects it
      instead of hanging forever.
- [ ] I can explain, with the real code of `Todos`, what each of its four pieces is for: goroutines,
      results channel, semaphore, and per-service `context`.
- [ ] I measured, with a real test (not from memory), that the `revisor`'s parallelism limit is
      respected.
- [ ] I can explain why `Falso.Contador` uses a mutex instead of a channel, and why that doesn't
      contradict the general rule from section 6.3.

---

## Further reading

1. [A Tour of Go: Concurrency](https://go.dev/tour/concurrency/1) — the official, interactive tour of
   goroutines, channels, and `sync.WaitGroup`.
2. [Go Blog: Share Memory By Communicating](https://go.dev/blog/codelab-share) — the original article
   that explains the motto quoted in section 6.3.
3. [Package context](https://pkg.go.dev/context) — the official documentation, with the four canonical
   use cases (`WithCancel`, `WithTimeout`, `WithDeadline`, `WithValue`).
4. [Data Race Detector](https://go.dev/doc/articles/race_detector) — the official documentation of
   `-race`: what it detects, what it does NOT detect (for example, it doesn't find deadlocks, only data
   races), and its cost in execution time.

### Terms from this lesson

| Term | What it means |
|---|---|
| goroutine | a function that runs concurrently with the rest of the program, much cheaper than an operating system thread |
| `sync.WaitGroup` | mechanism to wait for a group of goroutines to finish, counting `Add`/`Done` |
| data race | two goroutines access the same memory at the same time, at least one of them writing, without synchronization |
| channel (`chan`) | Go's mechanism for one goroutine to send data to another without shared memory |
| channel semaphore | a buffered channel used as a quota of available turns, not to carry data |
| `context` | standard mechanism to propagate time limits and cancellation between functions |
| deadlock | state in which all of a program's goroutines are asleep waiting for something that is never going to happen; Go's runtime detects it and ends the program |

---

**Previous:** [Lesson 5 — Modules and tests](05-modulos-y-pruebas.md) ·
**Next:** [Lesson 7 — The finished program](07-el-programa.md)
