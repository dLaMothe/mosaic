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
const imgwidth int = 20
const imgheight int = 20

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

func getDim(file *os.File) (int,int) {
    img, _, err := image.DecodeConfig(file)
    check(err)
    return img.Height, img.Width
}

func loadImage(w http.ResponseWriter, r *http.Request) {
    //TOFIX: Figure out file
    fileimg, err := os.Open("test.jpg")
    check(err)
    defer fileimg.Close()
    img, _, err := image.Decode(fileimg)
    check(err)
    fileimgtwo, err := os.Open("test.jpg")
    height, width := getDim(fileimgtwo)
    defer fileimgtwo.Close()
    //Get the tile sizes of the image
    tileh := height / imgheight
    tilew := width / imgwidth
    //Get the source rectangle
    sr := image.Rect(0,0,tilew,tileh)
    //Initialize Destination rectangle
    dst := image.NewRGBA(image.Rect(0,0,width,height))
    for i := 0; i < imgwidth; i++ {
        for j := 0; j < imgheight; j++ {
            dp := image.Point{tilew*i,tileh*j}
            rec := image.Rectangle{dp, dp.Add(sr.Size())}
            draw.Draw(dst,rec,img,sr.Min,draw.Src)
        }
    }
    newimg, _ :=  os.Create("tmp.jpg")
    defer newimg.Close()

    jpeg.Encode(newimg, dst, &jpeg.Options{jpeg.DefaultQuality})

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
