# Lição 4 — Coleções: slices e maps

**Duração:** 90 minutos (ou 2 sessões de 45).

**Ao terminar, você vai conseguir:**

- Distinguir um **array** de um **slice** e dizer por que você sempre vai usar o segundo.
- Explicar por que copiar um slice **não copia os dados**, e quando isso vai te morder.
- Usar um **map** e saber por que a sua iteração sai desordenada de propósito.
- Percorrer coleções com `for range` e saber o que significa o `_`.
- Ordenar resultados para que o seu programa produza sempre a mesma saída.
- Ler um arquivo e convertê-lo em uma lista de dados.

---

## Por que isso importa

**Esta é a lição em que quase todo mundo tropeça.** E não por ser difícil, mas porque há uma armadilha que
nenhum curso explica até que ela já tenha te mordido: em Go, copiar um slice não copia os seus dados, e esse
comportamento não produz nenhum erro — corrompe dados em silêncio.

O `revisor` precisa de coleções para duas coisas muito concretas: uma **lista** de serviços que você vai
consultar (isso são slices) e um **relatório** que associa cada serviço ao seu estado (isso são maps). Sem
slices nem maps não há programa: você não consegue ter "vários serviços" nem "o estado de cada um" com o que
viu até a Lição 3.

Este é o mapa da lição:

| | |
|---|---|
| **4.1** | Arrays: os que você quase não vai usar |
| **4.2** | Slices: os que você vai usar sempre |
| **4.3** | A armadilha da memória compartilhada |
| **4.4** | Maps |
| **4.5** | A ordem aleatória, e por que ela é deliberada |
| **4.6** | Ler um arquivo de verdade |

---

## Os conceitos

### 4.1 Arrays: os que você quase não vai usar

Um **array** é uma lista de tamanho **fixo**. O tamanho faz parte do tipo.

**Fig. 4.1** | Arrays, para que você os reconheça.

```go
 1  // fig04_01.go
 2  // Muestra arreglos, que en Go casi no se usan directamente.
 3  package main
 4
 5  import "fmt"
 6
 7  func main() {
 8      var puertos [3]int              // tres enteros, todos en 0
 9      fmt.Println("recien creado:", puertos)
10
11      puertos[0] = 443
12      puertos[1] = 8080
13      fmt.Println("con valores:", puertos)
14      fmt.Println("cuantos caben:", len(puertos))
15
16      otros := [3]int{80, 443, 8443}
17      fmt.Println("otro arreglo:", otros)
18
19      copia := otros                  // los arreglos SÍ se copian completos
20      copia[0] = 999
21      fmt.Println("original:", otros)
22      fmt.Println("copia:   ", copia)
23  }
```

```bash
$ go run fig04_01.go
recien creado: [0 0 0]
con valores: [443 8080 0]
cuantos caben: 3
otro arreglo: [80 443 8443]
original: [80 443 8443]
copia:    [999 443 8443]
```

Repare em duas coisas:

**Linha 8.** O array nasce com os valores zero, como tudo em Go. Não há lixo.

**Linhas 19-22.** Ao atribuir um array a outra variável, **ele é copiado por completo**: modificar a cópia
não toca no original. Lembre-se disso, porque com os slices **não** acontece o mesmo, e é aí que está a
armadilha desta lição.

> [!NOTE]
> 🔧 **Observação de engenharia de software 4.1**
> `[3]int` e `[4]int` são **tipos diferentes**. Uma função que recebe `[3]int` não aceita um `[4]int`. Isso
> torna os arrays pouco práticos, e é a razão pela qual em Go quase ninguém os usa diretamente: usam-se
> slices, que crescem. Os arrays estão aí porque os slices são construídos sobre eles.

### 4.2 Slices: os que você vai usar sempre

Um **slice** é uma lista de tamanho **variável**. É o que em outras linguagens você chamaria de lista ou
array dinâmico.

**Fig. 4.2** | Criar e fazer crescer um slice.

```go
 1  // fig04_02.go
 2  // Crea un slice y lo hace crecer con append.
 3  package main
 4
 5  import "fmt"
 6
 7  type Servicio struct {
 8      Nombre string
 9      Puerto int
10  }
11
12  func main() {
13      // forma 1: vacio, y crece
14      var nombres []string
15      fmt.Printf("vacio: %v · largo %d · es nil: %t\n", nombres, len(nombres), nombres == nil)
16
17      nombres = append(nombres, "catalogo")
18      nombres = append(nombres, "pagos", "inventario")
19      fmt.Printf("con datos: %v · largo %d\n", nombres, len(nombres))
20
21      // forma 2: con valores desde el inicio
22      servicios := []Servicio{
23          {Nombre: "catalogo", Puerto: 443},
24          {Nombre: "pagos", Puerto: 443},
25          {Nombre: "local", Puerto: 8080},
26      }
27      fmt.Println("cuantos servicios:", len(servicios))
28
29      // acceder por posicion, empezando en 0
30      fmt.Println("el primero:", servicios[0].Nombre)
31      fmt.Println("el ultimo: ", servicios[len(servicios)-1].Nombre)
32  }
```

