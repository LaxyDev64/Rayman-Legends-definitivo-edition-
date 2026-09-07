package main

// Respuestas de arranque en linea de SSBU para DataStore (0x73) y Utility (0x6E).
//
// SSBU pide una tanda de getters al entrar en linea y se queda colgado si alguno devuelve
// NotImplemented. Este fichero fue durante meses un TOCON que contestaba lista vacia a
// todo, y eso bastaba para que el juego arrancara... hasta que dejo de bastar.
//
// EL DIA QUE SE NOTO. El 2026-09-06 se reconstruyo el servidor desde este repositorio y
// quickplay dejo de funcionar, mientras las arenas seguian bien. La causa no estaba en el
// codigo sino en lo que faltaba: el binario en produccion llevaba respuestas CAPTURADAS de
// una consola real y este arbol no. El metodo 0x73.60 —que la consola llama unas cien veces
// cada quince minutos— devolvia cero elementos en vez de una lista de 25, y el juego
// mostraba "You can't play online because your connection is poor". Las arenas no usan ese
// metodo; por eso funcionaban y despistaban.
//
// POR QUE FALTABAN. No las omitio nadie: el .gitignore llevaba `*.bin` y git las descartaba
// en silencio. Ahora hay una excepcion explicita para captures/.
//
// Se recuperaron leyendo el binario de agosto, y viven en captures/ con su procedencia
// escrita. Van EMBEBIDAS en el binario a proposito: asi no dependen de que un fichero exista
// en el servidor, que es exactamente la clase de dependencia que produjo este incidente.

import (
	"embed"
	"fmt"
	"os"

	nex "github.com/NextendoNetwork/nextendo-nex"
)

//go:embed captures/*.bin
var capturas embed.FS

// claveMetodo identifica un metodo concreto de un protocolo concreto.
type claveMetodo struct {
	proto  uint16
	metodo uint32
}

// cuerpoVacio marca los metodos que respondian CERO bytes, no una lista vacia.
//
// La diferencia importa y no es teorica: una lista vacia son cuatro bytes a cero, y cero
// bytes son ninguno. Es la misma familia de fallo que la cadena vacia que costo dos noches
// en Mario Maker — un envoltorio de mas que el cliente no espera.
var cuerpoVacio = map[claveMetodo]bool{
	{0x73, 82}: true,
	{0x73, 83}: true,
}

// ficheroDe asocia cada metodo con su captura.
//
// ⚠️ LAS DOS DE CUATRO BYTES ESTAN SIN CONFIRMAR. La tabla del binario de agosto contiene
// dos respuestas de cuatro bytes, 00000000 y 64000000, que van a Utility 11 y DataStore 61
// sin que se sepa cual a cual: sus longitudes son iguales, asi que no se pudieron separar
// leyendo el binario. La asignacion de abajo es una suposicion.
//
// Se intercambian sin recompilar, poniendo SSBU_REPLAY_SWAP4=1. Si quickplay falla de una
// forma nueva, esto es lo primero que hay que probar.
func ficheroDe(proto uint16, metodo uint32) string {
	swap := os.Getenv("SSBU_REPLAY_SWAP4") == "1"

	switch (claveMetodo{proto, metodo}) {
	case claveMetodo{0x6E, 7}:
		return "captures/resp_0x6e_07.bin"
	case claveMetodo{0x6E, 9}:
		return "captures/resp_0x6e_09.bin"
	case claveMetodo{0x6E, 10}:
		return "captures/resp_0x6e_10.bin"
	case claveMetodo{0x73, 21}:
		return "captures/resp_0x73_21.bin"
	case claveMetodo{0x73, 60}:
		return "captures/resp_0x73_60.bin"
	case claveMetodo{0x73, 84}:
		return "captures/resp_0x73_84.bin"
	case claveMetodo{0x73, 95}:
		return "captures/resp_0x73_95.bin"

	case claveMetodo{0x6E, 11}:
		if swap {
			return "captures/resp_4bytes_B_64000000.bin"
		}

		return "captures/resp_4bytes_A_00000000.bin"
	case claveMetodo{0x73, 61}:
		if swap {
			return "captures/resp_4bytes_A_00000000.bin"
		}

		return "captures/resp_4bytes_B_64000000.bin"
	}

	return ""
}

// replayHandler contesta cada metodo con su captura si la hay, y si no con el repli de
// siempre: NotFound para 0x73.8, lista vacia para el resto.
//
// El repli se conserva a proposito. La consola llama a mas metodos de los que tenemos
// capturados, y una lista vacia los deja pasar; devolver NotImplemented cuelga el juego.
// Pero ahora el registro DISTINGUE los dos casos, porque confundirlos es justo lo que hizo
// que este fallo tardara una tarde en verse.
func replayHandler(proto uint16) nex.RMCHandler {
	return func(conn *nex.Connection, req *nex.RMCMessage) *nex.RMCMessage {
		s := conn.Settings

		if proto == 0x73 && req.Method == 8 {
			fmt.Printf("[SSBU DataStore] 0x73.8 -> NotFound 0x80690004\n")

			return nex.NewRMCError(s, proto, req.CallID, 0x80690004)
		}

		if cuerpoVacio[claveMetodo{proto, req.Method}] {
			fmt.Printf("[SSBU Replay] 0x%02x.%d -> cuerpo vacio (0 octetos, capturado)\n", proto, req.Method)

			return nex.NewRMCSuccess(s, proto, req.Method, req.CallID, nil)
		}

		if ruta := ficheroDe(proto, req.Method); ruta != "" {
			if b, err := capturas.ReadFile(ruta); err == nil {
				fmt.Printf("[SSBU Replay] 0x%02x.%d -> %d octetos capturados\n", proto, req.Method, len(b))

				return nex.NewRMCSuccess(s, proto, req.Method, req.CallID, b)
			}
			// Nunca deberia pasar: van embebidas. Si pasa, se dice y se sigue con el repli.
			fmt.Printf("[SSBU Replay] 0x%02x.%d -> captura %q ILEGIBLE, se usa el repli\n", proto, req.Method, ruta)
		}

		out := nex.NewStreamOut(s)
		out.U32(0)
		fmt.Printf("[SSBU Init] 0x%02x.%d callID=%d -> lista vacia (sin captura)\n", proto, req.Method, req.CallID)

		return nex.NewRMCSuccess(s, proto, req.Method, req.CallID, out.Bytes())
	}
}

// setupSSBUInitReplay registra los manejadores de DataStore (0x73) y Utility (0x6E).
func setupSSBUInitReplay(endpoint *nex.Endpoint) {
	endpoint.Register(0x73, replayHandler(0x73))
	endpoint.Register(0x6E, replayHandler(0x6E))

	n := 0
	for _, p := range []struct {
		proto  uint16
		metodo uint32
	}{{0x6E, 7}, {0x6E, 9}, {0x6E, 10}, {0x6E, 11}, {0x73, 21}, {0x73, 60}, {0x73, 61}, {0x73, 84}, {0x73, 95}} {
		if b, err := capturas.ReadFile(ficheroDe(p.proto, p.metodo)); err == nil && len(b) > 0 {
			n++
		}
	}
	fmt.Printf("[SSBU Init] DataStore(0x73) + Utility(0x6E) registrados, %d captura(s) embebida(s)\n", n)
}
