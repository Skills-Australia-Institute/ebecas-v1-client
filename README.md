# eBECAS V1 Client

A Go client for interacting with the eBECAS V1 API.

## Requirements

- Go 1.26.5 or later
- eBECAS V1 API credentials:
  - College code
  - Username
  - Authentication token

## Installation

```bash
go get github.com/skills-australia-institute/ebecas-v1-client
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	ebecasv1client "github.com/skills-australia-institute/ebecas-v1-client"
)

func main() {
	client, err := ebecasv1client.NewClient(
		ebecasv1client.Config{
			BaseURL:     "https://example.equatorit.net/ebecas.portal/api",
			CollegeCode: "YOUR_COLLEGE_CODE",
			Username:    "YOUR_USERNAME",
			AuthToken:   "YOUR_AUTH_TOKEN",
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	result, statusCode, err := client.GetAcademicClasses(
		ctx,
		ebecasv1client.GetAcademicClassesParams{
			LocationCode: "Perth",
			ClassDate:    "2026-07-06",
			PageSize:     100,
			Page:         1,
		},
	)
	if err != nil {
		log.Fatalf("request failed with status %d: %v", statusCode, err)
	}

	fmt.Println(statusCode)
	fmt.Println(result.Classes)
}
```

## Configuration

Create a client using `NewClient` and provide a `Config`:

```go
client, err := ebecasv1client.NewClient(
	ebecasv1client.Config{
		BaseURL:     "https://example.equatorit.net/ebecas.portal/api",
		CollegeCode: "YOUR_COLLEGE_CODE",
		Username:    "YOUR_USERNAME",
		AuthToken:   "YOUR_AUTH_TOKEN",
	},
)
if err != nil {
	log.Fatal(err)
}
```

### Configuration Options

| Option        | Required | Default           | Description                                                                          |
| ------------- | -------- | ----------------- | ------------------------------------------------------------------------------------ |
| `BaseURL`     | Yes      | -                 | Base URL of the eBECAS V1 API.                                                       |
| `CollegeCode` | Yes      | -                 | eBECAS college code used for authentication.                                         |
| `Username`    | Yes      | -                 | eBECAS username used for API requests.                                               |
| `AuthToken`   | Yes      | -                 | Base64-encoded authentication token used with Basic authentication.                  |
| `HTTPClient`  | No       | 30-second timeout | Optional custom `*http.Client`. If `nil`, a client with a 30-second timeout is used. |

### Custom HTTP Client

You can provide a custom `http.Client` when you need a different timeout or HTTP configuration:

```go
httpClient := &http.Client{
	Timeout: 60 * time.Second,
}

client, err := ebecasv1client.NewClient(
	ebecasv1client.Config{
		BaseURL:     "https://example.equatorit.net/ebecas.portal/api",
		CollegeCode: "YOUR_COLLEGE_CODE",
		Username:    "YOUR_USERNAME",
		AuthToken:   "YOUR_AUTH_TOKEN",
		HTTPClient:  httpClient,
	},
)
if err != nil {
	log.Fatal(err)
}
```

## Authentication

The client uses Basic authentication and the eBECAS-specific headers required by the V1 API.

Requests include:

```http
Authorization: Basic YOUR_AUTH_TOKEN
COLLEGECODE: YOUR_COLLEGE_CODE
USERNAME: YOUR_USERNAME
Accept: application/json
```

Do not commit API credentials to source control. Use environment variables or a secure secrets-management solution.

For example:

```go
baseURL := os.Getenv("EBECAS_BASE_URL")
collegeCode := os.Getenv("EBECAS_COLLEGE_CODE")
username := os.Getenv("EBECAS_USERNAME")
authToken := os.Getenv("EBECAS_AUTH_TOKEN")
```

Then:

```go
client, err := ebecasv1client.NewClient(
	ebecasv1client.Config{
		BaseURL:     baseURL,
		CollegeCode: collegeCode,
		Username:    username,
		AuthToken:   authToken,
	},
)
if err != nil {
	log.Fatal(err)
}
```

## Academic Classes

### Get Academic Classes

Use `GetAcademicClasses` to retrieve academic classes for a specified date and location.

```go
result, statusCode, err := client.GetAcademicClasses(
	ctx,
	ebecasv1client.GetAcademicClassesParams{
		LocationCode: "Perth",
		ClassDate:    "2026-07-06",
		PageSize:     100,
		Page:         1,
	},
)
if err != nil {
	log.Printf("request failed with status %d: %v", statusCode, err)
	return
}

for _, class := range result.Classes {
	fmt.Println(class.ClassID, class.ClassName)
}
```

To retrieve classes from all locations:

```go
result, statusCode, err := client.GetAcademicClasses(
	ctx,
	ebecasv1client.GetAcademicClassesParams{
		AllLocations: true,
		ClassDate:    "2026-07-06",
		PageSize:     100,
		Page:         1,
	},
)
```

### Pagination

Academic class results are paginated.

The default page size is `100`, and the maximum supported page size is `1000`.

`Page` represents the page number to request. Internally, eBECAS expects this value using the `PageCount` query parameter.

The response includes:

- `PageSize` - number of records per page
- `Page` - current page number returned by the API
- `TotalCount` - total number of matching academic classes

Example:

```go
result, statusCode, err := client.GetAcademicClasses(
	ctx,
	ebecasv1client.GetAcademicClassesParams{
		LocationCode: "Perth",
		ClassDate:    "2026-07-06",
		PageSize:     100,
		Page:         1,
	},
)
if err != nil {
	log.Printf("request failed with status %d: %v", statusCode, err)
	return
}

fmt.Println("Page:", result.Page)
fmt.Println("Page size:", result.PageSize)
fmt.Println("Total records:", result.TotalCount)
```

To retrieve all pages:

```go
page := 1
pageSize := 100

for {
	result, statusCode, err := client.GetAcademicClasses(
		ctx,
		ebecasv1client.GetAcademicClassesParams{
			LocationCode: "Perth",
			ClassDate:    "2026-07-06",
			PageSize:     pageSize,
			Page:         page,
		},
	)
	if err != nil {
		log.Fatalf("request failed with status %d: %v", statusCode, err)
	}

	for _, class := range result.Classes {
		fmt.Println(class.ClassID, class.ClassName)
	}

	if page*result.PageSize >= result.TotalCount {
		break
	}

	page++
}
```

## Error Handling

Client methods return the API response, HTTP status code, and error.

For example:

```go
result, statusCode, err := client.GetAcademicClasses(
	ctx,
	ebecasv1client.GetAcademicClassesParams{
		LocationCode: "Perth",
	},
)
```

Always check the returned error before using the response:

```go
if err != nil {
	log.Printf("request failed with status %d: %v", statusCode, err)
	return
}
```

Validation errors such as an invalid class ID or page size return `http.StatusBadRequest`.

For example, `PageSize` must not exceed `1000`.

## Context

API methods accept a `context.Context`, allowing requests to be cancelled or given a timeout.

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

result, statusCode, err := client.GetAcademicClasses(
	ctx,
	ebecasv1client.GetAcademicClassesParams{
		LocationCode: "Perth",
	},
)
```

## License

This project is licensed under the [MIT License](https://github.com/skills-australia-institute/ebecas-v1-client/blob/main/LICENSE).
