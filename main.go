package main

import (
  "log"
  "os"
  "os/signal"
  "syscall"
  "time"
  "./bme280"
)

func main() {
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

    // TODO: send to TG bot
    log.Printf( "temperature: %.2f C, humidity: %.2f RH, pressure: %.2f mmHg\n", temperature, humidity, pressure )

    time.Sleep( logPeriod )
  }

  log.Println( "Closing" )
}