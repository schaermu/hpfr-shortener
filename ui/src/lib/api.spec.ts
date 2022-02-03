import fetchMock from 'jest-fetch-mock'
import ApiClient from './api';

describe('ApiClient', () => {
    beforeEach(() => {
        fetchMock.resetMocks();
    })

    test('it does return a json object', async () => {
        const mockRes = { short_url: 'https://hpfr.ch/Gk9Fj9' }
        fetchMock.mockResponseOnce(JSON.stringify(mockRes))

        const json = await new ApiClient().shortenUrl('http://foobar.org')

        expect(json).toEqual(mockRes)
    })
})