```bash
$ go run fig04_02.go
vacio: [] · largo 0 · es nil: true
con datos: [catalogo pagos inventario] · largo 3
cuantos servicios: 3
el primero: catalogo
el ultimo:  local
```

**Linha 15.** Um slice sem inicializar vale **`nil`** e o seu tamanho é `0`. Mas repare em algo importante:

> [!TIP]
> ✅ **Boa prática 4.1**
> **Um slice `nil` pode ser usado com `append`, com `len` e com `for range` sem nenhum problema.** Por isso
> não é preciso inicializá-lo com `[]string{}`: `var nombres []string` e direto para o `append`. É a forma
> idiomática, e mais uma mostra de que o valor zero de Go foi pensado para ser útil.

**Linha 17: `nombres = append(nombres, "catalogo")`.** Repare que é preciso **reatribuir**. `append` não
modifica o slice: **retorna um novo**, e se você não guardar o resultado, ele se perde.

> [!WARNING]
> ⚠️ **Erro comum de programação 4.1**
> Escrever `append(nombres, "x")` sem atribuir o resultado. **O compilador detecta isso** porque o valor
> retornado não é usado (`append(...) evaluated but not used`), então aqui Go te salva. Mas se você o atribuir a
> outra variável por engano, compila e o slice original não muda.

**Linha 31: `servicios[len(servicios)-1]`.** Em Go não há índice negativo como em Python: para o último
elemento, ele é calculado. E se você passar do final, o programa quebra:

**Fig. 4.3** | Sair do intervalo.

```go
 1  // fig04_03.go — este programa COMPILA pero truena
 2  package main
 3
 4  import "fmt"
 5
 6  func main() {
 7      puertos := []int{443, 8080}
 8      fmt.Println(puertos[0])
 9      fmt.Println(puertos[5])      // solo hay 2 elementos
10  }
```

```bash
$ go run fig04_03.go
443
panic: runtime error: index out of range [5] with length 2

goroutine 1 [running]:
main.main()
	/tmp/fig04_03.go:9 +0x1c
exit status 2
```

> [!WARNING]
> ⚠️ **Erro comum de programação 4.2**
> O `index out of range` é o segundo panic mais frequente de Go, depois do ponteiro nulo. **Repare em como a
> mensagem é útil:** ela diz o índice que você pediu (`[5]`), o tamanho real (`length 2`) e a linha.
> Antes de indexar algo que venha de fora —um arquivo, uma requisição, um argumento— verifique `len`.

### 4.3 A armadilha da memória compartilhada

**Esta é a seção mais importante da lição.** É um comportamento que surpreende todo mundo, que não
produz nenhum erro, e que pode corromper dados em silêncio.

#### 4.3.1 O problema

**Fig. 4.4** | Copiar um slice **não** copia os dados.

```go
 1  // fig04_04.go
 2  // Demuestra que asignar un slice NO copia sus datos.
 3  package main
 4
 5  import "fmt"
 6
 7  func main() {
 8      original := []int{443, 8080, 3000}
 9      copia := original              // parece una copia...
10
11      copia[0] = 999                 // modifico la "copia"
12
13      fmt.Println("original:", original)
14      fmt.Println("copia:   ", copia)
15  }
```

```bash
$ go run fig04_04.go
original: [999 8080 3000]
copia:    [999 8080 3000]
```

**Você modificou `copia` e `original` mudou.** Não houve erro, não houve aviso, e os dois slices mostram o
mesmo.

Compare com a Fig. 4.1, onde o mesmo código com um **array** copiava de fato. A diferença é de fundo.

#### 4.3.2 Por que acontece

Um slice **não contém** os dados: é uma janela que aponta para eles. Por dentro, ele guarda três coisas:

| | |
|---|---|
| um **ponteiro** | para onde estão os dados de verdade |
| o **tamanho** (`len`) | quantos elementos ele tem agora |
| a **capacidade** (`cap`) | quantos cabem antes de ter de se mudar |

Quando você escreve `copia := original`, Go copia **essas três coisas** — não os dados. As duas variáveis
ficam apontando para o **mesmo** lugar.

É como te dar o endereço de uma casa em vez de construir uma igual para você: se você pinta a «sua» casa, a
minha muda, porque é a mesma.

#### 4.3.3 A solução

**Fig. 4.5** | Copiar de verdade.

```go
 1  // fig04_05.go
 2  // Las dos formas de copiar un slice de verdad.
 3  package main
 4
 5  import "fmt"
 6
 7  func main() {
 8      original := []int{443, 8080, 3000}
 9
10      // forma 1: make + copy
11      copia1 := make([]int, len(original))
12      copy(copia1, original)
13      copia1[0] = 111
14
15      // forma 2: append sobre un slice nil (mas corta)
16      copia2 := append([]int(nil), original...)
17      copia2[1] = 222
18
19      fmt.Println("original:", original)
20      fmt.Println("copia1:  ", copia1)
21      fmt.Println("copia2:  ", copia2)
22  }
```

