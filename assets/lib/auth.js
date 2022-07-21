const AUTH_TOKEN = 'AUTH_TOKEN'
const AUTH_EXP = 'AUTH_EXP'
const AUTH_SERVER_URL = "/signin";

class Authentication {

    static check() {
        const info = this.getInfo();

        if (!info) {
            return false;
        }
        
        if (!info.token) {
            return false;
        }

        if ((new Date()).getTime() > info.exp) {
            return false;
        }

        return true;
    }

    static getInfo() {
        const token = localStorage.getItem(AUTH_TOKEN);
        const exp = localStorage.getItem(AUTH_EXP);

        return { token, exp }
    }

    static signIn(username, password) {
        HttpRequest.post(AUTH_SERVER_URL, {username, password}, { 'Content-Type': 'application/json' }).then(result => {
            const { status, data } = result;
            if (!status) {
                alert("유저명, 혹은 비밀번호가 틀렸습니다.")
                return;
            }
            const { token, exp } = data;
            localStorage.setItem(AUTH_TOKEN, token);
            localStorage.setItem(AUTH_EXP, exp);

            Router.redirect();
        });
    }

    static signOut() {
        localStorage.removeItem(AUTH_TOKEN);
        localStorage.removeItem(AUTH_EXP);
    }

    static kakaoLogin() {

    }
}