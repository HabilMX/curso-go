# Leçon 6 — Concurrence

**Durée :** 2 séances de 60-90 minutes. C'est la raison pour laquelle Go existe, et la seule leçon du cours
qui se fait en plusieurs passes : les goroutines, puis les canaux et `WaitGroup`, puis `context` et la limite
de parallélisme appliquée au `revisor`. En une seule passe, cela ne prend pas — c'est ce que confirme l'ordre
même qu'utilise le cours de Go qui compte le plus d'apprenants au monde.

**À la fin, tu seras capable de :**

- Lancer une goroutine et expliquer, avec une vraie preuve, pourquoi le programme peut se terminer avant
  qu'elle ne se termine.
- Utiliser `sync.WaitGroup` pour attendre un groupe de goroutines, et reconnaître le message exact que donne
  Go quand les `Add`/`Done` ne concordent pas.
- Expliquer pourquoi le `revisor` distribue les résultats par un canal au lieu d'écrire dans une map
  partagée, et provoquer toi-même l'erreur qui se produirait s'il ne le faisait pas.
- Lire, littéralement, la sortie de `go test -race` quand il trouve une vraie data race.
- Utiliser `context.WithTimeout` pour imposer une limite de temps à une opération, et expliquer pourquoi
  `defer cancelar()` n'est pas facultatif.
- Limiter le nombre de goroutines qui tournent en même temps avec un sémaphore de canal, et mesurer que la
  limite est vraiment respectée.

---

## Pourquoi c'est important

Jusqu'à la leçon 5, le `revisor` interroge les services **un par un** : si tu as dix services et que chacun
met une seconde à répondre, le rapport complet prend dix secondes, quel que soit l'ordre. Avec cent services,
presque deux minutes. Et le temps d'attente, ce n'est pas toi qui le décides : c'est le service le plus lent de
la liste qui le décide, multiplié par le nombre de services.