```bash
$ go run fig04_05.go
original: [443 8080 3000]
copia1:   [111 8080 3000]
copia2:   [443 222 3000]
```

**Agora sim eles são independentes.**

**Linha 11: `make([]int, len(original))`.** `make` cria um slice com espaço reservado. É a forma de
dizer «quero um slice deste tamanho, com os seus valores zero».

**Linha 16: `original...`.** Esses três pontos significam «espalhe os elementos um por um». Sem eles
você estaria tentando adicionar o slice completo como um único elemento, e não compila.

#### 4.3.4 E a parte realmente traiçoeira: `append`

**Fig. 4.6** | O mesmo código se comporta de forma diferente conforme a capacidade.

```go
 1  // fig04_06.go
 2  // append comparte o no memoria segun la capacidad disponible.
 3  package main
 4
 5  import "fmt"
 6
 7  func main() {
 8      // CASO A: hay capacidad de sobra (cap 5, largo 3)
 9      a := make([]int, 3, 5)
10      a[0], a[1], a[2] = 1, 2, 3
11      fmt.Printf("A: len=%d cap=%d %v\n", len(a), cap(a), a)
12
13      b := append(a, 99)          // cabe: NO se muda, comparte memoria
14      b[0] = 777
15      fmt.Println("   tras modificar b, a vale:", a, "← cambió")
16
17      // CASO B: no hay capacidad (cap 3, largo 3)
18      c := make([]int, 3, 3)
19      c[0], c[1], c[2] = 1, 2, 3
20      fmt.Printf("B: len=%d cap=%d %v\n", len(c), cap(c), c)
21
22      d := append(c, 99)          // NO cabe: se muda a otra memoria
23      d[0] = 777
24      fmt.Println("   tras modificar d, c vale:", c, "← NO cambió")
25  }
```

```bash
$ go run fig04_06.go
A: len=3 cap=5 [1 2 3]
   tras modificar b, a vale: [777 2 3] ← cambió
B: len=3 cap=3 [1 2 3]
   tras modificar d, c vale: [1 2 3] ← NO cambió
```

**Leia essa saída duas vezes.** É **o mesmo código** —`append` e depois modificar— e o resultado é
diferente conforme a capacidade que havia de sobra.

> [!WARNING]
> 🔴 **Erro comum de programação 4.3 — o pior desta lição**
> Guardar um slice que você recebeu como parâmetro, ou ficar com o resultado de `append` supondo que ele é
> independente. **O comportamento depende da capacidade**, que você quase nunca controla e que muda conforme
> o quanto o slice cresceu antes. Ou seja: **o seu programa pode funcionar nos testes e corromper dados
> em produção**, com o mesmo código.
>
> **A regra que te salva:** se você vai **guardar** um slice que recebeu de outra pessoa, copie-o primeiro.
> Se só vai lê-lo, não é preciso.

> [!NOTE]
> 🚀 **Dica de desempenho 4.1**
> Quando você sabe quantos elementos vai adicionar, dê a capacidade desde o início:
> `make([]Estado, 0, len(servicios))`. Assim o `append` não tem de se mudar nem copiar nada enquanto cresce.
> Com listas pequenas é irrelevante; com milhares de elementos, faz diferença. **E meça antes de acreditar.**

### 4.4 Maps

Um **map** guarda pares de chave e valor. É o que em outras linguagens se chama dicionário ou tabela
associativa.

**Fig. 4.7** | Criar, escrever e ler um map.

```go
 1  // fig04_07.go
 2  // Operaciones basicas con un map.
 3  package main
 4
 5  import "fmt"
 6
 7  type Estado struct {
 8      Codigo int
 9      Ms     int
10  }
11
12  func main() {
13      // se crea con make, o vacio con {}
14      estados := make(map[string]Estado)
15
16      estados["catalogo"] = Estado{Codigo: 200, Ms: 142}
17      estados["pagos"] = Estado{Codigo: 200, Ms: 87}
18      estados["inventario"] = Estado{Codigo: 503, Ms: 2310}
19
20      fmt.Println("cuantos:", len(estados))
21
22      // leer una llave que SÍ existe
23      e := estados["catalogo"]
24      fmt.Println("catalogo:", e.Codigo, e.Ms)
25
26      // 🔑 leer una que NO existe: devuelve el valor cero, SIN error
27      fantasma := estados["no-existe"]
28      fmt.Printf("no-existe: %+v  ← el valor cero\n", fantasma)
29
30      // la forma correcta de preguntar: el segundo valor
31      if e, hay := estados["pagos"]; hay {
32          fmt.Println("pagos si esta:", e.Codigo)
33      }
34      if _, hay := estados["no-existe"]; !hay {
35          fmt.Println("no-existe NO esta")
36      }
37
38      // borrar
39      delete(estados, "inventario")
40      fmt.Println("tras borrar:", len(estados))
41  }
```

```bash
$ go run fig04_07.go
cuantos: 3
catalogo: 200 142
no-existe: {Codigo:0 Ms:0}  ← el valor cero
pagos si esta: 200
no-existe NO esta
tras borrar: 2
```

