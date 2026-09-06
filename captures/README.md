# Respuestas capturadas de una consola real

Estos ficheros son respuestas que una Switch de verdad recibio de los servidores de
Nintendo, y que el servidor reproduce tal cual. No se pueden deducir ni regenerar: o se
tienen o no se tienen.

| Fichero | Metodo | Bytes | Forma |
|---|---|---|---|
| `resp_0x6e_07.bin` | Utility 7, `GetIntegerSettings` | 2590 | u32 numero de entradas (431), luego pares clave u16 + valor u32 |
| `resp_0x73_60.bin` | DataStore 60 | 15704 | u32 numero de elementos (25), luego la lista |

## Por que estan aqui

El 2026-09-06 se reconstruyo el servidor desde este repositorio y **quickplay dejo de
funcionar**, mientras que las arenas seguian bien. La causa: el binario que llevaba un mes
en produccion era del 11 de agosto y este arbol tiene un solo commit, del 21 de julio. Tres
semanas de trabajo nunca subidas, y entre ellas estas capturas.

`init_replay.go` lo dice de si mismo: es un tocon que responde lista vacia a todo. Con eso,
el metodo 60 —que la consola llama unas cien veces cada quince minutos— devolvia cero
elementos, y el juego respondia *"You can't play online because your connection is poor"*.

Se recuperaron leyendo el binario de agosto: en un binario de Go una cadena embebida es un
par (puntero, longitud), asi que se busco la longitud exacta que el propio registro del
servidor imprimia (`0x73.60 -> 15704 captured bytes`) y se tradujo la direccion virtual a
desplazamiento con las cabeceras del ELF. La aritmetica de la tabla de reglajes cuadra al
byte, lo que confirma que la extraccion es correcta y no un trozo cualquiera.

## Lo que FALTA

El binario de agosto reproduce **once** metodos. Estos dos son los grandes; los otros nueve
son de 0 a 29 bytes y no se pudieron localizar por longitud sin ambiguedad:

    0x6e.9 (17)  0x6e.10 (10)  0x6e.11 (4)
    0x73.21 (8)  0x73.61 (4)   0x73.82 (0)  0x73.83 (0)  0x73.84 (29)  0x73.95 (6)

Son pequenos y probablemente triviales (una lista vacia son cuatro ceros), pero **no estan
comprobados**. Antes de volver a desplegar una reconstruccion hay que recuperarlos tambien,
o registrarlos desde el binario de agosto mientras siga en produccion.

## Regla que sale de todo esto

Reconstruir Smash desde este repositorio **rompe quickplay** mientras falten estas
respuestas. Vale hoy y valdra dentro de seis meses. El binario del 11 de agosto se conserva
en el servidor como `/opt/ssbu/ssbucs.bak-20260811` y es, por ahora, la unica copia completa
de lo que hace falta.
