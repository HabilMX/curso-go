# Урок 7 — Завършената програма

**Продължителност:** 90-120 минути или 2 сесии по 60. Това е последният урок за `revisor`: като приключиш, имаш
двоичен файл, който наистина работи, настройва се от командния ред, говори HTTP с истински времеви ограничения (timeouts),
създава два формата на изхода и се компилира за друга платформа, без да излизаш от машината си.

**Като приключиш, ще можеш:**

- Да напишеш HTTP клиент с изрично времево ограничение и да обясниш защо `http.DefaultClient` е опасен в
  продукция.
- Да отделиш логиката на програма (`ejecutar`) от нейната входна точка (`main`), за да можеш да я тестваш, без
  да докосваш `os.Exit`.
- Да дефинираш флагове от командния ред със стандартния пакет `flag`, без никаква външна библиотека.
- Да сериализираш структура в JSON с `encoding/json` и да обясниш кои полета се губят и защо.
- Да пуснеш `revisor` от началото до края срещу истински сървър (тестов, направен от теб) и да прочетеш
  изхода му.
- Да компилираш същия код за друга архитектура и друга операционна система, без да излизаш от машината си, и
  да обясниш защо това е възможно.

---

## Защо е важно

Вече имаш четирите части на `revisor`, изградени поотделно: `servicio` дефинира речника
(урок 2-3), `config` чете конфигурацията (урок 5), `revisar.Todos` проверява всичко едновременно (урок
6). Единственото, което липсва, е онова, което превръща тези части в **програма, която някой друг може да използва, без
да чете изходния код**: да говори HTTP наистина (досега го тествахме само с `Falso`), да получава
флагове от терминала вместо фиксирани стойности в кода, да създава формат, който друга
програма може да консумира, и да може да се компилира и разпространява като един-единствен файл.

Това е урокът, в който `revisor` престава да бъде "код, който работи, ако го пусна аз, на моята машина, с
моите тестови данни", и се превръща в двоичен файл, който можеш да копираш на някой друг с увереност, че ще
направи точно онова, което командният ред му поиска.

---

## Понятията

### 7.1 HTTP клиентът: защо никога `http.DefaultClient`

Ето как изглежда истинският `Revisor`, онзи, който наистина докосва мрежата (`programas/revisor/internal/revisar/revisar.go`):

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

🔴 **`http.DefaultClient` — онзи, който би използвал, ако напишеш направо `http.Get(url)` — няма никакво
времево ограничение.** Ако сървърът от другата страна приеме връзката и никога не отговори нищо, това извикване остава
да чака **завинаги**, без грешка, без panic, само програма, която един ден спира да напредва и никой
не знае защо. Това е една от най-скъпите грешки в Go в продукция точно защото не дава никакъв
симптом, докато вече не се случва.

`revisor` се защитава два пъти, а не веднъж: собственият `http.Client` има общо времево ограничение (ограничението на
целия отчет, подадено на `NuevoHTTP`), а освен това всяка отделна проверка се изпълнява под
`context.WithTimeout` за всяка услуга, който `Todos` сглоби в урок 6 — второто, по-кратко, е онова, което
командва на практика. Времевото ограничение на клиента е предпазна мрежа в краен случай, а не основният
механизъм.

⚠️ **`defer resp.Body.Close()` идва след проверката на грешката, никога преди нея.** Ако `err != nil`, `resp`
е `nil` и извикването на метод върху нулев указател предизвиква panic. А **`io.Copy(io.Discard, resp.Body)`
преди затварянето** не е украса: без да изчерпиш (drain) тялото на отговора, съответната TCP връзка не
може да се преизползва за следващата заявка към същия сървър и програмата накрая отваря
нова връзка всеки път, вместо да преизползва онези, които вече има.

**`traducir` заменя суровата грешка на `net/http` с такава, която се чете в таблица:**

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

Грешката, която `net/http` дава по подразбиране, носи целия URL и думата `Get`, които в таблица на
отчет само заемат място и не казват нищо ново на читателя — затова се превежда, преди да бъде показана.

### 7.2 `main` не се тества; `ejecutar` — да

Шаблонът, който отличава `revisor` от програма от един файл:

<!-- verificar:extracto:cmd/revisor/main.go -->
```go
func main() {
	os.Exit(ejecutar(os.Args[1:], os.Stdout, os.Stderr))
}
```

