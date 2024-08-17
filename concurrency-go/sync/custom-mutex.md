### Introduction:

To implement our readers–writer mutex, we need a system that, when a goroutine calls `ReadLock()`, blocks any access to 
the write part while allowing other goroutines to still call `ReadLock()` without blocking. We’ll block the write part 
by making sure that a goroutine calling `WriteLock()` suspends execution. Only when all the read goroutines call 
`ReadUnlock()` will we allow another goroutine to unblock from `WriteLock()`

To help us visualize this system, we can think of goroutines as entities trying to access a room with two entrances. 
This room signifies access to a shared resource. The reader goroutines use a specific entrance, and the writers use another.
Entrances only admit one goroutine at a time, although multiple goroutines can be in the room at the same time. 
We keep a counter that a reader goroutine increments by 1 when it enters via the reader’s entrance and reduces by 1 
when it leaves the room. The writer’s entrance can be locked from the inside using what we call a global lock.

The procedure is that when the first reader goroutine enters the room, it must lock the writers’ entrance. This ensures 
that a writer goroutine will find access impassable, blocking the goroutine’s execution. However, other reader 
goroutines will still have access through their own entrance. The reader goroutine knows that it’s the first one in 
the room because the counter has a value of 1.

The writer’s entrance here is just another mutex lock that we call the global lock. A writer needs to acquire this mutex
in order to hold the writer’s part of the readers-writer lock. When the first reader locks this mutex, it blocks any 
goroutine requesting the writer’s part of the lock.

We need to make sure that only one goroutine is using the readers’ entrance at any time because we don’t want two 
simultaneous read goroutines to enter at the same time and believe they are both the first in the room. 
This would result in both trying to lock the global lock and only one succeeding. Thus, to synchronize access so only 
one goroutine can use the readers’ entrance at any time, we can make use of another mutex. We’ll call this mutex 
`readersLock`. The readers counter is represented by the `readersCounter` variable, and we’ll call the writer’s lock 
`globalLock`.

```go
type ReadWriteMutex struct {
    readersCounter int
    readersLock    sync.Mutex
    globalLock     sync.Mutex
}
```

* Integer variable to count the number of reader goroutines currently in the critical section
* Mutex for synchronizing readers access
* Mutex for blocking any writers access

```go
func (rw *ReadWriteMutex) ReadLock() {
    rw.readersLock.Lock()
    rw.readersCounter++
    if rw.readersCounter == 1 {
        rw.globalLock.Lock()
    }
    rw.readersLock.Unlock()
}
```
* Synchronizes access so that only one goroutine is allowed at any time
* Reader goroutine increments readersCounter by 1
* If a reader goroutine is the first one in, it attempts to lock globalLock.
* Synchronizes access so that only one goroutine is allowed at any time
* Any writer access requires a lock on globalLock.

> Once the caller gets hold of the readersLock, it increments the readers’ counter by 1, signifying that another goroutine is about to have read access to the shared resource. If the goroutine realizes that it’s the first one to get read access, it tries to lock the globalLock so that it blocks access to any write goroutines. (The globalLock is used by the WriteLock() function when it needs to obtain the writer’s side of this mutex.) If the globalLock is free, it means that no writer is currently executing its critical section. In this case, the first reader obtains the globalLock, releases the readersLock, and goes ahead to execute its reader’s critical section.
>

> When a reader goroutine finishes executing its critical section, we can think of it as exiting through the same passageway. On its way out, it decreases the counter by 1. Using the same passageway simply means that it needs to get hold of the readersLock when updating the counter. The last one out of the room (when the counter is 0), unlocks the global lock so that a writer goroutine can finally access the shared resource
>

> While a writer goroutine is executing its critical section, accessing the room in our analogy, it holds a lock on the globalLock. This has two effects. First, it blocks other writers’ goroutines since writers need to acquire this lock before gaining access. Second, it also blocks the first reader goroutine when it tries to acquire the globalLock. The first reader goroutine will block and wait until the globalLock becomes available. Since the first reader goroutine also holds the readersLock, it will also block access to any other reader goroutine that follows while it waits. This is akin to the first reader goroutine not moving and thus blocking the readers’ entrance, not letting any other goroutines in.
Once the writer goroutine has finished executing its critical section, it releases the globalLock. This has the effect of unblocking the first reader goroutine and later allowing in any other blocked readers.
>

> We can implement this releasing logic in our two unlocking functions. Listing 4.14 shows both the ReadUnlock() and WriteUnlock() functions. ReadUnlock() again uses the readersLock to ensure that only one goroutine is executing this function at a time, protecting the shared readersCounter variable. Once the reader acquires the lock, it decrements the readersCounter count by 1, and if the count reaches 0, it also releases the globalLock. This allows the possibility of a writer gaining access. On the writer’s side, WriteUnlock() simply releases the globalLock, giving either readers or a single writer access.”
>

```go
func (rw *ReadWriteMutex) ReadUnlock() {
    rw.readersLock.Lock()
    rw.readersCounter--
    if rw.readersCounter == 0 {
        rw.globalLock.Unlock()
    }
    rw.readersLock.Unlock()
}

func (rw *ReadWriteMutex) WriteLock() {
    rw.globalLock.Lock()
}

func (rw *ReadWriteMutex) WriteUnlock() {
    rw.globalLock.Unlock()
}

```

