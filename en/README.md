# Go course — from zero to writing a good program

**By Dorian Chávez, founder of Hábil and integration architect.**

**Who it is for:** someone who is just starting out — first year of university, or fresh out of school.
**We take nothing for granted**: not that you have Go installed, not that you know what a pointer is, not
what compiling means. Each concept is explained when it appears, and we explain **why** it exists, not
just how it is written.

**What you will end up knowing:** how to write a complete program, understand what you wrote, and be able
to explain it to someone else. That last part is the real test.

**What you need before starting:** a computer with Linux Mint and knowing how to open a terminal. Nothing
else. [Lesson 1](01-instalacion.md) installs everything from scratch.

## The project you are going to build

A command-line program called **`revisor`**: it receives a list of services, queries them **all at
once**, and produces a report as a table or in JSON.

    $ revisor --config servicios.txt --formato tabla
    SERVICIO    ESTADO    TIEMPO  DETALLE
    catalogo    OK         142ms  200
    inventario  LENTO       2.3s  200
    pagos       OK          87ms  200
    reportes    FALLA          —  connection refused

It looks small and it isn't: to write it well you need structs, interfaces, errors, concurrency,
`context`, HTTP, JSON, command-line flags, tests, and a packaged binary.

**And it is the skeleton of what people really write in Go.** According to the
[Go Developer Survey 2025](https://go.dev/blog/survey2025) —the Go team's official survey, conducted in
September 2025 with **5,379 responses**—, the two most-built kinds of project are:

| What they build | |
|---|---|
| Command-line tools | **74 %** |
| API/RPC services | **73 %** |
| Libraries or frameworks | 49 % |

The `revisor` is the first, and in Lesson 7 you add the second to it.

⚠️ **And you are going to write it again in Rust**, in the sibling course. Writing the same program in
both languages is what really teaches how they differ; reading comparisons doesn't help.

## The eight lessons

| Lesson | | What you build | What you learn |
|---|---|---|---|
| 0 | [What Go is](00-introduccion.md) | nothing yet (reading) | where Go came from, what it is for and what it isn't for |
| 1 | [Installation](01-instalacion.md) | the environment and your first program | `GOROOT`, `GOPATH`, the `PATH`, and the real errors of installing Go |
| 2 | [Variables, functions, and types](02-fundamentos.md) | the base functions of the `revisor` | the compiler, variables, basic types, functions with `(resultado, error)` |
| 3 | [Structs, methods, errors, and interfaces](03-errores-interfaces.md) | the `revisor` with structs and interfaces | structs, methods, pointers, errors as values, interfaces |
| 4 | [Collections: slices and maps](04-colecciones.md) | the list of services and the report | slices, maps, `range`, stable order |
| 5 | [Modules and tests](05-modulos-y-pruebas.md) | the real project, with tests | `go mod`, table of cases, `-race`, coverage |
| 6 | [Concurrency](06-concurrencia.md) | **making it check everything at once** | goroutines, channels, `context`, `WaitGroup`, semaphore |
| 7 | [The finished program](07-el-programa.md) | CLI, HTTP, JSON, and binary | `net/http`, `flag`, `encoding/json`, cross-compilation |

## The order is not the obvious one, and it is on purpose

Structs, **errors, and interfaces come before** slices, maps, and pointers. It sounds odd —almost all
courses put collections first— but it is not a whim: it is the order used by two of the Go courses with
the most learners, the one from [Boot.dev / freeCodeCamp](https://www.boot.dev/courses/learn-golang)
(structs → interfaces → errors in positions 5, 6, and 7, **before** slices, maps, and pointers) and the
one from [Todd McLeod on Udemy](https://www.udemy.com/course/learn-how-to-code/).

The reason is a good one: **in Go, modeling data and handling errors IS the language.** A program with
perfect slices and ignored errors is not Go; it is C with a different syntax.

**Concurrency goes in Lesson 6, not in 3.** It is the reason Go exists, but you need structs, errors, and
interfaces so that the examples aren't toys. And it is taught in three passes —goroutines, then channels,
then applied to the project— because in a single one it doesn't sink in.

## How to use it

- **One 90-minute session per week**, or two of 45. Less doesn't sink in; more than two hours in a row
  gets forgotten.
- 🔴 **Write the code by hand, always.** Copying and pasting produces files, not knowledge. Your fingers
  learn things your eyes don't.
- **Break the code on purpose.** In several exercises we are going to ask you to cause an error and read
  what the compiler says. **It is not filler:** learning to read errors is half of knowing how to program.
- **Don't move on to the next lesson with open questions.** Write them down in the logbook and ask them.
  Everything that comes next rests on what came before.
- **Each lesson closes with a measurable "how I know I got it"**: it compiles and passes, or it doesn't.
- **[`bitacora.md`](bitacora.md) is yours**: questions, stumbles, and what surprised you. It is what makes
  this your course.

## The three sources worth it, in this order

1. **[A Tour of Go](https://go.dev/tour/)** — official, interactive, 2-3 h. **Completed, not skimmed.**
2. **[Go by Example](https://gobyexample.com/)** — each concept with minimal code. As a reference.
3. **[Effective Go](https://go.dev/doc/effective_go)** — the philosophy. **Read it in Lesson 3, not at the
   end.**

⚠️ **Learn the standard library before any framework.** In Go the standard library is enough for almost
everything, and that is a design decision, not a shortcoming. Whoever starts with Gin or Echo learns the
framework and not the language.

---

## Who writes this

This course is written by **Dorian Chávez**, CEO and principal consultant of **Hábil**, a Mexican software
engineering firm. Computer Systems Engineer from the IPN (ESCOM), with 26 years building and operating
systems in production.

We publish it openly, under a [Creative Commons](../LICENSE.md) license, because whoever explains well
builds well — and because a course that can't be copied, translated, or improved is of little use.

**If you find an error, a confusing explanation, or a figure that doesn't add up, say so.** This material
is made to be corrected. Before publishing each version, we check that all the course's code works and
that what each lesson shows is identical to the real program, so that no lesson teaches something that has
already changed.

### If this was useful to you

- **Do you want to work with people who write like this?** Write to us.
- **Is your team adopting Go or Rust?** Write to us too.

*(The contact links are added by the site when this page is published.)*
