1. ws one million connection
   1.1 open file limit: ulimit NOFILE
   1.2 reduce goroutine: epoll
   1.3 reduce buffer allocation: gobwas
   1.4 Conntrack table: concurrent in the OS

2. ws to millions connection
   1.1 get rid of reader goroutine: [netpoll]("github.com/mailru/easygo/netpoll")
   1.2 reuse goroutine: [gopool](https://github.com/gobwas/ws-examples/tree/master/src/gopool)
   1.3 zero copy upgrade
   1.4 library [gobwas](https://github.com/gobwas/ws)

3. my websocket library

   1. research with net package
      1. notice \r\n
   2. implement with http package
      1. hijack

4. chatapp implementation

5. monitor my chatapp

6. chat app load test
