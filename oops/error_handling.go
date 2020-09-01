/*
 * ----- Authors -----
 * marcelobezer
 * cahe7cb
 * washingt0
 */

package oops

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"go.uber.org/zap"

	"go-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
	"gopkg.in/go-playground/validator.v9"

	"github.com/pkg/errors"
)

const (
	pgxCode         = 1000
	jsonCode        = 2000
	internalCode    = 3000
	defaultCode     = 4000
	validationCode  = 5000
	grpcCode        = 6000
	timeParseError  = 7000
	httpRequestCode = 8000
)

// Error fit a error type for handling
type Error struct {
	Msg        string   `json:"msg"`
	Code       int      `json:"code"`
	Trace      []string `json:"-"`
	Err        error    `json:"-"`
	StatusCode int      `json:"-"`
}

// Error enable Error type to implements error type
func (e *Error) Error() string {
	return e.Msg
}

// Unwrap return the specific error cause for this error
func (e *Error) Unwrap() error {
	return e.Err
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type wrappedError interface {
	Unwrap() error
}

// fromErr wraps errors to provide user readable messages
func fromErr(rawError error) error {
	msg, code, responseStatus := "Erro desconhecido", 0, 400
	switch err := rawError.(type) {
	// input errors
	case *json.UnmarshalTypeError:
		msg, code = fmt.Sprintf("Tipo de valor %v não suportado no campo %v. Esperado tipo %v", err.Value, err.Field, err.Type.String()), jsonCode+1

	case validator.ValidationErrors:
		msg, code = parseValidationError(err)

	// internal errors
	case *reflect.ValueError:
		msg, code = fmt.Sprintf("Não é possível acessar o valor do tipo %v", err.Kind.String()), internalCode+1

	case *strconv.NumError:
		msg, code = fmt.Sprintf("Não é possível converter valor %v", err.Num), internalCode+2

	// data errors
	case pgx.PgError:
		msg, code = handlePgxError(&err)
		rawError = errors.Errorf("%s: %s", err.Error(), err.Hint)

	case *url.Error:
		msg, code = fmt.Sprintf("Falha no acesso à serviço. Operação: %v", err.Op), internalCode+3

	case *time.ParseError:
		msg, code = fmt.Sprintf("Impossível converter %v", err.Value), timeParseError+1

	case *Error:
		// this will create a deep copy of the Error struct
		rawError, msg, code, responseStatus = err, err.Msg, err.Code, err.StatusCode

	case *utils.HTTPError:
		msg, code = handleHTTPRequestError(err, httpRequestCode)

	case error:
		// Default errors
		switch err {
		case sql.ErrNoRows:
			msg, code = "Referência inválida", defaultCode+1
			responseStatus = http.StatusNotFound

		case io.EOF:
			msg, code = "Nenhum dado disponível para leitura", defaultCode+2
		}

		// external gRPC errors
		if s, ok := grpcStatus.FromError(err); ok {
			msg, code = s.Message(), grpcCode+int(s.Code())
			rawError = fmt.Errorf(fmt.Sprintf("%v", s.Details()))
			if s.Code() == grpcCodes.DeadlineExceeded {
				msg = "A consulta demorou mais do que o esperado, tente novamente."
			}
		}
	case nil:
		return nil
	}

	return &Error{
		Msg:        msg,
		Err:        rawError,
		Code:       code,
		StatusCode: responseStatus,
	}
}

// Err builds annotated error instance from any error value
func Err(err error) error {
	var e *Error
	if !errors.As(err, &e) {
		// annotate this error if it wasn't already
		err = fromErr(err)
	} else if err == e {
		// this is our error type but it hasn't been annotated
		err = fromErr(err)
	}
	return errors.WithStack(err)
}

// Wrap wraps an error adding an information message
func Wrap(err error, message string) error {
	return errors.Wrap(Err(err), message)
}

// NewErr creates an annotated error instance with default values
func NewErr(message string) error {
	return Err(&Error{
		Msg:        message,
		Err:        errors.Errorf("Inline error message: '%s'. See the stack trace of the error for additional information.", message),
		Code:       defaultCode,
		StatusCode: http.StatusBadRequest,
	})
}

// Handling handles an error by setting a message and a response status code
func Handling(err error, c *gin.Context) {
	var e *Error

	if !errors.As(err, &e) {
		Handling(Err(err), c)
		return
	}
	e.Msg = err.Error()
	e.Trace, _ = reconstructStackTrace(err, e)

	c.JSON(e.StatusCode, e)
	c.Set("error", err)
	c.Abort()
}

func reconstructStackTrace(err error, bound error) (output []string, traced bool) {
	var (
		wrapped wrappedError
		tracer  stackTracer
	)
	if errors.As(err, &wrapped) {
		internal := wrapped.Unwrap()
		// stop looking as we found our error instance
		if internal != bound {
			output, traced = reconstructStackTrace(internal, bound)
		}
		if !traced && errors.As(err, &tracer) {
			stack := tracer.StackTrace()
			for _, frame := range stack {
				output = append(output, fmt.Sprintf("%+v", frame))
			}
			traced = true
		}
	}
	return
}

func handleHTTPRequestError(err *utils.HTTPError, baseCode int) (string, int) {
	// define some not very informative defaults
	serviceDescription := "serviço externo"
	methodDescription := "executar ação solicitada"
	statusDescription := "ocorreu um problema inesperado"
	userInstruction := "tente novamente mais tarde"

	switch err.Method {
	case "GET":
		methodDescription = "buscar dados"
	case "POST":
		methodDescription = "cadastrar novos dados"
	case "PUT":
		methodDescription = "atualizar informações"
	case "DELETE":
		methodDescription = "remover registro"
	}

	// attempt to identify the external service
	if err.ClientName != "" {
		serviceDescription = "serviço de " + err.ClientName
	}

	// identify which action caused the problem
	if err.RequestTag != "" {
		methodDescription = err.RequestTag
	}

	switch {
	case err.StatusCode >= 300 && err.StatusCode < 400:
		statusDescription = "não está aceitando requisições"
	case err.StatusCode >= 400 && err.StatusCode < 500:
		statusDescription = "foi passado parâmetros incorretos"
	case err.StatusCode >= 500:
		statusDescription = "encontrou um problema interno"
	case err.StatusCode == -1:
		statusDescription = "encontra-se indisponível"
	}

	if cause, ok := err.Cause.(net.Error); ok {
		if cause.Timeout() {
			statusDescription = "demorou mais que o esperado"
			userInstruction = "tente novamente"
		}
	}

	return fmt.Sprintf("Não foi possível %s. O %s %s, %s.", methodDescription, serviceDescription, statusDescription, userInstruction), baseCode
}

func handlePgxError(err *pgx.PgError) (string, int) {
	switch err.Code {
	case "23505":
		return "Registro duplicado", pgxCode + 1
	case "23502":
		return "Dado requerido não foi especificado", pgxCode + 2
	case "23503":
		return "Dado indicado não é uma referência válida", pgxCode + 3
	case "42P01", "42703":
		return "Acesso incorreto de elementos nos registros de dados: erro de síntaxe", pgxCode + 4
	case "42601", "42803", "42883":
		return "Uso incorreto de função ou operador durante acesso aos registros de dados: erro de sintax", pgxCode + 5
	case "22001":
		return "Dado excede capacidade do registro no banco de dados", pgxCode + 6
	case "42702":
		return "Referência ambigua: erro de sintax", pgxCode + 7
	}

	// FIXME: This log should stay here to make easy handle unknown pgx errors.
	// It should be removed when we have an acceptable number of errors been handled
	log.Println(err.Code)
	return "Erro de dados desconhecido", pgxCode
}

func parseValidationError(err validator.ValidationErrors) (msg string, code int) {
	msg, code = "Não foi possível definir o erro de validação", validationCode

	if len(err) == 0 {
		return
	}

	switch err[0].ActualTag() {
	case "required":
		msg, code = "Campo "+err[0].Field()+" é obrigatorio", validationCode+1
	case "gt":
		msg, code = "Campo "+err[0].Field()+" deve ser maior que "+err[0].Param(), validationCode+2
	case "customerDocument":
		msg, code = "Documento inválido", validationCode+3
	case "gte":
		msg, code = "Campo "+err[0].Field()+" deve ser maior ou igual a "+err[0].Param(), validationCode+4
	case "stringField":
		msg, code = "Campo "+err[0].Field()+" não é uma string valida", validationCode+5
	case "required_with":
		msg, code = "Campo "+err[0].Field()+" é obrigatório quando campo "+err[0].Param()+" é enviado", validationCode+6
	case "required_without":
		msg, code = "Campo "+err[0].Field()+" é obrigatório se não for enviado o campo "+err[0].Param(), validationCode+7
	case "email":
		msg, code = "Campo "+err[0].Field()+" não contém email válido "+err[0].Param(), validationCode+8
	case "len":
		msg, code = "Campo "+err[0].Field()+" deve possuir tamanho igual a "+err[0].Param(), validationCode+9
	case "min":
		switch err[0].Kind() {
		case reflect.Int64, reflect.Int, reflect.Float64:
			msg, code = "Campo "+err[0].Field()+" deve possuir um valor de no mínimo "+err[0].Param(), validationCode+10
		case reflect.Array, reflect.Slice, reflect.String:
			msg, code = "Campo "+err[0].Field()+" deve possuir um tamanho de no mínimo "+err[0].Param(), validationCode+10
		default:
			msg, code = "Campo "+err[0].Field()+" deve possuir no mínimo "+err[0].Param(), validationCode+10
		}
	case "max":
		switch err[0].Kind() {
		case reflect.Int64, reflect.Int, reflect.Float64:
			msg, code = "Campo "+err[0].Field()+" deve possuir um valor de no máximo "+err[0].Param(), validationCode+11
		case reflect.Array, reflect.Slice, reflect.String:
			msg, code = "Campo "+err[0].Field()+" deve possuir um tamanho de no máximo "+err[0].Param(), validationCode+11
		default:
			msg, code = "Campo "+err[0].Field()+" deve possuir no máximo "+err[0].Param(), validationCode+11
		}
	}

	return
}

// PassRequired allow a request to pass even if a required field is missing
func PassRequired(err error) error {
	req := false
	if e, ok := err.(validator.ValidationErrors); ok {
		req = true
		for _, v := range e {
			if v.ActualTag() != "required" {
				req = false
			}
		}
	}

	if req {
		return nil
	}

	return err
}

func getErrorLocation(skip int) string {
	_, file, line, _ := runtime.Caller(skip + 1)
	return file + ":" + strconv.Itoa(line)
}

// Log just format and log one or more errors on
// the log file
func Log(title string, errs ...error) {
	zap.L().Error(
		title,
		zap.String("location", getErrorLocation(1)),
		zap.Errors("errors", errs),
	)
}
