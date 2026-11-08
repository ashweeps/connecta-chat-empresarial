package main

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

type Cliente struct {
	servidor string
	http     *http.Client
}

func NuevoCliente(servidor string) (*Cliente, error) {
	servidor = strings.TrimSpace(servidor)

	if servidor == "" {
		return nil, errors.New("la direccion del servidor es obligatoria")
	}

	cliente := &Cliente{
		servidor: servidor,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	return cliente, nil
}
