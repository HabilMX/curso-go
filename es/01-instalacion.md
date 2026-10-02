# Lección 1 — Instalar Go en tu Linux Mint y tu primer programa

**Duración:** 45-60 minutos (o 2 sesiones de 30). Es la única lección que no enseña Go como lenguaje:
enseña la herramienta y el terreno donde vas a trabajar — y ese terreno tiene más trampas de las que
parece a simple vista.

**Al terminar vas a poder:**

- Instalar la versión oficial y vigente de Go en Linux Mint, sin depender del paquete del sistema, y
  explicar por qué ese paquete está obsoleto por diseño, no por descuido.
- Explicar qué son `GOROOT`, `GOPATH` y el `PATH`, y por qué Go no funciona hasta que tocas el segundo.
- Reconocer la diferencia entre una terminal que lee tu configuración y una que no —la causa real de
  «lo agregué y de todos modos no funciona»— y saber cuál es la tuya.
- Crear un proyecto con `go mod init`, escribir un programa y ejecutarlo con `go run`.
- Producir un binario con `go build` y explicar, con una prueba real, por qué corre en otra máquina sin
  tener Go instalado.
- Usar `go fmt` para no discutir nunca de estilo de código.
- Reconocer, palabra por palabra, los mensajes de error más comunes de esta etapa —de compilador, de
  permisos y de ruta— y saber qué hacer con cada uno.

---

## Por qué importa

Lo primero que un tutorial viejo te va a sugerir es `apt install golang`, y es un error. No es un error
menor de "versión un poco vieja": es la clase de error que te hace perder una tarde entera sin entender
por qué, semanas después de instalar, porque el síntoma no aparece al instalar — aparece cuando ya estás
programando y algo que "debería funcionar" no funciona.

Esto no se opina, se mide. El 30-sep-2026 medí, en un contenedor con la misma base que Linux Mint 22.3
(Ubuntu 24.04 "Noble Numbat", que es de donde Mint toma sus paquetes):

```
$ cat /etc/os-release | grep VERSION
VERSION="24.04.5 LTS (Noble Numbat)"

$ apt-cache policy golang-go
golang-go:
  Installed: (none)
  Candidate: 2:1.22~2build1
  Version table:
     2:1.22~2build1 500
        500 http://ports.ubuntu.com/ubuntu-ports noble/main arm64 Packages

$ curl -s 'https://go.dev/VERSION?m=text'
go1.27.1
```

**El paquete del sistema ofrece la 1.22. La versión oficial vigente, medida en el mismo instante, es la
1.27.1.** Y esto no es un descuido que alguien vaya a corregir: es la política de Ubuntu y de Mint.

**¿Por qué se congela así, a propósito?** Ubuntu 24.04 es una versión LTS (*Long Term Support*): el día
que sale, sus paquetes principales se **congelan** y solo reciben parches de seguridad, nunca versiones
nuevas del programa. Es una decisión correcta para un servidor que no debe cambiar de comportamiento
solo por instalar actualizaciones — pero significa que un compilador de un lenguaje que libera versiones
nuevas dos veces al año se queda fijo en la que existía cuando Ubuntu 24.04 se publicó (abril de 2024),
para siempre, mientras esa versión de Ubuntu exista. Linux Mint no tiene sus propios paquetes de Go:
hereda literalmente los de Ubuntu, así que hereda también el congelamiento.

**¿Por qué te importa a ti, concretamente?** Porque vas a buscar cosas en internet, los ejemplos van a
usar funciones o comportamientos que tu Go de 1.22 no tiene, y los errores no van a decir *«te falta
versión»*: van a decir cosas que no tienen sentido para lo que estás viendo en pantalla. Es el tipo de
problema que parece un error tuyo y en realidad es un desfase de herramienta — y te pasa **antes** de
escribir tu primera línea de código, por eso esta lección existe antes que cualquier otra.

