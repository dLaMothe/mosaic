package main

import (
    //"fmt"
    //"log"
    "html/template"
    "image"
    "image/jpeg"
    _ "image/png"
    _ "image/gif"
    "image/color"
    //"image/draw"
    "net/http"
    "os"
)

//GLOBALS
var templates = template.Must((template.ParseFiles("form.html","image.html")))
const imgwidth float32 = 150
const imgheight float32 = 150

//ERROR CHECKING
func check(w http.ResponseWriter, r *http.Request, e error) {
    if e != nil {
        http.Redirect(w,r, "/", http.StatusFound)
        return
    }
}


//FUNCTIONS

func tileImage(height float32, width float32, img image.Image) *image.RGBA {

    //Initialize Destination rectangle
    dst := image.NewRGBA(image.Rect(0,0,int(width*imgwidth),int(height*imgheight)))
    for i := 0; i < int(imgwidth); i++ {
    //Iterate over all tiles
        for j := 0; j < int(imgheight); j++ {
            //Get the source rectangle
            //sr := image.Rect(width*i,height*j,width,height)
            //Destination point
            //dp := image.Point{width*i,height*j}
            //Destination rectangle
            //rec := image.Rectangle{dp, dp.Add(sr.Size())}
            //Colour counter for tile
            var red, green, blue, alpha float32 = 0, 0, 0, 0
            //Iterate over individual tile
            for k := 0; k < int(width); k++ {
                for l := 0; l < int(height); l++ {
                    tmpred, tmpgreen, tmpblue, tmpalpha := img.At(int(width*(float32(k)+float32(i))),int(height*(float32(l)+float32(j)))).RGBA()
                    //Need to divide by 256 to convert 16 bit integer range to 8 bit integer range
                    red += (float32(tmpred) / 256)
                    green += (float32(tmpgreen) / 256)
                    blue += (float32(tmpblue) / 256)
                    alpha += (float32(tmpalpha) / 256)
                }
            }
            //Calculate average colour
            avgred := (red / float32(width*height))
            avggreen := (green / float32(width*height))
            avgblue := (blue / float32(width*height))
            avgalpha := (alpha / float32(width*height))

            //Iterate over tile again to refill
            for k := 0; k < int(width); k++ {
                for l := 0; l < int(height); l++ {
                    dst.Set(((int(width*float32(i)+float32(k)))),((int(height*float32(j)+float32(l)))), color.RGBA{uint8(avgred),uint8(avggreen),uint8(avgblue),uint8(avgalpha)})
                }
            }

            //draw.Draw(dst,rec,img,sr.Min,draw.Src)
        }
    }
    return dst
}

func loadImage(w http.ResponseWriter, r *http.Request) {
    fileimg, _, err := r.FormFile("file")
    fileimg2, _, err := r.FormFile("file")
    check(w,r,err)
    img, _, err := image.Decode(fileimg)
    defer fileimg.Close()
    check(w,r,err)
    imgconf, _, err := image.DecodeConfig(fileimg2)
    check(w,r,err)
    defer fileimg2.Close()
    height := float32(imgconf.Height)
    width := float32(imgconf.Width)
    //Get the tile sizes of the image
    tileh := height / imgheight
    tilew := width / imgwidth
    //Tile the image
    dst := tileImage(tileh, tilew, img)
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

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, "form")
}

//MAIN
func main() {
    //Function Handlers
    http.HandleFunc("/image",imageHandler)
    http.HandleFunc("/",uploadHandler)
    //Begin server listening on port 8080
    http.ListenAndServe(":8080", nil)
}
