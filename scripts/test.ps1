$body = @'
{
    "data": [
       {
            "firstName": "Andres",
            "lastName": "Costa"
        }
    ]
}
'@

curl -X POST http://localhost:8080/events/m1/ev1 -H 'Content-Type: application/json' -d $body