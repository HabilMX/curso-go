# Lição 1 — Instalar Go no seu Linux Mint e seu primeiro programa

**Duração:** 45 a 60 minutos (ou 2 sessões de 30). É a única lição que não ensina Go como linguagem:
ensina a ferramenta e o terreno onde você vai trabalhar — e esse terreno tem mais armadilhas do que
parece à primeira vista.

**Ao terminar, você vai conseguir:**

- Instalar a versão oficial e vigente de Go no Linux Mint, sem depender do pacote do sistema, e
  explicar por que esse pacote está obsoleto por projeto, não por descuido.
- Explicar o que são `GOROOT`, `GOPATH` e o `PATH`, e por que Go não funciona até você mexer no segundo.
- Reconhecer a diferença entre um terminal que lê a sua configuração e um que não lê —a causa real de
  «eu adicionei e mesmo assim não funciona»— e saber qual é o seu.
- Criar um projeto com `go mod init`, escrever um programa e executá-lo com `go run`.
- Produzir um binário com `go build` e explicar, com uma prova real, por que ele roda em outra máquina
  sem ter Go instalado.
- Usar `go fmt` para nunca mais discutir estilo de código.
- Reconhecer, palavra por palavra, as mensagens de erro mais comuns desta etapa —de compilador, de
  permissões e de caminho— e saber o que fazer com cada uma.

---

## Por que isso importa

A primeira coisa que um tutorial antigo vai te sugerir é `apt install golang`, e isso é um erro. Não é um
erro pequeno de "versão um pouco velha": é o tipo de erro que te faz perder uma tarde inteira sem
entender por quê, semanas depois de instalar, porque o sintoma não aparece ao instalar — aparece quando
você já está programando e algo que "deveria funcionar" não funciona.

Isto não é opinião, é medição. Em 30-set-2026 medi, em um contêiner com a mesma base do Linux Mint 22.3
(Ubuntu 24.04 "Noble Numbat", de onde o Mint tira seus pacotes):

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

**O pacote do sistema oferece a 1.22. A versão oficial vigente, medida no mesmo instante, é a 1.27.1.**
E isso não é um descuido que alguém vá corrigir: é a política do Ubuntu e do Mint.

**Por que se congela assim, de propósito?** O Ubuntu 24.04 é uma versão LTS (*Long Term Support*): no dia
em que sai, seus pacotes principais são **congelados** e só recebem correções de segurança, nunca
versões novas do programa. É uma decisão correta para um servidor que não deve mudar de comportamento
só por instalar atualizações — mas significa que o compilador de uma linguagem que libera versões novas
duas vezes por ano fica fixo na que existia quando o Ubuntu 24.04 foi publicado (abril de 2024), para
sempre, enquanto essa versão do Ubuntu existir. O Linux Mint não tem pacotes próprios de Go: herda
literalmente os do Ubuntu, então herda também o congelamento.

**Por que isso importa para você, concretamente?** Porque você vai pesquisar coisas na internet, os
exemplos vão usar funções ou comportamentos que o seu Go 1.22 não tem, e os erros não vão dizer *«falta
a versão»*: vão dizer coisas que não fazem sentido para o que você está vendo na tela. É o tipo de
problema que parece erro seu e na verdade é uma defasagem de ferramenta — e acontece com você **antes**
de escrever a primeira linha de código, por isso esta lição existe antes de qualquer outra.

🔑 **A forma correta, e a que o próprio projeto Go recomenda:** baixar o pacote oficial diretamente de
`go.dev`, não do gerenciador de pacotes da sua distribuição. É isso que esta lição instala, passo a
passo, e o que vamos conferir em cada passo com a saída real dos comandos.

---

## Os conceitos

### 1.1 O que são GOROOT, GOPATH e por que antes importavam mais do que agora

Antes de instalar, vale a pena saber o que você vai ter depois, porque os nomes `GOROOT` e `GOPATH`
aparecem em quase qualquer erro de instalação que você procure na internet, e muitas respostas foram
escritas para uma versão de Go de dez anos atrás.

- **`GOROOT`** é a pasta onde vive **o próprio Go**: o compilador, o `gofmt`, a biblioteca padrão. É a
  pasta que você vai criar no passo 1.4 (`/usr/local/go`). Você nunca mexe nela à mão.