Cela n'a pas à être ainsi, parce qu'**interroger un service, c'est attendre, pas travailler** : pendant presque
tout le temps de cette attente, le CPU ne fait rien, il attend seulement une réponse du réseau. Go a été conçu
par des gens qui, chez Google, passaient leurs journées à attendre exactement cela —des réponses réseau entre
des milliers de services— et c'est pourquoi la concurrence n'est pas une bibliothèque ajoutée après coup :
elle fait partie du langage dès la première ligne (`go`, un mot réservé, pas une fonction d'une bibliothèque).

**Le résultat que tu vas construire dans cette leçon :** les mêmes dix services d'une seconde chacun,
interrogés **tous en même temps**, en un peu plus d'une seconde au total au lieu de dix. Ce chiffre —la
différence entre « en série » et « en parallèle »— est la récompense de la leçon, et tu vas le mesurer
toi-même, pas le croire sur parole.

🔑 **Et l'avertissement qui rend cette leçon différente des précédentes :** dans les leçons 2 à 5, un
programme qui compile et dont les tests passent est presque toujours correct. **En concurrence, non.** Un
programme avec une data race peut compiler, tourner, passer ses tests quatre-vingt-dix-neuf fois, et échouer
la centième — ou ne jamais échouer sur ta machine et échouer tous les jours en production, avec plus de cœurs
et plus de charge. Les erreurs de cette leçon sont celles qui **ne se voient pas à l'œil nu**, et c'est
pourquoi le point 4 de cette leçon —les vraies erreurs, provoquées et capturées— est le plus important de tout
le cours.

---

## Les concepts

### 6.1 Les goroutines : démarrer est trivial, attendre non

Lancer une goroutine, c'est un mot :

<!-- verificar:fragmento -->
```go
go revisar(s)
```

Et c'est tout. Cela suffit pour que `revisar(s)` tourne **de façon concurrente** avec le reste du programme,
sans attendre qu'elle se termine pour passer à la ligne suivante. Chacune coûte environ 2 Ko de mémoire au
départ (elles grandissent si nécessaire), pas les mégaoctets d'un thread du système d'exploitation — c'est
pourquoi un programme en Go peut avoir des centaines de milliers de goroutines vivantes sans rien prévoir de
spécial, ce qui serait impensable avec des threads du système.

Mais regarde ce programme, avec le bug le plus courant du premier jour de concurrence en Go (complet dans
`programas/revisor/ejemplos/06-goroutines-sin-esperar/main.go`) :

<!-- verificar:ejemplo:ejemplos/06-goroutines-sin-esperar:nodeterminista -->
```go
func main() {
	servicios := []string{"catalogo", "pagos", "inventario"}
	for _, s := range servicios {
		go fmt.Println("revisando", s)
	}
	// el programa termina AQUÍ, sin haber esperado a ninguna goroutine
}
```

Je l'ai lancé **260 fois** sur cette machine pour mesurer à quelle fréquence on VOIT le problème, pas pour le
deviner : 40 avec `go run`, 20 de plus en forçant un seul processeur logique (`GOMAXPROCS=1`, pour ôter au
programme toute chance qu'une goroutine parvienne à tourner sur un autre cœur pendant que `main` se termine), et
200 avec le binaire déjà compilé (`go build` + `./gor61`, pour écarter le temps que prend la compilation). Le
résultat :

```
$ go run main.go
$ go run main.go
$ go run main.go
```

**Vide. Les trois, et pratiquement toutes les 260** (une seule des 200 exécutions du binaire compilé a affiché
quelque chose — le reste, rien).

🔴 **Et il faut lire ce chiffre avec beaucoup de soin, parce qu'il est facile d'en tirer la mauvaise
conclusion. Le défaut n'est pas présent « 1 fois sur 260 » : il est présent 260 fois sur 260, sans
exception.** Dans **aucune** des 260 exécutions le programme n'a attendu ses goroutines — pas même dans la seule
qui a affiché quelque chose. Ce qui change d'une exécution à l'autre, ce n'est pas si le programme attend (il ne
le fait jamais) : c'est si, par pur hasard de la façon dont le système d'exploitation répartit le temps entre
les processus, une goroutine parvient à exécuter `fmt.Println` dans l'interstice de temps qu'il y a entre son
lancement et le moment où `main` tue le programme entier. Cet interstice ne suffit presque jamais — c'est
pourquoi on ne voit presque jamais rien — mais le programme est **tout aussi cassé** dans les 259 exécutions
silencieuses que dans celle qui a affiché quelque chose.

**Voici ce qu'il faut retenir : toujours cassé, presque jamais visible.** Ce n'est pas « il y a une petite
probabilité que cela échoue » —le lire ainsi est exactement l'erreur à éviter— : c'est que le programme ne fait
**jamais** ce qui est correct, et que la plupart du temps le symptôme n'a pas le temps de se montrer pour que tu
le remarques. C'est ainsi que ce type de bug s'infiltre : il passe la revue de code (cela compile, tourne, ne
plante pas), il passe les tests manuels (qui lance un programme 260 fois pour lui faire confiance ?), il passe
la CI (qui l'exécute une fois, peut-être deux) — et il ne se manifeste qu'une fois en production, avec plus de
charge, plus de cœurs, et un interstice de temps qui, un jour, finit par s'ouvrir, au pire moment possible pour
le découvrir.

⚠️ **Ne le prends pas non plus comme « cela n'affiche jamais rien, sur aucune machine ».** Avec plus de charge,
un autre système d'exploitation, ou simplement de la malchance, l'interstice peut s'ouvrir plus souvent — de
fait, il s'est ouvert une fois dans ces mêmes 260 exécutions. Ce qui ne change pas d'une machine à l'autre,
c'est que le programme n'attend **jamais** : c'est structurel, cela ne dépend pas de la chance. Lance toi-même
l'expérience sur ta machine, avec suffisamment de répétitions pour que le chiffre signifie quelque chose, et
compare.

**Ce n'est pas un bug intermittent du programme : c'est que `main` se termine dès qu'il arrive à la fin de son
corps, sans se soucier qu'il y ait encore des goroutines en cours** — et quand `main` se termine, le programme
entier se termine avec lui, goroutines vivantes comprises, sans avertissement et sans erreur. Lancer une
goroutine est facile ; le vrai travail est de t'assurer que le programme attend qu'elles se terminent avant de
continuer.

### 6.2 `sync.WaitGroup` : la bonne façon d'attendre

<!-- verificar:fragmento -->
```go
var wg sync.WaitGroup
for _, s := range servicios {
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("revisando", s)
	}()
}
wg.Wait() // aquí sí espera a que las tres hayan llamado Done
```

Le modèle est toujours le même, et il vaut mieux le mémoriser dans cet ordre :

1. **`wg.Add(1)` avant de lancer la goroutine**, pas à l'intérieur. Si tu le mets à l'intérieur, `Wait()` peut
   s'exécuter avant que la goroutine n'ait eu le temps de faire son propre `Add`, et alors il ne compte pas
   avec cette attente.
2. **`defer wg.Done()` comme première ligne de la goroutine.** Le `defer` garantit qu'il s'exécute quoi qu'il
   arrive à l'intérieur —même si la goroutine fait un panic—, et le mettre en premier évite qu'un `return`
   anticipé au milieu du code te fasse oublier de le compter.
3. **`wg.Wait()` là où tu as vraiment besoin du résultat**, normalement juste avant d'utiliser ce que les
   goroutines ont produit.

**Ce qui se passe si tu sautes l'étape 1 — `Add` manquant —, provoqué exprès** (complet dans
`programas/revisor/ejemplos/06-waitgroup-sin-add/main.go`) :

<!-- verificar:ejemplo:ejemplos/06-waitgroup-sin-add -->
```go
var wg sync.WaitGroup
for i := 0; i < 5; i++ {
	go func(n int) { // <- sin wg.Add(1) antes de esta línea
		defer wg.Done()
		time.Sleep(10 * time.Millisecond)
	}(i)
}
wg.Wait()
fmt.Println("listo")
```

Je l'ai lancé **six fois d'affilée, avec `-race`**, pour que tu voies ce qui se passe vraiment, pas seulement
l'erreur finale :

```
$ go run -race sinesperar.go
listo
panic: sync: negative WaitGroup counter
...
exit status 2

$ go run -race sinesperar.go
listo
panic: sync: negative WaitGroup counter
...
exit status 2

(cuatro corridas más, idénticas: "listo" primero, el panic después — 6 de 6, CON -race)
```

🔑 **Lis cela avec soin, parce que c'est l'erreur la plus trompeuse de la leçon : le programme affiche `listo`
AVANT de planter, les six fois.** Ce n'est pas un hasard de cette exécution : `wg.Wait()` n'attend rien parce que
le compteur n'a jamais été incrémenté (il n'y a jamais eu d'`Add(1)`), donc `Wait()` revient
**immédiatement**, `main` affiche `listo` comme si tout s'était bien passé, et **c'est seulement alors** que
l'une des goroutines retardataires —qui tournait bel et bien, simplement Go ne l'a jamais attendue— termine son
`Sleep` et appelle `Done()`, en retranchant un à un compteur qui était déjà à zéro. Le panic n'est pas la
première chose que tu vois : c'est la dernière, après que le programme t'a déjà dit que tout allait bien.

⚠️ **Et voici la partie que j'ai mal mesurée la première fois, et il vaut la peine de laisser la correction par
écrit : sans `-race`, ce même programme ne plante jamais.** Les mêmes six exécutions, sans `-race` :

```
$ go run sinesperar.go
listo

$ go run sinesperar.go
listo

(cuatro corridas más, idénticas: "listo" y nada más — 6 de 6, SIN -race, saliendo con código 0)
```

**Le `time.Sleep(10 * time.Millisecond)` n'est pas ce qui FAIT apparaître le panic — c'est ce qui
l'EMPÊCHE**, sans `-race`. `main` arrive à `wg.Wait()`, qui revient immédiatement, affiche `listo`, et le
programme entier se termine **bien avant** que les 10 millisecondes ne s'écoulent — si bien qu'aucune des cinq
goroutines retardataires n'a même le temps de se réveiller de son `Sleep` et d'appeler `Done()`. Le programme
a, littéralement, l'air « terminé avec succès » : code de sortie 0, sans aucune trace que quelque chose était
mal monté. Ce qui fait bel et bien apparaître le panic, de façon constante, c'est l'instrumentation de `-race` :
surveiller chaque accès mémoire pour détecter les data races rend le programme nettement plus lent, et ce
ralentissement supplémentaire est exactement ce qui donne à une goroutine retardataire le temps de terminer son
`Sleep` et d'appeler `Done()` avant que le processus ne se termine.

**C'est exactement le type d'erreur qui n'échoue pas à chaque fois, et c'est pourquoi c'est la plus précieuse de
cette leçon — et maintenant avec la bonne donnée : tu viens de mesurer qu'il n'est même pas nécessaire d'avoir
« un travail plus court » pour qu'elle passe inaperçue. En l'exécutant tel quel, SANS `-race`, elle passe déjà
inaperçue 6 fois sur 6.** Un apprenant qui lance cet exemple sans `-race` —ce qui est évident, si personne ne
l'avertit— va voir `listo` et rien d'autre, et va conclure, à juste titre, que le programme fonctionne. **Il ne
fonctionne pas : l'`Add(1)` manque toujours exactement de la même façon.** La seule chose qui a changé, c'est si
quelque chose a eu le temps de se montrer pour le trahir — le même modèle « toujours cassé, presque jamais
visible » de la section 6.1, ici avec un acteur différent : ce n'est pas la charge de la machine, c'est si tu as
exécuté avec le détecteur de data races activé ou non.

**`Done()` retranche un au compteur interne du `WaitGroup`.** Si tu n'as jamais fait `Add(1)`, le compteur
commence à zéro, et lui retrancher un le rend négatif — et Go, au lieu de le laisser passer en silence, fait un
panic avec un message qui dit exactement ce qui ne va pas : `negative WaitGroup counter`. Ce message littéral est
ton indice : si tu le vois, il manque presque toujours un `Add(1)` quelque part, ou il y a plus de `Done()`
qu'il ne devrait y en avoir. Mais la vraie leçon n'est pas le message du panic — c'est que **le programme avait
déjà dit `listo` avant qu'il n'apparaisse.**

### 6.3 Pourquoi le `revisor` ne partage pas de map : les canaux

La devise du langage, et elle vaut la peine d'être mémorisée telle quelle : **« ne communique pas en partageant
la mémoire ; partage la mémoire en communiquant ».** Au lieu que chaque goroutine écrive son résultat dans une
map ou un slice partagé —ce qui exigerait de protéger chaque accès avec un verrou—, chacune envoie son résultat
par un **canal**, et une seule goroutine (ou le code principal) les récupère de l'autre côté.

<!-- verificar:fragmento -->
```go
ch := make(chan servicio.Estado)         // sin buffer: cada envío espera a que alguien reciba
ch := make(chan servicio.Estado, 10)     // con buffer: hasta 10 caben sin esperar a que nadie reciba

ch <- estado        // enviar
e := <-ch           // recibir
close(ch)            // cerrar: después de cerrado, recibir sigue funcionando hasta vaciarlo
for e := range ch { ... }   // recibe hasta que el canal se cierre Y se vacíe
```

**J'ai vérifié ce qui se passe si, au lieu d'un canal, j'utilise un slice partagé sans protection**, exactement
l'erreur que les canaux évitent. Ce programme lance mille goroutines qui incrémentent la même variable (complet
dans `programas/revisor/ejemplos/06-carrera-de-datos/main.go`) :

