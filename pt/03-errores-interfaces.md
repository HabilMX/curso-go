# Lição 3 — Structs, métodos, erros e interfaces

> **Aqui Go vira Go.** Esta é a lição mais importante do curso.

**Duração:** duas sessões de 60 minutos. Não faça tudo de uma vez.

**Ao terminar, você vai conseguir:**

- Agrupar dados relacionados em um **struct** e explicar por que isso é melhor que variáveis soltas.
- Escrever **métodos** e decidir quando usar receptor de valor e quando de ponteiro.
- Explicar o que é um **ponteiro** sem se assustar.
- Tratar erros como valores, encapsulá-los e perguntar por eles.
- Definir uma **interface** e entender por que em Go não se declara que ela é cumprida.
- Escrever código que pode ser testado **sem se conectar a nada**.

---

## Por que isso importa

Na lição 2 você guardou os dados de um serviço em três variáveis separadas (`nombre`, `url`, `puerto`).
Funciona com **um** serviço. Com dez são trinta variáveis soltas, e nada no código diz quais vão
juntas — se você errar e combinar o nome de um com a porta de outro, **o programa compila do mesmo jeito**
e produz um relatório incorreto. Esse é o primeiro problema que esta lição resolve: o **struct**.

O segundo problema é mais de fundo. O `revisor` vai ter de consultar serviços de verdade, e as
consultas de verdade **falham**: a rede cai, o serviço demora, responde com um código que você não
esperava. Um programa que não sabe lidar com isso com ordem não serve para nada em produção. Go resolve
isso de uma forma que provavelmente você não viu se vem de outra linguagem: **os erros são valores**, não
exceções que interrompem o fluxo.

E o terceiro problema é o que faz as duas peças anteriores se encaixarem sem que o `revisor`
precise saber, de antemão, **como** cada serviço vai ser verificado (HTTP hoje, talvez um banco de dados
amanhã). Isso quem resolve é a **interface**, que é também o que vai te permitir, na lição 5, testar
todo o programa **sem se conectar a nada real**.

Structs, erros e interfaces são, nessa ordem, a coluna vertebral de tudo o que vem a seguir no curso.

---

## Os conceitos

### 3.1 Structs: agrupar o que vai junto

O que você precisa é dizer à linguagem: *«estas três coisas são um serviço»*.

**Fig. 3.1** | Definir e usar um struct.

```go
 1  // fig03_01.go
 2  // Agrupa los datos de un servicio en un solo tipo.
 3  package main
 4
 5  import "fmt"
 6
 7  // Servicio agrupa todo lo que describe a un servicio que vamos a revisar.
 8  type Servicio struct {
 9      Nombre string
10      URL    string
11      Puerto int
12  }
13
14  func main() {
15      s := Servicio{
16          Nombre: "catalogo",
17          URL:    "https://catalogo.example.com",
18          Puerto: 443,
19      }
20
21      fmt.Println("nombre:", s.Nombre)
22      fmt.Println("url:", s.URL)
23      fmt.Println("puerto:", s.Puerto)
24
25      s.Puerto = 8443          // se puede modificar
26      fmt.Println("nuevo puerto:", s.Puerto)
27
28      fmt.Println(s)           // e imprimir completo
29  }
```

```bash
$ go run fig03_01.go
nombre: catalogo
url: https://catalogo.example.com
puerto: 443
nuevo puerto: 8443
{catalogo https://catalogo.example.com 8443}
```

**Linha 8: `type Servicio struct {`.** Lê-se: *«defina um tipo novo chamado `Servicio`, que é uma
estrutura»*. **Você acabou de criar um tipo que não veio com a linguagem**, e a partir de agora ele é
usado como `int` ou `string`.

**Linhas 15-19.** Cria um valor desse tipo. Os nomes de campo com dois-pontos são opcionais —você poderia
escrever só os valores, em ordem— mas **coloque-os sempre**:

> [!TIP]
> ✅ **Boa prática 3.1**
> Escreva sempre os nomes dos campos ao criar um struct: `Servicio{Nombre: "x", Puerto: 443}` em vez de
> `Servicio{"x", "", 443}`. A versão curta quebra em silêncio se alguém acrescentar um campo ou mudar a
> ordem, e o compilador não consegue te avisar porque os tipos continuam batendo.

**Linha 21: `s.Nombre`.** O ponto acessa um campo. Lê-se «o `Nombre` de `s`».

#### 3.1.1 A maiúscula é permissão, não estética

Repare que os campos começam com **maiúscula**: `Nombre`, `URL`, `Puerto`. Em Go isso não é um estilo
escolhido: **é o controle de acesso da linguagem.**

| | |
|---|---|
| `Nombre` (maiúscula) | **exportado**: visível a partir de outros pacotes |
| `nombre` (minúscula) | **não exportado**: visível apenas dentro do seu próprio pacote |

Não existe `public`, `private` nem `protected`. **A letra inicial é a regra completa.**

#### 3.1.2 O valor zero de um struct

**Fig. 3.2** | Um struct sem inicializar.

```go
 1  // fig03_02.go
 2  package main
 3
 4  import "fmt"
 5
 6  type Servicio struct {
 7      Nombre string
 8      URL    string
 9      Puerto int
10  }
11
12  func main() {
13      var vacio Servicio
14      fmt.Printf("%+v\n", vacio)
15      fmt.Println("¿el nombre está vacío?", vacio.Nombre == "")
16      fmt.Println("puerto:", vacio.Puerto)
17  }
```

