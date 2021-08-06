package errcode

var (
	Success                   = NewError(0, "成功")
	ServerError               = NewError(10000000, "服务内部错误")
	InvalidParams             = NewError(10000001, "入参错误")
	NotFound                  = NewError(10000002, "找不到")
	UnauthorizedAuthNotExist  = NewError(10000003, "鉴权失败，找不到对应的 AppKey 和 AppSecret")
	UnauthorizedTokenError    = NewError(10000004, "鉴权失败，Token 错误")
	UnauthorizedTokenTimeout  = NewError(10000005, "鉴权失败，Token 超时")
	UnauthorizedTokenGenerate = NewError(10000006, "鉴权失败，Token 生成失败")
	TooManyRequests           = NewError(10000007, "请求过多")
)

var (
	ErrorGetTagListFail = NewError(20010001, "获取标签列表失败")
	ErrorCreateTagFail  = NewError(20010002, "创建标签失败")
	ErrorUpdateTagFail  = NewError(20010003, "更新标签失败")
	ErrorDeleteTagFail  = NewError(20010004, "删除标签失败")
	ErrorCountTagFail   = NewError(20010005, "统计标签失败")
)

var (
	ErrorGetArticleListFail = NewError(30010001, "获取文章列表失败")
	ErrorCreateArticleFail  = NewError(30010002, "创建文章失败")
	ErrorUpdateArticleFail  = NewError(30010003, "更新文章失败")
	ErrorDeleteArticleFail  = NewError(30010004, "删除文章失败")
	ErrorCountArticleFail   = NewError(30010005, "统计文章失败")
	ErrorGetArticleFail = NewError(30010006, "获取文章失败")
)
