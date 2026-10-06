# Leçon 3 — Structs, méthodes, erreurs et interfaces

> **C'est ici que Go devient Go.** C'est la leçon la plus importante du cours.

**Durée :** deux séances de 60 minutes. Ne la fais pas d'une traite.

**À la fin, tu seras capable de :**

- Regrouper des données apparentées dans un **struct** et expliquer pourquoi c'est mieux que des variables
  isolées.
- Écrire des **méthodes** et décider quand utiliser un receveur par valeur et quand un receveur par pointeur.
- Expliquer ce qu'est un **pointeur** sans t'effrayer.
- Gérer les erreurs comme des valeurs, les envelopper et les interroger.
- Définir une **interface** et comprendre pourquoi, en Go, on ne déclare pas qu'elle est satisfaite.
- Écrire du code que l'on peut tester **sans se connecter à rien**.

---

## Pourquoi c'est important

Dans la leçon 2, tu as stocké les données d'un service dans trois variables séparées (`nombre`, `url`,
`puerto`). Cela fonctionne avec **un** service. Avec dix, cela fait trente variables isolées, et rien dans le code
ne dit lesquelles vont ensemble — si tu te trompes et combines le nom de l'un avec le port d'un autre, **le
programme compile quand même** et produit un rapport incorrect. C'est le premier problème que résout cette
leçon : le **struct**.

Le deuxième problème est plus profond. Le `revisor` va devoir interroger de vrais services, et les vraies
requêtes **échouent** : le réseau tombe, le service est lent, il répond avec un code que tu n'attendais pas. Un
programme qui ne sait pas gérer cela avec ordre ne sert à rien en production. Go résout cela d'une manière que tu
n'as probablement pas vue si tu viens d'un autre langage : **les erreurs sont des valeurs**, pas des exceptions
qui interrompent le flux.

Et le troisième problème est celui qui fait que les deux pièces précédentes s'emboîtent sans que le `revisor` ait
à savoir, à l'avance, **comment** chaque service va être vérifié (HTTP aujourd'hui, peut-être une base de données
demain). C'est ce que résout l'**interface**, qui est aussi ce qui te permettra, dans la leçon 5, de tester tout
le programme **sans te connecter à rien de réel**.

Structs, erreurs et interfaces sont, dans cet ordre, l'épine dorsale de tout ce qui suit dans le cours.

---

## Les concepts

### 3.1 Structs : regrouper ce qui va ensemble

Ce dont tu as besoin, c'est de dire au langage : *« ces trois choses sont un service »*.

**Fig. 3.1** | Définir et utiliser un struct.

```go
 1  // fig03_01.go
 2  // Agrupa los datos de un servicio en un solo tipo.
 3  package main
 4
 5  import "fmt"
 6
 7  // Servicio agrupa todo lo que describe a un servicio que vamos a revisar.
 8  type Servicio struct {
 9      Nombre string
10      URL    string
11      Puerto int
12  }
13
14  func main() {
15      s := Servicio{
16          Nombre: "catalogo",
17          URL:    "https://catalogo.example.com",
18          Puerto: 443,
19      }
20
21      fmt.Println("nombre:", s.Nombre)
22      fmt.Println("url:", s.URL)
23      fmt.Println("puerto:", s.Puerto)
24
25      s.Puerto = 8443          // se puede modificar
26      fmt.Println("nuevo puerto:", s.Puerto)
27
28      fmt.Println(s)           // e imprimir completo
29  }
```

```bash
$ go run fig03_01.go
nombre: catalogo
url: https://catalogo.example.com
puerto: 443
nuevo puerto: 8443
{catalogo https://catalogo.example.com 8443}
```

**Ligne 8 : `type Servicio struct {`.** Cela se lit : *« définis un nouveau type appelé `Servicio`, qui est une
structure »*. **Tu viens de créer un type qui n'était pas fourni avec le langage**, et désormais il s'utilise
comme `int` ou `string`.

**Lignes 15-19.** Crée une valeur de ce type. Les noms de champ suivis de deux-points sont facultatifs —tu
pourrais n'écrire que les valeurs dans l'ordre— mais **mets-les toujours** :

> [!TIP]
> ✅ **Bonne pratique 3.1**
> Écris toujours les noms des champs en créant un struct : `Servicio{Nombre: "x", Puerto: 443}` plutôt que
> `Servicio{"x", "", 443}`. La version courte se casse en silence si quelqu'un ajoute un champ ou change l'ordre,
> et le compilateur ne peut pas t'avertir parce que les types continuent de correspondre.

**Ligne 21 : `s.Nombre`.** Le point accède à un champ. Cela se lit « le `Nombre` de `s` ».

#### 3.1.1 La majuscule est une permission, pas de l'esthétique

Remarque que les champs commencent par une **majuscule** : `Nombre`, `URL`, `Puerto`. En Go, ce n'est pas un style
choisi : **c'est le contrôle d'accès du langage.**

| | |
|---|---|
| `Nombre` (majuscule) | **exporté** : visible depuis d'autres paquets |
| `nombre` (minuscule) | **non exporté** : visible uniquement à l'intérieur de son propre paquet |

Il n'y a pas de `public`, `private` ni `protected`. **La lettre initiale est la règle complète.**

#### 3.1.2 La valeur zéro d'un struct

**Fig. 3.2** | Un struct non initialisé.

```go
 1  // fig03_02.go
 2  package main
 3
 4  import "fmt"
 5
 6  type Servicio struct {
 7      Nombre string
 8      URL    string
 9      Puerto int
10  }
11
12  func main() {
13      var vacio Servicio
14      fmt.Printf("%+v\n", vacio)
15      fmt.Println("¿el nombre está vacío?", vacio.Nombre == "")
16      fmt.Println("puerto:", vacio.Puerto)
17  }
```

