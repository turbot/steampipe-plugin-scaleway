package scaleway

import (
	"net/http"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"golang.org/x/xerrors"
)

func is404Error(err error) bool {
	notFoundError := &scw.ResourceNotFoundError{}
	return isHTTPCodeError(err, http.StatusNotFound) || xerrors.As(err, &notFoundError)
}

func isHTTPCodeError(err error, statusCode int) bool {
	if err == nil {
		return false
	}

	responseError := &scw.ResponseError{}
	if xerrors.As(err, &responseError) && responseError.StatusCode == statusCode {
		return true
	}
	return false
}
