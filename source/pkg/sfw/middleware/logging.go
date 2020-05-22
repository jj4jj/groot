package middleware

import (
	"bytes"
	"io/ioutil"
	"groot/proto/cserr"
	"groot/proto/csmsg"

	// "regexp"
	"time"

	proto "github.com/golang/protobuf/proto"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/willf/pad"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logging is a middleware function that logs the each request

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now().UTC()
		path := c.Request.URL.Path

		// reg := regexp.MustCompile("(/v1/user|/login)")
		// if !reg.MatchString(path) {
		// 	log.Infof("the path did not match (/v1/user|/login)")
		// 	return
		// }

		// skip for the health check request.
		if path == "/sd/health" || path == "/sd/ram" || path == "/sd/cpu" ||
			path == "/sd/disk" || path == "/upload_file" || path == "/upload_files" || path == "/download_file" {
			return
		}

		// Read the body content
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}

		//Restore the io.ReadCloser to it's original state
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		method := c.Request.Method
		ip := c.ClientIP()

		//log.Debugf("New Request come in,path: %s,Method: %s,body: '%s'",path,method,string(bodyBytes))

		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}

		c.Writer = blw

		//Continue
		c.Next()

		//Calculates the latency.
		end := time.Now().UTC()
		latency := end.Sub(start)

		code, message := -1, ""

		//get code and message
		response := csmsg.CSMsgHttpFrame{}
		if e := proto.Unmarshal(blw.body.Bytes(), &response); e != nil {
			log.Errorf(e, "response body can not unmarshal to model.Response struct,body: '%s'", blw.body.Bytes())
			code = int(cserr.CSErrCode_CS_ERR_INTERNAL)
			message = e.Error()
		} else {
			if response.RspHead != nil {
				code = int(response.RspHead.Code)
				message = response.RspHead.Msg
			}
		}

		log.Infof("%-13s | %-12s | %s %s | {code: %d,message: %s}", latency, ip,
			pad.Right(method, 5, ""), path, code, message)
	}
}