`ejecutar` е мястото, където живее истинската логика — раздели 7.3, 7.4 и 7.6 я показват на части —; ето
сигнатурата ѝ:

<!-- verificar:fragmento -->
```go
func ejecutar(args []string, salida, errores *os.File) int {
	// toda la lógica real vive aquí
}
```

**`go test` не може да тества функция, която извиква `os.Exit`**, защото `os.Exit` прекратява целия
процес незабавно — включително самия процес на тестовете, който никога не би стигнал до докладване на резултата.
Като се отдели "решаването кой код на изход съответства" (`ejecutar`, която **връща** `int`) от "истинското
излизане с този код" (`main`, единственият ред, който извиква `os.Exit`), цялата логика може да се тества
чрез директно извикване на `ejecutar`, с тестови аргументи, и проверка на числото, което връща — точно
това прави `cmd/revisor/main_test.go`:

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

Това е истинският изход на `--help` (генериран автоматично от `flag`, раздел 7.3), уловен по средата
на тест, без да се отваря терминал и без да се изпълнява компилираният двоичен файл.

### 7.3 Флагове от командния ред, без библиотеки

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

**`flag.NewFlagSet` вместо пакета `flag` на глобално ниво** (който би използвал с `flag.String(...)`
направо) е онова, което позволява да има функция `ejecutar(args []string, ...)`, която получава аргументите си
като параметър, вместо винаги да ги чете от `os.Args` — още един детайл, който съществува специално, за да
може тя да се тества, както в раздел 7.2.

`flag` идва в стандартната библиотека и стига за програма с такъв размер: четири флага, основни
типове (`string`, `int`, `time.Duration`), автоматична помощ. Ако някой ден `revisor` има нужда от
подкоманди (`revisor check`, `revisor list`, всяка със собствени флагове), тогава вече би си струвало
да се погледне библиотека като `cobra` — но не преди наистина да има нужда от нея.

### 7.4 Двата формата на изхода: таблица и JSON

Ето как изглежда `Tabla`, функцията, която създава четимия за хора изход, който виждаш през целия урок
(`programas/revisor/internal/reporte/tabla.go`):

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

**Ширината на първата колона не е фиксирана в кода — тя се изчислява.** Цикълът отгоре обхожда
всички състояния веднъж, преди да отпечата каквото и да е, и запазва най-дългото име (като започва от
ширината на самата дума "SERVICIO", в случай че някое име на услуга е по-кратко от нея). Това
число влиза като `*` в `%-*s`: форматиращ глагол с **променлива** ширина, при който самата ширина е още един
аргумент на `Fprintf`, а не число, написано на ръка. Това е причината истинската таблица в този
урок да има прави колони, независимо дали имената на услугите са кратки (`sano`), или дълги
(`inventario`): с фиксирана ширина един от двата случая винаги би изглеждал крив.

`ordenarPorNombre` (раздели 6.3 и 6.7 вече обясниха защо на реда на пристигане не може да се разчита: няколко
goroutines доставят резултати по канал и печели онзи, който отговори пръв) прави копие на slice-а и го
подрежда по азбучен ред, преди да отпечата — без това една и съща проверка би създала таблица в
различен ред всеки път, когато пуснеш програмата, дори данните да са идентични.

`etiquetaEstado`, `formatoTiempo` и `detalle` са трите малки функции, които решават коя дума отива във
всяка колона (`OK`/`LENTO`/`FALLA`, времето, закръглено до милисекунда, или тире, ако никога не е имало
отговор, и HTTP кодът или причината за провала) — те са онези, които вече видя в действие в истинската таблица от
раздел 7.5, сега с кода, който ги създава.

🔑 **Защо не беше използван `text/tabwriter`, който е инструментът, предлаган от стандартната библиотека
точно за подравняване на колони:** за таблица с четири колони, в която само една има променлива ширина
(името на услугата; другите три са кратки и предвидими), изчисляването на ширината на ръка е по-лесно
за четене, отколкото въвеждането на `tabwriter.Writer` със собствените му `Flush()` и разделители с табулация. При
таблица с повече колони с променлива ширина `tabwriter` наистина би бил правилният инструмент — струва си да го познаваш
(той е в раздела "За допълнително четене" на този урок), макар че `revisor` няма нужда от него.

След като таблицата е разбрана, другият формат на изхода — JSON — е другата половина на този раздел:

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

