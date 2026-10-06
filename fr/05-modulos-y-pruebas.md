# Leçon 5 — Modules et tests

**Durée :** 90 minutes, ou 2 séances de 45.

**À la fin, tu seras capable de :**

- Expliquer ce que résout `go.mod` et pourquoi le `revisor` n'a pas besoin de `go.sum`.
- Réorganiser un programme d'un seul fichier en paquets sous `internal/`, et expliquer ce que ce dossier
  t'apporte qu'un dossier normal n'apporte pas.
- Écrire des tests avec une table de cas, sans aucune bibliothèque externe, et lire ce qu'ils rapportent
  quand ils échouent.
- Lancer `go test` avec `-v`, `-run`, `-cover` et `-race`, et expliquer ce que mesure chaque option.
- Reconnaître le piège de « zéro test » qui a exactement l'air de « tous ont réussi ».
- Décider, avec un vrai chiffre de couverture sous les yeux, quelle partie de ce trou il importe de corriger
  et laquelle non.

---

## Pourquoi c'est important

Jusqu'à la leçon 4, tu avais un unique fichier, `main.go`, avec tout dedans : les structs `Servicio` et
`Estado`, l'interface `Revisor`, la fonction qui construit le rapport. Cela fonctionnait, et cela
fonctionnait bien —mais un seul fichier a un plafond. Dès que tu veux **tester** une pièce sans exécuter le
programme complet (sans toucher au réseau, sans lire un vrai fichier), un seul `main.go` ne te le permet
pas : tout est mélangé avec tout.

C'est la leçon où le `revisor` cesse d'être un exercice d'un seul fichier et devient **un vrai projet** :
plusieurs paquets, chacun avec une responsabilité, et une suite de tests qui démontre que chaque pièce fait
ce qu'elle dit sans avoir besoin des autres. C'est le même saut que tu as fait dans la leçon 1, de « un
programme qui tient dans la tête » à « un programme qui vit dans un dossier avec `go.mod` » — maintenant ce
`go.mod` va organiser plus d'un fichier.

🔑 **Et ce n'est pas un caprice d'organisation.** Un test qui a besoin du réseau, d'un fichier sur disque ou
d'un serveur démarré pour s'exécuter est un test lent, fragile et que personne ne lance souvent. Séparer en
paquets est ce qui te permet d'écrire des tests qui s'exécutent en millisecondes, sans toucher à rien
d'externe — et c'est ce qui fait que tu les lances vraiment, à chaque changement, pas seulement quand tu y
penses.

---

## Les concepts

### 5.1 `go.mod`, et pourquoi ce projet n'a pas de `go.sum`

Tu as déjà utilisé `go mod init` dans la leçon 1 pour le programme `hola`. Le `go.mod` du `revisor` est tout
aussi simple :

```
module github.com/habil/revisor

go 1.27
```

Deux lignes : le nom du module (c'est ainsi qu'un autre programme l'importerait, si un jour tu publies l'un de
ses paquets) et la version minimale de Go dont il a besoin.

**Si tu cherches des tutoriels sur les modules sur internet, presque tous vont mentionner `go.sum` tout de
suite** — le fichier avec les empreintes cryptographiques de chaque dépendance externe, pour que personne ne
te glisse une version différente d'un paquet que tu utilises. Le `revisor` **n'a pas de `go.sum`**, et ce
n'est ni une erreur ni un oubli :

```bash
$ ls go.sum
ls: go.sum: No such file or directory
```

**Il n'existe pas parce que le `revisor` n'importe pas un seul paquet externe.** Regarde les `import` de
n'importe quel fichier du projet et tu ne trouveras que des paquets de la bibliothèque standard : `net/http`,
`encoding/json`, `context`, `sync`, `flag`, `os`, `time`, `strings`, `sort`. C'est une décision délibérée,
pas une limitation : la leçon 0 l'annonçait déjà —*« apprends la bibliothèque standard avant n'importe quel
framework »*— et le `revisor` est la preuve que la bibliothèque standard suffit pour un programme complet,
avec concurrence, HTTP, JSON et tests, sans ajouter une seule dépendance tierce. Si un jour tu en ajoutes une
(par exemple, un vrai client YAML), à ce moment-là `go get` créera `go.sum` pour toi, et les deux fichiers
—`go.mod` et `go.sum`— vont dans le dépôt.

