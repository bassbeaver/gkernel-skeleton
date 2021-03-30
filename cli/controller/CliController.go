package controller

import (
	"fmt"
	cliKernel "github.com/bassbeaver/gkernel/cli"
	cliKernelError "github.com/bassbeaver/gkernel/cli/error"
)

const (
	cliControllerServiceAlias = "CliController"
)

type CliController struct {
}

func (c *CliController) Command1(args []string) cliKernelError.CliError {
	fmt.Printf("First command called with arguments: %+v\n", args)

	return nil
}

//--------------------

func newCliController() *CliController {
	return &CliController{}
}

func RegisterCli(kernelObj *cliKernel.Kernel) {
	err := kernelObj.RegisterService(cliControllerServiceAlias, newCliController, true)
	if nil != err {
		panic(fmt.Sprintf("failed to register "+cliControllerServiceAlias+" service, error: %s", err.Error()))
	}
}