```bash
$ go run fig03_02.go
{Nombre: URL: Puerto:0}
¿el nombre está vacío? true
puerto: 0
```

Chaque champ reçoit **sa** valeur zéro : les textes restent à `""` et les nombres à `0`. **Il n'y a ni déchets ni
« indéfini »**, et c'est pourquoi un struct tout juste créé peut déjà s'utiliser sans crainte.

**Et voici `%+v`**, très utile pour déboguer : il affiche le struct **avec les noms des champs**. Compare-le avec
`%v`, qui n'affiche que les valeurs.

> [!TIP]
> 🧪 **Astuce de test et de débogage 3.1**
> Quand tu ne comprends pas ce que contient un struct, affiche-le avec `%+v`. C'est la façon la plus rapide de
> voir tous ses champs avec leur nom, et cela te fera gagner énormément de temps.

### 3.2 Méthodes

Une **méthode** est une fonction qui appartient à un type. Au lieu d'écrire `etiqueta(s)`, tu écris
`s.Etiqueta()`, et la fonction est liée au type auquel elle correspond.

**Fig. 3.3** | Méthodes sur un struct.

```go
 1  // fig03_03.go
 2  // Define metodos que pertenecen al tipo Servicio.
 3  package main
 4
 5  import "fmt"
 6
 7  type Servicio struct {
 8      Nombre string
 9      URL    string
10      Puerto int
11  }
12
13  // Etiqueta devuelve una descripcion legible del servicio.
14  func (s Servicio) Etiqueta() string {
15      return fmt.Sprintf("%s (%s:%d)", s.Nombre, s.URL, s.Puerto)
16  }
17
18  // EsSeguro indica si el servicio usa el puerto de HTTPS.
19  func (s Servicio) EsSeguro() bool {
20      return s.Puerto == 443
21  }
22
23  func main() {
24      a := Servicio{Nombre: "catalogo", URL: "https://catalogo.example.com", Puerto: 443}
25      b := Servicio{Nombre: "local", URL: "http://localhost", Puerto: 8080}
26
27      fmt.Println(a.Etiqueta(), "— seguro:", a.EsSeguro())
28      fmt.Println(b.Etiqueta(), "— seguro:", b.EsSeguro())
29  }
```

```bash
$ go run fig03_03.go
catalogo (https://catalogo.example.com:443) — seguro: true
local (http://localhost:8080) — seguro: false
```

**Ligne 14 : `func (s Servicio) Etiqueta() string`.** Ce `(s Servicio)` entre `func` et le nom s'appelle le
**receveur**, et c'est ce qui transforme une fonction en méthode. Cela se lit : *« cette fonction appartient au
type `Servicio`, et à l'intérieur je vais désigner la valeur par `s` »*.

> [!NOTE]
> 🔧 **Observation de génie logiciel 3.1**
> Si tu viens de Java, C# ou Python, une méthode Go ressemble à une méthode de classe — avec une différence
> importante : **en Go, la méthode s'écrit en dehors du type.** Tu peux avoir le `type` dans un fichier et ses
> méthodes dans un autre. Et comme il n'y a pas de classes, il n'y a pas non plus d'héritage : en Go, on réutilise
> le code autrement, et tu le verras dans la section 3.6.

### 3.3 Pointeurs, en dix minutes

Les pointeurs ont la réputation d'être difficiles. En Go, ils sont bien plus simples qu'en C, et tu n'as besoin de
comprendre qu'**une seule idée** pour la suite.

#### 3.3.1 Le problème

**Fig. 3.4** | Une méthode qui **ne fonctionne pas** comme tu le croirais.

```go
 1  // fig03_04.go
 2  // Muestra por que un receptor de VALOR no puede modificar el original.
 3  package main
 4
 5  import "fmt"
 6
 7  type Servicio struct {
 8      Nombre string
 9      Puerto int
10  }
11
12  // CambiarPuerto intenta modificar el servicio... y no lo logra.
13  func (s Servicio) CambiarPuerto(nuevo int) {
14      s.Puerto = nuevo
15  }
16
17  func main() {
18      x := Servicio{Nombre: "catalogo", Puerto: 443}
19      fmt.Println("antes:", x.Puerto)
20
21      x.CambiarPuerto(8443)
22      fmt.Println("despues:", x.Puerto)      // ¿cambió?
23  }
24
```

```bash
$ go run fig03_04.go
antes: 443
despues: 443
```

**Rien n'a changé.** La méthode s'est exécutée, il n'y a pas eu d'erreur, et la valeur est restée la même.

**Pourquoi ?** Parce que lorsque le receveur est `(s Servicio)`, Go remet à la méthode **une copie** du struct.
La méthode modifie la copie, la copie est jetée à la fin, et l'original n'en a jamais rien su. C'est un
anti-patron si fréquent qu'il a sa propre entrée dans « Ce qu'on fait mal », plus bas.

#### 3.3.2 La solution : le receveur par pointeur

Un **pointeur** est une variable qui stocke **l'adresse** d'une autre, au lieu d'une copie de son contenu.

Pense à la différence entre te donner **une photocopie** d'un document et te donner **l'adresse du classeur où
se trouve l'original**. Avec la photocopie, tu peux écrire tout ce que tu veux : l'original ne change pas. Avec
l'adresse, tu peux aller modifier l'original.

- `Servicio` est la photocopie.
- `*Servicio` est l'adresse de l'original.

**Fig. 3.5** | La même méthode, maintenant avec un receveur par pointeur.

```go
 1  // fig03_05.go
 2  // Con receptor de PUNTERO, el metodo si modifica el original.
 3  package main
 4
 5  import "fmt"
 6
 7  type Servicio struct {
 8      Nombre string
 9      Puerto int
10  }
11
12  // CambiarPuerto modifica el servicio. El * es la diferencia.
13  func (s *Servicio) CambiarPuerto(nuevo int) {
14      s.Puerto = nuevo
15  }
16
17  func main() {
18      x := Servicio{Nombre: "catalogo", Puerto: 443}
19      fmt.Println("antes:", x.Puerto)
20
21      x.CambiarPuerto(8443)
22      fmt.Println("despues:", x.Puerto)
23  }
```

