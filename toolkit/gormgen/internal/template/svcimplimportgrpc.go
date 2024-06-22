package template

const SvcImplImportGrpc = `
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/copier"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "{{.TransportGrpcPackage}}"
`
