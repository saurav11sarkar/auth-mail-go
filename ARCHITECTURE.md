# Auth এবং user: কোড পড়ার guide

এই পরিবর্তন structure refactor; সম্পূর্ণ production security hardening নয়।

## Request কোথা দিয়ে যায়

```text
HTTP JSON → handler → request DTO (decode + validate)
                   → service input (explicit field mapping)
                   → service (existing business rules)
                   → repository (existing SQL + Scan)
                   → internal model / service result
                   → handler response DTO → JSON
```

Go struct cast করার বদলে field-by-field mapping ব্যবহার করা হয়েছে। এতে password/OTP ভুল করে response-এ যাওয়ার ঝুঁকি কমে।

## কোন ফাইলে কী

| File | দায়িত্ব |
| --- | --- |
| `internal/auth/dto/*.go` | Request-এ json/validate tags; response-এ JSON tags |
| `internal/auth/types.go` | RegisterInput, LoginInput, RefreshInput, ForgetPasswordInput, ResetPasswordInput, LoginResult |
| `internal/auth/model.go` | Account: repository/service-এর internal data |
| `internal/auth/handler.go` | Decode, validate, input mapping, service call, cookie, JSON |
| `internal/auth/service.go` | Registration, login, refresh, reset-এর বিদ্যমান business flow |
| `internal/auth/repository.go` | SQL এবং database error translation |
| `internal/auth/errors.go` | HTTP-independent পরিচিত errors |
| `internal/auth/http_mapping.go` | Response mapping এবং error → HTTP status |
| `internal/user/types.go` | Authenticated user ID-সহ ProfileInput |
| `internal/user/dto/profile_response.go` | ProfileResponse; পুরনো JSON casing বজায় আছে |
| `internal/user/http_mapping.go` | Profile model → response; not-found → 404 |

Profile GET-এ request body নেই। তাই খালি user_request.go সরানো হয়েছে; request struct প্রয়োজন নেই। User ID auth middleware-এর verified claims থেকে আসে।

## Error flow

Repository `pgx.ErrNoRows`-কে feature-এর not-found error-এ বদলায়। অন্য lookup error wrap হয়, তাই `errors.Is/As` দিয়ে underlying error চেনা যায়। Service business error ফেরায়। Handler-এর writeError পরিচিত error-কে আগের HTTP status/code/message-এ map করে। Unknown error-এর response generic 500 থাকে। Service-এ HTTP status বা DTO import নেই।

Login/refresh-এ missing account হলে 401, unexpected database failure হলে 500। Duplicate email constraint repository-তে ErrEmailExists হিসেবে translate হয়। Unexpected error-এর safe centralized logging এখনও যোগ করা হয়নি।

## Client compatibility এবং সীমা

- Route paths, request validation rules, success status/message এবং response envelope বদলানো হয়নি। Auth routes এখন POST-only; profile GET/HEAD এবং update PUT। ভুল method-এ Go ServeMux-এর default text 405 এবং Allow header ফেরে। Registration এখনও 200 ফেরায়।
- Auth timestamps `created_at/updated_at`; profile timestamps `createdAt/updatedAt`—আগের নামই আছে।
- Registration-এর পুরনো zero `otp_expires_at` field compatibility-এর জন্য রয়ে গেছে। Actual OTP/expiry response mapper কখনো কপি করে না। পরবর্তী API version-এ এই legacy fields বাদ দেওয়া যায়।
- Login/refresh token এখনও JSON-এ থাকে; refresh endpoint এখনও body-এর refresh_token পড়ে। Cookie-only auth-এ পরিবর্তন করা হয়নি।
- Cookie config inject করা হয়েছে; path `/api/v1/auth`, expiry seconds-এ এবং production-এ Secure। Production-এ HTTPS প্রয়োজন।
- Browser cross-site cookie আচরণ deployment origins ও SameSite policy-এর ওপর নির্ভরশীল। Angular project এখানে নেই; end-to-end integration যাচাই হয়নি।
- User-এর unfinished UpdateProfile method রাখা হয়েছে, তবে no-op success-এর বদলে explicit not-implemented error দেয়। কোনো update route যোগ করা হয়নি।

## পরবর্তী production কাজ

এই refactor দিয়ে production-ready দাবি করা যাবে না। আলাদা পরিবর্তনে পরীক্ষা করে এগুলো করতে হবে:

1. Tracked .env সরানো, secret rotation প্রয়োজন কি না যাচাই।
2. OTP-তে crypto/rand, নিরাপদ storage, attempt limit; refresh revocation/rotation।
3. Login inactive-account check।
4. Bcrypt byte-limit validation; refresh-token max=255 সীমা পর্যালোচনা।
5. Request body limit এবং strict JWT checks।
6. Graceful shutdown, database startup timeout, readiness, safe error logging।
7. PostgreSQL integration tests, CI, staging browser checks এবং deployment/backup verification।

## যাচাই

`go test ./...` দিয়ে packages compile এবং JSON contract/error mapping/cookie tests চালান। Tests live PostgreSQL, SMTP বা Angular চালায় না।

## Router ও middleware organization

- `internal/auth/router.go` এবং `internal/user/router.go`: feature-এর RegisterRoutes।
- `internal/routes/routes.go`: /api/v1 mount এবং global chain।
- Chain: RequestID → Logger → Recover → CORS → router (Go method pattern) → Auth (profile-এর জন্য) → handler।
- Middleware আলাদা ফাইলে: chain.go, cors.go, request_id.go, logger.go, recover.go, auth.go।
- Auth secret startup config থেকে আসে; প্রতি request-এ environment load হয় না। Valid token-এর role না মিললে 403।
- Request ID server তৈরি করে; CORS সেটি browser-এ expose করে। Logger query string log করে না।
- Recovery JSON 500 পাঠায়; ইতোমধ্যে response লেখা হয়ে গেলে HTTP status বদলানো যায় না। Streaming handler যোগ করলে recovery policy পুনর্বিবেচনা করতে হবে।
- AppError.Unwrap এবং repository-এর %w wrapping দিয়ে error chain বজায় থাকে।
- Existing handler error body (`status/code/message`) ও route-not-found envelope compatibility-এর জন্য এখনও পৃথক। Client migration ছাড়া global error schema বদলানো হয়নি।

Auth DTO naming: register_request.go, register_response.go, login_request.go, login_response.go, refresh_request.go, forget_password_request.go, reset_password_request.go।
Routing-এ GET/PUT/POST pattern সরাসরি ServeMux যাচাই করে। middlewares.Method সরানো হয়েছে। Custom method fallback নেই। API router-এর catch-all সরানো হয়েছে, তাই unknown API path-এ native text 404 ও unsupported method-এ native text 405 ফেরে।