```bash
$ go run fig03_05.go
antes: 443
despues: 8443
```

**La seule différence entre la Fig. 3.4 et la 3.5 est un astérisque à la ligne 13.** C'est tout.

Et remarque la ligne 21 : tu as écrit `x.CambiarPuerto(8443)` comme avant, **sans `&` ni rien de bizarre**. Go se
rend compte que la méthode a besoin d'un pointeur et le prend tout seul. C'est la raison pour laquelle les
pointeurs de Go donnent bien moins de travail que ceux de C.

> [!TIP]
> ✅ **Bonne pratique 3.2**
> La règle est simple : **si la méthode modifie le struct, receveur par pointeur (`*T`) ; si elle ne fait que
> lire, par valeur (`T`)**. Et sois cohérent au sein d'un même type : si la plupart de tes méthodes ont besoin
> d'un pointeur, utilise un pointeur partout, même si certaines n'en ont pas besoin. Les mélanger déroute celui
> qui lit le code.

> [!NOTE]
> 🚀 **Astuce de performance 3.1**
> Il y a une deuxième raison d'utiliser des pointeurs : **éviter de copier**. Si un struct a vingt champs, chaque
> appel avec receveur par valeur copie les vingt. Pour de petits structs c'est sans importance ; pour de grands
> ou dans des boucles de millions de tours, cela compte. **N'optimise pas cela sans mesurer :** la clarté vaut
> plus qu'une copie de 40 octets.

#### 3.3.3 Le seul pointeur dangereux : `nil`

Il reste une troisième situation avec les pointeurs, et celle-ci produit bien un vrai message d'erreur —si vrai
qu'il a sa propre section dédiée : **« L'erreur que tu vas voir »**, plus bas dans cette leçon. Nous en
anticipons l'idée : la valeur zéro d'un pointeur est **`nil`** (« je ne pointe sur rien »), et lire un champ à
travers un pointeur nul fait **planter** le programme au lieu de se comporter mal en silence comme dans la
Fig. 3.4.

### 3.4 Les erreurs sont des valeurs

Tu as déjà vu dans la leçon 2 qu'une fonction Go peut renvoyer `(résultat, error)`. Allons maintenant au fond des
choses, parce que **c'est ce qui distingue le plus Go de ce que tu as probablement déjà vu.**

#### 3.4.1 `error` est une interface, rien de magique

En Go, `error` est simplement un type avec une méthode :

```go
type error interface {
    Error() string
}
```

Cela signifie : *« toute chose qui a une méthode `Error()` renvoyant du texte est une erreur »*. Il n'y a pas de
hiérarchie de classes d'exception, pas de `throw`, pas de pile d'appels qui se déroule. **Une erreur est une
valeur ordinaire qui circule comme n'importe quelle autre.**

#### 3.4.2 Créer des erreurs

**Fig. 3.7** | Les trois façons de créer une erreur.

```go
 1  // fig03_07.go
 2  // Muestra las tres formas de producir un error.
 3  package main
 4
 5  import (
 6      "errors"
 7      "fmt"
 8  )
 9
10  // ErrPuertoInvalido es un error CENTINELA: se declara una vez y se compara.
11  var ErrPuertoInvalido = errors.New("el puerto debe estar entre 1 y 65535")
12
13  type Servicio struct {
14      Nombre string
15      Puerto int
16  }
17
18  // Validar revisa que el servicio tenga sentido.
19  func Validar(s Servicio) error {
20      if s.Nombre == "" {
21          // 1. un error sencillo, creado al momento
22          return errors.New("el nombre no puede estar vacio")
23      }
24      if s.Puerto < 1 || s.Puerto > 65535 {
25          // 2. un error centinela, para poder compararlo despues
26          return ErrPuertoInvalido
27      }
28      if s.Puerto == 80 {
29          // 3. un error con datos dentro
30          return fmt.Errorf("el puerto %d no usa cifrado, usa 443", s.Puerto)
31      }
32      return nil          // nil significa: todo bien
33  }
34
35  func main() {
36      casos := []Servicio{
37          {Nombre: "catalogo", Puerto: 443},
38          {Nombre: "", Puerto: 443},
39          {Nombre: "pagos", Puerto: 99999},
40          {Nombre: "viejo", Puerto: 80},
41      }
42
43      for _, s := range casos {
44          if err := Validar(s); err != nil {
45              fmt.Printf("%-10s ❌ %v\n", s.Nombre, err)
46          } else {
47              fmt.Printf("%-10s ✅ valido\n", s.Nombre)
48          }
49      }
50  }
```

```bash
$ go run fig03_07.go
catalogo   ✅ valido
           ❌ el nombre no puede estar vacio
pagos      ❌ el puerto debe estar entre 1 y 65535
viejo      ❌ el puerto 80 no usa cifrado, usa 443
```

**Ligne 11 : l'erreur sentinelle.** Elle est déclarée **une seule fois**, au niveau du paquet, avec le préfixe
`Err`. Elle sert à ce que celui qui appelle ta fonction puisse demander *« était-ce cette erreur en
particulier ? »*, comme tu le verras en 3.5.

**Ligne 30 : `fmt.Errorf`.** Comme `Printf`, mais il produit une erreur au lieu d'afficher. Utilise-le quand le
message a besoin de données.

**Ligne 44 : `if err := Validar(s); err != nil`.** Ce motif déclare `err` **à l'intérieur** du `if`, donc il
n'existe que là. C'est très idiomatique en Go et cela garde le code propre.

