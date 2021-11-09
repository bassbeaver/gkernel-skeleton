package controller

import (
	"fmt"
	webKernel "github.com/bassbeaver/gkernel/web"
	kernelResponse "github.com/bassbeaver/gkernel/web/response"
	"github.com/bassbeaver/logopher"
	"gkernel-skeleton/web/service/auth"
	csrfService "gkernel-skeleton/web/service/csrf"
	"gkernel-skeleton/web/service/user_provider"
	"net/http"
)

const (
	indexControllerServiceAlias = "IndexController"
)

type IndexController struct {
}

func (c *IndexController) Index(request *http.Request) kernelResponse.Response {
	var viewData struct {
		CsrfToken string
		User      auth.UserInterface
		Header    struct {
			Title string
		}
		H1      string
		Content struct {
			SomeText string
		}
	}

	viewData.CsrfToken = csrfService.GetTokenFromRequestContext(request)

	if auth.GetUser(request) != nil {
		viewData.User = auth.GetUser(request).(*user_provider.UserStub)
	}

	viewData.Header.Title = "gkernel skeleton site"
	viewData.H1 = "Hello world!"
	viewData.Content.SomeText = "This is index page of simple application based on gkernel framework."

	logger := request.Context().Value("logger").(*logopher.Logger)
	logger.Info("Hello world!", nil)

	response := kernelResponse.NewViewResponse("index/index.gohtml")
	response.SetData(viewData)

	return response
}

func (c *IndexController) PageWithParam(request *http.Request) kernelResponse.Response {
	var viewData struct {
		CsrfToken string
		User      auth.UserInterface
		Header    struct {
			Title string
		}
		H1    string
		Param string
	}

	viewData.Header.Title = "gkernel skeleton site"
	viewData.H1 = "This is page with URL parameter"
	viewData.CsrfToken = csrfService.GetTokenFromRequestContext(request)
	if auth.GetUser(request) != nil {
		viewData.User = auth.GetUser(request).(*user_provider.UserStub)
	}

	viewData.Param = request.URL.Query().Get(":parameterValue")

	response := kernelResponse.NewViewResponse("index/page-with-param.gohtml")
	response.SetData(viewData)

	return response
}

func (c *IndexController) PrivatePage(request *http.Request) kernelResponse.Response {
	var viewData struct {
		CsrfToken string
		User      auth.UserInterface
		Header    struct {
			Title string
		}
		H1 string
	}

	viewData.CsrfToken = csrfService.GetTokenFromRequestContext(request)

	if auth.GetUser(request) != nil {
		viewData.User = auth.GetUser(request).(*user_provider.UserStub)
	}

	viewData.Header.Title = "gkernel skeleton site"
	viewData.H1 = "Hello private!"

	logger := request.Context().Value("logger").(*logopher.Logger)
	logger.Info("Hello private!", nil)

	response := kernelResponse.NewViewResponse("index/private.gohtml")
	response.SetData(viewData)

	return response
}

func (c *IndexController) LoginPage(request *http.Request) kernelResponse.Response {
	var viewData struct {
		CsrfToken string
		User      auth.UserInterface
		Header    struct {
			Title string
		}
		Login        string
		ErrorMessage string
	}

	viewData.CsrfToken = csrfService.GetTokenFromRequestContext(request)

	if auth.GetUser(request) != nil {
		viewData.User = auth.GetUser(request).(*user_provider.UserStub)
	}

	viewData.Header.Title = "Login page"

	authError := auth.GetAuthenticationError(request)
	if nil != authError {
		viewData.Login = authError.Login
		viewData.ErrorMessage = authError.Message
	}

	response := kernelResponse.NewViewResponse("index/login.gohtml")
	response.SetData(viewData)

	return response
}

// Technical action that should not be reached by request, as it should be intercepted with auth middleware
func (c *IndexController) PerformLoginLogout(_ *http.Request) kernelResponse.Response {
	response := kernelResponse.NewBytesResponse()
	response.SetHttpStatus(http.StatusBadRequest)
	response.Body.Write([]byte("You should not see this. Run for your life."))

	return response
}

//--------------------

func newIndexController() *IndexController {
	return &IndexController{}
}

func RegisterIndex(kernelObj *webKernel.Kernel) {
	err := kernelObj.RegisterService(indexControllerServiceAlias, newIndexController, true)
	if nil != err {
		panic(fmt.Sprintf("failed to register "+indexControllerServiceAlias+" service, error: %s", err.Error()))
	}
}
