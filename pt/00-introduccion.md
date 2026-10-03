# Lição 0 — O que é Go, de onde veio e por que você deveria aprendê-lo

**Duração:** 30 a 45 minutos de leitura, sem escrever código ainda.

**Ao terminar, você vai conseguir:**

- Contar quem criou Go, quando e qual problema concreto eles estavam tentando resolver.
- Explicar por que existe uma linguagem nova se já havia dezenas.
- Citar três programas que você usa ou administra que são escritos em Go.
- Dizer em que Go é bom e em que ele **não** é, para não usá-lo onde não convém.
- Situar Go entre as demais linguagens que você já conhece ou já ouviu falar.

---

## Por que isso importa

É 2007. O Google tem uma das maiores bases de código do mundo, escrita em sua maior parte em C++. E
tem um problema muito concreto e muito chato: **compilar leva uma eternidade.**

Uma única mudança em um arquivo de cabeçalho podia disparar uma recompilação de **45 minutos ou mais**.

Pense no que isso significa para quem programa. Você muda uma linha. Espera 45 minutos. Descobre que
errou uma vírgula. Muda outra linha. Espera mais 45 minutos. **Em um dia de trabalho você consegue
tentar oito coisas.**

Em 21 de setembro de 2007, três engenheiros do Google pararam diante de um quadro-branco para projetar
uma linguagem que não os fizesse esperar. Essa linguagem é o tema de todo este curso, e o resto desta
lição explica como eles chegaram daquele quadro-branco à linguagem que você vai instalar na lição 1.

---

## Os conceitos

### 0.1 Quem a criou (e por que isso importa)

Não eram três programadores quaisquer:

| | Quem é |
|---|---|
| **Ken Thompson** | **Criou o Unix.** E criou a linguagem B, antecessora direta de C, que escreveu junto com Dennis Ritchie. Prêmio Turing —o equivalente ao Nobel na computação— em 1983 |
| **Rob Pike** | Trabalhou no Unix e criou o Plan 9, seu sucessor. **Co-inventou o UTF-8**, a forma como hoje se guarda texto em praticamente todo o mundo, inclusive este documento |
| **Robert Griesemer** | Trabalhou no motor JavaScript V8 (o do Chrome e do Node.js) e na máquina virtual do Java |

🔑 **Leia isso de novo: um dos criadores de Go é o criador do Unix e coautor de C.** Cinquenta anos
depois de ter inventado as ferramentas sobre as quais todo o software moderno é construído, a mesma
pessoa projetou Go.

Isso explica muito do caráter da linguagem: Go parece C —direto, pequeno, sem enfeites— mas sem as
armadilhas que faziam de C um campo minado.

### 0.2 Os três requisitos que ninguém cumpria

Os três queriam uma linguagem com **três** propriedades ao mesmo tempo:

1. **Que compile rápido.** O problema original.
2. **Que execute rápido.** O Google roda software em milhares de máquinas: velocidade custa dinheiro de verdade.
3. **Que seja fácil de programar.** Que uma pessoa nova na equipe consiga ler o código dos outros e
   entendê-lo no primeiro dia.

Eles examinaram as linguagens que já existiam. E aí está a descoberta interessante: **encontraram
linguagens que cumpriam dois dos três, mas nenhuma que cumprisse os três.**

| | Compila rápido | Executa rápido | Fácil de programar |
|---|---|---|---|
| **C / C++** | ❌ | ✅ | ❌ |
| **Java** | 🟡 | ✅ | 🟡 |
| **Python** | ✅ (não compila) | ❌ | ✅ |
| **Go** | ✅ | ✅ | ✅ |

> [!NOTE]
> 🔧 **Observação de engenharia de software 0.1**
> Uma linguagem de programação não é questão de gosto: é uma ferramenta com compromissos. Cada uma
> sacrificou algo para ganhar outra coisa. Entender **o que sacrificou** aquela que você está usando é a
> diferença entre programar com critério e programar de cor.

### 0.3 Como a simplicidade foi conquistada: tirando coisas

Aqui está a decisão mais contraintuitiva de Go, e a que mais o distingue.

