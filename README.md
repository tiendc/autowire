[![Go Version][gover-img]][gover] [![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov] [![GoReport][rpt-img]][rpt]

# Dependency injection for Go 1.18+ using Generics and reflection

## Functionalities

- Automatically wiring/injecting objects on creation.
- Easy to use (see Usage section).
- `Shared mode` by default that creates only one instance for a type (this option is configurable).
- Ability to overwrite objects of specific types on building which is convenient for unit testing.
- No code generation.

## Installation

```shell
go get github.com/tiendc/autowire
```

## Usage

### General usage

```go
    // Suppose you have ServiceA and a creator function
    type ServiceA interface {}
    func NewServiceA(srvB ServiceB, repoX RepoX) ServiceA {
        return <ServiceA-instance>
    }

    // and ServiceB
    type ServiceB interface {}
    func NewServiceB(ctx context.Context, repoX RepoX, repoY RepoY) (ServiceB, error) {
        return <ServiceB-instance>, nil
    }

    // and RepoX
    type RepoX interface {}
    func NewRepoX(s3Client S3Client) (RepoX, error) {
        return <RepoX-instance>, nil
    }

    // and RepoY
    type RepoY interface {}
    func NewRepoY(redisClient RedisClient) RepoY {
        return <RepoY-instance>
    }

    // and a struct provider
    type ProviderStruct struct {
        RedisClient RedisClient
        S3Client    S3Client
    }

    var (
        // init the struct value
        providerStruct = &ProviderStruct{...}

        // create a container with passing providers
        container = MustNewContainer([]any{
            // Service providers
            NewServiceA,
            NewServiceB,
            // Repo providers
            NewRepoX,
            NewRepoY,
            // Struct provider (must be a pointer)
            providerStruct,
        })
    )

    func main() {
        // Update some values of the struct provider
        providerStruct.RedisClient = newRedisClient()
        providerStruct.S3Client = newS3Client()

        // Create ServiceA
        serviceA, err := autowire.BuildWithCtx[ServiceA](ctx, container)
        // Create RepoX
        repoX, err := autowire.Build[RepoX](container)
    }
```

### Non-shared mode

```go
    // Set sharedMode when create a container
    container = MustNewContainer([]any{
        // your providers
    }, SetSharedMode(false))

    // Activate non-shared mode inline for a specific build only
    serviceA, err := Build[ServiceA](container, NonSharedMode())
```

### Overwrite values of specific types

This is convenient in unit testing to overwrite specific types only.

```go
    // In unit testing, you may want to overwrite some `repos` and `clients` with fake instances.
    // NOTE: you may need to use `non-shared mode` to make sure cached objects are not used.
    serviceA, err := Build[ServiceA](container, NonSharedMode(),
            ProviderOverwrite[RepoX](fakeRepoX),
            ProviderOverwrite[RepoY](fakeRepoY))
    serviceA, err := Build[ServiceA](container, NonSharedMode(),
            ProviderOverwrite[RedisClient](fakeRedisClient),
            ProviderOverwrite[S3Client](fakeS3Client))
```

### Reclaim memory after use

Typically, dependency injection is only used at the initialization phase of a program.
However, it can take some space in memory which will become wasted after the phase.
Use the below trick to reclaim the memory taken by the autowire variables.

```go
    var (
        // create a global container
        container = MustNewContainer(...)
    )

    func autowireCleanUp() {
        // Assign nil to allow the garbage collector to reclaim the taken memory
        container = nil
    }

    func main() {
        // Initialization phase
        // Create services and use them
        serviceA, _ := autowire.Build[ServiceA](container)
        serviceB, _ := autowire.Build[ServiceB](container)

        // Clean up the usage, make sure you won't use the vars any more
        autowireCleanUp()
    }
```

## Contributing

- You are welcome to make pull requests for new functions and bug fixes.

## License

- [MIT License](LICENSE)

[doc-img]: https://pkg.go.dev/badge/github.com/tiendc/autowire
[doc]: https://pkg.go.dev/github.com/tiendc/autowire
[gover-img]: https://img.shields.io/badge/Go-%3E%3D%201.18-blue
[gover]: https://img.shields.io/badge/Go-%3E%3D%201.18-blue
[ci-img]: https://github.com/tiendc/autowire/actions/workflows/go.yml/badge.svg
[ci]: https://github.com/tiendc/autowire/actions/workflows/go.yml
[cov-img]: https://codecov.io/gh/tiendc/autowire/branch/main/graph/badge.svg
[cov]: https://codecov.io/gh/tiendc/autowire
[rpt-img]: https://goreportcard.com/badge/github.com/tiendc/autowire
[rpt]: https://goreportcard.com/report/github.com/tiendc/autowire
