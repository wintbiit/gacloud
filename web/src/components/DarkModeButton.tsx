import { IconMoon, IconSun } from "@douyinfe/semi-icons";
import { Button } from "@douyinfe/semi-ui";
import { useState } from "react";

const switchMode = (darkMode: boolean) => {
    const body = document.body;
    if (darkMode) {
        body.setAttribute('theme-mode', 'dark');
    } else {
        body.removeAttribute('theme-mode');
    }
}

const systemPrefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;

const colorModeIcon = (darkMode: boolean) => {
    return darkMode ? <IconSun size="large" /> : <IconMoon size="large" />;
}

switchMode(true);

const DarkModeButton = () => {
    const [darkMode, setDarkMode] = useState(systemPrefersDark);

    const handleClick = () => {
        setDarkMode(!darkMode);
        switchMode(!darkMode);
    }

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
}

export default DarkModeButton;