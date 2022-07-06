const ROUTER_STORAGE_KEY = 'router';

class Router {
    static checkCached() {
        const route = localStorage.getItem(ROUTER_STORAGE_KEY);
        return route && route.length;
    }

    static cacheRoute(route) {
        localStorage.setItem(ROUTER_STORAGE_KEY, route);
    }

    static cacheCurrentRoute() {
        localStorage.setItem(ROUTER_STORAGE_KEY, window.location.href);
    }

    static goToCachedRoute() {
        const route = localStorage.getItem(ROUTER_STORAGE_KEY);

        if (!route || !route.length) {
            return;
        }

        window.location.href = route;
    }
}