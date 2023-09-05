/*
 * App template API
 *
 * API to access and configure the app template
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package apiserver

import (
	"context"
	"net/http"
)

// ConfigurationApiRouter defines the required methods for binding the api requests to a responses for the ConfigurationApi
// The ConfigurationApiRouter implementation should parse necessary information from the http request,
// pass the data to a ConfigurationApiServicer to perform the required actions, then write the service results to the http response.
type ConfigurationApiRouter interface {
	DeleteConfigurationById(http.ResponseWriter, *http.Request)
	GetConfigurationById(http.ResponseWriter, *http.Request)
	GetConfigurations(http.ResponseWriter, *http.Request)
	PostConfiguration(http.ResponseWriter, *http.Request)
	PutConfigurationById(http.ResponseWriter, *http.Request)
}

// CustomizationApiRouter defines the required methods for binding the api requests to a responses for the CustomizationApi
// The CustomizationApiRouter implementation should parse necessary information from the http request,
// pass the data to a CustomizationApiServicer to perform the required actions, then write the service results to the http response.
type CustomizationApiRouter interface {
	GetDashboardTemplateByName(http.ResponseWriter, *http.Request)
}

// VersionApiRouter defines the required methods for binding the api requests to a responses for the VersionApi
// The VersionApiRouter implementation should parse necessary information from the http request,
// pass the data to a VersionApiServicer to perform the required actions, then write the service results to the http response.
type VersionApiRouter interface {
	GetOpenAPI(http.ResponseWriter, *http.Request)
	GetVersion(http.ResponseWriter, *http.Request)
}

// ConfigurationApiServicer defines the api actions for the ConfigurationApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type ConfigurationApiServicer interface {
	DeleteConfigurationById(context.Context, int64) (ImplResponse, error)
	GetConfigurationById(context.Context, int64) (ImplResponse, error)
	GetConfigurations(context.Context) (ImplResponse, error)
	PostConfiguration(context.Context, Configuration) (ImplResponse, error)
	PutConfigurationById(context.Context, int64, Configuration) (ImplResponse, error)
}

// CustomizationApiServicer defines the api actions for the CustomizationApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type CustomizationApiServicer interface {
	GetDashboardTemplateByName(context.Context, string, string) (ImplResponse, error)
}

// VersionApiServicer defines the api actions for the VersionApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type VersionApiServicer interface {
	GetOpenAPI(context.Context) (ImplResponse, error)
	GetVersion(context.Context) (ImplResponse, error)
}