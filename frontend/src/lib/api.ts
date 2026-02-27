import { authState, setToken } from './auth.svelte';

const API_BASE = '';

export interface Todo {
    id: string;
    user_id: string;
    title: string;
    description: string;
    completed: boolean;
    created_at: string;
    updated_at: string;
    deadline: string;
}

async function request(endpoint: string, options: RequestInit = {}) {
    const headers: Record<string, string> = {
        'Content-Type': 'application/json',
        ...(options.headers as Record<string, string> || {})
    };

    if (authState.token) {
        headers['Authorization'] = `Bearer ${authState.token}`;
    }

    const response = await fetch(`${API_BASE}${endpoint}`, {
        ...options,
        headers,
    });

    if (response.status === 401) {
        setToken(null);
        throw new Error('Unauthorized');
    }

    let data;
    const text = await response.text();
    if (text) {
        try {
            data = JSON.parse(text);
        } catch (e) {
            data = text;
        }
    }

    if (!response.ok) {
        const rawError = data?.err || data?.error || data?.message || 'API request failed';
        throw new Error(normalizeErrorMessage(rawError));
    }

    return data;
}

function normalizeErrorMessage(raw: unknown): string {
    const msg = String(raw ?? '').trim();
    if (!msg) return 'Request failed. Please try again.';

    const lower = msg.toLowerCase();

    if (lower.includes("failed on the 'required' tag")) return 'Please fill in all required fields';
    if (lower.includes("registerinput.password") && lower.includes("failed on the 'min' tag")) {
        return 'Password must be at least 8 characters';
    }
    if (lower.includes("registerinput.username") && lower.includes("failed on the 'min' tag")) {
        return 'Username must be at least 3 characters';
    }
    if (lower.includes("verifyemailregisterinput.code") && lower.includes("failed on the 'len' tag")) {
        return 'Verification code must be exactly 6 digits';
    }
    if (lower.includes("failed on the 'email' tag")) return 'Enter a valid email address';

    if (lower.includes('users_email_key')) return 'Email is already in use';
    if (lower.includes('users_username_key')) return 'Username is already in use';
    if (lower.includes('duplicate key value violates unique constraint') || lower.includes('unique constraint')) {
        return 'This value is already in use';
    }
    if (lower.includes('violates foreign key constraint')) return 'Related record was not found';
    if (lower.includes('violates not-null constraint')) return 'Required field is missing';
    if (lower.includes('invalid input syntax')) return 'Invalid input format';
    if (lower.includes('pq:') || lower.includes('sqlstate') || lower.includes('sql:')) {
        return 'Database request failed. Please try again';
    }

    return msg;
}

export const api = {
    login: async (username: string, password: string) => {
        const res = await request('/auth/login', {
            method: 'POST',
            body: JSON.stringify({ username, password })
        });
        return res;
    },
    googleLogin: async (token: string) => {
        return request('/auth/google', {
            method: 'POST',
            body: JSON.stringify({ token })
        });
    },
    register: async (username: string, email: string, password: string) => {
        return request('/auth/register', {
            method: 'POST',
            body: JSON.stringify({ username, email, password })
        });
    },
    verifyEmail: async (email: string, code: string) => {
        return request('/auth/verify', {
            method: 'POST',
            body: JSON.stringify({ email, code })
        });
    },
    getUser: async () => {
        return request('/api/user/me');
    },
    updateUsername: async (username: string) => {
        return request('/api/user/me', {
            method: 'PATCH',
            body: JSON.stringify({ username })
        });
    },
    updatePassword: async (old_password: string, new_password: string) => {
        return request('/api/user/me/password', {
            method: 'PUT',
            body: JSON.stringify({ old_password, new_password })
        });
    },
    requestEmailUpdate: async (email: string) => {
        return request('/api/user/me/email', {
            method: 'POST',
            body: JSON.stringify({ email })
        });
    },
    verifyEmailUpdate: async (code: string) => {
        return request('/api/user/me/email', {
            method: 'PUT',
            body: JSON.stringify({ code })
        });
    },
    requestDeleteUser: async (password: string) => {
        return request('/api/user/me', {
            method: 'DELETE',
            body: JSON.stringify({ password })
        });
    },
    confirmDeleteUser: async (code: string) => {
        return request('/api/user/me/delete', {
            method: 'PUT',
            body: JSON.stringify({ code })
        });
    },
    getTodos: async (): Promise<Todo[]> => {
        return request('/api/todos');
    },
    createTodo: async (title: string, description: string, deadline: string): Promise<Todo> => {
        return request('/api/todos', {
            method: 'POST',
            body: JSON.stringify({ title, description, deadline })
        });
    },
    updateTodo: async (todo: Todo) => {
        return request(`/api/todos/${todo.id}`, {
            method: 'PUT',
            body: JSON.stringify(todo)
        });
    },
    deleteTodo: async (id: string) => {
        return request(`/api/todos/${id}`, {
            method: 'DELETE'
        });
    }
};
