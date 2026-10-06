# Lição 5 — Módulos e testes

**Duração:** 90 minutos, ou 2 sessões de 45.

**Ao terminar, você vai conseguir:**

- Explicar o que o `go.mod` resolve e por que o `revisor` não precisa de `go.sum`.
- Reorganizar um programa de um arquivo em pacotes sob `internal/`, e explicar o que essa pasta te dá
  que uma pasta normal não dá.
- Escrever testes com tabela de casos, sem nenhuma biblioteca externa, e ler o que eles reportam quando falham.
- Rodar `go test` com `-v`, `-run`, `-cover` e `-race`, e explicar o que cada flag mede.
- Reconhecer a armadilha de "zero testes" que parece idêntica a "todos passaram".
- Decidir, com um número real de cobertura à frente, que parte dessa lacuna importa corrigir e qual não.

---

## Por que isso importa

Até a lição 4 você tinha um único arquivo, `main.go`, com tudo dentro: os structs `Servicio` e
`Estado`, a interface `Revisor`, a função que monta o relatório. Funcionava, e funcionava bem —mas um
único arquivo tem um teto. Assim que você quer **testar** uma peça sem executar o programa completo
(sem tocar na rede, sem ler um arquivo real), um único `main.go` não deixa: tudo está misturado com tudo.

Esta é a lição em que o `revisor` deixa de ser um exercício de um arquivo e se converte em **um
projeto de verdade**: vários pacotes, cada um com uma responsabilidade, e uma suíte de testes que
demonstra que cada peça faz o que diz sem precisar das demais. É o mesmo salto que você deu na
lição 1 de "um programa que cabe na cabeça" para "um programa que vive em uma pasta com `go.mod`" —
agora esse `go.mod` vai organizar mais de um arquivo.

🔑 **E não é um capricho de organização.** Um teste que precisa de rede, de um arquivo em disco ou de um
servidor no ar para rodar é um teste lento, frágil e que ninguém roda com frequência. Separar em pacotes
é o que te permite escrever testes que rodam em milissegundos, sem tocar em nada externo — e é isso
que faz com que você de fato os rode, a cada mudança, não só quando se lembra.

---

## Os conceitos

### 5.1 `go.mod`, e por que este projeto não tem `go.sum`

Você já usou `go mod init` na lição 1 para o programa `hola`. O `go.mod` do `revisor` é igualmente
simples:

```
module github.com/habil/revisor

go 1.27
```

Duas linhas: o nome do módulo (é assim que outro programa o importaria, se algum dia você publicar algum
dos seus pacotes) e a versão mínima de Go de que ele precisa.

**Se você procurar tutoriais de módulos na internet, quase todos vão mencionar o `go.sum` logo de cara** — o
arquivo com os hashes criptográficos de cada dependência externa, para que ninguém te empurre uma versão
diferente de um pacote que você usa. O `revisor` **não tem `go.sum`**, e não é um erro nem um descuido:

```bash
$ ls go.sum
ls: go.sum: No such file or directory
```

**Ele não existe porque o `revisor` não importa nem um único pacote externo.** Revise os `import` de qualquer
arquivo do projeto e você vai encontrar unicamente pacotes da biblioteca padrão: `net/http`,
`encoding/json`, `context`, `sync`, `flag`, `os`, `time`, `strings`, `sort`. É uma decisão deliberada,
não uma limitação: a lição 0 já avisava —*"aprenda a biblioteca padrão antes de qualquer
framework"*— e o `revisor` é a prova de que a padrão basta para um programa completo, com
concorrência, HTTP, JSON e testes, sem adicionar uma única dependência de terceiros. Se algum dia você adicionar
uma (por exemplo, um cliente de YAML de verdade), nesse momento o `go get` vai criar o `go.sum` para você, e
os dois arquivos —`go.mod` e `go.sum`— vão para o repositório.

Para ver de verdade o que teria mudado, fiz o teste em um projeto à parte (não no `revisor`, que
continua sem dependências): um `go get` real de um pacote externo pequeno, `gopkg.in/yaml.v3`.

```
$ go get gopkg.in/yaml.v3
go: downloading gopkg.in/yaml.v3 v3.0.1
go: added gopkg.in/yaml.v3 v3.0.1

$ cat go.mod
module demo
go 1.27.1
require gopkg.in/yaml.v3 v3.0.1 // indirect

$ cat go.sum
gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
```

