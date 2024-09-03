package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"game-mining-server/entities"
	"github.com/google/uuid"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	LowwerAlphaNums = "abcdefghijklmnopqrstuvwxyz0123456789"
	UpperAlphaNums  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	Alphas          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AlphaNums       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	Letters         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#%^&*"
	Salt            = "Dk1$nLz*A2I@b0v.Ek39lG&M"
)

var randSrc = rand.New(rand.NewSource(time.Now().UnixNano()))

// Any return A if expr is true, else B
func Any[T any](expr bool, a, b T) T {
	if expr {
		return a
	}
	return b
}

func MustOpenFile(filename string) *os.File {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return file
}

func RandRangeInt(start int64, end int64) int64 {
	return randSrc.Int63n(end) + start
}

func Md5(input string) string {
	sum := md5.Sum([]byte(input))
	return hex.EncodeToString(sum[:])
}

func Sha256(input string) string {
	h256h := sha256.New()
	h256h.Write([]byte(input))
	return hex.EncodeToString(h256h.Sum(nil))
}

// RandStr generate a rand string by length
func RandStr(length int, alphanum bool) string {
	source := Any(alphanum, AlphaNums, Letters)
	return RandInSource(length, source)
}

func RandInSource(length int, source string) string {
	b := make([]byte, length)
	sourceLen := len(source)
	for i := 0; i < length; i++ {
		if i == 0 || i == length-1 {
			b[i] = source[randSrc.Intn(sourceLen)]
		} else {
			nextRand := source[randSrc.Intn(sourceLen)]
			if b[i-1] == nextRand { // make sure there has no duplicate
				b[i] = source[randSrc.Intn(sourceLen)]
			} else {
				b[i] = nextRand
			}
		}
	}
	return string(b)
}

// GenSession generate a user session format: base64({randStr[8]}:{uid}:{nowTs}:{expiresSec}:{uuid})
func GenSession(uid int64, expireSec int, encryptKey string) (string, error) {
	var builder strings.Builder
	builder.WriteString(RandStr(8, true))
	builder.WriteString(":")
	builder.WriteString(strconv.FormatInt(uid, 10))
	builder.WriteString(":")
	builder.WriteString(strconv.FormatInt(time.Now().Unix(), 10))
	builder.WriteString(":")
	builder.WriteString(strconv.Itoa(expireSec))
	builder.WriteString(":")
	builder.WriteString(uuid.New().String())
	if session, e := EncryptByAes(builder.String(), encryptKey); e != nil {
		return "", e
	} else {
		return session, nil
	}
}

// ParseSession Parse a base64 encoded session to UserSession
func ParseSession(encryptedSession string, encryptKey string) (*entities.UserSession, error) {
	session, e0 := DecryptByAes(encryptedSession, encryptKey)
	if e0 != nil {
		return nil, e0
	}
	splits := strings.Split(string(session), ":")
	if len(splits) != 5 {
		return nil, fmt.Errorf("session format error")
	}
	uid, e1 := strconv.ParseInt(splits[1], 10, 64)
	if e1 != nil {
		return nil, e1
	}
	issuedAt, e2 := strconv.ParseInt(splits[2], 10, 64)
	if e2 != nil {
		return nil, e2
	}
	expiresSec, e3 := strconv.Atoi(splits[3])
	if e3 != nil {
		return nil, e3
	}
	return &entities.UserSession{
		Uid:        uid,
		IssuedAt:   issuedAt,
		ExpiresSec: expiresSec,
	}, nil
}

// GenUserSecret Generate a secret string for client side to encrypt wallet
func GenUserSecret(uid int64) string {
	uidStr := strconv.FormatInt(uid, 10)
	return Sha256(Sha256("TMW:"+uidStr+Salt) + uidStr)
}

// GenReferralCode Generate an referral code for user
func GenReferralCode(uid int64) string {
	uidStr := strconv.FormatInt(uid, 10)
	randStr := RandStr(12, false)
	randUuid := uuid.New().String()
	nowTimeNanoStr := strconv.FormatInt(time.Now().UnixNano(), 10)
	return "RF" + strings.ToUpper(Sha256(uidStr + randStr + nowTimeNanoStr + randUuid)[4:8]) + RandInSource(3, LowwerAlphaNums) + uidStr[0:1] + RandInSource(2, Alphas)
}

// GetUTC0Ts Get today's 0:00 time in ts in UTC locale
func GetUTC0Ts() int64 {
	now := time.Now()
	utcLoc, _ := time.LoadLocation("UTC")
	utc0Time := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, utcLoc)
	return utc0Time.UnixMilli()
}

func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func CalPointForDailyCheckin(basePoint int64, continuousDays int) int64 {
	if continuousDays <= 1 {
		return basePoint
	} else if continuousDays == 2 {
		return basePoint + 15
	} else if continuousDays == 3 {
		return basePoint + 25
	} else if continuousDays == 4 {
		return basePoint + 30
	} else if continuousDays == 5 {
		return basePoint + 35
	} else if continuousDays == 6 {
		return basePoint + 40
	} else if continuousDays == 7 {
		return basePoint + 50
	} else {
		return basePoint + 50
	}
}

var levelList = []int64{0, 1, 5, 10, 20, 50, 100, 150, 500, 2000, 5000, 10000, 20000, 50000}

func indexOfList(ele int64) int {
	for i, v := range levelList {
		if v == ele {
			return i
		}
	}
	return -1
}

func CalPointForInvite(lastLevel int64, curLevel int64, invitedCount int64) int64 {
	if invitedCount < curLevel { // not enough invites
		return 0
	}
	lastIndex := indexOfList(lastLevel)
	curIndex := indexOfList(curLevel)
	if lastIndex < 0 || curIndex < 0 || curIndex < lastIndex {
		return 0
	}

	// increase all level's point between lastLevel to curLevel
	totalPoint := int64(0)
	for i := lastIndex + 1; i <= curIndex; i++ {
		totalPoint += getPointByLevel(levelList[i], invitedCount)
	}
	return totalPoint
}

func getPointByLevel(level int64, invitedCount int64) int64 {
	if level == 1 && invitedCount >= level {
		return 100
	} else if level == 5 && invitedCount >= level {
		return 600
	} else if level == 10 && invitedCount >= level {
		return 1500
	} else if level == 20 && invitedCount >= level {
		return 3000
	} else if level == 50 && invitedCount >= level {
		return 6000
	} else if level == 100 && invitedCount >= level {
		return 11000
	} else if level == 150 && invitedCount >= level {
		return 13500
	} else if level == 500 && invitedCount >= level {
		return 35000
	} else if level == 2000 && invitedCount >= level {
		return 120000
	} else if level == 5000 && invitedCount >= level {
		return 200000
	} else if level == 10000 && invitedCount >= level {
		return 300000
	} else if level == 20000 && invitedCount >= level {
		return 500000
	} else if level == 50000 && invitedCount >= level {
		return 1000000
	} else {
		return 0
	}
}
