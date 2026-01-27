package response

import (
	commonpb "github.com/nassabiq/golang-template/proto/common"
)

func Error(code int32, msg string) *commonpb.MetaData {
	return &commonpb.MetaData{
		Code:    code,
		Message: msg,
	}
}

func NotFound(msg ...string) *commonpb.MetaData {
	return Error(404, ResolveMessage("data not found", msg))
}

func Unauthorized(msg ...string) *commonpb.MetaData {
	return Error(401, ResolveMessage("unauthorized", msg))
}

func BadRequest(msg ...string) *commonpb.MetaData {
	return Error(400, ResolveMessage("bad request", msg))
}

func Validation(msg ...string) *commonpb.MetaData {
	return Error(422, ResolveMessage("validation error", msg))
}

func Conflict(msg ...string) *commonpb.MetaData {
	return Error(409, ResolveMessage("conflict", msg))
}

func Forbidden(msg ...string) *commonpb.MetaData {
	return Error(403, ResolveMessage("forbidden", msg))
}

func Internal(msg ...string) *commonpb.MetaData {
	return Error(500, ResolveMessage("internal server error", msg))
}

func ResolveMessage(defaultMsg string, msg []string) string {
	if len(msg) > 0 && msg[0] != "" {
		return msg[0]
	}
	return defaultMsg
}
