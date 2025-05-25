package com.mathservice.domain

import cats.implicits._

object VolatilityCalculator {
  // Using 247 as we filter out weekends
  private val AnnualizationFactor = math.sqrt(247)

  def calculateVolatility(returns: List[Double]): Either[String, Double] = {
    if (returns.isEmpty) {
      Left("Input sequence cannot be empty")
    } else if (returns.length < 2) {
      Left("At least 2 data points are required for volatility calculation")
    } else {
      // Calculate mean with better precision
      val mean = returns.sum / returns.length
      
      // Calculate variance with better precision
      val squaredDiffs = returns.map(r => {
        val diff = r - mean
        diff * diff  // More precise than math.pow
      })
      val variance = squaredDiffs.sum / (returns.length - 1)
      
      // Calculate standard deviation
      val stdDev = math.sqrt(variance)
      
      // Annualize volatility
      val volatility = stdDev * AnnualizationFactor
      
      if (volatility.isNaN || volatility.isInfinite) {
        Left("Cannot calculate volatility: result is not a valid number")
      } else {
        Right(volatility)
      }
    }
  }
} 