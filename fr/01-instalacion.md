# Leçon 1 — Installer Go sur ton Linux Mint et ton premier programme

**Durée :** 45 à 60 minutes (ou 2 séances de 30). C'est la seule leçon qui n'enseigne pas Go en tant que
langage : elle enseigne l'outil et le terrain où tu vas travailler — et ce terrain a plus de pièges qu'il n'y
paraît à première vue.

**À la fin, tu seras capable de :**

- Installer la version officielle et actuelle de Go sur Linux Mint, sans dépendre du paquet du système, et
  expliquer pourquoi ce paquet est obsolète par conception, et non par négligence.
- Expliquer ce que sont `GOROOT`, `GOPATH` et le `PATH`, et pourquoi Go ne fonctionne pas tant que tu n'as pas
  touché au second.
- Reconnaître la différence entre un terminal qui lit ta configuration et un qui ne la lit pas —la vraie cause du
  « je l'ai ajouté et ça ne marche quand même pas »— et savoir lequel est le tien.
- Créer un projet avec `go mod init`, écrire un programme et l'exécuter avec `go run`.
- Produire un binaire avec `go build` et expliquer, preuve réelle à l'appui, pourquoi il tourne sur une autre
  machine sans que Go y soit installé.
- Utiliser `go fmt` pour ne jamais débattre du style du code.
- Reconnaître, mot pour mot, les messages d'erreur les plus courants de cette étape —de compilateur, de
  permissions et de chemin— et savoir quoi faire pour chacun.

---

## Pourquoi c'est important

La première chose qu'un vieux tutoriel va te suggérer, c'est `apt install golang`, et c'est une erreur. Ce n'est
pas une erreur mineure du genre « version un peu ancienne » : c'est le genre d'erreur qui te fait perdre une
après-midi entière sans comprendre pourquoi, des semaines après l'installation, parce que le symptôme n'apparaît
pas à l'installation — il apparaît quand tu programmes déjà et que quelque chose qui « devrait marcher » ne marche
pas.

