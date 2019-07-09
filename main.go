package main

import (
	"flag"
	"fmt"
	"github.com/bassbeaver/gkernel"
	kernelResponse "github.com/bassbeaver/gkernel/response"
	"gkernel-skeleton/controller"
	authService "gkernel-skeleton/service/auth"
	csrfService "gkernel-skeleton/service/csrf"
	requestLoggerService "gkernel-skeleton/service/request_logger"
	sessionService "gkernel-skeleton/service/session"
	userProvider "gkernel-skeleton/service/user_provider"
	"html/template"
	"net/http"
	"net/http/pprof"
	"os"
)

func main() {
	flags := flag.NewFlagSet("flags", flag.PanicOnError)
	configPathFlag := flags.String("config", "", "Path to config file")
	profilingFlag := flags.Bool("profiling", false, "Flag to enable profiling (with pprof) and register profiling URLs")
	flagsErr := flags.Parse(os.Args[1:])
	if nil != flagsErr {
		panic(flagsErr)
	}

	//******************************** Creating application Kernel ********************************
	var configPath string
	if "" == *configPathFlag {
		curBinDir, curBinDirError := os.Getwd()
		if curBinDirError != nil {
			fmt.Println("Failed to determine path to own binary file")
			panic(curBinDirError)
		}
		configPath = curBinDir + "/config"
	} else {
		configPath = *configPathFlag
	}

	kernelObj, kernelError := gkernel.NewKernel(configPath)
	if nil != kernelError {
		panic(kernelError)
	}

	//******************************** Service's registration ********************************

	// ---- Controllers
	var controllerRegError error

	controllerRegError = kernelObj.RegisterService(
		"IndexController",
		func() *controller.IndexController {
			return &controller.IndexController{}
		},
		true,
	)
	if nil != controllerRegError {
		panic(fmt.Sprintf("failed to register IndexController service, error: %s", controllerRegError.Error()))
	}

	// ---- Services from packages
	authService.Register(kernelObj)
	userProvider.Register(kernelObj)
	csrfService.Register(kernelObj)
	requestLoggerService.Register(kernelObj)
	sessionService.Register(kernelObj)

	//******************************** Custom templates functions registration ********************************
	kernelObj.GetTemplates().Funcs(template.FuncMap{
		"sequence": func(size int) []int {
			sequence := make([]int, size)
			for i := 0; i < size; i++ {
				sequence[i] = i
			}

			return sequence
		},
		"addInt": func(a, b int) int {
			return a + b
		},
		"subInt": func(a, b int) int {
			return a - b
		},
	})

	//************** Profiler (pprof) configuration **************
	if *profilingFlag {
		registerPprofRoute := func(name, url string, handlerFunc http.HandlerFunc) {
			kernelObj.RegisterRoute(&gkernel.Route{
				Name:    name,
				Url:     url,
				Methods: []string{http.MethodGet},
				Controller: func(request *http.Request) kernelResponse.Response {
					w := kernelResponse.NewBytesResponseWriter()
					handlerFunc(w, request)
					return w
				},
			})
		}

		registerPprofRoute("pprof:index", "/debug/pprof/", pprof.Index)
		registerPprofRoute("pprof:cmdline", "/debug/pprof/cmdline", pprof.Cmdline)
		registerPprofRoute("pprof:profile", "/debug/pprof/profile", pprof.Profile)
		registerPprofRoute("pprof:symbol", "/debug/pprof/profile", pprof.Symbol)
		registerPprofRoute("pprof:trace", "/debug/pprof/trace", pprof.Trace)
		registerPprofRoute("pprof:heap", "/debug/pprof/heap", pprof.Index)
	}

	//************** Starting applications web-server **************
	kernelObj.Run()
}
