import { lazy, ReactNode } from "react";

const IconImage = lazy(() =>
  import("@douyinfe/semi-icons").then((module) => ({
    default: module.IconImage,
  })),
);
const IconVideo = lazy(() =>
  import("@douyinfe/semi-icons").then((module) => ({
    default: module.IconVideo,
  })),
);
const IconFile = lazy(() =>
  import("@douyinfe/semi-icons").then((module) => ({
    default: module.IconFile,
  })),
);
const IconMusic = lazy(() =>
  import("@douyinfe/semi-icons").then((module) => ({
    default: module.IconMusic,
  })),
);
const IconText = lazy(() =>
  import("@douyinfe/semi-icons").then((module) => ({
    default: module.IconText,
  })),
);

const iconMap = new Map<string, ReactNode>([
  ["image", <IconImage size="large" />],
  ["video", <IconVideo />],
  ["audio", <IconMusic />],
  ["text", <IconText />],
  ["default", <IconFile />],
]);

// Automatically find a suitable icon for a given mime type.
function FileIcon({ mimeType }: { mimeType: string }) {
  const type = mimeType.split("/")[0];
  return iconMap.get(type) || iconMap.get(mimeType) || iconMap.get("default");
}

export default FileIcon;
