import {
  Button,
  Checkbox,
  Divider,
  Form,
  Image,
  Space,
  SplitButtonGroup,
  Typography,
  Notification,
} from "@douyinfe/semi-ui";
import { login } from "../../api/auth.ts";
import { useNavigate } from "react-router-dom";

interface LoginValues {
  username: string;
  password: string;
}

function Login() {
  const { Title } = Typography;
  const navigate = useNavigate();

  const attemptLogin = async (values: LoginValues) => {
    await login(values.username, values.password)
      .then((resp) => {
        localStorage.setItem("authToken", resp.token);

        Notification.success({
          title: "登录成功",
          content: `欢迎回来，${values.username}`,
        });

        const redirect = new URLSearchParams(window.location.search).get(
          "redirect",
        );
        if (redirect) {
          navigate(redirect);
        } else {
          navigate("/");
        }
      })
      .catch(() => {
        Notification.error({
          title: "登录失败",
          content: "用户名或密码错误",
        });
      });
  };

  return (
    <div>
      <Space vertical>
        <Image src="/gacloud.svg" preview={false} />
        <Title>登录到 GaCloud </Title>
        <Form onSubmit={attemptLogin}>
          <Form.Input
            label={{ text: "用户名" }}
            field="username"
            placeholder="输入用户名或邮箱"
            autoComplete="username"
            trigger="blur"
            rules={[
              {
                required: true,
                message: "请输入用户名",
              },
            ]}
          />
          <Form.Input
            label={{ text: "密码" }}
            field="password"
            type="password"
            autoComplete="current-password"
            placeholder="输入密码"
            trigger="blur"
            rules={[
              {
                required: true,
                message: "请输入密码",
              },
            ]}
          />
          <Checkbox
            style={{
              alignSelf: "flex-start",
            }}
          >
            记住我
          </Checkbox>
          <SplitButtonGroup
            style={{
              display: "flex",
              justifyContent: "center",
              marginTop: "20px",
            }}
          >
            <Button theme="solid" htmlType="submit" size="large">
              登录
            </Button>
            <Button theme="outline" size="large">
              注册
            </Button>
          </SplitButtonGroup>
        </Form>
      </Space>
      <Divider />
    </div>
  );
}

export default Login;
