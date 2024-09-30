import {Await, Link, useAsyncValue, useLoaderData} from "react-router-dom";
import {DashboardFile} from "../loaders/dashboardLoader.ts";
import {Breadcrumb, Button, CardGroup, Empty, Skeleton, SplitButtonGroup} from "@douyinfe/semi-ui";
import {IconHome} from "@douyinfe/semi-icons";
import {IllustrationNoContent, IllustrationNoContentDark} from "@douyinfe/semi-illustrations";
import {Suspense} from "react";
import ListFile from "../../components/ListFile.tsx";

function FileList() {
    const { files } = useAsyncValue() as { files: DashboardFile[] };

    if (files.length === 0) {
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
        )
    }

    return (
        <CardGroup>
            {files.map((file, i) => (
                <ListFile key={i} file={file} />
            ))}
        </CardGroup>
    )
}

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

function Dashboard() {
    const loader = useLoaderData();
    const path = window.location.pathname.slice(1);

  return (
      <div>
          <Breadcrumb>
              <Breadcrumb.Item icon={<IconHome />} href="/" />
              {path.split("/").map((p, i) => (
                  <Breadcrumb.Item key={i} href={path.slice(0, path.indexOf(p) + p.length)}>
                        {p}
                  </Breadcrumb.Item>
              ))}
          </Breadcrumb>

          <Suspense fallback={<FileListSkeleton />}>
              <Await resolve={loader}>
                  <FileList />
              </Await>
          </Suspense>
      </div>
  );
}

export default Dashboard;
