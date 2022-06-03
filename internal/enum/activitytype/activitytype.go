package activitytype

type ActivityType string

const (
	CREATED         ActivityType = "CREATED"
	UPDATED                      = "UPDATED"
	ADMIN_NOTE                   = "ADMIN_NOTE"
	SIGN_IN                      = "SIGN_IN"
	KEYCHAIN_UPSERT              = "KEYCHAIN_UPSERT"
	KEYCHAIN_DELETE              = "KEYCHAIN_DELETE"
)

var types = [...]string{"CREATED", "UPDATED", "ADMIN_NOTE", "SIGN_IN", "KEYCHAIN_UPSERT", "KEYCHAIN_DELETE"}

func String(e ActivityType) string {
	x := string(e)
	for _, v := range types {
		if v == x {
			return x
		}
	}

	return ""
}