**Linhas 27-28: aqui está o que é preciso entender.** Ler uma chave que não existe **não dá erro**:
retorna o valor zero do tipo. Para um `int` isso é `0`, que poderia ser um valor legítimo.

> [!WARNING]
> ⚠️ **Erro comum de programação 4.4**
> Usar `m[llave]` sem verificar se existe, quando o valor zero é indistinguível de um dado real. Se
> você guarda `map[string]int` com tempos de resposta e lê uma chave inexistente, obtém `0` — que
> parece «respondeu instantaneamente» em vez de «não medimos». **Use sempre a forma de dois valores
> (`v, hay := m[k]`) quando a ausência significar algo diferente de zero.**

> [!TIP]
> ✅ **Boa prática 4.2**
> Chame a segunda variável de `hay`, `existe` ou `ok`. Em código Go você vai ver muito `ok`, e é a
> convenção: `if v, ok := m[k]; ok { … }`.

#### 4.4.1 O map nil: o único que morde

**Fig. 4.8** | Um map `nil` pode ser lido mas não escrito.

```go
 1  // fig04_08.go — este programa truena al escribir
 2  package main
 3
 4  import "fmt"
 5
 6  func main() {
 7      var m map[string]int          // nil: NO inicializado
 8
 9      fmt.Println("leer de un map nil:", m["x"])    // esto funciona
10      fmt.Println("su largo:", len(m))              // esto tambien
11
12      m["x"] = 1                                     // esto truena
13  }
```

```bash
$ go run fig04_08.go
leer de un map nil: 0
su largo: 0
panic: assignment to entry in nil map

goroutine 1 [running]:
main.main()
	/tmp/fig04_08.go:12 +0x3c
exit status 2
```

> [!WARNING]
> 🔴 **Erro comum de programação 4.5**
> **Esta é a grande assimetria de Go e é preciso memorizá-la:** um **slice** `nil` pode ser usado com `append`
> sem problema; um **map** `nil` quebra ao ser escrito. Os maps **precisam ser criados** com `make(map[K]V)` ou
> `map[K]V{}`. E repare no que tem de traiçoeiro: ler do map nil **funciona**, então o programa pode
> avançar um tempo antes de estourar.

### 4.5 A ordem aleatória, e por que ela é deliberada

**Fig. 4.9** | O mesmo programa, duas execuções, duas ordens.

```go
 1  // fig04_09.go
 2  // El recorrido de un map sale en orden distinto cada vez.
 3  package main
 4
 5  import "fmt"
 6
 7  func main() {
 8      puertos := map[string]int{
 9          "catalogo":   443,
10          "pagos":      443,
11          "inventario": 8080,
12          "reportes":   3001,
13      }
14
15      for nombre, puerto := range puertos {
16          fmt.Printf("%s=%d ", nombre, puerto)
17      }
18      fmt.Println()
19  }
```

```bash
$ go run fig04_09.go
reportes=3001 catalogo=443 pagos=443 inventario=8080

$ go run fig04_09.go
pagos=443 inventario=8080 reportes=3001 catalogo=443

$ go run fig04_09.go
inventario=8080 reportes=3001 catalogo=443 pagos=443
```

**Três execuções, três ordens diferentes.** E não é um defeito: **Go faz isso de propósito**, tornando
aleatório o ponto de partida em cada iteração.

> [!NOTE]
> 🔧 **Observação de engenharia de software 4.2**
> Por que incomodar o programador com isso? Porque **a ordem de um map nunca foi garantida**, nem em
> Go nem na maioria das linguagens. Se Go a deixasse «quase sempre igual», haveria programas que
> funcionam durante anos e um dia, quando os dados crescem, mudam de ordem e quebram — e ninguém entenderia
> por quê. **Tornando-a aleatória, Go te obriga a descobrir isso hoje.** É a mesma filosofia de se recusar a compilar com
> uma variável não usada: um erro precoce e incômodo vale mais que um tardio e incompreensível.

#### 4.5.1 Ordenar para que a saída seja estável

Se você precisa de ordem, ela é pedida explicitamente.

**Fig. 4.10** | Percorrer um map em ordem alfabética.

```go
 1  // fig04_10.go
 2  // Ordena las llaves para producir una salida estable.
 3  package main
 4
 5  import (
 6      "fmt"
 7      "sort"
 8  )
 9
10  type Estado struct {
11      Codigo int
12      Ms     int
13  }
14
15  func main() {
16      estados := map[string]Estado{
17          "reportes":   {Codigo: 0, Ms: 0},
18          "catalogo":   {Codigo: 200, Ms: 142},
19          "inventario": {Codigo: 503, Ms: 2310},
20          "pagos":      {Codigo: 200, Ms: 87},
21      }
22
23      // 1. saca las llaves a un slice
24      nombres := make([]string, 0, len(estados))
25      for nombre := range estados {
26          nombres = append(nombres, nombre)
27      }
28
29      // 2. ordenalas
30      sort.Strings(nombres)
31
32      // 3. recorre el slice ordenado, no el map
33      fmt.Println("SERVICIO     CODIGO  TIEMPO")
34      fmt.Println("-----------------------------")
35      for _, nombre := range nombres {
36          e := estados[nombre]
37          fmt.Printf("%-12s %6d  %5dms\n", nombre, e.Codigo, e.Ms)
38      }
39  }
```

