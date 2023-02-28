package govalidator

/*

This file originated from https://github.com/go-playground/validator/blob/master/regexes.go

The MIT License (MIT)

Copyright (c) 2015 Dean Karn

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
*/
import "regexp"

const (
	alphaRegexPattern               = "^[a-zA-Z]+$"
	alphaNumericRegexPattern        = "^[a-zA-Z0-9]+$"
	alphaUnicodeRegexPattern        = "^[\\p{L}]+$"
	alphaUnicodeNumericRegexPattern = "^[\\p{L}\\p{N}]+$"
	numericRegexPattern             = "^[-+]?[0-9]+(?:\\.[0-9]+)?$"
	numberRegexPattern              = "^[0-9]+$"
	hexadecimalRegexPattern         = "^(0[xX])?[0-9a-fA-F]+$"
	hexColorRegexPattern            = "^#(?:[0-9a-fA-F]{3}|[0-9a-fA-F]{4}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$"
	rgbRegexPattern                 = "^rgb\\(\\s*(?:(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])|(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%)\\s*\\)$"
	rgbaRegexPattern                = "^rgba\\(\\s*(?:(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])|(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%)\\s*,\\s*(?:(?:0.[1-9]*)|[01])\\s*\\)$"
	hslRegexPattern                 = "^hsl\\(\\s*(?:0|[1-9]\\d?|[12]\\d\\d|3[0-5]\\d|360)\\s*,\\s*(?:(?:0|[1-9]\\d?|100)%)\\s*,\\s*(?:(?:0|[1-9]\\d?|100)%)\\s*\\)$"
	hslaRegexPattern                = "^hsla\\(\\s*(?:0|[1-9]\\d?|[12]\\d\\d|3[0-5]\\d|360)\\s*,\\s*(?:(?:0|[1-9]\\d?|100)%)\\s*,\\s*(?:(?:0|[1-9]\\d?|100)%)\\s*,\\s*(?:(?:0.[1-9]*)|[01])\\s*\\)$"
	emailRegexPattern               = "^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	e164RegexPattern                = "^\\+[1-9]?[0-9]{7,14}$"
	localPhoneRegexPattern          = `^\(?\d{3}\)?[\s.-]\d{3}[\s.-]\d{4}$`
	base64RegexPattern              = "^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$"
	base64URLRegexPattern           = "^(?:[A-Za-z0-9-_]{4})*(?:[A-Za-z0-9-_]{2}==|[A-Za-z0-9-_]{3}=|[A-Za-z0-9-_]{4})$"
	iSBN10RegexPattern              = "^(?:[0-9]{9}X|[0-9]{10})$"
	iSBN13RegexPattern              = "^(?:(?:97(?:8|9))[0-9]{10})$"
	uUID3RegexPattern               = "^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$"
	uUID4RegexPattern               = "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	uUID5RegexPattern               = "^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	uUIDRegexPattern                = "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
	uUID3RFC4122RegexPattern        = "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-3[0-9a-fA-F]{3}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"
	uUID4RFC4122RegexPattern        = "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$"
	uUID5RFC4122RegexPattern        = "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-5[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$"
	uUIDRFC4122RegexPattern         = "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"
	uLIDRegexPattern                = "^[A-HJKMNP-TV-Z0-9]{26}$"
	md4RegexPattern                 = "^[0-9a-f]{32}$"
	md5RegexPattern                 = "^[0-9a-f]{32}$"
	sha256RegexPattern              = "^[0-9a-f]{64}$"
	sha384RegexPattern              = "^[0-9a-f]{96}$"
	sha512RegexPattern              = "^[0-9a-f]{128}$"
	ripemd128RegexPattern           = "^[0-9a-f]{32}$"
	ripemd160RegexPattern           = "^[0-9a-f]{40}$"
	tiger128RegexPattern            = "^[0-9a-f]{32}$"
	tiger160RegexPattern            = "^[0-9a-f]{40}$"
	tiger192RegexPattern            = "^[0-9a-f]{48}$"
	aSCIIRegexPattern               = "^[\x00-\x7F]*$"
	printableASCIIRegexPattern      = "^[\x20-\x7E]*$"
	multibyteRegexPattern           = "[^\x00-\x7F]"
	dataURIRegexPattern             = `^data:((?:\w+\/(?:([^;]|;[^;]).)+)?)`
	latitudeRegexPattern            = "^[-+]?([1-8]?\\d(\\.\\d+)?|90(\\.0+)?)$"
	longitudeRegexPattern           = "^[-+]?(180(\\.0+)?|((1[0-7]\\d)|([1-9]?\\d))(\\.\\d+)?)$"
	sSNRegexPattern                 = `^[0-9]{3}[ -]?(0[1-9]|[1-9][0-9])[ -]?([1-9][0-9]{3}|[0-9][1-9][0-9]{2}|[0-9]{2}[1-9][0-9]|[0-9]{3}[1-9])$`
	hostnameRFC952RegexPattern      = `^[a-zA-Z]([a-zA-Z0-9\-]+[\.]?)*[a-zA-Z0-9]$`                                                                   // https://tools.ietf.org/html/rfc952
	hostnameRFC1123RegexPattern     = `^([a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,62}){1}(\.[a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,62})*?$`                                 // accepts hostname starting with a digit https://tools.ietf.org/html/rfc1123
	fqdnRFC1123RegexPattern         = `^([a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,62})(\.[a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,62})*?(\.[a-zA-Z]{1}[a-zA-Z0-9]{0,62})\.?$` // same as hostnameRFC1123RegexPattern but must contain a non numerical TLD (possibly ending with '.')
	btcAddressRegexPattern          = `^[13][a-km-zA-HJ-NP-Z1-9]{25,34}$`                                                                             // bitcoin address
	uRLEncodedRegexPattern          = `^(?:[^%]|%[0-9A-Fa-f]{2})*$`
	hTMLEncodedRegexPattern         = `&#[x]?([0-9a-fA-F]{2})|(&gt)|(&lt)|(&quot)|(&amp)+[;]?`
	hTMLRegexPattern                = `<[/]?([a-zA-Z]+).*?>`
	jWTRegexPattern                 = "^[A-Za-z0-9-_]+\\.[A-Za-z0-9-_]+\\.[A-Za-z0-9-_]*$"
	bicRegexPattern                 = `^[A-Za-z]{6}[A-Za-z0-9]{2}([A-Za-z0-9]{3})?$`
	dnsRegexPatternRFC1035Label     = "^[a-z]([-a-z0-9]*[a-z0-9]){0,62}$"
)

