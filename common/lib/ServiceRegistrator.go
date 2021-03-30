package lib

type ServiceRegistrator interface {
	RegisterService(alias string, factoryMethod interface{}, enableCaching bool) error
}
