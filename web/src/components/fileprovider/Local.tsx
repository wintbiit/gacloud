import {Form} from "@douyinfe/semi-ui";
import {StorageProviderConfig} from "../../api/setup.ts";

function Local({config}: {config: StorageProviderConfig}) {
    const rootPath = config.rootPath;

    return (
        <div>
            <Form.Input label="Path" field="credential.mount_dir" placeholder={rootPath} />
        </div>
    )
}

export default Local;