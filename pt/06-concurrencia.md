# Lição 6 — Concorrência

**Duração:** 2 sessões de 60-90 minutos. É a razão pela qual Go existe, e a única lição do curso
que é dada em etapas: goroutines, depois canais e `WaitGroup`, depois `context` e o limite de paralelismo
aplicado ao `revisor`. Em uma única passada não assenta — confirma isso a mesma ordem usada pelo curso de Go com
mais alunos do mundo.

**Ao terminar, você vai conseguir:**

- Lançar uma goroutine e explicar, com uma prova real, por que o programa pode terminar antes que ela
  termine.
- Usar `sync.WaitGroup` para esperar um grupo de goroutines, e reconhecer a mensagem exata que Go dá
  quando o `Add`/`Done` não batem.
- Explicar por que o `revisor` distribui resultados por um canal em vez de escrever em um map compartilhado, e
  provocar você mesmo o erro que aconteceria se não fizesse isso.
- Ler, literalmente, a saída de `go test -race` quando ele encontra uma data race real.
- Usar `context.WithTimeout` para pôr um limite de tempo em uma operação, e explicar por que
  `defer cancelar()` não é opcional.
- Limitar quantas goroutines rodam ao mesmo tempo com um semáforo de canal, e medir que o limite é de fato
  respeitado.

---

## Por que isso importa

Até a lição 5 o `revisor` consulta os serviços **um por um**: se você tem dez serviços e cada
um leva um segundo para responder, o relatório completo leva dez segundos, seja qual for a ordem. Com
cem serviços, quase dois minutos. E o tempo de espera não é você quem decide: quem decide é o serviço mais lento
da lista, multiplicado por quantos houver.

Isso não precisa ser assim, porque **consultar um serviço é esperar, não trabalhar**: quase todo o
tempo dessa espera a CPU não está fazendo nada, só aguardando uma resposta de rede. Go foi projetado
por pessoas que no Google passavam os dias esperando exatamente isso —respostas de rede entre milhares de
serviços— e por isso a concorrência não é uma biblioteca que você adicionou depois: faz parte da linguagem
desde a primeira linha (`go`, uma palavra reservada, não uma função de uma biblioteca).

**O resultado que você vai construir nesta lição:** os mesmos dez serviços de um segundo cada um,
consultados **todos ao mesmo tempo**, em pouco mais de um segundo no total em vez de dez. Esse número —a diferença
entre "em série" e "em paralelo"— é o prêmio da lição, e você vai medi-lo por conta própria, não acreditar nele de
ouvido.

🔑 **E o aviso que torna esta lição diferente das anteriores:** nas lições 2 a 5, um
programa que compila e cujos testes passam quase sempre está certo. **Em concorrência, não.** Um programa
com uma data race pode compilar, rodar, passar nos seus testes noventa e nove vezes, e falhar na
centésima — ou nunca falhar na sua máquina e falhar todos os dias em produção, com mais núcleos e mais
carga. Os erros desta lição são os que **não se veem a olho nu**, e por isso o ponto 4 desta
lição —os erros reais, provocados e capturados— é o mais importante do curso inteiro.

---

## Os conceitos

### 6.1 Goroutines: iniciar é trivial, esperar não

Lançar uma goroutine é uma palavra:

<!-- verificar:fragmento -->
```go
go revisar(s)
```

E pronto. Isso basta para que `revisar(s)` rode **concorrentemente** com o resto do programa, sem esperar
que termine para seguir com a próxima linha. Elas custam aproximadamente 2 KB de memória cada uma para
começar (crescem se for preciso), não os megabytes de uma thread do sistema operacional — por isso um programa
em Go pode ter centenas de milhares de goroutines vivas sem planejar nada especial, algo que seria
impensável com threads do sistema.

