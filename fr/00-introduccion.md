# Leçon 0 — Ce qu'est Go, d'où il vient et pourquoi tu devrais l'apprendre

**Durée :** 30 à 45 minutes de lecture, sans écrire de code pour l'instant.

**À la fin, tu seras capable de :**

- Raconter qui a créé Go, quand, et quel problème concret ils essayaient de résoudre.
- Expliquer pourquoi un nouveau langage existe alors qu'il y en avait déjà des dizaines.
- Nommer trois programmes que tu utilises ou administres et qui sont écrits en Go.
- Dire en quoi Go est bon et en quoi il **ne** l'est **pas**, pour ne pas l'utiliser là où il ne convient pas.
- Situer Go parmi les autres langages que tu connais déjà ou dont tu as entendu parler.

---

## Pourquoi c'est important

Nous sommes en 2007. Google possède l'une des plus grandes bases de code du monde, écrite en majorité en C++. Et
il a un problème très concret et très ennuyeux : **la compilation prend une éternité.**

Une seule modification dans un fichier d'en-tête pouvait déclencher une recompilation de **45 minutes ou plus**.

Pense à ce que cela signifie pour celui qui programme. Tu changes une ligne. Tu attends 45 minutes. Tu découvres
que tu t'es trompé d'une virgule. Tu changes une autre ligne. Tu attends encore 45 minutes. **Dans une journée de
travail, tu as le temps d'essayer huit choses.**

Le 21 septembre 2007, trois ingénieurs de Google se sont plantés devant un tableau pour concevoir un langage qui
ne les ferait pas attendre. Ce langage est celui dont parle tout ce cours, et le reste de cette leçon explique
comment ils sont passés de ce tableau au langage que tu vas installer dans la leçon 1.

---

## Les concepts

### 0.1 Qui l'a créé (et pourquoi c'est important)

Ce n'étaient pas trois programmeurs ordinaires :

| | Qui est-ce |
|---|---|
| **Ken Thompson** | **A créé Unix.** Et il a créé le langage B, l'ancêtre direct de C, qu'il a écrit avec Dennis Ritchie. Prix Turing —l'équivalent du Nobel en informatique— en 1983 |
| **Rob Pike** | A travaillé sur Unix et a créé Plan 9, son successeur. **A co-inventé UTF-8**, la façon dont on stocke aujourd'hui le texte dans pratiquement le monde entier, y compris ce document |
| **Robert Griesemer** | A travaillé sur le moteur JavaScript V8 (celui de Chrome et de Node.js) et sur la machine virtuelle de Java |

🔑 **Relis ça : l'un des créateurs de Go est le créateur d'Unix et le co-auteur de C.** Cinquante ans après avoir
inventé les outils sur lesquels est construit tout le logiciel moderne, la même personne a conçu Go.

Cela explique beaucoup du caractère du langage : Go ressemble à C —direct, petit, sans fioritures— mais sans les
pièges qui faisaient de C un champ de mines.

### 0.2 Les trois exigences que personne ne remplissait

Tous les trois voulaient un langage possédant **trois** propriétés à la fois :

1. **Qu'il compile vite.** Le problème d'origine.
2. **Qu'il s'exécute vite.** Google fait tourner du logiciel sur des milliers de machines : la vitesse coûte de
   l'argent réel.
3. **Qu'il soit facile à programmer.** Qu'une personne nouvelle dans l'équipe puisse lire le code des autres et
   le comprendre dès le premier jour.

Ils ont passé en revue les langages existants. Et voici la découverte intéressante : **ils ont trouvé des langages
qui remplissaient deux des trois critères, mais aucun qui remplisse les trois.**

| | Compile vite | S'exécute vite | Facile à programmer |
|---|---|---|---|
| **C / C++** | ❌ | ✅ | ❌ |
| **Java** | 🟡 | ✅ | 🟡 |
| **Python** | ✅ (ne compile pas) | ❌ | ✅ |
| **Go** | ✅ | ✅ | ✅ |

> [!NOTE]
> 🔧 **Observation de génie logiciel 0.1**
> Un langage de programmation n'est pas une question de goût : c'est un outil fait de compromis. Chacun a
> sacrifié quelque chose pour gagner autre chose. Comprendre **ce qu'a sacrifié** celui que tu utilises, c'est la
> différence entre programmer avec discernement et programmer par cœur.

### 0.3 Comment la simplicité a été gagnée : en retirant des choses

Voici la décision la plus contre-intuitive de Go, et celle qui le distingue le plus.

