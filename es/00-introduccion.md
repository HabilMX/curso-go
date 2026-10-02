# Lección 0 — Qué es Go, de dónde salió y por qué deberías aprenderlo

**Duración:** 30-45 minutos de lectura, sin escribir código todavía.

**Al terminar vas a poder:**

- Contar quién creó Go, cuándo y qué problema concreto estaban tratando de resolver.
- Explicar por qué existe un lenguaje nuevo si ya había docenas.
- Nombrar tres programas que usas o administras que están escritos en Go.
- Decir en qué es bueno Go y en qué **no**, para no usarlo donde no conviene.
- Situar a Go entre los demás lenguajes que ya conoces o has oído.

---

## Por qué importa

Es 2007. Google tiene una de las bases de código más grandes del mundo, escrita mayormente en C++. Y
tiene un problema muy concreto y muy aburrido: **compilar tarda una eternidad.**

Un solo cambio en un archivo de encabezado podía disparar una recompilación de **45 minutos o más**.

Piensa en lo que eso significa para quien programa. Cambias una línea. Esperas 45 minutos. Descubres que
te equivocaste en una coma. Cambias otra línea. Esperas otros 45 minutos. **En un día de trabajo alcanzas
para intentar ocho cosas.**

El 21 de septiembre de 2007, tres ingenieros de Google se pararon frente a un pizarrón a diseñar un
lenguaje que no los hiciera esperar. Ese lenguaje es del que trata todo este curso, y el resto de esta
lección explica cómo llegaron de ese pizarrón al lenguaje que vas a instalar en la lección 1.

---

## Los conceptos

### 0.1 Quiénes lo hicieron (y por qué importa)

No eran tres programadores cualquiera:

| | Quién es |
|---|---|
| **Ken Thompson** | **Creó Unix.** Y creó el lenguaje B, el antecesor directo de C, que escribió junto con Dennis Ritchie. Premio Turing —el equivalente al Nobel en computación— en 1983 |
| **Rob Pike** | Trabajó en Unix y creó Plan 9, su sucesor. **Coinventó UTF-8**, la forma en que hoy se guarda texto en prácticamente todo el mundo, incluido este documento |
| **Robert Griesemer** | Trabajó en el motor JavaScript V8 (el de Chrome y Node.js) y en la máquina virtual de Java |

🔑 **Lee eso otra vez: uno de los creadores de Go es el creador de Unix y coautor de C.** Cincuenta años
después de haber inventado las herramientas sobre las que está construido todo el software moderno, la
misma persona diseñó Go.

Eso explica mucho del carácter del lenguaje: Go se siente como C —directo, pequeño, sin adornos— pero sin
las trampas que hacían de C un campo minado.

### 0.2 Los tres requisitos que nadie cumplía

Los tres querían un lenguaje con **tres** propiedades a la vez:

1. **Que compile rápido.** El problema original.
2. **Que ejecute rápido.** Google corre software en miles de máquinas: la velocidad cuesta dinero real.
3. **Que sea fácil de programar.** Que una persona nueva en el equipo pueda leer código ajeno y
   entenderlo el primer día.

Revisaron los lenguajes que ya existían. Y ahí está el hallazgo interesante: **encontraron lenguajes que
cumplían dos de los tres, pero ninguno que cumpliera los tres.**

| | Compila rápido | Ejecuta rápido | Fácil de programar |
|---|---|---|---|
| **C / C++** | ❌ | ✅ | ❌ |
| **Java** | 🟡 | ✅ | 🟡 |
| **Python** | ✅ (no compila) | ❌ | ✅ |
| **Go** | ✅ | ✅ | ✅ |

> [!NOTE]
> 🔧 **Observación de ingeniería de software 0.1**
> Un lenguaje de programación no es una cuestión de gusto: es una herramienta con compromisos. Cada uno
> sacrificó algo para ganar otra cosa. Entender **qué sacrificó** el que estás usando es la diferencia
> entre programar con criterio y programar de memoria.

### 0.3 Cómo se ganó la simplicidad: quitando cosas

Aquí está la decisión más contraintuitiva de Go, y la que más lo distingue.

