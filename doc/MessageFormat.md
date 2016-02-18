## Message format ID's

All commands are serialized with ProtoBuf [msgs.proto](https://github.com/v0l/wig/master/msgs.proto)

### Command
| id | direction | type           |
|----|-----------|----------------|
| 1  | Server    | ConnectCommand |
| 2  | Both      | ServerMessage  |
| 3  | Client    | StatusMessage  |

### StatusMessage
| msgtype | direction | statuscode | comment                  |
|---------|-----------|------------|--------------------------|
| 0       | Both      |            | Error                    |
| 1       |           |            |                          |
| 2       |           |            |                          |
| 3       | Client    | 1          | Connected to server      |
|         |           | 2          | Disconnected from server |