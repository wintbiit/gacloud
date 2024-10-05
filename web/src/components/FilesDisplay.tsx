import { File } from "../api/files.ts";
import { ComponentType, lazy } from "react";
import {
  IllustrationNoContent,
  IllustrationNoContentDark,
} from "@douyinfe/semi-illustrations";
import { Button, Empty, SplitButtonGroup } from "@douyinfe/semi-ui";
import { Link, NavigateFunction } from "react-router-dom";

const FileListDisplay = lazy(
  () => import("./FileListDisplay.tsx"),
) as ComponentType<{ files: File[]; navigate: NavigateFunction }>;
const FileGridDisplay = lazy(
  () => import("./FileGridDisplay.tsx"),
) as ComponentType<{ files: File[] }>;

type FileDisplayMode = "grid" | "list";

function FilesDisplay({
  files,
  displayMode,
  navigate,
}: {
  files: File[];
  displayMode: FileDisplayMode;
  navigate: NavigateFunction;
}) {
  if (files.length <= 0) {
    return (
      <Empty
        image={<IllustrationNoContent />}
        darkModeImage={<IllustrationNoContentDark />}
        title="找不到文件"
        description="当前文件夹为空"
        style={{
          marginTop: "20vh",
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
        }}
      >
        <SplitButtonGroup>
          <Link to="/">
            <Button type="primary">返回根目录</Button>
          </Link>
          <Button type="secondary">上传文件</Button>
        </SplitButtonGroup>
      </Empty>
    );
  }

  if (displayMode === "grid") {
    return <FileGridDisplay files={files} />;
  }

  return <FileListDisplay files={files} navigate={navigate} />;
}

export default FilesDisplay;
