#!/bin/bash
set -e
export PATH=/usr/local/go/bin:$PATH

APP=1
EXEC=1

echo "🔍 Checking..."

#command -v go >/dev/null 2>&1 || { echo "❌ Go not found"; exit 1; }

echo "ตรวจเช็คไฟล์"
sleep 1
[ -f "icon.png" ] || { echo "❌ icon.png missing"; exit 1; }
[ -f "main.go" ] || { echo "❌ main.go missing"; exit 1; }
[ -f "go.mod" ] || { echo "❌ main.go missing"; exit 1; }
[ -f "go.sum" ] || { echo "❌ main.go missing"; exit 1; }

echo "🔨 build..."
sleep 1
go mod tidy
go build -ldflags="-s -w" -o $EXEC

echo "📦 prepare...AppDir..."
sleep 1
rm -rf $APP.AppDir
mkdir -p $APP.AppDir
cp $EXEC $APP.AppDir/

cat > $APP.AppDir/AppRun << 'EOF'
#!/bin/bash
HERE="$(dirname "$(readlink -f "$0")")"
exec "$HERE/1"
EOF

chmod +x $APP.AppDir/AppRun

cat > $APP.AppDir/$APP.desktop << EOF
[Desktop Entry]
Name=1
Exec=1
Icon=1
Type=Application
Categories=Utility;Development;
Terminal=false
EOF

cp icon.png $APP.AppDir/$EXEC.png
cp icon.png $APP.AppDir/.DirIcon

echo "🚀 pack..."
./appimagetool-x86_64.AppImage $APP.AppDir
sleep 2
cp $APP-x86_64.AppImage $APP.AppDir/$APP-x86_64.AppImage 

echo "📦 tar..."
tar -czf $APP.tar.gz $APP.AppDir
sleep 2

echo "🧹 cleanup..."
rm -rf $APP.AppDir

echo "✅ DONE"
