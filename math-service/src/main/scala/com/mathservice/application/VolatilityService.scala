package com.mathservice.application

import cats.effect.Sync
import com.mathservice.domain.VolatilityCalculator
import com.typesafe.scalalogging.LazyLogging

class VolatilityService[F[_]: Sync] extends LazyLogging {
  def calculateVolatility(returns: List[Float]): F[Either[String, Double]] = {
    Sync[F].delay {
      logger.info(s"Calculating volatility for ${returns.length} data points")
      VolatilityCalculator.calculateVolatility(returns)
    }
  }
} 