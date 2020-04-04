package melrose

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// https://en.wikipedia.org/wiki/Chord_(music)
type Chord struct {
	start     Note
	inversion int
	interval  int
	quality   int
}

func zeroChord() Chord {
	return Chord{start: N("C"), inversion: Ground, quality: Major, interval: Triad}
}

func (c Chord) Storex() string {
	return fmt.Sprintf("chord('%v')", c.start)
}

func (c Chord) Modified(modifiers ...int) Chord {
	modified := c
	for _, each := range modifiers {
		switch each {
		case Major:
			modified.quality = Major
		case Minor:
			modified.quality = Minor
		case Ground:
			modified.inversion = Ground
		case Inversion1:
			modified.inversion = Inversion1
		case Inversion2:
			modified.inversion = Inversion2
		}
	}
	return modified
}

func (c Chord) S() Sequence {
	notes := []Note{c.start}
	var semitones []int
	if Major == c.quality {
		semitones = []int{4, 7}
	} else if Minor == c.quality {
		semitones = []int{3, 7}
	}
	for _, each := range semitones {
		next := c.start.Pitched(each)
		notes = append(notes, next)
	}
	return Sequence{[][]Note{notes}}
}

var chordRegexp = regexp.MustCompile("([MmDA]?)([67]?)")

//  C:D7:2 = C dominant 7, 2nd inversion
func ParseChord(s string) (Chord, error) {
	if len(s) == 0 {
		return Chord{}, errors.New("illegal chord: missing note")
	}
	parts := strings.Split(s, ":")
	start, err := ParseNote(parts[0])
	if err != nil {
		return Chord{}, err
	}
	if len(parts) == 1 {
		z := zeroChord()
		z.start = start
		return z, nil
	}
	// parts > 1
	chord := Chord{start: start, quality: Major}
	chord.inversion = readInversion(parts[1])

	matches := chordRegexp.FindStringSubmatch(parts[1])
	if matches == nil {
		return Chord{}, fmt.Errorf("illegal chord: [%s]", s)
	}
	switch matches[1] {
	case "M":
		chord.quality = Major
	case "m":
		chord.quality = Minor
	case "D":
		chord.quality = Dominant
	case "A":
		chord.quality = Augmented
	}
	switch matches[2] {
	case "6":
		chord.interval = Sixth
	case "7":
		chord.interval = Seventh
	default:
		chord.interval = Triad
	}

	// parts > 2
	if len(parts) > 2 {
		chord.inversion = readInversion(parts[2])
	}
	return chord, nil
}

func readInversion(s string) int {
	switch s {
	case "1":
		return Inversion1
	case "2":
		return Inversion2
	case "3":
		return Inversion3
	default:
		return Ground
	}
}
