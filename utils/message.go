package utils

func Error(desc string) map[string]interface{} {
	return map[string]interface{}{
		"ok":    false,
		"error": desc,
	}
}

func Message() map[string]interface{} {
	return map[string]interface{}{
		"ok": true,
	}
}