- **`GOPATH`** é a pasta onde Go guarda **coisas suas**: pacotes baixados da internet, e —em versões
  antigas de Go, antes de 2019— também onde você tinha de colocar **todo** o seu código, sem exceção,
  dentro de uma estrutura fixa (`$GOPATH/src/github.com/seu-usuario/seu-projeto`). Se alguma vez você
  vir um tutorial que pede para criar essa estrutura de pastas, ele é dessa época.

Hoje você quase não mexe nela porque, desde o Go 1.11 (2018), existem os **módulos** (`go.mod`, que você
vai criar no passo 1.6): seu projeto pode viver em qualquer pasta, com o nome que você quiser, e Go já
não precisa que você siga uma estrutura de pastas imposta. O `GOPATH` continua existindo, mas agora só
como cache de pacotes baixados, não como o único lugar onde o seu código pode viver.

Confira com `go env`, que mostra a configuração vigente da sua instalação (você pode rodar isso
**depois** de instalar, no passo 1.5):

```bash
go env GOROOT GOPATH GOBIN
```

Em uma instalação limpa, recém-feita, você vai ver algo como:

```
/usr/local/go
/home/tu-usuario/go
/home/tu-usuario/go/bin
```

Nenhuma dessas três pastas foi criada por você à mão: a primeira é criada pelo passo de instalação
(1.4), as outras duas o próprio Go decide, com valores padrão razoáveis.

### 1.2 O PATH: o que é, e por que "já está instalado" não significa "já funciona"

Quando você digita um comando no terminal, por exemplo `go`, o terminal não sabe magicamente onde está
esse programa: ele examina, uma por uma, uma lista de pastas guardada em uma variável chamada **PATH**, e
usa o primeiro programa que encontrar com esse nome. Se nenhuma pasta da lista tiver um programa chamado
`go`, responde com um erro — e esse erro é literal, não aproximado:

```bash
$ go version
bash: go: command not found
```

Verifiquei isso em um contêiner com Go **já extraído** em `/usr/local/go`, antes de mexer no PATH: o
programa existe no disco, mas o terminal não o encontra porque não sabe onde procurar. **"Instalado" e
"no PATH" são duas coisas diferentes**, e a confusão entre elas é a fonte de quase todos os tropeços
desta lição.

Isso também explica por que a mensagem muda um pouco conforme o terminal (`zsh` em vez de `bash`, que
você verá se usar macOS para acompanhar o curso enquanto pratica, ou se trocou o terminal padrão do seu
Mint):

```bash
$ go version
zsh: command not found: go
```

Mesmo problema, mesmo mecanismo, **ordem das palavras diferente**: identifique qual é o seu caso com
`echo $SHELL` antes de procurar o erro na internet, porque procurar a mensagem errada vai te levar a
respostas para o shell errado.

### 1.3 Descubra qual é a última versão

Não a copie daqui: este documento envelhece como os repositórios do Mint, e você já viu na seção
anterior o quanto pode pesar uma ferramenta ficar fixa no tempo.

```bash
curl -s 'https://go.dev/VERSION?m=text' | head -1
```

Ele vai responder algo como `go1.27.1`. **Essa é a que você vai instalar**, seja qual for quando você
fizer isso.

### 1.4 Baixe-a e instale-a

Substitua `go1.27.1` pelo que o comando anterior disse. `linux-amd64` é o correto para um PC ou notebook
normal (se a sua máquina fosse ARM, seria `linux-arm64`; para saber: `dpkg --print-architecture`).

```bash
cd /tmp
wget https://go.dev/dl/go1.27.1.linux-amd64.tar.gz
```

E agora a instalação:

```bash
sudo rm -rf /usr/local/go                          # borra una instalación anterior, si había
sudo tar -C /usr/local -xzf go1.27.1.linux-amd64.tar.gz
```

**O que você acabou de fazer**, porque convém entender e não apenas copiar:

- `tar` descompacta o arquivo.
- `-C /usr/local` diz a ele *«faça isso dentro desta pasta»*.
- `-xzf` é «extrair» (`x`), «está comprimido com gzip» (`z`), «deste arquivo» (`f`).
- O resultado é uma pasta `/usr/local/go` com todo o Go dentro — o `GOROOT` da seção 1.1.

**Por que `sudo`, se você nunca precisou dele para instalar algo com um gerenciador de pacotes?** Porque
`/usr/local` é uma pasta do sistema, não do seu usuário, e no Linux escrever ali exige permissões de
administrador. Sem `sudo`, é exatamente isto que você vai ver —provoquei de propósito, sem `sudo`, para
que você veja a mensagem real e não uma inventada—:

