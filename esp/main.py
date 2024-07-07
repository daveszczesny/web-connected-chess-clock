import time
import network
import urequests
import json
import time

from machine import Pin


""" Time control constants """
class TimeConstants:
    ONE_MINUTE = 60_000
    TWO_MINUTE = ONE_MINUTE * 2
    FIVE_MINUTE = ONE_MINUTE * 5
    TEN_MINUTE = ONE_MINUTE * 10

""" Player class  """
class Player:
    def __init__(self, button_pin, led_pin):
        self._timer_ms = TimeConstants.FIVE_MINUTE  # default time control of 5 minutes
        self._is_timer_running = False
        self._player_name = ""

        # Assign button and led to player (NON MUTABLE)
        self._button = Pin(button_pin, Pin.IN, Pin.PULL_UP)
        self._led = Pin(led_pin, Pin.OUT)

    @property
    def player_name(self) -> str:
         return self._player_name
    
    @player_name.setter
    def player_name(self, value: str):
         self._player_name = value

    @property
    def timer_ms(self):
        return self._timer_ms
    
    @timer_ms.setter
    def timer_ms(self, value):
        self._timer_ms = value

    @property
    def is_timer_running(self):
        return self._is_timer_running
    
    @is_timer_running.setter
    def is_timer_running(self, value):
        self._is_timer_running = value

    @property
    def button(self):
        return self._button
    
    @property
    def led(self):
        return self._led
    

    def read_from_json(self, res):
         assert isinstance(res, dict)
         self.timer_ms = res["game_length_ms"]


    def decrement_timer(self, value: int) -> None:
        self._timer_ms -= value

    def is_timer_up(self) -> bool:
        return self._timer_ms <= 0


"""
	Configure Pins
"""

# Ensure led connected to GPIO_38
wifi_led = Pin(38, Pin.OUT)

"""
    Player pins
    Player 1 has pin GPIO_40 for button and 37 for led
    Player 2 has pin GPIO_41 for button and 38 for led
    Ensure correct connection on board
"""
player1 = Player(40, 37)
player2 = Player(41, 38)


""" Test Connections """
"""
    WiFi LED, and both player LEDs should turn on,
    wait half a second and then turn off
    ensure working
"""
wifi_led.on()
player1.led.on()
player2.led.on()
time.sleep(0.5)
wifi_led.off()
player1.led.off()
player2.led.off()


"""
	Connect to WIFI
"""
ssid = 'GnaC_WiFi'
password = 'gnacaccess'

station = network.WLAN(network.STA_IF)
station.active(True)
station.connect(ssid, password)

while not station.isconnected():
	time.sleep(1)

# esp should now be connected to wifi
# turn on led for visual aid
wifi_led.on()


"""
	Fetch config file
"""
url = 'https://httpwc-nrfmpbzftq-ew.a.run.app/esp'
response = urequests.get(url)

# Set player settings
if response.status_code == 200:
    try:
        res = json.loads(response.text)
        player1.read_from_json(res["timecontrol"]["player_one"])
        player2.read_from_json(res["timecontrol"]["player_two"])
    except Exception as ex:
        pass # no way to view error

response.close()


"""
	Instaniate Displays
"""
# Connect Displays 1 and 2


"""
	Game Logic
"""

# Check if a user's time  has run out

def user_timer_up() -> bool:
	return player1.is_timer_up() or player2.is_timer_up()

def timer_set(player: Player, start_time):
    current_time = time.time()
    elapsed_time = current_time - start_time
    player.decrement_timer(elapsed_time)
    return current_time


start_time = time.time()

while True:

    if player1.is_timer_running:
         start_time = timer_set(player1, start_time)
    elif player2.is_timer_running:
         start_time = timer_set(player2, start_time)

    # Wait until a button is pressed by a player to start the game
    if player1.button.value() == True:
        player2.is_timer_running = False
        player1.is_timer_running = True
    elif player2.button.value() == True:
         player1.is_timer_running = False
         player2.is_timer_running = True


    if player1.is_timer_up or player2.is_timer_up:
         break
