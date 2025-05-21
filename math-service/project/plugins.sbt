resolvers += Resolver.sbtPluginRepo("releases")

addSbtPlugin("com.eed3si9n" % "sbt-assembly" % "2.1.1")
addSbtPlugin("com.thesamet" % "sbt-protoc" % "1.0.6")
addSbtPlugin("com.github.sbt" % "sbt-native-packager" % "1.9.16")

libraryDependencies += "com.thesamet.scalapb" %% "compilerplugin" % "0.11.13" 