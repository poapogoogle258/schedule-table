package constant

type ResponseStatus int

const (
	Success ResponseStatus = iota + 1
	DataNotFound
	UnknownError
	InvalidRequest
	Unauthorized
)

func (status ResponseStatus) GetResponseStatus() string {
	switch status {
	case Success:
		return "SUCCESS"
	case DataNotFound:
		return "DATA_NOT_FOUND"
	case UnknownError:
		return "UNKNOWN_ERROR"
	case InvalidRequest:
		return "INVALID_REQUEST"
	case Unauthorized:
		return "UNAUTHORIZED"
	default:
		panic("getResponseStatus: Not found status_code in constant_status")
	}
}

func (status ResponseStatus) GetResponseMessage() string {
	switch status {
	case Success:
		return "Success"
	case DataNotFound:
		return "Data Not Found"
	case UnknownError:
		return "Unknown Error"
	case InvalidRequest:
		return "Invalid Request"
	case Unauthorized:
		return "Unauthorized"
	default:
		panic("getResponseMessage: Not found status_code in constant_status")
	}
}

const (
	YEARLY int = iota
	MONTHLY
	WEEKLY
	DAILY
	HOURLY
	MINUTELY
	SECONDLY
)
