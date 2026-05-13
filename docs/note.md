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
