# Rosen

A graphical user interface for [Magic Wormhole](https://github.com/warner/magic-wormhole). 

Project is currently not actively developed. See [Wormhole GUI](https://github.com/jacalz/wormhole-gui) as a alternative.

## Screenshots

![Rosenbridge Send](/docs/rosenbridge-send.png?raw=true "Rosenbridge Send")
![Rosenbridge Receive](/docs/rosenbridge-receive.png?raw=true "Rosenbridge Receive")


## Compile

Linux: 
```
$(go env GOPATH)/bin/qtdeploy build desktop && ./deploy/linux/rosen 
```

Windows:
```
$(go env GOPATH)/bin/qtdeploy -docker build windows
```

Reference: [QT Getting Started](https://github.com/therecipe/qt/wiki/Getting-Started)


## Technology

* Golang
* GUI: QT with therecipe/qt library
* Backend: Magic Wormhole with psanford/wormhole-william library

## Name Origin

https://en.wikipedia.org/wiki/Nathan_Rosen

```
Nathan Rosen (Hebrew: נתן רוזן; March 22, 1909 – December 18, 1995) was an Israeli 
physicist noted for his study on the structure of the hydrogen atom and his work with
Albert Einstein and Boris Podolsky on entangled wave functions and the EPR paradox. 
The Einstein–Rosen bridge, later named the wormhole, was a theory of Nathan Rosen.
```
