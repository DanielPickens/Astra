package v1

// This file contains a collection of methods that can be used from go-restful to
// generate Swagger API documentation for its models. Please read this PR for more
// information on the implementation: https://github.com/emicklei/go-restful/pull/215
//
// Tastras are ignored from the parser (e.g. Tastra(andronat):... || Tastra:...) if and only if
// they are on one line! For multiple line or blocks that you want to ignore use ---.
// Any context after a --- is ignored.
//
// Those methods can be generated by using hack/update-swagger-docs.sh

// AUTO-GENERATED FUNCTIONS START HERE
var map_ClusterRoleScopeRestriction = map[string]string{
	"":                "ClusterRoleScopeRestriction describes restrictions on cluster role scopes",
	"roleNames":       "RoleNames is the list of cluster roles that can referenced.  * means anything",
	"namespaces":      "Namespaces is the list of namespaces that can be referenced.  * means any of them (including *)",
	"allowEscalation": "AllowEscalation indicates whether you can request roles and their escalating resources",
}

func (ClusterRoleScopeRestriction) SwaggerDoc() map[string]string {
	return map_ClusterRoleScopeRestriction
}

var map_OAuthAccessToken = map[string]string{
	"":                         "OAuthAccessToken describes an OAuth access token. The name of a token must be prefixed with a `sha256~` string, must not contain \"/\" or \"%\" characters and must be at least 32 characters long.\n\nThe name of the token is constructed from the actual token by sha256-hashing it and using URL-safe unpadded base64-encoding (as described in RFC4648) on the hashed result.\n\nCompatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).",
	"clientName":               "ClientName references the client that created this token.",
	"expiresIn":                "ExpiresIn is the seconds from CreationTime before this token expires.",
	"scopes":                   "Scopes is an array of the requested scopes.",
	"redirectURI":              "RedirectURI is the redirection associated with the token.",
	"userName":                 "UserName is the user name associated with this token",
	"userUID":                  "UserUID is the unique UID associated with this token",
	"authorizeToken":           "AuthorizeToken contains the token that authorized this token",
	"refreshToken":             "RefreshToken is the value by which this token can be renewed. Can be blank.",
	"inactivityTimeoutSeconds": "InactivityTimeoutSeconds is the value in seconds, from the CreationTimestamp, after which this token can no longer be used. The value is automatically incremented when the token is used.",
}

func (OAuthAccessToken) SwaggerDoc() map[string]string {
	return map_OAuthAccessToken
}

var map_OAuthAccessTokenList = map[string]string{
	"":      "OAuthAccessTokenList is a collection of OAuth access tokens\n\nCompatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).",
	"items": "Items is the list of OAuth access tokens",
}

func (OAuthAccessTokenList) SwaggerDoc() map[string]string {
	return map_OAuthAccessTokenList
}

var map_OAuthAuthorizeToken = map[string]string{
	"":                    "OAuthAuthorizeToken describes an OAuth authorization token\n\nCompatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).",
	"clientName":          "ClientName references the client that created this token.",
	"expiresIn":           "ExpiresIn is the seconds from CreationTime before this token expires.",
	"scopes":              "Scopes is an array of the requested scopes.",
	"redirectURI":         "RedirectURI is the redirection associated with the token.",
	"state":               "State data from request",
	"userName":            "UserName is the user name associated with this token",
	"userUID":             "UserUID is the unique UID associated with this token. UserUID and UserName must both match for this token to be valid.",
	"codeChallenge":       "CodeChallenge is the optional code_challenge associated with this authorization code, as described in rfc7636",
	"codeChallengeMethod": "CodeChallengeMethod is the optional code_challenge_method associated with this authorization code, as described in rfc7636",
}

func (OAuthAuthorizeToken) SwaggerDoc() map[string]string {
	return map_OAuthAuthorizeToken
}

var map_OAuthAuthorizeTokenList = map[string]string{
	"":      "OAuthAuthorizeTokenList is a collection of OAuth authorization tokens\n\nCompatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).",
	"items": "Items is the list of OAuth authorization tokens",
}

func (OAuthAuthorizeTokenList) SwaggerDoc() map[string]string {
	return map_OAuthAuthorizeTokenList
}

