package api

type Error struct {
	code int
	msg  string
}

func (e Error) Error() string {
	return e.msg
}

var (
	ErrorExt = Error{
		code: 1,
		msg:  "后缀名不符合要求",
	}
	ErrorSizeLimit = Error{
		code: 2,
		msg:  "文件过大",
	}
	ErrorCountLimit = Error{
		code: 3,
		msg:  "文件数量超过要求",
	}
	ErrorUnexpectedRequest = Error{
		code: 400,
		msg:  "无法理解的请求",
	}
	ErrorNoImp = Error{
		code: 502,
		msg:  "功能尚未实现",
	}
	ErrorNoRes = Error{
		code: 0,
	}
)