Pour voir pour de vrai ce qui aurait changé, j'ai fait l'essai dans un projet à part (pas dans le `revisor`,
qui reste sans dépendances) : un vrai `go get` d'un petit paquet externe, `gopkg.in/yaml.v3`.

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

Voilà ce qui apparaît dès que tu ajoutes **une seule** dépendance externe : `go.mod` gagne une ligne
`require`, et `go.sum` naît avec les empreintes cryptographiques (les longs textes qui commencent par `h1:`)
de cette dépendance et des siennes propres (`check.v1` est une dépendance indirecte de `yaml.v3`, pas quelque
chose que tu as demandé). Ces empreintes sont ce qui fait que, si quelqu'un essayait de te glisser une version
différente du paquet avec le même nom et la même version, `go build` refuserait de compiler — c'est une
garantie d'intégrité, pas seulement un registre. Le `revisor` ne les a pas parce qu'il n'en a pas besoin :
zéro dépendance externe, zéro surface pour ce genre de risque.

### 5.2 Des paquets par responsabilité, pas par couche

Voici comment le `revisor` est organisé dans cette leçon :

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

**Par responsabilité, pas par couche.** Il n'y a pas de paquet `modelos` avec tous les structs du programme
ni de paquet `utilidades` avec des fonctions éparses : chaque paquet a une seule question à laquelle il sait
répondre. `servicio` sait ce qu'est un service et un état. `config` sait lire la configuration. `revisar`
sait interroger. `reporte` sait afficher. Si demain tu changes l'apparence du tableau, tu touches un fichier,
pas cinq.

🔑 **`internal/` est une règle du compilateur, pas une convention de bonnes manières.** Tout paquet qui vit
sous un dossier appelé `internal/` ne peut être importé que par du code qui se trouve **à l'intérieur du même
module**, à n'importe quel niveau au-dessus de ce `internal/`. Vérifie-le : si un autre module Go —n'importe
lequel, pas seulement un des tiens— essaie `import "github.com/habil/revisor/internal/servicio"`, le
compilateur refuse de compiler, avec un message explicite indiquant que ce paquet est interne. Ce n'est pas
une recommandation que tu peux ignorer sous la pression : c'est une restriction réelle, le même genre de
garantie que la majuscule te donne à l'intérieur d'un struct (leçon 3), mais au niveau d'un paquet entier.

⚠️ **Ce que nous n'avons PAS fait, exprès : un paquet `utils`, `helpers` ou `common`.** C'est l'anti-patron le
plus répandu dans les vrais projets : quelqu'un crée un dossier pour « les choses qui ne vont nulle part
ailleurs », et ce dossier grandit sans limite jusqu'à ce que personne ne sache ce qu'il y a dedans ni
pourquoi. Chaque fois que tu es tenté de mettre quelque chose dans un paquet de ce genre, demande-toi de
quelle **responsabilité** relève cette fonction, et place-la dans le paquet propriétaire de cette
responsabilité — ou, si elle ne va vraiment dans aucun, c'est le signe qu'il manque un nouveau concept à
nommer, pas qu'il manque un fourre-tout.

### 5.3 Le fichier de configuration, et ses vraies erreurs

Le `config` de cette leçon lit un format texte simple, une ligne par service :

```
nombre  url  [tiempo-limite]
```

```
# las lineas que empiezan con # se ignoran
catalogo   https://catalogo.interno.mx
pagos      https://pagos.interno.mx      500ms
```

C'est ce que fait la fonction `Interpretar`, qui **ne lit jamais de fichier** : elle reçoit les octets déjà
lus, et c'est pourquoi on peut la tester avec cent variantes de contenu sans créer un seul fichier temporaire
(`Cargar`, qui, elle, touche au disque, est une fine couche par-dessus qui se contente de lire le fichier et
de passer le contenu à `Interpretar`). Chaque ligne mal écrite produit une vraie erreur, avec le numéro de
ligne et ce qui était attendu — testé, pas supposé :

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

