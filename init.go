package govalidator

func init() {

	Register("required", newRequiredCheck)
	Register("alpha", NewRegExprCheck(alphaRegex))
	Register("alphanum", NewRegExprCheck(alphaNumericRegex))
	Register("alphaunicode", NewRegExprCheck(alphaUnicodeRegex))
	Register("alphaUnicodeNumeric", NewRegExprCheck(alphaUnicodeNumericRegex))
	Register("numeric", NewRegExprCheck(numericRegex))
	Register("number", NewRegExprCheck(numberRegex))
	Register("hexadecimal", NewRegExprCheck(hexadecimalRegex))
	Register("hexColor", NewRegExprCheck(hexColorRegex))
	Register("rgb", NewRegExprCheck(rgbRegex))
	Register("rgb", NewRegExprCheck(rgbaRegex))
	Register("hsl", NewRegExprCheck(hslRegex))
	Register("hsla", NewRegExprCheck(hslaRegex))
	Register("e164", NewRegExprCheck(e164Regex))
	Register("localPhone", NewRegExprCheck(localPhoneRegex))
	Register("email", NewRegExprCheck(emailRegex))
	Register("base64", NewRegExprCheck(base64Regex))
	Register("base64URL", NewRegExprCheck(base64URLRegex))
	Register("iSBN10", NewRegExprCheck(iSBN10Regex))
	Register("iSBN13", NewRegExprCheck(iSBN13Regex))
	Register("uUID3", NewRegExprCheck(uUID3Regex))
	Register("uUID4", NewRegExprCheck(uUID4Regex))
	Register("uUID5", NewRegExprCheck(uUID5Regex))
	Register("uUID", NewRegExprCheck(uUIDRegex))
	Register("uUID3RFC4122", NewRegExprCheck(uUID3RFC4122Regex))
	Register("uUID4RFC4122", NewRegExprCheck(uUID4RFC4122Regex))
	Register("uUID5RFC4122", NewRegExprCheck(uUID5RFC4122Regex))
	Register("uUIDRFC4122", NewRegExprCheck(uUIDRFC4122Regex))

	Register("uLID", NewRegExprCheck(uLIDRegex))

	Register("md4", NewRegExprCheck(md4Regex))
	Register("md5", NewRegExprCheck(md5Regex))
	Register("sha256", NewRegExprCheck(sha256Regex))
	Register("sha384", NewRegExprCheck(sha384Regex))
	Register("sha512", NewRegExprCheck(sha512Regex))

	Register("ripemd128", NewRegExprCheck(ripemd128Regex))
	Register("ripemd160", NewRegExprCheck(ripemd160Regex))
	Register("tiger128", NewRegExprCheck(tiger128Regex))
	Register("tiger160", NewRegExprCheck(tiger160Regex))
	Register("tiger192", NewRegExprCheck(tiger192Regex))

	Register("aSCII", NewRegExprCheck(aSCIIRegex))
	Register("printableASCII", NewRegExprCheck(printableASCIIRegex))
	Register("multibyte", NewRegExprCheck(multibyteRegex))
	Register("dataURI", NewRegExprCheck(dataURIRegex))

	Register("latitude", NewRegExprCheck(latitudeRegex))
	Register("longitude", NewRegExprCheck(longitudeRegex))
	Register("ssn", NewRegExprCheck(sSNRegex))

	Register("hostnameRFC952", NewRegExprCheck(hostnameRFC952Regex))
	Register("hostnameRFC1123", NewRegExprCheck(hostnameRFC1123Regex))
	Register("fqdnRFC1123", NewRegExprCheck(fqdnRFC1123Regex))
	Register("btcAddress", NewRegExprCheck(btcAddressRegex))
	Register("uRLEncoded", NewRegExprCheck(uRLEncodedRegex))

	Register("hTMLEncoded", NewRegExprCheck(hTMLEncodedRegex))
	Register("hTML", NewRegExprCheck(hTMLRegex))
	Register("jWT", NewRegExprCheck(jWTRegex))
	Register("bic", NewRegExprCheck(bicRegex))
	Register("dnsRegexRFC1035Label", NewRegExprCheck(dnsRegexRFC1035Label))
	Register("iabcategory", NewRegExprCheck(iabCategory))
	Register("iabcategories", NewRepeatedRegExprCheck(iabCategory, ","))
	Register("domain", NewRegExprCheck(domain))
	Register("wwwdomain", NewRegExprCheck(worldWideWebDomain))
	Register("nonwwwdomain", NewNotRegExprCheck(worldWideWebDomain))
	Register("topdomain", NewNotRegExprCheck(worldWideWebDomain))
	Register("gt", NewGt())
	Register("lt", NewLt())
	Register("ge", NewGte())
	Register("le", NewLte())
	Register("choice", NewChoice())
	RegisterAlias("phone", "e164", "localPhone")
}
