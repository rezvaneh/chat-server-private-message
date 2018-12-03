package config

import "os"

type Configs struct {
  Hostname string
  Port string
}

var config = Configs{}

func LoadConfig() Configs {
  if config.Port != "" {
    return config
  }

  var conf = Configs {
    Hostname: "localhost",
    Port: "8080",
  }
  config = conf
  return conf
}

func CheckForError(err error, message string) {
  if err != nil {
      println(message + ".", err.Error())
      os.Exit(1)
  }
}