var (
	alphaRegex               = regexp.MustCompile(alphaRegexPattern)
	alphaNumericRegex        = regexp.MustCompile(alphaNumericRegexPattern)
	alphaUnicodeRegex        = regexp.MustCompile(alphaUnicodeRegexPattern)
	alphaUnicodeNumericRegex = regexp.MustCompile(alphaUnicodeNumericRegexPattern)
	numericRegex             = regexp.MustCompile(numericRegexPattern)
	numberRegex              = regexp.MustCompile(numberRegexPattern)
	hexadecimalRegex         = regexp.MustCompile(hexadecimalRegexPattern)
	hexColorRegex            = regexp.MustCompile(hexColorRegexPattern)
	rgbRegex                 = regexp.MustCompile(rgbRegexPattern)
	rgbaRegex                = regexp.MustCompile(rgbaRegexPattern)
	hslRegex                 = regexp.MustCompile(hslRegexPattern)
	hslaRegex                = regexp.MustCompile(hslaRegexPattern)
	e164Regex                = regexp.MustCompile(e164RegexPattern)
	localPhoneRegex          = regexp.MustCompile(localPhoneRegexPattern)
	emailRegex               = regexp.MustCompile(emailRegexPattern)
	base64Regex              = regexp.MustCompile(base64RegexPattern)
	base64URLRegex           = regexp.MustCompile(base64URLRegexPattern)
	iSBN10Regex              = regexp.MustCompile(iSBN10RegexPattern)
	iSBN13Regex              = regexp.MustCompile(iSBN13RegexPattern)
	uUID3Regex               = regexp.MustCompile(uUID3RegexPattern)
	uUID4Regex               = regexp.MustCompile(uUID4RegexPattern)
	uUID5Regex               = regexp.MustCompile(uUID5RegexPattern)
	uUIDRegex                = regexp.MustCompile(uUIDRegexPattern)
	uUID3RFC4122Regex        = regexp.MustCompile(uUID3RFC4122RegexPattern)
	uUID4RFC4122Regex        = regexp.MustCompile(uUID4RFC4122RegexPattern)
	uUID5RFC4122Regex        = regexp.MustCompile(uUID5RFC4122RegexPattern)
	uUIDRFC4122Regex         = regexp.MustCompile(uUIDRFC4122RegexPattern)
	uLIDRegex                = regexp.MustCompile(uLIDRegexPattern)
	md4Regex                 = regexp.MustCompile(md4RegexPattern)
	md5Regex                 = regexp.MustCompile(md5RegexPattern)
	sha256Regex              = regexp.MustCompile(sha256RegexPattern)
	sha384Regex              = regexp.MustCompile(sha384RegexPattern)
	sha512Regex              = regexp.MustCompile(sha512RegexPattern)
	ripemd128Regex           = regexp.MustCompile(ripemd128RegexPattern)
	ripemd160Regex           = regexp.MustCompile(ripemd160RegexPattern)
	tiger128Regex            = regexp.MustCompile(tiger128RegexPattern)
	tiger160Regex            = regexp.MustCompile(tiger160RegexPattern)
	tiger192Regex            = regexp.MustCompile(tiger192RegexPattern)
	aSCIIRegex               = regexp.MustCompile(aSCIIRegexPattern)
	printableASCIIRegex      = regexp.MustCompile(printableASCIIRegexPattern)
	multibyteRegex           = regexp.MustCompile(multibyteRegexPattern)
	dataURIRegex             = regexp.MustCompile(dataURIRegexPattern)
	latitudeRegex            = regexp.MustCompile(latitudeRegexPattern)
	longitudeRegex           = regexp.MustCompile(longitudeRegexPattern)
	sSNRegex                 = regexp.MustCompile(sSNRegexPattern)
	hostnameRFC952Regex      = regexp.MustCompile(hostnameRFC952RegexPattern)
	hostnameRFC1123Regex     = regexp.MustCompile(hostnameRFC1123RegexPattern)
	fqdnRFC1123Regex         = regexp.MustCompile(fqdnRFC1123RegexPattern)
	btcAddressRegex          = regexp.MustCompile(btcAddressRegexPattern)
	uRLEncodedRegex          = regexp.MustCompile(uRLEncodedRegexPattern)
	hTMLEncodedRegex         = regexp.MustCompile(hTMLEncodedRegexPattern)
	hTMLRegex                = regexp.MustCompile(hTMLRegexPattern)
	jWTRegex                 = regexp.MustCompile(jWTRegexPattern)
	bicRegex                 = regexp.MustCompile(bicRegexPattern)
	dnsRegexRFC1035Label     = regexp.MustCompile(dnsRegexPatternRFC1035Label)
)
