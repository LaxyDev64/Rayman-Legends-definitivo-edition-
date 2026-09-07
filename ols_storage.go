package rdv

// ols_storage.go
//
// Protocolo específico de juego "OLS Storage" (ID 200) usado por Rayman
// Legends. Spec tomado de la documentación pública de Kinnay's
// NintendoClients wiki (la versión documentada corresponde a Wii U/NEX;
// se asume igual en Switch a nivel de lógica de juego, ya que Ubisoft
// comparte esa capa entre plataformas — PENDIENTE de confirmar byte a
// byte con una captura de tráfico real de la versión Switch).
//
// IDs de método:
//   1  LoadVersion
//   2  SaveLocale
//   3  SaveProfile
//   4  LoadIDCard
//   5  QueryFriendProfiles
//   6  QueryUbisoftProfiles
//   7  CreateMessage
//   8  QueryMessage
//   9  QueryLeaderboard
//   10 QuerySmartSelection
//   11 SaveScore
//   12 SaveGhost
//   13 QueryCompetitionsInfos
//   14 QueryCompetitionsHistory
//   15 QueryCompetitionOfTheDay
//   16 QueryCompetition

const OLSStorageProtocolID = 200

const (
	MethodLoadVersion              = 1
	MethodSaveLocale               = 2
	MethodSaveProfile              = 3
	MethodLoadIDCard                = 4
	MethodQueryFriendProfiles      = 5
	MethodQueryUbisoftProfiles     = 6
	MethodCreateMessage            = 7
	MethodQueryMessage             = 8
	MethodQueryLeaderboard         = 9
	MethodQuerySmartSelection      = 10
	MethodSaveScore                = 11
	MethodSaveGhost                = 12
	MethodQueryCompetitionsInfos   = 13
	MethodQueryCompetitionsHistory = 14
	MethodQueryCompetitionOfTheDay = 15
	MethodQueryCompetition         = 16
)

// ---------------------------------------------------------------------
// Structs (tipos compuestos del protocolo)
// ---------------------------------------------------------------------

type OLSProfile struct {
	PID         int32
	PlatformID  string
	Name        string
	Costume     uint32
	Country     uint32
	Level       int8
}

func (s *Stream) ReadOLSProfile() (OLSProfile, error) {
	var p OLSProfile
	var err error
	if p.PID, err = s.ReadSint32(); err != nil {
		return p, err
	}
	if p.PlatformID, err = s.ReadString(); err != nil {
		return p, err
	}
	if p.Name, err = s.ReadString(); err != nil {
		return p, err
	}
	if p.Costume, err = s.ReadUint32(); err != nil {
		return p, err
	}
	if p.Country, err = s.ReadUint32(); err != nil {
		return p, err
	}
	if p.Level, err = s.ReadSint8(); err != nil {
		return p, err
	}
	return p, nil
}

func (s *Stream) WriteOLSProfile(p OLSProfile) {
	s.WriteSint32(p.PID)
	s.WriteString(p.PlatformID)
	s.WriteString(p.Name)
	s.WriteUint32(p.Costume)
	s.WriteUint32(p.Country)
	s.WriteSint8(p.Level)
}

func (s *Stream) WriteListOLSProfile(ps []OLSProfile) {
	s.WriteUint32(uint32(len(ps)))
	for _, p := range ps {
		s.WriteOLSProfile(p)
	}
}

type OLSRichProfile struct {
	PID                   int32
	Name                  string
	PlatformID            string
	Country               int16
	StatusIcon            uint32
	LastCostume           uint32
	TotalChallengePlayed  uint16
	DailyPlayed           bool
	WeeklyPlayed          bool
	DailyExpertPlayed     bool
	WeeklyExpertPlayed    bool
	DiamondMedals         uint16
	GoldMedals            uint16
	SilverMedals          uint16
	BronzeMedals          uint16
	GlobalMedalsRank      uint32
	GlobalMedalsMaxRank   uint32
	DistanceRun           float32
	RankDistanceRun       uint32
	Lums                  float32
	RankLums              uint32
	Pets                  float32
	RankPets              uint32
	Teensies              float32
	RankTeensies          uint32
	Jumps                 float32
	RankJumps             uint32
	Costumes              float32
	RankCostumes          uint32
	StatDaily             float32
	RankDaily             uint32
	UnitDaily             int8
	StatWeekly            float32
	RankWeekly            uint32
	UnitWeekly            int8
	StatDailyExpert       float32
	RankDailyExpert       uint32
	UnitDailyExpert       int8
	StatWeeklyExpert      float32
	RankWeeklyExpert      uint32
	UnitWeeklyExpert      int8
}

