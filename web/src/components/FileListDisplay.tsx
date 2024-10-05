import { File } from "../api/files.ts";
import { List, Radio, RadioGroup } from "@douyinfe/semi-ui";
import FileIcon from "./FileIcon.tsx";
import { IconFolder } from "@douyinfe/semi-icons";
import { NavigateFunction } from "react-router-dom";
import { useState } from "react";

function ListFolder(item: File, onOpenFolder: (item: File) => void) {
  return (
    <Radio value={item.path}>
      <List.Item
        header={<IconFolder />}
        main={<span>{item.path}</span>}
        onClick={() => onOpenFolder(item)}
      />
    </Radio>
  );
}

function ListFile(item: File, onOpenFile: (item: File) => void) {
  return (
    <Radio value={item.path}>
      <List.Item
        header={<FileIcon mimeType={item.mime} />}
        main={<span>{item.path}</span>}
        onClick={() => onOpenFile(item)}
      />
    </Radio>
  );
}

function FileListDisplay({
  files,
  navigate,
}: {
  files: File[];
  navigate: NavigateFunction;
}) {
  const [selected, setSelected] = useState<string | undefined>(undefined);

  const onOpenFolder = (item: File) => {
    if (selected !== item.path) {
      return;
    }

    navigate(item.path, { relative: "path" });
  };

  const onOpenFile = (item: File) => {
    if (selected !== item.path) {
      return;
    }

    console.log("Opening file", item);
  };

  return (
    <RadioGroup
      type="pureCard"
      direction="vertical"
      value={selected}
      onChange={(e) => setSelected(e.target.value)}
    >
      <List
        dataSource={files}
        renderItem={(item) =>
          item.path.endsWith("/")
            ? ListFolder(item, onOpenFolder)
            : ListFile(item, onOpenFile)
        }
      />
    </RadioGroup>
  );
}

export default FileListDisplay;
