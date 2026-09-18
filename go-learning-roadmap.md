# Go Development Learning Plan

## Learning Objective

Develop strong practical and conceptual expertise in Go, with the
ability to independently design, build, test, containerize, and deploy:

-   Command-line applications
-   Server-based applications
-   REST APIs
-   Concurrent applications
-   Dockerized applications
-   Applications running locally and on Linux servers

The learning journey will progress through clearly defined phases. Each
phase will introduce concepts, reinforce them through exercises, and
conclude with a practical deliverable.

------------------------------------------------------------------------

## Learning Approach

The learning process will follow these principles:

1.  Understand concepts before implementing them.
2.  Explore design decisions, constraints, and trade-offs.
3.  Study internal behavior when it is important for reliable
    development.
4.  Use diagrams and pseudocode before writing complete implementations.
5.  Practice through progressively more complex projects.
6.  Include testing, debugging, performance, and deployment throughout
    the journey.
7.  Complete each phase's learning objectives before progressing to the
    next phase.
8.  Introduce additional topics only when they support the current
    learning objective or are explicitly requested.

------------------------------------------------------------------------

# Phase 0 --- Go Philosophy and Programming Model

**Estimated duration:** 1--2 days

## Objectives

Understand Go's purpose, design philosophy, execution model, and core
programming principles.

## Topics

-   Why Go was created
-   Go's design philosophy
-   Compilation and execution model
-   Static typing
-   Garbage collection
-   Zero values
-   Interfaces
-   Explicit error handling
-   Concurrency model
-   Packages and modules
-   Standard library philosophy

## Key Areas of Understanding

-   How Go programs are compiled and executed
-   How memory is managed
-   Why zero values are important
-   How interfaces support abstraction
-   Why explicit error handling is central to Go
-   How goroutines and channels support concurrent programming
-   How packages and modules organize applications

## Deliverable

A concise technical overview covering Go's programming model, execution
flow, memory management, error handling, and concurrency principles.

------------------------------------------------------------------------

# Phase 1 --- Go Fundamentals

**Estimated duration:** 1 week

## Objectives

Build a strong foundation in Go syntax, data types, control flow,
functions, and core data structures.

## Topics

### Basic Syntax

-   `package main`
-   `func main()`
-   Variables and constants
-   Primitive types
-   Type inference
-   Type conversions
-   Zero values

### Control Flow

-   `if`
-   `for`
-   `switch`
-   `defer`

### Functions

-   Multiple return values
-   Named return values
-   Variadic functions
-   Anonymous functions
-   Closures

### Data Structures

-   Arrays
-   Slices
-   Maps
-   Structs

### Pointers

-   Pointer syntax
-   Value and pointer semantics
-   Mutability
-   Passing values to functions
-   Passing pointers to functions

## Practical Exercises

1.  Temperature converter
2.  Log file line counter
3.  Word frequency counter
4.  In-memory contact manager
5.  Configuration parser

## Key Questions

-   What is the difference between arrays and slices?
-   How does slice capacity work?
-   How do maps behave when a key is missing?
-   What are zero values?
-   When should a pointer be used?
-   How do functions return multiple values?

## Deliverable

A collection of small, tested Go programs demonstrating the fundamental
language features.

------------------------------------------------------------------------

# Phase 2 --- Go Core Programming Model

**Estimated duration:** 1 week

## Objectives

Learn idiomatic Go design, abstraction, package organization, and
structured error handling.

## Topics

### Structs and Methods

-   Value receivers
-   Pointer receivers
-   Method sets
-   Composition
-   Encapsulation through exported and unexported identifiers

### Interfaces

-   Implicit interface implementation
-   Interface values
-   Interface composition
-   Dependency inversion
-   Small interfaces
-   Avoiding unnecessary abstractions

### Error Handling

-   The `error` interface
-   Error wrapping
-   `errors.Is`
-   `errors.As`
-   Custom errors
-   Error propagation
-   Deciding when to return errors

### Packages

-   Package visibility
-   Exported and unexported identifiers
-   Package organization
-   Import cycles
-   Package responsibilities

### Modules

-   `go.mod`
-   `go.sum`
-   Dependency management
-   Module versioning
-   Dependency updates

## Mini-Project

### Configuration Management Library

Features:

-   Read configuration values
-   Validate configuration fields
-   Support environment variables
-   Return structured errors
-   Provide sensible defaults
-   Include unit tests
-   Include package documentation

## Deliverable

A reusable configuration package with tests, documentation, and clear
package boundaries.

------------------------------------------------------------------------

# Phase 3 --- CLI Application Development

**Estimated duration:** 1 week

## Objectives

Build a maintainable command-line application with structured commands,
configuration, persistence, and error handling.

## Topics

-   Command-line arguments
-   Standard input and output
-   File operations
-   JSON encoding and decoding
-   Configuration files
-   Exit codes
-   Logging
-   Command routing
-   CLI architecture
-   Building executable binaries

## Project: `gomanager`

A command-line application for managing local tasks or services.

## Example Commands

