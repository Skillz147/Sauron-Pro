package tls

// Subdomains for certificate coverage
var BaseSubdomains = []string{
	"login", "live", "outlook", "token",
	"live.login", "logincdn", "secure", "aadcdn", "aadcdn.msftauth",
	"login.microsoftonline", "outlook.microsoftonline", "token.microsoftonline",
	"live.login.microsoftonline", "logincdn.msauth", "secure.aadcdn",
	"aadcdn.msauth", "aadcdn.msftauth", "account.live",
	"login.microsoft", "outlook.microsoft", "token.microsoft",
	"live.login.microsoft", "logincdn.microsoft", "secure.microsoft",
	"aadcdn.microsoft", "aadcdn.msftauth.microsoft", "account.microsoft",
	"outlook.office",
}
