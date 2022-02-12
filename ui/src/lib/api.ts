class StatisticsResponse {
  hits: number;
  hitTimeData: [number, number][];
}

export default class ApiClient {
  baseUrl: string = '';

  constructor() {
    this.baseUrl = '/api';
  }

  private async call<T>(
    url: string,
    method: string = 'GET',
    body: any = null
  ): Promise<T> {
    let request: RequestInit = {
      method: method,
      headers: {
        'Content-Type': 'application/json',
      },
    };

    if (body != null) {
      request.body = JSON.stringify(body);
    }

    return new Promise((resolve, reject) => {
      fetch(url, request)
        .then((response) => {
          response.json().then((body) => {
            // produce errors for certain http status codes
            if ([400, 401, 403, 405, 500].indexOf(response.status) > -1) {
              reject(body.message);
            } else {
              resolve(body);
            }
          });
        })
        .catch((err) => {
          reject(err.message);
        });
    });
  }

  public async shortenUrl(url: string): Promise<{ short_url: string }> {
    return this.call(`${this.baseUrl}/shorten`, 'POST', { url });
  }

  public async getStatistics(code: string): Promise<StatisticsResponse> {
    return this.call(`${this.baseUrl}/statistics?code=${code}`);
  }
}
