package logs

import "wow/pkg/logger"

var baseLogsPath = "./.wow.logs"

var MainClientLogger = logger.NewWithKeeping(baseLogsPath).WithPrefix("main").WithSecondPrefix("client")

var MainServerLogger = logger.NewWithKeeping(baseLogsPath).WithPrefix("main").WithSecondPrefix("server")

var ServerLogger = logger.NewWithKeeping(baseLogsPath).WithPrefix("server")

var ClientLogger = logger.NewWithKeeping(baseLogsPath).WithPrefix("client")

var TransportLogger = logger.NewWithKeeping(baseLogsPath).WithPrefix("transport")

var EpollLogger = logger.NewWithKeeping(baseLogsPath).WithPrefix("transport").WithSecondPrefix("epoll")

var WorkersLogger = logger.NewWithKeeping(baseLogsPath).WithPrefix("transport").WithSecondPrefix("workers")

var WorkersPoolLogger = logger.NewWithKeeping(baseLogsPath).WithPrefix("transport").WithSecondPrefix("workers").WithThirdPrefix("pool")

var WorkersWorkerLogger = logger.NewWithKeeping(baseLogsPath).WithPrefix("transport").WithSecondPrefix("workers").WithThirdPrefix("worker")

var WorkersWriterLogger = logger.NewWithKeeping(baseLogsPath).WithPrefix("transport").WithSecondPrefix("workers").WithThirdPrefix("writer")
