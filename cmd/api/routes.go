package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	// Initialize a new httprouter router instance.
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	//======================================================================================================
	// healthcheck handler
	{
		router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	}

	//======================================================================================================
	// movies handler
	{
		router.HandlerFunc(http.MethodGet, "/v1/movies", app.requirePermission("movies:read", app.listMoviesHandler))
		router.HandlerFunc(http.MethodPost, "/v1/movies", app.requirePermission("movies:write", app.createMovieHandler))
		router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.requirePermission("movies:write", app.showMovieHandler))
		router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.requirePermission("movies:read", app.updateMovieHandler))
		router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.requirePermission("movies:write", app.deleteMovieHandler))
	}

	//======================================================================================================
	// users handler
	{
		router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
		router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
		router.HandlerFunc(http.MethodPut, "/v1/users/password", app.updateUserPasswordHandler)
	}

	//======================================================================================================
	// tokens handler
	{
		router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
		router.HandlerFunc(http.MethodPost, "/v1/tokens/password-reset", app.createPasswordResetTokenHandler)
		router.HandlerFunc(http.MethodPost, "/v1/tokens/activation", app.createActivationTokenHandler)
	}

	//======================================================================================================
	// metrics handler
	{
		router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())
	}

	// Return the httprouter instance.
	return app.metrics(app.recoverPanic(app.enalbeCORS(app.rateLimit(app.authentication(router)))))
}
