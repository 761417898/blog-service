package main

import (
	"blog-service/global"
	"blog-service/intelnal/model"
	"blog-service/intelnal/routers"
	"blog-service/pkg/logger"
	"blog-service/pkg/setting"
	"blog-service/pkg/tracer"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	err := setupSetting()
	if err != nil {
		log.Fatal("SetupSetting failed, ", err)
	}
	err = setupDBEngine()
	if err != nil {
		log.Fatal("SetupDBEngine failed, ", err)
	}
	err = setupLogger()
	if err != nil {
		log.Fatal("SetupLogger failed, ", err)
	}
	err = setupTracer()
	if err != nil {
		log.Fatal("SetupTracer failed, ", err)
	}
}

// @title 博客系统
// @version 1.0
func main() {
	gin.SetMode(global.ServerSetting.RunMode)
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	global.Logger.Info(global.ServerSetting, global.DatabaseSetting, global.AppSetting)
	global.Logger.Infof("%s:go-programming-language", "gaohongzhi")
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("s.ListenAndServer err: %v", err)
		}
	}()
	// 等待中断信号
	quit := make(chan os.Signal)
	// 接受SIGINT和SIGTERM信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Println("Server forced shutdown")
	}
	/*
		eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBfa2V5IjoiMjc1NjY4YmE2NTUwNDljZDczOWQxZDllNmIzMWNjZjEiLCJhcHBfc2VjcmV0IjoiN2M5NzI2NjMxNzBkNmJjMTg0ODRkMDViYzk4NzIyZjQiLCJleHAiOjE2MTg2Nzg3MjQsImlzcyI6ImJsb2ctc2VydmljZSJ9.JHioLzib
		VIPBMmpKANFYIokw3QmJGMqaS1eDjhCB6kQ
	*/

}

func setupSetting() error {
	setting, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = setting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("App", &global.AppSetting)
	if err != nil {
		return nil
	}
	err = setting.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("JWT", &global.JWTSetting)
	if err != nil {
		return err
	}
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	global.JWTSetting.Expire *= time.Second
	return nil
}

func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}

	return nil
}

func setupLogger() error {
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt,
		MaxSize:   global.AppSetting.MaxPageSize,
		MaxAge:    global.AppSetting.MaxAge,
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)
	return nil
}

func setupTracer() error {
	jaegerTracer, _, err := tracer.NewJaegerTracer("blog-service", "127.0.0.1:6831")
	if err != nil {
		return err
	}
	global.Tracer = jaegerTracer
	return nil
}

/*
var prefixs = []string{"hot", "clod"}
var fileTypes = []string{"text", "jpeg", "log", "flv"}

func GetRandomString(n int) string {
	randBytes := make([]byte, n)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

func generatePrefix()  {
	for i := 0; i < 1000; i++ {
		randStr := GetRandomString(50)
		prefixs = append(prefixs, randStr)
	}
}

func main() {
	f, err := os.Create("train.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	//f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	w := csv.NewWriter(f)

	//w.Write([]string{"FileName", "FileSize", "StorageType"})
	generatePrefix()
	for i := 0; i < 11000; i++ {
		isStandard := "0"
		prefixIdx := rand.Int() % 1000
		prefix := prefixs[prefixIdx]
		fileType := ""
		if prefixIdx < 5 {
			 isStandard = "1"
		}
		fileType = fileTypes[rand.Int() % 4]
		if fileType == "text" || fileType == "jpeg" {
			isStandard = "1"
		}
		fileSize := rand.Int() % 512 + 1  // 1~512KB
		if fileType == "flv" {
			fileSize = rand.Int() % 4000 + 1000  // 1~5MB
		}
		if fileType == "log" {
			fileSize = rand.Int() % 100000 + 10000 // 10~100MB
		}
		fileName := prefix + GetRandomString(rand.Int() % 50 + 1)
		fmt.Println(fileName)
		w.Write([]string{fileName, fileType, strconv.Itoa(fileSize), isStandard})
	}
	w.Flush()

}
*/
