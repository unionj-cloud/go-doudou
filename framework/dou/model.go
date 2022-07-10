package dou

import "encoding/json"

type HttpMethod int

const (
	UNKNOWN HttpMethod = iota
	POST
	GET
	PUT
	DELETE
)

func (k *HttpMethod) StringSetter(value string) {
	switch value {
	case "UNKNOWN":
		*k = UNKNOWN
	case "GET":
		*k = GET
	case "POST":
		*k = POST
	case "PUT":
		*k = PUT
	case "DELETE":
		*k = DELETE
	default:
		*k = UNKNOWN
	}
}

func (k *HttpMethod) StringGetter() string {
	switch *k {
	case UNKNOWN:
		return "UNKNOWN"
	case GET:
		return "GET"
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	case DELETE:
		return "DELETE"
	default:
		return "UNKNOWN"
	}
}

func (k *HttpMethod) UnmarshalJSON(bytes []byte) error {
	var _k string
	err := json.Unmarshal(bytes, &_k)
	if err != nil {
		return err
	}
	k.StringSetter(_k)
	return nil
}

func (k HttpMethod) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.StringGetter())
}
