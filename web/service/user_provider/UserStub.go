package user_provider

type UserStub struct {
	Id    string
	Login string
}

func (us *UserStub) GetId() string {
	return us.Id
}
