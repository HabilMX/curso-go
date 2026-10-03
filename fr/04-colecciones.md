# Leçon 4 — Collections : slices et maps

**Durée :** 90 minutes (ou 2 séances de 45).

**À la fin, tu seras capable de :**

- Distinguer un **tableau** d'un **slice** et dire pourquoi tu utiliseras toujours le second.
- Expliquer pourquoi copier un slice **ne copie pas les données**, et quand cela va te mordre.
- Utiliser une **map** et savoir pourquoi son parcours sort en désordre exprès.
- Parcourir des collections avec `for range` et savoir ce que signifie le `_`.
- Trier des résultats pour que ton programme produise toujours la même sortie.
- Lire un fichier et le convertir en une liste de données.

---

## Pourquoi c'est important

**C'est la leçon où presque tout le monde trébuche.** Et pas parce qu'elle est difficile, mais parce qu'il y a
un piège qu'aucun cours n'explique avant qu'il t'ait déjà mordu : en Go, copier un slice ne copie pas ses
données, et ce comportement ne produit aucune erreur — il corrompt les données en silence.

Le `revisor` a besoin de collections pour deux choses très concrètes : une **liste** de services que tu vas
interroger (ce sont des slices) et un **rapport** qui associe chaque service à son état (ce sont des maps).
Sans slices ni maps, il n'y a pas de programme : tu ne peux pas avoir « plusieurs services » ni « l'état de
chacun » avec ce que tu as vu jusqu'à la Leçon 3.

Voici le plan de la leçon :

| | |
|---|---|
| **4.1** | Les tableaux : ceux que tu n'utiliseras presque pas |
| **4.2** | Les slices : ceux que tu utiliseras toujours |
| **4.3** | Le piège de la mémoire partagée |
| **4.4** | Les maps |
| **4.5** | L'ordre aléatoire, et pourquoi il est délibéré |
| **4.6** | Lire un vrai fichier |

---

## Les concepts

### 4.1 Les tableaux : ceux que tu n'utiliseras presque pas

Un **tableau** est une liste de taille **fixe**. La taille fait partie du type.

**Fig. 4.1** | Les tableaux, pour que tu les reconnaisses.

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

Remarque deux choses :

**Ligne 8.** Le tableau naît avec les valeurs zéro, comme tout en Go. Il n'y a pas de déchets.

**Lignes 19-22.** Quand on affecte un tableau à une autre variable, **il est copié en entier** : modifier la
copie ne touche pas l'original. Retiens cela, parce qu'avec les slices **ce n'est pas** pareil, et c'est là
qu'est le piège de cette leçon.

> [!NOTE]
> 🔧 **Observation de génie logiciel 4.1**
> `[3]int` et `[4]int` sont des **types différents**. Une fonction qui reçoit `[3]int` n'accepte pas un
> `[4]int`. Cela rend les tableaux peu pratiques, et c'est la raison pour laquelle en Go presque personne ne
> les utilise directement : on utilise des slices, qui, eux, grandissent. Les tableaux sont là parce que les
> slices sont construits sur eux.

### 4.2 Les slices : ceux que tu utiliseras toujours

Un **slice** est une liste de taille **variable**. C'est ce que, dans d'autres langages, tu appellerais une
liste ou un tableau dynamique.

**Fig. 4.2** | Créer et faire grandir un slice.

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

**Ligne 15.** Un slice non initialisé vaut **`nil`** et sa longueur est `0`. Mais remarque quelque chose
d'important :

> [!TIP]
> ✅ **Bonne pratique 4.1**
> **Un slice `nil` peut s'utiliser avec `append`, avec `len` et avec `for range` sans aucun problème.** C'est
> pourquoi il n'est pas nécessaire de l'initialiser avec `[]string{}` : `var nombres []string` et directement
> `append`. C'est la forme idiomatique, et une preuve de plus que la valeur zéro de Go est pensée pour être
> utile.

**Ligne 17 : `nombres = append(nombres, "catalogo")`.** Remarque qu'il faut **réaffecter**. `append` ne
modifie pas le slice : **il en renvoie un nouveau**, et si tu ne gardes pas le résultat, tu le perds.

> [!WARNING]
> ⚠️ **Erreur courante de programmation 4.1**
> Écrire `append(nombres, "x")` sans affecter le résultat. **Le compilateur le détecte** parce que la valeur
> renvoyée n'est pas utilisée (`append(...) evaluated but not used`), donc ici Go te sauve. Mais si tu
> l'affectes à une autre variable par erreur, cela compile et le slice original ne change pas.

