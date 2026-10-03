# Lição 7 — O programa concluído

**Duração:** 90-120 minutos, ou 2 sessões de 60. É a última lição do `revisor`: ao terminar você tem
um binário que roda de verdade, é configurado pela linha de comando, fala HTTP com timeouts reais,
produz dois formatos de saída, e é compilado para outra plataforma sem sair da sua máquina.

**Ao terminar, você vai conseguir:**

- Escrever um cliente HTTP com timeout explícito, e explicar por que `http.DefaultClient` é perigoso em
  produção.
- Separar a lógica de um programa (`ejecutar`) do seu ponto de entrada (`main`), para poder testá-la sem
  tocar em `os.Exit`.
- Definir flags de linha de comando com o pacote padrão `flag`, sem nenhuma biblioteca externa.
- Serializar uma estrutura para JSON com `encoding/json`, e explicar que campos se perdem e por quê.
- Rodar o `revisor` de ponta a ponta contra um servidor real (um de teste, feito por você) e ler a sua
  saída.
- Compilar o mesmo código para outra arquitetura e outro sistema operacional sem sair da sua máquina, e
  explicar por que isso é possível.

---

## Por que isso importa

Você já tem as quatro peças do `revisor` construídas separadamente: `servicio` define o vocabulário
(lição 2-3), `config` lê a configuração (lição 5), `revisar.Todos` consulta tudo ao mesmo tempo (lição
6). A única coisa que falta é o que converte essas peças em **um programa que outra pessoa possa usar sem
ler o código-fonte**: que fale HTTP de verdade (até agora só o testamos com `Falso`), que receba
flags do terminal em vez de valores fixos no código, que produza um formato que outro
programa possa consumir, e que possa ser compilado e distribuído como um único arquivo.

Esta é a lição em que o `revisor` deixa de ser "código que funciona se eu o rodo, na minha máquina, com
os meus dados de teste" e se torna um binário que você pode copiar para outra pessoa, com a confiança de que ele vai
fazer exatamente o que a linha de comando pedir.

---

## Os conceitos

### 7.1 O cliente HTTP: por que nunca `http.DefaultClient`

Foi assim que ficou o `Revisor` de verdade, o que de fato toca a rede (`programas/revisor/internal/revisar/revisar.go`):

<!-- verificar:extracto:internal/revisar/revisar.go -->
```go
// HTTP es el Revisor de verdad: hace una petición GET y mira qué contesta.
type HTTP struct {
	Cliente *http.Client
}

// NuevoHTTP construye el Revisor con un cliente propio.
//
// 🔴 El cliente es propio y no http.DefaultClient a propósito: el cliente por
// omisión de Go NO tiene tiempo límite. Una petición contra un servicio que
// acepta la conexión y luego no contesta nada se queda colgada para siempre, sin
// error y sin síntoma.
//
// El Timeout del cliente es la red de seguridad de último recurso. El límite que
// de verdad manda es el del context, uno por servicio, que pone Todos.
func NuevoHTTP(limiteGeneral time.Duration) HTTP {
	return HTTP{Cliente: &http.Client{Timeout: limiteGeneral}}
}

// Revisar consulta un servicio. Nunca devuelve error: el fracaso es parte del
// resultado, porque «este servicio no responde» es justo lo que el programa
// quiere reportar, no una excepción al trabajo.
func (r HTTP) Revisar(ctx context.Context, s servicio.Servicio) servicio.Estado {
	cliente := r.Cliente
	if cliente == nil {
		cliente = &http.Client{Timeout: s.TimeoutEfectivo()}
	}

	inicio := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.URL, nil)
	if err != nil {
		return servicio.Estado{
			Servicio: s,
			Duracion: time.Since(inicio),
			Err:      fmt.Errorf("armando la peticion: %w", err),
		}
	}

	resp, err := cliente.Do(req)
	if err != nil {
		return servicio.Estado{
			Servicio: s,
			Duracion: time.Since(inicio),
			Err:      traducir(ctx, err),
		}
	}
	// Después de comprobar el error, nunca antes: si err != nil, resp es nil y
	// cerrarlo es un pánico.
	defer resp.Body.Close()

	// Se drena el cuerpo aunque no lo queramos leer. Sin esto la conexión no se
	// puede reutilizar y el programa abre una nueva por cada consulta.
	io.Copy(io.Discard, resp.Body)

	return servicio.Estado{
		Servicio: s,
		Codigo:   resp.StatusCode,
		Duracion: time.Since(inicio),
	}
}
```

🔴 **`http.DefaultClient` —o que você usaria se escrevesse `http.Get(url)` diretamente— não tem nenhum
timeout.** Se o servidor do outro lado aceita a conexão e nunca responde nada, essa chamada fica
esperando **para sempre**, sem erro, sem panic, só um programa que um dia deixa de avançar e ninguém
sabe por quê. É um dos erros mais caros de Go em produção justamente porque não dá nenhum
sintoma até que já esteja acontecendo.

