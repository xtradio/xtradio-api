FROM scratch
LABEL "maintainer"="XTRadio Ops <contact@xtradio.org"
LABEL "version"="0.1"
LABEL "description"="XTRadio API"

COPY bin/xtradio-api /xtradio-api

ENTRYPOINT /xtradio-api

EXPOSE 10000