# Lição 2 — Variáveis, funções e tipos

**Duração:** 90 minutos, ou duas sessões de 45.

**Ao terminar, você vai conseguir:**

- Explicar o que o compilador de Go faz com o seu arquivo de texto.
- Declarar variáveis e saber quando usar `:=` e quando `=`.
- Citar os quatro tipos básicos e dizer por que Go não deixa você misturá-los.
- Escrever funções que recebem e devolvem valores.
- Entender por que uma função de Go pode retornar **duas** coisas ao mesmo tempo.
- Ler as mensagens de erro do compilador e saber o que ele está pedindo.
- Escrever, compilar e executar um programa completo por conta própria.

---

## Por que isso importa

O `revisor` —o programa que você vai construir durante todo o curso— precisa, desde a sua primeira
linha, de três coisas: **dados** (o nome de um serviço, sua porta, quanto demorou para responder),
**lógica** que decida algo com esses dados (está rápido ou lento?) e uma forma de **avisar quando algo
deu errado** (o serviço não respondeu?). Isso é, nessa ordem, exatamente o que esta lição ensina:
variáveis para guardar dados, funções para a lógica, e o padrão de dois valores de retorno para os erros.

Não é por acaso que Go é **compilado** e de **tipagem estática**: as duas decisões existem para que o
computador pegue seus erros **antes** de o programa rodar, em vez de um usuário descobri-los em
produção. Uma linguagem interpretada deixa você escrever `puerto = "muchos"` e só estoura quando essa
linha é executada —o que pode ser dez minutos depois, ou nunca, se esse trecho do código quase não é
usado—. Go se recusa a produzir o executável. Essa diferença é a razão de fundo de quase tudo o que você
vai ver nesta lição: por que o compilador é estrito, por que existe o «valor zero», e por que uma função
que pode falhar **tem de dizê-lo na sua assinatura**, e não como uma surpresa.

No final da lição você vai escrever o primeiro passo real do `revisor` (exercício de projeto): uma
função que recebe o nome de um serviço, sua porta e seu tempo de resposta, e devolve uma linha de
relatório classificada. Ainda sem structs, sem concorrência e sem HTTP —isso vem depois— mas já é um
código que se parece com o programa final.

---

## Os conceitos

### 2.1 O que o compilador faz

Você escreve um arquivo de texto. Esse arquivo, por si só, não faz nada: é texto. Para que o computador
o execute, algo tem de traduzi-lo para as instruções que o processador entende.

Há duas maneiras de fazer essa tradução:

| | Como funciona | Linguagens |
|---|---|---|
| **Interpretada** | Um programa lê o seu texto linha por linha e faz o que ele diz, **cada vez** que você o executa | Python, JavaScript, PHP |
| **Compilada** | Um programa traduz **todo** o seu texto **uma única vez** e guarda o resultado em um arquivo executável | **Go**, C, C++, Rust |

**Go é compilado**, e isso tem três consequências que você vai notar desde hoje:

1. **Há uma etapa antes de executar: compilar.** Se você escreveu errado o nome de uma variável, descobre
   ali — antes de o programa rodar. Em uma linguagem interpretada você descobriria quando a execução
   chegasse a essa linha, o que pode ser dez minutos depois ou nunca.
2. **O resultado é um arquivo que funciona sozinho.** Não precisa que o outro computador tenha Go
   instalado.
3. **É rápido para executar**, porque a tradução já está feita.

> [!NOTE]
> 🔧 **Observação de engenharia de software 2.1**
> O compilador é a sua primeira linha de defesa, não um obstáculo. Cada erro que ele pega é um erro que
> você não vai estar procurando de noite com o programa já entregue. Quando Go te barrar, leia a
> mensagem com calma: ele está te poupando trabalho.

### 2.2 Seu primeiro programa, linha por linha

Vamos escrever, compilar e executar um programa completo. **Todos os programas deste curso são
completos e executáveis**: nada de fragmentos que não podem ser rodados.

Prepare a pasta:

```bash
mkdir -p ~/w/curso-go/cap01 && cd ~/w/curso-go/cap01
go mod init cap01
```

**Fig. 2.1** | Um programa que imprime uma mensagem.

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

