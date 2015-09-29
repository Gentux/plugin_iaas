package libnan

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const ()

// ===================================================================================================
// TYPES
// ===================================================================================================

type Err struct {
	Code    int
	Message string
	Details string

	jsonBytes []byte `json:"-"`
}

type errer interface {
	Error() string
}

func NewErr() *Err {
	return NewErrf("undefined")
}

func ErrFrom(e errer) *Err {
	if e == nil {
		return nil
	} else {
		return NewExitCode(1, e.Error())
	}
}

func NewErrf(_msg string, _args ...interface{}) *Err {
	msg := fmt.Sprintf(_msg, _args...)

	return &Err{Code: 0, Message: msg}
}

func NewExitCode(_code int, _msg string) *Err {
	p := &Err{Code: _code, Message: _msg}

	jsonBytes, err := json.Marshal(p)
	if err != nil {
		log.Printf("Error when json marshalling : { %d, %s }", _code, _msg)
	}

	p.jsonBytes = jsonBytes

	return p
}

func (o Err) Ok() bool {
	return o.Code == 1
}

func (o Err) Failed() bool {
	return o.Code != 1
}

func (p *Err) ToJson() string {
	return string(p.jsonBytes)
}

func (p *Err) ToString() string {
	return p.Message
}

func (p *Err) Unmarshal(s string) bool {
	if err := json.Unmarshal([]byte(s), p); err != nil {
		LogError("Failed to unmarshal exit code from string: %s", s)
		return false
	}
	return true
}

func PrintErrorJson(_pError *Err) {
	fmt.Println(_pError.ToJson())
}

func PrintOk(_pExitCode *Err) {
	Log(_pExitCode.Message)
	fmt.Println(_pExitCode.ToJson())
}

func ExitOk(_pExitCode *Err) {
	Log(_pExitCode.Message)
	fmt.Println(_pExitCode.ToJson())

	os.Exit(0)
}

func ExitError(_pExitCode *Err) {
	LogError(_pExitCode.Message)
	fmt.Println(_pExitCode.ToJson())

	os.Exit(-1)
}

func Errorf(_msg string, _args ...interface{}) *Err {
	msg := fmt.Sprintf(_msg, _args...)

	return &Err{Code: 0, Message: msg}
}

func ExitErrorf(_code int, _msg string, _args ...interface{}) {
	msg := fmt.Sprintf(_msg, _args...)

	LogError(msg)

	p := &Err{Code: _code, Message: msg}

	jsonBytes, err := json.Marshal(p)
	if err != nil {
		log.Printf("Error when json marshalling : { %d, %s }", _code, _msg)
		return
	}

	fmt.Println(string(jsonBytes))
	os.Exit(_code)
}

// ===================================================================================================

var (
	ErrUnset                = NewExitCode(0, "Error not set")
	ErrSomethingWrong       = NewExitCode(0, "Something went wrong")
	ErrOk                   = NewExitCode(1, "Operation succeeded")
	OkSuccess               = NewExitCode(1, "Operation succeeded")
	ErrConfigError          = NewExitCode(2, "Error in config file")
	ErrPluginError          = NewExitCode(2, "Plugin error")
	ErrUnknownUuid          = NewExitCode(3, "Unknown resource UUID")
	ErrOpFailed             = NewExitCode(4, "Operation failed: resource is not in the state required to perform the operation")
	ErrPbWithEmailFormat    = NewExitCode(2, "Problem with email format")
	ErrPasswordNonCompliant = NewExitCode(3, "The password does not respect the security policy")
	ErrPasswordNotUpdated   = NewExitCode(5, "Update password failed")
	ErrFilesystemError      = NewExitCode(16, "Filesystem error : failed to create/delete file/directory")
	ErrSystemError          = NewExitCode(17, "System error")
	ErrConsulNotFound       = NewExitCode(100, "Could not access Consul executable")
	ErrCouldNotPingVm       = NewExitCode(101, "Could not ping VM")
	ErrJsonParsingError     = NewExitCode(103, "Error when parsing JSON, see TAC log")
	ErrSshConnFailureonVm   = NewExitCode(104, "Failed to initiate SSH root on vm")

	ErrErrorWithExternalExe = NewExitCode(2, "Error returned by external executable, see TAC log")

	// VM creation
	ErrDuringVmCreation        = NewExitCode(0, "nc_create_vm did not return a valid VM uuid")
	ErrFailedToLocateVmProcess = NewExitCode(1, "Failed to locate VM pid in process list")

	// "message" : "Corrupt state, pool said to be non empty but has no VMs listed" }`)

	// Inside VM
	ErrCommandFailedInVm = NewExitCode(105, "Failed to run command on VM via SSH")
)