**Ligne 31 : `servicios[len(servicios)-1]`.** En Go, il n'y a pas d'indice négatif comme en Python : pour le
dernier élément, on le calcule. Et si tu dépasses la fin, le programme plante :

**Fig. 4.3** | Sortir de l'intervalle.

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
> ⚠️ **Erreur courante de programmation 4.2**
> Le `index out of range` est le deuxième panic le plus fréquent de Go, après le pointeur nul. **Remarque à
> quel point le message est utile :** il te dit l'indice que tu as demandé (`[5]`), la longueur réelle
> (`length 2`) et la ligne. Avant d'indexer quelque chose qui vient de l'extérieur —un fichier, une requête,
> un argument— vérifie `len`.

### 4.3 Le piège de la mémoire partagée

**C'est la section la plus importante de la leçon.** C'est un comportement qui surprend tout le monde, qui
ne produit aucune erreur, et qui peut corrompre des données en silence.

#### 4.3.1 Le problème

**Fig. 4.4** | Copier un slice **ne** copie **pas** les données.

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

**Tu as modifié `copia` et `original` a changé.** Il n'y a pas eu d'erreur, pas d'avertissement, et les deux
slices affichent la même chose.

Compare-le avec la Fig. 4.1, où le même code avec un **tableau** copiait bel et bien. La différence est de
fond.

#### 4.3.2 Pourquoi cela arrive

Un slice **ne contient pas** les données : c'est une fenêtre qui pointe vers elles. À l'intérieur, il stocke
trois choses :

| | |
|---|---|
| un **pointeur** | vers l'endroit où se trouvent vraiment les données |
| la **longueur** (`len`) | combien d'éléments il a maintenant |
| la **capacité** (`cap`) | combien en tiennent avant de devoir déménager |

Quand tu écris `copia := original`, Go copie **ces trois choses** — pas les données. Les deux variables
pointent vers le **même** endroit.

C'est comme te donner l'adresse d'une maison au lieu de t'en construire une identique : si tu peins « ta »
maison, la mienne change, parce que c'est la même.

#### 4.3.3 La solution

**Fig. 4.5** | Copier pour de vrai.

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

**Maintenant, ils sont bien indépendants.**

**Ligne 11 : `make([]int, len(original))`.** `make` crée un slice avec de l'espace réservé. C'est la façon de
dire « je veux un slice de cette longueur, avec ses valeurs zéro ».

**Ligne 16 : `original...`.** Ces trois points signifient « étale les éléments un par un ». Sans eux, tu
essaierais d'ajouter le slice complet comme un seul élément, et cela ne compile pas.

#### 4.3.4 Et la partie vraiment traîtresse : `append`

**Fig. 4.6** | Le même code se comporte différemment selon la capacité.

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

**Lis cette sortie deux fois.** C'est **le même code** —`append` puis modifier— et le résultat est différent
selon la capacité restante.

> [!WARNING]
> 🔴 **Erreur courante de programmation 4.3 — la pire de cette leçon**
> Conserver un slice que tu as reçu en paramètre, ou garder le résultat d'`append` en supposant qu'il est
> indépendant. **Le comportement dépend de la capacité**, que tu ne contrôles presque jamais et qui change
> selon combien le slice a grandi auparavant. Autrement dit : **ton programme peut fonctionner dans les tests
> et corrompre des données en production**, avec le même code.
>
> **La règle qui te sauve :** si tu vas **conserver** un slice que tu as reçu de quelqu'un d'autre,
> copie-le d'abord. Si tu vas seulement le lire, ce n'est pas nécessaire.

> [!NOTE]
> 🚀 **Astuce de performance 4.1**
> Quand tu sais combien d'éléments tu vas ajouter, donne-lui la capacité dès le départ :
> `make([]Estado, 0, len(servicios))`. Ainsi `append` n'a pas à déménager ni à copier quoi que ce soit
> pendant qu'il grandit. Avec de petites listes, c'est sans importance ; avec des milliers d'éléments, cela se
> remarque. **Et mesure-le avant de le croire.**

### 4.4 Les maps

Une **map** stocke des paires clé-valeur. C'est ce que, dans d'autres langages, on appelle un dictionnaire ou
un tableau associatif.

**Fig. 4.7** | Créer, écrire et lire une map.

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