Mas olhe este programa, com o bug mais comum do primeiro dia de concorrência em Go (completo em
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

Eu o rodei **260 vezes** nesta máquina para medir com que frequência o problema SE VÊ, não para adivinhar:
40 com `go run`, mais 20 forçando um único processador lógico (`GOMAXPROCS=1`, para tirar do programa
qualquer chance de que uma goroutine consiga rodar em outro núcleo enquanto `main` termina), e 200 com
o binário já compilado (`go build` + `./gor61`, para tirar do caminho o tempo que a compilação leva). O
resultado:

```
$ go run main.go
$ go run main.go
$ go run main.go
```

**Vazio. As três, e praticamente todas as 260** (uma única das 200 execuções do binário compilado
imprimiu algo — o resto, nada).

🔴 **E é preciso ler esse número com muito cuidado, porque é fácil tirar a conclusão errada. O
defeito não está presente «1 a cada 260 vezes»: está presente nas 260 de 260, sem exceção.** Em
**nenhuma** das 260 execuções o programa esperou as suas goroutines — nem mesmo na única que
imprimiu algo. O que muda entre execuções não é se o programa espera (ele nunca espera): é se, por pura
sorte da forma como o sistema operacional reparte o tempo entre processos, alguma goroutine consegue executar
`fmt.Println` na brecha de tempo que existe entre ser lançada e `main` matar o programa inteiro. Essa
brecha quase nunca basta — por isso quase nunca se vê nada — mas o programa está **igualmente quebrado** nas
259 execuções silenciosas e na que imprimiu.

**Então é preciso guardar isto: sempre quebrado, quase nunca visível.** Não é "há uma probabilidade pequena de
que falhe" —ler assim é exatamente o erro que é preciso evitar—: é que o programa **nunca** faz o
correto, e na maioria das vezes o sintoma não chega a aparecer a tempo de você notar. É assim que
esse tipo de bug passa despercebido: passa na revisão de código (compila, roda, não quebra), passa nos testes
manuais (quem roda um programa 260 vezes para confiar nele?), passa no CI (que roda uma vez, talvez
duas) — e só avisa quando já está em produção, com mais carga, mais núcleos, e uma brecha de tempo que um
dia consegue se abrir, no pior momento possível para descobri-lo.

⚠️ **Também não entenda isso como "isto nunca imprime nada, em nenhuma máquina".** Com mais carga, outro sistema
operacional, ou simples azar, a brecha pode se abrir com mais frequência — de fato ela se abriu uma vez nestas
mesmas 260 execuções. O que não muda entre máquinas é que o programa **nunca** espera: isso é
estrutural, não depende da sorte. Rode você mesmo o experimento no seu computador, com
repetições suficientes para que o número signifique algo, e compare.

**Não é um bug intermitente do programa: é que `main` termina assim que chega ao final do seu corpo,
sem se importar se ainda há goroutines rodando** — e quando `main` termina, o programa inteiro
termina com ele, goroutines vivas incluídas, sem aviso e sem erro. Lançar uma goroutine é fácil; o
trabalho de verdade é garantir que o programa espere que elas terminem antes de seguir.

### 6.2 `sync.WaitGroup`: a forma correta de esperar

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

O padrão é sempre o mesmo, e convém memorizá-lo nesta ordem:

1. **`wg.Add(1)` antes de lançar a goroutine**, não dentro dela. Se você o puser dentro, `Wait()` pode
   ser executado antes que a goroutine consiga fazer o seu próprio `Add`, e então não conta com essa
   espera.
2. **`defer wg.Done()` como primeira linha da goroutine.** O `defer` garante que ele seja executado aconteça o
   que acontecer lá dentro —mesmo se a goroutine der panic—, e colocá-lo primeiro evita que, com um `return`
   antecipado no meio do código, você se esqueça de contá-lo.
3. **`wg.Wait()` onde você de fato precisa do resultado**, normalmente logo antes de usar o que as
   goroutines produziram.

**O que acontece se você pular o passo 1 — `Add` faltando —, provocado de propósito** (completo em
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

Eu o rodei **seis vezes seguidas, com `-race`**, para que você veja o que de fato acontece, não só o erro
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

🔑 **Leia isso com cuidado, porque é o erro mais enganoso da lição: o programa imprime `listo`
ANTES de quebrar, nas seis vezes.** Não é uma coincidência desta execução: `wg.Wait()` não espera nada
porque o contador nunca foi incrementado (nunca houve um `Add(1)`), então `Wait()` retorna **de
imediato**, `main` imprime `listo` como se tudo tivesse dado certo, e **só então** uma das
goroutines atrasadas —que estava sim rodando, só que Go nunca esperou por ela— termina o seu `Sleep` e
chama `Done()`, subtraindo um de um contador que já estava em zero. O panic não é a primeira coisa que você vê:
é a última, depois que o programa já te disse que estava tudo bem.

⚠️ **E aqui vem a parte que eu medi errado da primeira vez, e vale a pena deixar a correção escrita: sem
`-race`, este mesmo programa nunca quebra.** As mesmas seis execuções, sem `-race`:

```
$ go run sinesperar.go
listo

$ go run sinesperar.go
listo

(cuatro corridas más, idénticas: "listo" y nada más — 6 de 6, SIN -race, saliendo con código 0)
```

**O `time.Sleep(10 * time.Millisecond)` não é o que FAZ o panic aparecer — é o que o IMPEDE**, sem
`-race`. `main` chega a `wg.Wait()`, que retorna de imediato, imprime `listo`, e o programa inteiro
termina **muito antes** de se cumprirem os 10 milissegundos — então nenhuma das cinco goroutines
atrasadas chega sequer a acordar do seu `Sleep` e chamar `Done()`. O programa parece, literalmente,
"terminado com sucesso": código de saída 0, sem nenhum rastro de que algo estava mal montado. O que de fato faz
o panic aparecer, de forma consistente, é a instrumentação do `-race`: vigiar cada acesso à memória
para detectar corridas faz o programa rodar notavelmente mais devagar, e essa lentidão extra é
exatamente o que dá tempo a alguma goroutine atrasada de completar o seu `Sleep` e chamar `Done()`
antes que o processo termine.

**Este é exatamente o tipo de erro que não falha sempre, e por isso é o mais valioso desta
lição — e agora com o dado correto: você acabou de medir que nem sequer é preciso "um trabalho mais curto"
para que ele passe despercebido. Rodando-o tal como está, SEM `-race`, ele já passa despercebido nas 6 de 6 vezes.** Um
aluno que rode este exemplo sem `-race` —o óbvio, se ninguém o avisar— vai ver `listo` e nada
mais, e vai concluir, com toda razão, que o programa funciona. **Não funciona: o `Add(1)` continua faltando
exatamente igual.** A única coisa que mudou é se algo conseguiu aparecer a tempo de denunciá-lo — o mesmo
padrão "sempre quebrado, quase nunca visível" da seção 6.1, aqui com um ator diferente: não é a carga da
máquina, é se você rodou com o detector de corridas ligado ou não.

**`Done()` subtrai um do contador interno do `WaitGroup`.** Se você nunca fez `Add(1)`, o contador
começa em zero, e subtrair um o torna negativo — e Go, em vez de deixar passar em silêncio, dá
panic com uma mensagem que diz exatamente o que está errado: `negative WaitGroup counter`. Essa mensagem literal
é a sua pista: se você a vir, quase sempre falta um `Add(1)` em algum lugar, ou há mais `Done()` do que
deveria haver. Mas a lição real não é a mensagem do panic — é que **o programa já tinha dito
`listo` antes que ela aparecesse.**

### 6.3 Por que o `revisor` não compartilha um map: canais

O lema da linguagem, e vale a pena memorizá-lo tal como está: **«não comunique compartilhando memória; compartilhe
memória comunicando».** Em vez de cada goroutine escrever o seu resultado em um map ou slice compartilhado
—o que exigiria proteger cada acesso com uma trava—, cada uma manda o seu resultado por um **canal**, e
uma única goroutine (ou o código principal) os recolhe do outro lado.

<!-- verificar:fragmento -->
```go
ch := make(chan servicio.Estado)         // sin buffer: cada envío espera a que alguien reciba
ch := make(chan servicio.Estado, 10)     // con buffer: hasta 10 caben sin esperar a que nadie reciba

ch <- estado        // enviar
e := <-ch           // recibir
close(ch)            // cerrar: después de cerrado, recibir sigue funcionando hasta vaciarlo
for e := range ch { ... }   // recibe hasta que el canal se cierre Y se vacíe
```

**Verifiquei o que acontece se, em vez de um canal, eu uso um slice compartilhado sem proteção**, exatamente o erro
que os canais evitam. Este programa lança mil goroutines que incrementam a mesma variável (completo em
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

Sem o detector de corridas, **o programa não quebra — só dá um resultado incorreto**, e diferente a cada
vez:

```
$ go run carrera.go
contador: 956
```

**956, não 1000.** Nenhum erro, nenhum panic, nenhum sinal de que algo deu errado — só um número menor
do que deveria ser, porque algumas das mil somas se perderam quando duas goroutines leram
`contador` ao mesmo tempo, ambas somaram um ao mesmo valor velho, e uma das duas escritas sobrescreveu
a outra sem que ninguém percebesse. **Este é o bug mais perigoso da concorrência: não falha, dá um
resultado plausível e errado.** Um programa assim pode passar na revisão, passar em testes superficiais, e
falhar em produção de forma esporádica durante meses antes que alguém perceba.

### 6.4 `go test -race`: ver o detector encontrar o bug

O mesmo programa, rodado com `-race`, o denuncia — com evidência exata de onde:

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

Leia de cima para baixo, porque cada bloco responde a uma pergunta diferente:

- **`Read at ... by goroutine 11`** e **`Previous write at ... by goroutine 8`**: duas goroutines
  diferentes tocaram o **mesmo endereço de memória** (`0x00c00013c018`, a variável `contador`), uma
  lendo e outra escrevendo, sem nenhum mecanismo que garanta que uma espere a outra.
- **`carrera.go:15`** em ambas: a linha exata do código onde aconteceu (`contador++`), a mesma linha para
  as duas, porque é a única linha que toca essa variável.
- **`Goroutine 11 (running) created at ... carrera.go:13`**: de onde saiu essa goroutine —a linha do
  `go func() {...}()` dentro do `for`—, para que você possa rastrear qual lançamento foi.
- **`Found 2 data race(s)`** e **`exit status 66`**: o detector não para na primeira corrida que
  encontra; continua rodando e reporta todas, e o programa termina com um código de saída diferente
  de 0 e de 1 (66 é o código que Go reserva para isso), para que um pipeline de integração contínua possa
  distingui-lo de uma falha normal.

**E o número final, `contador: 847`, mudou em relação à execução sem `-race` (956).** Não por
acaso: `-race` faz o programa rodar mais devagar e com mais instrumentação, o que muda a
ordem exata em que as goroutines se entrelaçam — outra razão pela qual este bug é tão traiçoeiro: o
número errado nem sequer é o mesmo errado a cada vez.

🔴 **A regra que decorre disso, sem exceção:** um programa concorrente que passa nos seus testes **sem**
`-race` não está testado. Rode `-race` desde o primeiro teste de concorrência que você escrever, não quando
"algo parecer estranho" — porque, como você acabou de ver, nada parece estranho até que já seja tarde.

### 6.5 Os outros dois modos de falhar: panic e deadlock

**Panic, se você fechar mal um canal:**

- Enviar para um canal já fechado dá panic: `panic: send on closed channel`.
- Fechar um canal duas vezes dá panic: `panic: close of closed channel`.

A regra que evita os dois: **quem fecha o canal é quem envia, nunca quem recebe**, e feche-o uma única vez,
normalmente a partir de uma goroutine dedicada que sabe quando não vai haver mais envios (você vai ver isso na
seção 6.7, com `wg.Wait()` seguido de `close`).

**Deadlock, se ninguém do outro lado estiver escutando.** Eu o provoquei com o programa mais curto possível
(completo em `programas/revisor/ejemplos/06-deadlock/main.go`):

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

**Isto corrige algo que muitas explicações dão como certo: um deadlock em Go nem sempre "trava para
sempre".** O próprio runtime de Go tem um detector de deadlocks: se em algum momento **todas** as
goroutines do programa estão dormindo esperando algo que nunca vai acontecer, Go percebe e mata o programa
com `fatal error: all goroutines are asleep - deadlock!`, com a linha exata onde ele empacou
(`[chan send]`, neste caso, porque estava enviando). Se o seu programa parece "travado para sempre" em
vez de terminar com esse erro, quase certamente **não** é um deadlock verdadeiro: é mais provável que você tenha
pelo menos uma goroutine viva fazendo outra coisa (por exemplo, um temporizador ou um servidor HTTP
escutando), e é ela que impede o runtime de declarar que "todas" estão dormindo.

### 6.6 `context`: como se cancela e se põe um limite de tempo

<!-- verificar:fragmento -->
```go
ctx, cancelar := context.WithTimeout(context.Background(), 5*time.Second)
defer cancelar() // SIEMPRE, incluso si terminas antes de que se cumplan los 5 segundos

req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
resp, err := http.DefaultClient.Do(req) // se aborta solo si pasan los 5 s
```

`context` é o padrão da casa em Go para duas coisas: **pôr um limite de tempo** em uma operação
que poderia demorar demais, e **propagar um cancelamento** para baixo (se o de cima é cancelado,
tudo o que depender dele é cancelado também, sem que cada função tenha de reinventar o seu próprio
mecanismo).

⚠️ **`defer cancelar()` não é opcional, mesmo quando a operação já terminou por conta própria.**
`context.WithTimeout` inicia um temporizador interno; se você nunca chamar a função de cancelar, esse
temporizador continua vivo até que se cumpra o prazo original, retendo memória durante todo esse tempo
— e se o seu programa cria contextos assim o tempo todo (um para cada serviço consultado, como o
`revisor`), sem `defer cancelar()` você acumula vazamentos de memória proporcionais a quantas consultas fez.
É o vazamento mais comum de programas em Go que usam `context`, e por isso convém escrever o `defer` na
mesma linha em que você cria o contexto, antes de escrever qualquer outra coisa.

### 6.7 `Todos`: a função que junta as três peças

Foi assim que ficou, de verdade, a função central do `revisor` (`programas/revisor/internal/revisar/todos.go`), depois de
juntar goroutines, canais, semáforo e `context`:

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

Leia-a com as seções anteriores frescas, porque cada peça responde a um problema que você já viu:

- **Uma goroutine por serviço** (6.1), contada com **`wg.Add(1)`/`defer wg.Done()`** (6.2) para que o
  programa saiba quando todas terminaram.
- **Um canal `resultados` com buffer** (6.3) recolhe os estados: ninguém escreve em um slice ou map
  compartilhado, então não é preciso nenhuma trava para essa parte.
- **`semaforo := make(chan struct{}, paralelo)`** é um canal usado como cota: tem espaço para
  `paralelo` valores, então a `paralelo + 1`-ésima goroutine que tenta escrever nele (`semaforo <-
  struct{}{}`) fica bloqueada até que outra libere o seu lugar (`<-semaforo`, no `defer`). É o mesmo
  canal com buffer da seção 6.3, usado não para levar dados, mas para levar a conta de quantos
  "turnos" restam.
- **Um `context.WithTimeout` próprio por serviço** (6.6), filho do `ctx` geral: se o de cima é
  cancelado (por exemplo, se o tempo total do relatório se esgota), todos os filhos são cancelados com ele.
- **`close(resultados)` a partir de OUTRA goroutine, depois de `wg.Wait()`** — e isto merece explicação,
  porque é a parte que menos se vê a olho nu: se fechássemos o canal na mesma goroutine que
  faz o `for e := range resultados` lá embaixo, travaríamos, porque `Wait()` precisa que todas as
  goroutines leiam do canal para liberar espaço e poder devolver o seu turno, mas o laço de leitura
  nunca começaria porque estaríamos esperando `Wait()` primeiro. Lançando-o à parte, o fechamento e a
  leitura acontecem **ao mesmo tempo**, não um depois do outro.

### 6.8 Testar sem rede: `Falso` e o limite de paralelismo medido

Testar `Todos` contra serviços de verdade seria lento e não determinístico. O `revisor` usa um `Revisor`
falso (`programas/revisor/internal/revisar/falso.go`) que simula respostas, demoras e até serviços que nunca
respondem, tudo em memória:

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

🔒 **Aqui sim é preciso um mutex, e é a exceção à regra de "prefira canais" da seção 6.3.**
`Contador` é um inteiro compartilhado que **muitas goroutines incrementam ao mesmo tempo** durante os
testes (`Falso.Revisar` é chamado concorrentemente, uma vez para cada serviço que `Todos` está
consultando). Um canal serviria para *mandar* resultados, mas para um contador simples compartilhado, um
`sync.Mutex` em volta da única linha que o toca é mais simples e mais claro. A regra não é "nunca use
um mutex": é "antes de usar um, pergunte-se se um canal expressa melhor o que você está fazendo" — e para
mandar resultados, quase sempre sim; para um contador compartilhado, quase sempre o mutex é a ferramenta
correta.

Com `Falso`, este teste mede algo que de outra forma seria quase impossível de verificar com confiança: que
o semáforo da seção 6.7 de fato limita quantas consultas rodam ao mesmo tempo, não só que "funciona
em geral" (versão completa, sem abreviar, em `programas/revisor/internal/revisar/todos_test.go`):

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

Execução real, com o detector de corridas ativo (para confirmar que nem mesmo a própria medição
introduz uma corrida):

```
$ go test ./internal/revisar/... -race -v -run TestTodos_elParaleloLimitaCuantasCorrenALaVez
=== RUN   TestTodos_elParaleloLimitaCuantasCorrenALaVez
--- PASS: TestTodos_elParaleloLimitaCuantasCorrenALaVez (0.11s)
PASS
ok  	github.com/habil/revisor/internal/revisar	1.435s
```

**20 serviços, limite de 3, e o teste confirma que nunca houve mais de 3 chamadas a `Revisar` rodando
ao mesmo tempo** — não porque o supomos a partir do código, mas porque um contador atômico o mediu enquanto rodava.

### 6.9 O tempo, medido: em série contra em paralelo

Com `Falso` configurado para simular demoras reais, o teste `TestTodos_respetaElTimeoutPorServicio`
confirma o outro lado da moeda: um serviço que nunca responde (`Colgado: true`) com um
`Timeout: 30 * time.Millisecond` não pode fazer `Todos` demorar mais do que isso:

```
$ go test ./internal/revisar/... -run TestTodos_respetaElTimeoutPorServicio -v
=== RUN   TestTodos_respetaElTimeoutPorServicio
--- PASS: TestTodos_respetaElTimeoutPorServicio (0.03s)
```

**0.03 segundos, não mais.** O `context.WithTimeout` da seção 6.7, filho do geral, cortou essa
consulta exatamente quando devia, sem que o resto do programa tivesse de ficar sabendo do que aconteceu.

---

## O erro que você vai ver

Todos estes são literais, provocados de propósito para esta lição:

| O sintoma | Mensagem / evidência literal | O que acontece e o que fazer |
|---|---|---|
| `main` termina antes que as goroutines rodem | (nenhuma saída, ou uma saída parcial e inconsistente entre execuções) | Falta um `sync.WaitGroup` (ou um canal) que faça `main` esperar. Seção 6.1 |
| `Add`/`Done` não batem | O programa imprime que **tudo deu certo primeiro**, e `panic: sync: negative WaitGroup counter` chega DEPOIS | Falta um `wg.Add(1)` antes de lançar alguma goroutine, ou sobra um `Done()`. `Wait()` retornou sem esperar nada. Seção 6.2 |
| Dados compartilhados sem proteção | `WARNING: DATA RACE` + `Found 2 data race(s)` + `exit status 66` (com `-race`); um número incorreto sem explicação, sem `-race` | Duas goroutines tocam a mesma memória sem sincronização. Use um canal para mandar o resultado, ou um mutex se de verdade você precisa de um contador compartilhado. Seções 6.3 e 6.4 |
| Enviar para um canal fechado | `panic: send on closed channel` | Outra pessoa já fechou o canal, ou quem o fechou não devia. Feche sempre a partir de quem envia, uma única vez |
| Fechar duas vezes o mesmo canal | `panic: close of closed channel` | Duas goroutines (ou dois caminhos de código) tentam fechar o mesmo canal. Centralize o fechamento em um único lugar |
| Ninguém do outro lado de um canal sem buffer | `fatal error: all goroutines are asleep - deadlock!` | O runtime de Go detecta que o programa inteiro ficou dormindo esperando algo que nunca vai acontecer, e o mata com esta linha exata. Seção 6.5 |
| Você esqueceu `defer cancelar()` | (sem erro imediato; vazamento de memória acumulado com o tempo) | Cada `context.WithTimeout` sem cancelar mantém vivo o seu temporizador até que expire por si só. Seção 6.6 |

**E a instrução que resume esta lição inteira:** se você vai escrever concorrência, rode `-race` desde
o seu primeiro teste, não quando algo cheirar mal — porque, como você viu na seção 6.3, nada cheira mal até
que já seja tarde.

---

## O que se faz errado

- **Lançar uma goroutine por elemento, sem nenhum limite, contra um recurso externo.** Com 5 serviços não
  se nota. Com 5,000, o programa abre 5,000 conexões de uma vez, e o gargalo deixa de ser o
  serviço remoto para ser a sua própria máquina (ou a rede, ou o próprio servidor sobre o qual chovem 5,000
  requisições simultâneas). O semáforo da seção 6.7 existe exatamente para isso — nunca deixe que
  "quantas goroutines eu lanço" dependa só de "quantos elementos eu tenho".
- **Ignorar o `context` que te dão, ou não propagá-lo.** Se uma função recebe um `ctx` e chama outra
  operação que pode demorar sem repassá-lo, essa operação não vai ser cancelada quando o `ctx` original
  for — você tem dois relógios que não conversam. Tudo o que puder demorar recebe o `ctx` de quem o chamou.
- **Não fechar o que se abre.** Um `context.WithTimeout` sem o seu `cancelar()`, uma conexão HTTP sem
  `resp.Body.Close()` (lição 7), um arquivo sem fechar: cada um é um vazamento diferente, mas a forma de
  evitá-los é a mesma — um `defer` logo depois de abrir, antes de escrever qualquer outra linha.
- **Compartilhar memória em vez de comunicá-la "porque é mais rápido de escrever".** Um map compartilhado com
  um mutex em volta do bloco inteiro pode parecer mais curto que montar um canal, mas é mais fácil
  esquecer um `Lock()` em um único ponto de acesso do que esquecer de mandar por um canal — e o primeiro esquecimento não
  se nota até que o detector de corridas (ou, pior, a produção) o encontre.
- **Confiar que "nunca falhou" significa "está certo".** Uma data race pode passar em noventa e
  nove execuções e falhar na centésima, ou só falhar com mais núcleos que os do seu laptop. A única
  prova de que um programa concorrente não tem corridas é rodá-lo com `-race`, não vê-lo passar várias
  vezes sem ele.

---

## Exercícios

1. Escreva o programa da seção 6.1 (goroutines sem esperar) e rode-o cinco vezes seguidas. Anote
   quantas vezes ele imprimiu algo e quantas não.
2. Adicione a ele um `sync.WaitGroup` correto (seção 6.2) e confirme que agora as três linhas são impressas
   sempre, nas cinco execuções.
3. Reproduza o bug de `Add` faltando da seção 6.2 e cole o `panic` completo no seu diário de bordo.
4. Reproduza a data race da seção 6.3 (um contador compartilhado sem proteção) e rode o
   mesmo programa com e sem `-race`. Compare os dois números finais e cole a saída completa do
   `-race` no seu diário de bordo.
5. Reproduza o deadlock da seção 6.5 com um canal sem buffer e sem receptor. Confirme que você vê o
   `fatal error: all goroutines are asleep - deadlock!`, não um travamento silencioso.
6. Pegue o seu próprio `revisor` (ou o deste curso) e rode `TestTodos_elParaleloLimitaCuantasCorrenALaVez`
   mudando o `limite` para 1 e depois para 10. Explique, com as suas palavras, por que o resultado do teste
   não muda (sempre passa) mas o **tempo** que ele leva muda.
7. (Um pouco mais difícil) Tire o semáforo de `Todos` (deixe que todas as goroutines sejam lançadas sem
   limite) e rode de novo `TestTodos_elParaleloLimitaCuantasCorrenALaVez`. Confirme que agora ele falha, e
   cole a mensagem de erro exata que o `t.Errorf` dá com o máximo que de fato foi observado.

### Soluções

1-2. Não há número único: depende da sua máquina. O que importa é a comparação — sem `WaitGroup`,
inconsistente; com ele, sempre as três linhas.

3. O programa imprime `listo` primeiro — `Wait()` retornou de imediato porque o contador nunca foi
   incrementado — e o `panic: sync: negative WaitGroup counter` chega depois, com um stack trace que inclui
   `sync.(*WaitGroup).Add(...)`. Rodado várias vezes seguidas, a ordem é sempre a mesma: `listo`,
   depois o panic.

4. Sem `-race`, um número menor que o esperado (por exemplo, 956 de 1000), sem nenhum erro. Com `-race`,
   o bloco `WARNING: DATA RACE` com a linha exata do código, mais `Found N data race(s)` e
   `exit status 66`.

5. `fatal error: all goroutines are asleep - deadlock!`, com `goroutine 1 [chan send]:` (ou `[chan
   receive]`, conforme o lado em que empacou) e a linha exata do canal.

6. O teste passa nos dois casos porque **mede** o máximo real e o compara com o limite que ele
   mesmo configurou (1 ou 10) — nunca com um número fixo esperado. O tempo total muda: com
   `limite=1` as 20 consultas são estritamente em série (uma espera a outra), com `limite=10` rodam
   em duas levas de 10 em vez de vinte de uma.

7. Sem semáforo, todas as 20 goroutines rodam ao mesmo tempo, então `maximoObservado` vai ficar perto de 20
   (não exatamente o limite do teste, que continua pedindo 3). O `t.Errorf` diz algo como
   `se observaron 20 consultas simultáneas; el límite era 3` — o teste SIM detecta a regressão, que é
   justamente para o que ele existe.

---

## Como sei que consegui

- [ ] Vi, com os meus próprios olhos, um programa terminar sem ter esperado as suas goroutines (seção 6.1).
- [ ] Provoquei o `panic: sync: negative WaitGroup counter` e vi que o programa imprimiu `listo` ANTES
      do panic — e consigo explicar por quê.
- [ ] Provoquei uma data race real e vi a diferença entre rodá-la com e sem `-race`.
- [ ] Consigo ler um bloco de `WARNING: DATA RACE` e dizer que linha, que goroutines e que variável estão
      em conflito.
- [ ] Provoquei o `fatal error: all goroutines are asleep - deadlock!` e sei por que Go o detecta em vez
      de ficar travado para sempre.
- [ ] Consigo explicar, com o código real de `Todos`, para que serve cada uma das suas quatro peças:
      goroutines, canal de resultados, semáforo e `context` por serviço.
- [ ] Medi, com um teste real (não de memória), que o limite de paralelismo do `revisor` é de fato
      respeitado.
- [ ] Sei explicar por que `Falso.Contador` usa um mutex em vez de um canal, e por que isso não contradiz
      a regra geral da seção 6.3.

---

## Para ler mais

1. [A Tour of Go: Concurrency](https://go.dev/tour/concurrency/1) — o tour oficial de goroutines,
   canais e `sync.WaitGroup`, interativo.
2. [Go Blog: Share Memory By Communicating](https://go.dev/blog/codelab-share) — o artigo original
   onde se explica o lema citado na seção 6.3.
3. [Package context](https://pkg.go.dev/context) — a documentação oficial, com os quatro casos de uso
   canônicos (`WithCancel`, `WithTimeout`, `WithDeadline`, `WithValue`).
4. [Data Race Detector](https://go.dev/doc/articles/race_detector) — a documentação oficial de `-race`:
   o que detecta, o que NÃO detecta (por exemplo, não encontra deadlocks, só data races), e o seu custo
   em tempo de execução.

### Termos desta lição

| Termo | O que significa |
|---|---|
| goroutine | uma função que roda concorrentemente com o resto do programa, muito mais barata que uma thread do sistema operacional |
| `sync.WaitGroup` | mecanismo para esperar que um grupo de goroutines termine, contando `Add`/`Done` |
| data race (*condição de corrida*) | duas goroutines acessam a mesma memória ao mesmo tempo, pelo menos uma escrevendo, sem sincronização |
| canal (`chan`) | o mecanismo de Go para que uma goroutine mande dados a outra sem memória compartilhada |
| semáforo de canal | um canal com buffer usado como cota de turnos disponíveis, não para levar dados |
| `context` | mecanismo padrão para propagar limites de tempo e cancelamento entre funções |
| deadlock | estado em que todas as goroutines de um programa estão dormindo esperando algo que nunca vai acontecer; o runtime de Go o detecta e termina o programa |

---

**Anterior:** [Lição 5 — Módulos e testes](05-modulos-y-pruebas.md) ·
**Próxima:** [Lição 7 — O programa terminado](07-el-programa.md)
