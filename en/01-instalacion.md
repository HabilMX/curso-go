# Lesson 1 — Installing Go on your Linux Mint and your first program

**Duration:** 45-60 minutes (or 2 sessions of 30). It is the only lesson that doesn't teach Go as a
language: it teaches the tool and the terrain where you will work — and that terrain has more traps than
it seems at first glance.

**By the end you will be able to:**

- Install the official, current version of Go on Linux Mint, without depending on the system package, and
  explain why that package is outdated by design, not by oversight.
- Explain what `GOROOT`, `GOPATH`, and the `PATH` are, and why Go doesn't work until you touch the second.
- Recognize the difference between a terminal that reads your configuration and one that doesn't —the
  real cause of "I added it and it still doesn't work"— and know which one is yours.
- Create a project with `go mod init`, write a program, and run it with `go run`.
- Produce a binary with `go build` and explain, with a real test, why it runs on another machine without
  Go installed.
- Use `go fmt` so you never argue about code style.
- Recognize, word for word, the most common error messages of this stage —from the compiler, about
  permissions, and about paths— and know what to do with each.

---

## Why it matters

The first thing an old tutorial is going to suggest is `apt install golang`, and it is a mistake. It is
not a minor "slightly old version" mistake: it is the kind of mistake that makes you lose an entire
afternoon without understanding why, weeks after installing, because the symptom doesn't show up when you
install — it shows up when you are already programming and something that "should work" doesn't.

This is not opinion, it is measured. On 2026-09-30 I measured, in a container with the same base as Linux
Mint 22.3 (Ubuntu 24.04 "Noble Numbat", which is where Mint takes its packages from):

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

**The system package offers 1.22. The current official version, measured at the same instant, is
1.27.1.** And this is not an oversight that someone is going to fix: it is Ubuntu's and Mint's policy.

**Why is it frozen like this, on purpose?** Ubuntu 24.04 is an LTS (*Long Term Support*) release: the day
it comes out, its main packages are **frozen** and only receive security patches, never new versions of
the program. It is the right decision for a server that must not change behavior just because you install
updates — but it means that a compiler for a language that releases new versions twice a year stays fixed
at the one that existed when Ubuntu 24.04 was published (April 2024), forever, for as long as that
Ubuntu version exists. Linux Mint has no Go packages of its own: it literally inherits Ubuntu's, so it
inherits the freeze too.

**Why does it matter to you, specifically?** Because you are going to look things up online, the examples
are going to use functions or behaviors that your Go 1.22 doesn't have, and the errors won't say *"you
need a newer version"*: they will say things that make no sense for what you are seeing on screen. It is
the kind of problem that looks like your mistake and is actually a tooling lag — and it happens to you
**before** you write your first line of code, which is why this lesson exists before any other.

🔑 **The right way, and the one the Go project itself recommends:** download the official package
directly from `go.dev`, not from your distribution's package manager. That is what this lesson installs,
step by step, and what we will check at every step with the real output of the commands.

---

## The concepts

### 1.1 What GOROOT and GOPATH are, and why they used to matter more than now

Before installing, it is worth knowing what you will have afterwards, because the names `GOROOT` and
`GOPATH` show up in almost any installation error you search for online, and many answers are written for
a version of Go from ten years ago.

- **`GOROOT`** is the folder where **Go itself** lives: the compiler, `gofmt`, the standard library. It is
  the folder you will create in step 1.4 (`/usr/local/go`). You never touch it by hand.
- **`GOPATH`** is the folder where Go keeps **your stuff**: packages downloaded from the internet, and
  —in old versions of Go, before 2019— also where you had to put **all** your code, without exception,
  inside a fixed structure (`$GOPATH/src/github.com/your-user/your-project`). If you ever see a tutorial
  that asks you to create that folder structure, it is from that era.

Today you hardly touch it because since Go 1.11 (2018) there are **modules** (`go.mod`, which you will
create in step 1.6): your project can live in any folder, with any name you want, and Go no longer needs
you to follow an imposed folder structure. `GOPATH` still exists, but now only as a cache of downloaded
packages, not as the only place where your code can live.

Check it with `go env`, which shows you the current configuration of your installation (you can run this
**after** installing, in step 1.5):