Casi todos los lenguajes crecen: cada versión agrega características. C++ acumuló décadas de ellas, y el
resultado es tan grande que **nadie lo conoce completo**. Hay gente que programa C++ veinte años y sigue
encontrando rincones que no conocía.

**Go hizo lo contrario: decidió qué dejar afuera.** No tiene:

- **Clases ni herencia.** La forma de reutilizar código es otra, y la vas a ver en la lección 3.
- **Excepciones** (`try / catch`). Los errores se manejan de otra manera, y la verás en la lección 2.
- **Sobrecarga de operadores.** `+` significa sumar, siempre, y no se puede cambiar.
- **Constructores ni destructores.**
- **Un montón de formas de hacer lo mismo.** Casi siempre hay **una** manera idiomática.

⚠️ **Esto va a molestarte alguna vez.** Vas a querer hacer algo que en otro lenguaje se hace en una línea
y en Go te va a tomar cinco. La compensación es enorme y se ve a los meses: **puedes leer código de Go
escrito por cualquier persona y entenderlo.** No hay dialectos, no hay trucos, no hay que aprender cómo
programa cada equipo.

> 🔑 **La cita que resume la filosofía**, de Rob Pike:
> *«La claridad es mejor que la astucia.»*

### 0.4 La otra razón: las computadoras dejaron de acelerar

Hay un segundo motivo detrás de Go, y es de hardware.

Hasta cerca de 2005, cada año salían procesadores más rápidos y tu programa corría más rápido **sin que
tocaras nada**. Eso se terminó: los procesadores dejaron de subir de velocidad y empezaron a multiplicar
**núcleos**. Tu laptop no tiene un procesador rapidísimo: tiene 4, 8 o 16 procesadores modestos.

**El problema:** un programa normal usa **un** núcleo. Los otros quince están ahí, mirando.

Para aprovecharlos hay que escribir programas que hagan varias cosas a la vez, y eso —la
**concurrencia**— era tradicionalmente una de las tareas más difíciles y propensas a errores de toda la
programación.

**Go se diseñó alrededor de eso.** Hacer que dos cosas ocurran a la vez en Go es literalmente escribir una
palabra:

```go
go revisarServicio()
```

Esa palabra `go` —la que le da el nombre al lenguaje— es la característica que lo hizo famoso. Lo vas a
aprender en la lección 6, y cuando llegues ahí vas a entender por qué tanta gente se cambió a este
lenguaje.

### 0.5 La línea de tiempo

| Fecha | Qué pasó |
|---|---|
| **21-sep-2007** | Griesemer, Pike y Thompson diseñan Go en un pizarrón |
| mediados de 2008 | ya hay un compilador que funciona |
| **10-nov-2009** | Google lo anuncia públicamente y lo libera como software libre |
| **28-mar-2012** | **Go 1.0**, la primera versión estable |
| 2012 en adelante | dos versiones grandes al año, cada febrero y agosto |
| hoy | Go 1.27 |

🔑 **Y algo que vale saber: la promesa de compatibilidad de Go 1.** Desde 2012, el equipo de Go prometió
que **un programa escrito para Go 1.0 sigue compilando hoy**, catorce años después. Han cumplido.

Eso no es normal. En muchos lenguajes, código de hace cinco años ya no compila. **En Go, lo que aprendas
hoy te va a servir dentro de diez años**, y eso hace que valga mucho más la pena el tiempo que vas a
invertir.

### 0.6 Qué está escrito en Go (probablemente ya lo usas)

Esto no es una lista de propaganda: es para que veas **dónde** se usa este lenguaje, porque dice mucho
sobre para qué sirve.

| | Qué es |
|---|---|
| **Docker** | la herramienta que empaqueta aplicaciones en contenedores |
| **Kubernetes** | el sistema que administra miles de contenedores en servidores. Lo usa media industria |
| **Terraform** | crea infraestructura en la nube escribiendo archivos |
| **Prometheus** · **Grafana** | recolectan y grafican métricas de servidores |
| **Traefik** · **Caddy** | servidores web y balanceadores de carga |
| **CockroachDB** · **InfluxDB** | bases de datos |
| **Hugo** | generador de sitios web estáticos, famoso por su velocidad |
| **ngrok**, **rclone**, **gh** (el CLI de GitHub) | herramientas de línea de comandos |