⚠️ **`servicio.Estado` не се сериализира директно — нарочно.** `Estado` носи поле `Err error`, а
`error` е интерфейс: `encoding/json` не знае как да го превърне в текст сам (опитът създава
празен обект `{}`, а не грешка при компилация, което е по-лошо: проваля се тихо). `revisor` решава
това с междинен тип, `lineaJSON`, който живее само в пакета `reporte`: превръща грешката
в нейния текст (`e.Motivo()`) преди кодирането и мимоходом разделя публичния формат (онова, което други
програми ще четат) от вътрешния модел (`Estado`) — ако утре `Estado` получи ново поле, JSON-ът, който
вече циркулира, не се променя само защото се е променило нещо вътрешно.

`omitempty` при `Codigo` и `Detalle` маха тези полета от JSON, когато са нула или празен низ — така
здравата услуга не влачи безсмислено `"detalle": ""`. Наистина кодирано:

<!-- verificar:extracto:internal/reporte/json.go -->
```go
codificador := json.NewEncoder(w)
codificador.SetIndent("", "  ")
return codificador.Encode(lineas)
```

### 7.5 От край до край: `revisor`, работещ срещу истински сървър

За да тествам цялата програма заедно — а не всяка част поотделно — изградих `cmd/servidor-demo`:
минимален HTTP сървър с четири маршрута, които се държат като четирите случая, които `revisor` трябва
да умее да докладва.

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

С този сървър, пуснат на `:8091`, и конфигурационен файл, сочещ към четирите му маршрута, пуснах
истинския двоичен файл:

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

**Всеки ред от тази таблица съответства точно на поведението, програмирано в тестовия
сървър:** `sano` отговаря бързо и с 200; `lento` наистина се бави 1.5 секунди и затова излиза `LENTO`;
`malo` отговаря мигновено, но с 500; `colgado` никога не отговаря, а собственото му времево ограничение от 200ms
— декларирано в конфигурационния файл, колона три — прекъсна чакането на 202ms, почти точно. **Кодът
на изход, 1, е истински**: поне една услуга не е била `OK`, така че `ejecutar` връща 1 вместо
0 (раздел 7.6) — същият двоичен файл, използван от скрипт, може да каже на онзи, който го извиква, дали е имало
проблеми, без някой да трябва да чете таблицата.

А в JSON — същият отчет:

```
$ /tmp/revisor --config /tmp/servicios-demo.txt --formato json
[
  {"servicio": "colgado", "ok": false, "tiempo_ms": 205, "detalle": "se acabo el tiempo de espera"},
  {"servicio": "lento", "ok": true, "codigo": 200, "tiempo_ms": 1503},
  {"servicio": "malo", "ok": false, "codigo": 500, "tiempo_ms": 1, "detalle": "codigo 500"},
  {"servicio": "sano", "ok": true, "codigo": 200, "tiempo_ms": 1}
]
```

(Тук без отстъпа на `SetIndent`, за да се побере на един ред за всяка услуга; истинската програма го
създава с нови редове, както видя в раздел 7.4.)

### 7.6 Кодът на изход, със значение

<!-- verificar:extracto:cmd/revisor/main.go -->
```go
for _, e := range estados {
	if !e.OK() {
		return 1 // al menos un servicio falló: el código de salida lo refleja
	}
}
return 0
```

Три възможни кода и всеки отговаря на различен въпрос за онзи, който извиква програмата от
скрипт или pipeline: **0** — всичко мина добре; **1** — програмата се изпълни докрай, но поне една
услуга не беше здрава; **2** — програмата дори не успя да стартира (липсва `--config`, файлът не
съществува, форматът на конфигурацията е написан погрешно). Това е разликата между "свърших работата и намерих
проблеми" и "не можах дори да започна работа" — и pipeline за непрекъсната интеграция може да реагира
различно на всеки от тях.

### 7.7 Кръстосана компилация: същият код, друга платформа

```bash
GOOS=linux GOARCH=amd64 go build -o revisor-linux-amd64 ./cmd/revisor
```

Потвърдено наистина, с компилиране от този Mac (arm64) към Linux x86-64:

```
$ GOOS=linux GOARCH=amd64 go build -o revisor-linux-amd64 ./cmd/revisor
$ file revisor-linux-amd64
revisor-linux-amd64: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, ...
```