<!-- verificar:ejemplo:ejemplos/06-carrera-de-datos:nodeterminista -->
```go
contador := 0
var wg sync.WaitGroup
for i := 0; i < 1000; i++ {
	wg.Add(1)
	go func() {
		defer wg.Done()
		contador++ // dos goroutines pueden leer el mismo valor antes de que la otra escriba
	}()
}
wg.Wait()
fmt.Println("contador:", contador)
```

Sans le détecteur de data races, **le programme ne plante pas — il donne seulement un résultat incorrect**, et
différent à chaque fois :

```
$ go run carrera.go
contador: 956
```

**956, pas 1000.** Aucune erreur, aucun panic, aucun signe que quelque chose s'est mal passé — seulement un
nombre plus petit que ce qu'il devrait être, parce que certaines des mille additions se sont perdues quand deux
goroutines ont lu `contador` en même temps, ont toutes deux ajouté un à la même ancienne valeur, et que l'une
des deux écritures a écrasé l'autre sans que personne ne s'en aperçoive. **C'est le bug le plus dangereux de la
concurrence : il n'échoue pas, il donne un résultat plausible et faux.** Un programme de ce genre peut passer la
revue, passer des tests superficiels, et échouer en production de façon sporadique pendant des mois avant que
quelqu'un ne le remarque.

