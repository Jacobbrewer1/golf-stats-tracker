// Package rounder provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package rounder

import (
	"fmt"
	"net/http"

	uhttp "github.com/Jacobbrewer1/golf-stats-tracker/pkg/utils/http"
	"github.com/gorilla/mux"
	"github.com/oapi-codegen/runtime"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Login
	// (POST /login)
	Login(w http.ResponseWriter, r *http.Request)
	// Get rounds
	// (GET /rounds)
	GetRounds(w http.ResponseWriter, r *http.Request)
	// Create a round
	// (POST /rounds)
	CreateRound(w http.ResponseWriter, r *http.Request)
	// Get courses to start a round
	// (GET /rounds/new/courses)
	GetNewRoundCourses(w http.ResponseWriter, r *http.Request, params GetNewRoundCoursesParams)
	// Get the marker used for a round
	// (GET /rounds/new/marker/{course_id})
	GetNewRoundMarker(w http.ResponseWriter, r *http.Request, courseId PathCourseId)
	// Get the stats for all rounds
	// (GET /rounds/stats/charts/line/averages)
	GetLineChartAverages(w http.ResponseWriter, r *http.Request, params GetLineChartAveragesParams)
	// Get the stats for all rounds
	// (GET /rounds/stats/charts/pie/averages)
	GetPieChartAverages(w http.ResponseWriter, r *http.Request, params GetPieChartAveragesParams)
	// Get the holes for a round
	// (GET /rounds/{round_id}/holes)
	GetRoundHoles(w http.ResponseWriter, r *http.Request, roundId PathRoundId)
	// Get the stats for a hole
	// (GET /rounds/{round_id}/holes/{hole_id}/stats)
	GetHoleStats(w http.ResponseWriter, r *http.Request, roundId PathRoundId, holeId PathHoleId)
	// Update the stats for a hole
	// (POST /rounds/{round_id}/holes/{hole_id}/stats)
	UpdateHoleStats(w http.ResponseWriter, r *http.Request, roundId PathRoundId, holeId PathHoleId)
	// Create a user
	// (POST /users)
	CreateUser(w http.ResponseWriter, r *http.Request)
}

type RateLimiterFunc = func(http.ResponseWriter, *http.Request) error
type MetricsMiddlewareFunc = http.HandlerFunc
type ErrorHandlerFunc = func(http.ResponseWriter, *http.Request, error)

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	authz             ServerInterface
	handler           ServerInterface
	rateLimiter       RateLimiterFunc
	metricsMiddleware MetricsMiddlewareFunc
	errorHandlerFunc  ErrorHandlerFunc
}

// WithAuthorization applies the passed authorization middleware to the server.
func WithAuthorization(authz ServerInterface) ServerOption {
	return func(s *ServerInterfaceWrapper) {
		s.authz = authz
	}
}

// WithRateLimiter applies the rate limiter middleware to routes with x-global-rate-limit.
func WithRateLimiter(rateLimiter RateLimiterFunc) ServerOption {
	return func(s *ServerInterfaceWrapper) {
		s.rateLimiter = rateLimiter
	}
}

// WithErrorHandlerFunc sets the error handler function for the server.
func WithErrorHandlerFunc(errorHandlerFunc ErrorHandlerFunc) ServerOption {
	return func(s *ServerInterfaceWrapper) {
		s.errorHandlerFunc = errorHandlerFunc
	}
}

// WithMetricsMiddleware applies the metrics middleware to the server.
func WithMetricsMiddleware(middleware MetricsMiddlewareFunc) ServerOption {
	return func(s *ServerInterfaceWrapper) {
		s.metricsMiddleware = middleware
	}
}

// ServerOption represents an optional feature applied to the server.
type ServerOption func(s *ServerInterfaceWrapper)

// Login operation middleware
func (siw *ServerInterfaceWrapper) Login(w http.ResponseWriter, r *http.Request) {
	cw := uhttp.NewClientWriter(w)
	ctx := r.Context()

	defer func() {
		if siw.metricsMiddleware != nil {
			siw.metricsMiddleware(cw, r)
		}
	}()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.handler.Login(cw, r)
	}))

	handler.ServeHTTP(cw, r.WithContext(ctx))
}

// GetRounds operation middleware
func (siw *ServerInterfaceWrapper) GetRounds(w http.ResponseWriter, r *http.Request) {
	cw := uhttp.NewClientWriter(w)
	ctx := r.Context()

	defer func() {
		if siw.metricsMiddleware != nil {
			siw.metricsMiddleware(cw, r)
		}
	}()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if siw.authz != nil {
			siw.authz.GetRounds(cw, r.WithContext(ctx))
			return
		}
	}))

	handler.ServeHTTP(cw, r.WithContext(ctx))
}

// CreateRound operation middleware
func (siw *ServerInterfaceWrapper) CreateRound(w http.ResponseWriter, r *http.Request) {
	cw := uhttp.NewClientWriter(w)
	ctx := r.Context()

	defer func() {
		if siw.metricsMiddleware != nil {
			siw.metricsMiddleware(cw, r)
		}
	}()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if siw.authz != nil {
			siw.authz.CreateRound(cw, r.WithContext(ctx))
			return
		}
	}))

	handler.ServeHTTP(cw, r.WithContext(ctx))
}

