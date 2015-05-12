package main

import (
    "fmt"
    "math"
    "html/template"
    "image"
    "image/jpeg"
    _ "image/png"
    _ "image/gif"
    //"image/color"
    "image/draw"
    "net/http"
    "os"
)

//GLOBALS
var templates = template.Must((template.ParseFiles("image.html")))
const imgwidth float32 = 40
const imgheight float32 = 40
const tilenum int = 250

//STRUCTS
type Color struct {
    r,g,b,a float32
}

//ERROR CHECKING
func check(w http.ResponseWriter, r *http.Request, e error) {
    if e != nil {
        http.Redirect(w,r, "/", http.StatusFound)
        return
    }
}


//FUNCTIONS

func getTiles(tiles *[tilenum]image.Image, tileh float32, tilew float32, w http.ResponseWriter, r *http.Request) {


    for i := 0; i < tilenum; i++ {
        //Initialize tile rectangle
        dst := image.NewRGBA(image.Rect(0,0,int(tilew),int(tileh)))
        val := r.FormValue(fmt.Sprint("photo",i))
        resp, err := http.Get(val)
        check(w,r,err)
        defer resp.Body.Close()

        m, _, err := image.Decode(resp.Body)
        check(w,r,err)
        rec := m.Bounds()

        imgheight := rec.Dy()
        imgwidth := rec.Dx()
        //Begin downsizing process
        var xratio float32 = float32(imgwidth) / tilew
        var yratio float32 = float32(imgheight) / tileh

        for j := 0; j < int(tileh); j++ {
            for k := 0; k < int(tilew); k++ {
                px := math.Floor(float64(k)*float64(xratio))
                py := math.Floor(float64(j)*float64(yratio))
                dst.Set(k,j,m.At(int(px),int(py)))
            }
        }
        tiles[i] = dst
    }
}

func getColors(tiles *[tilenum] image.Image, tilecolors *[tilenum] Color, width, height float32) {
    //Colour counter for tile
    for i := 0; i < tilenum; i++ {
        var red, green, blue, alpha float32 = 0, 0, 0, 0
        //Iterate over individual tile
        for j := 0; j < int(width); j++ {
            for k := 0; k < int(height); k++ {
                tmpred, tmpgreen, tmpblue, tmpalpha := tiles[i].At(j,k).RGBA()
                //Need to divide by 256 to convert 16 bit integer range to 8 bit integer range
                red += (float32(tmpred) / 256)
                green += (float32(tmpgreen) / 256)
                blue += (float32(tmpblue) / 256)
                alpha += (float32(tmpalpha) / 256)
            }
        }
        //Calculate average colour
        tilecolors[i].r = (red / float32(width*height))
        tilecolors[i].g = (green / float32(width*height))
        tilecolors[i].b = (blue / float32(width*height))
        tilecolors[i].a = (alpha / float32(width*height))
    }
}

func compareTiles(height,width,r,g,b,a float32,tilecolors [tilenum]Color) int {
    var best float64 = 10000
    bestindex := 0
    for i := 0; i < tilenum; i++ {
        //Compare the two averages
        difference := (math.Abs(float64(r-tilecolors[i].r))) + (math.Abs(float64(g-tilecolors[i].g)) + (math.Abs(float64(b-tilecolors[i].b)) + (math.Abs(float64(a - tilecolors[i].a)))))
        if difference < best {
            best = difference
            bestindex = i
        }
    }
    return bestindex
}

func tileImage(height float32, width float32, img image.Image, rec image.Rectangle, w http.ResponseWriter, r *http.Request) *image.RGBA {

    //Initialize Destination rectangle
    dst := image.NewRGBA(rec)
    //Initialize Tile array
    //Declare array of images to hold tiles
    var tileArr [tilenum]image.Image
    getTiles(&tileArr, height, width, w, r)
    //Declare array of structs to hold image color averages
    var tileColor [tilenum] Color
    getColors(&tileArr, &tileColor, height, width)
    for i := 0; i < int(imgwidth); i++ {
    //Iterate over all tiles
        for j := 0; j < int(imgheight); j++ {
            //Colour counter for tile
            var red, green, blue, alpha float32 = 0, 0, 0, 0
            var xoffset float32 = width*float32(i)
            var yoffset float32 = height*float32(j)
            //Iterate over individual tile
            for k := 0; k < int(width); k++ {
                for l := 0; l < int(height); l++ {
                    tmpred, tmpgreen, tmpblue, tmpalpha := img.At(int(xoffset+float32(i)),int(yoffset+float32(j))).RGBA()
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

            tileindex := compareTiles(height,width,avgred,avggreen,avgblue,avgalpha,tileColor)

            //Get the source rectangle
            sr := tileArr[tileindex].Bounds()
            //Destination point
            dp := image.Point{int(width*float32(i)),int(height*float32(j))}
            //Destination rectangle
            rec := image.Rectangle{dp, dp.Add(sr.Size())}

            draw.Draw(dst,rec,tileArr[tileindex],sr.Min,draw.Src)
        }
    }
    return dst
}

func loadImage(w http.ResponseWriter, r *http.Request) {
    fileimg, _, err := r.FormFile("imgfile")
    check(w,r,err)
    img, _, err := image.Decode(fileimg)
    defer fileimg.Close()
    check(w,r,err)
    rec := img.Bounds()
    height := rec.Dy()
    width := rec.Dx()
    //Get the tile sizes of the image
    tileh := float32(height) / imgheight
    tilew := float32(width) / imgwidth
    //Tile the image
    dst := tileImage(tileh, tilew, img, rec, w, r)
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
    err := os.Remove("tmp.jpg")
    check(w,r,err)
    //Render the page
    renderTemplate(w, "image")
}

//MAIN
func main() {
    fs := http.FileServer(http.Dir("static") )
    //Function Handlers
    http.Handle("/",fs)
    http.HandleFunc("/image",imageHandler)
    //Begin server listening on port 8080
    http.ListenAndServe(":8080", nil)
}
