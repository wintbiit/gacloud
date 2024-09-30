const switchDarkMode = (darkMode: boolean) => {
    const body = document.body;
    if (darkMode) {
        body.setAttribute("theme-mode", "dark");
    } else {
        body.removeAttribute("theme-mode");
    }
};

export default switchDarkMode;