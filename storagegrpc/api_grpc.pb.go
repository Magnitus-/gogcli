// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package storagegrpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// StorageServiceClient is the client API for StorageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StorageServiceClient interface {
	GetGameIds(ctx context.Context, in *GetGameIdsRequest, opts ...grpc.CallOption) (*GetGameIdsResponse, error)
	GetGameFiles(ctx context.Context, in *GetGameFilesRequest, opts ...grpc.CallOption) (*GetGameFilesResponse, error)
	IsSelfValidating(ctx context.Context, in *IsSelfValidatingRequest, opts ...grpc.CallOption) (*IsSelfValidatingResponse, error)
	GetPrintableSummary(ctx context.Context, in *GetPrintableSummaryRequest, opts ...grpc.CallOption) (*GetPrintableSummaryResponse, error)
	Exists(ctx context.Context, in *ExistsRequest, opts ...grpc.CallOption) (*ExistsResponse, error)
	Initialize(ctx context.Context, in *InitializeRequest, opts ...grpc.CallOption) (*InitializeResponse, error)
	HasManifest(ctx context.Context, in *HasManifestRequest, opts ...grpc.CallOption) (*HasManifestResponse, error)
	HasActions(ctx context.Context, in *HasActionsRequest, opts ...grpc.CallOption) (*HasActionsResponse, error)
	HasSource(ctx context.Context, in *HasSourceRequest, opts ...grpc.CallOption) (*HasSourceResponse, error)
	StoreManifest(ctx context.Context, opts ...grpc.CallOption) (StorageService_StoreManifestClient, error)
	StoreActions(ctx context.Context, opts ...grpc.CallOption) (StorageService_StoreActionsClient, error)
	StoreSource(ctx context.Context, in *StoreSourceRequest, opts ...grpc.CallOption) (*StoreSourceResponse, error)
	LoadManifest(ctx context.Context, in *LoadManifestRequest, opts ...grpc.CallOption) (StorageService_LoadManifestClient, error)
	LoadActions(ctx context.Context, in *LoadActionsRequest, opts ...grpc.CallOption) (StorageService_LoadActionsClient, error)
	LoadSource(ctx context.Context, in *LoadSourceRequest, opts ...grpc.CallOption) (*LoadSourceResponse, error)
	RemoveActions(ctx context.Context, in *RemoveActionsRequest, opts ...grpc.CallOption) (*RemoveActionsResponse, error)
	RemoveSource(ctx context.Context, in *RemoveSourceRequest, opts ...grpc.CallOption) (*RemoveSourceResponse, error)
	AddGame(ctx context.Context, in *AddGameRequest, opts ...grpc.CallOption) (*AddGameResponse, error)
	RemoveGame(ctx context.Context, in *RemoveGameRequest, opts ...grpc.CallOption) (*RemoveGameResponse, error)
	UploadFile(ctx context.Context, opts ...grpc.CallOption) (StorageService_UploadFileClient, error)
	RemoveFile(ctx context.Context, in *RemoveFileRequest, opts ...grpc.CallOption) (*RemoveFileResponse, error)
	DownloadFile(ctx context.Context, in *DownloadFileRequest, opts ...grpc.CallOption) (StorageService_DownloadFileClient, error)
}

type storageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStorageServiceClient(cc grpc.ClientConnInterface) StorageServiceClient {
	return &storageServiceClient{cc}
}