```bash
$ go run fig04_10.go
SERVICIO     CODIGO  TIEMPO
-----------------------------
catalogo        200    142ms
inventario      503   2310ms
pagos           200     87ms
reportes          0      0ms
```

**E agora sim: a mesma saída, sempre.** Rode o programa dez vezes e ela não muda.

**Linha 24: `make([]string, 0, len(estados))`.** Tamanho `0`, capacidade `len(estados)`. Você sabe quantas
chaves vai adicionar, então o espaço é reservado de uma vez.

**Linha 25: `for nombre := range estados`.** Sobre um map, `range` dá **a chave** na primeira variável.
Sobre um slice dá **o índice**. É uma diferença que convém ter presente:

| Coleção | primeira variável | segunda variável |
|---|---|---|
| slice | o **índice** (0, 1, 2…) | o elemento |
| map | a **chave** | o valor |

**Linha 35: `for _, nombre := range nombres`.** Aqui `nombres` é um slice, então a primeira variável é
o índice, e eu não preciso dele: por isso `_`.

> [!TIP]
> ✅ **Boa prática 4.3**
> **Se o seu programa imprime resultados, ordene-os.** Uma saída que muda de ordem entre execuções é
> impossível de comparar, impossível de testar automaticamente e confunde quem a lê. Este padrão de
> três passos —tirar as chaves, ordenar, percorrer— é idiomático e você vai usá-lo com frequência.

### 4.6 Ler um arquivo de verdade

Até agora os dados estiveram escritos no programa. Isso não serve: cada mudança exige recompilar.

**Fig. 4.11** | Ler serviços a partir de um arquivo de texto.

Primeiro o arquivo de dados, `servicios.txt`:

```
catalogo https://catalogo.example.com 443
pagos https://pagos.example.com 443
inventario https://inventario.example.com 8080
```

E o programa:

```go
 1  // fig04_11.go
 2  // Lee los servicios desde un archivo de texto.
 3  package main
 4
 5  import (
 6      "fmt"
 7      "os"
 8      "strconv"
 9      "strings"
10  )
11
12  type Servicio struct {
13      Nombre string
14      URL    string
15      Puerto int
16  }
17
18  // Cargar lee el archivo y devuelve los servicios, o un error.
19  func Cargar(ruta string) ([]Servicio, error) {
20      datos, err := os.ReadFile(ruta)
21      if err != nil {
22          return nil, fmt.Errorf("leyendo %s: %w", ruta, err)
23      }
24
25      var servicios []Servicio
26      lineas := strings.Split(strings.TrimSpace(string(datos)), "\n")
27
28      for i, linea := range lineas {
29          linea = strings.TrimSpace(linea)
30          if linea == "" || strings.HasPrefix(linea, "#") {
31              continue                       // salta vacias y comentarios
32          }
33
34          campos := strings.Fields(linea)     // separa por espacios
35          if len(campos) != 3 {
36              return nil, fmt.Errorf("%s linea %d: esperaba 3 campos, hay %d",
37                  ruta, i+1, len(campos))
38          }
39
40          puerto, err := strconv.Atoi(campos[2])
41          if err != nil {
42              return nil, fmt.Errorf("%s linea %d: puerto invalido %q: %w",
43                  ruta, i+1, campos[2], err)
44          }
45
46          servicios = append(servicios, Servicio{
47              Nombre: campos[0],
48              URL:    campos[1],
49              Puerto: puerto,
50          })
51      }
52      return servicios, nil
53  }
54
55  func main() {
56      servicios, err := Cargar("servicios.txt")
57      if err != nil {
58          fmt.Fprintln(os.Stderr, "error:", err)
59          os.Exit(1)
60      }
61
62      fmt.Printf("cargados %d servicios:\n", len(servicios))
63      for _, s := range servicios {
64          fmt.Printf("  %-12s %-38s :%d\n", s.Nombre, s.URL, s.Puerto)
65      }
66  }
```

```bash
$ go run fig04_11.go
cargados 3 servicios:
  catalogo     https://catalogo.example.com           :443
  pagos        https://pagos.example.com              :443
  inventario   https://inventario.example.com         :8080
```

E vamos testar o que acontece quando algo dá errado:

```bash
$ go run fig04_11.go            # con el archivo renombrado
error: leyendo servicios.txt: open servicios.txt: no such file or directory
$ echo $?
1
```

**Linha 20: `os.ReadFile`.** Lê o arquivo completo e retorna os seus bytes. Para arquivos de configuração
é perfeito; para um arquivo de vários gigabytes existem outras formas, porque isto carrega tudo na
memória.

**Linha 26: `string(datos)`.** Converte os bytes em texto. `strings.TrimSpace` remove espaços e quebras de
linha do início e do final — sem isso, a última linha vazia produziria um serviço fantasma.