Remarque la dernière : **elle enveloppe l'erreur de `time.ParseDuration` avec `%w`** (leçon 3) au lieu
d'inventer son propre texte — ainsi, si un jour tu as besoin de distinguer par programme « durée invalide »
d'un autre type d'erreur avec `errors.As`, l'information originale est toujours là.

### 5.4 Les tests : sans bibliothèques, avec une table de cas

Voici à quoi ressemble un vrai test du paquet `servicio` (le fichier complet se trouve dans
`programas/revisor/internal/servicio/servicio_test.go`) :

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

**Il n'y a pas d'`assert`, ni d'`expect`, ni aucune bibliothèque d'assertions — et c'est exprès.** Go compare
des valeurs avec un `if` normal et rapporte avec `t.Errorf`, en écrivant toi-même ce que tu attendais et ce
que tu as obtenu. Au début, cela paraît plus verbeux qu'un `assert.Equal(t, esperado, obtenido)` d'autres
langages ; le gain, c'est que le message d'échec, c'est toi qui le contrôles, au lieu d'hériter du format
générique d'une bibliothèque, et qu'il n'y a rien à installer ni à apprendre en plus pour écrire le test le
plus simple.

**`t.Run` donne un nom à chaque cas de la table**, et cela compte quand quelque chose échoue : au lieu d'un
générique « `TestTimeoutEfectivo` a échoué », le rapport dit exactement lequel des trois cas c'était :

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

Cette sortie est réelle : je l'ai lancée sur le code de ce même projet avant d'écrire cette ligne.
**Ajouter un nouveau cas à la table, c'est ajouter une ligne au slice `casos`** — pas une nouvelle fonction,
pas répéter le corps du test. C'est la façon idiomatique de tester une fonction avec beaucoup d'entrées en
Go, et tu vas l'utiliser dans chaque paquet à partir de maintenant.

### 5.5 Tester les erreurs, pas seulement les réussites

Une table de cas sert aussi à tester que quelque chose **échoue comme il faut**, pas seulement que cela
fonctionne (version abrégée ici, avec 3 des 6 cas réels et un struct littéral positionnel au lieu de noms
de champs, pour que cela tienne ; la version complète se trouve dans `internal/config/config_test.go`) :

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

**`t.Fatalf` au lieu de `t.Errorf`** dans le premier `if` : si `Interpretar` n'a pas renvoyé d'erreur alors
qu'elle le devait, continuer à vérifier `err.Error()` à la ligne suivante provoquerait un panic (`err` serait
`nil`). `Fatalf` arrête ce test en particulier sur-le-champ ; `Errorf` laisse le test continuer à tourner et
accumuler d'autres échecs avant de rapporter. La règle pratique : utilise `Fatalf` quand continuer n'a pas de
sens sans ce que tu viens de vérifier, `Errorf` quand cela en a.

### 5.6 Les commandes que tu vas utiliser tout le temps

```bash
go test ./...                          # todo el proyecto
go test ./internal/config/... -v       # verboso, un paquete
go test ./... -run TestInterpretar     # solo las pruebas cuyo nombre haga match
go test ./... -race                    # detector de carreras (se explica a fondo en la lección 6)
go test ./... -cover                   # porcentaje de líneas ejercitadas por las pruebas
```

Lancées pour de vrai sur le `revisor`, le 30-sep-2026 :

```
$ go test ./internal/servicio/... ./internal/config/... -cover
ok  	github.com/habil/revisor/internal/servicio	0.195s	coverage: 92.3% of statements
ok  	github.com/habil/revisor/internal/config	0.192s	coverage: 97.4% of statements
```

### 5.7 Que faire d'un chiffre de couverture

**92,3 % et 97,4 % ne sont pas des objectifs, ce sont des points de départ pour une question : que sont ces
8 % et ces 3 % qui n'ont pas été exécutés, et est-ce que cela m'importe ?** Avec `go test -coverprofile`, tu
peux voir exactement quelles lignes sont restées intouchées :

