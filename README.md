# Base URL
* The base URL is: http://localhost:8080

# Endpoint types
### Non-secure endpoints
All non-secure endpoints do not need authentication and use the method GET.
* [GET /api/articles](#get-articles)


# API documentation
Refer to the following for description of each endpoint

### GET /api/articles

#### Description:
Get endpoint status. When status is not `ok`, it is highly recommended to wait until the status changes back to `ok`.

#### Query:
-

#### Response:
```javascript
{
 "articles": {
        "items": [
            {
                "id": 52,
                "title": "Title#52",
                "excerpt": "Excerpt#52",
                "body": "Body#52",
                "image": "image#52",
                "categoryId": 0,
                "category": {
                    "ID": 0,
                    "Name": "categoryName#52"
                },
                "user": {
                    "Name": "userName#52",
                    "Avatar": user"Avatar#52"
                }
            }
       ],
        "paging": {
            "page": 1,
            "limit": 12,
            "prevPage": 0,
            "nextPage": 2,
            "count": 50,
            "totalPage": 5
        }
    }
}
```
