package melrose

import (
	"time"

	"github.com/emicklei/melrose/notify"
)

type Transformer interface {
	Transform(Sequence) Sequence
}

type Sequenceable interface {
	S() Sequence
}

type NoteConvertable interface {
	ToNote() Note
}

type Storable interface {
	Storex() string
}

type Indexable interface {
	At(i int) Sequenceable
}

type AudioDevice interface {
	// Per device specific commands
	Command(args []string) notify.Message

	// Play schedules all the notes on the timeline using a BPM (beats-per-minute).
	// Returns the end time of the last played Note.
	Play(seq Sequenceable, bpm float64, beginAt time.Time) (endingAt time.Time)
	Record(deviceID int, stopAfterInactivity time.Duration) (*Recording, error)
	Timeline() *Timeline
	SetEchoNotes(echo bool)
	Reset()
	Close()
}

type LoopController interface {
	Start()
	Stop()
	Reset()

	SetBPM(bpm float64)
	BPM() float64

	SetBIAB(biab int)
	BIAB() int

	Begin(l *Loop)
	End(l *Loop)
}

type MapFunc func(seq Sequenceable) Sequenceable

// TODO experiment
type Mappeable interface {
	Map(m MapFunc) Mappeable
}