> [!TIP]
> ✅ **Bonne pratique 3.3**
> Les messages d'erreur s'écrivent **en minuscules et sans point final** : `"no se pudo abrir el archivo"`, et
> non `"No se pudo abrir el archivo."`. La raison est pratique : les erreurs s'**enveloppent** les unes dans les
> autres (section 3.5), et en se concaténant elles donnent `"revisando catalogo: no se pudo abrir el archivo"`.
> Avec des majuscules et des points, le résultat paraîtrait cassé.

### 3.5 Envelopper les erreurs : `%w`, `errors.Is` et `errors.As`

Une erreur sans contexte est peu utile. Si ton programme dit `connection refused`, tu ne sais pas **quel**
service a échoué.

**Envelopper** une erreur, c'est lui ajouter du contexte **sans perdre l'originale**.

**Fig. 3.8** | Envelopper des erreurs et les interroger.

```go
 1  // fig03_08.go
 2  // Envuelve un error para agregar contexto sin perder el original.
 3  package main
 4
 5  import (
 6      "errors"
 7      "fmt"
 8  )
 9
10  var ErrNoResponde = errors.New("no responde")
11
12  type Servicio struct{ Nombre string }
13
14  // consultar simula la consulta de bajo nivel.
15  func consultar(s Servicio) error {
16      return ErrNoResponde
17  }
18
19  // Revisar agrega contexto al error de consultar.
20  func Revisar(s Servicio) error {
21      if err := consultar(s); err != nil {
22          // el %w ENVUELVE el error original
23          return fmt.Errorf("revisando %s: %w", s.Nombre, err)
24      }
25      return nil
26  }
27
28  func main() {
29      err := Revisar(Servicio{Nombre: "catalogo"})
30
31      fmt.Println("1. el mensaje completo:")
32      fmt.Println("  ", err)
33
34      fmt.Println("2. ¿es un ErrNoResponde, aunque este envuelto?")
35      fmt.Println("  ", errors.Is(err, ErrNoResponde))
36
37      fmt.Println("3. ¿y si pregunto por otro error?")
38      fmt.Println("  ", errors.Is(err, errors.New("otra cosa")))
39
40      fmt.Println("4. el error original, desenvuelto:")
41      fmt.Println("  ", errors.Unwrap(err))
42  }
```

```bash
$ go run fig03_08.go
1. el mensaje completo:
   revisando catalogo: no responde
2. ¿es un ErrNoResponde, aunque este envuelto?
   true
3. ¿y si pregunto por otro error?
   false
4. el error original, desenvuelto:
   no responde
```

**Ligne 23 : `%w`.** C'est le verbe clé. À l'affichage, il ressemble à `%v`, **mais il conserve l'erreur
originale à l'intérieur** pour qu'on puisse l'interroger ensuite.

**Ligne 35 : `errors.Is`.** Demande *« à un endroit quelconque de cette chaîne se trouve-t-il cette erreur ? »*.
Cela fonctionne même s'il y a cinq couches d'enveloppe.

> [!WARNING]
> 🔴 **Erreur courante de programmation 3.3 — la plus silencieuse de cette leçon**
> Écrire `%v` au lieu de `%w` en enveloppant :
>
> ```go
> return fmt.Errorf("revisando %s: %v", s.Nombre, err)   // ❌ con %v
> ```
>
> **Le message affiché est identique.** Il n'y a pas d'erreur, pas d'avertissement, tout semble fonctionner. Mais
> `errors.Is` cesse de trouver l'erreur originale et renvoie `false`, de sorte que le code qui décidait quoi
> faire selon le type de défaillance se met à prendre le mauvais chemin. **C'est le genre d'erreur qui n'échoue
> pas : elle renvoie une donnée pire.** Quand tu enveloppes, utilise `%w`. C'est le deuxième anti-patron de la
> section « Ce qu'on fait mal ».

#### 3.5.1 `errors.As`, quand tu as besoin des données de l'erreur

`errors.Is` répond « est-ce cette erreur ? ». `errors.As` répond « est-elle de ce **type** ? donne-la-moi pour
lire ses champs ».

**Fig. 3.9** | Une erreur personnalisée avec des données, récupérée avec `errors.As`.

```go
 1  // fig03_09.go
 2  // Define un tipo de error propio y recupera sus datos.
 3  package main
 4
 5  import (
 6      "errors"
 7      "fmt"
 8  )
 9
10  // ErrorHTTP es un error que lleva datos adentro.
11  type ErrorHTTP struct {
12      Codigo int
13      URL    string
14  }
15
16  // Error hace que ErrorHTTP cumpla la interfaz error.
17  func (e *ErrorHTTP) Error() string {
18      return fmt.Sprintf("el servidor respondio %d", e.Codigo)
19  }
20
21  func consultar(url string) error {
22      return &ErrorHTTP{Codigo: 503, URL: url}
23  }
24
25  func Revisar(nombre, url string) error {
26      if err := consultar(url); err != nil {
27          return fmt.Errorf("revisando %s: %w", nombre, err)
28      }
29      return nil
30  }
31
32  func main() {
33      err := Revisar("catalogo", "https://catalogo.example.com")
34      fmt.Println("mensaje:", err)
35
36      var errHTTP *ErrorHTTP
37      if errors.As(err, &errHTTP) {
38          fmt.Println("es un ErrorHTTP")
39          fmt.Println("  codigo:", errHTTP.Codigo)
40          fmt.Println("  url:   ", errHTTP.URL)
41
42          if errHTTP.Codigo >= 500 {
43              fmt.Println("  → es culpa del servidor, conviene reintentar")
44          }
45      }
46  }
```

```bash
$ go run fig03_09.go
mensaje: revisando catalogo: el servidor respondio 503
es un ErrorHTTP
  codigo: 503
  url:    https://catalogo.example.com
  → es culpa del servidor, conviene reintentar
```

