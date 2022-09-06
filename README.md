# 📦 box
A lightweight wrapper making Go interfaces serialisable so they can be used as arguments for Temporal workflows and activities.

## Usage
Suppose I have the following interface for a message request:
```go
type MessageRequest interface {
    SendMessage() error
    ScheduleMessage(time.Time) error
}
```
With two concrete implementations:
```go
type SlackMessageRequest struct {
    RecipientID string
    Content []byte
    ...
}

type EmailRequest struct {
    To string
    Subject string
    CC string
    Body string
    ...
}
```
As part of some Temporal workflow, I wish to send either kind of message using an activity. Ideally, I would want the activity to be type-agnostic w.r.t. the message request, and simply interact through the interface:
```go
func (h *CommHandler) HandleMessageDelivery(ctx context.Context, msgRequest MessageRequest) error {
    ...
}
```
However, arguments to Temporal activities are serialised and de-serialised behind the scenes, and so that intuitive solution would not work: a concrete JSON embedding of one of the implementations of the interface cannot be directly de-serialised into an interface instance.

By slightly changing the activity code, `Box` solves this problem:
```go
func (h *CommHandler) HandleMessageDelivery(ctx context.Context, msgRequest Box[MessageRequest]) error {
    err := msgRequest.Unbox(&SlackMessageRequest{}, &EmailRequest{})
    ...
    err = msgRequest.Data.SendMessage()
    ...
}
```
With `Box`, de-serialisisation of the representation of the struct implementing the interface `MessageRequest` is deferred until I explicitly `Unbox` it, providing concrete types that satisfy the interface and might reside inside the `Box`.

After `Unbox`-ing, my activity can work with the interface as if it received it as an argument directly.

### Creating a new `Box`
```go
box, err := NewBox[MessageRequest](&SlackMessageRequest{
    ...
})
```

### How a concrete type is chosen for `Unbox`-ing
`Unbox` iteratively attempts to de-serialise the raw JSON representation of the object inside the box into each of the concrete types given, stopping at the first success. 