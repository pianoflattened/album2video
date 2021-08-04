package main

type FormData struct {
    albumDirectory string `json:"albumDirectory"`
    coverPath      string `json:"coverPath"`
    detectCover    bool   `json:"detectCover"`
    separateFiles  bool   `json:"separateFiles"`
    outputPath     string `json:"outputPath"`
}

type VideoData struct {}