**Ligne 17.** En écrivant une méthode `Error() string`, ton type **est déjà une `error`**. Tu n'as rien déclaré :
le type satisfait l'interface parce qu'il a la méthode. C'est ce qu'explique la section suivante.

**Ligne 42.** Et voici la vraie valeur : le programme peut **décider** selon le type de défaillance. Un 503 se
réessaie ; un 404, non. Avec des erreurs en simple texte, ce serait impossible sans comparer des chaînes, ce qui
est fragile.

### 3.6 Interfaces : le cœur de Go

Une **interface** est une liste de méthodes. Tout type qui a ces méthodes **la satisfait**, et il n'y a rien à
déclarer nulle part.

**Fig. 3.10** | Une interface et deux implémentations.

```go
 1  // fig03_10.go
 2  // Una interfaz con dos implementaciones: la real y una de prueba.
 3  package main
 4
 5  import (
 6      "errors"
 7      "fmt"
 8  )
 9
10  type Servicio struct {
11      Nombre string
12      URL    string
13  }
14
15  type Estado struct {
16      Servicio Servicio
17      Codigo   int
18      Err      error
19  }
20
21  // Revisor es una INTERFAZ: cualquier cosa con este metodo la cumple.
22  type Revisor interface {
23      Revisar(s Servicio) Estado
24  }
25
26  // ---- primera implementacion: la de verdad (simplificada) ----
27  type RevisorHTTP struct{}
28
29  func (r RevisorHTTP) Revisar(s Servicio) Estado {
30      // En la leccion 7 esto hara una peticion HTTP real.
31      return Estado{Servicio: s, Codigo: 200}
32  }
33
34  // ---- segunda implementacion: para probar, sin red ----
35  type RevisorFalso struct {
36      Respuesta Estado
37  }
38
39  func (r RevisorFalso) Revisar(s Servicio) Estado {
40      r.Respuesta.Servicio = s
41      return r.Respuesta
42  }
43
44  // RevisarTodos acepta CUALQUIER Revisor. No sabe ni le importa cual.
45  func RevisarTodos(r Revisor, servicios []Servicio) []Estado {
46      var estados []Estado
47      for _, s := range servicios {
48          estados = append(estados, r.Revisar(s))
49      }
50      return estados
51  }
52
53  func main() {
54      servicios := []Servicio{
55          {Nombre: "catalogo", URL: "https://catalogo.example.com"},
56          {Nombre: "pagos", URL: "https://pagos.example.com"},
57      }
58
59      fmt.Println("--- con el revisor real ---")
60      for _, e := range RevisarTodos(RevisorHTTP{}, servicios) {
61          fmt.Printf("  %-10s codigo %d\n", e.Servicio.Nombre, e.Codigo)
62      }
63
64      fmt.Println("--- con el falso, simulando una falla ---")
65      falso := RevisorFalso{
66          Respuesta: Estado{Err: errors.New("no responde")},
67      }
68      for _, e := range RevisarTodos(falso, servicios) {
69          fmt.Printf("  %-10s error: %v\n", e.Servicio.Nombre, e.Err)
70      }
71  }
```

```bash
$ go run fig03_10.go
--- con el revisor real ---
  catalogo   codigo 200
  pagos      codigo 200
--- con el falso, simulando una falla ---
  catalogo   error: no responde
  pagos      error: no responde
```

**Lis-le lentement, parce qu'ici se trouve l'idée qui soutient tout Go.**

**Lignes 22-24 : l'interface.** Elle dit : *« un `Revisor` est toute chose qui a une méthode `Revisar` qui reçoit
un `Servicio` et renvoie un `Estado` »*.

**Lignes 29 et 39.** `RevisorHTTP` et `RevisorFalso` **ne déclarent nulle part** qu'ils satisfont `Revisor`. Il
n'y a pas de `implements`, pas de `: Revisor`, rien. **Ils satisfont l'interface parce qu'ils ont la méthode**, et
cela, le compilateur le vérifie tout seul.

**Ligne 45 : `func RevisarTodos(r Revisor, ...)`.** Cette fonction **ne sait pas** avec quoi elle travaille. Elle
sait seulement qu'elle peut appeler `.Revisar()`. Et c'est pourquoi les lignes 60 et 68 lui passent deux choses
complètement différentes **sans changer une seule ligne de `RevisarTodos`**.

> [!NOTE]
> 🔧 **Observation de génie logiciel 3.2 — la plus importante de la leçon**
> Cette satisfaction implicite a une conséquence qui n'existe ni en Java ni en C# : **tu peux définir une
> interface pour du code que tu n'as pas écrit.** Si une bibliothèque tierce a un type avec une méthode
> `Revisar`, ce type satisfait *ton* interface sans que l'auteur l'ait su ni ait eu à coopérer. Dans d'autres
> langages, si l'auteur n'a pas déclaré l'interface, il n'y a rien à faire.

> [!TIP]
> ✅ **Bonne pratique 3.4 — les deux règles d'or des interfaces en Go**
> **1. Fais-les petites.** Une ou deux méthodes. Une interface de dix méthodes vient presque toujours d'un autre
> langage — et c'est, de fait, le troisième anti-patron de la section « Ce qu'on fait mal ». Celles de la
> bibliothèque standard les plus utilisées —`io.Reader`, `io.Writer`— ont **une** méthode.
> **2. Définis-les là où on les UTILISE, pas là où on les implémente.** L'interface `Revisor` appartient au code
> qui a besoin de vérifier des choses, pas à celui qui sait comment les vérifier. Cela inverse la dépendance :
> celui qui consomme déclare ce dont il a besoin.
>
> Et le résumé que tu entendras souvent : **« accepte des interfaces, renvoie des structs »**.

### 3.7 Pourquoi cela rend ton code testable

Regarde encore la ligne 65 de la Fig. 3.10. Tu viens de tester `RevisarTodos` **en simulant un service en
panne**, sans rien éteindre, sans réseau, et en une milliseconde.

