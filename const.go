package dexcomshare

// All information for the API can be found in the link below
// https://gist.github.com/StephenBlackWasAlreadyTaken/adb0525344bedade1e25

// Application ID foud by a previous reverse engineering of the share app done by the user
// who wrote the above docs
const ApplicationId = "d8665ade-9673-4e27-9ff6-92db4ce13d13"

// API Endpoints
const BaseUrlUS = "https://share2.dexcom.com/ShareWebServices/Services"
const LoginPath = "/General/LoginPublisherAccountByName"
const AuthPath = "/General/AuthenticatePublisherAccount"
const CurrentEGVPath = "/ShareWebServices/Services/Publisher/ReadPublisherLatestGlucoseValues"