O `revisor` se protege duas vezes, não uma: o `http.Client` próprio tem um timeout geral (o limite de
todo o relatório, passado a `NuevoHTTP`), e além disso cada consulta individual roda sob o
`context.WithTimeout` por serviço que `Todos` montou na lição 6 — esse segundo, mais curto, é o que
manda na prática. O timeout do cliente é a rede de segurança de último recurso, não o mecanismo
principal.

⚠️ **`defer resp.Body.Close()` vem depois de verificar o erro, nunca antes.** Se `err != nil`, `resp`
é `nil`, e chamar um método sobre um ponteiro nulo dá panic. E **`io.Copy(io.Discard, resp.Body)`
antes de fechar** não é decorativo: sem drenar o corpo da resposta, a conexão TCP subjacente não
pode ser reutilizada para a próxima requisição ao mesmo servidor, e o programa acaba abrindo uma
conexão nova a cada vez em vez de reutilizar as que já tem.

**`traducir` troca o erro cru de `net/http` por um legível em uma tabela:**

<!-- verificar:extracto:internal/revisar/revisar.go -->
```go
func traducir(ctx context.Context, err error) error {
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("se acabo el tiempo de espera")
	}
	if ctx.Err() == context.Canceled {
		return fmt.Errorf("consulta cancelada")
	}
	return fmt.Errorf("no responde: %w", err)
}
```

O erro que o `net/http` dá de fábrica traz a URL completa e a palavra `Get`, que em uma tabela de
relatório só ocupam espaço e não dizem nada novo ao leitor — por isso ele é traduzido antes de ser mostrado.

### 7.2 `main` não se testa; `ejecutar` sim

O padrão que separa o `revisor` de um programa de um único arquivo:

<!-- verificar:extracto:cmd/revisor/main.go -->
```go
func main() {
	os.Exit(ejecutar(os.Args[1:], os.Stdout, os.Stderr))
}
```

`ejecutar` é onde vive a lógica real —as seções 7.3, 7.4 e 7.6 a vão mostrando por partes—; a sua
assinatura:

<!-- verificar:fragmento -->
```go
func ejecutar(args []string, salida, errores *os.File) int {
	// toda la lógica real vive aquí
}
```

**`go test` não consegue testar uma função que chama `os.Exit`**, porque `os.Exit` termina o processo
inteiro de imediato — incluindo o próprio processo de testes, que jamais chegaria a reportar o resultado.
Separando "decidir que código de saída corresponde" (`ejecutar`, que **retorna** um `int`) de "sair de
verdade com esse código" (`main`, a única linha que chama `os.Exit`), toda a lógica pode ser testada
chamando `ejecutar` diretamente, com argumentos de teste, e verificando o número que ele retorna — exatamente
o que `cmd/revisor/main_test.go` faz:

<!-- verificar:extracto:cmd/revisor/main_test.go -->
```go
func TestEjecutar_faltaConfig(t *testing.T) {
	var salida bytes.Buffer
	codigo := correr(t, []string{}, &salida)
	if codigo != 2 {
		t.Errorf("código de salida = %d, quería 2 (falta --config)", codigo)
	}
}
```

```
$ go test ./cmd/revisor/... -v -run TestEjecutar_faltaConfig
=== RUN   TestEjecutar_faltaConfig
revisor: falta --config
Usage of revisor:
  -config string
    	ruta al archivo de servicios (obligatorio)
  -formato string
    	tabla o json (default "tabla")
  -limite duration
    	tiempo máximo para todo el reporte (default 10s)
  -paralelo int
    	cuántos servicios consultar a la vez (0 = por omisión)
--- PASS: TestEjecutar_faltaConfig (0.00s)
```

Essa é a saída real de `--help` (gerada automaticamente por `flag`, seção 7.3) capturada no meio
de um teste, sem abrir um terminal nem executar o binário compilado.

### 7.3 Flags de linha de comando, sem bibliotecas

<!-- verificar:extracto:cmd/revisor/main.go -->
```go
banderas := flag.NewFlagSet("revisor", flag.ContinueOnError)
rutaConfig := banderas.String("config", "", "ruta al archivo de servicios (obligatorio)")
formato := banderas.String("formato", "tabla", "tabla o json")
paralelo := banderas.Int("paralelo", 0, "cuántos servicios consultar a la vez (0 = por omisión)")
limiteGeneral := banderas.Duration("limite", 10*time.Second, "tiempo máximo para todo el reporte")

if err := banderas.Parse(args); err != nil {
	return 2 // flag ya imprimió el error y el uso
}
```

**`flag.NewFlagSet` em vez do pacote `flag` no nível global** (que você usaria com `flag.String(...)`
diretamente) é o que permite ter uma função `ejecutar(args []string, ...)` que recebe os seus argumentos
como parâmetro, em vez de lê-los sempre de `os.Args` — outro detalhe que existe especificamente para
poder testá-la, como na seção 7.2.

