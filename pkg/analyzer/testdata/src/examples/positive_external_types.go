package examples

import (
	"go/token"
	"io"
	"net/http"
	"os"

	"somepkg"
)

func structPtrTypeExtPkg() (*os.File, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func primitivePtrTypeExtPkg() (*token.Token, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func channelTypeExtPkg() (somepkg.ChannelType, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func funcTypeExtPkg() (http.HandlerFunc, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func ifaceTypeExtPkg() (io.Closer, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

type closerAlias = io.Closer

func ifaceTypeAliasedExtPkg() (closerAlias, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}