Compile-o e execute-o:

```bash
$ go run fig02_01.go
Bienvenido a Go!
```

Agora **cada linha**, porque todas têm uma razão:

**Linhas 1-2.** Comentários. Começam com `//` e o compilador os ignora por completo. Estão ali para quem
ler o código — incluindo você daqui a três semanas.

**Linha 3: `package main`.** Todo arquivo de Go pertence a um *pacote*, que é um agrupamento de código.
**O pacote `main` é especial: é o único que produz um programa executável.** Se você escrevesse
`package utilidades`, teria uma biblioteca que outros programas podem usar, mas que não pode ser
executada sozinha.

**Linha 5: `import "fmt"`.** Diz ao compilador que você vai usar código do pacote `fmt` (de *format*),
que vem com Go e contém as funções para imprimir e formatar texto. **Go não carrega nada por padrão**:
tudo o que você usar, você pede.

**Linha 7: `func main() {`.** Declara a função `main`, que é **onde começa a execução do programa**.
Quando você executa um programa Go, o sistema procura exatamente esta função. Se você a chamasse de
`principal` ou `inicio`, o programa não iniciaria.

**Linha 8: `fmt.Println("Bienvenido a Go!")`.** Chama a função `Println` que está dentro do pacote
`fmt`. O ponto se lê como «de»: *«a função `Println` **de** `fmt`»*. `Println` imprime o que você der a
ela e pula para a linha seguinte (*print line*).

**Linha 9: `}`.** Fecha o corpo da função. Em Go as chaves delimitam blocos, como em C, C++, Java ou
JavaScript.

> [!TIP]
> ✅ **Boa prática 2.1**
> Coloque sempre um comentário no início do arquivo com o nome dele e o que ele faz. Leva cinco segundos
> e poupa minutos a quem o abrir depois.

#### 2.2.1 A diferença entre `go run` e `go build`

**Fig. 2.2** | As duas formas de executar o seu programa.

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

Repare nos tamanhos: o seu texto tem **104 bytes**; o programa compilado, **1,8 megabyte**. A diferença
é que o executável traz dentro tudo o que precisa para funcionar.

> [!NOTE]
> 🚀 **Dica de portabilidade 2.1**
> Esse arquivo `fig02_01` pode ser copiado para **qualquer** computador com Linux da mesma arquitetura e
> funciona, mesmo que essa máquina não tenha Go instalado. É a razão pela qual tantas ferramentas de
> servidores são escritas em Go.

> [!TIP]
> 🧪 **Dica de teste e depuração 2.1**
> Use `go run` enquanto programa: é mais rápido e não enche a pasta de executáveis. Use `go build`
> quando o programa já estiver servindo e você quiser entregá-lo.

### 2.3 Variáveis

Uma **variável** é um espaço de memória com um nome, onde você guarda um dado que o seu programa vai usar.

**Fig. 2.3** | Declarar variáveis e imprimi-las.

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

**O operador `:=`** se lê *«declare uma variável nova e guarde isto nela»*. Go olha o valor da direita
e **deduz o tipo sozinho**: `"catalogo"` está entre aspas, então é texto; `443` não as tem e não tem
ponto, então é inteiro.

#### 2.3.1 `:=` contra `=`

Quando a variável já existe, para mudar o valor dela você usa `=` **sem** os dois-pontos:

**Fig. 2.4** | Criar e modificar uma variável.

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
> ✅ **Boa prática 2.2**
> Se você realmente precisa declarar algo que ainda não vai usar, use o identificador em branco `_`. Em
> Go, `_` significa *«isto eu descarto de propósito»*, e o compilador o aceita.

#### 2.3.2 Go se recusa a compilar se você não usa uma variável

**Fig. 2.5** | Um programa que **não compila**, de propósito.

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

**Go não permite.** Não é um aviso que você possa ignorar: é um erro que interrompe a compilação.

**Por que tão estrito?** Porque uma variável que não é usada quase sempre significa uma de duas coisas:
você errou ao escrever o nome em outra linha, ou sobrou lixo de um código que você apagou pela metade.
As duas são problemas. Go prefere te incomodar hoje a deixar código morto se acumulando para sempre.

