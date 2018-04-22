package config

const LogFile = "LogFile"
const DBConnString = "DBConnString"
const Port = "Port"
const CertKey = "CertKey"
const PrivKey = "PrivKey"
const UseSSL = "UseSSL"
const SigningKey = "SigningKey"

var defaults = map[string]string{
	LogFile:      "logs.txt",
	DBConnString: "host=localhost port=5432 user=postgres sslmode=disable dbname=postgres password=postgres123test",
	Port:         "8000",
	UseSSL:       "false",
	SigningKey:   "mytestsigningkey",
}