**Lignes 27-28 : voici ce qu'il faut comprendre.** Lire une clé qui n'existe pas **ne donne pas d'erreur** :
cela renvoie la valeur zéro du type. Pour un `int`, c'est `0`, qui pourrait être une valeur légitime.

> [!WARNING]
> ⚠️ **Erreur courante de programmation 4.4**
> Utiliser `m[llave]` sans vérifier si elle existe, quand la valeur zéro est indiscernable d'une donnée
> réelle. Si tu stockes `map[string]int` avec des temps de réponse et que tu lis une clé inexistante, tu
> obtiens `0` — ce qui ressemble à « il a répondu instantanément » au lieu de « nous ne l'avons pas mesuré ».
> **Utilise toujours la forme à deux valeurs (`v, hay := m[k]`) quand l'absence signifie autre chose que
> zéro.**

> [!TIP]
> ✅ **Bonne pratique 4.2**
> Appelle la deuxième variable `hay`, `existe` ou `ok`. Dans le code Go, tu verras beaucoup `ok`, et c'est
> la convention : `if v, ok := m[k]; ok { … }`.

#### 4.4.1 La map nil : la seule qui mord

**Fig. 4.8** | Une map `nil` peut se lire mais pas s'écrire.

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
> 🔴 **Erreur courante de programmation 4.5**
> **C'est la grande asymétrie de Go et il faut la mémoriser :** un **slice** `nil` peut s'utiliser avec
> `append` sans problème ; une **map** `nil` plante quand on y écrit. Les maps **doivent être créées** avec
> `make(map[K]V)` ou `map[K]V{}`. Et remarque le côté piégeux : lire dans la map nil **fonctionne**, donc le
> programme peut avancer un moment avant d'exploser.

### 4.5 L'ordre aléatoire, et pourquoi il est délibéré

**Fig. 4.9** | Le même programme, deux exécutions, deux ordres.

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

**Trois exécutions, trois ordres différents.** Et ce n'est pas un défaut : **Go le fait exprès**, en rendant
aléatoire le point de départ de chaque parcours.

> [!NOTE]
> 🔧 **Observation de génie logiciel 4.2**
> Pourquoi embêter le programmeur avec cela ? Parce que **l'ordre d'une map n'a jamais été garanti**, ni en
> Go ni dans la plupart des langages. Si Go le laissait « presque toujours pareil », il y aurait des
> programmes qui fonctionnent pendant des années et qui, un jour, quand les données grandissent, changent
> d'ordre et cassent — et personne ne comprendrait pourquoi. **En le rendant aléatoire, il t'oblige à t'en
> rendre compte aujourd'hui.** C'est la même philosophie que refuser de compiler avec une variable inutilisée :
> une erreur précoce et agaçante vaut mieux qu'une erreur tardive et incompréhensible.

#### 4.5.1 Trier pour que la sortie soit stable

Si tu as besoin d'un ordre, il faut le demander explicitement.

**Fig. 4.10** | Parcourir une map dans l'ordre alphabétique.

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

**Et maintenant, oui : la même sortie, toujours.** Lance le programme dix fois et elle ne change pas.

**Ligne 24 : `make([]string, 0, len(estados))`.** Longueur `0`, capacité `len(estados)`. Tu sais combien de
clés tu vas ajouter, donc on réserve l'espace d'un coup.

**Ligne 25 : `for nombre := range estados`.** Sur une map, `range` donne **la clé** dans la première variable.
Sur un slice, il donne **l'indice**. C'est une différence qu'il vaut mieux garder en tête :

| Collection | première variable | deuxième variable |
|---|---|---|
| slice | l'**indice** (0, 1, 2…) | l'élément |
| map | la **clé** | la valeur |

**Ligne 35 : `for _, nombre := range nombres`.** Ici `nombres` est un slice, donc la première variable est
l'indice, et je n'en ai pas besoin : d'où le `_`.

> [!TIP]
> ✅ **Bonne pratique 4.3**
> **Si ton programme affiche des résultats, trie-les.** Une sortie qui change d'ordre entre deux exécutions
> est impossible à comparer, impossible à tester automatiquement et embrouille celui qui la lit. Ce modèle en
> trois étapes —extraire les clés, trier, parcourir— est idiomatique et tu vas l'utiliser souvent.

### 4.6 Lire un vrai fichier