**Linha 34: `strings.Fields`.** Separa por espaços, **colapsando os repetidos**. É mais robusto que
`strings.Split(linea, " ")`, que com dois espaços seguidos te daria um campo vazio.

**Linha 40: `strconv.Atoi`.** Converte texto em inteiro (*ASCII to integer*). Retorna erro se não for um
número, e **estamos verificando isso**.

**Linhas 58-59: a diferença entre um script e um programa.**

> [!TIP]
> ✅ **Boa prática 4.4**
> Os erros vão para **`os.Stderr`**, não para a saída normal, e o programa termina com **código diferente de
> zero**. Isso permite usá-lo em um pipe: `mi-programa | otro-programa` continua funcionando porque o
> erro não contamina a saída, e `mi-programa && echo ok` não imprime «ok» se falhou. **É o que separa um
> programa de um script.**

> [!TIP]
> 🧪 **Dica de teste e depuração 4.1**
> Repare nas mensagens de erro das linhas 36 e 42: elas dizem **o arquivo, o número da linha e o que
> esperavam**. Compare com um `invalid syntax` puro e simples. Quando alguém usar o seu programa com um arquivo de
> cinquenta linhas, a diferença entre as duas mensagens são vinte minutos da vida dessa pessoa.

---

## O erro que você vai ver

Esta lição deixa duas mensagens de `panic` gravadas na memória, porque são as duas mais frequentes de
todo o curso depois do ponteiro nulo:

**`panic: runtime error: index out of range [N] with length M`** (Fig. 4.3). Você pediu uma posição `N` que
não existe em uma coleção de tamanho `M`. **O que significa:** o índice válido mais alto é `M-1`, não `M`.
**Como se corrige:** verifique `len(coleccion)` antes de indexar, sobretudo se o índice vem de fora
(um arquivo, um argumento, uma requisição) e não foi você mesmo que o escreveu no código.

**`panic: assignment to entry in nil map`** (Fig. 4.8). Você tentou escrever em um map que nunca foi criado —
`var m map[string]Estado` deixa `m` em `nil`, e um map `nil` **pode ser lido mas não escrito**. **O que
significa:** falta o `make(map[K]V)` ou o `map[K]V{}` que reserva a tabela por dentro. **Como se
corrige:** crie o map antes de usá-lo como destino de uma atribuição; se você só vai lê-lo, `nil` não dá
nenhum problema.

E há um terceiro caso que **não** quebra, e por isso é o mais perigoso dos três: modificar um slice
compartilhado (Fig. 4.4 e Fig. 4.6) não produz nenhuma mensagem de erro. O programa continua rodando com dados
incorretos. Esse é exatamente o motivo de a seção 4.3 existir.

## O que se faz errado

- **Supor que `copia := original` copia os dados de um slice.** Copia o ponteiro, o tamanho e a
  capacidade — não os dados. As duas variáveis acabam compartilhando memória (Fig. 4.4).
- **Guardar um slice que você recebeu como parâmetro sem copiá-lo antes.** O comportamento de `append`
  depende da capacidade que o slice traz, que quase nunca é controlada por quem o recebe (Fig. 4.6). Se você vai
  **guardar** um slice alheio, copie-o; se só vai lê-lo, não é preciso.
- **Confiar que um `for range` sobre um map vai sair sempre na mesma ordem.** Nunca foi
  garantido, e Go o torna aleatório de propósito (Fig. 4.9) para que o erro apareça hoje e não no dia em que crescer
  o banco de dados em produção.
- **Ler `m[llave]` sem verificar `ok` quando o valor zero é ambíguo.** Um `map[string]int` que guarda
  tempos de resposta não consegue distinguir «respondeu em 0 ms» de «nunca foi consultado» se você não usa a forma
  de dois valores.
- **Indexar algo que vem de fora —um arquivo, `os.Args`, uma resposta de rede— sem verificar `len`
  primeiro.** É a causa mais comum do `index out of range` desta lição, e se evita com uma linha.
- **Escrever em um map antes de criá-lo com `make` ou com `{}`.** O compilador não detecta isso porque `var m
  map[K]V` é código válido; o `panic` só aparece quando o programa roda e tenta escrever.

---

## Exercícios

### Perguntas

**4.1** Qual é a diferença entre `[3]int` e `[]int`?

**4.2** Por que `[3]int` e `[4]int` são tipos diferentes, e que consequência prática isso tem?

**4.3** O que isto imprime e por quê?

```go
a := []int{1, 2, 3}
b := a
b[0] = 99
fmt.Println(a[0])
```

**4.4** Cite as duas formas de copiar um slice de verdade.

**4.5** Por que o mesmo `append` às vezes compartilha memória com o original e às vezes não?

**4.6** O que `m["no-existe"]` retorna em um `map[string]int`, e por que isso é perigoso?

**4.7** Qual é a assimetria entre um slice `nil` e um map `nil`?

**4.8** Por que a iteração de um map sai desordenada, e como você obtém uma ordem fixa?

**4.9** Em `for a, b := range x`, o que é `a` se `x` é um slice? E se é um map?

