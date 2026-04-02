#!/bin/bash
set -e

APP=1
EXEC=1

echo "🔍 checking dependencies..."

# ✅ check go
command -v go >/dev/null 2>&1 || { echo "❌ Go not installed"; exit 1; }

# ✅ check fyne CLI
if ! command -v fyne >/dev/null 2>&1; then
  echo "⚠️ fyne CLI not found → installing..."
  go install fyne.io/fyne/v2/cmd/fyne@latest
fi

# ✅ check appimagetool
if [ ! -f "./appimagetool-x86_64.AppImage" ]; then
  echo "⚠️ downloading appimagetool..."
  wget -q https://github.com/AppImage/AppImageKit/releases/latest/download/appimagetool-x86_64.AppImage
  chmod +x appimagetool-x86_64.AppImage
fi

#####################
[ -f "icon.png" ] || { echo "❌ icon.png missing"; exit 1; }
[ -f "main.go" ] || { echo "❌ main.go missing"; exit 1; }
#####################

echo "📦 preparing go modules..."
go mod tidy

#####################
echo "🎨 bundle icon..."
rm -f bundled.go
fyne bundle icon.png > bundled.go

echo "🔨 build..."
go build -o $EXEC

echo "📦 prepare..."
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
Categories=Utility;
Terminal=false
EOF

# icon
cp icon.png $APP.AppDir/$EXEC.png
cp icon.png $APP.AppDir/.DirIcon

mkdir -p $APP.AppDir/usr/share/icons/hicolor/256x256/apps
cp icon.png $APP.AppDir/usr/share/icons/hicolor/256x256/apps/$EXEC.png

echo "🚀 pack..."
./appimagetool-x86_64.AppImage $APP.AppDir

echo "✅ DONE"
