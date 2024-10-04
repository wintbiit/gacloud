import {File} from "../api/files.ts";
import {CardGroup} from "@douyinfe/semi-ui";
import ListFile from "./ListFile.tsx";

function FileListDisplay({files}: {files: File[]}) {
    return (
        <CardGroup>
            {files.map((file, i) => (
                <ListFile key={i} file={file} />
            ))}
        </CardGroup>
    )
}

export default FileListDisplay;