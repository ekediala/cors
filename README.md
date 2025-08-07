# cors

The `cors` package provides simple, yet powerful middleware for handling Cross-Origin Resource Sharing (CORS) in Go HTTP servers. It allows web applications to manage requests from different origins by specifying allowed origins, methods, headers, and the preflight request cache duration. This package is designed to be easily integrated into your server, while also promoting the Go ethos: "a little copying is better than a little dependency." You are encouraged to copy the code into your own projects, tailoring it to your needs without introducing dependencies.

## Features

- **Configurable Options:** Customize allowed origins, methods, headers, and cache duration (MaxAge) for preflight requests.
- **Preflight Request Handling:** Automatically manage preflight requests (HTTP OPTIONS) by setting the appropriate CORS headers.
- **Easy Integration:** Designed as simple middleware that can be used with any HTTP handler, providing a straightforward way to enable CORS on your server.

## Usage

Import the package and use the provided `CorsMiddleware` function to wrap your HTTP handlers:

```go
package main

import (
	"net/http"
	"yourmodule/cors"
)

func main() {
	mux := http.NewServeMux()

	// YourHandler is your HTTP handler implementation.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Wrap your handler with CorsMiddleware using DefaultOptions.
	corsHandler := cors.CorsMiddleware(handler, cors.DefaultOptions)
	mux.Handle("/", corsHandler)

	http.ListenAndServe(":8080", mux)
}
```

## Philosophy

In the spirit of the Go community, this package is intentionally simple and lightweight. Rather than relying on external dependencies, it encourages developers to copy and adapt the code to better suit their application's needs. This approach ensures that your projects remain minimal and easy to understand, while still benefiting from robust, production-ready functionality.

## License

This project is licensed under the MIT License. You are free to use, modify, and distribute the code as you see fit.

MIT License

Copyright (c) 2023

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.