🔑 **La forma correcta, y es la que recomienda el propio proyecto Go:** descargar el paquete oficial
directamente de `go.dev`, no del gestor de paquetes de tu distribución. Eso es lo que instala esta
lección, paso por paso, y lo que vamos a comprobar en cada paso con la salida real de los comandos.

---

## Los conceptos

### 1.1 Qué son GOROOT, GOPATH y por qué antes importaban más que ahora

Antes de instalar, vale la pena saber qué vas a tener después, porque los nombres `GOROOT` y `GOPATH`
aparecen en casi cualquier error de instalación que busques en internet, y muchas respuestas están
escritas para una versión de Go de hace diez años.

- **`GOROOT`** es la carpeta donde vive **el propio Go**: el compilador, `gofmt`, la biblioteca estándar.
  Es la carpeta que vas a crear en el paso 1.4 (`/usr/local/go`). No la tocas nunca a mano.
- **`GOPATH`** es la carpeta donde Go guarda **cosas tuyas**: paquetes descargados de internet, y —en
  versiones viejas de Go, antes de 2019— también donde tenías que poner **todo** tu código, sin
  excepción, dentro de una estructura fija (`$GOPATH/src/github.com/tu-usuario/tu-proyecto`). Si alguna
  vez ves un tutorial que te pide crear esa estructura de carpetas, es de esa época.

Hoy casi no la tocas porque desde Go 1.11 (2018) existen los **módulos** (`go.mod`, que vas a crear en el
paso 1.6): tu proyecto puede vivir en cualquier carpeta, con el nombre que quieras, y Go ya no necesita
que sigas una estructura de carpetas impuesta. `GOPATH` sigue existiendo, pero ahora solo como caché de
paquetes descargados, no como el único lugar donde puede vivir tu código.

Compruébalo con `go env`, que te muestra la configuración vigente de tu instalación (esto lo puedes
correr **después** de instalar, en el paso 1.5):

```bash
go env GOROOT GOPATH GOBIN
```

En una instalación limpia, recién hecha, vas a ver algo como:

```
/usr/local/go
/home/tu-usuario/go
/home/tu-usuario/go/bin
```

Ninguna de esas tres carpetas la creaste tú a mano: la primera la crea el paso de instalación (1.4), las
otras dos las decide Go solo, con valores por omisión razonables.

### 1.2 El PATH: qué es, y por qué "ya está instalado" no significa "ya funciona"

Cuando escribes un comando en la terminal, por ejemplo `go`, la terminal no sabe mágicamente dónde está
ese programa: revisa, una por una, una lista de carpetas guardada en una variable llamada **PATH**, y usa
el primer programa que encuentre con ese nombre. Si ninguna carpeta de la lista tiene un programa llamado
`go`, responde con un error — y ese error es literal, no aproximado:

```bash
$ go version
bash: go: command not found
```

Esto lo verifiqué en un contenedor con Go **ya extraído** en `/usr/local/go`, antes de tocar el PATH: el
programa existe en el disco, pero la terminal no lo encuentra porque no sabe dónde buscar. **"Instalado"
y "en el PATH" son dos cosas distintas**, y la confusión entre ambas es la fuente de casi todos los
tropiezos de esta lección.

Esto también explica por qué el mensaje cambia un poco según la terminal (`zsh` en vez de `bash`, que
verás si usas macOS para seguir el curso mientras practicas, o si cambiaste la terminal por defecto de tu
Mint):

```bash
$ go version
zsh: command not found: go
```

Mismo problema, mismo mecanismo, **orden de palabras distinto**: identifica cuál es tu caso con
`echo $SHELL` antes de buscar el error en internet, porque buscar el mensaje equivocado te va a llevar a
respuestas para el shell equivocado.

### 1.3 Averigua cuál es la última versión

No la copies de aquí: este documento envejece igual que los repositorios de Mint, y ya viste en la
sección anterior cuánto le puede pesar a una herramienta quedarse fija en el tiempo.

