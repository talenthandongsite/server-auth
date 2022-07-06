const AUTH_STORAGE_KEY = 'auth';
const TOKEN_EXPIRE_DURATION = 1000 * 60 * 60;

class Authentication {

    static check() {
        const info = this.getInfo();

        if (!info) {
            return false;
        }
        
        if (!info.token) {
            return false;
        }

        if ((new Date()).getTime() > info.expiration) {
            return false;
        }

        return true;
    }

    static getInfo() {
        const info = localStorage.getItem(AUTH_STORAGE_KEY);

        if (info) {
            return JSON.parse(info);
        }

        return;
    }

    static signIn(email, password) {
        // TODO: add login logic

        const expiration = (new Date()).getTime() + TOKEN_EXPIRE_DURATION;
        // TODO: replace with actual value
        const info = {
            token: "test token",
            email: email,
            expiration: expiration
        };
        localStorage.setItem(AUTH_STORAGE_KEY, JSON.stringify(info));
    }

    static signOut() {
        localStorage.removeItem(AUTH_STORAGE_KEY);
    }
}