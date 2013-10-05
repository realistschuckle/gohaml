package parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/realistschuckle/testify/mock"
	"io"
	"testing"
	"errors"
)

type mockRuneReader struct {
	mock.Mock
}

func (self *mockRuneReader) ReadRune() (rune, int, error) {
	args := self.Mock.Called()
	return args.Get(0).(rune), args.Int(1), args.Error(2)
}

func TestScannerIsARuneScanner(t *testing.T) {
	scanner := &scanner{}
	assert.Implements(t, (*io.RuneScanner)(nil), scanner)
}

func TestScannerReturnsRuneReadFromReader(t *testing.T) {
	mock := &mockRuneReader{}
	mock.On("ReadRune").Return('a', 1, nil)

	scanner := &scanner{mock, [8]rune{}, 0, 0}
	r, size, e := scanner.ReadRune()
	assert.Equal(t, 'a', r)
	assert.Equal(t, 1, size)
	assert.Nil(t, e)

	mock.AssertCalled(t, "ReadRune")
}

func TestScannerReturnsRunesInReadOrder(t *testing.T) {
	runes := [8]rune{'a', 'b', 'c', 'd', 0, 0, 0, 0}
	scanner := &scanner{nil, runes, 4, 0}

	for i := 0; i < 4; i += 1 {
		assert.Nil(t, scanner.UnreadRune())
	}

	for i := 0; i < 4; i += 1 {
		r, _, _ := scanner.ReadRune()
		assert.Equal(t, runes[i], r)
	}
}

func TestScannerReturnsRunesInReadOrderEvenAfterWrap(t *testing.T) {
	runes := [8]rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}
	scanner := &scanner{nil, runes, 13, 0}

	for i := 0; i < 8; i += 1 {
		assert.Nil(t, scanner.UnreadRune())
	}

	for i := 8; i > 0; i -= 1 {
		r, _, _ := scanner.ReadRune()
		e := runes[(13-i)%8]
		assert.Equal(t, e, r, "%c != %c", e, r)
	}
}

func TestCannotUnreadMoreThanEightRunes(t *testing.T) {
	runes := [8]rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}
	scanner := &scanner{nil, runes, 13, 0}

	for i := 0; i < 8; i += 1 {
		assert.Nil(t, scanner.UnreadRune())
	}
	assert.NotNil(t, scanner.UnreadRune())
}

func TestCannotUnreadMoreRunesThanRead(t *testing.T) {
	for i := 0; i < 9; i += 1 {
		mock := &mockRuneReader{}
		mock.On("ReadRune").Return('a', 1, nil)
		scanner := &scanner{mock, [8]rune{}, 0, 0}
		for j := 0; j < i; j += 1 {
			scanner.ReadRune()
		}

		for j := 0; j < i; j += 1 {
			assert.Nil(t, scanner.UnreadRune())
		}

		assert.NotNil(t, scanner.UnreadRune(), "%v", scanner)
	}
}

func TestWithReaderBeyondEightLimitRecovery(t *testing.T) {
	content := []rune("Hello, World")
	unreadable := content[len(content) - 8:]
	mock := &mockRuneReader{}
	var w int = 1
	var r rune
	for i := 0; i < len(content); i += 1 {
		r = content[i]
		mock.On("ReadRune").Return(r, w, nil).Once()
	}
	mock.On("ReadRune").Return('\000', w, errors.New("BOOM!"))
	scanner := &scanner{mock, [8]rune{}, 0, 0}

	for i := 0; i < len(content); i += 1 {
		scanner.ReadRune()
	}

	for i := 0; i < 8; i += 1 {
		if err := scanner.UnreadRune(); err != nil {
			assert.Fail(t, "Unreading rune failed.")
		}
	}

	for i := 0; i < 8; i += 1 {
		r, _, _ := scanner.ReadRune()
		c := unreadable[i]
		assert.Equal(t, c, r, "%c != %c", c, r)
	}
}

