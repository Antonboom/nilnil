# nilnil

[![CI](https://github.com/Antonboom/nilnil/actions/workflows/ci.yml/badge.svg)](https://github.com/Antonboom/nilnil/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Antonboom/nilnil)](https://goreportcard.com/report/github.com/Antonboom/nilnil)
[![Coverage](https://coveralls.io/repos/github/Antonboom/nilnil/badge.svg?branch=master)](https://coveralls.io/github/Antonboom/nilnil?branch=master)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

Checks that there is no simultaneous return of `nil` error and an invalid value.

## Installation & usage

```
$ go install github.com/Antonboom/nilnil@latest
$ nilnil ./...
```

## Motivation

`return nil, nil` is not idiomatic for Go. The developers are used to the fact that 
if there is no error, then the return value is valid and can be used without additional checks:
```go
user, err := getUser()
if err != nil {
    return err
}
if user != nil { // Ambiguous!
    // Use user.
}
```
In the worst case, code like this can lead to **panic**.
<br>

Rewrite the example for sentinel error:
```go
user, err := getUser()
if errors.Is(err, errUserNotFound) {
    // Do something and return.
}
if err != nil {
    return err
}

// Use user.
```

### What if I think it's bullshit?

I understand that each case needs to be analyzed separately, 
but I hope that the linter will make you think again -
is it necessary to use an ambiguous API or is it better to do it using a sentinel error?
<br>

In any case, you can just not enable the linter.

## Configuration

```shell
# command line (see help for full list of types)
$ nilnil --checked-types ptr,func ./...
```

```yaml
# https://golangci-lint.run/usage/configuration/
nilnil:
  checked-types:
    - ptr
    - func
    - iface
    - map
    - chan
```

## Examples

<details>
  <summary>parsePublicKey from crypto/tls</summary>

```go
// BEFORE

func parsePublicKey(algo PublicKeyAlgorithm, keyData *publicKeyInfo) (interface{}, error) {
    der := cryptobyte.String(keyData.PublicKey.RightAlign())
    switch algo {
    case RSA:
        // ...
        return pub, nil
    case ECDSA:
        // ...
        return pub, nil
    case Ed25519:
        // ...
        return ed25519.PublicKey(der), nil
    case DSA:
        // ...
        return pub, nil
    default:
        return nil, nil
    }
}

// AFTER

var errUnknownPublicKeyAlgo = errors.New("unknown public key algo")

func parsePublicKey(algo PublicKeyAlgorithm, keyData *publicKeyInfo) (interface{}, error) {
    der := cryptobyte.String(keyData.PublicKey.RightAlign())
    switch algo {
    case RSA:
        // ...
        return pub, nil
    case ECDSA:
        // ...
        return pub, nil
    case Ed25519:
        // ...
        return ed25519.PublicKey(der), nil
    case DSA:
        // ...
        return pub, nil
    default:
        return nil, fmt.Errorf("%w: %v", errUnknownPublicKeyAlgo, algo)
    }
}
```

</details>

<details>
  <summary>http2clientConnReadLoop from net/http</summary>

```go
// BEFORE

// As a special case, handleResponse may return (nil, nil) to skip the frame.
func (rl *http2clientConnReadLoop) handleResponse(/* ... */) (*Response, error) {
    if statusCode >= 100 && statusCode <= 199 {
        return nil, nil
    }
}

// ...
res, err := rl.handleResponse(cs, f)
if err != nil {
	return err
}
if res == nil {
    // (nil, nil) special case. See handleResponse docs.
    return nil
}

// AFTER

var errNeedSkipFrame = errors.New("need skip frame")

// As a special case, handleResponse may return errNeedSkipFrame to skip the frame.
func (rl *http2clientConnReadLoop) handleResponse(/* ... */) (*Response, error) {
    if statusCode >= 100 && statusCode <= 199 {
        return nil, errNeedSkipFrame
    }
}

// ...
res, err := rl.handleResponse(cs, f)
if errors.Is(err, errNeedSkipFrame) {
    return nil
}
if err != nil {
    return err
}
```

</details>

<details>
  <summary>Not implemented</summary>

```go
// BEFORE

func (s *Service) StartStream(ctx context.Context) (*Stream, error) {
    return nil, nil
}

// AFTER

func (s *Service) StartStream(ctx context.Context) (*Stream, error) {
    return nil, errors.New("not implemented")
}
```

</details>

<details>
  <summary>nil-safe type</summary>

```go
package ratelimiter

type RateLimiter struct {
    // ...
}

func New() (*RateLimiter, error) {
    // It's OK, RateLimiter is nil-safe.
    // But it's better not to do it anyway.
    return nil, nil
}

func (r *RateLimiter) Allow() bool {
    if r == nil {
        return true
    }
    return r.allow()
}
```

</details>

## Assumptions

<details>
  <summary>Click to expand</summary>

<br>

- Linter only checks funcs with two return arguments, the last of which has `error` type.
- Next types are checked:
  * pointers, functions & interfaces (`panic: invalid memory address or nil pointer dereference`);
  * maps (`panic: assignment to entry in nil map`);
  * channels (`fatal error: all goroutines are asleep - deadlock!`)
- `uinptr` & `unsafe.Pointer` are not checked as a special case.
- Supported only explicit `return nil, nil`.
- Types from external packages are not supported.

</details>

## Check Golang source code

<details>
  <summary>Click to expand</summary>

```shell
$ cd $GOROOT/src
$ nilnil ./...
/usr/local/go/src/net/sockopt_posix.go:48:3: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/crypto/x509/parser.go:321:3: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/crypto/tls/key_agreement.go:45:2: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/database/sql/driver/types.go:157:3: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/database/sql/driver/types.go:231:3: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/database/sql/driver/types.go:262:4: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/debug/dwarf/entry.go:882:3: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/debug/dwarf/line.go:146:3: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/debug/dwarf/line.go:153:3: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/debug/dwarf/typeunit.go:138:3: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/debug/pe/file.go:450:3: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/net/http/h2_bundle.go:8644:3: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/net/http/transfer.go:768:3: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/net/http/transfer.go:778:3: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/net/http/transfer.go:801:3: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/go/build/build.go:1404:3: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/go/build/build.go:1414:3: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/go/build/build.go:1419:3: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/go/build/build.go:1453:3: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/internal/profile/legacy_profile.go:1087:4: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/net/internal/socktest/switch.go:142:3: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/crypto/tls/handshake_server_test.go:411:4: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/crypto/tls/handshake_server_test.go:1012:3: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/crypto/tls/handshake_server_test.go:1470:4: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/crypto/tls/tls_test.go:747:4: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/crypto/tls/tls_test.go:751:4: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/crypto/tls/tls_test.go:755:4: return both the `nil` error and invalid value: use a sentinel error instead
/usr/local/go/src/encoding/xml/xml_test.go:92:4: return both the `nil` error and invalid value: use a sentinel error instead
```

</details>
