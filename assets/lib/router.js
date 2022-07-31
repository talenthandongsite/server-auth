const REDIRECT_TO = 'REDIRECT_TO';

class Router {
    static checkCached() {
        const route = localStorage.getItem(REDIRECT_TO);
        return route && route.length;
    }

    static cacheRoute(route) {
        localStorage.setItem(REDIRECT_TO, route);
    }

    static cacheCurrentRoute() {
        localStorage.setItem(REDIRECT_TO, window.location.href);
    }

    static cacheKakaoCode() {
        const key = 'code';
        const params = new URLSearchParams(window.location.search);
        if (!params.has(key)) return;
        localStorage.setItem(KAKAO_TOKEN, params.get(key));
    }

    static getUrlParam(key) {
        const params = new URLSearchParams(window.location.search);
        if (!params.has(key)) return;
        return params.get(key);
    }

    static redirect() {
        const route = localStorage.getItem(REDIRECT_TO);

        if (!route || !route.length) {
            return;
        }

        window.location.href = route;
    }

}