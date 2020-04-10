FROM centos:7

COPY conf/app.conf.sample /PPGo_ApiAdmin/conf/app.conf
COPY static /PPGo_ApiAdmin/static
COPY views /PPGo_ApiAdmin/views
COPY PPGo_ApiAdmin /PPGo_ApiAdmin/PPGo_ApiAdmin

WORKDIR /PPGo_ApiAdmin

RUN chmod +x /PPGo_ApiAdmin/PPGo_ApiAdmin

CMD ["/PPGo_ApiAdmin/PPGo_ApiAdmin"]