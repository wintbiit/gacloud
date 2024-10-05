import { Link, useRouteError } from "react-router-dom";
import { Button, Empty, SplitButtonGroup } from "@douyinfe/semi-ui";
import {
  IllustrationNoAccess,
  IllustrationNoAccessDark,
} from "@douyinfe/semi-illustrations";

function Error() {
  const error = useRouteError();

  const errorMessage = () => {
    if (error as Error) {
      return (error as Error).message;
    } else if (error as { data: string }) {
      return (error as { data: string }).data;
    }

    return "未知错误";
  };

  return (
    <Empty
      image={<IllustrationNoAccess />}
      darkModeImage={<IllustrationNoAccessDark />}
      title="布豪"
      description={errorMessage()}
      style={{
        marginTop: "20vh",
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
      }}
    >
      <SplitButtonGroup
        style={{
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
        }}
      >
        <Link to="/">
          <Button type="primary" size="large">
            返回首页
          </Button>
        </Link>
      </SplitButtonGroup>
    </Empty>
  );
}

export default Error;
