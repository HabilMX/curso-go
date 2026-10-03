# Glosario de traducciones del curso de Go (en · fr · pt-BR · bg)

Fuente única para quien traduzca `es/` a `en/`, `fr/`, `pt/` o `bg/`. Si un término está aquí, se usa tal cual.
Si un término falta, se agrega a esta tabla (con el candado `mkdir /tmp/glosario-go.lock`, avisando a los otros
traductores ANTES de usarlo) y luego se usa.

El español manda: la traducción se compara contra `es/` por huella (`registro-traducciones.tsv`).

## 1. Registro: el curso TUTEA al alumno, en masculino genérico

| Idioma | Forma | Notas |
|---|---|---|
| es (fuente) | *tú* | «tu computadora», «vas a ver». |
| en | «you / your» | Neutro, directo, sin calcos del español. |
| fr | **«tu / ton / ta / tes»** | Tuteo. Nunca «vous» dirigido al alumno. |
| pt-BR | **«você / seu / sua»** | Aquí SÍ se permite «você»: en el curso equivale al tuteo. Imperativo con «você» («Abra…», «Execute…»). |
| bg | **«ти / твой / твоя / твоите»** | Tuteo. Nunca «Вие/Ви». |

Masculino genérico: se conserva donde el español lo usa («el alumno», «el lector»); en inglés se usa «the learner / the reader / you».

## 2. Reglas generales

- **El código no se traduce.** Todo bloque ``` va BYTE A BYTE idéntico al de `es/` (identificadores, comentarios, salidas y
  errores reales). También van idénticos, sin tocar: los comentarios HTML `<!-- verificar:extracto:… -->`,
  `<!-- verificar:ejemplo:… -->`, `<!-- verificar:fragmento -->` y la primera línea `// figNN_NN.go`.
- **Código en línea** (`` `así` ``) tampoco se traduce: identificadores, comandos, rutas, nombres de archivo.
- Solo se traduce la prosa: títulos, párrafos, listas, tablas, preguntas, soluciones, resúmenes, textos de enlaces.
- **Se conserva la estructura**: mismas secciones `##` y `###` en el mismo orden y número; mismos nombres de archivo
  (`00-introduccion.md`, …, `bitacora.md`, `README.md`). Las anclas internas (`#...`) se ajustan al título traducido.
- **Se conserva TODA la profundidad.** No se resume, no se omite, no se agrega.
- Títulos: «Lesson N — …» / «Leçon N — …» / «Lição N — …» / «Урок N — …», sin emojis y sin mencionar IA. Los emojis que el es usa
  como marcadores dentro del texto (🔑 ⚠️ 🔧 ✅ ❌ 🟡 ☐) se CONSERVAN en la misma posición; no se agregan otros.
- Enlaces a documentación oficial (go.dev, pkg.go.dev, etc.): se cambian a la versión del idioma si existe; si no, se quedan.
- Los enlaces entre lecciones apuntan al archivo de su propio idioma (mismo nombre de archivo, misma carpeta).
- Cifras, nombres propios (Rob Pike, Ken Thompson, Robert Griesemer, Google, Hábil, Dorian Chávez) y versiones: tal cual.
- Los resultados de ejecución citados en la prosa (salidas, mensajes del compilador) se dejan en el idioma original en que
  Go los imprime (inglés); la prosa que los explica se traduce.
- Puntuación y comillas: las del idioma (en “ ” o " "; fr « » con espacio fino/normal; pt “ ”; bg „ “).

## 3. Nombres propios que NO se traducen

`revisor` (el programa del curso, siempre en código: `` `revisor` ``), `Go`, `Linux Mint`, `Hábil`, `bitacora`, los nombres de
archivo y paquetes. Si el español usa «revisor» como sustantivo común (fuera de código), se glosa en la primera mención de cada
archivo: en «checker» (`revisor`), fr «vérificateur» (`revisor`), pt «verificador» (`revisor`), bg «проверител» (`revisor`).

## 4. Términos técnicos y didácticos

«=» significa: no se traduce (se escribe igual; va en código en línea cuando es identificador).

