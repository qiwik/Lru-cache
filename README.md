# golru

The library provides work with a lru cache based on two-linked lists and a hash table. You just need to initialize 
the cache and start working with it through the provided methods.

Mutex is built into the cache, so it is safe from the point of view of competitive access. It is accepted that 
keys should be strings that can match any data type.

At the moment, the project is completely covered with tests.

## Import

````
go get github.com/qiwik/golru
````

## Example
You can create a cache instance in the following way:
````
cache, err := golru.NewCache(20)
if err != nil {
    log.Fatalf("can't create cache: %v", err)
}
````

And then start working with this structure:
````
ok := cache.Add("test", 42)
if !ok {
    log.Fatal("key with value weren't added")
}

value, ok := cache.Get("test")
if !ok {
    log.Fatal("this key doesn't exist")
}
````

## Comparison

According to benchmarks, the glory library is about 1.8-2 times faster than the existing library 
https://github.com/hashicorp/golang-lru
````
goos: linux
goarch: amd64
cpu: Intel(R) Core(TM) i5-4200M CPU @ 2.50GHz
BenchmarkReflectKeys-4            242337              6026 ns/op            3072 B/op         53 allocs/op
BenchmarkKeys-4                  1000000              1330 ns/op             896 B/op          1 allocs/op
BenchmarkGolangLruKeys-4         2102901               569.2 ns/op           896 B/op          1 allocs/op
BenchmarkGolangLruGet-4         19902330                54.83 ns/op            0 B/op          0 allocs/op
BenchmarkGet-4                  36394338                30.68 ns/op            0 B/op          0 allocs/op
BenchmarkGolangLruAdd-4          6463861               183.7 ns/op            96 B/op          2 allocs/op
BenchmarkAdd-4                   7743319               151.0 ns/op            96 B/op          2 allocs/op
BenchmarkGolangLruRemove-4      23369701                51.40 ns/op            0 B/op          0 allocs/op
BenchmarkRemove-4               31637878                31.88 ns/op            0 B/op          0 allocs/op

````