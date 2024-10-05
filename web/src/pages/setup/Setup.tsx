import DatabaseSetup from "./DataBaseSetup.tsx";
import StorageSetup from "./StorageSetup.tsx";
import AdminSetup from "./AdminSetup.tsx";
import Finish from "./Finish.tsx";
import ElasticsearchSetup from "./ElasticsearchSetup.tsx";
import { Avatar, Skeleton, Space, Steps, Typography } from "@douyinfe/semi-ui";
import { IconBolt, IconSave, IconServer, IconUser } from "@douyinfe/semi-icons";
import { Suspense, useState } from "react";
import { getSetupStatus, StorageProviderConfig } from "../../api/setup.ts";
import { Await, useAsyncValue, useLoaderData } from "react-router-dom";

const gaCloudIconUrl = "/gacloud.svg";

function SetupStep() {
  const { status, storageProviders } = useAsyncValue() as {
    status: number;
    storageProviders: StorageProviderConfig;
  };

  const [currentStep, setCurrentStep] = useState(status);

  const updateStep = async () => {
    const status = await getSetupStatus();

    setCurrentStep(status);
  };

  const setupSteps = [
    {
      title: "数据库设置",
      description: "设置数据库连接信息",
      component: <DatabaseSetup onFinish={updateStep} />,
      icon: <IconServer />,
    },
    {
      title: "Elasticsearch设置",
      description: "设置Elasticsearch连接信息",
      component: <ElasticsearchSetup onFinish={updateStep} />,
      icon: <IconServer />,
    },
    {
      title: "初始存储设置",
      description: "设置初始存储信息",
      component: (
        <StorageSetup config={storageProviders} onFinish={updateStep} />
      ),
      icon: <IconSave />,
    },
    {
      title: "管理员设置",
      description: "设置管理员账号信息",
      component: <AdminSetup onFinish={updateStep} />,
      icon: <IconUser />,
    },
    {
      title: "完成",
      description: "完成设置",
      component: <Finish />,
      icon: <IconBolt />,
    },
  ];

  return (
    <>
      <Steps
        type="basic"
        current={currentStep}
        onChange={(step) => {
          if (step <= status) setCurrentStep(step);
        }}
        style={{
          display: "flex",
          width: "70vw",
          marginTop: "30px",
          marginBottom: "24px",
        }}
      >
        {setupSteps.map((step, index) => (
          <Steps.Step
            key={index}
            title={step.title}
            status={
              status === index ? "process" : status > index ? "finish" : "wait"
            }
            description={step.description}
            icon={step.icon}
          />
        ))}
      </Steps>
      <div>{setupSteps[currentStep]?.component}</div>
    </>
  );
}

function SetupSkeleton() {
  const placeholder = (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        width: "300px",
        marginBottom: "10px",
      }}
    >
      <Skeleton.Paragraph rows={3} />
      <Skeleton.Button />
    </div>
  );

  return (
    <Skeleton
      placeholder={placeholder}
      loading={true}
      style={{ textAlign: "center" }}
    ></Skeleton>
  );
}

function Setup() {
  const loader = useLoaderData();
  const { Title } = Typography;

  return (
    <>
      <Space vertical style={{ marginTop: "24px", marginBottom: "24px" }}>
        <Avatar
          src={gaCloudIconUrl}
          style={{ marginRight: "8px" }}
          size="extra-large"
        />
        <Title>GaCloud Setup</Title>
        <Suspense fallback={<SetupSkeleton />}>
          <Await resolve={loader}>
            <SetupStep />
          </Await>
        </Suspense>
      </Space>
    </>
  );
}

export default Setup;
