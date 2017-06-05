FROM segment/base:v4

RUN curl -Ls https://github.com/segmentio/http_to_nsq/releases/download/0.2.0/http_to_nsq_linux_amd64 > /http_to_nsq

RUN chmod +x /http_to_nsq

ENTRYPOINT ["/http_to_nsq"]
