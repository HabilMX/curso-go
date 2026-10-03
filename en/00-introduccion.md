# Lesson 0 — What Go is, where it came from, and why you should learn it

**Duration:** 30-45 minutes of reading, no code to write yet.

**By the end you will be able to:**

- Tell who created Go, when, and what concrete problem they were trying to solve.
- Explain why a new language exists when there were already dozens.
- Name three programs you use or administer that are written in Go.
- Say what Go is good at and what it is **not** good at, so you don't use it where it doesn't fit.
- Place Go among the other languages you already know or have heard of.

---

## Why it matters

It is 2007. Google has one of the largest codebases in the world, written mostly in C++. And it has a
very concrete and very boring problem: **compiling takes forever.**

A single change in a header file could trigger a recompilation of **45 minutes or more**.

Think about what that means for the person writing code. You change one line. You wait 45 minutes. You
discover you made a mistake with a comma. You change another line. You wait another 45 minutes. **In a
working day you get to try eight things.**

On September 21, 2007, three Google engineers stood in front of a whiteboard to design a language that
wouldn't make them wait. That language is what this whole course is about, and the rest of this lesson
explains how they got from that whiteboard to the language you will install in lesson 1.

---

## The concepts

### 0.1 Who made it (and why it matters)

They were not just any three programmers:

| | Who they are |
|---|---|
| **Ken Thompson** | **Created Unix.** And created the B language, the direct ancestor of C, which he wrote together with Dennis Ritchie. Turing Award —the equivalent of the Nobel Prize in computing— in 1983 |
| **Rob Pike** | Worked on Unix and created Plan 9, its successor. **Co-invented UTF-8**, the way text is stored today in practically the whole world, this document included |
| **Robert Griesemer** | Worked on the V8 JavaScript engine (the one in Chrome and Node.js) and on the Java virtual machine |

🔑 **Read that again: one of the creators of Go is the creator of Unix and a co-author of C.** Fifty years
after inventing the tools on which all modern software is built, the same person designed Go.

That explains a lot about the character of the language: Go feels like C —direct, small, without
frills— but without the traps that made C a minefield.

### 0.2 The three requirements nobody met

The three of them wanted a language with **three** properties at the same time:

1. **That compiles fast.** The original problem.
2. **That runs fast.** Google runs software on thousands of machines: speed costs real money.
3. **That is easy to program in.** So that a new person on the team can read somebody else's code and
   understand it on the first day.

They reviewed the languages that already existed. And here is the interesting finding: **they found
languages that met two of the three, but none that met all three.**

