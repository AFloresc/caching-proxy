@echo off
setlocal

:: Configuración
set PORT=3000
set ORIGIN=http://dummyjson.com
set PROXY=http://localhost:%PORT%

echo 🚀 Iniciando Catching Proxy en el puerto %PORT% apuntando a %ORIGIN%...
start "" /B cmd /C "go run cmd --port %PORT% --origin %ORIGIN%"
timeout /t 2 >nul

echo.
echo 📦 Primera petición (esperamos MISS)...
curl -i %PROXY%/products

echo.
echo 📦 Segunda petición (esperamos HIT)...
curl -i %PROXY%/products

echo.
echo 📊 Estadísticas de caché...
curl -s %PROXY%/cache/stats

echo.
echo 🧹 Limpiando caché vía HTTP...
curl -s -X POST %PROXY%/cache/clear

echo.
echo 📊 Estadísticas después de limpiar...
curl -s %PROXY%/cache/stats

