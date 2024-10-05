import {
  Await,
  Link,
  useLoaderData,
  useNavigate,
  useParams,
} from "react-router-dom";
import {
  Breadcrumb,
  Button,
  ButtonGroup,
  Skeleton,
  Space,
} from "@douyinfe/semi-ui";
import { IconGridView1, IconHome, IconList } from "@douyinfe/semi-icons";
import { Suspense, useState } from "react";
import { File } from "../../api/files.ts";
import FilesDisplay from "../../components/FilesDisplay.tsx";

function FileListSkeleton() {
  const placeholder = (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        width: "300px",
        marginBottom: "10px",
      }}
    >
      <Skeleton.Paragraph rows={3} />
      <Skeleton.Button />
    </div>
  );

  return (
    <Skeleton
      placeholder={placeholder}
      loading={true}
      style={{ textAlign: "center" }}
    ></Skeleton>
  );
}

function Files() {
  const loader = useLoaderData();
  let path = useParams()["*"] as string;
  if (path.startsWith("/")) {
    path = path.slice(1);
  }
  if (path.endsWith("/")) {
    path = path.slice(0, -1);
  }
  const [displayMode, setDisplayMode] = useState<"list" | "grid">("list");
  const navigate = useNavigate();

  console.log("Files.tsx", path, loader, displayMode);

  return (
    <div>
      <Space
        style={{
          width: "100%",
          justifyContent: "space-between",
        }}
      >
        <Breadcrumb compact={false}>
          <Link to="/files/">
            <Breadcrumb.Item icon={<IconHome />} />
          </Link>
          {path !== "" &&
            path.split("/").map((part) => (
              <Link
                to={`/files/${path.slice(0, path.indexOf(part) + part.length)}`}
              >
                <Breadcrumb.Item>{part}</Breadcrumb.Item>
              </Link>
            ))}
        </Breadcrumb>
        <ButtonGroup>
          <Button
            icon={<IconList />}
            theme={displayMode == "list" ? "solid" : "light"}
            onClick={() => setDisplayMode("list")}
          />
          <Button
            icon={<IconGridView1 />}
            theme={displayMode == "grid" ? "solid" : "light"}
            onClick={() => setDisplayMode("grid")}
          />
        </ButtonGroup>
      </Space>

      <Suspense fallback={<FileListSkeleton />}>
        <Await resolve={loader}>
          {({ files }: { files: File[] }) => (
            <FilesDisplay
              files={files}
              displayMode={displayMode}
              navigate={navigate}
            />
          )}
        </Await>
      </Suspense>
    </div>
  );
}

export default Files;
