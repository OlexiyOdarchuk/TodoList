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
        throw new Error(data?.err || data?.error || data?.message || 'API request failed');
    }

    return data;
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
