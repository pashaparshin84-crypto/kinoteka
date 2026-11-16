@echo off
setlocal EnableExtensions EnableDelayedExpansion
chcp 65001 >nul
title Настройка Docker (WSL2) - Кинотека

:: ============================================================
::  setup-docker.bat
::  Готовит компьютер к запуску Docker Desktop:
::   1. Запрашивает права администратора (запуск от имени админа)
::   2. Включает компоненты Windows (WSL и Виртуальная машина)
::   3. Устанавливает/обновляет WSL2
::   4. Проверяет результат
::  После этого нужно ПЕРЕЗАГРУЗИТЬ компьютер,
::  затем запустить Docker Desktop и проект (см. ниже).
:: ============================================================

:: ---------- 1. Получение прав администратора ----------
net session >nul 2>&1
if %errorlevel% neq 0 (
    echo   Запрашиваю права администратора...
    powershell -Command "Start-Process -FilePath '%~f0' -Verb RunAs"
    exit /b
)

echo.
echo  ==============================================
echo   Настройка Docker (WSL2) для Windows
echo  ==============================================
echo.

:: ---------- 2. Проверка разрядности (нужен 64-bit) ----------
if "%PROCESSOR_ARCHITECTURE%"=="AMD64" (
    echo   [OK] 64-разрядная система
) else (
    echo   [ВНИМАНИЕ] WSL2 требует 64-битную Windows.
)

:: ---------- 3. Включение компонентов Windows ----------
echo.
echo   Включаю компоненты Windows (WSL и VirtualMachinePlatform)...
echo   Если они не включены - команды ниже ничего не сделают, это нормально.
dism.exe /online /enable-feature /featurename:Microsoft-Windows-Subsystem-Linux /all /norestart
dism.exe /online /enable-feature /featurename:VirtualMachinePlatform /all /norestart

:: ---------- 4. WSL2 ----------
echo.
echo   Устанавливаю/обновляю WSL2...
wsl --install --no-distribution
wsl --set-default-version 2

:: ---------- 5. Проверка ----------
echo.
echo   Проверяю состояние WSL:
wsl --status
echo.
echo  ==============================================
echo   ГОТОВО. Осталось:
echo    1. ПЕРЕЗАГРУЗИТЬ компьютер
echo    2. Запустить Docker Desktop и дождаться "Engine running"
echo    3. В папке проекта выполнить:
echo       docker compose up --build
echo    4. Открыть http://localhost:8080/
echo  ==============================================
echo.
pause
endlocal