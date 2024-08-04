namespace go simple_short_url

struct BaseResp {
    1:  string StatusMessage (api.body="message")
    2:  i32 StatusCode (api.body="status_code")
}

struct ShortURL {
    1: string Short
    2: string Long
    3: string Desp
    4: string Token
    5: string Code
    6: i64 CreateTime
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
    2: required string RedirectURL (api.body="long")
    3: optional string Desp (api.body="desp")
    255: required string Token (api.body="token")
}

struct AddShortURLResponse {
    1: string Short (api.body="short")
    2: string Code (api.body="code")
    3: string RedirectURL (api.body="long")
    255: optional BaseResp BaseResp
}

struct DeleteShortURLRequest {
    1: required string Short (api.path="url")
    2: required string Token (api.body="token")
}

struct DeleteShortURLResponse {
    1: string Short (api.body="short")
    2: string RedirectURL (api.body="long")
    255: optional BaseResp BaseResp
}

struct RefreshRequest {
    1: required string Token (api.body="token")
}

struct RefreshResponse {
    255: optional BaseResp BaseResp
}

struct ListShortURLRequest {
    1: required string Token (api.body="token")
    2: required i32 offset (api.query="offset")
    3: required i32 limit (api.query="limit")
}

struct ListShortURLResponse {
    1: list<ShortURL> ShortURLs (api.body="short_urls")
    2: bool hasMore (api.body="hasMore")
    255: optional BaseResp BaseResp
}

service ShortService {
    string Hello(1: string name);
    RedirectShortURLResponse RedirectShortURL(1: RedirectShortURLRequest request)(api.get="/s/:url");
    AddShortURLResponse AddShortURL(1: AddShortURLRequest request)(api.post="/s/short_url");
    DeleteShortURLResponse DeleteShortURL(1: DeleteShortURLRequest request)(api.delete="/s/:url");
    RefreshResponse Refresh(1: RefreshRequest request)(api.post="/s/refresh");
    ListShortURLResponse ListShortURL(1: ListShortURLRequest request)(api.get="/s/list");
}