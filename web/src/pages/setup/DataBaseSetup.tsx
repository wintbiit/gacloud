import { Button, Form, Notification, Space } from "@douyinfe/semi-ui";
import { useState } from "react";
import { IconTick } from "@douyinfe/semi-icons";
import {
  MySQLOptionsProps,
  PostgreSQLOptionsProps,
  SQLiteOptionsProps,
  testDatabase,
  setupDatabase,
} from "../../api/setup.ts";

const dbOptions = [
  {
    value: "mysql",
    label: "MySQL",
  },
  {
    value: "postgresql",
    label: "PostgreSQL",
  },
  {
    value: "sqlite",
    label: "SQLite",
  },
];

function MySQLOptions() {
  return (
    <>
      <Form.Input field="host" label="主机" placeholder="输入主机地址" />
      <Form.Input field="port" label="端口" placeholder="输入端口" />
      <Form.Input
        field="username"
        label="用户名"
        autoComplete="username"
        placeholder="输入用户名"
      />
      <Form.Input
        field="password"
        type="password"
        label="密码"
        autoComplete="password"
        placeholder="输入密码"
      />
      <Form.Input
        field="database"
        label="数据库"
        placeholder="输入数据库名称"
      />
    </>
  );
}

function PostgreSQLOptions() {
  return (
    <>
      <Form.Input field="host" label="主机" placeholder="输入主机地址" />
      <Form.Input field="port" label="端口" placeholder="输入端口" />
      <Form.Input
        field="username"
        label="用户名"
        autoComplete="username"
        placeholder="输入用户名"
      />
      <Form.Input
        field="password"
        type="password"
        autoComplete="password"
        label="密码"
        placeholder="输入密码"
      />
      <Form.Input
        field="database"
        label="数据库"
        placeholder="输入数据库名称"
      />
    </>
  );
}

function SQLiteOptions() {
  return (
    <>
      <Form.Input field="path" label="路径" placeholder="输入数据库路径" />
    </>
  );
}

function DatabaseSetup({ onFinish }: { onFinish: () => void }) {
  const [valid, setValid] = useState(false);
  const [loading, setLoading] = useState(false);

  const validate = async (
    formApi: any,
    values: MySQLOptionsProps | PostgreSQLOptionsProps | SQLiteOptionsProps,
  ) => {
    setLoading(true);
    setValid(false);

    await testDatabase(values)
      .then((res) => {
        setValid(res.success);

        if (!res.success) {
          formApi.setError("type", "连接失败，请检查配置");
          Notification.error({
            title: `${values.type} 连接失败`,
            type: "warning",
            content: res.reason,
            duration: 5,
          });
        }
      })
      .finally(() => {
        setLoading(false);
      });
  };

  const submit = async (
    values: MySQLOptionsProps | PostgreSQLOptionsProps | SQLiteOptionsProps,
  ) => {
    await setupDatabase(values).catch((err) => {
      Notification.error({
        title: "数据库设置失败",
        content: err,
        duration: 5,
      });
    });

    Notification.success({
      title: "数据库设置成功",
      duration: 5,
    });

    onFinish();
  };

  return (
    <Form
      onSubmit={(values) => submit(values)}
      style={{ width: 400, textAlign: "left" }}
      title="GaCloud Database Setup"
    >
      {({ values, formApi }) => (
        <>
          <Form.Select
            field="type"
            label="数据库"
            optionList={dbOptions}
            placeholder="选择数据库类型"
          />
          {values.type === "mysql" && <MySQLOptions />}
          {values.type === "postgresql" && <PostgreSQLOptions />}
          {values.type === "sqlite" && <SQLiteOptions />}
          <Space>
            <Button
              type="primary"
              loading={loading}
              onClick={() => validate(formApi, values)}
            >
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
  );
}

export default DatabaseSetup;
