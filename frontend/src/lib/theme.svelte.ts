export const themeState = $state({
    isDark: localStorage.getItem('theme') === 'dark',
});

export function toggleTheme() {
    themeState.isDark = !themeState.isDark;
    const theme = themeState.isDark ? 'dark' : 'light';
    localStorage.setItem('theme', theme);
    document.documentElement.setAttribute('data-bs-theme', theme);
}

export function initTheme() {
    const theme = themeState.isDark ? 'dark' : 'light';
    document.documentElement.setAttribute('data-bs-theme', theme);
}
