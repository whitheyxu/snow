// run
package snow

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/whitheyxu/snow/g/logs"
	"github.com/whitheyxu/snow/router"
)

type SnowApplication struct {
	Server        *http.Server
	DefaultLogger io.Writer
}

type ServerConfig struct {
	Addr           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

var defaultServer = &http.Server{
	Addr:           ":8080",
	ReadTimeout:    5 * time.Second,
	WriteTimeout:   5 * time.Second,
	MaxHeaderBytes: 1 << 20,
	Handler:        router.GetRouter(),
}

var snowApplication *SnowApplication

func init() {

	snowApplication = GetDefaultSnowApplication()
}

func Run() {
	logs.Infof("Snow is listening [ %s ] \n", snowApplication.Server.Addr)
	router.CheckRouter()
	err := snowApplication.Server.ListenAndServe()

	if err != nil {
		//		fmt.Println(err)
		logs.Crit(err)
	}
}

func GetDefaultSnowApplication() (sa *SnowApplication) {
	sa = new(SnowApplication)
	sa.Server = defaultServer
	sa.DefaultLogger = os.Stdout
	return
}

func SetServerAddr(addr string) {
	snowApplication.Server.Addr = addr
}

func SetServerReadTimeout(duration int) {
	snowApplication.Server.ReadTimeout = time.Duration(duration) * time.Second
}

func SetServerWriteTimeout(duration int) {
	snowApplication.Server.WriteTimeout = time.Duration(duration) * time.Second
}

func SetServerMaxHeaderBytes(maxBytes int) {
	snowApplication.Server.MaxHeaderBytes = maxBytes
}
