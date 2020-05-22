package app

import (
	"errors"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"groot/comm/conf"
	"net/http"
	"reflect"
)

type (
	ServiceCtx struct {
		name         string
		mpRpc        map[string]*ServiceRpcCtx
		mpMethod     map[string]*ServiceMethodCtx
		mpStream     map[string]*ServiceStreamCtx
		pathSet      map[string]bool
		mpTcsMethod  map[string]*ServiceTcsMethodCtx
		tcsWorkerNum int
		logic        ServiceLogic
		middlewares  []string
	}
	ServiceLogic interface {
		SyncDb() error
		Init(*ServiceCtx) error
	}
	Server struct {
		name      string
		mpService map[string]*ServiceCtx
		router    *gin.Engine
	}
)

var (
	mpServer      map[string]*Server         = make(map[string]*Server)
	mpMiddleWares map[string]gin.HandlerFunc = make(map[string]gin.HandlerFunc)
	ErrMiddleWareNotExist = errors.New("middleware not found")
)

func NewServer(name string) *Server {
	if mpServer[name] != nil {
		log.Errorf(nil, "new server name:%s is exist already !", name)
		return nil
	}
	svr := Server{
		name:      name,
		mpService: make(map[string]*ServiceCtx),
		router:    gin.New(),
	}
	mpServer[name] = &svr
	return &svr
}

func GetOrNewServer(name string) *Server {
	if s, ok := mpServer[name]; ok {
		return s
	}
	svr := Server{
		name:      name,
		mpService: make(map[string]*ServiceCtx),
		router:    gin.New(),
	}
	mpServer[name] = &svr
	log.Infof("register new server:%s", name)
	return &svr
}

func (svr *Server) AddService(name string, logic ServiceLogic, middlewares ...string) *ServiceCtx {
	var s, ok = svr.mpService[name]
	if ok {
		log.Errorf(nil, "new serivice error for already exist:%s!", name)
		return s
	}
	svc := &ServiceCtx{
		name:         name,
		mpRpc:        make(map[string]*ServiceRpcCtx),
		mpMethod:     make(map[string]*ServiceMethodCtx),
		mpStream:     make(map[string]*ServiceStreamCtx),
		pathSet:      make(map[string]bool),
		logic:        logic,
		middlewares:  middlewares,
		mpTcsMethod:  make(map[string]*ServiceTcsMethodCtx),
		tcsWorkerNum: 1,
	}
	svr.mpService[name] = svc
	log.Infof("add service:%s on server:%s ...", name, svr.name)
	return svc
}

func (svc *ServiceCtx) BindRpc(path string, handler ServiceRpcHandler, middlewares ...string) {
	if svc.pathSet[path] == true {
		log.Errorf(nil, "path:%s.%s already register .", svc.name,path)
		return
	}
	if !checkServiceRpcHandlerType(handler){
		return
	}

	svc.pathSet[path] = true
	svc.mpRpc[path] = &ServiceRpcCtx{
		handler:     handler,
		middlewares: middlewares,
		caller: reflect.ValueOf(handler),
		reqType: reflect.TypeOf(handler).In(1).Elem(),
		rspType: reflect.TypeOf(handler).In(2).Elem(),
	}
	log.Debugf("register rpc path:%s.%s ...", svc.name, path)
}

func (svc *ServiceCtx) Route(method string, path string, handler ServiceRouteHandler, middlewares ...string) {
	if svc.pathSet[path] == true {
		log.Errorf(nil, "path:%s.%s already register .", svc.name,path)
		return
	}
	checkServiceRouteHandlerType(handler)
	svc.pathSet[path] = true
	svc.mpMethod[path] = &ServiceMethodCtx{
		method:       method,
		routeHandler: handler,
		rawhandler:   nil,
		middlewares:  middlewares,
	}
	log.Debugf("register route path:%s.%s with meth:%s ...", svc.name, path, method)
}

func (svc *ServiceCtx) RouteRaw(method string, path string, handler gin.HandlerFunc, middlewares ...string) {
	if svc.pathSet[path] == true {
		log.Errorf(nil, "path:%s.%s already register .", svc.name,path)
		return
	}
	svc.pathSet[path] = true
	svc.mpMethod[path] = &ServiceMethodCtx{
		method:       method,
		routeHandler: nil,
		rawhandler:   handler,
		middlewares:  middlewares,
	}
	log.Debugf("register raw route path:%s.%s with meth:%s ...", svc.name, path, method)

}

func (svc *ServiceCtx) BindStream(path string, logic ServiceStreamLogic, reqInst interface{} ,
								heartBeatTimeOut int64, middlewares ...string) {
	if svc.pathSet[path] == true {
		log.Errorf(nil, "path:%s.%s already register .", svc.name,path)
		return
	}
	svc.pathSet[path] = true
	svc.mpStream[path] = &ServiceStreamCtx{
		logic:             logic,
		reqType:           reflect.TypeOf(reqInst),
		middlewares:       middlewares,
		heartBeatInterval: heartBeatTimeOut/2,
	}
	log.Debugf("register stream path:%s.%s ...", svc.name, path)
}
func (svc * ServiceCtx) TcsSetWorkerNum(tcsWokerNum int){
	if tcsWokerNum > 0 {
		svc.tcsWorkerNum  = tcsWokerNum
	}
}