### 2.4 Os tipos básicos

Um **tipo** é a resposta à pergunta *«que tipo de dado é este?»*. Os quatro que você vai usar o tempo
todo:

| Tipo | O que guarda | Exemplos | Valor zero |
|---|---|---|---|
| `string` | texto | `"hola"`, `"https://catalogo.example.com"` | `""` (vazio) |
| `int` | números inteiros | `42`, `-7`, `0` | `0` |
| `float64` | números com casas decimais | `3.14`, `0.142` | `0` |
| `bool` | verdadeiro ou falso | `true`, `false` | `false` |

#### 2.4.1 Go não mistura tipos

**Fig. 2.6** | Outro programa que **não compila**, de propósito.

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

Em Python isto funcionaria sem reclamação. **Em Go não compila**, e isso é uma vantagem enorme:

> [!NOTE]
> 🔧 **Observação de engenharia de software 2.2**
> Boa parte dos erros que chegam à produção é desta classe: alguém achou que uma variável tinha um
> número e ela tinha texto. Em uma linguagem interpretada isso estoura **quando um cliente está usando o
> programa**. Em Go nem chega a compilar: o erro aparece na sua máquina, não na do cliente.

#### 2.4.2 O valor zero: em Go nunca há lixo

**Fig. 2.7** | Variáveis declaradas sem valor.

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

A forma `var nombre tipo` declara uma variável **sem dar valor a ela**. Em outras linguagens isso
deixaria «lixo» —o que houvesse naquela memória— ou um estado especial de «indefinido». **Em Go ela
sempre recebe um valor válido**, chamado **valor zero**.

**E aqui aparece o `Printf`**, que é diferente do `Println`: recebe um modelo com lacunas e vai
preenchendo-as. Cada lacuna começa com `%`:

| | Para quê |
|---|---|
| `%d` | um número inteiro (*decimal*) |
| `%s` | texto (*string*) |
| `%q` | texto **com aspas** (*quoted*) — útil para ver se algo está vazio |
| `%t` | um `bool` (*true/false*) |
| `%v` | qualquer coisa, no seu formato natural (*value*) |
| `\n` | quebra de linha (`Printf` **não** pula sozinho, `Println` sim) |

> [!TIP]
> ✅ **Boa prática 2.3**
> Use `%q` quando imprimir texto que você esteja depurando. `%s` com um texto vazio não imprime nada e
> parece que a linha falhou; `%q` mostra `""` e fica claro que o texto está vazio.

### 2.5 Funções

Uma **função** é um bloco de código com nome, ao qual você pode dar dados e que pode devolver um
resultado. Você já usou duas: `fmt.Println` e `fmt.Printf`.

**Fig. 2.8** | Definir e chamar funções.

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

**Anatomia da linha 8**, que é onde está tudo:

```
func  sumar  (a int, b int)  int  {
 │      │          │          │
 │      │          │          └── lo que DEVUELVE
 │      │          └── lo que RECIBE (los parámetros), con su tipo
 │      └── el nombre
 └── palabra clave: «voy a definir una función»
```

**Linha 14: `fmt.Sprintf`.** É como o `Printf`, mas em vez de imprimir, **devolve** o texto montado. O
`S` é de *string*. É usado muitíssimo.

**Linha 19: `puerto == 443`.** O igual duplo **compara** e dá `true` ou `false`. Um único igual `=`
**atribui**. Confundi-los é clássico (veja «O que se faz errado», mais abaixo).

> [!TIP]
> ✅ **Boa prática 2.4**
> Quando dois parâmetros são do mesmo tipo você pode abreviar: `func sumar(a, b int) int`. É o idiomático
> e mais curto. Aqui escrevemos por extenso para que a estrutura fique visível.

### 2.6 Devolver dois valores: a assinatura de Go

**Isto é o mais distintivo de Go**, e você vai escrevê-lo milhares de vezes na sua carreira. Preste
atenção.

Pense em uma função que divide dois números. O que ela faz se o divisor for zero? Não pode devolver um
número, porque o resultado não existe.

