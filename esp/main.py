import time
import network
import urequests

from machine import Pin

# Ensure led connected to GPIO_38
led = Pin(38, Pin.OUT)
# Test led, led should turn on and then off
led.on()
time.sleep(0.5)
led.off()

# Wifi connection settings
ssid = 'WIFI_NAME'
password = 'WIFI PASSWORD'

station = network.WLAN(network.STA_IF)
station.active(True)
station.connect(ssid, password)

while not station.isconnected():
	time.sleep(1)

# esp should now be connected to wifi
# turn on led for visual aid
led.on()

# sample fetch
url = 'https://jsonplaceholder.typicode.com/todos/1'
response = urequests.get(url)
print(response.text)
response.close()