```bash
$ go run fig03_02.go
{Nombre: URL: Puerto:0}
¿el nombre está vacío? true
puerto: 0
```

Cada campo recebe **o seu** valor zero: os textos ficam em `""` e os números em `0`. **Não há lixo nem
«indefinido»**, e por isso um struct recém-criado já pode ser usado sem medo.

**E aqui aparece o `%+v`**, que é muito útil para depurar: imprime o struct **com os nomes dos
campos**. Compare com `%v`, que imprime apenas os valores.

> [!TIP]
> 🧪 **Dica de teste e depuração 3.1**
> Quando você não entender o que um struct tem, imprima-o com `%+v`. É a forma mais rápida de ver todos
> os seus campos com seu nome, e vai te poupar muitíssimo tempo.

### 3.2 Métodos

Um **método** é uma função que pertence a um tipo. Em vez de escrever `etiqueta(s)`, você escreve
`s.Etiqueta()`, e a função fica ligada ao tipo ao qual corresponde.

**Fig. 3.3** | Métodos sobre um struct.

```go
 1  // fig03_03.go
 2  // Define metodos que pertenecen al tipo Servicio.
 3  package main
 4
 5  import "fmt"
 6
 7  type Servicio struct {
 8      Nombre string
 9      URL    string
10      Puerto int
11  }
12
13  // Etiqueta devuelve una descripcion legible del servicio.
14  func (s Servicio) Etiqueta() string {
15      return fmt.Sprintf("%s (%s:%d)", s.Nombre, s.URL, s.Puerto)
16  }
17
18  // EsSeguro indica si el servicio usa el puerto de HTTPS.
19  func (s Servicio) EsSeguro() bool {
20      return s.Puerto == 443
21  }
22
23  func main() {
24      a := Servicio{Nombre: "catalogo", URL: "https://catalogo.example.com", Puerto: 443}
25      b := Servicio{Nombre: "local", URL: "http://localhost", Puerto: 8080}
26
27      fmt.Println(a.Etiqueta(), "— seguro:", a.EsSeguro())
28      fmt.Println(b.Etiqueta(), "— seguro:", b.EsSeguro())
29  }
```

```bash
$ go run fig03_03.go
catalogo (https://catalogo.example.com:443) — seguro: true
local (http://localhost:8080) — seguro: false
```

**Linha 14: `func (s Servicio) Etiqueta() string`.** Esse `(s Servicio)` entre `func` e o nome se chama
**receptor**, e é o que converte uma função em um método. Lê-se: *«esta função pertence ao tipo
`Servicio`, e dentro dela vou me referir ao valor como `s`»*.

> [!NOTE]
> 🔧 **Observação de engenharia de software 3.1**
> Se você vem de Java, C# ou Python, um método de Go é parecido com um método de classe — com uma
> diferença importante: **em Go o método é escrito fora do tipo.** Você pode ter o `type` em um arquivo e
> seus métodos em outro. E como não há classes, também não há herança: em Go o código é reutilizado de
> outra forma, e você vai vê-la na seção 3.6.

### 3.3 Ponteiros, em dez minutos

Os ponteiros têm fama de difíceis. Em Go são muito mais simples que em C, e você precisa entender
**uma única ideia** para o que vem a seguir.

#### 3.3.1 O problema

**Fig. 3.4** | Um método que **não funciona** como você esperaria.

```go
 1  // fig03_04.go
 2  // Muestra por que un receptor de VALOR no puede modificar el original.
 3  package main
 4
 5  import "fmt"
 6
 7  type Servicio struct {
 8      Nombre string
 9      Puerto int
10  }
11
12  // CambiarPuerto intenta modificar el servicio... y no lo logra.
13  func (s Servicio) CambiarPuerto(nuevo int) {
14      s.Puerto = nuevo
15  }
16
17  func main() {
18      x := Servicio{Nombre: "catalogo", Puerto: 443}
19      fmt.Println("antes:", x.Puerto)
20
21      x.CambiarPuerto(8443)
22      fmt.Println("despues:", x.Puerto)      // ¿cambió?
23  }
24
```

```bash
$ go run fig03_04.go
antes: 443
despues: 443
```

**Não mudou nada.** O método foi executado, não houve erro, e o valor continua igual.

**Por quê?** Porque quando o receptor é `(s Servicio)`, Go entrega ao método **uma cópia** do struct. O
método modifica a cópia, a cópia é descartada ao terminar, e o original nunca ficou sabendo. É um
antipadrão tão frequente que tem sua própria entrada em «O que se faz errado», mais abaixo.

#### 3.3.2 A solução: o receptor de ponteiro

Um **ponteiro** é uma variável que guarda **o endereço** de outra, em vez de uma cópia do seu conteúdo.

Pense na diferença entre te darem **uma fotocópia** de um documento e te darem **o endereço do arquivo
onde está o original**. Com a fotocópia você pode escrever o que quiser: o original não muda. Com o
endereço, você pode ir lá e modificar o original.

- `Servicio` é a fotocópia.
- `*Servicio` é o endereço do original.

**Fig. 3.5** | O mesmo método, agora com receptor de ponteiro.