C'est ce que les gens veulent dire quand ils parlent de « code testable », et en Go on y arrive **sans
bibliothèques de mocks, sans annotations et sans frameworks**. Seulement avec une petite interface.

> [!NOTE]
> 🔧 **Observation de génie logiciel 3.3**
> La question qu'il convient de se poser en concevant : *« puis-je tester cela sans que le monde extérieur
> existe ? »* Si la réponse est non, il manque normalement une interface. Dans la leçon 5, tu écriras de vrais
> tests, et tu seras content d'avoir fait cela maintenant.

---

## L'erreur que tu vas voir

L'erreur la plus fréquente de cette leçon est le **`nil pointer dereference`**, et contrairement à l'anti-patron
de la Fig. 3.4 (qui échoue en silence), celle-ci te prévient — avec un message qui effraie d'abord plus qu'il ne
devrait.

**Fig. 3.6** | Un programme qui **compile parfaitement** et **plante à l'exécution**.

```go
 1  // fig03_06.go — este programa COMPILA pero truena al correr
 2  package main
 3
 4  import "fmt"
 5
 6  type Servicio struct{ Puerto int }
 7
 8  func main() {
 9      var p *Servicio          // un puntero sin apuntar a nada: vale nil
10      fmt.Println(p)           // esto sí funciona
11      fmt.Println(p.Puerto)    // esto truena
12  }
```

```bash
$ go run fig03_06.go
<nil>
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x...]

goroutine 1 [running]:
main.main()
	/home/usuario/fig03_06.go:11 +0x18
exit status 2
```

**La valeur zéro d'un pointeur est `nil`** : « je ne pointe sur rien ». Tenter de lire un champ à travers un
pointeur nul provoque un **panic**, qui est la façon dont Go interrompt un programme quand quelque chose
d'irrécupérable se produit.

**Comment lire ce message, ligne par ligne :**

- `panic: runtime error: invalid memory address or nil pointer dereference` — le nom du problème. Quand tu le
  verras, la cause est presque toujours un pointeur, une map ou une interface à `nil` que tu as essayé d'utiliser
  comme si elle contenait quelque chose.
- `goroutine 1 [running]:` — quel fil d'exécution de Go tournait quand ça a explosé. Pour l'instant, dans un
  programme sans concurrence, ce sera toujours le 1 (la fonction `main`) ; dans la leçon 6, avec plusieurs
  goroutines, ce numéro commence à compter pour de bon.
- `main.main() … fig03_06.go:11` — **le fichier et la ligne exacts** où cela s'est produit. Commence toujours par
  là, pas par le message du haut : le message te dit *quel type* d'erreur c'était, cette ligne te dit *où*.
- `exit status 2` — le programme s'est terminé en erreur (un `exit status 0` est un succès).

> [!WARNING]
> ⚠️ **Erreur courante de programmation 3.2**
> Le `nil pointer dereference` est le panic le plus fréquent en Go. Remarque qu'il **compile parfaitement** : le
> compilateur ne peut pas savoir si un pointeur vaudra `nil` à l'exécution. Quand tu le vois, va directement au
> fichier et à la ligne que te donne la trace — dans l'exemple, `fig03_06.go:11`.

---

## Ce qu'on fait mal

**1. Écrire une méthode qui devrait modifier le struct, avec un receveur par valeur.** C'est l'erreur de la
Fig. 3.4 (section 3.3.1) : elle compile, s'exécute, et **ne fait absolument rien** — sans aucun message. C'est
l'un des faux pas les plus déconcertants pour qui débute, précisément parce qu'il n'y a aucun symptôme. La règle
qui l'évite est dans la Bonne pratique 3.2 : si la méthode modifie, receveur par pointeur.

**2. Envelopper une erreur avec `%v` au lieu de `%w`.** C'est l'erreur de la section 3.5 : le message affiché
reste identique, donc rien n'avertit, mais `errors.Is` cesse de reconnaître l'erreur originale et le code qui
décidait selon le type de défaillance se met à prendre le mauvais chemin. Des deux anti-patrons de cette leçon,
celui-ci est le plus coûteux, parce que personne ne remarque le symptôme avant que le programme ait déjà pris une
mauvaise décision avec lui.

**3. Écrire de grandes interfaces, avec beaucoup de méthodes, en copiant l'habitude d'un autre langage.** En Go,
les interfaces se font petites —une ou deux méthodes, comme tu l'as vu dans la Bonne pratique 3.4— et se
définissent du côté de celui qui les consomme, pas de celui qui les implémente. Une interface de dix méthodes
n'est presque jamais nécessaire : il suffit normalement de la ou des méthodes que la fonction qui la reçoit
appelle réellement.

---

## Exercices

### Exercices de révision

**3.1** Quel avantage un struct a-t-il sur trois variables isolées ?

**3.2** Que signifie le fait qu'un champ commence par une majuscule ?

**3.3** Qu'affiche `fmt.Printf("%+v\n", Servicio{})` si le struct a `Nombre string` et `Puerto int` ?

**3.4** Quelle est la différence entre `func (s Servicio) X()` et `func (s *Servicio) X()`, et quand utilise-t-on
chacun ?

**3.5** Cette méthode compile, s'exécute et ne fait rien. Pourquoi ?

```go
func (s Servicio) Renombrar(nuevo string) {
    s.Nombre = nuevo
}
```

**3.6** Qu'est-ce que `nil` pour un pointeur et que se passe-t-il si tu lis un champ à travers lui ?

**3.7** Que doit avoir un type pour être une `error` ?

**3.8** Quelle est la différence entre `%w` et `%v` en enveloppant une erreur, et pourquoi est-elle dangereuse ?

**3.9** Quand utilises-tu `errors.Is` et quand `errors.As` ?

**3.10** Que faut-il écrire pour qu'un type satisfasse une interface en Go ?

