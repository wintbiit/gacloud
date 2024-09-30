import { IconMoon, IconSun } from "@douyinfe/semi-icons";
import { Button } from "@douyinfe/semi-ui";
import { useState } from "react";
import switchDarkMode from "../pages/loaders/darkmode.ts";

const systemPrefersDark = window.matchMedia(
  "(prefers-color-scheme: dark)",
).matches;

const colorModeIcon = (darkMode: boolean) => {
  return darkMode ? <IconSun size="large" /> : <IconMoon size="large" />;
};

const DarkModeButton = () => {
  const [darkMode, setDarkMode] = useState(systemPrefersDark);

  const handleClick = () => {
    setDarkMode(!darkMode);
    switchDarkMode(!darkMode);
  };

  return (
    <Button
      theme="borderless"
      icon={colorModeIcon(darkMode)}
      onClick={handleClick}
      style={{
        color: "var(--semi-color-text-2)",
        marginRight: "12px",
      }}
    />
  );
};

export default DarkModeButton;