func (svc * ServiceCtx) TcsListen(methodName string, handler ServiceMethodHandler){
	_, ok := svc.mpTcsMethod[methodName]
	if ok {
		log.Errorf(nil, "tcs path:%s.%s already register .", svc.name, methodName)
		return
	}
	if !checkServiceTcsMethodHandler(handler) {
		log.Errorf(nil, "tcs:%s.%s method handler is invalid",svc.name, methodName)
		return
	}
	handlerVal := reflect.ValueOf(handler)
	svc.mpTcsMethod[methodName] = &ServiceTcsMethodCtx {
		handler: handler,
		caller: handlerVal,
		reqValType: handlerVal.Type().In(0).Elem(),
		name: methodName,
	}
	log.Debugf("register tcs method path:%s.%s ...", svc.name, methodName)
}

func RegisterService(svcName string, svcLogic ServiceLogic, serverName string, middlewares ...string) *ServiceCtx {
	if serverName == "" {
		serverName = "main"
	}
	server := GetOrNewServer(serverName)
	return server.AddService(svcName, svcLogic, middlewares...)
}

func GetServerRouter(name string) *gin.Engine {
	svr, ok := mpServer[name]
	if ok {
		return svr.router
	}
	return nil
}

func InitSyncDb(name string) error {
	svr, ok := mpServer[name]
	if !ok {
		log.Errorf(nil, "not found the server:%s", name)
		return errors.New("not found server")
	}
	for _, svc := range svr.mpService {
		if err := svc.logic.SyncDb(); err != nil {
			log.Errorf(err, "init service name:%s sync db error !", svc.name)
			return err
		}
	}
	return nil
}


func InitServices(name string, mw ...gin.HandlerFunc) error {
	svr, ok := mpServer[name]
	if !ok {
		log.Errorf(nil, "not found the server:%s", name)
		return errors.New("not found server")
	}

	svr.router.Use(mw...)
	svr.router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route \n")
	})
	config := conf.GetAppConfig()
	if config.RunEnv == "dev" {
		ginpprof.Wrap(svr.router)
	}
	for _, svc := range svr.mpService {
		if err := svc.logic.Init(svc); err != nil {
			log.Errorf(err, "init service name:%s error !", svc.name)
			return err
		}
		gr := svr.router.Group(svc.name)
		for _, mw := range svc.middlewares {
			ware := GetMiddleWare(mw)
			if ware == nil {
				log.Errorf(nil, "middleware:%s for service:%s is not exist", mw, svc.name)
				return ErrMiddleWareNotExist
			}
			gr.Use(ware)
		}
		for path, rpc := range svc.mpRpc {
			var wares []gin.HandlerFunc
			for _, mw := range rpc.middlewares {
				ware := GetMiddleWare(mw)
				if ware == nil {
					log.Errorf(nil, "middleware:%s for path:%s is not exist", mw, path)
					return ErrMiddleWareNotExist
				}
				wares = append(wares, ware)
			}
			//
			wares = append(wares, createRpcWrapper(rpc))
			//所有RPC通用POST
			gr.POST(path, wares...)
		}
		for path, method := range svc.mpMethod {
			var wares []gin.HandlerFunc
			for _, mw := range method.middlewares {
				ware := GetMiddleWare(mw)
				if ware == nil {
					log.Errorf(nil, "middleware:%s for path:%s is not exist", mw, path)
					return ErrMiddleWareNotExist
				}
				wares = append(wares, ware)
			}
			//
			//fmt.Println("register api wares:", len(wares))
			wares = append(wares, createRouteWrapper(method))
			//
			switch method.method {
			case "GET":
				gr.GET(path, wares...)
			case "POST":
				gr.POST(path, wares...)
			case "PUT":
				gr.PUT(path, wares...)
			case "DELETE":
				gr.DELETE(path, wares...)
			default:
				log.Errorf(nil, "not support method:%s for path:%s", method.method, path)
			}
		}
		for path, logic := range svc.mpStream {
			var wares []gin.HandlerFunc
			for _,wm := range logic.middlewares {
				ware := GetMiddleWare(wm)
				if ware == nil {
					log.Errorf(nil, "middleware:%s for path:%s is not exist", mw, path)
					return ErrMiddleWareNotExist
				}
				wares = append(wares, ware)
			}
			wares = append(wares, createStreamWrapper(logic))
			gr.GET(path, wares...)
		}

	}
	return nil
}

func AddMiddleWare(name string, middleware gin.HandlerFunc) {
	mpMiddleWares[name] = middleware
}

func GetMiddleWare(name string) gin.HandlerFunc {
	return mpMiddleWares[name]
}