Jusqu'ici, les données étaient écrites dans le programme. Cela ne sert à rien : chaque changement exige de
recompiler.

**Fig. 4.11** | Lire des services depuis un fichier texte.

D'abord le fichier de données, `servicios.txt` :

```
catalogo https://catalogo.example.com 443
pagos https://pagos.example.com 443
inventario https://inventario.example.com 8080
```

Et le programme :

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

Et voyons ce qui se passe quand quelque chose tourne mal :

```bash
$ go run fig04_11.go            # con el archivo renombrado
error: leyendo servicios.txt: open servicios.txt: no such file or directory
$ echo $?
1
```

**Ligne 20 : `os.ReadFile`.** Lit le fichier en entier et renvoie ses octets. Pour des fichiers de
configuration, c'est parfait ; pour un fichier de plusieurs gigaoctets, il existe d'autres façons, parce que
ceci charge tout en mémoire.

**Ligne 26 : `string(datos)`.** Convertit les octets en texte. `strings.TrimSpace` retire les espaces et les
sauts de ligne du début et de la fin — sans cela, la dernière ligne vide produirait un service fantôme.

**Ligne 34 : `strings.Fields`.** Sépare par les espaces, **en fusionnant ceux qui se répètent**. C'est plus
robuste que `strings.Split(linea, " ")`, qui, avec deux espaces consécutifs, te donnerait un champ vide.

**Ligne 40 : `strconv.Atoi`.** Convertit du texte en entier (*ASCII to integer*). Renvoie une erreur si ce
n'est pas un nombre, et **nous la vérifions**.

**Lignes 58-59 : la différence entre un script et un programme.**

> [!TIP]
> ✅ **Bonne pratique 4.4**
> Les erreurs vont vers **`os.Stderr`**, pas vers la sortie normale, et le programme se termine avec un
> **code différent de zéro**. Cela permet de l'utiliser dans un tube (pipe) : `mi-programa | otro-programa`
> continue de fonctionner parce que l'erreur ne contamine pas la sortie, et `mi-programa && echo ok`
> n'affiche pas « ok » s'il a échoué. **C'est ce qui sépare un programme d'un script.**

> [!TIP]
> 🧪 **Astuce de test et de débogage 4.1**
> Regarde les messages d'erreur des lignes 36 et 42 : ils disent **le fichier, le numéro de ligne et ce
> qu'ils attendaient**. Compare-les avec un simple `invalid syntax`. Quand quelqu'un utilisera ton programme
> avec un fichier de cinquante lignes, la différence entre les deux messages, ce sont vingt minutes de sa
> vie.

---

## L'erreur que tu vas voir

Cette leçon grave deux messages de `panic` dans la mémoire, parce que ce sont les deux plus fréquents de
tout le cours après le pointeur nul :

**`panic: runtime error: index out of range [N] with length M`** (Fig. 4.3). Tu as demandé une position `N`
qui n'existe pas dans une collection de longueur `M`. **Ce que cela signifie :** l'indice valide le plus élevé
est `M-1`, pas `M`. **Comment le corriger :** vérifie `len(coleccion)` avant d'indexer, surtout si l'indice
vient de l'extérieur (un fichier, un argument, une requête) et que tu ne l'as pas écrit toi-même dans le code.

**`panic: assignment to entry in nil map`** (Fig. 4.8). Tu as essayé d'écrire dans une map qui n'a jamais été
créée — `var m map[string]Estado` laisse `m` à `nil`, et une map `nil` **peut se lire mais pas s'écrire**.
**Ce que cela signifie :** il manque le `make(map[K]V)` ou le `map[K]V{}` qui réserve la table en interne.
**Comment le corriger :** crée la map avant de l'utiliser comme destination d'une affectation ; si tu vas
seulement la lire, `nil` ne pose aucun problème.

Et il y a un troisième cas qui **ne** plante **pas**, et c'est pour cela qu'il est le plus dangereux des
trois : modifier un slice partagé (Fig. 4.4 et Fig. 4.6) ne produit aucun message d'erreur. Le programme
continue de tourner avec des données incorrectes. C'est exactement la raison d'être de la section 4.3.

## Ce qu'on fait mal

- **Supposer que `copia := original` copie les données d'un slice.** Cela copie le pointeur, la longueur et
  la capacité — pas les données. Les deux variables finissent par partager la mémoire (Fig. 4.4).