🔑 **¿Ves el patrón?** Casi todo son **herramientas de infraestructura**: cosas que corren en servidores,
que tienen que ser rápidas, que se despliegan como un archivo suelto y que manejan muchas conexiones a la
vez. **Ahí es donde Go gana**, y no es casualidad: es exactamente el problema que Google tenía.

### 0.7 Y entonces, ¿por qué deberías aprenderlo tú?

Cuatro razones honestas:

**1. Se aprende rápido.** La especificación de Go se lee en una tarde. La de C++ tiene más de 1,800
páginas. **Vas a poder escribir programas útiles en semanas, no en años**, y eso importa mucho cuando
estás empezando.

**2. Te enseña cosas que sirven en cualquier lenguaje.** Al ser compilado y de tipos estrictos, Go te
obliga a pensar en tipos, en memoria y en errores. Esos conceptos **se transfieren**: si después programas
Java, C# o Rust, ya los tienes.

**3. Hay trabajo.** Todo lo que corre en servidores modernos tiene Go adentro. Si te interesa
infraestructura, nube, DevOps o backend, es una de las apuestas más seguras.

**4. Lo que aprendas no va a caducar.** Por la promesa de compatibilidad de Go 1.

### 0.8 Qué vas a construir en este curso

Un programa de línea de comandos llamado **`revisor`**: recibe una lista de servicios, los consulta
**todos a la vez**, y produce un reporte.

```
$ revisor --config servicios.txt --formato tabla
SERVICIO    ESTADO    TIEMPO  DETALLE
catalogo    OK         142ms  200
inventario  LENTO       2.3s  200
pagos       OK          87ms  200
reportes    FALLA          —  connection refused
```

Parece pequeño. **No lo es.** Para escribirlo bien vas a necesitar todo: tipos, structs, errores,
interfaces, listas, concurrencia, HTTP, archivos de configuración, pruebas y compilar un ejecutable.

Y va a crecer contigo: **cada lección agrega una pieza.** Al final tendrás un programa que de verdad
podrías usar, no un ejercicio de libro.

---

## El error que vas a ver

Esta lección no tiene código propio todavía, así que no hay un mensaje de compilador que mostrarte. Pero
sí hay un error de **pensamiento** que casi todo el mundo comete al leer la sección 0.3, y vale la pena
adelantártelo para que no te frene cuando te pase.

**El error:** leer la lista de lo que Go **no** tiene —clases, herencia, excepciones, sobrecarga de
operadores— y concluir *«entonces es un lenguaje pobre, le faltan cosas que necesito».*

**Por qué es un error:** confunde *tener menos herramientas* con *poder resolver menos problemas*. Go no
te quitó la capacidad de reutilizar código ni de manejar errores: te dio **una** forma de hacerlo en vez de
diez, y esa forma la vas a aprender en las lecciones 2 y 3. La frustración es real y es normal —tú también
la vas a sentir la primera semana—, pero se resuelve programando, no evitando el lenguaje. Cuando termines
la lección 3 vas a poder volver a leer la sección 0.3 y ver por qué cada renglón de esa lista es una
decisión, no una carencia.

---

## Lo que se hace mal

Tan importante como saber para qué sirve una herramienta es saber para qué no. **Ningún lenguaje es bueno
para todo**, y quien te diga lo contrario te está vendiendo algo. Usar Go donde no conviene es el primer
antipatrón de este curso, antes de haber escrito una sola línea:

| No es la mejor opción para | Qué se usa en su lugar | Por qué |
|---|---|---|
| Apps de iPhone o Android | Swift, Kotlin | las plataformas están hechas para esos |
| Páginas web (lo que corre en el navegador) | JavaScript, TypeScript | el navegador solo ejecuta JavaScript |
| Ciencia de datos, inteligencia artificial | Python | todas las bibliotecas del mundo están ahí |
| Videojuegos grandes | C++, C# | necesitan control absoluto de la memoria y el hardware |
| Sistemas donde un microsegundo importa | C, C++, **Rust** | Go tiene un recolector de basura que pausa el programa a veces |

