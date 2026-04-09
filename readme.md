oneimage
-ใช้สำหรับทำ .image (ไฟล์เดียว) 
-oneimage จำเป็นต้องมี ไฟล์ appimagetool-x86_64.AppImage อยู่ข้างๆเสมอ
-สามารถเอาตัว appimage version ที่ใหม่กว่ามาแทนได้เลย *แต่ชื่อต้องอ้างอิงด้านบน

**Go
**ใช้ได้กับภาษา Go
**Golang 
**fyne (gui)
**

เครื่องที่ใช้จะต้องมี 
-ภาษา go (golang) (ในเครื่อง)
-ImageMagick
-และอื่นๆที่ไฟล์ go.mod ต้องใช้งาน

แฟ้ม oneimage/
  ├── appimagetool-x86_64.AppImage (*ใช้รุ่นใหม่กว่าได้)
  └── oneimage_v1_0_0_0-x86_64.AppImage (ใช้งาน gui)

แฟ้ม**โปรเจคเป้าหมาย/
  ├── icon.png (ตั้งชื่อว่า icon.png) (Master)
  ├── main.go
  ├── go.mod
  └── go.sum

การใช้งาน
1. ใ้ส่ชื่อ app ,exec ,Display (โดยส่วนมาก ใช้ชื่อเดีนวกันหมด)
2. ช่อง categories ก็มีให้เลือกตามข้อความ ด้านล่างช่องกรอก ปิดท้ายด้วย ; เสมอ เช่น x; , x;x; , x;x;x;
3. เลือกแฟ้มโปรเจค (ระบบจะก็อปไฟล์ appimagetool-x86_64.AppImage ไปวางในแฟ้มโปรเจค)
4. Generate script (จะมีไฟล์ Build.sh ขึ้นที่แฟ้มโปรเจค)
5. Run build