```bash
go env GOROOT GOPATH GOBIN
```

On a clean, fresh installation, you will see something like:

```
/usr/local/go
/home/tu-usuario/go
/home/tu-usuario/go/bin
```

You didn't create any of those three folders by hand: the first is created by the installation step
(1.4), the other two are decided by Go itself, with reasonable default values.

### 1.2 The PATH: what it is, and why "already installed" doesn't mean "already works"

When you type a command in the terminal, for example `go`, the terminal doesn't magically know where that
program is: it checks, one by one, a list of folders stored in a variable called **PATH**, and uses the
first program it finds with that name. If no folder in the list has a program called `go`, it answers
with an error — and that error is literal, not approximate:

```bash
$ go version
bash: go: command not found
```

I verified this in a container with Go **already extracted** into `/usr/local/go`, before touching the
PATH: the program exists on disk, but the terminal can't find it because it doesn't know where to look.
**"Installed" and "in the PATH" are two different things**, and the confusion between the two is the
source of almost every stumble in this lesson.

This also explains why the message changes slightly depending on the terminal (`zsh` instead of `bash`,
which you will see if you use macOS to follow the course while you practice, or if you changed the
default terminal of your Mint):

```bash
$ go version
zsh: command not found: go
```

Same problem, same mechanism, **different word order**: identify which one is yours with `echo $SHELL`
before searching for the error online, because searching for the wrong message will lead you to answers
for the wrong shell.

### 1.3 Find out which is the latest version

Don't copy it from here: this document ages the same way Mint's repositories do, and you already saw in
the previous section how much it can cost a tool to stay fixed in time.

```bash
curl -s 'https://go.dev/VERSION?m=text' | head -1
```

It will answer something like `go1.27.1`. **That is the one you will install**, whatever it is when you do
it.

### 1.4 Download it and install it

Replace `go1.27.1` with what the previous command told you. `linux-amd64` is right for a normal PC or
laptop (if your machine were ARM, it would be `linux-arm64`; to find out: `dpkg --print-architecture`).

```bash
cd /tmp
wget https://go.dev/dl/go1.27.1.linux-amd64.tar.gz
```

And now the installation:

```bash
sudo rm -rf /usr/local/go                          # borra una instalación anterior, si había
sudo tar -C /usr/local -xzf go1.27.1.linux-amd64.tar.gz
```

**What you just did**, because it is worth understanding and not just copying:

- `tar` decompresses the archive.
- `-C /usr/local` tells it *"do it inside this folder"*.
- `-xzf` is "extract" (`x`), "it is compressed with gzip" (`z`), "from this file" (`f`).
- The result is a `/usr/local/go` folder with all of Go inside — the `GOROOT` from section 1.1.

**Why `sudo`, if you never needed it to install something with a package manager?** Because `/usr/local`
is a system folder, not one of your user's, and on Linux writing there requires administrator
permissions. Without `sudo`, this is exactly what you will see —I triggered it on purpose, without
`sudo`, so you see the real message and not an invented one—:

```
$ tar -C /usr/local -xzf go1.27.1.linux-amd64.tar.gz
tar: go: Cannot mkdir: Permission denied
tar: go/VERSION: Cannot open: No such file or directory
tar: go/api: Cannot mkdir: No such file or directory
tar: go/api/README: Cannot open: No such file or directory
... (se repite, una vez por cada archivo del paquete)
```

**It is not one error, it is hundreds.** `tar` tries to write each file in the package, one by one, and
each one fails the same way because none has permission to write to `/usr/local`. If you see this wall of
repeated lines, the cause is always the same and it is always fixed the same way: put `sudo` in front.

### 1.5 Tell your system where it is — and the trap of modern terminals

Now Go exists at `/usr/local/go/bin/go`, but as you saw in section 1.2, if you type `go` it will tell you
`command not found`: that folder is not in the PATH yet.

The instruction you will see in almost any tutorial, including Go's official documentation, is this one:

```bash
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ~/.profile
```

- `~/.profile` is a file that, in theory, your system reads **every time you log in**.
- `source ~/.profile` applies the change **in this terminal, right now**, without waiting for the next
  session.

