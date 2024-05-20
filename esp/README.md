# ESP code base environment

## Setup
1) Download firmware from micropython
<br> - https://micropython.org/download/ESP32_GENERIC_S3/
2) Download esptool
`pip install esptool`
3) Find port of esp
`ls /dev/tty.*`
4) Erase all from esp using esptool
`esptool.py --chip esp32s3 --port /dev/ttyACM0 erase_flash`
5) flash download .bin file onto esp
`esptool.py --chip esp32s3 --port /dev/ttyACM0 write_flash -z 0 board-20210902-v1.17.bin`
6) Download ampy
`pip install adafruit-ampy`
7) Load code onto esp
`ampy -pport /dev/ttyUSB0 put main.py
