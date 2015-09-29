package libnan

import (
	"fmt"

	"os"
	"path"
	"reflect"
	"time"
)

const ()

// ===================================================================================================
// TYPES
// ===================================================================================================

// ===================================================================================================
//
// Nano core utils + configuration loader
//
// ===================================================================================================

const (
	g_timeLayout = "2006-01-02 15:04 (CEST)"
)

var (
	DryRun  bool = false
	ModeRef bool = false

	g_sCommandLine string

	g_pLogFile *os.File

	g_StartTime time.Time

	NRETRIES uint = 5
)

func init() {

	g_StartTime = time.Now()

	for _, arg := range os.Args {
		g_sCommandLine += (arg + " ")
	}

	LoadConfig()

	// Setup logger

	_, errDirAccess := os.Stat(path.Dir(g_Config.LogFilePath))

	if errDirAccess != nil {
		fmt.Println("Could not access log file at : %s\n", g_Config.LogFilePath)
		fmt.Println("Will use : default.log")

		// revert to defaults in case of error or config file path not reachable
		g_Config.LogFilePath = "default.log"
	}

	var err error
	if g_pLogFile, err = os.OpenFile(g_Config.LogFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660); err != nil {
		fmt.Println("Error when opening log file:", g_Config.LogFilePath)
		os.Exit(-1)
	}

	//TODO OPTIONALize this

	// if _, errAccessingConsulExe := os.Stat(Config().ConsulPath); errAccessingConsulExe != nil {
	// 	LogError("Config error, Consul exe not found at: %s\n", Config().ConsulPath)
	// 	ExitError(ErrConsulNotFound)
	// }
}

func _log(sPrefix, sLogMsg string) {
	currentTime := time.Now().Format(g_timeLayout)

	sStartTime := fmt.Sprintf("%d:%d:%d", g_StartTime.Hour(), g_StartTime.Minute(), g_StartTime.Second())

	sLogLine := fmt.Sprintf("%v - %s: %s [%s] %s\n", currentTime, sPrefix, g_sCommandLine, sStartTime, sLogMsg)

	g_pLogFile.WriteString(sLogLine)
}

func Debug(_str string, _args ...interface{}) {
	str := fmt.Sprintf(_str, _args...)

	if Config().Debug {
		fmt.Println(str)
	}
}

func Log(_str string, _args ...interface{}) {
	str := fmt.Sprintf(_str, _args...)

	_log("INFO", str)

	if Config().Debug {
		fmt.Println(str)
	}
}

func LogError(_str string, _args ...interface{}) *Err {
	str := fmt.Sprintf(_str, _args...)

	_log("ERROR", str)

	if Config().Debug {
		fmt.Println(str)
	}

	return &Err{Message: str}
}

func LogErrorCode(pError *Err) *Err {
	_log("ERROR", pError.Message)
	return pError
}

type ProcedureStruct struct {
	Result *Err
}

func (o ProcedureStruct) GetResult() *Err {
	return o.Result
}

type Procedure interface {
	Do() *Err
	GetResult() *Err
	Undo() *Err
}

func UndoIfFailed(proc Procedure) {

	if proc.GetResult() == nil {
		return
	}

	val := reflect.Indirect(reflect.ValueOf(proc))
	LogError("Undoing", val.Type().Name(), "because it failed")
	proc.Undo()
}
