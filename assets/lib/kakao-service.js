const API_KEY = '3497f1f23c6af4f55a0bdf7c86c06998';
const REDIRECT_URL = 'http://localhost:8080/app/kakao-authorize';

const KAKAO_HOST = 'kauth.kakao.com';
const KAKAO_TOKEN_REQUEST_PATH = '/oauth/token';

class KakaoService {
    constructor() {
        Kakao.init(API_KEY);
        this.status = Kakao.isInitialized();
    } 

    kakaoAuthorize() {
        if (this.status) {
            const redirectUri = REDIRECT_URL;
            Kakao.Auth.authorize({ 
                redirectUri,
                scope: "openid"
            });
        }
    }


}