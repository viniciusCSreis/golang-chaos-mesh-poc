package api

type (
	Organization string
	AppName string
)

const DefaultOrganization = "default"


func (a AppName) String() string {
	return string(a)
}