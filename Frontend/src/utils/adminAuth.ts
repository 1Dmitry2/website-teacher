export const ADMIN_TOKEN_KEY = 'admin_token';

export const adminAuthService = {
  setToken(token: string): void {
    localStorage.setItem(ADMIN_TOKEN_KEY, token);
  },

  getToken(): string | null {
    return localStorage.getItem(ADMIN_TOKEN_KEY);
  },

  removeToken(): void {
    localStorage.removeItem(ADMIN_TOKEN_KEY);
  },

  isAuthenticated(): boolean {
    return !!this.getToken();
  },
};


