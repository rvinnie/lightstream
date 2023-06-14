FROM python:3.9-alpine

WORKDIR /usr/src/app

COPY ./ ./

RUN pip install -r requirements.txt

ENTRYPOINT python ./stress_test.py