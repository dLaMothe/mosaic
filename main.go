package main

import (
    //"fmt"
    "log"
    "html/template"
    "image"
    "image/jpeg"
    "image/color"
    //"image/draw"
    "net/http"
    "os"
)

//GLOBALS
var templates = template.Must((template.ParseFiles("image.html")))
const imgwidth int = 5
const imgheight int = 5

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

func tileImage(height int, width int, img image.Image) *image.RGBA {

    //Initialize Destination rectangle
    dst := image.NewRGBA(image.Rect(0,0,width*imgwidth,height*imgheight))
    //Iterate over all tiles
    for i := 0; i < imgwidth; i++ {
        for j := 0; j < imgheight; j++ {
            //Get the source rectangle
            //sr := image.Rect(width*i,height*j,width,height)
            //Destination point
            //dp := image.Point{width*i,height*j}
            //Destination rectangle
            //rec := image.Rectangle{dp, dp.Add(sr.Size())}
            //Colour counter for tile
            var red, green, blue, alpha uint32 = 0, 0, 0, 0
            //Iterate over individual tile
            for k := 0; k < width; k++ {
                for l := 0; l < height; l++ {
                    tmpred, tmpgreen, tmpblue, tmpalpha := img.At(width*(k+i),height*(l+j)).RGBA()
                    red += tmpred
                    green += tmpgreen
                    blue += tmpblue
                    alpha += tmpalpha
                }
            }
            //Calculate average colour
            avgred := red / uint32((width*height))
            avggreen := green / uint32((width*height))
            avgblue := blue / uint32((width*height))
            avgalpha := alpha / uint32((width*height))

            //Iterate over tile again to refill
            for k := 0; k < width; k++ {
                for l := 0; l < height; l++ {
                    dst.Set(width*(k+i),height*(l+j), color.RGBA{uint8(avgred),uint8(avggreen),uint8(avgblue),uint8(avgalpha)})
                }
            }

            //draw.Draw(dst,rec,img,sr.Min,draw.Src)
        }
    }
    return dst
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

//MAIN
func main() {
    //Function Handlers
    http.HandleFunc("/",imageHandler)
    //Begin server listening on port 8080
    http.ListenAndServe(":8080", nil)
}
