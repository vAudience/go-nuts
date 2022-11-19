package gonuts

import (
	"errors"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

const idAlphabet string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"

var UUIDRegEx *regexp.Regexp = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
var NotLegalIdCharacters *regexp.Regexp = regexp.MustCompile("[^A-Za-z0-9-_]")

var ErrBadUUID error = errors.New("uuid format error")
var ErrBadId error = errors.New("bad id format")
var ErrIllegalId error = errors.New("illegal id")
var ErrUnknownId error = errors.New("unknown id")
var ErrMalformedId error = errors.New("malformed id")

func NanoID(prefix string) (nid string) {
	nid, err := gonanoid.Generate(idAlphabet, 12)
	if err != nil {
		L.Error(err)
		nid = strconv.FormatInt(time.Now().UnixMicro(), 10)
	}
	nid = prefix + "_" + nid
	return nid
}

func NID(prefix string, length int) (nid string) {
	nid, err := gonanoid.Generate(idAlphabet, length)
	if err != nil {
		L.Error(err)
		nid = strconv.FormatInt(time.Now().UnixMicro(), 10)
	}
	if len(prefix) > 0 {
		nid = prefix + "_" + nid
	}
	return nid
}

const AllowedIdCharacters string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"

func CreateId(IdType string) (string, error) {
	var id string
	sid := NanoID("")
	switch IdType {
	case "NClient":
		id = "nc" + "_" + sid
	case "Subtoken":
		id = "st" + "_" + sid
	case "ApiKey":
		id = "ak" + "_" + sid
	case "RealmId":
		id = "rl" + "_" + sid
	case "TenantId":
		id = "tn" + "_" + sid
	case "GroupId":
		id = "gp" + "_" + sid
	case "BroadcastId":
		id = "bc" + "_" + sid
	case "ChannelId":
		id = "ch" + "_" + sid
	case "UserId":
		id = "us" + "_" + sid
	case "PinId":
		id = "pn" + "_" + sid
	case "FileId":
		id = "fi" + "_" + sid
	default:
		id = sid
	}
	return id, nil
}

func ValidateId(IdType string, value string) error {
	switch IdType {
	case "NClient":
		return IdOkay(value, "nc")
	case "Subtoken":
		return IdOkay(value, "st")
	case "ApiKey":
		return IdOkay(value, "ak")
	case "RealmId":
		return IdOkay(value, "rl")
	case "TenantId":
		return IdOkay(value, "tn")
	case "FusionAuthTenantId":
		isUUID := IsValidUUID(value)
		if !isUUID {
			return ErrBadUUID
		}
		return nil
	case "BroadcastId":
		return IdOkay(value, "bc")
	case "GroupId":
		return IdOkay(value, "gp")
	case "ChannelId":
		return IdOkay(value, "ch")
	case "UserId":
		return IdOkay(value, "us")
	case "PinId":
		return IdOkay(value, "pn")
	case "FileId":
		return IdOkay(value, "fi")
	}
	return ErrUnknownId
}

func IdOkay(value string, prefix string) error {
	if len(value) < 8 {
		return ErrMalformedId
	}
	if !strings.HasPrefix(value, prefix) {
		// L.Debugf("!!!!!!!!!!!!!! IdOkay fail for (%s) with prefix(%s) => shortened=(%s)", value, prefix, value[3:])
		return ErrBadId
	}
	// L.Debugf("--- checking allowed chars %s <-> %s", value, NotLegalIdCharacters.MatchString)
	if NotLegalIdCharacters.MatchString(value) {
		// L.Debugf("!!!!!!!!!!!!!! IdOkay.NotLegalIdCharacters fail for (%s) with prefix(%s) => shortened=(%s)", value, prefix, value[3:])
		return ErrIllegalId
	}
	return nil
}

func IsValidUUID(uuid string) bool {
	return UUIDRegEx.MatchString(uuid)
}

func GenerateRandomString(letters []rune, length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