Presque tous les langages grossissent : chaque version ajoute des fonctionnalités. C++ en a accumulé pendant des
décennies, et le résultat est si vaste que **personne ne le connaît en entier**. Des gens programment en C++
depuis vingt ans et continuent à découvrir des recoins qu'ils ne connaissaient pas.

**Go a fait le contraire : il a décidé ce qu'il laisserait de côté.** Il n'a pas :

- **De classes ni d'héritage.** La façon de réutiliser du code est autre, et tu la verras dans la leçon 3.
- **D'exceptions** (`try / catch`). Les erreurs se gèrent d'une autre manière, et tu la verras dans la leçon 2.
- **De surcharge d'opérateurs.** `+` veut dire additionner, toujours, et on ne peut pas le changer.
- **De constructeurs ni de destructeurs.**
- **Une multitude de façons de faire la même chose.** Il y a presque toujours **une** manière idiomatique.

⚠️ **Ça va t'agacer un jour.** Tu voudras faire quelque chose qui se fait en une ligne dans un autre langage, et
en Go ça t'en prendra cinq. La contrepartie est énorme et se voit au bout de quelques mois : **tu peux lire du
code Go écrit par n'importe qui et le comprendre.** Il n'y a pas de dialectes, pas d'astuces, pas besoin
d'apprendre comment programme chaque équipe.

> 🔑 **La citation qui résume la philosophie**, de Rob Pike :
> *« La clarté vaut mieux que l'astuce. »*

### 0.4 L'autre raison : les ordinateurs ont cessé d'accélérer

Il y a un deuxième motif derrière Go, et il tient au matériel.

Jusqu'aux environs de 2005, chaque année sortaient des processeurs plus rapides et ton programme s'exécutait plus
vite **sans que tu touches à rien**. C'est fini : les processeurs ont cessé d'augmenter en vitesse et se sont mis à
multiplier les **cœurs**. Ton portable n'a pas un processeur ultra-rapide : il en a 4, 8 ou 16, modestes.

**Le problème :** un programme normal utilise **un** cœur. Les quinze autres sont là, à regarder.

Pour les exploiter, il faut écrire des programmes qui font plusieurs choses à la fois, et cela —la
**concurrence**— était traditionnellement l'une des tâches les plus difficiles et les plus sujettes aux erreurs de
toute la programmation.

**Go a été conçu autour de ça.** Faire en sorte que deux choses se produisent en même temps en Go, c'est
littéralement écrire un mot :

```go
go revisarServicio()
```

Ce mot `go` —celui qui donne son nom au langage— est la fonctionnalité qui l'a rendu célèbre. Tu l'apprendras dans
la leçon 6, et quand tu y seras, tu comprendras pourquoi tant de gens sont passés à ce langage.

### 0.5 La chronologie

| Date | Ce qui s'est passé |
|---|---|
| **21-sep-2007** | Griesemer, Pike et Thompson conçoivent Go sur un tableau |
| mi-2008 | il existe déjà un compilateur qui fonctionne |
| **10-nov-2009** | Google l'annonce publiquement et le publie comme logiciel libre |
| **28-mar-2012** | **Go 1.0**, la première version stable |
| 2012 et après | deux versions majeures par an, chaque février et chaque août |
| aujourd'hui | Go 1.27 |

🔑 **Et un point à connaître : la promesse de compatibilité de Go 1.** Depuis 2012, l'équipe de Go a promis qu'**un
programme écrit pour Go 1.0 compile toujours aujourd'hui**, quatorze ans plus tard. Ils ont tenu parole.

Ce n'est pas normal. Dans beaucoup de langages, du code vieux de cinq ans ne compile plus. **En Go, ce que tu
apprends aujourd'hui te servira dans dix ans**, ce qui rend le temps que tu vas investir beaucoup plus rentable.

### 0.6 Ce qui est écrit en Go (tu l'utilises probablement déjà)

Ce n'est pas une liste de propagande : c'est pour que tu voies **où** ce langage est utilisé, car cela en dit long
sur son utilité.

| | Ce que c'est |
|---|---|
| **Docker** | l'outil qui empaquette les applications dans des conteneurs |
| **Kubernetes** | le système qui administre des milliers de conteneurs sur des serveurs. La moitié de l'industrie l'utilise |
| **Terraform** | crée de l'infrastructure dans le cloud en écrivant des fichiers |
| **Prometheus** · **Grafana** | collectent et représentent graphiquement les métriques des serveurs |
| **Traefik** · **Caddy** | serveurs web et répartiteurs de charge |
| **CockroachDB** · **InfluxDB** | bases de données |
| **Hugo** | générateur de sites web statiques, réputé pour sa vitesse |
| **ngrok**, **rclone**, **gh** (le CLI de GitHub) | outils en ligne de commande |

