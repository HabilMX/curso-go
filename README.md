# Curso de Go — de cero a escribir un buen programa

[![Verificar programas](https://github.com/HabilMX/curso-go/actions/workflows/verificar.yml/badge.svg)](https://github.com/HabilMX/curso-go/actions/workflows/verificar.yml)

**Por Dorian Chávez, fundador de Hábil y arquitecto de integración.**

**Cada programa de este curso se compila y se ejecuta automáticamente en cada cambio; el sello verde lo comprueba y cualquiera puede ver la corrida.**
Haz clic en el sello para abrir la última corrida y ver, paso por paso, qué se ejecutó y qué salió.

> Por ahora el curso está en español; las traducciones vienen en camino.

## Dónde está el contenido

- 📘 **[Español — el curso completo](es/README.md)** ← empieza aquí
- 💻 **[`programas/`](programas/)** — todos los programas del curso, listos para ejecutar

## Cómo ver que los programas funcionan

Sin instalar nada: abre el sello de arriba. Cada corrida muestra los pasos que se ejecutaron y su resultado.

En tu computadora, con [Go](https://go.dev/dl/) instalado (la versión exacta está en `programas/revisor/go.mod`):

```bash
cd programas/03-errores-interfaces
go run fig03_06.go          # compara lo que imprime con fig03_06.salida.txt
```

Cada programa de las lecciones 2 a 4 es un archivo `figNN_NN.go` dentro de la carpeta de su lección, con su
salida esperada al lado (`figNN_NN.salida.txt`). Los que fallan a propósito, porque la lección enseña justo ese error (no compilan o terminan en un `panic`), traen `figNN_NN.error-esperado.txt` en lugar de la salida. Se ejecutan desde la carpeta de su lección: los que leen un archivo de
datos (como `servicios.txt` en la lección 4) lo traen ahí mismo. El proyecto completo que se construye en las
lecciones 5 a 7 está en [`programas/revisor/`](programas/revisor/), con sus pruebas (`go test ./...`).

**Las lecciones son la fuente; `programas/` es una copia que se genera de ellas.** Los programas de las lecciones 2 a 4
se extraen de los bloques de código de `es/*.md` con `herramientas/generar-programas.sh`, y en cada cambio la
verificación comprueba que la copia es idéntica a lo que dicen las lecciones. Así lo que lees y lo que ejecutas
no pueden diferir.

## Qué hay en el repositorio

| | |
|---|---|
| `es/` | el curso en español, una lección por archivo |
| `en/`, `fr/`, `pt/`, `bg/` | **futuras:** las traducciones todavía no existen |
| `programas/` | los programas de las lecciones 2 a 4 (uno por archivo, con su salida esperada) y `revisor/`, el proyecto real de las lecciones 5 a 7 |
| `herramientas/` | los scripts que verifican el curso (ver abajo) |
| `.github/workflows/verificar.yml` | la verificación automática que muestra el sello |
| `verificar-publicable.sh` | revisa que el material no contenga rutas internas ni claves antes de publicarlo |
| `LICENSE.md` | CC BY-SA 4.0 |

Dentro de `herramientas/`:

| | |
|---|---|
| `verificar-programas.sh` | compila y corre cada programa completo de las lecciones 2 a 4 (los marcados `// figNN_NN.go`) y compara su salida real contra la documentada |
| `verificar-extractos.sh` | lo mismo para las lecciones 5 a 7: compara cada bloque marcado **extracto** contra el archivo real de `programas/revisor/`, byte a byte, y compila y corre cada **ejemplo** real |
| `generar-programas.sh` | arma `programas/` desde las lecciones; con `--comprobar` verifica que esté al día |
| `medir-profundidad.sh` | mide las líneas de explicación por lección |
| `verificar-traducciones.sh`, `registrar-traduccion.sh`, `registro-traducciones.tsv` | llevan el control de qué traducciones están al día, para cuando existan |

## Qué es un programa, qué es un extracto, y por qué importa la diferencia

**No todo bloque de código de este curso es un programa completo, y confundirlos engaña al lector.** Hay tres tipos, y cada uno se verifica distinto:

- **Programa completo** — corre solo, de principio a fin. En las lecciones 0-4 son los bloques marcados
  `// figNN_NN.go`; en las lecciones 5-7 son los que viven de verdad en `programas/revisor/ejemplos/`. Los dos se
  compilan y se ejecutan en cada cambio, y su salida documentada se compara contra la real.
- **Extracto** — una porción exacta, copiada tal cual, de un archivo real de `programas/revisor/` (que sí es un
  programa completo, con sus propias pruebas). No corre por sí solo fuera de ese archivo, pero
  `herramientas/verificar-extractos.sh` confirma que sigue siendo copia fiel, byte a byte, del archivo real —para que
  un cambio en el código no deje a la lección enseñando algo que ya no existe.
- **Fragmento** — una ilustración sintáctica o un extracto deliberadamente abreviado (con `// ...` en vez
  de código real), para enseñar un patrón sin el ruido de una función completa. No pretende ser copia
  exacta de nada, y no se verifica en automático — se declara así, en vez de dejar que alguien lo confunda
  con una promesa que no se cumple.
