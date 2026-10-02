# Curso de Go — de cero a escribir un buen programa

**Por Dorian Chávez, fundador de Hábil y arquitecto de integración.**

> ⚠️ **ESTE CURSO ESTÁ EN REVISIÓN Y NO ESTÁ TERMINADO.** No se publica hasta que la dirección comercial dé su
> aceptación. Medido con `./medir-profundidad.sh` el 29-sep-2026: **no cumple todavía la paridad de profundidad.**

## Dónde está el contenido

**El material vive por idioma.** El original es el español:

- 📘 **[Español — el curso completo](es/README.md)** ← empieza aquí
- Traducciones: `en/`, `fr/`, `pt/`, `bg/`

## Qué hay en la raíz

| | |
|---|---|
| `es/`, `en/`, `fr/`, `pt/`, `bg/` | el curso, un directorio por idioma |
| `revisor/` | el proyecto real del curso: los mismos archivos `.go` que las lecciones 5-7 citan, con sus propias pruebas, `go vet` y `gofmt` |
| `verificar-publicable.sh` | 🔴 **se corre ANTES de publicar.** El material es público: una ruta interna o un token quedan indexados y no se retiran. **Falla cerrado** y trae autoprueba |
| `medir-profundidad.sh` | mide las líneas de explicación por lección. **Piso 150, mediana del curso 250.** Es el árbitro de la paridad: nadie cuenta a mano. Falla cerrado si un archivo tiene un número impar de vallas de código (\`\`\`), que es como una sección entera se vuelve invisible para el conteo |
| `verificar-programas.sh` | compila y corre cada programa completo de las lecciones 0-4 (los marcados `// figNN_NN.go`), y compara su salida real contra la documentada |
| `verificar-extractos.sh` | lo mismo, pero para 5-7: compara cada bloque de código marcado **extracto** contra el archivo real de `revisor/` del que se copió, byte a byte, y compila y corre cada **ejemplo** real bajo `revisor/ejemplos/` |
| `LICENSE.md` | CC BY-SA 4.0 |

## Qué es un programa, qué es un extracto, y por qué importa la diferencia

**No todo bloque de código de este curso es un programa completo, y decirlo mal fue un defecto real de
una entrega anterior.** Hay tres tipos, y cada uno se verifica distinto:

- **Programa completo** — corre solo, de principio a fin. En las lecciones 0-4 son los bloques marcados
  `// figNN_NN.go`; en las lecciones 5-7 son los que viven de verdad en `revisor/ejemplos/`. Los dos se
  compilan y se ejecutan antes de publicarse, y su salida documentada se compara contra la real.
- **Extracto** — una porción exacta, copiada tal cual, de un archivo real de `revisor/` (que sí es un
  programa completo, con sus propias pruebas). No corre por sí solo fuera de ese archivo, pero
  `verificar-extractos.sh` confirma que sigue siendo copia fiel, byte a byte, del archivo real —para que
  un cambio en el código no deje a la lección enseñando algo que ya no existe.
- **Fragmento** — una ilustración sintáctica o un extracto deliberadamente abreviado (con `// ...` en vez
  de código real), para enseñar un patrón sin el ruido de una función completa. No pretende ser copia
  exacta de nada, y no se verifica en automático — se declara así, en vez de dejar que alguien lo confunda
  con una promesa que no se cumple.

## Cómo se acepta este curso

**Nada se da por terminado sin la aceptación de la dirección comercial**, y se entrega con evidencia: la corrida
del medidor, la del verificador con su control positivo, y la salida real de los ejemplos ejecutados —los
programas completos y los extractos verificados, con los fragmentos declarados aparte, no confundidos con
ninguno de los dos.

**El molde de cada lección** —los ocho puntos, sin emojis en los encabezados— es el del README de
`hb-curso-typescript`, que es el acordado para los tres cursos.
