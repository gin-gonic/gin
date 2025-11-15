//go:build gomock || generate

package quic

//go:generate sh -c "go tool mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/quic-go/quic-go -destination mock_send_conn_test.go github.com/quic-go/quic-go SendConn"
type SendConn = sendConn

//go:generate sh -c "go tool mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/quic-go/quic-go -destination mock_raw_conn_test.go github.com/quic-go/quic-go RawConn"
type RawConn = rawConn

//go:generate sh -c "go tool mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/quic-go/quic-go -destination mock_sender_test.go github.com/quic-go/quic-go Sender"
type Sender = sender

//go:generate sh -c "go tool mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/quic-go/quic-go -destination mock_stream_sender_test.go github.com/quic-go/quic-go StreamSender"
type StreamSender = streamSender

//go:generate sh -c "go tool mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/quic-go/quic-go -destination mock_stream_control_frame_getter_test.go github.com/quic-go/quic-go StreamControlFrameGetter"
type StreamControlFrameGetter = streamControlFrameGetter

//go:generate sh -c "go tool mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/quic-go/quic-go -destination mock_stream_frame_getter_test.go github.com/quic-go/quic-go StreamFrameGetter"
type StreamFrameGetter = streamFrameGetter

//go:generate sh -c "go tool mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/quic-go/quic-go -destination mock_frame_source_test.go github.com/quic-go/quic-go FrameSource"
type FrameSource = frameSource

//go:generate sh -c "go tool mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/quic-go/quic-go -destination mock_ack_frame_source_test.go github.com/quic-go/quic-go AckFrameSource"
type AckFrameSource = ackFrameSource

//go:generate sh -c "go tool mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/quic-go/quic-go -destination mock_sealing_manager_test.go github.com/quic-go/quic-go SealingManager"
type SealingManager = sealingManager

//go:generate sh -c "go tool mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/quic-go/quic-go -destination mock_unpacker_test.go github.com/quic-go/quic-go Unpacker"
type Unpacker = unpacker

//go:generate sh -c "go tool mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/quic-go/quic-go -destination mock_packer_test.go github.com/quic-go/quic-go Packer"
type Packer = packer

//go:generate sh -c "go tool mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/quic-go/quic-go -destination mock_mtu_discoverer_test.go github.com/quic-go/quic-go MTUDiscoverer"
type MTUDiscoverer = mtuDiscoverer

//go:generate sh -c "go tool mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/quic-go/quic-go -destination mock_conn_runner_test.go github.com/quic-go/quic-go ConnRunner"
type ConnRunner = connRunner

//go:generate sh -c "go tool mockgen -typed -build_flags=\"-tags=gomock\" -package quic -self_package github.com/quic-go/quic-go -destination mock_packet_handler_test.go github.com/quic-go/quic-go PacketHandler"
type PacketHandler = packetHandler

//go:generate sh -c "go tool mockgen -typed -package quic -self_package github.com/quic-go/quic-go -self_package github.com/quic-go/quic-go -destination mock_packetconn_test.go net PacketConn"