* Synchronizes access so that only one goroutine is allowed at any time
* The reader goroutine decrements readersCounter by 1.
* If the reader goroutine is the last one out, it unlocks the global lock.
* Synchronizes access so that only one goroutine is allowed at any time
* The writer goroutine, finishing its critical section, releases the global lock.

This implementation of the readers–writer lock is read-preferring. This means that if we have a consistent number of 
readers goroutines hogging the read part of the mutex, a writer goroutine would be unable to acquire the mutex. 
In technical terms, we say that the reader goroutines are starving the writer ones, not allowing them access to the 
shared resource.

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	rwMutex := ReadWriteMutex{}
	for i := 0; i < 2; i++ {
		go func() {
			for {
				rwMutex.ReadLock()
				time.Sleep(1 * time.Second)
				fmt.Println("Read done")
				rwMutex.ReadUnlock()
			}
		}()
	}
	time.Sleep(1 * time.Second)
	rwMutex.WriteLock()
	fmt.Println("Write finished")
}
```

Even though we have an infinite loop in our goroutines, we expect that eventually the `main()` goroutine will acquire 
a hold on the writer’s lock, output the message Write finished, and terminate. This should happen because in Go, 
whenever the `main()` goroutine terminates, the entire process exits. However, when we run, Our two goroutines 
constantly hold the reader part of our mutex, which prevents our `main()` goroutine from ever acquiring the writer’s 
part of the lock. If we are lucky, the readers might release the readers lock at the same time, enabling the writer 
goroutine to acquire it. However, in practice, it is unlikely that both reader threads will release the lock at the 
same time. This leads to the writer-starvation of our `main()` goroutine.

Starvation is a situation where an execution is blocked from gaining access to a shared resource because the 
resource is made unavailable for a long time (or indefinitely) by other greedy executions.

### Improved Custom RWMutex:
We need a different design for a readers–writer lock that is not read-preferred—one that doesn’t starve our writer 
goroutines. We could block new readers from acquiring the read lock as soon as a writer calls the `WriteLock()` function.
To achieve this, instead of having the goroutines block on a mutex, we could have them suspended using a condition 
variable. With a condition variable, we can have different conditions on when to block readers and writers. 
To design a write-preferred lock, we need a few properties

1. Readers’ counter—Initially set to 0, this tells us how many reader goroutines are actively accessing the shared resources.
2. Writers’ waiting counter—Initially set to 0, this tells us how many writer goroutines are suspended waiting to access the shared resource.
3. Writer active indicator—Initially set to false, this flag tells us if the resource is currently being updated by a writer goroutine.
4. Condition variable with mutex—This allows us to set various conditions on the preceding properties, suspending execution when the conditions aren’t met.

#### Go's RWMutex:
The RWMutex bundled with Go is write-preferring. This is highlighted in Go’s documentation 
(from https://pkg.go.dev/sync#RWMutex; calling Lock() acquires the writer’s part of the mutex)
If a goroutine holds a RWMutex for reading and another goroutine might call Lock, no goroutine should expect to be 
able to acquire a read lock until the initial read lock is released. In particular, this prohibits recursive read 
locking. This is to ensure that the lock eventually becomes available; a blocked Lock call excludes new readers from 
acquiring the lock.

#### Solution:
1. Readers can access the shared resource when no writers are active or waiting.
2. We block writers from accessing the shared resource when readers or a writer are using it
3. We also block new readers when writers are waiting.

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type ReadWriteMutex struct {
    readersCounter int
    writersWaiting int
    writerActive   bool
    cond           *sync.Cond
}

func NewReadWriteMutex() *ReadWriteMutex {
    return &ReadWriteMutex{cond: sync.NewCond(&sync.Mutex{})}
}

func (rw *ReadWriteMutex) ReadLock() {
    rw.cond.L.Lock()
    for rw.writersWaiting > 0 || rw.writerActive {
        rw.cond.Wait()
    }
    rw.readersCounter++
    rw.cond.L.Unlock()
}

func (rw *ReadWriteMutex) WriteLock() {
    rw.cond.L.Lock()
    rw.writersWaiting++
    for rw.readersCounter > 0 || rw.writerActive {
        rw.cond.Wait()
    }
    rw.writersWaiting--
    rw.writerActive = true
    rw.cond.L.Unlock()
}

func (rw *ReadWriteMutex) ReadUnlock() {
    rw.cond.L.Lock()
    rw.readersCounter--
    if rw.readersCounter == 0 {
        rw.cond.Broadcast()
    }
    rw.cond.L.Unlock()
}

func (rw *ReadWriteMutex) WriteUnlock() {
    rw.cond.L.Lock()
    rw.writerActive = false
    rw.cond.Broadcast()
    rw.cond.L.Unlock()
}

func main() {
    rwMutex := NewReadWriteMutex()
    for i := 0; i < 2; i++ {
        go func() {
            for {
                rwMutex.ReadLock()
                time.Sleep(1 * time.Second)
                fmt.Println("Read done")
                rwMutex.ReadUnlock()
            }
        }()
    }
    time.Sleep(1 * time.Second)
    rwMutex.WriteLock()
    fmt.Println("Write finished")
}
```