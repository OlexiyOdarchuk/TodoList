export const authState = $state({
    token: localStorage.getItem('token') || null,
    isAuthenticated: !!localStorage.getItem('token'),
});

export function setToken(token: string | null) {
    if (token) {
        localStorage.setItem('token', token);
        authState.token = token;
        authState.isAuthenticated = true;
    } else {
        localStorage.removeItem('token');
        authState.token = null;
        authState.isAuthenticated = false;
    }
}