**4.10** O que `strings.Fields` faz que `strings.Split(s, " ")` não faz?

**4.11** Por que os erros vão para `os.Stderr` e não para a saída normal?

**4.12** Este programa tem **dois** problemas. Encontre-os.

```go
func main() {
    var m map[string]int
    m["catalogo"] = 443
    lista := []int{1, 2, 3}
    fmt.Println(lista[3])
}
```

### Exercícios de código

**4.13** Escreva `func Contar(estados map[string]Estado) (ok, fallas int)` que conte quantos estados
têm código 200-299 e quantos não. Go retorna vários valores: use isso.

**4.14** Escreva `func Nombres(servicios []Servicio) []string` que retorne só os nomes, **ordenados**.

**4.15 (A armadilha)** Escreva uma função `func Guardar(s []int)` que guarde o slice em uma variável global
e depois o imprima. Chame-a assim:

```go
datos := make([]int, 3, 10)
datos[0], datos[1], datos[2] = 1, 2, 3
Guardar(datos)
datos = append(datos, 4)
datos[0] = 999
// ahora imprime lo que guardó Guardar
```

**Preveja o que vai imprimir antes de rodar.** Depois rode. Se acertou, você entendeu a seção 4.3;
se não, leia-a de novo — é a que sai mais cara.

**4.16** Estenda `Cargar` da Fig. 4.11 para que aceite um quarto campo opcional com o timeout em
milissegundos. Se não vier, use 5000. **Cuide para que um arquivo com três campos continue funcionando.**

**4.17** Escreva `func Agrupar(servicios []Servicio) map[int][]Servicio` que agrupe por porta. Teste que
uma porta com três serviços tenha todos eles.

**4.18 (Encontre o erro)** Diga o que está errado em cada um:

```go
// (a)
var m map[string]Estado
m["catalogo"] = Estado{}

// (b)
append(servicios, nuevo)

// (c)
for i := range servicios {
    fmt.Println(servicios[i+1].Nombre)
}

// (d)
ms := estados["catalogo"].Ms
if ms == 0 {
    fmt.Println("respondio instantaneo")
}
```

**4.19 (Projeto do curso)** Leve o seu programa ao próximo passo:
1. `Cargar(ruta string) ([]Servicio, error)` que leia o arquivo, como a Fig. 4.11.
2. `RevisarTodos` que retorne `map[string]Estado`.
3. `Reporte(estados map[string]Estado) string` que produza a tabela **ordenada alfabeticamente**.
4. Uma contagem final: «3 de 5 responderam».
5. Que o programa aceite o caminho do arquivo como argumento: `os.Args[1]`. 🔴 **Verifique `len(os.Args)`
   antes de lê-lo**, ou você vai ter um `index out of range` quando alguém o rodar sem argumentos.
6. Que saia com código 0 se tudo respondeu e 1 se algo falhou.

### Soluções

**4.1** `[3]int` é um **array**: tamanho fixo, parte do tipo, e é copiado ao ser atribuído. `[]int` é um
**slice**: tamanho variável, e ao ser atribuído **compartilha os dados**.

**4.2** Porque o tamanho faz parte do tipo. A consequência: uma função que recebe `[3]int` não aceita um
`[4]int`, o que torna os arrays pouco práticos. Por isso se usam slices.

**4.3** Imprime **99**. `b := a` não copia os dados: as duas variáveis apontam para a mesma memória.

**4.4** `make` + `copy`, ou `append([]T(nil), original...)`.

**4.5** Porque depende da **capacidade**. Se o slice tem espaço de sobra, `append` escreve ali mesmo e
continua compartilhando memória; se não cabe, ele se muda para memória nova e deixa de compartilhar. Como você quase nunca
controla a capacidade, **o mesmo código pode se comportar de forma diferente**.

**4.6** Retorna `0`, o valor zero, **sem erro**. É perigoso porque `0` pode ser um dado legítimo, então
você não consegue distinguir «vale zero» de «não está». Resolve-se com `v, ok := m[k]`.

**4.7** Um slice `nil` aceita `append` sem problema. Um map `nil` **quebra ao ser escrito**
(`assignment to entry in nil map`), embora possa ser lido. Os maps precisam ser criados com `make`.

**4.8** Porque a sua ordem **nunca foi garantida**, e Go a torna aleatória de propósito para que você não escreva
programas que dependam de uma ordem acidental. Para ordem fixa: tirar as chaves para um slice, ordená-las com
`sort.Strings`, e percorrer o slice.

**4.9** Se `x` é um slice, `a` é o **índice**. Se é um map, `a` é a **chave**.

**4.10** `Fields` separa por espaços **colapsando os repetidos**, então dois espaços seguidos não
produzem um campo vazio. Também trata as tabulações como separador.

**4.11** Para que a saída normal fique limpa e possa ser usada em um pipe, e porque quem chama o
programa espera encontrar os erros ali. Junto com o código de saída diferente de zero, é o que faz
com que um programa possa ser automatizado.

