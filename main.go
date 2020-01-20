package main

import (
  "fmt"
  "log"
  "net/http"
  "os"
  "os/signal"
  "syscall"
  "time"
  "encoding/json"
  "io/ioutil"
  "errors"
  "./bme280"
)

type Config struct {
  BotAddress string
}

func ReadConfig( cfgPath string ) ( *Config, error ) {
  file, err := ioutil.ReadFile( cfgPath )
  if err != nil {
    return nil, err
  }

  var cfg Config
  err = json.Unmarshal( file, &cfg )
  if err != nil {
    return nil, err
  }

  return &cfg, nil
}

func main() {
  cfg, err := ReadConfig( "config.json" )
  if err != nil {
    log.Fatal( err )
    return
  }

  isWorking := true
  sigs := make( chan os.Signal, 1 )
  go func() {
    sig := <-sigs
    log.Printf( "%v\n", sig )
    isWorking = false;
  } ()

  signal.Notify( sigs, syscall.SIGINT, syscall.SIGTERM )

  log.Println( "Starting..." )

  bme280conn, err := bme280.Connect( 118, 1 )
  if err != nil {
    log.Fatal( err )
  }

  defer bme280conn.Disconnect()

  log.Println( "BME280 init OK..." )

  logPeriod := time.Duration( 1 ) * time.Minute
  for isWorking {
    temperature, humidity, pressure, err := bme280conn.ReadData()
    if err != nil {
      log.Fatal( err )
    }

    timestamp := time.Now().Local().Format( "02012006150405MST" )
    url := fmt.Sprintf( "%s%s%sts=%s&t=%f&h=%f&p=%f", "http://", cfg.BotAddress, "/bme280/?", string( timestamp ), temperature, humidity, pressure )
    _, err = http.Get( url )
    if err != nil {
      log.Printf( "Error sending request: %v", err )
      // Write data for debug purposes
      log.Printf( "temperature: %.2f C, humidity: %.2f RH, pressure: %.2f mmHg\n", temperature, humidity, pressure )
    }

    // TODO: add local CSV log

    time.Sleep( logPeriod )
  }

  log.Println( "Closing" )
}