### 6.4 `go test -race` : voir le détecteur trouver le bug

Le même programme, exécuté avec `-race`, le trahit bel et bien — avec la preuve exacte de l'endroit :

```
$ go run -race carrera.go
==================
WARNING: DATA RACE
Read at 0x00c00013c018 by goroutine 11:
  main.main.func1()
      carrera.go:15 +0x68

Previous write at 0x00c00013c018 by goroutine 8:
  main.main.func1()
      carrera.go:15 +0x78

Goroutine 11 (running) created at:
  main.main()
      carrera.go:13 +0x6c

Goroutine 8 (finished) created at:
  main.main()
      carrera.go:13 +0x6c
==================
contador: 847
Found 2 data race(s)
exit status 66
```

Lis-le de haut en bas, parce que chaque bloc répond à une question différente :

- **`Read at ... by goroutine 11`** et **`Previous write at ... by goroutine 8`** : deux goroutines différentes
  ont touché la **même adresse mémoire** (`0x00c00013c018`, la variable `contador`), l'une en lecture et
  l'autre en écriture, sans aucun mécanisme garantissant que l'une attende l'autre.
- **`carrera.go:15`** dans les deux : la ligne exacte du code où cela s'est produit (`contador++`), la même
  ligne pour les deux, parce que c'est la seule ligne qui touche cette variable.
- **`Goroutine 11 (running) created at ... carrera.go:13`** : d'où vient cette goroutine —la ligne du
  `go func() {...}()` à l'intérieur du `for`—, pour que tu puisses retrouver quel lancement c'était.
- **`Found 2 data race(s)`** et **`exit status 66`** : le détecteur ne s'arrête pas à la première data race
  qu'il trouve ; il continue de tourner et les rapporte toutes, et le programme se termine avec un code de
  sortie différent de 0 et de 1 (66 est le code que Go réserve à cela), pour qu'un pipeline d'intégration
  continue puisse le distinguer d'un échec normal.

**Et le chiffre final, `contador: 847`, a changé par rapport à l'exécution sans `-race` (956).** Pas par
hasard : `-race` fait tourner le programme plus lentement et avec plus d'instrumentation, ce qui change l'ordre
exact dans lequel les goroutines s'entrelacent — une autre raison pour laquelle ce bug est si traître : le
chiffre faux n'est même pas le même faux à chaque fois.

🔴 **La règle qui en découle, sans exception :** un programme concurrent qui passe ses tests **sans** `-race`
n'est pas testé. Lance `-race` dès le premier test de concurrence que tu écris, pas quand « quelque chose a
l'air bizarre » — parce que, comme tu viens de le voir, rien n'a l'air bizarre jusqu'à ce qu'il soit trop tard.

### 6.5 Les deux autres façons d'échouer : panic et deadlock

**Panic, si tu fermes mal un canal :**

- Envoyer sur un canal déjà fermé fait un panic : `panic: send on closed channel`.
- Fermer un canal deux fois fait un panic : `panic: close of closed channel`.

La règle qui évite les deux : **c'est celui qui envoie qui ferme le canal, jamais celui qui reçoit**, et il le
ferme une seule fois, normalement depuis une goroutine dédiée qui sait quand il n'y aura plus d'envois (tu le
verras dans la section 6.7, avec `wg.Wait()` suivi de `close`).

**Deadlock, si personne n'écoute de l'autre côté.** Je l'ai provoqué avec le programme le plus court possible
(complet dans `programas/revisor/ejemplos/06-deadlock/main.go`) :

<!-- verificar:ejemplo:ejemplos/06-deadlock:nodeterminista -->
```go
func main() {
	canal := make(chan int)
	canal <- 1 // nadie del otro lado está leyendo: se bloquea
	fmt.Println(<-canal)
}
```

```
$ go run candado.go
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send]:
main.main()
	candado.go:7 +0x38
exit status 2
```

