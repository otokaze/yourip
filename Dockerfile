FROM scratch
MAINTAINER Otokaze <admin@otokaze.cn>
COPY yourip /usr/local/bin/
EXPOSE 80
CMD ["yourip","--http"]