// GetNewRoundCourses operation middleware
func (siw *ServerInterfaceWrapper) GetNewRoundCourses(w http.ResponseWriter, r *http.Request) {
	cw := uhttp.NewClientWriter(w)
	ctx := r.Context()

	defer func() {
		if siw.metricsMiddleware != nil {
			siw.metricsMiddleware(cw, r)
		}
	}()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetNewRoundCoursesParams

	// ------------- Optional query parameter "name" -------------

	err = runtime.BindQueryParameter("form", true, false, "name", r.URL.Query(), &params.Name)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "name", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if siw.authz != nil {
			siw.authz.GetNewRoundCourses(cw, r.WithContext(ctx), params)
			return
		}
	}))

	handler.ServeHTTP(cw, r.WithContext(ctx))
}

// GetNewRoundMarker operation middleware
func (siw *ServerInterfaceWrapper) GetNewRoundMarker(w http.ResponseWriter, r *http.Request) {
	cw := uhttp.NewClientWriter(w)
	ctx := r.Context()

	defer func() {
		if siw.metricsMiddleware != nil {
			siw.metricsMiddleware(cw, r)
		}
	}()

	var err error

	// ------------- Path parameter "course_id" -------------
	var courseId PathCourseId

	err = runtime.BindStyledParameterWithOptions("simple", "course_id", mux.Vars(r)["course_id"], &courseId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "course_id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if siw.authz != nil {
			siw.authz.GetNewRoundMarker(cw, r.WithContext(ctx), courseId)
			return
		}
	}))

	handler.ServeHTTP(cw, r.WithContext(ctx))
}

// GetLineChartAverages operation middleware
func (siw *ServerInterfaceWrapper) GetLineChartAverages(w http.ResponseWriter, r *http.Request) {
	cw := uhttp.NewClientWriter(w)
	ctx := r.Context()

	defer func() {
		if siw.metricsMiddleware != nil {
			siw.metricsMiddleware(cw, r)
		}
	}()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetLineChartAveragesParams

	// ------------- Required query parameter "average_type" -------------

	if paramValue := r.URL.Query().Get("average_type"); paramValue != "" {

	} else {
		siw.errorHandlerFunc(cw, r, &RequiredParamError{ParamName: "average_type"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "average_type", r.URL.Query(), &params.AverageType)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "average_type", Err: err})
		return
	}

	// ------------- Optional query parameter "from_date" -------------

	err = runtime.BindQueryParameter("form", true, false, "from_date", r.URL.Query(), &params.FromDate)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "from_date", Err: err})
		return
	}

	// ------------- Optional query parameter "since" -------------

	err = runtime.BindQueryParameter("form", true, false, "since", r.URL.Query(), &params.Since)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "since", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if siw.authz != nil {
			siw.authz.GetLineChartAverages(cw, r.WithContext(ctx), params)
			return
		}
	}))

	handler.ServeHTTP(cw, r.WithContext(ctx))
}

// GetPieChartAverages operation middleware
func (siw *ServerInterfaceWrapper) GetPieChartAverages(w http.ResponseWriter, r *http.Request) {
	cw := uhttp.NewClientWriter(w)
	ctx := r.Context()

	defer func() {
		if siw.metricsMiddleware != nil {
			siw.metricsMiddleware(cw, r)
		}
	}()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetPieChartAveragesParams

	// ------------- Required query parameter "average_type" -------------

	if paramValue := r.URL.Query().Get("average_type"); paramValue != "" {

	} else {
		siw.errorHandlerFunc(cw, r, &RequiredParamError{ParamName: "average_type"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "average_type", r.URL.Query(), &params.AverageType)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "average_type", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if siw.authz != nil {
			siw.authz.GetPieChartAverages(cw, r.WithContext(ctx), params)
			return
		}
	}))

	handler.ServeHTTP(cw, r.WithContext(ctx))
}

// GetRoundHoles operation middleware
func (siw *ServerInterfaceWrapper) GetRoundHoles(w http.ResponseWriter, r *http.Request) {
	cw := uhttp.NewClientWriter(w)
	ctx := r.Context()

	defer func() {
		if siw.metricsMiddleware != nil {
			siw.metricsMiddleware(cw, r)
		}
	}()

	var err error

	// ------------- Path parameter "round_id" -------------
	var roundId PathRoundId

	err = runtime.BindStyledParameterWithOptions("simple", "round_id", mux.Vars(r)["round_id"], &roundId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "round_id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if siw.authz != nil {
			siw.authz.GetRoundHoles(cw, r.WithContext(ctx), roundId)
			return
		}
	}))

	handler.ServeHTTP(cw, r.WithContext(ctx))
}

