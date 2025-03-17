# nilnil

![Latest release](https://img.shields.io/github/v/release/Antonboom/nilnil)
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

At best this leads to comments [like](https://github.com/golang/go/blob/07fc59199b9522bfe0d14f35c4391394efc336c9/src/debug/dwarf/line.go#L142)

```go
// ...
// If this compilation unit has no line table, it returns nil, nil.
func (d *Data) LineReader(cu *Entry) (*LineReader, error) {
```

which is not a _true_ constraint. And in the worst case such code can lead to **panic**.
<br>

Rewrite the example above with sentinel error:
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

### Opposite situation

Sometimes people consider the opposite situation (returning a non-nil error and a valid value at the same time) to be
an anti-pattern too, since it can lead to hard-to-find bugs.

For example, this **kubernetes** code

```go
func (fh *fakeHistory) UpdateControllerRevision(revision *apps.ControllerRevision, newRevision int64) (*apps.ControllerRevision, error) {
    clone := revision.DeepCopy()
    clone.Revision = newRevision
    return clone, fh.indexer.Update(clone)
}
```

can be rewritten as follows

```go
func (fh *fakeHistory) UpdateControllerRevision(revision *apps.ControllerRevision, newRevision int64) (*apps.ControllerRevision, error) {
    clone := revision.DeepCopy()
    clone.Revision = newRevision
    if err := fh.indexer.Update(clone); err != nil {
        return nil, fmt.Errorf("update index: %w", err)
    }
    return clone, nil
}
```

Use the `--detect-opposite` flag to cover such cases too.

### What if I think it's bullshit?

I understand that each case needs to be analyzed separately, 
but I hope that the linter will make you think again –
is it necessary to use an **ambiguous API** or is it better to do it using a sentinel error?
<br>

In any case, you can just not enable the linter.

## Configuration

### CLI

```shell
$ nilnil --checked-types chan,func,iface,map,ptr,uintptr,unsafeptr ./...
$ nilnil --detect-opposite ./...
$ nilnil --detect-opposite --checked-types ptr ./..
$ nilnil --only-two=false ./...
$ nilnil --detect-opposite --only-two=false --checked-types iface,map ./..
```

### golangci-lint

https://golangci-lint.run/usage/linters/#nilnil

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

- Linter only checks functions with two return arguments, the last of which implements `error`. To disable this behavior
and enable linter for N returned values, use `--only-two=false` flag. 
- Next types are checked:
  * pointers (including `uinptr` and `unsafe.Pointer`), functions and interfaces (`panic: invalid memory address or nil pointer dereference`);
  * maps (`panic: assignment to entry in nil map`);
  * channels (`fatal error: all goroutines are asleep - deadlock!`)
- Only explicit nils (e.g. `return nil, nil` or `return nil, nil, ..., nil`) are supported.

## Check Go 1.22.2 source code

<details>
  <summary>Click to expand</summary>

```shell
$ cd $GOROOT/src
$ nilnil ./... 2>&1 | grep -v "_test.go"  
/usr/local/go/src/internal/bisect/bisect.go:196:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/net/fd_unix.go:71:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/net/fd_unix.go:79:4: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/net/fd_unix.go:156:4: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/net/iprawsock_posix.go:36:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/net/tcpsock_posix.go:38:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/net/udpsock_posix.go:37:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/net/unixsock_posix.go:92:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/crypto/tls/key_agreement.go:46:2: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/crypto/tls/ticket.go:355:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/crypto/tls/ticket.go:359:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/database/sql/driver/types.go:157:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/database/sql/driver/types.go:232:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/database/sql/driver/types.go:263:4: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/database/sql/convert.go:548:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/database/sql/sql.go:205:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/database/sql/sql.go:231:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/database/sql/sql.go:257:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/database/sql/sql.go:284:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/database/sql/sql.go:311:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/database/sql/sql.go:337:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/database/sql/sql.go:363:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/database/sql/sql.go:389:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/database/sql/sql.go:422:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/debug/dwarf/entry.go:884:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/debug/dwarf/line.go:146:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/debug/dwarf/line.go:153:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/debug/dwarf/typeunit.go:138:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/debug/pe/file.go:470:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/net/http/h2_bundle.go:9530:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/net/http/transfer.go:765:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/net/http/transfer.go:775:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/net/http/transfer.go:798:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/go/build/build.go:1442:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/go/build/build.go:1453:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/go/build/build.go:1457:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/go/build/build.go:1491:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/go/internal/gccgoimporter/ar.go:125:3: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/image/jpeg/reader.go:622:5: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/image/png/reader.go:434:4: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/internal/profile/legacy_profile.go:1089:4: return both a `nil` error and an invalid value: use a sentinel error instead
/usr/local/go/src/net/internal/socktest/switch.go:142:3: return both a `nil` error and an invalid value: use a sentinel error instead
```

</details>