```go
 1  // fig03_05.go
 2  // Con receptor de PUNTERO, el metodo si modifica el original.
 3  package main
 4
 5  import "fmt"
 6
 7  type Servicio struct {
 8      Nombre string
 9      Puerto int
10  }
11
12  // CambiarPuerto modifica el servicio. El * es la diferencia.
13  func (s *Servicio) CambiarPuerto(nuevo int) {
14      s.Puerto = nuevo
15  }
16
17  func main() {
18      x := Servicio{Nombre: "catalogo", Puerto: 443}
19      fmt.Println("antes:", x.Puerto)
20
21      x.CambiarPuerto(8443)
22      fmt.Println("despues:", x.Puerto)
23  }
```

```bash
$ go run fig03_05.go
antes: 443
despues: 8443
```

**A única diferença entre a Fig. 3.4 e a 3.5 é um asterisco na linha 13.** Só isso.

E repare na linha 21: você escreveu `x.CambiarPuerto(8443)` igual a antes, **sem `&` nem nada estranho**. Go
percebe que o método precisa de um ponteiro e o pega sozinho. Essa é a razão pela qual os ponteiros
de Go dão muito menos trabalho que os de C.

> [!TIP]
> ✅ **Boa prática 3.2**
> A regra é simples: **se o método modifica o struct, receptor de ponteiro (`*T`); se apenas lê, de
> valor (`T`)**. E seja consistente dentro de um mesmo tipo: se a maioria dos seus métodos precisa de
> ponteiro, use ponteiro em todos, mesmo que alguns não o exijam. Misturá-los confunde quem lê o código.

> [!NOTE]
> 🚀 **Dica de desempenho 3.1**
> Há uma segunda razão para usar ponteiros: **evitar cópias**. Se um struct tem vinte campos, cada
> chamada com receptor de valor copia os vinte. Para structs pequenos é irrelevante; para grandes ou em
> laços de milhões de voltas, importa. **Não otimize isso sem medir:** a clareza vale mais que uma
> cópia de 40 bytes.

#### 3.3.3 O único ponteiro perigoso: `nil`

Ainda há uma terceira situação com ponteiros, e esta sim produz uma mensagem de erro real —tão real que
tem sua própria seção dedicada: **«O erro que você vai ver»**, mais abaixo nesta lição. Adiantamos a
ideia: o valor zero de um ponteiro é **`nil`** («não aponto para nada»), e ler um campo através de um
ponteiro nulo faz o programa **estourar** em vez de se comportar mal em silêncio como na Fig. 3.4.

### 3.4 Os erros são valores

Você já viu na lição 2 que uma função de Go pode devolver `(resultado, error)`. Agora vamos ao fundo,
porque **isto é o que mais distingue Go do que você provavelmente já viu.**

#### 3.4.1 `error` é uma interface, não algo mágico

Em Go, `error` é simplesmente um tipo com um método:

```go
type error interface {
    Error() string
}
```

Isso significa: *«qualquer coisa que tenha um método `Error()` que devolva texto, é um erro»*. Não há
hierarquia de classes de exceção, não há `throw`, não há pilha de chamadas que se desenrola. **Um erro é
um valor comum que viaja como qualquer outro.**

#### 3.4.2 Criar erros

**Fig. 3.7** | As três formas de criar um erro.

```go
 1  // fig03_07.go
 2  // Muestra las tres formas de producir un error.
 3  package main
 4
 5  import (
 6      "errors"
 7      "fmt"
 8  )
 9
10  // ErrPuertoInvalido es un error CENTINELA: se declara una vez y se compara.
11  var ErrPuertoInvalido = errors.New("el puerto debe estar entre 1 y 65535")
12
13  type Servicio struct {
14      Nombre string
15      Puerto int
16  }
17
18  // Validar revisa que el servicio tenga sentido.
19  func Validar(s Servicio) error {
20      if s.Nombre == "" {
21          // 1. un error sencillo, creado al momento
22          return errors.New("el nombre no puede estar vacio")
23      }
24      if s.Puerto < 1 || s.Puerto > 65535 {
25          // 2. un error centinela, para poder compararlo despues
26          return ErrPuertoInvalido
27      }
28      if s.Puerto == 80 {
29          // 3. un error con datos dentro
30          return fmt.Errorf("el puerto %d no usa cifrado, usa 443", s.Puerto)
31      }
32      return nil          // nil significa: todo bien
33  }
34
35  func main() {
36      casos := []Servicio{
37          {Nombre: "catalogo", Puerto: 443},
38          {Nombre: "", Puerto: 443},
39          {Nombre: "pagos", Puerto: 99999},
40          {Nombre: "viejo", Puerto: 80},
41      }
42
43      for _, s := range casos {
44          if err := Validar(s); err != nil {
45              fmt.Printf("%-10s ❌ %v\n", s.Nombre, err)
46          } else {
47              fmt.Printf("%-10s ✅ valido\n", s.Nombre)
48          }
49      }
50  }
```

```bash
$ go run fig03_07.go
catalogo   ✅ valido
           ❌ el nombre no puede estar vacio
pagos      ❌ el puerto debe estar entre 1 y 65535
viejo      ❌ el puerto 80 no usa cifrado, usa 443
```

**Linha 11: o erro sentinela.** É declarado **uma vez**, no nível do pacote, com o prefixo `Err`. Serve
para que quem chama a sua função possa perguntar *«foi este erro em particular?»*, como você verá em 3.5.

**Linha 30: `fmt.Errorf`.** Como o `Printf`, mas produz um erro em vez de imprimir. Use-o quando a
mensagem precisa de dados.

**Linha 44: `if err := Validar(s); err != nil`.** Esse padrão declara `err` **dentro** do `if`, então
só existe ali. É muito idiomático em Go e mantém o código limpo.