`flag` vem na biblioteca padrão e basta para um programa deste tamanho: quatro flags, tipos
básicos (`string`, `int`, `time.Duration`), ajuda automática. Se algum dia o `revisor` precisasse de
subcomandos (`revisor check`, `revisor list`, cada um com as suas próprias flags), aí sim valeria a pena
olhar uma biblioteca como `cobra` — mas não antes de precisar dela de verdade.

### 7.4 Os dois formatos de saída: tabela e JSON

Foi assim que ficou `Tabla`, a função que produz a saída legível por humanos que você viu em toda a lição
(`programas/revisor/internal/reporte/tabla.go`):

<!-- verificar:extracto:internal/reporte/tabla.go -->
```go
func Tabla(w io.Writer, estados []servicio.Estado) {
	ordenados := ordenarPorNombre(estados)

	anchoNombre := len("SERVICIO")
	for _, e := range ordenados {
		if n := len(e.Servicio.Etiqueta()); n > anchoNombre {
			anchoNombre = n
		}
	}

	fmt.Fprintf(w, "%-*s  %-6s  %8s  %s\n", anchoNombre, "SERVICIO", "ESTADO", "TIEMPO", "DETALLE")
	for _, e := range ordenados {
		fmt.Fprintf(w, "%-*s  %-6s  %8s  %s\n",
			anchoNombre,
			e.Servicio.Etiqueta(),
			etiquetaEstado(e),
			formatoTiempo(e),
			detalle(e),
		)
	}
}
```

**A largura da primeira coluna não é fixa no código — ela é calculada.** O laço acima percorre
todos os estados uma vez antes de imprimir qualquer coisa, e fica com o nome mais longo (começando pela
largura da própria palavra "SERVICIO", caso algum nome de serviço seja mais curto que isso). Esse
número entra como o `*` em `%-*s`: um verbo de formatação de largura **variável**, em que a própria largura é
mais um argumento de `Fprintf`, não um número escrito à mão. É a razão pela qual a tabela real desta
lição tem colunas retas não importa se os nomes dos serviços são curtos (`sano`) ou longos
(`inventario`): com uma largura fixa, um dos dois casos sempre ficaria torto.

`ordenarPorNombre` (as seções 6.3 e 6.7 já explicaram por que a ordem de chegada não é confiável: várias
goroutines entregam resultados por um canal, e ganha quem responder primeiro) faz uma cópia do slice e a
ordena alfabeticamente antes de imprimir — sem isso, a mesma consulta produziria uma tabela em uma ordem
diferente cada vez que você rodasse o programa, ainda que os dados fossem idênticos.

`etiquetaEstado`, `formatoTiempo` e `detalle` são as três funções pequenas que decidem que palavra vai em
cada coluna (`OK`/`LENTO`/`FALLA`, o tempo arredondado ao milissegundo ou um hífen se nunca houve
resposta, e o código HTTP ou o motivo da falha) — são as que você já viu atuando na tabela real da
seção 7.5, agora com o código que as produz.

🔑 **Por que não se usou `text/tabwriter`, que é a ferramenta que a biblioteca padrão oferece
justamente para alinhar colunas:** para uma tabela de quatro colunas em que só uma tem largura variável
(o nome do serviço; as outras três são curtas e previsíveis), calcular a largura à mão é mais simples
de ler do que introduzir um `tabwriter.Writer` com os seus próprios `Flush()` e separadores por tabulação. Com
uma tabela de mais colunas variáveis, `tabwriter` seria sim a ferramenta correta — vale a pena conhecê-lo
(está na seção "Para ler mais" desta lição) ainda que o `revisor` não precise dele.

Com a tabela já entendida, o outro formato de saída — JSON — é a outra metade desta seção:

<!-- verificar:extracto:internal/reporte/json.go -->
```go
type lineaJSON struct {
	Servicio string `json:"servicio"`
	OK       bool   `json:"ok"`
	Codigo   int    `json:"codigo,omitempty"`
	TiempoMs int64  `json:"tiempo_ms"`
	Detalle  string `json:"detalle,omitempty"`
}
```

⚠️ **`servicio.Estado` não é serializado diretamente — de propósito.** `Estado` traz um campo `Err error`, e
`error` é uma interface: `encoding/json` não sabe como convertê-la em texto por si só (tentar isso produz
um objeto vazio `{}`, não um erro de compilação, o que é pior: falha em silêncio). O `revisor` resolve
isso com um tipo intermediário, `lineaJSON`, que só vive dentro do pacote `reporte`: converte o erro
no seu texto (`e.Motivo()`) antes de codificar, e de quebra desacopla o formato público (o que outros
programas vão ler) do modelo interno (`Estado`) — se amanhã `Estado` ganhar um campo novo, o JSON que
já circula não muda só porque algo interno mudou.