Isso é o que aparece assim que você adiciona **uma única** dependência externa: o `go.mod` ganha uma linha
`require`, e o `go.sum` nasce com os hashes criptográficos (os textos longos que começam com `h1:`) dessa
dependência e das dependências dela (`check.v1` é uma dependência indireta de `yaml.v3`, não algo que
você pediu). Esses hashes são o que faz com que, se alguém tentasse te empurrar uma versão diferente do
pacote com o mesmo nome e versão, o `go build` se recuse a compilar — é uma garantia de integridade,
não só um registro. O `revisor` não os tem porque não precisa deles: zero dependências externas, zero
superfície para essa classe de risco.

### 5.2 Pacotes por responsabilidade, não por camada

Foi assim que o `revisor` ficou organizado nesta lição:

```
revisor/
  go.mod
  cmd/
    revisor/            → package main: arranca, parsea banderas, imprime
    servidor-demo/       → package main: un servidor de prueba, no parte del programa
  internal/
    servicio/           → package servicio: los tipos (Servicio, Estado) y sus métodos
    config/             → package config: lee y valida el archivo de servicios
    revisar/            → package revisar: la interfaz Revisor y quien la implementa
    reporte/            → package reporte: convierte estados en tabla o en JSON
```

**Por responsabilidade, não por camada.** Não há um pacote `modelos` com todos os structs do programa nem
um pacote `utilidades` com funções soltas: cada pacote tem uma única pergunta que sabe responder.
`servicio` sabe o que é um serviço e um estado. `config` sabe ler a configuração. `revisar` sabe
consultar. `reporte` sabe imprimir. Se amanhã você mudar a aparência da tabela, toca em um arquivo, não em cinco.

🔑 **`internal/` é uma regra do compilador, não uma convenção de boas maneiras.** Qualquer pacote
que viva sob uma pasta chamada `internal/` só pode ser importado por código que esteja **dentro do
mesmo módulo**, em qualquer nível acima desse `internal/`. Comprove: se outro módulo de Go —qualquer um,
não só um seu— tentar `import "github.com/habil/revisor/internal/servicio"`, o compilador se recusa
a compilar, com uma mensagem explícita de que esse pacote é interno. Não é uma recomendação que você possa
ignorar sob pressão: é uma restrição real, a mesma classe de garantia que a maiúscula te dá dentro
de um struct (lição 3), mas no nível do pacote inteiro.

⚠️ **O que NÃO fizemos, de propósito: um pacote `utils`, `helpers` ou `common`.** É o antipadrão mais
repetido em projetos reais: alguém cria uma pasta para "coisas que não se encaixam em outro lugar", e essa
pasta cresce sem limite até que ninguém sabe o que há dentro nem por quê. Cada vez que você sentir a tentação
de colocar algo em um pacote assim, pergunte-se de que **responsabilidade** é essa função, e coloque-a no
pacote dono dessa responsabilidade — ou, se de verdade ela não se encaixa em nenhum, é um sinal de que falta
nomear um conceito novo, não de que falta uma gaveta de tralhas.

### 5.3 O arquivo de configuração, e os seus erros reais

O `config` desta lição lê um formato de texto simples, uma linha por serviço:

```
nombre  url  [tiempo-limite]
```

```
# las lineas que empiezan con # se ignoran
catalogo   https://catalogo.interno.mx
pagos      https://pagos.interno.mx      500ms
```

Quem faz isso é a função `Interpretar`, que **nunca lê um arquivo**: recebe os bytes já lidos, e por
isso pode ser testada com cem variantes de conteúdo sem criar um único arquivo temporário (`Cargar`, que sim
toca no disco, é uma camada fina por cima que só lê o arquivo e passa o conteúdo para `Interpretar`).
Cada linha mal escrita produz um erro real, com o número da linha e o que se esperava — testado, não
suposto:

```
$ (línea: "catalogo")
prueba.txt:1: esperaba «nombre url [tiempo]», hay 1 campo(s)

$ (línea: "catalogo catalogo.interno.mx")
prueba.txt:1: la URL "catalogo.interno.mx" debe empezar con http:// o https://

$ (línea: "catalogo https://a.mx nombas")
prueba.txt:1: tiempo limite "nombas" invalido (usa 500ms, 2s, 1m): time: invalid duration "nombas"

$ (dos líneas con el mismo nombre "catalogo")
prueba.txt:2: el nombre "catalogo" ya estaba en la linea 1
```

