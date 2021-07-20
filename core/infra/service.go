package infra

type Service interface {
	Execute(Request) Response
}
