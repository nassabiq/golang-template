package metadata

type MetaData struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func ResolveMessage(defaultMsg string, msg []string) string {
	if len(msg) > 0 && msg[0] != "" {
		return msg[0]
	}
	return defaultMsg
}
