import {Await, useLoaderData, useParams} from "react-router-dom";
import {Breadcrumb, Button, ButtonGroup, Skeleton, Space} from "@douyinfe/semi-ui";
import {IconGridView1, IconHome, IconList} from "@douyinfe/semi-icons";
import {Suspense, useState} from "react";
import {File} from "../../api/files.ts";
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

function pathBreadcrumbRoutes(path: string, prefix: string) {
    const routes = path.split("/").map((p) => {
        return {
            path: `/${p}`,
            href: `/${prefix}/` + path.slice(0, path.indexOf(p) + p.length),
            name: p,
        }
    })

    return [
        {
            href: `/${prefix}/`,
            icon: <IconHome />,
        },
        ...routes
    ]
}

function Files() {
    const loader = useLoaderData();
    const path = useParams()["*"] as string;
    const [displayMode, setDisplayMode] = useState<"list" | "grid">("list");

  return (
      <div>
          <Space style={{
              width: "100%",
              justifyContent: "space-between",
          }}>
              <Breadcrumb routes={pathBreadcrumbRoutes(path, "files")} />
              <ButtonGroup>
                  <Button icon={<IconList />} theme={displayMode == "list" ? "solid" : "light"} onClick={() => setDisplayMode("list")} />
                  <Button icon={<IconGridView1 />} theme={displayMode == "grid" ? "solid" : "light"} onClick={() => setDisplayMode("grid")} />
              </ButtonGroup>
          </Space>

          <Suspense fallback={<FileListSkeleton />}>
              <Await resolve={loader}>
                  {({files}: {files: File[]}) => (
                      <FilesDisplay files={files} displayMode={displayMode} />
                  )}
              </Await>
          </Suspense>
      </div>
  );
}

export default Files;