Repare no último: **ele encapsula o erro de `time.ParseDuration` com `%w`** (lição 3) em vez de
inventar o próprio texto — assim, se algum dia você precisar distinguir programaticamente "duração inválida"
de outro tipo de erro com `errors.As`, a informação original continua ali.

### 5.4 Testes: sem bibliotecas, com tabela de casos

É assim que fica um teste real do pacote `servicio` (o arquivo completo vive em
`programas/revisor/internal/servicio/servicio_test.go`):

<!-- verificar:extracto:internal/servicio/servicio_test.go -->
```go
func TestTimeoutEfectivo(t *testing.T) {
	casos := []struct {
		nombre  string
		timeout time.Duration
		quiere  time.Duration
	}{
		{"declarado", 5 * time.Second, 5 * time.Second},
		{"cero (valor por omision)", 0, TimeoutPorOmision},
		{"negativo (dato corrupto, no debe pasar)", -1 * time.Second, TimeoutPorOmision},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			s := Servicio{Timeout: c.timeout}
			if got := s.TimeoutEfectivo(); got != c.quiere {
				t.Errorf("TimeoutEfectivo() = %v, quería %v", got, c.quiere)
			}
		})
	}
}
```

**Não há `assert`, nem `expect`, nem nenhuma biblioteca de asserções — e é de propósito.** Go compara
valores com um `if` normal e reporta com `t.Errorf`, e é você mesmo quem escreve o que esperava e o que obteve.
No começo parece mais verboso que um `assert.Equal(t, esperado, obtenido)` de outras linguagens; o
ganho é que a mensagem de falha é controlada por você, em vez de herdar o formato genérico de uma
biblioteca, e que não é preciso instalar nem aprender nada extra para escrever o teste mais simples.

**`t.Run` dá nome a cada caso da tabela**, e isso importa quando algo falha: em vez de um genérico
"`TestTimeoutEfectivo` falhou", o relatório diz exatamente qual dos três casos foi:

```
$ go test ./internal/servicio/... -run TestTimeoutEfectivo -v
=== RUN   TestTimeoutEfectivo
=== RUN   TestTimeoutEfectivo/declarado
=== RUN   TestTimeoutEfectivo/cero_(valor_por_omision)
=== RUN   TestTimeoutEfectivo/negativo_(dato_corrupto,_no_debe_pasar)
--- PASS: TestTimeoutEfectivo (0.00s)
    --- PASS: TestTimeoutEfectivo/declarado (0.00s)
    --- PASS: TestTimeoutEfectivo/cero_(valor_por_omision) (0.00s)
    --- PASS: TestTimeoutEfectivo/negativo_(dato_corrupto,_no_debe_pasar) (0.00s)
PASS
ok  	github.com/habil/revisor/internal/servicio	0.382s
```

Essa saída é real: eu a rodei sobre o código deste mesmo projeto antes de escrever esta linha.
**Adicionar um caso novo à tabela é adicionar uma linha ao slice `casos`** — não uma função nova, não
repetir o corpo do teste. É a forma idiomática de testar uma função com muitas entradas em Go, e
você vai usá-la em cada pacote daqui em diante.

### 5.5 Testar os erros, não só os acertos

Uma tabela de casos também serve para testar que algo **falha como deve**, não só que funciona (versão
abreviada aqui, com 3 dos 6 casos reais e struct literal posicional em vez de com nome de campo,
para que caiba; a completa vive em `internal/config/config_test.go`):

<!-- verificar:fragmento -->
```go
func TestInterpretar_casosDeError(t *testing.T) {
	casos := []struct {
		nombre   string
		entrada  string
		contexto string // fragmento que el mensaje de error debe contener
	}{
		{"nombre repetido", "catalogo https://a.mx\ncatalogo https://b.mx\n", "ya estaba en la linea"},
		{"url sin esquema", "catalogo catalogo.interno.mx\n", "debe empezar con http"},
		{"tiempo limite invalido", "catalogo https://a.mx nombas\n", "invalido"},
		// ...
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			_, err := Interpretar([]byte(c.entrada), "prueba.txt")
			if err == nil {
				t.Fatalf("Interpretar() no devolvió error, y debía contener %q", c.contexto)
			}
			if !strings.Contains(err.Error(), c.contexto) {
				t.Errorf("Interpretar() error = %q, quería que contuviera %q", err.Error(), c.contexto)
			}
		})
	}
}
```