`omitempty` em `Codigo` e `Detalle` remove esses campos do JSON quando valem zero ou string vazia — assim um
serviço saudável não arrasta um `"detalle": ""` sem sentido. Codificado de verdade:

<!-- verificar:extracto:internal/reporte/json.go -->
```go
codificador := json.NewEncoder(w)
codificador.SetIndent("", "  ")
return codificador.Encode(lineas)
```

### 7.5 De ponta a ponta: o `revisor` rodando contra um servidor real

Para testar todo o programa junto —não cada peça separadamente— construí `cmd/servidor-demo`: um
servidor HTTP mínimo com quatro rotas que se comportam como os quatro casos que o `revisor` precisa
saber reportar.

<!-- verificar:extracto:cmd/servidor-demo/main.go -->
```go
mux.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
})
mux.HandleFunc("/lento", func(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1500 * time.Millisecond)
	w.WriteHeader(http.StatusOK)
})
mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
})
mux.HandleFunc("/nunca-contesta", func(w http.ResponseWriter, r *http.Request) {
	<-r.Context().Done() // no responde nunca por su cuenta; solo cede cuando el cliente se cansa
})
```

Com esse servidor no ar em `:8091` e um arquivo de configuração apontando para as suas quatro rotas, rodei
o binário real:

```
$ /tmp/servidor-demo -puerto 8091 &
servidor-demo escuchando en :8091 (/ok, /lento, /error, /nunca-contesta)

$ cat /tmp/servicios-demo.txt
sano      http://127.0.0.1:8091/ok
lento     http://127.0.0.1:8091/lento
malo      http://127.0.0.1:8091/error
colgado   http://127.0.0.1:8091/nunca-contesta   200ms

$ /tmp/revisor --config /tmp/servicios-demo.txt --formato tabla
SERVICIO  ESTADO    TIEMPO  DETALLE
colgado   FALLA      202ms  se acabo el tiempo de espera
lento     LENTO     1.503s  200
malo      FALLA        1ms  codigo 500
sano      OK           1ms  200
$ echo $?
1
```

**Cada linha dessa tabela corresponde exatamente ao comportamento que foi programado no servidor de
teste:** `sano` responde rápido e com 200; `lento` leva 1.5 segundos de verdade e por isso sai `LENTO`;
`malo` responde na hora mas com 500; `colgado` nunca responde, e o seu timeout individual de 200ms
—declarado no arquivo de configuração, coluna três— cortou a espera em 202ms, quase exato. **O
código de saída, 1, é real**: pelo menos um serviço não estava `OK`, então `ejecutar` retorna 1 em vez
de 0 (seção 7.6) — o mesmo binário, usado a partir de um script, pode dizer a quem o invocar se houve
problemas sem que ninguém precise ler a tabela.

E em JSON, o mesmo relatório:

```
$ /tmp/revisor --config /tmp/servicios-demo.txt --formato json
[
  {"servicio": "colgado", "ok": false, "tiempo_ms": 205, "detalle": "se acabo el tiempo de espera"},
  {"servicio": "lento", "ok": true, "codigo": 200, "tiempo_ms": 1503},
  {"servicio": "malo", "ok": false, "codigo": 500, "tiempo_ms": 1, "detalle": "codigo 500"},
  {"servicio": "sano", "ok": true, "codigo": 200, "tiempo_ms": 1}
]
```

(Aqui sem a indentação de `SetIndent` para que caiba em uma linha por serviço; o programa real a
produz com quebras de linha, como você viu na seção 7.4.)

### 7.6 O código de saída, com significado

<!-- verificar:extracto:cmd/revisor/main.go -->
```go
for _, e := range estados {
	if !e.OK() {
		return 1 // al menos un servicio falló: el código de salida lo refleja
	}
}
return 0
```

Três códigos possíveis, e cada um responde a uma pergunta diferente para quem invocar o programa a partir de um
script ou de um pipeline: **0** — tudo deu certo; **1** — o programa rodou por completo, mas pelo menos um
serviço não estava saudável; **2** — o programa nem sequer conseguiu iniciar (falta `--config`, o arquivo não
existe, o formato de configuração está mal escrito). É a diferença entre "fiz o trabalho e encontrei
problemas" e "não consegui nem começar a trabalhar" — e um pipeline de integração contínua pode reagir
de forma diferente a cada um.

### 7.7 Compilação cruzada: o mesmo código, outra plataforma

```bash
GOOS=linux GOARCH=amd64 go build -o revisor-linux-amd64 ./cmd/revisor
```

Confirmado de verdade, compilando a partir deste Mac (arm64) para Linux x86-64:

```
$ GOOS=linux GOARCH=amd64 go build -o revisor-linux-amd64 ./cmd/revisor
$ file revisor-linux-amd64
revisor-linux-amd64: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, ...
```

