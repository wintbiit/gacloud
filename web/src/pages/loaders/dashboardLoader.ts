import {defer} from "react-router-dom";
import {listFiles} from "../../api/files.ts";

export interface DashboardFile {
    name: string
    isDir: boolean
}

const dashboardLoader = async () => {
    const path = window.location.pathname.slice(1);

    console.log("Loading files from", path);

    const files = await listFiles(path);

    return defer({
        files: files.map((file) => ({
            name: file.path.endsWith("/") ? file.path.slice(0, -1) : file.path,
            isDir: file.path.endsWith("/"),
        } as DashboardFile))
    })
}

export default dashboardLoader;