```bash
go test ./internal/servicio/... -coverprofile=/tmp/cobertura.out
go tool cover -func=/tmp/cobertura.out
```

Dans le `revisor`, le trou de `servicio` est la branche de `Motivo()` qui construit le message quand le code
n'est ni 0 ni un succès avec un texte spécifique (`fmt.Sprintf("codigo %d", ...)` pour un 404 sans plus de
contexte) — une branche que les tests existants n'exercent pas avec ce code exact. **La bonne décision n'est
pas de poursuivre les 100 %** en remplissant chaque branche d'un test forcé qui n'apprend rien de nouveau :
c'est de regarder le trou, de décider s'il importe (ici, un peu : tu ajouterais un cas avec 404 dans la table)
et de le noter, au lieu de faire comme s'il n'existait pas.

🔴 **Et un vrai piège, mesuré dans ce même projet : la couverture par paquet peut sous-estimer une fonction
centrale sans te prévenir.** Lancé seul, le paquet `revisar` du `revisor` (que tu vas connaître à fond dans la
leçon 6) rapporte :

```
$ go test ./internal/revisar/... -cover
ok  	github.com/habil/revisor/internal/revisar	1.33s	coverage: 61.9% of statements
```

**61,9 % donne l'impression que plus d'un tiers de ce paquet ne s'exécute jamais dans aucun test — et c'est
faux.** Sa fonction la plus importante, `Revisar` (celle qui parle vraiment HTTP), est bel et bien testée :
simplement, le test qui l'exerce vraiment —`TestEjecutar_reportaOKyFalla`, avec un vrai serveur `httptest`,
dans la leçon 7— vit dans le paquet `cmd/revisor`, pas dans `internal/revisar`. `go test ./internal/revisar/...`
**ne compte que ce que les tests DE CE PAQUET exercent** ; il ne voit pas ce qu'un test d'un autre paquet
parcourt en chemin, même s'il passe par le même code. Avec `-coverpkg`, qui demande à Go de mesurer la
couverture d'un paquet en comptant **tous** les tests du projet, pas seulement les siens :

```
$ go test ./... -coverpkg=./... -coverprofile=/tmp/cov.out
$ go tool cover -func=/tmp/cov.out | grep 'revisar.go.*Revisar'
github.com/habil/revisor/internal/revisar/revisar.go:48:  Revisar    86.4%

$ go tool cover -func=/tmp/cov.out | tail -1
total:                                                    (statements)    81.5%
```

**86,4 % pour `Revisar`, 81,5 % pour le projet complet — pas 61,9 %.** La leçon n'est pas « ignore le chiffre
par paquet » : c'est qu'un chiffre de couverture répond toujours à une question implicite —couverture de
quoi, mesurée par rapport aux tests de où ?— et `go test ./paquete/... -cover` tait cette seconde moitié de
la question. Avant de décider que quelque chose « n'est pas testé » à cause d'un chiffre bas, lance
`-coverpkg=./...` sur tout le projet et compare.

### 5.8 Provoquer un échec, pour savoir que le test sert à quelque chose

Un test que tu n'as jamais vu échouer est un test dont tu ne sais pas s'il fonctionne — il peut être en train
de comparer deux choses qui sont toujours égales par accident. Casse-le exprès : modifie
`TimeoutEfectivo()` pour qu'elle renvoie toujours `s.Timeout`, sans le `if`, et lance le test :

```
$ go test ./internal/servicio/... -run TestTimeoutEfectivo -v
=== RUN   TestTimeoutEfectivo/cero_(valor_por_omision)
    servicio_test.go:64: TimeoutEfectivo() = 0s, quería 2s
--- FAIL: TestTimeoutEfectivo (0.00s)
    --- FAIL: TestTimeoutEfectivo/cero_(valor_por_omision) (0.00s)
FAIL
```

**Ce `FAIL`, avec le cas exact qui a cassé et la valeur obtenue face à celle attendue, est la preuve que le
test fonctionne.** Remets le `if` et confirme qu'il repasse au vert avant de continuer.

### 5.9 Concevoir pour pouvoir tester : séparer la logique des E/S

