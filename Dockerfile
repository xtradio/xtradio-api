FROM scratch

LABEL "maintainer"="XTRadio Ops <contact@xtradio.org"
LABEL "version"="0.1"
LABEL "description"="XTRadio API"

ADD ./bin/xtradio-api /xtradio-api

#ADD ./bin/xtradio-api /xtradio-api
EXPOSE 10000

CMD ["/xtradio-api"]