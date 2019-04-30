package auth

import (
	"fmt"
	"github.com/bassbeaver/gkernel"
	"github.com/bassbeaver/gkernel/event_bus/event"
	"github.com/bassbeaver/gkernel/response"
	"github.com/kataras/go-sessions"
	loggerService "gkernel-skeleton/service/request_logger"
	sessionService "gkernel-skeleton/service/session"
	"net/http"
)

const (
	MiddlewareServiceAlias = "AuthService"
	UserSessionKey         = "auth::userId"
	RequestContextKey      = "Auth::user"
	flashMessageErrorKey   = "Auth::error"
)

type UserInterface interface {
	GetId() string
}

type UserProvider interface {
	LoadById(string) UserInterface
	LoadByLogPass(string, string) UserInterface
}

type AuthenticationError struct {
	Login   string
	Message string
}

//--------------------

type Service struct {
	userProvider                      UserProvider
	loginPageUrl                      string
	ifAlreadyAuthenticatedFallbackUrl string
}

func (s *Service) AuthenticateBySession(eventObj *event.RequestReceived) {
	userId, user := s.loadUserBySession(sessionService.GetFromRequestContext(eventObj.GetRequest()))
	if nil == user {
		if "" != userId {
			logger := loggerService.GetFromRequestContext(eventObj.GetRequest())
			logger.Critical("Failed to authenticate user by session. User "+userId+" not found", nil)

			s.LogOut(eventObj)
		}

		return
	}

	eventObj.RequestContextAppend(RequestContextKey, user)
}

func (s *Service) AuthenticateByLogPass(eventObj *event.RequestReceived) {
	sessionObj := sessionService.GetFromRequestContext(eventObj.GetRequest())

	login := eventObj.GetRequest().FormValue("login[login]")
	password := eventObj.GetRequest().FormValue("login[password]")

	user := s.userProvider.LoadByLogPass(login, password)
	if nil == user {
		sessionObj.SetFlash(
			flashMessageErrorKey,
			&AuthenticationError{
				Login:   login,
				Message: "Invalid login/password",
			},
		)

		eventObj.SetResponse(
			response.NewRedirectResponse(eventObj.GetRequest(), s.loginPageUrl, http.StatusSeeOther),
		)

		return
	}

	sessionObj.Set(UserSessionKey, user.GetId())

	eventObj.RequestContextAppend(RequestContextKey, user)

	eventObj.SetResponse(
		response.NewRedirectResponse(eventObj.GetRequest(), s.ifAlreadyAuthenticatedFallbackUrl, http.StatusSeeOther),
	)
}

func (s *Service) RedirectIfAuthenticated(eventObj *event.RequestReceived) {
	if GetUser(eventObj.GetRequest()) != nil {
		eventObj.SetResponse(
			response.NewRedirectResponse(eventObj.GetRequest(), s.ifAlreadyAuthenticatedFallbackUrl, http.StatusSeeOther),
		)
	}
}

func (s *Service) RedirectToLoginIfNotAuthenticated(eventObj *event.RequestReceived) {
	if GetUser(eventObj.GetRequest()) == nil {
		eventObj.SetResponse(
			response.NewRedirectResponse(eventObj.GetRequest(), s.loginPageUrl, http.StatusSeeOther),
		)
	}
}

func (s *Service) LogOut(eventObj *event.RequestReceived) {
	sessionObj := sessionService.GetFromRequestContext(eventObj.GetRequest())
	sessionObj.Clear()
	eventObj.SetResponse(
		response.NewRedirectResponse(eventObj.GetRequest(), s.loginPageUrl, http.StatusSeeOther),
	)
}

func (s *Service) loadUserBySession(sessionObj *sessions.Session) (string, UserInterface) {
	userId := sessionObj.GetString(UserSessionKey)
	if "" == userId {
		return "", nil
	}

	user := s.userProvider.LoadById(userId)

	return userId, user
}

//--------------------

func newAuthService(userLoader UserProvider, loginPageUrl, ifAlreadyAuthenticatedFallbackUrl string) *Service {
	return &Service{
		userProvider:                      userLoader,
		loginPageUrl:                      loginPageUrl,
		ifAlreadyAuthenticatedFallbackUrl: ifAlreadyAuthenticatedFallbackUrl,
	}
}

func Register(kernelObj *gkernel.Kernel) {
	err := kernelObj.RegisterService(MiddlewareServiceAlias, newAuthService, true)
	if nil != err {
		panic(fmt.Sprintf("failed to register %s service, error: %s", MiddlewareServiceAlias, err.Error()))
	}
}

//--------------------

func GetAuthenticationError(request *http.Request) *AuthenticationError {
	sessionObj := sessionService.GetFromRequestContext(request)

	err := sessionObj.GetFlash(flashMessageErrorKey)
	if nil == err {
		return nil
	}

	return err.(*AuthenticationError)
}

func GetUser(request *http.Request) UserInterface {
	if userObj := request.Context().Value(RequestContextKey); userObj != nil {
		return userObj.(UserInterface)
	}

	return nil
}
