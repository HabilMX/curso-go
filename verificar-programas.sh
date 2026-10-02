#!/usr/bin/env bash
# Compila y ejecuta CADA programa de ejemplo del curso, y compara su salida
# real contra la documentada.
#
# Por qué existe: el curso promete que «cada programa se compila y se ejecuta
# antes de publicarse». Una promesa verificable que nadie verifica es peor que
# no prometer nada. Un QA independiente encontró tres salidas escritas a mano
# —una de ellas una tabla que el texto llamaba «derecha» estando torcida— y
# este script existe para que eso no vuelva a pasar.
#
# Uso:  ./verificar-programas.sh [idioma]      (por omisión: es)
#       ./verificar-programas.sh es --mostrar  (imprime la salida real de cada uno)
#
# Sale 0 si todo cuadra, 1 si alguna salida no coincide.
set -uo pipefail
cd "$(dirname "$0")"
IDI="${1:-es}"
MOSTRAR="${2:-}"
TRABAJO=$(mktemp -d)
trap 'rm -rf "$TRABAJO"' EXIT

command -v go >/dev/null || { echo "falta Go"; exit 2; }
echo "Go: $(go version)"
echo

total=0; ok=0; malos=0; sin_salida=0

# extrae los programas: bloques ```go cuyo codigo trae // figNN_NN.go,
# quitando la numeracion de linea con que se imprimen en el curso.
for cap in "$IDI"/[0-9][0-9]-*.md; do
  python3 - "$cap" "$TRABAJO" <<'PY'
import re, sys, os
cap, trabajo = sys.argv[1], sys.argv[2]
texto = open(cap, encoding="utf8").read()
# cada bloque go, con lo que venga despues (para encontrar su salida documentada)
for m in re.finditer(r'```go\n(.*?)```(.*?)(?=\n##|\n```go|\Z)', texto, re.S):
    codigo, resto = m.group(1), m.group(2)
    # quita la numeracion « 12  » con que el curso imprime los programas.
    # Una linea en blanco del codigo aparece como « 4 » (numero sin los dos
    # espacios), asi que el patron tiene que aceptar las dos formas.
    lineas = codigo.split("\n")
    NUM = re.compile(r'^\s{0,4}\d{1,3}(\s\s|\s*$)')
    conNum = [l for l in lineas if l.strip()]
    if conNum and sum(1 for l in conNum if NUM.match(l)) >= len(conNum) * 0.9:
        lineas = [NUM.sub('', l, count=1) for l in lineas]
    codigo = "\n".join(lineas)
    nom = re.search(r'//\s*(fig\d+_\d+)\.go', codigo)
    if not nom:
        continue
    fig = nom.group(1)
    d = os.path.join(trabajo, fig)
    os.makedirs(d, exist_ok=True)
    open(os.path.join(d, fig + ".go"), "w", encoding="utf8").write(codigo)
    open(os.path.join(d, ".capitulo"), "w").write(os.path.basename(cap))
    # la salida documentada: primer bloque bash que empiece con $ go run
    sal = re.search(r'```bash\n\$ go run [^\n]*\n(.*?)```', resto, re.S)
    if sal:
        esperado = sal.group(1)
        # Si el curso muestra VARIAS ejecuciones del mismo programa (para
        # ilustrar salidas no deterministas, como el orden de un map), se toma
        # solo la primera: comparar contra tres corridas distintas no tiene
        # sentido. El programa se marca como no determinista.
        if re.search(r'^\$ go run', esperado, re.M):
            esperado = re.split(r'^\$ go run', esperado, maxsplit=1, flags=re.M)[0]
            open(os.path.join(d, ".nodeterminista"), "w").write("1")
        open(os.path.join(d, ".esperado"), "w", encoding="utf8").write(esperado)
    # Archivos de datos que el programa necesita: el curso los muestra en un
    # bloque de texto ANTES del codigo, precedido de su nombre entre comillas.
    prev = texto[:m.start()]
    for dat in re.finditer(r'`([a-z0-9_.-]+\.(?:txt|yaml|yml|json))`[^\n]*:\n\n```\n(.*?)```', prev, re.S):
        open(os.path.join(d, dat.group(1)), "w", encoding="utf8").write(dat.group(2))
PY
done

