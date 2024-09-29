import styles from "../../styles/login.module.scss";
import {Button, Checkbox, Form, Typography} from "@douyinfe/semi-ui";

function Login() {
    const { Title } = Typography;

    return (
        <div className={styles.login}>
            <div className={styles.component66}>
                <img
                    src="/gacloud.svg"
                    className={styles.logo}
                />
                <Title>登录到 GaCloud </Title>
            </div>
            <div className={styles.form}>
                <Form className={styles.inputs}>
                    <Form.Input
                        label={{text: "用户名"}}
                        field="input"
                        placeholder="输入用户名"
                        style={{width: "100%"}}
                        fieldStyle={{alignSelf: "stretch", padding: 0}}
                    />
                    <Form.Input
                        label={{text: "密码"}}
                        field="field1"
                        placeholder="输入密码"
                        style={{width: "100%"}}
                        fieldStyle={{alignSelf: "stretch", padding: 0}}
                    />
                </Form>
                <Checkbox type="default">记住我</Checkbox>
                <Button theme="solid" className={styles.button}>
                    登录
                </Button>
            </div>
        </div>
    )
}

export default Login;