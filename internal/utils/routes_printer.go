package utils

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"com.tom-ludwig/go-server-template/internal/api"
	"github.com/fatih/color"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
)

// PrintRoutes prints all registered routes to stdout with colors and query parameters
func PrintRoutes(r chi.Router) {
	swagger, err := api.GetSwagger()
	if err != nil {
		// If we can't load swagger, fall back to simple listing
		printRoutesSimple(r)
		return
	}

	fmt.Println()
	color.New(color.FgCyan, color.Bold).Println("Registered Routes:")
	fmt.Println(strings.Repeat("─", 100))

	type routeInfo struct {
		method      string
		route       string
		params      []string
		description string
	}

	var routes []routeInfo

	err = chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		params := []string{}
		description := ""

		// Try to find this route in the OpenAPI spec
		if swagger.Paths != nil {
			// Try exact match first
			pathItem := swagger.Paths.Find(route)
			if pathItem == nil {
				// If exact match fails, try to find by matching pattern
				// (chi routes might differ slightly from OpenAPI paths)
				pathMap := swagger.Paths.Map()
				for specPath, item := range pathMap {
					// Simple matching: if route contains specPath or vice versa
					if strings.Contains(route, specPath) || strings.Contains(specPath, route) {
						pathItem = item
						break
					}
				}
			}

			if pathItem != nil {
				var operation *openapi3.Operation
				switch strings.ToUpper(method) {
				case "GET":
					operation = pathItem.Get
				case "POST":
					operation = pathItem.Post
				case "PUT":
					operation = pathItem.Put
				case "DELETE":
					operation = pathItem.Delete
				case "PATCH":
					operation = pathItem.Patch
				}

				if operation != nil {
					if operation.Summary != "" {
						description = operation.Summary
					}
					// Extract query parameters
					for _, param := range operation.Parameters {
						if param != nil && param.Value != nil && param.Value.In == "query" {
							paramStr := param.Value.Name
							if param.Value.Required {
								paramStr = paramStr + "*"
							}
							params = append(params, paramStr)
						}
					}
				}
			}
		}

		routes = append(routes, routeInfo{
			method:      method,
			route:       route,
			params:      params,
			description: description,
		})
		return nil
	})

	if err != nil {
		color.New(color.FgRed).Printf("Error walking routes: %v\n", err)
		return
	}

	// Sort routes by method, then by path
	sort.Slice(routes, func(i, j int) bool {
		if routes[i].method != routes[j].method {
			return routes[i].method < routes[j].method
		}
		return routes[i].route < routes[j].route
	})

	// Print routes with colors
	for _, route := range routes {
		methodColor := getMethodColor(route.method)
		methodColor.Printf("%-8s", route.method)

		fmt.Print(" ")
		color.New(color.FgWhite).Print(route.route)

		if len(route.params) > 0 {
			fmt.Print(" ")
			color.New(color.FgYellow).Print("?")
			for i, param := range route.params {
				if i > 0 {
					fmt.Print(",")
				}
				if strings.HasSuffix(param, "*") {
					color.New(color.FgYellow, color.Bold).Print(strings.TrimSuffix(param, "*"))
				} else {
					color.New(color.FgYellow).Print(param)
				}
			}
		}

		if route.description != "" {
			fmt.Print("  ")
			color.New(color.FgHiBlack).Printf("// %s", route.description)
		}

		fmt.Println()
	}

	fmt.Println(strings.Repeat("─", 100))
	fmt.Println()
}

func printRoutesSimple(r chi.Router) {
	fmt.Println("\nRegistered Routes:")
	fmt.Println(strings.Repeat("-", 80))

	err := chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%-8s %s\n", method, route)
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking routes: %v\n", err)
	}

	fmt.Println(strings.Repeat("-", 80))
	fmt.Println()
}

func getMethodColor(method string) *color.Color {
	switch strings.ToUpper(method) {
	case "GET":
		return color.New(color.FgGreen, color.Bold)
	case "POST":
		return color.New(color.FgBlue, color.Bold)
	case "PUT":
		return color.New(color.FgYellow, color.Bold)
	case "DELETE":
		return color.New(color.FgRed, color.Bold)
	case "PATCH":
		return color.New(color.FgMagenta, color.Bold)
	default:
		return color.New(color.FgWhite, color.Bold)
	}
}
