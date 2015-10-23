# standup with gunicorn!

FROM ubuntu:14.04
MAINTAINER Daniel Riti <dmriti@gmail.com>

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update
RUN apt-get install -y python python-pip python-virtualenv gunicorn

# Setup flask application
RUN mkdir -p /deploy
COPY ./ /deploy
RUN pip install -r /deploy/requirements.txt
WORKDIR /deploy

EXPOSE 5000

# Start gunicorn
CMD ["/usr/bin/gunicorn", "--preload", "--config", "/deploy/gunicorn_config.py", "standup:app"]
