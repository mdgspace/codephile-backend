FROM library/golang

MAINTAINER Mobile Development Group <mdg@iitr.ac.in>

RUN apt-get update && apt-get -y install supervisor

# Recompile the standard library without CGO
RUN CGO_ENABLED=0 go install -a std

ENV APP_DIR $GOPATH/src/github.com/mdg-iitr/Codephile

RUN mkdir -p $APP_DIR

WORKDIR $APP_DIR

COPY . .

COPY supervisord.conf /etc/supervisor/conf.d/

# Compile the binary and statically link
RUN cd $APP_DIR && CGO_ENABLED=0 go build -ldflags '-d -w -s'

ENTRYPOINT ["supervisord", "-n"]

EXPOSE 8080
