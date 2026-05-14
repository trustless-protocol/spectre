## 1 Cosmos escrow stuck forever: does not submit MsgTimeout back to Cosmos

khi gửi apcket với timeout ngắn, relayer chưa kịp gửi đến eth, tuy nhiên hết timeout không hoàn lại tiền

### giải pháp 
thêm logic 

New flow for Cosmos → ETH packet timeout:

    - UpdateEthClient
    Update the ETH light client on Cosmos to the latest finalized slot.

    - GetEthereumClientState
    Fetch proofBlockNumber and proofSlot from the updated client state.

    - Build receipt path
    destClientID + [0x02] + seqBytes
    (0x02 is the receipt type byte, symmetric to ack type 0x03 used in WriteAck).

    - GetEthMembershipProof
    Call eth_getProof on Ethereum. Since the receipt was never written, the storage value is 0x0, which serves as the non-membership proof.

    - SendCosmosTx(MsgTimeout)
    Submit the timeout packet to Cosmos to release the escrow.

## 2a relayer chỉ subscribe events từ block hiện tại, không query historical events khi restart.
test 2a



Tóm tắt những gì đã xảy ra:

    -Relayer cũ relay seq 6 thành công → ETH nhận được (balance = 4000)
    -Relayer cũ bị kill (Ctrl+C)
    -Relayer mới start 1 phút sau → miss WriteAcknowledgement event ở block 1425
    -MsgAcknowledgement cho seq 6 không được submit → Cosmos escrow vẫn khóa funds
    -Token trên ETH đã mint nhưng Cosmos không biết → inconsistent state
Lý do fail: relayer chỉ subscribe events từ block hiện tại, không query historical events khi restart.


### giải pháp

 Khi start, relayer giờ sẽ:

lấy latestBlock trên ETH
quét lại WriteAcknowledgement logs trong window lookback trước đó
chỉ enqueue các event còn pending trên Cosmos bằng cách kiểm tra packet commitment vẫn còn tồn tại
sau đó mới mở WS watch từ latestBlock + 1 để không có gap giữa recovery và live subscription

Điểm chính cho case của bạn là relayer mới không còn chỉ “nghe từ hiện tại”, nên nếu WriteAcknowledgement đã xảy ra lúc relayer down thì vẫn được nhặt lại và submit MsgAcknowledgement về Cosmos. Logic mới cũng tránh replay đống ack cũ đã xử lý xong bằng cách skip event lịch sử khi escrow/commitment bên Cosmos đã được clear.

Mặc định relayer scan lại 256 block ETH gần nhất khi khởi động. Có thể chỉnh bằng env ETH_STARTUP_LOOKBACK_BLOCKS nếu muốn window lớn hơn cho các lần restart lâu hơn
## 2b

sau khi stop eth node và restart thì relayer không connect lại ngay

 Đúng, log bạn đưa ra cho thấy nguyên nhân là relayer đã consume batch rồi làm rơi nó khi ETH RPC down.

Cụ thể ở relayer/services/main.go (line 147):


nhận batch

gọi UpdateCosmosClient

nếu lỗi thì continue

batch cũ biến mất, không có retry


 ### giải pháp

Tôi đã vá chỗ này:


khi UpdateCosmosClient fail hoặc trả nil, relayer sẽ requeue lại batch thay vì làm mất nó

code ở relayer/services/main.go (line 149) và helper mới ở relayer/services/main.go (line 612)


Sau khi ETH lên lại, khoảng vài giây sau batch sẽ được flush lại và retry.

Một lưu ý còn lại: nếu bạn stop cả EL/CL, WS subscription ETH có thể vẫn chết và không tự resubscribe. Nghĩa là fix này xử lý phần “mất batch khi ETH RPC outage”, nhưng nếu sau đó bạn vẫn miss WriteAcknowledgement thì bước tiếp theo sẽ là vá auto-resubscribe cho relayer/subscriber/event.go (line 393).

Cho test ổn định hơn, tôi vẫn khuyên 2b chỉ stop/start el-1-geth-lighthouse.


