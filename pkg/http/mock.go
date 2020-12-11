package http

type (
	MockHTTP struct {
	}

	MockHandlerContext struct {
		MockBind   func(interface{}) error
		MockJSON   func(int, interface{}) error
		MockString func(int, string) error
		MockGet    func(string) interface{}
		MockParam  func(string) string
	}
)

func (h MockHandlerContext) JSON(statusCode int, i interface{}) error {
	return h.MockJSON(statusCode, i)
}

func (h MockHandlerContext) Bind(i interface{}) error {
	return h.MockBind(i)
}

func (h MockHandlerContext) String(code int, message string) error {
	return h.MockString(code, message)
}

func (h MockHandlerContext) Get(key string) interface{} {
	return h.MockGet(key)
}

func (h MockHandlerContext) Param(key string) string {
	return h.MockParam(key)
}

func (h *MockHTTP) Listen(address string) error {
	return nil
}

func (h *MockHTTP) ListenMetrics(string, string) {
}

func (h *MockHTTP) AddRoute(method HTTPMethodType, path string, handler HandlerFunc) {
}
