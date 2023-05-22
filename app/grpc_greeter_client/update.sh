#!/bin/bash

appname="grpc_greeter_client.out"
go build -gcflags="all=-N -l" -o "$appname"
# shellcheck disable=SC2164
pid=$(pgrep $appname)
if [ "$pid" = "" ]; then
	./$appname
else
#	kill -SIGTERM "${pid}"
	kill -SIGUSR2 "${pid}"
fi
