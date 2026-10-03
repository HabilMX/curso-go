# Cours de Go — de zéro à l'écriture d'un bon programme

**Par Dorian Chávez, fondateur de Hábil et architecte d'intégration.**

**Pour qui :** quelqu'un qui débute — première année d'université, ou tout juste sorti de l'école. **Nous ne
tenons rien pour acquis** : ni que tu as Go installé, ni que tu sais ce qu'est un pointeur, ni ce que signifie
compiler. Chaque concept est expliqué quand il apparaît, et on explique **pourquoi** il existe, pas seulement
comment il s'écrit.

**Ce que tu sauras faire à la fin :** écrire un programme complet, comprendre ce que tu as écrit, et pouvoir
l'expliquer à quelqu'un d'autre. Ce dernier point est la vraie épreuve.

**Ce dont tu as besoin avant de commencer :** un ordinateur avec Linux Mint et savoir ouvrir un terminal. Rien
de plus. La [Leçon 1](01-instalacion.md) installe tout à partir de zéro.

## Le projet que tu vas construire

Un programme en ligne de commande appelé **`revisor`** : il reçoit une liste de services, les interroge **tous
en même temps**, et produit un rapport sous forme de tableau ou en JSON.

    $ revisor --config servicios.txt --formato tabla
    SERVICIO    ESTADO    TIEMPO  DETALLE
    catalogo    OK         142ms  200
    inventario  LENTO       2.3s  200
    pagos       OK          87ms  200
    reportes    FALLA          —  connection refused

Il a l'air petit et il ne l'est pas : pour bien l'écrire, tu as besoin de structs, d'interfaces, d'erreurs, de
concurrence, de `context`, de HTTP, de JSON, d'options de ligne de commande, de tests et d'un binaire empaqueté.

