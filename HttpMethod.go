package unirest

type HttpMethod int

const (
	GET HttpMethod = 1 + iota
	POST
	PUT
	PATCH
	DELETE
)

func (method HttpMethod) ToString() string {
	switch method {
	case POST:
		return "POST"

	case PUT:
		return "PUT"

	case PATCH:
		return "PATCH"

	case DELETE:
		return "DELETE"

	default:
		return "GET"
	}
}