Outras linguagens usam **exceções**: a função «lança» um erro e alguém mais acima o «captura» com
`try / catch`. O problema é que **não é óbvio quem o captura nem onde**, e é facílimo que ninguém o
faça e o programa morra.

**Go faz outra coisa: devolve dois valores.** O resultado, e um erro.

**Fig. 2.9** | Uma função que devolve um resultado e um possível erro.

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

**Linha 11: `(int, error)`.** Os parênteses com dois tipos significam *«esta função devolve duas
coisas»*: um inteiro e um erro.

**Linha 13: `return 0, errors.New(...)`.** Quando há problema, devolve um valor qualquer (aqui `0`, que não
vai ser usado) **e** um erro descrevendo o que aconteceu.

**Linha 15: `return a / b, nil`.** Quando tudo vai bem, devolve o resultado **e `nil`**.

🔑 **`nil` significa «nada», «vazio», «não há».** Então a pergunta `if err != nil` se lê:
*«o erro é diferente de nada?»*, ou em português normal: **«houve erro?»**

**Linhas 20-21: o padrão que define Go.**

```go
resultado, err := dividir(10, 2)
if err != nil {
    // algo salió mal: aquí se atiende
}
// si llegaste aquí, todo bien
```

Esse padrão `if err != nil` é **a linha mais escrita em toda a história de Go**. Aparece por toda parte, e
há gente que o critica por ser repetitivo.

> [!NOTE]
> 🔧 **Observação de engenharia de software 2.3**
> A crítica é verdadeira: é verboso. A vantagem é que **é impossível ignorar um erro por descuido**,
> porque ele está ali, na sua cara, na linha seguinte à chamada. Com exceções você pode esquecer um
> `catch` e só descobrir quando o programa morrer em produção. Go trocou brevidade por segurança, de
> propósito.

> [!TIP]
> 🧪 **Dica de teste e depuração 2.2**
> Repare na linha 28: diz `resultado, err = dividir(10, 0)` com `=`, não `:=`. É porque as duas
> variáveis **já existem** desde a linha 20. Se você colocasse `:=` ali, o compilador avisaria.

### 2.7 Juntando tudo

**Fig. 2.10** | Programa completo que usa tudo o que há na lição.

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

Este programa já é o antepassado direto do `revisor`: `clasificar` e `reporte` são, em miniatura, o que
na lição 3 vai se converter em um método sobre um struct `Servicio`, e na lição 6 vai rodar para
vários serviços **ao mesmo tempo**.

**Linha 21: `%-10s` e `%-5d`.** O número diz **quantos espaços de largura** reservar, e o sinal de menos
significa **alinhado à esquerda**. Assim as colunas ficam retas. Sem isso, a tabela sairia torta.

> [!TIP]
> ✅ **Boa prática 2.5**
> Quando uma chamada for muito longa, quebre-a em várias linhas como nas linhas 21-22. Go permite,
> desde que a vírgula fique no fim da linha anterior.

> [!NOTE]
> 🚀 **Dica de desempenho 2.1**
> `Sprintf` é cômodo mas não é de graça: monta um texto novo na memória a cada vez. Para cinco linhas
> não importa em absoluto. Se um dia você montar milhares, existe o `strings.Builder`, que é muito mais
> rápido. **Não otimize antes de medir**: isto é informação para depois, não algo a fazer hoje.

---

## O erro que você vai ver

Quase todas as mensagens de erro de Go têm a mesma forma:

```
<archivo>:<línea>:<columna>: <qué esperaba el compilador y qué encontró>
```

Aprender a ler essa linha é metade do trabalho. Dois exemplos que você já viu nesta lição:

```
./fig02_05.go:8:5: declared and not used: edad
```

Lê-se: *«no arquivo `fig02_05.go`, linha 8, coluna 5, você declarou `edad` e nunca a usou»*. A
correção é simples: ou você usa a variável, ou a apaga, ou —se você realmente a precisa declarada mas
ainda não a usa— a troca pelo identificador em branco `_`.

```
./fig02_06.go:8:14: cannot use "muchos" (untyped string constant) as int value in assignment
```

