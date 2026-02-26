<script lang="ts">
    import { api } from './api';
    import { setToken } from './auth.svelte';
    import { themeState, toggleTheme } from './theme.svelte';

    let { email, onSwitchToLogin } = $props<{ email: string, onSwitchToLogin: () => void }>();

    let code = $state('');
    let error = $state('');
    let success = $state('');
    let loading = $state(false);

    async function handleVerify() {
        error = '';
        success = '';
        loading = true;
        try {
            const res = await api.verifyEmail(email, code);
            success = res.message || 'Verification successful!';
            if (res.token) {
                setTimeout(() => {
                    setToken(res.token);
                }, 1000);
            }
        } catch (err: any) {
            error = err.message || 'Verification failed';
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
                    <i class="bi bi-envelope-check display-4 mb-2"></i>
                    <h2 class="fw-bold mb-0">Verify Email</h2>
                    <p class="opacity-75 small">Check your inbox</p>
                </div>
                <div class="card-body p-4 p-md-5">
                    <p class="text-center text-muted mb-4 small">We sent a verification code to <strong>{email}</strong>. Please enter it below.</p>
                    
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
                    
                    <form onsubmit={(e) => { e.preventDefault(); handleVerify(); }}>
                        <div class="mb-4">
                            <input type="text" class="form-control form-control-lg text-center fs-4 letter-spacing-2" bind:value={code} required maxlength="6" placeholder="123456" />
                        </div>
                        <button type="submit" class="btn btn-primary w-100 mb-3 py-2 fs-5" disabled={loading || code.length < 6}>
                            {loading ? 'Verifying...' : 'Verify'}
                        </button>
                    </form>
                    <div class="mt-4 text-center">
                        <button class="btn btn-link p-0 fw-bold text-decoration-none" onclick={onSwitchToLogin}>Back to Sign In</button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>

<style>
    .letter-spacing-2 {
        letter-spacing: 0.5em;
    }
</style>