Remarque une chose que nous avons déjà mentionnée en passant dans la section 5.3 et qui mérite son propre
espace : `config` a **deux** fonctions, pas une.

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

`Interpretar` (la signature seulement ; le corps complet, avec toute la logique d'analyse, est dans la
section 5.3, avec ses vraies erreurs) :

<!-- verificar:fragmento -->
```go
func Interpretar(datos []byte, origen string) ([]servicio.Servicio, error) {
	// ... toda la lógica de parseo, sin tocar el disco ...
}
```

**S'il n'y avait qu'une seule fonction `Cargar(ruta string)` qui lise le fichier et analyse tout d'un coup**,
chaque test de « ligne mal écrite », « nom répété » ou « délai d'expiration invalide » devrait commencer par
créer un fichier temporaire sur le disque avec `os.CreateTemp`, y écrire le contenu du cas, lui passer le
chemin, et le supprimer à la fin. Cela fonctionne, mais c'est lent (cela touche le vrai système de fichiers)
et cela salit le test avec du code qui n'a rien à voir avec ce qu'on teste réellement : **si ta configuration
est bien ou mal interprétée**, pas si tu sais créer des fichiers temporaires.

En séparant **la partie qui décide** (`Interpretar`, une fonction pure : mêmes octets en entrée, même
résultat toujours, sans accès à rien d'externe) de **la partie qui obtient les octets** (`Cargar`, la seule
qui touche au disque), les neuf tests de la section 5.5 s'exécutent en microsecondes et sans créer un seul
fichier. `Cargar` elle-même n'a presque pas besoin de tests propres : seulement vérifier qu'elle renvoie une
erreur si le fichier n'existe pas (exercice 6), parce que toute la logique intéressante est déjà dans
`Interpretar` et déjà testée.

🔑 **La règle générale, utile bien au-delà de ce projet :** quand une fonction est difficile à tester, c'est
presque toujours parce qu'elle mélange « décider quelque chose » et « toucher au monde extérieur » (un
fichier, le réseau, l'horloge). Les séparer n'est pas une règle de style : c'est ce qui détermine si tu vas
pouvoir écrire le test en trois lignes ou en vingt.

### 5.10 Benchmarks : mesurer, pas deviner

En plus de `Test...`, Go reconnaît des fonctions `Benchmark...` qui mesurent combien de temps prend ton code,
pas s'il est correct :

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

`b.N`, ce n'est pas toi qui le choisis : Go exécute la boucle avec des valeurs de `N` de plus en plus grandes
jusqu'à ce que la mesure soit stable, et rapporte le temps par opération. Lancé pour de vrai sur le `revisor`
(Apple M5, 200 000 répétitions) :

```
$ go test ./internal/config/... -bench=. -run '^$'
goos: darwin
goarch: arm64
pkg: github.com/habil/revisor/internal/config
cpu: Apple M5
BenchmarkInterpretar-10    	  200000	       432.5 ns/op
PASS
```

**`-run '^$'`** dit à `go test` de ne lancer aucun test normal (une expression régulière qui ne correspond à
aucun nom), pour que le rapport du benchmark ne se mélange pas avec celui des tests. Le chiffre —432,5
nanosecondes par appel, sur cette machine, à ce moment— n'est pas à mémoriser : il sert à le comparer **avec
lui-même** après un changement. Si demain tu réécris `Interpretar` et que le benchmark monte à 4 000 ns/op, tu
as un signal objectif que quelque chose est devenu plus lent, sans avoir besoin d'avoir un avis là-dessus.

### 5.11 `go vet` : celui qui trouve ce qui compile mais qui est faux

`go test` te dit si ta logique fait ce que tu attendais. **`go vet` te dit si ton code contient une erreur que
le compilateur n'attrape pas parce que, techniquement, il est valide** — mais ce n'est presque sûrement pas ce
que tu voulais écrire. Le cas le plus courant est un verbe de formatage (leçon 2) qui ne correspond pas au
type de l'argument :

