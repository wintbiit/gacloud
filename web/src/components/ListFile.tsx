import {DashboardFile} from "../pages/loaders/dashboardLoader.ts";
import {Card} from "@douyinfe/semi-ui";

function ListFile({file} : {file: DashboardFile}) {
    return (
        <Card title={file.name} style={{width: "300px"}}>
            {file.isDir ? "文件夹" : "文件"}
        </Card>
    );
}

export default ListFile;