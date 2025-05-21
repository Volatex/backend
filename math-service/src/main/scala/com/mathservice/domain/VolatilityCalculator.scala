package com.mathservice.domain

import cats.implicits._

object VolatilityCalculator {
  private val AnnualizationFactor = math.sqrt(12)

  def calculateVolatility(returns: List[Float]): Either[String, Double] = {
    if (returns.isEmpty) {
      Left("Input sequence cannot be empty")
    } else if (returns.length < 2) {
      Left("At least 2 data points are required for volatility calculation")
    } else {
      val mean = returns.sum / returns.length
      val variance = returns.map(r => math.pow(r - mean, 2)).sum / (returns.length - 1)
      val stdDev = math.sqrt(variance)
      Right(stdDev * AnnualizationFactor)
    }
  }
} 