Lê-se: *«na linha 8, coluna 12, você tentou usar o texto `"muchos"` onde se esperava um `int`, e eu não
deixo»*. A correção também não é fácil de adivinhar se você nunca a viu: **não existe uma conversão
automática** entre `string` e `int` em uma atribuição assim. Você teria de converter explicitamente (com
`strconv`, que veremos mais adiante) ou, mais provável, perceber que misturou o tipo por engano.

E um terceiro, o mais comum da primeira semana, que aparece se você reusa `:=` sobre uma variável que já
existe:

```
./fig02_04.go:10:9: no new variables on left side of :=
```

Lê-se: *«do lado esquerdo do `:=` não há nenhuma variável nova»* — todas já existiam, então Go não sabe
o que declarar. A correção: use `=` em vez de `:=` quando você já declarou a variável antes.

**A regra geral:** o compilador de Go quase nunca diz «algo deu errado» em abstrato. Diz o arquivo, a
linha, a coluna e uma frase que —ainda que soe estranha da primeira vez— descreve exatamente o problema.
Lê-la completa, sem se assustar, resolve a maioria dos casos sem precisar pesquisar nada na internet.

## O que se faz errado

- **Escrever `println` em minúscula em vez de `Println`.** Go **diferencia maiúsculas de minúsculas**:
  `Println` e `println` são nomes diferentes (`println` é uma função interna do compilador, pensada para
  depurar o próprio Go, não o seu programa). O erro é `undefined: fmt.println`.
- **Usar `import "fmt"` e não usá-lo, ou usar `fmt.Println` sem o `import`.** O primeiro dá
  `"fmt" imported and not used`; o segundo, `undefined: fmt`. Os dois são o mesmo tipo de erro: você
  disse ao compilador algo que não coincide com o que fez.
- **Confundir `=` (atribuir) com `==` (comparar)** dentro de um `if`. Em C isso compilaria e faria algo
  inesperado (atribuir o valor e seguir em frente). **Em Go é um erro de compilação** —
  `cannot use puerto = 443 (...) as value` —, o que é uma boa notícia: te barra antes que o programa
  faça algo que você não pediu.
- **Esquecer o `\n` no `Printf`.** O programa **compila e roda**, mas tudo sai grudado em uma única
  linha. `Println` pula de linha automaticamente; `Printf` não. É um erro silencioso, não um que o
  compilador aponte.
- **Ignorar o erro com `_`**, assim: `resultado, _ := dividir(10, 0)`. **Isto também compila**, e o
  programa segue em frente com `resultado` valendo `0` como se tudo estivesse bem. É a forma mais rápida
  de criar um erro impossível de encontrar depois, porque não há nenhuma mensagem nem queda: o programa
  simplesmente segue com um dado incorreto. **Se alguma vez você descartar um erro com `_`, que seja uma
  decisão consciente e comentada, nunca um reflexo para que o compilador pare de reclamar.**

## Exercícios

### Perguntas de revisão

Responda sem ver as respostas. Elas estão no final desta seção.

**2.1** Qual é a diferença entre uma linguagem compilada e uma interpretada, e qual das duas é Go?

**2.2** Por que a função onde o programa começa tem de se chamar exatamente `main`?

**2.3** Que diferença há entre `:=` e `=`?

**2.4** Que erro Go dá se você declara uma variável e não a usa? Por que você acha que ele faz isso?

**2.5** O que é o *valor zero* e qual é o de `string`, `int` e `bool`?

**2.6** O que este programa imprime?

```go
package main

import "fmt"

func main() {
    var n int
    var s string
    fmt.Printf("[%d] [%q]\n", n, s)
}
```

**2.7** Que diferença há entre `Println`, `Printf` e `Sprintf`?

**2.8** O que significa `nil` e o que se pergunta com `if err != nil`?

**2.9** Este programa tem **três** erros que impedem que compile. Encontre-os sem rodá-lo.

```go
package main

func main() {
    nombre := "catalogo"
    puerto := 443
    puerto := 8080
    fmt.Println(nombre)
}
```

#### Respostas

**2.1** Uma linguagem **interpretada** é traduzida linha por linha a cada execução; uma **compilada** é
traduzida por completo uma única vez e produz um executável. **Go é compilado.**

