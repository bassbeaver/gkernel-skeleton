package user_provider

import (
	"fmt"
	"github.com/bassbeaver/gkernel"
	"gkernel-skeleton/service/auth"
)

const (
	MiddlewareServiceAlias = "UserProvider"
)

type UserProvider struct {
	userStorage []*UserStub
}

func (up *UserProvider) LoadById(userId string) auth.UserInterface {
	// Mocking of user loading from storage

	if "user1" == userId {
		return &UserStub{
			Id:    "user1",
			Login: "login1",
		}
	}

	return nil
}

func (up *UserProvider) LoadByLogPass(login string, password string) auth.UserInterface {
	// Mocking of user loading from storage

	if "login1" == login && "password1" == password {
		return &UserStub{
			Id:    "user1",
			Login: "login1",
		}
	}

	return nil
}

//--------------------

func newUserProvider() *UserProvider {
	return &UserProvider{}
}

func Register(kernelObj *gkernel.Kernel) {
	err := kernelObj.RegisterService(MiddlewareServiceAlias, newUserProvider, true)
	if nil != err {
		panic(fmt.Sprintf("failed to register %s service, error: %s", MiddlewareServiceAlias, err.Error()))
	}
}

//--------------------

type UserStub struct {
	Id    string
	Login string
}

func (us *UserStub) GetId() string {
	return us.Id
}