```bash
curl -s 'https://go.dev/VERSION?m=text' | head -1
```

Te va a responder algo como `go1.27.1`. **Ésa es la que vas a instalar**, sea cual sea cuando tú lo hagas.

### 1.4 Descárgala e instálala

Sustituye `go1.27.1` por lo que te dijo el comando anterior. `linux-amd64` es lo correcto para una PC o
laptop normal (si tu máquina fuera ARM, sería `linux-arm64`; para saberlo: `dpkg --print-architecture`).

```bash
cd /tmp
wget https://go.dev/dl/go1.27.1.linux-amd64.tar.gz
```

Y ahora la instalación:

```bash
sudo rm -rf /usr/local/go                          # borra una instalación anterior, si había
sudo tar -C /usr/local -xzf go1.27.1.linux-amd64.tar.gz
```

**Qué acabas de hacer**, porque conviene entenderlo y no solo copiarlo:

- `tar` descomprime el archivo.
- `-C /usr/local` le dice *«hazlo dentro de esta carpeta»*.
- `-xzf` es «extrae» (`x`), «está comprimido con gzip» (`z`), «de este archivo» (`f`).
- El resultado es una carpeta `/usr/local/go` con todo Go adentro — el `GOROOT` de la sección 1.1.

**¿Por qué `sudo`, si nunca lo habías necesitado para instalar algo con un gestor de paquetes?** Porque
`/usr/local` es una carpeta del sistema, no de tu usuario, y en Linux escribir ahí requiere permisos de
administrador. Sin `sudo`, esto es exactamente lo que vas a ver —lo provoqué a propósito, sin `sudo`, para
que veas el mensaje real y no uno inventado—:

```
$ tar -C /usr/local -xzf go1.27.1.linux-amd64.tar.gz
tar: go: Cannot mkdir: Permission denied
tar: go/VERSION: Cannot open: No such file or directory
tar: go/api: Cannot mkdir: No such file or directory
tar: go/api/README: Cannot open: No such file or directory
... (se repite, una vez por cada archivo del paquete)
```

**No es un error, es cientos.** `tar` intenta escribir cada archivo del paquete, uno por uno, y cada uno
falla igual porque ninguno tiene permiso de escribir en `/usr/local`. Si ves esta pared de líneas
repetidas, la causa es siempre la misma y siempre se arregla igual: antepón `sudo`.

### 1.5 Dile a tu sistema dónde está — y la trampa de las terminales modernas

Ahora Go existe en `/usr/local/go/bin/go`, pero como viste en la sección 1.2, si escribes `go` te va a
decir `command not found`: esa carpeta no está en el PATH todavía.

La instrucción que vas a ver en casi cualquier tutorial, incluida la documentación oficial de Go, es
ésta:

```bash
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ~/.profile
```

- `~/.profile` es un archivo que, en teoría, tu sistema lee **cada vez que inicias sesión**.
- `source ~/.profile` aplica el cambio **en esta terminal, ahora mismo**, sin esperar a la próxima
  sesión.

Y aquí viene lo que la mayoría de los tutoriales no te dicen, y que comprobé con una prueba real: **una
terminal nueva no siempre es una "sesión nueva".**

`~/.profile` solo lo lee lo que Linux llama una **shell de login** — la que arranca cuando inicias
sesión en el sistema (por ejemplo, al encender la máquina y entrar con tu usuario y contraseña). Pero la
mayoría de las aplicaciones de terminal (la terminal de Mint incluida, en su configuración por omisión)
al abrir una ventana o pestaña nueva **no** arrancan una shell de login: arrancan una shell interactiva
normal, y ésas leen otro archivo distinto, `~/.bashrc`, **no** `~/.profile`.

Lo comprobé así, simulando exactamente ese escenario —"ya edité el archivo, cierro la terminal, abro
otra"— en un contenedor recién instalado:

```
# después de agregar el export SOLO a ~/.profile, en una terminal nueva:
$ go version
bash: go: command not found          ← sigue sin funcionar

# después de agregar la MISMA línea también a ~/.bashrc, en una terminal nueva:
$ go version
go version go1.27.1 linux/arm64      ← ahora sí
```

**La instrucción "cierra la terminal, abre otra" no siempre alcanza**, y cuando no alcanza, parece que
hiciste algo mal cuando en realidad seguiste el tutorial al pie de la letra. Por eso la recomendación de
esta lección, más robusta que la de la mayoría de las guías, es agregar la línea **a los dos archivos**:

```bash
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

`~/.profile` cubre el caso de una sesión de login de verdad (por ejemplo, si usas SSH o cambias de
usuario); `~/.bashrc` cubre el caso, mucho más común en el día a día, de abrir una ventana o pestaña de
terminal nueva sobre una sesión que ya estaba abierta.

💡 **Si usas `zsh`** en vez de `bash` (lo sabes si tu terminal se ve distinta o si te lo cambiaron), el
archivo equivalente a `.bashrc` es `~/.zshrc`. Para saber cuál shell usas: `echo $SHELL`.

🔑 **Cómo distinguir si tu terminal abre una shell de login, sin adivinar:** corre `echo $0` recién
abierta. Si la respuesta empieza con un guion (`-bash`, `-zsh`), es una shell de login y sí lee
`~/.profile`. Si no lleva el guion (`bash`, `zsh`), no lo es, y necesitas el cambio en `~/.bashrc` (o
`~/.zshrc`) para que persista.

### 1.6 Comprueba que funcionó

```bash
go version
```

Debe responder algo como `go version go1.27.1 linux/amd64` (o `linux/arm64` si tu máquina es ARM). Si
sigue diciendo `command not found`, revisa con `echo $0` de qué tipo de shell se trata y en qué archivo
pusiste el cambio, siguiendo la sección 1.5.

### 1.7 Tu primer programa

```bash
mkdir -p ~/w/curso-go/hola && cd ~/w/curso-go/hola
go mod init hola
```

`go mod init` crea un archivo `go.mod`. Es la ficha de identidad del proyecto: dice cómo se llama y qué
versión de Go usa. Es el mecanismo de módulos que ya mencionamos en la sección 1.1, y es el mismo con el
que vas a arrancar el proyecto `revisor` en la lección 5.

Crea `main.go` con esto:

```go
package main

import "fmt"

