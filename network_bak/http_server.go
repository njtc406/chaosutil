package network_bak

import (
	"crypto/tls"
	"errors"
	"github.com/njtc406/chaosutil/log"

	"net/http"
	"time"
)

// TODO http后面看能不能换成gin框架来做，东西处理起来会简单点

var DefaultMaxHeaderBytes int = 1 << 20

// CAFile 证书文件
type CAFile struct {
	CertFile string
	Keyfile  string
}

// HttpServer http服务器
type HttpServer struct {
	listenAddr   string
	readTimeout  time.Duration
	writeTimeout time.Duration

	handler    http.Handler
	caFileList []CAFile

	httpServer *http.Server

	logger log.ILogger
}

func (slf *HttpServer) Init(listenAddr string, handler http.Handler, readTimeout time.Duration, writeTimeout time.Duration, logger log.ILogger) {
	slf.listenAddr = listenAddr
	slf.handler = handler
	slf.readTimeout = readTimeout
	slf.writeTimeout = writeTimeout
	if logger == nil {
		panic("http server required a logger")
	}
	slf.logger = logger

}

func (slf *HttpServer) Start() {
	go slf.startListen()
}

func (slf *HttpServer) startListen() error {
	if slf.httpServer != nil {
		return errors.New("Duplicate start not allowed")
	}

	var tlsCaList []tls.Certificate
	var tlsConfig *tls.Config
	for _, caFile := range slf.caFileList {
		cer, err := tls.LoadX509KeyPair(caFile.CertFile, caFile.Keyfile)
		if err != nil {
			slf.logger.Infof("Load CA  [%s]-[%s] file is fail:%s", caFile.CertFile, caFile.Keyfile, err.Error())
			return err
		}
		tlsCaList = append(tlsCaList, cer)
	}

	if len(tlsCaList) > 0 {
		tlsConfig = &tls.Config{Certificates: tlsCaList}
	}

	slf.httpServer = &http.Server{
		Addr:           slf.listenAddr,
		Handler:        slf.handler,
		ReadTimeout:    slf.readTimeout,
		WriteTimeout:   slf.writeTimeout,
		MaxHeaderBytes: DefaultMaxHeaderBytes,
		TLSConfig:      tlsConfig,
	}

	var err error
	if len(tlsCaList) > 0 {
		err = slf.httpServer.ListenAndServeTLS("", "")
	} else {
		err = slf.httpServer.ListenAndServe()
	}

	if err != nil {
		slf.logger.Infof("Listen for address %s failure: %s", slf.listenAddr, err.Error())
		return err
	}

	return nil
}

func (slf *HttpServer) SetCAFile(caFile []CAFile) {
	slf.caFileList = caFile
}