**`t.Fatalf` em vez de `t.Errorf`** no primeiro `if`: se `Interpretar` não retornou erro quando devia,
continuar verificando `err.Error()` na linha seguinte causaria um panic (`err` seria `nil`). `Fatalf` interrompe
esse teste em particular ali mesmo; `Errorf` deixa o teste continuar rodando e acumular mais falhas antes
de reportar. A regra prática: use `Fatalf` quando continuar não faz sentido sem o que você acabou de
verificar, `Errorf` quando faz.

### 5.6 Os comandos que você vai usar sempre

```bash
go test ./...                          # todo el proyecto
go test ./internal/config/... -v       # verboso, un paquete
go test ./... -run TestInterpretar     # solo las pruebas cuyo nombre haga match
go test ./... -race                    # detector de carreras (se explica a fondo en la lección 6)
go test ./... -cover                   # porcentaje de líneas ejercitadas por las pruebas
```

Rodados de verdade sobre o `revisor`, em 30-set-2026:

```
$ go test ./internal/servicio/... ./internal/config/... -cover
ok  	github.com/habil/revisor/internal/servicio	0.195s	coverage: 92.3% of statements
ok  	github.com/habil/revisor/internal/config	0.192s	coverage: 97.4% of statements
```

### 5.7 O que fazer com um número de cobertura

**92.3% e 97.4% não são metas, são pontos de partida para uma pergunta: o que são esses 8% e esses 3% que não foram
executados, e eles me importam?** Com `go test -coverprofile` você pode ver exatamente quais linhas ficaram sem
ser tocadas:

```bash
go test ./internal/servicio/... -coverprofile=/tmp/cobertura.out
go tool cover -func=/tmp/cobertura.out
```

No `revisor`, a lacuna de `servicio` é o ramo de `Motivo()` que monta a mensagem quando o código não
é nem 0 nem um sucesso com um texto específico (`fmt.Sprintf("codigo %d", ...)` para um 404 sem mais
contexto) — um ramo que os testes existentes não exercitam com esse código exato. **A decisão correta
não é perseguir os 100%** preenchendo cada ramo com um teste forçado que não ensina nada novo: é olhar a
lacuna, decidir se importa (aqui, um pouco: você adicionaria um caso com 404 na tabela) e anotá-la, em vez de
fingir que ela não existe.

🔴 **E uma armadilha real, medida neste mesmo projeto: a cobertura por pacote pode subestimar uma
função central sem te avisar.** Rodado sozinho, o pacote `revisar` do `revisor` (que você vai conhecer a
fundo na lição 6) reporta:

```
$ go test ./internal/revisar/... -cover
ok  	github.com/habil/revisor/internal/revisar	1.33s	coverage: 61.9% of statements
```

**61.9% soa como se mais de um terço desse pacote nunca fosse executado em nenhum teste — e é falso.** A sua
função mais importante, `Revisar` (a que de verdade fala HTTP), está sim bem testada: só que o teste
que a exercita de verdade —`TestEjecutar_reportaOKyFalla`, com um servidor `httptest` real, na lição
7— vive no pacote `cmd/revisor`, não em `internal/revisar`. `go test ./internal/revisar/...` **só
conta o que os testes DESSE PACOTE exercitam**; não vê o que um teste de outro pacote percorre pelo
caminho, ainda que passe pelo mesmo código. Com `-coverpkg`, que pede ao Go para medir a cobertura de um
pacote contando **todos** os testes do projeto, não só os dele:

```
$ go test ./... -coverpkg=./... -coverprofile=/tmp/cov.out
$ go tool cover -func=/tmp/cov.out | grep 'revisar.go.*Revisar'
github.com/habil/revisor/internal/revisar/revisar.go:48:  Revisar    86.4%

$ go tool cover -func=/tmp/cov.out | tail -1
total:                                                    (statements)    81.5%
```

**86.4% para `Revisar`, 81.5% para o projeto completo — não 61.9%.** A lição não é "ignore o número
por pacote": é que um número de cobertura sempre responde a uma pergunta implícita —cobertura de quê,
medida contra os testes de onde?— e `go test ./paquete/... -cover` cala essa segunda metade da
pergunta. Antes de decidir que algo "não está testado" por um número baixo, rode `-coverpkg=./...` sobre
todo o projeto e compare.