🔑 **Tu vois le motif ?** Presque tout ce sont des **outils d'infrastructure** : des choses qui tournent sur des
serveurs, qui doivent être rapides, qui se déploient comme un fichier unique et qui gèrent beaucoup de connexions
à la fois. **C'est là que Go l'emporte**, et ce n'est pas un hasard : c'est exactement le problème que Google
avait.

### 0.7 Alors, pourquoi devrais-tu l'apprendre, toi ?

Quatre raisons honnêtes :

**1. Il s'apprend vite.** La spécification de Go se lit en une après-midi. Celle de C++ compte plus de 1 800
pages. **Tu pourras écrire des programmes utiles en quelques semaines, pas en quelques années**, et cela compte
beaucoup quand on débute.

**2. Il t'enseigne des choses utiles dans n'importe quel langage.** Étant compilé et à typage strict, Go t'oblige
à penser aux types, à la mémoire et aux erreurs. Ces concepts **se transfèrent** : si tu programmes ensuite en
Java, C# ou Rust, tu les as déjà.

**3. Il y a du travail.** Tout ce qui tourne sur les serveurs modernes contient du Go. Si l'infrastructure, le
cloud, le DevOps ou le backend t'intéressent, c'est l'un des paris les plus sûrs.

**4. Ce que tu apprendras ne périmera pas.** Grâce à la promesse de compatibilité de Go 1.

### 0.8 Ce que tu vas construire dans ce cours

Un programme en ligne de commande appelé **`revisor`** : il reçoit une liste de services, les interroge **tous en
même temps**, et produit un rapport.

```
$ revisor --config servicios.txt --formato tabla
SERVICIO    ESTADO    TIEMPO  DETALLE
catalogo    OK         142ms  200
inventario  LENTO       2.3s  200
pagos       OK          87ms  200
reportes    FALLA          —  connection refused
```

Ça a l'air petit. **Ça ne l'est pas.** Pour l'écrire correctement, tu auras besoin de tout : types, structs,
erreurs, interfaces, listes, concurrence, HTTP, fichiers de configuration, tests et compilation d'un exécutable.

Et il va grandir avec toi : **chaque leçon ajoute une pièce.** À la fin, tu auras un programme que tu pourrais
vraiment utiliser, pas un exercice de manuel.

---

## L'erreur que tu vas voir

Cette leçon n'a pas encore de code propre, donc il n'y a pas de message de compilateur à te montrer. Mais il y a
bien une erreur de **raisonnement** que presque tout le monde commet en lisant la section 0.3, et il vaut la peine
de te la signaler à l'avance pour qu'elle ne te freine pas quand elle t'arrivera.

**L'erreur :** lire la liste de ce que Go **n'**a **pas** —classes, héritage, exceptions, surcharge
d'opérateurs— et conclure *« alors c'est un langage pauvre, il lui manque des choses dont j'ai besoin ».*

**Pourquoi c'est une erreur :** elle confond *avoir moins d'outils* et *pouvoir résoudre moins de problèmes*. Go ne
t'a pas retiré la capacité de réutiliser du code ni de gérer les erreurs : il t'a donné **une** façon de le faire
au lieu de dix, et tu apprendras cette façon dans les leçons 2 et 3. La frustration est réelle et normale —toi
aussi tu la ressentiras la première semaine—, mais elle se résout en programmant, pas en évitant le langage. Quand
tu auras terminé la leçon 3, tu pourras relire la section 0.3 et voir pourquoi chaque ligne de cette liste est une
décision, pas un manque.

---

## Ce qu'on fait mal

Aussi important que de savoir à quoi sert un outil, c'est de savoir à quoi il ne sert pas. **Aucun langage n'est
bon pour tout**, et celui qui te dit le contraire essaie de te vendre quelque chose. Utiliser Go là où il ne
convient pas est le premier anti-patron de ce cours, avant même d'avoir écrit une seule ligne :

