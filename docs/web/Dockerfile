FROM node:20.12.2-alpine3.19

#COPY ./docs/web/package.json /app/package.json

WORKDIR /app

RUN apk add --no-cache curl git openssh supervisor bash && \
    mkdir -p /var/log/supervisor /etc/supervisor/conf.d && \
    chown -R node:node /app

ARG APP_ENV=dev
COPY ./docs/web/docker/supervisor/supervisord.conf /etc/supervisord.conf
COPY ./docs/web/docker/supervisor/${APP_ENV} /etc/supervisor/conf.d/

ARG UID=1000
ARG GID=1000

RUN if [ "$GID" != "1000" ]; then apk add shadow && groupmod -g $GID node && usermod -u $UID -g $GID node; fi

USER node

COPY --chown=node:node ./docs/web /app

RUN if [ "$APP_ENV" != "dev" ]; then \
      yarn --immutable && cp .env.${APP_ENV} .env && yarn run build; \
    fi;

EXPOSE 8080

CMD ["/usr/bin/supervisord"]
