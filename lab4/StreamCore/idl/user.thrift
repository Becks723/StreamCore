namespace go user

include "common.thrift"

struct RegisterReq {
    1: required string username
    2: required string password
}

struct RegisterResp {
    1: required common.BaseResp base
}

struct LoginReq {
    1: required string username
    2: required string password
}

struct LoginResp {
    1: required common.BaseResp base
    2: required common.UserInfo data
    3: required common.AuthenticationInfo auth
}

struct InfoQuery {
    1: required string user_id
}

struct InfoResp {
    1: required common.BaseResp base
    2: required common.UserInfo data
}

struct AvatarReq {
    1: optional binary data  // required
}

struct AvatarResp {
    1: required common.BaseResp base
    2: required common.UserInfo data
}

struct MFAQrcodeReq {}

struct MFAQrcodeResp {
    1: required common.BaseResp base
    2: required common.MFAInfo data
}

struct MFABindReq {
    1: required string code
    2: required string secret
}

struct MFABindResp {
    1: required common.BaseResp base
}

struct MFAVerifyReq {
    1: required string mfa_token // MFA token
    2: required string code      // 六位校验码
}

struct MFAVerifyResp {
    1: required common.BaseResp base
}

service UserService {
    RegisterResp  Register(1: required RegisterReq req)
    LoginResp     Login(1: required LoginReq req)
    InfoResp      GetInfo(1: required InfoQuery req)
    AvatarResp    UploadAvatar(1: required AvatarReq req)
    MFAQrcodeResp MFAQrcode(1: required MFAQrcodeReq req)
    MFABindResp   MFABind(1: required MFABindReq req)
    MFAVerifyResp MFAVerify(1: required MFAVerifyReq req)
}