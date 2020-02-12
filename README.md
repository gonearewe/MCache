# MCache
A simple cache system.

# Dependencies

In the persistence section, we use RockSdb to provide supports. It's a library developed by Facebook
that provides an embeddable, persistent key-value store for fast storage. You may learn more by clicking
[here](https://github.com/facebook/rocksdb).

Download and build this library by

> git clone https://github.com/facebook/rocksdb.git
>
> cd rocksdb && make static_library

Then you'll find a 'librocksdb.a' in current directory. Move it to rocksdb/ in the project root
directory so that build will succeed.

Actually, make and g++ are required to build librocksdb, but I bet you've already got them. 

Also, extra libraries libsnappy and libz are required. You may easily get them through

> apt install libz libsnappy
