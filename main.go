package main

import (
  "fmt"
  "log"
  "net/http"
  "os"
  "os/signal"
  "syscall"
  "time"
  "./bme280"
)

func main() {
  // TODO: make configurable
  botHost := "0.0.0.0:8080"

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

    url := fmt.Sprintf( "%s%s%st=%f&h=%f&p=%f", "http://", botHost, "/bme280/?", temperature, humidity, pressure )
    resp, err := http.Get( url )
    if err != nil {
      log.Printf( "Error sending request: %v", err )
      // Write data for debug purposes
      log.Printf( "temperature: %.2f C, humidity: %.2f RH, pressure: %.2f mmHg\n", temperature, humidity, pressure )
    } else {
      log.Printf( "Response Status: %v", resp.Status )
    }

    // TODO: add local CSV log

    time.Sleep( logPeriod )
  }

  log.Println( "Closing" )
}
