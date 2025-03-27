package constants

// Query Parameters
const (
	ParamAppID       = "appID"
	Limit            = "limit"
	Offset           = "page"
	ParamFilterPrice = "price"
)

// Error Messages
const (
	ErrorInvalidLimit  = "Invalid limit value"
	ErrorInvalidOffset = "Invalid page value"
	ErrorInvalidAppID  = "Invalid App ID"
	ErrorAppNotFound   = "App Not Found"
	ErrorLoadingCache  = "Error loading app data into cache"
)

// Defaults
const (
	DefaultLimit  = "30"
	DefaultPage   = "1"
	DefaultOffset = "1"
)
const (
	// Query Parameter Names
	ParamAppName     = "appname"
	ParamSentiment   = "sentiment"
	ParamPolarityMin = "polarity_min"
	ParamPolarityMax = "polarity_max"

	// Default Query Values
	DefaultAppName     = "10 Best Foods for You"
	DefaultSentiment   = "Positive"
	DefaultPolarityMin = "0.5"
	DefaultPolarityMax = "1.0"

	// Error Messages
	ErrorInvalidPolarityMin = "Invalid minimum polarity value"
	ErrorInvalidPolarityMax = "Invalid maximum polarity value"
)
