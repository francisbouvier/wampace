FROM microbox/scratch

ADD https://s3-us-west-1.amazonaws.com/pipesdocker/wampace.v0.1.1 bin/wampace
RUN chmod 0755 bin/wampace

ENTRYPOINT ["wampace"]
