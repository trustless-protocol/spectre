# IBC v2

## Abstraction layers

IBC v1 linked two chains through a **connection**, and each application through a **port** bound to its own **channel**, the connection and channel each set up by a multi-step handshake, and packets (application-specific) move on their application respective channel, each **routed by port and channel**. 

```text
                    Chain A                                               Chain B
application layer:  sending/handling packet                               handling/sending packet              per packet (packet send/receive/ack flow)
                    │ bind to                                             │ bind to
channel layer:      packet flow/management; app binding ◄──packet flow──► packet flow/management; app binding  one instance per application (channel handshake)
                    │ built on                                            │ built on
connection layer:   pair the two clients                                  pair the two clients                 chain-pair link (connection handshake)
                    │ wraps                                               │ wraps
light client layer: light client of B                                     light client of A                    each tracks the counterparty's consensus
                    (packet verification)                                 (packet verification)                               
```

IBC v2 drops both **connection** and the **channel** semantic from the stack, everything reduce to just light client pairs. Chains register each other's client IDs once, and every packet is routed by its IBC module (application) name. The verification machinery is reused from v1 unchanged. v2 adds only a thin `02-client/v2` layer recording which client on the other chain is our counterparty. The IBC v2 code lives in the `/v2` subpackages, chiefly [`04-channel/v2`](https://github.com/cosmos/ibc-go/tree/v10.3.0/modules/core/04-channel/v2) for packets and [`02-client/v2`](https://github.com/cosmos/ibc-go/tree/v10.3.0/modules/core/02-client/v2) for pairing.
```text
                    Chain A                                                 Chain B
application layer:  handling source/dest packet                             handling source/dest packet     per packet (send/receive/ack flow)
                    │                                                       │ 
router:             routing packet to app                                   routing packet to app        
					|                                                       |
channel layer:      packet flow/management                ◄──packet flow──► packet flow management   packet flow/mgmt only, no per app instance
                    │ call to                                               │ call to 
light client layer: pair the two clients;                                   pair the two clients;        chain-pair link (no handshake)
				    client of B                                             client of A
                    (packet verification)                                   (packet verification)
```


|                      | Description                                             | Abstraction layer                                       |
| -------------------- | ------------------------------------------------------- | ------------------------------------------------------- |
| Packet management    | Keeping track of the packet                             | 04-channel/v2                                           |
| Packet flow          | The packet lifecycle: send / recv / ack / timeout       | 04-channel/v2                                           |
| Verification         | Consensus verification of packet                        | 02-client keeper → light client (08-wasm) -> eth client |
| Routing              | Route packet to its correct IBC module                  | router (IBC module, by `DestinationPort`)               |
| Token transfer logic | handle burn and minting IBC token, keep track of supply | transfer module                                         |
