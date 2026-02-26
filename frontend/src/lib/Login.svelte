<script lang="ts">
    import { onMount } from 'svelte';
    import { api } from './api';
    import { setToken } from './auth.svelte';
    import { themeState, toggleTheme } from './theme.svelte';

    let username = $state('');
    let password = $state('');
    let error = $state('');
    let loading = $state(false);
    let needsVerification = $state(false);

    let { onSwitchToRegister, onSwitchToVerify } = $props<{ onSwitchToRegister: () => void, onSwitchToVerify: (email: string) => void }>();

    const GOOGLE_CLIENT_ID = import.meta.env.VITE_GOOGLE_CLIENT_ID;

    onMount(() => {
        const initializeGoogle = () => {
            const google = (window as any).google;
            if (typeof google !== 'undefined' && google.accounts) {
                google.accounts.id.initialize({
                    client_id: GOOGLE_CLIENT_ID,
                    callback: handleGoogleResponse
                });

                const gTheme = themeState.isDark ? 'filled_black' : 'outline';

                google.accounts.id.renderButton(
                    document.getElementById('googleSignInDiv'),
                    { theme: gTheme, size: 'large', width: '100%', shape: 'pill' }
                );
            } else {
               setTimeout(initializeGoogle, 100);
            }
        };
        
        initializeGoogle();
    });

    $effect(() => {
        const google = (window as any).google;
        if (typeof google !== 'undefined' && google.accounts) {
            const btnContainer = document.getElementById('googleSignInDiv');
            if (btnContainer) {
                btnContainer.innerHTML = ''; 
                const gTheme = themeState.isDark ? 'filled_black' : 'outline';
                google.accounts.id.renderButton(
                    btnContainer,
                    { theme: gTheme, size: 'large', width: '100%', shape: 'pill' }
                );
            }
        }
    });

    async function handleGoogleResponse(response: any) {
        error = '';
        loading = true;
        try {
            const res = await api.googleLogin(response.credential);
            if (res.token) {
                setToken(res.token);
            } else {
                error = 'No token received from Google Login';
            }
        } catch (err: any) {
            error = err.message || 'Google Login failed';
        } finally {
            loading = false;
        }
    }

    async function handleLogin() {
        error = '';
        needsVerification = false;
        loading = true;
        try {
            const res = await api.login(username, password);
            if (res.token) {
                setToken(res.token);
            } else {
                error = 'No token received';
            }
        } catch (err: any) {
            error = err.message || 'Login failed';
            if (error.includes('not verified')) {
                needsVerification = true;
            }
        } finally {
            loading = false;
        }
    }
</script>

<div class="container">
    <div class="row justify-content-center align-items-center min-vh-100">
        <div class="col-11 col-sm-8 col-md-6 col-lg-4">
            <div class="card auth-card shadow">
                <div class="auth-header position-relative">
                    <button class="btn btn-link text-white position-absolute top-0 end-0 m-3 p-0 fs-5 border-0" onclick={toggleTheme} title="Toggle Theme">
                        {#if themeState.isDark}
                            <i class="bi bi-sun-fill"></i>
                        {:else}
                            <i class="bi bi-moon-stars-fill"></i>
                        {/if}
                    </button>
                    <i class="bi bi-check2-circle display-4 mb-2"></i>
                    <h2 class="fw-bold mb-0">Welcome Back</h2>
                    <p class="opacity-75 small">Log in to manage your tasks</p>
                </div>
                <div class="card-body p-4 p-md-5">
                    {#if error}
                        <div class="alert alert-danger border-0 rounded-3 small py-2 d-flex align-items-center">
                            <i class="bi bi-exclamation-triangle-fill me-2"></i>
                            <div>
                                {error}
                                {#if needsVerification}
                                    <button class="btn btn-link btn-sm p-0 d-block text-danger fw-bold text-decoration-none mt-1" onclick={() => onSwitchToVerify(username)}>
                                        Verify Account Now →
                                    </button>
                                {/if}
                            </div>
                        </div>
                    {/if}

                    <form onsubmit={(e) => { e.preventDefault(); handleLogin(); }}>
                        <div class="mb-3">
                            <label class="form-label small fw-semibold text-muted" for="usernameInput">Username or Email</label>
                            <div class="input-group">
                                <span class="input-group-text bg-light border-end-0 text-muted"><i class="bi bi-person"></i></span>
                                <input id="usernameInput" type="text" class="form-control border-start-0 ps-0" bind:value={username} required placeholder="john_doe" />
                            </div>
                        </div>
                        <div class="mb-4">
                            <label class="form-label small fw-semibold text-muted" for="passwordInput">Password</label>
                            <div class="input-group">
                                <span class="input-group-text bg-light border-end-0 text-muted"><i class="bi bi-lock"></i></span>
                                <input id="passwordInput" type="password" class="form-control border-start-0 ps-0" bind:value={password} required placeholder="••••••••" />
                            </div>
                        </div>
                        <button type="submit" class="btn btn-primary w-100 mb-3 py-2 fs-5" disabled={loading}>
                            {loading ? 'Logging in...' : 'Sign In'}
                        </button>
                    </form>

                    <div class="d-flex align-items-center my-4">
                        <hr class="flex-grow-1 text-muted opacity-25" />
                        <span class="mx-3 text-muted small fw-medium">OR CONTINUE WITH</span>
                        <hr class="flex-grow-1 text-muted opacity-25" />
                    </div>

                    <div id="googleSignInDiv" class="d-flex justify-content-center mb-3"></div>

                    <div class="mt-3 text-center">
                        <p class="text-muted small mb-0">Don't have an account?</p>
                        <button class="btn btn-link p-0 fw-bold text-decoration-none" onclick={onSwitchToRegister}>Create an account</button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>
