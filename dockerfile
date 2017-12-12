FROM ubuntu

# container labels
LABEL author="Rogier Lommers <rogier@lommers.org>"
LABEL description="pingback server"

# install dependencies
RUN apt-get update  
RUN apt-get install -y ca-certificates

# add binary
COPY bin/pingback-linux-amd64 /app/

# expose port
EXPOSE 8080

# change to data dir and run bianry
WORKDIR "/app"
CMD ["/app/pingback-linux-amd64", "-mode", "server"]
