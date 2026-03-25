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
    2: optional common.UserInfo data
    3: optional MFAInfo auth
    4: optional common.TokenInfo token // 令牌（可选，当未开启MFA时）
}

struct MFAInfo {
    1: required bool mfa_required     // 是否需要MFA校验
    2: required string mfa_token      // MFA token
}

struct InfoQuery {
    1: required string user_id
}

struct InfoResp {
    1: required common.BaseResp base
    2: optional common.UserInfo data
}

struct AvatarReq {
    1: optional binary data  // required
}

struct AvatarResp {
    1: required common.BaseResp base
    2: optional common.UserInfo data
}

struct RefreshTokenReq {
    1: required string token  // old refresh token
}

struct RefreshTokenResp {
    1: required common.BaseResp base
    2: optional common.TokenInfo data
}

struct MFAQrcodeReq {}

struct MFAQrcodeResp {
    1: required common.BaseResp base
    2: optional MFAQrcodeInfo data
}

struct MFAQrcodeInfo {
    1: required string secret
    2: required string qrcode
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
    2: optional common.TokenInfo data
}

service UserService {
    RegisterResp  Register(1: required RegisterReq req)
    LoginResp     Login(1: required LoginReq req)
    InfoResp      GetInfo(1: required InfoQuery req)
    AvatarResp    UploadAvatar(1: required AvatarReq req)
    RefreshTokenResp RefreshToken(1: required RefreshTokenReq req)
    MFAQrcodeResp MFAQrcode(1: required MFAQrcodeReq req)
    MFABindResp   MFABind(1: required MFABindReq req)
    MFAVerifyResp MFAVerify(1: required MFAVerifyReq req)
}