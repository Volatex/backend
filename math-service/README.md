# Math Service

Микросервис для расчета исторической волатильности, на Scala с использованием gRPC.

```bash
grpcurl -plaintext -d '{"returns": [0.01, -0.02, 0.03, -0.01, 0.02]}' localhost:50055 mathservice.MathService/CalculateVolatility
```
