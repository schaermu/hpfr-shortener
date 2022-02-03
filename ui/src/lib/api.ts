export default class ApiClient {
    baseUrl: string = ''

    constructor() {
        this.baseUrl = '/api'
    }

    private async call<T>(url: string, method: string, body: any = null) : Promise<T> {
        let request: RequestInit = {
            method: method,
            headers: {
                'Content-Type': 'application/json',
            }
        };

        if (body != null) {
            request.body = JSON.stringify(body)
        }

        return fetch(url, request)
            .then(response => {
                // produce errors for certain http status codes
                if ([400, 401, 403, 500].indexOf(response.status) > -1) {
                    throw new Error(response.statusText)
                }
                return response.json() as Promise<T>
            })
            .catch((err) => {
                return err
            })
            
    }

    public async shortenUrl(url: string) : Promise<{ short_url: string }> {
        return this.call(`${this.baseUrl}/shorten`, 'POST', { url })
    }
}