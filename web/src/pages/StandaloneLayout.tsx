import { Outlet, useNavigation } from "react-router-dom";
import { Avatar, Layout, Nav, Skeleton } from "@douyinfe/semi-ui";
import DarkModeButton from "../components/DarkModeButton.tsx";
import AppFooter from "../components/AppFooter.tsx";

const Placeholder = () => {
  return (
    <Skeleton
      loading={true}
      placeholder={
        <div style={{ display: "flex", alignItems: "center" }}>
          <Skeleton.Avatar style={{ marginRight: 12 }} />
          <Skeleton.Title style={{ width: 120 }} />
        </div>
      }
    ></Skeleton>
  );
};

const gaCloudIconUrl = "/gacloud.svg";

function StandaloneLayout() {
  const navigation = useNavigation();

  const { Header, Footer, Content } = Layout;

  return (
    <Layout
      style={{ border: "1px solid var(--semi-color-border)", height: "100vh" }}
    >
      <Header style={{ backgroundColor: "var(--semi-color-bg-1)" }}>
        <Nav mode="horizontal">
          <Nav.Header>
            <Avatar src={gaCloudIconUrl} style={{ marginRight: "8px" }} />
            <span style={{ color: "var(--semi-color-text-2)" }}>GaCloud</span>
          </Nav.Header>
          <Nav.Footer>
            <DarkModeButton />
          </Nav.Footer>
        </Nav>
      </Header>
      <Content
        style={{
          padding: "24px",
          backgroundColor: "var(--semi-color-bg-0)",
          height: "100%",
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
        }}
      >
        <div
          style={{
            borderRadius: "10px",
            border: "1px solid var(--semi-color-border)",
            padding: "32px",
            width: "90%",
            height: "80%",
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
          }}
        >
          {navigation.state === "loading" ? <Placeholder /> : <Outlet />}
        </div>
      </Content>
      <Footer
        style={{
          display: "flex",
          justifyContent: "space-between",
          padding: "8px",
          color: "var(--semi-color-text-2)",
          backgroundColor: "rgba(var(--semi-grey-0), 1)",
        }}
      >
        <AppFooter />
      </Footer>
    </Layout>
  );
}

export default StandaloneLayout;
