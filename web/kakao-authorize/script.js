const KAKAO_ACCESS_TOKEN_KEY = "code";

const kakaoService = new KakaoService();

const accessToken = Router.getUrlParam(KAKAO_ACCESS_TOKEN_KEY);
if (!accessToken) {
    alert("잘못된 요청입니다.");
}
kakaoService.kakaoAuthorize(accessToken);