| es | en | fr | pt-BR | bg |
|---|---|---|---|---|
| lección | lesson | leçon | lição | урок |
| curso | course | cours | curso | курс |
| alumno | learner | apprenant | aluno | учещ / ученик |
| ejercicio | exercise | exercice | exercício | упражнение |
| solución (de ejercicio) | solution | solution (corrigé en títulos de sección) | solução / gabarito (títulos: «Soluções») | решение |
| pregunta de repaso | review question | question de révision | pergunta de revisão | въпрос за преговор |
| Por qué importa | Why it matters | Pourquoi c'est important | Por que isso importa | Защо е важно |
| Los conceptos | The concepts | Les concepts | Os conceitos | Понятията |
| El error que vas a ver | The error you will see | L'erreur que tu vas voir | O erro que você vai ver | Грешката, която ще видиш |
| Lo que se hace mal | What goes wrong | Ce qu'on fait mal | O que se faz errado | Какво се прави погрешно |
| Ejercicios | Exercises | Exercices | Exercícios | Упражнения |
| Cómo sé que lo logré | How I know I got it | Comment savoir que j'y suis arrivé | Como sei que consegui | Как разбирам, че съм успял |
| Resumen | Summary | Résumé | Resumo | Резюме |
| Para leer más | Further reading | Pour aller plus loin | Para ler mais | За допълнително четене |
| Para ir más allá | Going further | Pour aller plus loin | Para ir além | Отвъд основите |
| Términos de esta lección | Terms from this lesson | Termes de cette leçon | Termos desta lição | Термини от този урок |
| Siguiente | Next | Suivant | Próxima | Следващ |
| bitácora | logbook | journal de bord | diário de bordo | дневник |
| compilar | to compile | compiler | compilar | компилирам |
| compilador | compiler | compilateur | compilador | компилатор |
| error de compilación | compilation error | erreur de compilation | erro de compilação | грешка при компилация |
| ejecutable | executable | exécutable | executável | изпълним файл |
| código fuente | source code | code source | código-fonte | изходен код |
| lenguaje compilado / interpretado | compiled / interpreted language | langage compilé / interprété | linguagem compilada / interpretada | компилиран / интерпретиран език |
| recolector de basura | garbage collector | ramasse-miettes | coletor de lixo | събирач на боклук (garbage collector) |
| concurrencia | concurrency | concurrence | concorrência | конкурентност |
| paralelismo | parallelism | parallélisme | paralelismo | паралелизъм |
| goroutine | goroutine | goroutine | goroutine | goroutine |
| canal (de comunicación) | channel | canal (pl. canaux) | canal (pl. canais) | канал |
| canal con búfer / sin búfer | buffered / unbuffered channel | canal avec tampon / sans tampon | canal com buffer / sem buffer | канал с буфер / без буфер |
| `select` | `select` | `select` | `select` | `select` |
| `sync.WaitGroup`, `Mutex`, `defer` | = | = | = | = |
| condición de carrera | race condition / data race | situation de compétition / data race | condição de corrida / data race | състезание за данни (data race) |
| interbloqueo | deadlock | interblocage (deadlock) | deadlock (impasse) | взаимно блокиране (deadlock) |
| paquete | package | paquet | pacote | пакет |
| módulo | module | module | módulo | модул |
| `go.mod` | = | = | = | = |
| dependencia | dependency | dépendance | dependência | зависимост |
| importar | to import | importer | importar | импортирам |
| exportado / no exportado | exported / unexported | exporté / non exporté | exportado / não exportado | експортиран / неекспортиран |
| función | function | fonction | função | функция |
| método | method | méthode | método | метод |
| receptor | receiver | receveur | receptor | получател (receiver) |
| parámetro | parameter | paramètre | parâmetro | параметър |
| argumento | argument | argument | argumento | аргумент |
| valor de retorno / devolver | return value / to return | valeur de retour / retourner | valor de retorno / retornar | върната стойност / връщам |
| variable | variable | variable | variável | променлива |
| constante | constant | constante | constante | константа |
| tipo | type | type | tipo | тип |
| valor cero | zero value | valeur zéro | valor zero | нулева стойност |
| identificador vacío (`_`) | blank identifier | identifiant vide | identificador em branco | празен идентификатор |
| `nil` | = | = | = | = |
| ausencia de valor | absence of a value | absence de valeur | ausência de valor | липса на стойност |
| puntero | pointer | pointeur | ponteiro | указател |
| struct | struct | struct | struct | struct |
| campo | field | champ | campo | поле |
| interfaz | interface | interface | interface | интерфейс |
| satisfacción implícita | implicit satisfaction | satisfaction implicite | satisfação implícita | имплицитно удовлетворяване |
| error (valor `error`) | error | erreur | erro | грешка |
| error centinela | sentinel error | erreur sentinelle | erro sentinela | сигнален (sentinel) error |
| envolver un error | to wrap an error | envelopper une erreur | encapsular (embrulhar) um erro | обвивам грешка |
| desenvolver | to unwrap | déballer | desencapsular | разопаковам |
| panic | panic | panic | panic | panic |
| slice | slice | slice (tranche en prosa solo si el es la explica) | slice | slice |
| arreglo | array | tableau | array | масив |
| map | map | map | map | map |
| llave (de un map) | key | clé | chave | ключ |
| largo / capacidad (`len`, `cap`) | length / capacity | longueur / capacité | tamanho / capacidade | дължина / капацитет |
| rebanar (cortar un slice) | to slice | découper | fatiar | нарязвам |
| cadena de texto | string | chaîne de caractères | string | низ (string) |
| entero | integer | entier | inteiro | цяло число |
| número de punto flotante | floating-point number | nombre à virgule flottante | número de ponto flutuante | число с плаваща запетая |
| booleano | boolean | booléen | booleano | булев |
| verbo de formato | formatting verb | verbe de formatage | verbo de formatação | форматиращ глагол (verb) |
| plantilla (de `Printf`) | template (format string) | modèle (chaîne de format) | modelo (string de formato) | шаблон |
| salida estándar / de errores | standard output / standard error | sortie standard / erreur standard | saída padrão / saída de erro | стандартен изход / изход за грешки |
| argumentos de línea de comandos | command-line arguments | arguments de ligne de commande | argumentos de linha de comando | аргументи от командния ред |
| bandera (`flag`) | flag | option (flag) | flag | флаг |
| subcomando | subcommand | sous-commande | subcomando | подкоманда |
| prueba (automatizada) | test | test | teste | тест |
| prueba de tabla | table-driven test | test piloté par table | teste orientado a tabela | тест, управляван от таблица |
| subprueba | subtest | sous-test | subteste | подтест |
| cobertura | coverage | couverture | cobertura | покритие |
| doble de prueba / simulacro | test double / fake | doublure de test / faux | dublê de teste / fake | двойник за тест (fake) |
| servidor de prueba | test server | serveur de test | servidor de teste | тестов сървър |
| núcleo (procesador) | core (processor core) | cœur | núcleo | ядро |
| encabezado (C/C++) | header | en-tête | cabeçalho | заглавен файл (header) |
| promesa de compatibilidad (Go 1) | compatibility promise | promesse de compatibilité | promessa de compatibilidade | обещание за съвместимост |
| Premio Turing | Turing Award | prix Turing | Prêmio Turing | награда «Тюринг» |
| terminal | terminal | terminal | terminal | терминал |
| línea de comandos | command line | ligne de commande | linha de comando | команден ред |
| archivo / carpeta | file / folder | fichier / dossier | arquivo / pasta | файл / папка |
| ruta | path | chemin | caminho | път |
| variable de entorno | environment variable | variable d'environnement | variável de ambiente | променлива на средата |
| extracto (de código real) | excerpt | extrait | trecho | откъс |
| fragmento | fragment | fragment | fragmento | фрагмент |
| programa completo | complete program | programme complet | programa completo | пълна програма |