Quase todas as linguagens crescem: cada versão acrescenta recursos. C++ acumulou décadas deles, e o
resultado é tão grande que **ninguém o conhece por completo**. Há gente que programa em C++ há vinte
anos e continua encontrando cantos que não conhecia.

**Go fez o contrário: decidiu o que deixar de fora.** Ele não tem:

- **Classes nem herança.** A forma de reutilizar código é outra, e você vai vê-la na lição 3.
- **Exceções** (`try / catch`). Os erros são tratados de outra maneira, e você vai vê-la na lição 2.
- **Sobrecarga de operadores.** `+` significa somar, sempre, e não dá para mudar isso.
- **Construtores nem destrutores.**
- **Um monte de formas de fazer a mesma coisa.** Quase sempre há **uma** maneira idiomática.

⚠️ **Isso vai te incomodar alguma vez.** Você vai querer fazer algo que em outra linguagem se faz em uma
linha e em Go vai te custar cinco. A compensação é enorme e aparece depois de meses: **você consegue ler
código Go escrito por qualquer pessoa e entendê-lo.** Não há dialetos, não há truques, não é preciso
aprender como cada equipe programa.

> 🔑 **A citação que resume a filosofia**, de Rob Pike:
> *«Clareza é melhor do que esperteza.»*

### 0.4 A outra razão: os computadores pararam de acelerar

Há um segundo motivo por trás de Go, e é de hardware.

Até cerca de 2005, a cada ano saíam processadores mais rápidos e o seu programa rodava mais rápido
**sem que você mexesse em nada**. Isso acabou: os processadores deixaram de subir de velocidade e
passaram a multiplicar **núcleos**. O seu notebook não tem um processador rapidíssimo: tem 4, 8 ou 16
processadores modestos.

**O problema:** um programa normal usa **um** núcleo. Os outros quinze estão ali, só olhando.

Para aproveitá-los é preciso escrever programas que façam várias coisas ao mesmo tempo, e isso —a
**concorrência**— era tradicionalmente uma das tarefas mais difíceis e sujeitas a erros de toda a
programação.

**Go foi projetado em torno disso.** Fazer duas coisas acontecerem ao mesmo tempo em Go é literalmente
escrever uma palavra:

```go
go revisarServicio()
```

Essa palavra `go` —a que dá nome à linguagem— é o recurso que a tornou famosa. Você vai aprendê-lo na
lição 6, e quando chegar lá vai entender por que tanta gente migrou para esta linguagem.

### 0.5 A linha do tempo

| Data | O que aconteceu |
|---|---|
| **21-set-2007** | Griesemer, Pike e Thompson projetam Go em um quadro-branco |
| meados de 2008 | já existe um compilador que funciona |
| **10-nov-2009** | O Google o anuncia publicamente e o libera como software livre |
| **28-mar-2012** | **Go 1.0**, a primeira versão estável |
| 2012 em diante | duas versões grandes por ano, todo fevereiro e agosto |
| hoje | Go 1.27 |

🔑 **E algo que vale saber: a promessa de compatibilidade do Go 1.** Desde 2012, a equipe de Go
prometeu que **um programa escrito para o Go 1.0 continua compilando hoje**, catorze anos depois. Eles
cumpriram.

Isso não é normal. Em muitas linguagens, código de cinco anos atrás já não compila. **Em Go, o que você
aprender hoje vai servir daqui a dez anos**, e isso faz o tempo que você vai investir valer muito mais a
pena.

### 0.6 O que está escrito em Go (provavelmente você já usa)

Esta não é uma lista de propaganda: serve para você ver **onde** esta linguagem é usada, porque isso diz
muito sobre para que ela serve.

| | O que é |
|---|---|
| **Docker** | a ferramenta que empacota aplicações em contêineres |
| **Kubernetes** | o sistema que administra milhares de contêineres em servidores. Metade da indústria o usa |
| **Terraform** | cria infraestrutura na nuvem escrevendo arquivos |
| **Prometheus** · **Grafana** | coletam e plotam métricas de servidores |
| **Traefik** · **Caddy** | servidores web e balanceadores de carga |
| **CockroachDB** · **InfluxDB** | bancos de dados |
| **Hugo** | gerador de sites estáticos, famoso pela velocidade |
| **ngrok**, **rclone**, **gh** (o CLI do GitHub) | ferramentas de linha de comando |

