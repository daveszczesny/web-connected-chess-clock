import time
import network
import urequests

from machine import Pin

"""
	Configure Pins
"""

# Ensure led connected to GPIO_38
led = Pin(38, Pin.OUT)

# User buttons to start game and stop ticker timer
user_button_1 = Pin(39, Pin.IN)
user_button_2 = Pin(40, Pin.IN)

# Test led, led should turn on and then off
led.on()
time.sleep(0.5)
led.off()


"""
	Connect to WIFI
"""

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


"""
	Fetch config file from Cloud Run
"""

# sample fetch
url = 'https://jsonplaceholder.typicode.com/todos/1'
response = urequests.get(url)
print(response.text)
response.close()


"""
	Instantiate Timer variables
"""

user_1_timer_seconds: int = 5000 # default value
user_2_timer_seconds: int = 5000

user_1_timer_running: bool = False
user_2_timer_running: bool = False

game_time_increment: int = 100


"""
	Instaniate Displays
"""
# Connect Displays 1 and 2


"""
	Game Logic
"""

# Check if a user's time  has run out
def user_timer_up() -> bool:
	if (user_1_timer_running and user_1_timer_seconds <= 0) or (user_2_timer_running and user_2_timer_seconds <= 0):
		return True
	

# Decrement user timer by time
def decrement_timer(user_timer: int) -> None:
	# Decrement time in milliseconds
	pass


# Add game increment time to user timer
def add_increment_to_timer(user_timer: int, time_increment: int) -> int:
	return user_timer - time_increment


while not user_timer_up():
	if user_1_timer_running:
		user_1_timer_seconds = decrement_timer(game_time_increment)
	elif user_2_timer_running:
		user_2_timer_seconds = decrement_timer(game_time_increment)

