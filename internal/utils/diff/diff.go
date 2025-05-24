package diff

import (
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

//go:generate stringer -type=Operation -trimprefix=Operation
type Operation int

const (
	// OperationDelete item represents a delete diff.
	OperationDelete Operation = -1
	// OperationInsert item represents an insert diff.
	OperationInsert Operation = 1
	// OperationEqual item represents an equal diff.
	OperationEqual Operation = 0
)

type Diff struct {
	Operation Operation
	Value     []string
}

func Slice(a, b []string) []Diff {
	dmp := diffmatchpatch.New()

	text1 := strings.Join(a, "\n")
	text2 := strings.Join(b, "\n")

	// Map lines to runes and diff them
	chars1, chars2, lineArray := dmp.DiffLinesToRunes(text1, text2)
	diffs := dmp.DiffMainRunes(chars1, chars2, false)
	diffs = dmp.DiffCharsToLines(diffs, lineArray)

	result := make([]Diff, len(diffs))
	for i, diff := range diffs {
		result[i] = Diff{
			Operation: Operation(diff.Type),
			Value:     strings.Split(diff.Text, "\n"),
		}
	}

	return result
}
