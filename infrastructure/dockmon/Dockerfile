FROM czarsimon/dockmon:1.0
COPY serviceConf.yml /etc/dockmon/serviceConf.yml

ARG APP_USER=dockmon
RUN echo "Add runtime group and  user.." \
    && addgroup -S ${APP_USER} \
    && adduser -S ${APP_USER} -G ${APP_USER} \
    && echo "Done."

USER ${APP_USER}