func (c *storageServiceClient) GetGameIds(ctx context.Context, in *GetGameIdsRequest, opts ...grpc.CallOption) (*GetGameIdsResponse, error) {
	out := new(GetGameIdsResponse)
	err := c.cc.Invoke(ctx, "/grpc_storage.StorageService/GetGameIds", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageServiceClient) GetGameFiles(ctx context.Context, in *GetGameFilesRequest, opts ...grpc.CallOption) (*GetGameFilesResponse, error) {
	out := new(GetGameFilesResponse)
	err := c.cc.Invoke(ctx, "/grpc_storage.StorageService/GetGameFiles", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageServiceClient) IsSelfValidating(ctx context.Context, in *IsSelfValidatingRequest, opts ...grpc.CallOption) (*IsSelfValidatingResponse, error) {
	out := new(IsSelfValidatingResponse)
	err := c.cc.Invoke(ctx, "/grpc_storage.StorageService/IsSelfValidating", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageServiceClient) GetPrintableSummary(ctx context.Context, in *GetPrintableSummaryRequest, opts ...grpc.CallOption) (*GetPrintableSummaryResponse, error) {
	out := new(GetPrintableSummaryResponse)
	err := c.cc.Invoke(ctx, "/grpc_storage.StorageService/GetPrintableSummary", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageServiceClient) Exists(ctx context.Context, in *ExistsRequest, opts ...grpc.CallOption) (*ExistsResponse, error) {
	out := new(ExistsResponse)
	err := c.cc.Invoke(ctx, "/grpc_storage.StorageService/Exists", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageServiceClient) Initialize(ctx context.Context, in *InitializeRequest, opts ...grpc.CallOption) (*InitializeResponse, error) {
	out := new(InitializeResponse)
	err := c.cc.Invoke(ctx, "/grpc_storage.StorageService/Initialize", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageServiceClient) HasManifest(ctx context.Context, in *HasManifestRequest, opts ...grpc.CallOption) (*HasManifestResponse, error) {
	out := new(HasManifestResponse)
	err := c.cc.Invoke(ctx, "/grpc_storage.StorageService/HasManifest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageServiceClient) HasActions(ctx context.Context, in *HasActionsRequest, opts ...grpc.CallOption) (*HasActionsResponse, error) {
	out := new(HasActionsResponse)
	err := c.cc.Invoke(ctx, "/grpc_storage.StorageService/HasActions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageServiceClient) HasSource(ctx context.Context, in *HasSourceRequest, opts ...grpc.CallOption) (*HasSourceResponse, error) {
	out := new(HasSourceResponse)
	err := c.cc.Invoke(ctx, "/grpc_storage.StorageService/HasSource", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageServiceClient) StoreManifest(ctx context.Context, opts ...grpc.CallOption) (StorageService_StoreManifestClient, error) {
	stream, err := c.cc.NewStream(ctx, &StorageService_ServiceDesc.Streams[0], "/grpc_storage.StorageService/StoreManifest", opts...)
	if err != nil {
		return nil, err
	}
	x := &storageServiceStoreManifestClient{stream}
	return x, nil
}

type StorageService_StoreManifestClient interface {
	Send(*StoreManifestRequest) error
	CloseAndRecv() (*StoreManifestResponse, error)
	grpc.ClientStream
}

type storageServiceStoreManifestClient struct {
	grpc.ClientStream
}

func (x *storageServiceStoreManifestClient) Send(m *StoreManifestRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *storageServiceStoreManifestClient) CloseAndRecv() (*StoreManifestResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(StoreManifestResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *storageServiceClient) StoreActions(ctx context.Context, opts ...grpc.CallOption) (StorageService_StoreActionsClient, error) {
	stream, err := c.cc.NewStream(ctx, &StorageService_ServiceDesc.Streams[1], "/grpc_storage.StorageService/StoreActions", opts...)
	if err != nil {
		return nil, err
	}
	x := &storageServiceStoreActionsClient{stream}
	return x, nil
}

type StorageService_StoreActionsClient interface {
	Send(*StoreActionsRequest) error
	CloseAndRecv() (*StoreActionsResponse, error)
	grpc.ClientStream
}

type storageServiceStoreActionsClient struct {
	grpc.ClientStream
}

func (x *storageServiceStoreActionsClient) Send(m *StoreActionsRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *storageServiceStoreActionsClient) CloseAndRecv() (*StoreActionsResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(StoreActionsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *storageServiceClient) StoreSource(ctx context.Context, in *StoreSourceRequest, opts ...grpc.CallOption) (*StoreSourceResponse, error) {
	out := new(StoreSourceResponse)
	err := c.cc.Invoke(ctx, "/grpc_storage.StorageService/StoreSource", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageServiceClient) LoadManifest(ctx context.Context, in *LoadManifestRequest, opts ...grpc.CallOption) (StorageService_LoadManifestClient, error) {
	stream, err := c.cc.NewStream(ctx, &StorageService_ServiceDesc.Streams[2], "/grpc_storage.StorageService/LoadManifest", opts...)
	if err != nil {
		return nil, err
	}
	x := &storageServiceLoadManifestClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type StorageService_LoadManifestClient interface {
	Recv() (*LoadManifestResponse, error)
	grpc.ClientStream
}

type storageServiceLoadManifestClient struct {
	grpc.ClientStream
}

func (x *storageServiceLoadManifestClient) Recv() (*LoadManifestResponse, error) {
	m := new(LoadManifestResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *storageServiceClient) LoadActions(ctx context.Context, in *LoadActionsRequest, opts ...grpc.CallOption) (StorageService_LoadActionsClient, error) {
	stream, err := c.cc.NewStream(ctx, &StorageService_ServiceDesc.Streams[3], "/grpc_storage.StorageService/LoadActions", opts...)
	if err != nil {
		return nil, err
	}
	x := &storageServiceLoadActionsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type StorageService_LoadActionsClient interface {
	Recv() (*LoadActionsResponse, error)
	grpc.ClientStream
}

type storageServiceLoadActionsClient struct {
	grpc.ClientStream
}

func (x *storageServiceLoadActionsClient) Recv() (*LoadActionsResponse, error) {
	m := new(LoadActionsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *storageServiceClient) LoadSource(ctx context.Context, in *LoadSourceRequest, opts ...grpc.CallOption) (*LoadSourceResponse, error) {
	out := new(LoadSourceResponse)
	err := c.cc.Invoke(ctx, "/grpc_storage.StorageService/LoadSource", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageServiceClient) RemoveActions(ctx context.Context, in *RemoveActionsRequest, opts ...grpc.CallOption) (*RemoveActionsResponse, error) {
	out := new(RemoveActionsResponse)
	err := c.cc.Invoke(ctx, "/grpc_storage.StorageService/RemoveActions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageServiceClient) RemoveSource(ctx context.Context, in *RemoveSourceRequest, opts ...grpc.CallOption) (*RemoveSourceResponse, error) {
	out := new(RemoveSourceResponse)
	err := c.cc.Invoke(ctx, "/grpc_storage.StorageService/RemoveSource", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageServiceClient) AddGame(ctx context.Context, in *AddGameRequest, opts ...grpc.CallOption) (*AddGameResponse, error) {
	out := new(AddGameResponse)
	err := c.cc.Invoke(ctx, "/grpc_storage.StorageService/AddGame", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageServiceClient) RemoveGame(ctx context.Context, in *RemoveGameRequest, opts ...grpc.CallOption) (*RemoveGameResponse, error) {
	out := new(RemoveGameResponse)
	err := c.cc.Invoke(ctx, "/grpc_storage.StorageService/RemoveGame", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageServiceClient) UploadFile(ctx context.Context, opts ...grpc.CallOption) (StorageService_UploadFileClient, error) {
	stream, err := c.cc.NewStream(ctx, &StorageService_ServiceDesc.Streams[4], "/grpc_storage.StorageService/UploadFile", opts...)
	if err != nil {
		return nil, err
	}
	x := &storageServiceUploadFileClient{stream}
	return x, nil
}

type StorageService_UploadFileClient interface {
	Send(*UploadFileRequest) error
	CloseAndRecv() (*UploadFileResponse, error)
	grpc.ClientStream
}

type storageServiceUploadFileClient struct {
	grpc.ClientStream
}

func (x *storageServiceUploadFileClient) Send(m *UploadFileRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *storageServiceUploadFileClient) CloseAndRecv() (*UploadFileResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadFileResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *storageServiceClient) RemoveFile(ctx context.Context, in *RemoveFileRequest, opts ...grpc.CallOption) (*RemoveFileResponse, error) {
	out := new(RemoveFileResponse)
	err := c.cc.Invoke(ctx, "/grpc_storage.StorageService/RemoveFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageServiceClient) DownloadFile(ctx context.Context, in *DownloadFileRequest, opts ...grpc.CallOption) (StorageService_DownloadFileClient, error) {
	stream, err := c.cc.NewStream(ctx, &StorageService_ServiceDesc.Streams[5], "/grpc_storage.StorageService/DownloadFile", opts...)
	if err != nil {
		return nil, err
	}
	x := &storageServiceDownloadFileClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type StorageService_DownloadFileClient interface {
	Recv() (*DownloadFileResponse, error)
	grpc.ClientStream
}

type storageServiceDownloadFileClient struct {
	grpc.ClientStream
}

func (x *storageServiceDownloadFileClient) Recv() (*DownloadFileResponse, error) {
	m := new(DownloadFileResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// StorageServiceServer is the server API for StorageService service.
// All implementations must embed UnimplementedStorageServiceServer
// for forward compatibility
type StorageServiceServer interface {
	GetGameIds(context.Context, *GetGameIdsRequest) (*GetGameIdsResponse, error)
	GetGameFiles(context.Context, *GetGameFilesRequest) (*GetGameFilesResponse, error)
	IsSelfValidating(context.Context, *IsSelfValidatingRequest) (*IsSelfValidatingResponse, error)
	GetPrintableSummary(context.Context, *GetPrintableSummaryRequest) (*GetPrintableSummaryResponse, error)
	Exists(context.Context, *ExistsRequest) (*ExistsResponse, error)
	Initialize(context.Context, *InitializeRequest) (*InitializeResponse, error)
	HasManifest(context.Context, *HasManifestRequest) (*HasManifestResponse, error)
	HasActions(context.Context, *HasActionsRequest) (*HasActionsResponse, error)
	HasSource(context.Context, *HasSourceRequest) (*HasSourceResponse, error)
	StoreManifest(StorageService_StoreManifestServer) error
	StoreActions(StorageService_StoreActionsServer) error
	StoreSource(context.Context, *StoreSourceRequest) (*StoreSourceResponse, error)
	LoadManifest(*LoadManifestRequest, StorageService_LoadManifestServer) error
	LoadActions(*LoadActionsRequest, StorageService_LoadActionsServer) error
	LoadSource(context.Context, *LoadSourceRequest) (*LoadSourceResponse, error)
	RemoveActions(context.Context, *RemoveActionsRequest) (*RemoveActionsResponse, error)
	RemoveSource(context.Context, *RemoveSourceRequest) (*RemoveSourceResponse, error)
	AddGame(context.Context, *AddGameRequest) (*AddGameResponse, error)
	RemoveGame(context.Context, *RemoveGameRequest) (*RemoveGameResponse, error)
	UploadFile(StorageService_UploadFileServer) error
	RemoveFile(context.Context, *RemoveFileRequest) (*RemoveFileResponse, error)
	DownloadFile(*DownloadFileRequest, StorageService_DownloadFileServer) error
	mustEmbedUnimplementedStorageServiceServer()
}

// UnimplementedStorageServiceServer must be embedded to have forward compatible implementations.
type UnimplementedStorageServiceServer struct {
}

func (UnimplementedStorageServiceServer) GetGameIds(context.Context, *GetGameIdsRequest) (*GetGameIdsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGameIds not implemented")
}
func (UnimplementedStorageServiceServer) GetGameFiles(context.Context, *GetGameFilesRequest) (*GetGameFilesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGameFiles not implemented")
}
func (UnimplementedStorageServiceServer) IsSelfValidating(context.Context, *IsSelfValidatingRequest) (*IsSelfValidatingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsSelfValidating not implemented")
}
func (UnimplementedStorageServiceServer) GetPrintableSummary(context.Context, *GetPrintableSummaryRequest) (*GetPrintableSummaryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPrintableSummary not implemented")
}
func (UnimplementedStorageServiceServer) Exists(context.Context, *ExistsRequest) (*ExistsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Exists not implemented")
}
func (UnimplementedStorageServiceServer) Initialize(context.Context, *InitializeRequest) (*InitializeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Initialize not implemented")
}
func (UnimplementedStorageServiceServer) HasManifest(context.Context, *HasManifestRequest) (*HasManifestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HasManifest not implemented")
}
func (UnimplementedStorageServiceServer) HasActions(context.Context, *HasActionsRequest) (*HasActionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HasActions not implemented")
}
func (UnimplementedStorageServiceServer) HasSource(context.Context, *HasSourceRequest) (*HasSourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HasSource not implemented")
}
func (UnimplementedStorageServiceServer) StoreManifest(StorageService_StoreManifestServer) error {
	return status.Errorf(codes.Unimplemented, "method StoreManifest not implemented")
}
func (UnimplementedStorageServiceServer) StoreActions(StorageService_StoreActionsServer) error {
	return status.Errorf(codes.Unimplemented, "method StoreActions not implemented")
}
func (UnimplementedStorageServiceServer) StoreSource(context.Context, *StoreSourceRequest) (*StoreSourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoreSource not implemented")
}
func (UnimplementedStorageServiceServer) LoadManifest(*LoadManifestRequest, StorageService_LoadManifestServer) error {
	return status.Errorf(codes.Unimplemented, "method LoadManifest not implemented")
}
func (UnimplementedStorageServiceServer) LoadActions(*LoadActionsRequest, StorageService_LoadActionsServer) error {
	return status.Errorf(codes.Unimplemented, "method LoadActions not implemented")
}
func (UnimplementedStorageServiceServer) LoadSource(context.Context, *LoadSourceRequest) (*LoadSourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadSource not implemented")
}
func (UnimplementedStorageServiceServer) RemoveActions(context.Context, *RemoveActionsRequest) (*RemoveActionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveActions not implemented")
}
func (UnimplementedStorageServiceServer) RemoveSource(context.Context, *RemoveSourceRequest) (*RemoveSourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveSource not implemented")
}
func (UnimplementedStorageServiceServer) AddGame(context.Context, *AddGameRequest) (*AddGameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddGame not implemented")
}
func (UnimplementedStorageServiceServer) RemoveGame(context.Context, *RemoveGameRequest) (*RemoveGameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveGame not implemented")
}
func (UnimplementedStorageServiceServer) UploadFile(StorageService_UploadFileServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadFile not implemented")
}
func (UnimplementedStorageServiceServer) RemoveFile(context.Context, *RemoveFileRequest) (*RemoveFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveFile not implemented")
}
func (UnimplementedStorageServiceServer) DownloadFile(*DownloadFileRequest, StorageService_DownloadFileServer) error {
	return status.Errorf(codes.Unimplemented, "method DownloadFile not implemented")
}
func (UnimplementedStorageServiceServer) mustEmbedUnimplementedStorageServiceServer() {}

// UnsafeStorageServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StorageServiceServer will
// result in compilation errors.
type UnsafeStorageServiceServer interface {
	mustEmbedUnimplementedStorageServiceServer()
}

func RegisterStorageServiceServer(s grpc.ServiceRegistrar, srv StorageServiceServer) {
	s.RegisterService(&StorageService_ServiceDesc, srv)
}

func _StorageService_GetGameIds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGameIdsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServiceServer).GetGameIds(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc_storage.StorageService/GetGameIds",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServiceServer).GetGameIds(ctx, req.(*GetGameIdsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StorageService_GetGameFiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGameFilesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServiceServer).GetGameFiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc_storage.StorageService/GetGameFiles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServiceServer).GetGameFiles(ctx, req.(*GetGameFilesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StorageService_IsSelfValidating_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsSelfValidatingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServiceServer).IsSelfValidating(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc_storage.StorageService/IsSelfValidating",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServiceServer).IsSelfValidating(ctx, req.(*IsSelfValidatingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StorageService_GetPrintableSummary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPrintableSummaryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServiceServer).GetPrintableSummary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc_storage.StorageService/GetPrintableSummary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServiceServer).GetPrintableSummary(ctx, req.(*GetPrintableSummaryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StorageService_Exists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExistsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServiceServer).Exists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc_storage.StorageService/Exists",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServiceServer).Exists(ctx, req.(*ExistsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StorageService_Initialize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitializeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServiceServer).Initialize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc_storage.StorageService/Initialize",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServiceServer).Initialize(ctx, req.(*InitializeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StorageService_HasManifest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HasManifestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServiceServer).HasManifest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc_storage.StorageService/HasManifest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServiceServer).HasManifest(ctx, req.(*HasManifestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StorageService_HasActions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HasActionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServiceServer).HasActions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc_storage.StorageService/HasActions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServiceServer).HasActions(ctx, req.(*HasActionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StorageService_HasSource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HasSourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServiceServer).HasSource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc_storage.StorageService/HasSource",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServiceServer).HasSource(ctx, req.(*HasSourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StorageService_StoreManifest_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StorageServiceServer).StoreManifest(&storageServiceStoreManifestServer{stream})
}

type StorageService_StoreManifestServer interface {
	SendAndClose(*StoreManifestResponse) error
	Recv() (*StoreManifestRequest, error)
	grpc.ServerStream
}

type storageServiceStoreManifestServer struct {
	grpc.ServerStream
}

func (x *storageServiceStoreManifestServer) SendAndClose(m *StoreManifestResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *storageServiceStoreManifestServer) Recv() (*StoreManifestRequest, error) {
	m := new(StoreManifestRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _StorageService_StoreActions_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StorageServiceServer).StoreActions(&storageServiceStoreActionsServer{stream})
}

type StorageService_StoreActionsServer interface {
	SendAndClose(*StoreActionsResponse) error
	Recv() (*StoreActionsRequest, error)
	grpc.ServerStream
}

type storageServiceStoreActionsServer struct {
	grpc.ServerStream
}

func (x *storageServiceStoreActionsServer) SendAndClose(m *StoreActionsResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *storageServiceStoreActionsServer) Recv() (*StoreActionsRequest, error) {
	m := new(StoreActionsRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _StorageService_StoreSource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoreSourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServiceServer).StoreSource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc_storage.StorageService/StoreSource",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServiceServer).StoreSource(ctx, req.(*StoreSourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StorageService_LoadManifest_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(LoadManifestRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StorageServiceServer).LoadManifest(m, &storageServiceLoadManifestServer{stream})
}

type StorageService_LoadManifestServer interface {
	Send(*LoadManifestResponse) error
	grpc.ServerStream
}

type storageServiceLoadManifestServer struct {
	grpc.ServerStream
}

func (x *storageServiceLoadManifestServer) Send(m *LoadManifestResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _StorageService_LoadActions_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(LoadActionsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StorageServiceServer).LoadActions(m, &storageServiceLoadActionsServer{stream})
}

type StorageService_LoadActionsServer interface {
	Send(*LoadActionsResponse) error
	grpc.ServerStream
}

type storageServiceLoadActionsServer struct {
	grpc.ServerStream
}

func (x *storageServiceLoadActionsServer) Send(m *LoadActionsResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _StorageService_LoadSource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoadSourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServiceServer).LoadSource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc_storage.StorageService/LoadSource",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServiceServer).LoadSource(ctx, req.(*LoadSourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StorageService_RemoveActions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveActionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServiceServer).RemoveActions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc_storage.StorageService/RemoveActions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServiceServer).RemoveActions(ctx, req.(*RemoveActionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StorageService_RemoveSource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveSourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServiceServer).RemoveSource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc_storage.StorageService/RemoveSource",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServiceServer).RemoveSource(ctx, req.(*RemoveSourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StorageService_AddGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddGameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServiceServer).AddGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc_storage.StorageService/AddGame",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServiceServer).AddGame(ctx, req.(*AddGameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StorageService_RemoveGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveGameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServiceServer).RemoveGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc_storage.StorageService/RemoveGame",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServiceServer).RemoveGame(ctx, req.(*RemoveGameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StorageService_UploadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StorageServiceServer).UploadFile(&storageServiceUploadFileServer{stream})
}

type StorageService_UploadFileServer interface {
	SendAndClose(*UploadFileResponse) error
	Recv() (*UploadFileRequest, error)
	grpc.ServerStream
}

type storageServiceUploadFileServer struct {
	grpc.ServerStream
}

func (x *storageServiceUploadFileServer) SendAndClose(m *UploadFileResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *storageServiceUploadFileServer) Recv() (*UploadFileRequest, error) {
	m := new(UploadFileRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _StorageService_RemoveFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServiceServer).RemoveFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc_storage.StorageService/RemoveFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServiceServer).RemoveFile(ctx, req.(*RemoveFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StorageService_DownloadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadFileRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StorageServiceServer).DownloadFile(m, &storageServiceDownloadFileServer{stream})
}

type StorageService_DownloadFileServer interface {
	Send(*DownloadFileResponse) error
	grpc.ServerStream
}

type storageServiceDownloadFileServer struct {
	grpc.ServerStream
}

func (x *storageServiceDownloadFileServer) Send(m *DownloadFileResponse) error {
	return x.ServerStream.SendMsg(m)
}

// StorageService_ServiceDesc is the grpc.ServiceDesc for StorageService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StorageService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc_storage.StorageService",
	HandlerType: (*StorageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetGameIds",
			Handler:    _StorageService_GetGameIds_Handler,
		},
		{
			MethodName: "GetGameFiles",
			Handler:    _StorageService_GetGameFiles_Handler,
		},
		{
			MethodName: "IsSelfValidating",
			Handler:    _StorageService_IsSelfValidating_Handler,
		},
		{
			MethodName: "GetPrintableSummary",
			Handler:    _StorageService_GetPrintableSummary_Handler,
		},
		{
			MethodName: "Exists",
			Handler:    _StorageService_Exists_Handler,
		},
		{
			MethodName: "Initialize",
			Handler:    _StorageService_Initialize_Handler,
		},
		{
			MethodName: "HasManifest",
			Handler:    _StorageService_HasManifest_Handler,
		},
		{
			MethodName: "HasActions",
			Handler:    _StorageService_HasActions_Handler,
		},
		{
			MethodName: "HasSource",
			Handler:    _StorageService_HasSource_Handler,
		},
		{
			MethodName: "StoreSource",
			Handler:    _StorageService_StoreSource_Handler,
		},
		{
			MethodName: "LoadSource",
			Handler:    _StorageService_LoadSource_Handler,
		},
		{
			MethodName: "RemoveActions",
			Handler:    _StorageService_RemoveActions_Handler,
		},
		{
			MethodName: "RemoveSource",
			Handler:    _StorageService_RemoveSource_Handler,
		},
		{
			MethodName: "AddGame",
			Handler:    _StorageService_AddGame_Handler,
		},
		{
			MethodName: "RemoveGame",
			Handler:    _StorageService_RemoveGame_Handler,
		},
		{
			MethodName: "RemoveFile",
			Handler:    _StorageService_RemoveFile_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StoreManifest",
			Handler:       _StorageService_StoreManifest_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "StoreActions",
			Handler:       _StorageService_StoreActions_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "LoadManifest",
			Handler:       _StorageService_LoadManifest_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "LoadActions",
			Handler:       _StorageService_LoadActions_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "UploadFile",
			Handler:       _StorageService_UploadFile_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "DownloadFile",
			Handler:       _StorageService_DownloadFile_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api.proto",
}