- **Conserver un slice reçu en paramètre sans le copier d'abord.** Le comportement d'`append` dépend de la
  capacité qu'apporte le slice, que celui qui le reçoit ne contrôle presque jamais (Fig. 4.6). Si tu vas
  **conserver** un slice qui n'est pas à toi, copie-le ; si tu vas seulement le lire, ce n'est pas
  nécessaire.
- **Compter sur le fait qu'un `for range` sur une map sortira toujours dans le même ordre.** Cela n'a jamais
  été garanti, et Go le rend aléatoire exprès (Fig. 4.9) pour que l'erreur apparaisse aujourd'hui et non le
  jour où la base de données grandira en production.
- **Lire `m[llave]` sans vérifier `ok` quand la valeur zéro est ambiguë.** Une `map[string]int` qui stocke
  des temps de réponse ne peut pas distinguer « a répondu en 0 ms » de « n'a jamais été interrogé » si tu
  n'utilises pas la forme à deux valeurs.
- **Indexer quelque chose qui vient de l'extérieur —un fichier, `os.Args`, une réponse réseau— sans vérifier
  `len` d'abord.** C'est la cause la plus courante du `index out of range` de cette leçon, et cela s'évite
  avec une ligne.
- **Écrire dans une map avant de la créer avec `make` ou avec `{}`.** Le compilateur ne le détecte pas parce
  que `var m map[K]V` est du code valide ; le `panic` n'apparaît que lorsque le programme tourne et essaie
  d'écrire.

---

## Exercices

### Questions

**4.1** Quelle est la différence entre `[3]int` et `[]int` ?

**4.2** Pourquoi `[3]int` et `[4]int` sont-ils des types différents, et quelle conséquence pratique cela
a-t-il ?

**4.3** Qu'affiche ceci et pourquoi ?

```go
a := []int{1, 2, 3}
b := a
b[0] = 99
fmt.Println(a[0])
```

**4.4** Nomme les deux façons de copier un slice pour de vrai.

**4.5** Pourquoi le même `append` partage-t-il parfois la mémoire avec l'original et parfois non ?

**4.6** Que renvoie `m["no-existe"]` dans une `map[string]int`, et pourquoi est-ce dangereux ?

**4.7** Quelle est l'asymétrie entre un slice `nil` et une map `nil` ?

**4.8** Pourquoi le parcours d'une map sort-il en désordre, et comment obtiens-tu un ordre fixe ?

**4.9** Dans `for a, b := range x`, qu'est-ce que `a` si `x` est un slice ? Et si c'est une map ?

**4.10** Que fait `strings.Fields` que `strings.Split(s, " ")` ne fait pas ?

**4.11** Pourquoi les erreurs vont-elles vers `os.Stderr` et pas vers la sortie normale ?

**4.12** Ce programme a **deux** problèmes. Trouve-les.

```go
func main() {
    var m map[string]int
    m["catalogo"] = 443
    lista := []int{1, 2, 3}
    fmt.Println(lista[3])
}
```

### Exercices de code

**4.13** Écris `func Contar(estados map[string]Estado) (ok, fallas int)` qui compte combien d'états ont un
code 200-299 et combien non. Go renvoie plusieurs valeurs : sers-t'en.

**4.14** Écris `func Nombres(servicios []Servicio) []string` qui renvoie uniquement les noms, **triés**.

**4.15 (Le piège)** Écris une fonction `func Guardar(s []int)` qui stocke le slice dans une variable globale
puis l'affiche. Appelle-la ainsi :

```go
datos := make([]int, 3, 10)
datos[0], datos[1], datos[2] = 1, 2, 3
Guardar(datos)
datos = append(datos, 4)
datos[0] = 999
// ahora imprime lo que guardó Guardar
```

**Prédis ce qu'elle va afficher avant de l'exécuter.** Ensuite, exécute-la. Si tu as vu juste, tu as compris
la section 4.3 ; sinon, relis-la — c'est celle qui coûte le plus cher.

**4.16** Étends `Cargar` de la Fig. 4.11 pour qu'il accepte un quatrième champ facultatif avec le délai
d'expiration en millisecondes. S'il n'est pas fourni, utilise 5000. **Veille à ce qu'un fichier à trois
champs continue de fonctionner.**

**4.17** Écris `func Agrupar(servicios []Servicio) map[int][]Servicio` qui regroupe par port. Vérifie qu'un
port avec trois services les contient tous.