func main() {
    fmt.Println("hola, ya tengo Go")
}
```

Línea por línea, porque cada una tiene su razón:

| | |
|---|---|
| `package main` | *«este archivo pertenece al paquete `main`»*. **El paquete `main` es especial: es el único que produce un programa ejecutable.** Sin esta línea tendrías una biblioteca, no un programa |
| `import "fmt"` | *«voy a usar cosas del paquete `fmt`»* (de *format*), que trae lo necesario para imprimir. Go **no** trae nada cargado por omisión: lo que uses, lo pides |
| `func main()` | **la función donde empieza tu programa.** Cuando lo ejecutas, Go busca exactamente esta función. Si se llama de otra forma, no arranca |
| `fmt.Println(...)` | imprime y salta de línea. El punto significa *«la función `Println` que está dentro de `fmt`»* |

Córrelo:

```bash
$ go run main.go
hola, ya tengo Go
```

Esa salida es real: la corrí antes de escribir esta línea.

### 1.8 Los comandos que vas a usar siempre — y la sutileza de `go.mod`

```bash
go run main.go     # compila y ejecuta de una vez, sin dejar archivo. Para probar mientras trabajas
go build           # crea el programa ejecutable y lo deja ahí
go fmt ./...       # ordena tu código
go test ./...      # corre las pruebas (lección 5)
```

**Una sutileza real que vale la pena probar tú mismo, porque cambió entre versiones de Go:** con un
programa de un solo archivo como el de arriba, `go run main.go` funciona **incluso sin `go.mod`**. Lo
comprobé en una carpeta nueva, sin `go mod init`:

```
$ go run main.go
hola, ya tengo Go
```

Pero `go build`, en esa misma carpeta sin `go.mod`, sí exige el módulo:

```
$ go build
go: go.mod file not found in current directory or any parent directory; see 'go help modules'
```

**¿Por qué la diferencia?** `go run` sobre un solo archivo puede resolverlo todo sin necesitar saber el
nombre del módulo, porque no hay nada que otro archivo pueda importar de él. En cuanto tu programa
necesita más de un archivo — por ejemplo, si intentas importar un paquete propio, como vas a hacer en la
lección 5 con el `revisor` — Go necesita el `go.mod` para saber cómo se llama tu módulo y así resolver
esos imports. Lo comprobé también:

```
$ go run main.go            # main.go importa "hola/utilidades", sin go.mod
main.go:4:5: package hola/utilidades is not in std (/usr/local/go/src/hola/utilidades)
```

**La lección práctica:** corre `go mod init` siempre, desde el primer archivo, aunque `go run` a veces
funcione sin él — te ahorras este error exacto en cuanto tu programa necesite importar un paquete propio,
que va a pasar ya en la lección 5.

**Prueba la diferencia entre `run` y `build`:**

```bash
$ go build
$ ls -la
-rwxr-xr-x 1 tu-usuario tu-usuario 2341973 ... hola
$ ./hola
hola, ya tengo Go
```

### 1.9 Por qué el binario es portátil: enlazado estático

🔑 **El archivo `hola` que produjo `go build` es un programa completo y autosuficiente.** Y esto no hay
que creerlo de fe: lo comprobé copiando ese binario exacto, compilado en una máquina con Go instalado, a
un contenedor **limpio, sin Go**:

```
$ which go
go no esta instalado
$ ./hola
hola, ya tengo Go
```

Funcionó. **No necesita intérprete, ni máquina virtual, ni bibliotecas externas instaladas aparte**,
porque Go hace *enlazado estático*: en vez de decir «cuando corras, busca la biblioteca `fmt` en algún
lado del sistema» (que es como funcionan muchos programas en C o en Python), Go copia dentro del propio
binario todo lo que tu programa necesita para correr. El archivo pesa más por eso —poco más de 2 MB para
un programa de una línea— pero a cambio puedes copiarlo a cualquier máquina Linux compatible y va a
funcionar tal cual, sin instalar nada más.

Ésa es la característica más práctica de Go, y por eso tantas herramientas de servidores están escritas
en él — incluido el `revisor` que vas a compilar y publicar en la lección 7.

### 1.10 Un editor que te ayude

No es obligatorio, pero te va a ahorrar mucho tiempo. Si usas **VS Code**, instala la extensión oficial
**Go** (de `golang.go`). Medido el 30-sep-2026 en el Marketplace de Visual Studio: la versión publicada
es la **0.57.2**. Ese número también va a envejecer — lo relevante es que instales la que ofrezca el
Marketplace el día que lo hagas, no que anotes este número.

La extensión te va a subrayar los errores mientras escribes, en lugar de que te enteres al compilar, y te
deja saltar a la definición de cualquier cosa con `F12`. La primera vez te va a pedir instalar unas
herramientas extra (`gopls`, el servidor de lenguaje de Go, entre otras): dile que sí a todas.

---

## El error que vas a ver

Todos estos son mensajes reales, provocados a propósito para esta lección — no descripciones
aproximadas de lo que "probablemente" dirían:

| El síntoma | Mensaje literal | Qué pasa y qué hacer |
|---|---|---|
| Terminal `bash` sin Go en el PATH | `bash: go: command not found` | El programa no está en ninguna carpeta del PATH. Revisa la sección 1.5 |
| Terminal `zsh` sin Go en el PATH | `zsh: command not found: go` | Mismo problema, orden de palabras distinto porque es otro shell |
| `tar` sin `sudo` contra `/usr/local` | `tar: go: Cannot mkdir: Permission denied` (repetido por archivo) | Faltó `sudo` antes del comando de extracción (sección 1.4) |
| Agregaste el PATH solo a `~/.profile` y abriste una terminal nueva | `bash: go: command not found` (persiste) | Tu terminal no abre una shell de login; agrega la línea también a `~/.bashrc` (sección 1.5) |
| `go build` en una carpeta sin `go.mod` | `go: go.mod file not found in current directory or any parent directory; see 'go help modules'` | Corre `go mod init <nombre>` en esa carpeta |
| Import de un paquete propio sin `go.mod` | `main.go:4:5: package hola/utilidades is not in std (...)` | Igual que el anterior: sin módulo, Go no sabe resolver tus propios paquetes |
| Olvidaste `import "fmt"` y usaste `fmt.Println` | `# command-line-arguments`<br>`./main.go:4:2: undefined: fmt` | Go conoce todos los símbolos disponibles; si no importaste el paquete, no existe para él |
| `go version` responde **1.22** en vez de la vigente | (sin error, pero versión equivocada) | Se está usando el de `apt`. Quítalo con `sudo apt remove golang-go` y confirma que `/usr/local/go/bin` esté en tu PATH |

