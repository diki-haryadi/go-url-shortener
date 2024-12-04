package shortener

// The shortener is just over HTTP, so we just have a single transport.go.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

// MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
// Useful in a shortener server.
func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// GET     /:url                       retrieves the given profile by id
	// POST    /api/v1                          adds another profile

	r.Methods("GET").Path("/{url}").Handler(httptransport.NewServer(
		e.ResolveEndpoint,
		decodeGetResolveRequest,
		encodeResponseRedirect,
		options...,
	))

	r.Methods("POST").Path("/api/v1/shorten").Handler(httptransport.NewServer(
		e.ShortenEndpoint,
		decodePostShortenRequest,
		encodeResponse,
		options...,
	))
	return r
}

func decodePostShortenRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postShortenRequest
	if e := json.NewDecoder(r.Body).Decode(&req.Shorten); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetResolveRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["url"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getResolveRequest{ID: id}, nil
}

func encodePostShortenRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("POST").Path("/api/v1")
	req.URL.Path = "/api/v1/shorten"
	return encodeRequest(ctx, req, request)
}

func encodeGetResolveRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("GET").Path("/{url}")
	r := request.(getResolveRequest)
	url := url.QueryEscape(r.Url)
	req.URL.Path = "/" + url
	return encodeRequest(ctx, req, request)
}

func decodePostShortenResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response postShortenResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetResolveResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response getResolveResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error. For more information, read the
// big comment in endpoints.go.
type errorer interface {
	error() error
}

// Update your encode function
func encodeResponseRedirect(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	// First check for errors
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}

	// Check if this is a redirect response
	if r, ok := response.(redirectResponse); ok {
		w.Header().Set("Location", r.getRedirectURL())
		w.WriteHeader(http.StatusTemporaryRedirect) // or StatusPermanentRedirect (308)
		return nil
	}

	// Fall back to normal JSON response for non-redirect cases
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// encodeRequest likewise JSON-encodes the request to the HTTP request body.
// Don't use it directly as a transport/http.Client EncodeRequestFunc:
// shortener endpoints require mutating the HTTP method and request path.
func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrAlreadyExists, ErrInconsistentIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
