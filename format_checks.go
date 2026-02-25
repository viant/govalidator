package govalidator

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	uuidv7Regex = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-7[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$")
	slugRegex   = regexp.MustCompile("^[a-z0-9]+(?:-[a-z0-9]+)*$")
	semverRegex = regexp.MustCompile("^(0|[1-9]\\d*)\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)(?:-(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*)?(?:\\+[0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*)?$")
)

func NewURL() func(field *Field, check *Check) (IsValid, error) {
	return newFormatCheck("url", func(actual string) bool {
		u, err := url.Parse(actual)
		if err != nil {
			return false
		}
		return u.Scheme != "" && u.Host != ""
	})
}

func NewURI() func(field *Field, check *Check) (IsValid, error) {
	return newFormatCheck("uri", func(actual string) bool {
		u, err := url.ParseRequestURI(actual)
		if err != nil {
			return false
		}
		return u.Scheme != ""
	})
}

func NewHTTPURL() func(field *Field, check *Check) (IsValid, error) {
	return newFormatCheck("http_url", func(actual string) bool {
		u, err := url.Parse(actual)
		if err != nil {
			return false
		}
		if u.Host == "" {
			return false
		}
		scheme := strings.ToLower(u.Scheme)
		return scheme == "http" || scheme == "https"
	})
}

func NewIP() func(field *Field, check *Check) (IsValid, error) {
	return newFormatCheck("ip", func(actual string) bool {
		return net.ParseIP(actual) != nil
	})
}

func NewIPv4() func(field *Field, check *Check) (IsValid, error) {
	return newFormatCheck("ipv4", func(actual string) bool {
		ip := net.ParseIP(actual)
		return ip != nil && ip.To4() != nil
	})
}

func NewIPv6() func(field *Field, check *Check) (IsValid, error) {
	return newFormatCheck("ipv6", func(actual string) bool {
		ip := net.ParseIP(actual)
		return ip != nil && ip.To4() == nil
	})
}

func NewCIDR() func(field *Field, check *Check) (IsValid, error) {
	return newFormatCheck("cidr", func(actual string) bool {
		_, _, err := net.ParseCIDR(actual)
		return err == nil
	})
}

func NewHostname() func(field *Field, check *Check) (IsValid, error) {
	return newFormatCheck("hostname", func(actual string) bool {
		return hostnameRFC1123Regex.MatchString(actual)
	})
}

func NewMAC() func(field *Field, check *Check) (IsValid, error) {
	return newFormatCheck("mac", func(actual string) bool {
		_, err := net.ParseMAC(actual)
		return err == nil
	})
}

func NewPort() func(field *Field, check *Check) (IsValid, error) {
	return func(field *Field, check *Check) (IsValid, error) {
		kind, elemKind := typeKinds(field)
		switch kind {
		case reflect.String:
			return func(ctx context.Context, value interface{}) (bool, error) {
				actual, ok := asStringValue(value)
				if !ok {
					return false, nil
				}
				return isValidPort(actual), nil
			}, nil
		case reflect.Slice:
			if elemKind == reflect.String {
				return func(ctx context.Context, value interface{}) (bool, error) {
					actual, ok := asStringValue(value)
					if !ok {
						return false, nil
					}
					return isValidPort(actual), nil
				}, nil
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return func(ctx context.Context, value interface{}) (bool, error) {
				numeric, ok := numericValue(value)
				if !ok {
					return false, nil
				}
				return numeric >= 1 && numeric <= 65535 && numeric == float64(int64(numeric)), nil
			}, nil
		}
		return nil, fmt.Errorf("unsupported port type: %s", field.Type.String())
	}
}

func NewUUIDv7() func(field *Field, check *Check) (IsValid, error) {
	return newFormatCheck("uuidv7", func(actual string) bool {
		return uuidv7Regex.MatchString(actual)
	})
}

func NewSlug() func(field *Field, check *Check) (IsValid, error) {
	return newFormatCheck("slug", func(actual string) bool {
		return slugRegex.MatchString(actual)
	})
}

func NewSemver() func(field *Field, check *Check) (IsValid, error) {
	return newFormatCheck("semver", func(actual string) bool {
		return semverRegex.MatchString(actual)
	})
}

func NewJSON() func(field *Field, check *Check) (IsValid, error) {
	return func(field *Field, check *Check) (IsValid, error) {
		kind, elemKind := typeKinds(field)
		switch kind {
		case reflect.String:
			return func(ctx context.Context, value interface{}) (bool, error) {
				actual, ok := asStringValue(value)
				if !ok {
					return false, nil
				}
				return json.Valid([]byte(actual)), nil
			}, nil
		case reflect.Slice:
			if elemKind == reflect.String {
				return func(ctx context.Context, value interface{}) (bool, error) {
					actual, ok := asStringValue(value)
					if !ok {
						return false, nil
					}
					return json.Valid([]byte(actual)), nil
				}, nil
			}
			if elemKind == reflect.Uint8 {
				return func(ctx context.Context, value interface{}) (bool, error) {
					bytes, ok := value.([]byte)
					if !ok {
						return false, nil
					}
					return json.Valid(bytes), nil
				}, nil
			}
		}
		return nil, fmt.Errorf("unsupported json type: %s", field.Type.String())
	}
}

func newFormatCheck(name string, predicate func(actual string) bool) func(field *Field, check *Check) (IsValid, error) {
	return func(field *Field, check *Check) (IsValid, error) {
		kind, elemKind := typeKinds(field)
		switch kind {
		case reflect.String:
		case reflect.Slice:
			if elemKind != reflect.String {
				return nil, fmt.Errorf("unsupported %s type: %s", name, field.Type.String())
			}
		default:
			return nil, fmt.Errorf("unsupported %s type: %s", name, field.Type.String())
		}
		return func(ctx context.Context, value interface{}) (bool, error) {
			actual, ok := asStringValue(value)
			if !ok {
				return false, nil
			}
			return predicate(actual), nil
		}, nil
	}
}

func isValidPort(actual string) bool {
	numeric, err := strconv.Atoi(actual)
	if err != nil {
		return false
	}
	return numeric >= 1 && numeric <= 65535
}
