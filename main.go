package main

import (
    //"fmt"
    "image"
    _ "image/jpeg"
    "image/draw"
    "os"
)

func check(e error) {
    if e != nil {
        log.Fatal(e)
    }
}


func main() {
  fileimg, err := os.Open("test.jpg")
  check(err)
  defer fileimg.Close()
  img, _, err := image.Decode(fileimg)
  check(err)

  m := image.NewRGBA(image.Rect(0,0,800,600))

  draw.Draw(m,m.Bounds(), img, image.Point{0,0},draw.Src)

}
