# xesende

A client for the [Esendex REST API][Esendex].


## dream

``` go
func main() {
  client := xesende.New(user, pass)

  paging := xesende.Paging(0, 10)

  messages, err := client.Messages.Sent(paging)
  // ...
  paging = paging.NextPage()
  moremessages, _ := client.Messages.Sent(paging)

  err := client.Send(xesende.Messages{
    {
      To: "...",
      Body: "...",
    },
  })

  // etc.
}
```


[Esendex]: http://developers.esendex.com/APIs/REST-API
