import {useState} from "react";
import {
    setupElasticsearch,
    testElasticsearch,
    ElasticsearchOptionsProps
} from "../../api/setup.ts";
import {Button, Form, Notification, Space} from "@douyinfe/semi-ui";
import {IconTick} from "@douyinfe/semi-icons";

function ElasticsearchSetup({onFinish}: {onFinish: () => void}) {
    const [valid, setValid] = useState(false);
    const [loading, setLoading] = useState(false);

    const validate = async (formApi: any, values: ElasticsearchOptionsProps) => {
        setLoading(true);
        setValid(false);

        await testElasticsearch(values).then(res => {
            setValid(res.success);

            if (!res.success) {
                formApi.setError("host", "连接失败，请检查配置");
                Notification.error({
                    title: "Elasticsearch 连接失败",
                    content: res.reason,
                    duration: 5
                })
            }
        }).finally(() => {
            setLoading(false);
        })
    }

    const submit = async (values: ElasticsearchOptionsProps) => {
        await setupElasticsearch(values).catch(err => {
            Notification.error({
                title: "Elasticsearch 设置失败",
                content: err,
                duration: 5
            })
        })

        Notification.success({
            title: "Elasticsearch 设置成功",
            duration: 5
        })

        onFinish();
    }

    return (
        <Form onSubmit={values => submit(values)} style={{ width: 400, textAlign: "left" }} title="GaCloud Elasticsearch Setup">
            {({values, formApi}) => (
                <>
                    <Form.Input field="host" label="地址" placeholder="http://localhost:9200" />
                    <Form.Input field="user" label="用户名" autoComplete="username" placeholder="输入用户名" />
                    <Form.Input field="password" type="password" autoComplete="current-password" label="密码" placeholder="输入密码" />
                    <Space>
                        <Button type="primary" loading={loading} onClick={() => validate(formApi, values)}>
                            {valid && <IconTick />}
                            测试连接
                        </Button>
                        <Button type="tertiary" htmlType="submit" disabled={!valid}>
                            确认
                        </Button>
                        <Button type="warning" htmlType="reset">
                            重置
                        </Button>
                    </Space>
                </>
            )}
        </Form>
    )
}

export default ElasticsearchSetup