<!-- verificar:ejemplo:ejemplos/05-vet-printf -->
```go
puerto := 443
fmt.Printf("servicio %s en el puerto %s\n", nombre, puerto) // %s para un int
```

Ceci **compile** et **s'exécute**, sans panic ni erreur — et produit une sortie cassée (programme complet
dans `programas/revisor/ejemplos/05-vet-printf/main.go`) :

```
$ go run ./ejemplos/05-vet-printf/
servicio catalogo en el puerto %!s(int=443)
```

Et `go vet`, sur ce même fichier, le détecte bien :

```
$ go vet ./ejemplos/05-vet-printf/
ejemplos/05-vet-printf/main.go:11:39: fmt.Printf format %s has arg puerto of wrong type int
```

`go vet` le détecte bien, parce qu'il analyse la chaîne de format par rapport aux types réels des arguments,
une étape que le compilateur de Go ne fait pas à lui seul. Lance `go vet ./...` avec `go test ./...` comme une
routine : la plupart des éditeurs avec l'extension Go (leçon 1) le font déjà pour toi pendant que tu écris,
en soulignant le problème avant que tu n'arrives à exécuter quoi que ce soit.

### 5.12 Tests aux limites : là où vivent vraiment les bugs

`Estado.OK()` décide qu'un code est bon s'il tombe **entre 200 et 299**. Il est tentant de le tester avec un
cas « évident » (200) et un cas « évident » d'échec (500) et de s'en contenter. La vraie table du `revisor`
teste aussi les deux valeurs qui sont **juste à la limite de l'intervalle** :

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

**Pourquoi est-ce important, si 199 et 300 ne sont « évidemment » pas OK ?** Parce qu'une erreur d'un seul
caractère dans la condition —`>=` au lieu de `>`, ou `<=` au lieu de `<`— est exactement le type de bug qu'un
cas « évident » ne détecte jamais, et qu'un cas aux limites détecte toujours. Si quelqu'un changeait `OK()`
en `e.Codigo >= 200 && e.Codigo <= 300` (en incluant le 300 par erreur), les cas 200/299/404/500 continueraient
de passer tout aussi bien — **seul le cas de 300 le trahirait.** C'est la raison profonde derrière « teste les
limites, pas seulement le centre » : les limites sont l'endroit où se cachent les erreurs de comparaison, et
elles sont invisibles pour tout test qui n'utilise que des valeurs bien à l'intérieur ou bien à l'extérieur de
l'intervalle.

---

## L'erreur que tu vas voir

| Le symptôme | Message littéral | Ce qui se passe et quoi faire |
|---|---|---|
| Import d'un paquet `internal` étranger | `use of internal package github.com/habil/revisor/internal/servicio not allowed` | Le compilateur empêche d'importer quelque chose sous `internal/` depuis l'extérieur du module. Ce n'est pas une permission que tu peux accorder : il faut exposer le type depuis un paquet public si c'est vraiment nécessaire |
| Service en double dans la configuration | `prueba.txt:2: el nombre "catalogo" ya estaba en la linea 1` | Deux services avec le même nom ; corrige le fichier de configuration |
| URL sans schéma | `prueba.txt:1: la URL "catalogo.interno.mx" debe empezar con http:// o https://` | Il manque `http://` ou `https://` au début de l'URL |
| Délai d'expiration mal écrit | `prueba.txt:1: tiempo limite "nombas" invalido (usa 500ms, 2s, 1m): time: invalid duration "nombas"` | Le format de durée de Go n'est pas libre : utilise un nombre suivi d'une unité (`ms`, `s`, `m`, `h`) |
| Un test que tu as cassé exprès | `servicio_test.go:64: TimeoutEfectivo() = 0s, quería 2s` | Voir la section 5.8 : voilà à quoi ressemble un `t.Errorf` qui signale exactement quel cas a échoué et pourquoi |
| Zéro test dans un paquet | `?   	github.com/habil/revisor/cmd/servidor-demo	[no test files]` | Ce n'est pas un échec : `go test` prévient explicitement quand un paquet n'a aucun fichier `_test.go`, au lieu de faire semblant d'avoir exécuté quelque chose |

