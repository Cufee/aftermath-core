package permissions

import (
	"fmt"
	"strconv"
	"strings"
)

const version = "v2"

type Permissions uint64

func (p Permissions) Encode() string {
	return fmt.Sprintf("%s/%d", version, p)
}

func Parse(input string) Permissions {
	split := strings.Split(input, "/")
	if !strings.HasPrefix(input, version+"/") || len(split) != 2 {
		return Blank
	}

	perms, err := strconv.ParseInt(split[1], 10, 64)
	if err != nil {
		return Blank
	}

	return Permissions(perms)
}

const (
	Blank Permissions = 0

	BasicUserActions Permissions = 1 << iota
	ModerationActions
	AdminActions
)

func (p Permissions) Has(permission Permissions) bool {
	return p&permission == permission
}

func (p Permissions) Add(permission Permissions) Permissions {
	return p | permission
}

func (p Permissions) Remove(permission Permissions) Permissions {
	return p &^ permission
}

func (p Permissions) Toggle(permission Permissions) Permissions {
	return p ^ permission
}

func (p Permissions) Set(permission Permissions, enabled bool) Permissions {
	if enabled {
		return p.Add(permission)
	}
	return p.Remove(permission)
}

func (p Permissions) String() string {
	return strconv.FormatInt(int64(p), 10)
}

func ParseString(s string) (Permissions, error) {
	perms, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return Blank, err
	}
	return Permissions(perms), nil
}
