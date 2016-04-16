::получаем curpath:
@FOR /f %%i IN ("%0") DO SET curpath=%~dp0
::задаем основные переменные окружения
@CALL "%curpath%/set_path.bat"


@del app.exe
@CLS

@echo === install ===================================================================
go get github.com/gorilla/mux
go get github.com/satori/go.uuid
go get "github.com/parnurzeal/gorequest"
go get "github.com/palantir/stacktrace"

::go install

@echo ==== end ======================================================================
@PAUSE