> [!NOTE]
> 🔧 **Observación de ingeniería de software 0.2**
> La última fila es la razón por la que este curso tiene un **segundo curso, de Rust**. Rust resuelve
> exactamente eso: velocidad de C sin recolector de basura y sin los errores de memoria de C. Son
> herramientas para problemas distintos, y vas a aprender las dos.

---

## Ejercicios

### Preguntas de repaso

**0.1** ¿Qué problema concreto motivó la creación de Go, y en qué año?

**0.2** Nombra a los tres creadores de Go y di por qué es relevante quién es Ken Thompson.

**0.3** ¿Cuáles eran los tres requisitos que buscaban y por qué no les servía ningún lenguaje existente?

**0.4** Menciona tres cosas que Go **no** tiene, a propósito. ¿Qué se gana al quitarlas?

**0.5** ¿Qué cambió en el hardware alrededor de 2005 y qué tiene que ver con Go?

**0.6** ¿Qué es la promesa de compatibilidad de Go 1 y por qué te conviene?

**0.7** Nombra tres programas escritos en Go y di qué tienen en común.

**0.8** Da dos casos en los que **no** usarías Go, y di qué usarías.

### Para ir más allá

**0.9** Busca en internet la charla o el artículo original *«Go at Google: Language Design in the Service
of Software Engineering»* de Rob Pike. Lee la introducción y anota **una** razón de diseño que no esté en
esta lección.

**0.10** Elige **dos** de los programas de la sección 0.6 que no conozcas. Averigua en una frase qué hace
cada uno y anótalo en tu bitácora.

**0.11** Busca el índice de la especificación del lenguaje Go (*The Go Programming Language
Specification*) y cuenta cuántas páginas o secciones tiene. Compáralo con el estándar de C++. Anota los
dos números: es la forma más concreta de ver qué significa «simple».

**0.12 (Para pensar, sin respuesta correcta)** Go quitó las excepciones, las clases y la herencia —cosas
que otros lenguajes consideran indispensables. ¿Se te ocurre alguna razón por la que **quitar** una
característica pueda hacer mejor a un lenguaje? Escribe tu opinión en la bitácora **antes** de empezar el
curso, y vuelve a leerla al terminar. Es interesante ver si cambió.

### Soluciones

**0.1** Los tiempos de compilación de C++ en Google: un cambio en un encabezado podía costar **45 minutos**
de recompilación. Empezaron el **21 de septiembre de 2007**.

**0.2** Robert Griesemer, Rob Pike y **Ken Thompson**. Thompson **creó Unix** y el lenguaje B, antecesor
de C, que escribió con Dennis Ritchie. Es decir: uno de los autores de las herramientas sobre las que se
construyó el software moderno diseñó también Go, cincuenta años después.

**0.3** Compilación rápida, ejecución rápida y facilidad de programar. Los lenguajes existentes cumplían
**dos de los tres**: C++ era rápido al ejecutar pero lento al compilar y difícil; Python era fácil pero
lento al ejecutar.

**0.4** Clases y herencia, excepciones, sobrecarga de operadores, constructores. **Se gana legibilidad**:
código de Go escrito por cualquiera se puede leer y entender, porque no hay dialectos ni muchas formas de
hacer lo mismo.

**0.5** Los procesadores dejaron de acelerar y empezaron a multiplicar **núcleos**. Un programa normal usa
uno solo, así que hacía falta un lenguaje donde escribir programas concurrentes fuera fácil. De ahí la
palabra `go`.

**0.6** La promesa de que un programa escrito para Go 1.0 (2012) **sigue compilando hoy**. Te conviene
porque lo que aprendas no caduca y el código que escribas va a seguir funcionando durante años.

**0.7** Docker, Kubernetes, Terraform, Prometheus, Grafana, Hugo… Todos son **herramientas de
infraestructura**: corren en servidores, se distribuyen como un ejecutable suelto y manejan muchas
conexiones a la vez.

**0.8** Apps móviles (Swift/Kotlin), código de navegador (JavaScript), ciencia de datos (Python),
videojuegos grandes (C++), sistemas de tiempo real estricto (C o Rust).

**0.9-0.12** No tienen una única respuesta correcta: son de investigación y reflexión propia. Compara lo
que anotaste con un compañero o en la bitácora del curso, no con una respuesta fija.

