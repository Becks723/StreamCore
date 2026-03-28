package util

func IsPageParamsValid(limit, page int) bool {
	if page < 0 {
		return false
	}
	if limit <= 0 {
		return false
	}
	return true
}