**Et c'est le squelette de ce que les gens écrivent vraiment en Go.** Selon la
[Go Developer Survey 2025](https://go.dev/blog/survey2025) —l'enquête officielle de l'équipe de Go, menée en
septembre 2025 avec **5 379 réponses**—, les deux types de projets les plus construits sont :

| Ce qu'ils construisent | |
|---|---|
| Outils en ligne de commande | **74 %** |
| Services API/RPC | **73 %** |
| Bibliothèques ou frameworks | 49 % |

Le `revisor` relève du premier, et dans la Leçon 7 tu y ajoutes le second.

⚠️ **Et tu vas l'écrire une nouvelle fois en Rust**, dans le cours frère. Écrire le même programme dans les
deux langages, c'est ce qui t'apprend vraiment en quoi ils diffèrent ; lire des comparaisons ne sert à rien.

## Les huit leçons

| Leçon | | Ce que tu construis | Ce que tu apprends |
|---|---|---|---|
| 0 | [Ce qu'est Go](00-introduccion.md) | rien encore (lecture) | d'où vient Go, à quoi il sert et à quoi il ne sert pas |
| 1 | [Installation](01-instalacion.md) | l'environnement et ton premier programme | `GOROOT`, `GOPATH`, le `PATH`, et les vraies erreurs de l'installation de Go |
| 2 | [Variables, fonctions et types](02-fundamentos.md) | les fonctions de base du `revisor` | le compilateur, les variables, les types de base, les fonctions avec `(resultado, error)` |
| 3 | [Structs, méthodes, erreurs et interfaces](03-errores-interfaces.md) | le `revisor` avec des structs et des interfaces | structs, méthodes, pointeurs, erreurs comme valeurs, interfaces |
| 4 | [Collections : slices et maps](04-colecciones.md) | la liste des services et le rapport | slices, maps, `range`, ordre stable |
| 5 | [Modules et tests](05-modulos-y-pruebas.md) | le vrai projet, avec des tests | `go mod`, table de cas, `-race`, couverture |
| 6 | [Concurrence](06-concurrencia.md) | **qu'il vérifie tout en même temps** | goroutines, canaux, `context`, `WaitGroup`, sémaphore |
| 7 | [Le programme terminé](07-el-programa.md) | CLI, HTTP, JSON et binaire | `net/http`, `flag`, `encoding/json`, compilation croisée |

## L'ordre n'est pas l'ordre évident, et c'est exprès

Les structs, **les erreurs et les interfaces viennent avant** les slices, les maps et les pointeurs. Cela
paraît étrange —presque tous les cours mettent les collections en premier— mais ce n'est pas un caprice : c'est
l'ordre qu'utilisent deux des cours de Go qui comptent le plus d'apprenants, celui de
[Boot.dev / freeCodeCamp](https://www.boot.dev/courses/learn-golang) (structs → interfaces → erreurs aux
places 5, 6 et 7, **avant** les slices, les maps et les pointeurs) et celui de
[Todd McLeod sur Udemy](https://www.udemy.com/course/learn-how-to-code/).

La raison est bonne : **en Go, modéliser les données et gérer les erreurs, C'EST le langage.** Un programme avec
des slices parfaits et des erreurs ignorées n'est pas du Go ; c'est du C avec une autre syntaxe.

**La concurrence vient dans la Leçon 6, pas dans la 3.** C'est la raison pour laquelle Go existe, mais tu as
besoin des structs, des erreurs et des interfaces pour que les exemples ne soient pas des jouets. Et elle se fait
en trois passes —les goroutines, puis les canaux, puis appliquée au projet— parce qu'en une seule, cela ne
prend pas.

## Comment l'utiliser

- **Une séance de 90 minutes par semaine**, ou deux de 45. Moins, cela ne prend pas ; plus de deux heures
  d'affilée, cela s'oublie.
- 🔴 **Écris le code à la main, toujours.** Copier-coller produit des fichiers, pas des connaissances. Tes
  doigts apprennent des choses que tes yeux n'apprennent pas.
- **Casse le code exprès.** Dans plusieurs exercices pratiques, nous allons te demander de provoquer une erreur
  et de lire ce que dit le compilateur. **Ce n'est pas du remplissage :** apprendre à lire les erreurs, c'est la
  moitié de savoir programmer.
- **Ne passe pas à la leçon suivante avec des doutes non résolus.** Note-les dans le journal de bord et
  pose-les. Tout ce qui suit s'appuie sur ce qui précède.
- **Chaque leçon se termine par un « comment savoir que j'y suis arrivé » mesurable** : cela compile et
  passe, ou non.
- **[`bitacora.md`](bitacora.md) est à toi** : doutes, faux pas et ce qui t'a surpris. C'est ce qui fait de
  ceci ton cours.

## Les trois sources qui valent la peine, dans cet ordre

1. **[A Tour of Go](https://go.dev/tour/)** — officiel, interactif, 2-3 h. **En entier, pas survolé.**
2. **[Go by Example](https://gobyexample.com/)** — chaque concept avec un code minimal. Comme référence.
3. **[Effective Go](https://go.dev/doc/effective_go)** — la philosophie. **Lis-le dans la Leçon 3, pas à la
   fin.**

⚠️ **Apprends la bibliothèque standard avant n'importe quel framework.** En Go, la bibliothèque standard suffit
pour presque tout, et c'est une décision de conception, pas une lacune. Celui qui commence avec Gin ou Echo
apprend le framework et pas le langage.

---

## Qui écrit ceci

Ce cours est écrit par **Dorian Chávez**, CEO et consultant principal de **Hábil**, cabinet mexicain
d'ingénierie logicielle. Ingénieur en systèmes informatiques de l'IPN (ESCOM), avec 26 ans passés à construire
et à exploiter des systèmes en production.

Nous le publions en accès libre, sous licence [Creative Commons](../LICENSE.md), parce que celui qui explique
bien construit bien — et parce qu'un cours qu'on ne peut ni copier, ni traduire, ni améliorer ne sert pas à
grand-chose.

**Si tu trouves une erreur, une explication confuse ou un chiffre qui ne colle pas, dis-le.** Ce matériel est
fait pour être corrigé. Avant de publier chaque version, nous vérifions que tout le code du cours fonctionne et
que ce que montre chaque leçon est identique au vrai programme, pour qu'aucune leçon n'enseigne quelque chose
qui a déjà changé.

### Si cela t'a été utile

- **Tu veux travailler avec des gens qui écrivent ainsi ?** Écris-nous.
- **Ton équipe adopte Go ou Rust ?** Écris-nous aussi.

*(Les liens de contact sont ajoutés par le site lors de la publication de cette page.)*