---

## Cómo sé que lo logré

No hay código que compilar en esta lección, así que el checklist es de comprensión. Marca cada punto solo
si puedes hacerlo **sin volver a ver el texto**:

- ☐ Explicas en tus palabras, en menos de un minuto, por qué nació Go y qué problema resolvía.
- ☐ Nombras a los tres creadores y dices por qué que Ken Thompson sea uno de ellos no es un dato de
  relleno.
- ☐ Dices qué son la compilación rápida, la ejecución rápida y la facilidad de programar, y por qué
  ningún lenguaje anterior a Go tenía las tres.
- ☐ Nombras, de memoria, al menos tres cosas que Go decidió no tener, y explicas qué se gana al quitarlas.
- ☐ Explicas la relación entre los núcleos de un procesador, la concurrencia y la palabra `go`.
- ☐ Nombras tres programas escritos en Go y dices qué tienen en común.
- ☐ Das un ejemplo de un problema para el que **no** usarías Go, y dices qué usarías en su lugar.
- ☐ Resolviste las ocho preguntas de repaso sin ver las soluciones antes de comparar.

Si te faltó alguno, no sigas a la lección 1 todavía: vuelve a leer la sección correspondiente. Todo lo que
viene se apoya en esto.

---

## Resumen

- Go nació en **2007** en Google por la frustración con los **tiempos de compilación de C++**: 45 minutos
  por un cambio en un encabezado.
- Lo crearon **Robert Griesemer, Rob Pike y Ken Thompson**; Thompson creó Unix y coescribió C.
- Buscaban **tres** propiedades juntas —compilar rápido, ejecutar rápido, ser fácil— que ningún lenguaje
  de entonces cumplía a la vez.
- Go logró la simplicidad **quitando** cosas: sin clases, sin herencia, sin excepciones, sin sobrecarga de
  operadores.
- La segunda motivación fue el hardware: los procesadores multiplicaron **núcleos** en vez de acelerar, y
  aprovecharlos exigía que la **concurrencia** fuera fácil. De ahí la palabra `go`.
- Fue anunciado en **2009** y la versión **1.0** salió en **2012**; hoy va por 1.27, con dos versiones
  grandes al año.
- La **promesa de compatibilidad de Go 1** garantiza que el código de 2012 sigue compilando: lo que
  aprendes no caduca.
- Go domina la **infraestructura**: Docker, Kubernetes, Terraform, Prometheus, Grafana.
- **No** es la mejor opción para móviles, navegador, ciencia de datos, videojuegos ni tiempo real
  estricto.

---

## Para leer más

1. **La especificación del lenguaje** — *The Go Programming Language Specification*, en `go.dev/ref/spec`.
   Es la fuente oficial y, como viste en el ejercicio 0.11, se lee en una tarde. No hace falta entenderla
   toda hoy; vale la pena saber que existe y es corta.
2. **El blog oficial de Go** — `go.dev/blog`, donde el propio equipo publica las novedades de cada versión
   y artículos de diseño.
3. **Rob Pike, «Go at Google: Language Design in the Service of Software Engineering»** — la charla/artículo
   original donde uno de los creadores explica las decisiones de diseño con el contexto de Google. Es la
   lectura del ejercicio 0.9, y la más directa para entender el **porqué** detrás de cada decisión que
   viste en esta lección.

### Términos de esta lección

| | |
|---|---|
| **compilar** | traducir código fuente a instrucciones de máquina |
| **concurrencia** | que un programa haga varias cosas a la vez |
| **encabezado (*header*)** | en C/C++, archivo con declaraciones que otros archivos incluyen |
| **Go 1 (promesa de compatibilidad)** | compromiso de que el código viejo siga compilando |
| **núcleo (*core*)** | cada procesador independiente dentro de un mismo chip |
| **Premio Turing** | el mayor reconocimiento en ciencias de la computación |
| **recolector de basura** | parte del lenguaje que libera memoria automáticamente |
| **UTF-8** | forma estándar de representar texto, coinventada por Rob Pike y Ken Thompson |

---

**Siguiente:** [Lección 1 — Instalar Go en tu Linux Mint](01-instalacion.md)
