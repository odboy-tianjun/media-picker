rm -f OdImageWidthPicker.exe
rm -f OdImageWidthPickerRelease.exe
# build
go build -o OdImageWidthPicker.exe main.go
# upx compress
./upx -o OdImageWidthPickerRelease.exe OdImageWidthPicker.exe
rm -f OdImageWidthPicker.exe
