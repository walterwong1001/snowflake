Snowflake
=========

Snowflake is a distributed unique ID generator inspired by [Twitter's Snowflake](https://blog.twitter.com/2010/announcing-snowflake).

A Snowflake ID is composed of

    41 bits for time in units of 1 msec
    10 bits for a work id
    12 bits for a sequence number

The start time of the Snowflake is set to "2024-07-15 00:00:00 +0000 UTC".

Installation
------------
```
go get github.com/walterwong1001/snowflake
```

Usage
-----
```go
s, err := snowflake.NewSnowflake(1)
if err != nil {
    log.Println(err)
    return
}
s.NextID()
```