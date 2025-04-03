package consts

const (
	Port       string = ":8080"
	Expiry     int64  = 1 // in minute
	Period     int    = 5 // in minute
	HeaderKey  string = "Authorization"
	EmailRegex string = `^\w+[-.+]?\w+@[\w]+([-.]\w+)*\.[a-zA-Z]{2,}$`
)
