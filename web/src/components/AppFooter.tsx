import { Link } from "react-router-dom";
import { IconGithubLogo } from "@douyinfe/semi-icons";
import { Divider } from "@douyinfe/semi-ui";

const githubUrl = "https://github.com";

function AppFooter() {
  return (
    <Divider margin="5px">
      <span style={{ marginRight: "16px" }}>GaCloud</span>
      <span>
        <Link to={githubUrl}>
          <IconGithubLogo />
        </Link>
      </span>
    </Divider>
  );
}

export default AppFooter;