🔑 **Percebe o padrão?** Quase tudo são **ferramentas de infraestrutura**: coisas que rodam em
servidores, que precisam ser rápidas, que são implantadas como um único arquivo e que lidam com muitas
conexões ao mesmo tempo. **É aí que Go ganha**, e não por acaso: é exatamente o problema que o Google
tinha.

### 0.7 E então, por que você deveria aprendê-lo?

Quatro razões honestas:

**1. Aprende-se rápido.** A especificação de Go se lê em uma tarde. A de C++ tem mais de 1.800
páginas. **Você vai conseguir escrever programas úteis em semanas, não em anos**, e isso importa muito
quando você está começando.

**2. Ele ensina coisas que servem em qualquer linguagem.** Por ser compilado e de tipagem estrita, Go
obriga você a pensar em tipos, em memória e em erros. Esses conceitos **se transferem**: se depois você
programar em Java, C# ou Rust, já os terá.

**3. Há trabalho.** Tudo o que roda em servidores modernos tem Go dentro. Se você se interessa por
infraestrutura, nuvem, DevOps ou backend, é uma das apostas mais seguras.

**4. O que você aprender não vai caducar.** Pela promessa de compatibilidade do Go 1.

### 0.8 O que você vai construir neste curso

Um programa de linha de comando chamado **`revisor`**: recebe uma lista de serviços, consulta **todos ao
mesmo tempo** e produz um relatório.

```
$ revisor --config servicios.txt --formato tabla
SERVICIO    ESTADO    TIEMPO  DETALLE
catalogo    OK         142ms  200
inventario  LENTO       2.3s  200
pagos       OK          87ms  200
reportes    FALLA          —  connection refused
```

Parece pequeno. **Não é.** Para escrevê-lo bem você vai precisar de tudo: tipos, structs, erros,
interfaces, listas, concorrência, HTTP, arquivos de configuração, testes e compilar um executável.

E ele vai crescer com você: **cada lição acrescenta uma peça.** No final você terá um programa que de
fato poderia usar, não um exercício de livro.

---

## O erro que você vai ver

Esta lição ainda não tem código próprio, então não há uma mensagem de compilador para mostrar. Mas há
sim um erro de **raciocínio** que quase todo mundo comete ao ler a seção 0.3, e vale a pena adiantá-lo
para que ele não te trave quando acontecer com você.

**O erro:** ler a lista do que Go **não** tem —classes, herança, exceções, sobrecarga de operadores— e
concluir *«então é uma linguagem pobre, faltam coisas de que preciso».*

**Por que é um erro:** confunde *ter menos ferramentas* com *conseguir resolver menos problemas*. Go não
tirou de você a capacidade de reutilizar código nem de tratar erros: deu a você **uma** forma de fazê-lo
em vez de dez, e essa forma você vai aprender nas lições 2 e 3. A frustração é real e é normal —você
também vai senti-la na primeira semana—, mas ela se resolve programando, não evitando a linguagem.
Quando terminar a lição 3, você vai poder reler a seção 0.3 e ver por que cada linha dessa lista é uma
decisão, não uma carência.

---

## O que se faz errado

Tão importante quanto saber para que serve uma ferramenta é saber para que ela não serve. **Nenhuma
linguagem é boa para tudo**, e quem disser o contrário está te vendendo algo. Usar Go onde não convém é
o primeiro antipadrão deste curso, antes mesmo de você ter escrito uma única linha:

| Não é a melhor opção para | O que se usa no lugar | Por quê |
|---|---|---|
| Apps de iPhone ou Android | Swift, Kotlin | as plataformas foram feitas para elas |
| Páginas web (o que roda no navegador) | JavaScript, TypeScript | o navegador só executa JavaScript |
| Ciência de dados, inteligência artificial | Python | todas as bibliotecas do mundo estão lá |
| Videogames grandes | C++, C# | precisam de controle absoluto da memória e do hardware |
| Sistemas em que um microssegundo importa | C, C++, **Rust** | Go tem um coletor de lixo que às vezes pausa o programa |

