package com.mathservice.infrastructure

import cats.effect.{ExitCode, IO, IOApp, Resource}
import com.mathservice.application.VolatilityService
import com.mathservice.grpc.math_service.{MathServiceGrpc, VolatilityRequest, VolatilityResponse}
import com.typesafe.scalalogging.LazyLogging
import io.grpc.{Server, ServerBuilder}
import io.grpc.protobuf.services.ProtoReflectionService
import scala.concurrent.Future

class MathServiceImpl(volatilityService: VolatilityService[IO]) extends MathServiceGrpc.MathService with LazyLogging {
  override def calculateVolatility(request: VolatilityRequest): Future[VolatilityResponse] = {
    implicit val runtime: cats.effect.unsafe.IORuntime = cats.effect.unsafe.IORuntime.global
    val returns = request.returns.toList
    volatilityService.calculateVolatility(returns).map {
      case Right(volatility) =>
        logger.info(s"Successfully calculated volatility: $volatility")
        VolatilityResponse(volatility = volatility)
      case Left(error) =>
        logger.error(s"Error calculating volatility: $error")
        throw new io.grpc.StatusRuntimeException(io.grpc.Status.INVALID_ARGUMENT.withDescription(error.toString))
    }.unsafeToFuture()
  }
}

object GrpcServer extends IOApp with LazyLogging {
  private val port = 50055

  private def createServer: Resource[IO, Server] = {
    val volatilityService = new VolatilityService[IO]
    val serviceImpl = new MathServiceImpl(volatilityService)

    Resource.make(
      IO {
        val server = ServerBuilder
          .forPort(port)
          .addService(MathServiceGrpc.bindService(serviceImpl, scala.concurrent.ExecutionContext.global))
          .addService(ProtoReflectionService.newInstance())
          .build()
        server.start()
        logger.info(s"Starting gRPC server on port $port")
        server
      }
    )(server => IO(server.shutdown()))
  }

  override def run(args: List[String]): IO[ExitCode] = {
    createServer.use { server =>
      IO {
        server.awaitTermination()
        ExitCode.Success
      }
    }
  }
} 