# MCache

A simple cache system.

I recorded my learning process here on my way to build a cache system,
it ought to be a distributed one but I gave up implementing that part. 
My learning reference lies [here](https://github.com/stuarthu/go-implement-your-cache-server).

For now, it is ready for using with both in memory and persistence options.
For persistence, I made use of Rocksdb.

`This system runs on TCP`, you may find the protocol in the comments in

> server/tcp/read.go

`A simple TCP client application depending on this protocol` can be found in 

> client/cli/cli.go

It can handle Set/Get requests and may be useful when debugging the cache server.
`Benchmark package also provides a performance tool depending on Client package.`
It's a CLI application providing multiple options and usage details can be found in the source file. 
It also supports redis benchmark for comparision if you have redis-server installed; If not,

> apt install redis-server

`MCache supports GC with FIFO strategy.` Each KV-Cache has a limited life time specified by users,
and it will be deleted after its life circle. There are other strategies for cache GC. For example,
if you extend the cache's life time each time you access it by Get, it's called LRU
(Least Recently Used). However, updating cache requires acquiring LOCK when Get instead of just
RLOCK, which means a big performance loss when caches are frequently Get and rarely Set. 
If you records cache's usage count instead of expiration time, and update the count each time
this cache is used while deleting caches having less count during GC, it's called LFU
(Least Frequently Used). It causes similar issues as LRU strategy. So for now, FIFO is
simple, convenient and works just fine enough. `Of course, you may turn off GC function by
giving TTL a zero value(or negative value, anyway).`

# Dependencies

`In the persistence section, we use RockSdb to provide supports.` It's a library developed by Facebook
that provides an embeddable, persistent key-value store for fast storage. You may learn more by clicking
[here](https://github.com/facebook/rocksdb).

Download and build this library by

> git clone https://github.com/facebook/rocksdb.git
>
> cd rocksdb && make static_library

Then you'll find a `'librocksdb.a'` in current directory. **Move it to rocksdb/ in the project root
directory so that build will succeed.**

Actually, make and g++ are required to build librocksdb, but I bet you've already got them. 

Also, `extra libraries libsnappy and libz are required`. You may easily get them through

> apt install libz libsnappy
