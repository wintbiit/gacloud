import {Space, Spin, Typography} from "@douyinfe/semi-ui";

function Loading() {
    const {Title} = Typography;

    return (
        <div style={{
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            height: "100vh",
            width: "100vw"
        }}>
            <Space vertical>
                <Spin size="large"/>
                <Title color="white">GaCloud is Loading...</Title>
            </Space>
        </div>
    )
}

export default Loading;