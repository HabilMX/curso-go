# Leçon 2 — Variables, fonctions et types

**Durée :** 90 minutes, ou deux séances de 45.

**À la fin, tu seras capable de :**

- Expliquer ce que fait le compilateur de Go avec ton fichier texte.
- Déclarer des variables et savoir quand utiliser `:=` et quand utiliser `=`.
- Nommer les quatre types de base et dire pourquoi Go ne te laisse pas les mélanger.
- Écrire des fonctions qui reçoivent et renvoient des valeurs.
- Comprendre pourquoi une fonction Go peut renvoyer **deux** choses à la fois.
- Lire les messages d'erreur du compilateur et savoir ce qu'il te demande.
- Écrire, compiler et exécuter un programme complet par toi-même.

---

## Pourquoi c'est important

Le `revisor` —le programme que tu vas construire tout au long du cours— a besoin, dès sa première ligne, de trois
choses : des **données** (le nom d'un service, son port, le temps qu'il a mis à répondre), de la **logique** qui
décide quelque chose avec ces données (est-il rapide ou lent ?) et un moyen de **signaler que quelque chose s'est
mal passé** (le service n'a pas répondu ?). C'est, dans cet ordre, exactement ce qu'enseigne cette leçon : des
variables pour stocker des données, des fonctions pour la logique, et le motif des deux valeurs de retour pour les
erreurs.

Ce n'est pas un hasard si Go est **compilé** et à **typage statique** : ces deux décisions existent pour que
l'ordinateur attrape tes erreurs **avant** que le programme ne s'exécute, au lieu qu'un utilisateur les découvre
en production. Un langage interprété te laisse écrire `puerto = "muchos"` et ne plante que lorsque cette ligne
s'exécute —ce qui peut être dix minutes plus tard, ou jamais, si ce chemin du code n'est presque pas utilisé—. Go
refuse de produire l'exécutable. Cette différence est la raison de fond de presque tout ce que tu vas voir dans
cette leçon : pourquoi le compilateur est strict, pourquoi la « valeur zéro » existe, et pourquoi une fonction
qui peut échouer **doit le dire dans sa signature**, et non comme une surprise.

À la fin de la leçon, tu écriras le premier pas réel du `revisor` (exercice de projet) : une fonction qui reçoit le
nom d'un service, son port et son temps de réponse, et renvoie une ligne de rapport classée. Toujours sans
structs, sans concurrence et sans HTTP —cela viendra plus tard— mais c'est déjà du code qui ressemble au
programme final.

---

## Les concepts

### 2.1 Ce que fait le compilateur

Tu écris un fichier texte. Ce fichier, à lui seul, ne fait rien : c'est du texte. Pour que l'ordinateur
l'exécute, quelque chose doit le traduire en instructions que le processeur comprend.

Il y a deux façons de faire cette traduction :

| | Comment ça fonctionne | Langages |
|---|---|---|
| **Interprété** | Un programme lit ton texte ligne par ligne et fait ce qu'il dit, **à chaque fois** que tu l'exécutes | Python, JavaScript, PHP |
| **Compilé** | Un programme traduit **tout** ton texte **une seule fois** et range le résultat dans un fichier exécutable | **Go**, C, C++, Rust |

**Go est compilé**, et cela a trois conséquences que tu remarqueras dès aujourd'hui :

1. **Il y a une étape avant l'exécution : compiler.** Si tu as mal écrit le nom d'une variable, tu le découvres à
   ce moment-là — avant que le programme ne s'exécute. Dans un langage interprété, tu le découvrirais quand
   l'exécution arriverait à cette ligne, ce qui peut être dix minutes plus tard ou jamais.
2. **Le résultat est un fichier qui fonctionne seul.** Il n'a pas besoin que l'autre ordinateur ait Go installé.
3. **Il est rapide à l'exécution**, parce que la traduction est déjà faite.

> [!NOTE]
> 🔧 **Observation de génie logiciel 2.1**
> Le compilateur est ta première ligne de défense, pas un obstacle. Chaque erreur qu'il attrape est une erreur
> que tu n'auras pas à chercher de nuit avec le programme déjà livré. Quand Go t'arrête, lis le message
> calmement : il t'épargne du travail.

### 2.2 Ton premier programme, ligne par ligne

Nous allons écrire, compiler et exécuter un programme complet. **Tous les programmes de ce cours sont complets et
exécutables** : pas de fragments qu'on ne peut pas lancer.

Prépare le dossier :

```bash
mkdir -p ~/w/curso-go/cap01 && cd ~/w/curso-go/cap01
go mod init cap01
```

**Fig. 2.1** | Un programme qui affiche un message.

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

Compile-le et exécute-le :

```bash
$ go run fig02_01.go
Bienvenido a Go!
```

Maintenant **chaque ligne**, parce que toutes ont une raison :

**Lignes 1-2.** Commentaires. Ils commencent par `//` et le compilateur les ignore totalement. Ils sont là pour
celui qui lira le code — toi compris dans trois semaines.

**Ligne 3 : `package main`.** Tout fichier Go appartient à un *paquet*, qui est un regroupement de code. **Le
paquet `main` est spécial : c'est le seul qui produit un programme exécutable.** Si tu écrivais
`package utilidades`, tu aurais une bibliothèque que d'autres programmes peuvent utiliser, mais qui ne peut pas
s'exécuter seule.

**Ligne 5 : `import "fmt"`.** Dit au compilateur que tu vas utiliser du code du paquet `fmt` (de *format*), qui
est fourni avec Go et contient les fonctions pour afficher et mettre en forme du texte. **Go ne charge rien par
défaut** : tout ce que tu utilises, tu le demandes.

**Ligne 7 : `func main() {`.** Déclare la fonction `main`, qui est **l'endroit où commence l'exécution du
programme**. Quand tu exécutes un programme Go, le système cherche exactement cette fonction. Si tu l'appelais
`principal` ou `inicio`, le programme ne démarrerait pas.

**Ligne 8 : `fmt.Println("Bienvenido a Go!")`.** Appelle la fonction `Println` qui se trouve dans le paquet
`fmt`. Le point se lit comme « de » : *« la fonction `Println` **de** `fmt` »*. `Println` affiche ce que tu lui
donnes et passe à la ligne suivante (*print line*).

**Ligne 9 : `}`.** Ferme le corps de la fonction. En Go, les accolades délimitent des blocs, comme en C, C++, Java
ou JavaScript.

> [!TIP]
> ✅ **Bonne pratique 2.1**
> Mets toujours un commentaire au début du fichier avec son nom et ce qu'il fait. Cela prend cinq secondes et fait
> gagner des minutes à celui qui l'ouvrira ensuite.

#### 2.2.1 La différence entre `go run` et `go build`

**Fig. 2.2** | Les deux façons d'exécuter ton programme.

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

Regarde les tailles : ton texte fait **104 octets** ; le programme compilé, **1,8 mégaoctet**. La différence est
que l'exécutable embarque tout ce dont il a besoin pour fonctionner.

> [!NOTE]
> 🚀 **Astuce de portabilité 2.1**
> Ce fichier `fig02_01` peut être copié sur **n'importe quel** ordinateur Linux de la même architecture et il
> fonctionne, même si cette machine n'a pas Go installé. C'est la raison pour laquelle tant d'outils serveur sont
> écrits en Go.

> [!TIP]
> 🧪 **Astuce de test et de débogage 2.1**
> Utilise `go run` pendant que tu programmes : c'est plus rapide et ça ne remplit pas ton dossier d'exécutables.
> Utilise `go build` quand le programme est au point et que tu veux le livrer.

### 2.3 Variables

Une **variable** est un espace mémoire doté d'un nom, où tu ranges une donnée que ton programme va utiliser.

**Fig. 2.3** | Déclarer des variables et les afficher.

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

**L'opérateur `:=`** se lit *« déclare une nouvelle variable et range-y ceci »*. Go regarde la valeur de droite et
**déduit le type tout seul** : `"catalogo"` est entre guillemets, donc c'est du texte ; `443` n'en a pas et n'a pas
de point, donc c'est un entier.

#### 2.3.1 `:=` contre `=`

Une fois que la variable existe, pour changer sa valeur tu utilises `=` **sans** les deux-points :

**Fig. 2.4** | Créer et modifier une variable.

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
> ✅ **Bonne pratique 2.2**
> Si tu as vraiment besoin de déclarer quelque chose que tu n'utiliseras pas tout de suite, utilise
> l'identifiant vide `_`. En Go, `_` signifie *« ceci, je le jette exprès »*, et le compilateur l'accepte.

#### 2.3.2 Go refuse de compiler si tu n'utilises pas une variable

**Fig. 2.5** | Un programme qui **ne compile pas**, à dessein.

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

**Go ne le permet pas.** Ce n'est pas un avertissement que tu peux ignorer : c'est une erreur qui arrête la
compilation.

**Pourquoi si strict ?** Parce qu'une variable inutilisée signifie presque toujours l'une de deux choses : tu t'es
trompé en écrivant le nom dans une autre ligne, ou il est resté des déchets d'un code que tu as supprimé à moitié.
Les deux sont des problèmes. Go préfère t'embêter aujourd'hui plutôt que de laisser du code mort s'accumuler pour
toujours.

### 2.4 Les types de base

Un **type** est la réponse à la question *« quelle sorte de donnée est-ce ? »*. Les quatre que tu vas utiliser en
permanence :

| Type | Ce qu'il stocke | Exemples | Valeur zéro |
|---|---|---|---|
| `string` | du texte | `"hola"`, `"https://catalogo.example.com"` | `""` (vide) |
| `int` | des nombres entiers | `42`, `-7`, `0` | `0` |
| `float64` | des nombres à virgule | `3.14`, `0.142` | `0` |
| `bool` | vrai ou faux | `true`, `false` | `false` |

#### 2.4.1 Go ne mélange pas les types

**Fig. 2.6** | Un autre programme qui **ne compile pas**, à dessein.

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

En Python, cela fonctionnerait sans broncher. **En Go, ça ne compile pas**, et c'est un énorme avantage :

> [!NOTE]
> 🔧 **Observation de génie logiciel 2.2**
> Une bonne partie des erreurs qui arrivent en production sont de cette classe : quelqu'un a cru qu'une variable
> contenait un nombre alors qu'elle contenait du texte. Dans un langage interprété, cela plante **quand un client
> utilise le programme**. En Go, cela n'arrive pas jusqu'à la compilation : l'erreur apparaît sur ta machine, pas
> sur celle du client.

#### 2.4.2 La valeur zéro : en Go, il n'y a jamais de déchets

**Fig. 2.7** | Variables déclarées sans valeur.

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

La forme `var nombre tipo` déclare une variable **sans lui donner de valeur**. Dans d'autres langages, cela
laisserait des « déchets » —ce qui se trouvait dans cette mémoire— ou un état spécial « indéfini ». **En Go, elle
reçoit toujours une valeur valide**, appelée **valeur zéro**.

**Et voici `Printf`**, qui est différent de `Println` : il reçoit un modèle avec des trous et les remplit. Chaque
trou commence par `%` :

| | À quoi ça sert |
|---|---|
| `%d` | un nombre entier (*decimal*) |
| `%s` | du texte (*string*) |
| `%q` | du texte **avec guillemets** (*quoted*) — utile pour voir si quelque chose est vide |
| `%t` | un `bool` (*true/false*) |
| `%v` | n'importe quoi, dans son format naturel (*value*) |
| `\n` | saut de ligne (`Printf` ne saute **pas** de ligne tout seul, `Println` si) |

> [!TIP]
> ✅ **Bonne pratique 2.3**
> Utilise `%q` quand tu affiches du texte que tu es en train de déboguer. `%s` avec un texte vide n'affiche rien
> et on dirait que la ligne a échoué ; `%q` affiche `""` et on voit clairement que le texte est vide.

### 2.5 Fonctions

Une **fonction** est un bloc de code nommé, auquel tu peux donner des données et qui peut te renvoyer un
résultat. Tu en as déjà utilisé deux : `fmt.Println` et `fmt.Printf`.

**Fig. 2.8** | Définir et appeler des fonctions.

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

**Anatomie de la ligne 8**, où tout se trouve :

```
func  sumar  (a int, b int)  int  {
 │      │          │          │
 │      │          │          └── lo que DEVUELVE
 │      │          └── lo que RECIBE (los parámetros), con su tipo
 │      └── el nombre
 └── palabra clave: «voy a definir una función»
```

**Ligne 14 : `fmt.Sprintf`.** C'est comme `Printf` mais au lieu d'afficher, il **renvoie** le texte construit. Le
`S` vient de *string*. On s'en sert énormément.

**Ligne 19 : `puerto == 443`.** Le double égal **compare** et donne `true` ou `false`. Un seul égal `=`
**affecte**. Les confondre est un classique (voir « Ce qu'on fait mal », plus bas).

> [!TIP]
> ✅ **Bonne pratique 2.4**
> Quand deux paramètres sont du même type, tu peux abréger : `func sumar(a, b int) int`. C'est idiomatique et plus
> court. Ici nous l'écrivons en entier pour que la structure se voie.

### 2.6 Renvoyer deux valeurs : la signature de Go

**C'est le trait le plus distinctif de Go**, et tu l'écriras des milliers de fois dans ta carrière. Fais bien
attention.

Pense à une fonction qui divise deux nombres. Que fait-elle si le diviseur est zéro ? Elle ne peut pas renvoyer
un nombre, parce que le résultat n'existe pas.

D'autres langages utilisent des **exceptions** : la fonction « lance » une erreur et quelqu'un plus haut
l'« attrape » avec `try / catch`. Le problème est qu'**on ne voit pas clairement qui l'attrape ni où**, et il est
très facile que personne ne le fasse et que le programme meure.

**Go fait autre chose : il renvoie deux valeurs.** Le résultat, et une erreur.

**Fig. 2.9** | Une fonction qui renvoie un résultat et une erreur possible.

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

**Ligne 11 : `(int, error)`.** Les parenthèses avec deux types signifient *« cette fonction renvoie deux
choses »* : un entier et une erreur.

**Ligne 13 : `return 0, errors.New(...)`.** Quand il y a un problème, elle renvoie une valeur quelconque (ici `0`,
qui ne sera pas utilisée) **et** une erreur décrivant ce qui s'est passé.

**Ligne 15 : `return a / b, nil`.** Quand tout va bien, elle renvoie le résultat **et `nil`**.

🔑 **`nil` signifie « rien », « vide », « il n'y en a pas ».** Donc la question `if err != nil` se lit :
*« l'erreur est-elle différente de rien ? »*, ou en français normal : **« y a-t-il eu une erreur ? »**

**Lignes 20-21 : le motif qui définit Go.**

```go
resultado, err := dividir(10, 2)
if err != nil {
    // algo salió mal: aquí se atiende
}
// si llegaste aquí, todo bien
```

Ce motif `if err != nil` est **la ligne la plus écrite de toute l'histoire de Go**. Il apparaît partout, et
certains le critiquent comme répétitif.

> [!NOTE]
> 🔧 **Observation de génie logiciel 2.3**
> La critique est juste : c'est verbeux. L'avantage est qu'**il est impossible d'ignorer une erreur par
> inadvertance**, parce qu'elle est là, sous ton nez, à la ligne suivant l'appel. Avec des exceptions, tu peux
> oublier un `catch` et ne rien savoir jusqu'à ce que le programme meure en production. Go a échangé la concision
> contre la sécurité, à dessein.

> [!TIP]
> 🧪 **Astuce de test et de débogage 2.2**
> Regarde la ligne 28 : elle dit `resultado, err = dividir(10, 0)` avec `=`, pas `:=`. C'est parce que les deux
> variables **existent déjà** depuis la ligne 20. Si tu mettais `:=` là, le compilateur te le dirait.

### 2.7 Tout assembler

**Fig. 2.10** | Programme complet qui utilise tout le contenu de la leçon.

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

Ce programme est déjà l'ancêtre direct du `revisor` : `clasificar` et `reporte` sont, en miniature, ce qui dans
la leçon 3 deviendra une méthode sur un struct `Servicio`, et dans la leçon 6 s'exécutera pour plusieurs services
**à la fois**.

**Ligne 21 : `%-10s` et `%-5d`.** Le nombre dit **combien d'espaces de largeur** réserver, et le signe moins
signifie **aligné à gauche**. Ainsi les colonnes restent droites. Sans cela, le tableau sortirait de travers.

> [!TIP]
> ✅ **Bonne pratique 2.5**
> Quand un appel est très long, coupe-le en plusieurs lignes comme aux lignes 21-22. Go le permet tant que la
> virgule reste à la fin de la ligne précédente.

> [!NOTE]
> 🚀 **Astuce de performance 2.1**
> `Sprintf` est pratique mais il n'est pas gratuit : il construit un nouveau texte en mémoire à chaque fois. Pour
> cinq lignes, cela n'a absolument aucune importance. Si un jour tu en construis des milliers, il existe
> `strings.Builder`, qui est beaucoup plus rapide. **N'optimise pas avant de mesurer** : c'est une information
> pour plus tard, pas quelque chose à faire aujourd'hui.

---

## L'erreur que tu vas voir

Presque tous les messages d'erreur de Go ont la même forme :

```
<archivo>:<línea>:<columna>: <qué esperaba el compilador y qué encontró>
```

Apprendre à lire cette ligne, c'est la moitié du travail. Deux exemples que tu as déjà vus dans cette leçon :

```
./fig02_05.go:8:5: declared and not used: edad
```

Cela se lit : *« dans le fichier `fig02_05.go`, ligne 8, colonne 5, tu as déclaré `edad` et tu ne l'as jamais
utilisée »*. La correction est simple : soit tu utilises la variable, soit tu la supprimes, soit —si tu as
vraiment besoin qu'elle soit déclarée mais que tu ne l'utilises pas encore— tu la remplaces par l'identifiant vide
`_`.

```
./fig02_06.go:8:14: cannot use "muchos" (untyped string constant) as int value in assignment
```

Cela se lit : *« à la ligne 8, colonne 12, tu as essayé d'utiliser le texte `"muchos"` là où un `int` était
attendu, et je ne te le permets pas »*. La correction n'est pas non plus facile à deviner si tu ne l'as jamais
vue : **il n'existe pas de conversion automatique** entre `string` et `int` dans une affectation comme celle-ci.
Il faudrait convertir explicitement (avec `strconv`, qu'on verra plus tard) ou, plus probablement, te rendre
compte que tu as mélangé le type par erreur.

Et un troisième, le plus courant de la première semaine, qui apparaît si tu réutilises `:=` sur une variable qui
existe déjà :

```
./fig02_04.go:10:9: no new variables on left side of :=
```

Cela se lit : *« du côté gauche du `:=` il n'y a aucune nouvelle variable »* — toutes existaient déjà, donc Go ne
sait pas quoi déclarer. La correction : utilise `=` au lieu de `:=` quand tu as déjà déclaré la variable
auparavant.

**La règle générale :** le compilateur de Go ne dit presque jamais « quelque chose s'est mal passé » dans
l'abstrait. Il dit le fichier, la ligne, la colonne et une phrase qui —même si elle sonne bizarrement la première
fois— décrit exactement le problème. La lire en entier, sans s'effrayer, résout la plupart des cas sans avoir à
chercher quoi que ce soit sur internet.

## Ce qu'on fait mal

- **Écrire `println` en minuscule au lieu de `Println`.** Go **distingue majuscules et minuscules** : `Println`
  et `println` sont des noms différents (`println` est une fonction interne du compilateur, pensée pour déboguer
  Go lui-même, pas ton programme). L'erreur est `undefined: fmt.println`.
- **Faire `import "fmt"` sans l'utiliser, ou utiliser `fmt.Println` sans l'`import`.** Le premier donne
  `"fmt" imported and not used` ; le second, `undefined: fmt`. Les deux sont le même type d'erreur : tu as dit au
  compilateur quelque chose qui ne correspond pas à ce que tu as fait.
- **Confondre `=` (affecter) avec `==` (comparer)** dans un `if`. En C, cela compilerait et ferait quelque chose
  d'inattendu (affecter la valeur et passer son chemin). **En Go, c'est une erreur de compilation** —
  `cannot use puerto = 443 (...) as value` —, ce qui est une bonne nouvelle : il t'arrête avant que le programme
  ne fasse quelque chose que tu n'as pas demandé.
- **Oublier le `\n` dans `Printf`.** Le programme **compile bien et s'exécute bien**, mais tout sort collé sur une
  seule ligne. `Println` passe à la ligne automatiquement ; `Printf` non. C'est une erreur silencieuse, pas une
  que le compilateur te signale.
- **Ignorer l'erreur avec `_`**, comme ceci : `resultado, _ := dividir(10, 0)`. **Cela compile aussi**, et le
  programme continue avec `resultado` valant `0` comme si tout allait bien. C'est la façon la plus rapide de créer
  une erreur impossible à retrouver ensuite, parce qu'il n'y a aucun message ni plantage : le programme continue
  simplement avec une donnée incorrecte. **Si tu jettes un jour une erreur avec `_`, que ce soit une décision
  consciente et commentée, jamais un réflexe pour que le compilateur cesse de se plaindre.**

## Exercices

### Questions de révision

Réponds sans regarder les réponses. Elles sont à la fin de cette section.

**2.1** Quelle est la différence entre un langage compilé et un langage interprété, et lequel des deux est Go ?

**2.2** Pourquoi la fonction où commence le programme doit-elle s'appeler exactement `main` ?

**2.3** Quelle différence y a-t-il entre `:=` et `=` ?

**2.4** Quelle erreur Go donne-t-il si tu déclares une variable et ne l'utilises pas ? Pourquoi penses-tu qu'il le
fait ?

**2.5** Qu'est-ce que la *valeur zéro* et quelle est celle de `string`, `int` et `bool` ?

**2.6** Qu'affiche ce programme ?

```go
package main

import "fmt"

func main() {
    var n int
    var s string
    fmt.Printf("[%d] [%q]\n", n, s)
}
```

**2.7** Quelle différence y a-t-il entre `Println`, `Printf` et `Sprintf` ?

**2.8** Que signifie `nil` et que demande `if err != nil` ?

**2.9** Ce programme a **trois** erreurs qui l'empêchent de compiler. Trouve-les sans l'exécuter.

```go
package main

func main() {
    nombre := "catalogo"
    puerto := 443
    puerto := 8080
    fmt.Println(nombre)
}
```

#### Réponses

**2.1** Un langage **interprété** se traduit ligne par ligne à chaque exécution ; un langage **compilé** se
traduit en entier une seule fois et produit un exécutable. **Go est compilé.**

**2.2** Parce que c'est une convention du langage : en exécutant un programme, Go cherche la fonction `main` du
paquet `main` pour démarrer. Avec un autre nom, il ne trouve pas par où commencer.

**2.3** `:=` **déclare** une nouvelle variable et lui affecte une valeur ; `=` ne fait qu'**affecter** à une
variable qui existe déjà.

**2.4** `declared and not used`. Il le fait parce qu'une variable inutilisée indique presque toujours une erreur
—un nom mal écrit ou des restes de code supprimé— et Go préfère t'arrêter plutôt que de laisser du code mort.

**2.5** C'est la valeur que reçoit automatiquement une variable déclarée sans valeur. `string` → `""`,
`int` → `0`, `bool` → `false`. Elle garantit qu'il n'y a **jamais** de déchets ni d'« indéfini ».

**2.6** `[0] [""]` — la valeur zéro de `int` et de `string`, et `%q` affiche les guillemets.

**2.7** `Println` affiche et passe à la ligne. `Printf` affiche avec un modèle de `%` et ne passe **pas** à la
ligne tout seul. `Sprintf` utilise le même modèle mais **renvoie** le texte au lieu de l'afficher.

**2.8** `nil` signifie « rien » / « vide ». `if err != nil` demande *« l'erreur n'est-elle pas vide ? »*, c'est-à-dire
**« y a-t-il eu une erreur ? »**.

**2.9** Les trois :
1. Il manque `import "fmt"`, et `fmt.Println` est utilisé.
2. `puerto := 8080` utilise `:=` sur une variable qui existe déjà → ce doit être `=`.
3. `puerto` est déclaré et **jamais utilisé** (seul `nombre` est affiché) → `declared and not used`.

### Exercices de code

> 🔴 **Ces sept exercices n'ont pas encore de solution de référence** — elle sera ajoutée lors d'une livraison
> ultérieure. On n'a pas inventé de code de solution afin de ne pas publier un exemple sans l'avoir d'abord
> exécuté.

**2.10** Écris un programme qui déclare ton nom, ton âge et si tu étudies ou travailles, et les affiche sur trois
lignes avec `Printf`, en utilisant le bon verbe pour chaque type.

**2.11** Écris `func celsiusAFahrenheit(c float64) float64` et teste-la avec 0, 37 et 100. La formule est
`f = c*9/5 + 32`. **Attention :** si tu écris `c*9/5` avec des entiers, le résultat est tronqué. Pourquoi cela
n'arrive-t-il pas ici ?

**2.12** Écris `func esPar(n int) bool`. Indice : l'opérateur `%` donne le reste d'une division, donc
`n % 2 == 0` est vrai pour les nombres pairs.

**2.13** Écris `func raizCuadrada(n float64) (float64, error)` qui renvoie une erreur si `n` est négatif.
Utilise `math.Sqrt` (tu as besoin de `import "math"`). Teste-la avec 16 et avec -4, en traitant l'erreur dans les
deux cas.

**2.14** Prends la **Fig. 2.10** et ajoute-lui une colonne qui dise `SI` ou `NO` selon que le port est 443 ou non.
Il te faut une nouvelle fonction et il faut ajuster le modèle de `Sprintf`.

**2.15 (Trouve l'erreur)** Chacun de ces fragments a un problème. Dis-le sans compiler, puis compile pour
confirmer :

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

**2.16 (Projet du cours)** C'est le premier pas du programme que tu vas construire tout au long du cours. Écris un
programme qui :

1. Ait une fonction `revisar(nombre string, puerto int, ms int) string` qui renvoie une ligne de rapport.
2. Ait une fonction `clasificar(ms int) string`.
3. Affiche un en-tête et cinq services.
4. À la fin, affiche combien des cinq étaient `"rapido"`. Tu auras besoin d'une variable compteur et d'un `if`.

Garde-le : dans la leçon 3 tu vas le réorganiser avec des structs.

## Comment savoir que j'y suis arrivé

- Tu as compilé et exécuté, toi-même et sans copier/coller, les figures 2.1 à 2.10 de cette leçon, et tu as obtenu
  exactement la sortie montrée.
- Tu peux expliquer, sans regarder le texte, pourquoi Go est un langage compilé et quelle conséquence pratique
  cela a pour attraper les erreurs avant la production.
- Tu sais par cœur quand utiliser `:=` et quand `=`, et quelle erreur donne le compilateur si tu te trompes.
- Tu peux nommer les quatre types de base (`string`, `int`, `float64`, `bool`) et la valeur zéro de chacun, sans
  consulter le tableau.
- Tu as écrit et exécuté ta propre version de l'exercice **2.16** (le premier pas du `revisor`) : il compile,
  affiche l'en-tête, les cinq lignes et le décompte des services rapides.
- Tu peux expliquer en deux phrases pourquoi `if err != nil` apparaît partout dans le code Go, et quel problème il
  évite par rapport aux exceptions d'autres langages.

Ouvre [`bitacora.md`](bitacora.md) et note **deux choses** : l'erreur du compilateur qui t'a le plus coûté à
comprendre, et quelque chose qui t'a surpris. Dans trois semaines cela te paraîtra évident et tu ne te
rappelleras plus pourquoi cela t'avait coûté — et c'est justement ce qu'il convient d'avoir écrit.

---

## Résumé

- Go est un langage **compilé** : il traduit tout ton code une fois et produit un exécutable indépendant.
- Tout fichier appartient à un **paquet** ; le paquet **`main`** est le seul qui produise un exécutable.
- L'exécution commence dans la fonction **`main`**.
- **`import`** apporte des paquets ; Go ne charge rien par défaut.
- **`:=`** déclare et affecte ; **`=`** ne fait qu'affecter.
- Go **ne compile pas** si tu déclares une variable sans l'utiliser, ni si tu mélanges les types.
- Les quatre types de base sont **`string`**, **`int`**, **`float64`** et **`bool`**.
- La **valeur zéro** garantit que toute variable naît avec une valeur valide : `""`, `0`, `0`, `false`.
- **`Println`** affiche avec saut de ligne ; **`Printf`** utilise un modèle avec `%` ; **`Sprintf`** renvoie le
  texte au lieu de l'afficher.
- Une fonction Go peut renvoyer **plusieurs valeurs**, et le motif standard est de renvoyer
  **`(résultat, error)`**.
- **`nil`** signifie « rien ». **`if err != nil`** est la façon idiomatique de vérifier s'il y a eu une erreur, et
  la ligne la plus écrite en Go.
- `go run` compile et exécute sans laisser de fichier ; `go build` laisse l'exécutable.

---

## Pour aller plus loin

1. **[A Tour of Go](https://go.dev/tour/)** — officiel, interactif. Fais la partie « Basics » jusqu'où tu en es
   dans cette leçon.
2. **[Go by Example — Variables](https://gobyexample.com/variables)** et
   **[Go by Example — Multiple Return Values](https://gobyexample.com/multiple-return-values)** — le même motif
   `(résultat, error)` avec des programmes minimaux qui s'exécutent.
3. **[Effective Go](https://go.dev/doc/effective_go)** — il n'est pas encore nécessaire de le lire en entier,
   mais consulte-le si quelque chose dans cette leçon t'a semblé une règle arbitraire : le pourquoi s'y trouve.

Si quelque chose n'est pas clair : lis le message d'erreur en entier (Go dit en général la ligne, la colonne et ce
qu'il attendait), cherche le concept dans Go by Example, et note le doute dans le journal de bord **même si tu ne
le résous pas**. Un doute écrit peut être résolu plus tard ; un doute oublié, non.

### Termes de cette leçon

| | |
|---|---|
| **compilateur** | programme qui traduit du code source en instructions machine |
| **erreur de compilation** | erreur qui empêche de produire l'exécutable ; elle est détectée avant l'exécution |
| **fonction** | bloc de code nommé qui reçoit des paramètres et peut renvoyer des valeurs |
| **identifiant vide (`_`)** | symbole qui jette une valeur à dessein |
| **`int`, `float64`, `string`, `bool`** | les quatre types de base |
| **langage compilé / interprété** | traduit tout une fois / ligne par ligne à chaque exécution |
| **`main` (fonction)** | point d'entrée du programme |
| **`main` (paquet)** | le seul paquet qui produit un exécutable |
| **`nil`** | absence de valeur |
| **paquet** | regroupement de code apparenté |
| **paramètre** | donnée qu'une fonction reçoit |
| **`Printf` / `Println` / `Sprintf`** | afficher avec modèle / afficher avec saut de ligne / renvoyer du texte |
| **type** | classe de donnée qu'une variable peut stocker |
| **valeur zéro** | valeur valide que reçoit toute variable déclarée sans valeur |
| **variable** | espace mémoire nommé qui stocke une donnée |
| **verbe de formatage (`%d`, `%s`, `%q`, `%t`, `%v`)** | trou dans un modèle de `Printf` |

---

**Précédent :** [Leçon 1 — Installer Go](01-instalacion.md) ·
**Suivant :** [Leçon 3 — Structs, méthodes, erreurs et interfaces](03-errores-interfaces.md)