**Без Docker, без виртуална машина, без да се инсталира нищо допълнително: две променливи на средата стигат.** Това е
възможно, защото компилаторът на Go носи по подразбиране кода, нужен за създаване на двоични файлове за
всяка комбинация от операционна система и архитектура, която поддържа — не зависи от това да работи на
целевата система, за да компилира за нея, обратно на начина, по който работят много други компилирани езици.

**Размерът на двоичния файл, измерен, със и без символи за дебъгване:**

```
$ go build -o revisor-normal ./cmd/revisor
$ wc -c revisor-normal
9464130 revisor-normal

$ go build -ldflags="-s -w" -o revisor-chico ./cmd/revisor
$ wc -c revisor-chico
6386194 revisor-chico
```

**От 9.46 MB на 6.39 MB, истинско намаление от 32%**, като се махнат таблицата със символи (`-s`) и
информацията за дебъгване на DWARF (`-w`), които нормалният двоичен файл включва, за да може дебъгер да
го инспектира. За двоичен файл, който ще разпространяваш в продукция и няма да дебъгваш точно там, тези данни
не служат за нищо и само заемат място — за такъв, който активно разработваш, ги запази.

### 7.8 Тестване на код, който прави HTTP, без истинска мрежа: `httptest`

До урок 6 тестовете на конкурентност използваха `Falso` (`Revisor`, който никога не докосва мрежата). За
да се тества **цялата** програма — флагове, четене на конфигурацията, истинският HTTP клиент, отчетът —
без зависимост от външна услуга, нито от това нещо да слуша на фиксиран порт, `net/http/httptest`
пуска истински сървър, на порт, който операционната система назначава сама, вътре в самия процес на
тестовете:

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

Истинско пускане:

```
$ go test ./cmd/revisor/... -run TestEjecutar_reportaOKyFalla -v
=== RUN   TestEjecutar_reportaOKyFalla
--- PASS: TestEjecutar_reportaOKyFalla (0.00s)
```

**`servidor.URL` е истински URL** (`http://127.0.0.1:PUERTO`, със свободен порт, който системата е избрала),
така че HTTP `Revisor` от раздел 7.1 го проверява точно така, както би проверил всяка
услуга в продукция — единствената разлика е, че "услугата" е `http.HandlerFunc` от петнадесет
реда, който живее в същия процес като теста. Това е, което позволява да се тества от край до край
— флагове, конфигурация, HTTP клиент, отчет, код на изход — за милисекунди, без да се отваря нито един
фиксиран порт, който би могъл да се сблъска с нещо друго, работещо на машината, и без да остава нищо да работи след
приключването на теста (`defer servidor.Close()` се грижи за това).

**`t.TempDir()`** създава временна папка, която Go изтрива сам в края на теста (за разлика от
голото `os.CreateTemp("/tmp", ...)`, което би оставило файла там завинаги, ако никой не го изтрие на ръка) —
комбинацията от `httptest` за мрежата и `t.TempDir()` за диска е онова, което позволява
`TestEjecutar_reportaOKyFalla` да не оставя никаква следа в системата, след като се изпълни.

### 7.9 Защо две времеви ограничения, а не едно

`revisor` приема `--limite` (максималното време за **целия** отчет) и всеки ред от конфигурацията
може да декларира собствено времево ограничение **за всяка услуга**. Това не е излишно: те решават два различни
въпроса. `--limite` отговаря на *„колко време най-много съм готов да чакам целия
отчет?“* — полезно, ако `revisor` работи вътре в друг процес със собствен срок (проверка на здравето,
която оркестратор очаква на всеки определен период, например). Времето за всяка услуга отговаря на *„колко е
разумно да се чака ТАЗИ конкретна услуга?“* — вътрешна услуга с ниска латентност и друга, която минава
към външен доставчик, не би трябвало да споделят едно и също ограничение.

Демонстрацията от раздел 7.5 го потвърждава с истински данни: `colgado` имаше собствено ограничение от
200ms, декларирано в конфигурационния файл, и беше прекъсната на 202ms — много преди общият `--limite`
(10 секунди по подразбиране) да има възможност да се намеси. По-краткото от двете ограничения е
винаги онова, което командва, и точно това искаш: една-единствена бавна услуга да не изяде целия
бюджет от време на целия отчет.

### 7.10 Двоичният файл "няма нужда от нищо" — освен от едно нещо: сертификати

