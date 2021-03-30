package main

import (
	"flag"
	"fmt"
	"github.com/bassbeaver/gkernel/cli"
	cliController "gkernel-skeleton/cli/controller"
	redisService "gkernel-skeleton/common/service/redis"
	"os"
)

func main() {
	flagsSet := flag.NewFlagSet("", flag.ContinueOnError)

	configPathFlag := flagsSet.String("config", "", "Config folder path")

	helpFlag := flagsSet.Bool("help", false, "Output help")
	flagsSet.BoolVar(helpFlag, "h", false, "Output help (shorthand)")

	flagsErr := flagsSet.Parse(os.Args[1:])
	if nil != flagsErr {
		panic(flagsErr)
	}

	var configPath string
	if "" == *configPathFlag {
		curBinDir, curBinDirError := os.Getwd()
		if curBinDirError != nil {
			fmt.Println("Failed to determine working dir")
			panic(curBinDirError)
		}
		configPath = curBinDir + "/config"
	} else {
		configPath = *configPathFlag
	}

	kernelObj, kernelCreationError := cli.NewKernel(configPath)
	if nil != kernelCreationError {
		panic(kernelCreationError)
	}

	//******************************** Service's registration ********************************

	// ---- Controllers
	cliController.RegisterCli(kernelObj)

	// ---- Other services
	redisService.Register(kernelObj)

	//************** Starting applications **************

	cliArgsWithoutConfigFlag := flagsSet.Args()
	// Restoring common help flag to pass it to Kernel
	if *helpFlag {
		cliArgsWithoutConfigFlag = append([]string{"-h"}, cliArgsWithoutConfigFlag...)
	}

	cliErr := kernelObj.Run(cliArgsWithoutConfigFlag)
	if nil != cliErr {
		fmt.Printf("Command finished with error: %s", cliErr.Message())
		os.Exit(cliErr.Status())
	}
}