## 5. Registro de términos nuevos

(Agrega aquí, bajo candado, los términos que no estaban: `es | en | fr | pt-BR | bg | quién lo agregó`.)


| es | en | fr | pt-BR | bg | quién |
|---|---|---|---|---|---|
| antipatrón | antipattern | anti-patron | antipadrão | антипатърн (antipattern) | trad-bg |
| Buena práctica (recuadro) | Good practice | Bonne pratique | Boa prática | Добра практика | trad-en |
| Error común de programación (recuadro) | Common programming error | Erreur courante de programmation | Erro comum de programação | Често срещана грешка при програмиране | trad-en |
| Observación de ingeniería de software (recuadro) | Software engineering observation | Observation de génie logiciel | Observação de engenharia de software | Наблюдение от софтуерното инженерство | trad-en |
| Tip de portabilidad (recuadro) | Portability tip | Astuce de portabilité | Dica de portabilidade | Съвет за преносимост | trad-en |
| Tip de prueba y depuración (recuadro) | Testing and debugging tip | Astuce de test et de débogage | Dica de teste e depuração | Съвет за тестване и отстраняване на грешки | trad-en |
| Tip de rendimiento (recuadro) | Performance tip | Astuce de performance | Dica de desempenho | Съвет за производителност | trad-en |
| shell de *login* | login shell | shell de connexion (login shell) | shell de login | login shell (обвивка при вход) | trad-bg |
| binario (ejecutable) | binary | binaire | binário | двоичен файл (binary) | trad-bg |
| enlazado estático | static linking | liaison statique | linkagem estática | статично свързване (static linking) | trad-bg |
| tabla de casos (patrón de prueba) | table of cases | table de cas | tabela de casos | таблица със случаи | trad-fr (aviso: en uso desde la lección 5) |
| tubería (`\|` del shell) | pipe / pipeline | tube (pipe) | pipe | конвейер (pipe) | trad-fr (aviso: en uso desde la lección 4) |
| código de salida | exit code | code de sortie | código de saída | код на изход | trad-fr (aviso: en uso desde la lección 4) |
| biblioteca estándar | standard library | bibliothèque standard | biblioteca padrão | стандартна библиотека | trad-fr (aviso: en uso desde la lección 5) |
| huella criptográfica | cryptographic hash (checksum) | empreinte cryptographique | hash criptográfico | криптографски отпечатък | trad-fr (aviso: en uso desde la lección 5) |
| benchmark | benchmark | benchmark | benchmark | бенчмарк (benchmark) | trad-fr (aviso: en uso desde la lección 5) |
| prueba de frontera | boundary test | test aux limites | teste de fronteira | гранични тестове | trad-fr (aviso: en uso desde la lección 5) |
| tiempo límite (timeout) | timeout | délai d'expiration (timeout) | tempo limite (timeout) | времево ограничение (timeout) | trad-fr (aviso: en uso desde la lección 5) |
| semáforo (de canal) | semaphore | sémaphore | semáforo | семафор | trad-fr (aviso: en uso desde la lección 6) |
| carrera de datos | data race | data race (situation de compétition) | data race (condição de corrida) | състезание за данни (data race) | trad-fr (aviso: en uso desde la lección 6) |
| cancelación (`context`) | cancellation | annulation | cancelamento | отмяна | trad-fr (aviso: en uso desde la lección 6) |
| serializar (a JSON) | to serialize | sérialiser | serializar | сериализирам | trad-fr (aviso: en uso desde la lección 7) |
| punto de entrada (`main`) | entry point | point d'entrée | ponto de entrada | входна точка | trad-fr (aviso: en uso desde la lección 7) |
| compilación cruzada | cross-compilation | compilation croisée | compilação cruzada | кръстосана компилация | trad-fr (aviso: en uso desde la lección 7) |
| detector de carreras (`-race`) | race detector | détecteur de data races | detector de corridas (race detector) | детектор на състезания (race detector) | trad-en (aviso: en uso desde la lección 5) |
| fuga de memoria | memory leak | fuite de mémoire | vazamento de memória | изтичане на памет | trad-en (aviso: en uso desde la lección 6) |
| candado (mutex / `Lock`) | lock | verrou | trava (lock) | заключване (lock) | trad-en (aviso: en uso desde la lección 6) |
| recorrido (de un map/slice con `range`) | iteration | parcours | iteração | обхождане | trad-en (aviso: en uso desde la lección 4) |
| integración continua | continuous integration (CI) | intégration continue | integração contínua | непрекъсната интеграция | trad-en (aviso: en uso desde la lección 6) |
| etiqueta (de struct, `json:"..."`) | struct tag | tag de struct | tag de struct | таг на struct | trad-en (aviso: en uso desde la lección 7) |
| drenar (el cuerpo de una respuesta) | to drain | vider (drainer) | drenar | изчерпвам (drain) | trad-en (aviso: en uso desde la lección 7) |
| símbolos / información de depuración | debug symbols / debugging information | symboles / informations de débogage | símbolos / informações de depuração | символи / информация за дебъгване | trad-en (aviso: en uso desde la lección 7) |
| autoridad certificadora | certificate authority | autorité de certification | autoridade certificadora | удостоверяващ орган (CA) | trad-en (aviso: en uso desde la lección 7) |
| receptor de valor / de puntero | value receiver / pointer receiver | receveur valeur / pointeur | receptor de valor / de ponteiro | получател по стойност / по указател | trad-bg (aviso: en uso en bg/03) |
| hilo (del sistema operativo) | thread (OS thread) | thread (fil d'exécution) | thread | нишка (thread) | trad-pt (aviso: en uso desde la lección 6) |
| traza (de pila, de un panic) | stack trace | trace d'appels (stack trace) | stack trace | трасиране на стека (stack trace) | trad-pt (aviso: en uso desde la lección 6) |
