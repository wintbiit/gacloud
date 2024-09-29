import { setAdmin, SetAdminProps } from "../../api/setup.ts";
import { Button, Form, Notification, Space } from "@douyinfe/semi-ui";

function AdminSetup({ onFinish }: { onFinish: () => void }) {
  const submit = async (values: SetAdminProps) => {
    await setAdmin(values).catch((err) => {
      Notification.error({
        title: "管理员设置失败",
        content: err,
        duration: 5,
      });
    });

    Notification.success({
      title: "管理员设置成功",
      duration: 5,
    });

    onFinish();
  };

  return (
    <Form
      onSubmit={(values) => submit(values)}
      style={{ width: 400, textAlign: "left" }}
      title="GaCloud Elasticsearch Setup"
    >
      {({}) => (
        <>
          <Form.Input
            field="username"
            label="用户名"
            autoComplete="username"
            placeholder="输入用户名"
          />
          <Form.Input
            field="password"
            type="password"
            autoComplete="current-password"
            label="密码"
            placeholder="输入密码"
          />
          <Form.Input
            field="email"
            label="邮箱"
            autoComplete="email"
            placeholder="输入邮箱"
          />
          <Space>
            <Button type="tertiary" htmlType="submit">
              确认
            </Button>
            <Button type="warning" htmlType="reset">
              重置
            </Button>
          </Space>
        </>
      )}
    </Form>
  );
}

export default AdminSetup;