**Sem Docker, sem máquina virtual, sem instalar nada adicional: duas variáveis de ambiente bastam.** Isso é
possível porque o compilador de Go traz, de fábrica, o código necessário para gerar binários de
qualquer combinação de sistema operacional e arquitetura que ele suporte — não depende de estar rodando no
sistema de destino para compilar para ele, ao contrário de como funcionam muitas outras linguagens compiladas.

**O tamanho do binário, medido, com e sem símbolos de depuração:**

```
$ go build -o revisor-normal ./cmd/revisor
$ wc -c revisor-normal
9464130 revisor-normal

$ go build -ldflags="-s -w" -o revisor-chico ./cmd/revisor
$ wc -c revisor-chico
6386194 revisor-chico
```

**De 9.46 MB para 6.39 MB, uma redução real de 32%**, removendo a tabela de símbolos (`-s`) e as
informações de depuração DWARF (`-w`) que o binário normal inclui para que um depurador possa
inspecioná-lo. Para um binário que você vai distribuir para produção e não vai depurar ali mesmo, esses dados
não servem para nada e só ocupam espaço — para um que você está desenvolvendo ativamente, conserve-os.

### 7.8 Testar código que faz HTTP, sem rede de verdade: `httptest`

Até a lição 6, os testes de concorrência usavam `Falso` (um `Revisor` que nunca toca a rede). Para
testar o programa **completo** —flags, leitura de configuração, o cliente HTTP real, o relatório—
sem depender de um serviço externo nem de que algo esteja escutando em uma porta fixa, `net/http/httptest`
sobe um servidor real, em uma porta que o sistema operacional atribui sozinho, dentro do próprio processo de
testes:

<!-- verificar:extracto:cmd/revisor/main_test.go -->
```go
func TestEjecutar_reportaOKyFalla(t *testing.T) {
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/mal") {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer servidor.Close()

	archivoConfig, err := os.CreateTemp(t.TempDir(), "servicios-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	contenido := "sano " + servidor.URL + "/bien\nmalo " + servidor.URL + "/mal\n"
	if _, err := archivoConfig.WriteString(contenido); err != nil {
		t.Fatal(err)
	}
	archivoConfig.Close()

	var salida, errores bytes.Buffer
	codigo := correr(t, []string{"--config", archivoConfig.Name()}, &salida)

	if codigo != 1 {
		t.Errorf("código de salida = %d, quería 1 (hay un servicio con falla)", codigo)
	}
	if !strings.Contains(salida.String(), "sano") || !strings.Contains(salida.String(), "malo") {
		t.Errorf("la tabla no menciona a los dos servicios:\n%s", salida.String())
	}
	_ = errores
}
```

Execução real:

```
$ go test ./cmd/revisor/... -run TestEjecutar_reportaOKyFalla -v
=== RUN   TestEjecutar_reportaOKyFalla
--- PASS: TestEjecutar_reportaOKyFalla (0.00s)
```

**`servidor.URL` é uma URL real** (`http://127.0.0.1:PUERTO`, com uma porta livre que o sistema escolheu),
então o `Revisor` HTTP da seção 7.1 a consulta exatamente como consultaria qualquer
serviço de produção — a única diferença é que o "serviço" é um `http.HandlerFunc` de quinze
linhas que vive no mesmo processo que o teste. É isso que permite testar de ponta a ponta
—flags, configuração, cliente HTTP, relatório, código de saída— em milissegundos, sem abrir nenhuma
porta fixa que pudesse colidir com outra coisa rodando na máquina, e sem deixar nada rodando depois
que o teste termina (`defer servidor.Close()` cuida disso).

**`t.TempDir()`** cria uma pasta temporária que o Go apaga sozinho ao terminar o teste (diferente de
`os.CreateTemp("/tmp", ...)` puro e simples, que deixaria o arquivo ali para sempre se ninguém o apagasse à mão) —
a combinação de `httptest` para a rede e `t.TempDir()` para o disco é o que permite que
`TestEjecutar_reportaOKyFalla` não deixe nenhum rastro no sistema depois de rodar.

### 7.9 Por que dois limites de tempo, não um

O `revisor` aceita `--limite` (o tempo máximo para **todo** o relatório) e cada linha de configuração
pode declarar o seu próprio tempo limite **por serviço**. Não é redundante: eles resolvem duas perguntas
diferentes. `--limite` responde *«quanto tempo, no máximo, estou disposto a esperar pelo relatório
completo?»* — útil se o `revisor` roda dentro de outro processo com o seu próprio prazo (um health check
que um orquestrador espera de tempos em tempos, por exemplo). O tempo por serviço responde *«quanto é
razoável esperar por ESTE serviço em particular?»* — um serviço interno de baixa latência e outro que vai
até um fornecedor externo não deveriam compartilhar o mesmo limite.

