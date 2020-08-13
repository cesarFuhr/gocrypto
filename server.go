package main

import (
	"fmt"
	"net/http"
)

func KeyServer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{ 
		"publicKey": "-----BEGIN RSA PUBLIC KEY-----\n\tMIGJAoGBAKE4I3DWqAAd8z6SsrmibIneKeHXP8gfvKuHPi7ePv3ZnK1AlfToHQ16\n\tirvEMcEoa+BNQEAbBot8aB0UrvwuhMXyWHQ9ugf7AShfMNteTbsdbvKv10W0Hee9\n\tLnU+xYxmRgUNUPwba3xJG8oIshkwYXEQbQW+ys1sbTv6ohmGXdv7AgMBAAE=\n\t-----END RSA PUBLIC KEY-----"
		}`)
}
