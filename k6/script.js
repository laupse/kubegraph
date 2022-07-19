import { describe } from 'https://jslib.k6.io/expect/0.0.4/index.js';
import http from 'k6/http';

export const options = {
  thresholds: {
    checks: [{threshold: 'rate == 1.00',abortOnFail: true }],
  }
};

export default function () {

    describe('Health', (t) => {
        const response = http.get('http://localhost:30000/api/health');

        t.expect(response.status)
        .as('response status')
        .toEqual(200)

    });

    describe('Fetch fields', (t) => {
        const response = http.get('http://localhost:30000/api/graph/fields');

        t.expect(response.status)
        .as('response status')
        .toEqual(200)
        .and(response)
        .toHaveValidJson()

    });

    describe('Fetch data', (t) => {
        const response = http.get('http://localhost:30000/api/graph/data');

        t.expect(response.status)
        .as('response status')
        .toEqual(200)
        .and(response)
        .toHaveValidJson()

    });
}