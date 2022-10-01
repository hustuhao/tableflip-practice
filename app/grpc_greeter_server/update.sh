#!/bin/bash

appname="grpc_greeter_server.out"
go build -o "$appname"
# shellcheck disable=SC2164
pid=$(pgrep $appname)
if [ "$pid" = "" ]; then
	./$appname
else
	kill -SIGUSR2 "${pid}"
fi
