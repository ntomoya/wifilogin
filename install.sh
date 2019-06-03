#!/bin/bash

plutil -lint wifilogin.plist && launchctl load wifilogin.plist