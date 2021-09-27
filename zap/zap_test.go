package zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"testing"
	"time"
)

var url = "www.google.com"

func TestZap1(t *testing.T){
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)
}

func TestZap2(t *testing.T){
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("failed to fetch URL",
		// Structured context as strongly typed Field values.
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
}
//////////////////////////////////////////////////////////////////////
var sugarLogger *zap.SugaredLogger

func getEncoder()zapcore.Encoder{
	//return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()) //修改输出格式
	return  zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter()zapcore.WriteSyncer{
	file,_ := os.Create("./test.log")
	return zapcore.AddSync(file)
}

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core)
	sugarLogger = logger.Sugar()
}

func simpleHttpGet(url string) {
	sugarLogger.Debugf("Trying to hit GET request for %s", url)
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
	} else {
		sugarLogger.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
		resp.Body.Close()
	}
}

func TestSimpleHttpGet(t *testing.T){
	InitLogger()
	defer sugarLogger.Sync()
	simpleHttpGet("www.google.com")
}
