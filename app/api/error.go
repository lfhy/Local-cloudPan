package api

type Error struct {
	Code int
	Msg  string
}

func (e Error) Error() string {
	return e.Msg
}

var (
	ErrorExt = Error{
		Code: 1,
		Msg:  "后缀名不符合要求",
	}
	ErrorSizeLimit = Error{
		Code: 2,
		Msg:  "文件过大",
	}
	ErrorCountLimit = Error{
		Code: 3,
		Msg:  "文件数量超过要求",
	}
	ErrorUnexpectedRequest = Error{
		Code: 400,
		Msg:  "无法理解的请求",
	}
	ErrorNoImp = Error{
		Code: 502,
		Msg:  "功能尚未实现",
	}
	ErrorNoRes = Error{
		Code: 0,
	}
	ErrorCheckFileFailed = Error{
		Code: 201,
		Msg:  "部分文件不存在",
	}
	ErrorRenameFailed = Error{
		Code: 201,
		Msg:  "重命名目标已存在",
	}
)
