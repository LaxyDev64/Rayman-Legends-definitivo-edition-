package rdv

// simple_auth.go
//
// Protocolo "Simple Authentication" (ID 16) — parte de la lista "Not
// provided by NEX" de Kinnay (protocolos que sólo aparecen en juegos
// Ubisoft, no en la NEX de Nintendo). Sustituye al Ticket
// Granting/Kerberos completo que usa NEX: los juegos Ubisoft
// tradicionalmente autentican contra un servicio HTTP/REST propio
// (viejo esquema: onlineconfigservice.ubi.com, como en acb-rdv) y sólo
// usan esta capa RDV para el login "ligero" dentro de la sesión de
// juego una vez ya tienes credenciales válidas.
//
// OJO: esto es un placeholder mínimo, NO la especificación exacta —
// a diferencia de OLS Storage (documentado por Kinnay con detalle),
// Simple Authentication no tiene una wiki pública tan completa. La
// forma real de los mensajes (campos exactos, orden) hay que
// confirmarla con acb-rdv (QuazalWV, en C#) o con captura de tráfico
// real. Se deja aquí la estructura mínima razonable para que el
// prototipo compile y sirva de punto de partida.

const SimpleAuthenticationProtocolID = 16

const (
	MethodLogin        = 1 // nombre/orden exactos por confirmar
	MethodLoginEx       = 2
	MethodRequestTicket = 3
)

type LoginRequest struct {
	Username string
}

func (s *Stream) ReadLoginRequest() (LoginRequest, error) {
	var r LoginRequest
	var err error
	r.Username, err = s.ReadString()
	return r, err
}

type LoginResponse struct {
	Result       uint32 // código de resultado NEX/RDV (0 = éxito, por confirmar convención exacta)
	PID          int32
	ConnectionID uint32
	// El login real de NEX también devuelve un "pidConnectionData" con
	// IP/puerto del servidor seguro y una ticket cifrada — Simple
	// Authentication de Ubisoft probablemente tiene algo equivalente
	// pero más simple. Por confirmar.
}

func (s *Stream) WriteLoginResponse(r LoginResponse) {
	s.WriteUint32(r.Result)
	s.WriteSint32(r.PID)
	s.WriteUint32(r.ConnectionID)
}

// AuthStore es el contrato mínimo que necesita el server para resolver
// un login. La implementación real validará contra tu base de datos de
// cuentas (o contra una sesión ya validada vía el servicio HTTP de
// Ubisoft, si decides mantener ese flujo en vez de reemplazarlo).
type AuthStore interface {
	Authenticate(username string) (pid int32, ok bool, err error)
}

func DispatchAuth(store AuthStore, methodID uint32, payload []byte) ([]byte, error) {
	in := NewReadStream(payload)
	out := NewWriteStream()

	switch methodID {
	case MethodLogin:
		req, err := in.ReadLoginRequest()
		if err != nil {
			return nil, err
		}
		pid, ok, err := store.Authenticate(req.Username)
		if err != nil {
			return nil, err
		}
		resp := LoginResponse{PID: pid}
		if !ok {
			resp.Result = 1 // placeholder: "fallo" — sustituir por el código real de error RDV
		}
		out.WriteLoginResponse(resp)
		return out.Bytes(), nil

	default:
		return out.Bytes(), nil
	}
}
