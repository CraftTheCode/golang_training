# Module 05: Concurrency and Parallel Programming Deep-Dive

## 1. Conceptual Foundation

"Do not communicate by sharing memory; instead, share memory by communicating." — Rob Pike

Go's concurrency is founded on Tony Hoare's **Communicating Sequential Processes (CSP)**. In traditional languages (C++, Java, Python), concurrency often relies on operating system threads and shared memory locked by mutexes. In Go:
- **Goroutines** are extremely lightweight user-space threads managed entirely by the Go runtime.
- A goroutine starts with only **2 KB of stack space** (compared to 1 MB to 8 MB for an OS thread), growing and shrinking dynamically on the heap as needed.
- **Channels** provide safe synchronization and data transfer without manual lock contention.

---

## 2. Internal Behavior: The GMP Scheduler

The Go runtime uses an $M:N$ multiplexing scheduler known as the **GMP Model**:

```text
[ G1 ] [ G2 ] [ G3 ] (Global Run Queue)
         |
         v
     +-------+                    +-------+
     |   P0  |                    |   P1  |  <--- P: Logical Processor (GOMAXPROCS)
     +-------+                    +-------+       Holds local run-queue of Goroutines (G)
     |   M0  |                    |   M1  |  <--- M: Machine (OS Thread)
     +-------+                    +-------+
         |                            |
  [ OS Thread 1 ]              [ OS Thread 2 ]
```

1. **G (Goroutine)**: Represents the goroutine, its stack, program counter, and state (runnable, running, waiting).
2. **M (Machine / OS Thread)**: Physical OS thread created by the OS kernel.
3. **P (Processor)**: Logical context representing execution resources. By default, `P = runtime.NumCPU()`.

**Work Stealing & Syscall Preemption**:
- If a processor's local queue is empty, it steals half the goroutines from another processor's queue.
- If a goroutine makes a blocking system call (e.g. reading from disk or network), the runtime disassociates the thread $M$ from $P$, creating or assigning a new thread $M'$ to keep $P$ executing other runnable goroutines.

---

## 3. Design Discussion: Channels vs. Mutexes

| Feature | Channels | Mutexes (`sync.Mutex`, `sync.RWMutex`) |
| :--- | :--- | :--- |
| **Best For** | Passing data ownership, event distribution, pipelines, coordination | Protecting internal struct state, counters, in-memory caches |
| **Mental Model** | Conveyor belt / message queue | Door lock on shared room |
| **Complexity** | Can cause deadlocks or leaks if unclosed or unbuffered improperly | Prone to race conditions if lock is forgotten or unlocked late |

---

## 4. Module Directory Structure

```text
05-concurrency-deep-dive/
├── 01_goroutines_and_channels/main.go     # Unbuffered vs buffered, select, directions
├── 02_sync_primitives/main.go             # Mutex, RWMutex, WaitGroup, Once, Atomic
├── 03_context_and_cancellation/main.go    # context.WithTimeout, cancellation cascades
├── 04_concurrency_patterns/main.go        # Worker pool, Fan-out/Fan-in, Pipeline
├── 05_race_detector_and_leaks/main.go     # Race demo, fix, and leak prevention
└── project-file-processor/                # Concurrent File Processor
    ├── processor.go                       # Bounded worker pool, progress reporting
    ├── processor_test.go                  # Race-safe unit & stress tests
    └── README.md                          # Usage instructions
```