A demonstração da seção 7.5 confirma isso com dados reais: `colgado` tinha um limite próprio de
200ms declarado no arquivo de configuração, e foi cortado em 202ms — muito antes que o `--limite`
geral (10 segundos por padrão) tivesse oportunidade de intervir. O mais curto dos dois limites é
sempre o que manda, e isso é exatamente o que você quer: que um único serviço lento não consuma todo o
orçamento de tempo do relatório completo.

### 7.10 O binário "não precisa de nada" — exceto uma coisa: certificados

A lição 1 demonstrou que um binário de Go roda em uma máquina Linux sem Go instalado, graças à
linkagem estática. É tentador estender essa ideia para "não precisa de nada do sistema, ponto" — e eu verifiquei isso
levando-a ao extremo: coloquei o binário do `revisor` em uma imagem Docker construída `FROM scratch`,
a base mais vazia que existe (não tem sequer um shell nem componentes do sistema operacional, só o binário).

```dockerfile
FROM scratch
COPY revisor /revisor
COPY servicios.txt /servicios.txt
ENTRYPOINT ["/revisor", "--config", "/servicios.txt"]
```

Com um serviço real apontando para `https://example.com`:

```
$ docker build -t revisor-demo:scratch .
$ docker run --rm revisor-demo:scratch --formato tabla
SERVICIO  ESTADO    TIEMPO  DETALLE
ejemplo   FALLA      176ms  no responde: Get "https://example.com": tls: failed to verify certificate: x509: certificate signed by unknown authority
```

**Falhou — e não por causa do `revisor`, mas por algo que nenhuma imagem `scratch` traz: a lista de autoridades
certificadoras.** Para validar um certificado TLS (o cadeado do `https://`), o sistema operacional
normalmente fornece uma lista de quem tem permissão para assinar certificados válidos — tipicamente em
`/etc/ssl/certs/`. Uma imagem `scratch` não tem nem essa pasta, então Go não consegue verificar nenhum
certificado e se recusa a continuar, com toda razão: seguir sem verificar seria aceitar qualquer
certificado, válido ou falso.

**A solução, verificada, é copiar só esse arquivo** a partir de uma imagem que o tenha, sem arrastar
o resto do sistema operacional:

```dockerfile
FROM alpine:3.20 AS certs
RUN apk add --no-cache ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY revisor /revisor
COPY servicios.txt /servicios.txt
ENTRYPOINT ["/revisor", "--config", "/servicios.txt"]
```

```
$ docker build -f Dockerfile.certs -t revisor-demo:scratch-certs .
$ docker run --rm revisor-demo:scratch-certs --formato tabla
SERVICIO  ESTADO    TIEMPO  DETALLE
ejemplo   OK         238ms  200

$ docker images revisor-demo:scratch-certs --format "{{.Size}}"
14.7MB
```

**14.7 MB no total, contra 14.4 MB sem os certificados — menos de 300 KB para resolver o problema
corretamente.** O binário é sim autossuficiente para tudo o que é *código Go*; o que ele nunca traz por si
só é a confiança sobre quais autoridades são legítimas, porque isso é informação do mundo (que
certificadoras existem e continuam vigentes), não algo que o compilador possa decidir por você.

---

## O erro que você vai ver

| O sintoma | Mensagem / evidência literal | O que acontece e o que fazer |
|---|---|---|
| Falta a flag obrigatória | `revisor: falta --config` seguido do uso automático de `flag` | `ejecutar` valida explicitamente que `--config` não venha vazio antes de fazer qualquer outra coisa, e sai com código 2 |
| Formato de saída inválido | `revisor: --formato debe ser "tabla" o "json", no "xml"` | Só existem dois formatos; qualquer outro valor é rejeitado antes de tentar qualquer coisa |
| Serviço que nunca responde | Na tabela: `FALLA` com detalhe `se acabo el tiempo de espera`, tempo próximo do timeout configurado | O `context.WithTimeout` por serviço (lição 6) cortou a espera; não é um bug, é o mecanismo funcionando |
| Fechar o corpo de uma resposta nula | (se fosse feito errado) `panic: runtime error: invalid memory address or nil pointer dereference` | Aconteceria se `defer resp.Body.Close()` fosse colocado antes de verificar `err`; o `revisor` evita isso verificando o erro primeiro (seção 7.1) |
| `error` serializado diretamente para JSON sem traduzir | `{}` (objeto vazio, sem nenhum dado útil, sem nenhum erro de compilação que avise) | `encoding/json` não sabe converter uma interface `error`; por isso `reporte` usa `lineaJSON` com o motivo já convertido em texto (seção 7.4) |
| Código de saída não verificado em um script | O script segue como se nada fosse, ainda que um serviço tenha falhado | `ejecutar` distingue sim 0/1/2 (seção 7.6); quem integra o `revisor` em um pipeline deve ler `$?`, não só a saída impressa |
| Binário em uma imagem `FROM scratch`, contra HTTPS | `no responde: Get "https://...": tls: failed to verify certificate: x509: certificate signed by unknown authority` | Falta a lista de autoridades certificadoras (`ca-certificates.crt`); copie-a a partir de uma imagem que a tenha (seção 7.10) |

