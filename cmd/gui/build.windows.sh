#!/bin/bash

fyne package --os windows --tags ui --icon ../../assets/logo.png
mv gui.exe ../../builds/agos.exe
cp -R ../../bin ../../builds