Урок 1 показа, че двоичен файл на Go работи на Linux машина без инсталиран Go, благодарение на
статичното свързване. Изкушаващо е тази идея да се разшири до "няма нужда от нищо от системата, точка" — и го проверих,
като я доведох до крайност: сложих двоичния файл на `revisor` в Docker образ, изграден `FROM scratch`,
най-празната основа, която съществува (няма дори обвивка на операционна система, само двоичния файл).

```dockerfile
FROM scratch
COPY revisor /revisor
COPY servicios.txt /servicios.txt
ENTRYPOINT ["/revisor", "--config", "/servicios.txt"]
```

С истинска услуга, сочеща към `https://example.com`:

```
$ docker build -t revisor-demo:scratch .
$ docker run --rm revisor-demo:scratch --formato tabla
SERVICIO  ESTADO    TIEMPO  DETALLE
ejemplo   FALLA      176ms  no responde: Get "https://example.com": tls: failed to verify certificate: x509: certificate signed by unknown authority
```

**Провали се — и не заради `revisor`, а заради нещо, което никой образ `scratch` не носи: списъка с удостоверяващи
органи (CA).** За да валидира TLS сертификат (катинарът на `https://`), операционната система
обикновено предоставя списък на онези, които имат право да подписват валидни сертификати — обикновено в
`/etc/ssl/certs/`. Образът `scratch` няма дори тази папка, така че Go не може да провери нито един
сертификат и отказва да продължи, с пълно право: да продължи без проверка би означавало да приеме всеки
сертификат, валиден или фалшив.

**Решението, проверено, е да се копира само този файл** от образ, който го има, без да се влачи
останалата част от операционната система:

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

**14.7 MB общо, срещу 14.4 MB без сертификатите — по-малко от 300 KB, за да се реши проблемът
правилно.** Двоичният файл наистина е самодостатъчен за всичко, което е *код на Go*; онова, което никога не носи
сам, е доверието кои органи са легитимни, защото това е информация за света (кои
удостоверяващи органи съществуват и все още са валидни), а не нещо, което компилаторът може да реши вместо теб.

---

## Грешката, която ще видиш

| Симптомът | Дословно съобщение / доказателство | Какво става и какво да направиш |
|---|---|---|
| Липсва задължителният флаг | `revisor: falta --config`, последвано от автоматичната помощ за употреба на `flag` | `ejecutar` проверява изрично, че `--config` не е празен, преди да направи каквото и да е друго, и излиза с код 2 |
| Невалиден формат на изхода | `revisor: --formato debe ser "tabla" o "json", no "xml"` | Съществуват само два формата; всяка друга стойност се отхвърля, преди да се опита каквото и да е |
| Услуга, която никога не отговаря | В таблицата: `FALLA` с детайл `se acabo el tiempo de espera`, време, близко до настроеното времево ограничение | `context.WithTimeout` за всяка услуга (урок 6) е прекъснал чакането; това не е бъг, а механизмът, който работи |
| Затваряне на тялото на нулев отговор | (ако се направи погрешно) `panic: runtime error: invalid memory address or nil pointer dereference` | Би се случило, ако `defer resp.Body.Close()` се сложи преди проверката на `err`; `revisor` го избягва, като проверява грешката първо (раздел 7.1) |
| `error`, сериализиран директно в JSON без превод | `{}` (празен обект, без никакви полезни данни, без никаква грешка при компилация, която да предупреди) | `encoding/json` не знае как да превърне интерфейс `error`; затова `reporte` използва `lineaJSON` с причината, вече превърната в текст (раздел 7.4) |
| Непроверен код на изход в скрипт | Скриптът продължава, сякаш нищо не е станало, макар че някоя услуга се е провалила | `ejecutar` наистина различава 0/1/2 (раздел 7.6); който интегрира `revisor` в pipeline, трябва да чете `$?`, а не само отпечатания изход |
| Двоичен файл в образ `FROM scratch`, срещу HTTPS | `no responde: Get "https://...": tls: failed to verify certificate: x509: certificate signed by unknown authority` | Липсва списъкът с удостоверяващи органи (`ca-certificates.crt`); копирай го от образ, който го има (раздел 7.10) |

---

## Какво се прави погрешно

- **Използване на `http.DefaultClient` или `http.Get` директно в продукция.** Раздел 7.1: без собствено времево ограничение
  връзка, която никога не отговаря, увисва програмата завинаги, без никакъв предварителен симптом.
