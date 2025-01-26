package logs

import "wow/pkg/logger"

var baseLogsPath = "./wow.logs"

var MainLogger = logger.NewWithKeeping(baseLogsPath).WithPrefix("main")

var ServerLogger = logger.NewWithKeeping(baseLogsPath).WithPrefix("server")

var TransportLogger = logger.NewWithKeeping(baseLogsPath).WithPrefix("transport")

var EpollLogger = logger.NewWithKeeping(baseLogsPath).WithPrefix("transport").WithSecondPrefix("epoll")

var WorkersLogger = logger.NewWithKeeping(baseLogsPath).WithPrefix("transport").WithSecondPrefix("workers")

var WorkersPoolLogger = logger.NewWithKeeping(baseLogsPath).WithPrefix("transport").WithSecondPrefix("workers").WithThirdPrefix("pool")

var WorkersWorkerLogger = logger.NewWithKeeping(baseLogsPath).WithPrefix("transport").WithSecondPrefix("workers").WithThirdPrefix("worker")

var WorkersWriterLogger = logger.NewWithKeeping(baseLogsPath).WithPrefix("transport").WithSecondPrefix("workers").WithThirdPrefix("writer")
