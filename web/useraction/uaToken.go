package useraction

import (
	"math/rand/v2"
	"strconv"
	"strings"
	"time"
)

type Token struct {
	Id, Time uint64
}

func (t *Token) Valid() bool {
	return t.Id != 0 && t.Time != 0
}

func (t *Token) String() string {
	return toString(t.Id) + "-" + toString(t.Time)
}

func timeNow() uint64 {
	return uint64(time.Now().UnixMicro())
}

func timeWhen(t uint64) (ret time.Time) {
	return time.UnixMicro(int64(t))
}

func toString(v uint64) string {
	return strconv.FormatUint(v, 36)
}

// ignores any parsing errors
func fromString(s string) (ret uint64) {
	if v, e := strconv.ParseUint(s, 36, 64); e == nil {
		ret = v
	}
	return
}

// ignores any parsing errors and returns an empty token
func MakeToken() (ret Token) {
	return Token{
		Id:   rand.Uint64(),
		Time: timeNow(),
	}
}

// ignores any parsing errors and returns an empty token
func ReadToken(s string) (ret Token) {
	if split := strings.IndexRune(s, '-'); split > 0 {
		lhs, rhs := s[:split], s[split+1:]
		ret = Token{
			Id:   fromString(lhs),
			Time: fromString(rhs),
		}
	}
	return
}
