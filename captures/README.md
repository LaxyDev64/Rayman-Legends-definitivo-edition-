# Respuestas capturadas de una consola real

Respuestas que una Switch recibio de los servidores de Nintendo y que el servidor
reproduce tal cual. No se deducen ni se regeneran: o se tienen o no se tienen.

## El mapa

| Fichero | Metodo | Bytes | Certeza |
|---|---|---|---|
| `resp_0x6e_07.bin` | Utility 7, `GetIntegerSettings` | 2590 | segura |
| `resp_0x6e_09.bin` | Utility 9 | 17 | segura |
| `resp_0x6e_10.bin` | Utility 10 | 10 | segura |
| `resp_0x73_21.bin` | DataStore 21 | 8 | segura |
| `resp_0x73_60.bin` | DataStore 60 | 15704 | segura |
| `resp_0x73_84.bin` | DataStore 84 | 29 | segura |
| `resp_0x73_95.bin` | DataStore 95 | 6 | segura |
| `resp_4bytes_A_00000000.bin` | Utility 11 **o** DataStore 61 | 4 | ⚠️ AMBIGUA |
| `resp_4bytes_B_64000000.bin` | la otra de las dos | 4 | ⚠️ AMBIGUA |

`0x73.82` y `0x73.83` responden CERO bytes: no hay nada que guardar, pero hay que
contestarlos igualmente con cuerpo vacio.

Las siete primeras son seguras porque su longitud es unica en la tabla. Las dos de cuatro
bytes valen `00000000` y `64000000` (o sea 0 y 100) y no se pudo decidir cual va a cual:
probar las dos combinaciones en consola cuesta menos que seguir leyendo el binario.

## Por que estan aqui

El 2026-09-06 se reconstruyo el servidor desde este repositorio y **quickplay dejo de
funcionar** mientras las arenas seguian bien. La causa no era el codigo: este arbol tiene un
commit del 21 de julio y el binario en produccion era del 11 de agosto. Tres semanas nunca
subidas, y entre ellas estas capturas.

`init_replay.go` lo dice de si mismo: es un tocon que responde lista vacia a todo. Con eso
el metodo 0x73.60 —que la consola llama unas cien veces cada quince minutos— devolvia cero
elementos y el juego mostraba *"You can't play online because your connection is poor"*.
Las arenas no usan ese metodo; por eso funcionaban y despistaban.

## Como se recuperaron

En un binario de Go una cadena embebida es un par (puntero, longitud) en la seccion de
datos. Se busco la longitud que el propio registro del servidor imprimia
(`0x73.60 -> 15704 captured bytes`), y desde ese par se leyo la tabla entera, traduciendo
direcciones virtuales a desplazamientos con las cabeceras del ELF.

Dos trampas por el camino, anotadas para quien repita esto:

Los pares estaban desplazados ocho bytes respecto a una rejilla de dieciseis, asi que leerlos
alineados daba basura convincente. Y la mayoria de coincidencias por longitud eran nombres de
cabecera HTTP de la biblioteca estandar (`Host`, `Cookie`, `Etag`...), no capturas: hay que
mirar el contenido, no solo el tamano.

La extraccion se comprobo por aritmetica y no por fe: la tabla de reglajes declara 431
entradas y mide `4 + 431*6 = 2590` bytes, exactamente su tamano.

## Regla

Reconstruir Smash desde este repositorio rompe quickplay mientras estos ficheros no esten
CABLEADOS en el codigo. Estan guardados, que era lo urgente, pero `init_replay.go` sigue
siendo un tocon y no los lee todavia.

El binario del 11 de agosto se conserva en el servidor como `/opt/ssbu/ssbucs.bak-20260811`
y sigue siendo la referencia contra la que comprobar cualquier reconstruccion.