**Y la regla general, que aplica a los ocho:** copia el mensaje de error completo y búscalo tal cual, sin
resumirlo con tus palabras primero. Casi siempre alguien ya lo tuvo, y el mensaje casi siempre dice
exactamente qué falta — adivinar sin leerlo completo es la forma más lenta de resolverlo.

---

## Lo que se hace mal

- **Instalar Go con el gestor de paquetes del sistema (`apt install golang`).** Es la sugerencia más
  común en tutoriales viejos, y es exactamente el problema de la sección "Por qué importa": una versión
  LTS congelada te deja años de atraso sin avisarte con ningún error, solo con comportamientos raros más
  adelante.
- **Dar por hecho que "cerrar y abrir la terminal" siempre recarga la configuración.** Es la trampa
  medida en la sección 1.5: si tu terminal no abre una shell de *login*, cerrarla y abrirla de nuevo no
  vuelve a leer `~/.profile`. El síntoma —"lo hice bien y no funciona"— no significa que hiciste algo
  mal, significa que ese archivo no era el que tu terminal lee.
- **Extraer la versión nueva encima de una instalación vieja, sin borrar antes.** El
  `sudo rm -rf /usr/local/go` del paso 1.4 no es decorativo: si extraes encima de una instalación vieja,
  quedan archivos de las dos mezclados y el resultado es un Go que falla de formas incomprensibles.
  Borra primero, siempre.
- **Confiar en que `go run` funciona sin `go.mod` y por eso saltarte `go mod init`.** Es cierto para un
  archivo suelto (sección 1.8), pero deja de serlo en cuanto tu programa tiene más de un archivo — y para
  entonces ya escribiste código que vas a tener que reorganizar. Corre `go mod init` desde el principio,
  siempre.
- **Discutir el estilo del código a mano.** En Go el estilo no se discute: lo decide la herramienta.
  Prueba escribir esto a propósito, todo torcido:

  ```go
  package main
  import "fmt"
  func main(){
  fmt.Println( "hola" )
  }
  ```

  Corre `go fmt ./...` y vuelve a abrir el archivo: **quedó ordenado solo.** No existen las peleas sobre
  dónde va la llave o cuántos espacios lleva la indentación, porque `gofmt` tiene una sola respuesta y
  todo el mundo la usa. Acostúmbrate a correrlo antes de guardar, en vez de formatear a mano.

---

## Ejercicios

