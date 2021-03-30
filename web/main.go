package main

import (
	"flag"
	"fmt"
	webKernel "github.com/bassbeaver/gkernel/web"
	kernelResponse "github.com/bassbeaver/gkernel/web/response"
	loggerFactoryService "gkernel-skeleton/common/service/logger_factory"
	redisService "gkernel-skeleton/common/service/redis"
	webController "gkernel-skeleton/web/controller"
	authService "gkernel-skeleton/web/service/auth"
	csrfService "gkernel-skeleton/web/service/csrf"
	requestLoggerService "gkernel-skeleton/web/service/request_logger"
	requestSizeValidatorService "gkernel-skeleton/web/service/request_size_validator"
	sessionService "gkernel-skeleton/web/service/session"
	userProvider "gkernel-skeleton/web/service/user_provider"
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

	kernelObj, kernelError := webKernel.NewKernel(configPath)
	if nil != kernelError {
		panic(kernelError)
	}

	//******************************** Service's registration ********************************

	// ---- Controllers
	webController.RegisterIndex(kernelObj)

	// ---- Other services
	authService.Register(kernelObj)
	userProvider.Register(kernelObj)
	csrfService.Register(kernelObj)
	loggerFactoryService.Register(kernelObj)
	requestSizeValidatorService.Register(kernelObj)
	requestLoggerService.Register(kernelObj)
	redisService.Register(kernelObj)
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
			kernelObj.RegisterRoute(&webKernel.Route{
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
