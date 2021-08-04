package main

type FormData struct {
    albumDirectory string `json:"albumDirectory"`
    coverPath      string `json:"coverPath"`
    detectCover    bool   `json:"detectCover"`
    separateFiles  bool   `json:"separateFiles"`
    outputPath     string `json:"outputPath"`
}

// make sure u look up ffprobe test / example data so u can make the struct
// unless u just wanna use gjson which is fine too considering how much effort that'll be

type VideoData struct {}
