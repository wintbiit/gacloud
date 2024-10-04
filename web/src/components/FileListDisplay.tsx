import {File} from "../api/files.ts";
import {List, Radio, RadioGroup} from "@douyinfe/semi-ui";
import FileIcon from "./FileIcon.tsx";
import {IconFolder} from "@douyinfe/semi-icons";

function ListFolder(item: File) {
    return (
        <Radio value={item.path} >
            <List.Item
                header={<IconFolder /> }
                main={ <span>{item.path}</span> }
            />
        </Radio>
    )
}

function ListFile(item: File) {
    return (
        <Radio value={item.path}>
            <List.Item
                header={<FileIcon mimeType={item.mime} /> }
                main={ <span>{item.path}</span> }
            />
        </Radio>
    )
}

function FileListDisplay({files}: {files: File[]}) {
    return (
        <RadioGroup type="pureCard" direction="vertical">
            <List dataSource={files} renderItem={item => item.path.endsWith("/") ? ListFolder(item) : ListFile(item)} />
        </RadioGroup>
    )
}

export default FileListDisplay;