**Et le plus important à apprendre à lire**, parce que ce n'est pas une erreur mais l'absence d'une erreur :

```
$ go test ./...
ok  	github.com/habil/revisor/internal/vacio	0.001s
```

Si `internal/vacio` n'avait **aucune** fonction `Test...`, cette ligne aurait exactement la même allure :
`ok`, en vert, sans aucune marque indiquant que rien n'a été exécuté. La seule façon de distinguer « tout est
passé, pour de vrai » de « il n'y avait rien à exécuter » est de regarder le décompte avec `-v` (qui, lui,
affiche chaque `RUN`) ou, mieux, de ne jamais faire confiance à un paquet qui n'a aucun fichier `_test.go` —
cela, `go test` le dit bien, comme dans la ligne du tableau ci-dessus.

---

## Ce qu'on fait mal

- **Créer un paquet `utils`, `helpers` ou `common`.** Tu l'as déjà vu dans la section 5.2 : c'est l'anti-patron
  le plus répandu et celui qui perd le plus vite son objectif. Nomme la responsabilité, pas le fait que cela
  « ne va nulle part ailleurs ».
- **Confondre « ça a compilé » avec « les tests sont passés ».** `go build` vérifie que le code est valide ; il
  n'exécute pas un seul test. Ce sont deux commandes différentes avec deux questions différentes.
- **Lire le code de sortie de `go test` au lieu du décompte.** Un paquet sans fichiers de test se termine avec
  le code 0 (succès), exactement comme un paquet avec 50 tests qui sont bel et bien passés. Le code de sortie
  répond à « quelque chose a-t-il échoué ? », pas à « quelque chose a-t-il été testé ? » — pour la seconde
  question, il faut lire la sortie, pas seulement le `$?`.
- **Poursuivre 100 % de couverture comme si c'était l'objectif.** Une couverture élevée avec des assertions
  faibles (vérifier que quelque chose ne plante pas, sans vérifier ce que cela renvoie) donne un joli chiffre
  et un test qui ne détecte presque rien. Le chiffre est un guide pour savoir où regarder, pas un but en soi —
  la section 5.7 le montre avec un vrai trou du `revisor` lui-même.
- **Écrire un test et ne jamais le voir échouer.** Si tu n'as jamais cassé le code exprès pour confirmer que
  le test passe au rouge, tu ne sais pas si ce test teste quelque chose. La section 5.8 l'a fait avec des
  données réelles du projet : fais-le toi aussi avec au moins un de tes tests avant de le considérer comme bon.

---

## Exercices

1. Clone la structure de paquets de la section 5.2 (`internal/servicio`, `internal/config`) pour ta propre
   copie du `revisor`, en déplaçant le code que tu avais déjà des leçons 2 à 4.
2. Écris la table de cas de `TestEtiqueta` pour la méthode `Etiqueta()` de `Servicio`, avec au moins un cas
   de nom normal et un de struct vide.
3. Ajoute un cas à `TestInterpretar_casosDeError` pour une ligne avec **quatre** champs (plus que les trois
   que le format permet). Vérifie le message exact par rapport au code de `config.go`.
4. Lance `go test ./... -cover` sur ta copie et note le pourcentage de chaque paquet. Choisis **un** trou de
   couverture et décide, par écrit dans ton journal de bord, s'il t'importe de le combler et pourquoi.
5. (Comme dans la section 5.8) Casse exprès une fonction que tu as déjà testée, lance le test, lis le `FAIL`
   complet, et répare-la. Colle les deux résultats —le rouge et le vert— dans ton journal de bord.
6. (Un peu plus difficile) Écris `TestCargar_archivoInexistente`, qui confirme que `Cargar` (et non
   `Interpretar`) renvoie une erreur quand le chemin n'existe pas. Indice : tu n'as besoin de créer aucun
   fichier pour ce test, seulement de passer un chemin dont tu sais qu'il n'existe pas.

### Corrigés

1 et 2 n'ont pas de corrigé de référence unique : cela dépend de la façon dont tu avais organisé ton propre
programme de la leçon 4. Compare ton résultat avec le vrai code de `programas/revisor/internal/servicio/servicio.go` du
projet de ce cours.