And here comes what most tutorials don't tell you, and which I verified with a real test: **a new
terminal is not always a "new session".**

`~/.profile` is only read by what Linux calls a **login shell** — the one that starts when you log in to
the system (for example, when you turn on the machine and sign in with your username and password). But
most terminal applications (Mint's terminal included, in its default configuration) when opening a new
window or tab do **not** start a login shell: they start a normal interactive shell, and those read a
different file, `~/.bashrc`, **not** `~/.profile`.

I verified it like this, simulating exactly that scenario —"I already edited the file, I close the
terminal, I open another one"— in a freshly installed container:

```
# después de agregar el export SOLO a ~/.profile, en una terminal nueva:
$ go version
bash: go: command not found          ← sigue sin funcionar

# después de agregar la MISMA línea también a ~/.bashrc, en una terminal nueva:
$ go version
go version go1.27.1 linux/arm64      ← ahora sí
```

**The instruction "close the terminal, open another" is not always enough**, and when it isn't, it looks
like you did something wrong when in fact you followed the tutorial to the letter. That is why this
lesson's recommendation, more robust than most guides', is to add the line **to both files**:

```bash
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

`~/.profile` covers the case of a real login session (for example, if you use SSH or switch users);
`~/.bashrc` covers the much more common day-to-day case of opening a new terminal window or tab on a
session that was already open.

💡 **If you use `zsh`** instead of `bash` (you know it if your terminal looks different or if someone
changed it for you), the equivalent of `.bashrc` is `~/.zshrc`. To find out which shell you use:
`echo $SHELL`.

🔑 **How to tell whether your terminal opens a login shell, without guessing:** run `echo $0` right after
opening it. If the answer starts with a dash (`-bash`, `-zsh`), it is a login shell and it does read
`~/.profile`. If it has no dash (`bash`, `zsh`), it is not, and you need the change in `~/.bashrc` (or
`~/.zshrc`) for it to persist.

### 1.6 Check that it worked

```bash
go version
```

It should answer something like `go version go1.27.1 linux/amd64` (or `linux/arm64` if your machine is
ARM). If it still says `command not found`, check with `echo $0` what kind of shell it is and in which
file you put the change, following section 1.5.

### 1.7 Your first program

```bash
mkdir -p ~/w/curso-go/hola && cd ~/w/curso-go/hola
go mod init hola
```

`go mod init` creates a `go.mod` file. It is the project's identity card: it says what it is called and
which version of Go it uses. It is the module mechanism we already mentioned in section 1.1, and it is the
same one you will use to start the `revisor` project in lesson 5.

Create `main.go` with this:

```go
package main

import "fmt"

func main() {
    fmt.Println("hola, ya tengo Go")
}
```

Line by line, because each one has its reason:

| | |
|---|---|
| `package main` | *"this file belongs to the `main` package"*. **The `main` package is special: it is the only one that produces an executable program.** Without this line you would have a library, not a program |
| `import "fmt"` | *"I am going to use things from the `fmt` package"* (from *format*), which brings what is needed to print. Go loads **nothing** by default: what you use, you ask for |
| `func main()` | **the function where your program starts.** When you run it, Go looks for exactly this function. If it is named anything else, it doesn't start |
| `fmt.Println(...)` | prints and moves to a new line. The dot means *"the `Println` function that is inside `fmt`"* |

Run it:

```bash
$ go run main.go
hola, ya tengo Go
```

That output is real: I ran it before writing this line.

### 1.8 The commands you will always use — and the subtlety of `go.mod`

```bash
go run main.go     # compila y ejecuta de una vez, sin dejar archivo. Para probar mientras trabajas
go build           # crea el programa ejecutable y lo deja ahí
go fmt ./...       # ordena tu código
go test ./...      # corre las pruebas (lección 5)
```

**A real subtlety worth testing yourself, because it changed between Go versions:** with a single-file
program like the one above, `go run main.go` works **even without `go.mod`**. I checked it in a new
folder, without `go mod init`:

```
$ go run main.go
hola, ya tengo Go
```

But `go build`, in that same folder without `go.mod`, does demand the module:

```
$ go build
go: go.mod file not found in current directory or any parent directory; see 'go help modules'
```

**Why the difference?** `go run` on a single file can resolve everything without needing to know the
module name, because there is nothing another file could import from it. As soon as your program needs
more than one file — for example, if you try to import a package of your own, as you will do in lesson 5
with the `revisor` — Go needs `go.mod` to know what your module is called and thus resolve those imports.
I checked this too:

```
$ go run main.go            # main.go importa "hola/utilidades", sin go.mod
main.go:4:5: package hola/utilidades is not in std (/usr/local/go/src/hola/utilidades)
```

**The practical lesson:** always run `go mod init`, from the first file, even if `go run` sometimes works
without it — you save yourself this exact error as soon as your program needs to import a package of its
own, which will already happen in lesson 5.

**Try the difference between `run` and `build`:**

```bash
$ go build
$ ls -la
-rwxr-xr-x 1 tu-usuario tu-usuario 2341973 ... hola
$ ./hola
hola, ya tengo Go
```

### 1.9 Why the binary is portable: static linking

🔑 **The `hola` file that `go build` produced is a complete, self-sufficient program.** And you don't have
to take this on faith: I verified it by copying that exact binary, compiled on a machine with Go
installed, to a **clean container, without Go**:

```
$ which go
go no esta instalado
$ ./hola
hola, ya tengo Go
```

It worked. **It needs no interpreter, no virtual machine, no external libraries installed separately**,
because Go does *static linking*: instead of saying "when you run, look for the `fmt` library somewhere on
the system" (which is how many C or Python programs work), Go copies inside the binary itself everything
your program needs to run. The file weighs more because of that —a little over 2 MB for a one-line
program— but in exchange you can copy it to any compatible Linux machine and it will work as is, without
installing anything else.

That is Go's most practical feature, and it is why so many server tools are written in it — including the
`revisor` that you will build and publish in lesson 7.

### 1.10 An editor to help you

It is not mandatory, but it will save you a lot of time. If you use **VS Code**, install the official
**Go** extension (from `golang.go`). Measured on 2026-09-30 on the Visual Studio Marketplace: the
published version is **0.57.2**. That number will age too — what matters is that you install the one the
Marketplace offers the day you do it, not that you write down this number.

The extension will underline errors as you type, instead of you finding out when you compile, and lets you
jump to the definition of anything with `F12`. The first time it will ask you to install some extra tools
(`gopls`, Go's language server, among others): say yes to all of them.

---

## The error you will see

All of these are real messages, triggered on purpose for this lesson — not approximate descriptions of
what they "would probably" say:

| The symptom | Literal message | What is happening and what to do |
|---|---|---|
| `bash` terminal without Go in the PATH | `bash: go: command not found` | The program is not in any folder of the PATH. Check section 1.5 |
| `zsh` terminal without Go in the PATH | `zsh: command not found: go` | Same problem, different word order because it is another shell |
| `tar` without `sudo` against `/usr/local` | `tar: go: Cannot mkdir: Permission denied` (repeated per file) | `sudo` was missing before the extraction command (section 1.4) |
| You added the PATH only to `~/.profile` and opened a new terminal | `bash: go: command not found` (persists) | Your terminal doesn't open a login shell; add the line to `~/.bashrc` too (section 1.5) |
| `go build` in a folder without `go.mod` | `go: go.mod file not found in current directory or any parent directory; see 'go help modules'` | Run `go mod init <name>` in that folder |
| Import of a package of your own without `go.mod` | `main.go:4:5: package hola/utilidades is not in std (...)` | Same as the previous one: without a module, Go doesn't know how to resolve your own packages |
| You forgot `import "fmt"` and used `fmt.Println` | `# command-line-arguments`<br>`./main.go:4:2: undefined: fmt` | Go knows all the available symbols; if you didn't import the package, it doesn't exist for it |
| `go version` answers **1.22** instead of the current one | (no error, but wrong version) | The `apt` one is being used. Remove it with `sudo apt remove golang-go` and confirm that `/usr/local/go/bin` is in your PATH |

**And the general rule, which applies to all eight:** copy the complete error message and search for it
as is, without summarizing it in your own words first. Almost always someone has already had it, and the
message almost always says exactly what is missing — guessing without reading it in full is the slowest
way to solve it.

---

## What goes wrong

- **Installing Go with the system package manager (`apt install golang`).** It is the most common
  suggestion in old tutorials, and it is exactly the problem from the "Why it matters" section: a frozen
  LTS version leaves you years behind without warning you with any error, only with strange behaviors
  later on.
- **Taking for granted that "closing and opening the terminal" always reloads the configuration.** It is
  the trap measured in section 1.5: if your terminal doesn't open a *login* shell, closing it and opening
  it again doesn't re-read `~/.profile`. The symptom —"I did it right and it doesn't work"— doesn't mean
  you did something wrong, it means that file wasn't the one your terminal reads.
- **Extracting the new version on top of an old installation, without deleting first.** The
  `sudo rm -rf /usr/local/go` in step 1.4 is not decorative: if you extract on top of an old installation,
  files from both end up mixed and the result is a Go that fails in incomprehensible ways. Delete first,
  always.
- **Trusting that `go run` works without `go.mod` and therefore skipping `go mod init`.** It is true for a
  single loose file (section 1.8), but it stops being true as soon as your program has more than one
  file — and by then you have already written code that you will have to reorganize. Run `go mod init`
  from the start, always.
- **Arguing about code style by hand.** In Go style is not argued: the tool decides it. Try writing this
  on purpose, all crooked:

  ```go
  package main
  import "fmt"
  func main(){
  fmt.Println( "hola" )
  }
  ```

  Run `go fmt ./...` and reopen the file: **it ordered itself.** There are no fights about where the
  brace goes or how many spaces the indentation takes, because `gofmt` has a single answer and everyone
  uses it. Get used to running it before saving, instead of formatting by hand.

---

## Exercises

1. Install Go following sections 1.3 to 1.6 and confirm the version with `go version`.
2. Run `echo $0` in your terminal **before** touching the PATH. Depending on what it answers, decide
   whether you need to touch `~/.profile`, `~/.bashrc`, or both (section 1.5) — and explain why, in your
   own words, before continuing.
3. Create the `hola` project from section 1.7, run `go run main.go` and then `go build` + `./hola`.
4. (As in section 1.9) Copy your `hola` binary to another Linux machine, or to a virtual machine /
   container without Go installed, and confirm that it runs the same.
5. Change the `fmt.Println` message to something of your own, mess up the indentation on purpose, and run
   `go fmt ./...`. Verify that the file ended up properly formatted.
6. (A bit harder) Create a new folder, **without** `go mod init`, with a single-file `main.go`. Confirm
   that `go run main.go` works anyway. Then run `go build` in that same folder and compare the error with
   the one in section 1.8. Explain, in your own words, why one works and the other doesn't.
7. (A bit harder) Delete the `import "fmt"` line on purpose and run `go run main.go`. Read the complete
   error, without searching for it yet, and try to explain in your own words what it is telling you
   before continuing.

### Solutions

1. `go version` must print the same version you saw in `curl -s 'https://go.dev/VERSION?m=text'`, never
   `go1.22.x` (that is Mint's `apt` one).
2. If `echo $0` answers with a dash at the start (`-bash`), your terminal opens login shells and
   `~/.profile` is enough. If it answers without a dash (`bash`), it is not a login shell, and only
   `~/.bashrc` (or `~/.zshrc` in zsh) will persist between new windows — add the change there too, as in
   section 1.5.
3. `go run main.go` prints `hola, ya tengo Go`. After `go build`, a `hola` file appears in the folder;
   `./hola` runs it directly, without compiling again.
4. The binary runs the same on the machine without Go, because Go does static linking (section 1.9):
   everything the program needs is already copied inside the file itself.
5. Before `go fmt`, the file looks exactly as it was written (crooked). After, `gofmt` rearranges it with
   Go's standard indentation and spacing — without you deciding any of that.
6. `go run main.go` works because, with a single file, Go doesn't need to resolve any import of your own
   to compile and run it in one go. `go build`, on the other hand, demands to know the module name from
   the very first moment and answers `go: go.mod file not found...`. The difference is that `build` leaves
   the binary ready to be used outside that one invocation, and for that it needs a project identity —
   `run` doesn't.
7. The compiler answers:

   ```
   # command-line-arguments
   ./main.go:4:5: undefined: fmt
   ```

   It says you used `fmt.Println` without having imported the `fmt` package: the compiler knows all the
   symbols you can use, and if you didn't import it, it doesn't exist for it. It is fixed by putting the
   `import "fmt"` line back.

---

## How I know I got it

- [ ] `go version` answers with the version I saw at `go.dev/VERSION`, not with 1.22.
- [ ] I know whether my terminal opens login shells or not (`echo $0`), and in which file(s) I put the PATH
      change accordingly.
- [ ] `go run main.go` prints my message.
- [ ] `go build` produced a file and `./hola` works.
- [ ] I copied my binary to another machine (or container) without Go and it ran the same.
- [ ] I tried `go fmt` on badly written code and it ordered it.
- [ ] I triggered the `go build` error without `go.mod` and I can explain why `go run` didn't have it.
- [ ] I triggered the missing `import` error and understood the message without needing the solution.
- [ ] I can explain, without looking at the text, the difference between `GOROOT`, `GOPATH`, and the
      `PATH`.

---

## Summary

- Never install Go with the system package manager (`apt install golang`): Linux LTS versions are frozen,
  and the package ends up years behind without warning you with any error.
- **`GOROOT`** is where Go lives; **`GOPATH`** is where Go keeps downloaded packages; you touch neither
  by hand in a modern project with modules (`go.mod`).
- The **`PATH`** is the list of folders where the terminal looks for programs. Having Go installed and
  having Go in the `PATH` are two different things.
- A new terminal is **not always** a *login* session: most only read `~/.bashrc` (or `~/.zshrc`), not
  `~/.profile`. Add the `PATH` change to both files.
- `go mod init` creates the project's identity (`go.mod`); every Go project needs it, even though `go run`
  sometimes works without it with a single file.
- `go run` compiles and runs without leaving a file; `go build` leaves the binary; `go fmt` formats
  automatically, with no possible argument about style.
- The binary that `go build` produces is self-sufficient thanks to **static linking**: it runs on another
  compatible Linux machine without having Go installed.
- Every error message in this lesson is literal, triggered on purpose: `command not found`,
  `Permission denied`, `go.mod file not found`, `undefined: fmt` — copy them as is when you search for
  them.

---

## Further reading

1. [Download and install](https://go.dev/doc/install) — the official installation guide, the source of
   truth if anything in this lesson ages. It includes the same `~/.profile` recommendation we use here,
   with the same warning that the change "may not take effect until the next time you log in" — which is
   exactly what was measured and explained in depth in section 1.5.
2. [A Tour of Go](https://go.dev/tour/) — the official interactive tour; it starts right where this
   lesson ends.
3. [Effective Go](https://go.dev/doc/effective_go) — includes the section on `gofmt` and why formatting
   is not argued about in Go.
4. [GNU Bash — Manual: Bash Startup Files](https://www.gnu.org/software/bash/manual/html_node/Bash-Startup-Files.html) —
   the official reference for when `~/.profile` is read and when `~/.bashrc` is; it is worth reading in
   full once in your life, because the confusion between the two is not exclusive to Go.

### Terms from this lesson

| Term | What it means |
|---|---|
| `GOROOT` | the folder where the Go installation lives: the compiler, `gofmt`, the standard library |
| `GOPATH` | the folder where Go keeps downloaded packages; before modules, all your code also had to live there |
| `PATH` | the list of folders where the terminal looks for executable programs |
| *login* shell | the terminal session that starts when you log in to the system; it is the only one that reads `~/.profile` |
| `go.mod` | the file that identifies a Go project: its name and the Go version it uses |
| binary | the executable file that `go build` produces, self-sufficient, not depending on Go being installed |
| static linking | the technique by which Go copies inside the binary everything the program needs, instead of looking for it on the system at run time |
| `gofmt` | the tool that automatically formats Go code, with no options to argue about |

---

**Previous:** [Lesson 0 — What Go is](00-introduccion.md) ·
**Next:** [Lesson 2 — Variables, functions, and types](02-fundamentos.md)
