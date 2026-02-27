<script lang="ts">
    import { onMount } from 'svelte';
    import { fade, slide, fly } from 'svelte/transition';
    import { flip } from 'svelte/animate';
    import { api, type Todo } from './api';
    import { setToken } from './auth.svelte';
    import { themeState, toggleTheme } from './theme.svelte';

    let todos: Todo[] = $state([]);
    let loading = $state(true);
    let error = $state('');

    let newTitle = $state('');
    let newDescription = $state('');
    let newDate = $state('');
    let newTime = $state('12:00');
    let isSubmitting = $state(false);

    let editingId: string | null = $state(null);
    let editTitle = $state('');
    let editDescription = $state('');
    let editDate = $state('');
    let editTime = $state('12:00');

    let filter: 'all' | 'pending' | 'completed' = $state('all');

    let filteredTodos = $derived(
        todos.filter(todo => {
            if (filter === 'completed') return todo.completed;
            if (filter === 'pending') return !todo.completed;
            return true;
        }).sort((a, b) => new Date(a.deadline).getTime() - new Date(b.deadline).getTime())
    );

    async function fetchTodos() {
        loading = true;
        try {
            const data = await api.getTodos();
            todos = Array.isArray(data) ? data : [];
        } catch (err: any) {
            error = err.message || 'Failed to fetch todos';
            if (error === 'Unauthorized') {
                setToken(null);
            }
        } finally {
            loading = false;
        }
    }

    onMount(() => {
        fetchTodos();
    });

    async function handleCreateTodo() {
        if (!newTitle.trim() || !newDate.trim() || !newTime.trim()) return;
        isSubmitting = true;
        try {
            const formattedDeadline = new Date(`${newDate}T${newTime}`).toISOString();
            const newTodo = await api.createTodo(newTitle, newDescription, formattedDeadline);
            todos = [newTodo, ...todos];
            newTitle = '';
            newDescription = '';
            newDate = '';
            newTime = '12:00';
        } catch (err: any) {
            alert('Failed to create todo: ' + err.message);
        } finally {
            isSubmitting = false;
        }
    }

    async function handleDeleteTodo(id: string) {
        if (!confirm('Are you sure you want to delete this task?')) return;
        try {
            await api.deleteTodo(id);
            todos = todos.filter(t => t.id !== id);
        } catch (err: any) {
            alert('Failed to delete: ' + err.message);
        }
    }

    async function toggleComplete(todo: Todo) {
        try {
            const updatedTodo = { ...todo, completed: !todo.completed };
            await api.updateTodo(updatedTodo);
            todos = todos.map(t => t.id === todo.id ? updatedTodo : t);
        } catch (err: any) {
            alert('Failed to update: ' + err.message);
        }
    }

    function startEditing(todo: Todo) {
        editingId = todo.id;
        editTitle = todo.title;
        editDescription = todo.description;
        const d = new Date(todo.deadline);
        d.setMinutes(d.getMinutes() - d.getTimezoneOffset());
        const isoStr = d.toISOString();
        editDate = isoStr.slice(0, 10);
        editTime = isoStr.slice(11, 16);
    }

    function cancelEditing() {
        editingId = null;
    }

    async function saveEdit(todo: Todo) {
        if (!editTitle.trim() || !editDate.trim() || !editTime.trim()) return;
        try {
            const formattedDeadline = new Date(`${editDate}T${editTime}`).toISOString();
            const updatedTodo = {
                ...todo,
                title: editTitle,
                description: editDescription,
                deadline: formattedDeadline
            };
            await api.updateTodo(updatedTodo);
            todos = todos.map(t => t.id === todo.id ? updatedTodo : t);
            editingId = null;
        } catch (err: any) {
            alert('Failed to update todo: ' + err.message);
        }
    }

    function logout() {
        setToken(null);
    }

    function isOverdue(deadlineStr: string, completed: boolean): boolean {
        if (completed) return false;
        return new Date(deadlineStr).getTime() < new Date().getTime();
    }

    let currentTab: 'tasks' | 'profile' = $state('tasks');

    let userProfile: any = $state(null);
    let isFetchingProfile = $state(false);

    let profileUsername = $state('');
    let profileEmail = $state('');
    let pendingEmail = $state('');
    let isOauthUser = $state(false);
    let hasPassword = $state(true);

    let oldPassword = $state('');
    let newPassword = $state('');
    let confirmPassword = $state('');

    let emailVerificationCode = $state('');
    let isWaitingForEmailCode = $state(false);
    let deletePassword = $state('');
    let deleteVerificationCode = $state('');
    let isWaitingForDeleteCode = $state(false);

    let profileMessage = $state({ type: '', text: '' });
    
    async function loadProfile() {
        if (isFetchingProfile) return;
        isFetchingProfile = true;
        try {
            const data = await api.getUser();
            userProfile = data;
            profileUsername = data.username || '';
            profileEmail = data.email || '';
            pendingEmail = data.pending_email || '';
            isOauthUser = !!data.oauth_provider;
            hasPassword = data.has_password;
            if (pendingEmail) {
                isWaitingForEmailCode = true;
            }
        } catch (err: any) {
            console.error("Failed to load profile", err);
        } finally {
            isFetchingProfile = false;
        }
    }

    function switchTab(tab: 'tasks' | 'profile') {
        currentTab = tab;
        if (tab === 'profile') {
            loadProfile();
            profileMessage = { type: '', text: '' };
        }
    }

    async function handleUpdateUsername() {
        if (!profileUsername.trim()) return;
        try {
            await api.updateUsername(profileUsername);
            profileMessage = { type: 'success', text: 'Username updated successfully!' };
            loadProfile();
        } catch (err: any) {
            profileMessage = { type: 'danger', text: err.message || 'Failed to update username' };
        }
    }

    async function handleUpdatePassword() {
        if ((hasPassword && !oldPassword) || !newPassword || !confirmPassword) {
            profileMessage = { type: 'danger', text: 'Required fields are missing' };
            return;
        }
        if (newPassword !== confirmPassword) {
            profileMessage = { type: 'danger', text: 'New passwords do not match' };
            return;
        }
        try {
            await api.updatePassword(oldPassword, newPassword);
            profileMessage = { type: 'success', text: hasPassword ? 'Password updated successfully!' : 'Password set successfully!' };
            oldPassword = '';
            newPassword = '';
            confirmPassword = '';
            loadProfile();
        } catch (err: any) {
            profileMessage = { type: 'danger', text: err.message || 'Failed to update password' };
        }
    }

    async function handleRequestEmailUpdate() {
        if (!profileEmail.trim() || profileEmail === userProfile.email) return;
        try {
            await api.requestEmailUpdate(profileEmail);
            profileMessage = { type: 'success', text: 'Verification code sent to new email' };
            isWaitingForEmailCode = true;
            loadProfile();
        } catch (err: any) {
            profileMessage = { type: 'danger', text: err.message || 'Failed to request email update' };
        }
    }

    async function handleVerifyEmailUpdate() {
        if (!emailVerificationCode.trim()) return;
        try {
            await api.verifyEmailUpdate(emailVerificationCode);
            profileMessage = { type: 'success', text: 'Email updated successfully!' };
            isWaitingForEmailCode = false;
            emailVerificationCode = '';
            loadProfile();
        } catch (err: any) {
            profileMessage = { type: 'danger', text: err.message || 'Failed to verify email' };
        }
    }

    async function handleRequestDeleteUser() {
        if (hasPassword && !deletePassword.trim()) {
            profileMessage = { type: 'danger', text: 'Password is required to delete account' };
            return;
        }
        if (!confirm('This action will permanently delete your account and todos. Continue?')) {
            return;
        }
        try {
            await api.requestDeleteUser(deletePassword);
            isWaitingForDeleteCode = true;
            profileMessage = { type: 'warning', text: 'Verification code sent to your email. Enter it to complete deletion.' };
            deletePassword = '';
        } catch (err: any) {
            profileMessage = { type: 'danger', text: err.message || 'Failed to request account deletion' };
        }
    }

    async function handleConfirmDeleteUser() {
        if (!deleteVerificationCode.trim()) {
            profileMessage = { type: 'danger', text: 'Enter verification code to delete account' };
            return;
        }
        if (!confirm('Delete account permanently? This cannot be undone.')) {
            return;
        }
        try {
            await api.confirmDeleteUser(deleteVerificationCode);
            alert('Account deleted successfully');
            setToken(null);
        } catch (err: any) {
            profileMessage = { type: 'danger', text: err.message || 'Failed to delete account' };
        }
    }