> [!TIP]
> ✅ **Boa prática 3.3**
> As mensagens de erro se escrevem **em minúscula e sem ponto final**: `"no se pudo abrir el archivo"`,
> e não `"No se pudo abrir el archivo."`. A razão é prática: os erros são **encapsulados** uns nos outros
> (seção 3.5), e ao serem concatenados ficam como `"revisando catalogo: no se pudo abrir el archivo"`.
> Com maiúsculas e pontos, o resultado ficaria quebrado.

### 3.5 Encapsular erros: `%w`, `errors.Is` e `errors.As`

Um erro sem contexto é pouco útil. Se o seu programa diz `connection refused`, você não sabe **qual**
serviço falhou.

**Encapsular** um erro é acrescentar contexto a ele **sem perder o original**.

**Fig. 3.8** | Encapsular erros e perguntar por eles.

```go
 1  // fig03_08.go
 2  // Envuelve un error para agregar contexto sin perder el original.
 3  package main
 4
 5  import (
 6      "errors"
 7      "fmt"
 8  )
 9
10  var ErrNoResponde = errors.New("no responde")
11
12  type Servicio struct{ Nombre string }
13
14  // consultar simula la consulta de bajo nivel.
15  func consultar(s Servicio) error {
16      return ErrNoResponde
17  }
18
19  // Revisar agrega contexto al error de consultar.
20  func Revisar(s Servicio) error {
21      if err := consultar(s); err != nil {
22          // el %w ENVUELVE el error original
23          return fmt.Errorf("revisando %s: %w", s.Nombre, err)
24      }
25      return nil
26  }
27
28  func main() {
29      err := Revisar(Servicio{Nombre: "catalogo"})
30
31      fmt.Println("1. el mensaje completo:")
32      fmt.Println("  ", err)
33
34      fmt.Println("2. ¿es un ErrNoResponde, aunque este envuelto?")
35      fmt.Println("  ", errors.Is(err, ErrNoResponde))
36
37      fmt.Println("3. ¿y si pregunto por otro error?")
38      fmt.Println("  ", errors.Is(err, errors.New("otra cosa")))
39
40      fmt.Println("4. el error original, desenvuelto:")
41      fmt.Println("  ", errors.Unwrap(err))
42  }
```

```bash
$ go run fig03_08.go
1. el mensaje completo:
   revisando catalogo: no responde
2. ¿es un ErrNoResponde, aunque este envuelto?
   true
3. ¿y si pregunto por otro error?
   false
4. el error original, desenvuelto:
   no responde
```

**Linha 23: `%w`.** Este é o verbo-chave. Parece igual a `%v` ao imprimir, **mas conserva o erro
original dentro** para que se possa perguntar por ele depois.

**Linha 35: `errors.Is`.** Pergunta *«em algum ponto desta cadeia está esse erro?»*. Funciona mesmo que
haja cinco camadas de encapsulamento.

> [!WARNING]
> 🔴 **Erro comum de programação 3.3 — o mais silencioso desta lição**
> Escrever `%v` em vez de `%w` ao encapsular:
>
> ```go
> return fmt.Errorf("revisando %s: %v", s.Nombre, err)   // ❌ con %v
> ```
>
> **A mensagem impressa é idêntica.** Não há erro, não há aviso, tudo parece funcionar. Mas
> `errors.Is` deixa de encontrar o erro original e devolve `false`, de modo que o código que decidia o que
> fazer conforme o tipo de falha passa a tomar o caminho errado. **É o tipo de erro que não falha:
> devolve um dado pior.** Ao encapsular, use `%w`. Este é o segundo antipadrão da seção «O que se faz
> errado».

#### 3.5.1 `errors.As`, quando você precisa dos dados do erro

`errors.Is` responde «é este erro?». `errors.As` responde «é deste **tipo**? me dê para eu ler seus
campos».

**Fig. 3.9** | Um erro próprio com dados, recuperado com `errors.As`.

```go
 1  // fig03_09.go
 2  // Define un tipo de error propio y recupera sus datos.
 3  package main
 4
 5  import (
 6      "errors"
 7      "fmt"
 8  )
 9
10  // ErrorHTTP es un error que lleva datos adentro.
11  type ErrorHTTP struct {
12      Codigo int
13      URL    string
14  }
15
16  // Error hace que ErrorHTTP cumpla la interfaz error.
17  func (e *ErrorHTTP) Error() string {
18      return fmt.Sprintf("el servidor respondio %d", e.Codigo)
19  }
20
21  func consultar(url string) error {
22      return &ErrorHTTP{Codigo: 503, URL: url}
23  }
24
25  func Revisar(nombre, url string) error {
26      if err := consultar(url); err != nil {
27          return fmt.Errorf("revisando %s: %w", nombre, err)
28      }
29      return nil
30  }
31
32  func main() {
33      err := Revisar("catalogo", "https://catalogo.example.com")
34      fmt.Println("mensaje:", err)
35
36      var errHTTP *ErrorHTTP
37      if errors.As(err, &errHTTP) {
38          fmt.Println("es un ErrorHTTP")
39          fmt.Println("  codigo:", errHTTP.Codigo)
40          fmt.Println("  url:   ", errHTTP.URL)
41
42          if errHTTP.Codigo >= 500 {
43              fmt.Println("  → es culpa del servidor, conviene reintentar")
44          }
45      }
46  }
```

