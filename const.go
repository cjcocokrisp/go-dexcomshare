package dexcomshare

// All information for the API can be found in the link below
// https://gist.github.com/StephenBlackWasAlreadyTaken/adb0525344bedade1e25

// Application ID foud by a previous reverse engineering of the share app done by the user
// who wrote the above docs
const ApplicationId = "d89443d2-327c-4a6f-89e5-496bbb0317db"

// API Endpoints
const BaseUrlUS = "https://share2.dexcom.com/ShareWebServices/Services"
const BaseUrlOUS = "https://shareous1.dexcom.com/ShareWebServices/Services"
const LoginPath = "/General/LoginPublisherAccountByName"
const AuthPath = "/General/AuthenticatePublisherAccount"
const CurrentEGVPath = "/Publisher/ReadPublisherLatestGlucoseValues"

// TrendArrow Map
// To check the possible states I used documentation from the pydexcom module
// https://github.com/gagebenne/pydexcom/tree/main
// The ones that are unconfirmed are noted on there as well
// I'm not getting my blood sugar purposely that high to test this
var TrendArrowMap = map[string]string{
	"None":           "", // Unconfirmed
	"Flat":           "→",
	"SingleUp":       "↑",
	"DoubleUp":       "↑↑",
	"FortyFiveUp":    "↗",
	"SingleDown":     "↓",
	"DoubleDown":     "↓↓",
	"FortyFiveDown":  "↘",
	"NotComputable":  "?", // Unconfirmed
	"RateOutOfRange": "-", // Unconfirmed
}

// Default Values
const DefaultMinutes = 1440