**4.18 (Trouve l'erreur)** Dis ce qui ne va pas dans chacun :

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

**4.19 (Projet du cours)** Fais passer ton programme à l'étape suivante :
1. `Cargar(ruta string) ([]Servicio, error)` qui lit le fichier, comme la Fig. 4.11.
2. `RevisarTodos` qui renvoie `map[string]Estado`.
3. `Reporte(estados map[string]Estado) string` qui produit le tableau **trié par ordre alphabétique**.
4. Un décompte final : « 3 de 5 respondieron ».
5. Que le programme accepte le chemin du fichier comme argument : `os.Args[1]`. 🔴 **Vérifie
   `len(os.Args)` avant de le lire**, sinon tu auras un `index out of range` quand quelqu'un le lancera sans
   arguments.
6. Fais en sorte qu'il se termine avec le code 0 si tout a répondu et 1 si quelque chose a échoué.

### Corrigés

**4.1** `[3]int` est un **tableau** : taille fixe, qui fait partie du type, et il est copié lors de
l'affectation. `[]int` est un **slice** : taille variable, et lors de l'affectation il **partage les
données**.

**4.2** Parce que la taille fait partie du type. La conséquence : une fonction qui reçoit `[3]int` n'accepte
pas un `[4]int`, ce qui rend les tableaux peu pratiques. C'est pourquoi on utilise des slices.

**4.3** Cela affiche **99**. `b := a` ne copie pas les données : les deux variables pointent vers la même
mémoire.

**4.4** `make` + `copy`, ou `append([]T(nil), original...)`.

**4.5** Parce que cela dépend de la **capacité**. Si le slice a de l'espace en trop, `append` écrit à cet
endroit même et continue de partager la mémoire ; si cela ne tient pas, il déménage vers une nouvelle mémoire
et cesse de partager. Comme tu ne contrôles presque jamais la capacité, **le même code peut se comporter
différemment**.

**4.6** Cela renvoie `0`, la valeur zéro, **sans erreur**. C'est dangereux parce que `0` peut être une donnée
légitime, donc tu ne peux pas distinguer « vaut zéro » de « n'est pas là ». Cela se résout avec
`v, ok := m[k]`.

**4.7** Un slice `nil` accepte `append` sans problème. Une map `nil` **plante quand on y écrit**
(`assignment to entry in nil map`), bien qu'on puisse la lire. Les maps doivent être créées avec `make`.

**4.8** Parce que leur ordre **n'a jamais été garanti**, et Go le rend aléatoire exprès pour que tu n'écrives
pas de programmes qui dépendent d'un ordre accidentel. Pour un ordre fixe : extraire les clés dans un slice,
les trier avec `sort.Strings`, et parcourir le slice.

**4.9** Si `x` est un slice, `a` est l'**indice**. Si c'est une map, `a` est la **clé**.

**4.10** `Fields` sépare par les espaces **en fusionnant ceux qui se répètent**, donc deux espaces consécutifs
ne produisent pas de champ vide. Il traite aussi les tabulations comme séparateur.

**4.11** Pour que la sortie normale reste propre et puisse être utilisée dans un tube, et parce que celui qui
appelle le programme s'attend à y trouver les erreurs. Avec le code de sortie différent de zéro, c'est ce
qui fait qu'un programme peut être automatisé.

**4.12** Les deux : (1) `m` est une map `nil` et y écrire provoque `panic: assignment to entry in nil map` —
il manque `m := make(map[string]int)` ; (2) `lista[3]` sort de l'intervalle, parce que les indices valides
sont 0, 1 et 2 — `panic: index out of range [3] with length 3`.

**Exercices de code (4.13 à 4.19) :** leurs corrigés vérifiés —compilés et exécutés— seront ajoutés lors de
la livraison qui clôturera le chapitre des exercices avec corrigé du cours complet ; ils ne sont pas publiés
sans avoir lancé `go build` sur chacun.

---

## Comment savoir que j'y suis arrivé

- Tu peux expliquer, sans regarder le texte, pourquoi `copia := original` ne copie pas les données d'un
  slice — et le dessiner (pointeur, longueur, capacité).
- Tu as prédit correctement ce qu'affiche l'exercice 4.15 **avant** de l'exécuter. Si tu ne l'as pas fait, tu
  as repris la section 4.3 jusqu'à y arriver.
- Ton `Cargar` (exercice 4.16) fonctionne toujours avec des fichiers à trois champs et accepte le quatrième
  facultatif.
- Ton `Reporte` produit **la même sortie, dans le même ordre**, sur dix exécutions consécutives.
- Ton programme ne plante pas avec `index out of range` si tu le lances sans arguments : tu as vérifié
  `len(os.Args)` avant de lire `os.Args[1]`.
- Tu as noté dans [`bitacora.md`](bitacora.md) le résultat de l'exercice 4.15 : ce que tu avais prédit et ce
  qui s'est réellement passé. C'est ce qui sépare comprendre les slices de croire qu'on les comprend.

---

## Résumé

- Un **tableau** (`[3]int`) a une taille fixe, la taille fait partie du type, et il **est copié** lors de
  l'affectation. On ne l'utilise presque pas directement.
- Un **slice** (`[]int`) a une taille variable et c'est ce qu'on utilise toujours.
- Un slice stocke **pointeur, longueur et capacité** : il ne contient pas les données, il pointe vers elles.
- 🔴 **Affecter un slice NE copie PAS les données** : les deux variables partagent la mémoire. Pour copier
  pour de vrai, `make` + `copy` ou `append([]T(nil), s...)`.
- 🔴 **`append` partage ou non la mémoire selon la capacité disponible**, donc le même code peut se comporter
  différemment. Si tu vas **conserver** un slice qui n'est pas à toi, copie-le.
- Un slice `nil` fonctionne avec `append`, `len` et `range`. **Une map `nil` plante quand on y écrit.**
- Lire une clé inexistante d'une map **renvoie la valeur zéro sans erreur**. Utilise `v, ok := m[k]` quand
  l'absence compte.
- **Le parcours d'une map est aléatoire exprès.** Pour un ordre fixe : les clés dans un slice, `sort`, et
  parcourir le slice.
- Dans `range`, la première variable est l'**indice** pour les slices et la **clé** pour les maps. `_`
  écarte.
- `os.ReadFile` lit un fichier en entier ; `strings.Fields` sépare par les espaces en fusionnant les
  répétitions ; `strconv.Atoi` convertit du texte en entier et **renvoie une erreur**.
- Les erreurs vont vers **`os.Stderr`** et le programme se termine avec un **code différent de zéro**.

---

## Pour aller plus loin

1. **[Go Slices: usage and internals](https://go.dev/blog/slices-intro)** — le blog officiel de l'équipe de
   Go, explique le pointeur/longueur/capacité avec des dessins. Commence ici.
2. **[Go maps in action](https://go.dev/blog/maps)** — blog officiel, couvre l'ordre aléatoire et la map
   `nil` avec plus de détails que n'en contient cette leçon.
3. **[Go by Example: Slices](https://gobyexample.com/slices) et [Maps](https://gobyexample.com/maps)** —
   code minimal, bon comme référence rapide.
4. **[Effective Go — Slices](https://go.dev/doc/effective_go#slices)** — la section consacrée au modèle en
   deux étapes (`make` + `copy`) et à la raison pour laquelle `append` se comporte comme il se comporte.

### Termes de cette leçon

| | |
|---|---|
| **`append`** | ajoute des éléments à un slice et **renvoie** le résultat |
| **tableau** | liste de taille fixe ; la taille fait partie du type |
| **capacité (`cap`)** | combien d'éléments tiennent dans un slice avant qu'il déménage |
| **`copy`** | copie des éléments entre slices |
| **`delete`** | retire une clé d'une map |
| **`index out of range`** | panic pour avoir demandé une position qui n'existe pas |
| **longueur (`len`)** | combien d'éléments il a maintenant |
| **map** | collection de paires clé→valeur |
| **`make`** | crée des slices, des maps et des canaux avec de l'espace réservé |
| **`nil map`** | map non créée ; elle peut se lire mais **pas** s'écrire |
| **`os.Stderr`** | sortie des erreurs, séparée de la sortie normale |
| **slice** | liste de taille variable ; une fenêtre sur les données |
| **`sort.Strings`** | trie un slice de textes |
| **`strings.Fields`** | sépare un texte par les espaces, en fusionnant les répétitions |
| **`strconv.Atoi`** | convertit du texte en entier |
| **valeur à deux résultats (`v, ok`)** | façon de lire une map en distinguant l'absence |

---

**Précédent :** [Leçon 3 — Structs, erreurs et interfaces](03-errores-interfaces.md) ·
**Suivant :** [Leçon 5 — Modules et tests](05-modulos-y-pruebas.md)
