import { useState } from "react";
import {
  setupStorage,
  StorageOptionsProps,
  StorageProviderConfig,
  testStorage,
} from "../../api/setup.ts";
import { Button, Form, Notification, Space } from "@douyinfe/semi-ui";
import { IconTick } from "@douyinfe/semi-icons";
import FileProvider from "../../components/fileprovider/FileProvider.tsx";

function StorageSetup({
  config,
  onFinish,
}: {
  config: StorageProviderConfig;
  onFinish: () => void;
}) {
  const [valid, setValid] = useState(false);
  const [loading, setLoading] = useState(false);

  console.log(config);

  const validate = async (formApi: any, values: StorageOptionsProps) => {
    setLoading(true);
    setValid(false);

    await testStorage(values)
      .then((res) => {
        setValid(res.success);

        if (!res.success) {
          formApi.setError("type", "连接失败，请检查配置");
          Notification.error({
            title: "Elasticsearch 连接失败",
            content: res.reason,
            duration: 5,
          });
        }
      })
      .finally(() => {
        setLoading(false);
      });
  };

  const submit = async (values: StorageOptionsProps) => {
    await setupStorage(values).catch((err) => {
      Notification.error({
        title: "存储器设置失败",
        content: err,
        duration: 5,
      });
    });

    Notification.success({
      title: "存储器设置成功",
      duration: 5,
    });

    onFinish();
  };

  return (
    <Form
      onSubmit={(values) => submit(values)}
      style={{ width: 400, textAlign: "left" }}
      title="GaCloud Initialil Storage Setup"
    >
      {({ values, formApi }) => (
        <>
          <Form.Select
            field="type"
            label="类型"
            optionList={config.providers.map((p) => {
              return {
                value: p,
                label: p.toUpperCase(),
              };
            })}
            placeholder="选择存储器类型"
          />
          <FileProvider name={values.type} config={config} />

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

export default StorageSetup;
