# enterprise
Atılgan aka Enterprise is a new ORM for postgresql databases. For speed and lightness.

## Features
- A small number of auto generated database model structs and methods
- Create, Get, Update, Delete and List functions on table
- It can load Relation tables and filter them with single query
- Developer friendly helper methods on models
- Simple aggregate functionality
- Simple transaction commit and rollback
- Customizable hooks
- Atlas migrations
- Unique Relation system
- Fastest ORM in golang for postgresql


## Generating Models
First create a package. Create one file for every table. The function name and file name should match.
After write your models. We recommend to create a new package with name generate and put a generate.go for generating models you want.
It will go to create a models package and write files in to it. For examples [Enterprise Example](https://github.com/MrSametBurgazoglu/enterprise_example).

### Migration
For migration we recommend to create a migrate named package and put this script into it.
Execute the script and it will create a directory called migrations and write raw sql for migration.
Further questions please look to Atlas.

### RoadMap

Version 1.0.0 Roadmap

V0.1.0-alpha

- [x] Add simple unit tests
- [x] Hook system functions to every model db operation.
- [ ] Pk and fk on migration table

V0.1.0
- [ ] Add bulks(Insert, Update, Delete)

V0.2.0
- [ ] Choose Join Type

V0.3.0
- [ ] Json field
- [ ] []byte field

V0.4.0
- [ ] Use different clients on Read and Write db operations

V0.5.0
- [ ] Custom go type on DB

V0.6.0
- [ ] Nested Transactions, Save Point, RollbackTo to Saved Point

V0.7.0
- [ ] Composite primary key

V0.8.0
- [ ] Nice and beautiful Debug Mode with logger

V0.9.0
- [ ] %100 covered unit tests
- [ ] Fully Documentation
- [ ] Github Actions
- [ ] Add, Update and Delete Constraints

V1.0.0


### Benchmark Results
````text
goos: linux
goarch: amd64
pkg: github.com/FournyP/go-orm-benchmarks/benchmarks
cpu: 12th Gen Intel(R) Core(TM) i7-1255U

BenchmarkEntCreate-4                                 105          11393218 ns/op            7571 B/op        218 allocs/op
BenchmarkEnterpriseCreate-4                          175           6458168 ns/op            9264 B/op        193 allocs/op
BenchmarkGORMCreate-4                                138           8334585 ns/op           22754 B/op        334 allocs/op
BenchmarkSqlxCreate-4                                139           8516631 ns/op            3088 B/op         87 allocs/op

BenchmarkEntUpdate-4                                 100          10051190 ns/op            5302 B/op        137 allocs/op
BenchmarkEnterpriseUpdate-4                          507           2359239 ns/op            1652 B/op         39 allocs/op
BenchmarkGORMUpdate-4                                468           2712004 ns/op            6657 B/op         86 allocs/op
BenchmarkSqlxUpdate-4                                502           2518410 ns/op             535 B/op         16 allocs/op

BenchmarkEntDelete-4                                 367           2897743 ns/op            1904 B/op         45 allocs/op
BenchmarkEnterpriseDelete-4                          513           2089324 ns/op            1157 B/op         21 allocs/op
BenchmarkGORMDelete-4                                492           3073590 ns/op            5634 B/op         85 allocs/op
BenchmarkSqlxDelete-4                                513           2385171 ns/op             304 B/op          9 allocs/op

BenchmarkEntRead-4                                  1107           1179827 ns/op            3744 B/op         93 allocs/op
BenchmarkEnterpriseRead-4                           2575            428796 ns/op            3005 B/op         73 allocs/op
BenchmarkGORMRead-4                                 2296            519182 ns/op            5406 B/op         94 allocs/op
BenchmarkSqlxRead-4                                 1273            943025 ns/op            1200 B/op         32 allocs/op

BenchmarkEntReadWithRelations-4                      424           3185571 ns/op           11888 B/op        295 allocs/op
BenchmarkEnterpriseReadWithRelations-4              2185            544165 ns/op           10781 B/op        219 allocs/op
BenchmarkGORMReadWithRelations-4                     692           2137871 ns/op           34967 B/op        372 allocs/op
BenchmarkSqlxReadWithRelations-4                     456           2687558 ns/op            3856 B/op        105 allocs/op

BenchmarkEntReadSingleField-4                       1153            973295 ns/op            2776 B/op         73 allocs/op
BenchmarkEnterpriseReadSingleField-4                2908            494583 ns/op            2686 B/op         62 allocs/op
BenchmarkGormReadSingleField-4                      2036            515719 ns/op            4738 B/op         73 allocs/op
BenchmarkSqlxReadSingleField-4                      1242           1015917 ns/op             792 B/op         22 allocs/op
````