```
$ tar -C /usr/local -xzf go1.27.1.linux-amd64.tar.gz
tar: go: Cannot mkdir: Permission denied
tar: go/VERSION: Cannot open: No such file or directory
tar: go/api: Cannot mkdir: No such file or directory
tar: go/api/README: Cannot open: No such file or directory
... (se repite, una vez por cada archivo del paquete)
```

**Não é um erro, são centenas.** O `tar` tenta escrever cada arquivo do pacote, um por um, e cada um
falha do mesmo jeito porque nenhum tem permissão para escrever em `/usr/local`. Se você vir essa parede
de linhas repetidas, a causa é sempre a mesma e sempre se resolve do mesmo jeito: coloque `sudo` antes.

### 1.5 Diga ao seu sistema onde ele está — e a armadilha dos terminais modernos

Agora o Go existe em `/usr/local/go/bin/go`, mas, como você viu na seção 1.2, se digitar `go` ele vai
dizer `command not found`: essa pasta ainda não está no PATH.

A instrução que você vai ver em quase qualquer tutorial, inclusive na documentação oficial de Go, é
esta:

```bash
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ~/.profile
```

- `~/.profile` é um arquivo que, em teoria, o seu sistema lê **toda vez que você inicia sessão**.
- `source ~/.profile` aplica a mudança **neste terminal, agora mesmo**, sem esperar a próxima sessão.

E aqui vem o que a maioria dos tutoriais não conta, e que conferi com um teste real: **um terminal novo
nem sempre é uma "sessão nova".**

O `~/.profile` só é lido pelo que o Linux chama de **shell de login** — a que inicia quando você começa
uma sessão no sistema (por exemplo, ao ligar a máquina e entrar com seu usuário e senha). Mas a maioria
dos aplicativos de terminal (o terminal do Mint incluído, na configuração padrão) ao abrir uma janela ou
aba nova **não** iniciam uma shell de login: iniciam uma shell interativa normal, e essas leem outro
arquivo, `~/.bashrc`, **não** o `~/.profile`.

Conferi assim, simulando exatamente esse cenário —"já editei o arquivo, fecho o terminal, abro outro"—
em um contêiner recém-instalado:

```
# después de agregar el export SOLO a ~/.profile, en una terminal nueva:
$ go version
bash: go: command not found          ← sigue sin funcionar

# después de agregar la MISMA línea también a ~/.bashrc, en una terminal nueva:
$ go version
go version go1.27.1 linux/arm64      ← ahora sí
```

**A instrução "feche o terminal, abra outro" nem sempre basta**, e quando não basta, parece que você
fez algo errado quando na verdade seguiu o tutorial à risca. Por isso a recomendação desta lição, mais
robusta que a da maioria dos guias, é adicionar a linha **aos dois arquivos**:

```bash
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

O `~/.profile` cobre o caso de uma sessão de login de verdade (por exemplo, se você usa SSH ou troca de
usuário); o `~/.bashrc` cobre o caso, muito mais comum no dia a dia, de abrir uma janela ou aba de
terminal nova sobre uma sessão que já estava aberta.

💡 **Se você usa `zsh`** em vez de `bash` (você sabe se o seu terminal parece diferente ou se trocaram
para você), o arquivo equivalente ao `.bashrc` é o `~/.zshrc`. Para saber qual shell você usa:
`echo $SHELL`.

🔑 **Como distinguir se o seu terminal abre uma shell de login, sem adivinhar:** rode `echo $0` logo
depois de abri-lo. Se a resposta começar com um hífen (`-bash`, `-zsh`), é uma shell de login e sim lê o
`~/.profile`. Se não tiver o hífen (`bash`, `zsh`), não é, e você precisa da mudança no `~/.bashrc` (ou
`~/.zshrc`) para que persista.

### 1.6 Confira que funcionou

```bash
go version
```

Deve responder algo como `go version go1.27.1 linux/amd64` (ou `linux/arm64` se a sua máquina for ARM).
Se continuar dizendo `command not found`, verifique com `echo $0` de que tipo de shell se trata e em
qual arquivo você colocou a mudança, seguindo a seção 1.5.

### 1.7 Seu primeiro programa

```bash
mkdir -p ~/w/curso-go/hola && cd ~/w/curso-go/hola
go mod init hola
```

`go mod init` cria um arquivo `go.mod`. É a ficha de identidade do projeto: diz como ele se chama e qual
versão de Go usa. É o mecanismo de módulos que já mencionamos na seção 1.1, e é o mesmo com o qual você
vai iniciar o projeto `revisor` na lição 5.

Crie o `main.go` com isto:

```go
package main

