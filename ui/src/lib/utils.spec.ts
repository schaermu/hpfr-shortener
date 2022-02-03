import {delayFor} from './utils'

describe('utils', () => {
    describe('delayFor', () => {
        beforeAll(() => {
            jest.useFakeTimers();
            jest.spyOn(global, 'setTimeout');
        })

        test('it does return a json object', async () => {
            delayFor(1000);

            expect(setTimeout).toHaveBeenCalledTimes(1);
            expect(setTimeout).toHaveBeenLastCalledWith(expect.any(Function), 1000);
        })
    })
})