</script>

<nav class="navbar navbar-expand navbar-dark mb-4 mb-md-5 shadow-sm dashboard-navbar" style="background: linear-gradient(135deg, #6366f1 0%, #a855f7 100%);">
    <div class="container">
        <a class="navbar-brand d-flex align-items-center gap-2" href="/">
            <i class="bi bi-check2-all fs-3"></i>
            <span>TaskFlow</span>
        </a>
        <div class="ms-auto d-flex align-items-center gap-2 dashboard-nav-actions">
            <button class="btn btn-link text-white p-0 fs-5 border-0" onclick={toggleTheme} title="Toggle Theme">
                {#if themeState.isDark}
                    <i class="bi bi-sun-fill"></i>
                {:else}
                    <i class="bi bi-moon-stars-fill"></i>
                {/if}
            </button>
            <button class="btn btn-light btn-sm text-primary fw-bold rounded-pill px-3 py-1 shadow-sm {currentTab === 'tasks' ? 'opacity-50' : ''}" onclick={() => switchTab('tasks')} title="Tasks">
                <i class="bi bi-list-task me-1"></i> Tasks
            </button>
            <button class="btn btn-light btn-sm text-primary fw-bold rounded-pill px-3 py-1 shadow-sm {currentTab === 'profile' ? 'opacity-50' : ''}" onclick={() => switchTab('profile')} title="Profile">
                <i class="bi bi-person-fill me-1"></i> Profile
            </button>
            <button class="btn btn-outline-light btn-sm fw-bold rounded-pill px-3 py-1 border-0 text-white" style="background: rgba(255, 255, 255, 0.2);" onclick={logout}>
                <i class="bi bi-box-arrow-right me-1"></i> Logout
            </button>
        </div>
    </div>
</nav>

<div class="container pb-4 pb-md-5 dashboard-shell" style="max-width: 900px;">
    {#if error && error !== 'Unauthorized'}
        <div class="alert alert-danger shadow-sm border-0 rounded-4 mb-4" transition:fade>{error}</div>
    {/if}

    {#if currentTab === 'tasks'}
        <div class="card shadow-sm border-0 rounded-4 mb-5 overflow-hidden" transition:fly={{ y: -20, duration: 400 }}>
            <div class="card-body p-4 p-md-5">
                <h5 class="fw-bold text-dark mb-4 d-flex align-items-center">
                    <span class="bg-primary bg-opacity-10 p-2 rounded-3 me-3">
                        <i class="bi bi-pencil-square text-primary"></i>
                    </span>
                    Add New Task
                </h5>
                <form class="row g-3" onsubmit={(e) => { e.preventDefault(); handleCreateTodo(); }}>
                    <div class="col-md-12 mb-2">
                        <label class="form-label small fw-bold text-muted">Title</label>
                        <input type="text" class="form-control form-control-lg bg-light" placeholder="What needs to be done?" bind:value={newTitle} required />
                    </div>
                    <div class="col-md-7">
                        <label class="form-label small fw-bold text-muted">Description (optional)</label>
                        <input type="text" class="form-control form-control-lg bg-light" placeholder="Details..." bind:value={newDescription} />
                    </div>
                    <div class="col-md-3">
                        <label class="form-label small fw-bold text-muted">Date</label>
                        <input type="date" class="form-control form-control-lg bg-light" bind:value={newDate} required />
                    </div>
                    <div class="col-md-2">
                        <label class="form-label small fw-bold text-muted">Time</label>
                        <input type="time" class="form-control form-control-lg bg-light" bind:value={newTime} required />
                    </div>
                    <div class="col-12 text-end mt-4">
                        <button type="submit" class="btn btn-primary btn-lg rounded-pill px-5 fw-bold shadow" disabled={isSubmitting}>
                            {#if isSubmitting}
                                <span class="spinner-border spinner-border-sm me-2" role="status"></span> Adding...
                            {:else}
                                <i class="bi bi-plus-lg me-1"></i> Add Task
                            {/if}
                        </button>
                    </div>
                </form>
            </div>
        </div>

        <div class="d-flex flex-column flex-sm-row justify-content-between align-items-start align-items-sm-center gap-3 mb-4">
            <h4 class="fw-bold m-0 text-dark">Your Tasks</h4>
            <div class="nav nav-pills bg-white p-1 rounded-pill shadow-sm task-filter">
                <button type="button" class="nav-link py-1 px-3 rounded-pill {filter === 'all' ? 'active shadow-sm' : 'text-muted'}" onclick={() => filter = 'all'}>All</button>
                <button type="button" class="nav-link py-1 px-3 rounded-pill {filter === 'pending' ? 'active shadow-sm bg-warning text-dark' : 'text-muted'}" onclick={() => filter = 'pending'}>Pending</button>
                <button type="button" class="nav-link py-1 px-3 rounded-pill {filter === 'completed' ? 'active shadow-sm bg-success' : 'text-muted'}" onclick={() => filter = 'completed'}>Completed</button>
            </div>
        </div>

        {#if loading}
            <div class="text-center py-5 my-5" transition:fade>
                <div class="spinner-border text-primary" style="width: 3rem; height: 3rem;" role="status"></div>
                <p class="text-muted mt-3 fw-medium">Syncing with server...</p>
            </div>
        {:else if filteredTodos.length === 0}
            <div class="text-center py-5 bg-white rounded-4 shadow-sm border border-dashed" transition:fade>
                <div class="bg-light rounded-circle d-inline-flex p-4 mb-3">
                    <i class="bi bi-clipboard-x text-muted fs-1"></i>
                </div>
                <h5 class="text-dark fw-bold">No tasks here</h5>
                <p class="text-muted px-4">It seems you don't have any {filter !== 'all' ? filter : ''} tasks. Relax or create a new one!</p>
            </div>
        {:else}
            <div class="row row-cols-1 g-3">
                {#each filteredTodos as todo (todo.id)}
                    <div class="col" animate:flip={{ duration: 400 }} transition:slide|local>
                        <div class="card border-0 todo-card shadow-sm {todo.completed ? 'bg-light' : 'bg-white'}">
                            <div class="card-body p-4">
                                
                                {#if editingId === todo.id}
                                    <div transition:slide|local>
                                        <form onsubmit={(e) => { e.preventDefault(); saveEdit(todo); }}>
                                            <div class="row g-3 mb-3">
                                                <div class="col-md-12">
                                                    <input type="text" class="form-control fw-bold" bind:value={editTitle} required placeholder="Task title" />
                                                </div>
                                                <div class="col-md-12">
                                                    <textarea class="form-control" bind:value={editDescription} rows="2" placeholder="Task description"></textarea>
                                                </div>
                                                <div class="col-md-6">
                                                    <input type="date" class="form-control" bind:value={editDate} required />
                                                </div>
                                                <div class="col-md-6">
                                                    <input type="time" class="form-control" bind:value={editTime} required />
                                                </div>
                                            </div>
                                            <div class="d-flex justify-content-end gap-2 edit-actions">
                                                <button type="button" class="btn btn-light btn-sm text-muted rounded-pill px-3" onclick={cancelEditing}>Cancel</button>
                                                <button type="submit" class="btn btn-success btn-sm px-4 rounded-pill fw-bold"><i class="bi bi-check-lg me-1"></i> Save Changes</button>
                                            </div>
                                        </form>
                                    </div>
                                {:else}
                                    <div class="d-flex align-items-start gap-3">
                                        <button 
                                            class="btn p-0 text-{todo.completed ? 'success' : 'secondary'} mt-1 border-0 bg-transparent" 
                                            onclick={() => toggleComplete(todo)}
                                            style="font-size: 1.75rem;"
                                        >
                                            <i class="bi {todo.completed ? 'bi-check-circle-fill' : 'bi-circle'}"></i>
                                        </button>

                                        <div class="flex-grow-1" style="min-width: 0;">
                                            <div class="d-flex justify-content-between align-items-start">
                                                <h5 class="mb-1 fw-bold text-break pe-3 {todo.completed ? 'text-decoration-line-through text-muted' : 'text-dark'}">
                                                    {todo.title}
                                                </h5>
                                                <div class="d-flex gap-1">
                                                    <button class="btn btn-light btn-sm rounded-pill p-1 px-2 border-0" onclick={() => startEditing(todo)} title="Edit">
                                                        <i class="bi bi-pencil-fill text-primary small"></i>
                                                    </button>
                                                    <button class="btn btn-light btn-sm rounded-pill p-1 px-2 border-0" onclick={() => handleDeleteTodo(todo.id)} title="Delete">
                                                        <i class="bi bi-trash-fill text-danger small"></i>
                                                    </button>
                                                </div>
                                            </div>
                                            
                                            {#if todo.description}
                                                <p class="mb-3 text-muted text-break {todo.completed ? 'text-decoration-line-through' : ''}" style="white-space: pre-wrap;">{todo.description}</p>
                                            {/if}

                                            <div class="d-flex flex-wrap gap-2 align-items-center">
                                                <span class="badge {todo.completed ? 'bg-success bg-opacity-10 text-success' : 'bg-warning bg-opacity-10 text-warning'} rounded-pill px-3 py-1 fw-bold small">
                                                    {todo.completed ? 'COMPLETED' : 'PENDING'}
                                                </span>
                                                
                                                <span class="badge {isOverdue(todo.deadline, todo.completed) ? 'bg-danger bg-opacity-10 text-danger' : 'bg-light text-muted'} rounded-pill px-3 py-1 fw-medium small">
                                                    <i class="bi bi-clock-fill me-1"></i> {new Date(todo.deadline).toLocaleString([], { dateStyle: 'medium', timeStyle: 'short' })}
                                                    {#if isOverdue(todo.deadline, todo.completed)}
                                                        <span class="ms-1 fw-bold">(! OVERDUE)</span>
                                                    {/if}
                                                </span>
                                            </div>
                                        </div>
                                    </div>
                                {/if}

                            </div>
                        </div>
                    </div>
                {/each}
            </div>
        {/if}

    {:else if currentTab === 'profile'}
        <div class="row" transition:fade>
            <div class="col-12">
                <h4 class="fw-bold mb-4">Account Settings</h4>
                
                {#if profileMessage.text}
                    <div class="alert alert-{profileMessage.type} border-0 shadow-sm rounded-4 mb-4">
                        {profileMessage.text}
                    </div>
                {/if}

                {#if isFetchingProfile}
                    <div class="text-center py-5">
                        <div class="spinner-border text-primary" role="status"></div>
                    </div>
                {:else if userProfile}
                    <div class="card border-0 shadow-sm rounded-4 mb-4">
                        <div class="card-header bg-white border-bottom-0 pt-4 pb-0 px-4">
                            <h5 class="fw-bold m-0"><i class="bi bi-person me-2 text-primary"></i> Personal Info</h5>
                        </div>
                        <div class="card-body p-4">
                            <form onsubmit={(e) => { e.preventDefault(); handleUpdateUsername(); }}>
                                <div class="mb-3">
                                    <label class="form-label small fw-bold text-muted">Username</label>
                                    <div class="input-group mobile-input-group">
                                        <input type="text" class="form-control bg-light" bind:value={profileUsername} required />
                                        <button class="btn btn-outline-primary fw-bold px-4" type="submit" disabled={profileUsername === userProfile.username}>Update</button>
                                    </div>
                                </div>
                            </form>
                        </div>
                    </div>

                    <div class="card border-0 shadow-sm rounded-4 mb-4">
                        <div class="card-header bg-white border-bottom-0 pt-4 pb-0 px-4">
                            <h5 class="fw-bold m-0"><i class="bi bi-envelope me-2 text-primary"></i> Email Address</h5>
                        </div>
                        <div class="card-body p-4">
                            <form onsubmit={(e) => { e.preventDefault(); handleRequestEmailUpdate(); }}>
                                <div class="mb-3">
                                    <label class="form-label small fw-bold text-muted">Current Email</label>
                                    <div class="input-group mobile-input-group">
                                        <input type="email" class="form-control bg-light" bind:value={profileEmail} required />
                                        <button class="btn btn-outline-primary fw-bold px-4" type="submit" disabled={profileEmail === userProfile.email || isOauthUser}>
                                            Request Change
                                        </button>
                                    </div>
                                    {#if isOauthUser}
                                        <div class="form-text text-muted"><i class="bi bi-info-circle me-1"></i>Email change is disabled for OAuth users.</div>
                                    {/if}
                                    {#if userProfile.pending_email && !isWaitingForEmailCode}
                                        <div class="form-text text-warning mt-2"><i class="bi bi-exclamation-triangle me-1"></i>You have a pending email change to <b>{userProfile.pending_email}</b>. Please check your inbox.</div>
                                    {/if}
                                </div>
                            </form>

                            {#if isWaitingForEmailCode || userProfile.pending_email}
                                <div class="bg-primary bg-opacity-10 p-3 rounded-3 mt-3 border border-primary border-opacity-25" transition:slide>
                                    <label class="form-label small fw-bold text-primary">Verification Code (Sent to {userProfile.pending_email || profileEmail})</label>
                                    <form class="d-flex gap-2 mobile-inline-form" onsubmit={(e) => { e.preventDefault(); handleVerifyEmailUpdate(); }}>
                                        <input type="text" class="form-control bg-white" bind:value={emailVerificationCode} placeholder="Enter 6-digit code" required />
                                        <button class="btn btn-primary fw-bold px-4" type="submit">Verify</button>
                                    </form>
                                </div>
                            {/if}
                        </div>
                    </div>

                    <div class="card border-0 shadow-sm rounded-4 mb-4">
                        <div class="card-header bg-white border-bottom-0 pt-4 pb-0 px-4">
                            <h5 class="fw-bold m-0"><i class="bi bi-shield-lock me-2 text-primary"></i> {hasPassword ? 'Security' : 'Set Password'}</h5>
                        </div>
                        <div class="card-body p-4">
                            {#if !hasPassword}
                                <div class="alert alert-info py-2 border-0 shadow-sm mb-4 small">
                                    <i class="bi bi-info-circle me-2"></i>You logged in via {userProfile.oauth_provider}. Set a password to log in with your email later.
                                </div>
                            {/if}
                            <form onsubmit={(e) => { e.preventDefault(); handleUpdatePassword(); }}>
                                {#if hasPassword}
                                    <div class="mb-3">
                                        <label class="form-label small fw-bold text-muted">Current Password</label>
                                        <input type="password" class="form-control bg-light" bind:value={oldPassword} required />
                                    </div>
                                {/if}
                                <div class="row g-3 mb-4">
                                    <div class="col-md-6">
                                        <label class="form-label small fw-bold text-muted">New Password</label>
                                        <input type="password" class="form-control bg-light" bind:value={newPassword} required />
                                    </div>
                                    <div class="col-md-6">
                                        <label class="form-label small fw-bold text-muted">Confirm New Password</label>
                                        <input type="password" class="form-control bg-light" bind:value={confirmPassword} required />
                                    </div>
                                </div>
                                <div class="text-end">
                                    <button type="submit" class="btn btn-primary fw-bold px-4 rounded-pill shadow-sm" disabled={(hasPassword && !oldPassword) || !newPassword || !confirmPassword}>
                                        {hasPassword ? 'Update Password' : 'Set Password'}
                                    </button>
                                </div>
                            </form>
                        </div>
                    </div>

                    <div class="card border-danger border-opacity-25 shadow-sm rounded-4 mb-4">
                        <div class="card-header bg-danger bg-opacity-10 border-bottom-0 py-3 px-4">
                            <h5 class="fw-bold m-0 text-danger"><i class="bi bi-exclamation-triangle-fill me-2"></i>Danger Zone</h5>
                        </div>
                        <div class="card-body p-4">
                            <p class="text-muted mb-3">Deleting your account will permanently remove your profile and all tasks.</p>
                            <form onsubmit={(e) => { e.preventDefault(); handleRequestDeleteUser(); }}>
                                <div class="mb-3">
                                    <label class="form-label small fw-bold text-muted">
                                        {hasPassword ? 'Confirm with Password' : 'Confirm Deletion'}
                                    </label>
                                    <div class="input-group mobile-input-group">
                                        <input
                                            type="password"
                                            class="form-control bg-light"
                                            bind:value={deletePassword}
                                            placeholder={hasPassword ? 'Enter current password' : 'No password required for OAuth account'}
                                            required={hasPassword}
                                            disabled={!hasPassword}
                                        />
                                        <button class="btn btn-outline-danger fw-bold px-4" type="submit">Request Deletion</button>
                                    </div>
                                </div>
                            </form>

                            {#if isWaitingForDeleteCode}
                                <div class="bg-danger bg-opacity-10 p-3 rounded-3 border border-danger border-opacity-25" transition:slide>
                                    <label class="form-label small fw-bold text-danger">Deletion Verification Code</label>
                                    <form class="d-flex gap-2 mobile-inline-form" onsubmit={(e) => { e.preventDefault(); handleConfirmDeleteUser(); }}>
                                        <input type="text" class="form-control bg-white" bind:value={deleteVerificationCode} placeholder="Enter 6-digit code" required />
                                        <button class="btn btn-danger fw-bold px-4" type="submit">Delete Account</button>
                                    </form>
                                </div>
                            {/if}
                        </div>
                    </div>
                {/if}
            </div>
        </div>
    {/if}
</div>

<style>
    .active {
        background-color: #6366f1 !important;
        color: white !important;
    }
    .border-dashed {
        border: 2px dashed #e2e8f0 !important;
    }

    @media (max-width: 575.98px) {
        .dashboard-navbar .container {
            align-items: flex-start;
            gap: 0.75rem;
        }

        .dashboard-nav-actions {
            width: 100%;
            justify-content: flex-end;
            flex-wrap: wrap;
        }

        .dashboard-nav-actions .btn {
            padding-left: 0.75rem;
            padding-right: 0.75rem;
        }

        .task-filter {
            width: 100%;
            display: grid;
            grid-template-columns: 1fr 1fr 1fr;
            gap: 0.35rem;
            border-radius: 0.75rem !important;
        }

        .task-filter .nav-link {
            text-align: center;
            padding-left: 0.4rem !important;
            padding-right: 0.4rem !important;
            font-size: 0.8rem;
        }

        .mobile-input-group {
            flex-direction: column;
            gap: 0.5rem;
        }

        .mobile-input-group > .form-control {
            width: 100%;
            border-radius: 0.75rem !important;
            border-left: 1px solid #e2e8f0 !important;
        }

        .mobile-input-group > .btn {
            width: 100%;
            border-radius: 0.75rem !important;
        }

        .mobile-inline-form {
            flex-direction: column;
        }

        .mobile-inline-form .btn {
            width: 100%;
        }

        .edit-actions {
            justify-content: stretch !important;
            flex-direction: column;
        }

        .edit-actions .btn {
            width: 100%;
        }
    }
</style>