var map_OAuthClient = map[string]string{
	"":                                    "OAuthClient describes an OAuth client\n\nCompatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).",
	"secret":                              "Secret is the unique secret associated with a client",
	"additionalSecrets":                   "AdditionalSecrets holds other secrets that may be used to identify the client.  This is useful for rotation and for service account token validation",
	"respondWithChallenges":               "RespondWithChallenges indicates whether the client wants authentication needed responses made in the form of challenges instead of redirects",
	"redirectURIs":                        "RedirectURIs is the valid redirection URIs associated with a client",
	"grantMethod":                         "GrantMethod is a required field which determines how to handle grants for this client. Valid grant handling methods are:\n - auto:   always approves grant requests, useful for trusted clients\n - prompt: prompts the end user for approval of grant requests, useful for third-party clients",
	"scopeRestrictions":                   "ScopeRestrictions describes which scopes this client can request.  Each requested scope is checked against each restriction.  If any restriction matches, then the scope is allowed. If no restriction matches, then the scope is denied.",
	"accessTokenMaxAgeSeconds":            "AccessTokenMaxAgeSeconds overrides the default access token max age for tokens granted to this client. 0 means no expiration.",
	"accessTokenInactivityTimeoutSeconds": "AccessTokenInactivityTimeoutSeconds overrides the default token inactivity timeout for tokens granted to this client. The value represents the maximum amount of time that can occur between consecutive uses of the token. Tokens become invalid if they are not used within this temporal window. The user will need to acquire a new token to regain access once a token times out. This value needs to be set only if the default set in configuration is not appropriate for this client. Valid values are: - 0: Tokens for this client never time out - X: Tokens time out if there is no activity for X seconds The current minimum allowed value for X is 300 (5 minutes)\n\nWARNING: existing tokens' timeout will not be affected (lowered) by changing this value",
}

func (OAuthClient) SwaggerDoc() map[string]string {
	return map_OAuthClient
}

var map_OAuthClientAuthorization = map[string]string{
	"":           "OAuthClientAuthorization describes an authorization created by an OAuth client\n\nCompatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).",
	"clientName": "ClientName references the client that created this authorization",
	"userName":   "UserName is the user name that authorized this client",
	"userUID":    "UserUID is the unique UID associated with this authorization. UserUID and UserName must both match for this authorization to be valid.",
	"scopes":     "Scopes is an array of the granted scopes.",
}

func (OAuthClientAuthorization) SwaggerDoc() map[string]string {
	return map_OAuthClientAuthorization
}

var map_OAuthClientAuthorizationList = map[string]string{
	"":      "OAuthClientAuthorizationList is a collection of OAuth client authorizations\n\nCompatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).",
	"items": "Items is the list of OAuth client authorizations",
}

func (OAuthClientAuthorizationList) SwaggerDoc() map[string]string {
	return map_OAuthClientAuthorizationList
}

var map_OAuthClientList = map[string]string{
	"":      "OAuthClientList is a collection of OAuth clients\n\nCompatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).",
	"items": "Items is the list of OAuth clients",
}

func (OAuthClientList) SwaggerDoc() map[string]string {
	return map_OAuthClientList
}

var map_OAuthRedirectReference = map[string]string{
	"":          "OAuthRedirectReference is a reference to an OAuth redirect object.\n\nCompatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).",
	"reference": "The reference to an redirect object in the current namespace.",
}

func (OAuthRedirectReference) SwaggerDoc() map[string]string {
	return map_OAuthRedirectReference
}

var map_RedirectReference = map[string]string{
	"":      "RedirectReference specifies the target in the current namespace that resolves into redirect URIs.  Only the 'Route' kind is currently allowed.",
	"group": "The group of the target that is being referred to.",
	"kind":  "The kind of the target that is being referred to.  Currently, only 'Route' is allowed.",
	"name":  "The name of the target that is being referred to. e.g. name of the Route.",
}

func (RedirectReference) SwaggerDoc() map[string]string {
	return map_RedirectReference
}

var map_ScopeRestriction = map[string]string{
	"":            "ScopeRestriction describe one restriction on scopes.  Exactly one option must be non-nil.",
	"literals":    "ExactValues means the scope has to match a particular set of strings exactly",
	"clusterRole": "ClusterRole describes a set of restrictions for cluster role scoping.",
}

func (ScopeRestriction) SwaggerDoc() map[string]string {
	return map_ScopeRestriction
}

var map_UserOAuthAccessTokenList = map[string]string{
	"": "UserOAuthAccessTokenList is a collection of access tokens issued on behalf of the requesting user\n\nCompatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).",
}

func (UserOAuthAccessTokenList) SwaggerDoc() map[string]string {
	return map_UserOAuthAccessTokenList
}

// AUTO-GENERATED FUNCTIONS END HERE