### 5.8 Provocar uma falha, para saber que o teste serve

Um teste que você nunca viu falhar é um teste do qual você não sabe se funciona — pode estar
comparando duas coisas que são sempre iguais por acidente. Quebre-o de propósito: mude
`TimeoutEfectivo()` para que retorne `s.Timeout` sempre, sem o `if`, e rode o teste:

```
$ go test ./internal/servicio/... -run TestTimeoutEfectivo -v
=== RUN   TestTimeoutEfectivo/cero_(valor_por_omision)
    servicio_test.go:64: TimeoutEfectivo() = 0s, quería 2s
--- FAIL: TestTimeoutEfectivo (0.00s)
    --- FAIL: TestTimeoutEfectivo/cero_(valor_por_omision) (0.00s)
FAIL
```

**Esse `FAIL`, com o caso exato que quebrou e o valor que obteve contra o que esperava, é a prova
de que o teste funciona.** Devolva o `if` e confirme que ele volta a ficar verde antes de seguir.

### 5.9 Projetar para poder testar: separar a lógica da E/S

Repare em algo que já mencionamos de passagem na seção 5.3 e que merece o seu próprio espaço: `config`
tem **duas** funções, não uma.

<!-- verificar:extracto:internal/config/config.go -->
```go
func Cargar(ruta string) ([]servicio.Servicio, error) {
	datos, err := os.ReadFile(ruta)
	if err != nil {
		return nil, fmt.Errorf("leyendo la configuracion: %w", err)
	}
	return Interpretar(datos, ruta)
}
```

`Interpretar` (apenas a assinatura; o corpo completo, com toda a lógica de parsing, está na seção
5.3, com os seus erros reais):

<!-- verificar:fragmento -->
```go
func Interpretar(datos []byte, origen string) ([]servicio.Servicio, error) {
	// ... toda la lógica de parseo, sin tocar el disco ...
}
```

**Se houvesse uma única função `Cargar(ruta string)` que lesse o arquivo e fizesse o parsing de tudo junto**, cada
teste de "linha mal escrita", "nome repetido" ou "tempo limite inválido" teria de começar criando
um arquivo temporário em disco com `os.CreateTemp`, escrever nele o conteúdo do caso, passar o caminho, e
apagá-lo no final. Funciona, mas é lento (toca no sistema de arquivos real) e suja o teste com
código que não tem nada a ver com o que de fato está sendo testado: **se a sua configuração é
interpretada bem ou mal**, não se você sabe criar arquivos temporários.

Separando **a parte que decide** (`Interpretar`, uma função pura: mesmos bytes de entrada, mesmo
resultado sempre, sem acesso a nada externo) da **parte que obtém os bytes** (`Cargar`, a única
que toca no disco), os nove testes da seção 5.5 rodam em microssegundos e sem criar um único arquivo.
`Cargar` em si quase não precisa de testes próprios: só verificar que retorna erro se o arquivo não
existe (exercício 6), porque toda a lógica interessante já está em `Interpretar` e já está testada.

🔑 **A regra geral, útil muito além deste projeto:** quando uma função é difícil de testar,
quase sempre é porque mistura "decidir algo" com "tocar no mundo exterior" (um arquivo, a rede, o
relógio). Separá-las não é uma regra de estilo: é o que determina se você vai poder escrever o teste em
três linhas ou em vinte.

### 5.10 Benchmarks: medir, não adivinhar

Além de `Test...`, Go reconhece funções `Benchmark...` que medem quanto tempo o seu código leva, não se ele é
correto:

<!-- verificar:extracto:internal/config/config_bench_test.go -->
```go
func BenchmarkInterpretar(b *testing.B) {
	entrada := []byte(`
catalogo   https://catalogo.interno.mx
pagos      https://pagos.interno.mx      500ms
inventario https://inventario.interno.mx 1s
reportes   https://reportes.interno.mx
`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Interpretar(entrada, "bench.txt"); err != nil {
			b.Fatal(err)
		}
	}
}
```

`b.N` não é você quem escolhe: Go roda o laço com valores de `N` cada vez maiores até que a medição seja
estável, e reporta o tempo por operação. Rodado de verdade sobre o `revisor` (Apple M5, 200,000
repetições):