**4.12** Os dois: (1) `m` é um map `nil` e escrever nele provoca `panic: assignment to entry in nil map` —
falta `m := make(map[string]int)`; (2) `lista[3]` sai do intervalo, porque os índices válidos são 0, 1 e
2 — `panic: index out of range [3] with length 3`.

**Exercícios de código (4.13 a 4.19):** suas soluções verificadas —compiladas e executadas— serão adicionadas
na entrega em que for fechado o capítulo de exercícios com solução do curso completo; não são publicadas
sem ter rodado `go build` sobre cada uma.

---

## Como sei que consegui

- Você consegue explicar, sem ver o texto, por que `copia := original` não copia os dados de um slice — e
  desenhá-lo (ponteiro, tamanho, capacidade).
- Você previu corretamente o que o exercício 4.15 imprime **antes** de rodá-lo. Se não acertou, repetiu
  a seção 4.3 até acertar.
- O seu `Cargar` (exercício 4.16) continua funcionando com arquivos de três campos e aceita o quarto opcional.
- O seu `Reporte` produz **a mesma saída, na mesma ordem**, em dez execuções seguidas.
- O seu programa não quebra com `index out of range` se você o roda sem argumentos: você verificou `len(os.Args)`
  antes de ler `os.Args[1]`.
- Você anotou no [`bitacora.md`](bitacora.md) o resultado do exercício 4.15: o que previu e o que aconteceu
  de verdade. É o que separa entender os slices de achar que os entende.

---

## Resumo

- Um **array** (`[3]int`) tem tamanho fixo, o tamanho faz parte do tipo, e **é copiado** ao ser atribuído.
  Quase não se usa diretamente.
- Um **slice** (`[]int`) tem tamanho variável e é o que se usa sempre.
- Um slice guarda **ponteiro, tamanho e capacidade**: não contém os dados, aponta para eles.
- 🔴 **Atribuir um slice NÃO copia os dados**: as duas variáveis compartilham memória. Para copiar de verdade,
  `make` + `copy` ou `append([]T(nil), s...)`.
- 🔴 **`append` compartilha ou não memória conforme a capacidade disponível**, então o mesmo código pode
  se comportar de forma diferente. Se você vai **guardar** um slice alheio, copie-o.
- Um slice `nil` funciona com `append`, `len` e `range`. **Um map `nil` quebra ao ser escrito.**
- Ler uma chave inexistente de um map **retorna o valor zero sem erro**. Use `v, ok := m[k]` quando a
  ausência importar.
- **A iteração de um map é aleatória de propósito.** Para ordem fixa: chaves para um slice, `sort`, e
  percorrer o slice.
- Em `range`, a primeira variável é o **índice** em slices e a **chave** em maps. `_` descarta.
- `os.ReadFile` lê um arquivo completo; `strings.Fields` separa por espaços colapsando repetidos;
  `strconv.Atoi` converte texto em inteiro e **retorna erro**.
- Os erros vão para **`os.Stderr`** e o programa sai com **código diferente de zero**.

---

## Para ler mais

1. **[Go Slices: usage and internals](https://go.dev/blog/slices-intro)** — o blog oficial da equipe de
   Go, explica o ponteiro/tamanho/capacidade com desenhos. Comece por aqui.
2. **[Go maps in action](https://go.dev/blog/maps)** — blog oficial, cobre a ordem aleatória e o map
   `nil` com mais detalhe do que cabe nesta lição.
3. **[Go by Example: Slices](https://gobyexample.com/slices) e [Maps](https://gobyexample.com/maps)** —
   código mínimo, bom como referência rápida.
4. **[Effective Go — Slices](https://go.dev/doc/effective_go#slices)** — a seção específica sobre o
   padrão de dois passos (`make` + `copy`) e por que `append` se comporta como se comporta.

### Termos desta lição

| | |
|---|---|
| **`append`** | adiciona elementos a um slice e **retorna** o resultado |
| **array** | lista de tamanho fixo; o tamanho faz parte do tipo |
| **capacidade (`cap`)** | quantos elementos cabem em um slice antes de ele se mudar |
| **`copy`** | copia elementos entre slices |
| **`delete`** | remove uma chave de um map |
| **`index out of range`** | panic por pedir uma posição que não existe |
| **tamanho (`len`)** | quantos elementos ele tem agora |
| **map** | coleção de pares chave→valor |
| **`make`** | cria slices, maps e canais com espaço reservado |
| **`nil map`** | map não criado; pode ser lido mas **não** escrito |
| **`os.Stderr`** | saída de erro, separada da normal |
| **slice** | lista de tamanho variável; uma janela sobre os dados |
| **`sort.Strings`** | ordena um slice de texto |
| **`strings.Fields`** | separa um texto por espaços, colapsando repetidos |
| **`strconv.Atoi`** | converte texto em inteiro |
| **valor de dois resultados (`v, ok`)** | forma de ler um map distinguindo a ausência |

---

**Anterior:** [Lição 3 — Structs, erros e interfaces](03-errores-interfaces.md) ·
**Próxima:** [Lição 5 — Módulos e testes](05-modulos-y-pruebas.md)
