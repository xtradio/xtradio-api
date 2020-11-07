# xtradio-api

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/4f18592c47374d23982426853fb9a6ad)](https://app.codacy.com/app/puck/xtradio-api?utm_source=github.com&utm_medium=referral&utm_content=xtradio/xtradio-api&utm_campaign=badger)
[![Build Status](https://github.com/xtradio/xtradio-api/workflows/build/badge.svg?branch=master)](https://github.com/xtradio/xtradio-api/actions)

XTRadio JSON Api to retrieve currently playing song and additional metadata.

## Build instructions

Build using :
``` CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/xtradio-api . ```

## Golang tutorials

The following tutorials helped us build our code, we've started with no prior knowledge of golang, maybe they will be helpfull for others as well.

* [Creating restful API with golang](https://tutorialedge.net/post/golang/creating-restful-api-with-golang/)