```bash
$ go run fig03_09.go
mensaje: revisando catalogo: el servidor respondio 503
es un ErrorHTTP
  codigo: 503
  url:    https://catalogo.example.com
  → es culpa del servidor, conviene reintentar
```

**Linha 17.** Ao escrever um método `Error() string`, o seu tipo **já é um `error`**. Você não declarou
nada: o tipo cumpre a interface porque tem o método. É isso que explica a seção seguinte.

**Linha 42.** E aqui está o valor real: o programa pode **decidir** conforme o tipo de falha. Um 503 é
repetido; um 404, não. Com erros só de texto isso seria impossível sem comparar strings, o que é frágil.

### 3.6 Interfaces: o coração de Go

Uma **interface** é uma lista de métodos. Qualquer tipo que tenha esses métodos **a cumpre**, e não é
preciso declará-lo em lugar nenhum.

**Fig. 3.10** | Uma interface e duas implementações.

```go
 1  // fig03_10.go
 2  // Una interfaz con dos implementaciones: la real y una de prueba.
 3  package main
 4
 5  import (
 6      "errors"
 7      "fmt"
 8  )
 9
10  type Servicio struct {
11      Nombre string
12      URL    string
13  }
14
15  type Estado struct {
16      Servicio Servicio
17      Codigo   int
18      Err      error
19  }
20
21  // Revisor es una INTERFAZ: cualquier cosa con este metodo la cumple.
22  type Revisor interface {
23      Revisar(s Servicio) Estado
24  }
25
26  // ---- primera implementacion: la de verdad (simplificada) ----
27  type RevisorHTTP struct{}
28
29  func (r RevisorHTTP) Revisar(s Servicio) Estado {
30      // En la leccion 7 esto hara una peticion HTTP real.
31      return Estado{Servicio: s, Codigo: 200}
32  }
33
34  // ---- segunda implementacion: para probar, sin red ----
35  type RevisorFalso struct {
36      Respuesta Estado
37  }
38
39  func (r RevisorFalso) Revisar(s Servicio) Estado {
40      r.Respuesta.Servicio = s
41      return r.Respuesta
42  }
43
44  // RevisarTodos acepta CUALQUIER Revisor. No sabe ni le importa cual.
45  func RevisarTodos(r Revisor, servicios []Servicio) []Estado {
46      var estados []Estado
47      for _, s := range servicios {
48          estados = append(estados, r.Revisar(s))
49      }
50      return estados
51  }
52
53  func main() {
54      servicios := []Servicio{
55          {Nombre: "catalogo", URL: "https://catalogo.example.com"},
56          {Nombre: "pagos", URL: "https://pagos.example.com"},
57      }
58
59      fmt.Println("--- con el revisor real ---")
60      for _, e := range RevisarTodos(RevisorHTTP{}, servicios) {
61          fmt.Printf("  %-10s codigo %d\n", e.Servicio.Nombre, e.Codigo)
62      }
63
64      fmt.Println("--- con el falso, simulando una falla ---")
65      falso := RevisorFalso{
66          Respuesta: Estado{Err: errors.New("no responde")},
67      }
68      for _, e := range RevisarTodos(falso, servicios) {
69          fmt.Printf("  %-10s error: %v\n", e.Servicio.Nombre, e.Err)
70      }
71  }
```

```bash
$ go run fig03_10.go
--- con el revisor real ---
  catalogo   codigo 200
  pagos      codigo 200
--- con el falso, simulando una falla ---
  catalogo   error: no responde
  pagos      error: no responde
```

**Leia devagar, porque aqui está a ideia que sustenta todo o Go.**

**Linha 22-24: a interface.** Diz: *«um `Revisor` é qualquer coisa que tenha um método `Revisar` que
receba um `Servicio` e devolva um `Estado`»*.

**Linhas 29 e 39.** `RevisorHTTP` e `RevisorFalso` **não declaram em lugar nenhum** que cumprem `Revisor`.
Não há `implements`, não há `: Revisor`, nada. **Cumprem a interface porque têm o método**, e isso o
compilador verifica sozinho.

**Linha 45: `func RevisarTodos(r Revisor, ...)`.** Esta função **não sabe** com o que está trabalhando.
Só sabe que pode chamar `.Revisar()`. E por isso as linhas 60 e 68 passam a ela duas coisas completamente
diferentes **sem mudar uma única linha de `RevisarTodos`**.

> [!NOTE]
> 🔧 **Observação de engenharia de software 3.2 — a mais importante da lição**
> Essa satisfação implícita tem uma consequência que não existe em Java nem em C#: **você pode definir uma
> interface para código que não foi escrito por você.** Se uma biblioteca alheia tem um tipo com um método
> `Revisar`, esse tipo cumpre a *sua* interface sem que o autor soubesse nem precisasse cooperar. Em outras
> linguagens, se o autor não declarou a interface, não há nada a fazer.

> [!TIP]
> ✅ **Boa prática 3.4 — as duas regras de ouro das interfaces em Go**
> **1. Faça-as pequenas.** Um ou dois métodos. Uma interface de dez métodos quase sempre vem de outra
> linguagem — e é, de fato, o terceiro antipadrão da seção «O que se faz errado». As da biblioteca
> padrão mais usadas —`io.Reader`, `io.Writer`— têm **um** método.
> **2. Defina-as onde são USADAS, não onde são implementadas.** A interface `Revisor` pertence ao código que
> precisa verificar coisas, não ao que sabe como verificá-las. Isso inverte a dependência: quem consome
> declara o que precisa.
>
> E o resumo que você vai ouvir muito: **«aceite interfaces, devolva structs»**.

