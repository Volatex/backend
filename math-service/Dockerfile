# Этап 1: Кэширование зависимостей sbt
FROM eclipse-temurin:11-jdk as deps

# Установка sbt и создание пользователя
RUN apt-get update && \
    apt-get install -y curl unzip && \
    curl -L -o sbt.zip https://github.com/sbt/sbt/releases/download/v1.10.0/sbt-1.10.0.zip && \
    unzip sbt.zip -d /usr/local && \
    rm sbt.zip && \
    ln -s /usr/local/sbt/bin/sbt /usr/local/bin/sbt && \
    useradd -m -s /bin/bash sbtuser && \
    mkdir -p /home/sbtuser/.sbt /home/sbtuser/.ivy2 && \
    chown -R sbtuser:sbtuser /home/sbtuser

USER sbtuser
WORKDIR /app
COPY --chown=sbtuser:sbtuser project project
COPY --chown=sbtuser:sbtuser build.sbt .
RUN sbt update

# Этап 2: Сборка приложения
FROM eclipse-temurin:11-jdk as builder

# Установка sbt, protoc и создание пользователя
RUN apt-get update && \
    apt-get install -y curl unzip protobuf-compiler && \
    curl -L -o sbt.zip https://github.com/sbt/sbt/releases/download/v1.10.0/sbt-1.10.0.zip && \
    unzip sbt.zip -d /usr/local && \
    rm sbt.zip && \
    ln -s /usr/local/sbt/bin/sbt /usr/local/bin/sbt && \
    useradd -m -s /bin/bash sbtuser && \
    mkdir -p /home/sbtuser/.sbt /home/sbtuser/.ivy2 && \
    chown -R sbtuser:sbtuser /home/sbtuser

USER sbtuser
WORKDIR /app
COPY --from=deps --chown=sbtuser:sbtuser /home/sbtuser/.sbt /home/sbtuser/.sbt
COPY --from=deps --chown=sbtuser:sbtuser /home/sbtuser/.ivy2 /home/sbtuser/.ivy2
COPY --chown=sbtuser:sbtuser . .

RUN sbt clean assembly

# Этап 3: Создание итогового образа
FROM eclipse-temurin:11-jre

WORKDIR /app

# Копируем собранный JAR и конфигурацию
COPY --from=builder /app/target/scala-2.13/math-service.jar /app/math-service.jar
COPY --from=builder /app/src/main/resources /app/config

# Устанавливаем переменные окружения
ENV JAVA_OPTS="-Xmx512m -Xms256m"

# Открываем порт для gRPC
EXPOSE 50055

# Запускаем приложение
ENTRYPOINT ["java", "-jar", "/app/math-service.jar"] 