**Ceci corrige une chose que beaucoup d'explications tiennent pour acquise : un deadlock en Go ne « reste pas
toujours bloqué pour toujours ».** Le runtime de Go lui-même a un détecteur de deadlocks : si à un moment donné
**toutes** les goroutines du programme sont endormies à attendre quelque chose qui n'arrivera jamais, Go s'en
aperçoit et tue le programme avec `fatal error: all goroutines are asleep - deadlock!`, avec la ligne exacte où
il s'est coincé (`[chan send]`, dans ce cas, parce qu'il était en train d'envoyer). Si ton programme semble
« bloqué pour toujours » au lieu de se terminer avec cette erreur, ce n'est presque sûrement **pas** un vrai
deadlock : il est plus probable que tu aies au moins une goroutine vivante qui fait autre chose (par exemple, un
minuteur ou un serveur HTTP à l'écoute), et c'est elle qui empêche le runtime de déclarer que « toutes » sont
endormies.

### 6.6 `context` : comment on annule et comment on impose une limite de temps

<!-- verificar:fragmento -->
```go
ctx, cancelar := context.WithTimeout(context.Background(), 5*time.Second)
defer cancelar() // SIEMPRE, incluso si terminas antes de que se cumplan los 5 segundos

req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
resp, err := http.DefaultClient.Do(req) // se aborta solo si pasan los 5 s
```

`context` est le standard maison en Go pour deux choses : **imposer une limite de temps** à une opération qui
pourrait prendre trop longtemps, et **propager une annulation** vers le bas (si celui du dessus est annulé, tout
ce qui dépend de lui est annulé aussi, sans que chaque fonction ait à réinventer son propre mécanisme).

⚠️ **`defer cancelar()` n'est pas facultatif, même quand l'opération s'est déjà terminée d'elle-même.**
`context.WithTimeout` démarre un minuteur interne ; si tu n'appelles jamais la fonction d'annulation, ce
minuteur reste vivant jusqu'à l'expiration du délai d'origine, en retenant de la mémoire pendant tout ce temps
— et si ton programme crée des contextes ainsi en permanence (un par service interrogé, comme le `revisor`),
sans `defer cancelar()` tu accumules des fuites de mémoire proportionnelles au nombre de requêtes que tu as
faites. C'est la fuite la plus courante des programmes Go qui utilisent `context`, et c'est pourquoi il vaut
mieux écrire le `defer` sur la même ligne que celle où tu crées le contexte, avant d'écrire quoi que ce soit
d'autre.

### 6.7 `Todos` : la fonction qui réunit les trois pièces

Voici, pour de vrai, à quoi ressemble la fonction centrale du `revisor`
(`programas/revisor/internal/revisar/todos.go`), après avoir réuni goroutines, canaux, sémaphore et
`context` :

<!-- verificar:extracto:internal/revisar/todos.go -->
```go
func Todos(ctx context.Context, r Revisor, servicios []servicio.Servicio, paralelo int) []servicio.Estado {
	if len(servicios) == 0 {
		return nil
	}
	if paralelo <= 0 {
		paralelo = ParaleloPorOmision
	}

	resultados := make(chan servicio.Estado, len(servicios))
	semaforo := make(chan struct{}, paralelo)

	var wg sync.WaitGroup
	for _, s := range servicios {
		wg.Add(1)
		go func() {
			defer wg.Done()

			semaforo <- struct{}{}        // pide turno; se bloquea si ya hay «paralelo» corriendo
			defer func() { <-semaforo }() // devuelve el turno pase lo que pase

			// Cada servicio tiene su propio tiempo límite, hijo del general. Si
			// el de arriba se cancela, este muere con él.
			propio, cancelar := context.WithTimeout(ctx, s.TimeoutEfectivo())
			defer cancelar() // sin esto, el temporizador no se libera: es la fuga más común de Go

			resultados <- r.Revisar(propio, s)
		}()
	}

	// Quien envía cierra, nunca quien recibe. Y se cierra desde otra goroutine
	// porque Wait tiene que poder esperar mientras el bucle de abajo ya está
	// recibiendo: si cerráramos aquí mismo, con el canal lleno nos trabaríamos.
	go func() {
		wg.Wait()
		close(resultados)
	}()

	estados := make([]servicio.Estado, 0, len(servicios))
	for e := range resultados {
		estados = append(estados, e)
	}
	return estados
}
```

Lis-la avec les sections précédentes en tête, parce que chaque pièce répond à un problème que tu as déjà vu :

- **Une goroutine par service** (6.1), comptée avec **`wg.Add(1)`/`defer wg.Done()`** (6.2) pour que le
  programme sache quand toutes se sont terminées.
- **Un canal `resultados` avec tampon** (6.3) recueille les états : personne n'écrit dans un slice ou une map
  partagés, donc aucun verrou n'est nécessaire pour cette partie.
- **`semaforo := make(chan struct{}, paralelo)`** est un canal utilisé comme quota : il a de la place pour
  `paralelo` valeurs, donc la `paralelo + 1`-ième goroutine qui essaie d'y écrire (`semaforo <-
  struct{}{}`) se bloque jusqu'à ce qu'une autre libère sa place (`<-semaforo`, dans le `defer`). C'est le
  même canal avec tampon que dans la section 6.3, utilisé non pas pour transporter des données mais pour
  tenir le compte des « tours » qui restent.
- **Un `context.WithTimeout` propre à chaque service** (6.6), enfant du `ctx` général : si celui du dessus est
  annulé (par exemple, si le temps total du rapport est épuisé), tous les enfants sont annulés avec lui.
- **`close(resultados)` depuis UNE AUTRE goroutine, après `wg.Wait()`** — et cela mérite une explication, parce
  que c'est la partie qui se voit le moins à l'œil nu : si nous fermions le canal dans la même goroutine qui fait
  le `for e := range resultados` plus bas, nous nous bloquerions, parce que `Wait()` a besoin que toutes les
  goroutines lisent dans le canal pour libérer de la place et pouvoir rendre leur tour, mais la boucle de
  lecture ne démarrerait jamais parce que nous attendrions d'abord `Wait()`. En le lançant à part, la fermeture
  et la lecture se produisent **en même temps**, pas l'une après l'autre.

### 6.8 Le tester sans réseau : `Falso` et la limite de parallélisme mesurée

Tester `Todos` contre de vrais services serait lent et non déterministe. Le `revisor` utilise un `Revisor`
factice (`programas/revisor/internal/revisar/falso.go`) qui simule des réponses, des délais et même des services
qui ne répondent jamais, tout en mémoire :

<!-- verificar:fragmento -->
```go
type Falso struct {
	Respuestas map[string]RespuestaFalsa
	mu         sync.Mutex
	Contador   int
}

func (f *Falso) Revisar(ctx context.Context, s servicio.Servicio) servicio.Estado {
	f.mu.Lock()
	f.Contador++
	f.mu.Unlock()
	// ...
}
```

🔒 **Ici, un mutex est bel et bien nécessaire, et c'est l'exception à la règle « préfère les canaux » de la
section 6.3.** `Contador` est un entier partagé que **beaucoup de goroutines incrémentent en même temps**
pendant les tests (`Falso.Revisar` est appelé de façon concurrente, une fois pour chaque service que `Todos`
est en train d'interroger). Un canal servirait à *envoyer* des résultats, mais pour un simple compteur partagé,
un `sync.Mutex` autour de l'unique ligne qui le touche est plus simple et plus clair. La règle n'est pas
« n'utilise jamais de mutex » : c'est « avant d'en utiliser un, demande-toi si un canal exprime mieux ce que tu
es en train de faire » — et pour envoyer des résultats, presque toujours oui ; pour un compteur partagé, le
mutex est presque toujours le bon outil.

Avec `Falso`, ce test mesure quelque chose qu'il serait autrement presque impossible de vérifier avec
confiance : que le sémaphore de la section 6.7 limite vraiment le nombre de requêtes qui tournent en même
temps, pas seulement qu'il « fonctionne en général » (version complète, non abrégée, dans
`programas/revisor/internal/revisar/todos_test.go`) :

<!-- verificar:fragmento -->
```go
func TestTodos_elParaleloLimitaCuantasCorrenALaVez(t *testing.T) {
	const totalServicios = 20
	const limite = 3
	var enVuelo, maximoObservado int32
	// medidorDeConcurrencia cuenta, con un contador atómico, cuántas llamadas
	// a Revisar están abiertas AL MISMO TIEMPO, y se queda con el máximo.
	medidor := &medidorDeConcurrencia{enVuelo: &enVuelo, maximoObservado: &maximoObservado, espera: 15 * time.Millisecond}
	servicios := make([]servicio.Servicio, totalServicios)
	// ... llena servicios ...
	Todos(context.Background(), medidor, servicios, limite)
	if max := atomic.LoadInt32(&maximoObservado); max > int32(limite) {
		t.Errorf("se observaron %d consultas simultáneas; el límite era %d", max, limite)
	}
}
```

Exécution réelle, avec le détecteur de data races actif (pour confirmer que même la mesure elle-même
n'introduit pas de data race) :

```
$ go test ./internal/revisar/... -race -v -run TestTodos_elParaleloLimitaCuantasCorrenALaVez
=== RUN   TestTodos_elParaleloLimitaCuantasCorrenALaVez
--- PASS: TestTodos_elParaleloLimitaCuantasCorrenALaVez (0.11s)
PASS
ok  	github.com/habil/revisor/internal/revisar	1.435s
```

**20 services, limite de 3, et le test confirme qu'il n'y a jamais eu plus de 3 appels à `Revisar` en cours en
même temps** — non pas parce que nous le supposons à partir du code, mais parce qu'un compteur atomique l'a
mesuré pendant l'exécution.

### 6.9 Le temps, mesuré : en série contre en parallèle

Avec `Falso` configuré pour simuler des délais réels, le test `TestTodos_respetaElTimeoutPorServicio` confirme
l'autre face de la médaille : un service qui ne répond jamais (`Colgado: true`) avec un
`Timeout: 30 * time.Millisecond` ne peut pas faire durer `Todos` plus longtemps que cela :

```
$ go test ./internal/revisar/... -run TestTodos_respetaElTimeoutPorServicio -v
=== RUN   TestTodos_respetaElTimeoutPorServicio
--- PASS: TestTodos_respetaElTimeoutPorServicio (0.03s)
```

**0,03 seconde, pas plus.** Le `context.WithTimeout` de la section 6.7, enfant du contexte général, a coupé
cette requête exactement quand il le fallait, sans que le reste du programme ait à savoir que cela s'était
produit.

---

## L'erreur que tu vas voir

Tous ceux-ci sont littéraux, provoqués exprès pour cette leçon :

| Le symptôme | Message / preuve littérale | Ce qui se passe et quoi faire |
|---|---|---|
| `main` se termine avant que les goroutines ne tournent | (aucune sortie, ou une sortie partielle et incohérente d'une exécution à l'autre) | Il manque un `sync.WaitGroup` (ou un canal) qui fasse attendre `main`. Section 6.1 |
| Les `Add`/`Done` ne concordent pas | Le programme affiche que **tout s'est bien passé d'abord**, et `panic: sync: negative WaitGroup counter` arrive APRÈS | Il manque un `wg.Add(1)` avant de lancer une goroutine, ou il y a un `Done()` de trop. `Wait()` est revenu sans rien attendre. Section 6.2 |
| Données partagées sans protection | `WARNING: DATA RACE` + `Found 2 data race(s)` + `exit status 66` (avec `-race`) ; un nombre incorrect sans explication, sans `-race` | Deux goroutines touchent la même mémoire sans synchronisation. Utilise un canal pour envoyer le résultat, ou un mutex si tu as vraiment besoin d'un compteur partagé. Sections 6.3 et 6.4 |
| Envoyer sur un canal fermé | `panic: send on closed channel` | Quelqu'un d'autre a déjà fermé le canal, ou c'est celui qui ne devait pas le faire qui l'a fermé. Ferme toujours depuis celui qui envoie, une seule fois |
| Fermer deux fois le même canal | `panic: close of closed channel` | Deux goroutines (ou deux chemins de code) essaient de fermer le même canal. Centralise la fermeture en un seul endroit |
| Personne de l'autre côté d'un canal sans tampon | `fatal error: all goroutines are asleep - deadlock!` | Le runtime de Go détecte que tout le programme s'est endormi à attendre quelque chose qui n'arrivera jamais, et le tue avec cette ligne exacte. Section 6.5 |
| Tu as oublié `defer cancelar()` | (pas d'erreur immédiate ; fuite de mémoire accumulée avec le temps) | Chaque `context.WithTimeout` non annulé garde son minuteur vivant jusqu'à ce qu'il expire de lui-même. Section 6.6 |

**Et la consigne qui résume toute cette leçon :** si tu vas écrire de la concurrence, lance `-race` dès ton
premier test, pas quand quelque chose sent mauvais — parce que, comme tu l'as vu dans la section 6.3, rien ne
sent mauvais jusqu'à ce qu'il soit trop tard.

---

## Ce qu'on fait mal

- **Lancer une goroutine par élément, sans aucune limite, contre une ressource externe.** Avec 5 services, cela
  ne se remarque pas. Avec 5 000, le programme ouvre 5 000 connexions d'un coup, et le goulot d'étranglement
  cesse d'être le service distant pour devenir ta propre machine (ou le réseau, ou le serveur lui-même qui
  reçoit une pluie de 5 000 requêtes simultanées). Le sémaphore de la section 6.7 existe exactement pour cela —
  ne laisse jamais « combien de goroutines je lance » dépendre uniquement de « combien d'éléments j'ai ».
- **Ignorer le `context` qu'on te donne, ou ne pas le propager.** Si une fonction reçoit un `ctx` et appelle
  une autre opération qui peut prendre du temps sans le lui passer, cette opération ne sera pas annulée quand
  le `ctx` original le sera — tu as deux horloges qui ne se parlent pas. Tout ce qui peut prendre du temps
  reçoit le `ctx` de celui qui l'a appelé.
- **Ne pas fermer ce qu'on ouvre.** Un `context.WithTimeout` sans son `cancelar()`, une connexion HTTP sans
  `resp.Body.Close()` (leçon 7), un fichier non fermé : chacun est une fuite différente, mais la façon de les
  éviter est la même — un `defer` juste après l'ouverture, avant d'écrire la moindre autre ligne.
- **Partager la mémoire au lieu de la communiquer « parce que c'est plus rapide à écrire ».** Une map partagée
  avec un mutex autour de tout le bloc peut sembler plus court que de monter un canal, mais il est plus facile
  d'oublier un `Lock()` à un seul point d'accès que d'oublier d'envoyer par un canal — et le premier oubli ne se
  remarque pas jusqu'à ce que le détecteur de data races (ou, pire, la production) le trouve.
- **Croire que « cela n'a jamais échoué » signifie « c'est correct ».** Une data race peut passer
  quatre-vingt-dix-neuf exécutions et échouer à la centième, ou n'échouer qu'avec plus de cœurs que ceux de ton
  ordinateur portable. La seule preuve qu'un programme concurrent n'a pas de data races est de l'exécuter avec
  `-race`, pas de le voir passer plusieurs fois sans.

---

## Exercices

1. Écris le programme de la section 6.1 (goroutines sans attente) et lance-le cinq fois d'affilée. Note
   combien de fois il a affiché quelque chose et combien de fois non.
2. Ajoute-lui un `sync.WaitGroup` correct (section 6.2) et confirme que maintenant les trois lignes
   s'affichent toujours, dans les cinq exécutions.
3. Reproduis le bug de l'`Add` manquant de la section 6.2 et colle le `panic` complet dans ton journal de bord.
4. Reproduis la data race de la section 6.3 (un compteur partagé sans protection) et lance le même programme
   avec et sans `-race`. Compare les deux nombres finaux et colle la sortie complète de `-race` dans ton journal
   de bord.
5. Reproduis le deadlock de la section 6.5 avec un canal sans tampon et sans récepteur. Confirme que tu vois le
   `fatal error: all goroutines are asleep - deadlock!`, et pas un blocage silencieux.
6. Prends ton propre `revisor` (ou celui de ce cours) et lance `TestTodos_elParaleloLimitaCuantasCorrenALaVez`
   en changeant la `limite` à 1 puis à 10. Explique, avec tes mots, pourquoi le résultat du test ne change pas
   (il passe toujours) mais le **temps** qu'il prend, si.
7. (Un peu plus difficile) Retire le sémaphore de `Todos` (laisse toutes les goroutines se lancer sans limite)
   et relance `TestTodos_elParaleloLimitaCuantasCorrenALaVez`. Confirme que maintenant il échoue, et colle le
   message d'erreur exact que donne `t.Errorf` avec le maximum qui a été effectivement observé.

### Corrigés

1-2. Il n'y a pas de nombre unique : cela dépend de ta machine. Ce qui compte, c'est la comparaison — sans
`WaitGroup`, incohérent ; avec lui, toujours les trois lignes.

3. Le programme affiche d'abord `listo` — `Wait()` est revenu immédiatement parce que le compteur n'a jamais été
   incrémenté — et le `panic: sync: negative WaitGroup counter` arrive après, avec une trace qui inclut
   `sync.(*WaitGroup).Add(...)`. Exécuté plusieurs fois d'affilée, l'ordre est toujours le même : `listo`,
   puis le panic.

4. Sans `-race`, un nombre inférieur à celui attendu (par exemple, 956 sur 1000), sans aucune erreur. Avec
   `-race`, le bloc `WARNING: DATA RACE` avec la ligne exacte du code, plus `Found N data race(s)` et
   `exit status 66`.

5. `fatal error: all goroutines are asleep - deadlock!`, avec `goroutine 1 [chan send]:` (ou `[chan
   receive]`, selon le côté où cela s'est coincé) et la ligne exacte du canal.

6. Le test passe dans les deux cas parce qu'il **mesure** le maximum réel et le compare à la limite qu'il a
   lui-même configurée (1 ou 10) — jamais à un nombre fixe attendu. Le temps total, lui, change : avec
   `limite=1` les 20 requêtes sont strictement en série (l'une attend l'autre), avec `limite=10` elles
   tournent en deux fournées de 10 au lieu de vingt d'une seule.

7. Sans sémaphore, les 20 goroutines tournent toutes en même temps, donc `maximoObservado` sera proche de 20
   (pas exactement la limite du test, qui continue de demander 3). Le `t.Errorf` dit quelque chose comme
   `se observaron 20 consultas simultáneas; el límite era 3` — le test détecte BIEN la régression, ce qui est
   justement sa raison d'être.

---

## Comment savoir que j'y suis arrivé

- [ ] J'ai vu, de mes propres yeux, un programme se terminer sans avoir attendu ses goroutines (section 6.1).
- [ ] J'ai provoqué le `panic: sync: negative WaitGroup counter` et j'ai vu que le programme a affiché `listo`
      AVANT le panic — et je peux expliquer pourquoi.
- [ ] J'ai provoqué une vraie data race et j'ai vu la différence entre l'exécuter avec et sans `-race`.
- [ ] Je peux lire un bloc de `WARNING: DATA RACE` et dire quelle ligne, quelles goroutines et quelle variable
      sont en conflit.
- [ ] J'ai provoqué le `fatal error: all goroutines are asleep - deadlock!` et je sais pourquoi Go le détecte au
      lieu de rester bloqué pour toujours.
- [ ] Je peux expliquer, avec le vrai code de `Todos`, à quoi sert chacune de ses quatre pièces : goroutines,
      canal de résultats, sémaphore et `context` par service.
- [ ] J'ai mesuré, avec un vrai test (pas de mémoire), que la limite de parallélisme du `revisor` est bien
      respectée.
- [ ] Je sais expliquer pourquoi `Falso.Contador` utilise un mutex au lieu d'un canal, et pourquoi cela ne
      contredit pas la règle générale de la section 6.3.

---

## Pour aller plus loin

1. [A Tour of Go: Concurrency](https://go.dev/tour/concurrency/1) — le parcours officiel des goroutines, des
   canaux et de `sync.WaitGroup`, interactif.
2. [Go Blog: Share Memory By Communicating](https://go.dev/blog/codelab-share) — l'article original où est
   expliquée la devise citée dans la section 6.3.
3. [Package context](https://pkg.go.dev/context) — la documentation officielle, avec les quatre cas
   d'utilisation canoniques (`WithCancel`, `WithTimeout`, `WithDeadline`, `WithValue`).
4. [Data Race Detector](https://go.dev/doc/articles/race_detector) — la documentation officielle de `-race` :
   ce qu'il détecte, ce qu'il NE détecte PAS (par exemple, il ne trouve pas les deadlocks, seulement les data
   races), et son coût en temps d'exécution.

### Termes de cette leçon

| Terme | Ce qu'il signifie |
|---|---|
| goroutine | une fonction qui tourne de façon concurrente avec le reste du programme, bien moins coûteuse qu'un thread du système d'exploitation |
| `sync.WaitGroup` | mécanisme pour attendre qu'un groupe de goroutines se termine, en comptant les `Add`/`Done` |
| data race (*situation de compétition*) | deux goroutines accèdent à la même mémoire en même temps, au moins l'une en écriture, sans synchronisation |
| canal (`chan`) | le mécanisme de Go pour qu'une goroutine envoie des données à une autre sans mémoire partagée |
| sémaphore de canal | un canal avec tampon utilisé comme quota de tours disponibles, pas pour transporter des données |
| `context` | mécanisme standard pour propager des limites de temps et l'annulation entre fonctions |
| deadlock | état dans lequel toutes les goroutines d'un programme sont endormies à attendre quelque chose qui n'arrivera jamais ; le runtime de Go le détecte et termine le programme |

---

**Précédent :** [Leçon 5 — Modules et tests](05-modulos-y-pruebas.md) ·
**Suivant :** [Leçon 7 — Le programme terminé](07-el-programa.md)
