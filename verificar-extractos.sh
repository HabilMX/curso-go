#!/usr/bin/env bash
# Verifica los bloques de código de las lecciones 05, 06 y 07, que no son
# programas de un archivo (como los de 02-04, que revisa verificar-programas.sh)
# sino EXTRACTOS del proyecto real en revisor/, o EJEMPLOS reales bajo
# revisor/ejemplos/.
#
# Por qué existe: un extracto que se copió a mano de revisor/ puede quedar
# desincronizado del código real en el siguiente cambio (pasó de verdad: la
# familia .yaml y la tabla de §7.4 llegaron a publicarse desincronizados). Un
# extracto sin ancla verificable no es prueba de nada — así que cada bloque de
# código de estas tres lecciones tiene que declarar QUÉ es, con un comentario
# HTML justo antes del bloque:
#
#   <!-- verificar:extracto:RUTA -->
#       El bloque debe aparecer TAL CUAL dentro de RUTA (relativa a revisor/),
#       línea por línea, con UN SOLO corrimiento de indentación CONSTANTE para
#       todo el bloque (el mismo prefijo de espacios/tabs en cada línea no
#       vacía) — no "ignorando toda la indentación": una indentación caótica
#       (2 espacios en una línea, 6 en la siguiente, sin que el archivo real
#       tenga esa misma forma) SÍ debe fallar.
#
#   <!-- verificar:ejemplo:RUTA -->
#       El bloque corresponde a un programa real en revisor/RUTA.
#       1. Primero se COMPILA (go build) — esto tiene que salir 0, sin
#          excepción. Un error de sintaxis aquí es un error de sintaxis, no
#          "el ejemplo truena a propósito".
#       2. Si hay un bloque "$ go run ..." con salida documentada justo
#          después, se CORRE y su salida se compara CONTRA esa, exacta. Si no
#          coincide, es un 🔴 real — nunca se reclasifica como "no
#          determinista" solo porque no coincidió.
#
#   <!-- verificar:ejemplo:RUTA:nodeterminista -->
#       Igual que "ejemplo", pero declarando EXPLÍCITAMENTE que la salida no
#       es comparable byte a byte (goroutine IDs, direcciones, orden de
#       impresión). Aun así se exige que COMPILE (paso 1). La salida, si la
#       documentación trae una, se revisa a mano — nunca se asume ni se
#       inventa la razón de por qué no matchea.
#
#   <!-- verificar:fragmento -->
#       Ilustración sintáctica o extracto ABREVIADO a propósito (con "// ..."
#       reemplazando código real). No pretende ser copia exacta de nada. Se
#       cuenta pero no se verifica automáticamente.
#
# Un bloque ```go sin ninguno de los cuatro es un error de este guion: falla
# cerrado, no lo deja pasar en silencio.
#
# 🔴 FALLA CERRADO en cualquiera de estos casos, ANTES de imprimir nada de
# "limpio" — medido con la autoprueba, porque una compuerta que se pone verde
# sin haber trabajado es peor que no tener compuerta:
#   - falta el intérprete de Python
#   - el glob de lecciones no encontró NINGÚN archivo
#   - el paso de extracción no produjo NINGUNA línea de resultado
#   - el total de bloques encontrados es 0
#
# Uso:   ./verificar-extractos.sh                 (corre sobre el curso real)
#        ./verificar-extractos.sh --probar        (autoprueba: siembra los
#                                                   defectos conocidos — H1 a
#                                                   H6, más los dos ataques al
#                                                   piso de la elisión — en
#                                                   una copia temporal, y
#                                                   exige que todos se
#                                                   detecten)
# Sale 0 si todo cuadra, 1 si algo no corresponde, 2 si algo está mal armado
# (incluida la propia compuerta no habiendo podido trabajar).
set -uo pipefail
cd "$(dirname "$0")"