func (s *Stream) WriteOLSRichProfile(p OLSRichProfile) {
	s.WriteSint32(p.PID)
	s.WriteString(p.Name)
	s.WriteString(p.PlatformID)
	s.WriteUint16(uint16(p.Country)) // Sint16 -> escrito como 2 bytes crudos
	s.WriteUint32(p.StatusIcon)
	s.WriteUint32(p.LastCostume)
	s.WriteUint16(p.TotalChallengePlayed)
	s.WriteBool(p.DailyPlayed)
	s.WriteBool(p.WeeklyPlayed)
	s.WriteBool(p.DailyExpertPlayed)
	s.WriteBool(p.WeeklyExpertPlayed)
	s.WriteUint16(p.DiamondMedals)
	s.WriteUint16(p.GoldMedals)
	s.WriteUint16(p.SilverMedals)
	s.WriteUint16(p.BronzeMedals)
	s.WriteUint32(p.GlobalMedalsRank)
	s.WriteUint32(p.GlobalMedalsMaxRank)
	s.WriteFloat32(p.DistanceRun)
	s.WriteUint32(p.RankDistanceRun)
	s.WriteFloat32(p.Lums)
	s.WriteUint32(p.RankLums)
	s.WriteFloat32(p.Pets)
	s.WriteUint32(p.RankPets)
	s.WriteFloat32(p.Teensies)
	s.WriteUint32(p.RankTeensies)
	s.WriteFloat32(p.Jumps)
	s.WriteUint32(p.RankJumps)
	s.WriteFloat32(p.Costumes)
	s.WriteUint32(p.RankCostumes)
	s.WriteFloat32(p.StatDaily)
	s.WriteUint32(p.RankDaily)
	s.WriteSint8(p.UnitDaily)
	s.WriteFloat32(p.StatWeekly)
	s.WriteUint32(p.RankWeekly)
	s.WriteSint8(p.UnitWeekly)
	s.WriteFloat32(p.StatDailyExpert)
	s.WriteUint32(p.RankDailyExpert)
	s.WriteSint8(p.UnitDailyExpert)
	s.WriteFloat32(p.StatWeeklyExpert)
	s.WriteUint32(p.RankWeeklyExpert)
	s.WriteSint8(p.UnitWeeklyExpert)
}

type OLSAttribute struct {
	AttributeType  int8
	AttributeValue uint32
}

func (s *Stream) ReadOLSAttribute() (OLSAttribute, error) {
	var a OLSAttribute
	var err error
	if a.AttributeType, err = s.ReadSint8(); err != nil {
		return a, err
	}
	if a.AttributeValue, err = s.ReadUint32(); err != nil {
		return a, err
	}
	return a, nil
}