---

## O que se faz errado

- **Usar `http.DefaultClient` ou `http.Get` diretamente em produção.** Seção 7.1: sem timeout próprio, uma
  conexão que nunca responde trava o programa para sempre, sem nenhum sintoma prévio.
- **Colocar toda a lógica dentro de `main` e chamar `os.Exit` no meio do código.** Torna a função
  impossível de testar com `go test`, porque `os.Exit` termina o processo de testes junto com o
  programa. Separe "decidir o código de saída" de "sair de verdade" (seção 7.2).
- **Serializar `error` diretamente para JSON, esperando que "algo saia".** Sai um objeto vazio, sem nenhum
  aviso de que algo não pôde ser convertido — o tipo de falha silenciosa mais difícil de detectar, porque
  o programa não quebra nem reclama.
- **Fechar o corpo de uma resposta HTTP antes de verificar o erro.** Se a requisição falhou, `resp` é
  `nil`, e chamar um método sobre ele dá panic. Verifique o erro primeiro, sempre.
- **Não drenar o corpo da resposta antes de fechá-lo.** A conexão não pode ser reutilizada, e um
  programa que faz muitas requisições ao mesmo servidor acaba abrindo uma conexão nova a cada vez, mais
  lento do que o necessário sem nenhum erro que o indique.
- **Ignorar o código de saída do binário a partir de um script ou de um pipeline.** O `revisor` distingue
  "rodei bem, tudo saudável" (0), "rodei bem, algo falhou" (1) e "não consegui nem iniciar" (2) — desperdiçar isso
  lendo só a saída impressa é jogar fora informação que o programa já está te dando de graça.

---

## Exercícios

1. Implemente (ou revise, se você já o tem) o `Revisor` HTTP real da seção 7.1, com o seu próprio cliente
   e timeout. Confirme que compila e que `go vet ./...` não reclama.
2. Suba `cmd/servidor-demo` em uma porta livre e rode o seu `revisor` contra as suas quatro rotas, como na
   seção 7.5. Cole a tabela real que você obteve, não uma inventada.
3. Adicione a `Tabla` (seção 7.4) uma quinta coluna, `PROTOCOLO`, que diga `https` ou `http` conforme
   `Servicio.EsSeguro()` (o método que foi definido na lição 3). Ajuste a largura fixa dessa coluna
   à mão, sem precisar do cálculo dinâmico que `anchoNombre` já tem.
4. Rode o mesmo relatório com `--formato json` e valide que o resultado é JSON legítimo (por exemplo,
   passe-o por `jq .` ou cole-o em um validador de JSON).
5. Provoque o erro de `--formato` inválido (seção 7.6, tabela) e confirme o código de saída com
   `echo $?`.
6. Compile o seu `revisor` para `GOOS=linux GOARCH=amd64` a partir da sua máquina atual e confirme com `file` que
   o binário resultante é o da plataforma correta.
7. (Um pouco mais difícil) Meça o tamanho do seu binário com e sem `-ldflags="-s -w"`, como na seção
   7.7, e calcule a porcentagem de redução.
8. (Fechamento do curso) Rode a suíte completa do projeto — `go vet ./...`, `gofmt -l .` e
   `go test ./... -race -cover` — e cole a saída completa no seu diário de bordo. Se algo não estiver verde,
   corrija antes de dar o curso por encerrado.

### Soluções

1-7. Não há uma única saída de referência porque depende da sua própria implementação e da sua própria
máquina; compare a forma do seu resultado com as seções correspondentes (7.1, 7.4, 7.5, 7.6, 7.7).

8. A execução real, sobre o `revisor` deste curso, em 30-set-2026:

   ```
   $ gofmt -l .
   (sin salida: nada por formatear)

   $ go vet ./...
   (sin salida: limpio)

   $ go test ./... -race -cover
   ok  	github.com/habil/revisor/cmd/revisor	1.21s	coverage: 79.2% of statements
   	github.com/habil/revisor/cmd/servidor-demo		coverage: 0.0% of statements
   ok  	github.com/habil/revisor/internal/config	1.17s	coverage: 97.4% of statements
   ok  	github.com/habil/revisor/internal/reporte	1.18s	coverage: 100.0% of statements
   ok  	github.com/habil/revisor/internal/revisar	1.33s	coverage: 61.9% of statements
   ok  	github.com/habil/revisor/internal/servicio	1.17s	coverage: 92.3% of statements
   ```

   `cmd/servidor-demo` sai com 0.0% sem `ok` nem `FAIL`, não porque algo esteja quebrado: com `-cover`, um
   pacote sem nenhum arquivo `_test.go` mesmo assim é instrumentado e reportado, mas não há nenhum teste
   que exercite esse código. É uma ferramenta de apoio para testar o programa completo, não parte do
   `revisor` que é distribuído, e não tem lógica própria que decida nada — só responde o que foi
   programado para responder — por isso não precisa de testes próprios.

   ⚠️ **Os 61.9% de `internal/revisar` nesta tabela são o artefato que a seção 5.7 explica: por
   pacote, não se conta o que `TestEjecutar_reportaOKyFalla` (em `cmd/revisor`, um pouco mais acima nesta
   mesma execução) exercita de `Revisar` ao falar HTTP de verdade contra o servidor `httptest`.** Medido
   com `-coverpkg=./...` sobre todo o projeto: `Revisar` sobe para 86.4% e o total do projeto é 81.5%,
   não 61.9%. Se você vai citar um número de cobertura para decidir algo, rode `-coverpkg=./...`, não confie
   só no número por pacote.

