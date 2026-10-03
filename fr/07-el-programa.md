# Leçon 7 — Le programme terminé

**Durée :** 90-120 minutes, ou 2 séances de 60. C'est la dernière leçon du `revisor` : à la fin, tu as un
binaire qui tourne vraiment, se configure depuis la ligne de commande, parle HTTP avec de vrais délais
d'expiration, produit deux formats de sortie, et se compile pour une autre plateforme sans quitter ta
machine.

**À la fin, tu seras capable de :**

- Écrire un client HTTP avec un délai d'expiration explicite, et expliquer pourquoi `http.DefaultClient` est
  dangereux en production.
- Séparer la logique d'un programme (`ejecutar`) de son point d'entrée (`main`), pour pouvoir la tester sans
  toucher à `os.Exit`.
- Définir des options de ligne de commande avec le paquet standard `flag`, sans aucune bibliothèque externe.
- Sérialiser une structure en JSON avec `encoding/json`, et expliquer quels champs se perdent et pourquoi.
- Exécuter le `revisor` de bout en bout contre un vrai serveur (un serveur de test, fait par toi) et lire sa
  sortie.
- Compiler le même code pour une autre architecture et un autre système d'exploitation sans quitter ta
  machine, et expliquer pourquoi c'est possible.

---

## Pourquoi c'est important

Tu as déjà les quatre pièces du `revisor` construites séparément : `servicio` définit le vocabulaire (leçon
2-3), `config` lit la configuration (leçon 5), `revisar.Todos` interroge tout en même temps (leçon 6). La seule
chose qui manque, c'est ce qui transforme ces pièces en **un programme que quelqu'un d'autre peut utiliser
sans lire le code source** : qu'il parle vraiment HTTP (jusqu'ici, nous ne l'avons testé qu'avec `Falso`),
qu'il reçoive des options depuis le terminal au lieu de valeurs fixes dans le code, qu'il produise un format
qu'un autre programme peut consommer, et qu'il puisse être compilé et distribué sous la forme d'un seul
fichier.

C'est la leçon où le `revisor` cesse d'être « du code qui fonctionne si c'est moi qui le lance, sur ma machine,
avec mes données de test » et devient un binaire que tu peux copier à quelqu'un d'autre, avec la certitude
qu'il fera exactement ce que la ligne de commande lui demande.

---

## Les concepts

### 7.1 Le client HTTP : pourquoi jamais `http.DefaultClient`

Voici le vrai `Revisor`, celui qui touche vraiment au réseau (`programas/revisor/internal/revisar/revisar.go`) :

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

🔴 **`http.DefaultClient` —celui que tu utiliserais si tu écrivais directement `http.Get(url)`— n'a aucun délai
d'expiration.** Si le serveur de l'autre côté accepte la connexion et ne répond jamais rien, cet appel reste en
attente **pour toujours**, sans erreur, sans panic, juste un programme qui un jour cesse d'avancer sans que
personne ne sache pourquoi. C'est l'une des erreurs les plus coûteuses de Go en production, précisément parce
qu'elle ne donne aucun symptôme avant d'être déjà en train de se produire.

Le `revisor` se protège deux fois, pas une : le `http.Client` propre a un délai d'expiration général (la limite
de tout le rapport, passée à `NuevoHTTP`), et en plus chaque requête individuelle tourne sous le
`context.WithTimeout` par service que `Todos` a mis en place dans la leçon 6 — ce second délai, plus court, est
celui qui commande en pratique. Le délai d'expiration du client est le filet de sécurité de dernier recours, pas
le mécanisme principal.

⚠️ **`defer resp.Body.Close()` vient après la vérification de l'erreur, jamais avant.** Si `err != nil`,
`resp` est `nil`, et appeler une méthode sur un pointeur nul fait un panic. Et **`io.Copy(io.Discard,
resp.Body)` avant de fermer** n'est pas décoratif : sans vider le corps de la réponse, la connexion TCP
sous-jacente ne peut pas être réutilisée pour la requête suivante vers le même serveur, et le programme finit
par ouvrir une nouvelle connexion à chaque fois au lieu de réutiliser celles qu'il a déjà.

