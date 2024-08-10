package gonuts

import (
	"fmt"
	"net/url"
	"strings"
)

// URLBuilder provides a fluent interface for constructing URLs
type URLBuilder struct {
	scheme   string
	username string
	password string
	host     string
	port     string
	path     string
	query    url.Values
	fragment string
}

// NewURLBuilder creates a new URLBuilder
//
// Example:
//
//	builder := NewURLBuilder("https://api.example.com")
func NewURLBuilder(baseURL string) (*URLBuilder, error) {
	builder := &URLBuilder{
		query: make(url.Values),
	}

	if baseURL != "" {
		parsed, err := url.Parse(baseURL)
		if err != nil {
			return nil, fmt.Errorf("invalid base URL: %w", err)
		}
		builder.scheme = parsed.Scheme
		builder.host = parsed.Hostname()
		builder.port = parsed.Port()
		builder.path = parsed.Path
		builder.query = parsed.Query()
		builder.fragment = parsed.Fragment
		if parsed.User != nil {
			builder.username = parsed.User.Username()
			builder.password, _ = parsed.User.Password()
		}
	}

	return builder, nil
}

// SetScheme sets the scheme (protocol) of the URL
//
// Example:
//
//	builder.SetScheme("https")
func (b *URLBuilder) SetScheme(scheme string) *URLBuilder {
	b.scheme = scheme
	return b
}

// SetHost sets the host of the URL
//
// Example:
//
//	builder.SetHost("api.example.com")
func (b *URLBuilder) SetHost(host string) *URLBuilder {
	b.host = host
	return b
}

// SetPort sets the port of the URL
//
// Example:
//
//	builder.SetPort("8080")
func (b *URLBuilder) SetPort(port string) *URLBuilder {
	b.port = port
	return b
}

// SetCredentials sets the username and password for basic auth
//
// Example:
//
//	builder.SetCredentials("username", "password")
func (b *URLBuilder) SetCredentials(username, password string) *URLBuilder {
	b.username = username
	b.password = password
	return b
}

// AddPath adds a path segment to the URL
//
// Example:
//
//	builder.AddPath("v1").AddPath("users")
func (b *URLBuilder) AddPath(segment string) *URLBuilder {
	b.path = strings.TrimRight(b.path, "/") + "/" + strings.Trim(segment, "/")
	return b
}

// SetPath sets the complete path, overwriting any existing path
//
// Example:
//
//	builder.SetPath("/v1/users")
func (b *URLBuilder) SetPath(path string) *URLBuilder {
	b.path = path
	return b
}

// AddQuery adds a query parameter to the URL
//
// Example:
//
//	builder.AddQuery("page", "1").AddQuery("limit", "10")
func (b *URLBuilder) AddQuery(key, value string) *URLBuilder {
	b.query.Add(key, value)
	return b
}

// SetQuery sets a query parameter, overwriting any existing values for the key
//
// Example:
//
//	builder.SetQuery("page", "1")
func (b *URLBuilder) SetQuery(key, value string) *URLBuilder {
	b.query.Set(key, value)
	return b
}

// RemoveQuery removes a query parameter from the URL
//
// Example:
//
//	builder.RemoveQuery("page")
func (b *URLBuilder) RemoveQuery(key string) *URLBuilder {
	b.query.Del(key)
	return b
}

// SetFragment sets the fragment (hash) of the URL
//
// Example:
//
//	builder.SetFragment("section1")
func (b *URLBuilder) SetFragment(fragment string) *URLBuilder {
	b.fragment = fragment
	return b
}

// Build constructs and returns the final URL as a string
//
// Example:
//
//	url := builder.Build()
//	fmt.Println(url)  // https://api.example.com/v1/users?page=1&limit=10#section1
func (b *URLBuilder) Build() string {
	var sb strings.Builder

	if b.scheme != "" {
		sb.WriteString(b.scheme)
		sb.WriteString("://")
	}

	if b.username != "" {
		sb.WriteString(url.UserPassword(b.username, b.password).String())
		sb.WriteString("@")
	}

	sb.WriteString(b.host)

	if b.port != "" {
		sb.WriteString(":")
		sb.WriteString(b.port)
	}

	if b.path != "" && b.path != "/" {
		if !strings.HasPrefix(b.path, "/") {
			sb.WriteString("/")
		}
		sb.WriteString(b.path)
	}

	if len(b.query) > 0 {
		sb.WriteString("?")
		sb.WriteString(b.query.Encode())
	}

	if b.fragment != "" {
		sb.WriteString("#")
		sb.WriteString(b.fragment)
	}

	return sb.String()
}

// BuildURL constructs and returns the final URL as a *url.URL
//
// Example:
//
//	urlObj, err := builder.BuildURL()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(urlObj.String())
func (b *URLBuilder) BuildURL() (*url.URL, error) {
	return url.Parse(b.Build())
}

// Clone creates a deep copy of the URLBuilder
//
// Example:
//
//	newBuilder := builder.Clone()
func (b *URLBuilder) Clone() *URLBuilder {
	newBuilder := &URLBuilder{
		scheme:   b.scheme,
		username: b.username,
		password: b.password,
		host:     b.host,
		port:     b.port,
		path:     b.path,
		query:    make(url.Values),
		fragment: b.fragment,
	}
	for k, v := range b.query {
		newBuilder.query[k] = v
	}
	return newBuilder
}

// ParseURL parses a URL string and returns a new URLBuilder
//
// Example:
//
//	builder, err := ParseURL("https://api.example.com/v1/users?page=1")
//	if err != nil {
//	    log.Fatal(err)
//	}
func ParseURL(rawURL string) (*URLBuilder, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	builder := &URLBuilder{
		scheme:   parsed.Scheme,
		host:     parsed.Hostname(),
		port:     parsed.Port(),
		path:     parsed.Path,
		query:    parsed.Query(),
		fragment: parsed.Fragment,
	}

	if parsed.User != nil {
		builder.username = parsed.User.Username()
		builder.password, _ = parsed.User.Password()
	}

	return builder, nil
}
