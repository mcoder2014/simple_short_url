namespace go simple_short_url

struct BaseResp {
    1:  string StatusMessage (api.body="message")
    2:  i32 StatusCode (api.body="status_code")
}

struct RedirectShortURLRequest {
    1: required string URL (api.path="url")
}

struct RedirectShortURLResponse {
    1: string URL (api.header="Location")
    255: optional BaseResp BaseResp
}

struct AddShortURLRequest {
    1: optional string Short (api.body="short")
    2: required string RedirectURL (api.body="redirect_url")
    3: optional string Desp (api.body="desp")
    255: required string Token (api.body="token")
}

struct AddShortURLResponse {
    1: string Short (api.body="short")
    2: string RedirectURL (api.body="redirect_url")
    255: optional BaseResp BaseResp
}

struct DeleteShortURLRequest {
    1: required string Short (api.path="url")
    2: required string Token (api.body="token")
}

struct DeleteShortURLResponse {
    1: string Short (api.body="short")
    2: string RedirectURL (api.body="redirect_url")
    255: optional BaseResp BaseResp
}

struct RefreshRequest {
    1: required string Token (api.body="token")
}

struct RefreshResponse {
    255: optional BaseResp BaseResp
}

service ShortService {
    string Hello(1: string name);
    RedirectShortURLResponse RedirectShortURL(1: RedirectShortURLRequest request)(api.get="/s/:url");
    AddShortURLResponse AddShortURL(1: AddShortURLRequest request)(api.post="/s/short_url");
    DeleteShortURLResponse DeleteShortURL(1: DeleteShortURLRequest request)(api.delete="/s/:url");
    RefreshResponse Refresh(1: RefreshRequest request)(api.post="/s/refresh");
}