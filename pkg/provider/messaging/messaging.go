package messaging

import "plataform/pkg/api"

type Message struct {
	Header  map[string]string `json:"header"`
	Content []byte            `json:"content"`
}

const OrganizationHeader = "x-org"

var SubjectBuildBooK = Subject("build_book")

func (m Message) Org() api.Organization {
	return api.Organization(m.Header[OrganizationHeader])
}

type Subject string

type MessageHandler func(m Message)
type ErrorHandler func(err error)

type Publisher interface {
	Publish(org api.Organization, s Subject, m Message, he ErrorHandler)
}

type PublisherSync interface {
	PublishSync(org api.Organization, s Subject, m Message) error
}

type Subscriber interface {
	Subscribe(s Subject, hm MessageHandler, he ErrorHandler)
}

type Messenger interface {
	PublisherSync
	Publisher
	Subscriber
}
