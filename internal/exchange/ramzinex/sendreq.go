package ramzinex

import (
	"fmt"

	"github.com/shoshtari/paroo/internal/pkg"
)

type sendReqRequest struct {
	path    string
	reqbody any
	resbody any
	// optionals
	auth         bool
	usePublicApi bool
}

func (w ramzinexClientImp) sendReq(r sendReqRequest) error {
	var url string
	if r.usePublicApi {
		url = fmt.Sprintf("%v/%v", w.basePublicAddress, r.path)
	} else {
		url = fmt.Sprintf("%v/%v", w.baseAddress, r.path)
	}
	if r.auth {
		return pkg.SendHTTPRequest(w.httpClient, url, r.reqbody, r.resbody,
			pkg.WithHeader("Authorization", w.token),
		)
	}
	return pkg.SendHTTPRequest(w.httpClient, url, r.reqbody, r.resbody)
}