```
$ go test ./internal/config/... -bench=. -run '^$'
goos: darwin
goarch: arm64
pkg: github.com/habil/revisor/internal/config
cpu: Apple M5
BenchmarkInterpretar-10    	  200000	       432.5 ns/op
PASS
```

**`-run '^$'`** diz ao `go test` que não rode nenhum teste normal (uma expressão regular que não
corresponde a nenhum nome), para que o relatório do benchmark não se misture com o dos testes. O
número —432.5 nanossegundos por chamada, nesta máquina, neste momento— não é para ser memorizado: é para
compará-lo **com ele mesmo** depois de uma mudança. Se amanhã você reescrever `Interpretar` e o benchmark
subir para 4,000 ns/op, você tem um sinal objetivo de que algo ficou mais lento, sem precisar opinar a
respeito.

### 5.11 `go vet`: o que encontra o que compila mas está errado

`go test` te diz se a sua lógica faz o que você esperava. **`go vet` te diz se o seu código tem um erro que
o compilador não pega porque, tecnicamente, ele é válido** — mas quase certamente não é o que você quis
escrever. O caso mais comum é um verbo de formatação (lição 2) que não corresponde ao tipo do argumento:

<!-- verificar:ejemplo:ejemplos/05-vet-printf -->
```go
puerto := 443
fmt.Printf("servicio %s en el puerto %s\n", nombre, puerto) // %s para un int
```

Isto **compila** e **roda**, sem panic nem erro — e produz uma saída quebrada (programa completo em
`programas/revisor/ejemplos/05-vet-printf/main.go`):

```
$ go run ./ejemplos/05-vet-printf/
servicio catalogo en el puerto %!s(int=443)
```

E o `go vet`, sobre esse mesmo arquivo, detecta o problema:

```
$ go vet ./ejemplos/05-vet-printf/
ejemplos/05-vet-printf/main.go:11:39: fmt.Printf format %s has arg puerto of wrong type int
```

O `go vet` detecta, porque analisa a string de formato contra os tipos reais dos argumentos, um
passo que o compilador de Go não dá por si só. Rode `go vet ./...` junto com `go test ./...` como
rotina: a maioria dos editores com a extensão de Go (lição 1) já faz isso por você enquanto você
escreve, sublinhando o problema antes que você chegue a executar qualquer coisa.

### 5.12 Testes de fronteira: onde os bugs vivem de verdade

`Estado.OK()` decide que um código está bem se cai **entre 200 e 299**. É tentador testá-lo com um
caso "óbvio" (200) e um caso "óbvio" de falha (500) e dar por encerrado. A tabela real do `revisor`
testa também os dois valores que estão **bem na borda do intervalo**:

<!-- verificar:fragmento -->
```go
{"199, justo debajo del rango", Estado{Codigo: 199}, false},
{"300, justo arriba del rango", Estado{Codigo: 300}, false},
```

```
$ go test ./internal/servicio/... -run TestEstadoOK -v
=== RUN   TestEstadoOK/199,_justo_debajo_del_rango
=== RUN   TestEstadoOK/300,_justo_arriba_del_rango
--- PASS: TestEstadoOK (0.00s)
    --- PASS: TestEstadoOK/199,_justo_debajo_del_rango (0.00s)
    --- PASS: TestEstadoOK/300,_justo_arriba_del_rango (0.00s)
```

**Por que isso importa, se 199 e 300 "obviamente" não estão OK?** Porque um erro de um único caractere na
condição —`>=` em vez de `>`, ou `<=` em vez de `<`— é exatamente o tipo de bug que um caso "óbvio"
nunca detecta, e que um caso de fronteira detecta sempre. Se alguém mudasse `OK()` para
`e.Codigo >= 200 && e.Codigo <= 300` (incluindo o 300 por engano), os casos 200/299/404/500 continuariam
passando igualmente bem — **só o caso de 300 o denunciaria.** Essa é a razão de fundo por trás de "teste
as bordas, não só o centro": as bordas são onde os erros de comparação se escondem, e são
invisíveis para qualquer teste que só use valores muito dentro ou muito fora do intervalo.

---

## O erro que você vai ver