1. Instala Go siguiendo las secciones 1.3 a 1.6 y confirma la versión con `go version`.
2. Corre `echo $0` en tu terminal **antes** de tocar el PATH. Según lo que responda, decide si necesitas
   tocar `~/.profile`, `~/.bashrc`, o los dos (sección 1.5) — y explica por qué, con tus palabras, antes
   de seguir.
3. Crea el proyecto `hola` de la sección 1.7, corre `go run main.go` y luego `go build` + `./hola`.
4. (Como en la sección 1.9) Copia tu binario `hola` a otra máquina Linux, o a una máquina virtual /
   contenedor sin Go instalado, y confirma que corre igual.
5. Cambia el mensaje de `fmt.Println` por algo tuyo, desordena la sangría a propósito, y corre
   `go fmt ./...`. Verifica que el archivo quedó bien formateado.
6. (Un poco más difícil) Crea una carpeta nueva, **sin** `go mod init`, con un `main.go` de un solo
   archivo. Confirma que `go run main.go` funciona de todos modos. Después corre `go build` en esa misma
   carpeta y compara el error con el de la sección 1.8. Explica, con tus palabras, por qué uno funciona y
   el otro no.
7. (Un poco más difícil) Borra a propósito la línea `import "fmt"` y corre `go run main.go`. Lee el
   error completo, sin buscarlo todavía, e intenta explicar con tus palabras qué te está diciendo antes
   de seguir.

### Soluciones

1. `go version` debe imprimir la misma versión que viste en `curl -s 'https://go.dev/VERSION?m=text'`,
   nunca `go1.22.x` (ésa es la del `apt` de Mint).
2. Si `echo $0` responde con guion al inicio (`-bash`), tu terminal abre shells de login y `~/.profile`
   basta. Si responde sin guion (`bash`), no es una shell de login, y solo `~/.bashrc` (o `~/.zshrc` en
   zsh) va a persistir entre ventanas nuevas — agrega el cambio ahí también, como en la sección 1.5.
3. `go run main.go` imprime `hola, ya tengo Go`. Tras `go build`, aparece un archivo `hola` en la
   carpeta; `./hola` lo ejecuta directo, sin volver a compilar.
4. El binario corre igual en la máquina sin Go, porque Go hace enlazado estático (sección 1.9): todo lo
   que el programa necesita ya está copiado dentro del propio archivo.
5. Antes de `go fmt`, el archivo se ve exactamente como se escribió (torcido). Después, `gofmt` lo
   reacomoda con la indentación y los espacios estándar de Go — sin que tú decidas nada de eso.
6. `go run main.go` funciona porque, con un solo archivo, Go no necesita resolver ningún import propio
   para compilarlo y ejecutarlo de una vez. `go build`, en cambio, exige saber el nombre del módulo desde
   el primer momento y responde `go: go.mod file not found...`. La diferencia es que `build` deja el
   binario listo para usarse fuera de esa invocación puntual, y para eso necesita una identidad de
   proyecto — `run` no.
7. El compilador responde:

   ```
   # command-line-arguments
   ./main.go:4:5: undefined: fmt
   ```

   Dice que usaste `fmt.Println` sin haber importado el paquete `fmt`: el compilador conoce todos los
   símbolos que puedes usar, y si no lo importaste, no existe para él. Se arregla devolviendo la línea
   `import "fmt"`.

---

## Cómo sé que lo logré

- [ ] `go version` responde con la versión que vi en `go.dev/VERSION`, no con 1.22.
- [ ] Sé si mi terminal abre shells de login o no (`echo $0`), y en qué archivo(s) puse el cambio de PATH
      en consecuencia.
- [ ] `go run main.go` imprime mi mensaje.
- [ ] `go build` produjo un archivo y `./hola` funciona.
- [ ] Copié mi binario a otra máquina (o contenedor) sin Go y corrió igual.
- [ ] Probé `go fmt` sobre código mal escrito y lo ordenó.
- [ ] Provoqué el error de `go build` sin `go.mod` y puedo explicar por qué `go run` no lo tuvo.
- [ ] Provoqué el error de `import` faltante y entendí el mensaje sin necesitar la solución.
- [ ] Puedo explicar, sin ver el texto, qué diferencia hay entre `GOROOT`, `GOPATH` y el `PATH`.

