# tableflip-practice

some example with [tableflip](https://github.com/cloudflare/tableflip)


## 1.categories
- app/cron: show how to use tableflip with cron server
- app/grpc_greeter_server:  show how to use tableflip with gRPC server.

## 2.have a try

#### 2.1 test cron server with tableflip
Enter the directory:
> cd app/cron/

Start the cron service:
> ./update.sh

Update the cron service with the same sh file:
> ./update.sh

You can check the log at log/cron/

#### 2.2 test gRPC server with tableflip
Enter the directory of gRPC server:
> cd app/grpc_greeter_server/

Start the gRPC server:
> ./update.sh

Enter the directory of gRPC client:
> cd app/grpc_greeter_client/

Start the gRPC client:
> ./update.sh
