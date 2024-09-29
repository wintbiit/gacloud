import {ComponentType, lazy} from "react";
import {StorageProviderConfig} from "../../api/setup.ts";

const fileProviders = {
    "local": lazy(() => import("./Local.tsx")),
} as {[key: string]: ComponentType<{config: StorageProviderConfig}>};

function FileProvider({name, config}: {name: string | undefined, config: StorageProviderConfig}) {
    if (name === undefined) {
        return null;
    }

    const Provider = fileProviders[name];

    return <Provider config={config} />;

}

export default FileProvider;