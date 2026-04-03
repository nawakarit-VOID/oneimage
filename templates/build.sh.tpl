#!/bin/bash
set -e

export PATH=/usr/local/go/bin:$PATH

APP={{.AppName}}
EXEC={{.ExecName}}
OFFLINE=${OFFLINE:-0}

APPIMAGE_TOOL=appimagetool-x86_64.AppImage
FYNE_BIN="$(go env GOPATH)/bin/fyne"
FYNE="fyne"

# ------------------------
log() { echo -e "\n🔹 $1"; }
fail() { echo "❌ $1"; exit 1; }

# ------------------------
check_go() {
  command -v go >/dev/null 2>&1 || fail "Go not installed"
  export PATH=$PATH:$(go env GOPATH)/bin
}

# ------------------------
setup_fyne() {
  if command -v fyne >/dev/null 2>&1; then
    return
  fi

  if [ -x "$FYNE_BIN" ]; then
    log "using fyne from GOPATH"
    FYNE="$FYNE_BIN"
    return
  fi

  [ "$OFFLINE" = "1" ] && fail "fyne not found (offline)"

  log "installing fyne..."
  go install fyne.io/fyne/v2/cmd/fyne@latest
}

# ------------------------
setup_appimagetool() {
  if [ -f "$APPIMAGE_TOOL" ]; then
    return
  fi

  [ "$OFFLINE" = "1" ] && fail "appimagetool missing (offline)"

  log "downloading appimagetool..."

  if command -v wget >/dev/null 2>&1; then
    wget -q https://github.com/AppImage/AppImageKit/releases/latest/download/$APPIMAGE_TOOL
  else
    curl -L -o $APPIMAGE_TOOL https://github.com/AppImage/AppImageKit/releases/latest/download/$APPIMAGE_TOOL
  fi

  chmod +x $APPIMAGE_TOOL
}

# ------------------------
check_files() {
  [ -f "icon.png" ] || fail "icon.png missing"
  [ -f "main.go" ] || fail "main.go missing"
}

# ------------------------
prepare_modules() {
  log "preparing modules..."
  if [ "$OFFLINE" = "1" ]; then
    go mod tidy -e
  else
    go mod tidy
  fi
}

# ------------------------
bundle_icon() {
  log "bundling icon..."
  rm -f bundled.go
  $FYNE bundle icon.png > bundled.go
}

# ------------------------
build_binary() {
  log "building..."
  go build -ldflags="-s -w" -o $EXEC || fail "build failed"
}

# ------------------------
prepare_appdir() {
  log "preparing AppDir..."

  rm -rf $APP.AppDir
  mkdir -p $APP.AppDir

  cp $EXEC $APP.AppDir/

  cat > $APP.AppDir/AppRun << 'EOF'
#!/bin/bash
HERE="$(dirname "$(readlink -f "$0")")"
exec "$HERE/{{.ExecName}}"
EOF

  chmod +x $APP.AppDir/AppRun

  cat > $APP.AppDir/$APP.desktop << EOF
[Desktop Entry]
Name={{.DisplayName}}
Exec={{.ExecName}}
Icon={{.ExecName}}
Type={{.Type}}
Categories={{.Categories}}
Terminal=false
EOF

  cp icon.png $APP.AppDir/$EXEC.png
  cp icon.png $APP.AppDir/.DirIcon

  mkdir -p $APP.AppDir/usr/share/icons/hicolor/256x256/apps
  cp icon.png $APP.AppDir/usr/share/icons/hicolor/256x256/apps/$EXEC.png
}

# ------------------------
pack_appimage() {
  log "packing AppImage..."
  "./$APPIMAGE_TOOL" $APP.AppDir
}

# ------------------------
# MAIN FLOW

log "checking dependencies..."
check_go
setup_fyne
setup_appimagetool

check_files
prepare_modules
bundle_icon
build_binary
prepare_appdir
pack_appimage

log "DONE ✅"