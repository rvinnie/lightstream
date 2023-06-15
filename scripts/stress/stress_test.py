import requests
import random
import time

def do_stress():
    while True:
        path = 'http://gateway:8080/images/%s' % random.randint(0, 2)
        try:
            requests.get(path)
        except:
            print("An exception occurred")
        time.sleep(0.5)

if __name__ == "__main__":
    do_stress()