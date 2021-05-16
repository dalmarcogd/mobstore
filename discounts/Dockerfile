FROM python:3.9.5-slim-buster

WORKDIR /srv

ADD Pipfile* ./

RUN pip install --no-cache -U pip pipenv && pipenv install --system

ADD . .

EXPOSE 8082

ENTRYPOINT python -m src

