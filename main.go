package main

import (
    //"fmt"
    "log"
    "html/template"
    "image"
    "image/jpeg"
    "image/draw"
    "net/http"
    "os"
)

//GLOBALS
var templates = template.Must((template.ParseFiles("image.html")))

/*STRUCTS
type Page struct {
    Title string
    Body []byte
}
*/
//ERROR CHECKING
func check(e error) {
    if e != nil {
        log.Fatal(e)
    }
}

//FUNCTIONS

func loadImage(w http.ResponseWriter, r *http.Request) {
  fileimg, err := os.Open("test.jpg")
  check(err)
  defer fileimg.Close()
  img, _, err := image.Decode(fileimg)
  check(err)

  m := image.NewRGBA(image.Rect(0,0,800,600))

  draw.Draw(m,m.Bounds(), img, image.Point{0,0},draw.Src)

  newimg, _ :=  os.Create("tmp.jpg")
  defer newimg.Close()

  jpeg.Encode(newimg, m, &jpeg.Options{jpeg.DefaultQuality})

}

func renderTemplate(w http.ResponseWriter, tmpl string) {
    err := templates.Execute(w, tmpl+".html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

//PAGE HANDLERS
func imageHandler(w http.ResponseWriter, r *http.Request) {
    //Load the image
    loadImage(w,r)
    //Serve file to load the image
    http.ServeFile(w, r, "tmp.jpg")
    //Render the page
    renderTemplate(w, "image")
}

//MAIN
func main() {
    //Function Handlers
    http.HandleFunc("/",imageHandler)
    //Begin server listening on port 8080
    http.ListenAndServe(":8080", nil)
}
