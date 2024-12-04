package shortener

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
	"net/url"
	"strings"
)

// Endpoints collects all of the endpoints that compose a profile service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
//
// In a server, it's useful for functions that need to operate on a per-endpoint
// basis. For example, you might pass an Endpoints to a function that produces
// an http.Handler, with each method (endpoint) wired up to a specific path. (It
// is probably a mistake in design to invoke the Service methods on the
// Endpoints struct in a server.)
//
// In a client, it's useful to collect individually constructed endpoints into a
// single type that implements the Service interface. For example, you might
// construct individual endpoints using transport/http.NewClient, combine them
// into an Endpoints, and return it to the caller as a Service.
type Endpoints struct {
	ResolveEndpoint endpoint.Endpoint
	ShortenEndpoint endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided service. Useful in a shortener
// server.
func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		ShortenEndpoint: MakePostShortenEndpoint(s),
		ResolveEndpoint: MakeGetResolveEndpoint(s),
	}
}

// MakeClientEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the remote instance, via a transport/http.Client.
// Useful in a shortener client.
func MakeClientEndpoints(instance string) (Endpoints, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	tgt, err := url.Parse(instance)
	if err != nil {
		return Endpoints{}, err
	}
	tgt.Path = ""

	options := []httptransport.ClientOption{}

	// Note that the request encoders need to modify the request URL, changing
	// the path. That's fine: we simply need to provide specific encoders for
	// each endpoint.

	return Endpoints{
		ShortenEndpoint: httptransport.NewClient("POST", tgt, encodePostShortenRequest, decodePostShortenResponse, options...).Endpoint(),
		ResolveEndpoint: httptransport.NewClient("GET", tgt, encodeGetResolveRequest, decodeGetResolveResponse, options...).Endpoint(),
	}, nil
}

//
//// PostProfilee implements Service. Primarily useful in a client.
//func (e Endpoints) PostProfile(ctx context.Context, p Shorten) error {
//	request := postShortenRequest{Shorten: p}
//	response, err := e.ShortenEndpoint(ctx, request)
//	if err != nil {
//		return err
//	}
//	resp := response.(postShortenResponse)
//	return resp.Err
//}
//
//// GetProfilee implements Service. Primarily useful in a client.
//func (e Endpoints) Resolve(ctx context.Context, id string) (Shorten, error) {
//	request := getResolveRequest{ID: id}
//	response, err := e.ResolveEndpoint(ctx, request)
//	if err != nil {
//		return Shorten{}, err
//	}
//	resp := response.(getResolveResponse)
//	return Shorten{}, resp.Err
//}

// MakePostShortenEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePostShortenEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postShortenRequest)
		shorten, e := s.PostShorten(ctx, req.Shorten)
		return postShortenResponse{Shorten: shorten, Err: e}, nil
	}
}

// MakeGetResolveEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetResolveEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getResolveRequest)
		url, err := s.Resolve(ctx, req.Url)
		return getResolveResponse{
			URL: url,
			Err: err}, nil
	}
}

// We have two options to return errors from the business logic.
//
// We could return the error via the endpoint itself. That makes certain things
// a little bit easier, like providing non-200 HTTP responses to the client. But
// Go kit assumes that endpoint errors are (or may be treated as)
// transport-domain errors. For example, an endpoint error will count against a
// circuit breaker error count.
//
// Therefore, it's often better to return service (business logic) errors in the
// response object. This means we have to do a bit more work in the HTTP
// response encoder to detect e.g. a not-found error and provide a proper HTTP
// status code. That work is done with the errorer interface, in transport.go.
// Response types that may contain business-logic errors implement that
// interface.

type postShortenRequest struct {
	Shorten Shorten
}

type postShortenResponse struct {
	Shorten ShortenResp `json:"shorten,omitempty"`
	Err     error       `json:"err,omitempty"`
}

func (r postShortenResponse) error() error { return r.Err }

type getResolveRequest struct {
	ID  string
	Url string
}

type getResolveResponse struct {
	URL string `json:"url,omitempty"`
	Err error  `json:"err,omitempty"`
}

// If you need custom status codes
type redirectResponse interface {
	getRedirectURL() string
	getStatusCode() int
	error() error
}

func (r getResolveResponse) error() error { return r.Err }

func (r getResolveResponse) getRedirectURL() string {
	return r.URL
}

func (r getResolveResponse) getStatusCode() int {
	return http.StatusTemporaryRedirect // or custom status code
}