import "fmt"

func main() {
    fmt.Println("hola, ya tengo Go")
}
```

Linha por linha, porque cada uma tem sua razão:

| | |
|---|---|
| `package main` | *«este arquivo pertence ao pacote `main`»*. **O pacote `main` é especial: é o único que produz um programa executável.** Sem esta linha você teria uma biblioteca, não um programa |
| `import "fmt"` | *«vou usar coisas do pacote `fmt`»* (de *format*), que traz o necessário para imprimir. Go **não** traz nada carregado por padrão: o que você usar, você pede |
| `func main()` | **a função onde o seu programa começa.** Quando você o executa, Go procura exatamente esta função. Se tiver outro nome, ele não inicia |
| `fmt.Println(...)` | imprime e pula de linha. O ponto significa *«a função `Println` que está dentro de `fmt`»* |

Rode:

```bash
$ go run main.go
hola, ya tengo Go
```

Essa saída é real: rodei antes de escrever esta linha.

### 1.8 Os comandos que você vai usar sempre — e a sutileza do `go.mod`

```bash
go run main.go     # compila y ejecuta de una vez, sin dejar archivo. Para probar mientras trabajas
go build           # crea el programa ejecutable y lo deja ahí
go fmt ./...       # ordena tu código
go test ./...      # corre las pruebas (lección 5)
```

**Uma sutileza real que vale a pena testar você mesmo, porque mudou entre versões de Go:** com um
programa de um único arquivo como o de cima, `go run main.go` funciona **mesmo sem `go.mod`**. Conferi
em uma pasta nova, sem `go mod init`:

```
$ go run main.go
hola, ya tengo Go
```

Mas o `go build`, nessa mesma pasta sem `go.mod`, sim exige o módulo:

```
$ go build
go: go.mod file not found in current directory or any parent directory; see 'go help modules'
```

**Por que a diferença?** O `go run` sobre um único arquivo consegue resolver tudo sem precisar saber o
nome do módulo, porque não há nada que outro arquivo possa importar dele. Assim que o seu programa
precisa de mais de um arquivo — por exemplo, se você tentar importar um pacote próprio, como vai fazer
na lição 5 com o `revisor` — Go precisa do `go.mod` para saber como o seu módulo se chama e assim
resolver esses imports. Também conferi:

```
$ go run main.go            # main.go importa "hola/utilidades", sin go.mod
main.go:4:5: package hola/utilidades is not in std (/usr/local/go/src/hola/utilidades)
```

**A lição prática:** rode `go mod init` sempre, desde o primeiro arquivo, mesmo que o `go run` às vezes
funcione sem ele — você evita este erro exato assim que o seu programa precisar importar um pacote
próprio, o que vai acontecer já na lição 5.

**Teste a diferença entre `run` e `build`:**

```bash
$ go build
$ ls -la
-rwxr-xr-x 1 tu-usuario tu-usuario 2341973 ... hola
$ ./hola
hola, ya tengo Go
```

### 1.9 Por que o binário é portátil: linkagem estática

🔑 **O arquivo `hola` que o `go build` produziu é um programa completo e autossuficiente.** E isso não
precisa ser acreditado por fé: conferi copiando exatamente esse binário, compilado em uma máquina com Go
instalado, para um contêiner **limpo, sem Go**:

```
$ which go
go no esta instalado
$ ./hola
hola, ya tengo Go
```

Funcionou. **Não precisa de interpretador, nem de máquina virtual, nem de bibliotecas externas instaladas
à parte**, porque Go faz *linkagem estática*: em vez de dizer «quando rodar, procure a biblioteca `fmt`
em algum lugar do sistema» (que é como funcionam muitos programas em C ou em Python), Go copia para
dentro do próprio binário tudo o que o seu programa precisa para rodar. O arquivo fica maior por isso
—pouco mais de 2 MB para um programa de uma linha— mas em troca você pode copiá-lo para qualquer máquina
Linux compatível e ele vai funcionar do jeito que está, sem instalar mais nada.

Essa é a característica mais prática de Go, e por isso tantas ferramentas de servidores são escritas
nele — inclusive o `revisor` que você vai compilar e publicar na lição 7.

### 1.10 Um editor que te ajude

Não é obrigatório, mas vai te poupar muito tempo. Se você usa o **VS Code**, instale a extensão oficial
**Go** (de `golang.go`). Medido em 30-set-2026 no Marketplace do Visual Studio: a versão publicada é a
**0.57.2**. Esse número também vai envelhecer — o relevante é que você instale a que o Marketplace
oferecer no dia em que fizer isso, não que anote este número.

A extensão vai sublinhar os erros enquanto você escreve, em vez de você descobri-los ao compilar, e deixa
você pular para a definição de qualquer coisa com `F12`. Na primeira vez ela vai pedir para instalar
algumas ferramentas extras (`gopls`, o servidor de linguagem de Go, entre outras): diga que sim a todas.

---

## O erro que você vai ver

Todas estas são mensagens reais, provocadas de propósito para esta lição — não descrições aproximadas do
que "provavelmente" diriam:

| O sintoma | Mensagem literal | O que acontece e o que fazer |
|---|---|---|
| Terminal `bash` sem Go no PATH | `bash: go: command not found` | O programa não está em nenhuma pasta do PATH. Veja a seção 1.5 |
| Terminal `zsh` sem Go no PATH | `zsh: command not found: go` | Mesmo problema, ordem das palavras diferente porque é outro shell |
| `tar` sem `sudo` contra `/usr/local` | `tar: go: Cannot mkdir: Permission denied` (repetido por arquivo) | Faltou `sudo` antes do comando de extração (seção 1.4) |
| Você adicionou o PATH só ao `~/.profile` e abriu um terminal novo | `bash: go: command not found` (persiste) | Seu terminal não abre uma shell de login; adicione a linha também ao `~/.bashrc` (seção 1.5) |
| `go build` em uma pasta sem `go.mod` | `go: go.mod file not found in current directory or any parent directory; see 'go help modules'` | Rode `go mod init <nome>` nessa pasta |
| Import de um pacote próprio sem `go.mod` | `main.go:4:5: package hola/utilidades is not in std (...)` | Igual ao anterior: sem módulo, Go não sabe resolver seus próprios pacotes |
| Você esqueceu `import "fmt"` e usou `fmt.Println` | `# command-line-arguments`<br>`./main.go:4:2: undefined: fmt` | Go conhece todos os símbolos disponíveis; se você não importou o pacote, ele não existe para ele |
| `go version` responde **1.22** em vez da vigente | (sem erro, mas versão errada) | Está sendo usado o do `apt`. Remova-o com `sudo apt remove golang-go` e confirme que `/usr/local/go/bin` está no seu PATH |