> [!NOTE]
> 🔧 **Observação de engenharia de software 0.2**
> A última linha é a razão pela qual este curso tem um **segundo curso, de Rust**. Rust resolve
> exatamente isso: velocidade de C sem coletor de lixo e sem os erros de memória de C. São ferramentas
> para problemas diferentes, e você vai aprender as duas.

---

## Exercícios

### Perguntas de revisão

**0.1** Que problema concreto motivou a criação de Go, e em que ano?

**0.2** Cite os três criadores de Go e diga por que é relevante quem é Ken Thompson.

**0.3** Quais eram os três requisitos que buscavam e por que nenhuma linguagem existente servia?

**0.4** Cite três coisas que Go **não** tem, de propósito. O que se ganha ao tirá-las?

**0.5** O que mudou no hardware por volta de 2005 e o que isso tem a ver com Go?

**0.6** O que é a promessa de compatibilidade do Go 1 e por que ela é boa para você?

**0.7** Cite três programas escritos em Go e diga o que têm em comum.

**0.8** Dê dois casos em que você **não** usaria Go, e diga o que usaria.

### Para ir além

**0.9** Procure na internet a palestra ou o artigo original *«Go at Google: Language Design in the Service
of Software Engineering»*, de Rob Pike. Leia a introdução e anote **uma** razão de projeto que não esteja
nesta lição.

**0.10** Escolha **dois** dos programas da seção 0.6 que você não conheça. Descubra em uma frase o que
cada um faz e anote no seu diário de bordo.

**0.11** Procure o índice da especificação da linguagem Go (*The Go Programming Language
Specification*) e conte quantas páginas ou seções ela tem. Compare com o padrão do C++. Anote os dois
números: é a forma mais concreta de ver o que significa «simples».

**0.12 (Para pensar, sem resposta correta)** Go tirou as exceções, as classes e a herança —coisas que
outras linguagens consideram indispensáveis. Ocorre a você alguma razão pela qual **tirar** um recurso
possa tornar uma linguagem melhor? Escreva sua opinião no diário de bordo **antes** de começar o curso, e
releia-a ao terminar. É interessante ver se mudou.

### Soluções

**0.1** Os tempos de compilação do C++ no Google: uma mudança em um cabeçalho podia custar **45 minutos**
de recompilação. Começaram em **21 de setembro de 2007**.

**0.2** Robert Griesemer, Rob Pike e **Ken Thompson**. Thompson **criou o Unix** e a linguagem B,
antecessora de C, que escreveu com Dennis Ritchie. Ou seja: um dos autores das ferramentas sobre as quais
foi construído o software moderno projetou também Go, cinquenta anos depois.

**0.3** Compilação rápida, execução rápida e facilidade de programar. As linguagens existentes cumpriam
**dois dos três**: C++ era rápido na execução, mas lento na compilação e difícil; Python era fácil, mas
lento na execução.

**0.4** Classes e herança, exceções, sobrecarga de operadores, construtores. **Ganha-se legibilidade**:
código Go escrito por qualquer pessoa pode ser lido e entendido, porque não há dialetos nem muitas formas
de fazer a mesma coisa.

**0.5** Os processadores pararam de acelerar e passaram a multiplicar **núcleos**. Um programa normal usa
apenas um, então fazia falta uma linguagem em que escrever programas concorrentes fosse fácil. Daí a
palavra `go`.

**0.6** A promessa de que um programa escrito para o Go 1.0 (2012) **continua compilando hoje**. É boa
para você porque o que você aprender não caduca e o código que escrever vai continuar funcionando por
anos.

**0.7** Docker, Kubernetes, Terraform, Prometheus, Grafana, Hugo… Todos são **ferramentas de
infraestrutura**: rodam em servidores, são distribuídas como um executável único e lidam com muitas
conexões ao mesmo tempo.

**0.8** Apps móveis (Swift/Kotlin), código de navegador (JavaScript), ciência de dados (Python),
videogames grandes (C++), sistemas de tempo real estrito (C ou Rust).

**0.9-0.12** Não têm uma única resposta correta: são de pesquisa e reflexão própria. Compare o que você
anotou com um colega ou no diário de bordo do curso, não com uma resposta fixa.

