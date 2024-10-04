
import {Card} from "@douyinfe/semi-ui";
import {File} from "../api/files";

function ListFile({file} : {file: File}) {
    return (
        <Card title={file.path} style={{width: "300px"}}>
        </Card>
    );
}

export default ListFile;