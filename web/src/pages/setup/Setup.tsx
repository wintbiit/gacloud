import {Avatar, Space, Steps, Typography} from "@douyinfe/semi-ui";
import DatabaseSetup from "./DataBaseSetup.tsx";
import StorageSetup from "./StorageSetup.tsx";
import AdminSetup from "./AdminSetup.tsx";
import Finish from "./Finish.tsx";
import {IconBolt, IconSave, IconServer, IconUser} from "@douyinfe/semi-icons";
import {useState} from "react";
import {getSetupStatus} from "../../api/setup.ts";

const gaCloudIconUrl = "/gacloud.svg";

function Setup() {
    const { Title } = Typography;

    const [backendCurrentStep, setBackendCurrentStep] = useState(0);
    const [currentStep, setCurrentStep] = useState(0);

    const updateStep = async () => {
        const status = await getSetupStatus();

        setCurrentStep(status);
        setBackendCurrentStep(status);
    }

    const setupSteps = [
        {
            title: "数据库设置",
            description: "设置数据库连接信息",
            component: <DatabaseSetup onFinish={updateStep} />,
            icon: <IconServer />
        },
        {
            title: "初始存储设置",
            description: "设置初始存储信息",
            component: <StorageSetup />,
            icon: <IconSave />
        },
        {
            title: "管理员设置",
            description: "设置管理员账号信息",
            component: <AdminSetup />,
            icon: <IconUser />
        },
        {
            title: "完成",
            description: "完成设置",
            component: <Finish />,
            icon: <IconBolt />
        }
    ]

    return (
        <>
            <Space vertical style={{marginTop: "24px", marginBottom: "24px"}}>
                <Avatar src={gaCloudIconUrl} style={{marginRight: "8px"}} size="extra-large" />
                <Title>GaCloud Setup</Title>
                <Steps type="basic"
                       current={currentStep}
                       onChange={(step) => {
                           if (step <= backendCurrentStep)
                               setCurrentStep(step)
                       }}
                       style={{display: "flex", width: "70vw", marginTop: "30px", marginBottom: "24px"}}>
                    {setupSteps.map((step) => (
                        <Steps.Step key={step.title} title={step.title}
                                    description={step.description} icon={step.icon}/>
                    ))}
                </Steps>
                <div>
                    {setupSteps[currentStep].component}
                </div>
            </Space>
        </>
    )
}

export default Setup;