---

## Como sei que consegui

Não há código para compilar nesta lição, então o checklist é de compreensão. Marque cada item somente se
você conseguir fazê-lo **sem voltar a ver o texto**:

- ☐ Você explica com suas palavras, em menos de um minuto, por que Go nasceu e que problema resolvia.
- ☐ Você cita os três criadores e diz por que Ken Thompson ser um deles não é um dado de enchimento.
- ☐ Você diz o que são a compilação rápida, a execução rápida e a facilidade de programar, e por que
  nenhuma linguagem anterior a Go tinha as três.
- ☐ Você cita, de memória, pelo menos três coisas que Go decidiu não ter, e explica o que se ganha ao
  tirá-las.
- ☐ Você explica a relação entre os núcleos de um processador, a concorrência e a palavra `go`.
- ☐ Você cita três programas escritos em Go e diz o que têm em comum.
- ☐ Você dá um exemplo de problema para o qual **não** usaria Go, e diz o que usaria no lugar.
- ☐ Você resolveu as oito perguntas de revisão sem ver as soluções antes de comparar.

Se algum ficou faltando, não passe ainda para a lição 1: volte a ler a seção correspondente. Tudo o que
vem depois se apoia nisto.

---

## Resumo

- Go nasceu em **2007** no Google pela frustração com os **tempos de compilação do C++**: 45 minutos por
  uma mudança em um cabeçalho.
- Foi criado por **Robert Griesemer, Rob Pike e Ken Thompson**; Thompson criou o Unix e coescreveu C.
- Buscavam **três** propriedades juntas —compilar rápido, executar rápido, ser fácil— que nenhuma
  linguagem da época cumpria ao mesmo tempo.
- Go alcançou a simplicidade **tirando** coisas: sem classes, sem herança, sem exceções, sem sobrecarga
  de operadores.
- A segunda motivação foi o hardware: os processadores multiplicaram **núcleos** em vez de acelerar, e
  aproveitá-los exigia que a **concorrência** fosse fácil. Daí a palavra `go`.
- Foi anunciado em **2009** e a versão **1.0** saiu em **2012**; hoje está no 1.27, com duas versões
  grandes por ano.
- A **promessa de compatibilidade do Go 1** garante que o código de 2012 continua compilando: o que você
  aprende não caduca.
- Go domina a **infraestrutura**: Docker, Kubernetes, Terraform, Prometheus, Grafana.
- **Não** é a melhor opção para mobile, navegador, ciência de dados, videogames nem tempo real estrito.

---

## Para ler mais

1. **A especificação da linguagem** — *The Go Programming Language Specification*, em `go.dev/ref/spec`.
   É a fonte oficial e, como você viu no exercício 0.11, se lê em uma tarde. Não é preciso entendê-la
   toda hoje; vale a pena saber que existe e que é curta.
2. **O blog oficial de Go** — `go.dev/blog`, onde a própria equipe publica as novidades de cada versão
   e artigos de projeto.
3. **Rob Pike, «Go at Google: Language Design in the Service of Software Engineering»** — a
   palestra/artigo original em que um dos criadores explica as decisões de projeto com o contexto do
   Google. É a leitura do exercício 0.9, e a mais direta para entender o **porquê** de cada decisão que
   você viu nesta lição.

### Termos desta lição

| | |
|---|---|
| **compilar** | traduzir código-fonte em instruções de máquina |
| **concorrência** | fazer um programa realizar várias coisas ao mesmo tempo |
| **cabeçalho (*header*)** | em C/C++, arquivo com declarações que outros arquivos incluem |
| **Go 1 (promessa de compatibilidade)** | compromisso de que o código antigo continue compilando |
| **núcleo (*core*)** | cada processador independente dentro de um mesmo chip |
| **Prêmio Turing** | o maior reconhecimento em ciências da computação |
| **coletor de lixo** | parte da linguagem que libera memória automaticamente |
| **UTF-8** | forma padrão de representar texto, co-inventada por Rob Pike e Ken Thompson |

---

**Próxima:** [Lição 1 — Instalar Go no seu Linux Mint](01-instalacion.md)
