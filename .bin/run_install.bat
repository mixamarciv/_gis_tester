::получаем curpath:
@FOR /f %%i IN ("%0") DO SET curpath=%~dp0
::задаем основные переменные окружения
@CALL "%curpath%/set_path.bat"


@del app.exe
@CLS

@echo === install ===================================================================
go get -u "github.com/gorilla/mux"
go get -u "github.com/satori/go.uuid"
go get -u "github.com/parnurzeal/gorequest"
go get -u "github.com/palantir/stacktrace"
go get -u "github.com/gosuri/uilive"


::библиотека для работы с XMLками
go get -u "github.com/jteeuwen/go-pkg-xmlx"


go install

@echo ==== end ======================================================================
@PAUSE
