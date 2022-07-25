class HttpRequest {
    static async get(url, headers) {
        return await fetch(url, {
            method: 'GET',
            headers: headers
        }).then(result => result.json());
    }

    static async post(url, body, headers) {
        return await fetch(url, {
            method: 'POST',
            headers: headers,
            body: body
        }).then(result => result.json());
    }
}