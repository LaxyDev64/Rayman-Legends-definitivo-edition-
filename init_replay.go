package main

// SSBU online-init handlers for DataStore (0x73) and Utility (0x6E) — stub.
//
// SSBU calls a set of DataStore and Utility getters during the online bring-up ("OnlineCore"
// population) and soft-locks if they return NotImplemented. A full server-side implementation
// of these protocols is not yet part of this tree, so each method returns a minimal valid
// response — an empty list (count = 0) — enough for the game to proceed. The live log shows
// which methods it calls so they can be implemented properly over time.

import (
	"fmt"

	nex "github.com/NextendoNetwork/nextendo-nex"
)

// replayHandler answers every method of the given protocol (0x73 / 0x6E) with a minimal
// valid response. 0x73.8 returns DataStore::NotFound (0x80690004), which SSBU handles
// gracefully; everything else returns an empty list (count = 0).
func replayHandler(proto uint16) nex.RMCHandler {
	return func(conn *nex.Connection, req *nex.RMCMessage) *nex.RMCMessage {
		s := conn.Settings

		if proto == 0x73 && req.Method == 8 {
			fmt.Printf("[SSBU DataStore] 0x73.8 -> NotFound 0x80690004\n")
			return nex.NewRMCError(s, proto, req.CallID, 0x80690004)
		}

		out := nex.NewStreamOut(s)
		out.U32(0) // empty list / count = 0
		fmt.Printf("[SSBU Init] 0x%02x.%d callID=%d -> empty list (stub)\n", proto, req.Method, req.CallID)
		return nex.NewRMCSuccess(s, proto, req.Method, req.CallID, out.Bytes())
	}
}

// setupSSBUInitReplay registers the DataStore (0x73) + Utility (0x6E) stub handlers.
func setupSSBUInitReplay(endpoint *nex.Endpoint) {
	endpoint.Register(0x73, replayHandler(0x73))
	endpoint.Register(0x6E, replayHandler(0x6E))
	fmt.Println("[SSBU Init] DataStore(0x73) + Utility(0x6E) stub registered")
}