| O sintoma | Mensagem literal | O que acontece e o que fazer |
|---|---|---|
| Import de um pacote `internal` alheio | `use of internal package github.com/habil/revisor/internal/servicio not allowed` | O compilador impede importar algo sob `internal/` de fora do módulo. Não é uma permissão que você possa dar: é preciso expor o tipo a partir de um pacote público se de verdade for necessário |
| Serviço duplicado na configuração | `prueba.txt:2: el nombre "catalogo" ya estaba en la linea 1` | Dois serviços com o mesmo nome; corrija o arquivo de configuração |
| URL sem esquema | `prueba.txt:1: la URL "catalogo.interno.mx" debe empezar con http:// o https://` | Falta `http://` ou `https://` no início da URL |
| Tempo limite mal escrito | `prueba.txt:1: tiempo limite "nombas" invalido (usa 500ms, 2s, 1m): time: invalid duration "nombas"` | O formato de duração de Go não é livre: use um número seguido de uma unidade (`ms`, `s`, `m`, `h`) |
| Um teste que você quebrou de propósito | `servicio_test.go:64: TimeoutEfectivo() = 0s, quería 2s` | Veja a seção 5.8: é assim que fica um `t.Errorf` apontando exatamente qual caso falhou e por quê |
| Zero testes em um pacote | `?   	github.com/habil/revisor/cmd/servidor-demo	[no test files]` | Não é uma falha: o `go test` avisa explicitamente quando um pacote não tem nenhum arquivo `_test.go`, em vez de fingir que rodou algo |

**E o mais importante de aprender a ler**, porque não é um erro, mas a ausência de um:

```
$ go test ./...
ok  	github.com/habil/revisor/internal/vacio	0.001s
```

Se `internal/vacio` não tivesse **nenhuma** função `Test...`, essa linha ficaria exatamente igual:
`ok`, em verde, sem nenhuma marca de que nada foi executado. A única forma de distinguir "tudo passou, de
verdade" de "não havia nada para rodar" é olhar a contagem com `-v` (que imprime cada `RUN`) ou, melhor,
nunca confiar em um pacote que não tem nenhum arquivo `_test.go` — isso o `go test` diz, como na
linha da tabela acima.

---

## O que se faz errado

- **Criar um pacote `utils`, `helpers` ou `common`.** Você já viu isso na seção 5.2: é o antipadrão mais
  repetido e o que mais rápido perde o seu propósito. Nomeie a responsabilidade, não o fato de que "não
  se encaixa em outro lugar".
- **Confundir "compilou" com "os testes passaram".** `go build` verifica que o código é válido; não
  executa nem um único teste. São dois comandos diferentes com duas perguntas diferentes.
- **Ler o código de saída do `go test` em vez da contagem.** Um pacote sem arquivos de teste sai com
  código 0 (sucesso), igual a um com 50 testes que de fato passaram. O código de saída responde "algo
  falhou?", não "algo foi testado?" — para a segunda pergunta é preciso ler a saída, não só o `$?`.
- **Perseguir 100% de cobertura como se fosse o objetivo.** Uma cobertura alta com asserções fracas
  (verificar que algo não quebra, sem verificar o que retorna) dá um número bonito e um teste que não
  detecta quase nada. O número é um guia de onde olhar, não uma meta em si mesma — a seção 5.7 mostra
  isso com uma lacuna real do próprio `revisor`.
- **Escrever um teste e nunca vê-lo falhar.** Se você nunca quebrou o código de propósito para confirmar
  que o teste fica vermelho, não sabe se esse teste testa alguma coisa. A seção 5.8 fez isso com dados reais
  do projeto: faça você também com pelo menos um teste seu antes de dá-lo por bom.

---

## Exercícios

1. Clone a estrutura de pacotes da seção 5.2 (`internal/servicio`, `internal/config`) para a sua
   própria cópia do `revisor`, movendo o código que você já tinha das lições 2 a 4.
2. Escreva a tabela de casos de `TestEtiqueta` para o método `Etiqueta()` de `Servicio`, com pelo menos um
   caso de nome normal e um de struct vazio.
3. Adicione um caso a `TestInterpretar_casosDeError` para uma linha com **quatro** campos (mais do que os três
   que o formato permite). Verifique a mensagem exata contra o código de `config.go`.
4. Rode `go test ./... -cover` sobre a sua cópia e anote a porcentagem de cada pacote. Escolha **uma** lacuna
   de cobertura e decida, por escrito no seu diário de bordo, se te importa fechá-la e por quê.