### 3.7 Por que isso torna o seu código testável

Olhe de novo a linha 65 da Fig. 3.10. Você acabou de testar `RevisarTodos` **simulando um serviço
fora do ar**, sem desligar nada, sem rede, e em um milissegundo.

Isso é o que as pessoas querem dizer quando falam de «código testável», e em Go se consegue **sem
bibliotecas de mocks, sem anotações e sem frameworks**. Apenas com uma interface pequena.

> [!NOTE]
> 🔧 **Observação de engenharia de software 3.3**
> A pergunta que convém se fazer ao projetar: *«posso testar isto sem que o mundo exterior exista?»* Se
> a resposta for não, normalmente falta uma interface. Na lição 5 você vai escrever testes de verdade, e
> vai agradecer por ter feito isto agora.

---

## O erro que você vai ver

O erro mais frequente desta lição é o **`nil pointer dereference`**, e ao contrário do antipadrão
da Fig. 3.4 (que falha em silêncio), este sim te avisa — com uma mensagem que no início assusta mais do
que deveria.

**Fig. 3.6** | Um programa que **compila perfeitamente** e **estoura ao executar**.

```go
 1  // fig03_06.go — este programa COMPILA pero truena al correr
 2  package main
 3
 4  import "fmt"
 5
 6  type Servicio struct{ Puerto int }
 7
 8  func main() {
 9      var p *Servicio          // un puntero sin apuntar a nada: vale nil
10      fmt.Println(p)           // esto sí funciona
11      fmt.Println(p.Puerto)    // esto truena
12  }
```

```bash
$ go run fig03_06.go
<nil>
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x...]

goroutine 1 [running]:
main.main()
	/home/usuario/fig03_06.go:11 +0x18
exit status 2
```

**O valor zero de um ponteiro é `nil`**: «não aponto para nada». Tentar ler um campo através de um
ponteiro nulo provoca um **panic**, que é a forma como Go aborta um programa quando acontece algo
irrecuperável.

**Como ler esta mensagem, linha por linha:**

- `panic: runtime error: invalid memory address or nil pointer dereference` — o nome do problema.
  Quando você vir isso, a causa quase sempre é um ponteiro, um map ou uma interface em `nil` que você
  tentou usar como se tivesse algo dentro.
- `goroutine 1 [running]:` — qual linha de execução de Go estava rodando quando estourou. Por ora, em um
  programa sem concorrência, sempre vai ser a 1 (a função `main`); na lição 6, com várias
  goroutines, esse número passa a importar de verdade.
- `main.main() … fig03_06.go:11` — **o arquivo e a linha exatos** onde aconteceu. Comece sempre por aí,
  não pela mensagem de cima: a mensagem diz *que tipo* de erro foi, esta linha diz *onde*.
- `exit status 2` — o programa terminou com erro (um `exit status 0` é sucesso).

> [!WARNING]
> ⚠️ **Erro comum de programação 3.2**
> O `nil pointer dereference` é o panic mais frequente em Go. Repare que **compila perfeitamente**:
> o compilador não tem como saber se um ponteiro vai valer `nil` na execução. Quando você o vir, vá direto
> ao arquivo e à linha que o rastreamento indica — no exemplo, `fig03_06.go:11`.

---

## O que se faz errado

**1. Escrever um método que deveria modificar o struct, com receptor de valor.** É o erro da
Fig. 3.4 (seção 3.3.1): compila, executa, e **não faz absolutamente nada** — sem nenhuma mensagem. É
um dos tropeços mais desconcertantes para quem está começando, justamente porque não há sintoma algum. A
regra que o evita está na Boa prática 3.2: se o método modifica, receptor de ponteiro.

**2. Encapsular um erro com `%v` em vez de `%w`.** É o erro da seção 3.5: a mensagem impressa fica
idêntica, então nada avisa, mas `errors.Is` deixa de reconhecer o erro original e o código que decidia
conforme o tipo de falha passa a tomar o caminho errado. Dos dois antipadrões desta lição, este
é o mais caro, porque ninguém nota o sintoma até que o programa já tenha tomado uma má decisão com ele.

**3. Escrever interfaces grandes, com muitos métodos, copiando o hábito de outra linguagem.** Em Go
as interfaces são pequenas —um ou dois métodos, como você viu na Boa prática 3.4— e são definidas do
lado de quem as consome, não de quem as implementa. Uma interface de dez métodos quase nunca é
necessária: normalmente basta o método ou os métodos que a função que a recebe realmente chama.

---

## Exercícios

### Exercícios de revisão

**3.1** Que vantagem um struct tem sobre três variáveis soltas?

**3.2** O que significa um campo começar com maiúscula?

**3.3** O que `fmt.Printf("%+v\n", Servicio{})` imprime se o struct tem `Nombre string` e `Puerto int`?

**3.4** Qual é a diferença entre `func (s Servicio) X()` e `func (s *Servicio) X()`, e quando se usa
cada um?

**3.5** Este método compila, executa e não faz nada. Por quê?

```go
func (s Servicio) Renombrar(nuevo string) {
    s.Nombre = nuevo
}
```

**3.6** O que é `nil` para um ponteiro e o que acontece se você ler um campo através de um?

**3.7** O que um tipo precisa ter para ser um `error`?

**3.8** Qual é a diferença entre `%w` e `%v` ao encapsular um erro, e por que é perigosa?