## 2b
 fail vì:

không có ack (MsgAcknowledgement), cũng không có MsgTimeout trên Cosmos
nghĩa là trạng thái đang là: ETH đã mint, nhưng Cosmos chưa nhận ack

## 3
ok

## 4a
ok

## 4b
ok

## 4c

ok

## 5a 
pass

## 5b

pass

## 5c 
skip

trên setup hiện tại. Payload on-chain gần nhất có bucket=4, active_count=3, unique_active=3, nên không phải case “exact boundary”.

## 5d
skip

trên setup hiện tại. MAX_BUCKET=4, min_signers_for_2of3=3, nên không rơi vào case “requires more than largest bucket”.

## 6a

Nếu không sửa thì không test được case này theo cách thực tế.

Nguyên nhân là:

- `./relayer create-clients` ban đầu không support `--trusting-period`
- runtime relayer cũng chưa đọc `trusting_period` từ config test
- vì vậy client luôn dùng trusting period mặc định dài theo chain thật, dẫn tới muốn test expiry phải chờ rất lâu

### giải pháp

thêm support cho `trusting_period` ở flow tạo client và ở runtime config của relayer

Sau đó có thể tạo một client test riêng với trusting period ngắn, ví dụ `60s`, rồi chạy theo flow:

- start relayer với config test
- stop relayer
- chờ quá trusting period
- gửi một packet Cosmos -> ETH mới
- start relayer lại

Kết quả mong đợi và cũng là kết quả local:

- relayer nhận packet, cố `updateClient`
- ETH revert với:
  `invalid block: untrusted state is outside of trusting period`
- log cuối:
  `[UpdateCosmosClient] SendEthTx failed: tx ... reverted`
  `[StartLoop] Failed to update cosmos light client: tx ... reverted`

pass

## 6b

Nếu không sửa thì path misbehaviour hiện tại không freeze được client.

Nguyên nhân là:

- `Misbehaviour.sol` trả `trustedHeight1/2` theo `header.signedHeader.header.height`
- nhưng phía `Groth16ICS07Tendermint` lại verify `trustedConsensusState1/2` theo `header.trustedHeight`
- tức là input và output đang lệch semantics, làm cho path misbehaviour gần như không thể thỏa đúng điều kiện verify

Ngoài ra còn một lỗi nữa:

- check trusting period trong `Misbehaviour.sol` đang so timestamp nanosecond trực tiếp với `trustingPeriod` tính theo second
- vì vậy payload hợp lệ vẫn có thể revert do hết trusting period giả

### giải pháp

sửa logic misbehaviour như sau:

- `trustedHeight1/2` phải trả theo `header.trustedHeight`
- khi check trusting period phải đổi timestamp từ nanosecond sang second trước khi so với `trustingPeriod`

Sau khi sửa, path misbehaviour có thể freeze client thành công.

Kết quả local:

- submit misbehaviour thành công
- chạy lại kiểm tra thì client đã ở trạng thái frozen
- các call tiếp theo revert với `FrozenClientState`

Lưu ý:

- case này hiện pass theo implementation on-chain hiện tại
- test local đang dùng 2 header hợp lệ khác height nhưng khác `appHash`, chưa phải same-height equivocation đúng nghĩa trên một Gaia honest chain

pass

## 7

pass

Kết quả local:

- happy path `Cosmos -> ETH` mint ra `1000` wrapped `stake` trên ETH
- sau đó `ETH -> Cosmos` gửi ngược `500`
- balance wrapped token trên ETH còn `500`
- phía Cosmos nhận lại `500stake` thành công qua `MsgRecvPacket`

Lưu ý khi chạy local:

- `timeoutTimestamp` trong setup này phải dùng `unix seconds`
- nếu dùng account `test1` làm signer cho relayer thì số dư ví Cosmos sẽ bị nhiễu do relayer cũng trả gas từ cùng account đó
- vì vậy với case `ETH -> Cosmos`, bằng chứng sạch nhất là event `coin_received 500stake` và `fungible_token_packet success=true` trong tx recv trên Cosmos

