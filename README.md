# govalidator (Flexible validator for GoLang)

[![GoReportCard](https://goreportcard.com/badge/github.com/viant/govalidator)](https://goreportcard.com/report/github.com/viant/govalidator)
[![GoDoc](https://godoc.org/github.com/viant/govalidator?status.svg)](https://godoc.org/github.com/viant/govalidator)

This library is compatible with Go 1.17+

Please refer to [`CHANGELOG.md`](CHANGELOG.md) if you encounter breaking changes.

- [Motivation](#motivation)
- [Usage](#usage)
- [Contribution](#contributing-to-govalidator)
- [License](#license)

## Motivation

The goal of this library is to provide flexible go validator for both regular struct and struct with presence set.
Presence set if a special part of go struct holding flags about what struct field has been set.
For example to use go struct for patch operation, user may only set subset of the whole struct,
thus validator has to have ability to discriminate fields that need to be validated dynamically.

Other great project for validating go struct.

- [validator](https://github.com/asaskevich/govalidator)
- [govalidator](https://github.com/asaskevich/govalidator)


## Usage

```go
validator := govalidator.New()
validation, err := validator.Validate(someStruct)
```


### The following check have been implemented

- required
- alpha
- alphanum
- alphaunicode
- alphaunicodenumeric
- numeric
- number
- hexadecimal
- hexcolor
- rgb
- rgba
- hsl
- hsla
- e164
- localphone
- phone
- email
- base64
- base64url
- isbn10
- isbn13
- uuid3
- uuid4
- uuid5
- uuid
- uuid3rfc4122
- uuid4rfc4122
- uuid5rfc4122
- uuidrfc4122
- ulid
- md4
- md5
- sha256
- sha384
- sha512
- ripemd128
- ripemd160
- tiger128
- tiger160
- tiger192
- ascii
- printableascii
- multibyte
- datauri
- latitude
- longitude
- ssn
- hostnamerfc952
- hostnamerfc1123
- fqdnrfc1123
- btcaddress
- urlencoded
- htmlencoded
- html
- jwt
- bic
- dnsregexrfc1035label
- iabCategory
- iabCategories
- domain
- wwwdomain
- nonwwwDomain
- gt(N)
- ge(N)
- gte(N)
- lt(N)
- le(N)
- lte(N)
- min(N)
- max(N)
- between(min,max)
- choice(coma separated list of allowed int or string values)
- oneof(coma separated list of allowed int or string values)
- contains(text)
- notcontains(text)
- startswith(prefix)
- endswith(suffix)
- eqfield(OtherField)
- nefield(OtherField)
- gtfield(OtherField)
- required_if(OtherField,value)
- required_unless(OtherField,value)
- required_with(OtherField)
- required_without(OtherField)
- past
- future
- url
- uri
- http_url
- ip
- ipv4
- ipv6
- cidr
- hostname
- mac
- port
- uuidv7
- slug
- semver
- json

### Validation matrix (tag -> Go kinds -> example)

| Tag | Supported Go kinds | Example |
|---|---|---|
| `required` | `string`, `bool`, `int*`, `uint*`, `float*`, `slice`, `ptr`, `time.Time` | ``Name string `validate:"required"` `` |
| `gt/ge/gte/lt/le/lte` | numbers (`int*`,`uint*`,`float*`), `string` (length), primitive slice elements (`[]int`,`[]string`) | ``Age int `validate:"gte(18),lte(65)"` `` |
| `min/max/between` | numbers, `string` (length), slice length, primitive slice elements | ``Code string `validate:"between(3,10)"` `` |
| `choice/oneof` | `string`, `int*`, pointers to them, primitive slice elements | ``State string `validate:"oneof(AZ,AK,CA)"` `` |
| `contains/notcontains/startswith/endswith` | `string`, `*string`, `[]string` elements | ``Email string `validate:"contains(@),endswith(.com)"` `` |
| `eqfield/nefield/gtfield` | compares current field to another field in same struct | ``Confirm string `validate:"eqfield(Password)"` `` |
| `required_if/required_unless` | any field type using emptiness check, based on another field value | ``Phone string `validate:"required_if(Type,mobile)"` `` |
| `required_with/required_without` | any field type using emptiness check, based on presence/absence of other fields | ``Phone string `validate:"required_without(Email)"` `` |
| `past/future` | `time.Time`, `*time.Time`, `string`, `*string` (RFC3339/RFC3339Nano/`2006-01-02`) | ``StartAt string `validate:"future"` `` |
| `url/uri/http_url` | `string`, `*string`, `[]string` elements | ``Link string `validate:"http_url"` `` |
| `ip/ipv4/ipv6/cidr` | `string`, `*string`, `[]string` elements | ``Network string `validate:"cidr"` `` |
| `hostname/mac` | `string`, `*string`, `[]string` elements | ``Host string `validate:"hostname"` `` |
| `port` | `string`, `*string`, `int*`, `uint*`, `[]string` elements | ``Port int `validate:"port"` `` |
| `uuidv7` | `string`, `*string`, `[]string` elements | ``ID string `validate:"uuidv7"` `` |
| `slug` | `string`, `*string`, `[]string` elements | ``Slug string `validate:"slug"` `` |
| `semver` | `string`, `*string`, `[]string` elements | ``Version string `validate:"semver"` `` |
| `json` | `string`, `*string`, `[]string` elements, `[]byte` | ``Payload string `validate:"json"` `` |
| Regex family (`email`, `alpha`, `domain`, `uuid4`, etc.) | `string`, `*string` | ``Email string `validate:"email"` `` |

### Additional tag
- omitempty
- skipPath - remove path from location
- presence - presence field

### Message template placeholders
- `$field` - current field name
- `$value` - current field value
- `$param` - check parameters joined with comma
- `$otherField` - related field for cross/conditional checks


## Validation option 
- WithShallow  - shallow check
- WithSetMarker - check only fields marked as present


## Contributing to govalidator

govalidator is an open source project and contributors are welcome!

See [TODO](TODO.md) list

## Credits and Acknowledgements

**Library Author:** Adrian Witas
