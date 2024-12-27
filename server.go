package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

func (ds *switchdefinitions) hauth(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		auth := r.FormValue("auth")
		if auth == ds.auth {
			// reset timer
			fmt.Println("Resetting timers")
			ds.resetTimers()
		}
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		render(w, `<!DOCTYPE html>
		<html>
			<head>
				<meta name="HandheldFriendly" content="true" />
				<meta name="MobileOptimized" content="320" />
				<meta name="viewport" content="initial-scale=1.0, maximum-scale=1.0, width=device-width, user-scalable=no" />
				<script src="https://unpkg.com/htmx.org@1.9.2" integrity="sha384-L6OqL9pRWyyFU3+/bjdSri+iIphTN/bvYyM37tICVyOJkWZLpP2vGn6VUEXgzg6h" crossorigin="anonymous"></script>
				<title>Dead Switch</title>
			</head>
			<body>
				<div class="container">
					<div class="centertext" id="inputbox">
						Input auth:<br>
						<input id="auth" name="auth" type="password">
						<button hx-post="/" hx-target="html" hx-swap="none" hx-include="#auth" hx-trigger="click, keydown[keyCode==13&&shiftKey!=true] from:#inputbox">Reset</button>
						<button hx-delete="/" hx-target="html" hx-swap="none" hx-include="#auth" hx-trigger="click">Stop</button>
					</div>
				</div>
			</body>
		</html>`, nil)
	}
	if r.Method == http.MethodDelete {
		auth := r.FormValue("auth")
		if auth == ds.auth {
			// shutdown
			fmt.Println("Shutting down...")
			os.Exit(0)
		}
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func render(w http.ResponseWriter, html string, data any) {
	// Render the HTML template
	// fmt.Println("Rendering...")
	w.WriteHeader(http.StatusOK)
	tmpl, err := template.New(html).Parse(html)
	if err != nil {
		fmt.Println(err)
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
