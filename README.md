# govalidator (Flexible validator for GoLang)

[![GoReportCard](https://goreportcard.com/badge/github.com/viant/godiff)](https://goreportcard.com/report/github.com/viant/godiff)
[![GoDoc](https://godoc.org/github.com/viant/godiff?status.svg)](https://godoc.org/github.com/viant/godiff)

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
- rgb
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
- gt(N)
- gte(N)


### Additional tag
- omitempty
- skipPath - remove path from location
- presence - presence field


## Validation option 
- WithShallow  - shallow check
- WithPresence - check only field set


## Contributing to govalidator

govalidator is an open source project and contributors are welcome!

See [TODO](TODO.md) list

## Credits and Acknowledgements

**Library Author:** Adrian Witas