**`traducir` remplace l'erreur brute de `net/http` par une erreur lisible dans un tableau :**

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

L'erreur que donne `net/http` d'origine contient l'URL complète et le mot `Get`, qui dans un tableau de rapport
ne font qu'occuper de la place et n'apprennent rien de nouveau au lecteur — c'est pourquoi on la traduit avant
de l'afficher.

### 7.2 `main` ne se teste pas ; `ejecutar`, si

Le modèle qui distingue le `revisor` d'un programme d'un seul fichier :

<!-- verificar:extracto:cmd/revisor/main.go -->
```go
func main() {
	os.Exit(ejecutar(os.Args[1:], os.Stdout, os.Stderr))
}
```

`ejecutar` est l'endroit où vit la vraie logique —les sections 7.3, 7.4 et 7.6 la montrent par parties— ; sa
signature :

<!-- verificar:fragmento -->
```go
func ejecutar(args []string, salida, errores *os.File) int {
	// toda la lógica real vive aquí
}
```

**`go test` ne peut pas tester une fonction qui appelle `os.Exit`**, parce que `os.Exit` termine le processus
entier immédiatement — y compris le processus de test lui-même, qui n'arriverait jamais à rapporter le
résultat. En séparant « décider quel code de sortie correspond » (`ejecutar`, qui **renvoie** un `int`) de
« sortir vraiment avec ce code » (`main`, la seule ligne qui appelle `os.Exit`), toute la logique peut être
testée en appelant directement `ejecutar`, avec des arguments de test, et en vérifiant le nombre qu'elle
renvoie — exactement ce que fait `cmd/revisor/main_test.go` :

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

C'est la vraie sortie de `--help` (générée automatiquement par `flag`, section 7.3) capturée au milieu d'un
test, sans ouvrir de terminal ni exécuter le binaire compilé.

### 7.3 Options de ligne de commande, sans bibliothèques

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

**`flag.NewFlagSet` au lieu du paquet `flag` au niveau global** (que tu utiliserais directement avec
`flag.String(...)`) est ce qui permet d'avoir une fonction `ejecutar(args []string, ...)` qui reçoit ses
arguments en paramètre, au lieu de les lire toujours dans `os.Args` — un autre détail qui existe précisément
pour pouvoir la tester, comme dans la section 7.2.

`flag` est fourni avec la bibliothèque standard et suffit pour un programme de cette taille : quatre options,
des types de base (`string`, `int`, `time.Duration`), une aide automatique. Si un jour le `revisor` avait besoin
de sous-commandes (`revisor check`, `revisor list`, chacune avec ses propres options), alors il vaudrait la
peine de regarder une bibliothèque comme `cobra` — mais pas avant d'en avoir vraiment besoin.

### 7.4 Les deux formats de sortie : tableau et JSON

Voici `Tabla`, la fonction qui produit la sortie lisible par des humains que tu as vue dans toute la leçon
(`programas/revisor/internal/reporte/tabla.go`) :

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

**La largeur de la première colonne n'est pas fixée dans le code — elle est calculée.** La boucle ci-dessus
parcourt tous les états une fois avant d'afficher quoi que ce soit, et retient le nom le plus long (en partant
de la largeur du mot « SERVICIO » lui-même, au cas où un nom de service serait plus court que cela). Ce nombre
entre comme le `*` dans `%-*s` : un verbe de formatage à largeur **variable**, où la largeur elle-même est un
argument de plus de `Fprintf`, pas un nombre écrit à la main. C'est la raison pour laquelle le vrai tableau de
cette leçon a des colonnes droites, que les noms des services soient courts (`sano`) ou longs (`inventario`) :
avec une largeur fixe, l'un des deux cas aurait toujours l'air de travers.