**3.9** Quando você usa `errors.Is` e quando `errors.As`?

**3.10** O que é preciso escrever para que um tipo cumpra uma interface em Go?

**3.11** O que este programa imprime?

```go
type Estado struct{ Codigo int }

func (e Estado) OK() bool { return e.Codigo == 200 }

func main() {
    var e Estado
    fmt.Println(e.OK())
}
```

### Exercícios de código

**3.12** Acrescente ao struct `Servicio` um campo `TimeoutMs int` e um método `Timeout() time.Duration` que o
converta. Dica: `time.Duration(s.TimeoutMs) * time.Millisecond`.

**3.13** Escreva um método `func (s *Servicio) Normalizar()` que: remova espaços do nome com
`strings.TrimSpace`, passe-o para minúsculas com `strings.ToLower`, e se a porta for 0 coloque-a em 443.
**Explique a si mesmo por que este método precisa de receptor de ponteiro.**

**3.14** Defina um tipo `Estado` com os campos `Servicio`, `Codigo`, `Duracion` e `Err`. Escreva para ele um
método `OK() bool` que devolva verdadeiro somente se não houver erro **e** o código estiver entre 200 e 299.
Teste-o com cinco casos, incluindo o `Estado{}` vazio.

**3.15** Crie um erro sentinela `ErrTimeout` e uma função que o devolva encapsulado com contexto.
Verifique com `errors.Is` que ele é detectado. **Depois troque o `%w` por `%v` e verifique que `errors.Is`
devolve `false`.** Anote no diário de bordo que a mensagem impressa não mudou.

**3.16** Defina um erro próprio `ErrorValidacion` com os campos `Campo string` e `Motivo string`, faça-o
cumprir a interface `error`, e recupere-o com `errors.As` para imprimir qual campo falhou.

**3.17** Defina a interface `Notificador` com um método `Notificar(mensaje string) error`. Implemente
`NotificadorConsola` (que imprime) e `NotificadorFalso` (que guarda as mensagens em um slice para poder
examiná-las). Escreva uma função que receba um `Notificador` e use-a com as duas.

**3.18 (Encontre o erro)** Diga o que está errado em cada um **sem compilar**, e depois compile para
confirmar:

```go
// (a)
func (s Servicio) Renombrar(n string) { s.Nombre = n }

// (b)
return fmt.Errorf("fallo al revisar %s: %v", nombre, err)

// (c)
var p *Servicio
fmt.Println(p.Nombre)

// (d)
type Revisor interface {
    Revisar(s Servicio) Estado
}
type MiRevisor struct{}
func (m MiRevisor) revisar(s Servicio) Estado { return Estado{} }
// y luego:  var r Revisor = MiRevisor{}
```

**3.19 (Projeto do curso)** Reorganize o seu programa da lição 2 usando o que há nesta lição:
1. Os structs `Servicio` e `Estado`.
2. Os métodos `Etiqueta()`, `EsSeguro()` e `OK()`.
3. A interface `Revisor` com `RevisorHTTP` (que por enquanto devolve dados inventados) e `RevisorFalso`.
4. A função `RevisarTodos(r Revisor, servicios []Servicio) []Estado`.
5. Um `main` que imprima o relatório usando o **falso**, com um serviço que falhe.

🔑 **Quando terminar, repare em algo: o seu programa já pode ser testado por completo sem se conectar a
nada, e você ainda não escreveu um único teste.** Isso é o que você acabou de ganhar nesta lição.

### Soluções

**Exercícios de revisão (3.1 a 3.11):**

**3.1** Agrupa os dados que vão juntos em uma única variável, de modo que a linguagem —e quem lê o
código— sabe que pertencem à mesma coisa. E não é possível combinar por engano dados de dois serviços
diferentes.

**3.2** Que está **exportado**: é visível a partir de outros pacotes. Com minúscula só é visto dentro do
seu pacote. É todo o controle de acesso que Go tem.

**3.3** `{Nombre: Puerto:0}` — os valores zero, com os nomes dos campos porque é `%+v`.

**3.4** O primeiro é **receptor de valor**: recebe uma cópia e não pode modificar o original. O segundo
é **receptor de ponteiro**: recebe o endereço e pode. Usa-se ponteiro quando o método modifica, ou
quando o struct é grande e copiá-lo custa.

**3.5** Porque o receptor é de **valor**: o método modifica uma cópia que é descartada ao terminar. É
preciso mudá-lo para `func (s *Servicio) Renombrar(...)`.

**3.6** `nil` é o valor zero de um ponteiro e significa «não aponto para nada». Ler um campo através de um
ponteiro `nil` provoca um **panic** (`nil pointer dereference`) e o programa aborta.

**3.7** Um método `Error() string`. Nada mais: `error` é uma interface com esse único método.

**3.8** `%w` **encapsula** o erro original e o conserva dentro; `%v` apenas o converte em texto. É
perigosa porque **a mensagem impressa é idêntica**, então não há sintoma — mas `errors.Is` deixa de
encontrar o erro original e o código que decidia conforme o tipo de falha começa a errar.

**3.9** `errors.Is` para perguntar *«é este erro em particular?»*. `errors.As` para perguntar *«é deste
tipo? me dê»*, quando você precisa ler os dados que o erro leva dentro.

**3.10** **Nada.** Basta ter os métodos que a interface pede; o compilador verifica sozinho. Não
existe `implements`.