3. Avec `"catalogo https://a.mx 500ms extra\n"`, le message attendu est
   `prueba.txt:1: esperaba «nombre url [tiempo]», hay 4 campo(s)` — le même chemin de code qui gère déjà
   « il manque des champs », parce que `len(campos) > 3` couvre les deux cas avec une seule vérification.

4. Il n'y a pas de réponse unique : ce qui compte, c'est que la décision soit écrite avec sa raison, pas le
   pourcentage en lui-même.

5. Voir la section 5.8 en entier : le modèle est toujours « modifie le code, lance le test, lis le `FAIL`
   avec le cas exact, répare, relance ».

6. Ainsi :

   <!-- verificar:extracto:internal/config/config_test.go -->
   ```go
   func TestCargar_archivoInexistente(t *testing.T) {
   	_, err := Cargar("/ruta/que/no/existe.txt")
   	if err == nil {
   		t.Fatal("Cargar() no devolvió error con una ruta inexistente")
   	}
   }
   ```

   Ce test existe bel et bien dans le vrai projet (`config_test.go`) et il passe parce que `os.ReadFile`, à
   l'intérieur de `Cargar`, renvoie une erreur du système d'exploitation que `Cargar` enveloppe avec `%w`
   avant de la propager.

---

## Comment savoir que j'y suis arrivé

- [ ] Mon `revisor` est organisé en paquets sous `internal/`, chacun avec une seule responsabilité.
- [ ] `go test ./...` s'exécute et **le décompte, pas seulement la couleur**, me dit combien de tests ont été
  exécutés.
- [ ] J'ai écrit au moins une table de cas avec `t.Run`, et je sais lire quel cas a échoué quand quelque chose
  casse.
- [ ] J'ai lancé `go test -cover` et je peux dire quel pourcentage est sorti et ce qu'il y a dans le trou.
- [ ] J'ai cassé une fonction exprès, j'ai vu le `FAIL` exact, et je l'ai réparée — c'est dans mon journal de
  bord.
- [ ] Je sais expliquer pourquoi le `revisor` n'a pas de `go.sum` et ce qui le générerait si un jour il en
  avait besoin.
- [ ] Je sais pourquoi je ne devrais pas créer de paquet `utils`.

---

## Pour aller plus loin

1. [Writing tests](https://go.dev/doc/tutorial/add-a-test) — le tutoriel officiel sur les tests, table de cas
   comprise.
2. [Go Wiki: TableDrivenTests](https://go.dev/wiki/TableDrivenTests) — la référence canonique du modèle
   utilisé dans toute cette leçon.
3. [Package internal](https://go.dev/doc/go1.4#internalpackages) — l'annonce originale de la règle de
   `internal/`, tirée directement des notes de version de Go.
4. [Go Blog: The Cover Story](https://go.dev/blog/cover) — comment fonctionne `go test -cover` en interne, et
   pourquoi un chiffre élevé ne signifie pas toujours de bons tests.
5. [Livrer sans crainte sur votre propre infrastructure](https://www.habil.mx/fr/blog/cicd-devsecops-infrastructure-propre/) — article sur un chemin vers la production où les tests, parmi d'autres portes, peuvent arrêter une livraison.

### Termes de cette leçon

| Terme | Ce qu'il signifie |
|---|---|
| `go.sum` | fichier avec les empreintes cryptographiques des dépendances externes ; le `revisor` ne l'a pas parce qu'il n'en utilise aucune |
| `internal/` | dossier spécial que le compilateur de Go empêche d'importer depuis l'extérieur du module |
| table de cas | modèle de test où une liste d'entrées et de sorties attendues est parcourue avec un seul corps de test |
| `t.Run` | exécute un sous-cas avec son propre nom, pour que le rapport dise exactement lequel a échoué |
| couverture | pourcentage des lignes du code que les tests ont exercées pendant leur exécution |

---

**Précédent :** [Leçon 4 — Collections](04-colecciones.md) ·
**Suivant :** [Leçon 6 — Concurrence](06-concurrencia.md)