**E a regra geral, que vale para os oito:** copie a mensagem de erro completa e pesquise-a tal qual, sem
resumi-la com suas palavras primeiro. Quase sempre alguém já passou por isso, e a mensagem quase sempre
diz exatamente o que falta — adivinhar sem lê-la por inteiro é a forma mais lenta de resolver.

---

## O que se faz errado

- **Instalar Go com o gerenciador de pacotes do sistema (`apt install golang`).** É a sugestão mais
  comum em tutoriais antigos, e é exatamente o problema da seção "Por que isso importa": uma versão LTS
  congelada te deixa anos atrasado sem avisar com nenhum erro, só com comportamentos estranhos mais
  adiante.
- **Dar por certo que "fechar e abrir o terminal" sempre recarrega a configuração.** É a armadilha
  medida na seção 1.5: se o seu terminal não abre uma shell de *login*, fechá-lo e abri-lo de novo não
  volta a ler o `~/.profile`. O sintoma —"fiz certo e não funciona"— não significa que você fez algo
  errado, significa que esse arquivo não era o que o seu terminal lê.
- **Extrair a versão nova por cima de uma instalação antiga, sem apagar antes.** O
  `sudo rm -rf /usr/local/go` do passo 1.4 não é decorativo: se você extrair por cima de uma instalação
  antiga, ficam arquivos das duas misturados e o resultado é um Go que falha de formas incompreensíveis.
  Apague primeiro, sempre.
- **Confiar que o `go run` funciona sem `go.mod` e por isso pular o `go mod init`.** É verdade para um
  arquivo solto (seção 1.8), mas deixa de ser assim assim que o seu programa tem mais de um arquivo — e
  para então você já escreveu código que vai ter de reorganizar. Rode `go mod init` desde o começo,
  sempre.
