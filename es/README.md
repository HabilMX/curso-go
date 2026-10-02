# Curso de Go — de cero a escribir un buen programa

**Por Dorian Chávez, fundador de Hábil y arquitecto de integración.**

**Para quién es:** alguien que va empezando — primer año de universidad, o recién salido de la
escuela. **No damos por hecho nada**: ni que tienes Go instalado, ni que sabes qué es un puntero, ni qué
significa compilar. Cada concepto se explica cuando aparece, y se explica **por qué** existe, no solo
cómo se escribe.

**Qué vas a terminar sabiendo:** escribir un programa completo, entender lo que escribiste, y poder
explicárselo a alguien más. Eso último es la verdadera prueba.

**Lo que necesitas antes de empezar:** una computadora con Linux Mint y saber abrir una terminal. Nada
más. La [Lección 1](01-instalacion.md) instala todo desde cero.

## El proyecto que vas a construir

Un programa de línea de comandos llamado **`revisor`**: recibe una lista de servicios, los consulta
**todos a la vez**, y produce un reporte en tabla o en JSON.

    $ revisor --config servicios.txt --formato tabla
    SERVICIO    ESTADO    TIEMPO  DETALLE
    catalogo    OK         142ms  200
    inventario  LENTO       2.3s  200
    pagos       OK          87ms  200
    reportes    FALLA          —  connection refused

Parece pequeño y no lo es: para escribirlo bien necesitas structs, interfaces, errores, concurrencia,
`context`, HTTP, JSON, banderas de línea de comandos, pruebas y un binario empaquetado.

**Y es el esqueleto de lo que la gente escribe en Go de verdad.** Según la
[Go Developer Survey 2025](https://go.dev/blog/survey2025) —la encuesta oficial del equipo de Go,
levantada en septiembre de 2025 con **5,379 respuestas**—, los dos tipos de proyecto más construidos son:

| Qué construyen | |
|---|---|
| Herramientas de línea de comandos | **74 %** |
| Servicios API/RPC | **73 %** |
| Bibliotecas o marcos de trabajo | 49 % |

El `revisor` es lo primero, y en la Lección 7 le agregas lo segundo.

⚠️ **Y lo vas a escribir otra vez en Rust**, en el curso hermano. Escribir el mismo programa en los dos
lenguajes es lo que de verdad enseña en qué se diferencian; leer comparaciones no sirve.

## Las ocho lecciones

| Lección | | Qué construyes | Qué aprendes |
|---|---|---|---|
| 0 | [Qué es Go](00-introduccion.md) | nada todavía (lectura) | de dónde salió Go, para qué sirve y para qué no |
| 1 | [Instalación](01-instalacion.md) | el entorno y tu primer programa | `GOROOT`, `GOPATH`, el `PATH`, y los errores reales de instalar Go |
| 2 | [Variables, funciones y tipos](02-fundamentos.md) | las funciones base del `revisor` | el compilador, variables, tipos básicos, funciones con `(resultado, error)` |
| 3 | [Structs, métodos, errores e interfaces](03-errores-interfaces.md) | el `revisor` con structs e interfaces | structs, métodos, punteros, errores como valores, interfaces |
| 4 | [Colecciones: slices y maps](04-colecciones.md) | la lista de servicios y el reporte | slices, maps, `range`, orden estable |
| 5 | [Módulos y pruebas](05-modulos-y-pruebas.md) | el proyecto de verdad, con tests | `go mod`, tabla de casos, `-race`, cobertura |
| 6 | [Concurrencia](06-concurrencia.md) | **que revise todo a la vez** | goroutines, canales, `context`, `WaitGroup`, semáforo |
| 7 | [El programa terminado](07-el-programa.md) | CLI, HTTP, JSON y binario | `net/http`, `flag`, `encoding/json`, compilación cruzada |

## El orden no es el obvio, y es a propósito

Structs, **errores e interfaces van antes** que slices, maps y punteros. Suena raro —casi todos los
cursos ponen las colecciones primero— pero no es un capricho: es el orden que usan dos de los cursos de Go
con más alumnos, el de [Boot.dev / freeCodeCamp](https://www.boot.dev/courses/learn-golang) (structs →
interfaces → errores en los puestos 5, 6 y 7, **antes** de slices, maps y punteros) y el de
[Todd McLeod en Udemy](https://www.udemy.com/course/learn-how-to-code/).

La razón es buena: **en Go, modelar datos y manejar errores ES el lenguaje.** Un programa con slices
perfectos y errores ignorados no es Go; es C con otra sintaxis.

**La concurrencia va en la Lección 6, no en la 3.** Es la razón por la que Go existe, pero necesitas
structs, errores e interfaces para que los ejemplos no sean de juguete. Y se da en tres pasadas
—goroutines, luego canales, luego aplicada al proyecto— porque en una sola no cuaja.

## Cómo usarlo

- **Una sesión de 90 minutos por semana**, o dos de 45. Menos no cuaja; más de dos horas seguidas se
  olvida.
- 🔴 **Escribe el código a mano, siempre.** Copiar y pegar produce archivos, no conocimiento. Tus dedos
  aprenden cosas que tus ojos no.
- **Rompe el código a propósito.** En varias prácticas te vamos a pedir que provoques un error y leas lo
  que dice el compilador. **No es relleno:** aprender a leer errores es la mitad de saber programar.
- **No pases de lección con dudas abiertas.** Anótalas en la bitácora y pregúntalas. Todo lo que viene se
  apoya en lo anterior.
- **Cada lección cierra con un «cómo sé que lo logré» medible**: compila y pasa, o no.
- **[`bitacora.md`](bitacora.md) es tuyo**: dudas, tropiezos y lo que te sorprendió. Es lo que vuelve
  esto tu curso.

## Las tres fuentes que valen, en este orden

1. **[A Tour of Go](https://go.dev/tour/)** — oficial, interactivo, 2-3 h. **Completo, no hojeado.**
2. **[Go by Example](https://gobyexample.com/)** — cada concepto con código mínimo. Como referencia.
3. **[Effective Go](https://go.dev/doc/effective_go)** — la filosofía. **Léelo en la Lección 3, no al final.**

⚠️ **Aprende la biblioteca estándar antes de cualquier framework.** En Go el estándar alcanza para casi
todo, y eso es una decisión de diseño, no una carencia. Quien empieza con Gin o Echo aprende el
framework y no el lenguaje.

---

## Quién escribe esto

Este curso lo escribe **Dorian Chávez**, CEO y consultor principal de **Hábil**, firma mexicana de
ingeniería de software. Ingeniero en Sistemas Computacionales por el IPN (ESCOM), con 26 años construyendo
y operando sistemas en producción.

Lo publicamos abierto, con licencia [Creative Commons](../LICENSE.md), porque quien explica bien construye
bien — y porque un curso que no se puede copiar, traducir ni mejorar sirve de poco.

**Si encuentras un error, una explicación confusa o una cifra que no cuadra, dilo.** Este material está
hecho para corregirse. Antes de publicar cada versión, comprobamos que todo el código del curso funciona
y que lo que muestra cada lección es idéntico al programa real, para que ninguna lección enseñe algo que
ya cambió.

### Si esto te resultó útil

- **¿Quieres trabajar con gente que escribe así?** Escríbenos.
- **¿Tu equipo está adoptando Go o Rust?** También.

*(Los enlaces de contacto los pone el sitio al publicar esta página.)*
