<script lang="ts">
    import { onMount } from 'svelte';
    import { authState } from './lib/auth.svelte';
    import { initTheme } from './lib/theme.svelte';
    import Login from './lib/Login.svelte';
    import Register from './lib/Register.svelte';
    import Verify from './lib/Verify.svelte';
    import Dashboard from './lib/Dashboard.svelte';

    let currentView: 'login' | 'register' | 'verify' = $state('login');
    let verifyEmail: string = $state('');

    onMount(() => {
        initTheme();
    });

    function switchToRegister() {
        currentView = 'register';
    }

    function switchToLogin() {
        currentView = 'login';
    }

    function switchToVerify(email: string) {
        verifyEmail = email;
        currentView = 'verify';
    }
</script>

<main>
    {#if authState.isAuthenticated}
        <Dashboard />
    {:else}
        {#if currentView === 'login'}
            <Login onSwitchToRegister={switchToRegister} onSwitchToVerify={switchToVerify} />
        {:else if currentView === 'register'}
            <Register onSwitchToLogin={switchToLogin} onSwitchToVerify={switchToVerify} />
        {:else if currentView === 'verify'}
            <Verify email={verifyEmail} onSwitchToLogin={switchToLogin} />
        {/if}
    {/if}
</main>

<style>
    :global(body) {
        margin: 0;
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
        background-color: #f8f9fa;
    }
</style>