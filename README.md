# coap

    import "github.com/cloudwalkio/omg-coap"

CoAP Client and Server in Go [![Build
Status](https://drone.io/github.com/cloudwalkio/omg-coap/status.png)](https://drone.io/github.com/cloudwalkio/omg-coap/latest)

### Index

* [Constants](#constants)
* [Variables](#variables)
* [func  Handle](#func--handle)
* [func  HandleFunc](#func--handlefunc)
* [func  ListenAndServe](#func--listenandserve)
* [func  Notify](#func--notify)
* [func  Serve](#func--serve)
* [func  StringToTime](#func--stringtotime)
* [func  TimeToString](#func--timetostring)
* [func  Transmit](#func--transmit)
* [func  TransmitToObserver](#func--transmittoobserver)
* [type COAPCode](#type-coapcode)
    + [func (COAPCode) String](#func-coapcode-string)
* [type COAPType](#type-coaptype)
    + [func (COAPType) String](#func-coaptype-string)
* [type Conn](#type-conn)
    + [func  Dial](#func--dial)
    + [func (*Conn) Receive](#func-conn-receive)
    + [func (*Conn) Send](#func-conn-send)
* [type Handler](#type-handler)
* [type HandlerFunc](#type-handlerfunc)
    + [func (HandlerFunc) ServeCOAP](#func-handlerfunc-servecoap)
* [type MediaType](#type-mediatype)
* [type Message](#type-message)
    + [func  Receive](#func--receive)
    + [func  ReceiveTimeout](#func--receivetimeout)
    + [func (*Message) AddOption](#func-message-addoption)
    + [func (Message) IsConfirmable](#func-message-isconfirmable)
    + [func (Message) IsObserver](#func-message-isobserver)
    + [func (Message) Option](#func-message-option)
    + [func (Message) Options](#func-message-options)
    + [func (Message) Path](#func-message-path)
    + [func (Message) PathString](#func-message-pathstring)
    + [func (*Message) RemoveOption](#func-message-removeoption)
    + [func (*Message) SetOption](#func-message-setoption)
    + [func (*Message) SetPath](#func-message-setpath)
    + [func (*Message) SetPathString](#func-message-setpathstring)
* [type OptionID](#type-optionid)
* [type RemoteAddr](#type-remoteaddr)
* [type ServeMux](#type-servemux)
    + [func  NewServeMux](#func--newservemux)
    + [func (*ServeMux) Handle](#func-servemux-handle)
    + [func (*ServeMux) HandleFunc](#func-servemux-handlefunc)
    + [func (*ServeMux) ServeCOAP](#func-servemux-servecoap)
* [type Server](#type-server)
    + [func (*Server) ListenAndServe](#func-server-listenandserve)
    + [func (*Server) Serve](#func-server-serve)


#### Constants
```go
const (
	// ResponseTimeout is the amount of time to wait for a
	// response.
	ResponseTimeout = time.Second * 1
	// ResponseRandomFactor is a multiplier for response backoff.
	ResponseRandomFactor = 1.5
	// MaxRetransmit is the maximum number of times a message will
	// be retransmitted.
	MaxRetransmit = 0
	// Print debug messages
	Verbose = false
)
```

```go
const (
	// Confirmable messages require acknowledgements.
	Confirmable = COAPType(0)
	// NonConfirmable messages do not require acknowledgements.
	NonConfirmable = COAPType(1)
	// Acknowledgement is a message type indicating a response to
	// a confirmable message.
	Acknowledgement = COAPType(2)
	// Reset indicates a permanent negative acknowledgement.
	Reset = COAPType(3)
)
```

```go
const (
	IfMatch       = OptionID(1)
	URIHost       = OptionID(3)
	ETag          = OptionID(4)
	IfNoneMatch   = OptionID(5)
	Observe       = OptionID(6)
	URIPort       = OptionID(7)
	LocationPath  = OptionID(8)
	URIPath       = OptionID(11)
	ContentFormat = OptionID(12)
	MaxAge        = OptionID(14)
	URIQuery      = OptionID(15)
	Accept        = OptionID(17)
	LocationQuery = OptionID(20)
	ProxyURI      = OptionID(35)
	ProxyScheme   = OptionID(39)
	Size1         = OptionID(60)
)
```

```go
const (
	TextPlain     = MediaType(0)  // text/plain;charset=utf-8
	AppLinkFormat = MediaType(40) // application/link-format
	AppXML        = MediaType(41) // application/xml
	AppOctets     = MediaType(42) // application/octet-stream
	AppExi        = MediaType(47) // application/exi
	AppJSON       = MediaType(50) // application/json
)
```
Content types.

#### Variables

```go
var (
	ErrInvalidTokenLen   = errors.New("invalid token length")
	ErrOptionTooLong     = errors.New("option is too long")
	ErrOptionGapTooLarge = errors.New("option gap too large")
)
```
Message encoding errors.

```go
var DefaultServeMux = NewServeMux()
```

```go
var DefaultServer = new(Server)
```

#### func  [Handle](#handle)

```go
func Handle(pattern string, handler Handler)
```

#### func  [HandleFunc](#handlefunc)

```go
func HandleFunc(pattern string, handler func(r *RemoteAddr, m *Message) *Message)
```

#### func  [ListenAndServe](#listenandserve)

```go
func ListenAndServe(addr string, handler Handler) error
```
Bind to the given address and serve requests forever.

#### func  [Notify](#notify)

```go
func Notify(resource string, m *Message)
```
Notify observers of resource.

#### func  [Serve](#serve)

```go
func Serve(l *net.UDPConn, handler Handler) error
```

#### func  [StringToTime](#stringtotime)

```go
func StringToTime(s string) (uint32, error)
```
StringToTime translates the RRSIG's incep. and expir. times from string values
like "20110403154150" to an 32 bit integer. It takes serial arithmetic (RFC
1982) into observe option.

#### func  [TimeToString](#timetostring)

```go
func TimeToString(t uint32) string
```
TimeToString translates the RRSIG's incep. and expir. times to the string
representation used when printing the record. It takes serial arithmetic (RFC
1982) into observe option.

#### func  [Transmit](#transmit)

```go
func Transmit(r *RemoteAddr, m Message) error
```
Unicast message to remote address

#### func  [TransmitToObserver](#transmittoobserver)

```go
func TransmitToObserver(resource, id string, m *Message) (done chan bool)
```
Transmit to an observer

#### type [COAPCode](#coapcode)

```go
type COAPCode uint8
```

COAPCode is the type used for both request and response codes.

```go
const (
	GET    COAPCode = 1
	POST   COAPCode = 2
	PUT    COAPCode = 3
	DELETE COAPCode = 4
	// deprecated
	SUBSCRIBE COAPCode = 5
)
```
Request Codes

```go
const (
	Created               COAPCode = 65
	Deleted               COAPCode = 66
	Valid                 COAPCode = 67
	Changed               COAPCode = 68
	Content               COAPCode = 69
	BadRequest            COAPCode = 128
	Unauthorized          COAPCode = 129
	BadOption             COAPCode = 130
	Forbidden             COAPCode = 131
	NotFound              COAPCode = 132
	MethodNotAllowed      COAPCode = 133
	NotAcceptable         COAPCode = 134
	PreconditionFailed    COAPCode = 140
	RequestEntityTooLarge COAPCode = 141
	UnsupportedMediaType  COAPCode = 143
	InternalServerError   COAPCode = 160
	NotImplemented        COAPCode = 161
	BadGateway            COAPCode = 162
	ServiceUnavailable    COAPCode = 163
	GatewayTimeout        COAPCode = 164
	ProxyingNotSupported  COAPCode = 165
)
```
Response Codes

#### func (COAPCode) [String](#string)

```go
func (c COAPCode) String() string
```

#### type [COAPType](#coaptype)

```go
type COAPType uint8
```

COAPType represents the message type.

#### func (COAPType) [String](#string)

```go
func (t COAPType) String() string
```

#### type [Conn](#conn)

```go
type Conn struct {
}
```

Conn is a CoAP client connection.

#### func  [Dial](#dial)

```go
func Dial(n, addr string) (*Conn, error)
```
Dial connects a CoAP client.

#### func (*Conn) [Receive](#receive)

```go
func (c *Conn) Receive() (*Message, error)
```
Receive a message.

#### func (*Conn) [Send](#send)

```go
func (c *Conn) Send(req Message) (*Message, error)
```
Send a message. Get a response if there is one.

#### type [Handler](#handler)

```go
type Handler interface {
	// Handle the message and optionally return a response message.
	ServeCOAP(r *RemoteAddr, m *Message) *Message
}
```

Handler is a type that handles CoAP messages.

#### type [HandlerFunc](#handlerfunc)

```go
type HandlerFunc func(r *RemoteAddr, m *Message) *Message
```


#### func (HandlerFunc) [ServeCOAP](#servecoap)

```go
func (f HandlerFunc) ServeCOAP(r *RemoteAddr, m *Message) *Message
```

#### type [MediaType](#mediatype)

```go
type MediaType byte
```

MediaType specifies the content type of a message.

#### type [Message](#message)

```go
type Message struct {
	Type      COAPType
	Code      COAPCode
	MessageID uint16

	Token, Payload []byte
}
```

Message is a CoAP message.

#### func  [Receive](#receive)

```go
func Receive(l *net.UDPConn, buf []byte) (Message, error)
```
Receive a message.

#### func  [ReceiveTimeout](#receivetimeout)

```go
func ReceiveTimeout(l *net.UDPConn, rt time.Duration, buf []byte) (Message, error)
```
Receive a message with timeout.

#### func (*Message) [AddOption](#addoption)

```go
func (m *Message) AddOption(opId OptionID, val interface{})
```
AddOption adds an option.

#### func (Message) [IsConfirmable](#isconfirmable)

```go
func (m Message) IsConfirmable() bool
```
IsConfirmable returns true if this message is confirmable.

#### func (Message) [IsObserver](#isobserver)

```go
func (m Message) IsObserver() bool
```
IsObserver returns true if this message is for observe.

#### func (Message) [Option](#option)

```go
func (m Message) Option(o OptionID) interface{}
```
Option gets the first value for the given option ID.

#### func (Message) [Options](#options)

```go
func (m Message) Options(o OptionID) []interface{}
```
Get all the values for the given option.

#### func (Message) [Path](#path)

```go
func (m Message) Path() []string
```
Path gets the Path set on this message if any.

#### func (Message) [PathString](#pathstring)

```go
func (m Message) PathString() string
```
PathString gets a path as a / separated string.

#### func (*Message) [RemoveOption](#removeoption)

```go
func (m *Message) RemoveOption(opId OptionID)
```
RemoveOption removes all references to an option

#### func (*Message) [SetOption](#setoption)

```go
func (m *Message) SetOption(opId OptionID, val interface{})
```
SetOption sets an option, discarding any previous value

#### func (*Message) [SetPath](#setpath)

```go
func (m *Message) SetPath(s []string)
```
SetPath updates or adds a LocationPath attribute on this message.

#### func (*Message) [SetPathString](#setpathstring)

```go
func (m *Message) SetPathString(s string)
```
SetPathString sets a path by a / separated string.

#### type [OptionID](#optionid)

```go
type OptionID uint8
```

Option IDs.

#### type [RemoteAddr](#remoteaddr)

```go
type RemoteAddr struct {
	*net.UDPAddr
}
```


#### type [ServeMux](#servemux)

```go
type ServeMux struct {
}
```

ServeMux provides mappings from a common endpoint to handlers by request path.

#### func  [NewServeMux](#newservemux)

```go
func NewServeMux() *ServeMux
```
NewServeMux creates a new ServeMux.

#### func (*ServeMux) [Handle](#handle)

```go
func (mux *ServeMux) Handle(pattern string, handler Handler)
```
Handle configures a handler for the given path.

#### func (*ServeMux) [HandleFunc](#handlefunc)

```go
func (mux *ServeMux) HandleFunc(pattern string, handler func(r *RemoteAddr, m *Message) *Message)
```
HandleFunc configures a handler for the given path.

#### func (*ServeMux) [ServeCOAP](#servecoap)

```go
func (mux *ServeMux) ServeCOAP(r *RemoteAddr, m *Message) *Message
```
ServeCOAP handles a single COAP message. The message arrives from the given
listener having originated from the given UDPAddr.

#### type [Server](#server)

```go
type Server struct {
	Addr           string        // UDP address to listen on, ":coap" if empty
	Handler        Handler       // handler to invoke, coap.ServeMux if nil
	ReadTimeout    time.Duration // maximum duration before timing out read of the request
	WriteTimeout   time.Duration // maximum duration before timing out write of the response
	MaxHeaderBytes int           // maximum size of request headers, DefaultMaxHeaderBytes if 0
}
```


#### func (*Server) [ListenAndServe](#listenandserve)

```go
func (srv *Server) ListenAndServe() error
```

#### func (*Server) [Serve](#serve)

```go
func (srv *Server) Serve(l *net.UDPConn) error
```