Ce n'est pas une opinion, c'est une mesure. Le 30-sep-2026 j'ai mesuré, dans un conteneur ayant la même base que
Linux Mint 22.3 (Ubuntu 24.04 « Noble Numbat », d'où Mint tire ses paquets) :

```
$ cat /etc/os-release | grep VERSION
VERSION="24.04.5 LTS (Noble Numbat)"

$ apt-cache policy golang-go
golang-go:
  Installed: (none)
  Candidate: 2:1.22~2build1
  Version table:
     2:1.22~2build1 500
        500 http://ports.ubuntu.com/ubuntu-ports noble/main arm64 Packages

$ curl -s 'https://go.dev/VERSION?m=text'
go1.27.1
```

**Le paquet du système propose la 1.22. La version officielle actuelle, mesurée au même instant, est la 1.27.1.**
Et ce n'est pas une négligence que quelqu'un va corriger : c'est la politique d'Ubuntu et de Mint.

**Pourquoi se fige-t-elle ainsi, à dessein ?** Ubuntu 24.04 est une version LTS (*Long Term Support*) : le jour de
sa sortie, ses paquets principaux sont **figés** et ne reçoivent que des correctifs de sécurité, jamais de
nouvelles versions du programme. C'est une décision correcte pour un serveur qui ne doit pas changer de
comportement du seul fait d'installer des mises à jour — mais cela signifie que le compilateur d'un langage qui
publie de nouvelles versions deux fois par an reste fixé à celle qui existait à la publication d'Ubuntu 24.04
(avril 2024), pour toujours, tant que cette version d'Ubuntu existe. Linux Mint n'a pas ses propres paquets Go :
il hérite littéralement de ceux d'Ubuntu, donc il hérite aussi du gel.

**Pourquoi cela t'importe-t-il, concrètement ?** Parce que tu vas chercher des choses sur internet, les exemples
vont utiliser des fonctions ou des comportements que ton Go 1.22 n'a pas, et les erreurs ne diront pas *« il te
manque une version »* : elles diront des choses qui n'ont aucun sens par rapport à ce que tu vois à l'écran. C'est
le genre de problème qui ressemble à une erreur de ta part alors que c'est en réalité un décalage d'outil — et il
t'arrive **avant** d'écrire ta première ligne de code, c'est pourquoi cette leçon existe avant toutes les autres.

🔑 **La bonne méthode, et c'est celle que recommande le projet Go lui-même :** télécharger le paquet officiel
directement depuis `go.dev`, pas depuis le gestionnaire de paquets de ta distribution. C'est ce qu'installe cette
leçon, pas à pas, et ce que nous allons vérifier à chaque étape avec la sortie réelle des commandes.

---

## Les concepts

### 1.1 Ce que sont GOROOT et GOPATH, et pourquoi ils comptaient plus avant que maintenant

Avant d'installer, il vaut la peine de savoir ce que tu auras ensuite, parce que les noms `GOROOT` et `GOPATH`
apparaissent dans presque toute erreur d'installation que tu chercheras sur internet, et beaucoup de réponses sont
écrites pour une version de Go d'il y a dix ans.

- **`GOROOT`** est le dossier où vit **Go lui-même** : le compilateur, `gofmt`, la bibliothèque standard. C'est le
  dossier que tu vas créer à l'étape 1.4 (`/usr/local/go`). Tu n'y touches jamais à la main.
- **`GOPATH`** est le dossier où Go range **tes affaires** : les paquets téléchargés depuis internet, et —dans les
  anciennes versions de Go, avant 2019— c'était aussi là qu'il fallait mettre **tout** ton code, sans exception,
  dans une structure fixe (`$GOPATH/src/github.com/tu-usuario/tu-proyecto`). Si tu vois un jour un tutoriel qui te
  demande de créer cette structure de dossiers, il date de cette époque.

Aujourd'hui tu n'y touches presque plus parce que, depuis Go 1.11 (2018), il existe les **modules** (`go.mod`, que
tu vas créer à l'étape 1.6) : ton projet peut vivre dans n'importe quel dossier, avec le nom que tu veux, et Go n'a
plus besoin que tu suives une structure de dossiers imposée. `GOPATH` existe toujours, mais désormais seulement
comme cache de paquets téléchargés, et non comme l'unique endroit où ton code peut vivre.

Vérifie-le avec `go env`, qui t'affiche la configuration en vigueur de ton installation (tu peux le lancer
**après** l'installation, à l'étape 1.5) :

```bash
go env GOROOT GOPATH GOBIN
```

Dans une installation propre, toute fraîche, tu verras quelque chose comme :

```
/usr/local/go
/home/tu-usuario/go
/home/tu-usuario/go/bin
```

Aucun de ces trois dossiers, tu ne l'as créé à la main : le premier est créé par l'étape d'installation (1.4), les
deux autres sont décidés par Go tout seul, avec des valeurs par défaut raisonnables.

### 1.2 Le PATH : ce que c'est, et pourquoi « déjà installé » ne veut pas dire « ça marche déjà »

Quand tu tapes une commande dans le terminal, par exemple `go`, le terminal ne sait pas par magie où se trouve ce
programme : il passe en revue, une par une, une liste de dossiers stockée dans une variable appelée **PATH**, et
utilise le premier programme qu'il trouve portant ce nom. Si aucun dossier de la liste ne contient un programme
appelé `go`, il répond par une erreur — et cette erreur est littérale, pas approximative :

```bash
$ go version
bash: go: command not found
```

Je l'ai vérifié dans un conteneur avec Go **déjà extrait** dans `/usr/local/go`, avant de toucher au PATH : le
programme existe sur le disque, mais le terminal ne le trouve pas parce qu'il ne sait pas où chercher.
**« Installé » et « dans le PATH » sont deux choses différentes**, et la confusion entre les deux est la source
de presque tous les faux pas de cette leçon.

Cela explique aussi pourquoi le message change un peu selon le terminal (`zsh` au lieu de `bash`, que tu verras si
tu utilises macOS pour suivre le cours pendant que tu pratiques, ou si tu as changé le terminal par défaut de ton
Mint) :

```bash
$ go version
zsh: command not found: go
```

Même problème, même mécanisme, **ordre des mots différent** : identifie quel est ton cas avec `echo $SHELL` avant
de chercher l'erreur sur internet, parce que chercher le mauvais message te mènera à des réponses pour le mauvais
shell.

### 1.3 Découvre quelle est la dernière version

Ne la recopie pas d'ici : ce document vieillit comme les dépôts de Mint, et tu as déjà vu dans la section
précédente combien cela peut peser sur un outil de rester figé dans le temps.

```bash
curl -s 'https://go.dev/VERSION?m=text' | head -1
```

Il te répondra quelque chose comme `go1.27.1`. **C'est celle-là que tu vas installer**, quelle qu'elle soit au
moment où tu le fais.

### 1.4 Télécharge-la et installe-la

Remplace `go1.27.1` par ce que t'a dit la commande précédente. `linux-amd64` est le bon choix pour un PC ou un
portable normal (si ta machine était ARM, ce serait `linux-arm64` ; pour le savoir : `dpkg --print-architecture`).

```bash
cd /tmp
wget https://go.dev/dl/go1.27.1.linux-amd64.tar.gz
```

Et maintenant l'installation :

```bash
sudo rm -rf /usr/local/go                          # borra una instalación anterior, si había
sudo tar -C /usr/local -xzf go1.27.1.linux-amd64.tar.gz
```

**Ce que tu viens de faire**, parce qu'il vaut mieux le comprendre que le copier :

- `tar` décompresse l'archive.
- `-C /usr/local` lui dit *« fais-le dans ce dossier »*.
- `-xzf` signifie « extrais » (`x`), « est compressé avec gzip » (`z`), « à partir de ce fichier » (`f`).
- Le résultat est un dossier `/usr/local/go` avec tout Go dedans — le `GOROOT` de la section 1.1.

**Pourquoi `sudo`, si tu n'en avais jamais eu besoin pour installer quelque chose avec un gestionnaire de
paquets ?** Parce que `/usr/local` est un dossier du système, pas de ton utilisateur, et sous Linux y écrire
demande des droits d'administrateur. Sans `sudo`, voici exactement ce que tu vas voir —je l'ai provoqué
volontairement, sans `sudo`, pour que tu voies le vrai message et non un message inventé— :

```
$ tar -C /usr/local -xzf go1.27.1.linux-amd64.tar.gz
tar: go: Cannot mkdir: Permission denied
tar: go/VERSION: Cannot open: No such file or directory
tar: go/api: Cannot mkdir: No such file or directory
tar: go/api/README: Cannot open: No such file or directory
... (se repite, una vez por cada archivo del paquete)
```

**Ce n'est pas une erreur, ce sont des centaines.** `tar` essaie d'écrire chaque fichier du paquet, un par un, et
chacun échoue de la même façon parce qu'aucun n'a la permission d'écrire dans `/usr/local`. Si tu vois ce mur de
lignes répétées, la cause est toujours la même et se corrige toujours de la même façon : fais précéder la commande
de `sudo`.

### 1.5 Dis à ton système où il se trouve — et le piège des terminaux modernes

Maintenant Go existe dans `/usr/local/go/bin/go`, mais comme tu l'as vu à la section 1.2, si tu tapes `go` il va te
répondre `command not found` : ce dossier n'est pas encore dans le PATH.

L'instruction que tu verras dans presque tous les tutoriels, y compris la documentation officielle de Go, est
celle-ci :

```bash
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ~/.profile
```

- `~/.profile` est un fichier que, en théorie, ton système lit **à chaque ouverture de session**.
- `source ~/.profile` applique la modification **dans ce terminal, tout de suite**, sans attendre la prochaine
  session.

Et voici ce que la plupart des tutoriels ne te disent pas, et que j'ai vérifié par un test réel : **un nouveau
terminal n'est pas toujours une « nouvelle session ».**

`~/.profile` n'est lu que par ce que Linux appelle un **shell de connexion** (*login shell*) — celui qui démarre quand tu ouvres une
session sur le système (par exemple, en allumant la machine et en te connectant avec ton utilisateur et ton mot de
passe). Mais la plupart des applications de terminal (celui de Mint compris, dans sa configuration par défaut),
quand elles ouvrent une nouvelle fenêtre ou un nouvel onglet, **ne** démarrent **pas** un shell de connexion : elles
démarrent un shell interactif normal, et ceux-là lisent un autre fichier, `~/.bashrc`, **pas** `~/.profile`.

Je l'ai vérifié ainsi, en simulant exactement ce scénario —« j'ai déjà modifié le fichier, je ferme le terminal,
j'en ouvre un autre »— dans un conteneur fraîchement installé :

```
# después de agregar el export SOLO a ~/.profile, en una terminal nueva:
$ go version
bash: go: command not found          ← sigue sin funcionar

# después de agregar la MISMA línea también a ~/.bashrc, en una terminal nueva:
$ go version
go version go1.27.1 linux/arm64      ← ahora sí
```

**L'instruction « ferme le terminal, ouvre-en un autre » ne suffit pas toujours**, et quand elle ne suffit pas, on
dirait que tu as mal fait quelque chose alors qu'en réalité tu as suivi le tutoriel à la lettre. C'est pourquoi la
recommandation de cette leçon, plus robuste que celle de la plupart des guides, est d'ajouter la ligne **aux deux
fichiers** :

```bash
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

`~/.profile` couvre le cas d'une vraie session de connexion (par exemple, si tu utilises SSH ou si tu changes
d'utilisateur) ; `~/.bashrc` couvre le cas, bien plus courant au quotidien, d'ouvrir une fenêtre ou un onglet de
terminal neuf au sein d'une session déjà ouverte.

💡 **Si tu utilises `zsh`** au lieu de `bash` (tu le sais si ton terminal a un aspect différent ou si on te l'a
changé), le fichier équivalent à `.bashrc` est `~/.zshrc`. Pour savoir quel shell tu utilises : `echo $SHELL`.

🔑 **Comment distinguer si ton terminal ouvre un shell de connexion, sans deviner :** lance `echo $0` juste après
l'ouverture. Si la réponse commence par un tiret (`-bash`, `-zsh`), c'est un shell de connexion et il lit bien
`~/.profile`. S'il n'y a pas de tiret (`bash`, `zsh`), ce n'en est pas un, et il te faut la modification dans
`~/.bashrc` (ou `~/.zshrc`) pour qu'elle persiste.

### 1.6 Vérifie que ça a marché

```bash
go version
```

Il doit répondre quelque chose comme `go version go1.27.1 linux/amd64` (ou `linux/arm64` si ta machine est ARM).
S'il continue de dire `command not found`, vérifie avec `echo $0` de quel type de shell il s'agit et dans quel
fichier tu as mis la modification, en suivant la section 1.5.

### 1.7 Ton premier programme

```bash
mkdir -p ~/w/curso-go/hola && cd ~/w/curso-go/hola
go mod init hola
```

`go mod init` crée un fichier `go.mod`. C'est la carte d'identité du projet : il dit comment il s'appelle et
quelle version de Go il utilise. C'est le mécanisme de modules que nous avons déjà mentionné à la section 1.1, et
c'est le même avec lequel tu vas démarrer le projet `revisor` dans la leçon 5.

Crée `main.go` avec ceci :

```go
package main

import "fmt"

func main() {
    fmt.Println("hola, ya tengo Go")
}
```

Ligne par ligne, parce que chacune a sa raison :

| | |
|---|---|
| `package main` | *« ce fichier appartient au paquet `main` »*. **Le paquet `main` est spécial : c'est le seul qui produit un programme exécutable.** Sans cette ligne, tu aurais une bibliothèque, pas un programme |
| `import "fmt"` | *« je vais utiliser des choses du paquet `fmt` »* (de *format*), qui fournit ce qu'il faut pour afficher. Go ne charge **rien** par défaut : ce que tu utilises, tu le demandes |
| `func main()` | **la fonction où commence ton programme.** Quand tu l'exécutes, Go cherche exactement cette fonction. Si elle s'appelle autrement, il ne démarre pas |
| `fmt.Println(...)` | affiche et passe à la ligne. Le point signifie *« la fonction `Println` qui se trouve dans `fmt` »* |

Lance-le :

```bash
$ go run main.go
hola, ya tengo Go
```

Cette sortie est réelle : je l'ai exécutée avant d'écrire cette ligne.

### 1.8 Les commandes que tu utiliseras toujours — et la subtilité de `go.mod`

```bash
go run main.go     # compila y ejecuta de una vez, sin dejar archivo. Para probar mientras trabajas
go build           # crea el programa ejecutable y lo deja ahí
go fmt ./...       # ordena tu código
go test ./...      # corre las pruebas (lección 5)
```

**Une subtilité réelle qui vaut la peine d'être testée par toi-même, parce qu'elle a changé entre les versions de
Go :** avec un programme d'un seul fichier comme celui ci-dessus, `go run main.go` fonctionne **même sans
`go.mod`**. Je l'ai vérifié dans un dossier neuf, sans `go mod init` :

```
$ go run main.go
hola, ya tengo Go
```

Mais `go build`, dans ce même dossier sans `go.mod`, exige bien le module :

```
$ go build
go: go.mod file not found in current directory or any parent directory; see 'go help modules'
```

**Pourquoi cette différence ?** `go run` sur un seul fichier peut tout résoudre sans avoir besoin de connaître le
nom du module, car il n'y a rien qu'un autre fichier puisse importer depuis lui. Dès que ton programme a besoin de
plus d'un fichier — par exemple, si tu essaies d'importer un paquet à toi, comme tu le feras dans la leçon 5 avec
le `revisor` — Go a besoin du `go.mod` pour savoir comment s'appelle ton module et ainsi résoudre ces imports. Je
l'ai vérifié aussi :

```
$ go run main.go            # main.go importa "hola/utilidades", sin go.mod
main.go:4:5: package hola/utilidades is not in std (/usr/local/go/src/hola/utilidades)
```

**La leçon pratique :** lance toujours `go mod init`, dès le premier fichier, même si `go run` marche parfois sans
lui — tu t'épargnes cette erreur précise dès que ton programme aura besoin d'importer un paquet à toi, ce qui
arrivera dès la leçon 5.

**Teste la différence entre `run` et `build` :**

```bash
$ go build
$ ls -la
-rwxr-xr-x 1 tu-usuario tu-usuario 2341973 ... hola
$ ./hola
hola, ya tengo Go
```

### 1.9 Pourquoi le binaire est portable : la liaison statique

🔑 **Le fichier `hola` produit par `go build` est un programme complet et autosuffisant.** Et il ne faut pas le
croire sur parole : je l'ai vérifié en copiant ce binaire exact, compilé sur une machine avec Go installé, vers un
conteneur **propre, sans Go** :

```
$ which go
go no esta instalado
$ ./hola
hola, ya tengo Go
```

Ça a fonctionné. **Il n'a besoin ni d'interpréteur, ni de machine virtuelle, ni de bibliothèques externes
installées à part**, parce que Go fait de la *liaison statique* : au lieu de dire « quand tu t'exécutes,
va chercher la bibliothèque `fmt` quelque part dans le système » (comme le font beaucoup de programmes en C ou en
Python), Go copie dans le binaire lui-même tout ce dont ton programme a besoin pour tourner. Le fichier est plus
lourd pour cette raison —un peu plus de 2 Mo pour un programme d'une ligne— mais en échange tu peux le copier sur
n'importe quelle machine Linux compatible et il fonctionnera tel quel, sans rien installer d'autre.

C'est la caractéristique la plus pratique de Go, et c'est pourquoi tant d'outils serveur sont écrits avec lui — y
compris le `revisor` que tu vas compiler et publier dans la leçon 7.

### 1.10 Un éditeur qui t'aide

Ce n'est pas obligatoire, mais cela te fera gagner beaucoup de temps. Si tu utilises **VS Code**, installe
l'extension officielle **Go** (de `golang.go`). Mesuré le 30-sep-2026 sur le Marketplace de Visual Studio : la
version publiée est la **0.57.2**. Ce numéro vieillira lui aussi — l'important est que tu installes celle que
propose le Marketplace le jour où tu le fais, pas que tu notes ce numéro.

L'extension te soulignera les erreurs pendant que tu écris, au lieu que tu les découvres à la compilation, et te
permet de sauter à la définition de n'importe quoi avec `F12`. La première fois, elle te demandera d'installer
quelques outils supplémentaires (`gopls`, le serveur de langage de Go, entre autres) : réponds oui à tous.

---

## L'erreur que tu vas voir

Tous ces messages sont réels, provoqués à dessein pour cette leçon — pas des descriptions approximatives de ce
qu'ils « diraient probablement » :

| Le symptôme | Message littéral | Ce qui se passe et quoi faire |
|---|---|---|
| Terminal `bash` sans Go dans le PATH | `bash: go: command not found` | Le programme n'est dans aucun dossier du PATH. Vois la section 1.5 |
| Terminal `zsh` sans Go dans le PATH | `zsh: command not found: go` | Même problème, ordre des mots différent parce que c'est un autre shell |
| `tar` sans `sudo` vers `/usr/local` | `tar: go: Cannot mkdir: Permission denied` (répété par fichier) | Il manquait `sudo` avant la commande d'extraction (section 1.4) |
| Tu as ajouté le PATH uniquement à `~/.profile` et ouvert un nouveau terminal | `bash: go: command not found` (persiste) | Ton terminal n'ouvre pas de shell de connexion ; ajoute la ligne aussi à `~/.bashrc` (section 1.5) |
| `go build` dans un dossier sans `go.mod` | `go: go.mod file not found in current directory or any parent directory; see 'go help modules'` | Lance `go mod init <nom>` dans ce dossier |
| Import d'un paquet à toi sans `go.mod` | `main.go:4:5: package hola/utilidades is not in std (...)` | Comme le précédent : sans module, Go ne sait pas résoudre tes propres paquets |
| Tu as oublié `import "fmt"` et utilisé `fmt.Println` | `# command-line-arguments`<br>`./main.go:4:2: undefined: fmt` | Go connaît tous les symboles disponibles ; si tu n'as pas importé le paquet, il n'existe pas pour lui |
| `go version` répond **1.22** au lieu de la version actuelle | (pas d'erreur, mais mauvaise version) | C'est celui d'`apt` qui est utilisé. Retire-le avec `sudo apt remove golang-go` et vérifie que `/usr/local/go/bin` est dans ton PATH |

**Et la règle générale, valable pour les huit :** copie le message d'erreur complet et cherche-le tel quel, sans
le résumer d'abord avec tes mots. Presque toujours quelqu'un l'a déjà rencontré, et le message dit presque
toujours exactement ce qui manque — deviner sans le lire en entier est la façon la plus lente de le résoudre.

---

## Ce qu'on fait mal

- **Installer Go avec le gestionnaire de paquets du système (`apt install golang`).** C'est la suggestion la plus
  courante dans les vieux tutoriels, et c'est exactement le problème de la section « Pourquoi c'est important » :
  une version LTS figée te laisse des années de retard sans te prévenir par la moindre erreur, seulement par des
  comportements bizarres plus tard.
- **Considérer comme acquis que « fermer et rouvrir le terminal » recharge toujours la configuration.** C'est le
  piège mesuré à la section 1.5 : si ton terminal n'ouvre pas un shell de connexion, le fermer et le rouvrir ne
  relit pas `~/.profile`. Le symptôme —« j'ai bien fait et ça ne marche pas »— ne veut pas dire que tu as mal
  fait, il veut dire que ce fichier n'était pas celui que lit ton terminal.
- **Extraire la nouvelle version par-dessus une ancienne installation, sans supprimer d'abord.** Le
  `sudo rm -rf /usr/local/go` de l'étape 1.4 n'est pas décoratif : si tu extrais par-dessus une ancienne
  installation, il reste des fichiers des deux mélangés et le résultat est un Go qui échoue de façons
  incompréhensibles. Supprime d'abord, toujours.
- **Compter sur le fait que `go run` marche sans `go.mod` et donc sauter `go mod init`.** C'est vrai pour un
  fichier isolé (section 1.8), mais cela cesse de l'être dès que ton programme a plus d'un fichier — et d'ici là
  tu as déjà écrit du code qu'il faudra réorganiser. Lance `go mod init` dès le début, toujours.
- **Débattre à la main du style du code.** En Go, le style ne se débat pas : c'est l'outil qui le décide. Essaie
  d'écrire ceci exprès, tout de travers :

  ```go
  package main
  import "fmt"
  func main(){
  fmt.Println( "hola" )
  }
  ```

  Lance `go fmt ./...` et rouvre le fichier : **il s'est ordonné tout seul.** Les disputes sur l'emplacement de
  l'accolade ou le nombre d'espaces de l'indentation n'existent pas, parce que `gofmt` n'a qu'une seule réponse et
  que tout le monde l'utilise. Habitue-toi à le lancer avant d'enregistrer, au lieu de formater à la main.

---

## Exercices

1. Installe Go en suivant les sections 1.3 à 1.6 et confirme la version avec `go version`.
2. Lance `echo $0` dans ton terminal **avant** de toucher au PATH. Selon ce qu'il répond, décide si tu dois
   modifier `~/.profile`, `~/.bashrc`, ou les deux (section 1.5) — et explique pourquoi, avec tes mots, avant de
   continuer.
3. Crée le projet `hola` de la section 1.7, lance `go run main.go` puis `go build` + `./hola`.
4. (Comme à la section 1.9) Copie ton binaire `hola` vers une autre machine Linux, ou vers une machine virtuelle /
   un conteneur sans Go installé, et confirme qu'il tourne pareil.
5. Change le message de `fmt.Println` pour quelque chose de ton cru, déforme l'indentation à dessein, et lance
   `go fmt ./...`. Vérifie que le fichier est bien formaté.
6. (Un peu plus difficile) Crée un nouveau dossier, **sans** `go mod init`, avec un `main.go` d'un seul fichier.
   Confirme que `go run main.go` fonctionne quand même. Puis lance `go build` dans ce même dossier et compare
   l'erreur avec celle de la section 1.8. Explique, avec tes mots, pourquoi l'un fonctionne et l'autre non.
7. (Un peu plus difficile) Supprime exprès la ligne `import "fmt"` et lance `go run main.go`. Lis l'erreur en
   entier, sans la chercher encore, et essaie d'expliquer avec tes mots ce qu'elle te dit avant de continuer.

### Corrigés

1. `go version` doit afficher la même version que celle que tu as vue avec `curl -s 'https://go.dev/VERSION?m=text'`,
   jamais `go1.22.x` (c'est celle de l'`apt` de Mint).
2. Si `echo $0` répond avec un tiret au début (`-bash`), ton terminal ouvre des shells de connexion et `~/.profile`
   suffit. S'il répond sans tiret (`bash`), ce n'est pas un shell de connexion, et seul `~/.bashrc` (ou `~/.zshrc`
   sous zsh) persistera entre les nouvelles fenêtres — ajoute la modification là aussi, comme à la section 1.5.
3. `go run main.go` affiche `hola, ya tengo Go`. Après `go build`, un fichier `hola` apparaît dans le dossier ;
   `./hola` l'exécute directement, sans recompiler.
4. Le binaire tourne pareil sur la machine sans Go, parce que Go fait de la liaison statique (section
   1.9) : tout ce dont le programme a besoin est déjà copié dans le fichier lui-même.
5. Avant `go fmt`, le fichier se voit exactement tel qu'il a été écrit (de travers). Après, `gofmt` le réarrange
   avec l'indentation et les espaces standard de Go — sans que tu décides quoi que ce soit.
6. `go run main.go` fonctionne parce que, avec un seul fichier, Go n'a besoin de résoudre aucun import propre pour
   le compiler et l'exécuter d'un coup. `go build`, en revanche, exige de connaître le nom du module dès le
   départ et répond `go: go.mod file not found...`. La différence est que `build` laisse le binaire prêt à être
   utilisé en dehors de cette invocation ponctuelle, et pour cela il a besoin d'une identité de projet — `run`
   non.
7. Le compilateur répond :

   ```
   # command-line-arguments
   ./main.go:4:5: undefined: fmt
   ```

   Il dit que tu as utilisé `fmt.Println` sans avoir importé le paquet `fmt` : le compilateur connaît tous les
   symboles que tu peux utiliser, et si tu ne l'as pas importé, il n'existe pas pour lui. On corrige en
   remettant la ligne `import "fmt"`.

---

## Comment savoir que j'y suis arrivé

- [ ] `go version` répond avec la version que j'ai vue sur `go.dev/VERSION`, pas avec la 1.22.
- [ ] Je sais si mon terminal ouvre des shells de connexion ou non (`echo $0`), et dans quel(s) fichier(s) j'ai mis la
      modification du PATH en conséquence.
- [ ] `go run main.go` affiche mon message.
- [ ] `go build` a produit un fichier et `./hola` fonctionne.
- [ ] J'ai copié mon binaire vers une autre machine (ou conteneur) sans Go et il a tourné pareil.
- [ ] J'ai essayé `go fmt` sur du code mal écrit et il l'a ordonné.
- [ ] J'ai provoqué l'erreur de `go build` sans `go.mod` et je peux expliquer pourquoi `go run` ne l'a pas eue.
- [ ] J'ai provoqué l'erreur d'`import` manquant et j'ai compris le message sans avoir besoin du corrigé.
- [ ] Je peux expliquer, sans regarder le texte, quelle différence il y a entre `GOROOT`, `GOPATH` et le `PATH`.

---

## Résumé

- N'installe jamais Go avec le gestionnaire de paquets du système (`apt install golang`) : les versions LTS de
  Linux se figent, et le paquet reste des années en arrière sans prévenir par la moindre erreur.
- **`GOROOT`** est l'endroit où vit Go ; **`GOPATH`** est l'endroit où Go range les paquets téléchargés ; tu ne
  touches à aucun des deux à la main dans un projet moderne avec modules (`go.mod`).
- Le **`PATH`** est la liste de dossiers où le terminal cherche les programmes. Avoir Go installé et avoir Go dans
  le `PATH` sont deux choses différentes.
- Un nouveau terminal n'est **pas toujours** une session de connexion : la plupart ne lisent que `~/.bashrc` (ou
  `~/.zshrc`), pas `~/.profile`. Ajoute la modification du `PATH` aux deux fichiers.
- `go mod init` crée l'identité du projet (`go.mod`) ; tout projet Go en a besoin, même si `go run` marche parfois
  sans elle avec un seul fichier.
- `go run` compile et exécute sans laisser de fichier ; `go build` laisse le binaire ; `go fmt` formate
  automatiquement, sans débat possible sur le style.
- Le binaire que produit `go build` est autosuffisant grâce à la **liaison statique** : il tourne sur une
  autre machine Linux compatible sans que Go y soit installé.
- Chaque message d'erreur de cette leçon est littéral, provoqué à dessein : `command not found`,
  `Permission denied`, `go.mod file not found`, `undefined: fmt` — copie-les tels quels quand tu les cherches.

---

## Pour aller plus loin

1. [Download and install](https://go.dev/doc/install) — le guide officiel d'installation, la source de vérité si
   quelque chose dans cette leçon vieillit. Il inclut la même recommandation de `~/.profile` que nous utilisons
   ici, avec le même avertissement que la modification « peut ne pas s'appliquer avant la prochaine ouverture de
   session » — ce qui est justement ce qui a été mesuré et expliqué en profondeur à la section 1.5.
2. [A Tour of Go](https://go.dev/tour/) — le parcours interactif officiel ; il commence exactement là où se
   termine cette leçon.
3. [Effective Go](https://go.dev/doc/effective_go) — inclut la section sur `gofmt` et sur pourquoi le format ne se
   débat pas en Go.
4. [GNU Bash — Manual: Bash Startup Files](https://www.gnu.org/software/bash/manual/html_node/Bash-Startup-Files.html) —
   la référence officielle sur quand `~/.profile` est lu et quand `~/.bashrc` l'est ; elle vaut la peine d'être
   lue en entier une seule fois dans sa vie, parce que la confusion entre les deux n'est pas propre à Go.

### Termes de cette leçon

| Terme | Ce que ça signifie |
|---|---|
| `GOROOT` | le dossier où vit l'installation de Go : le compilateur, `gofmt`, la bibliothèque standard |
| `GOPATH` | le dossier où Go range les paquets téléchargés ; avant les modules, tout ton code devait aussi y vivre |
| `PATH` | la liste de dossiers où le terminal cherche les programmes exécutables |
| shell de connexion | la session de terminal qui démarre à l'ouverture de session sur le système ; c'est la seule qui lit `~/.profile` |
| `go.mod` | le fichier qui identifie un projet Go : son nom et la version de Go qu'il utilise |
| binaire | le fichier exécutable que produit `go build`, autosuffisant, sans dépendre de la présence de Go |
| liaison statique | la technique par laquelle Go copie dans le binaire tout ce dont le programme a besoin, au lieu de le chercher dans le système à l'exécution |
| `gofmt` | l'outil qui formate automatiquement le code Go, sans options à débattre |

---

**Précédent :** [Leçon 0 — Ce qu'est Go](00-introduccion.md) ·
**Suivant :** [Leçon 2 — Variables, fonctions et types](02-fundamentos.md)
