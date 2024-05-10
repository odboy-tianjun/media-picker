rm -f OdMediaPicker.exe
rm -f OdMediaPickerRelease.exe
# build
go build -o OdMediaPicker.exe main.go
# upx compress
./upx -o OdMediaPickerRelease.exe OdMediaPicker.exe
rm -f OdMediaPicker.exe