**3.11** `false` — o valor zero de `Estado` tem `Codigo: 0`, que não é 200. Note que o método funciona
perfeitamente sobre um struct vazio: isso é o valor zero sendo útil.

**Exercícios de código (3.12 a 3.19):** suas soluções verificadas —compiladas e executadas— serão
adicionadas na entrega em que for fechado o capítulo de exercícios com solução do curso completo; não são
publicadas sem ter rodado `go build` sobre cada uma.

---

## Como sei que consegui

- Você compilou e rodou as figuras 3.1 a 3.10 e a sua saída coincide com a mostrada.
- Você consegue explicar, sem ver o texto, que diferença há entre receptor de valor e receptor de
  ponteiro, e dar um exemplo de quando usar cada um.
- Você fez o experimento do exercício 3.15 (trocar `%w` por `%v`) e viu com seus próprios olhos que
  `errors.Is` deixa de encontrar o erro **sem que a mensagem impressa mude**.
- Você consegue explicar a outra pessoa, sem tecnicismos, por que em Go «cumprir uma interface» não é
  declarado em lugar nenhum.
- Você terminou o exercício 3.19: o seu programa do `revisor` já usa structs, métodos, a interface
  `Revisor` e pode ser rodado com o `RevisorFalso` sem tocar na rede.

**Antes de fechar:** no [`bitacora.md`](bitacora.md) anote o que mais custou a você entre ponteiros, erros
e interfaces, e o resultado do exercício 3.15 —o do `%w` contra `%v`—. Esse experimento é o que mais se
esquece e o que mais caro sai em um programa real.

---

## Resumo

- Um **struct** agrupa dados relacionados em um tipo próprio. Define-se com `type Nombre struct { … }`.
- Os campos com **maiúscula** são exportados; com minúscula, privados ao pacote.
- O **valor zero** de um struct preenche cada campo com o seu próprio valor zero: nunca há lixo.
- **`%+v`** imprime um struct com os nomes dos seus campos: a melhor ferramenta para depurar.
- Um **método** é uma função com **receptor**: `func (s Servicio) X()`.
- **Receptor de valor** recebe uma cópia e **não pode modificar** o original; **receptor de ponteiro**
  (`*T`) pode.
- Um **ponteiro** guarda o endereço de outra variável. Go insere o `&` e o `*` para você ao chamar
  métodos.
- O valor zero de um ponteiro é **`nil`**; ler através dele provoca um **panic**.
- **`error` é uma interface** com um único método `Error() string`. Um erro é um valor comum.
- São criados com **`errors.New`** (simples), como **sentinelas** (`var ErrX = errors.New(...)`) ou com
  **`fmt.Errorf`** (com dados).
- **`%w`** encapsula um erro conservando o original; **`%v` o achata e quebra `errors.Is` sem avisar**.
- **`errors.Is`** pergunta se um erro está na cadeia; **`errors.As`** recupera o erro de um tipo
  concreto para ler seus dados.
- Uma **interface** é uma lista de métodos. Um tipo a cumpre **por ter os métodos**, sem declará-lo.
- As interfaces de Go são feitas **pequenas** e definidas **onde são usadas**.
- Uma interface pequena permite **testar sem o mundo exterior**: a implementação real é substituída por
  uma falsa.

---

## Para ler mais

1. **[A Tour of Go — Methods](https://go.dev/tour/methods/1)** — oficial e interativo: receptores,
   ponteiros e interfaces com exercícios no navegador.
2. **[Effective Go — Interfaces and other types](https://go.dev/doc/effective_go#interfaces_and_types)**
   — a explicação oficial de por que a satisfação de interfaces é implícita.
3. **[Error handling and Go](https://go.dev/blog/error-handling-and-go)** — a entrada do blog oficial de
   Go sobre por que os erros são valores e como são tratados de forma idiomática.
4. **[Pacote `errors` — documentação oficial](https://pkg.go.dev/errors)** — a referência exata de
   `errors.Is`, `errors.As` e `errors.Unwrap`.
5. **[APIs prontas para agentes de IA](https://www.habil.mx/pt/blog/apis-prontas-para-agentes-de-ia/)** — artigo sobre erros de API que dizem ao cliente se deve tentar de novo (um 503) ou não (um 404), a distinção da seção 3.5.

### Termos desta lição

| | |
|---|---|
| **campo** | cada dado que um struct contém |
| **encapsular (um erro)** | acrescentar contexto a ele conservando o original, com `%w` |
| **erro sentinela** | erro declarado uma vez no nível do pacote, para comparar com `errors.Is` |
| **exportado / não exportado** | visível fora do pacote (maiúscula) ou não (minúscula) |
| **interface** | lista de métodos; um tipo a cumpre ao ter esses métodos |
| **método** | função associada a um tipo mediante um receptor |
| **`nil`** | ausência de valor; o valor zero de ponteiros, interfaces e erros |
| **`nil pointer dereference`** | panic por ler através de um ponteiro nulo |
| **panic** | aborto do programa por um erro irrecuperável |
| **ponteiro** | variável que guarda o endereço de outra; seu tipo se escreve `*T` |
| **receptor** | o `(s Servicio)` que liga uma função a um tipo |
| **satisfação implícita** | cumprir uma interface sem declará-lo |
| **struct** | tipo que agrupa vários campos |

---

**Anterior:** [Lição 2 — Variáveis, funções e tipos](02-fundamentos.md) ·
**Próxima:** [Lição 4 — Coleções](04-colecciones.md)