`ordenarPorNombre` (les sections 6.3 et 6.7 ont déjà expliqué pourquoi l'ordre d'arrivée n'est pas fiable :
plusieurs goroutines livrent leurs résultats par un canal, et c'est celle qui répond la première qui gagne) fait
une copie du slice et la trie par ordre alphabétique avant d'afficher — sans cela, la même requête produirait un
tableau dans un ordre différent à chaque exécution du programme, même si les données étaient identiques.

`etiquetaEstado`, `formatoTiempo` et `detalle` sont les trois petites fonctions qui décident quel mot va dans
chaque colonne (`OK`/`LENTO`/`FALLA`, le temps arrondi à la milliseconde ou un tiret s'il n'y a jamais eu de
réponse, et le code HTTP ou le motif de l'échec) — ce sont celles que tu as déjà vues à l'œuvre dans le vrai
tableau de la section 7.5, maintenant avec le code qui les produit.

🔑 **Pourquoi on n'a pas utilisé `text/tabwriter`, qui est l'outil que la bibliothèque standard propose
justement pour aligner des colonnes :** pour un tableau de quatre colonnes dont une seule a une largeur variable
(le nom du service ; les trois autres sont courtes et prévisibles), calculer la largeur à la main est plus simple
à lire que d'introduire un `tabwriter.Writer` avec ses propres `Flush()` et ses séparateurs par tabulation. Avec
un tableau comportant davantage de colonnes variables, `tabwriter` serait bien le bon outil — il vaut la peine
de le connaître (il est dans la section « Pour aller plus loin » de cette leçon) même si le `revisor` n'en a pas
besoin.

Une fois le tableau compris, l'autre format de sortie — JSON — est l'autre moitié de cette section :

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

⚠️ **`servicio.Estado` n'est pas sérialisé directement — exprès.** `Estado` contient un champ `Err error`, et
`error` est une interface : `encoding/json` ne sait pas la convertir en texte tout seul (essayer produit un
objet vide `{}`, pas une erreur de compilation, ce qui est pire : cela échoue en silence). Le `revisor` résout
cela avec un type intermédiaire, `lineaJSON`, qui ne vit qu'à l'intérieur du paquet `reporte` : il convertit
l'erreur en son texte (`e.Motivo()`) avant d'encoder, et au passage découple le format public (ce que
d'autres programmes vont lire) du modèle interne (`Estado`) — si demain `Estado` gagne un nouveau champ, le
JSON qui circule déjà ne change pas simplement parce que quelque chose d'interne a changé.

`omitempty` sur `Codigo` et `Detalle` retire ces champs du JSON quand ils valent zéro ou une chaîne vide —
ainsi, un service en bonne santé ne traîne pas un `"detalle": ""` dénué de sens. Encodé pour de vrai :

<!-- verificar:extracto:internal/reporte/json.go -->
```go
codificador := json.NewEncoder(w)
codificador.SetIndent("", "  ")
return codificador.Encode(lineas)
```

### 7.5 De bout en bout : le `revisor` contre un vrai serveur

Pour tester tout le programme ensemble —et pas chaque pièce séparément— j'ai construit `cmd/servidor-demo` :
un serveur HTTP minimal avec quatre routes qui se comportent comme les quatre cas que le `revisor` doit savoir
rapporter.

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

Avec ce serveur démarré sur `:8091` et un fichier de configuration pointant vers ses quatre routes, j'ai lancé
le vrai binaire :

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

**Chaque ligne de ce tableau correspond exactement au comportement programmé dans le serveur de test :**
`sano` répond vite et avec 200 ; `lento` met vraiment 1,5 seconde et c'est pourquoi il apparaît `LENTO` ;
`malo` répond instantanément mais avec 500 ; `colgado` ne répond jamais, et son délai d'expiration individuel de
200ms —déclaré dans le fichier de configuration, colonne trois— a coupé l'attente à 202ms, presque exactement.
**Le code de sortie, 1, est réel** : au moins un service n'était pas `OK`, donc `ejecutar` renvoie 1 au lieu de
0 (section 7.6) — le même binaire, utilisé depuis un script, peut dire à celui qui l'invoque s'il y a eu des
problèmes sans que personne n'ait à lire le tableau.

Et en JSON, le même rapport :

```
$ /tmp/revisor --config /tmp/servicios-demo.txt --formato json
[
  {"servicio": "colgado", "ok": false, "tiempo_ms": 205, "detalle": "se acabo el tiempo de espera"},
  {"servicio": "lento", "ok": true, "codigo": 200, "tiempo_ms": 1503},
  {"servicio": "malo", "ok": false, "codigo": 500, "tiempo_ms": 1, "detalle": "codigo 500"},
  {"servicio": "sano", "ok": true, "codigo": 200, "tiempo_ms": 1}
]
```

(Ici sans l'indentation de `SetIndent` pour que cela tienne sur une ligne par service ; le vrai programme la
produit avec des sauts de ligne, comme tu l'as vu dans la section 7.4.)

### 7.6 Le code de sortie, avec une signification

<!-- verificar:extracto:cmd/revisor/main.go -->
```go
for _, e := range estados {
	if !e.OK() {
		return 1 // al menos un servicio falló: el código de salida lo refleja
	}
}
return 0
```

Trois codes possibles, et chacun répond à une question différente pour celui qui invoque le programme depuis
un script ou un pipeline : **0** — tout s'est bien passé ; **1** — le programme a tourné jusqu'au bout, mais au
moins un service n'était pas en bonne santé ; **2** — le programme n'a même pas pu démarrer (il manque
`--config`, le fichier n'existe pas, le format de configuration est mal écrit). C'est la différence entre
« j'ai fait le travail et j'ai trouvé des problèmes » et « je n'ai même pas pu commencer à travailler » — et un
pipeline d'intégration continue peut réagir différemment à chacun.

### 7.7 Compilation croisée : le même code, une autre plateforme

```bash
GOOS=linux GOARCH=amd64 go build -o revisor-linux-amd64 ./cmd/revisor
```

Confirmé pour de vrai, en compilant depuis ce Mac (arm64) vers Linux x86-64 :

```
$ GOOS=linux GOARCH=amd64 go build -o revisor-linux-amd64 ./cmd/revisor
$ file revisor-linux-amd64
revisor-linux-amd64: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, ...
```

**Sans Docker, sans machine virtuelle, sans rien installer de plus : deux variables d'environnement
suffisent.** C'est possible parce que le compilateur de Go embarque, d'origine, le code nécessaire pour générer
des binaires pour n'importe quelle combinaison de système d'exploitation et d'architecture qu'il prend en charge
— il n'a pas besoin de tourner sur le système cible pour compiler vers lui, contrairement à la façon dont
fonctionnent beaucoup d'autres langages compilés.

**La taille du binaire, mesurée, avec et sans symboles de débogage :**

```
$ go build -o revisor-normal ./cmd/revisor
$ wc -c revisor-normal
9464130 revisor-normal

$ go build -ldflags="-s -w" -o revisor-chico ./cmd/revisor
$ wc -c revisor-chico
6386194 revisor-chico
```

**De 9,46 Mo à 6,39 Mo, une réduction réelle de 32 %**, en retirant la table des symboles (`-s`) et les
informations de débogage DWARF (`-w`) que le binaire normal inclut pour qu'un débogueur puisse l'inspecter.
Pour un binaire que tu vas distribuer en production et que tu ne vas pas déboguer sur place, ces données ne
servent à rien et ne font qu'occuper de la place — pour un binaire que tu es en train de développer activement,
garde-les.

### 7.8 Tester du code qui fait du HTTP, sans vrai réseau : `httptest`

Jusqu'à la leçon 6, les tests de concurrence utilisaient `Falso` (un `Revisor` qui ne touche jamais au réseau).
Pour tester le programme **complet** —options, lecture de la configuration, le vrai client HTTP, le rapport—
sans dépendre d'un service externe ni de quelque chose qui écoute sur un port fixe, `net/http/httptest` démarre
un vrai serveur, sur un port que le système d'exploitation attribue tout seul, à l'intérieur du processus de
test lui-même :

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

Exécution réelle :

```
$ go test ./cmd/revisor/... -run TestEjecutar_reportaOKyFalla -v
=== RUN   TestEjecutar_reportaOKyFalla
--- PASS: TestEjecutar_reportaOKyFalla (0.00s)
```

**`servidor.URL` est une vraie URL** (`http://127.0.0.1:PUERTO`, avec un port libre choisi par le système), donc
le `Revisor` HTTP de la section 7.1 l'interroge exactement comme il interrogerait n'importe quel service de
production — la seule différence, c'est que le « service » est un `http.HandlerFunc` de quinze lignes qui vit
dans le même processus que le test. C'est ce qui permet de tester de bout en bout —options, configuration,
client HTTP, rapport, code de sortie— en millisecondes, sans ouvrir aucun port fixe qui pourrait entrer en
conflit avec autre chose qui tourne sur la machine, et sans rien laisser tourner une fois le test terminé
(`defer servidor.Close()` s'en charge).

**`t.TempDir()`** crée un dossier temporaire que Go supprime tout seul à la fin du test (contrairement à un
simple `os.CreateTemp("/tmp", ...)`, qui laisserait le fichier là pour toujours si personne ne le supprimait à
la main) — la combinaison de `httptest` pour le réseau et de `t.TempDir()` pour le disque est ce qui permet à
`TestEjecutar_reportaOKyFalla` de ne laisser aucune trace sur le système après son exécution.

### 7.9 Pourquoi deux limites de temps, et pas une

Le `revisor` accepte `--limite` (le temps maximal pour **tout** le rapport) et chaque ligne de configuration
peut déclarer son propre délai d'expiration **par service**. Ce n'est pas redondant : ils répondent à deux
questions différentes. `--limite` répond à *« combien de temps, au maximum, suis-je prêt à attendre le rapport
complet ? »* — utile si le `revisor` tourne à l'intérieur d'un autre processus avec sa propre échéance (un
contrôle de santé qu'un orchestrateur attend à intervalles réguliers, par exemple). Le temps par service répond
à *« combien est-il raisonnable d'attendre CE service en particulier ? »* — un service interne à faible
latence et un autre qui passe par un fournisseur externe ne devraient pas partager la même limite.

La démonstration de la section 7.5 le confirme avec des données réelles : `colgado` avait une limite propre de
200ms déclarée dans le fichier de configuration, et il a été coupé à 202ms — bien avant que la `--limite`
générale (10 secondes par défaut) n'ait l'occasion d'intervenir. La plus courte des deux limites est toujours
celle qui commande, et c'est exactement ce que tu veux : qu'un seul service lent ne consomme pas tout le budget
de temps du rapport complet.

### 7.10 Le binaire « n'a besoin de rien » — sauf d'une chose : les certificats

La leçon 1 a démontré qu'un binaire Go tourne sur une machine Linux sans Go installé, grâce à la liaison
statique. Il est tentant d'étendre cette idée à « il n'a besoin de rien du système, point » — et je l'ai vérifié
en la poussant à l'extrême : j'ai mis le binaire du `revisor` dans une image Docker construite `FROM scratch`,
la base la plus vide qui existe (elle n'a même pas de coquille de système d'exploitation, seulement le
binaire).

```dockerfile
FROM scratch
COPY revisor /revisor
COPY servicios.txt /servicios.txt
ENTRYPOINT ["/revisor", "--config", "/servicios.txt"]
```

Avec un vrai service pointant vers `https://example.com` :

```
$ docker build -t revisor-demo:scratch .
$ docker run --rm revisor-demo:scratch --formato tabla
SERVICIO  ESTADO    TIEMPO  DETALLE
ejemplo   FALLA      176ms  no responde: Get "https://example.com": tls: failed to verify certificate: x509: certificate signed by unknown authority
```

**Échec — et pas à cause du `revisor`, mais à cause d'une chose qu'aucune image `scratch` n'apporte : la liste
des autorités de certification.** Pour valider un certificat TLS (le cadenas du `https://`), le système
d'exploitation fournit normalement une liste de ceux qui ont le droit de signer des certificats valides —
typiquement dans `/etc/ssl/certs/`. Une image `scratch` n'a même pas ce dossier, donc Go ne peut vérifier aucun
certificat et refuse de continuer, à juste titre : continuer sans vérifier reviendrait à accepter n'importe quel
certificat, valide ou faux.

**La solution, vérifiée, est de copier uniquement ce fichier** depuis une image qui l'a, sans entraîner le reste
du système d'exploitation :

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

**14,7 Mo au total, contre 14,4 Mo sans les certificats — moins de 300 Ko pour résoudre le problème
correctement.** Le binaire est bel et bien autosuffisant pour tout ce qui est *code Go* ; ce qu'il n'apporte
jamais tout seul, c'est la confiance sur les autorités légitimes, parce que c'est une information du monde
(quelles autorités de certification existent et sont toujours valides), pas quelque chose que le compilateur
peut décider pour toi.

---

## L'erreur que tu vas voir

| Le symptôme | Message / preuve littérale | Ce qui se passe et quoi faire |
|---|---|---|
| Il manque l'option obligatoire | `revisor: falta --config` suivi de l'usage automatique de `flag` | `ejecutar` vérifie explicitement que `--config` n'est pas vide avant de faire quoi que ce soit d'autre, et se termine avec le code 2 |
| Format de sortie invalide | `revisor: --formato debe ser "tabla" o "json", no "xml"` | Seuls deux formats existent ; toute autre valeur est rejetée avant de tenter quoi que ce soit |
| Service qui ne répond jamais | Dans le tableau : `FALLA` avec le détail `se acabo el tiempo de espera`, un temps proche du délai d'expiration configuré | Le `context.WithTimeout` par service (leçon 6) a coupé l'attente ; ce n'est pas un bug, c'est le mécanisme qui fonctionne |
| Fermer le corps d'une réponse nulle | (si c'était mal fait) `panic: runtime error: invalid memory address or nil pointer dereference` | Cela arriverait si `defer resp.Body.Close()` était placé avant de vérifier `err` ; le `revisor` l'évite en vérifiant l'erreur d'abord (section 7.1) |
| `error` sérialisée directement en JSON sans traduction | `{}` (objet vide, sans aucune donnée utile, sans aucune erreur de compilation pour prévenir) | `encoding/json` ne sait pas convertir une interface `error` ; c'est pourquoi `reporte` utilise `lineaJSON` avec le motif déjà converti en texte (section 7.4) |
| Code de sortie non vérifié dans un script | Le script continue comme si de rien n'était, même si un service a échoué | `ejecutar` distingue bien 0/1/2 (section 7.6) ; celui qui intègre le `revisor` dans un pipeline doit lire `$?`, pas seulement la sortie affichée |
| Binaire dans une image `FROM scratch`, contre HTTPS | `no responde: Get "https://...": tls: failed to verify certificate: x509: certificate signed by unknown authority` | Il manque la liste des autorités de certification (`ca-certificates.crt`) ; copie-la depuis une image qui l'a (section 7.10) |

---

## Ce qu'on fait mal

- **Utiliser `http.DefaultClient` ou `http.Get` directement en production.** Section 7.1 : sans délai
  d'expiration propre, une connexion qui ne répond jamais bloque le programme pour toujours, sans aucun
  symptôme préalable.
- **Mettre toute la logique dans `main` et appeler `os.Exit` au milieu du code.** Cela rend la fonction
  impossible à tester avec `go test`, parce que `os.Exit` termine le processus de test en même temps que le
  programme. Sépare « décider du code de sortie » de « sortir vraiment » (section 7.2).
- **Sérialiser `error` directement en JSON, en espérant que « quelque chose sorte ».** Il sort un objet vide,
  sans aucun avertissement que quelque chose n'a pas pu être converti — le type d'échec silencieux le plus
  difficile à détecter, parce que le programme ne plante pas et ne se plaint pas.
- **Fermer le corps d'une réponse HTTP avant de vérifier l'erreur.** Si la requête a échoué, `resp` est `nil`,
  et appeler une méthode dessus fait un panic. Vérifie l'erreur d'abord, toujours.
- **Ne pas vider le corps de la réponse avant de le fermer.** La connexion ne peut pas être réutilisée, et un
  programme qui fait beaucoup de requêtes vers le même serveur finit par ouvrir une nouvelle connexion à chaque
  fois, plus lentement que nécessaire, sans aucune erreur pour le signaler.
- **Ignorer le code de sortie du binaire depuis un script ou un pipeline.** Le `revisor` distingue « j'ai bien
  tourné, tout est sain » (0), « j'ai bien tourné, quelque chose a échoué » (1) et « je n'ai même pas pu
  démarrer » (2) — le gaspiller en ne lisant que la sortie affichée, c'est jeter une information que le
  programme te donne déjà gratuitement.

---

## Exercices

1. Implémente (ou révise, si tu l'as déjà) le vrai `Revisor` HTTP de la section 7.1, avec son propre client et
   son délai d'expiration. Confirme qu'il compile et que `go vet ./...` ne se plaint pas.
2. Démarre `cmd/servidor-demo` sur un port libre et lance ton `revisor` contre ses quatre routes, comme dans la
   section 7.5. Colle le vrai tableau que tu as obtenu, pas un tableau inventé.
3. Ajoute à `Tabla` (section 7.4) une cinquième colonne, `PROTOCOLO`, qui indique `https` ou `http` selon
   `Servicio.EsSeguro()` (la méthode définie dans la leçon 3). Ajuste la largeur fixe de cette colonne à la
   main, sans avoir besoin du calcul dynamique dont dispose déjà `anchoNombre`.
4. Lance le même rapport avec `--formato json` et valide que le résultat est du JSON légitime (par exemple,
   passe-le dans `jq .` ou colle-le dans un validateur JSON).
5. Provoque l'erreur de `--formato` invalide (section 7.6, tableau) et confirme le code de sortie avec
   `echo $?`.
6. Compile ton `revisor` pour `GOOS=linux GOARCH=amd64` depuis ta machine actuelle et confirme avec `file` que
   le binaire obtenu est celui de la bonne plateforme.
7. (Un peu plus difficile) Mesure la taille de ton binaire avec et sans `-ldflags="-s -w"`, comme dans la
   section 7.7, et calcule le pourcentage de réduction.
8. (Clôture du cours) Lance la suite complète du projet — `go vet ./...`, `gofmt -l .` et
   `go test ./... -race -cover` — et colle la sortie complète dans ton journal de bord. Si quelque chose n'est
   pas au vert, corrige-le avant de considérer le cours comme terminé.

### Corrigés

1-7. Il n'y a pas de sortie de référence unique, car cela dépend de ta propre implémentation et de ta propre
machine ; compare la forme de ton résultat avec les sections correspondantes (7.1, 7.4, 7.5, 7.6, 7.7).

8. L'exécution réelle, sur le `revisor` de ce cours, le 30-sep-2026 :

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

   `cmd/servidor-demo` apparaît à 0.0% sans `ok` ni `FAIL`, non pas parce que quelque chose est cassé : avec
   `-cover`, un paquet sans aucun fichier `_test.go` est quand même instrumenté et rapporté, mais il n'y a
   aucun test qui exerce ce code. C'est un outil d'appui pour tester le programme complet, pas une partie du
   `revisor` qui est distribué, et il n'a pas de logique propre qui décide quoi que ce soit — il se contente de
   répondre ce qu'on l'a programmé pour répondre — c'est pourquoi il n'a pas besoin de tests propres.

   ⚠️ **Le 61.9% de `internal/revisar` dans ce tableau est l'artefact qu'explique la section 5.7 : par paquet,
   il ne compte pas ce que `TestEjecutar_reportaOKyFalla` (dans `cmd/revisor`, un peu plus haut dans cette même
   exécution) exerce de `Revisar` en parlant vraiment HTTP contre le serveur `httptest`.** Mesuré avec
   `-coverpkg=./...` sur tout le projet : `Revisar` monte à 86,4 % et le total du projet est de 81,5 %, pas
   61,9 %. Si tu vas citer un chiffre de couverture pour décider quelque chose, lance `-coverpkg=./...`, ne te
   fie pas au seul chiffre par paquet.

---

## Comment savoir que j'y suis arrivé

- [ ] Mon `Revisor` HTTP a son propre client, avec un délai d'expiration, et n'utilise jamais
      `http.DefaultClient`.
- [ ] J'ai séparé `main` (qui ne fait qu'appeler `os.Exit`) d'une fonction qui fait le travail et renvoie un
      `int`.
- [ ] Mon binaire accepte `--config`, `--formato`, `--paralelo` et `--limite` depuis la ligne de commande.
- [ ] J'ai lancé le vrai `revisor` contre un vrai serveur (le mien ou le `servidor-demo` du cours) et le tableau
      obtenu reflète exactement ce que fait ce serveur.
- [ ] `--formato json` produit du JSON valide, vérifié avec un outil qui n'est pas moi en train de le lire.
- [ ] `echo $?` après avoir lancé le `revisor` donne 0, 1 ou 2, et je sais expliquer chacun.
- [ ] J'ai compilé pour une autre plateforme (`GOOS`/`GOARCH`) et j'ai confirmé avec `file` que le binaire est le
      bon.
- [ ] `go vet ./...`, `gofmt -l .` et `go test ./... -race -cover` sont propres sur ma propre copie du projet —
      je ne le suppose pas, je l'ai lancé.

---

## Pour aller plus loin

1. [net/http package](https://pkg.go.dev/net/http) — la documentation officielle complète du client et du
   serveur HTTP de la bibliothèque standard.
2. [Command flag](https://pkg.go.dev/flag) — la documentation officielle du paquet d'options utilisé dans
   cette leçon.
3. [encoding/json: JSON and Go](https://go.dev/blog/json) — l'article officiel sur la façon dont Go traduit
   entre structs et JSON, y compris les règles des étiquettes (`json:"..."`, `omitempty`, `-`).
4. [Build constraints et compilation croisée](https://pkg.go.dev/cmd/go#hdr-Environment_variables) — la
   référence de `GOOS`/`GOARCH` et des autres variables d'environnement qui contrôlent `go build`.
5. [text/tabwriter](https://pkg.go.dev/text/tabwriter) — l'outil standard pour aligner des colonnes quand le
   calcul manuel de la section 7.4 ne suffit plus (plusieurs colonnes à largeur variable).

### Termes de cette leçon

| Terme | Ce qu'il signifie |
|---|---|
| `http.Client` | le type qui fait les requêtes HTTP en Go ; non configuré, il n'a pas de limite de temps |
| `flag.FlagSet` | ensemble d'options de ligne de commande qui peut être analysé à partir d'un slice propre, pas seulement de `os.Args` |
| `omitempty` | option d'étiquette JSON qui retire un champ de la sortie quand il vaut zéro ou est vide |
| code de sortie | le nombre entier qu'un programme renvoie en se terminant ; 0 signifie succès par convention universelle |
| `GOOS` / `GOARCH` | variables d'environnement qui indiquent à `go build` pour quel système d'exploitation et quelle architecture compiler |
| `-ldflags="-s -w"` | options de compilation qui retirent les symboles et les données de débogage pour réduire la taille du binaire |

---

## Et maintenant, la suite

Tu as un programme complet : il interroge de vrais services, tous en même temps, avec des limites de temps, et
produit un rapport qu'un humain ou un programme peut lire. Trois choses le rendent plus solide, dans l'ordre où
il vaut le mieux les attaquer :

1. **`golangci-lint`** — réunit des dizaines d'analyseurs statiques en une seule exécution. Passe-le sur ton
   propre `revisor` et lis ce qu'il dit : il t'apprend presque toujours quelque chose que ni `go vet` ni les
   tests n'attrapent.
2. **[Effective Go](https://go.dev/doc/effective_go) encore une fois.** Maintenant que tu as écrit un
   programme complet, tu vas le lire autrement que dans la leçon 0.
3. **Lis du bon code écrit par d'autres.** Le paquet `net/http` lui-même, de la bibliothèque standard, est un
   bon point de départ : tu as déjà de quoi t'orienter pour comprendre pourquoi il est écrit comme il l'est.

**Et ce qui était promis depuis le README :** tu vas écrire ce même programme une nouvelle fois, en Rust, dans
le cours frère. Pas pour comparer la syntaxe, mais pour voir le même problème résolu avec d'autres outils —c'est
ce qui t'apprend vraiment en quoi les deux langages diffèrent.

---

**Précédent :** [Leçon 6 — Concurrence](06-concurrencia.md)