5. (Como na seção 5.8) Quebre de propósito uma função que você já testou, rode o teste, leia o
   `FAIL` completo, e conserte. Cole os dois resultados —o vermelho e o verde— no seu diário de bordo.
6. (Um pouco mais difícil) Escreva `TestCargar_archivoInexistente`, que confirme que `Cargar` (não
   `Interpretar`) retorna erro quando o caminho não existe. Dica: você não precisa criar nenhum arquivo para
   este teste, só passar um caminho que você sabe que não existe.

### Soluções

1 e 2 não têm solução de referência única: depende de como você tinha organizado o seu próprio programa da
lição 4. Compare o seu resultado com o código real de `programas/revisor/internal/servicio/servicio.go` do
projeto deste curso.

3. Com `"catalogo https://a.mx 500ms extra\n"`, a mensagem esperada é
   `prueba.txt:1: esperaba «nombre url [tiempo]», hay 4 campo(s)` — o mesmo caminho de código que já
   trata "faltam campos", porque `len(campos) > 3` cobre ambos os casos com uma única verificação.

4. Não há uma resposta única: o que importa é que a decisão fique escrita com a sua razão, não a
   porcentagem em si.

5. Veja a seção 5.8 completa: o padrão é sempre "mude o código, rode o teste, leia o `FAIL`
   com o caso exato, conserte, rode de novo".

6. Assim:

   <!-- verificar:extracto:internal/config/config_test.go -->
   ```go
   func TestCargar_archivoInexistente(t *testing.T) {
   	_, err := Cargar("/ruta/que/no/existe.txt")
   	if err == nil {
   		t.Fatal("Cargar() no devolvió error con una ruta inexistente")
   	}
   }
   ```

   Este teste existe de fato no projeto real (`config_test.go`) e passa porque `os.ReadFile`, dentro de
   `Cargar`, retorna um erro de sistema operacional que `Cargar` encapsula com `%w` antes de propagá-lo.

---

## Como sei que consegui

- [ ] O meu `revisor` está organizado em pacotes sob `internal/`, cada um com uma única responsabilidade.
- [ ] `go test ./...` roda e **a contagem, não só a cor**, me diz quantos testes foram executados.
- [ ] Escrevi pelo menos uma tabela de casos com `t.Run`, e sei ler qual caso falhou quando algo quebra.
- [ ] Rodei `go test -cover` e consigo dizer que porcentagem saiu e o que há na lacuna.
- [ ] Quebrei uma função de propósito, vi o `FAIL` exato, e a consertei — está no meu diário de bordo.
- [ ] Sei explicar por que o `revisor` não tem `go.sum` e o que o geraria se algum dia ele precisasse.
- [ ] Sei por que não deveria criar um pacote `utils`.

---

## Para ler mais

1. [Writing tests](https://go.dev/doc/tutorial/add-a-test) — o tutorial oficial de testes, com tabela de
   casos incluída.
2. [Go Wiki: TableDrivenTests](https://go.dev/wiki/TableDrivenTests) — a referência canônica do padrão
   usado em toda esta lição.
3. [Package internal](https://go.dev/doc/go1.4#internalpackages) — o anúncio original da regra de
   `internal/`, direto das notas de versão de Go.
4. [Go Blog: The Cover Story](https://go.dev/blog/cover) — como o `go test -cover` funciona por dentro, e
   por que um número alto nem sempre significa testes bons.
5. [Liberar sem medo na infraestrutura própria](https://www.habil.mx/pt/blog/cicd-devsecops-infraestrutura-propria/) — artigo sobre um caminho até a produção em que os testes, entre outros portões, podem interromper uma liberação.

### Termos desta lição

| Termo | O que significa |
|---|---|
| `go.sum` | arquivo com os hashes criptográficos das dependências externas; o `revisor` não o tem porque não usa nenhuma |
| `internal/` | pasta especial que o compilador de Go impede de importar de fora do módulo |
| tabela de casos | padrão de teste em que uma lista de entradas e saídas esperadas é percorrida com um único corpo de teste |
| `t.Run` | executa um subcaso com o seu próprio nome, para que o relatório diga exatamente qual falhou |
| cobertura | porcentagem de linhas do código que os testes exercitaram ao rodar |

---

**Anterior:** [Lição 4 — Coleções](04-colecciones.md) ·
**Próxima:** [Lição 6 — Concorrência](06-concurrencia.md)