## 8

pass

Kết quả local:

- set rate limit `1500` cho wrapped `stake` trên escrow
- packet đầu gửi `1000stake` thành công
- sau packet đầu, `dailyUsage=1000` và wrapped balance của faucet tăng lên `1500`
- packet tiếp theo gửi `600stake` không được credit trên ETH
- sau packet thứ hai, `dailyUsage` vẫn `1000` và wrapped balance vẫn `1500`
- phía Cosmos nhận `MsgAcknowledgement` error cho `seq=3`

Lưu ý khi chạy local:

- deployment E2E hiện chưa wire sẵn `RATE_LIMITER_ROLE` cho escrow, nên cần bật role này trên `AccessManager` rồi mới `setRateLimit(...)` được
- nếu relayer và các tx manual cùng dùng account `test1` trên Cosmos thì có thể gặp `account sequence mismatch`; restart relayer là đủ để recovery và gửi lại ack pending

## 9

pass

Kết quả local:

- gửi 5 packet liên tiếp với `seq=4..8`
- phía ETH nhận đủ cả 5 packet
- balance cuối của receiver ETH mới là `500`
- phía Cosmos có đủ 5 `MsgAcknowledgement` tương ứng `seq=4..8`

Lưu ý khi chạy local:

- nếu relayer signer đang dùng `test1` trên Cosmos thì không nên dùng lại `test1` để bắn packet nhanh; dùng account khác như `test2`, `test3` sẽ tránh lỗi `account sequence mismatch`
- khi reuse môi trường cũ, nên dùng một receiver ETH mới có balance `0` để kết quả cuối `500` dễ kiểm tra hơn

## 10

pass

Kết quả local:

- `create-clients` fail sớm ở bước preflight checksum
- lỗi trả về là:
  `wasm checksum 0x0000000000000000000000000000000000000000000000000000000000000000 has not been previously stored on Cosmos`
- không có deploy ICS07 mới trên ETH
- config tạm không bị rewrite `ics07_client`
- Cosmos không tạo thêm wasm client mới, tổng số client state vẫn là `1`

Lưu ý khi chạy local:

- flow đã được sửa để validate checksum trên Cosmos trước khi mutate state bên ETH/config

## 11

pass

Kết quả local:

- start relayer với `ETH_PRIVATE_KEY` rỗng fail ngay với:
  `ETH_PRIVATE_KEY environment variable is required in .env file`
- start relayer với `COSMOS_PRIVATE_KEY=notvalidhex` fail ngay với:
  `failed to decode COSMOS_PRIVATE_KEY`

Lưu ý khi chạy local:

- flow startup đã được sửa để validate key ngay sau `godotenv.Load()`
- relayer fail trước khi dial RPC Cosmos hoặc Ethereum

## 12

pass

Kết quả local:

- generate bộ artifact prover mới vào thư mục tạm và start relayer với
  `PROVER_BIN_DIR=/private/tmp/test12-bin`
- không redeploy `Groth16Verifier_N{N}.sol` và không update `WrapperVerifier`
- gửi 1 packet Cosmos -> ETH mới, tx hash Cosmos là
  `2BA3AFC0FF4C9FFDF37B5345360201F55C043B0708D4D3F0943F2B38946E9C74`,
  `seq=9`
- relayer build proof local thành công nhưng `updateClient` trên ETH revert, tx hash là
  `0x9151f794bc6bfaa5171b86b25ff1df0d5976fc7240b7bf402c495d4676cc1314`
- revert data trả về là `0xd611c318`
- receiver ETH `0x2222222222222222222222222222222222222222` giữ nguyên balance `0`

Lưu ý khi chạy local:

- case này xác nhận mismatch giữa proving key mới và verifying key cũ trên chain sẽ làm mọi lần
  `updateClient` fail
- fix vận hành là redeploy toàn bộ verifier tương ứng và re-register bucket trong
  `WrapperVerifier` trước khi dùng bộ artifact prover mới
