ThisBuild / version := "0.1.0-SNAPSHOT"
ThisBuild / scalaVersion := "2.13.11"

lazy val root = (project in file("."))
  .settings(
    name := "math-service",
    libraryDependencies ++= Seq(
      "io.grpc" % "grpc-netty-shaded" % "1.62.2",
      "io.grpc" % "grpc-protobuf" % "1.62.2",
      "io.grpc" % "grpc-stub" % "1.62.2",
      "io.grpc" % "grpc-services" % "1.62.2",
      "com.thesamet.scalapb" %% "scalapb-runtime-grpc" % scalapb.compiler.Version.scalapbVersion,
      "org.typelevel" %% "cats-effect" % "3.5.4",
      "com.typesafe.scala-logging" %% "scala-logging" % "3.9.5",
      "ch.qos.logback" % "logback-classic" % "1.5.3",
      "org.scalatest" %% "scalatest" % "3.2.18" % Test
    ),
    assembly / assemblyMergeStrategy := {
      case PathList("META-INF", "io.netty.versions.properties") => MergeStrategy.first
      case PathList("META-INF", "versions", "9", "module-info.class") => MergeStrategy.first
      case PathList("module-info.class") => MergeStrategy.discard
      case x =>
        val oldStrategy = (assembly / assemblyMergeStrategy).value
        oldStrategy(x)
    },
    assembly / assemblyJarName := "math-service.jar"
  )

Compile / PB.targets := Seq(
  scalapb.gen(grpc = true) -> (Compile / sourceManaged).value / "scalapb"
)

Compile / PB.protoSources := Seq(
  (Compile / sourceDirectory).value / "protobuf"
)

scalacOptions ++= Seq(
  "-Ymacro-annotations",
  "-deprecation",
  "-feature",
  "-unchecked",
  "-language:implicitConversions",
  "-language:higherKinds"
) 