| | Compiles fast | Runs fast | Easy to program |
|---|---|---|---|
| **C / C++** | ❌ | ✅ | ❌ |
| **Java** | 🟡 | ✅ | 🟡 |
| **Python** | ✅ (doesn't compile) | ❌ | ✅ |
| **Go** | ✅ | ✅ | ✅ |

> [!NOTE]
> 🔧 **Software engineering observation 0.1**
> A programming language is not a matter of taste: it is a tool with trade-offs. Each one sacrificed
> something to gain something else. Understanding **what the one you are using sacrificed** is the
> difference between programming with judgment and programming by rote.

### 0.3 How simplicity was won: by removing things

Here is the most counterintuitive decision in Go, and the one that sets it apart the most.

Almost all languages grow: each version adds features. C++ accumulated decades of them, and the result
is so large that **nobody knows all of it**. There are people who have programmed in C++ for twenty years
and still find corners they didn't know about.

**Go did the opposite: it decided what to leave out.** It has no:

- **Classes or inheritance.** The way to reuse code is different, and you will see it in lesson 3.
- **Exceptions** (`try / catch`). Errors are handled another way, and you will see it in lesson 2.
- **Operator overloading.** `+` means add, always, and it can't be changed.
- **Constructors or destructors.**
- **A pile of ways to do the same thing.** Almost always there is **one** idiomatic way.

⚠️ **This is going to annoy you at some point.** You will want to do something that in another language
takes one line and in Go will take you five. The payoff is enormous and shows up months later: **you can
read Go code written by anyone and understand it.** There are no dialects, no tricks, no need to learn how
each team programs.

> 🔑 **The quote that sums up the philosophy**, from Rob Pike:
> *"Clear is better than clever."*

### 0.4 The other reason: computers stopped getting faster

There is a second motive behind Go, and it is about hardware.

Until about 2005, every year faster processors came out and your program ran faster **without you
touching anything**. That ended: processors stopped getting faster and started multiplying **cores**. Your
laptop doesn't have one blazing-fast processor: it has 4, 8, or 16 modest ones.

**The problem:** a normal program uses **one** core. The other fifteen are just sitting there, watching.

To take advantage of them you have to write programs that do several things at once, and that
—**concurrency**— was traditionally one of the hardest and most error-prone tasks in all of programming.

**Go was designed around that.** Making two things happen at once in Go is literally writing one word:

```go
go revisarServicio()
```

That word `go` —the one that gives the language its name— is the feature that made it famous. You will
learn it in lesson 6, and when you get there you will understand why so many people switched to this
language.

### 0.5 The timeline

| Date | What happened |
|---|---|
| **Sep 21, 2007** | Griesemer, Pike, and Thompson design Go on a whiteboard |
| mid-2008 | a working compiler exists |
| **Nov 10, 2009** | Google announces it publicly and releases it as free software |
| **Mar 28, 2012** | **Go 1.0**, the first stable version |
| 2012 onward | two major releases a year, every February and August |
| today | Go 1.27 |

🔑 **And something worth knowing: the Go 1 compatibility promise.** Since 2012, the Go team has promised
that **a program written for Go 1.0 still compiles today**, fourteen years later. They have kept it.

That is not normal. In many languages, code from five years ago no longer compiles. **In Go, what you
learn today will serve you ten years from now**, and that makes the time you are about to invest well
worth it.

### 0.6 What is written in Go (you probably already use it)

This is not a propaganda list: it is so you can see **where** this language is used, because it says a
lot about what it is for.

| | What it is |
|---|---|
| **Docker** | the tool that packages applications into containers |
| **Kubernetes** | the system that manages thousands of containers on servers. Half the industry uses it |
| **Terraform** | creates cloud infrastructure by writing files |
| **Prometheus** · **Grafana** | collect and graph server metrics |
| **Traefik** · **Caddy** | web servers and load balancers |
| **CockroachDB** · **InfluxDB** | databases |
| **Hugo** | static website generator, famous for its speed |
| **ngrok**, **rclone**, **gh** (the GitHub CLI) | command-line tools |

🔑 **Do you see the pattern?** Almost all of them are **infrastructure tools**: things that run on
servers, that have to be fast, that are deployed as a single file, and that handle many connections at the
same time. **That is where Go wins**, and it is no coincidence: it is exactly the problem Google had.

### 0.7 So, why should you learn it?

Four honest reasons:

**1. It is learned quickly.** The Go specification can be read in an afternoon. The C++ one has more than
1,800 pages. **You will be able to write useful programs in weeks, not years**, and that matters a lot
when you are starting out.

**2. It teaches you things that are useful in any language.** Being compiled and strictly typed, Go
forces you to think about types, memory, and errors. Those concepts **transfer**: if you later program in
Java, C#, or Rust, you already have them.

**3. There are jobs.** Everything that runs on modern servers has Go inside. If you are interested in
infrastructure, cloud, DevOps, or backend, it is one of the safest bets.

**4. What you learn won't expire.** Because of the Go 1 compatibility promise.

### 0.8 What you will build in this course

A command-line program called **`revisor`**: it takes a list of services, queries **all of them at the
same time**, and produces a report.

```
$ revisor --config servicios.txt --formato tabla
SERVICIO    ESTADO    TIEMPO  DETALLE
catalogo    OK         142ms  200
inventario  LENTO       2.3s  200
pagos       OK          87ms  200
reportes    FALLA          —  connection refused
```

It looks small. **It is not.** To write it well you will need everything: types, structs, errors,
interfaces, lists, concurrency, HTTP, configuration files, tests, and building an executable.

And it will grow with you: **each lesson adds one piece.** In the end you will have a program you could
actually use, not a textbook exercise.

---

## The error you will see

This lesson has no code of its own yet, so there is no compiler message to show you. But there is an
error of **thinking** that almost everyone makes when reading section 0.3, and it is worth warning you
about so it doesn't slow you down when it happens to you.

**The error:** reading the list of what Go does **not** have —classes, inheritance, exceptions, operator
overloading— and concluding *"so it's a poor language, it lacks things I need."*

**Why it is an error:** it confuses *having fewer tools* with *being able to solve fewer problems*. Go
didn't take away your ability to reuse code or handle errors: it gave you **one** way to do it instead of
ten, and you will learn that way in lessons 2 and 3. The frustration is real and normal —you will feel
it too in the first week— but it is solved by programming, not by avoiding the language. When you finish
lesson 3 you will be able to go back and read section 0.3 again and see why each line of that list is a
decision, not a shortcoming.

---

## What goes wrong

As important as knowing what a tool is for is knowing what it is not for. **No language is good at
everything**, and whoever tells you otherwise is selling you something. Using Go where it doesn't fit is
the first antipattern of this course, before you have written a single line:

| Not the best choice for | What is used instead | Why |
|---|---|---|
| iPhone or Android apps | Swift, Kotlin | the platforms are built for those |
| Web pages (what runs in the browser) | JavaScript, TypeScript | the browser only runs JavaScript |
| Data science, artificial intelligence | Python | all the libraries in the world are there |
| Large video games | C++, C# | they need absolute control of memory and hardware |
| Systems where a microsecond matters | C, C++, **Rust** | Go has a garbage collector that sometimes pauses the program |

> [!NOTE]
> 🔧 **Software engineering observation 0.2**
> The last row is the reason this course has a **second course, on Rust**. Rust solves exactly that: the
> speed of C without a garbage collector and without C's memory errors. They are tools for different
> problems, and you will learn both.

---

## Exercises

### Review questions

**0.1** What concrete problem motivated the creation of Go, and in what year?

**0.2** Name the three creators of Go and say why it is relevant who Ken Thompson is.

**0.3** What were the three requirements they were after, and why did no existing language fit?

**0.4** Mention three things Go deliberately does **not** have. What is gained by removing them?

**0.5** What changed in hardware around 2005 and what does it have to do with Go?

**0.6** What is the Go 1 compatibility promise and why does it benefit you?

**0.7** Name three programs written in Go and say what they have in common.

**0.8** Give two cases in which you would **not** use Go, and say what you would use.

### Going further

**0.9** Look online for the original talk or article *"Go at Google: Language Design in the Service of
Software Engineering"* by Rob Pike. Read the introduction and write down **one** design reason that is
not in this lesson.

**0.10** Pick **two** of the programs in section 0.6 that you don't know. Find out in one sentence what
each one does and write it down in your logbook.

**0.11** Find the table of contents of the Go language specification (*The Go Programming Language
Specification*) and count how many pages or sections it has. Compare it with the C++ standard. Write down
the two numbers: it is the most concrete way to see what "simple" means.

**0.12 (To think about, no right answer)** Go removed exceptions, classes, and inheritance —things other
languages consider indispensable. Can you think of any reason why **removing** a feature could make a
language better? Write your opinion in the logbook **before** starting the course, and read it again when
you finish. It is interesting to see whether it changed.

### Solutions

**0.1** C++ compilation times at Google: a change in a header could cost **45 minutes** of recompilation.
They started on **September 21, 2007**.

**0.2** Robert Griesemer, Rob Pike, and **Ken Thompson**. Thompson **created Unix** and the B language,
ancestor of C, which he wrote with Dennis Ritchie. That is: one of the authors of the tools on which
modern software was built also designed Go, fifty years later.

**0.3** Fast compilation, fast execution, and ease of programming. The existing languages met **two of
the three**: C++ was fast to run but slow to compile and hard; Python was easy but slow to run.

**0.4** Classes and inheritance, exceptions, operator overloading, constructors. **Readability is
gained**: Go code written by anyone can be read and understood, because there are no dialects or many
ways to do the same thing.

**0.5** Processors stopped getting faster and started multiplying **cores**. A normal program uses just
one, so a language was needed where writing concurrent programs was easy. Hence the word `go`.

**0.6** The promise that a program written for Go 1.0 (2012) **still compiles today**. It benefits you
because what you learn doesn't expire and the code you write will keep working for years.

**0.7** Docker, Kubernetes, Terraform, Prometheus, Grafana, Hugo… All of them are **infrastructure
tools**: they run on servers, are distributed as a single executable, and handle many connections at the
same time.

**0.8** Mobile apps (Swift/Kotlin), browser code (JavaScript), data science (Python), large video games
(C++), strict real-time systems (C or Rust).

**0.9-0.12** They have no single right answer: they are for your own research and reflection. Compare
what you wrote down with a classmate or in the course logbook, not with a fixed answer.

---

## How I know I got it

There is no code to compile in this lesson, so the checklist is about understanding. Check each item only
if you can do it **without looking at the text again**:

- ☐ You explain in your own words, in under a minute, why Go was born and what problem it solved.
- ☐ You name the three creators and say why Ken Thompson being one of them is not a filler fact.
- ☐ You say what fast compilation, fast execution, and ease of programming are, and why no language
  before Go had all three.
- ☐ You name, from memory, at least three things Go decided not to have, and explain what is gained by
  removing them.
- ☐ You explain the relationship between processor cores, concurrency, and the word `go`.
- ☐ You name three programs written in Go and say what they have in common.
- ☐ You give an example of a problem for which you would **not** use Go, and say what you would use
  instead.
- ☐ You solved the eight review questions without looking at the solutions before comparing.

If you missed any, don't move on to lesson 1 yet: go back and reread the corresponding section.
Everything that follows builds on this.

---

## Summary

- Go was born in **2007** at Google out of frustration with **C++ compilation times**: 45 minutes for a
  change in a header.
- It was created by **Robert Griesemer, Rob Pike, and Ken Thompson**; Thompson created Unix and co-wrote C.
- They were looking for **three** properties together —compile fast, run fast, be easy— that no
  language at the time met all at once.
- Go achieved simplicity by **removing** things: no classes, no inheritance, no exceptions, no operator
  overloading.
- The second motivation was hardware: processors multiplied **cores** instead of getting faster, and
  taking advantage of them demanded that **concurrency** be easy. Hence the word `go`.
- It was announced in **2009** and version **1.0** came out in **2012**; today it is at 1.27, with two
  major releases a year.
- The **Go 1 compatibility promise** guarantees that code from 2012 still compiles: what you learn
  doesn't expire.
- Go dominates **infrastructure**: Docker, Kubernetes, Terraform, Prometheus, Grafana.
- It is **not** the best choice for mobile, browser, data science, video games, or strict real time.

---

## Further reading

1. **The language specification** — *The Go Programming Language Specification*, at `go.dev/ref/spec`.
   It is the official source and, as you saw in exercise 0.11, it can be read in an afternoon. You don't
   need to understand all of it today; it is worth knowing that it exists and is short.
2. **The official Go blog** — `go.dev/blog`, where the team itself publishes the news of each release
   and design articles.
3. **Rob Pike, "Go at Google: Language Design in the Service of Software Engineering"** — the original
   talk/article where one of the creators explains the design decisions in the context of Google. It is
   the reading for exercise 0.9, and the most direct one for understanding the **why** behind every
   decision you saw in this lesson.

### Terms from this lesson

| | |
|---|---|
| **compile** | translate source code into machine instructions |
| **concurrency** | a program doing several things at once |
| **header** | in C/C++, a file with declarations that other files include |
| **Go 1 (compatibility promise)** | commitment that old code keeps compiling |
| **core** | each independent processor inside a single chip |
| **Turing Award** | the highest recognition in computer science |
| **garbage collector** | the part of the language that frees memory automatically |
| **UTF-8** | standard way of representing text, co-invented by Rob Pike and Ken Thompson |

---

**Next:** [Lesson 1 — Installing Go on your Linux Mint](01-instalacion.md)
