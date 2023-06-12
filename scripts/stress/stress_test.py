import requests
import random
import time

def do_stress():
    while True:
        path = 'http://gateway:8080/image/%s' % random.randint(0, 2)
        try:
            requests.get(path)
        except:
            print("An exception occurred")
        time.sleep(1)

if __name__ == "__main__":
    do_stress()