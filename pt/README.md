# Curso de Go — do zero a escrever um bom programa

**Por Dorian Chávez, fundador da Hábil e arquiteto de integração.**

**Para quem é:** alguém que está começando — primeiro ano de universidade, ou recém-saído da
escola. **Não damos nada como certo**: nem que você tem Go instalado, nem que sabe o que é um ponteiro, nem o que
significa compilar. Cada conceito é explicado quando aparece, e se explica **por que** ele existe, não só
como se escreve.

**O que você vai saber ao terminar:** escrever um programa completo, entender o que escreveu, e conseguir
explicá-lo a outra pessoa. Essa última é a verdadeira prova.

**O que você precisa antes de começar:** um computador com Linux Mint e saber abrir um terminal. Nada
mais. A [Lição 1](01-instalacion.md) instala tudo do zero.

## O projeto que você vai construir

Um programa de linha de comando chamado **`revisor`**: recebe uma lista de serviços, consulta
**todos ao mesmo tempo**, e produz um relatório em tabela ou em JSON.

    $ revisor --config servicios.txt --formato tabla
    SERVICIO    ESTADO    TIEMPO  DETALLE
    catalogo    OK         142ms  200
    inventario  LENTO       2.3s  200
    pagos       OK          87ms  200
    reportes    FALLA          —  connection refused

Parece pequeno e não é: para escrevê-lo bem você precisa de structs, interfaces, erros, concorrência,
`context`, HTTP, JSON, flags de linha de comando, testes e um binário empacotado.

**E é o esqueleto do que as pessoas escrevem em Go de verdade.** Segundo a
[Go Developer Survey 2025](https://go.dev/blog/survey2025) —a pesquisa oficial da equipe de Go,
realizada em setembro de 2025 com **5,379 respostas**—, os dois tipos de projeto mais construídos são:

| O que constroem | |
|---|---|
| Ferramentas de linha de comando | **74 %** |
| Serviços API/RPC | **73 %** |
| Bibliotecas ou frameworks | 49 % |

O `revisor` é o primeiro, e na Lição 7 você acrescenta o segundo.

⚠️ **E você vai escrevê-lo outra vez em Rust**, no curso irmão. Escrever o mesmo programa nas duas
linguagens é o que de verdade ensina em que elas se diferenciam; ler comparações não serve.

## As oito lições

| Lição | | O que você constrói | O que você aprende |
|---|---|---|---|
| 0 | [O que é Go](00-introduccion.md) | nada ainda (leitura) | de onde Go veio, para que serve e para que não serve |
| 1 | [Instalação](01-instalacion.md) | o ambiente e o seu primeiro programa | `GOROOT`, `GOPATH`, o `PATH`, e os erros reais de instalar Go |
| 2 | [Variáveis, funções e tipos](02-fundamentos.md) | as funções base do `revisor` | o compilador, variáveis, tipos básicos, funções com `(resultado, error)` |
| 3 | [Structs, métodos, erros e interfaces](03-errores-interfaces.md) | o `revisor` com structs e interfaces | structs, métodos, ponteiros, erros como valores, interfaces |
| 4 | [Coleções: slices e maps](04-colecciones.md) | a lista de serviços e o relatório | slices, maps, `range`, ordem estável |
| 5 | [Módulos e testes](05-modulos-y-pruebas.md) | o projeto de verdade, com testes | `go mod`, tabela de casos, `-race`, cobertura |
| 6 | [Concorrência](06-concurrencia.md) | **que ele verifique tudo ao mesmo tempo** | goroutines, canais, `context`, `WaitGroup`, semáforo |
| 7 | [O programa terminado](07-el-programa.md) | CLI, HTTP, JSON e binário | `net/http`, `flag`, `encoding/json`, compilação cruzada |

## A ordem não é a óbvia, e é de propósito

Structs, **erros e interfaces vêm antes** de slices, maps e ponteiros. Soa estranho —quase todos os
cursos colocam as coleções primeiro— mas não é um capricho: é a ordem usada por dois dos cursos de Go
com mais alunos, o do [Boot.dev / freeCodeCamp](https://www.boot.dev/courses/learn-golang) (structs →
interfaces → erros nas posições 5, 6 e 7, **antes** de slices, maps e ponteiros) e o de
[Todd McLeod na Udemy](https://www.udemy.com/course/learn-how-to-code/).

A razão é boa: **em Go, modelar dados e tratar erros É a linguagem.** Um programa com slices
perfeitos e erros ignorados não é Go; é C com outra sintaxe.

**A concorrência fica na Lição 6, não na 3.** É a razão pela qual Go existe, mas você precisa de
structs, erros e interfaces para que os exemplos não sejam de brinquedo. E ela é dada em três etapas
—goroutines, depois canais, depois aplicada ao projeto— porque em uma só não assenta.

## Como usá-lo

- **Uma sessão de 90 minutos por semana**, ou duas de 45. Menos não assenta; mais de duas horas seguidas se
  esquece.
- 🔴 **Escreva o código à mão, sempre.** Copiar e colar produz arquivos, não conhecimento. Os seus dedos
  aprendem coisas que os seus olhos não aprendem.
- **Quebre o código de propósito.** Em várias práticas vamos te pedir que provoque um erro e leia o
  que o compilador diz. **Não é enchimento:** aprender a ler erros é metade de saber programar.
- **Não passe de lição com dúvidas em aberto.** Anote-as no diário de bordo e pergunte. Tudo o que vem a seguir se
  apoia no anterior.
- **Cada lição fecha com um «como sei que consegui» mensurável**: compila e passa, ou não.
- **O [`bitacora.md`](bitacora.md) é seu**: dúvidas, tropeços e o que te surpreendeu. É o que torna
  isto o seu curso.

## As três fontes que valem a pena, nesta ordem

1. **[A Tour of Go](https://go.dev/tour/)** — oficial, interativo, 2-3 h. **Completo, não folheado.**
2. **[Go by Example](https://gobyexample.com/)** — cada conceito com código mínimo. Como referência.
3. **[Effective Go](https://go.dev/doc/effective_go)** — a filosofia. **Leia-o na Lição 3, não no final.**

⚠️ **Aprenda a biblioteca padrão antes de qualquer framework.** Em Go a padrão basta para quase
tudo, e isso é uma decisão de projeto, não uma carência. Quem começa com Gin ou Echo aprende o
framework e não a linguagem.

---

## Quem escreve isto

Este curso é escrito por **Dorian Chávez**, CEO e consultor principal da **Hábil**, empresa mexicana de
engenharia de software. Engenheiro em Sistemas Computacionais pelo IPN (ESCOM), com 26 anos construindo
e operando sistemas em produção.

Nós o publicamos aberto, com licença [Creative Commons](../LICENSE.md), porque quem explica bem constrói
bem — e porque um curso que não pode ser copiado, traduzido nem melhorado serve para pouco.

**Se você encontrar um erro, uma explicação confusa ou um número que não bate, diga.** Este material foi
feito para ser corrigido. Antes de publicar cada versão, verificamos que todo o código do curso funciona
e que o que cada lição mostra é idêntico ao programa real, para que nenhuma lição ensine algo que
já mudou.

### Se isto foi útil para você

- **Quer trabalhar com gente que escreve assim?** Escreva para nós.
- **A sua equipe está adotando Go ou Rust?** Também.

*(Os links de contato são colocados pelo site ao publicar esta página.)*