**3.11** Qu'affiche ce programme ?

```go
type Estado struct{ Codigo int }

func (e Estado) OK() bool { return e.Codigo == 200 }

func main() {
    var e Estado
    fmt.Println(e.OK())
}
```

### Exercices de code

**3.12** Ajoute au struct `Servicio` un champ `TimeoutMs int` et une méthode `Timeout() time.Duration` qui le
convertisse. Indice : `time.Duration(s.TimeoutMs) * time.Millisecond`.

**3.13** Écris une méthode `func (s *Servicio) Normalizar()` qui : retire les espaces du nom avec
`strings.TrimSpace`, le passe en minuscules avec `strings.ToLower`, et si le port est 0 le met à 443.
**Explique-toi pourquoi cette méthode a besoin d'un receveur par pointeur.**

**3.14** Définis un type `Estado` avec les champs `Servicio`, `Codigo`, `Duracion` et `Err`. Écris-lui une
méthode `OK() bool` qui renvoie vrai uniquement s'il n'y a pas d'erreur **et** que le code est entre 200 et 299.
Teste-la avec cinq cas, y compris l'`Estado{}` vide.

**3.15** Crée une erreur sentinelle `ErrTimeout` et une fonction qui la renvoie enveloppée avec du contexte.
Vérifie avec `errors.Is` qu'elle est détectée. **Ensuite, remplace le `%w` par `%v` et vérifie que `errors.Is`
renvoie `false`.** Note dans le journal de bord que le message affiché n'a pas changé.

**3.16** Définis une erreur personnalisée `ErrorValidacion` avec les champs `Campo string` et `Motivo string`,
fais-lui satisfaire l'interface `error`, et récupère-la avec `errors.As` pour afficher quel champ a échoué.

**3.17** Définis l'interface `Notificador` avec une méthode `Notificar(mensaje string) error`. Implémente
`NotificadorConsola` (qui affiche) et `NotificadorFalso` (qui stocke les messages dans un slice pour pouvoir les
examiner). Écris une fonction qui reçoit un `Notificador` et utilise-la avec les deux.

**3.18 (Trouve l'erreur)** Dis ce qui ne va pas dans chacun **sans compiler**, puis compile pour confirmer :

```go
// (a)
func (s Servicio) Renombrar(n string) { s.Nombre = n }

// (b)
return fmt.Errorf("fallo al revisar %s: %v", nombre, err)

// (c)
var p *Servicio
fmt.Println(p.Nombre)

// (d)
type Revisor interface {
    Revisar(s Servicio) Estado
}
type MiRevisor struct{}
func (m MiRevisor) revisar(s Servicio) Estado { return Estado{} }
// y luego:  var r Revisor = MiRevisor{}
```

**3.19 (Projet du cours)** Réorganise ton programme de la leçon 2 en utilisant le contenu de cette leçon :
1. Les structs `Servicio` et `Estado`.
2. Les méthodes `Etiqueta()`, `EsSeguro()` et `OK()`.
3. L'interface `Revisor` avec `RevisorHTTP` (qui pour l'instant renvoie des données inventées) et `RevisorFalso`.
4. La fonction `RevisarTodos(r Revisor, servicios []Servicio) []Estado`.
5. Un `main` qui affiche le rapport en utilisant le **faux**, avec un service qui échoue.

🔑 **Quand tu auras terminé, remarque quelque chose : ton programme peut déjà être testé en entier sans se
connecter à rien, et tu n'as pas encore écrit un seul test.** C'est ce que tu viens de gagner dans cette leçon.

### Corrigés

**Exercices de révision (3.1 à 3.11) :**

**3.1** Il regroupe les données qui vont ensemble en une seule variable, de sorte que le langage —et celui qui lit
le code— sait qu'elles appartiennent à la même chose. Et on ne peut pas combiner par erreur des données de deux
services différents.

**3.2** Qu'il est **exporté** : il est visible depuis d'autres paquets. Avec une minuscule, il n'est visible qu'à
l'intérieur de son paquet. C'est tout le contrôle d'accès que possède Go.

**3.3** `{Nombre: Puerto:0}` — les valeurs zéro, avec les noms des champs parce que c'est `%+v`.

**3.4** Le premier est un **receveur par valeur** : il reçoit une copie et ne peut pas modifier l'original. Le
second est un **receveur par pointeur** : il reçoit l'adresse et le peut. On utilise un pointeur quand la méthode
modifie, ou quand le struct est grand et que le copier coûte.

**3.5** Parce que le receveur est **par valeur** : la méthode modifie une copie qui est jetée à la fin. Il faut le
changer en `func (s *Servicio) Renombrar(...)`.

**3.6** `nil` est la valeur zéro d'un pointeur et signifie « je ne pointe sur rien ». Lire un champ à travers un
pointeur `nil` provoque un **panic** (`nil pointer dereference`) et le programme s'interrompt.

**3.7** Une méthode `Error() string`. Rien de plus : `error` est une interface avec cette unique méthode.

**3.8** `%w` **enveloppe** l'erreur originale et la conserve à l'intérieur ; `%v` ne fait que la convertir en
texte. C'est dangereux parce que **le message affiché est identique**, donc il n'y a pas de symptôme — mais
`errors.Is` cesse de trouver l'erreur originale et le code qui décidait selon le type de défaillance se met à se
tromper.

**3.9** `errors.Is` pour demander *« est-ce cette erreur en particulier ? »*. `errors.As` pour demander *« est-elle
de ce type ? donne-la-moi »*, quand tu as besoin de lire les données que l'erreur porte en elle.

**3.10** **Rien.** Il suffit d'avoir les méthodes que demande l'interface ; le compilateur le vérifie tout seul.
Le `implements` n'existe pas.

**3.11** `false` — la valeur zéro d'`Estado` a `Codigo: 0`, qui n'est pas 200. Remarque que la méthode fonctionne
parfaitement sur un struct vide : c'est la valeur zéro qui est utile.