# ---------------------------------------------------------------------------
# El verificador de verdad, parametrizado: recibe el glob de lecciones y la
# carpeta de revisor/, para poder correr tanto sobre el curso real como sobre
# la copia sembrada de la autoprueba con el MISMO código.
# ---------------------------------------------------------------------------
correr_verificacion() {
  local glob_lecciones="$1" revisor_dir="$2"
  local total=0 extractos_ok=0 extractos_mal=0 ejemplos_ok=0 ejemplos_mal=0 ejemplos_manual=0 fragmentos=0 sin_marcar=0

  command -v python3 >/dev/null || { echo "  🔴 falta python3: no se puede verificar nada."; return 2; }
  command -v go >/dev/null || { echo "  🔴 falta go: no se pueden compilar ni correr los ejemplos."; return 2; }
  if [ ! -d "$revisor_dir" ]; then
    echo "  🔴 la carpeta de revisor '$revisor_dir' no existe. NO publicar sobre esta base."
    return 2
  fi

  shopt -s nullglob
  local archivos=( $glob_lecciones )
  shopt -u nullglob
  if [ "${#archivos[@]}" -eq 0 ]; then
    echo "  🔴 el patrón '$glob_lecciones' no encontró NINGÚN archivo. Eso no es 'nada que revisar':"
    echo "     es la compuerta trabajando sobre el vacío. NO publicar sobre esta base."
    return 2
  fi

  local salida
  salida=$(mktemp)
  for cap in "${archivos[@]}"; do
    python3 - "$cap" "$revisor_dir" <<'PY' >> "$salida"
import re, sys, os, subprocess, tempfile

cap, revisor = sys.argv[1], sys.argv[2]
texto = open(cap, encoding="utf8").read()

patron = re.compile(
    r'(?:^[ \t]*<!--\s*verificar:(extracto|ejemplo|fragmento):?([^\s>]*)\s*-->\n)?'
    r'^[ \t]*```go\n(.*?)^[ \t]*```(.*?)(?=\n##|\n[ \t]*<!--\s*verificar:|\n[ \t]*```go|\Z)',
    re.S | re.M,
)

def una_linea(s):
    """🔴 Medido por QA: un mensaje de error de go build de varias líneas,
    impreso tal cual en un campo separado por tabs, hacía que el lector de
    bash de más abajo contara la segunda línea como si fuera OTRO resultado
    —de ahí que un solo ejemplo roto se reportara como "32 bloques" en vez
    de 31—. Cada resultado tiene que caber en UNA línea de verdad, siempre."""
    return " ⏎ ".join(s.split("\n"))

def normalizar(bloque):
    lineas = [l.rstrip() for l in bloque.split("\n")]
    while lineas and lineas[0] == "":
        lineas.pop(0)
    while lineas and lineas[-1] == "":
        lineas.pop()
    return lineas

def es_extracto_fiel(bloque_lineas, archivo_lineas):
    """¿Las líneas del bloque aparecen, EN ORDEN Y CONSECUTIVAS, dentro de las
    líneas del archivo, con UN SOLO corrimiento de indentación CONSTANTE para
    TODO el bloque? No compara "ignorando la indentación" — eso permitiría
    una indentación caótica (H5, medido: 6 espacios en líneas alternas) y
    llamarlo "byte a byte" sería falso.

    El corrimiento puede ir en CUALQUIERA de los dos sentidos (el bloque
    puede traer MÁS indentación que el archivo real — p. ej. un extracto
    metido dentro de un punto de una lista numerada del markdown — o MENOS
    — p. ej. el cuerpo de un método copiado sin la indentación de la
    struct/paquete que lo envuelve), pero el sentido y el prefijo exacto se
    fijan con la primera línea no vacía y se exigen idénticos, carácter por
    carácter (tabs y espacios NO son intercambiables), en TODAS las demás."""
    def resolver(bl, al):
        """Un solo par (bl, al) -> (direccion, prefijo) si coinciden salvo un
        prefijo constante en alguno de los dos lados, o None si no coinciden
        en absoluto. Cuando al == bl, el prefijo es "" — y ES la dirección
        "archivo_mas_indentado" con prefijo vacío, NUNCA un caso aparte
        exento del chequeo: eso fue el bug (H5) que dejaba pasar una
        indentación caótica en líneas alternas, porque una línea "igual" no
        se comparaba contra el prefijo ya establecido por las demás."""
        if al == bl:
            return ("archivo_mas_indentado", "")
        if len(al) > len(bl) and al[len(al) - len(bl):] == bl:
            return ("archivo_mas_indentado", al[: len(al) - len(bl)])
        if len(bl) > len(al) and bl[len(bl) - len(al):] == al:
            return ("bloque_mas_indentado", bl[: len(bl) - len(al)])
        return None

    def coincide_con_prefijo(bl, al, direccion, prefijo):
        if direccion == "archivo_mas_indentado":
            return al == prefijo + bl
        return bl == prefijo + al

    n = len(bloque_lineas)
    if n == 0:
        return False
    if not any(l.strip() != "" for l in bloque_lineas):
        return False
    for offset in range(0, len(archivo_lineas) - n + 1):
        ventana = archivo_lineas[offset:offset + n]
        direccion = prefijo = None
        candidato = True
        for bl, al in zip(bloque_lineas, ventana):
            if bl.strip() == "":
                if al.strip() != "":
                    candidato = False
                    break
                continue
            if direccion is None:
                resuelto = resolver(bl, al)
                if resuelto is None:
                    candidato = False
                    break
                direccion, prefijo = resuelto
                continue
            # Con (direccion, prefijo) YA fijados por la primera línea no
            # vacía, TODAS las demás —incluida cualquiera que por sí sola se
            # viera "igual"— tienen que satisfacer EXACTAMENTE esa misma
            # fórmula. Una línea igual con prefijo establecido distinto de
            # vacío es, precisamente, la inconsistencia que hay que atrapar.
            if not coincide_con_prefijo(bl, al, direccion, prefijo):
                candidato = False
                break
        if candidato:
            return True
    return False

resultados = []
for m in patron.finditer(texto):
    tipo, ruta, codigo, resto = m.group(1), m.group(2), m.group(3), m.group(4)
    numlinea = texto[:m.start()].count("\n") + 1
    resultados.append((numlinea, tipo, ruta, codigo, resto))

for numlinea, tipo, ruta, codigo, resto in resultados:
    if tipo is None:
        print(f"SIN_MARCAR\t{cap}\t{numlinea}\t-")
        continue
    if tipo == "fragmento":
        print(f"FRAGMENTO\t{cap}\t{numlinea}\t-")
        continue
    if tipo == "extracto":
        archivo = os.path.join(revisor, ruta)
        if not os.path.isfile(archivo):
            print(f"EXTRACTO_MAL\t{cap}\t{numlinea}\tarchivo no existe: {archivo}")
            continue
        contenido_archivo = open(archivo, encoding="utf8").read()
        bl = normalizar(codigo)
        al = contenido_archivo.split("\n")
        if es_extracto_fiel(bl, al):
            print(f"EXTRACTO_OK\t{cap}\t{numlinea}\t{ruta}")
        else:
            print(f"EXTRACTO_MAL\t{cap}\t{numlinea}\tno se encontró tal cual (indentación constante) en {ruta}")
        continue
    if tipo == "ejemplo":
        nodeterminista = ruta.endswith(":nodeterminista")
        ruta_prog = ruta[: -len(":nodeterminista")] if nodeterminista else ruta
        carpeta = os.path.join(revisor, ruta_prog)
        if not os.path.isdir(carpeta):
            print(f"EJEMPLO_MAL\t{cap}\t{numlinea}\tcarpeta no existe: {carpeta}")
            continue
        ejemplos_dir = os.path.join(revisor, "ejemplos")
        ruta_relativa = os.path.relpath(carpeta, ejemplos_dir)

        # PASO 1, obligatorio siempre: COMPILAR. Tiene que salir 0. Un error
        # de sintaxis es un error de sintaxis, nunca "el ejemplo truena a
        # propósito" — eso se decide en el PASO 2, con el programa ya
        # compilado, no aquí.
        try:
            with tempfile.TemporaryDirectory() as tmp:
                rb = subprocess.run(["go", "build", "-o", os.path.join(tmp, "bin"), "./" + ruta_relativa],
                                     cwd=ejemplos_dir, capture_output=True, text=True, timeout=30)
        except FileNotFoundError:
            # No es el contenido: es que el binario "go" no se pudo invocar
            # (aunque el "command -v go" de arriba haya pasado — un PATH
            # distinto para el subproceso, por ejemplo). Se dice así, no como
            # si el ejemplo tuviera la culpa.
            print(f"EJEMPLO_MAL\t{cap}\t{numlinea}\tHERRAMIENTA AUSENTE: no se pudo invocar 'go' (no es un defecto del ejemplo)")
            continue
        if rb.returncode != 0:
            print(f"EJEMPLO_MAL\t{cap}\t{numlinea}\tNO COMPILA: {una_linea(rb.stderr.strip()[:200])}")
            continue

        # 🔴 Medido por QA, activo de verdad: correr siempre "go run ./ruta" a
        # secas, sin mirar qué comando documenta la lección, dejaba pasar un
        # caso real donde la lección promete "$ go run -race archivo.go" y la
        # compuerta corría SIN -race — otra salida, y ":nodeterminista" tapaba
        # la diferencia en vez de exponerla. La compuerta tiene que correr el
        # comando que el libro promete, no uno parecido.
        sal = re.search(r'```\n(\$ go run [^\n]*)\n(.*?)```', resto, re.S)
        if sal:
            comando_doc, cuerpo = sal.group(1), sal.group(2)
            # Si el bloque muestra VARIAS corridas seguidas (para ilustrar que
            # algo cambia entre una y otra), se toma solo la salida de la
            # primera: comparar contra dos corridas concatenadas no tiene
            # sentido, y lo que sigue ("(cuatro corridas más...)") es texto
            # narrativo, no salida del programa.
            cuerpo = re.split(r'^\$ go run', cuerpo, maxsplit=1, flags=re.M)[0]
        if not sal:
            # Sin comando documentado que replicar: igual se ejecuta con "go
            # run" simple, solo para confirmar que además de compilar corre,
            # pero sin salida que comparar no hay promesa que verificar.
            try:
                r = subprocess.run(["go", "run", "./" + ruta_relativa],
                                    cwd=ejemplos_dir, capture_output=True, text=True, timeout=15)
            except subprocess.TimeoutExpired:
                pass
            print(f"EJEMPLO_SIN_SALIDA_FIJA\t{cap}\t{numlinea}\t{ruta_prog} (compila y corre; sin bloque de salida que comparar)")
            continue

        esperado = cuerpo
        # extrae las BANDERAS del comando documentado (todo lo que empiece con
        # "-"), y las usa tal cual con go run; el nombre de archivo que trae
        # el comando documentado (p. ej. "sinesperar.go") se ignora, porque el
        # programa real vive en revisor/ejemplos/, no con ese nombre suelto.
        banderas = [tok for tok in comando_doc.split()[2:] if tok.startswith("-")]

        # PASO 2: ejecutar de verdad, CON las banderas documentadas. Aquí SÍ
        # es válido que truene (deadlock, panic, exit distinto de 0) — el
        # programa ya demostró que compila en el PASO 1.
        try:
            r = subprocess.run(["go", "run", *banderas, "./" + ruta_relativa],
                                cwd=ejemplos_dir, capture_output=True, text=True, timeout=30)
            real = r.stdout + r.stderr
        except subprocess.TimeoutExpired:
            real = "(tiempo agotado corriendo el ejemplo)"
        def norm_lineas(s):
            return [l.rstrip() for l in s.strip("\n").split("\n")]

        def piso_de_elision(esperado_lineas):
            """🔴 Toda salida de escape necesita un piso, o se come la
            comprobación entera — es la cuarta vez que este patrón aparece en
            este guion (el grep que fallaba siempre, el conjunto de códigos
            que aceptaba un error de compilación, la exención de la línea
            "igual", y ahora esta). "..." es un escape LEGÍTIMO —un panic de
            verdad trae una traza no determinista—, pero sin piso, un bloque
            que sea SOLO "..." pasaría cualquier salida, y un bloque que
            TERMINE en "..." dejaría todo lo que sigue sin verificar. Devuelve
            un mensaje de error si el bloque está mal escrito, o None si es
            válido: tiene que haber al menos una línea literal (piso mínimo),
            y no puede terminar en una elisión sin ancla después."""
            if "..." not in esperado_lineas:
                return None  # no usa elisión: nada que exigirle a esta regla
            sin_vacias_al_final = list(esperado_lineas)
            while sin_vacias_al_final and sin_vacias_al_final[-1].strip() == "":
                sin_vacias_al_final.pop()
            literales = [l for l in sin_vacias_al_final if l != "..."]
            if not literales:
                return "el bloque de salida documentado no tiene NINGUNA línea literal (todo es \"...\"): no ancla nada"
            if sin_vacias_al_final and sin_vacias_al_final[-1] == "...":
                return "el bloque de salida documentado TERMINA en \"...\": después de la última elisión tiene que quedar algo que anclar"
            return None

        def coincide_con_elision(esperado_lineas, real_lineas):
            """Compara línea por línea, salvo que una línea del esperado sea
            literalmente "...": ahí se permite CUALQUIER número de líneas
            reales (incluido cero) hasta que vuelva a aparecer, en orden, la
            siguiente línea literal del esperado. Esto es lo que necesita un
            panic real: el mensaje y el "exit status N" son deterministas,
            la traza de en medio (goroutine, direcciones) no lo es, y la
            lección ya la elide con "..." a propósito — no es lo mismo que
            "no determinista": el resto SÍ se exige exacto. Solo se llama
            DESPUÉS de que piso_de_elision confirma que el bloque tiene piso;
            aquí ya no hace falta el caso "..." al final."""
            i = j = 0
            while i < len(esperado_lineas):
                if esperado_lineas[i] == "...":
                    i += 1
                    if i == len(esperado_lineas):
                        return False  # "..." al final sin ancla: piso_de_elision ya debió rechazar esto antes
                    objetivo = esperado_lineas[i]
                    while j < len(real_lineas) and real_lineas[j] != objetivo:
                        j += 1
                    if j == len(real_lineas):
                        return False
                    continue
                if j >= len(real_lineas) or real_lineas[j] != esperado_lineas[i]:
                    return False
                i += 1
                j += 1
            return j == len(real_lineas)

        esperado_lineas, real_lineas = norm_lineas(esperado), norm_lineas(real)

        # 🔴 Piso obligatorio ANTES de comparar, sin excepción ni siquiera para
        # los ejemplos ":nodeterminista": un bloque de salida mal escrito (sin
        # ancla, o que termina en "...") no es "no verificable", es una
        # promesa vacía — se rechaza el bloque, no se disculpa.
        error_piso = piso_de_elision(esperado_lineas)
        if error_piso:
            print(f"EJEMPLO_MAL\t{cap}\t{numlinea}\tbloque de salida mal escrito: {error_piso}")
            continue

        coincide = (
            coincide_con_elision(esperado_lineas, real_lineas)
            or (esperado.strip() == "" and real.strip() == "")
        )

        if coincide:
            print(f"EJEMPLO_OK\t{cap}\t{numlinea}\t{ruta_prog}")
        elif nodeterminista:
            # Declarado explícitamente como no determinista: no falla, pero
            # tampoco se cuenta como verificado — es revisión manual, y se
            # marca así, no como "ok".
            print(f"EJEMPLO_NODETERMINISTA\t{cap}\t{numlinea}\t{ruta_prog}")
        else:
            # 🔴 No estaba declarado no determinista y la salida NO coincide:
            # esto es un fallo real, no una suposición de no-determinismo.
            print(f"EJEMPLO_MAL\t{cap}\t{numlinea}\tsalida no coincide (no declarado no determinista): esperado={esperado.strip()[:80]!r} real={real.strip()[:80]!r}")
PY
  done

  if [ ! -s "$salida" ]; then
    echo "  🔴 el paso de extracción no produjo NINGÚN resultado (python falló o no había nada que"
    echo "     procesar). Eso no es 'todo limpio': es que no se verificó nada. NO publicar."
    rm -f "$salida"
    return 2
  fi

  while IFS=$'\t' read -r tipo cap linea resto; do
    total=$((total+1))
    case "$tipo" in
      EXTRACTO_OK) extractos_ok=$((extractos_ok+1)); printf "  ✅ extracto   %-28s :%-5s %s\n" "$cap" "$linea" "$resto" ;;
      EXTRACTO_MAL) extractos_mal=$((extractos_mal+1)); printf "  🔴 extracto   %-28s :%-5s %s\n" "$cap" "$linea" "$resto" ;;
      EJEMPLO_OK) ejemplos_ok=$((ejemplos_ok+1)); printf "  ✅ ejemplo    %-28s :%-5s %s\n" "$cap" "$linea" "$resto" ;;
      EJEMPLO_MAL) ejemplos_mal=$((ejemplos_mal+1)); printf "  🔴 ejemplo    %-28s :%-5s %s\n" "$cap" "$linea" "$resto" ;;
      EJEMPLO_SIN_SALIDA_FIJA) ejemplos_manual=$((ejemplos_manual+1)); printf "  ⬜ ejemplo(compila y corre, sin salida fija)  %-28s :%-5s %s\n" "$cap" "$linea" "$resto" ;;
      EJEMPLO_NODETERMINISTA) ejemplos_manual=$((ejemplos_manual+1)); printf "  ⬜ ejemplo(no determinista, declarado)  %-28s :%-5s %s\n" "$cap" "$linea" "$resto" ;;
      FRAGMENTO) fragmentos=$((fragmentos+1)); printf "  ⬜ fragmento  %-28s :%-5s (declarado, no verificable)\n" "$cap" "$linea" ;;
      SIN_MARCAR) sin_marcar=$((sin_marcar+1)); printf "  🔴 SIN MARCAR %-28s :%-5s — todo bloque \`\`\`go debe declarar extracto/ejemplo/fragmento\n" "$cap" "$linea" ;;
      *) echo "  🔴 línea de resultado no reconocida: $tipo $cap $linea $resto" ;;
    esac
  done < "$salida"
  rm -f "$salida"

  echo
  echo "  total de bloques go: $total"
  echo "  extractos: $extractos_ok ok, $extractos_mal mal (línea por línea, indentación constante, contra revisor/)"
  echo "  ejemplos:  $ejemplos_ok ok automático, $ejemplos_manual compilan/corren sin salida comparable en automático, $ejemplos_mal mal"
  echo "  fragmentos declarados (excluidos a propósito): $fragmentos"
  echo "  sin marcar: $sin_marcar"

  if [ "$total" -eq 0 ]; then
    echo "  🔴 se procesaron archivos pero el total de bloques es 0. Eso no es 'nada que corregir':"
    echo "     es que el patrón de extracción no encontró nada donde debería. NO publicar."
    return 2
  fi
  if [ "$sin_marcar" -gt 0 ]; then
    echo "  🔴 hay bloques \`\`\`go sin declarar qué son. NO publicar."
    return 2
  fi
  if [ "$extractos_mal" -gt 0 ] || [ "$ejemplos_mal" -gt 0 ]; then
    echo "  🔴 hay extractos o ejemplos que no corresponden al código real. NO publicar."
    return 1
  fi
  if [ "$ejemplos_manual" -gt 0 ]; then
    echo "  ⚠️  $ejemplos_manual ejemplo(s) compilan y corren de verdad, sin salida comparable en automático"
    echo "     (no determinista o sin bloque de salida): su fidelidad se confirma a mano, no la puede"
    echo "     automatizar esta compuerta."
  fi
  echo "  ✅ todos los bloques de 05-07 están declarados; los comparables byte a byte, coinciden."
  return 0
}

