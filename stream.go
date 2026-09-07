package rdv

// stream.go
//
// Lector/escritor binario mínimo para los tipos comunes que usan tanto NEX
// como Quazal Rendez-Vous (RMC comparte formato: son la misma familia de
// protocolo, NEX es un derivado de Quazal RDV).
//
// Esto es standalone para que el prototipo compile solo. En tu fork real
// de nextendo-nex ya tienes stream.go / rmc.go / types.go con lectura y
// escritura equivalentes (probablemente más completos) — sustituye este
// archivo por imports a tu propio core en cuanto integres esto.

import (
	"encoding/binary"
	"errors"
	"math"
)

var ErrShortBuffer = errors.New("rdv: buffer too short")

// Stream envuelve un []byte para lectura/escritura secuencial estilo RMC.
type Stream struct {
	buf []byte
	pos int
}

func NewReadStream(b []byte) *Stream {
	return &Stream{buf: b}
}

func NewWriteStream() *Stream {
	return &Stream{buf: make([]byte, 0, 256)}
}

func (s *Stream) Bytes() []byte {
	return s.buf
}

// ---- lectura ----

func (s *Stream) need(n int) error {
	if s.pos+n > len(s.buf) {
		return ErrShortBuffer
	}
	return nil
}

func (s *Stream) ReadUint8() (uint8, error) {
	if err := s.need(1); err != nil {
		return 0, err
	}
	v := s.buf[s.pos]
	s.pos++
	return v, nil
}

func (s *Stream) ReadSint8() (int8, error) {
	v, err := s.ReadUint8()
	return int8(v), err
}

func (s *Stream) ReadUint16() (uint16, error) {
	if err := s.need(2); err != nil {
		return 0, err
	}
	v := binary.LittleEndian.Uint16(s.buf[s.pos:])
	s.pos += 2
	return v, nil
}

func (s *Stream) ReadUint32() (uint32, error) {
	if err := s.need(4); err != nil {
		return 0, err
	}
	v := binary.LittleEndian.Uint32(s.buf[s.pos:])
	s.pos += 4
	return v, nil
}

func (s *Stream) ReadSint32() (int32, error) {
	v, err := s.ReadUint32()
	return int32(v), err
}

func (s *Stream) ReadUint64() (uint64, error) {
	if err := s.need(8); err != nil {
		return 0, err
	}
	v := binary.LittleEndian.Uint64(s.buf[s.pos:])
	s.pos += 8
	return v, nil
}

func (s *Stream) ReadFloat32() (float32, error) {
	v, err := s.ReadUint32()
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(v), nil
}

func (s *Stream) ReadBool() (bool, error) {
	v, err := s.ReadUint8()
	return v != 0, err
}

// String NEX/RDV: uint16 length (incluye null terminator) + bytes + NUL.
func (s *Stream) ReadString() (string, error) {
	l, err := s.ReadUint16()
	if err != nil {
		return "", err
	}
	if l == 0 {
		return "", nil
	}
	if err := s.need(int(l)); err != nil {
		return "", err
	}
	b := s.buf[s.pos : s.pos+int(l)-1] // sin el NUL final
	s.pos += int(l)
	return string(b), nil
}

// DateTime NEX/RDV: empaquetado en un uint64 (bitfield); se deja como
// entero crudo por ahora — implementar el desempaquetado si se necesita
// leer año/mes/día/hora en el server.
type DateTime uint64

func (s *Stream) ReadDateTime() (DateTime, error) {
	v, err := s.ReadUint64()
	return DateTime(v), err
}

// ---- escritura ----

func (s *Stream) WriteUint8(v uint8) {
	s.buf = append(s.buf, v)
}

func (s *Stream) WriteSint8(v int8) {
	s.WriteUint8(uint8(v))
}

func (s *Stream) WriteUint16(v uint16) {
	var b [2]byte
	binary.LittleEndian.PutUint16(b[:], v)
	s.buf = append(s.buf, b[:]...)
}

func (s *Stream) WriteUint32(v uint32) {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], v)
	s.buf = append(s.buf, b[:]...)
}

func (s *Stream) WriteSint32(v int32) {
	s.WriteUint32(uint32(v))
}

func (s *Stream) WriteUint64(v uint64) {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], v)
	s.buf = append(s.buf, b[:]...)
}

func (s *Stream) WriteFloat32(v float32) {
	s.WriteUint32(math.Float32bits(v))
}

func (s *Stream) WriteBool(v bool) {
	if v {
		s.WriteUint8(1)
	} else {
		s.WriteUint8(0)
	}
}

func (s *Stream) WriteString(v string) {
	if v == "" {
		s.WriteUint16(0)
		return
	}
	b := append([]byte(v), 0) // NUL terminator
	s.WriteUint16(uint16(len(b)))
	s.buf = append(s.buf, b...)
}

func (s *Stream) WriteDateTime(v DateTime) {
	s.WriteUint64(uint64(v))
}

// ---- listas genéricas ----
//
// NEX/RDV codifica List<T> como uint32 count + T repetido. Como Go no
// tiene genéricos "libres" para leer arbitrariamente sin reflection
// pesada, se exponen helpers por tipo primitivo y se dejan las listas de
// structs para que cada protocolo las lea con su propio bucle (ver
// ols_storage.go).

func (s *Stream) ReadListUint32() ([]uint32, error) {
	n, err := s.ReadUint32()
	if err != nil {
		return nil, err
	}
	out := make([]uint32, 0, n)
	for i := uint32(0); i < n; i++ {
		v, err := s.ReadUint32()
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, nil
}

func (s *Stream) WriteListUint32(vs []uint32) {
	s.WriteUint32(uint32(len(vs)))
	for _, v := range vs {
		s.WriteUint32(v)
	}
}

func (s *Stream) ReadListSint32() ([]int32, error) {
	n, err := s.ReadUint32()
	if err != nil {
		return nil, err
	}
	out := make([]int32, 0, n)
	for i := uint32(0); i < n; i++ {
		v, err := s.ReadSint32()
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, nil
}

func (s *Stream) WriteListSint32(vs []int32) {
	s.WriteUint32(uint32(len(vs)))
	for _, v := range vs {
		s.WriteSint32(v)
	}
}

func (s *Stream) ReadListFloat32() ([]float32, error) {
	n, err := s.ReadUint32()
	if err != nil {
		return nil, err
	}
	out := make([]float32, 0, n)
	for i := uint32(0); i < n; i++ {
		v, err := s.ReadFloat32()
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, nil
}

func (s *Stream) WriteListFloat32(vs []float32) {
	s.WriteUint32(uint32(len(vs)))
	for _, v := range vs {
		s.WriteFloat32(v)
	}
}

func (s *Stream) ReadListString() ([]string, error) {
	n, err := s.ReadUint32()
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, n)
	for i := uint32(0); i < n; i++ {
		v, err := s.ReadString()
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, nil
}

func (s *Stream) WriteListString(vs []string) {
	s.WriteUint32(uint32(len(vs)))
	for _, v := range vs {
		s.WriteString(v)
	}
}