- **Discutir o estilo do código à mão.** Em Go o estilo não se discute: quem decide é a ferramenta.
  Experimente escrever isto de propósito, todo torto:

  ```go
  package main
  import "fmt"
  func main(){
  fmt.Println( "hola" )
  }
  ```

  Rode `go fmt ./...` e reabra o arquivo: **ficou organizado sozinho.** Não existem brigas sobre onde
  vai a chave ou quantos espaços leva a indentação, porque o `gofmt` tem uma única resposta e todo mundo
  a usa. Acostume-se a rodá-lo antes de salvar, em vez de formatar à mão.

---

## Exercícios

1. Instale Go seguindo as seções 1.3 a 1.6 e confirme a versão com `go version`.
2. Rode `echo $0` no seu terminal **antes** de mexer no PATH. Conforme a resposta, decida se você precisa
   mexer no `~/.profile`, no `~/.bashrc`, ou nos dois (seção 1.5) — e explique por quê, com suas
   palavras, antes de continuar.
3. Crie o projeto `hola` da seção 1.7, rode `go run main.go` e depois `go build` + `./hola`.
4. (Como na seção 1.9) Copie o seu binário `hola` para outra máquina Linux, ou para uma máquina virtual /
   contêiner sem Go instalado, e confirme que roda igual.
5. Troque a mensagem do `fmt.Println` por algo seu, bagunce o recuo de propósito, e rode
   `go fmt ./...`. Verifique que o arquivo ficou bem formatado.
6. (Um pouco mais difícil) Crie uma pasta nova, **sem** `go mod init`, com um `main.go` de um único
   arquivo. Confirme que `go run main.go` funciona mesmo assim. Depois rode `go build` nessa mesma
   pasta e compare o erro com o da seção 1.8. Explique, com suas palavras, por que um funciona e o outro
   não.
7. (Um pouco mais difícil) Apague de propósito a linha `import "fmt"` e rode `go run main.go`. Leia o
   erro completo, sem pesquisá-lo ainda, e tente explicar com suas palavras o que ele está te dizendo
   antes de continuar.

### Soluções

1. `go version` deve imprimir a mesma versão que você viu em `curl -s 'https://go.dev/VERSION?m=text'`,
   nunca `go1.22.x` (essa é a do `apt` do Mint).
2. Se `echo $0` responder com hífen no início (`-bash`), seu terminal abre shells de login e o
   `~/.profile` basta. Se responder sem hífen (`bash`), não é uma shell de login, e só o `~/.bashrc` (ou
   `~/.zshrc` no zsh) vai persistir entre janelas novas — adicione a mudança ali também, como na seção
   1.5.
3. `go run main.go` imprime `hola, ya tengo Go`. Depois do `go build`, aparece um arquivo `hola` na
   pasta; `./hola` o executa direto, sem compilar de novo.
4. O binário roda igual na máquina sem Go, porque Go faz linkagem estática (seção 1.9): tudo o que o
   programa precisa já está copiado dentro do próprio arquivo.
5. Antes do `go fmt`, o arquivo fica exatamente como foi escrito (torto). Depois, o `gofmt` o reorganiza
   com a indentação e os espaços padrão de Go — sem que você decida nada disso.
6. `go run main.go` funciona porque, com um único arquivo, Go não precisa resolver nenhum import próprio
   para compilá-lo e executá-lo de uma vez. O `go build`, por outro lado, exige saber o nome do módulo
   desde o primeiro momento e responde `go: go.mod file not found...`. A diferença é que o `build` deixa
   o binário pronto para ser usado fora dessa invocação pontual, e para isso precisa de uma identidade
   de projeto — o `run` não.
7. O compilador responde:

   ```
   # command-line-arguments
   ./main.go:4:5: undefined: fmt
   ```

   Diz que você usou `fmt.Println` sem ter importado o pacote `fmt`: o compilador conhece todos os
   símbolos que você pode usar, e se você não o importou, ele não existe para ele. Resolve-se devolvendo
   a linha `import "fmt"`.

---

## Como sei que consegui

- [ ] `go version` responde com a versão que vi em `go.dev/VERSION`, não com 1.22.
- [ ] Sei se o meu terminal abre shells de login ou não (`echo $0`), e em qual(is) arquivo(s) coloquei a
      mudança de PATH em consequência.