// GetHoleStats operation middleware
func (siw *ServerInterfaceWrapper) GetHoleStats(w http.ResponseWriter, r *http.Request) {
	cw := uhttp.NewClientWriter(w)
	ctx := r.Context()

	defer func() {
		if siw.metricsMiddleware != nil {
			siw.metricsMiddleware(cw, r)
		}
	}()

	var err error

	// ------------- Path parameter "round_id" -------------
	var roundId PathRoundId

	err = runtime.BindStyledParameterWithOptions("simple", "round_id", mux.Vars(r)["round_id"], &roundId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "round_id", Err: err})
		return
	}

	// ------------- Path parameter "hole_id" -------------
	var holeId PathHoleId

	err = runtime.BindStyledParameterWithOptions("simple", "hole_id", mux.Vars(r)["hole_id"], &holeId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "hole_id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if siw.authz != nil {
			siw.authz.GetHoleStats(cw, r.WithContext(ctx), roundId, holeId)
			return
		}
	}))

	handler.ServeHTTP(cw, r.WithContext(ctx))
}

// UpdateHoleStats operation middleware
func (siw *ServerInterfaceWrapper) UpdateHoleStats(w http.ResponseWriter, r *http.Request) {
	cw := uhttp.NewClientWriter(w)
	ctx := r.Context()

	defer func() {
		if siw.metricsMiddleware != nil {
			siw.metricsMiddleware(cw, r)
		}
	}()

	var err error

	// ------------- Path parameter "round_id" -------------
	var roundId PathRoundId

	err = runtime.BindStyledParameterWithOptions("simple", "round_id", mux.Vars(r)["round_id"], &roundId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "round_id", Err: err})
		return
	}

	// ------------- Path parameter "hole_id" -------------
	var holeId PathHoleId

	err = runtime.BindStyledParameterWithOptions("simple", "hole_id", mux.Vars(r)["hole_id"], &holeId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "hole_id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if siw.authz != nil {
			siw.authz.UpdateHoleStats(cw, r.WithContext(ctx), roundId, holeId)
			return
		}
	}))

	handler.ServeHTTP(cw, r.WithContext(ctx))
}

// CreateUser operation middleware
func (siw *ServerInterfaceWrapper) CreateUser(w http.ResponseWriter, r *http.Request) {
	cw := uhttp.NewClientWriter(w)
	ctx := r.Context()

	defer func() {
		if siw.metricsMiddleware != nil {
			siw.metricsMiddleware(cw, r)
		}
	}()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.handler.CreateUser(cw, r)
	}))

	handler.ServeHTTP(cw, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// wrapHandler will wrap the handler with middlewares in the other specified
// making the execution order the inverse of the parameter declaration
func wrapHandler(handler http.HandlerFunc, middlewares ...mux.MiddlewareFunc) http.Handler {
	var wrappedHandler http.Handler = handler
	for _, middleware := range middlewares {
		if middleware == nil {
			continue
		}
		wrappedHandler = middleware(wrappedHandler)
	}
	return wrappedHandler
}

// RegisterHandlers registers the api handlers.
func RegisterHandlers(router *mux.Router, si ServerInterface, opts ...ServerOption) {
	wrapper := ServerInterfaceWrapper{
		handler: si,
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(&wrapper)
	}

	router.Use(uhttp.AuthHeaderToContextMux())
	router.Use(uhttp.RequestIDToContextMux())

	router.Methods(http.MethodGet).Path("/rounds").Handler(wrapHandler(wrapper.GetRounds))

	router.Methods(http.MethodPost).Path("/rounds").Handler(wrapHandler(wrapper.CreateRound))

	router.Methods(http.MethodGet).Path("/rounds/new/courses").Handler(wrapHandler(wrapper.GetNewRoundCourses))

	router.Methods(http.MethodGet).Path("/rounds/new/marker/{course_id}").Handler(wrapHandler(wrapper.GetNewRoundMarker))

	router.Methods(http.MethodGet).Path("/rounds/stats/charts/line/averages").Handler(wrapHandler(wrapper.GetLineChartAverages))

	router.Methods(http.MethodGet).Path("/rounds/stats/charts/pie/averages").Handler(wrapHandler(wrapper.GetPieChartAverages))

	router.Methods(http.MethodGet).Path("/rounds/{round_id}/holes").Handler(wrapHandler(wrapper.GetRoundHoles))

	router.Methods(http.MethodGet).Path("/rounds/{round_id}/holes/{hole_id}/stats").Handler(wrapHandler(wrapper.GetHoleStats))

	router.Methods(http.MethodPost).Path("/rounds/{round_id}/holes/{hole_id}/stats").Handler(wrapHandler(wrapper.UpdateHoleStats))
}

// RegisterUnauthedHandlers registers any api handlers which do not have any authentication on them. Most services will not have any.
func RegisterUnauthedHandlers(router *mux.Router, si ServerInterface, opts ...ServerOption) {
	wrapper := ServerInterfaceWrapper{
		handler: si,
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(&wrapper)
	}

	router.Use(uhttp.AuthHeaderToContextMux())
	router.Use(uhttp.RequestIDToContextMux())

	// We do not have a gateway preparer here as no auth is sent.

	router.Methods(http.MethodPost).Path("/login").Handler(wrapHandler(wrapper.Login))

	router.Methods(http.MethodPost).Path("/users").Handler(wrapHandler(wrapper.CreateUser))
}