| Ce n'est pas la meilleure option pour | Ce qu'on utilise à la place | Pourquoi |
|---|---|---|
| Applications iPhone ou Android | Swift, Kotlin | les plateformes sont faites pour ceux-là |
| Pages web (ce qui s'exécute dans le navigateur) | JavaScript, TypeScript | le navigateur n'exécute que du JavaScript |
| Science des données, intelligence artificielle | Python | toutes les bibliothèques du monde sont là |
| Gros jeux vidéo | C++, C# | ils ont besoin d'un contrôle absolu de la mémoire et du matériel |
| Systèmes où une microseconde compte | C, C++, **Rust** | Go a un ramasse-miettes qui met parfois le programme en pause |

> [!NOTE]
> 🔧 **Observation de génie logiciel 0.2**
> La dernière ligne est la raison pour laquelle ce cours a un **deuxième cours, sur Rust**. Rust résout
> exactement cela : la vitesse de C sans ramasse-miettes et sans les erreurs mémoire de C. Ce sont des outils pour
> des problèmes différents, et tu apprendras les deux.

---

## Exercices

### Questions de révision

**0.1** Quel problème concret a motivé la création de Go, et en quelle année ?

**0.2** Nomme les trois créateurs de Go et dis pourquoi il est pertinent de savoir qui est Ken Thompson.

**0.3** Quelles étaient les trois exigences qu'ils recherchaient et pourquoi aucun langage existant ne leur
convenait-il ?

**0.4** Cite trois choses que Go **n'**a **pas**, à dessein. Que gagne-t-on à les retirer ?

**0.5** Qu'est-ce qui a changé dans le matériel vers 2005 et quel rapport avec Go ?

**0.6** Qu'est-ce que la promesse de compatibilité de Go 1 et pourquoi est-elle avantageuse pour toi ?

**0.7** Nomme trois programmes écrits en Go et dis ce qu'ils ont en commun.

**0.8** Donne deux cas où tu **n'**utiliserais **pas** Go, et dis ce que tu utiliserais.

### Pour aller plus loin

**0.9** Cherche sur internet la conférence ou l'article original *« Go at Google: Language Design in the Service
of Software Engineering »* de Rob Pike. Lis l'introduction et note **une** raison de conception qui ne figure pas
dans cette leçon.

**0.10** Choisis **deux** des programmes de la section 0.6 que tu ne connais pas. Découvre en une phrase ce que
fait chacun et note-le dans ton journal de bord.

**0.11** Cherche la table des matières de la spécification du langage Go (*The Go Programming Language
Specification*) et compte combien de pages ou de sections elle comporte. Compare avec la norme de C++. Note les
deux chiffres : c'est la façon la plus concrète de voir ce que « simple » veut dire.

**0.12 (À réfléchir, sans réponse correcte)** Go a retiré les exceptions, les classes et l'héritage —des choses
que d'autres langages jugent indispensables. Te vient-il une raison pour laquelle **retirer** une fonctionnalité
puisse améliorer un langage ? Écris ton opinion dans ton journal de bord **avant** de commencer le cours, et
relis-la à la fin. Il est intéressant de voir si elle a changé.

### Corrigés

**0.1** Les temps de compilation de C++ chez Google : une modification dans un en-tête pouvait coûter **45
minutes** de recompilation. Ils ont commencé le **21 septembre 2007**.

**0.2** Robert Griesemer, Rob Pike et **Ken Thompson**. Thompson a **créé Unix** et le langage B, ancêtre de C,
qu'il a écrit avec Dennis Ritchie. C'est-à-dire : l'un des auteurs des outils sur lesquels s'est construit le
logiciel moderne a aussi conçu Go, cinquante ans plus tard.

**0.3** Compilation rapide, exécution rapide et facilité de programmation. Les langages existants en remplissaient
**deux sur trois** : C++ était rapide à l'exécution mais lent à compiler et difficile ; Python était facile mais
lent à l'exécution.

**0.4** Classes et héritage, exceptions, surcharge d'opérateurs, constructeurs. **On gagne en lisibilité** : du
code Go écrit par n'importe qui peut être lu et compris, parce qu'il n'y a ni dialectes ni multiples façons de
faire la même chose.

**0.5** Les processeurs ont cessé d'accélérer et se sont mis à multiplier les **cœurs**. Un programme normal n'en
utilise qu'un seul, il fallait donc un langage où écrire des programmes concurrents soit facile. D'où le mot
`go`.

**0.6** La promesse qu'un programme écrit pour Go 1.0 (2012) **compile toujours aujourd'hui**. Elle t'avantage
parce que ce que tu apprends ne périme pas et que le code que tu écris continuera de fonctionner pendant des
années.

**0.7** Docker, Kubernetes, Terraform, Prometheus, Grafana, Hugo… Tous sont des **outils d'infrastructure** : ils
tournent sur des serveurs, se distribuent comme un exécutable unique et gèrent beaucoup de connexions à la fois.

**0.8** Applications mobiles (Swift/Kotlin), code de navigateur (JavaScript), science des données (Python), gros
jeux vidéo (C++), systèmes à temps réel strict (C ou Rust).

**0.9-0.12** Elles n'ont pas de réponse correcte unique : ce sont des exercices de recherche et de réflexion
personnelle. Compare ce que tu as noté avec un camarade ou dans le journal de bord du cours, pas avec une réponse
fixe.

---

## Comment savoir que j'y suis arrivé

Il n'y a pas de code à compiler dans cette leçon, donc la liste de contrôle porte sur la compréhension. Coche
chaque point seulement si tu peux le faire **sans relire le texte** :

- ☐ Tu expliques avec tes mots, en moins d'une minute, pourquoi Go est né et quel problème il résolvait.
- ☐ Tu nommes les trois créateurs et tu dis pourquoi le fait que Ken Thompson soit l'un d'eux n'est pas un détail
  accessoire.
- ☐ Tu dis ce que sont la compilation rapide, l'exécution rapide et la facilité de programmation, et pourquoi
  aucun langage antérieur à Go n'avait les trois.
- ☐ Tu nommes, de mémoire, au moins trois choses que Go a décidé de ne pas avoir, et tu expliques ce qu'on gagne à
  les retirer.
- ☐ Tu expliques la relation entre les cœurs d'un processeur, la concurrence et le mot `go`.
- ☐ Tu nommes trois programmes écrits en Go et tu dis ce qu'ils ont en commun.
- ☐ Tu donnes un exemple de problème pour lequel tu n'utiliserais **pas** Go, et tu dis ce que tu utiliserais à la
  place.
- ☐ Tu as résolu les huit questions de révision sans voir les corrigés avant de comparer.

S'il t'en manque un, ne passe pas encore à la leçon 1 : relis la section correspondante. Tout ce qui vient
ensuite s'appuie sur ceci.

---

## Résumé

- Go est né en **2007** chez Google de la frustration face aux **temps de compilation de C++** : 45 minutes pour
  une modification dans un en-tête.
- Il a été créé par **Robert Griesemer, Rob Pike et Ken Thompson** ; Thompson a créé Unix et co-écrit C.
- Ils cherchaient **trois** propriétés réunies —compiler vite, s'exécuter vite, être facile— qu'aucun langage de
  l'époque ne remplissait à la fois.
- Go a atteint la simplicité en **retirant** des choses : pas de classes, pas d'héritage, pas d'exceptions, pas de
  surcharge d'opérateurs.
- La seconde motivation était le matériel : les processeurs ont multiplié les **cœurs** au lieu d'accélérer, et
  les exploiter exigeait que la **concurrence** soit facile. D'où le mot `go`.
- Il a été annoncé en **2009** et la version **1.0** est sortie en **2012** ; il en est aujourd'hui à 1.27, avec
  deux versions majeures par an.
- La **promesse de compatibilité de Go 1** garantit que le code de 2012 compile toujours : ce que tu apprends ne
  périme pas.
- Go domine l'**infrastructure** : Docker, Kubernetes, Terraform, Prometheus, Grafana.
- Ce **n'**est **pas** la meilleure option pour le mobile, le navigateur, la science des données, les jeux vidéo
  ni le temps réel strict.

---

## Pour aller plus loin

1. **La spécification du langage** — *The Go Programming Language Specification*, sur `go.dev/ref/spec`.
   C'est la source officielle et, comme tu l'as vu dans l'exercice 0.11, elle se lit en une après-midi. Inutile de
   tout comprendre aujourd'hui ; il suffit de savoir qu'elle existe et qu'elle est courte.
2. **Le blog officiel de Go** — `go.dev/blog`, où l'équipe elle-même publie les nouveautés de chaque version et
   des articles de conception.
3. **Rob Pike, « Go at Google: Language Design in the Service of Software Engineering »** — la conférence/article
   original où l'un des créateurs explique les décisions de conception avec le contexte de Google. C'est la
   lecture de l'exercice 0.9, et la plus directe pour comprendre le **pourquoi** de chaque décision que tu as vue
   dans cette leçon.

### Termes de cette leçon

| | |
|---|---|
| **compiler** | traduire du code source en instructions machine |
| **concurrence** | qu'un programme fasse plusieurs choses à la fois |
| **en-tête (*header*)** | en C/C++, fichier de déclarations que d'autres fichiers incluent |
| **Go 1 (promesse de compatibilité)** | engagement que le code ancien continue de compiler |
| **cœur (*core*)** | chaque processeur indépendant au sein d'une même puce |
| **prix Turing** | la plus haute distinction en informatique |
| **ramasse-miettes** | partie du langage qui libère la mémoire automatiquement |
| **UTF-8** | forme standard de représenter le texte, co-inventée par Rob Pike et Ken Thompson |

---

**Suivant :** [Leçon 1 — Installer Go sur ton Linux Mint](01-instalacion.md)
