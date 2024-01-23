package permissions

import "strconv"

type Permissions int64

const (
	Blank Permissions = 0
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
