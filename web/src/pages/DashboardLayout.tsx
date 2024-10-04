import {Link, Outlet, useNavigation} from "react-router-dom";
import {Layout, Nav, Button, Avatar,} from "@douyinfe/semi-ui";
import {
    IconBell, IconFile,
    IconHelpCircle,
    IconLikeHeart,
    IconShareStroked,
    IconUserGroup
} from "@douyinfe/semi-icons";
import DarkModeButton from "../components/DarkModeButton.tsx";
import AppFooter from "../components/AppFooter.tsx";
import Spinning from "../components/Spinning.tsx";

function DashboardLayout() {
    const navigation = useNavigation();

    const { Header, Footer, Sider, Content } = Layout;
    const selected = () => {
        const path = window.location.pathname;
        if (path.startsWith("/likes")) return "Likes";
        if (path.startsWith("/shares")) return "Share";
        if (path.startsWith("/groups")) return "Group";
        if (path.startsWith("/files")) return "Files";

        return "";
    }

    return (
      <Layout style={{
          border: '1px solid var(--semi-color-border)',
          height: "100vh",
          width: "100vw",
      }}>
          <Sider style={{ backgroundColor: 'var(--semi-color-bg-1)' }}>
              <Nav
                  selectedKeys={[selected()]}
                  style={{ maxWidth: 200, height: '100%' }}
                  items={[
                      { itemKey: 'Files', text: '文件', icon: <IconFile size="large" /> },
                      { itemKey: 'Likes', text: '收藏', icon: <IconLikeHeart size="large" />},
                      { itemKey: 'Share', text: '共享', icon: <IconShareStroked size="large" /> },
                      { itemKey: 'Group', text: '组文件夹', icon: <IconUserGroup size="large" /> },
                  ]}
                  header={{
                      logo: <Avatar src="/gacloud.svg" />,
                      text: 'GaCloud',
                  }}
                  footer={{
                      collapseButton: true,
                  }}
                  renderWrapper={({ itemElement, props }) => {
                      const routerMap = {
                            Files: '/',
                            Likes: '/likes',
                            Share: '/shares',
                            Group: '/groups',
                      } as Record<string, string>;
                      return (
                          <Link
                              style={{ textDecoration: "none" }}
                              to={routerMap[props.itemKey as string]}
                          >
                              {itemElement}
                          </Link>
                      );
                  }}
              />
          </Sider>
          <Layout>
              <Header style={{ backgroundColor: 'var(--semi-color-bg-1)' }}>
                  <Nav
                      mode="horizontal"
                      footer={
                          <>
                              <DarkModeButton />
                              <Button
                                  theme="borderless"
                                  icon={<IconBell size="large" />}
                                  style={{
                                      color: 'var(--semi-color-text-2)',
                                      marginRight: '12px',
                                  }}
                              />
                              <Button
                                  theme="borderless"
                                  icon={<IconHelpCircle size="large" />}
                                  style={{
                                      color: 'var(--semi-color-text-2)',
                                      marginRight: '12px',
                                  }}
                              />
                              <Avatar color="orange" size="small">
                                  YJ
                              </Avatar>
                          </>
                      }
                  ></Nav>
              </Header>
              <Content
                  style={{
                      padding: '24px',
                      backgroundColor: 'var(--semi-color-bg-0)',
                  }}
              >
                  {navigation.state == "loading" && <Spinning title="加载中" />}
                  <Outlet />
              </Content>
              <Footer
                  style={{
                      color: 'var(--semi-color-text-2)',
                      backgroundColor: 'rgba(var(--semi-grey-0), 1)',
                  }}
              >
                    <AppFooter />
              </Footer>
          </Layout>
      </Layout>
  );
}

export default DashboardLayout;
