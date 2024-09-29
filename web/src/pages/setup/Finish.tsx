import {Button, Space, Typography} from "@douyinfe/semi-ui";
import {Link} from "react-router-dom";
import {finishSetup} from "../../api/setup.ts";

function Finish() {
    const { Title } = Typography;

    return (
        <Space vertical>
            <Title heading={3}>完成!</Title>
            <Title heading={4}>现在启动GaCloud!</Title>
            <Link to={"/login"}>
                <Button type="primary" size="large" onClick={finishSetup}>登录</Button>
            </Link>
        </Space>
    );
}

export default Finish;