- **Слагане на цялата логика в `main` и извикване на `os.Exit` по средата на кода.** Това прави функцията
  невъзможна за тестване с `go test`, защото `os.Exit` прекратява процеса на тестовете заедно с
  програмата. Отдели "решаването за кода на изход" от "истинското излизане" (раздел 7.2).
- **Сериализиране на `error` директно в JSON с надеждата, че "нещо ще излезе".** Излиза празен обект, без никакво
  предупреждение, че нещо не е могло да се преобразува — най-трудно откриваемият вид тих провал, защото
  програмата нито гърми, нито се оплаква.
- **Затваряне на тялото на HTTP отговор преди проверката на грешката.** Ако заявката се е провалила, `resp` е
  `nil` и извикването на метод върху него предизвиква panic. Проверявай грешката първо, винаги.
- **Неизчерпване на тялото на отговора преди затварянето му.** Връзката не може да се преизползва и
  програма, която прави много заявки към същия сървър, накрая отваря нова връзка всеки път, по-бавно
  от необходимото, без никаква грешка, която да го посочи.
- **Пренебрегване на кода на изход на двоичния файл от скрипт или pipeline.** `revisor` различава
  "изпълних се добре, всичко е здраво" (0), "изпълних се добре, нещо се провали" (1) и "не можах дори да стартирам" (2) — да се похаби това,
  като се чете само отпечатаният изход, означава да изхвърлиш информация, която програмата вече ти дава безплатно.

---

## Упражнения

1. Реализирай (или прегледай, ако вече го имаш) истинския HTTP `Revisor` от раздел 7.1, със собствен клиент
   и времево ограничение. Потвърди, че се компилира и че `go vet ./...` не се оплаква.
2. Пусни `cmd/servidor-demo` на свободен порт и пусни своя `revisor` срещу четирите му маршрута, както в
   раздел 7.5. Постави истинската таблица, която получи, а не измислена.
3. Добави към `Tabla` (раздел 7.4) пета колона, `PROTOCOLO`, която казва `https` или `http` според
   `Servicio.EsSeguro()` (методът, който беше дефиниран в урок 3). Задай фиксираната ширина на тази колона
   на ръка, без да ти трябва динамичното изчисление, което вече има `anchoNombre`.
4. Пусни същия отчет с `--formato json` и провери, че резултатът е легитимен JSON (например
   прекарай го през `jq .` или го постави във валидатор за JSON).
5. Предизвикай грешката с невалиден `--formato` (раздел 7.6, таблица) и потвърди кода на изход с
   `echo $?`.
6. Компилирай своя `revisor` за `GOOS=linux GOARCH=amd64` от текущата си машина и потвърди с `file`, че
   полученият двоичен файл е за правилната платформа.
7. (Малко по-трудно) Измери размера на двоичния си файл със и без `-ldflags="-s -w"`, както в раздел
   7.7, и изчисли процента на намаление.
8. (Закриване на курса) Пусни пълния набор от проверки на проекта — `go vet ./...`, `gofmt -l .` и
   `go test ./... -race -cover` — и постави целия изход в дневника си. Ако нещо не е зелено,
   поправи го, преди да смяташ курса за приключен.

### Решения

1-7. Няма един-единствен еталонен изход, защото зависи от собствената ти реализация и от собствената ти
машина; сравни формата на резултата си със съответните раздели (7.1, 7.4, 7.5, 7.6, 7.7).

8. Истинското пускане, върху `revisor` от този курс, на 30 септ. 2026 г.:

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

   `cmd/servidor-demo` излиза с 0.0% без `ok` или `FAIL`, но не защото нещо е счупено: с `-cover`
   пакет без нито един файл `_test.go` все пак се инструментира и се докладва, но няма нито един тест,
   който да упражнява този код. Това е помощен инструмент за тестване на цялата програма, а не част от
   `revisor`, който се разпространява, и няма собствена логика, която да решава каквото и да е — само отговаря онова, което е
   програмиран да отговаря — затова не се нуждае от собствени тестове.

   ⚠️ **61.9% на `internal/revisar` в тази таблица е артефактът, който обяснява раздел 5.7: по
   пакети не се брои онова, което `TestEjecutar_reportaOKyFalla` (в `cmd/revisor`, малко по-нагоре в същото това
   пускане) упражнява от `Revisar`, когато наистина говори HTTP срещу сървъра `httptest`.** Измерено
   с `-coverpkg=./...` върху целия проект: `Revisar` се качва на 86.4%, а общото за проекта е 81.5%,
   а не 61.9%. Ако ще цитираш число за покритие, за да решиш нещо, пусни `-coverpkg=./...` и не се доверявай
   само на числото по пакети.