**2.2** Porque é uma convenção da linguagem: ao executar um programa, Go procura a função `main` do
pacote `main` para começar. Com outro nome ele não encontra por onde iniciar.

**2.3** `:=` **declara** uma variável nova e atribui valor a ela; `=` apenas **atribui** a uma que já
existe.

**2.4** `declared and not used`. Faz isso porque uma variável sem uso quase sempre indica um erro —um
nome escrito errado ou restos de código apagado— e Go prefere te barrar a deixar código morto.

**2.5** É o valor que uma variável declarada sem valor recebe automaticamente. `string` → `""`,
`int` → `0`, `bool` → `false`. Garante que **nunca** há lixo nem «indefinido».

**2.6** `[0] [""]` — o valor zero de `int` e de `string`, e `%q` mostra as aspas.

**2.7** `Println` imprime e pula de linha. `Printf` imprime com um modelo de `%` e **não** pula
sozinho. `Sprintf` usa o mesmo modelo, mas **devolve** o texto em vez de imprimi-lo.

**2.8** `nil` significa «nada» / «vazio». `if err != nil` pergunta *«o erro não está vazio?»*, ou seja,
**«houve um erro?»**.

**2.9** Os três:
1. Falta `import "fmt"`, e se usa `fmt.Println`.
2. `puerto := 8080` usa `:=` sobre uma variável que já existe → deve ser `=`.
3. `puerto` é declarada e **nunca usada** (só se imprime `nombre`) → `declared and not used`.

### Exercícios de código

> 🔴 **Estes sete exercícios ainda não trazem solução de referência** — ela será adicionada em uma
> entrega posterior. Não se inventou código de solução para não publicar um exemplo sem executá-lo
> primeiro.

**2.10** Escreva um programa que declare o seu nome, a sua idade e se você estuda ou trabalha, e os
imprima em três linhas com `Printf`, usando o verbo correto para cada tipo.

**2.11** Escreva `func celsiusAFahrenheit(c float64) float64` e teste-a com 0, 37 e 100. A fórmula é
`f = c*9/5 + 32`. **Cuidado:** se você escrever `c*9/5` com inteiros o resultado é truncado. Por que
aqui isso não acontece?

**2.12** Escreva `func esPar(n int) bool`. Dica: o operador `%` dá o resto de uma divisão, então
`n % 2 == 0` é verdadeiro para os pares.

**2.13** Escreva `func raizCuadrada(n float64) (float64, error)` que devolva erro se `n` for negativo.
Use `math.Sqrt` (você precisa de `import "math"`). Teste-a com 16 e com -4, tratando o erro nos dois
casos.

**2.14** Pegue a **Fig. 2.10** e acrescente uma coluna que diga `SI` ou `NO` conforme a porta seja 443.
Você precisa de uma função nova e de ajustar o modelo do `Sprintf`.

**2.15 (Encontre o erro)** Cada um destes fragmentos tem um problema. Diga qual é sem compilar, e depois
compile para confirmar:

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

**2.16 (Projeto do curso)** Este é o primeiro passo do programa que você vai construir durante todo o
curso. Escreva um programa que:

1. Tenha uma função `revisar(nombre string, puerto int, ms int) string` que devolva uma linha de relatório.
2. Tenha uma função `clasificar(ms int) string`.
3. Imprima um cabeçalho e cinco serviços.
4. No final, imprima quantos dos cinco foram `"rapido"`. Você vai precisar de uma variável contador e de um `if`.

Guarde-o: na lição 3 você vai reorganizá-lo com structs.

## Como sei que consegui

- Você compilou e rodou, você mesmo e sem copiar/colar, as figuras 2.1 a 2.10 desta lição, e obteve
  exatamente a saída mostrada.
- Você consegue explicar, sem ver o texto, por que Go é uma linguagem compilada e que consequência
  prática isso tem para pegar erros antes da produção.
- Você sabe de cor quando usar `:=` e quando `=`, e que erro o compilador dá se você se enganar.
- Você consegue citar os quatro tipos básicos (`string`, `int`, `float64`, `bool`) e o valor zero de
  cada um, sem consultar a tabela.
