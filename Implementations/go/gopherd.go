package main

import (
  "bytes"
  "fmt"
  "log"
  "net"
  "os"
  "path/filepath"
  "strings"
)

type errorString struct {
  s string
}

func (e *errorString) Error() string {
  return e.s
}

func main() {
  wd, _ := os.Getwd()
  line, err := net.Listen("tcp", ":70")
  if err != nil {
    log.Fatal(err)
  } else {
    for {
      conn, err := line.Accept()
      if err != nil {
        log.Fatal(err)
      }

      go handleConnection(wd, conn)
    }
  }
}

func get_real_path(base string, path string) (string, error) {
  cleaned := filepath.Clean(path)
  joined := filepath.Join(base, cleaned)
  abspath, _ := filepath.Abs(joined)
  if ! strings.HasPrefix(abspath, base) {
    return abspath, &errorString{ "path not contained within base" }
  } else {
    return abspath, nil
  }
}

func handleConnection(server_dir string, conn net.Conn) {
  log.Printf("got connection from %v\n", conn.RemoteAddr())
  buffer := make([]byte, 4096)
  n, err := conn.Read(buffer)
  if err != nil {
    log.Fatal(err)
  }

  log.Printf("read %d bytes from %v\n", n, conn.RemoteAddr())
  str := strings.Split(string(buffer[:4096]), "\r\n")[0]

  path, err := get_real_path(server_dir, str)
  if err != nil {
    log.Printf("cannot resolve path: %s", err)
    conn.Close()
    return
  }

  output, err := getPath(path)
  if err != nil {
    log.Printf("cannot retrieve path: %s", err)
    conn.Close()
    return
  }

  conn.Write([]byte(output.String()))
  conn.Close()
}

func getPath(path string) (bytes.Buffer, error) {
  buffer := bytes.Buffer{}

  err := filepath.Walk(path, func(p string, info os.FileInfo, e error) error {
    buffer.WriteString(fmt.Sprintf("i%s\n", p))
    return nil
  })

  return buffer, err
}

