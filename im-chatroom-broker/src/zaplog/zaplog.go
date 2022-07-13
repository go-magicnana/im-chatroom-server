package zaplog

import (
	"github.com/go-redis/redis/v8"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"sync"
	"time"
)

var Logger *zap.SugaredLogger


func getEncoder() zapcore.Encoder {
	//return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	//return zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())


	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.FunctionKey = "printMethodName"
	encoderConfig.ConsoleSeparator = " "

	return zapcore.NewConsoleEncoder(encoderConfig)

}

func getLogWriter() zapcore.WriteSyncer {

	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./embru.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}
	//return zapcore.AddSync(lumberJackLogger)
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberJackLogger), zapcore.AddSync(os.Stdout))


	//file, _ := os.Create("./training.log")
	//return zapcore.AddSync(file)
}

type Log struct {
	Logger *zap.SugaredLogger
}

var once sync.Once

//var Logger *zap.SugaredLogger
var client *redis.Client

func Singleton() *Log {
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     "47.95.148.121:6379",
			Password: "o1trUmeh", // no password set
			DB:       1,          // use default DB
		})
	})

	return nil
}


func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	//_logger := zap.New(core)
	_logger := zap.New(core, zap.AddCaller())

	Logger = _logger.Sugar()

	defer Logger.Sync()

}

func simpleHttpGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		Logger.Error(
			"Error fetching url..",
			zap.String("url", url),
			zap.Error(err))
	} else {
		Logger.Info("Success..",
			zap.String("statusCode", resp.Status),
			zap.String("url", url))
		resp.Body.Close()
	}
}

func Test1() {
	InitLogger()
	defer Logger.Sync()

	for i := 0; i < 4; i++ {
		Logger.Debugf("Trying to hit GET request for %s", "http://www.baidu.com")
		Logger.Infof("Success! statusCode = %s for URL %s", "http://www.baidu.com",200)
		Logger.Errorf("Error fetching URL %s : Error = %s", "http://www.baidu.com",400)
	}


	//simpleHttpGet("www.google.com")
	//simpleHttpGet("http://www.google.com")
}

func Debugf(template string, args ...interface{}) {

	Logger.Debugf(template,args,
		// Structured context as strongly typed Field values.
		zap.String("f1", `http://foo.com`),
		zap.Int("f2", 3),
		zap.Duration("f3", time.Second))

}

const (
	name = "im-chatroom-broker"
	ip = "192.168.1.1"
	port = "8080"
	trace = "xxxxx"
	span = "span"
	url = "www.baidu.com"

)

func Infof(template string, args ...interface{}) {

	//Logger.Info("Success..",
	//	zap.String("statusCode", resp.Status),
	//	zap.String("url", url))

	Logger.Infof("%s %s %s %s %s - %s"+template,name,ip,port,trace,span,url,args)

}

func Warnf(template string, args ...interface{}) {

	Logger.Warnf("%s %s %s %s %s - %s"+template,name,ip,port,trace,span,url,args)


}

func Errorf(template string, args ...interface{}) {

	Logger.Errorf("%s %s %s %s %s - %s"+template,name,ip,port,trace,span,url,args)


}