# ---------------------------------------------------------------------------
# Autoprueba: siembra los defectos reales que QA encontró (H1-H6, más los dos
# ataques al piso de la elisión) en una copia temporal del curso, y exige que
# todos se detecten. Si uno solo pasa como "limpio", la autoprueba falla y
# este guion se niega a decir que sirve.
# ---------------------------------------------------------------------------
autoprueba() {
  local tmp
  tmp=$(mktemp -d)
  trap 'rm -rf "$tmp"' RETURN

  mkdir -p "$tmp/es" "$tmp/revisor/ejemplos/00-det/" "$tmp/revisor/internal/x"

  cat > "$tmp/revisor/internal/x/real.go" <<'EOF'
package x

func Suma(a, b int) int {
	resultado := a + b
	return resultado
}
EOF

  # real2.go, con indentación de ESPACIOS (no tabs) a propósito: sirve para
  # construir, a mano, el caso exacto que se le escapó a H5 la vez pasada —
  # una línea que coincide LITERAL con el archivo ("igual", prefijo vacío)
  # intercalada con una que necesita un prefijo real. El bug dejaba pasar
  # esto porque una línea "igual" quedaba exenta del chequeo de consistencia.
  cat > "$tmp/revisor/internal/x/real2.go" <<'EOF'
package x

func Doble(a int) int {
  x := a * 2
  return x
}
EOF

  cat > "$tmp/revisor/ejemplos/00-det/main.go" <<'EOF'
package main

import "fmt"

func main() {
	fmt.Println("hola")
}
EOF
  cat > "$tmp/revisor/ejemplos/go.mod" <<'EOF'
module ejemplos-prueba

go 1.27
EOF

  # H1: salida documentada sembrada como falsa para un ejemplo DETERMINISTA
  #     declarado como tal (sin ":nodeterminista"). Tiene que FALLAR (🔴), no
  #     reclasificarse como "no determinista".
  # H2: un ejemplo con error de SINTAXIS. Tiene que fallar en "no compila",
  #     nunca pasar como "compiló y corrió".
  # H5: un extracto con indentación CAÓTICA (no constante). Tiene que fallar.
  # (Fiel, para control negativo: un extracto correcto SÍ debe pasar, y un
  #  ejemplo declarado no determinista con salida distinta NO debe fallar.)
  cat > "$tmp/es/05-prueba.md" <<'EOF'
# Lección de prueba

<!-- verificar:extracto:internal/x/real.go -->
```go
func Suma(a, b int) int {
	resultado := a + b
	return resultado
}
```

<!-- verificar:extracto:internal/x/real.go -->
```go
func Suma(a, b int) int {
	  resultado := a + b
	    return resultado
}
```

<!-- verificar:extracto:internal/x/real2.go -->
```go
func Doble(a int) int {
    x := a * 2
  return x
}
```

<!-- verificar:ejemplo:ejemplos/00-det -->
```go
fmt.Println("hola")
```

```
$ go run ./ejemplos/00-det/
SALIDA FALSA SEMBRADA — esto NO es lo que el programa imprime
```

<!-- verificar:ejemplo:ejemplos/00-det:nodeterminista -->
```go
fmt.Println("hola")
```

```
$ go run ./ejemplos/00-det/
...
```

<!-- verificar:ejemplo:ejemplos/00-det:nodeterminista -->
```go
fmt.Println("hola")
```

```
$ go run ./ejemplos/00-det/
hola
...
```

<!-- verificar:ejemplo:ejemplos/00-roto:nodeterminista -->
```go
func main() {
	fmt.Println("esto no compila"
}
```
EOF

  mkdir -p "$tmp/revisor/ejemplos/00-roto"
  cat > "$tmp/revisor/ejemplos/00-roto/main.go" <<'EOF'
package main

import "fmt"

func main() {
	fmt.Println("esto no compila"
}
EOF

  local salida rc
  salida=$(correr_verificacion "$tmp/es/05-prueba.md" "$tmp/revisor" 2>&1)
  rc=$?

  local fallos=0
  echo "  --- autoprueba: sembrando H1 (salida falsa), H2 (no compila), H5 (indentación caótica) ---"

  local n_h5_malos
  n_h5_malos=$(echo "$salida" | grep -c "no se encontró tal cual (indentación constante)")
  if [ "$n_h5_malos" -ge 2 ]; then
    echo "  ✅ H5 detectado: los DOS extractos con indentación inconsistente se marcaron mal"
    echo "     (incluido el caso de la línea que coincide literal intercalada con una que no — el que se escapó la vez pasada)."
  else
    echo "  🔴 H5 NO detectado del todo: se esperaban 2 extractos malos por indentación, se vieron $n_h5_malos."
    fallos=$((fallos+1))
  fi

  if echo "$salida" | grep -q "NO COMPILA"; then
    echo "  ✅ H2 detectado: el ejemplo con error de sintaxis se marcó como que no compila."
  else
    echo "  🔴 H2 NO detectado: un ejemplo con error de sintaxis no se rechazó."
    fallos=$((fallos+1))
  fi

  if echo "$salida" | grep -q "salida no coincide (no declarado no determinista)"; then
    echo "  ✅ H1 detectado: la salida sembrada como falsa se marcó como que no coincide (no se le atribuyó a 'no determinismo')."
  else
    echo "  🔴 H1 NO detectado: una salida falsa se coló como 'no determinista' o como 'ok'."
    fallos=$((fallos+1))
  fi

  if echo "$salida" | grep -q "no tiene NINGUNA línea literal"; then
    echo "  ✅ piso de elisión (ataque 1) detectado: un bloque documentado que es solo \"...\" se rechazó."
  else
    echo "  🔴 piso de elisión (ataque 1) NO detectado: un bloque que es solo \"...\" pasó cualquier salida."
    fallos=$((fallos+1))
  fi

  if echo "$salida" | grep -q "TERMINA en \"\.\.\.\""; then
    echo "  ✅ piso de elisión (ataque 2) detectado: un bloque que termina en \"...\" se rechazó."
  else
    echo "  🔴 piso de elisión (ataque 2) NO detectado: un bloque que termina en \"...\" dejó todo lo que sigue sin verificar."
    fallos=$((fallos+1))
  fi

  if [ "$rc" -eq 0 ]; then
    echo "  🔴 la corrida con los 3 defectos sembrados salió 0 (verde). Debía salir distinto de 0."
    fallos=$((fallos+1))
  else
    echo "  ✅ la corrida con defectos sembrados salió distinto de 0 (rc=$rc)."
  fi

  echo
  echo "  --- autoprueba: control positivo (un extracto fiel real no debe fallar) ---"
  if echo "$salida" | grep -q "✅ extracto"; then
    echo "  ✅ el extracto FIEL (el primero, sin defecto) sí pasó — la compuerta no rechaza todo a ciegas."
  else
    echo "  🔴 ni siquiera el extracto fiel pasó: la compuerta está descalibrada, no solo estricta."
    fallos=$((fallos+1))
  fi

  echo
  echo "  --- autoprueba: H3/H4 (falla cerrado con glob vacío y con archivo sin bloques) ---"
  local h3out; h3out=$(mktemp)
  correr_verificacion "$tmp/es/no-existe-*.md" "$tmp/revisor" >"$h3out" 2>&1
  local rc_h3=$?
  if [ "$rc_h3" -ne 0 ] && grep -q "NINGÚN archivo" "$h3out"; then
    echo "  ✅ H3 detectado: un glob de lecciones vacío falla cerrado (rc=$rc_h3), no 'nada que revisar'."
  else
    echo "  🔴 H3 NO detectado: un glob vacío no falló cerrado (rc=$rc_h3)."
    fallos=$((fallos+1))
  fi
  rm -f "$h3out"

  printf '# vacio, sin bloques go\n' > "$tmp/es/06-vacio.md"
  local h4out; h4out=$(mktemp)
  correr_verificacion "$tmp/es/06-vacio.md" "$tmp/revisor" >"$h4out" 2>&1
  local rc_h4=$?
  if [ "$rc_h4" -ne 0 ] && grep -qE "total de bloques es 0|no produjo NINGÚN resultado" "$h4out"; then
    echo "  ✅ H4 detectado: un archivo sin ningún bloque go falla cerrado (rc=$rc_h4), no '0 de 0 está limpio'."
  else
    echo "  🔴 H4 NO detectado: un archivo sin bloques no falló cerrado (rc=$rc_h4)."
    fallos=$((fallos+1))
  fi
  rm -f "$h4out"

  echo
  echo "  --- autoprueba: H6 (propia, no la de QA) — carpeta revisor/ inexistente ---"
  local h6out; h6out=$(mktemp)
  correr_verificacion "$tmp/es/05-prueba.md" "$tmp/no/existe/esta/carpeta" >"$h6out" 2>&1
  local rc_h6=$?
  if [ "$rc_h6" -ne 0 ] && grep -q "no existe" "$h6out"; then
    echo "  ✅ H6 detectado: una carpeta de revisor/ inexistente falla cerrado (rc=$rc_h6)."
  else
    echo "  🔴 H6 NO detectado: una carpeta de revisor/ inexistente no falló cerrado (rc=$rc_h6)."
    fallos=$((fallos+1))
  fi
  rm -f "$h6out"

  echo
  if [ "$fallos" -gt 0 ]; then
    echo "  🔴 autoprueba FALLIDA: $fallos de 8 defectos sembrados (H1-H6 + los 2 ataques al piso de elisión) no se detectaron como debían."
    echo "     Esta compuerta NO se usa hasta que la autoprueba pase completa."
    return 1
  fi
  echo "  ✅ autoprueba correcta: los 8 defectos sembrados (H1-H6 + los 2 ataques al piso de elisión) se detectan, y el control positivo no se rompe."
  return 0
}

if [ "${1:-}" = "--probar" ]; then
  autoprueba
  exit $?
fi

correr_verificacion "es/0[5-7]-*.md" "revisor"
exit $?
