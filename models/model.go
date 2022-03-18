package models

import "github.com/scmn-dev/tran/ios"

type TranOptions struct {
	TranxAddress string
	TranxPort    int
	Auth         AuthLogin
}

type AuthLogin struct {
	Token    string
	Hostname string
	IO 	     *ios.IOStreams
}

type Password string