- [ ] `go run main.go` imprime a minha mensagem.
- [ ] `go build` produziu um arquivo e `./hola` funciona.
- [ ] Copiei o meu binário para outra máquina (ou contêiner) sem Go e rodou igual.
- [ ] Testei o `go fmt` sobre código mal escrito e ele o organizou.
- [ ] Provoquei o erro do `go build` sem `go.mod` e consigo explicar por que o `go run` não o teve.
- [ ] Provoquei o erro de `import` faltando e entendi a mensagem sem precisar da solução.
- [ ] Consigo explicar, sem ver o texto, que diferença há entre `GOROOT`, `GOPATH` e o `PATH`.

---

## Resumo

- Nunca instale Go com o gerenciador de pacotes do sistema (`apt install golang`): as versões LTS do
  Linux são congeladas, e o pacote fica anos atrás sem avisar com nenhum erro.
- **`GOROOT`** é onde Go vive; **`GOPATH`** é onde Go guarda pacotes baixados; nenhum dos dois você mexe
  à mão em um projeto moderno com módulos (`go.mod`).
- O **`PATH`** é a lista de pastas onde o terminal procura programas. Ter Go instalado e ter Go no
  `PATH` são duas coisas diferentes.
- Um terminal novo **nem sempre** é uma sessão de *login*: a maioria só lê o `~/.bashrc` (ou o
  `~/.zshrc`), não o `~/.profile`. Adicione a mudança do `PATH` aos dois arquivos.
- `go mod init` cria a identidade do projeto (`go.mod`); todo projeto Go precisa dela, mesmo que o
  `go run` às vezes funcione sem ela com um único arquivo.
- `go run` compila e executa sem deixar arquivo; `go build` deixa o binário; `go fmt` dá formatação
  automática, sem discussão possível sobre estilo.
- O binário que o `go build` produz é autossuficiente por **linkagem estática**: roda em outra máquina
  Linux compatível sem ter Go instalado.
- Cada mensagem de erro desta lição é literal, provocada de propósito: `command not found`,
  `Permission denied`, `go.mod file not found`, `undefined: fmt` — copie-as tal qual quando for
  pesquisá-las.

---

## Para ler mais

1. [Download and install](https://go.dev/doc/install) — o guia oficial de instalação, a fonte da
   verdade se algo nesta lição envelhecer. Inclui a mesma recomendação do `~/.profile` que usamos
   aqui, com o mesmo aviso de que a mudança "pode não ser aplicada até o próximo início de sessão" — que
   é justamente o que foi medido e explicado a fundo na seção 1.5.
2. [A Tour of Go](https://go.dev/tour/) — o passeio interativo oficial; começa exatamente onde esta
   lição termina.
3. [Effective Go](https://go.dev/doc/effective_go) — inclui a seção sobre o `gofmt` e por que a
   formatação não se discute em Go.
4. [GNU Bash — Manual: Bash Startup Files](https://www.gnu.org/software/bash/manual/html_node/Bash-Startup-Files.html) —
   a referência oficial de quando se lê o `~/.profile` e quando o `~/.bashrc`; vale a pena lê-la inteira
   uma única vez na vida, porque a confusão entre os dois não é exclusiva de Go.

### Termos desta lição

| Termo | O que significa |
|---|---|
| `GOROOT` | a pasta onde vive a instalação de Go: o compilador, o `gofmt`, a biblioteca padrão |
| `GOPATH` | a pasta onde Go guarda pacotes baixados; antes dos módulos, todo o seu código também tinha de viver ali |
| `PATH` | a lista de pastas onde o terminal procura programas executáveis |
| shell de *login* | a sessão de terminal que inicia ao começar uma sessão no sistema; é a única que lê o `~/.profile` |
| `go.mod` | o arquivo que identifica um projeto Go: seu nome e a versão de Go que ele usa |
| binário | o arquivo executável que o `go build` produz, autossuficiente, sem depender de ter Go instalado |
| linkagem estática | a técnica pela qual Go copia para dentro do binário tudo o que o programa precisa, em vez de procurá-lo no sistema ao rodar |
| `gofmt` | a ferramenta que dá formatação automática ao código Go, sem opções para discutir |

---

**Anterior:** [Lição 0 — O que é Go](00-introduccion.md) ·
**Próxima:** [Lição 2 — Variáveis, funções e tipos](02-fundamentos.md)
