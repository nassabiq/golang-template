package response

import (
	commonpb "github.com/nassabiq/golang-template/proto/common"
)

func Success(code int32, msg string) *commonpb.MetaData {
	return &commonpb.MetaData{
		Code:    code,
		Message: msg,
	}
}

func Created(msg ...string) *commonpb.MetaData {
	return Success(201, ResolveMessage("data created successfully", msg))
}

func Deleted(msg ...string) *commonpb.MetaData {
	return Success(201, ResolveMessage("data deleted successfully", msg))
}