**Exercices de code (3.12 à 3.19) :** leurs corrigés vérifiés —compilés et exécutés— seront ajoutés lors de la
livraison qui clôturera le chapitre des exercices avec corrigé du cours complet ; ils ne sont pas publiés sans
avoir lancé `go build` sur chacun.

---

## Comment savoir que j'y suis arrivé

- Tu as compilé et exécuté les figures 3.1 à 3.10 et ta sortie correspond à celle montrée.
- Tu peux expliquer, sans regarder le texte, quelle différence il y a entre receveur par valeur et receveur par
  pointeur, et donner un exemple de quand utiliser chacun.
- Tu as fait l'expérience de l'exercice 3.15 (remplacer `%w` par `%v`) et tu as vu de tes propres yeux
  qu'`errors.Is` cesse de trouver l'erreur **sans que le message affiché change**.
- Tu peux expliquer à quelqu'un d'autre, sans jargon, pourquoi en Go « satisfaire une interface » ne se déclare
  nulle part.
- Tu as terminé l'exercice 3.19 : ton programme `revisor` utilise déjà des structs, des méthodes, l'interface
  `Revisor` et peut s'exécuter avec le `RevisorFalso` sans toucher au réseau.

**Avant de clore :** dans [`bitacora.md`](bitacora.md), note ce qui t'a le plus coûté entre pointeurs, erreurs et
interfaces, et le résultat de l'exercice 3.15 —celui de `%w` contre `%v`—. Cette expérience est celle qu'on oublie
le plus et celle qui coûte le plus cher dans un vrai programme.

---

## Résumé

- Un **struct** regroupe des données apparentées en un type propre. Il se définit avec
  `type Nombre struct { … }`.
- Les champs avec **majuscule** sont exportés ; avec minuscule, privés au paquet.
- La **valeur zéro** d'un struct remplit chaque champ de sa propre valeur zéro : il n'y a jamais de déchets.
- **`%+v`** affiche un struct avec les noms de ses champs : le meilleur outil de débogage.
- Une **méthode** est une fonction avec **receveur** : `func (s Servicio) X()`.
- Le **receveur par valeur** reçoit une copie et **ne peut pas modifier** l'original ; le **receveur par
  pointeur** (`*T`) le peut.
- Un **pointeur** stocke l'adresse d'une autre variable. Go insère le `&` et le `*` pour toi lors de l'appel des
  méthodes.
- La valeur zéro d'un pointeur est **`nil`** ; lire à travers lui provoque un **panic**.
- **`error` est une interface** avec une seule méthode `Error() string`. Une erreur est une valeur ordinaire.
- On les crée avec **`errors.New`** (simples), comme **sentinelles** (`var ErrX = errors.New(...)`) ou avec
  **`fmt.Errorf`** (avec des données).
- **`%w`** enveloppe une erreur en conservant l'originale ; **`%v` l'écrase et casse `errors.Is` sans
  avertir**.
- **`errors.Is`** demande si une erreur est dans la chaîne ; **`errors.As`** récupère l'erreur d'un type concret
  pour lire ses données.
- Une **interface** est une liste de méthodes. Un type la satisfait **en ayant les méthodes**, sans le déclarer.
- Les interfaces de Go se font **petites** et se définissent **là où on les utilise**.
- Une petite interface permet de **tester sans le monde extérieur** : on remplace l'implémentation réelle par
  une fausse.

---

## Pour aller plus loin

1. **[A Tour of Go — Methods](https://go.dev/tour/methods/1)** — officiel et interactif : receveurs, pointeurs et
   interfaces avec des exercices dans le navigateur.
2. **[Effective Go — Interfaces and other types](https://go.dev/doc/effective_go#interfaces_and_types)**
   — l'explication officielle de pourquoi la satisfaction des interfaces est implicite.
3. **[Error handling and Go](https://go.dev/blog/error-handling-and-go)** — l'article du blog officiel de Go sur
   pourquoi les erreurs sont des valeurs et comment on les gère de façon idiomatique.
4. **[Paquet `errors` — documentation officielle](https://pkg.go.dev/errors)** — la référence exacte de
   `errors.Is`, `errors.As` et `errors.Unwrap`.
5. **[Des API prêtes pour les agents d'IA](https://www.habil.mx/fr/blog/apis-pretes-pour-les-agents-d-ia/)** — article sur les erreurs d'API qui indiquent au client s'il faut réessayer (un 503) ou non (un 404), la distinction de la section 3.5.

### Termes de cette leçon

| | |
|---|---|
| **champ** | chaque donnée que contient un struct |
| **envelopper (une erreur)** | lui ajouter du contexte en conservant l'originale, avec `%w` |
| **erreur sentinelle** | erreur déclarée une fois au niveau du paquet, pour comparer avec `errors.Is` |
| **exporté / non exporté** | visible hors du paquet (majuscule) ou non (minuscule) |
| **interface** | liste de méthodes ; un type la satisfait en ayant ces méthodes |
| **méthode** | fonction associée à un type au moyen d'un receveur |
| **`nil`** | absence de valeur ; la valeur zéro des pointeurs, interfaces et erreurs |
| **`nil pointer dereference`** | panic pour lecture à travers un pointeur nul |
| **panic** | arrêt du programme dû à une erreur irrécupérable |
| **pointeur** | variable qui stocke l'adresse d'une autre ; son type s'écrit `*T` |
| **receveur** | le `(s Servicio)` qui lie une fonction à un type |
| **satisfaction implicite** | satisfaire une interface sans le déclarer |
| **struct** | type qui regroupe plusieurs champs |

---

**Précédent :** [Leçon 2 — Variables, fonctions et types](02-fundamentos.md) ·
**Suivant :** [Leçon 4 — Collections](04-colecciones.md)
