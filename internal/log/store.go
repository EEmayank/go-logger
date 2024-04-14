package log

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

var (
	enc = binary.BigEndian
)

const (
	// lenWidth defines the number of bytes used to store the record's lenght, which in this case is 8-bytes
	lenWidth = 8
)

// BaseStore is the interface that any log store needs to implement
type BaseStore interface {
	Append(p []byte) (n uint64, pos uint64, err error)
}

// store is a simple wrapper around a file with two APIs to append and read butyes to and from the file
type store struct {
	*os.File
	mu   sync.Mutex
	buf  *bufio.Writer
	size uint64
}

// newStore creates a store for the given file.
// Calls os.State(name string) to get the file's current size in case we're re-creationg the store
// from a file that has existing data, which would happen if, for example our service had restarted
func newStore(f *os.File) (*store, error) {
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	size := uint64(fi.Size())
	return &store{
		File: f,
		buf:  bufio.NewWriter(f),
		size: size,
	}, nil
}

// Append persits the give bytes to the stroe. We write the record so that, wehn we read the recoed, we know how many bytes to read
// We write to the buffered writer instead of directly to the file to reduce the number of system calls and improve performance. If a user wrote a lot of small recoed, this would help a lot.
// We then return the number of bytes written, which similar to GO APIs conventionally do, and the position where the store hold the record in its file.
// The segment will use this position when it creates an assosciated index entry for the record
func (s *store) Append(p []byte) (n uint64, pos uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	pos = s.size
	if err := binary.Write(s.buf, enc, uint64(len(p))); err != nil {
		return 0, 0, err
	}
	w, err := s.buf.Write(p)
	if err != nil {
		return 0, 0, err
	}
	w += lenWidth
	s.size += uint64(w)
	return uint64(w), pos, nil
}