---

## Как разбирам, че съм успял

- [ ] Моят HTTP `Revisor` има собствен клиент, с времево ограничение, и никога не използва `http.DefaultClient`.
- [ ] Отделих `main` (който само извиква `os.Exit`) от функция, която върши работата и връща `int`.
- [ ] Моят двоичен файл приема `--config`, `--formato`, `--paralelo` и `--limite` от командния ред.
- [ ] Пуснах истинския `revisor` срещу истински сървър (моя или `servidor-demo` от курса) и таблицата,
      която получих, отразява точно онова, което този сървър прави.
- [ ] `--formato json` създава валиден JSON, проверен с инструмент, а не с моето четене.
- [ ] `echo $?` след пускането на `revisor` дава 0, 1 или 2 и знам как да обясня всеки от тях.
- [ ] Компилирах за друга платформа (`GOOS`/`GOARCH`) и потвърдих с `file`, че двоичният файл е правилният.
- [ ] `go vet ./...`, `gofmt -l .` и `go test ./... -race -cover` са чисти в моето собствено копие на
      проекта — не го предполагам, пуснах го.

---

## За допълнително четене

1. [net/http package](https://pkg.go.dev/net/http) — пълната официална документация на HTTP клиента и
   сървъра от стандартната библиотека.
2. [Command flag](https://pkg.go.dev/flag) — официалната документация на пакета за флагове, използван в
   този урок.
3. [encoding/json: JSON and Go](https://go.dev/blog/json) — официалната статия за това как Go превежда
   между structs и JSON, включително правилата за таговете (`json:"..."`, `omitempty`, `-`).
4. [Build constraints и кръстосана компилация](https://pkg.go.dev/cmd/go#hdr-Environment_variables) —
   справочникът за `GOOS`/`GOARCH` и останалите променливи на средата, които управляват `go build`.
5. [text/tabwriter](https://pkg.go.dev/text/tabwriter) — стандартният инструмент за подравняване на колони,
   когато ръчното изчисление от раздел 7.4 не стига (няколко колони с променлива ширина).

### Термини от този урок

| Термин | Какво означава |
|---|---|
| `http.Client` | типът, който прави HTTP заявки в Go; ненастроен, няма времево ограничение |
| `flag.FlagSet` | набор от флагове от командния ред, който може да се разбере от собствен slice, а не само от `os.Args` |
| `omitempty` | опция на JSON тага, която маха поле от изхода, когато то е нула или е празно |
| код на изход | цялото число, което програмата връща при приключване; 0 означава успех по универсална конвенция |
| `GOOS` / `GOARCH` | променливи на средата, които казват на `go build` за коя операционна система и архитектура да компилира |
| `-ldflags="-s -w"` | опции за компилация, които махат символите и данните за дебъгване, за да намалят размера на двоичния файл |

---

## И сега — какво следва

Имаш пълна програма: проверява истински услуги, всички едновременно, с времеви ограничения, и създава
отчет, който човек или програма могат да прочетат. Три неща, които я правят по-солидна, в реда, в
който си струва най-много да се захванеш с тях:

1. **`golangci-lint`** — събира десетки статични анализатори в едно-единствено пускане. Прекарай през него
   собствения си `revisor` и прочети какво казва: почти винаги учи на нещо, което нито `go vet`, нито тестовете хващат.
2. **[Effective Go](https://go.dev/doc/effective_go) отново.** Сега, след като написа пълна
   програма, ще го прочетеш по различен начин, отколкото в урок 0.
3. **Чети добър чужд код.** Самият пакет `net/http` от стандартната библиотека е добра отправна
   точка: вече имаш с какво да се ориентираш, за да разбереш защо е написан така, както е написан.

**И обещаното още от README:** ще напишеш същата тази програма отново, на Rust, в сестринския
курс. Не за да сравняваш синтаксис, а за да видиш същия проблем, решен с други инструменти — това
е, което наистина учи по какво се различават двата езика.

---

**Предишен:** [Урок 6 — Конкурентност](06-concurrencia.md)