---

## Resumen

- Nunca instales Go con el gestor de paquetes del sistema (`apt install golang`): las versiones LTS de
  Linux se congelan, y el paquete queda años atrás sin avisar con ningún error.
- **`GOROOT`** es dónde vive Go; **`GOPATH`** es dónde Go guarda paquetes descargados; ninguno de los dos
  lo tocas a mano en un proyecto moderno con módulos (`go.mod`).
- El **`PATH`** es la lista de carpetas donde la terminal busca programas. Tener Go instalado y tener Go
  en el `PATH` son dos cosas distintas.
- Una terminal nueva **no siempre** es una sesión de *login*: la mayoría solo lee `~/.bashrc` (o
  `~/.zshrc`), no `~/.profile`. Agrega el cambio de `PATH` a los dos archivos.
- `go mod init` crea la identidad del proyecto (`go.mod`); todo proyecto de Go la necesita, aunque
  `go run` a veces funcione sin ella con un solo archivo.
- `go run` compila y ejecuta sin dejar archivo; `go build` deja el binario; `go fmt` da formato
  automático, sin discusión posible sobre estilo.
- El binario que produce `go build` es autosuficiente por **enlazado estático**: corre en otra máquina
  Linux compatible sin tener Go instalado.
- Cada mensaje de error de esta lección es literal, provocado a propósito: `command not found`,
  `Permission denied`, `go.mod file not found`, `undefined: fmt` — cópialos tal cual cuando los busques.

---

## Para leer más

1. [Download and install](https://go.dev/doc/install) — la guía oficial de instalación, la fuente de la
   verdad si algo en esta lección envejece. Incluye la misma recomendación de `~/.profile` que usamos
   aquí, con la misma advertencia de que el cambio "puede no aplicarse hasta el siguiente inicio de
   sesión" — que es justo lo que se midió y se explicó a fondo en la sección 1.5.
2. [A Tour of Go](https://go.dev/tour/) — el recorrido interactivo oficial; empieza justo donde termina
   esta lección.
3. [Effective Go](https://go.dev/doc/effective_go) — incluye la sección sobre `gofmt` y por qué el
   formato no se discute en Go.
4. [GNU Bash — Manual: Bash Startup Files](https://www.gnu.org/software/bash/manual/html_node/Bash-Startup-Files.html) —
   la referencia oficial de cuándo se lee `~/.profile` y cuándo `~/.bashrc`; vale la pena leerla completa
   una sola vez en la vida, porque la confusión entre los dos no es exclusiva de Go.

### Términos de esta lección

| Término | Qué significa |
|---|---|
| `GOROOT` | la carpeta donde vive la instalación de Go: el compilador, `gofmt`, la biblioteca estándar |
| `GOPATH` | la carpeta donde Go guarda paquetes descargados; antes de los módulos, también tenía que vivir ahí todo tu código |
| `PATH` | la lista de carpetas donde la terminal busca programas ejecutables |
| shell de *login* | la sesión de terminal que arranca al iniciar sesión en el sistema; es la única que lee `~/.profile` |
| `go.mod` | el archivo que identifica un proyecto de Go: su nombre y la versión de Go que usa |
| binario | el archivo ejecutable que produce `go build`, autosuficiente, sin depender de tener Go instalado |
| enlazado estático | la técnica por la que Go copia dentro del binario todo lo que el programa necesita, en vez de buscarlo en el sistema al correr |
| `gofmt` | la herramienta que da formato automático al código Go, sin opciones que discutir |

---

**Anterior:** [Lección 0 — Qué es Go](00-introduccion.md) ·
**Siguiente:** [Lección 2 — Variables, funciones y tipos](02-fundamentos.md)
