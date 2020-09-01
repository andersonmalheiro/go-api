package oops

import "errors"

var (
	// ErrDataNotFound defines that an expected
	// data cannot be found
	ErrDataNotFound = Error{
		Msg:        "Erro interno: não foi possível obter dado requerido para continuar operações",
		Code:       internalCode,
		StatusCode: 400,
		Err:        errors.New("Erro interno: não foi possível obter dado requerido para continuar operações"),
	}

	// ErrInvalidFilter indicates that a filter was set with
	// a invalid value
	ErrInvalidFilter = Error{
		Msg:        "Valor especificado em filtro de rota tem tipo inválido",
		Code:       defaultCode,
		StatusCode: 400,
		Err:        errors.New("Valor especificado em filtro de rota tem tipo inválido"),
	}

	// ErrMemcachedConn indicates that was not possible
	// connect to memcached
	ErrMemcachedConn = Error{
		Msg:        "Nenhuma conexão com o memcached aberta",
		Code:       internalCode,
		StatusCode: 500,
		Err:        errors.New("Nenhuma conexão com o memcached aberta"),
	}

	// ErrInvalidPlatform indicates that the user try
	// to access some system functionality by an
	// invalid platform
	ErrInvalidPlatform = Error{
		Msg:        "Acesso a partir de plataforma inválida",
		Code:       defaultCode,
		StatusCode: 409,
		Err:        errors.New("Acesso a partir de plataforma inválida"),
	}
)
