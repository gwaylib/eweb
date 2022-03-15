
Fork and modify from:
```
https://github.com/ot24net/eweb
https://github.com/labstack/echo
```

Changes:
```
* Make up global route of echo.Echo 
* Make up log output and add H reference gin, it is really good.
* Be careful! We use text/template instead the html/template
* Some new features, such as quic
```

Why base on echo?
```
* I like the return design in routing. It can make the compilation more reassuring. When return is missing, it will prompt compilation errors.
* I need to replace some features, it make me easy to changed it.
* These are the easiest to substitute when I tried to use them many years ago.
```