``` bash
gomanager init
gomanager config show
gomanager task create backup
gomanager task list
gomanager task delete backup
```

## Suggested Architecture

``` text
CLI
 |
Command Handler
 |
Application Service
 |
Repository
 |
Storage
```

## Required Capabilities

-   Command routing
-   Application configuration
-   Service layer
-   Repository abstraction
-   Error handling
-   Input validation
-   Unit testing
-   User-friendly command output

## Deliverable

A working CLI application with clean separation of responsibilities,
documented commands, and automated tests.

------------------------------------------------------------------------

# Phase 4 --- Concurrency and Parallel Programming

**Estimated duration:** 2 weeks

## Objectives

Understand Go's concurrency model and apply concurrency safely in
practical applications.

## Topics

### Goroutines

-   Goroutine lifecycle
-   Scheduling fundamentals
-   Goroutine coordination
-   Goroutine leaks

### Channels

-   Unbuffered channels
-   Buffered channels
-   Sending and receiving
-   Channel ownership
-   Channel direction
-   Closing channels
-   `select`

### Synchronization

-   `sync.Mutex`
-   `sync.RWMutex`
-   `sync.WaitGroup`
-   `sync.Once`
-   `sync.Cond`
-   Atomic operations

### Context

-   Cancellation
-   Deadlines
-   Request-scoped values
-   Propagating cancellation
-   Context ownership

### Concurrency Patterns

-   Worker pools
-   Fan-in and fan-out
-   Pipelines
-   Rate limiting
-   Graceful shutdown
-   Backpressure

### Debugging

-   Race detector
-   Goroutine leaks
-   Deadlocks
-   Data races
-   Shared state protection

## Project: Concurrent File Processor

The application must:

1.  Scan a directory
2.  Find matching files
3.  Process files concurrently
4.  Limit the number of workers
5.  Report progress
6.  Support cancellation
7.  Handle processing errors
8.  Produce a final summary

## Key Questions

-   When should a mutex be used instead of a channel?
-   How can a worker pool be cancelled safely?
-   How do goroutine leaks occur?
-   How can concurrent writes be protected?
-   How does the race detector identify problems?
-   How should errors from multiple workers be collected?

## Deliverable

A concurrent file-processing application with controlled worker
execution, cancellation support, error handling, tests, and race
detection.

------------------------------------------------------------------------

# Phase 5 --- Networking and HTTP Servers

**Estimated duration:** 2 weeks

## Objectives

Understand networking fundamentals and build reliable server-side
applications using Go.

## Part A --- Networking Fundamentals

-   TCP and UDP
-   Client-server communication
-   Sockets
-   HTTP fundamentals
-   Request and response lifecycle
-   Timeouts
-   Connection handling
-   Connection termination

## Part B --- Go Networking APIs

-   `net`
-   `net/http`
-   HTTP clients
-   HTTP servers
-   Middleware
-   JSON APIs
-   Context propagation
-   Graceful shutdown

## Project: Go REST API

Build a REST API with:

-   CRUD endpoints
-   Request validation
-   JSON serialization
-   Structured errors
-   Logging
-   Request timeouts
-   Graceful shutdown
-   Unit tests
-   Health endpoints

## Suggested Architecture

``` text
Client
 |
HTTP Server
 |
Router
 |
Middleware
 |
Handler
 |
Service
 |
Repository
 |
Database
```

## Learning Strategy

Begin with Go's standard library to understand HTTP server and client
behavior. Introduce additional libraries only when their benefits and
trade-offs are understood.

## Deliverable

A tested REST API with clear separation between transport, application
logic, and persistence.

------------------------------------------------------------------------

# Phase 6 --- Database Integration and Application Architecture

**Estimated duration:** 1--2 weeks

## Objectives

Build a database-backed application with reliable persistence and
maintainable architecture.

## Database Topics

-   `database/sql`
-   PostgreSQL connectivity
-   Connection pools
-   Transactions
-   Prepared statements
-   Query timeouts
-   Migration strategies
-   Transaction boundaries
-   Handling database errors

## Architecture Topics

-   Layered architecture
-   Dependency injection
-   Repository pattern
-   Service layer
-   Configuration management
-   Dependency boundaries
-   Package responsibilities
-   Managing application lifecycle

## API Topics

-   Authentication concepts
-   Authorization
-   Validation
-   Pagination
-   Idempotency
-   Consistent error responses
-   Resource ownership
-   API versioning fundamentals

## Project Extension

Extend the REST API to use PostgreSQL.

## Suggested Architecture

``` text
HTTP API
 |
Application Layer
 |
Domain Logic
 |
Repository Interface
 |
PostgreSQL Implementation
```

## Design Focus

Evaluate where abstractions improve maintainability and where they
introduce unnecessary complexity.

## Deliverable

A PostgreSQL-backed REST API with transactions, validation, database
error handling, and integration tests.

------------------------------------------------------------------------

# Phase 7 --- Testing and Production Practices

**Estimated duration:** 1 week

## Objectives

Learn how to validate, debug, measure, and maintain Go applications.

## Topics

