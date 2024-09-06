package util

func DefaultRowEndFunc(s []string) bool {
	return len(s) == 0 || len(s[0]) == 0
}
