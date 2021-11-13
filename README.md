# Lru-cache

This library provides the ability to use information caching using the least recently used mechanism. It is based on doubly linked lists and hash tables.

You can create a cache with the desired capacity, and use the provided methods to add a new value to the cache, call it, or delete it by key.
This cache implementation is thread safe out of the box, so you don't have to worry about it. 