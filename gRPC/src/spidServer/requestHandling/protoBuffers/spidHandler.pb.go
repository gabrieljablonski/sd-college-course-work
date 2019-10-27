// Code generated by protoc-gen-go. DO NOT EDIT.
// source: spidHandler.proto

package protoBuffers

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Spid struct {
	Id                   string          `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	BatteryLevel         uint32          `protobuf:"varint,2,opt,name=batteryLevel,proto3" json:"batteryLevel,omitempty"`
	LockInfo             string          `protobuf:"bytes,3,opt,name=lockInfo,proto3" json:"lockInfo,omitempty"`
	Location             *GlobalPosition `protobuf:"bytes,4,opt,name=location,proto3" json:"location,omitempty"`
	LastUpdated          string          `protobuf:"bytes,5,opt,name=lastUpdated,proto3" json:"lastUpdated,omitempty"`
	CurrentUserID        string          `protobuf:"bytes,6,opt,name=currentUserID,proto3" json:"currentUserID,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Spid) Reset()         { *m = Spid{} }
func (m *Spid) String() string { return proto.CompactTextString(m) }
func (*Spid) ProtoMessage()    {}
func (*Spid) Descriptor() ([]byte, []int) {
	return fileDescriptor_97215654e7c7179a, []int{0}
}

func (m *Spid) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Spid.Unmarshal(m, b)
}
func (m *Spid) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Spid.Marshal(b, m, deterministic)
}
func (m *Spid) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Spid.Merge(m, src)
}
func (m *Spid) XXX_Size() int {
	return xxx_messageInfo_Spid.Size(m)
}
func (m *Spid) XXX_DiscardUnknown() {
	xxx_messageInfo_Spid.DiscardUnknown(m)
}

var xxx_messageInfo_Spid proto.InternalMessageInfo

func (m *Spid) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Spid) GetBatteryLevel() uint32 {
	if m != nil {
		return m.BatteryLevel
	}
	return 0
}

func (m *Spid) GetLockInfo() string {
	if m != nil {
		return m.LockInfo
	}
	return ""
}

func (m *Spid) GetLocation() *GlobalPosition {
	if m != nil {
		return m.Location
	}
	return nil
}

func (m *Spid) GetLastUpdated() string {
	if m != nil {
		return m.LastUpdated
	}
	return ""
}

func (m *Spid) GetCurrentUserID() string {
	if m != nil {
		return m.CurrentUserID
	}
	return ""
}

type SpidMinimal struct {
	Id                   string          `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	BatteryLevel         uint32          `protobuf:"varint,2,opt,name=batteryLevel,proto3" json:"batteryLevel,omitempty"`
	Location             *GlobalPosition `protobuf:"bytes,3,opt,name=location,proto3" json:"location,omitempty"`
	LockState            string          `protobuf:"bytes,4,opt,name=lockState,proto3" json:"lockState,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *SpidMinimal) Reset()         { *m = SpidMinimal{} }
func (m *SpidMinimal) String() string { return proto.CompactTextString(m) }
func (*SpidMinimal) ProtoMessage()    {}
func (*SpidMinimal) Descriptor() ([]byte, []int) {
	return fileDescriptor_97215654e7c7179a, []int{1}
}

func (m *SpidMinimal) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SpidMinimal.Unmarshal(m, b)
}
func (m *SpidMinimal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SpidMinimal.Marshal(b, m, deterministic)
}
func (m *SpidMinimal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SpidMinimal.Merge(m, src)
}
func (m *SpidMinimal) XXX_Size() int {
	return xxx_messageInfo_SpidMinimal.Size(m)
}
func (m *SpidMinimal) XXX_DiscardUnknown() {
	xxx_messageInfo_SpidMinimal.DiscardUnknown(m)
}

var xxx_messageInfo_SpidMinimal proto.InternalMessageInfo

func (m *SpidMinimal) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *SpidMinimal) GetBatteryLevel() uint32 {
	if m != nil {
		return m.BatteryLevel
	}
	return 0
}