func (s *Stream) ReadListOLSAttribute() ([]OLSAttribute, error) {
	n, err := s.ReadUint32()
	if err != nil {
		return nil, err
	}
	out := make([]OLSAttribute, 0, n)
	for i := uint32(0); i < n; i++ {
		a, err := s.ReadOLSAttribute()
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

type OLSMessage struct {
	MessageType       int8
	MessagePrompt     bool
	MessageDRC        bool
	MessageBloomberg  bool
	MessageDate       DateTime
	MessageDuration   uint32
	MessageTitle      string
	MessageBody       string
	MessageButtons    []string
	MessageAttributes []OLSAttribute
}

func (s *Stream) WriteOLSMessage(m OLSMessage) {
	s.WriteSint8(m.MessageType)
	s.WriteBool(m.MessagePrompt)
	s.WriteBool(m.MessageDRC)
	s.WriteBool(m.MessageBloomberg)
	s.WriteDateTime(m.MessageDate)
	s.WriteUint32(m.MessageDuration)
	s.WriteString(m.MessageTitle)
	s.WriteString(m.MessageBody)
	s.WriteListString(m.MessageButtons)
	s.WriteUint32(uint32(len(m.MessageAttributes)))
	for _, a := range m.MessageAttributes {
		s.WriteSint8(a.AttributeType)
		s.WriteUint32(a.AttributeValue)
	}
}

type OLSCompetitionResult struct {
	IDLeaderboard uint32
	Name          string
	Begin         DateTime
	End           DateTime
	Level         uint32
	Mode          uint32
	Rank          uint32
	MaxRank       uint32
}

func (s *Stream) WriteOLSCompetitionResult(r OLSCompetitionResult) {
	s.WriteUint32(r.IDLeaderboard)
	s.WriteString(r.Name)
	s.WriteDateTime(r.Begin)
	s.WriteDateTime(r.End)
	s.WriteUint32(r.Level)
	s.WriteUint32(r.Mode)
	s.WriteUint32(r.Rank)
	s.WriteUint32(r.MaxRank)
}

type OLSCompetition struct {
	Result            OLSCompetitionResult
	Message           string
	Seed              uint32
	Objective         float32
	ScoreValidation   float32
	IDBricks          uint32
	Score             float32
}

func (s *Stream) WriteOLSCompetition(c OLSCompetition) {
	s.WriteOLSCompetitionResult(c.Result)
	s.WriteString(c.Message)
	s.WriteUint32(c.Seed)
	s.WriteFloat32(c.Objective)
	s.WriteFloat32(c.ScoreValidation)
	s.WriteUint32(c.IDBricks)
	s.WriteFloat32(c.Score)
}

type OLSSelectionRow struct {
	ID        uint32
	Name      string
	IDGhost   uint64
	IDCostume uint32
	Country   uint32
	Level     uint32
	Score     float32
}

func (s *Stream) WriteOLSSelectionRow(r OLSSelectionRow) {
	s.WriteUint32(r.ID)
	s.WriteString(r.Name)
	s.WriteUint64(r.IDGhost)
	s.WriteUint32(r.IDCostume)
	s.WriteUint32(r.Country)
	s.WriteUint32(r.Level)
	s.WriteFloat32(r.Score)
}

func (s *Stream) WriteListOLSSelectionRow(rs []OLSSelectionRow) {
	s.WriteUint32(uint32(len(rs)))
	for _, r := range rs {
		s.WriteOLSSelectionRow(r)
	}
}

type OLSCompetitionInfos struct {
	IDCompetition     uint32
	Participants      uint32
	Friends           []uint32
	LevelID           uint32
	Mode              uint32
	MyRank            uint32
	RemainingSeconds  uint32
	Competitors       []OLSSelectionRow
	Unit              uint32
}

func (s *Stream) WriteOLSCompetitionInfos(c OLSCompetitionInfos) {
	s.WriteUint32(c.IDCompetition)
	s.WriteUint32(c.Participants)
	s.WriteListUint32(c.Friends)
	s.WriteUint32(c.LevelID)
	s.WriteUint32(c.Mode)
	s.WriteUint32(c.MyRank)
	s.WriteUint32(c.RemainingSeconds)
	s.WriteListOLSSelectionRow(c.Competitors)
	s.WriteUint32(c.Unit)
}

type OLSLdbRow struct {
	ID         uint32
	Name       string
	Value      float32
	Costume    uint32
	StatusIcon uint32
	Country    uint32
}

func (s *Stream) WriteOLSLdbRow(r OLSLdbRow) {
	s.WriteUint32(r.ID)
	s.WriteString(r.Name)
	s.WriteFloat32(r.Value)
	s.WriteUint32(r.Costume)
	s.WriteUint32(r.StatusIcon)
	s.WriteUint32(r.Country)
}

func (s *Stream) WriteListOLSLdbRow(rs []OLSLdbRow) {
	s.WriteUint32(uint32(len(rs)))
	for _, r := range rs {
		s.WriteOLSLdbRow(r)
	}
}

type OLSTomb struct {
	PID       int32
	Name      string
	IDCostume uint32
	X, Y, Z   float32
}

func (s *Stream) WriteOLSTomb(t OLSTomb) {
	s.WriteSint32(t.PID)
	s.WriteString(t.Name)
	s.WriteUint32(t.IDCostume)
	s.WriteFloat32(t.X)
	s.WriteFloat32(t.Y)
	s.WriteFloat32(t.Z)
}

func (s *Stream) WriteListOLSTomb(ts []OLSTomb) {
	s.WriteUint32(uint32(len(ts)))
	for _, t := range ts {
		s.WriteOLSTomb(t)
	}
}

type OLSFriend struct {
	PID          int32
	Relationship uint32
}

func (s *Stream) ReadOLSFriend() (OLSFriend, error) {
	var f OLSFriend
	var err error
	if f.PID, err = s.ReadSint32(); err != nil {
		return f, err
	}
	if f.Relationship, err = s.ReadUint32(); err != nil {
		return f, err
	}
	return f, nil
}

func (s *Stream) ReadListOLSFriend() ([]OLSFriend, error) {
	n, err := s.ReadUint32()
	if err != nil {
		return nil, err
	}
	out := make([]OLSFriend, 0, n)
	for i := uint32(0); i < n; i++ {
		f, err := s.ReadOLSFriend()
		if err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, nil
}

// ---------------------------------------------------------------------
// Handler
// ---------------------------------------------------------------------
//
// Store es la interfaz que el server real implementa (contra tu base de
// datos). Aquí sólo se define el contrato + un handler de stubs que
// responde con valores vacíos/por defecto — igual que hiciste en el
// server de SSBU para DataStore/Utility, para que el juego avance en el
// bring-up online mientras se implementa la lógica real.

type Store interface {
	LoadVersion() (version int32, sandboxName string, err error)
	SaveLocale(pid int32, localeCode string) error
	SaveProfile(pid int32, req SaveProfileRequest) (SaveProfileResponse, error)
	LoadIDCard(pid int32, target int32) (OLSRichProfile, error)
	QueryLeaderboard(pid int32, idLeaderboard uint32) (QueryLeaderboardResponse, error)
	SaveScore(pid int32, req SaveScoreRequest) (SaveScoreResponse, error)
	QueryCompetitionOfTheDay(pid int32, idCompetitionMeta uint32) (QueryCompetitionOfTheDayResponse, error)
	// ... completar según se vaya implementando cada método real.
}

type SaveProfileRequest struct {
	UpdateBitfield  uint32
	Level           int8
	Currency        int32
	Costume         uint32
	BronzeMedals    uint16
	SilverMedals    uint16
	GoldMedals      uint16
	DiamondMedals   uint16
	RunDistance     uint32
	TeensiesFreed   uint16
	Jumps           uint32
	UnlockedPets    uint16
	Pets            uint64
	UnlockedCostumes uint16
}

func (s *Stream) ReadSaveProfileRequest() (SaveProfileRequest, error) {
	var r SaveProfileRequest
	var err error
	if r.UpdateBitfield, err = s.ReadUint32(); err != nil {
		return r, err
	}
	if r.Level, err = s.ReadSint8(); err != nil {
		return r, err
	}
	if csi, err2 := s.ReadUint32(); err2 != nil {
		return r, err2
	} else {
		r.Currency = int32(csi)
	}
	if r.Costume, err = s.ReadUint32(); err != nil {
		return r, err
	}
	if r.BronzeMedals, err = s.ReadUint16(); err != nil {
		return r, err
	}
	if r.SilverMedals, err = s.ReadUint16(); err != nil {
		return r, err
	}
	if r.GoldMedals, err = s.ReadUint16(); err != nil {
		return r, err
	}
	if r.DiamondMedals, err = s.ReadUint16(); err != nil {
		return r, err
	}
	if r.RunDistance, err = s.ReadUint32(); err != nil {
		return r, err
	}
	if r.TeensiesFreed, err = s.ReadUint16(); err != nil {
		return r, err
	}
	if r.Jumps, err = s.ReadUint32(); err != nil {
		return r, err
	}
	if r.UnlockedPets, err = s.ReadUint16(); err != nil {
		return r, err
	}
	if r.Pets, err = s.ReadUint64(); err != nil {
		return r, err
	}
	if r.UnlockedCostumes, err = s.ReadUint16(); err != nil {
		return r, err
	}
	return r, nil
}

type SaveProfileResponse struct {
	CompetitionMedals [4]uint16
}

func (s *Stream) WriteSaveProfileResponse(r SaveProfileResponse) {
	s.WriteUint16(r.CompetitionMedals[0])
	s.WriteUint16(r.CompetitionMedals[1])
	s.WriteUint16(r.CompetitionMedals[2])
	s.WriteUint16(r.CompetitionMedals[3])
}

type QueryLeaderboardResponse struct {
	Result       []OLSLdbRow
	Graduations  []float32
	Envelope     []uint32
	Unit         uint32
	MyCountry    uint32
	Participants uint32
	Cacheable    bool
}

func (s *Stream) WriteQueryLeaderboardResponse(r QueryLeaderboardResponse) {
	s.WriteListOLSLdbRow(r.Result)
	s.WriteListFloat32(r.Graduations)
	s.WriteListUint32(r.Envelope)
	s.WriteUint32(r.Unit)
	s.WriteUint32(r.MyCountry)
	s.WriteUint32(r.Participants)
	s.WriteBool(r.Cacheable)
}

type SaveScoreRequest struct {
	IDLeaderboard       uint32
	IsObjectiveReached  bool
	Score               float32
	TombX, TombY, TombZ float32
	IDCostume           uint32
}

func (s *Stream) ReadSaveScoreRequest() (SaveScoreRequest, error) {
	var r SaveScoreRequest
	var err error
	if r.IDLeaderboard, err = s.ReadUint32(); err != nil {
		return r, err
	}
	if r.IsObjectiveReached, err = s.ReadBool(); err != nil {
		return r, err
	}
	if r.Score, err = s.ReadFloat32(); err != nil {
		return r, err
	}
	if r.TombX, err = s.ReadFloat32(); err != nil {
		return r, err
	}
	if r.TombY, err = s.ReadFloat32(); err != nil {
		return r, err
	}
	if r.TombZ, err = s.ReadFloat32(); err != nil {
		return r, err
	}
	if r.IDCostume, err = s.ReadUint32(); err != nil {
		return r, err
	}
	return r, nil
}

type SaveScoreResponse struct {
	SaveScore      bool
	Retry          bool
	Medal          uint32
	MessageMedal   string
	MessageFriends string
}

func (s *Stream) WriteSaveScoreResponse(r SaveScoreResponse) {
	s.WriteBool(r.SaveScore)
	s.WriteBool(r.Retry)
	s.WriteUint32(r.Medal)
	s.WriteString(r.MessageMedal)
	s.WriteString(r.MessageFriends)
}

type QueryCompetitionOfTheDayResponse struct {
	Competition      OLSCompetition
	RemainingSeconds uint32
	Tomorrow         OLSCompetition
}

func (s *Stream) WriteQueryCompetitionOfTheDayResponse(r QueryCompetitionOfTheDayResponse) {
	s.WriteOLSCompetition(r.Competition)
	s.WriteUint32(r.RemainingSeconds)
	s.WriteOLSCompetition(r.Tomorrow)
}

// Dispatch enruta un método entrante (methodID) al Store real. pid es el
// PID del jugador ya autenticado (viene de la capa de sesión/secure
// connection, no de este protocolo). Devuelve el payload de respuesta ya
// serializado, listo para envolver en un paquete RMC de respuesta.
//
// De momento sólo están cableados los métodos con más probabilidad de
// aparecer en el bring-up inicial (LoadVersion, SaveProfile,
// QueryLeaderboard, SaveScore, QueryCompetitionOfTheDay). El resto
// (QueryFriendProfiles, QueryUbisoftProfiles, mensajería, ghosts,
// smart selection, historial de competiciones...) quedan como
// siguiente paso, una vez confirmados con captura de tráfico real.
func Dispatch(store Store, pid int32, methodID uint32, payload []byte) ([]byte, error) {
	in := NewReadStream(payload)
	out := NewWriteStream()

	switch methodID {
	case MethodLoadVersion:
		version, sandbox, err := store.LoadVersion()
		if err != nil {
			return nil, err
		}
		out.WriteSint32(version)
		out.WriteString(sandbox)
		return out.Bytes(), nil

	case MethodSaveLocale:
		locale, err := in.ReadString()
		if err != nil {
			return nil, err
		}
		if err := store.SaveLocale(pid, locale); err != nil {
			return nil, err
		}
		return out.Bytes(), nil // sin respuesta

	case MethodSaveProfile:
		req, err := in.ReadSaveProfileRequest()
		if err != nil {
			return nil, err
		}
		resp, err := store.SaveProfile(pid, req)
		if err != nil {
			return nil, err
		}
		out.WriteSaveProfileResponse(resp)
		return out.Bytes(), nil

	case MethodLoadIDCard:
		targetPID, err := in.ReadSint32()
		if err != nil {
			return nil, err
		}
		profile, err := store.LoadIDCard(pid, targetPID)
		if err != nil {
			return nil, err
		}
		out.WriteOLSRichProfile(profile)
		return out.Bytes(), nil

	case MethodQueryLeaderboard:
		idLdb, err := in.ReadUint32()
		if err != nil {
			return nil, err
		}
		resp, err := store.QueryLeaderboard(pid, idLdb)
		if err != nil {
			return nil, err
		}
		out.WriteQueryLeaderboardResponse(resp)
		return out.Bytes(), nil

	case MethodSaveScore:
		req, err := in.ReadSaveScoreRequest()
		if err != nil {
			return nil, err
		}
		resp, err := store.SaveScore(pid, req)
		if err != nil {
			return nil, err
		}
		out.WriteSaveScoreResponse(resp)
		return out.Bytes(), nil

	case MethodQueryCompetitionOfTheDay:
		idMeta, err := in.ReadUint32()
		if err != nil {
			return nil, err
		}
		resp, err := store.QueryCompetitionOfTheDay(pid, idMeta)
		if err != nil {
			return nil, err
		}
		out.WriteQueryCompetitionOfTheDayResponse(resp)
		return out.Bytes(), nil

	default:
		// Métodos aún no cableados: responder vacío/OK para no bloquear
		// el bring-up, igual que el stub de DataStore/Utility del server
		// de SSBU. Sustituir caso por caso a medida que se confirme el
		// comportamiento real con captura de tráfico.
		return out.Bytes(), nil
	}
}