- Você escreveu e rodou a sua própria versão do exercício **2.16** (o primeiro passo do `revisor`):
  compila, imprime o cabeçalho, as cinco linhas e a contagem de serviços rápidos.
- Você consegue explicar em duas frases por que `if err != nil` aparece em toda parte no código Go, e que
  problema evita em comparação com as exceções de outras linguagens.

Abra o [`bitacora.md`](bitacora.md) e anote **duas coisas**: o erro do compilador que mais custou a
você entender, e algo que te surpreendeu. Daqui a três semanas isso vai parecer óbvio e você não vai se
lembrar por que custou — e é justamente isso que convém ter escrito.

---

## Resumo

- Go é uma linguagem **compilada**: traduz todo o seu código uma vez e produz um executável independente.
- Todo arquivo pertence a um **pacote**; o pacote **`main`** é o único que produz um executável.
- A execução começa na função **`main`**.
- **`import`** traz pacotes; Go não carrega nada por padrão.
- **`:=`** declara e atribui; **`=`** apenas atribui.
- Go **não compila** se você declara uma variável e não a usa, nem se mistura tipos.
- Os quatro tipos básicos são **`string`**, **`int`**, **`float64`** e **`bool`**.
- O **valor zero** garante que toda variável nasce com um valor válido: `""`, `0`, `0`, `false`.
- **`Println`** imprime com quebra de linha; **`Printf`** usa modelo com `%`; **`Sprintf`** devolve o
  texto em vez de imprimi-lo.
- Uma função de Go pode devolver **vários valores**, e o padrão é devolver **`(resultado, error)`**.
- **`nil`** significa «nada». **`if err != nil`** é a forma idiomática de verificar se houve erro, e a
  linha mais escrita em Go.
- `go run` compila e executa sem deixar arquivo; `go build` deixa o executável.

---

## Para ler mais

1. **[A Tour of Go](https://go.dev/tour/)** — oficial, interativo. Faça a parte de «Basics» até onde você
   chegou nesta lição.
2. **[Go by Example — Variables](https://gobyexample.com/variables)** e
   **[Go by Example — Multiple Return Values](https://gobyexample.com/multiple-return-values)** — o
   mesmo padrão `(resultado, error)` com programas mínimos que rodam.
3. **[Effective Go](https://go.dev/doc/effective_go)** — ainda não é preciso lê-lo por completo, mas
   consulte-o se algo desta lição pareceu uma regra arbitrária: ali está o porquê.

Se algo não ficou claro: leia a mensagem de erro completa (Go costuma dizer a linha, a coluna e o que
esperava), procure o conceito no Go by Example, e anote a dúvida no diário de bordo **mesmo que não a
resolva**. Uma dúvida escrita pode ser resolvida depois; uma esquecida, não.

### Termos desta lição

| | |
|---|---|
| **compilador** | programa que traduz código-fonte em instruções de máquina |
| **erro de compilação** | erro que impede produzir o executável; é detectado antes de executar |
| **função** | bloco de código com nome que recebe parâmetros e pode devolver valores |
| **identificador em branco (`_`)** | símbolo que descarta um valor de propósito |
| **`int`, `float64`, `string`, `bool`** | os quatro tipos básicos |
| **linguagem compilada / interpretada** | traduz tudo uma vez / linha por linha a cada execução |
| **`main` (função)** | ponto de entrada do programa |
| **`main` (pacote)** | o único pacote que produz um executável |
| **`nil`** | ausência de valor |
| **pacote** | agrupamento de código relacionado |
| **parâmetro** | dado que uma função recebe |
| **`Printf` / `Println` / `Sprintf`** | imprimir com modelo / imprimir com quebra de linha / devolver texto |
| **tipo** | tipo de dado que uma variável pode guardar |
| **valor zero** | valor válido que toda variável declarada sem valor recebe |
| **variável** | espaço de memória com nome que guarda um dado |
| **verbo de formatação (`%d`, `%s`, `%q`, `%t`, `%v`)** | lacuna em um modelo de `Printf` |

---

**Anterior:** [Lição 1 — Instalar Go](01-instalacion.md) ·
**Próxima:** [Lição 3 — Structs, métodos, erros e interfaces](03-errores-interfaces.md)