---

## Como sei que consegui

- [ ] O meu `Revisor` HTTP tem o seu próprio cliente, com timeout, e nunca usa `http.DefaultClient`.
- [ ] Separei `main` (que só chama `os.Exit`) de uma função que faz o trabalho e retorna um `int`.
- [ ] O meu binário aceita `--config`, `--formato`, `--paralelo` e `--limite` pela linha de comando.
- [ ] Rodei o `revisor` real contra um servidor real (o meu ou o `servidor-demo` do curso) e a tabela
      que obtive reflete exatamente o que esse servidor faz.
- [ ] `--formato json` produz JSON válido, verificado com uma ferramenta que não sou eu lendo.
- [ ] `echo $?` depois de rodar o `revisor` dá 0, 1 ou 2, e sei explicar cada um.
- [ ] Compilei para outra plataforma (`GOOS`/`GOARCH`) e confirmei com `file` que o binário é o correto.
- [ ] `go vet ./...`, `gofmt -l .` e `go test ./... -race -cover` estão limpos na minha própria cópia do
      projeto — não suponho, eu rodei.

---

## Para ler mais

1. [net/http package](https://pkg.go.dev/net/http) — a documentação oficial completa do cliente e
   servidor HTTP da biblioteca padrão.
2. [Command flag](https://pkg.go.dev/flag) — a documentação oficial do pacote de flags usado
   nesta lição.
3. [encoding/json: JSON and Go](https://go.dev/blog/json) — o artigo oficial sobre como Go traduz
   entre structs e JSON, incluindo as regras das tags (`json:"..."`, `omitempty`, `-`).
4. [Build constraints e compilação cruzada](https://pkg.go.dev/cmd/go#hdr-Environment_variables) — a
   referência de `GOOS`/`GOARCH` e das demais variáveis de ambiente que controlam `go build`.
5. [text/tabwriter](https://pkg.go.dev/text/tabwriter) — a ferramenta padrão para alinhar colunas
   quando o cálculo manual da seção 7.4 não basta (várias colunas de largura variável).

### Termos desta lição

| Termo | O que significa |
|---|---|
| `http.Client` | o tipo que faz requisições HTTP em Go; sem configurar, não tem limite de tempo |
| `flag.FlagSet` | conjunto de flags de linha de comando que pode ser analisado a partir de um slice próprio, não só de `os.Args` |
| `omitempty` | opção de tag JSON que remove um campo da saída quando vale zero ou está vazio |
| código de saída | o número inteiro que um programa retorna ao terminar; 0 significa sucesso por convenção universal |
| `GOOS` / `GOARCH` | variáveis de ambiente que dizem ao `go build` para que sistema operacional e arquitetura compilar |
| `-ldflags="-s -w"` | opções de compilação que removem símbolos e dados de depuração para reduzir o tamanho do binário |

---

## E agora, o que vem a seguir

Você tem um programa completo: consulta serviços reais, todos ao mesmo tempo, com limites de tempo, e produz
um relatório que um humano ou um programa podem ler. Três coisas que o tornam mais sólido, na ordem em
que mais convém atacá-las:

1. **`golangci-lint`** — junta dezenas de analisadores estáticos em uma única execução. Passe-o pelo seu
   próprio `revisor` e leia o que ele diz: quase sempre ensina algo que nem `go vet` nem os testes pegam.
2. **[Effective Go](https://go.dev/doc/effective_go) outra vez.** Agora que você escreveu um programa
   completo, vai lê-lo de forma diferente da lição 0.
3. **Leia bom código escrito por outras pessoas.** O próprio pacote `net/http` da biblioteca padrão é um bom ponto de
   partida: você já tem com o que se orientar para entender por que ele está escrito como está.

**E o prometido desde o README:** você vai escrever este mesmo programa outra vez, em Rust, no curso
irmão. Não para comparar sintaxe, mas para ver o mesmo problema resolvido com outras ferramentas —isso
é o que de verdade ensina em que as duas linguagens se diferenciam.

---

**Anterior:** [Lição 6 — Concorrência](06-concurrencia.md)