for d in "$TRABAJO"/fig*; do
  [ -d "$d" ] || continue
  fig=$(basename "$d"); cap=$(cat "$d/.capitulo" 2>/dev/null)
  total=$((total+1))
  ( cd "$d" && go mod init "$fig" >/dev/null 2>&1 )
  # el archivo se corre por su nombre real (fig.go), no como "main.go": si el
  # programa truena o no compila, el mensaje real trae ese nombre, y tiene que
  # poder compararse contra el que el curso documenta con ESE MISMO nombre.
  real=$( cd "$d" && go run "$fig.go" 2>&1 )
  if [ -n "$MOSTRAR" ]; then
    echo "──── $fig ($cap) ────"; echo "$real"; echo
    continue
  fi
  if [ ! -f "$d/.esperado" ]; then
    printf "  %-12s %-26s ⬜ sin salida documentada\n" "$fig" "$cap"; sin_salida=$((sin_salida+1)); continue
  fi
  # compara ignorando lineas con placeholders (0x..., rutas) y espacios al final
  # 🔴 ORDEN QUE IMPORTA, medido por QA: la regla de SIGSEGV tiene que ir ANTES
  # que cualquier regla que pudiera borrar o alterar esa línea de forma
  # distinta en el esperado y en lo real. La versión anterior borraba la línea
  # completa SOLO cuando contenía el placeholder literal "0x..." — que el
  # curso escribe en el texto documentado para el `pc=`, pero que una corrida
  # real JAMÁS produce (el valor real es hexadecimal, nunca tres puntos). El
  # resultado: la línea desaparecía del esperado y sobrevivía en lo real,
  # dejando ese caso rojo para siempre, sin importar si el programa estaba
  # bien. La corrección normaliza la línea COMPLETA a una forma fija en los
  # dos lados — código, dirección y contador de programa varían por
  # plataforma y no son parte de lo que este curso promete verificar.
  norm() { sed \
      -e 's/[[:space:]]*$//' \
      -e 's/^\[signal SIGSEGV.*/[signal SIGSEGV: ...]/' \
      -e '/0x\.\.\./d' \
      -e 's|^[[:space:]]*/.*\.go:\([0-9]*\).*|\tTRAZA:\1|' \
      -e '/^goroutine/d' -e '/^main\.main/d' -e '/^exit status/d'; }
  if [ -f "$d/.nodeterminista" ]; then
    # su salida cambia entre ejecuciones a proposito (p.ej. el orden de un map):
    # se exige que compile y que sus lineas sean las mismas SIN importar el orden
    if [ "$(echo "$real" | tr ' ' '\n' | sort | norm)" = "$(cat "$d/.esperado" | tr ' ' '\n' | sort | norm)" ]; then
      printf "  %-12s %-26s ✅ (no determinista: comparado sin orden)\n" "$fig" "$cap"; ok=$((ok+1))
    else
      printf "  %-12s %-26s 🔴 NO COINCIDE (ni sin orden)\n" "$fig" "$cap"
      diff <(cat "$d/.esperado" | tr ' ' '\n' | sort | norm) <(echo "$real" | tr ' ' '\n' | sort | norm) | head -8 | sed 's/^/        /'
      malos=$((malos+1))
    fi
  elif [ "$(echo "$real" | norm)" = "$(cat "$d/.esperado" | norm)" ]; then
    printf "  %-12s %-26s ✅\n" "$fig" "$cap"; ok=$((ok+1))
  else
    printf "  %-12s %-26s 🔴 NO COINCIDE\n" "$fig" "$cap"
    diff <(cat "$d/.esperado" | norm) <(echo "$real" | norm) | head -12 | sed 's/^/        /'
    malos=$((malos+1))
  fi
done

[ -n "$MOSTRAR" ] && exit 0
echo
echo "  programas: $total · coinciden: $ok · NO coinciden: $malos · sin salida documentada: $sin_salida"
[ "$malos" -gt 0 ] && { echo "  🔴 Hay salidas documentadas que no corresponden. NO publicar."; exit 1; }
# 🔴 FALLA CERRADO, sembrado y medido por QA: un programa sin salida documentada
# NO es lo mismo que un programa verificado. Antes, si a una figura se le
# quitaba (por accidente o por vandalismo) su bloque de salida, el conteo de
# "sin_salida" subía pero el veredicto final seguía diciendo que TODO
# correspondía, con código de salida 0 — el mismo defecto de forma que ya se
# corrigió en la comparación de figuras, pero en la línea de cierre. Una
# figura sin salida documentada es una promesa rota ("cada programa se
# compila y se ejecuta antes de publicarse") tanto como una que no coincide.
[ "$sin_salida" -gt 0 ] && { echo "  🔴 $sin_salida figura(s) sin salida documentada. NO publicar."; exit 1; }
echo "  ✅ Todas las salidas documentadas corresponden a la ejecución real."