-   Go testing package
-   Table-driven tests
-   Test fixtures
-   Mocks and fakes
-   HTTP testing
-   Integration tests
-   Benchmarks
-   Profiling
-   Race detection
-   Static analysis
-   Code formatting
-   Documentation
-   Test organization
-   Test reliability

## Tools

``` bash
go test ./...
go test -race ./...
go test -bench=.
go vet ./...
gofmt
```

## Testing Strategy

``` text
Unit Tests
 |
Service Tests
 |
HTTP Handler Tests
 |
Integration Tests
```

## Quality Areas

-   Correctness
-   Error handling
-   Concurrency safety
-   Performance
-   Maintainability
-   Observability
-   Reproducibility

## Deliverable

A complete testing and quality strategy applied to the REST API,
including unit tests, integration tests, race detection, and basic
performance analysis.

------------------------------------------------------------------------

# Phase 8 --- Docker and Deployment

**Estimated duration:** 1 week

## Objectives

Build, package, and deploy Go applications locally and on Linux servers.

## Local Execution

-   Build Go binaries
-   Environment configuration
-   Cross-compilation
-   Linux executable deployment
-   Application configuration at runtime
-   Process lifecycle management

## Docker Topics

-   Dockerfile
-   Multi-stage builds
-   Minimal runtime images
-   Container configuration
-   Volumes
-   Networking
-   Health checks
-   Container logs
-   Image tagging

## Deployment Topics

-   Deploying to a Linux server
-   Docker Compose
-   Configuration management
-   Logging
-   Restart policies
-   Graceful shutdown
-   Service health monitoring
-   Basic deployment troubleshooting

## Deployment Flow

``` text
Developer Machine
 |
Go Application
 |
Docker Image
 |
Container Registry
 |
Linux Server
 |
Running Application
```

## Project

Deploy the Go REST API with:

-   Go API
-   PostgreSQL
-   Docker Compose
-   Environment configuration
-   Health checks
-   Persistent database storage
-   Application logs

Understand both:

-   Running the application directly on Linux
-   Running the application inside a container

## Deliverable

A reproducible local and server deployment process with documentation
and operational instructions.

------------------------------------------------------------------------

# Phase 9 --- Capstone Project

**Estimated duration:** 2--3 weeks

## Project: Go Service Manager

A production-oriented application containing both CLI and server
components.

## Proposed Functionality

``` text
Go Service Manager
|
+-- CLI
|   +-- Service registration
|   +-- Service status
|   +-- Configuration
|   +-- Health checks
|
+-- HTTP API
|   +-- Service CRUD
|   +-- Status endpoints
|   +-- Health checks
|   +-- Graceful shutdown
|
+-- Storage
|   +-- PostgreSQL
|
+-- Concurrency
|   +-- Background workers
|
+-- Testing
|   +-- Unit tests
|   +-- Integration tests
|
+-- Deployment
    +-- Local binary
    +-- Docker
```

## Skills Combined

-   Project structure
-   Interfaces
-   Concurrency
-   Networking
-   Database integration
-   Testing
-   Docker
-   Deployment
-   Documentation
-   Operational troubleshooting

## Deliverable

A documented, tested, Dockerized Go application with both CLI and HTTP
interfaces.

------------------------------------------------------------------------

# Standard Lesson Format

Each lesson should follow this sequence.

## 1. Conceptual Foundation

Explain the concept clearly and deeply, including why it exists and when
it should be used.

## 2. Internal Behavior

Explain runtime, memory, compiler, concurrency, networking, or
operating-system behavior when relevant.

## 3. Design Discussion

Discuss constraints, trade-offs, alternatives, and common design
mistakes.

## 4. Pseudocode

Present the solution structure in pseudocode before implementation.

## 5. Implementation

Provide complete code when implementation is explicitly requested or
when the learning exercise requires it.

## 6. Testing and Debugging

Explain how to validate the implementation and diagnose failures.

## 7. Understanding Check

Use targeted questions, short exercises, or assessments before
progressing to the next concept.

------------------------------------------------------------------------

# Progression Rules

Follow the phases in order:

``` text
Phase 0: Go Philosophy
    |
Phase 1: Fundamentals
    |
Phase 2: Core Programming Model
    |
Phase 3: CLI Development
    |
Phase 4: Concurrency
    |
Phase 5: Networking and HTTP
    |
Phase 6: Databases and Architecture
    |
Phase 7: Testing and Production
    |
Phase 8: Docker and Deployment
    |
Phase 9: Capstone
```

Do not skip foundational concepts without confirming understanding.

Do not introduce unrelated advanced topics prematurely.

Adjust the depth and pace of a phase based on demonstrated
understanding, practical exercises, and assessment results.

------------------------------------------------------------------------

# Starting Point

Begin with:

## Phase 0 --- Go Philosophy and Programming Model

The first lesson should cover:

1.  Why Go was created
2.  The problems Go aims to solve
3.  How Go programs compile and execute
4.  Go's memory management model
5.  Interfaces and composition
6.  Error handling principles
7.  Go's concurrency model
8.  The trade-offs of using Go
9.  The role of packages, modules, and the standard library
