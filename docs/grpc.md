# gRPC implentation


 Metamorph API and Controller communicate each other using gRPC proto
 
 
 we can divide gRPC implementation to 3 parts
 
 1. Define service in proto 
 
 	  	Refer: `metamorph/proto/metamorph.proto` . Once we have the .proto file ready 
		next step is convert it into go code, so that other components can use it
		
 2. Write Server componenet
     Refer: `metamorph/pkg/grpc/controller/server.go`
	 
 3. Write Client Component 
     Refer: `metamorph/pkg/grpc/controller/client.go`
	 
	 
	 
## Dev environment setup

   1. run grpc server component in one terminal like this `go run metamorph/pkg/grpc/controller/client.go`
   2. run grpc client component in another terminal like this `go run metamorph/pkg/grpc/controller/client.go`
   3. User browser or postman to access the client at `localhost:8080`
   
   
   	  Example request `http://localhost:8080/node/2B0110D0-B648-4BB6-9F4F-92BDE0245892`
	  
	  Example response:
	  
	  			`{"result":"2B0110D0-B648-4BB6-9F4F-92BDE0245892"}`