func (m *SpidMinimal) GetLocation() *GlobalPosition {
	if m != nil {
		return m.Location
	}
	return nil
}

func (m *SpidMinimal) GetLockState() string {
	if m != nil {
		return m.LockState
	}
	return ""
}

type GetSpidRequest struct {
	SpidID               string   `protobuf:"bytes,1,opt,name=spidID,proto3" json:"spidID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetSpidRequest) Reset()         { *m = GetSpidRequest{} }
func (m *GetSpidRequest) String() string { return proto.CompactTextString(m) }
func (*GetSpidRequest) ProtoMessage()    {}
func (*GetSpidRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_97215654e7c7179a, []int{2}
}

func (m *GetSpidRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSpidRequest.Unmarshal(m, b)
}
func (m *GetSpidRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSpidRequest.Marshal(b, m, deterministic)
}
func (m *GetSpidRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSpidRequest.Merge(m, src)
}
func (m *GetSpidRequest) XXX_Size() int {
	return xxx_messageInfo_GetSpidRequest.Size(m)
}
func (m *GetSpidRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSpidRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetSpidRequest proto.InternalMessageInfo

func (m *GetSpidRequest) GetSpidID() string {
	if m != nil {
		return m.SpidID
	}
	return ""
}

type GetSpidResponse struct {
	Message              string       `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	Spid                 *SpidMinimal `protobuf:"bytes,2,opt,name=spid,proto3" json:"spid,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *GetSpidResponse) Reset()         { *m = GetSpidResponse{} }
func (m *GetSpidResponse) String() string { return proto.CompactTextString(m) }
func (*GetSpidResponse) ProtoMessage()    {}
func (*GetSpidResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_97215654e7c7179a, []int{3}
}

func (m *GetSpidResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSpidResponse.Unmarshal(m, b)
}
func (m *GetSpidResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSpidResponse.Marshal(b, m, deterministic)
}
func (m *GetSpidResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSpidResponse.Merge(m, src)
}
func (m *GetSpidResponse) XXX_Size() int {
	return xxx_messageInfo_GetSpidResponse.Size(m)
}
func (m *GetSpidResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSpidResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetSpidResponse proto.InternalMessageInfo

func (m *GetSpidResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *GetSpidResponse) GetSpid() *SpidMinimal {
	if m != nil {
		return m.Spid
	}
	return nil
}

type RegisterSpidRequest struct {
	BatteryLevel         uint32          `protobuf:"varint,2,opt,name=batteryLevel,proto3" json:"batteryLevel,omitempty"`
	Location             *GlobalPosition `protobuf:"bytes,1,opt,name=location,proto3" json:"location,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *RegisterSpidRequest) Reset()         { *m = RegisterSpidRequest{} }
func (m *RegisterSpidRequest) String() string { return proto.CompactTextString(m) }
func (*RegisterSpidRequest) ProtoMessage()    {}
func (*RegisterSpidRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_97215654e7c7179a, []int{4}
}

func (m *RegisterSpidRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterSpidRequest.Unmarshal(m, b)
}
func (m *RegisterSpidRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterSpidRequest.Marshal(b, m, deterministic)
}
func (m *RegisterSpidRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterSpidRequest.Merge(m, src)
}
func (m *RegisterSpidRequest) XXX_Size() int {
	return xxx_messageInfo_RegisterSpidRequest.Size(m)
}
func (m *RegisterSpidRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterSpidRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterSpidRequest proto.InternalMessageInfo

func (m *RegisterSpidRequest) GetBatteryLevel() uint32 {
	if m != nil {
		return m.BatteryLevel
	}
	return 0
}

func (m *RegisterSpidRequest) GetLocation() *GlobalPosition {
	if m != nil {
		return m.Location
	}
	return nil
}

type RegisterSpidResponse struct {
	Message              string       `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	Spid                 *SpidMinimal `protobuf:"bytes,2,opt,name=spid,proto3" json:"spid,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *RegisterSpidResponse) Reset()         { *m = RegisterSpidResponse{} }
func (m *RegisterSpidResponse) String() string { return proto.CompactTextString(m) }
func (*RegisterSpidResponse) ProtoMessage()    {}
func (*RegisterSpidResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_97215654e7c7179a, []int{5}
}

func (m *RegisterSpidResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterSpidResponse.Unmarshal(m, b)
}
func (m *RegisterSpidResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterSpidResponse.Marshal(b, m, deterministic)
}
func (m *RegisterSpidResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterSpidResponse.Merge(m, src)
}
func (m *RegisterSpidResponse) XXX_Size() int {
	return xxx_messageInfo_RegisterSpidResponse.Size(m)
}
func (m *RegisterSpidResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterSpidResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterSpidResponse proto.InternalMessageInfo

func (m *RegisterSpidResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *RegisterSpidResponse) GetSpid() *SpidMinimal {
	if m != nil {
		return m.Spid
	}
	return nil
}

type UpdateBatteryRequest struct {
	SpidID               string   `protobuf:"bytes,1,opt,name=spidID,proto3" json:"spidID,omitempty"`
	BatteryLevel         uint32   `protobuf:"varint,2,opt,name=batteryLevel,proto3" json:"batteryLevel,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateBatteryRequest) Reset()         { *m = UpdateBatteryRequest{} }
func (m *UpdateBatteryRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateBatteryRequest) ProtoMessage()    {}
func (*UpdateBatteryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_97215654e7c7179a, []int{6}
}

func (m *UpdateBatteryRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateBatteryRequest.Unmarshal(m, b)
}
func (m *UpdateBatteryRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateBatteryRequest.Marshal(b, m, deterministic)
}
func (m *UpdateBatteryRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateBatteryRequest.Merge(m, src)
}
func (m *UpdateBatteryRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateBatteryRequest.Size(m)
}
func (m *UpdateBatteryRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateBatteryRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateBatteryRequest proto.InternalMessageInfo

func (m *UpdateBatteryRequest) GetSpidID() string {
	if m != nil {
		return m.SpidID
	}
	return ""
}

func (m *UpdateBatteryRequest) GetBatteryLevel() uint32 {
	if m != nil {
		return m.BatteryLevel
	}
	return 0
}

type UpdateBatteryResponse struct {
	Message              string       `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	Spid                 *SpidMinimal `protobuf:"bytes,2,opt,name=spid,proto3" json:"spid,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *UpdateBatteryResponse) Reset()         { *m = UpdateBatteryResponse{} }
func (m *UpdateBatteryResponse) String() string { return proto.CompactTextString(m) }
func (*UpdateBatteryResponse) ProtoMessage()    {}
func (*UpdateBatteryResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_97215654e7c7179a, []int{7}
}

func (m *UpdateBatteryResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateBatteryResponse.Unmarshal(m, b)
}
func (m *UpdateBatteryResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateBatteryResponse.Marshal(b, m, deterministic)
}
func (m *UpdateBatteryResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateBatteryResponse.Merge(m, src)
}
func (m *UpdateBatteryResponse) XXX_Size() int {
	return xxx_messageInfo_UpdateBatteryResponse.Size(m)
}
func (m *UpdateBatteryResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateBatteryResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateBatteryResponse proto.InternalMessageInfo

func (m *UpdateBatteryResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *UpdateBatteryResponse) GetSpid() *SpidMinimal {
	if m != nil {
		return m.Spid
	}
	return nil
}

type UpdateSpidLocationRequest struct {
	SpidID               string          `protobuf:"bytes,1,opt,name=spidID,proto3" json:"spidID,omitempty"`
	Location             *GlobalPosition `protobuf:"bytes,2,opt,name=location,proto3" json:"location,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *UpdateSpidLocationRequest) Reset()         { *m = UpdateSpidLocationRequest{} }
func (m *UpdateSpidLocationRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateSpidLocationRequest) ProtoMessage()    {}
func (*UpdateSpidLocationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_97215654e7c7179a, []int{8}
}

func (m *UpdateSpidLocationRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateSpidLocationRequest.Unmarshal(m, b)
}
func (m *UpdateSpidLocationRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateSpidLocationRequest.Marshal(b, m, deterministic)
}
func (m *UpdateSpidLocationRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateSpidLocationRequest.Merge(m, src)
}
func (m *UpdateSpidLocationRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateSpidLocationRequest.Size(m)
}
func (m *UpdateSpidLocationRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateSpidLocationRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateSpidLocationRequest proto.InternalMessageInfo

func (m *UpdateSpidLocationRequest) GetSpidID() string {
	if m != nil {
		return m.SpidID
	}
	return ""
}

func (m *UpdateSpidLocationRequest) GetLocation() *GlobalPosition {
	if m != nil {
		return m.Location
	}
	return nil
}

type UpdateSpidLocationResponse struct {
	Message              string       `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	Spid                 *SpidMinimal `protobuf:"bytes,2,opt,name=spid,proto3" json:"spid,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *UpdateSpidLocationResponse) Reset()         { *m = UpdateSpidLocationResponse{} }
func (m *UpdateSpidLocationResponse) String() string { return proto.CompactTextString(m) }
func (*UpdateSpidLocationResponse) ProtoMessage()    {}
func (*UpdateSpidLocationResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_97215654e7c7179a, []int{9}
}

func (m *UpdateSpidLocationResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateSpidLocationResponse.Unmarshal(m, b)
}
func (m *UpdateSpidLocationResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateSpidLocationResponse.Marshal(b, m, deterministic)
}
func (m *UpdateSpidLocationResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateSpidLocationResponse.Merge(m, src)
}
func (m *UpdateSpidLocationResponse) XXX_Size() int {
	return xxx_messageInfo_UpdateSpidLocationResponse.Size(m)
}
func (m *UpdateSpidLocationResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateSpidLocationResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateSpidLocationResponse proto.InternalMessageInfo

func (m *UpdateSpidLocationResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *UpdateSpidLocationResponse) GetSpid() *SpidMinimal {
	if m != nil {
		return m.Spid
	}
	return nil
}

type DeleteSpidRequest struct {
	SpidID               string   `protobuf:"bytes,1,opt,name=spidID,proto3" json:"spidID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteSpidRequest) Reset()         { *m = DeleteSpidRequest{} }
func (m *DeleteSpidRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteSpidRequest) ProtoMessage()    {}
func (*DeleteSpidRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_97215654e7c7179a, []int{10}
}

func (m *DeleteSpidRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteSpidRequest.Unmarshal(m, b)
}
func (m *DeleteSpidRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteSpidRequest.Marshal(b, m, deterministic)
}
func (m *DeleteSpidRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteSpidRequest.Merge(m, src)
}
func (m *DeleteSpidRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteSpidRequest.Size(m)
}
func (m *DeleteSpidRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteSpidRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteSpidRequest proto.InternalMessageInfo

func (m *DeleteSpidRequest) GetSpidID() string {
	if m != nil {
		return m.SpidID
	}
	return ""
}

type DeleteSpidResponse struct {
	Message              string       `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	Spid                 *SpidMinimal `protobuf:"bytes,2,opt,name=spid,proto3" json:"spid,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *DeleteSpidResponse) Reset()         { *m = DeleteSpidResponse{} }
func (m *DeleteSpidResponse) String() string { return proto.CompactTextString(m) }
func (*DeleteSpidResponse) ProtoMessage()    {}
func (*DeleteSpidResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_97215654e7c7179a, []int{11}
}

func (m *DeleteSpidResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteSpidResponse.Unmarshal(m, b)
}
func (m *DeleteSpidResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteSpidResponse.Marshal(b, m, deterministic)
}
func (m *DeleteSpidResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteSpidResponse.Merge(m, src)
}
func (m *DeleteSpidResponse) XXX_Size() int {
	return xxx_messageInfo_DeleteSpidResponse.Size(m)
}
func (m *DeleteSpidResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteSpidResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteSpidResponse proto.InternalMessageInfo

func (m *DeleteSpidResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *DeleteSpidResponse) GetSpid() *SpidMinimal {
	if m != nil {
		return m.Spid
	}
	return nil
}

func init() {
	proto.RegisterType((*Spid)(nil), "spidProtoBuffers.Spid")
	proto.RegisterType((*SpidMinimal)(nil), "spidProtoBuffers.SpidMinimal")
	proto.RegisterType((*GetSpidRequest)(nil), "spidProtoBuffers.GetSpidRequest")
	proto.RegisterType((*GetSpidResponse)(nil), "spidProtoBuffers.GetSpidResponse")
	proto.RegisterType((*RegisterSpidRequest)(nil), "spidProtoBuffers.RegisterSpidRequest")
	proto.RegisterType((*RegisterSpidResponse)(nil), "spidProtoBuffers.RegisterSpidResponse")
	proto.RegisterType((*UpdateBatteryRequest)(nil), "spidProtoBuffers.UpdateBatteryRequest")
	proto.RegisterType((*UpdateBatteryResponse)(nil), "spidProtoBuffers.UpdateBatteryResponse")
	proto.RegisterType((*UpdateSpidLocationRequest)(nil), "spidProtoBuffers.UpdateSpidLocationRequest")
	proto.RegisterType((*UpdateSpidLocationResponse)(nil), "spidProtoBuffers.UpdateSpidLocationResponse")
	proto.RegisterType((*DeleteSpidRequest)(nil), "spidProtoBuffers.DeleteSpidRequest")
	proto.RegisterType((*DeleteSpidResponse)(nil), "spidProtoBuffers.DeleteSpidResponse")
}

func init() { proto.RegisterFile("spidHandler.proto", fileDescriptor_97215654e7c7179a) }

var fileDescriptor_97215654e7c7179a = []byte{
	// 522 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x55, 0x41, 0x6f, 0xd3, 0x4c,
	0x10, 0xfd, 0x9c, 0xe4, 0x0b, 0x64, 0xdc, 0x16, 0xb2, 0x14, 0x64, 0x2c, 0x90, 0xcc, 0x52, 0x20,
	0x52, 0x51, 0x24, 0xc2, 0x95, 0x53, 0x15, 0xa9, 0x44, 0x2a, 0x52, 0xe5, 0xd2, 0x03, 0x17, 0xd0,
	0x26, 0x1e, 0x47, 0x2b, 0x5c, 0xdb, 0xd9, 0xdd, 0x80, 0xf8, 0x31, 0xfc, 0x33, 0x7e, 0x09, 0x27,
	0xb4, 0x6b, 0xd7, 0xb5, 0x6b, 0x53, 0x37, 0x48, 0xb9, 0x65, 0xc7, 0x6f, 0x66, 0xde, 0x7b, 0x3b,
	0xb3, 0x81, 0xa1, 0x4c, 0x79, 0xf0, 0x9e, 0xc5, 0x41, 0x84, 0x62, 0x9c, 0x8a, 0x44, 0x25, 0xe4,
	0xbe, 0x0e, 0x9d, 0xea, 0x9f, 0x47, 0xeb, 0x30, 0x44, 0x21, 0x5d, 0x67, 0x89, 0x31, 0x0a, 0xbe,
	0x98, 0x62, 0xc8, 0x63, 0xae, 0x78, 0x12, 0xcb, 0x0c, 0x4b, 0x7f, 0x59, 0xd0, 0x3b, 0x4b, 0x79,
	0x40, 0xf6, 0xa0, 0xc3, 0x03, 0xc7, 0xf2, 0xac, 0xd1, 0xc0, 0xef, 0xf0, 0x80, 0x50, 0xd8, 0x99,
	0x33, 0xa5, 0x50, 0xfc, 0x38, 0xc1, 0x6f, 0x18, 0x39, 0x1d, 0xcf, 0x1a, 0xed, 0xfa, 0x95, 0x18,
	0x71, 0xe1, 0x6e, 0x94, 0x2c, 0xbe, 0xce, 0xe2, 0x30, 0x71, 0xba, 0x26, 0xb3, 0x38, 0x93, 0x77,
	0xe6, 0x1b, 0xd3, 0xbd, 0x9c, 0x9e, 0x67, 0x8d, 0xec, 0x89, 0x37, 0xbe, 0xce, 0x6b, 0x7c, 0x1c,
	0x25, 0x73, 0x16, 0x9d, 0x26, 0xd2, 0x70, 0xf2, 0x8b, 0x0c, 0xe2, 0x81, 0x1d, 0x31, 0xa9, 0xce,
	0xd3, 0x80, 0x29, 0x0c, 0x9c, 0xff, 0x4d, 0xf1, 0x72, 0x88, 0x1c, 0xc0, 0xee, 0x62, 0x2d, 0x04,
	0xc6, 0xea, 0x5c, 0xa2, 0x98, 0x4d, 0x9d, 0xbe, 0xc1, 0x54, 0x83, 0xf4, 0xa7, 0x05, 0xb6, 0x96,
	0xf7, 0x81, 0xc7, 0xfc, 0x82, 0x45, 0xff, 0xa4, 0xb2, 0xac, 0xa4, 0xbb, 0xb1, 0x92, 0x27, 0x30,
	0xd0, 0x9e, 0x9c, 0x29, 0xa6, 0xd0, 0x18, 0x31, 0xf0, 0xaf, 0x02, 0x74, 0x04, 0x7b, 0xc7, 0xa8,
	0x34, 0x43, 0x1f, 0x57, 0x6b, 0x94, 0x8a, 0x3c, 0x82, 0xbe, 0x2e, 0x3e, 0x9b, 0xe6, 0x2c, 0xf3,
	0x13, 0xfd, 0x0c, 0xf7, 0x0a, 0xa4, 0x4c, 0x93, 0x58, 0x22, 0x71, 0xe0, 0xce, 0x05, 0x4a, 0xc9,
	0x96, 0x98, 0x63, 0x2f, 0x8f, 0xe4, 0x0d, 0xf4, 0x74, 0x9a, 0x91, 0x63, 0x4f, 0x9e, 0xd6, 0xe9,
	0x96, 0x3c, 0xf1, 0x0d, 0x94, 0x7e, 0x87, 0x07, 0x3e, 0x2e, 0xb9, 0x54, 0x28, 0xca, 0x74, 0x36,
	0x35, 0xc8, 0xda, 0xd4, 0x20, 0xba, 0x80, 0xfd, 0x6a, 0xe3, 0x6d, 0xa8, 0xf3, 0x61, 0x3f, 0x1b,
	0x9c, 0xa3, 0x8c, 0x78, 0x8b, 0xdb, 0xb7, 0x91, 0x4d, 0x03, 0x78, 0x78, 0xad, 0xe6, 0x36, 0x98,
	0xaf, 0xe0, 0x71, 0xd6, 0x45, 0x7f, 0x3a, 0xc9, 0x4d, 0x6b, 0xa3, 0x5f, 0xbe, 0x91, 0xce, 0xc6,
	0x37, 0xc2, 0xc1, 0x6d, 0x6a, 0xb9, 0x0d, 0x75, 0x87, 0x30, 0x9c, 0x62, 0x84, 0x59, 0xab, 0xb6,
	0x15, 0x60, 0x40, 0xca, 0xe0, 0x2d, 0xf0, 0x99, 0xfc, 0xee, 0x66, 0xef, 0x45, 0xfe, 0xa0, 0x92,
	0x8f, 0x60, 0xe7, 0x5b, 0x67, 0x1e, 0xb5, 0x26, 0x17, 0x2b, 0xeb, 0xeb, 0x3e, 0xbb, 0x01, 0x91,
	0x11, 0xa6, 0xff, 0x91, 0x2f, 0xb0, 0x53, 0x1e, 0x79, 0xf2, 0xa2, 0x9e, 0xd4, 0xb0, 0x8b, 0xee,
	0xcb, 0x36, 0x58, 0xd1, 0x20, 0x84, 0x61, 0x65, 0x34, 0x0d, 0xf9, 0x86, 0xf4, 0xa6, 0x9d, 0x70,
	0x5f, 0xb5, 0xe2, 0x8a, 0x3e, 0x2b, 0x20, 0xf5, 0x49, 0x21, 0x87, 0x7f, 0x2b, 0xd0, 0x30, 0xc2,
	0xee, 0xeb, 0xdb, 0x81, 0x8b, 0x96, 0x9f, 0x00, 0xae, 0x86, 0x80, 0x3c, 0xaf, 0x67, 0xd7, 0xe6,
	0xc9, 0x3d, 0xb8, 0x19, 0x74, 0x59, 0x7a, 0xde, 0x37, 0x7f, 0x89, 0x6f, 0xff, 0x04, 0x00, 0x00,
	0xff, 0xff, 0x94, 0x9f, 0x39, 0xe7, 0x53, 0x07, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SpidHandlerClient is the client API for SpidHandler service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SpidHandlerClient interface {
	GetSpidInfo(ctx context.Context, in *GetSpidRequest, opts ...grpc.CallOption) (*GetSpidResponse, error)
	RegisterSpid(ctx context.Context, in *RegisterSpidRequest, opts ...grpc.CallOption) (*RegisterSpidResponse, error)
	UpdateBatteryInfo(ctx context.Context, in *UpdateBatteryRequest, opts ...grpc.CallOption) (*UpdateBatteryResponse, error)
	UpdateSpidLocation(ctx context.Context, in *UpdateSpidLocationRequest, opts ...grpc.CallOption) (*UpdateSpidLocationResponse, error)
	DeleteSpid(ctx context.Context, in *DeleteSpidRequest, opts ...grpc.CallOption) (*DeleteSpidResponse, error)
}

type spidHandlerClient struct {
	cc *grpc.ClientConn
}

func NewSpidHandlerClient(cc *grpc.ClientConn) SpidHandlerClient {
	return &spidHandlerClient{cc}
}

func (c *spidHandlerClient) GetSpidInfo(ctx context.Context, in *GetSpidRequest, opts ...grpc.CallOption) (*GetSpidResponse, error) {
	out := new(GetSpidResponse)
	err := c.cc.Invoke(ctx, "/spidProtoBuffers.SpidHandler/GetSpidInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *spidHandlerClient) RegisterSpid(ctx context.Context, in *RegisterSpidRequest, opts ...grpc.CallOption) (*RegisterSpidResponse, error) {
	out := new(RegisterSpidResponse)
	err := c.cc.Invoke(ctx, "/spidProtoBuffers.SpidHandler/RegisterSpid", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *spidHandlerClient) UpdateBatteryInfo(ctx context.Context, in *UpdateBatteryRequest, opts ...grpc.CallOption) (*UpdateBatteryResponse, error) {
	out := new(UpdateBatteryResponse)
	err := c.cc.Invoke(ctx, "/spidProtoBuffers.SpidHandler/UpdateBatteryInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *spidHandlerClient) UpdateSpidLocation(ctx context.Context, in *UpdateSpidLocationRequest, opts ...grpc.CallOption) (*UpdateSpidLocationResponse, error) {
	out := new(UpdateSpidLocationResponse)
	err := c.cc.Invoke(ctx, "/spidProtoBuffers.SpidHandler/UpdateSpidLocation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *spidHandlerClient) DeleteSpid(ctx context.Context, in *DeleteSpidRequest, opts ...grpc.CallOption) (*DeleteSpidResponse, error) {
	out := new(DeleteSpidResponse)
	err := c.cc.Invoke(ctx, "/spidProtoBuffers.SpidHandler/DeleteSpid", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SpidHandlerServer is the server API for SpidHandler service.
type SpidHandlerServer interface {
	GetSpidInfo(context.Context, *GetSpidRequest) (*GetSpidResponse, error)
	RegisterSpid(context.Context, *RegisterSpidRequest) (*RegisterSpidResponse, error)
	UpdateBatteryInfo(context.Context, *UpdateBatteryRequest) (*UpdateBatteryResponse, error)
	UpdateSpidLocation(context.Context, *UpdateSpidLocationRequest) (*UpdateSpidLocationResponse, error)
	DeleteSpid(context.Context, *DeleteSpidRequest) (*DeleteSpidResponse, error)
}

// UnimplementedSpidHandlerServer can be embedded to have forward compatible implementations.
type UnimplementedSpidHandlerServer struct {
}

func (*UnimplementedSpidHandlerServer) GetSpidInfo(ctx context.Context, req *GetSpidRequest) (*GetSpidResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSpidInfo not implemented")
}
func (*UnimplementedSpidHandlerServer) RegisterSpid(ctx context.Context, req *RegisterSpidRequest) (*RegisterSpidResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterSpid not implemented")
}
func (*UnimplementedSpidHandlerServer) UpdateBatteryInfo(ctx context.Context, req *UpdateBatteryRequest) (*UpdateBatteryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateBatteryInfo not implemented")
}
func (*UnimplementedSpidHandlerServer) UpdateSpidLocation(ctx context.Context, req *UpdateSpidLocationRequest) (*UpdateSpidLocationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSpidLocation not implemented")
}
func (*UnimplementedSpidHandlerServer) DeleteSpid(ctx context.Context, req *DeleteSpidRequest) (*DeleteSpidResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSpid not implemented")
}

func RegisterSpidHandlerServer(s *grpc.Server, srv SpidHandlerServer) {
	s.RegisterService(&_SpidHandler_serviceDesc, srv)
}

func _SpidHandler_GetSpidInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSpidRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpidHandlerServer).GetSpidInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/spidProtoBuffers.SpidHandler/GetSpidInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpidHandlerServer).GetSpidInfo(ctx, req.(*GetSpidRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SpidHandler_RegisterSpid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterSpidRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpidHandlerServer).RegisterSpid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/spidProtoBuffers.SpidHandler/RegisterSpid",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpidHandlerServer).RegisterSpid(ctx, req.(*RegisterSpidRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SpidHandler_UpdateBatteryInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateBatteryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpidHandlerServer).UpdateBatteryInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/spidProtoBuffers.SpidHandler/UpdateBatteryInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpidHandlerServer).UpdateBatteryInfo(ctx, req.(*UpdateBatteryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SpidHandler_UpdateSpidLocation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateSpidLocationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpidHandlerServer).UpdateSpidLocation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/spidProtoBuffers.SpidHandler/UpdateSpidLocation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpidHandlerServer).UpdateSpidLocation(ctx, req.(*UpdateSpidLocationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SpidHandler_DeleteSpid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteSpidRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpidHandlerServer).DeleteSpid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/spidProtoBuffers.SpidHandler/DeleteSpid",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpidHandlerServer).DeleteSpid(ctx, req.(*DeleteSpidRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _SpidHandler_serviceDesc = grpc.ServiceDesc{
	ServiceName: "spidProtoBuffers.SpidHandler",
	HandlerType: (*SpidHandlerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSpidInfo",
			Handler:    _SpidHandler_GetSpidInfo_Handler,
		},
		{
			MethodName: "RegisterSpid",
			Handler:    _SpidHandler_RegisterSpid_Handler,
		},
		{
			MethodName: "UpdateBatteryInfo",
			Handler:    _SpidHandler_UpdateBatteryInfo_Handler,
		},
		{
			MethodName: "UpdateSpidLocation",
			Handler:    _SpidHandler_UpdateSpidLocation_Handler,
		},
		{
			MethodName: "DeleteSpid",
			Handler:    _SpidHandler_DeleteSpid_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "spidHandler.proto",
}