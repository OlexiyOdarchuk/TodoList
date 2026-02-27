<script lang="ts">
    import { onMount } from 'svelte';
    import { api } from './api';
    import { setToken } from './auth.svelte';
    import { themeState, toggleTheme } from './theme.svelte';

    let username = $state('');
    let email = $state('');
    let password = $state('');
    let error = $state('');
    let success = $state('');
    let loading = $state(false);

    let { onSwitchToLogin, onSwitchToVerify } = $props<{ onSwitchToLogin: () => void, onSwitchToVerify: (email: string) => void }>();

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
                    document.getElementById('googleSignUpDiv'),
                    { theme: gTheme, size: 'large', width: '100%', shape: 'pill', text: 'signup_with' }
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
            const btnContainer = document.getElementById('googleSignUpDiv');
            if (btnContainer) {
                btnContainer.innerHTML = ''; 
                const gTheme = themeState.isDark ? 'filled_black' : 'outline';
                google.accounts.id.renderButton(
                    btnContainer,
                    { theme: gTheme, size: 'large', width: '100%', shape: 'pill', text: 'signup_with' }
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
                error = 'No token received from Google';
            }
        } catch (err: any) {
            error = err.message || 'Google signup failed';
        } finally {
            loading = false;
        }
    }

    async function handleRegister() {
        error = '';
        success = '';
        loading = true;
        try {
            const res = await api.register(username, email, password);
            success = res.message || 'Registration successful. Please verify your email.';
            setTimeout(() => {
                onSwitchToVerify(email);
            }, 1500);
        } catch (err: any) {
            error = err.message || 'Registration failed';
        } finally {
            loading = false;
        }
    }
</script>

<div class="container auth-container">
    <div class="row justify-content-center align-items-center auth-page">
        <div class="col-12 col-sm-9 col-md-7 col-lg-5 col-xl-4 auth-col">
            <div class="card auth-card shadow">
                <div class="auth-header position-relative">
                    <button class="btn btn-link text-white position-absolute top-0 end-0 m-3 p-0 fs-5 border-0" onclick={toggleTheme} title="Toggle Theme">
                        {#if themeState.isDark}
                            <i class="bi bi-sun-fill"></i>
                        {:else}
                            <i class="bi bi-moon-stars-fill"></i>
                        {/if}
                    </button>
                    <i class="bi bi-person-plus display-4 mb-2"></i>
                    <h2 class="fw-bold mb-0">Create Account</h2>
                    <p class="opacity-75 small">Join TaskFlow today</p>
                </div>
                <div class="card-body p-4 p-md-5">
                    {#if error}
                        <div class="alert alert-danger border-0 rounded-3 small py-2">
                            <i class="bi bi-exclamation-circle-fill me-2"></i>{error}
                        </div>
                    {/if}
                    {#if success}
                        <div class="alert alert-success border-0 rounded-3 small py-2">
                            <i class="bi bi-check-circle-fill me-2"></i>{success}
                        </div>
                    {/if}

                    <form onsubmit={(e) => { e.preventDefault(); handleRegister(); }}>
                        <div class="mb-3">
                            <label class="form-label small fw-semibold text-muted">Username</label>
                            <div class="input-group">
                                <span class="input-group-text bg-light border-end-0 text-muted"><i class="bi bi-at"></i></span>
                                <input type="text" class="form-control border-start-0 ps-0" bind:value={username} required placeholder="username" />
                            </div>
                        </div>
                        <div class="mb-3">
                            <label class="form-label small fw-semibold text-muted">Email Address</label>
                            <div class="input-group">
                                <span class="input-group-text bg-light border-end-0 text-muted"><i class="bi bi-envelope"></i></span>
                                <input type="email" class="form-control border-start-0 ps-0" bind:value={email} required placeholder="name@example.com" />
                            </div>
                        </div>
                        <div class="mb-4">
                            <label class="form-label small fw-semibold text-muted">Password</label>
                            <div class="input-group">
                                <span class="input-group-text bg-light border-end-0 text-muted"><i class="bi bi-key"></i></span>
                                <input type="password" class="form-control border-start-0 ps-0" bind:value={password} required placeholder="••••••••" />
                            </div>
                        </div>
                        <button type="submit" class="btn btn-primary w-100 mb-3 py-2 fs-5" disabled={loading}>
                            {loading ? 'Processing...' : 'Register'}
                        </button>
                    </form>

                    <div class="d-flex align-items-center my-4">
                        <hr class="flex-grow-1 text-muted opacity-25" />
                        <span class="mx-3 text-muted small fw-medium">OR CONTINUE WITH</span>
                        <hr class="flex-grow-1 text-muted opacity-25" />
                    </div>

                    <div id="googleSignUpDiv" class="d-flex justify-content-center mb-3 w-100"></div>

                    <div class="mt-4 text-center">
                        <p class="text-muted small mb-0">Already have an account?</p>
                        <button class="btn btn-link p-0 fw-bold text-decoration-none" onclick={onSwitchToLogin}>Back to Sign In</button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>
